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

from unittest import TestCase, main

from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.registry import ConfigurableRegistry
from canopsis.configuration.driver.file import FileConfigurationDriver

from tempfile import NamedTemporaryFile

from os import remove


class TestConfigurable(Configurable):
    pass


class TestRegistry(ConfigurableRegistry):

    def _get_conf_paths(self, *args, **kwargs):

        result = super(TestRegistry, self)._get_conf_paths(*args, **kwargs)
        result.append(NamedTemporaryFile().name)
        return result


class ManagerTest(TestCase):

    def test_apply_configuration(self):

        driver = TestRegistry()

        conf_path = FileConfigurationDriver.get_path(driver.conf_paths[-1])

        configurable_name = 'test'
        full_configurable_name = '{0}{1}'.format(
            configurable_name, ConfigurableRegistry.CONFIGURABLE_SUFFIX)
        configurable_type_name = '{0}{1}'.format(
            configurable_name, ConfigurableRegistry.CONFIGURABLE_TYPE_SUFFIX)
        # ensure configurable doesn't exis
        self.assertFalse(configurable_name in driver)

        self.assertEqual(len(driver.configurables), 0)
        self.assertEqual(len(driver.configurable_types), 0)

        driver[configurable_name] = Configurable()

        self.assertTrue(driver[configurable_name].auto_conf)

        configurable_path = "canopsis.configuration.configurable.Configurable"

        with open(conf_path, 'w+') as conf_file:
            conf_file.write("[MANAGER]")
            # set configurable
            conf_file.write("\n{0}=".format(full_configurable_name))
            conf_file.write(configurable_path)
            # set configurable type
            conf_file.write("\n{0}=".format(configurable_type_name))
            conf_file.write(configurable_path)
            # set sub-configurable auto_conf to false
            configurable_category = \
                ConfigurableRegistry.get_configurable_category(
                    configurable_name)
            conf_file.write("\n[{0}]".format(configurable_category))
            conf_file.write("\nauto_conf=false")

        driver.apply_configuration()

        remove(conf_path)

        self.assertEqual(len(driver.configurables), 1)
        self.assertEqual(len(driver.configurable_types), 1)

        # check if configurable and short attribute exist
        self.assertFalse(driver[configurable_name].auto_conf)
        self.assertTrue(isinstance(driver[configurable_name], Configurable))

        # check if configurable type exist
        self.assertEqual(
            driver.configurable_types[configurable_name], Configurable)

        # change type with str or Configurable class
        driver.configurable_types[configurable_name] = configurable_path
        self.assertEqual(
            driver.configurable_types[configurable_name], Configurable)
        driver.configurable_types[configurable_name] = Configurable
        self.assertEqual(
            driver.configurable_types[configurable_name], Configurable)

        # change type in order to do not remove old value
        driver.test_configurable_type = Configurable
        self.assertTrue(isinstance(driver[configurable_name], Configurable))
        driver.test = configurable_path
        self.assertTrue(isinstance(driver[configurable_name], Configurable))
        driver.test = Configurable
        self.assertTrue(isinstance(driver[configurable_name], Configurable))
        driver.test = Configurable()
        self.assertTrue(isinstance(driver[configurable_name], Configurable))

        self.assertEqual(len(driver.configurables), 1)
        self.assertEqual(len(driver.configurable_types), 1)

        # change type value in order to remove old value
        driver.configurable_types[configurable_name] = TestConfigurable
        self.assertFalse(configurable_name in driver)
        self.assertEqual(
            driver.configurable_types[configurable_name], TestConfigurable)

        self.assertEqual(len(driver.configurables), 0)
        self.assertEqual(len(driver.configurable_types), 1)

        driver[configurable_name] = Configurable()
        self.assertFalse(configurable_name in driver)

        self.assertEqual(len(driver.configurables), 0)
        self.assertEqual(len(driver.configurable_types), 1)

if __name__ == '__main__':
    main()
