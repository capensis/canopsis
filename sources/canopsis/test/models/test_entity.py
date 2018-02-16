#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import unittest

from canopsis.models.entity import Entity

import xmlrunner

db_entity = {}


class EntityTest(unittest.TestCase):
    """
    Test the entity model.
    """

    def test_entity(self):
        entity = Entity(**Entity.convert_keys(db_entity))
        print(entity)
        self.assertEqual(entity._id, db_entity['_id'])
        self.assertTrue(entity.enabled)
        self.assertTrue(len(entity.enable_history) > 0)

        self.assertDictEqual(entity.to_dict(), db_entity)

if __name__ == '__main__':
    output = root_path + "/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
