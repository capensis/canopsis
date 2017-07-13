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

import ConfigParser


class IniParser(object):
    """
    Parse a configuration file.
    """

    def __init__(self, path, logger):
        self.logger = logger
        self.path = path

        self.config = ConfigParser.ConfigParser()
        self.config.read(self.path)

    def get_sections(self):
        """
        Get all section names.

        :rtype: list
        """
        return self.config.sections()

    def get(self, section):
        """
        Get a specific section, as a dict.

        :param str section: a section name
        :rtype: dict
        """
        result = {}
        # TODO: if the section doesnt exist ?!
        for key, item in self.config.items(section):
            if key in result:
                self.logger.warning('Duplicated key {} in section {}'
                                    .format(key, section))
                continue

            result[key] = item

        return result
