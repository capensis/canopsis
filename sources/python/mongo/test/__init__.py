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

from unittest import TestCase, main

from canopsis.storage.core import Storage
from canopsis.mongo.core import MongoDataBase, MongoStorage

from tempfile import NamedTemporaryFile


class DataBaseTest(TestCase):

    def setUp(self):
        self.database = MongoDataBase(
            data_scope="test_store", auto_connect=False
        )

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

        indexes = [
            'a',
            ('b', Storage.DESC),
            ['c', ('d', Storage.ASC)]
        ]
        data = {
            'e': '',
            'f': {
                Storage.KEY: Storage.DESC
            }
        }
        final_indexes = [
            [(Storage.DATA_ID, Storage.ASC)],
            [('a', Storage.ASC)],
            [('b', Storage.DESC)],
            [('c', Storage.ASC), ('d', Storage.ASC)],
            [('f', Storage.DESC)],
            [('_id', Storage.ASC)]
        ]

        with open(conf_path, 'w') as _file:
            _file.write('[{}]'.format(MongoStorage.CATEGORY))
            _file.write('\nindexes={}'.format(indexes))
            _file.write('\ndata={}'.format(data))

        self.storage.apply_configuration(conf_paths=conf_path)

        _indexes = self.storage.all_indexes()

        self.assertEqual(final_indexes, _indexes)

    def test_CRUD(self):
        document = {'a': 'b'}

        self.storage.drop()

        request = self.storage._find(document)
        self.assertEqual(request.count(), 0)

        self.storage._insert(document)

        request = self.storage._find(document)

        self.assertEqual(request.count(), 1)
        self.assertEqual(request[0], document)

        self.storage._remove(document)

        request = self.storage._find(document)
        self.assertEqual(request.count(), 0)

        request = self.storage._find()
        self.assertEqual(request.count(), 0)

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

    def test_cache(self):

        cache_size = 25
        count = cache_size + 1  # play with 2 * cache_size + 1

        # activate cache
        self.storage.cache_autocommit = 0.1
        self.storage.cache_size = cache_size
        data = list({'_id': str(i)} for i in range(count))

        # check insert
        for i, d in enumerate(data):
            self.storage._insert(d, cache=True)
            count = self.storage._count(document=d)
            # documents are inserted when cache_size insertions are done
            self.assertEqual(count, 0 if i != cache_size - 1 else 1)
        for d in data[:-1]:
            count = self.storage._count(document=d)
            self.assertEqual(count, 1)
        count = self.storage._count(document=data[-1])
        self.assertEqual(count, 0)
        self.storage.execute_cache()
        count = self.storage._count(document=data[-1])
        self.assertEqual(count, 1)
        self.assertEqual(self.storage._cache_count, 0)

        # check update
        for i, d in enumerate(data):
            self.storage._update(
                spec=d, document={'$set': {'a': 1}}, cache=True
            )
            elt = self.storage._find(document=d)[0]
            # documents are updated when cache_size updates are done
            if i != cache_size - 1:
                self.assertNotIn('a', elt)
            else:
                self.assertIn('a', elt)
        for d in data[:-1]:
            elt = self.storage._find(document=d)[0]
            self.assertIn('a', elt)
        elt = self.storage._find(document=data[-1])[0]
        self.assertNotIn('a', elt)
        self.storage.execute_cache()
        elt = self.storage._find(document=data[-1])[0]
        self.assertIn('a', elt)

        # check remove
        for i, d in enumerate(data):
            self.storage._remove(document=d, cache=True)
            count = self.storage._count(document=d)
            # documents are removed when cache_size removes are done
            self.assertEqual(count, 1 if i != cache_size - 1 else 0)
        for d in data[:-1]:
            count = self.storage._count(document=d)
            self.assertEqual(count, 0)
        count = self.storage._count(document=data[-1])
        self.assertEqual(count, 1)
        self.storage.execute_cache()
        count = self.storage._count(document=data[-1])
        self.assertEqual(count, 0)

    def test_thread(self):

        # ensure cache auto commit is short
        self.storage._cache_autocommit = 0.01
        # check the cache thread is started
        self.storage.remove_elements(cache=True)
        self.assertTrue(self.storage._cache_thread.isAlive())

        # check the cache thread is halted
        self.storage.halt_cache_thread()
        self.assertFalse(self.storage._cache_thread.isAlive())

        # check the cache thread is started twice
        self.storage.remove_elements(cache=True)
        self.assertTrue(self.storage._cache_thread.isAlive())

if __name__ == '__main__':
    main()
