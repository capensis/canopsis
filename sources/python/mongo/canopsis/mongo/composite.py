# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

from canopsis.common.init import basestring
from canopsis.common.utils import isiterable
from canopsis.mongo.core import MongoStorage
from canopsis.storage.core import Storage
from canopsis.storage.composite import CompositeStorage


class MongoCompositeStorage(MongoStorage, CompositeStorage):

    def all_indexes(self, *args, **kwargs):

        result = super(MongoCompositeStorage, self).all_indexes(
            *args, **kwargs
        )

        # add all sub_paths concatened with id
        for n, _ in enumerate(self.path):

            sub_path = self.path[:n + 1]
            index = [(path_name, Storage.ASC) for path_name in sub_path]
            index.append((CompositeStorage.NAME, Storage.ASC))

            result.append(index)

        result.append([(CompositeStorage.NAME, Storage.ASC)])

        return result

    def get(
            self,
            path, names=None, _filter=None, shared=False,
            limit=0, skip=0, sort=None, with_count=False,
            *args, **kwargs
    ):

        result = []

        # create a get query which is a copy of input path plus _filter
        query = path.copy()
        if _filter is not None:
            query.update(_filter)
        # check if only one element is asked
        one_result = isinstance(names, basestring)
        # if names is given
        if names is not None:
            # if one element is asked
            if one_result:
                query[CompositeStorage.NAME] = names
            else:
                query[CompositeStorage.NAME] = {"$in": names}

        # get the right hint
        self_path = self.path
        hint = []

        for _path in self_path:
            if _path in path:
                hint.append((_path, 1))

        hint.append((CompositeStorage.NAME, 1))

        # get elements
        result = self.find_elements(
            query=query,
            limit=limit,
            skip=skip,
            sort=sort,
            with_count=with_count,
            hint=hint
        )
        if with_count:
            count = result[1]
            result = result[0]

        if result is not None and shared:
            # if result is a list of data, use data_to_extend
            result, data_to_extend = [], result
            # for all data in result
            for data in data_to_extend:

                # if data is shared, get shared data
                if CompositeStorage.SHARED in data:
                    shared_id = data[CompositeStorage.SHARED]
                    data = list(self.get_shared_data(shared_ids=shared_id))

                else:
                    data = [data]
                # append data in result
                result.append(data)
        if result is not None and one_result:
            result = result[0] if result else [] if shared else None

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
            limit=limit, skip=skip, sort=sort, with_count=with_count
        )

        return result

    def put(
            self, path, name, data, share_id=None, cache=False, *args, **kwargs
    ):

        # get unique id
        _id = self.get_absolute_path(path=path, name=name)

        data_to_put = data.copy()

        if share_id is not None:
            data_to_put[CompositeStorage.SHARED] = share_id

        query = {MongoStorage.ID: _id}
        query.update(path)
        query[CompositeStorage.NAME] = name

        _set = {
            '$set': data_to_put
        }

        self._update(spec=query, document=_set, multi=False, cache=cache)

    def remove(
            self, path, names=None, shared=False, cache=False, *args, **kwargs
    ):

        query = path.copy()

        parameters = {}

        if names is not None:
            if isiterable(names, is_str=False):
                query[CompositeStorage.NAME] = {'$in': names}

            else:
                parameters = {'justOne': 1}
                query[CompositeStorage.NAME] = names

        self._remove(document=query, cache=cache, **parameters)

        # remove extended data
        if shared:
            _ids = []
            data_to_remove = self.get(path=path, names=names, shared=True)

            for dtr in data_to_remove:
                path, name = self.get_path_with_name(data=dtr)
                extended = self.get(path=path, name=name, shared=True)
                _ids.append([data[MongoStorage.ID] for data in extended])

            document = {MongoStorage.ID: {'$in': _ids}}
            self._remove(document=document, cache=cache)
