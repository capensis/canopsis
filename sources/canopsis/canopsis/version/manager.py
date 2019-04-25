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


class CanopsisVersionManager(object):
    """
    Canopsis version manager abstraction.

    """

    COLLECTION = "configuration"
    EDITION_FIELD = "edition"
    STACK_FIELD = "stack"
    VERSION_FIELD = "version"
    __DOCUMENT_ID = "canopsis_version"

    def __init__(self, collection):
        """

        :param collection: `pymongo.collection.Collection` object.
        """
        self.__collection = MongoCollection(collection)

    def find_canopsis_document(self):
        """
        Find Canopsis version document.

        :returns: Canopsis version document or None if not found.

        :raises: (`pymongo.errors.PyMongoError`, ).
        """
        return self.__collection.find_one({
            '_id': self.__DOCUMENT_ID
        })

    def put_canopsis_document(self, version, stack, edition):
        """
        Put Canopsis version document (upsert).

        :param version: `str` Canopsis version.

        :raises: (`canopsis.common.collection.CollectionError`, ).
        """
        self.__collection.update(
            {
                '_id': self.__DOCUMENT_ID
            },
            {
                '_id': self.__DOCUMENT_ID,
                self.VERSION_FIELD: version,
                self.EDITION_FIELD: edition,
                self.STACK_FIELD: stack
            },
            upsert=True
        )
