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

try:
    from json import loads, dump
except ImportError:
    from simplejson import loads, dump


class ConfigurationManager(ConfigurationManager):
    """
    Manage json configuration.
    """

    """
    Register it automatically among managers.
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
            try:
                with open(configuration_file, 'r') as handle:
                    result = loads(handle.read())

            except Exception:
                pass

        return result

    def _get_parameter(
        self, config_resource, category, parameter_name, *args, **kwargs
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

        try:
            with open(configuration_file, 'a') as handle:
                dump(config_resource, handle)

        except Exception:
            pass
