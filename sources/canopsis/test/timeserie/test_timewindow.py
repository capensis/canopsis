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

from datetime import datetime
from time import time, mktime, timezone
from calendar import monthrange

from random import random, randint

from canopsis.timeserie.timewindow import Period, Interval, TimeWindow


class PeriodTest(TestCase):

    def test_new(self):

        args = int(36)

        period = Period.new(args)

        self.assertEqual(period.total_seconds(), args)

        args = float(args)

        period = Period.new(args)

        self.assertEqual(period.total_seconds(), args)

        period = Period.new({'second': args})

        self.assertEqual(period.total_seconds(), args)

    @staticmethod
    def _new_period():

        unit_values = {}

        for unit in Period.UNITS:
            unit_values[unit] = random() * 10

        result = Period(**unit_values)

        return result

    def test_copy(self):

        period = PeriodTest._new_period()

        copy = period.copy()

        self.assertEqual(copy, period)
        self.assertFalse(copy is period)

    def test_delta(self):
        period = PeriodTest._new_period()

        delta = period.get_delta()

        for unit in period.unit_values:
            value = period.unit_values[unit]

            if unit in (Period.WEEK, Period.DAY):
                value_to_compare = int(delta.days)
                value = int(
                    period.unit_values.get('day', 0) + 7 * period.unit_values.get('week', 0)
                )

            else:
                unit = '{0}s'.format(unit)
                value_to_compare = getattr(delta, unit)

            self.assertEqual(value, value_to_compare)

    def test_round_datetime(self, ts=None):

        now = datetime.now() if ts is None else datetime.utcfromtimestamp(ts)

        for unitindex, unit in enumerate(Period.UNITS):

            if unit in (Period.MONTH, Period.YEAR):
                continue

            for i in range(0, 5):

                period = Period(**{unit: i})
                round_dt = period.round_datetime(now)
                round_value = getattr(round_dt, unit, None)

                unitval = getattr(now, unit, None)

                if unitval is not None:

                    if unit is Period.DAY:

                        roundday = unitval - (unitval % max(1, i))
                        if roundday < 0:
                            month = now.month - 1

                            if month == 0:
                                month = 12

                            _, roundday = monthrange(now.year, month)

                        self.assertEqual(round_value, roundday)

                    else:
                        maxunitval = Period.MAX_UNIT_VALUES[unitindex]

                        rval = (unitval - (unitval % max(1, i))) % maxunitval

                        self.assertEqual(rval, round_value)

                    for _unit in Period.UNITS[:max(0, unitindex - 1)]:
                        _round_value = getattr(round_dt, _unit, None)

                        if _unit in (Period.MONTH, Period.DAY):
                            _round_value -= 1

                        self.assertFalse(_round_value)

    def test_round_timestamp(self):

        self.test_round_datetime(int(time()))

    def test_get_max_unit(self):

        period = self._new_period()

        max_unit = period.get_max_unit()

        self.assertTrue(max_unit[Period.UNIT], Period.YEAR)

    def test_units_and_values(self):

        seconds = 1

        value = randint(0, 10**9-1)

        total_seconds = Period(microsecond=value).total_seconds()

        self.assertEqual(total_seconds, value * 10**-9)

        for index, unit in enumerate(Period.UNITS[1:-1]):

            value = randint(1, 10**10)

            total_seconds = Period(**{unit: value}).total_seconds()

            self.assertEqual(total_seconds, seconds * value, '{0}:{1}:{2}'.format(value, total_seconds, unit))

            seconds *= Period.MAX_UNIT_VALUES[index + 1]

    def test_total_seconds_mix(self):
        """
        Test total seconds with all units
        """

        kwargs = {
            Period.MICROSECOND: 1,
            Period.SECOND: 1,
            Period.MINUTE: 1,
            Period.HOUR: 1,
            Period.DAY: 1,
            Period.WEEK: 1,
            Period.MONTH: 1,
            Period.YEAR: 1
        }

        period = Period(**kwargs)

        self.assertEqual(
            period.total_seconds(),
            10**-9 + 1 + 60 + 3600 + 86400 + 86400*7 + 86400*7*4 + 86400*7*4*12
        )

    def test_mul(self):

        p = Period(s=5, mn=10)

        p1 = p * 5

        self.assertEqual(p1.unit_values, {'s': 5 * 5, 'mn': 10 * 5})


