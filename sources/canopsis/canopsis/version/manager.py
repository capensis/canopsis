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

from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection


class CanopsisVersionManager(object):

    __COLLECTION = "configuration"
    __VERSION_FIELD = "version"
    __DOCUMENT_ID = "canopsis_version"

    def __init__(self):
        store = MongoStore.get_default()
        collection = store.get_collection(name=self.__COLLECTION)
        self.__collection = MongoCollection(collection)

    @property
    def version_field(self):
        return self.__VERSION_FIELD

    def find_canopsis_version_document(self):
        return self.__collection.find_one({
            '_id': self.__DOCUMENT_ID
        })

    def put_canopsis_version_document(self, version):
        """

        :param version: `str` Canopsis version.
        """
        self.__collection.update(
            {
                '_id': self.__DOCUMENT_ID
            },
            {
                '_id': self.__DOCUMENT_ID,
                self.__VERSION_FIELD: version
            },
            upsert=True
        )
