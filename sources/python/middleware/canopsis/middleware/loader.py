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

from importlib import import_module
from canopsis.configuration import add_category, conf_paths

from canopsis.configuration import Configurable, Parameter

LOADER_CONF_PATH = 'middleware/loader.conf'
LOADER_CATEGORY = 'LOADER'


@add_category(LOADER_CATEGORY)
@conf_paths(LOADER_CONF_PATH)
class Loader(Configurable):

    LIBRARIES = 'libraries'

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
            for library in value.split(','):
                import_module(library)

    def _configure(self, unified_conf, *args, **kwargs):

        super(Loader, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        self._update_property(
            unified_conf=unified_conf, param_name=Loader.LIBRARIES,
            public=True)

    def _conf(self, *args, **kwargs):

        result = super(Loader, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Loader.CATEGORY,
            new_content=(
                Parameter(Loader.LIBRARIES, self.libraries)))

        return result
