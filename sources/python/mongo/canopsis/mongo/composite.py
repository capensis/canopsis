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
from canopsis.storage import Storage
from canopsis.storage.composite import CompositeStorage


class MongoCompositeStorage(MongoStorage, CompositeStorage):

    def _get_indexes(self, *args, **kwargs):

        result = super(MongoCompositeStorage, self)._get_indexes(
            *args, **kwargs)

        # add all sub_paths concatened with id
        for n, _ in enumerate(self.path):
            sub_path = self.path[:n + 1]
            index = [(path_name, Storage.ASC) for path_name in sub_path]
            index.append((CompositeStorage.ID, Storage.ASC))
            result.append(index)

        # add an index to the shared property
        result.append([(CompositeStorage.SHARED, 1)])

        return result

    def get(
        self, path, ids=None, _filter=None, limit=0, skip=0, sort=None,
        *args, **kwargs
    ):

        query = path.copy()
        if _filter is not None:
            query.update(_filter)

        result = self.get_elements(
            ids=ids, query=query, limit=limit, skip=skip, sort=sort)

        return result

    def find(self, path, _filter, limit=0, skip=0, sort=None):

        result = self.get(
            path=path, _filter=_filter, limit=limit, skip=skip, sort=sort)

        return result

    def put(self, path, _id, data, *args, **kwargs):

        # get unique id
        _id = self.get_absolute_path(path=path, _id=_id)

        # get query such as a sum of path and _id
        query = path.copy()
        query[MongoStorage.ID] = _id

        _set = {
            '$set': data
        }

        self._update(_id=query, document=_set, multi=False)

    def remove(self, path, ids=None, *args, **kwargs):

        query = path.copy()

        parameters = {}

        if ids is not None:
            if isiterable(ids, is_str=False):
                query[CompositeStorage.ID] = {'$in': ids}
            else:
                parameters = {'justOne': 1}
                query[CompositeStorage.ID] = ids

        self._remove(document=query, **parameters)
