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


class Threshold(object):
    """
    Threshold.
    """

    UNITS = {
        'Yocto': 10 ** -24,
        'Zepto': 10 ** -21,
        'Atto': 10 ** -18,
        'Femto': 10 ** -15,
        'Pico': 10 ** -12,
        'Nano': 10 ** -9,
        'Micro': 10 ** -6,
        'Milli': 10 ** -3,
        'Centi': 10 ** -2,
        'Deci': 10 ** -1,
        None: 1,
        'Deca': 10,
        'Hecto': 10 ** 2,
        'Kilo': 10 ** 3,
        'Mega': 10 ** 6,
        'Ton': 10 ** 6,
        'Giga': 10 ** 9,
        'Tera': 10 ** 12,
        'Peta': 10 ** 15,
        'Exa': 10 ** 18,
        'Zelta': 10 ** 21,
        'Yotta': 10 ** 24
    }

    POURCENT = '%'

    def __init__(self, value=10, unit=POURCENT):

        super(Threshold, self).__init__()

        self.value = value
        self.unit = unit

    def _get_value(self, value, operation):

        result = value + self.value

        if self.unit == Threshold.POURCENT:
            delta = self.value * value / 100
        else:
            delta = self.value
            if self.unit in Threshold.UNITS:
                delta *= Threshold.UNITS[self.unit]

        result = getattr(value, operation)(delta)

        return result

    def get_add_value(self, value=10, unit=POURCENT):

        result = self._get_value(value, unit, '__add__')
        return result

    def get_sub_value(self, value=10, unit=POURCENT):

        result = self._get_value(value, unit, '__sub__')
        return result
