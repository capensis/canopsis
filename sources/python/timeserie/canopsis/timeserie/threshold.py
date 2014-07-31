#!/usr/bin/env python
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


class Threshold(object):
    """
    Threshold.
    """

    POURCENT = '%'
    KILO = 'Kilo'
    MEGA = 'Mega'
    GIGA = 'Giga'
    TERA = 'Tera'

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
            if self.unit == Threshold.KILO:
                delta *= 1000
            elif self.unit == Threshold.MEGA:
                delta *= 1000 * 1000
            elif self.unit == Threshold.GIGA:
                delta *= 1000 * 1000 * 1000
            elif self.unit == Threshold.TERA:
                delta *= 1000 * 1000 * 1000 * 1000

        result = getattr(value, operation)(delta)

        return result

    def get_add_value(self, value=10, unit=POURCENT):

        result = self._get_value(value, unit, '__add__')
        return result

    def get_sub_value(self, value=10, unit=POURCENT):

        result = self._get_value(value, unit, '__sub__')
        return result
