#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

"""
Action manager.
"""

from __future__ import unicode_literals

from canopsis.common.collection import MongoCollection
from canopsis.common.mongo_store import MongoStore
from canopsis.logger import Logger
from canopsis.models.action import Action


class ActionManager(object):
    """
    Action managment.
    """
    LOG_PATH = 'var/log/action.log'

    ACTION_COLLECTION = 'default_action'

    def __init__(self, logger, mongo_collection):
        self.logger = logger
        self.collection = mongo_collection

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
        collection = store.get_collection(name=cls.ACTION_COLLECTION)
        mongo_collection = MongoCollection(collection)

        return (logger, mongo_collection)

    def get_id(self, id_):
        """
        Helper to find just an object from his _id.
        """
        return self.get(query={Action._ID: id_})

    def get(self, query):
        """
        Read an action.

        :param str query: an action query string
        :rtype: Action or None
        """
        record = self.collection.find_one(query=query)
        if not record:
            return

        action = Action(**Action.convert_keys(record))
        return action

    def create(self, action):
        """
        Create an action.

        :param dict action: an action as a dict
        :rtype: bool
        """
        id_ = self.collection.insert(document=action)

        return id_ is not None

    def update_id(self, id_, action):
        """
        Update an action.

        :param str id_: an action _id
        :param dict action: an action as a dict
        :rtype: bool
        """
        query = {Action._ID: id_}
        resp = self.collection.update(query=query, document=action)

        return self.collection.is_successfull(resp)

    def delete_id(self, id_):
        """
        Delete an action.

        :param str id_: an action _id
        :rtype: bool
        """
        query = {Action._ID: id_}
        resp = self.collection.remove(query=query)

        return self.collection.is_successfull(resp)
