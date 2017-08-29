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

from canopsis.logger import Logger


class MongoCollection(object):
    """
    A mongodb collection handeling class, to ease access to mongodb.
    """

    def __init__(self, collection, logger=None):
        self.collection = collection

        if logger is not None:
            self.logger = logger
        else:
            self.logger = Logger.get('collection', 'var/log/collection.log')

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
            result = self.collection.update(query, document, upsert=upsert)
        except TypeError:
            if not isinstance(query, dict):
                self.logger.error('query is not a dict')
            if not isinstance(document, dict):
                self.logger.error('document is not a dict')
            if not isinstance(upsert, bool):
                self.logger.error('upsert is not a boolean')
            return None
        except:
            self.logger.error('Unkown exception on collection update')
            return None

        return result

    def remove(self, query):
        """
        Remove an element in the collection.
        """
        return NotImplementedError()

    @staticmethod
    def is_mongo_successfull(dico):
        """
        Check if a pymongo dict response report a success ({'ok': 1.0, 'n': 2})

        :param dict dico: a pymongo dict response on update, remove...
        :rtype: bool
        """
        return 'ok' in dico and dico['ok'] == 1.0
