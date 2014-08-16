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

from canopsis.configuration import Configurable


def conf_paths(*conf_paths):
    """
    Configurable decorator which adds conf_path paths to a Configurable.

    :param paths: conf resource pathes to add to a Configurable
    :type paths: list of str

    Example:
    >>>@conf_paths('myexample/example.conf')
    >>>class MyConfigurable(Configurable):
    >>>    pass
    >>>assert MyConfigurable().conf_paths == \
        (Configurable().conf_paths + ['myexample/example.conf'])
    """

    def _get_conf_paths(self):
        # get super result and append conf_paths
        result = super(type(self), self)._get_conf_paths()
        result += conf_paths

        return result

    def add_conf_paths(cls):
        # add _get_conf_paths method to configurable classes
        if issubclass(cls, Configurable):
            cls._get_conf_paths = _get_conf_paths

        else:
            raise Configurable.Error(
                "class %s is not a Configurable class" % cls)

        return cls

    return add_conf_paths
