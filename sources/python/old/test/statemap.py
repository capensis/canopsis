#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
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

from unittest import main, TestCase

from canopsis.old.statemap import Statemap


class KnownValues(TestCase):
    def setUp(self):
        pass

    def test_01_mapped_state(self):
        statemap = Statemap(statemap=[0, 0, 1, 1, 1, 2, 2, 2])

        self.assertEqual(0, statemap.get_mapped_state(0))
        self.assertEqual(0, statemap.get_mapped_state(1))
        self.assertEqual(1, statemap.get_mapped_state(2))
        self.assertEqual(1, statemap.get_mapped_state(3))
        self.assertEqual(1, statemap.get_mapped_state(4))
        self.assertEqual(2, statemap.get_mapped_state(5))
        self.assertEqual(2, statemap.get_mapped_state(6))
        self.assertEqual(2, statemap.get_mapped_state(7))

if __name__ == "__main__":
    main(verbosity=2)
