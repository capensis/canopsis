#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

from pymongo.collection import Collection
import unittest

from canopsis.common.mongo_store import MongoStore
from canopsis.confng import Configuration, Ini


class TestMongoStore(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.db_name = 'canopsis'
        self.collection_name = 'test_mongostorage'

        self.conf = {
            MongoStore.CONF_CAT: {
                'db': self.db_name
            }
        }
        self.cred_conf = Configuration.load(MongoStore.CRED_CONF_PATH, Ini)

        self.ms = MongoStore(config=self.conf,
                             cred_config=self.cred_conf)
        self.collection = self.ms.get_collection(self.collection_name)

    @classmethod
    def tearDownClass(self):
        """Teardown"""
        self.collection.drop()

    def tearDown(self):
        self.collection.remove()

    def test_get_collection(self):
        coll = self.ms.get_collection(self.collection_name)
        self.assertTrue(isinstance(coll, Collection))
        self.assertEqual(coll.full_name, '{}.{}'.format(self.db_name,
                                                        self.collection_name))

if __name__ == '__main__':
    unittest.main()
