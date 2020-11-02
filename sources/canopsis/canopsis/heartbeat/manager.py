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
from canopsis.common.mongo_store import MongoStore
from canopsis.logger import Logger
from canopsis.models.heartbeat import HeartBeat
import time
import re


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

    def create(self, heartbeat):
        """
        Create a new Heartbeat.

        :param `HeartBeat` heartbeat: a Heartbeat model.

        :returns: a created Heartbeat ID.
        :rtype: `str`.

        :raises: (`.HeartbeatPatternExistsError`,
                  `pymongo.errors.PyMongoError`,
                  `~.common.collection.CollectionError`, ).
        """
        if self.get(heartbeat.id):
            raise HeartbeatPatternExistsError()
        now = int(time.time())
        heartbeat = heartbeat.to_dict()
        heartbeat[HeartBeat.CREATED_KEY] = now
        heartbeat[HeartBeat.UPDATED_KEY] = now
        return self.__collection.insert(heartbeat)

    def get(self, heartbeat_id=None, page=None, limit=None, search=None, sort=None, sort_by=None):
        """
        Get Heartbeat by ID or a list of Heartbeats
        when calling with default arguments.

        :param `Optional[str]` heartbeat_id:
        :returns: list of Heartbeat documents if **heartbeat_id** is None
                  else single Heartbeat document or None if not found.
        :raises: (`pymongo.errors.PyMongoError`, ).
        """

        if heartbeat_id:
            return self.__collection.find_one({"_id": heartbeat_id})

        pipeline = []
        if search is not None:
            or_query = [
                {"name": re.compile(str(search), re.IGNORECASE)},
                {"description": re.compile(str(search), re.IGNORECASE)},
                {"author": re.compile(str(search), re.IGNORECASE)}
            ]
            pipeline.append({"$match": {"$or": or_query}})
        else:
            pipeline.append({"$match": {}})

        total_count_data = list(self.__collection.aggregate(
            pipeline + [{'$count': 'total_count'}]))
        total_count = 0
        if len(total_count_data) == 1:
            try:
                total_count = total_count_data[0]["total_count"]
            except (IndexError, KeyError):
                self.__logger.error(
                    "Exception while trying to reach total_count")
                return {"meta": {"page": 0, "page_count": 0, "per_page": 0, "total_count": 0}, "data": []}

        sort_by = sort_by or "created"
        sort = sort or "desc"
        sort = -1 if sort == "desc" else 1
        pipeline.append({"$sort": {sort_by: sort}})

        page = int(page or 1)
        limit = int(limit or 10)
        pipeline.append({"$skip": (page - 1) * limit})
        pipeline.append({"$limit": limit})

        data = list(self.__collection.aggregate(pipeline))
        page_count = len(data)/limit + 1
        return {
            "meta": {
                "page": page,
                "page_count": page_count,
                "per_page": limit,
                "total_count": total_count
            },
            "data": data
        }

    def delete(self, heartbeat_id):
        """
        Delete Heartbeat by ID.

        :param `str` heartbeat_id: Heartbeat ID.
        :return:
        :raises: (`~.common.collection.CollectionError`, ).
        """
        return self.__collection.remove({"_id": heartbeat_id})

    def update(self, _id, heartbeat):
        try:
            current_heartbeat = self.__collection.find_one({"_id": _id})
            heartbeat = heartbeat.to_dict()
            if HeartBeat.CREATED_KEY in current_heartbeat:
                heartbeat[HeartBeat.CREATED_KEY] = current_heartbeat[HeartBeat.CREATED_KEY]
            heartbeat[HeartBeat.UPDATED_KEY] = int(time.time())
            heartbeat[HeartBeat.ID_KEY] = _id
            resp = self.__collection.update(query={'_id': _id}, document=heartbeat)
        except Exception as e:
            raise e
        return self.__collection.is_successfull(resp)
