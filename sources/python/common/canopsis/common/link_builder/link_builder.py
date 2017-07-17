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

from canopsis.common.associative_table.manager import AssociativeTableManager

PACKAGE_NAME = "canopsis.common.link_builder.{0}"


class HypertextLinkManager:
    """
    Try to load an instanciate HypertextLinkBuilder classes, according to a
    configuration.

    The configuration associate a filename in link_builder folder, to an
    AssociativeTable table name.
    Thus, classes in the same file will use the same configuration dict.
    """

    def __init__(self, config, logger):
        """
        :param <AssociativeTable> config: association of a target filename with
        another AssociativeTable table_name
        :param Logger logger: where to log things
        """
        self.config = config
        self.logger = logger
        self.at_manager = AssociativeTableManager()

        self.builders = []

        # build all available builders in the config collection
        builders_info = config.get_all()
        for fname, options in builders_info.items():

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
            for classe_name, classe in members:
                if classe is not HypertextLinkBuilder and \
                   HypertextLinkBuilder in inspect.getmro(classe):

                    founded.append(classe_name)
                    assoc_table = self.at_manager.get(options)
                    self.builders.append(classe(assoc_table.get_all()))

            if len(founded) == 0:
                msg = "Any classes of {} is a subclass of {}. Ignoring it..."
                logger.warning(msg.format(mod_name, HypertextLinkBuilder))

    def links_for_entity(self, entity, options=None):
        """Generate links for the entity with the builder specify in the
        configuration.

        :param dict entity: the entity to handle
        :param dict options: the options
        :return list: a list of links as a string.
        """
        result = []
        for builder in self.builders:
            result.append(builder.build(entity, options))

        return result


class HypertextLinkBuilder:

    """
    Abstract class for all LinkBuilder classes.
    """

    __metaclass__ = ABCMeta

    @abstractmethod
    def build(self, entity, options=None):
        """Build links from an entity and the given option.

        :param dict entity: the entity to handle
        :param dict options: the options table
        :return list: a list of links as strings.
        """
        raise NotImplementedError()
