from unittest import main, TestCase

from canopsis.timeserie import TimeSerie
from canopsis.timeserie.timewindow import TimeWindow, Period
from random import random, randint
from time import time, mktime
from dateutil.relativedelta import relativedelta
from datetime import datetime


class TimeSerieTest(TestCase):

    def setUp(self):
        self.timeserie = TimeSerie()
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

                timeserie = TimeSerie(round_time=round_time, **kwargs)

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


if __name__ == '__main__':
    main()
