#!/usr/bin/env python
# -*- coding: utf-8  -*-

from __future__ import unicode_literals

import unittest

from canopsis.common.event import Event


class TestEvent(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.connector = 'connector'
        self.connector_name = 'connector_name'
        self.component = 'component'
        self.resource = 'resource'

        self.ev_ = Event(connector=self.connector,
                         connector_name=self.connector_name,
                         component=self.component,
                         resource=self.resource)

    @classmethod
    def tearDownClass(self):
        """Teardown"""

    def test_to_dict(self):
        res = self.ev_.to_dict()
        expected = {
            'component': 'component',
            'connector': 'connector',
            'connector_name': 'connector_name',
            'resource': 'resource'
        }

        self.assertDictEqual(res, expected)

if __name__ == '__main__':
    unittest.main()
