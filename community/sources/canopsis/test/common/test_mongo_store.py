#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

import configparser
import unittest
import os
from pymongo.collection import Collection

from canopsis.common.mongo_store import MongoStore
from canopsis.confng import Configuration, Ini
from canopsis.common import root_path
import xmlrunner

DEFAULT_CONF_FILE = "etc/common/mongo_store.conf"


class TestMongoStore(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.db_name = 'canopsis'
        cls.collection_name = 'test_mongostorage'

        config = configparser.RawConfigParser()
        config.read(os.path.join(root_path, DEFAULT_CONF_FILE))
        cls.conf = {
            MongoStore.CONF_CAT: {
                "host": config["DATABASE"]["host"],
                "port": config["DATABASE"]["port"],
                'db': config["DATABASE"]["db"],
                'user': config["DATABASE"]["user"],
                'pwd': config["DATABASE"]["pwd"]
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
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
