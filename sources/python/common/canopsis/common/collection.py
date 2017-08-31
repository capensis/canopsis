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
from pymongo.errors import PyMongoError, OperationFailure

from canopsis.logger import Logger

LOG_NAME = 'collection'
LOG_PATH = 'var/log/collection.log'


class CollectionError(Exception):
    """
    Generic error on a MongoCollection.
    """
    pass


class CollectionSetError(Exception):
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
        self.collection = collection

        if logger is not None:
            self.logger = logger
        else:
            self.logger = Logger.get(LOG_NAME, LOG_PATH)

    def find(self, query):
        """
        Find elements in the collection.

        :param dict query: a query search
        :rtype: MongoCursor
        """
        return self.collection.find(query)

    def find_one(self, query):
        """
        Find one element in the collection.

        :param dict query: a query search
        :rtype: MongoCursor
        """
        return self.collection.find_one(query)

    def insert(self, document):
        """
        Update an element in the collection.

        :param dict document: the document to insert
        :returns: the _id in the document or a generated one
        :rtype: str
        """
        try:
            return self.collection.insert(document)

        except OperationFailure as of_err:
            message = 'Operation failure while doing insert: {}'.format(of_err)
        except Exception:
            message = 'Unknown exception on collection insert'

        self.logger.error(message)
        raise CollectionError(message)

    def update(self, query, document, upsert=False):
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
            return self.collection.update(query, document, upsert=upsert)

        except BSONError as ex:
            message = 'document error: {}'.format(ex)
        except PyMongoError as ex:
            message = 'pymongo error: {}'.format(ex)
        except OperationFailure as of_err:
            message = 'Operation failure while doing update: {}'.format(of_err)
        except TypeError:
            message = []
            if not isinstance(query, dict):
                message.append('query is not a dict')
            if not isinstance(document, dict):
                message.append('document is not a dict')
            if not isinstance(upsert, bool):
                message.append('upsert is not a boolean')
            message = ' ; '.join(message)
        except Exception:
            message = 'Unknown exception on collection update'

        self.logger.error(message)
        raise CollectionError(message)

    def remove(self, query={}):
        """
        Remove an element in the collection.

        :param dict query: a mongo search query
        :raises: OperationFailure
        :rtype: dict
        """
        try:
            return self.collection.remove(query)

        except OperationFailure as of_err:
            message = 'Operation failure while doing remove: {}'.format(of_err)
        except Exception:
            message = 'Unknown error while doing remove'

        self.logger.error(message)
        raise CollectionError(message)

    @staticmethod
    def is_successfull(dico):
        """
        Check if a pymongo dict response is a success ({'ok': 1.0, 'n': 2})

        :param dict dico: a pymongo dict response on update and remove
        :rtype: bool
        """
        return 'ok' in dico and dico['ok'] == 1.0
