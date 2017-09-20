#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

from time import sleep
import unittest

from canopsis.common.ethereal_data import EtherealData
from canopsis.common.mongo_store import MongoStore, DEFAULT_DB_NAME
from canopsis.confng import Configuration, Ini


class TestEtherealData(unittest.TestCase):

    def setUp(self):
        self.conf = {
            MongoStore.CONF_CAT: {
                'db': DEFAULT_DB_NAME
            }
        }
        self.cred_conf = Configuration.load(MongoStore.CRED_CONF_PATH, Ini)
        self.client = MongoStore(config=self.conf,
                                 cred_config=self.cred_conf)
        self.collection = self.client.get_collection('any_collection')

        self.ed = EtherealData(collection=self.collection,
                               filter_={},
                               timeout=2)

    def tearDown(self):
        """Teardown"""
        self.collection.drop()

    def test_get_set(self):
        self.assertIsNone(self.ed.get('mario'))
        self.ed.set('mario', 'bros')
        self.assertEqual(self.ed.get('mario'), 'bros')

    def test_cache(self):
        self.ed.set('sonic', 'hedgehog')
        self.assertEqual(self.ed.get('sonic'), 'hedgehog')

        self.collection.update({}, {'sonic': 'tails'})
        self.assertEqual(self.ed.get('sonic'), 'hedgehog')

        sleep(3)
        self.assertEqual(self.ed.get('sonic'), 'tails')

if __name__ == '__main__':
    unittest.main()
