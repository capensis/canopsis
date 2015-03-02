#!/usr/bin/env python
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

from unittest import TestCase, main
import math
from canopsis.common.math_parser import Formulas

class FormulasTest(TestCase):
    """docstring for ClassName"""

    def setUp(self):
        _dict = {'x':4, 'y':5}
        self.formula = Formulas(_dict)

    def test(self):
        '''abs'''
        expressions = {"x^2 + 9*x + 5":4**2 + 9*4 + 5, "x^y + x + 2*y":4**5 + 4 + 2*5,"-9": -9, "-E":-math.e, "9 + 3 + 6":18, "2*3.14159": 2*3.14159, "PI * PI / 10": math.pi * math.pi / 10, "PI^2": math.pi**2, "E^PI": math.e**math.pi, "2^3^2": 2**3**2, "sgn(-2)": -1}

        for k, v in expressions.iteritems():
            self.assertEqual(self.formula.evaluate(k), v)


if __name__ == '__main__':
    main()
