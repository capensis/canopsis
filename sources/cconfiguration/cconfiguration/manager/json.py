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

from cconfiguration.language import ConfigurationManager

from json import load, dump


class ConfigurationManager(ConfigurationManager):
    """
    Manage json configuration.
    """

    """
    Register it automatically among managers.
    """
    __register__ = True

    def _has_category(config_resource, category, logger):
        return category in config_resource

    def _has_parameter(config_resource, category, parameter_name, logger):
        return parameter_name in config_resource[category]

    def _get_config_resource(configuration_file, logger):
        result = None

        try:
            result = load(open(configuration_file))

        except Exception:
            pass

        return result

    def _get_parameter(config_resource, category, parameter_name):
        return config_resource[category][parameter_name]

    def _set_category(config_resource, category, logger):
        config_resource.setdefault(category, dict())

    def _set_parameter(
        config_resource, category, parameter_name, parameter, logger
    ):
        config_resource[category][parameter_name] = parameter

    def _write_config_resource(config_resource, configuration_file):
        dump(open(config_resource, 'w+'))
