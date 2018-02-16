#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

import unittest

from canopsis.common import root_path

from canopsis.webcore.services.weather import WatcherFilter

import xmlrunner


class TestWatcherFilter(unittest.TestCase):

    def test_filters(self):
        doc1 = {
            "$and": [
                {
                    "hasallactivepbehaviorinentities": False
                },
                {
                    "$or": [
                        {
                            "hasallactivepbehaviorinentities": False
                        },
                    ],
                },
                {
                    "$or": [
                        {
                            "hasallactivepbehaviorinentities": False
                        },
                        {"IWillBeBack":2}
                    ],
                },
                {
                    "SarahConnor": {
                        "hasallactivepbehaviorinentities": True
                    },
                    "hasallactivepbehaviorinentities": True
                }
            ]
        }
        fdoc1 = {'$and': [{'$or': [{'IWillBeBack': 2}]}]}

        doc2 = {
            "$and": [
                {"SarahConnor":{"$eq": 'Terminated'}},
                {"hasactivepbehaviorinentities": True},
            ]
        }
        fdoc2 = {'$and': [{'SarahConnor': {'$eq': 'Terminated'}}]}

        doc3 = {
            "$and": [
                {},
            ]
        }
        fdoc3 = {'$and': [{}]}

        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc1), fdoc1)
        self.assertTrue(wf.all())
        self.assertIsNone(wf.some())

        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc2), fdoc2)
        self.assertIsNone(wf.all())
        self.assertTrue(wf.some())

        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc3), fdoc3)
        self.assertIsNone(wf.all())
        self.assertIsNone(wf.some())


if __name__ == '__main__':
    output = root_path + "/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
