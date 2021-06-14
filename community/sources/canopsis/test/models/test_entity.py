#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import unittest

from canopsis.common import root_path
from canopsis.models.entity import Entity

import xmlrunner

db_entity = {
    "_id": "Tanya/Adams",
    "impact": ["Adams"],
    "name": "Tanya",
    "enable_history": [],
    "measurements": {},
    "enabled": True,
    "depends": ["red/alert"],
    "infos": {},
    "type": "resource",
    "last_state_change": 0
}


class EntityTest(unittest.TestCase):
    """
    Test the entity model.
    """

    def test_entity(self):
        entity = Entity(**Entity.convert_keys(db_entity))
        self.assertEqual(entity._id, db_entity['_id'])
        self.assertTrue(entity.enabled)
        self.assertTrue(len(entity.enable_history) > 0)

        entity = entity.to_dict()
        del entity['enable_history']
        del db_entity['enable_history']
        self.assertDictEqual(entity, db_entity)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
