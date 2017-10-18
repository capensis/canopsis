#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

from influxdb import InfluxDBClient
from time import time
import unittest

try:
    from io import StringIO
except ImportError:
    from StringIO import StringIO

from canopsis.common.influx_store import InfluxStore
from canopsis.logger import Logger, OutputStream


class TestInfluxStore(unittest.TestCase):
    """
    Test case for InfluxStore.
    """

    @classmethod
    def setUpClass(self):
        self.db = 'canopsis_test'
        self.logger = Logger.get('test', StringIO(), OutputStream)

        self.conf = {
            InfluxStore.CONF_CAT: {
                'database': self.db,
                'host': InfluxStore.DEFAULT_HOST,
                'udp_port': InfluxStore.DEFAULT_PORT,
                'username': InfluxStore.DEFAULT_USER,
                'password': InfluxStore.DEFAULT_PWD
            }
        }

        self.client = InfluxDBClient(**self.conf[InfluxStore.CONF_CAT])
        self.client.create_database(self.db)

        self.is_ = InfluxStore(logger=self.logger, config=self.conf).client

        self.points = [
            {
                'measurement': 'cpu_load_short',
                #'time': '2017-09-15T10:00:00Z',
                'time': int(time() * 1e9),
                'fields': {
                    'value': 0.42,
                }
            }
        ]
        self.tags = {
            'duncan': 'idaho'
        }

    @classmethod
    def tearDownClass(self):
        """Teardown"""
        self.client.drop_database(self.db)

    def test_get_set(self):
        res = self.is_.get_list_database()
        self.assertTrue(self.db in [r['name'] for r in res])

        res = self.is_.write_points(self.points, tags=self.tags)
        self.assertTrue(res)

        query = 'select value from cpu_load_short;'
        res = self.is_.query(query)
        for point in res.get_points():
            self.assertEqual(point['value'], 0.42)
            break

if __name__ == '__main__':
    unittest.main()
