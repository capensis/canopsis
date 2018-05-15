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
from six import string_types

AUTORIZED_TYPES = (string_types, int, list, dict)


class AssociativeTable(object):
    """
    Generic associative table of key/value pairs, grouped in a single
    collection and indexed with a table name.
    """

    def __init__(self, table_name, content):
        self.table_name = table_name
        self.content = content

    def get(self, key):
        """
        Search stored value.

        :param str key: the key to retreive
        :rtype: an AUTHORIZED_TYPES
        """
        return self.content.get(key, None)

    def get_all(self):
        """
        Return all the content
        """
        return self.content

    def set(self, key, value):
        """
        Update a specific element.

        :param str key: the key to access
        :param object value: the value to update
        """
        if not isinstance(value, AUTORIZED_TYPES):
            raise ValueError('Unauthorized insertion type {}'
                             .format(type(value)))

        self.content[key] = value
