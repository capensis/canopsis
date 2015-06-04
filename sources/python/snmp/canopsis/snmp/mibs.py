# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from canopsis.context.manager import Context
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
import os

#: snmp manager configuration category
CATEGORY = 'MIBS'

#: snmp manager configuration path
CONF_PATH = 'snmp/mibs.conf'


class MibsImportException(Exception):
    pass


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class MibsManager(MiddlewareRegistry):

    # The configuration key to read in the configuration file
    MIBS_STORAGE = 'mibs_storage'

    def __init__(self, *args, **kwargs):
        super(MibsManager, self).__init__(*args, **kwargs)

    def put(self, oid, info):
        self[MibsManager.MIBS_STORAGE].put_element(
            _id=oid, element=info)

    def get(self, oids=None, limit=None, query={}, projection=None):
        return self[MibsManager.MIBS_STORAGE].get_elements(
            ids=oids,
            query=query,
            limit=limit,
            projection=projection
        )

    def remove(self, oids=None):
        self[MibsManager.MIBS_STORAGE].remove_elements(ids=oids)

    def check_mib(self, filename):
        r = os.system('smilint {} 2>/dev/null'.format(filename))
        return r == 0

    def import_mib(self, filename):
        from subprocess import Popen, PIPE
        from os.path import exists
        print("-> parse {}".format(filename))

        if not exists(filename):
            raise MibsImportException("File not found: {}".format(filename))

        process = Popen(["smidump", "-k", "-f", "python", filename],
                        stdout=PIPE)
        stdout, _ = process.communicate()
        if process.returncode != 0:
            raise MibsImportException("Unable to convert the mib to json")
        if not stdout:
            raise MibsImportException("Empty or non-existent file")

        # ensure it's a python file
        try:
            code = compile(stdout, "generated-mib.py", "exec")
            mod = {}
            exec code in mod
        except:
            raise MibsImportException("Invalid python generated from smidump")
        if "MIB" not in mod:
            raise MibsImportException("No MIB to use!")
        mib = mod["MIB"]

        # insert the mib info
        mib_name = mib["moduleName"]
        mib_info = mib[mib_name]
        mib_key = "{}".format(mib_name)
        print("-> import mib {}".format(mib_name))
        self.put(mib_key, mib_info)

        # insert all notifications
        mib = mod["MIB"]
        notifications = mib.get("notifications", {})
        print("-> import {} notifications".format(len(notifications)))
        for notification_name, notification in notifications.items():
            # insert a name to recover later
            notification["name"] = notification_name
            notification_oid = notification["oid"]
            notification_objects = notification.get("objects", {})
            # prefix all the objects with module name.
            notification_objects = {
                (object["module"], name): object
                for name, object in notification_objects.items()}

            print("   - {}: {}".format(notification_name, notification_oid))
            if notification_objects:
                for object_module, object_name in notification_objects:
                    print("      - reference {}::{}".format(
                        object_module, object_name))
            self.put(notification_oid, notification)

            # save a relation between source MIB <-> oid
            field_id = "{}::{}".format(mib_name, notification_name)
            self.put(field_id, {"oid": notification_oid})

        # insert all objects
        objects = mib.get("nodes", {}).items()
        print("-> import {} objects".format(len(objects)))
        for object_name, object in objects:
            object = mib["nodes"][object_name]
            object_oid = object["oid"]
            print("   - {}::{}: {}".format(mib_name, object_name, object_oid))

            # check if the object exists in the database already
            ids = self.get(object_oid)
            if ids:
                print "    > already inserted"
                continue

            self.put(object_oid, object)

            # save relation between the source MIB <-> field
            field_id = "{}::{}".format(mib_name, object_name)
            self.put(field_id, {"oid": object_oid})

        return (len(notifications), len(objects))


if __name__ == "__main__":
    import sys
    import argparse
    import traceback

    parser = argparse.ArgumentParser(
        description="Import mibs to canopsis database")
    parser.add_argument("mibs", nargs="*",
        help="files to import")
    parser.add_argument("-k", action="store_true",
        help="keep going even on error")
    parser.add_argument("-D", action="store_true",
        help="drop all the database before import")
    parser.add_argument("--query",
        help="search an oid in the database")

    args = parser.parse_args()

    manager = MibsManager()

    if args.D:
        manager.remove()
        sys.exit(0)

    if args.query:
        from pprint import pprint
        ret = manager.get(args.query)
        pprint(ret)
        sys.exit(0)

    if args.mibs:
        ret = 0
        notif_count = object_count = 0
        errors = []
        for filename in args.mibs:
            try:
                counts = manager.import_mib(filename)
                notif_count += counts[0]
                object_count += counts[1]
            except Exception as e:
                traceback.print_exc()
                errors.append((filename, e))
                if not args.k:
                    ret = 1
                    break

        print("-" * 70)
        print("Import summary")
        print("- {} notifications definitions".format(notif_count))
        print("- {} objects definitions".format(object_count))
        if errors:
            ret = 1
            print("- {} error{}".format(len(errors), "s" if len(errors) > 1 else ""))
            for error in errors:
                print("  - {}: {}".format(*error))

    sys.exit(ret)
