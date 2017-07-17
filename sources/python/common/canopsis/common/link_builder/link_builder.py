#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

from abc import ABCMeta, abstractmethod
import importlib
import inspect
import re

from canopsis.common.ini_parser import IniParser

CONF_FILE = "etc/common/link_builder.conf"
BUILDERS_CAT = "LINK_BUILDERS"
DEFAULT_BUILDER_CAT = "DEFAULT_BUILDER"
PACKAGE_NAME = "canopsis.common.link_builder.{0}"


class HypertextLinkManager:

    SEPARATOR = ','

    def __init__(self, config, logger):
        self.config = config
        parser = IniParser(CONF_FILE, logger)
        builders_info = parser.get(BUILDERS_CAT)

        for key in builders_info:
            builders_info[key] = re.split(self.SEPARATOR, builders_info[key])

            if builders_info[key][-1] == '':  # in case of a trailing separator
                del builders_info[-1]

        self.builders = []

        # build the list builders
        for fname in builders_info:

            mod_name = PACKAGE_NAME.format(fname)
            mod = None
            try:
                mod = importlib.import_module(mod_name)
            except ImportError:
                logger.warning("Cannot import {0}.".format(mod_name))
                continue

            members = inspect.getmembers(mod, inspect.isclass)

            classes = {name: obj for name, obj in members}

            # if the class_name is a subclass of HypertextLinkBuilder, and add
            # an instance in the builders list
            for class_name in builders_info[fname]:
                if class_name not in classes:
                    logger.warning("Cannot find {0} "
                                   "class in {1}.".format(class_name, fname))
                    continue

                class_obj = classes[class_name]
                if HypertextLinkBuilder in inspect.getmro(class_obj):
                    # TODO add the option here
                    self.builders.append(class_obj(None))
                else:
                    msg = "Class {0} is not a subclass of {1}"
                    logger.warning(msg.format(class_name, HypertextLinkBuilder))

    def links_for_entity(self, entity, options):
        """Generate links for the entity with the builder specify in the
        configuration.

        :param dict entity: the entity to handle
        :param options: the options
        :return list: a list of links as a string.
        """
        result = []
        for builder in self.builders:
            result.append(builder.build(entity, options))

        return result


class HypertextLinkBuilder:

    __metaclass__ = ABCMeta

    @abstractmethod
    def build(self, entity, options):
        """Build links from an entity and the given option.

        :param dict entity: the entity to handle
        :param options: the options
        :return list: a list of links as a string.
        """
        raise NotImplementedError()


class BasicLinkBuilder(HypertextLinkBuilder):

    def __init__(self, options):
        self.options = options

    def build(self, entity, options):
        pass
