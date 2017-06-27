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

    F_ENABLED = "enabled"
    F_ENABLED_HIST = "enabled_history"
    F_DISABLED_HIST = "disabled_history"


SCHEMA = {
    "schema": {
        Keys.K_ID.value: "schema_infos",
        Keys.K_TYPE.value: Keys.T_OBJECT.value,
        Keys.K_REQUIRED.value: [Keys.F_ENABLED.value,
                                Keys.F_ENABLED_HIST.value,
                                Keys.F_DISABLED_HIST.value],
        Keys.K_PROPERTIES.value: {
            Keys.F_ENABLED.value: {
                Keys.K_TYPE.value: Keys.T_BOOL.value
            },
            Keys.F_ENABLED_HIST.value: {
                Keys.K_TYPE.value: Keys.T_ARRAY.value,
                Keys.K_ITEMS.value: {
                    Keys.K_TYPE.value: Keys.T_NUMBER.value
                }
            },
            Keys.F_DISABLED_HIST.value: {
                Keys.K_TYPE.value: Keys.T_ARRAY.value,
                Keys.K_ITEMS.value: {
                    Keys.K_TYPE.value: Keys.T_NUMBER.value
                }
            }
        }
    }
}

TEMPLATE_INFOS = {Keys.F_DISABLED_HIST.value: None,
                  Keys.F_ENABLED_HIST.value: None,
                  Keys.F_DISABLED_HIST.value: None}


class MockLogger:

    def __init__(self, test_instance):
        self.test_instance = test_instance
        self.called = False

    def warning(self, _):
        self.called = True


class BaseTest(TestCase):

    def setUp(self):
        self.logger = MockLogger(self)
        self.infosfilter = InfosFilter(self.logger)
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


class TestFilter(BaseTest):

    def test_not_match_schema(self):
        infos = TEMPLATE_INFOS.copy()

        infos[Keys.F_ENABLED] = 1
        infos[Keys.F_DISABLED_HIST] = [1]
        infos[Keys.F_ENABLED_HIST] = [1]
        self.infosfilter.filter(infos)

        self.assertTrue(self.logger.called)
        self.logger.called = True

        infos[Keys.F_ENABLED] = True
        infos[Keys.F_DISABLED_HIST] = "string"
        infos[Keys.F_ENABLED_HIST] = [1]
        self.infosfilter.filter(infos)

        self.assertTrue(self.logger.called)
        self.logger.called = True

        infos[Keys.F_ENABLED] = True
        infos[Keys.F_DISABLED_HIST] = "string"
        infos[Keys.F_ENABLED_HIST] = [1]
        self.infosfilter.filter(infos)

        self.assertTrue(self.logger.called)
        self.logger.called = True

    def test_match_schema(self):
        infos = TEMPLATE_INFOS.copy()

        infos[Keys.F_ENABLED] = True
        infos[Keys.F_DISABLED_HIST] = [1]
        infos[Keys.F_ENABLED_HIST] = [1]
        self.infosfilter.filter(infos)
        self.logger.called = False


if __name__ == '__main__':
    main()
