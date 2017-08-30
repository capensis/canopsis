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


class CollectionSetError(Exception):
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
        :rtype: dict or None
        """
        return NotImplementedError()

    def update(self, query, document, upsert=False):
        """
        Update an element in the collection.

        :param dict query: a query search
        :param dict document: the document to update
        :param bool upsert: do insert if the document does not already exist
        :rtype: dict or None
        """
        try:
            return self.collection.update(query, document, upsert=upsert)
        except BSONError as ex:
            raise CollectionSetError('document error: {}'.format(ex))
        except PyMongoError as ex:
            raise CollectionSetError('pymongo error: {}'.format(ex))
        except OperationFailure:
            self.logger.error('Operation failure while doing update')
        except TypeError:
            if not isinstance(query, dict):
                self.logger.error('query is not a dict')
            if not isinstance(document, dict):
                self.logger.error('document is not a dict')
            if not isinstance(upsert, bool):
                self.logger.error('upsert is not a boolean')
        except Exception:
            self.logger.error('Unkown exception on collection update')

        return None

    def remove(self, query={}):
        """
        Remove an element in the collection.

        :param dict query:
        :rtype: dict or None
        """
        try:
            return self.collection.remove(query)
        except OperationFailure:
            self.logger.error('Operation failure while doing remove')
        except Exception:
            self.logger.error('Unkown error while doing remove')

        return None

    @staticmethod
    def is_successfull(dico):
        """
        Check if a pymongo dict response is a success ({'ok': 1.0, 'n': 2})

        :param dict dico: a pymongo dict response on update, remove...
        :rtype: bool
        """
        return 'ok' in dico and dico['ok'] == 1.0
