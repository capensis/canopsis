#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

import unittest

from canopsis.common.influx_store import InfluxStore
from canopsis.logger.logger import Logger
from canopsis.metrics.manager import MetricsManager


class TestMetricsManager(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.logger = Logger.get('metrics', MetricsManager.LOG_PATH)

        self.db = 'canopsis_test'
        self.conf = {
            InfluxStore.CONF_CAT: {
                'database': self.db,
                'host': InfluxStore.DEFAULT_HOST,
                'udp_port': InfluxStore.DEFAULT_PORT,
                'username': InfluxStore.DEFAULT_USER,
                'password': InfluxStore.DEFAULT_PWD
            }
        }
        self.store = InfluxStore(config=self.conf).client

        self.metrics = MetricsManager(logger=self.logger,
                                      store=self.store)

    @classmethod
    def tearDownClass(self):
        """Teardown"""
        self.store.drop_database(self.db)

    def test_get_set(self):
        res = self.metrics.get_all_metrics()
        print(res)

if __name__ == '__main__':
    unittest.main()
