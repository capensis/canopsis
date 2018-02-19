#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

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
    "eids": ["Tanya/Adams"]
}


class PBehaviorTest(unittest.TestCase):
    """
    Test the pbehavior model.
    """

    def test_pbehavior(self):
        pbehavior = PBehavior(**PBehavior.convert_keys(db_pbehavior))
        print(pbehavior)
        self.assertEqual(pbehavior._id, db_pbehavior['_id'])
        self.assertTrue(pbehavior.enabled)

        self.assertDictEqual(pbehavior.to_dict(), db_pbehavior)

if __name__ == '__main__':
    output = root_path + "/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
