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

from canopsis.mongo import MongoDataBase, MongoStorage

from tempfile import NamedTemporaryFile


class DataBaseTest(TestCase):

    def setUp(self):
        self.database = MongoDataBase(
            data_scope="test_store", auto_connect=False)

    def test_connect(self):
        self.database.connect()
        self.assertTrue(self.database.connected())
        self.database.disconnect()
        self.assertFalse(self.database.connected())
        self.database.reconnect()
        self.assertTrue(self.database.connected())


class TestStorage(MongoStorage):

    def _get_storage_type(self, *args, **kwargs):
        return 'test'


class StorageTest(TestCase):

    def setUp(self):
        self.storage = TestStorage(data_type='test')

    def tearDown(self):
        self.storage.drop()

    def test_connect(self):
        self.assertTrue(self.storage.connected())
        self.storage.disconnect()
        self.assertFalse(self.storage.connected())
        self.storage.reconnect()
        self.assertTrue(self.storage.connected())

    def test_indexes(self):
        indexes = self.storage.all_indexes()

        collection_index = self.storage._backend.index_information()

        for index in indexes:
            key = ''
            for i in index:
                key += '{0}_{1}'.format(i[0], i[1] if i[0] != '_id' else '')
            self.assertIn(key, collection_index)

        conf_path = NamedTemporaryFile().name

        indexes = ['a', 'b']

        with open(conf_path, 'w') as _file:
            _file.write('[%s]' % MongoStorage.CATEGORY)
            _file.write('\nindexes=%s' % (indexes))

        self.storage.apply_configuration(conf_paths=conf_path)

        _indexes = self.storage.indexes

        self.assertEqual(indexes, _indexes)

    def test_CRUD(self):
        document = {'a': 'b'}

        self.storage.drop()

        request = list(self.storage._find(document))
        self.assertEqual(len(request), 0)

        self.storage._insert(document)

        request = list(self.storage._find(document))

        self.assertEqual(len(request), 1)
        self.assertEqual(request[0], document)

        self.storage._remove(document)

        request = list(self.storage._find(document))
        self.assertEqual(len(request), 0)

        request = list(self.storage._find())
        self.assertEqual(len(request), 0)

        _id = 'test'
        document['_id'] = _id

        request = self.storage.get_elements(ids=[_id])

        self.assertEqual(len(request), 0)

        self.storage.put_element(_id=_id, element=document)

        request = self.storage.get_elements()
        self.assertEqual(len(request), 1)
        self.assertEqual(request[0], document)

        request = self.storage.get_elements(ids=[_id])
        self.assertEqual(len(request), 1)
        self.assertEqual(request[0], document)

        self.storage.remove_elements(ids=[_id])

        request = self.storage.get_elements(ids=[_id])

        self.assertEqual(len(request), 0)

if __name__ == '__main__':
    main()
