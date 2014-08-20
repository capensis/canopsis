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

from canopsis.common.utils import isiterable
from canopsis.mongo import MongoStorage
from canopsis.storage.scoped import ScopedStorage


class MongoScopedStorage(MongoStorage, ScopedStorage):

    def _get_indexes(self, *args, **kwargs):

        result = super(MongoScopedStorage, self)._get_indexes(*args, **kwargs)

        scope = self.scope

        if scope is not None:
            # add all sub_scopes concatened with id
            for scope_count in range(1, len(scope)):
                sub_scope = scope[:scope_count]
                index = [(scope_name, 'text') for scope_name in sub_scope]
                index.append((ScopedStorage.ID, 'text'))
                result.append(index)

        return result

    def get(
        self, scope, ids=None, _filter=None, limit=0, skip=0, sort=None,
        *args, **kwargs
    ):

        query = scope.copy()
        if _filter is not None:
            query.update(_filter)

        result = self.get_elements(
            ids=ids, query=query, limit=limit, skip=skip, sort=sort)

        return result

    def find(self, scope, _filter, limit=0, skip=0, sort=None):

        return self.get(
            scope=scope, _filter=_filter, limit=limit, skip=skip, sort=sort)

    def put(self, scope, _id, data, *args, **kwargs):

        # get unique id
        _id = self.get_absolute_path(scope=scope, _id=_id)

        # get query such as a sum of scope and _id
        query = scope.copy()
        query[MongoStorage.ID] = _id

        _set = {
            '$set': {ScopedStorage.VALUE: data}
        }

        self._update(_id=query, document=_set, multi=False)

    def remove(self, scope, ids=None, *args, **kwargs):

        query = scope.copy()

        parameters = {}

        if ids is not None:
            if isiterable(ids, is_str=False):
                query[ScopedStorage.ID] = {'$in': ids}
            else:
                parameters = {'justOne': 1}
                query[ScopedStorage.ID] = ids

        self._remove(document=query, **parameters)
