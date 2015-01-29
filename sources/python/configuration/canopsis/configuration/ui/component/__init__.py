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
Module dedicated to define visual components
"""

from canopsis.configuration.configurable import Configurable


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
