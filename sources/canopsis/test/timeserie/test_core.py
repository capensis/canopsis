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

import unittest
from dateutil.relativedelta import relativedelta
from datetime import datetime
from math import isnan
from random import random, randint
from time import time, mktime
from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.timeserie.core import TimeSerie
from canopsis.timeserie.timewindow import TimeWindow, Period
import xmlrunner


class TimeSerieTest(unittest.TestCase):

    def setUp(self):
        self.conf = Configuration.load(TimeSerie.CONF_PATH, Ini)
        self.timeserie = TimeSerie(self.conf)
        points = [
            (ts, 1) for ts in range(0, 24 * 3600, 3600)
        ]
        self.timewindow = TimeWindow(start=points[0][0], stop=points[-1][0])
        self.points = points

    def _test_agg_per_x(self, period, len_points):

        self.timeserie.period = period

        agg_points = self.timeserie.calculate(
            points=self.points, timewindow=self.timewindow
        )

        self.assertEqual(len(agg_points), len_points)

    def test_aggregation_per_day(self):

        period = Period(day=1)

        self._test_agg_per_x(period, 1)

    def test_aggregation_per_6hours(self):

        period = Period(hour=6)

        self._test_agg_per_x(period, 4)

    def _five_years_timewidow(self):

        now = time()
        rd = relativedelta(years=5)
        now = datetime.now()
        past = now - rd
        past_ts = mktime(past.timetuple())
        now_ts = mktime(now.timetuple())
        result = TimeWindow(start=past_ts, stop=now_ts)

        return result

    def test_intervals(self):
        """Test calculate on different intervals."""

        now = time()

        # let a period of 1 day
        period = Period(day=1)
        oneday = period.total_seconds()

        rnow = period.round_timestamp(now)

        # let a timewindow of 10+1/4 days
        timewindow = TimeWindow(start=now - oneday, stop=now + 45/4 * oneday)

        nan = float('nan')

        points = [
            # the first interval is empty
            (rnow, nan),  # the second interval contains nan at start
            (rnow + oneday + 1, nan),  # the third interval contains nan at start + 1
            (rnow + 2 * oneday, 1),  # the fourth interval contains 1 at start
            (rnow + 3 * oneday + 1, 1),  # the fourth interval contains 1 at start + 1
            (rnow + 4 * oneday, nan), (rnow + 4 * oneday + 1, 1),  # the fith interval contains 1 and nan
            (rnow + 5 * oneday, 1), (rnow + 5 * oneday + 1, 1),  # the sixth interval contains 1 and 1
            (rnow + 6 * oneday, 1), (rnow + 6 * oneday, 1),  # the sixth interval contains 1 and 1 at the same time
            (rnow + 7 * oneday, nan), (rnow + 7 * oneday, nan),  # the sixth interval contains nan and nan at the same time
        ]

        timeserie = TimeSerie(
            config=self.conf,
            aggregation='sum',
            period=period,
            round_time=True
        )

        _points = timeserie.calculate(points, timewindow)

        for i in [0, 1, 2, 5, 8, 9, 10, 11, 12]:
            self.assertEqual(_points[i][0], rnow + (i - 1) * oneday)
            self.assertTrue(isnan(_points[i][1]))

        for i in [3, 4]:
            self.assertEqual(_points[i][0], rnow + (i - 1) * oneday)
            self.assertEqual(_points[i][1], 1)

        for i in [6, 7]:
            self.assertEqual(_points[i][0], rnow + (i - 1) * oneday)
            self.assertEqual(_points[i][1], 2)

        self.assertEqual(len(_points), len(points) + 1)

    def test_scenario(self):
        """
        Calculate aggregations over 5 years
        """

        timewindow = self._five_years_timewidow()

        # for all round_time values
        for round_time in (True, False):

            unit_length = 3600

            # for all units
            for index, unit in enumerate(Period.UNITS):

                max_value_unit = Period.MAX_UNIT_VALUES[index]

                if unit in (
                    Period.MICROSECOND,
                    Period.SECOND,
                    Period.MINUTE,
                    Period.WEEK,
                    Period.MONTH,
                    Period.YEAR
                ):
                    continue

                value = randint(2, max_value_unit)

                period = Period(**{unit: value})
                kwargs = {'period': period}
                period_length = unit_length * value

                timeserie = TimeSerie(config=self.conf,
                                      round_time=round_time,
                                      **kwargs)

                timesteps = timeserie.timesteps(timewindow)

                timesteps_gap = timesteps[1] - timesteps[0]

                self.assertEqual(timesteps_gap, period_length)

                for i in range(5):
                    points = [
                        (t, random()) for t in
                        range(
                            int(timewindow.start()),
                            int(timewindow.stop()),
                            Period(**{unit: 1}).total_seconds()
                        )
                    ]

                    aggregated_points = timeserie.calculate(points, timewindow)
                    len_aggregated_points = len(aggregated_points)
                    self.assertIn(
                        len(timesteps) - 1,
                        (
                            len_aggregated_points,
                            len_aggregated_points + 1
                        )
                    )

                unit_length *= max_value_unit

    def test_timesteps(self):
        timewindow = TimeWindow(start=1508420521, stop=1508506921, timezone=-120)

        steps = self.timeserie.timesteps(timewindow=timewindow)

        self.assertEqual(len(steps), 3)
        self.assertEqual(steps[0], 1508371200)

        # TODO: do more tests or rewrite perfdatas

    def test__get_period(self):
        timewindow = TimeWindow(start=1508420521, stop=1508506921, timezone=-120)

        period = self.timeserie._get_period(timewindow=timewindow)

        self.assertEqual(period.total_seconds(), 86400)

        # TODO: do more tests or rewrite perfdatas

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
