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

from cconfiguration import Configurable, Configuration, Category, Parameter

from os import remove


class ConfigurableTest(TestCase):

    def setUp(self):

        self.conf_files = (
            '/tmp/ConfigurableTest0',
            '/tmp/ConfigurableTest1'
        )

        self.configurable = Configurable()

        self.configuration = Configuration(
            Category('A',
                Parameter('a', value='a'),
                Parameter('2', value=2, parser=int),
                Parameter('error', value='error', parser=float)),
            Category('B',
                Parameter('a', value='b'),
                Parameter('b', value='b')))

    def test_configuration_files(self):

        configurable = Configurable()
        configurable.conf_files = self.conf_files

        self.assertEquals(
            configurable.conf_files,
            self.conf_files)

        configurable = Configurable(
            conf_files=self.conf_files)

        self.assertEquals(
            configurable.conf_files,
            self.conf_files)

    def test_auto_conf(self):

        configurable = Configurable()

        self.assertTrue(configurable.auto_conf)

        configurable.auto_conf = False

        self.assertFalse(configurable.auto_conf)

    def test_logging_level(self):

        configurable = Configurable()

        self.assertTrue(configurable.log_lvl, 'INFO')

        configurable = Configurable(log_lvl='DEBUG')

        self.assertTrue(configurable.log_lvl, 'DEBUG')

        configurable.log_lvl = 'INFO'

        self.assertTrue(configurable.log_lvl, 'INFO')

    def test_configuration(self):

        # test to get from no file
        configurable = Configurable()

        configuration = configurable.get_configuration()

        self.assertEquals(len(configuration), len(self.configuration))

        # test to get from files which do not exist
        configurable.conf_files = self.conf_files

        for conf_file in self.conf_files:
            try:
                remove(conf_file)
            except OSError:
                pass

        configuration = configurable.get_configuration()

        self.assertEquals(len(configuration), len(self.configuration))

        # get parameters from empty files
        for conf_file in self.conf_files:
            open(conf_file, 'w').close()

        configuration = configurable.get_configuration()

        self.assertEquals(len(configuration), len(self.configuration))

        # get parameters from empty files and empty parsing_rules
        configuration = Configuration()
        configurable.get_configuration(configuration=configuration)

        self.assertEquals(len(configuration), 0)

        # fill files
        configurable = Configurable(
            conf_files=self.conf_files,
            configuration=self.configuration)

        # add first category in conf file[0]
        configurable.set_configuration(
            conf_file=self.conf_files[0],
            configuration=Configuration(self.configuration['A']),
            conf_manager=tuple(configurable.managers)[0])

        # add second category in conf file[1]
        configurable.set_configuration(
            conf_file=self.conf_files[1],
            configuration=Configuration(self.configuration['B']),
            conf_manager=tuple(configurable.managers)[1])

        configuration = configurable.get_configuration(fill=True)

        parameters, errors = configuration.get_parameters()

        self.assertTrue('a' in parameters and 'a' not in errors)
        self.assertTrue('2' in parameters and '2' not in errors)
        self.assertTrue('b' in parameters and 'b' not in errors)
        self.assertTrue('error' in errors and 'error' not in parameters)

    def test_reconfigure(self):
        raise NotImplementedError()


if __name__ == '__main__':
    main()
