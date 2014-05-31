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

from unittest import TestCase, main

from cconfiguration.manager.ini import ConfigurationManager


class ConfigurationManagerTest(TestCase):

    PATH = './ini.conf'

    def setUp(self):
        self.manager = ConfigurationManager()
        self.resource = self.manager._get_config_resource(
            ConfigurationManagerTest.PATH, logger=None)

    def test_get_parameters(self):

        self.manager.get_parameters(
            config_resource=self.resource,
            configuration_file=PATH, )

    def test_set_parameters(self):
        raise NotImplementedError()

if __name__ == '__main__':
    main()
