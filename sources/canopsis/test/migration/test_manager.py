#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from canopsis.migration.manager import MigrationTool, MigrationModule


class MigrationToolTest(TestCase):

    def setUp(self):
        self.migration_tool = MigrationTool()

    def test_fill(self):
        #self.migration_tool.fill()
        pass


class MigrationModuleTest(TestCase):

    def setUp(self):
        self.migration_module = MigrationModule()

    #def test_get_version(self):
    #    res = self.migration_module.get_version('perfdata')
    #    self.assertEqual(res, 0)

if __name__ == '__main__':
    main()
