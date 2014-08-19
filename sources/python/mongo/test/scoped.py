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

from canopsis.mongo.scoped import MongoScopedStorage


class MongoScopedStorageTest(TestCase):
    """
    MongoScopedStorage UT on data_scope = "test_store"
    """

    def setUp(self):
        # create a storage on test_store collections
        self.scope = ['a', 'b', 'c']
        self.storage = MongoScopedStorage(
            data_scope="test", safe=True, scope=self.scope)

    def test_connect(self):
        self.assertTrue(self.storage.connected())

        self.storage.disconnect()

        self.assertFalse(self.storage.connected())

        self.storage.connect()

        self.assertTrue(self.storage.connected())

    def test_CRUD(self):

        self.storage.drop()

        request = self.storage.get()
        self.assertEqual(len(request), 0)

        self.storage.put(_id=0, data_type=0, data={'1': 2})
        self.storage.put(_id=1, data_type=0, data={'1': 2})
        self.storage.put(_id=2, data_type=1, data={'1': 2})
        self.storage.put(_id=3, data_type=1, data={'1': 2})

        request = self.storage.get()
        self.assertEqual(len(request), 4)

        request = self.storage.get(data_type=0)
        self.assertEqual(len(request), 2)

        request = self.storage.get(ids=[0])
        self.assertEqual(len(request), 1)

        self.storage.remove(data_type=1)
        request = self.storage.get()
        self.assertEqual(len(request), 2)

        self.storage.remove(ids=[0])
        request = self.storage.get()
        self.assertEqual(len(request), 1)

        request = self.storage.get(ids=[1])
        self.assertEqual(len(request), 1)

        request = self.storage.get(ids=[0])
        self.assertEqual(len(request), 0)

        request = self.storage.get(data_type=0)
        self.assertEqual(len(request), 1)
        self.storage.remove()

        request = self.storage.get()
        self.assertEqual(len(request), 0)

if __name__ == '__main__':
    main()
