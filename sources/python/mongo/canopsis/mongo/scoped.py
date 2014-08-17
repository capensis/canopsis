# -*- coding: utf-8 -*-
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

from canopsis.mongo import MongoStorage
from canopsis.storage.scoped import ScopedStorage


class MongoScopedStorage(MongoStorage, ScopedStorage):

    __datatype__ = 'scoped'  #: register this class to scoped data types

    def _get_indexes(self, *args, **kwargs):

        result = super(MongoScopedStorage, self)._get_indexes(*args, **kwargs)

        scope = self.scope

        # add all sub_scopes concatened with id
        for scope_count in range(1, len(scope) + 1):
            sub_scope = scope[:scope_count]
            index = [(scope_name, 'text') for scope_name in sub_scope]
            index.append((ScopedStorage.ID, 'text'))
            result.append(index)

        return result

    def get(self, scope, data_id, *args, **kwargs):

        query = {}

        query.update(scope)

        query[ScopedStorage.ID]

        cursor = self._find(document=query)

        index = self._get_index(query)

        cursor.hint(index)

        self._get_generic_result
        result = list(cursor)

        return result

    def put(self, _id, data, data_type, *args, **kwargs):

        query = {
            '_id': _id
        }

        _set = {
            '$set': data
        }

        query[ScopedStorage.Key.TYPE] = data_type

        self._update(_id=query, document=_set, multi=False)

    def remove(self, ids=None, data_type=None, *args, **kwargs):

        query = {}

        if ids is not None:
            query['_id'] = {'$in': ids}

        if data_type is not None:
            query[ScopedStorage.Key.TYPE] = data_type

        self._remove(document=query)
