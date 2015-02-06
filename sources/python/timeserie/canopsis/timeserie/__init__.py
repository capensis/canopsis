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

__version__ = "0.1"

# provide only TimeSerie
__all__ = ('TimeSerie')

from canopsis.common.init import basestring
from canopsis.timeserie.timewindow import Period, TimeWindow
from canopsis.timeserie.aggregation import get_aggregations
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.parameters import Parameter

CONF_PATH = 'timeserie/timeserie.conf'


@conf_paths(CONF_PATH)
class TimeSerie(Configurable):
    """
    Time serie management. Contain a period and operation of aggregation,
    and a round time and a fill boolean properties.

    - period: interval of steps of aggregated points.
    - aggregation: aggregation operation name.
    - round_time: round_time of input timewindow during calculations.
    - fill: change None values by 0.
    """

    __slots__ = ('period', 'max_points', 'round_time', 'aggregation', 'fill')

    MAX_POINTS = 'max_points'
    AGGREGATION = 'aggregation'
    FILL = 'fill'
    PERIOD = 'period'
    ROUND_TIME = 'round_time'

    VMAX_POINTS = 500
    VDEFAULT_AGGREGATION = 'MEAN'
    VPERIOD = Period(day=1)
    VFILL = False
    VROUND_TIME = True

    CATEGORY = 'TIMESERIE'

    def __init__(
        self,
        aggregation=VDEFAULT_AGGREGATION,
        max_points=VMAX_POINTS,
        period=None,
        round_time=VROUND_TIME,
        fill=VFILL,
        *args, **kwargs
    ):

        super(TimeSerie, self).__init__(*args, **kwargs)

        self.period = period
        self.max_points = max_points
        self.round_time = round_time
        self.aggregation = aggregation
        self.fill = fill

    def __repr__(self):

        message = "TimeSerie(period: {0}, round_time: {1}" + \
            ", aggregation: {2}, fill: {3}, max_points: {4})"
        result = message.format(
            self.period, self.round_time, self.aggregation, self.fill,
            self.max_points)

        return result

    def __eq__(self, other):

        result = isinstance(other, TimeSerie) and repr(self) == repr(other)

        return result

    @property
    def max_points(self):
        return self._max_points

    @max_points.setter
    def max_points(self, value):
        self._max_points = value

    @property
    def aggregation(self):
        return self._aggregation.upper()

    @aggregation.setter
    def aggregation(self, value):
        self._aggregation = value

    @property
    def fill(self):
        return self._fill

    @fill.setter
    def fill(self, value):
        self._fill = value

    @property
    def period(self):
        return self._period

    @period.setter
    def period(self, value):
        if isinstance(value, basestring):
            value = Period.from_str(value)
        self._period = value

    @property
    def round_time(self):
        return self._round_time

    @round_time.setter
    def round_time(self, value):
        self._round_time = value

    def timesteps(self, timewindow):
        """
        Get a list of same longer intervals inside timewindow.
        The upper bound is timewindow.stop_datetime()
        """

        result = []
        # get the right period to apply on timewindow
        period = self._get_period(timewindow=timewindow)

        # set start and stop datetime
        start_datetime = timewindow.start_datetime()
        stop_datetime = timewindow.stop_datetime()

        if self.round_time:  # normalize if round time is True
            start_datetime = period.round_datetime(
                datetime=start_datetime, normalize=True
            )

        current_datetime = start_datetime
        delta = period.get_delta()

        while current_datetime < stop_datetime:
            ts = TimeWindow.get_timestamp(current_datetime)
            result.append(ts)
            current_datetime += delta

        result.append(timewindow.stop())

        return result

    def calculate(self, points, timewindow, meta=None):
        """
        Do an operation on all points with input timewindow.

        Return points su as follow:
        Let fn self aggregation function and
        input points of the form: [(T0, V0), ..., (Tn, Vn)]
        then the result is [(T0, fn(V0, V1)), (T2, fn(V2, V3), ...].
        """

        result = []

        # start to exclude points not in timewindow
        # in taking care about round time
        if self.round_time:
            period = self._get_period(timewindow)
            round_starttimestamp = period.round_timestamp(
                timestamp=timewindow.start(),
                normalize=True
            )
            timewindow = timewindow.reduce(
                start=round_starttimestamp,
                stop=timewindow.stop()
            )

        # start to exclude points which are not in timewindow
        points = [point for point in points if point[0] in timewindow]

        if not meta:
            meta = {}

        transform_method = meta.get('value', {}).get('type', None)
        points = self.apply_transform(points, method=transform_method)
        points_len = len(points)

        fn = None

        # if no period and max_points > len(points)
        if (not points) or self.period is None\
                and self.max_points > points_len:
            result = points  # result is points

        else:  # else get the right aggregation function
            fn = get_aggregations()[self.aggregation]

        # if an aggregation is required
        if fn is not None:

            # get timesteps
            timesteps = self.timesteps(timewindow)
            # initialize variables for loop
            i = 0
            values_to_aggregate = []
            last_point = None

            # iterate on all timesteps in order to get points in [prev_ts, ts[
            for index in range(1, len(timesteps)):
                # initialize values_to_aggregate
                values_to_aggregate = []
                # set timestamp and previous_timestamp
                previous_timestamp = timesteps[index - 1]
                timestamp = timesteps[index]

                # if no point to process in [previous_timestamp, timestamp]
                if points[i][0] > timestamp:
                    continue  # go to the next iteration

                # fill the values_to_aggregate array
                while i < points_len:  # while there are points to process

                    _timestamp, value = points[i]
                    i += 1

                    # leave the loop if _timestamp is for a future aggregation
                    if _timestamp > timestamp:
                        break

                    # if _timestamp is in timewindow and value is not None
                    # add value to list of values to aggregate
                    if value is not None:
                        values_to_aggregate.append(value)

                # TODO: understand what it means :D
                if self.aggregation == "DELTA" and last_point:
                    values_to_aggregate.insert(0, last_point)

                # get the aggregated value related to values_to_aggregate
                _aggregation_value = self._aggregation_value(
                    values_to_aggregate, fn)

                # new point to add to result
                aggregation_point = previous_timestamp, _aggregation_value
                result.append(aggregation_point)

                # save last_point for future DELTA checking
                if len(values_to_aggregate) > 0:
                    last_point = values_to_aggregate[-1]

                if i >= points_len:
                    break

        else:
            result = points

        return result

    def apply_transform(self, points, method=None):
        """
        Apply DERIVE, ABSOLUTE, COUNTER, GAUGE transforms to points.

        :param list points: list of points.
        :param str method: method (it's the "type" metadata of perfdata).
        :returns: list of points.
        :rtype: list
        """

        def gauge(pts):
            return pts

        def absolute(pts):
            points = []

            for pt in pts:
                ts, val = pt[0], pt[1]

                if val < 0:
                    val = -val

                points.append([ts, val])

            return points

        def derive(pts):
            points = []

            for i in range(1, len(pts)):
                ts, val = pts[i][0], pts[i][1]
                prevts, prevval = pts[i - 1][0], pts[i - 1][1]

                if val > prevval:
                    val -= prevval

                interval = abs(ts - prevts)
                if interval:
                    val = round(float(val) / interval, 3)

                points.append([ts, val])

            return points

        def counter(pts):
            points = []

            val = 0

            for pt in pts:
                ts, increment = pt[0], pt[1]
                val += increment

                points.append([ts, val])

            return points

        methods = {
            'GAUGE': gauge,
            'ABSOLUTE': absolute,
            'DERIVE': derive,
            'COUNTER': counter
        }

        result = points

        if method and method in methods:
            result = methods[method](points)

        return result

    def _conf(self, *args, **kwargs):

        result = super(TimeSerie, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=TimeSerie.CATEGORY,
            new_content=(
                Parameter(TimeSerie.AGGREGATION),
                Parameter(TimeSerie.PERIOD, parser=Period.from_str),
                Parameter(TimeSerie.FILL, parser=Parameter.bool),
                Parameter(TimeSerie.ROUND_TIME, parser=Parameter.bool),
                Parameter(TimeSerie.MAX_POINTS, parser=int)
            )
        )

        return result

    def _get_period(self, timewindow):
        """
        Get a period related to input max_points or a period.
        """

        result = self.period

        if result is None:
            seconds = \
                (timewindow.stop() - timewindow.start()) / self.max_points
            result = Period(second=seconds)

        return result

    def _aggregation_value(self, values_to_aggregate, fn=None):
        """
        Get the aggregated value related to input values_to_aggregate,
        a specific function and a timestamp.
        """

        if fn is None:
            fn = get_aggregations()[self.aggregation]

        if len(values_to_aggregate) > 0:
            result = round(fn(values_to_aggregate), 2)

        else:
            result = 0 if self.fill else None

        return result
