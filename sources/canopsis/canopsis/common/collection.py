# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from bson.errors import BSONError
from pymongo.errors import PyMongoError, OperationFailure, AutoReconnect

from time import sleep

from canopsis.logger import Logger
from canopsis.common.mongo_store import MongoStore

LOG_NAME = 'collection'
LOG_PATH = 'var/log/collection.log'


class CollectionError(Exception):
    """
    Generic error on a MongoCollection.
    """
    pass


class CollectionSetError(CollectionError):
    """
    Error on a set in a MongoCollection.
    """
    pass


class MongoCollection(object):
    """
    A mongodb collection handeling class, to ease access to mongodb.

    For futur generation: Behold ! Don't rebuild a storage layer like
    the old one, with over engineered classes/functions.
    """

    def __init__(self, collection, logger=None):
        """
        :param pymongo.collection.Collection collection: mongo Collection object
        :param logging.Logger logger: python logger instance.
            If None, a new logger is provided.
        """
        self.collection = collection
        self._hr = MongoStore.hr

        if logger is not None:
            self.logger = logger
        else:
            self.logger = Logger.get(LOG_NAME, LOG_PATH)

    def find(self, query, *args, **kwargs):
        """
        Find elements in the collection.

        :param dict query: a query search
        :rtype: pymongo.cursor.Cursor
        """
        return self._hr(self.collection.find, query, *args, **kwargs)

    def find_one(self, query, *args, **kwargs):
        """
        Find one element in the collection.

        :param dict query: a query search
        :rtype: pymongo.cursor.Cursor
        """
        return self._hr(self.collection.find_one, query, *args, **kwargs)

    def insert(self, document, *args, **kwargs):
        """
        Update an element in the collection.

        :param dict document: the document to insert
        :returns: the _id in the document or a generated one
        :rtype: str
        """
        try:
            return self._hr(self.collection.insert, document, *args, **kwargs)

        except OperationFailure as of_err:
            message = 'Operation failure while doing insert: {}'.format(of_err)
        except Exception as ex:
            message = 'Unknown exception on collection insert: {}'.format(ex)

        self.logger.error(message)
        raise CollectionError(message)

    def update(self, query, document, upsert=False, *args, **kwargs):
        """
        Update an element in the collection.

        Be carefull ! It behaves like mongos updates: the document
        parameter will totally replace the old one.
        Use { '$set': document } to softly update the document

        :param dict query: a mongo search query
        :param dict document: the document to update
        :param bool upsert: do insert if the document does not already exist
        :raises: BSONError, PyMongoError, OperationFailure, TypeError
        :rtype: dict
        """
        try:
            return self._hr(
                self.collection.update, query, document, upsert=upsert,
                *args, **kwargs
            )

        except BSONError as ex:
            message = 'document error: {}'.format(ex)
        except PyMongoError as ex:
            message = 'pymongo error: {}'.format(ex)
        except OperationFailure as of_err:
            message = 'Operation failure while doing update: {}'.format(of_err)
        except TypeError as ex:
            message = []
            if not isinstance(query, dict):
                message.append('query is not a dict')
            if not isinstance(document, dict):
                message.append('document is not a dict')
            if not isinstance(upsert, bool):
                message.append('upsert is not a boolean')
            message = ' ; '.join(message)
            message = '{}: {}'.format(ex, message)
        except Exception as ex:
            message = 'Unknown exception on collection update: {}'.format(ex)

        self.logger.error(message)
        raise CollectionError(message)

    def remove(self, query={}, *args, **kwargs):
        """
        Remove an element in the collection.

        :param dict query: a mongo search query
        :raises: OperationFailure
        :rtype: dict
        """
        try:
            return self._hr(self.collection.remove, query, *args, **kwargs)

        except OperationFailure as of_err:
            message = 'Operation failure while doing remove: {}'.format(of_err)
        except Exception as ex:
            message = 'Unknown error while doing remove: {}'.format(ex)

        self.logger.error(message)
        raise CollectionError(message)

    def find_and_modify(self, *args, **kwargs):
        return self._hr(
            self.collection.find_and_modify, *args, **kwargs
        )

    def save(self, *args, **kwargs):
        return self._hr(
            self.collection.save, *args, **kwargs
        )

    def count(self):
        """
        Counts the number of items in the current collection.
        """
        return self._hr(self.collection.count)

    def drop_indexes(self):
        return self._hr(self.collection.drop_indexes)

    def ensure_index(self, *args, **kwargs):
        return self._hr(self.collection.ensure_index, *args, **kwargs)

    def wrap_callable(self, func):
        def wrapped(*args, **kwargs):
            while True:
                try:
                    return func(*args, **kwargs)
                except AutoReconnect:
                    sleep(1)

        return wrapped

    def __getattr__(self, name):
        while True:
            try:
                res = None
                if hasattr(self, name):
                    res = super(MongoCollection, self).__getattribute__(name)
                else:
                    res = getattr(self.collection, name)

                if callable(res):
                    return self.wrap_callable(res)
                return res

            except AutoReconnect:
                sleep(1)

    @staticmethod
    def is_successfull(dico):
        """
        Check if a pymongo dict response is a success ({'ok': 1.0, 'n': 2})

        :param dict dico: a pymongo dict response on update and remove
        :rtype: bool
        """
        return 'ok' in dico and dico['ok'] == 1.0