class IntervalTest(TestCase):

    def test_copy(self):

        sub_intervals = list()
        for i in range(randint(1, 99)):
            sub_intervals += (i - random(), i + random())

        interval = Interval(*sub_intervals)

        copy = Interval(interval)

        self.assertEqual(copy, interval)

    def test_is_empty(self):
        interval = Interval()

        self.assertTrue(interval.is_empty())

        interval = Interval(10 ** -99)

        self.assertFalse(interval.is_empty())

        interval = Interval(0)

        self.assertFalse(interval.is_empty())

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


class IntervalReductionTest(TestCase):
    """
    Test to reduce an interval.
    """

    def setUp(self):
        """
        Initialize lower, upper
        """

        self.lower, self.upper = 0, 5000
        self.length = self.upper - self.lower

    def test_empty(self):
        """
        Test to reduce empty interval.
        """

        # start with empty intervals
        interval = Interval()
        reduced = interval.reduce(self.lower, self.upper)
        self.assertEqual(len(reduced), 0)

    def test_one(self):
        """
        Test to reduce one interval.
        """

        # test same interval
        interval = Interval((self.lower, self.upper))
        reduced = interval.reduce(self.lower, self.upper)
        self.assertEqual(len(reduced), self.length)

        # remove 2 points
        interval = Interval((self.lower + 1, self.upper - 1))
        reduced = interval.reduce(self.lower, self.upper)
        self.assertEqual(len(reduced), self.length - 2)

        # add 2 points
        interval = Interval((self.lower - 1, self.upper + 1))
        reduced = interval.reduce(self.lower, self.upper)
        self.assertEqual(len(reduced), self.length)

    def test_intervals(self):
        """
        Test intervals.
        """

        # use five intervals where union & reduced interval = reduced interval
        interval = Interval(
            (self.lower - 5, self.lower - 2),
            (self.lower - 1, self.lower + 1),
            (self.lower + 2, self.upper - 2),
            (self.upper - 1, self.upper + 1),
            (self.upper + 2, self.upper + 5)
        )
        reduced = interval.reduce(self.lower, self.upper)
        self.assertEqual(len(reduced), self.length - 2)


class TimeWindowTest(TestCase):

    def setUp(self):
        self.timewindow = TimeWindow()

    def test_copy(self):

        copy = self.timewindow.copy()

        self.assertEqual(copy, self.timewindow)

    def test_total_seconds(self):
        self.assertEqual(
            self.timewindow.total_seconds(),
            TimeWindow.DEFAULT_DURATION)

    def test_start_stop(self):

        start = random() * 10000
        stop = start + random() * 10000
        timewindow = TimeWindow(start=start, stop=stop)
        TimeWindow()

        self.assertEqual(timewindow.start(), int(start))
        self.assertEqual(timewindow.stop(), int(stop))

    def test_get_datetime(self):

        now = time()

        dt = TimeWindow.get_datetime(now)
        ts_now = mktime(dt.timetuple())

        ri = randint(1, 500000)

        dt = TimeWindow.get_datetime(now + ri)
        self.assertEqual(ts_now + ri, mktime(dt.timetuple()))

        dt = TimeWindow.get_datetime(now, timezone)
        ts = mktime(dt.timetuple())

        self.assertEqual(ts, ts_now + timezone)

    def test_no_startstop(self):

        start, stop = 5, 10
        intervals = Interval((start, stop))
        timewindow = TimeWindow(intervals=intervals)
        self.assertEqual(timewindow.start(), start)
        self.assertEqual(timewindow.stop(), stop)


if __name__ == '__main__':
    main()
