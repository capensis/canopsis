#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.configuration.decorator import conf_paths
from canopsis.configuration.configurable import Configurable


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

if __name__ == '__main__':
    main()
