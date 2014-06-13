from unittest import TestCase, main

from ctimeserie.aggregation import get_aggregation_value
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
			'DELTA', lambda points: (max(points) - min(points)) / 2)

	def test_sum(self):
		self._test_aggregation('SUM', lambda points: sum(points))

	def test_max(self):
		self._test_aggregation('MAX', lambda points: max(points))

	def test_min(self):
		self._test_aggregation('MIN', lambda points: min(points))


if __name__ == "__main__":
	main()
