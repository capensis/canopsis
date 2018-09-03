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

from __future__ import unicode_literals


class PersistentKeyValue(object):
    """
    This class is used to store named values persistently in the MongoDB
    database.
    """

    def __init__(self, mongo_store, collection):
        self.collection = mongo_store.get_collection(collection)

    def get(self, key):
        """
        Get a value from the storage.

        :param str key: The name of the value
        :rtype: Any
        """
        document = self.collection.find_one({
            '_id': key
        })
        if document is None:
            return None

        return document['value']

    def set(self, key, value):
        """
        Set a value in the storage.

        :param str key: The name of the value
        :param Any value: The value
        """
        return self.collection.update_one({
            '_id': key
        }, {
            '$set': {
                'value': value
            }
        }, upsert=True)
