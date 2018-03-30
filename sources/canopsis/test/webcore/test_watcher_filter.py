#!/usr/bin/env python
# -*- coding: utf-8  -*-
from __future__ import unicode_literals

import unittest

from canopsis.common import root_path

from canopsis.watcher.filtering import WatcherFilter

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

        doc9 = {
            "$and": [
                {"active_pb_watcher": True}
            ]
        }
        fdoc9 = {}

        doc = {
            'active_pb_exclude_type': 'ti'
        }
        fdoc = {}
        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc), fdoc)
        # do not exclude unfiltered pb types
        self.assertTrue(wf.match(False, True, False, pb_types=['ti_unknown']))

        doc10 = {
            '$and': [
                {
                    'active_pb_type': 'pb_ti1',
                    'active_pb_include_type': 'pb_ti2',
                    'active_pb_exclude_type': 'pb_te1'
                }
            ]
        }
        fdoc10 = {}
        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc10), fdoc10)
        # no type check
        self.assertTrue(wf.match(True, True, False))
        # unknown pb type
        self.assertFalse(wf.match(False, True, False, pb_types=['pb_ti_unknown']))
        # type check
        self.assertTrue(wf.match(True, True, False, pb_types=['pb_ti1']))
        # type check 2
        self.assertTrue(wf.match(True, True, False, pb_types=['pb_ti2']))
        # type exclude and include, exclude prevails
        self.assertFalse(wf.match(True, True, False, pb_types=['pb_te1', 'pb_ti1']))
        self.assertFalse(wf.match(False, True, False, pb_types=['pb_te1', 'pb_ti1']))
        self.assertTrue(wf.match(True, True, False, pb_types=['pb_ti1']))
        # type with pbehavior watcher
        self.assertTrue(wf.match(False, False, True, pb_types=['pb_ti1']))
        # type with pbehavior entities
        self.assertTrue(wf.match(False, True, False, pb_types=['pb_ti1']))
        # type check without explicit active_pb_some|all implies active_pb_some=True
        self.assertTrue(wf.match(False, True, False, pb_types=['pb_ti1']))
        self.assertTrue(wf.match(False, True, True, pb_types=['pb_ti1']))
        self.assertTrue(wf.match(True, True, False, pb_types=['pb_ti1']))
        self.assertTrue(wf.match(False, False, True, pb_types=['pb_ti1']))
        self.assertTrue(wf.match(True, True, True, pb_types=['pb_ti1']))

        doc11 = {
            '$and': [
                {
                    'active_pb_type': 'pb_ti1',
                    'active_pb_watcher': True
                }
            ]
        }
        fdoc11 = {}

        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc11), fdoc11)
        # type with only pbehavior watcher True
        self.assertTrue(wf.match(False, False, True, pb_types=['pb_ti1']))
        self.assertTrue(wf.match(False, False, True))
        # type with only pbehavior watcher False
        self.assertFalse(wf.match(False, True, False, pb_types=['pb_ti1']))
        self.assertFalse(wf.match(False, True, False))


        wf = WatcherFilter()
        with self.assertRaises(ValueError):
            wf.match(True, False, False)

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
        self.assertTrue(wf.match(False, True, False))
        self.assertTrue(wf.match(True, True, False))
        self.assertFalse(wf.match(False, False, False))

        wf = WatcherFilter()
        wf.filter(doc6)
        self.assertTrue(wf.match(False, True, False))
        self.assertFalse(wf.match(False, False, False))
        self.assertFalse(wf.match(True, True, False))

        wf = WatcherFilter()
        wf.filter(doc7)
        self.assertTrue(wf.match(False, True, False))
        self.assertTrue(wf.match(False, False, False))
        self.assertTrue(wf.match(True, True, False))

        doc8 = {
            "$and": [
                {"active_pb_type": "cotcot"},
                {"active_pb_type": "CooooT"} # case is important in this test, do not modify.
            ]
        }
        fdoc8 = {}
        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc8), fdoc8)
        self.assertTrue(wf.match(True, True, False)) # no pb type given, all supported
        self.assertTrue(wf.match(True, True, False, pb_types=[]))
        self.assertTrue(wf.match(True, True, False, pb_types=["CooooT"]))
        self.assertTrue(wf.match(True, True, False, pb_types=["cOOOOt"]))
        self.assertTrue(wf.match(True, True, False, pb_types=["cOtcOt", "cOOOOt"]))
        self.assertTrue(wf.match(True, True, False, pb_types=["cUtcUt", "cOtcOt", "cOOOOt"]))
        self.assertFalse(wf.match(True, True, False, pb_types=["cUtcUt"]))
        self.assertFalse(wf.match(True, True, False, pb_types=["Courou"]))
        with self.assertRaises(TypeError):
            wf.match(True, True, "haha-nelson.com")

        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc9), fdoc9)
        self.assertTrue(wf.match(True, True, True))
        self.assertFalse(wf.match(True, True, False))


        # some|watcher tests
        doc12 = {
            'active_pb_some': True,
            'active_pb_watcher': True
        }
        fdoc12 = {}
        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc12), fdoc12)
        self.assertTrue(wf.match(True, True, True))
        self.assertTrue(wf.match(False, True, True))
        self.assertFalse(wf.match(False, False, True))
        self.assertFalse(wf.match(True, True, False))
        self.assertFalse(wf.match(False, True, False))
        self.assertFalse(wf.match(False, False, False))

        doc13 = {
            'active_pb_some': False,
            'active_pb_watcher': True
        }
        fdoc13 = {}

        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc13), fdoc13)
        self.assertFalse(wf.match(True, True, True))
        self.assertFalse(wf.match(False, True, True))
        self.assertTrue(wf.match(False, False, True))
        self.assertFalse(wf.match(True, True, False))
        self.assertFalse(wf.match(False, True, False))
        self.assertFalse(wf.match(False, False, False))

        doc14 = {
            'active_pb_some': True
        }
        fdoc14 = {}
        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc14), fdoc14)
        self.assertTrue(wf.match(True, True, True))
        self.assertTrue(wf.match(False, True, True))
        self.assertTrue(wf.match(False, False, True))
        self.assertTrue(wf.match(True, True, False))
        self.assertTrue(wf.match(False, True, False))
        self.assertFalse(wf.match(False, False, False))

        doc15 = {
            'active_pb_watcher': True
        }
        fdoc15 = {}
        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc15), fdoc15)
        self.assertTrue(wf.match(True, True, True))
        self.assertTrue(wf.match(False, True, True))
        self.assertTrue(wf.match(False, False, True))
        self.assertFalse(wf.match(False, False, False))

        doc16 = {
            'active_pb_watcher': False
        }
        fdoc16 = {}
        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc16), fdoc16)
        self.assertFalse(wf.match(True, True, True))
        self.assertFalse(wf.match(False, True, True))
        self.assertFalse(wf.match(False, False, True))
        self.assertTrue(wf.match(False, False, False))
        self.assertTrue(wf.match(False, True, False))
        self.assertTrue(wf.match(True, True, False))

        doc17 = {
            'active_pb_some': False
        }
        fdoc17 = {}
        wf = WatcherFilter()
        self.assertDictEqual(wf.filter(doc17), fdoc17)
        self.assertFalse(wf.match(True, True, True))
        self.assertFalse(wf.match(False, True, True))
        self.assertFalse(wf.match(False, False, True))
        self.assertTrue(wf.match(False, False, False))
        self.assertFalse(wf.match(False, True, False))
        self.assertFalse(wf.match(True, True, False))

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
