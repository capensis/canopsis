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

from importlib import import_module

from canopsis.configuration.model import Parameter
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths)

LOADER_CONF_PATH = 'middleware/loader.conf'  #: loader library conf path
LOADER_CATEGORY = 'LOADER'  #: loader category
LOADER_LIBRARIES = 'libraries'  #: libraries parameter name


@add_category(LOADER_CATEGORY, content=Parameter(LOADER_LIBRARIES))
@conf_paths(LOADER_CONF_PATH)
class Loader(Configurable):
    """
    Middleware library loader
    """

    def __init__(self, libraries=None, *args, **kwargs):

        super(Loader, self).__init__(*args, **kwargs)

        self._libraries = libraries

    @property
    def libraries(self):
        return self._libraries

    @libraries.setter
    def libraries(self, value):
        self._libraries = value
        if value is not None:
            if isinstance(value, basestring):
                value = value.split(',')
            for library in value:
                import_module(library)
