# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2019 "Capensis" [http://www.capensis.fr]
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

from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection

class WatcherManager(object):
    """
    Manager for the watchers
    """

    COLLECTION = "default_entities"

    def __init__(self, mongo_collection):
        """
        :param mongo_collection: `pymongo.collection.Collection` object.
        """
        super(WatcherManager, self).__init__()
        self.__collection = mongo_collection

    @classmethod
    def default_collection(cls):
        """
        Returns the default collection for the manager.

        ! Do not use in tests !

        :rtype: canopsis.common.collection.MongoCollection
        """
        store = MongoStore.get_default()
        collection = store.get_collection(name=cls.COLLECTION)
        return MongoCollection(collection)

    def get_watcher_list(self):
        """
        Return a list of all the watchers.

        :rtype: List[Dict[str, Any]]
        """
        return list(self.__collection.find({'type':'watcher'}))

    def get_watcher_by_id(self, wid):
        """
        Get a watcher given its id.

        :param str rule_id: the id of the watcher.
        :rtype: Dict[str, Any]
        """
        return self.__collection.find_one({'_id': wid})

    def create_watcher(self, watcher):
        """
        Create a watcher and return its id.

        :param Dict[str, Any] watcher:
        :rtype: str
        :raises: CollectionError if the creation fails.
        """
        return self.__collection.insert(watcher)

    def update_watcher_by_id(self, watcher, wid):
        """
        Update a watcher given its id.
        Return a boolean if the operation is successful.

        :param str wid: the id of the watcher.
        :param Dict[str, Any] watcher:
        :rtype: bool
        :raises: CollectionError if the update fails.
        """
        resp = self.__collection.update(query={'_id': wid}, document=watcher)
        return self.__collection.is_successfull(resp)

    def delete_watcher_by_id(self, wid):
        """
        Remove a watcher given its id.
        Return a boolean if the operation is successful.

        :param str wid: the id of the rule.
        :rtype: bool
        :raises: CollectionError if the deletion fails.
        """
        resp = self.__collection.remove({'_id': wid})
        return self.__collection.is_successfull(resp)
