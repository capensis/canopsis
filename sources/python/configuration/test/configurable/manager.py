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

from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.manager import Manager
from canopsis.configuration.driver.file import FileConfigurationDriver

from tempfile import NamedTemporaryFile

from os import remove


class TestConfigurable(Configurable):
    pass


class TestManager(Manager):

    def _get_conf_paths(self, *args, **kwargs):

        result = super(TestManager, self)._get_conf_paths(*args, **kwargs)
        result.append(NamedTemporaryFile().name)
        return result


class ManagerTest(TestCase):

    def setUp(self):
        pass

    def test_apply_configuration(self):

        manager = TestManager()

        conf_path = FileConfigurationDriver.get_path(manager.conf_paths[-1])

        configurable_name = 'test'
        full_configurable_name = '%s%s' % (
            configurable_name, Manager.CONFIGURABLE_SUFFIX)
        configurable_type_name = '%s%s' % (
            configurable_name, Manager.CONFIGURABLE_TYPE_SUFFIX)
        # ensure configurable doesn't exis
        self.assertFalse(configurable_name in manager)

        self.assertEqual(len(manager.configurables), 0)
        self.assertEqual(len(manager.configurable_types), 0)

        manager[configurable_name] = Configurable()

        self.assertTrue(manager[configurable_name].auto_conf)

        configurable_path = "canopsis.configuration.configurable.Configurable"

        with open(conf_path, 'w+') as conf_file:
            conf_file.write("[MANAGER]")
            # set configurable
            conf_file.write("\n%s=" % (full_configurable_name))
            conf_file.write(configurable_path)
            # set configurable type
            conf_file.write("\n%s=" % (configurable_type_name))
            conf_file.write(configurable_path)
            # set sub-configurable auto_conf to false
            configurable_category = Manager.get_configurable_category(
                configurable_name)
            conf_file.write("\n[%s]" % configurable_category)
            conf_file.write("\nauto_conf=false")

        manager.apply_configuration()

        remove(conf_path)

        self.assertEqual(len(manager.configurables), 1)
        self.assertEqual(len(manager.configurable_types), 1)

        # check if configurable and short attribute exist
        self.assertFalse(manager[configurable_name].auto_conf)
        self.assertTrue(isinstance(manager[configurable_name], Configurable))

        # check if configurable type exist
        self.assertEqual(
            manager.configurable_types[configurable_name], Configurable)

        # change type with str or Configurable class
        manager.configurable_types[configurable_name] = configurable_path
        self.assertEqual(
            manager.configurable_types[configurable_name], Configurable)
        manager.configurable_types[configurable_name] = Configurable
        self.assertEqual(
            manager.configurable_types[configurable_name], Configurable)

        # change type in order to do not remove old value
        manager.test_configurable_type = Configurable
        self.assertTrue(isinstance(manager[configurable_name], Configurable))
        manager.test = configurable_path
        self.assertTrue(isinstance(manager[configurable_name], Configurable))
        manager.test = Configurable
        self.assertTrue(isinstance(manager[configurable_name], Configurable))
        manager.test = Configurable()
        self.assertTrue(isinstance(manager[configurable_name], Configurable))

        self.assertEqual(len(manager.configurables), 1)
        self.assertEqual(len(manager.configurable_types), 1)

        # change type value in order to remove old value
        manager.configurable_types[configurable_name] = TestConfigurable
        self.assertFalse(configurable_name in manager)
        self.assertEqual(
            manager.configurable_types[configurable_name], TestConfigurable)

        self.assertEqual(len(manager.configurables), 0)
        self.assertEqual(len(manager.configurable_types), 1)

        manager[configurable_name] = Configurable()
        self.assertFalse(configurable_name in manager)

        self.assertEqual(len(manager.configurables), 0)
        self.assertEqual(len(manager.configurable_types), 1)

if __name__ == '__main__':
    main()
