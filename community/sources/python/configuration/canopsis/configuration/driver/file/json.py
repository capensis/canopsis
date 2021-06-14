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

from . import FileConfigurationDriver

try:
    from json import loads, dump
except ImportError:
    from simplejson import loads, dump


class JSONConfigurationDriver(FileConfigurationDriver):
    """
    Manage json configuration.
    """

    """
    Register it automatically among managers.
    """
    __register__ = True

    def _has_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        return category.name in conf_resource

    def _has_parameter(
        self, conf_resource, category, param, logger,
        *args, **kwargs
    ):
        return param.name in conf_resource[category.name]

    def _get_conf_resource(
        self, logger, conf_path=None, *args, **kwargs
    ):
        result = {}

        if conf_path is not None:
            result = None

            path = FileConfigurationDriver.get_path(conf_path)

            try:
                with open(path, 'r') as handle:
                    content = handle.read()
                    result = loads(content)

            except Exception:
                pass

        return result

    def _get_categories(self, conf_resource, logger, *args, **kwargs):
        return conf_resource.keys()

    def _get_parameters(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        return conf_resource[category.name].keys()

    def _get_value(
        self, conf_resource, category, param, *args, **kwargs
    ):
        return conf_resource[category.name][param.name]

    def _set_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        conf_resource.setdefault(category.name, {})

    def _set_parameter(
        self, conf_resource, category, param, logger,
        *args, **kwargs
    ):
        conf_resource[category.name][param.name] = param.value

    def _update_conf_resource(
        self, conf_resource, conf_path, *args, **kwargs
    ):

        path = FileConfigurationDriver.get_path(conf_path)

        try:
            with open(path, 'w') as handle:
                dump(conf_resource, handle)

        except Exception:
            pass
