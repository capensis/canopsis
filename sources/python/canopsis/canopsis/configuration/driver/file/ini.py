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

from os.path import join
from sys import prefix as sys_prefix

from . import FileConfigurationDriver

from ConfigParser import RawConfigParser, DuplicateSectionError,\
    MissingSectionHeaderError


class INIConfigurationDriver(FileConfigurationDriver):
    """
    Manage ini configuration.
    """

    """
    Register it automatically among global managers.
    """
    __register__ = True

    def _has_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        return conf_resource.has_section(category.name)

    def _has_parameter(
        self, conf_resource, category, param, logger,
        *args, **kwargs
    ):
        return conf_resource.has_option(category.name, param.name)

    def _get_conf_resource(
        self, logger, conf_path=None, *args, **kwargs
    ):
        result = RawConfigParser()

        if conf_path is not None:

            files = []

            path = FileConfigurationDriver.get_path(conf_path)

            try:
                files = result.read(path)

            except MissingSectionHeaderError:
                logger.warning('Missing section header')

            if not files:
                result = None

        return result

    def _get_categories(self, conf_resource, logger, *args, **kwargs):
        return conf_resource.sections()

    def _get_parameters(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        return conf_resource.options(category.name)

    def _get_value(
        self, conf_resource, category, param, *args, **kwargs
    ):
        return conf_resource.get(category.name, param.name)

    def _set_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        try:
            conf_resource.add_section(category.name)
        except (DuplicateSectionError, ValueError):
            pass

    def _set_parameter(
        self, conf_resource, category, param, logger,
        *args, **kwargs
    ):
        conf_resource.set(category.name, param.name, param.value)

    def _update_conf_resource(
        self, conf_resource, conf_path, *args, **kwargs
    ):
        with open(join(sys_prefix, 'etc', conf_path), 'w') as f:
            conf_resource.write(f)
