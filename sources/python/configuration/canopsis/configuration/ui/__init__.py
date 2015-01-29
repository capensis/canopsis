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

"""
Module dedicated to manage configuration edition with a dedicated UI.
"""

from canopsis.common.utils import lookup
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.parameters import (
    Configuration, Category, Parameter)


class ConfigurationEditor(Configurable):
    """
    Configuration editor.
    """

    FULLSCREEN = 'fullscreen'
    CONF_FILE = 'conf_file'

    COMPONENT_PREFIX = 'component_'

    DEFAULT_CONFIGURATION = Configuration(
        Category('UI',
            Parameter(FULLSCREEN, parser=bool),
            Parameter(CONF_FILE)))

    def __init__(
        self, fullscreen=False, conf_file=None, components_by_parsers=None,
        *args, **kwargs
    ):
        super(ConfigurationEditor, self).__init__(*args, **kwargs)

        self.fullscreen = fullscreen
        self.conf_file = conf_file
        self.components_by_parsers = components_by_parsers

    def _configure(self, parameters, error_parameters, *args, **kwargs):

        self.fullscreen = parameters.get(ConfigurationEditor.FULLSCREEN)

        for name in parameters:
            parameter = parameters[name]
            if name.startswith(Configurable.COMPONENT_PREFIX):
                parser = name[:len(Configurable.COMPONENT_PREFIX)]
                component = lookup(parameter)
                self.components_by_parsers[parser] = component

    def display(self, *args, **kwargs):
        """
        Display content of the editor.
        """

        raise NotImplementedError()


class Component(Configurable):
    """
    Visual component
    """

    pass


class TextComponent(Component):
    """
    Textual parameter component
    """

    pass


class MultiTextComponent(Component):
    """
    Multiple choice parameter component
    """

    pass


class CheckboxComponent(Configurable):
    """
    Binary parameter component
    """

    pass


class ColorComponent(Configurable):
    """
    Colorized parameter component
    """

    pass
