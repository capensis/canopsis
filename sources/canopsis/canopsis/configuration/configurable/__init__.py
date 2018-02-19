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

from logging import Formatter, getLogger, FileHandler, Filter

from os.path import join, sep

from inspect import isclass

from canopsis.common import root_path
from canopsis.common.init import basestring
from canopsis.configuration.model import Configuration, Category, Parameter
from canopsis.configuration.driver import ConfigurationDriver


class MetaConfigurable(type):
    """Meta class for Configurable."""
    pass

class ConfigurableError(Exception):
    """Handle Configurable errors."""
    pass

class Configurable(object):
    """Manages class conf synchronisation with conf resources."""
    pass