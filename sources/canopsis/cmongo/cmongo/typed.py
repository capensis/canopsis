#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from cmongo import Storage

from cstorage.typed import TypedStorage


class TypedStorage(Storage, TypedStorage):

    class Key:

        VALUE = 'v'
        TYPE = 't'

    TYPE_INDEX = [(Key.TYPE, Storage.ASC)]

    def _get_indexes(self, *args, **kwargs):

        result = super(TypedStorage, self)._get_indexes(*args, **kwargs)

        result.append(TypedStorage.TYPE_INDEX)

        return result

    def get(
        self, ids=None, data_type=None, limit=0, skip=0, sort=None,
        *args, **kwargs
    ):

        query = dict()

        if ids is not None:
            query['_id'] = {'$in': ids}

        if data_type is not None:
            query[TypedStorage.Key.TYPE] = data_type

        cursor = self._find(document=query)

        if ids is not None:
            cursor.hint([('_id', Storage.ASC)])

        elif data_type is not None:
            cursor.hint(TypedStorage.TYPE_INDEX)

        if limit:
            cursor.limit(limit)
        if skip:
            cursor.start(skip)
        if sort:
            Storage._update_sort(sort)
            cursor.sort(sort)

        result = list(cursor)

        return result

    def put(self, _id, data, data_type, *args, **kwargs):

        query = {
            '_id': _id
        }

        _set = {
            '$set': data
        }

        query[TypedStorage.Key.TYPE] = data_type

        self._update(_id=query, document=_set, multi=False)

    def remove(self, ids=None, data_type=None, *args, **kwargs):

        query = dict()

        if ids is not None:
            query['_id'] = {'$in': ids}

        if data_type is not None:
            query[TypedStorage.Key.TYPE] = data_type

        self._remove(document=query)
