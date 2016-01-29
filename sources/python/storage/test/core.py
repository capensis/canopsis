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

from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.storage.core import Storage

from .base import BaseTestConfiguration, BaseStorageTest

from tempfile import NamedTemporaryFile


@conf_paths('storage/test-core.conf')
class TestConfiguration(BaseTestConfiguration):
    """Default test configuration."""


class StorageTest(BaseStorageTest):

    def _testconfcls(self):

        return TestConfiguration

    def _test(self, storage):

        self._test_connection(storage)
        self.test_indexes(storage)
        self.test_CRUD(storage)
        self.test_cache(storage)
        self.test_thread(storage)

    def _test_connection(self, storage):

        storage.connect()
        self.assertTrue(storage.connected())
        storage.disconnect()
        self.assertFalse(storage.connected())
        storage.reconnect()
        self.assertTrue(storage.connected())

    def _test_indexes(self, storage):

        indexes = storage.all_indexes()

        collection_index = storage._backend.index_information()

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
            _file.write('[{}]'.format('STORAGE'))
            _file.write('\nindexes={}'.format(indexes))
            _file.write('\ndata={}'.format(data))

        storage.apply_configuration(conf_paths=conf_path)

        _indexes = storage.all_indexes()

        self.assertEqual(final_indexes, _indexes)

    def _test_CRUD(self, storage):

        document = {'a': 'b'}

        storage.drop()

        request = storage._find(document)
        self.assertEqual(request.count(), 0)

        storage._insert(document)

        request = storage._find(document)

        self.assertEqual(request.count(), 1)
        self.assertEqual(request[0], document)

        storage._remove(document)

        request = storage._find(document)
        self.assertEqual(request.count(), 0)

        request = storage._find()
        self.assertEqual(request.count(), 0)

        _id = 'test'
        document['_id'] = _id

        request = storage.get_elements(ids=[_id])

        self.assertEqual(len(request), 0)

        storage.put_element(_id=_id, element=document)

        request = storage.get_elements()
        self.assertEqual(len(request), 1)
        self.assertEqual(request[0], document)

        request = storage.get_elements(ids=[_id])
        self.assertEqual(len(request), 1)
        self.assertEqual(request[0], document)

        storage.remove_elements(ids=[_id])

        request = storage.get_elements(ids=[_id])

        self.assertEqual(len(request), 0)

    def _test_cache(self, storage):

        cache_size = 25
        count = cache_size + 1  # play with 2 * cache_size + 1

        # activate cache
        storage.cache_autocommit = 0.1
        storage.cache_size = cache_size
        data = list({'_id': str(i)} for i in range(count))

        # check insert
        for i, d in enumerate(data):
            storage._insert(d, cache=True)
            count = storage._count(document=d)
            # documents are inserted when cache_size insertions are done
            self.assertEqual(count, 0 if i != cache_size - 1 else 1)

        for d in data[:-1]:
            count = storage._count(document=d)
            self.assertEqual(count, 1)
        count = storage._count(document=data[-1])
        self.assertEqual(count, 0)
        storage.execute_cache()
        count = storage._count(document=data[-1])
        self.assertEqual(count, 1)
        self.assertEqual(storage._cache_count, 0)

        # check update
        for i, d in enumerate(data):
            storage._update(
                spec=d, document={'$set': {'a': 1}}, cache=True
            )
            elt = storage._find(document=d)[0]

            # documents are updated when cache_size updates are done
            if i != cache_size - 1:
                self.assertNotIn('a', elt)
            else:
                self.assertIn('a', elt)

        for d in data[:-1]:
            elt = storage._find(document=d)[0]
            self.assertIn('a', elt)

        elt = storage._find(document=data[-1])[0]
        self.assertNotIn('a', elt)
        storage.execute_cache()
        elt = storage._find(document=data[-1])[0]
        self.assertIn('a', elt)

        # check remove
        for i, d in enumerate(data):
            storage._remove(document=d, cache=True)
            count = storage._count(document=d)
            # documents are removed when cache_size removes are done
            self.assertEqual(count, 1 if i != cache_size - 1 else 0)

        for d in data[:-1]:
            count = storage._count(document=d)
            self.assertEqual(count, 0)

        count = storage._count(document=data[-1])
        self.assertEqual(count, 1)
        storage.execute_cache()
        count = storage._count(document=data[-1])
        self.assertEqual(count, 0)

    def _test_thread(self, storage):

        # ensure cache auto commit is short
        storage._cache_autocommit = 0.01
        # check the cache thread is started
        storage.remove_elements(cache=True)
        self.assertTrue(storage._cached_thread.isAlive())

        # check the cache thread is halted
        storage.halt_cache_thread()
        self.assertFalse(storage._cached_thread.isAlive())

        # check the cache thread is started twice
        storage.remove_elements(cache=True)
        self.assertTrue(storage._cached_thread.isAlive())


if __name__ == '__main__':
    main()
