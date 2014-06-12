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

from cconfiguration import Configurable
from cconfiguration.manager import ConfigurationManager

from os import remove


class ConfigurableTest(TestCase):

    def setUp(self):

        self.conf_files = (
            '/tmp/ConfigurableTest0',
            '/tmp/ConfigurableTest1'
        )

        self.configurable = Configurable()
        self.parameters0 = {
            'A': {
                'a': 'a',
                '2': 2,
                'None': None
            }
        }
        self.parameters1 = {
            'B': {
                'a': 'b',
                'b': 'b'
            }
        }

        self.parsing_rules = [
            {
                'A': {
                    'a': str,
                    '2': int,
                    'None': float
                },
                'B': {
                    'a': str,
                    'b': str
                }
            }
        ]

        self.configurable.set_parameters

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

        self.assertTrue(configurable.logging_level, 'INFO')

        configurable = Configurable(logging_level='DEBUG')

        self.assertTrue(configurable.logging_level, 'DEBUG')

        configurable.logging_level = 'INFO'

        self.assertTrue(configurable.logging_level, 'INFO')

    def test_parameters(self):

        # test to get from no file
        configurable = Configurable()

        parameters, errors = configurable.get_parameters()

        self.assertEquals(len(parameters), 0)
        self.assertEquals(len(errors), 0)

        # test to get from files which do not exist
        configurable.conf_files = self.conf_files

        for conf_file in self.conf_files:
            try:
                remove(conf_file)
            except OSError:
                pass

        parameters, errors = configurable.get_parameters()

        self.assertEquals(len(parameters), 0)
        self.assertEquals(len(errors), 0)

        # get parameters from empty files
        for conf_file in self.conf_files:
            open(conf_file, 'w').close()

        parameters, errors = configurable.get_parameters()

        self.assertEquals(len(parameters), 0)
        self.assertEquals(len(errors), 0)

        # get parameters from empty files and empty parsing_rules
        parameters, errors = configurable.get_parameters(
            parsing_rules=dict())

        self.assertEquals(len(parameters), 0)
        self.assertEquals(len(errors), 0)

        # fill files
        configurable = Configurable(
            conf_files=self.conf_files,
            parsing_rules=self.parsing_rules)

        # add A section in conf file[0]
        configurable.set_parameters(
            conf_file=self.conf_files[0],
            parameter_by_categories=self.parameters0,
            conf_manager=tuple(configurable.managers)[0])

        # add A section in conf file[1]
        configurable.set_parameters(
            conf_file=self.conf_files[1],
            parameter_by_categories=self.parameters0,
            conf_manager=tuple(configurable.managers)[1])
        # add B section in conf file[1]
        configurable.set_parameters(
            conf_file=self.conf_files[1],
            parameter_by_categories=self.parameters1,
            conf_manager=tuple(configurable.managers)[1])

        parameters, errors = configurable.get_parameters()

        print parameters, errors, configurable.managers

        self.assertTrue('a' in parameters)
        self.assertTrue('2' in parameters)
        self.assertTrue('b' in parameters)
        self.assertTrue('None' in errors)

    def test_reconfigure(self):
        raise NotImplementedError()


if __name__ == '__main__':
    main()
