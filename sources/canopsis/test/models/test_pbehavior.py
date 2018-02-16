#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import unittest

from canopsis.models.pbehavior import PBehavior

import xmlrunner

db_pbehavior = {}


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
