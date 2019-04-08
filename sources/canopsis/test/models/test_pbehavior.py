#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import time
import unittest

from canopsis.common import root_path
from canopsis.models.pbehavior import PBehavior

import xmlrunner

db_pbehavior = {
    "_id": "xxx",
    "filter": "{\"_id\": \"Tanya/Adams\"}",
    "name": "Albert",
    "author": "Einstein",
    "enabled": True,
    "type_": "pause",
    "comments": [],
    "connector": "canopsis",
    "reason": "no",
    "connector_name": "canopsis",
    "rrule": "",
    "tstart": 0,
    "tstop": 2147483647,
    "eids": ["Tanya/Adams"],
    "source": "nagioslike",
    "exdate": [],
    "timezone": "Europe/Paris"
}


class PBehaviorTest(unittest.TestCase):
    """
    Test the pbehavior model.
    """

    def test_pbehavior(self):
        pbehavior = PBehavior(**PBehavior.convert_keys(db_pbehavior))
        self.assertEqual(pbehavior._id, db_pbehavior['_id'])
        self.assertTrue(pbehavior.enabled)

        self.assertDictEqual(pbehavior.to_dict(), db_pbehavior)

    def test_pbehavior_is_active(self):
        pb = PBehavior(**PBehavior.convert_keys(db_pbehavior))
        pb.tstop = None
        self.assertTrue(pb.is_active)

        pb.tstop = pb.tstart
        self.assertFalse(pb.is_active)

        pb.tstop = int(time.time()) + 1000
        self.assertTrue(pb.is_active)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
