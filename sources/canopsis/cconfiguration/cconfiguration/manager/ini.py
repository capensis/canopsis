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

from cconfiguration.manager import ConfigurationManager

from ConfigParser import RawConfigParser, DuplicateSectionError


class ConfigurationManager(ConfigurationManager):
    """
    Manage ini configuration.
    """

    """
    Register it automatically among global managers.
    """
    __register__ = True

    def _has_category(
        self, config_resource, category, logger, *args, **kwargs
    ):
        return config_resource.has_section(category)

    def _has_parameter(
        self, config_resource, category, parameter_name, logger,
        *args, **kwargs
    ):
        return config_resource.has_option(category, parameter_name)

    def _get_config_resource(
        self, logger, configuration_file=None, *args, **kwargs
    ):
        result = RawConfigParser()

        if configuration_file is not None:

            files = result.read(configuration_file)

            if not files:
                result = None

        return result

    def _get_parameter(
        self, config_resource, category, parameter_name, *args, **kwargs
    ):
        return config_resource.get(category, parameter_name)

    def _set_category(
        self, config_resource, category, logger, *args, **kwargs
    ):
        try:
            config_resource.add_section(category)
        except (DuplicateSectionError, ValueError):
            pass

    def _set_parameter(
        self, config_resource, category, parameter_name, parameter, logger,
        *args, **kwargs
    ):
        config_resource.set(category, parameter_name, parameter)

    def _write_config_resource(
        self, config_resource, configuration_file, *args, **kwargs
    ):
        config_resource.write(open(configuration_file, 'a'))
