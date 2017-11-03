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

from logging import getLogger

from unittest import main, TestCase

from canopsis.configuration.model import Configuration, Category, Parameter
from canopsis.configuration.driver.file import FileConfigurationDriver

from pickle import loads, dump

from os import remove


class TestConfigurationDriver(FileConfigurationDriver):
    """
    Configuration Manager for test.
    """

    __register__ = True

    def _has_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):

        return category in conf_resource

    def _has_parameter(
        self, conf_resource, category, param, logger,
        *args, **kwargs
    ):

        return param.name in conf_resource[category.name]

    def _get_categories(self, conf_resource, logger, *args, **kwargs):
        return conf_resource.keys()

    def _get_parameters(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        return conf_resource[category.name].keys()

    def _get_conf_resource(
        self, logger, conf_path=None, *args, **kwargs
    ):

        result = dict()

        if conf_path is not None:

            with open(conf_path, 'r') as handle:

                try:
                    result = loads(handle.read())

                except Exception:
                    pass

        return result

    def _get_value(
        self, conf_resource, category, param, logger,
        *args, **kwargs
    ):

        return conf_resource[category.name][param.name]

    def _set_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):

        conf_resource.setdefault(category.name, dict())

    def _set_parameter(
        self, conf_resource, category, param, logger,
        *args, **kwargs
    ):
        conf_resource[category.name][param.name] = param.value

    def _update_conf_resource(
        self, conf_resource, conf_path, *args, **kwargs
    ):

        with open(conf_path, 'w') as handle:

            try:
                dump(conf_resource, handle)

            except Exception:
                pass


class ConfigurationDriverTest(TestCase):
    """
    Configuration Manager unittest class.
    """

    ERROR_PARAMETER = 'foo4'

    def setUp(self):
        self.logger = getLogger()

        self.manager = self._get_configuration_manager()

        self.conf = Configuration(
            Category('A',
                Parameter('a', value=0, parser=int),  # a is 0
                Parameter('b', value=True, parser=Parameter.bool)),
                # b is overriden
            Category('B',
                Parameter('b', value=1, parser=int),  # b is 1
                Parameter('c', value='er', parser=int)))  # error

        self.conf_path = self.get_configuration_file()

    def get_configuration_file(self):

        return '/tmp/canopsis.configuration.conf'

    def _remove(self):
        try:
            remove(self.conf_path)
        except OSError:
            pass

    def _open(self):

        try:
            open(self.conf_path, 'w').close()
        except OSError:
            pass

    def test_configuration(self):

        # try to get conf from not existing file
        self._remove()

        conf = self.manager.get_configuration(
            conf_path=self.conf_path,
            logger=self.logger)

        self.assertEqual(conf, None)

        # get conf from an empty media
        self._open()

        conf = self.manager.get_configuration(
            conf_path=self.conf_path,
            logger=self.logger)

        self.assertTrue(conf is None)

        # get full conf
        self.manager.set_configuration(
            conf_path=self.conf_path,
            conf=self.conf,
            logger=self.logger)

        conf = self.manager.get_configuration(
            conf_path=self.conf_path,
            conf=self.conf,
            logger=self.logger)

        self.assertFalse(conf is None)
        self.assertEqual(len(conf), 2)

        unified_conf = conf.unify()

        parameters = unified_conf[Configuration.VALUES]
        errors = unified_conf[Configuration.ERRORS]

        self.assertTrue('a' in parameters and 'a' not in errors)
        self.assertEqual(parameters['a'].value, 0)
        self.assertTrue('b' in parameters and 'b' not in errors)
        self.assertEqual(parameters['b'].value, 1)
        self.assertTrue('c' in errors and 'c' not in parameters)

        # get some conf
        conf = Configuration(self.conf['B'])

        conf = self.manager.get_configuration(
            conf_path=self.conf_path,
            conf=conf,
            logger=self.logger)

        unified_conf = conf.unify()

        parameters = unified_conf[Configuration.VALUES]
        errors = unified_conf[Configuration.ERRORS]

        self.assertTrue('a' not in parameters and 'a' not in errors)
        self.assertTrue('b' in parameters and 'b' not in errors)
        self.assertEqual(parameters['b'].value, 1)
        self.assertTrue('c' in errors and 'c' not in parameters)

    def _get_configuration_manager(self):
        """
        Only one method to override by sub tests
        """
        return TestConfigurationDriver()

    def _get_manager_path(self):

        return 'test.manager.file.TestConfigurationDriver'

    def _get_manager(self):

        return TestConfigurationDriver

if __name__ == '__main__':
    main()
