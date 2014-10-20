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

    def all_indexes(self, *args, **kwargs):

        result = super(MongoCompositeStorage, self).all_indexes(
            *args, **kwargs)

        # add all sub_paths concatened with id
        for n, _ in enumerate(self.path):

            sub_path = self.path[:n + 1]
            index = [(path_name, Storage.ASC) for path_name in sub_path]
            index.append((Storage.DATA_ID, Storage.ASC))

            result.append(index)

        return result

    def get(
        self,
        path, data_ids=None, _filter=None, shared=False,
        limit=0, skip=0, sort=None, with_count=False,
        *args, **kwargs
    ):

        result = []

        # create a get query which is a copy of input path plus _filter
        query = path.copy()
        if _filter is not None:
            query.update(_filter)

        one_result = False

        # if data_ids is given
        if data_ids is not None:

            # add absolute pathes into ids
            if isiterable(data_ids, is_str=False):
                query[Storage.DATA_ID] = {"$in": data_ids}
            else:
                query[Storage.DATA_ID] = data_ids
                one_result = True

        # get elements
        result = self.find_elements(
            query=query,
            limit=limit,
            skip=skip,
            sort=sort,
            with_count=with_count)

        if with_count:
            count = result[1]
            result = result[0]

        if result is not None and shared:
            # if result is one value
            if isinstance(result, dict):
                # if result is shared
                if CompositeStorage.SHARED in result:
                    shared_id = result[CompositeStorage.SHARED]
                    # result equals shared data
                    result = [self.get_shared_data(shared_ids=shared_id)]
                else:
                    # else, result is a list of itself
                    result = [[result]]

            else:
                # if result is a list of data, use data_to_extend
                data_to_extend = result[:]
                # and initialize result such as an empty list
                result = []
                # for all data in result
                for data in data_to_extend:
                    # if data is shared, get shared data
                    if CompositeStorage.SHARED in data:
                        shared_id = data[CompositeStorage.SHARED]
                        data = self.get_shared_data(shared_ids=shared_id)
                    else:
                        data = [data]
                    # append data in result
                    result.append(data)

        if result is not None and one_result:
            if result:
                result = result[0]
            else:
                result = None

        if with_count:
            result = result, count

        return result

    def find(
        self,
        path,
        _filter, shared=False, limit=0, skip=0, sort=None, with_count=False
    ):

        result = self.get(
            path=path, _filter=_filter, shared=shared,
            limit=limit, skip=skip, sort=sort, with_count=with_count)

        return result

    def put(self, path, data_id, data, share_id=None, *args, **kwargs):

        # get unique id
        _id = self.get_absolute_path(path=path, data_id=data_id)

        data_to_put = data.copy()

        if share_id is not None:
            data_to_put[CompositeStorage.SHARED] = share_id

        query = {MongoStorage.ID: _id}
        query.update(path)
        query[Storage.DATA_ID] = data_id

        _set = {
            '$set': data_to_put
        }

        self._update(_id=query, document=_set, multi=False)

    def remove(self, path, data_ids=None, shared=False, *args, **kwargs):

        query = path.copy()

        parameters = {}

        if data_ids is not None:
            if isiterable(data_ids, is_str=False):
                query[Storage.DATA_ID] = {'$in': data_ids}
            else:
                parameters = {'justOne': 1}
                query[Storage.DATA_ID] = data_ids

        self._remove(document=query, **parameters)

        # remove extended data
        if shared:
            _ids = []
            data_to_remove = self.get(
                path=path, data_ids=data_ids, shared=True)
            for dtr in data_to_remove:
                path, data_id = self.get_path_with_id(dtr)
                extended = self.get(path=path, data_id=data_id, shared=True)
                _ids.append([data[MongoStorage.ID] for data in extended])
            document = {MongoStorage.ID: {'$in': _ids}}
            self._remove(document=document)
