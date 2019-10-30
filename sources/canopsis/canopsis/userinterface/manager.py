#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2019 "Capensis" [http://www.capensis.com]
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
User Interface Configuration manager.
"""

from __future__ import unicode_literals

from canopsis.common.collection import MongoCollection
from canopsis.common.mongo_store import MongoStore
from canopsis.logger import Logger
from canopsis.models.userinterface import UserInterface


class UserInterfaceManager(object):
    """
    UserInterface managment.
    """
    LOG_PATH = 'var/log/configuration.log'
    COLLECTION = 'configuration'
    __DOCUMENT_ID = "user_interface"

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
        logger = Logger.get('configuration', cls.LOG_PATH)
        store = MongoStore.get_default()
        collection = store.get_collection(name=cls.COLLECTION)
        mongo_collection = MongoCollection(collection)

        return (logger, mongo_collection)

    def get(self):
        """
        Read a user interface config.

        :param str query: an user interface config query string
        :rtype: UserInterface or None
        """
        record = self.collection.find_one(query={"_id": self.__DOCUMENT_ID})
        if not record:
            return

        user_interface = UserInterface(**record)
        return user_interface

    def update(self, interface):
        """
        Update a user interface config.

        :param str id_: an user interface config _id
        :param dict user_interface: an user interface config as a dict
        :rtype: bool
        """
        resp = self.collection.update({"_id": self.__DOCUMENT_ID}, {
                                      '$set': interface}, upsert=True)

        return self.collection.is_successfull(resp)

    def delete(self):
        """
        Update a user interface config.

        :param str id_: an user interface config _id
        :param dict user_interface: an user interface config as a dict
        :rtype: bool
        """
        resp = self.collection.remove({"_id": self.__DOCUMENT_ID})
        return self.collection.is_successfull(resp)
