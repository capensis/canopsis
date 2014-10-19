#!/usr/bin/env python
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

from unittest import TestCase, main

from functools import reduce

from canopsis.storage import Storage
from canopsis.storage.composite import CompositeStorage
from canopsis.mongo.composite import MongoCompositeStorage


class MongoCompositeStorageTest(TestCase):
    """
    MongoCompositeStorage UT on data_scope = "test"
    """

    def setUp(self):
        # create a storage on test_store collections
        self.path = ['a', 'b', 'c']
        self.storage = MongoCompositeStorage(
            data_scope="test", path=self.path, safe=True)

    def test_connect(self):
        self.assertTrue(self.storage.connected())

        self.storage.disconnect()

        self.assertFalse(self.storage.connected())

        self.storage.connect()

        self.assertTrue(self.storage.connected())

    def test_CRUD(self):

        self.storage.drop()

        indexes = self.storage.all_indexes()
        # check number of indexes
        self.assertEqual(len(indexes), len(self.path) + 3)

        # add path data with name which are path last value
        for n, _ in enumerate(self.path):
            # get prefixed path
            _path = {path: path for path in self.path[:n + 1]}

            # data name is just n
            name = str(n)

            # compare absolute path
            absolute_path = self.storage.get_absolute_path(
                path=_path, data_id=name)
            __path = [path for path in self.storage.path if path in _path]
            _absolute_path = reduce(
                lambda x, y: '%s%s%s' % (
                    x, CompositeStorage.PATH_SEPARATOR, y), __path)
            _absolute_path = '%s%s%s' % (
                _absolute_path, CompositeStorage.PATH_SEPARATOR, name)
            _absolute_path = '%s%s' % (
                CompositeStorage.PATH_SEPARATOR, _absolute_path)
            self.assertEqual(absolute_path, _absolute_path)
            # put new entry
            self.storage.put(path=_path, data_id=name, data={'value': n})

        # get all data related to path[n-1]
        for n, _ in enumerate(self.path):
            _path = {path: path for path in self.path[:n + 1]}
            elements = self.storage.get(path=_path)
            self.assertEqual(len(elements), len(self.path) - n)

        for n in range(len(self.path), 0, -1):
            _path = {path: path for path in self.path[:n]}
            self.storage.remove(path=_path)
            elements = self.storage.get(path={self.path[0]: self.path[0]})
            self.assertEqual(len(elements), n - 1)

    def test_shared(self):

        self.storage.drop()

        path = {_path: _path for _path in self.path}

        data = []

        for n in range(10):
            d = path.copy()
            d[Storage.DATA_ID] = n
            data.append(d)

        # check unary data sharing
        for n, d in enumerate(data):
            self.storage.put(path=path, data_id=n, data=d)
            ds = self.storage.get(path=path, data_ids=n, shared=True)
            self.assertEqual(len(ds), 1)

            shared_id = self.storage.share_data(data=d)
            shared_data = self.storage.get_shared_data(shared_ids=shared_id)
            self.assertEqual(len(shared_data), 1)
            self.assertTrue(isinstance(shared_data[0], dict))

            self.storage.unshare_data(d)
            self.assertNotEqual(shared_id, d[CompositeStorage.SHARED])

            shared_data = self.storage.get_shared_data(shared_ids=shared_id)
            self.assertEqual(len(shared_data), 0)

        shared_id = 1

        # check shared data
        for n, d in enumerate(data):
            self.storage.share_data(data=d, shared_id=shared_id)
            shared_data = self.storage.get_shared_data(shared_ids=shared_id)
            self.assertEqual(len(shared_data), n + 1)

            shared_id += 1
            self.storage.share_data(
                data=d, shared_id=shared_id, share_extended=True)
            shared_data = self.storage.get_shared_data(shared_ids=shared_id)
            self.assertEqual(len(shared_data), n + 1)

        # unshare the first data and check if number of shared data is data - 1
        self.storage.unshare_data(data=data[0])
        shared_data = self.storage.get_shared_data(shared_ids=shared_id)
        self.assertEqual(len(shared_data), len(data) - 1)

if __name__ == '__main__':
    main()
