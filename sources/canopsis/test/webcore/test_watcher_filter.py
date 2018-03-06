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
                    "active_pb_all": False
                },
                {
                    "$or": [
                        {
                            "active_pb_all": False
                        },
                    ],
                },
                {
                    "$or": [
                        {
                            "active_pb_all": False
                        },
                        {"IWillBeBack":2}
                    ],
                },
                {
                    "SarahConnor": {
                        "active_pb_all": True
                    },
                    "active_pb_all": True
                }
            ]
        }
        fdoc1 = {'$and': [{'$or': [{'IWillBeBack': 2}]}]}

        doc2 = {
            "$and": [
                {"SarahConnor":{"$eq": 'Terminated'}},
                {"active_pb_some": True},
            ]
        }
        fdoc2 = {'$and': [{'SarahConnor': {'$eq': 'Terminated'}}]}

        doc3 = {
            "$and": [
                {},
            ]
        }
        fdoc3 = {'$and': [{}]}

        doc4 = {
            "$and": [
                {"T-800": {"$contains": None}}
            ]
        }
        fdoc4 = {'$and': [{'T-800': {'$contains': None}}]}

        doc5 = {
            "$and": [
                {"active_pb_some": True}
            ]
        }

        doc6 = {
            "$and": [
                {"active_pb_some": True},
                {"active_pb_all": False}
            ]
        }

        doc7 = {}

        doc8 = {
            "$and": [
                {"active_pb_type": "cotcot"},
                {"active_pb_type": "CooooT"} # case is important in this test, do not modify.
            ]
        }
        fdoc8 = {}

        doc9 = {
            "$and": [
                {"active_pb_watcher": True}
            ]
        }
        fdoc9 = {}

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

        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc4), fdoc4)
        self.assertIsNone(wf.all())
        self.assertIsNone(wf.some())

        wf = WatcherFilter()
        wf.filter(doc5)
        self.assertTrue(wf.appendable(False, True, False))
        self.assertTrue(wf.appendable(True, True, False))
        self.assertFalse(wf.appendable(True, False, False))
        self.assertFalse(wf.appendable(False, False, False))

        wf = WatcherFilter()
        wf.filter(doc6)
        self.assertTrue(wf.appendable(False, True, False))
        self.assertFalse(wf.appendable(False, False, False))
        self.assertFalse(wf.appendable(True, False, False))
        self.assertFalse(wf.appendable(True, True, False))

        wf = WatcherFilter()
        wf.filter(doc7)
        self.assertTrue(wf.appendable(False, True, False))
        self.assertTrue(wf.appendable(False, False, False))
        self.assertTrue(wf.appendable(True, False, False))
        self.assertTrue(wf.appendable(True, True, False))

        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc8), fdoc8)
        self.assertTrue(wf.appendable(True, True, False, )) # no pb type given, all supported
        self.assertTrue(wf.appendable(True, True, False, pb_types=[]))
        self.assertTrue(wf.appendable(True, True, False, pb_types=["CooooT"]))
        self.assertTrue(wf.appendable(True, True, False, pb_types=["cOOOOt"]))
        self.assertTrue(wf.appendable(True, True, False, pb_types=["cOtcOt", "cOOOOt"]))
        self.assertTrue(wf.appendable(True, True, False, pb_types=["cUtcUt", "cOtcOt", "cOOOOt"]))
        self.assertFalse(wf.appendable(True, True, False, pb_types=["cUtcUt"]))
        self.assertFalse(wf.appendable(True, True, False, pb_types=["Courou"]))
        with self.assertRaises(ValueError):
            wf.appendable(True, True, "haha-nelson.com")

        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc9), fdoc9)
        self.assertTrue(wf.appendable(True, True, True))
        self.assertFalse(wf.appendable(True, True, False))

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
