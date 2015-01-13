#!/usr/bin/env python
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.old.record import Record

import os
import sys
import json

logger = None

cservices_path = os.path.expanduser('~/opt/mongodb/load.d/cservices')

##set root account
root = Account(user="root", group="root")
storage = get_storage(account=root, namespace='object')


def init():

    for path, folders, files in os.walk(cservices_path):
        for filename in files:
            filepath = os.path.join(path, filename)

            with open(filepath) as f:
                data = json.loads(f.read())

                try:
                    _id = data.pop('id')
                except KeyError as err:
                    print >>sys.stderr, "Can't parse cservice, missing key:", err
                    sys.exit(1)

                create_cservice(_id, filename, data)


def update():

    init()


def create_cservice(_id, name, data, mod='o+r', autorm=True, internal=False):
    #Delete old cservice

    try:
        record = storage.get('cservice.%s' % _id)
        if autorm:
            storage.remove(record)
        else:
            return record
    except:
        pass

    logger.info(" + Create cservice '%s'" % name)
    record = Record(data, _type='cservice', name=name)
    storage.put(record)
    return record
