#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

from time import sleep
import configparser
import os
import unittest
from canopsis.common import root_path
from canopsis.common.ethereal_data import EtherealData
from canopsis.common.mongo_store import MongoStore, DEFAULT_DB_NAME
from canopsis.confng import Configuration, Ini
import xmlrunner

DEFAULT_CONF_FILE = "etc/common/mongo_store.conf"


class TestEtherealData(unittest.TestCase):

    def setUp(self):
        config = configparser.RawConfigParser()
        config.read(os.path.join(root_path, DEFAULT_CONF_FILE))
        self.conf = {
            MongoStore.CONF_CAT: {
                "host": config["DATABASE"]["host"],
                "port": config["DATABASE"]["port"],
                'db': config["DATABASE"]["db"],
                'user': config["DATABASE"]["user"],
                'pwd': config["DATABASE"]["pwd"]
            }
        }
        self.client = MongoStore(config=self.conf)
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
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
