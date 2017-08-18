#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from unittest import main, TestCase

from canopsis.configuration.model import Parameter, Category
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)


class DecoratorTest(TestCase):
    """
    Configuration Manager unittest class.
    """

    def test_conf_paths(self):

        test_conf_paths = ["test1", "test2"]

        @conf_paths(*test_conf_paths)
        class TestConfigurable(Configurable):
            pass

        testConfigurable = TestConfigurable()

        configurable_conf_paths = testConfigurable.conf_paths

        for i in range(1, len(test_conf_paths)):
            self.assertEqual(test_conf_paths[-i], configurable_conf_paths[-i])

    def test_add_category(self):

        CATEGORY = 'TEST'

        @add_category(name=CATEGORY)
        class TestConfigurable(Configurable):
            pass

        tc = TestConfigurable()

        self.assertTrue(CATEGORY in tc.conf)
        self.assertTrue(len(tc.conf) > 0)
        self.assertTrue(len(tc.conf[CATEGORY]) > 0)

        category_len = len(tc.conf[CATEGORY])

        parameters = [Parameter('a'), Parameter('b')]

        @add_category(name=CATEGORY, content=parameters)
        class TestConfigurable(Configurable):
            pass

        tc = TestConfigurable()

        self.assertTrue(CATEGORY in tc.conf)
        self.assertTrue(len(tc.conf) > 0)
        self.assertEqual(
            len(tc.conf[CATEGORY]), category_len + len(parameters))

        @add_category(name=CATEGORY, unified=False)
        class TestConfigurable(Configurable):
            pass

        tc = TestConfigurable()

        self.assertTrue(CATEGORY in tc.conf)
        self.assertTrue(len(tc.conf) > 0)
        self.assertEqual(len(tc.conf[CATEGORY]), 0)

        @add_category(name=CATEGORY, unified=False, content=parameters)
        class TestConfigurable(Configurable):
            pass

        tc = TestConfigurable()

        self.assertTrue(CATEGORY in tc.conf)
        self.assertTrue(len(tc.conf) > 0)
        self.assertEqual(len(tc.conf[CATEGORY]), len(parameters))

        category = Category(CATEGORY, *parameters)

        @add_category(name=CATEGORY, unified=False, content=category)
        class TestConfigurable(Configurable):
            pass

        tc = TestConfigurable()

        self.assertTrue(CATEGORY in tc.conf)
        self.assertTrue(len(tc.conf) > 0)
        self.assertEqual(len(tc.conf[CATEGORY]), len(category))

if __name__ == '__main__':
    main()
