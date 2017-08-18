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

from canopsis.timeserie.aggregation import get_aggregation_value
from random import random


class AggregationTest(TestCase):

	def setUp(self):

		self.points = [v * random() for v in range(int(random() * 1000))]

	def _test_aggregation(self, aggregation, fn):
		aggregation_value = get_aggregation_value(aggregation, self.points)
		fn_result = fn(self.points)

		self.assertEqual(fn_result, aggregation_value)

	def test_mean(self):
		self._test_aggregation('MEAN', lambda points: sum(points) / len(points))

	def test_last(self):
		self._test_aggregation('LAST', lambda points: points[-1])

	def test_first(self):
		self._test_aggregation('FIRST', lambda points: points[0])

	def test_delta(self):
		self._test_aggregation(
			'DELTA', lambda points: (max(points) - min(points)) / 2
		)

	def test_sum(self):
		self._test_aggregation('SUM', lambda points: sum(points))

	def test_max(self):
		self._test_aggregation('MAX', lambda points: max(points))

	def test_min(self):
		self._test_aggregation('MIN', lambda points: min(points))


if __name__ == "__main__":
	main()
