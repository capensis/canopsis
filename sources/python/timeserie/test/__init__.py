from unittest import main, TestCase

from canopsis.timeserie import TimeSerie
from canopsis.timeserie.timewindow import TimeWindow, Period
from random import random, randint
from time import time


class AggregationTest(TestCase):

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
        # create an interval of 5 year every 30 minutes
        five_year = 60 * 60 * 24 * 365 * 5  # five years in seconds
        result = TimeWindow(start=now - five_year, stop=now)

        return result

    def test_scenario(self):

        timewindow = self._five_years_timewidow()

        total_seconds = timewindow.total_seconds()

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
                    Period.MONTH,
                    Period.YEAR
                ):
                    continue

                if unit is Period.WEEK:
                    value = randint(2, 500)
                    kwargs = {'max_points': value}
                    period_length = total_seconds / value
                else:
                    value = randint(2, max_value_unit)
                    kwargs = {'period': Period(**{unit: value})}
                    period_length = unit_length * value

                timeserie = TimeSerie(
                    aggregation="MEAN", round_time=round_time, **kwargs
                )

                timesteps = timeserie.timesteps(timewindow)

                self.assertEqual(timesteps[1] - timesteps[0], period_length)

                all_aggregated_points = list()

                for i in range(5):

                    points = [
                        (t, random()) for t in range(
                            int(timewindow.start()),
                            int(timewindow.stop()), 3600)]

                    aggregated_points = timeserie.calculate(points, timewindow)

                    all_aggregated_points.append(aggregated_points)

                    len_aggregated_points = len(aggregated_points)
                    self.assertTrue((len(timesteps) - 1) in (
                        len_aggregated_points, len_aggregated_points + 1))

                unit_length *= max_value_unit


if __name__ == '__main__':
    main()
