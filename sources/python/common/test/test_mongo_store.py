#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

from pymongo import MongoClient
from pymongo.collection import Collection
import unittest

from canopsis.common.mongo_store import MongoStore
from canopsis.confng import Configuration, Ini


class TestMongoStore(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.db_name = 'canopsis'
        self.coll = 'test_mongostorage'
        self.collection = MongoClient()[self.db_name][self.coll]

        self.conf = {
            MongoStore.CONF_CAT: {
                'db': self.db_name
            }
        }
        self.mid_conf = {MongoStore.MID_CONF_CAT: {}}
        self.cred_conf = Configuration.load(MongoStore.CRED_CONF_PATH, Ini)

        self.ms = MongoStore(config=self.conf,
                             mid_config=self.mid_conf,
                             cred_config=self.cred_conf)

    @classmethod
    def tearDownClass(self):
        """Teardown"""
        self.collection.drop()

    def tearDown(self):
        self.collection.remove()

    def test_get_collection(self):
        coll = self.ms.get_collection(self.coll)
        self.assertTrue(isinstance(coll, Collection))
        self.assertEqual(coll.full_name, '{}.{}'.format(self.db_name,
                                                        self.coll))

if __name__ == '__main__':
    unittest.main()
