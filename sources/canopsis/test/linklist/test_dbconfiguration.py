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

import unittest
from canopsis.common import root_path
from canopsis.linklist.dbconfigurationmanager import DBConfiguration
import xmlrunner


class DBConfigurationTest(unittest.TestCase):
    """
    Base class for all check manager tests.
    """

    def setUp(self):
        self.dbconfigurationmanager = DBConfiguration()


class DBTest(DBConfigurationTest):

    def test_instanciate(self):
        pass


if __name__ == '__main__':
    output = root_path + "tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
