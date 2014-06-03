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

from logging import getLogger

from unittest import main, TestCase

from cconfiguration.manager import ConfigurationManager

from pickle import loads, dump


class ConfigurationManager(ConfigurationManager):
    """
    Configuration Manager for test.
    """

    __register__ = True

    def _has_category(
        self, config_resource, category, logger, *args, **kwargs
    ):

        return category in config_resource

    def _has_parameter(
        self, config_resource, category, parameter_name, logger,
        *args, **kwargs
    ):

        return parameter_name in config_resource[category]

    def _get_config_resource(
        self, logger, configuration_file=None, *args, **kwargs
    ):

        result = dict()

        if configuration_file is not None:

            with open(configuration_file, 'a+') as handle:
                try:
                    result = loads(handle.read())
                except Exception:
                    pass

        return result

    def _get_parameter(
        self, config_resource, category, parameter_name, logger,
        *args, **kwargs
    ):

        return config_resource[category][parameter_name]

    def _set_category(
        self, config_resource, category, logger, *args, **kwargs
    ):

        config_resource.setdefault(category, dict())

    def _set_parameter(
        self, config_resource, category, parameter_name, parameter, logger,
        *args, **kwargs
    ):

        config_resource[category][parameter_name] = parameter

    def _write_config_resource(
        self, config_resource, configuration_file, *args, **kwargs
    ):

        with open(configuration_file, 'a+') as handle:
            try:
                dump(config_resource, handle)
            except Exception:
                pass

from os import remove


class ConfigurationManagerTest(TestCase):
    """
    Configuration Manager unittest class.
    """

    ERROR_PARAMETER = 'foo4'

    def setUp(self):
        self.logger = getLogger()

        self.manager = self._get_configuration_manager()

        # configuration files content
        self.full_parameters = {
            'FOO_0': {
                'foo_0.0': 0,
                'foo_0.1': True,
                'foo_0.2': 'foo'
            },
            'FOO_1': {
                'foo_0.0': 1,
                'foo_0.1': False,
                'foo_0.2': 'foo'
            }
        }

        # parameters to read/write
        self.parameters = {
            'foo': 1,
            'foo2': True,
            'foo3': 'foo'
        }

        # parsing rules
        self.parsing_rules = [
            {
                'FOO': {
                    name: type(value)
                    for name, value in self.parameters.iteritems()}
            }, {
                'FOOFOO': {
                    name: type(value)
                    for name, value in self.parameters.iteritems()}
            }
        ]
        # add parameters in parsing_rules to full_parameters
        for parsing_rule in self.parsing_rules:
            for category in parsing_rule:
                self.full_parameters[category] = self.parameters.copy()

        # introduce an error in parsing_rules
        self.parsing_rules[0]['FOO'][ConfigurationManagerTest.ERROR_PARAMETER]\
            = int
        self.full_parameters['FOO'][ConfigurationManagerTest.ERROR_PARAMETER]\
            = 'er'

        self.configuration_file = self.get_configuration_file()

        # empty configuration file
        try:
            open(self.configuration_file, 'w').close()
        except OSError as ose:  # do nothing if file does not exist
            print(ose)

        # fill configuration file with set_parameters
        self.manager.set_parameters(
            configuration_file=self.configuration_file,
            parameter_by_categories=self.full_parameters,
            logger=self.logger)

    def tearDown(self):
        # remove self configuration file
        try:
            remove(self.configuration_file)
        except OSError:
            pass

    def get_configuration_file(self):

        return '/tmp/cconfiguration.conf'

    def test_get_parameters(self):

        parameters, error_parameters = self.manager.get_parameters(
            configuration_file=self.configuration_file,
            parsing_rules=self.parsing_rules,
            logger=self.logger)

        self.assertEqual(parameters, self.parameters)

        self.assertEqual(len(error_parameters), 1)
        self.assertTrue(
            ConfigurationManagerTest.ERROR_PARAMETER in error_parameters)

    def test_set_parameters(self):

        parameters_by_categories = {
            '_': self.parameters
        }

        self.manager.set_parameters(
            configuration_file=self.configuration_file,
            parameter_by_categories=parameters_by_categories,
            logger=self.logger)

        self.parsing_rules.append({
            '_': {
                name: type(value) for name, value in
                    self.parameters.iteritems()}
        })

    def _get_configuration_manager(self):
        """
        Only one method to override by sub tests
        """
        return ConfigurationManager()

if __name__ == '__main__':
    main()
