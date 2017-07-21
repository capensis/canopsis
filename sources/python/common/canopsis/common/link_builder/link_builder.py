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

PACKAGE_NAME = "canopsis.common.link_builder.{0}"


class HypertextLinkManager:
    """
    Try to load an instanciate HypertextLinkBuilder classes, according to
    a configuration.

    The configuration associate a filename in link_builder folder, to a dict.
    Thus, classes in the same file will use the same configuration dict and
    we cannot instanciate the same class with other parameters.
    """

    def __init__(self, config, logger):
        """
        :param dict config: association of a target filename with a
        configuration dict
        :param Logger logger: where to log things
        """
        self.config = config
        self.logger = logger

        self.builders = []

        # instanciate all available builders in the config
        for fname, options in self.config.items():

            mod_name = PACKAGE_NAME.format(fname)
            mod = None
            try:
                mod = importlib.import_module(mod_name)
            except ImportError:
                logger.warning("Cannot import {}.".format(mod_name))
                continue

            members = inspect.getmembers(mod, inspect.isclass)

            founded = []
            # we search for classes that inherit from HypertextLinkBuilder
            for member, member_type in members:
                if member_type is not HypertextLinkBuilder and \
                   HypertextLinkBuilder in inspect.getmro(member_type):

                    founded.append(member)
                    self.builders.append(member_type(options))

            if len(founded) == 0:
                msg = "Any classes of {} is a subclass of {}. Ignoring it..."
                logger.warning(msg.format(mod_name, HypertextLinkBuilder))

    def links_for_entity(self, entity, options={}):
        """
        Generate links for the entity with the builder specify in the
        configuration

        :param dict entity: the entity to handle
        :param dict options: the options
        :returns: an association of categories with a list of strings
        :rtype: dict of list
        """
        result = {}
        for builder in self.builders:
            for cat, build in builder.build(entity, options).items():
                if cat not in result:
                    result[cat] = []
                result[cat] = result[cat] + build
        return result


class HypertextLinkBuilder:

    """
    Abstract class for all LinkBuilder classes.
    """

    __metaclass__ = ABCMeta

    CATEGORY_KEY = "category"
    DEFAULT_CATEGORY = "links"

    def __init__(self, options={}):
        self.options = options

        # The category is set on instanciation (only)
        if self.CATEGORY_KEY in options:
            self.category = options[self.CATEGORY_KEY]
        else:
            self.category = self.DEFAULT_CATEGORY

    @abstractmethod
    def build(self, entity, options={}):
        """
        Build links from an entity and the given option

        :param dict entity: the entity to handle
        :param dict options: the options table
        :returns: an association of categories with a list of strings
        :rtype: dict of list
        """
        raise NotImplementedError()
