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

from canopsis.common.collection import MongoCollection
from canopsis.models.heartbeat import HeartBeat
from canopsis.common.mongo_store import MongoStore
from canopsis.logger import Logger


class HeartbeatError(Exception):
    """
    Base Heartbeat error.
    """


class HeartbeatPatternExistsError(HeartbeatError):
    """
    Heartbeat pattern exists error.
    """


class HeartbeatManager(object):
    """
    Heartbeat service manager abstraction.
    """

    COLLECTION = 'heartbeat'

    LOG_PATH = 'var/log/heartbeat.log'

    @classmethod
    def provide_default_basics(cls):
        """
        Provide logger, config, storages...

        ! Do not use in tests !

        :rtype: Union[logging.Logger,
                      canopsis.common.collection.MongoCollection]
        """
        store = MongoStore.get_default()
        collection = store.get_collection(name=cls.COLLECTION)
        return (Logger.get('action', cls.LOG_PATH),
                MongoCollection(collection))

    def __init__(self, logger, collection):
        """

        :param `~.logger.Logger` logger: object.
        :param `~.common.collection.MongoCollection` collection: object.
        """
        self.__logger = logger
        self.__collection = MongoCollection(collection)

    def create_heartbeat(self, heartbeat):
        """
        Create a new Heartbeat.

        :param `HeartBeat` heartbeat: a Heartbeat model.

        :returns: a created Heartbeat ID.
        :rtype: `str`.

        :raises: (`ValueError`,
                  `.HeartbeatPatternExistsError`,
                  `pymongo.errors.PyMongoError`,
                  `~.common.collection.CollectionError`, ).
        """
        if self.find_heartbeat_document(heartbeat.id):
            raise HeartbeatPatternExistsError()

        return self.__collection.insert(heartbeat.to_dict())

    def find_heartbeat_document(self, heartbeat_id):
        """

        :param heartbeat_id:
        :return:
        :raises: (`pymongo.errors.PyMongoError`, ).
        """
        return self.__collection.find_one({"_id": heartbeat_id})

    def remove_heartbeat_document(self, heartbeat_id):
        """
        Remove Heartbeat by ID.

        :param `str` heartbeat_id: Heartbeat ID.
        :return:
        :raises: (`~.common.collection.CollectionError`, ).
        """
        return self.__collection.remove({"_id": heartbeat_id})

    def list_heartbeat_collection(self):
        """
        Get Heartbeats list.

        :returns: list of heartbeat documents.
        :raises: (`pymongo.errors.PyMongoError`, ).
        """
        return self.__collection.find({})
