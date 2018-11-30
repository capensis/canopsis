# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

import json

from canopsis.common.mongo_store import MongoStore
from canopsis.heartbeat import HeartbeatManager, HeartbeatPatternExistsError
from canopsis.models.heartbeat import HeartBeat
from canopsis.common.collection import MongoCollection

from .manager import MigrationModule


class HeartbeatMigrationSource(object):
    """
    Heartbeat migration source abstraction.

    """
    COLLECTION = "configuration"
    ID = "_id"
    GLOBAL_CONF_ID = "global_config"
    HEARTBEAT_SECTION = "heartbeat"
    ITEMS_KEY = "items"
    MAPPINGS_KEY = "mappings"
    MAX_DUR_KEY = "maxduration"

    @classmethod
    def provide_default_basics(cls):
        """
        Provide mongo collection.

        ! Do not use in tests !

        :rtype: `~.common.collection.MongoCollection`.
        """
        store = MongoStore.get_default()
        return (MongoCollection(store.get_collection(cls.COLLECTION)), )

    def __init__(self, mongo_collection):
        self.__collection = mongo_collection

    def get_old_heartbeat_items(self):
        """
        Get old Heartbeat items from the **configuration** collection.

        :returns: list of an old heartbeat items.
        :rtype: `list`.
        """
        global_config = \
            self.__collection.find_one({self.ID: self.GLOBAL_CONF_ID})
        if global_config:
            try:
                return global_config[self.HEARTBEAT_SECTION][self.ITEMS_KEY]
            except KeyError:
                pass
        return []

    def get_new_model_from_old_item(self, heartbeat_item):
        """
        Convert an old Heartbeat item to a new Heartbeat model.

        :param `dict` heartbeat_item: old Heartbeat item.
        :returns: a Heartbeat model else None if conversion fails.
        :rtype: `Optional[HeartBeat]`.
        """
        try:
            return HeartBeat({
                HeartBeat.PATTERN_KEY:
                    heartbeat_item[self.MAPPINGS_KEY],
                HeartBeat.EXPECTED_INTERVAL_KEY:
                    heartbeat_item[self.MAX_DUR_KEY]
            })
        except (ValueError, KeyError):
            return


class HeartbeatModule(MigrationModule):
    """
    Heartbeat migration module.

    """
    def init(self, yes=None):
        pass

    def update(self, yes=None):
        migration_source = HeartbeatMigrationSource(
            *HeartbeatMigrationSource.provide_default_basics())
        manager = HeartbeatManager(
            *HeartbeatManager.provide_default_basics())
        print("Looking for old Heartbeat items..")
        items = migration_source.get_old_heartbeat_items()
        if not items:
            print("No previously Heartbeat items found.")
            print("Heartbeat migration was skipped.")
            return
        total = len(items)
        print("{} old Heartbeat items found.".format(total))
        print("Started Heartbeat items migration..")
        failed = 0
        for heartbeat_item in items:
            model = migration_source\
                .get_new_model_from_old_item(heartbeat_item)
            if not model:
                print("Heartbeat item conversion failed: \n{}"
                      .format(json.dumps(heartbeat_item, indent=3,
                                         sort_keys=True)))
                failed += 1
                continue
            try:
                manager.create(model)
            except HeartbeatPatternExistsError:
                print("Duplicate Heartbeat mapping occured: \n{}"
                      .format(json.dumps(heartbeat_item, indent=3,
                                         sort_keys=True)))
        print("Heartbeat migration done:")
        print("  {} items was updated successfully".format(total-failed))
        print("  {} items could not updated".format(failed))
        print("Note! The old Heartbeat documents was not removed or "
              "modified for backward compatibility reason.")
