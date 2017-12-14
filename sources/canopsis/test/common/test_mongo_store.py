#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

import unittest

from pymongo.collection import Collection

from canopsis.common.mongo_store import MongoStore


class TestMongoStore(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.db_name = 'canopsis'
        cls.collection_name = 'test_mongostorage'

        cls.conf = {
            MongoStore.CONF_CAT: {
                'db': cls.db_name,
                'user': 'cpsmongo',
                'pwd': 'canopsis'
            }
        }

        cls.ms = MongoStore(config=cls.conf)
        cls.collection = cls.ms.get_collection(cls.collection_name)

    @classmethod
    def tearDownClass(cls):
        """Teardown"""
        cls.collection.drop()

    def tearDown(self):
        self.collection.remove()

    def test_get_collection(self):
        coll = self.ms.get_collection(self.collection_name)
        self.assertTrue(isinstance(coll, Collection))
        self.assertEqual(coll.full_name, '{}.{}'.format(self.db_name,
                                                        self.collection_name))

if __name__ == '__main__':
    unittest.main()
