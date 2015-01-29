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

from stat import ST_SIZE

from os import stat
from os.path import exists, join
import sys

from canopsis.configuration.driver import ConfigurationDriver


class FileConfigurationDriver(ConfigurationDriver):
    """
    Configuration Manager dedicated to files.
    """

    CONF_DIR = join(sys.prefix, 'etc')

    def exists(self, conf_path, *args, **kwargs):

        path = FileConfigurationDriver.get_path(conf_path)

        result = exists(path) and stat(path)[ST_SIZE]

        return result

    @staticmethod
    def get_path(conf_path):

        result = join(FileConfigurationDriver.CONF_DIR, conf_path)

        return result
