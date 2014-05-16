import unittest

from datetime import datetime
from time import time

import sys
import os.path

sys.path.append(os.path.abspath('.'))

from pyperfstore3.timewindow import Period

import calendar


class PeriodTest(unittest.TestCase):

	def test_sliding_datetime(self):

		# get current datetime
		dt = datetime.now()

		for unit in Period.UNITS:

			period = Period(**{unit: 1})

			round_dt = period.sliding_datetime(dt)
			self.assertEqual(round_dt, dt)

			value = getattr(dt, unit, None)
			if value is not None:
				period.unit_values[unit] = value + 1 if unit is not Period.YEAR else 2000
				round_dt = period.sliding_datetime(dt)
				round_value = getattr(round_dt, unit)

				if round_value is not None:
					if unit is Period.YEAR:
						self.assertEqual(round_value, 2000)
					elif unit is Period.DAY:
						_, monthday = calendar.monthrange(dt.year, dt.month-1)
						self.assertEqual(round_value, monthday)
					elif unit is Period.MONTH:
						self.assertEqual(round_value, 12)
					else:
						self.assertEqual(round_value, 0)

			if Period.MICROSECOND is not unit:
				normalized_dt = period.sliding_datetime(dt, normalize=True)
				for _unit in Period.UNITS[0:Period.UNITS.index(unit)-1]:
					if _unit is not Period.WEEK:
						if _unit is Period.MONTH or _unit is Period.DAY:
							self.assertEqual(getattr(normalized_dt, _unit), 1)
						else:
							self.assertEqual(getattr(normalized_dt, _unit), 0)

	def test_sliding_timestamp(self):

		t = time()

		for unit in Period.UNITS:
			period = Period(**{unit: 1})
			st = period.sliding_timestamp(t)
			self.assertEqual(t, st)

from pyperfstore3.timewindow import Interval


class IntervalTest(unittest.TestCase):

	def test_min_max_empty(self):

		interval = Interval()

		self.assertEqual(None, interval.min())
		self.assertEqual(None, interval.max())

	def test_min_max_point(self):

		interval = Interval(2)

		self.assertEqual(0, interval.min())
		self.assertEqual(2, interval.max())

	def test_min_max_points(self):

		interval = Interval(2, 3)

		self.assertEqual(0, interval.min())
		self.assertEqual(3, interval.max())

	def test_min_max_interval(self):

		interval = Interval((2, 3))

		self.assertEqual(2, interval.min())
		self.assertEqual(3, interval.max())

	def test_min_max_intervals(self):

		interval = Interval((2, 3), (4, 6))

		self.assertEqual(2, interval.min())
		self.assertEqual(6, interval.max())

	def test_empty_sub_interval(self):

		interval = Interval()

		self.assertEqual(len(interval.sub_intervals), 0)

	def test_sub_interval_simple_point(self):

		interval = Interval(1)

		self.assertEqual(len(interval.sub_intervals), 1)

	def test_sub_interval_multi_points(self):

		interval = Interval(2, 3)

		self.assertEqual(len(interval.sub_intervals), 1)

	def test_sub_interval_interval(self):

		interval = Interval((2, 3))

		self.assertEqual(len(interval.sub_intervals), 1)

	def test_sub_interval_multi_interval(self):

		interval = Interval((2, 3), (4, 5))

		self.assertEqual(len(interval.sub_intervals), 2)

	def test_sub_interval_multi_interval_with_intersection(self):

		interval = Interval((2, 5), (4, 6))

		self.assertEqual(len(interval.sub_intervals), 1)

	def test_contains_empty(self):

		interval = Interval()

		self.assertFalse(2 in interval)

		self.assertFalse((0, 2) in interval)

	def test_contains_simple_interval(self):

		interval = Interval(2)

		self.assertTrue(2 in interval)

		self.assertTrue(1.5 in interval)

		self.assertFalse(-1 in interval)

		self.assertTrue((0, 2) in interval)

		self.assertFalse((-1, 2) in interval)

	def test_contains_simple_negative_interval(self):

		interval = Interval(-2)

		self.assertTrue(-2 in interval)

		self.assertTrue(-1.5 in interval)

		self.assertFalse(1 in interval)

		self.assertTrue((0, -2) in interval)

		self.assertFalse((-1, 2) in interval)

	def test_contains_multi_simple_point(self):

		interval = Interval(1, 2)

		self.assertTrue(2 in interval)

		self.assertTrue(1.5 in interval)

		self.assertFalse(-1 in interval)

		self.assertTrue((1, 2) in interval)

		self.assertFalse((-1, 2) in interval)

	def test_contains_interval(self):

		interval = Interval((1, 2))

		self.assertTrue(2 in interval)

		self.assertTrue(1.5 in interval)

		self.assertFalse(-1 in interval)

		self.assertTrue((1, 2) in interval)

		self.assertFalse((0, 2) in interval)

	def test_contains_multi_interval(self):

		interval = Interval((1, 2), (6, 8))

		self.assertTrue(2 in interval)

		self.assertTrue(7 in interval)

		self.assertFalse(3 in interval)

		self.assertTrue((1, 2, 7) in interval)

		self.assertFalse((0, 2, 7) in interval)

	def test_len_empty(self):

		interval = Interval()

		self.assertEqual(len(interval), 0)

	def test_simple_len(self):

		interval = Interval(10)

		self.assertEqual(len(interval), 10)

	def test_simple_negative_len(self):

		interval = Interval(-10)

		self.assertEqual(len(interval), 10)

	def test_multi_simple_len(self):

		interval = Interval(2, 4)

		self.assertEqual(len(interval), 4)

	def test_multi_simple_negative_len(self):

		interval = Interval(-2, -4)

		self.assertEqual(len(interval), 4)

	def test_interval_len(self):

		interval = Interval((2, 4))

		self.assertEqual(len(interval), 2)

	def test_negative_interval_len(self):

		interval = Interval((-2, 4))

		self.assertEqual(len(interval), 6)

	def test_multi_interval(self):

		interval = Interval((2, 4), (5, 6))

		self.assertEqual(len(interval), 3)

	def test_multi_interval_with_intersection(self):

		interval = Interval((2, 5), (5, 6))

		self.assertEqual(len(interval), 4)

from pyperfstore3.timewindow import TimeWindow


class TimeWindowTest(unittest.TestCase):

	def setUp(self):
		self.timewindow = TimeWindow()

	def test_total_seconds(self):
		self.assertEqual(self.timewindow.total_seconds(), TimeWindow.DEFAULT_DURATION)

if __name__ == '__main__':
	unittest.main()
