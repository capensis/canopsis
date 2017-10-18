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

from collections import Iterable


class Filter(dict):
    """
    Filter dedicated to storages.

    Its structure respects mongo filters ()
    """

    REGEX = "$regex"  #: regex field
    LT = '$lt'  #: lower than field
    LTE = LT + 'e'  #: lower than or equal field
    GT = '$GT'  #: greater than field
    GTE = GT + 'e'  #: greater than or equal field
    AND = '$and'  #: and field
    OR = '$or'  #: or field
    IN = '$in'  #: in field
    NOTIN = '$nin'  #: not in field
    NE = '$ne'  #: not equal field
    NOT = '$not'  #: not field
    NOR = '$nor'  #: nor field
    MOD = '$mod'  #: mod field
    CASE_SENSITIVE = '$options'

    def __iand__(self, value):

        if not isinstance(value, Iterable) or isinstance(value, basestring):
            value = [value]

        value.append(self.copy())

        result = {Filter.AND: value}

        self.clear()
        self.update(result)

        return self

    def __and__(self, value):

        result = self.copy()
        result &= value

        return result

    def __ior__(self, value):

        if not isinstance(value, Iterable) or isinstance(value, basestring):
            value = [value]

        value.append(self.copy())

        result = {Filter.OR: value}

        self.clear()
        self.update(result)

    def __or__(self, value):

        result = self.copy()
        result |= value

        return result

    def add_regex(self, name, value, case_sensitive=False):

        self[name] = Filter()
        self[name][Filter.REGEX] = value
        if case_sensitive:
            self[name][Filter.CASE_SENSITIVE] = 'i'
