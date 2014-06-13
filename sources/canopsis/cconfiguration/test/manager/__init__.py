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

from cconfiguration import Configuration, Category, Parameter
from cconfiguration.manager import ConfigurationManager

from pickle import loads, dump

from os import remove


class ConfigurationManager(ConfigurationManager):
    """
    Configuration Manager for test.
    """

    __register__ = True

    def _has_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):

        return category in conf_resource

    def _has_parameter(
        self, conf_resource, category, parameter, logger,
        *args, **kwargs
    ):

        return parameter.name in conf_resource[category.name]

    def _get_categories(self, conf_resource, logger, *args, **kwargs):
        return conf_resource.keys()

    def _get_parameters(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        return conf_resource[category.name].keys()

    def _get_conf_resource(
        self, logger, conf_file=None, *args, **kwargs
    ):

        result = dict()

        if conf_file is not None:

            with open(conf_file, 'r') as handle:

                try:
                    result = loads(handle.read())

                except Exception:
                    pass

        return result

    def _get_value(
        self, conf_resource, category, parameter, logger,
        *args, **kwargs
    ):

        return conf_resource[category.name][parameter.name]

    def _set_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):

        conf_resource.setdefault(category.name, dict())

    def _set_parameter(
        self, conf_resource, category, parameter, logger,
        *args, **kwargs
    ):
        conf_resource[category.name][parameter.name] = parameter.value

    def _update_conf_file(
        self, conf_resource, conf_file, *args, **kwargs
    ):

        with open(conf_file, 'w') as handle:

            try:
                dump(conf_resource, handle)

            except Exception:
                pass


class ConfigurationManagerTest(TestCase):
    """
    Configuration Manager unittest class.
    """

    ERROR_PARAMETER = 'foo4'

    def setUp(self):
        self.logger = getLogger()

        self.manager = self._get_configuration_manager()

        self.configuration = Configuration(
            Category('A',
                Parameter('a', value=0, parser=int),  # a is 0
                Parameter('b', value=True, parser=bool)),  # b is overriden
            Category('B',
                Parameter('b', value=1, parser=int),  # b is 1
                Parameter('c', value='er', parser=int)))  # error

        self.conf_file = self.get_configuration_file()

    def get_configuration_file(self):

        return '/tmp/cconfiguration.conf'

    def test_configuration(self):

        # try to get configuration from not existing file
        try:
            remove(self.conf_file)
        except OSError:
            pass

        configuration = self.manager.get_configuration(
            conf_file=self.conf_file,
            logger=self.logger)

        self.assertEqual(configuration, None)

        # get configuration from an empty file
        try:
            open(self.conf_file, 'w').close()
        except OSError:
            pass

        configuration = self.manager.get_configuration(
            conf_file=self.conf_file,
            logger=self.logger)

        self.assertTrue(configuration is None)

        # get full configuration
        self.manager.set_configuration(
            conf_file=self.conf_file,
            configuration=self.configuration,
            logger=self.logger)

        configuration = self.manager.get_configuration(
            conf_file=self.conf_file,
            configuration=self.configuration,
            logger=self.logger,
            fill=True)

        self.assertFalse(configuration is None)
        self.assertEqual(len(configuration), 2)

        parameters, errors = configuration.get_parameters()

        self.assertTrue('a' in parameters and 'a' not in errors)
        self.assertEqual(parameters['a'], 0)
        self.assertTrue('b' in parameters and 'b' not in errors)
        self.assertEqual(parameters['b'], 1)
        self.assertTrue('c' in errors and 'c' not in parameters)

        # get some configuration
        configuration = Configuration(
            self.configuration['B'])

        configuration = self.manager.get_configuration(
            conf_file=self.conf_file,
            configuration=configuration,
            logger=self.logger)

        parameters, errors = configuration.get_parameters()

        self.assertTrue('a' not in parameters and 'a' not in errors)
        self.assertTrue('b' in parameters and 'b' not in errors)
        self.assertEqual(parameters['b'], 1)
        self.assertTrue('c' in errors and 'c' not in parameters)

    def _get_configuration_manager(self):
        """
        Only one method to override by sub tests
        """
        return ConfigurationManager()

    def test_manager(self):

        manager = ConfigurationManager.get_manager(self._get_manager_path())

        self.assertTrue(manager is self._get_manager())

    def _get_manager_path(self):

        return 'test.manager.ConfigurationManager'

    def _get_manager(self):

        return ConfigurationManager

if __name__ == '__main__':
    main()
