#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase
from canopsis.context_graph.manager import InfosFilter
from enum import Enum

class Keys(Enum):

    K_ID = "_id"
    K_TYPE = "type"
    K_REQUIRED = "required"
    K_PROPERTIES = "properties"
    K_ITEMS = "items"

    T_OBJECT = "object"
    T_BOOL = "boolean"
    T_ARRAY = "array"
    T_NUMBER = "number"
    T_STRING = "string"

SCHEMA = {
    Keys.K_ID.value: "schema_infos",
    Keys.K_TYPE.value: Keys.T_OBJECT.value,
    Keys.K_REQUIRED.value: ["enabled", "enabled_history", "disabled_history"],
    Keys.K_PROPERTIES.value: {
        "enabled": {
            Keys.K_TYPE.value: Keys.T_BOOL.value
        },
        "enabled_history": {
            Keys.K_TYPE.value: Keys.T_ARRAY.value,
            Keys.K_ITEMS.value: {
                Keys.K_TYPE.value: Keys.T_NUMBER.value
            }
        },
        "disabled_history": {
            Keys.K_TYPE.value: Keys.T_ARRAY.value,
            Keys.K_ITEMS.value: {
                Keys.K_TYPE.value: Keys.T_NUMBER.value
            }
        }
    }
}



class MockLogger:

    def __init__(self, test_instance):
        self.test_instance = test_instance

    def warning(self, _):
        self.test_instance.assertTrue(True)

class BaseTest(TestCase):

    def setUp(self):
        self.infosfilter = InfosFilter(MockLogger(self))
        self.infosfilter._schema = SCHEMA

class TestReloadSchema(BaseTest):

    def test_wrong_id(self):
        self.infosfilter._schema_id = "I am not an ID"
        desc = "No infos schema found in database."
        with self.assertRaisesRegexp(ValueError, desc):
            self.infosfilter.reload_schema()

    def test_good_id(self):
        self.infosfilter.reload_schema()
        self.assertTrue(True)


if __name__ == '__main__':
    main()
