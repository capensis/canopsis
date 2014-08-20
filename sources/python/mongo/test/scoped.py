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

from canopsis.storage.scoped import ScopedStorage
from canopsis.mongo.scoped import MongoScopedStorage


class MongoScopedStorageTest(TestCase):
    """
    MongoScopedStorage UT on data_scope = "test_store"
    """

    def setUp(self):
        # create a storage on test_store collections
        self.scope = ['a', 'b', 'c']
        self.storage = MongoScopedStorage(
            data_scope="test", scope=self.scope, safe=True)

    def test_connect(self):
        self.assertTrue(self.storage.connected())

        self.storage.disconnect()

        self.assertFalse(self.storage.connected())

        self.storage.connect()

        self.assertTrue(self.storage.connected())

    def test_CRUD(self):

        self.storage.drop()

        indexes = self.storage.indexes
        # check number of indexes
        self.assertEqual(len(indexes), len(self.scope) + 1)

        # add scope data with name which are scope last value
        for n, _ in enumerate(self.scope):
            # get prefixed scope
            _scope = {scope: scope for scope in self.scope[:n + 1]}

            # data name is just n
            name = str(n)

            # compare absolute path
            absolute_path = self.storage.get_absolute_path(
                scope=_scope, _id=name)
            _absolute_path = reduce(
                lambda x, y: '%s%s%s' % (x, ScopedStorage.SCOPE_SEPARATOR, y),
                _scope)
            _absolute_path = '%s%s%s' % (
                _absolute_path, ScopedStorage.SCOPE_SEPARATOR, name)
            _absolute_path = '%s%s' % (
                ScopedStorage.SCOPE_SEPARATOR, _absolute_path)
            self.assertEqual(absolute_path, _absolute_path)

            # put new entry
            self.storage.put(_scope, name, {'value': n})

        # get all data related to scope[n-1]
        for n, _ in enumerate(self.scope):
            _scope = {scope: scope for scope in self.scope[:n + 1]}
            elements = self.storage.get(scope=_scope)
            self.assertEqual(len(elements), len(self.scope) - n)

        for n in range(len(self.scope), 0, -1):
            _scope = {scope: scope for scope in self.scope[:n]}
            self.storage.remove(scope=_scope)
            elements = self.storage.get(scope={self.scope[0]: self.scope[0]})
            self.assertEqual(len(elements), n - 1)

if __name__ == '__main__':
    main()
