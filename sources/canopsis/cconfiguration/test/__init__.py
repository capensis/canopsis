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


class ConfigurableTest(TestCase):

    def setUp(self):

        self.configuration_files = (
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
        configurable.configuration_files = self.configuration_files

        self.assertEquals(
            configurable.configuration_files,
            self.configuration_files)

        configurable = Configurable(
            configuration_files=self.configuration_files)

        self.assertEquals(
            configurable.configuration_files,
            self.configuration_files)

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

    def test_get_parameters(self):

        configurable = Configurable()

        parameters, errors = configurable.get_parameters()

        self.assertEquals(len(parameters), 0)
        self.assertEquals(len(errors), 0)

        configurable.configuration_files = self.configuration_files

        parameters, errors = configurable.get_parameters()

        self.assertEquals(len(parameters), 0)
        self.assertEquals(len(errors), 0)

        parameters, errors = configurable.get_parameters(
            parsing_rules=dict())

        self.assertEquals(len(parameters), 0)
        self.assertEquals(len(errors), 0)

        configurable = Configurable(
            configuration_files=self.configuration_files,
            parsing_rules=self.parsing_rules)

        # add A section in conf file[0]
        configurable.set_parameters(
            configuration_file=self.configuration_files[0],
            parameter_by_categories=self.parameters0)

        # add A section in conf file[1]
        configurable.set_parameters(
            configuration_file=self.configuration_files[1],
            parameter_by_categories=self.parameters0)
        # add B section in conf file[1]
        configurable.set_parameters(
            configuration_file=self.configuration_files[1],
            parameter_by_categories=self.parameters1)

        parameters, errors = configurable.get_parameters()

        self.assertTrue('a' in parameters)
        self.assertTrue('2' in parameters)
        self.assertTrue('b' in parameters)
        self.assertTrue('None' in errors)

        print parameters, errors

    def test_set_parameters(self):
        raise NotImplementedError()

    def test_reconfigure(self):
        raise NotImplementedError()


if __name__ == '__main__':
    main()
