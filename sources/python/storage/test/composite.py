#!/usr/bin/env python
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

from unittest import main

from functools import reduce

from canopsis.storage.composite import CompositeStorage

from canopsis.configuration.configurable.decorator import conf_paths

from .base import BaseTestConfiguration, BaseStorageTest


@conf_paths('storage/test-composite.conf')
class TestConfiguration(BaseTestConfiguration):
    """Default test configuration."""


class CompositeStorageTest(BaseStorageTest):
    """CompositeStorage UT on data_scope = "test" """

    def _test(self, storage):

        self.path = storage.path

        self._test_CRUD(storage)
        self._test_shared(storage)
        self._test_distinct(storage)

    def _test_CRUD(self, storage):

        storage.drop()

        indexes = storage.all_indexes()
        # check number of indexes
        self.assertEqual(len(indexes), len(self.path) + 4)

        # add path data with name which are path last value
        for n, _ in enumerate(self.path):
            # get prefixed path
            _path = {path: path for path in self.path[:n + 1]}

            # data name is just n
            name = str(n)

            # compare absolute path
            absolute_path = storage.get_absolute_path(
                path=_path, name=name
            )
            __path = [path for path in storage.path if path in _path]
            _absolute_path = reduce(
                lambda x, y: '%s%s%s' % (
                    x, CompositeStorage.PATH_SEPARATOR, y), __path
            )
            _absolute_path = '%s%s%s' % (
                _absolute_path, CompositeStorage.PATH_SEPARATOR, name
            )
            _absolute_path = '%s%s' % (
                CompositeStorage.PATH_SEPARATOR, _absolute_path
            )
            self.assertEqual(absolute_path, _absolute_path)
            # put new entry
            storage.put(path=_path, name=name, data={'value': n})

        # get all data related to path[n-1]
        for n, _ in enumerate(self.path):
            _path = {path: path for path in self.path[:n + 1]}
            elements = storage.get(path=_path)
            self.assertEqual(len(elements), len(self.path) - n)

        for n in range(len(self.path), 0, -1):
            _path = {path: path for path in self.path[:n]}
            storage.remove(path=_path)
            elements = storage.get(path={self.path[0]: self.path[0]})
            self.assertEqual(len(elements), n - 1)

    def _test_shared(self, storage):

        storage.drop()

        path = {_path: _path for _path in self.path}

        data = []

        for n in range(10):
            d = path.copy()
            d[CompositeStorage.NAME] = str(n)
            data.append(d)

        # check unary data sharing
        for n, d in enumerate(data):
            storage.put(path=path, name=str(n), data=d)
            ds = storage.get(path=path, names=str(n), shared=True)
            self.assertEqual(len(ds), 1)

            shared_id = storage.share_data(data=d)
            shared_data = storage.get_shared_data(shared_ids=shared_id)
            self.assertEqual(len(shared_data), 1)
            self.assertTrue(isinstance(shared_data[0], dict))

            storage.unshare_data(d)
            self.assertNotEqual(shared_id, d[CompositeStorage.SHARED])

            shared_data = storage.get_shared_data(
                shared_ids=str(shared_id)
            )
            self.assertEqual(len(shared_data), 0)

        shared_id = 1

        # check shared data
        for n, d in enumerate(data):
            storage.share_data(data=d, shared_id=str(shared_id))
            shared_data = storage.get_shared_data(
                shared_ids=str(shared_id)
            )
            self.assertEqual(len(shared_data), n + 1)
            shared_id += 1
            storage.share_data(
                data=d, shared_id=str(shared_id), share_extended=True
            )
            shared_data = storage.get_shared_data(
                shared_ids=str(shared_id)
            )
            self.assertEqual(len(shared_data), n + 1)

        # unshare the first data and check if number of shared data is data - 1
        storage.unshare_data(data=data[0])
        shared_data = storage.get_shared_data(shared_ids=str(shared_id))
        self.assertEqual(len(shared_data), len(data) - 1)

    def _test_distinct(self, storage):

        storage.put_element(
            _id='id_1', element={'group': 'a', 'test': 'value1'}
        )
        storage.put_element(
            _id='id_2', element={'group': 'a', 'test': 'value2'}
        )
        storage.put_element(
            _id='id_3', element={'group': 'a', 'test': 'value3'}
        )

        storage.put_element(
            _id='id_4', element={'group': 'b', 'test': 'value1'}
        )
        storage.put_element(
            _id='id_5', element={'group': 'b', 'test': 'value2'}
        )
        storage.put_element(
            _id='id_6', element={'group': 'b', 'test': 'value4'}
        )

        result = storage.distinct('test', {'group': 'b'})
        self.assertEqual(result, ['value1', 'value2', 'value4'])

        result = storage.distinct('test', {'group': 'a'})
        self.assertEqual(result, ['value1', 'value2', 'value3'])


if __name__ == '__main__':
    main()
