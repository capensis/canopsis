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

from __future__ import unicode_literals

from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
from canopsis.logger import Logger


class HeartBeatServiceException(Exception):
    pass


class HeartBeatService:
    """HeartBeat mapping management."""

    HEARTBEAT_COLLECTION = "configuration"
    LOG_PATH = 'var/log/heartbeat.log'

    ID = "_id"
    GLOBAL_CONF_ID = "global_config"
    HEARTBEAT_SECTION = "heartbeat"
    MAPPINGS_KEY = "MAPPINGS"
    ITEMS_KEY = "items"

    @classmethod
    def provide_default_basics(cls):
        """
        Provide logger, config, storages...

        ! Do not use in tests !

        :rtype: Union[logging.Logger,
                      canopsis.common.collection.MongoCollection]
        """
        logger = Logger.get('action', cls.LOG_PATH)
        store = MongoStore.get_default()
        collection = store.get_collection(name=cls.HEARTBEAT_COLLECTION)
        mongo_collection = MongoCollection(collection)

        return logger, mongo_collection

    def __init__(self, logger, mongo_collection):
        self.logger = logger
        self.collection = mongo_collection

    def __get_conf(self):
        return self.collection.find_one({self.ID: self.GLOBAL_CONF_ID})

    def get_heartbeats(self):
        global_config = self.__get_conf()
        return global_config[self.HEARTBEAT_SECTION]

    def create(self, heartbeat):
        """
        Create a new heartbeat in the database from a heartbeat model instance.

        :param heartbeat: a heartbeat model instance.
        :raises: CollectionError if an error occured while the heartBeat is
        stored into the database, HeartBeatServiceException if the given
        heartbeat is not valid.
        """
        valid, error_message = heartbeat.isValid()
        if not valid:
            raise HeartBeatServiceException(error_message)

        global_config = self.__get_conf()
        hb_Section = global_config[self.HEARTBEAT_SECTION]
        hb_Section[self.ITEMS_KEY].append(heartbeat.to_dict())

        self.collection.update({"_id": self.GLOBAL_CONF_ID},
                               {"$set": {self.HEARTBEAT_SECTION: hb_Section}})
