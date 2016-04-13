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

"""Timeserie module."""

# provide only TimeSerie
__all__ = ['TimeSerie']

from canopsis.common.init import basestring
from canopsis.timeserie.timewindow import Period, TimeWindow
from canopsis.timeserie.aggregation import get_aggregation, DELTA, NONE
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.model import Parameter

from math import isnan

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

        # set protected attributes
        self._period = None
        self._aggregation = None

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
            self.max_points
        )

        return result

    def __eq__(self, other):

        result = isinstance(other, TimeSerie) and repr(self) == repr(other)

        return result

    @property
    def aggregation(self):
        """Get this aggregation method."""

        return self._aggregation.upper()

    @aggregation.setter
    def aggregation(self, value):
        """Change of aggregation method."""

        self._aggregation = value

    @property
    def period(self):
        """Get this period."""

        return self._period

    @period.setter
    def period(self, value):
        """Change of period."""

        if isinstance(value, basestring):
            value = Period.from_str(value)

        self._period = value

    def timesteps(self, timewindow):
        """Get a list of same longer intervals inside timewindow.

        The upper bound is timewindow.stop_datetime()
        """

        result = []
        # get the right period to apply on timewindow
        period = self._get_period(timewindow=timewindow)

        # set start and stop datetime
        start_datetime = timewindow.start_datetime()
        stop_datetime = timewindow.stop_datetime()

        if self.round_time:  # normalize if round time is True
            start_datetime = period.round_datetime(datetime=start_datetime)

        current_datetime = start_datetime
        delta = period.get_delta()

        while current_datetime < stop_datetime:
            timestamp = TimeWindow.get_timestamp(current_datetime)
            result.append(timestamp)
            current_datetime += delta

        result.append(timewindow.stop())

        return result

    def calculate(self, points, timewindow, meta=None, usenan=True):
        """Do an operation on all points with input timewindow.

        :param bool usenan: if False (True by default) remove nan point values.
        :return: points such as follow:
            Let func self aggregation function and
            input points of the form: [(T0, V0), ..., (Tn, Vn)]
            then the result is [(T0, func(V0, V1)), (T2, func(V2, V3), ...].
        """

        result = []

        nan = float('nan')

        # start to exclude points not in timewindow
        # in taking care about round time
        if self.round_time:
            period = self._get_period(timewindow)
            round_starttimestamp = period.round_timestamp(
                timestamp=timewindow.start()
            )
            timewindow = timewindow.reduce(
                start=round_starttimestamp,
                stop=timewindow.stop()
            )

        # start to exclude points which are not in timewindow
        points = [
            point for point in points
            if point[0] in timewindow and (usenan or not isnan(point[1]))
        ]

        if not meta:
            meta = {}

        transform_method = meta.get('type')
        points = apply_transform(points, method=transform_method)
        points_len = len(points)

        func = None

        # if no period and max_points > len(points)
        if (
                get_aggregation(self.aggregation) is None or
                (
                    (not points) or self.period is None
                    and self.max_points > points_len
                )
        ):
            result = points  # result is points

        else:  # else calculate points with the right aggregation function
            func = get_aggregation(name=self.aggregation)

            # get timesteps
            timesteps = self.timesteps(timewindow)[:-1]

            # initialize variables for loop
            i = 0
            values_to_aggregate = []
            last_point = None

            len_timesteps = len(timesteps)

            # iterate on timesteps to get points in [prev_ts, timestamp[
            for index, timestamp in enumerate(timesteps):
                # initialize values_to_aggregate
                values_to_aggregate = []
                # set timestamp and previous_timestamp

                if index < (len_timesteps - 1):  # calculate the upper bound
                    next_timestamp = timesteps[index + 1]

                else:
                    next_timestamp = None

                # fill the values_to_aggregate array
                for i in range(i, points_len):  # while points to process

                    pt_ts, pt_val = points[i]

                    # leave the loop if _timestamp is for a future aggregation
                    if next_timestamp is not None and pt_ts >= next_timestamp:
                        break

                    else:
                        # add value to list of values to aggregate
                        values_to_aggregate.append(pt_val)

                else:  # leave the loop whatever timesteps
                    i += 1

                # TODO: understand what it means :D
                if self.aggregation == DELTA and last_point:
                    values_to_aggregate.insert(0, last_point)

                if values_to_aggregate:

                    # get the aggregated value related to values_to_aggregate
                    aggregation_value = self._aggregation_value(
                        values_to_aggregate, func
                    )

                    # new point to add to result
                    if usenan or not isnan(aggregation_value):
                        aggregation_point = timestamp, aggregation_value
                        result.append(aggregation_point)
                        # save last_point for future DELTA checking
                        last_point = aggregation_point[-1]

                elif usenan:
                    result.append((timestamp, nan))

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
        """Get a period related to input max_points or a period."""

        result = self.period

        if result is None:
            seconds = (
                (timewindow.stop() - timewindow.start()) / self.max_points
            )
            result = Period(second=seconds)

        return result

    def _aggregation_value(self, values_to_aggregate, func=None):
        """Get the aggregated value related to input values_to_aggregate, a
        specific function and a timestamp.
        """

        if func is None:
            func = get_aggregation(name=self.aggregation)

        if len(values_to_aggregate) > 0:
            result = round(func(values_to_aggregate), 2)

        else:
            result = 0 if self.fill else None

        return result


def gauge(pts):
    """Calculate gauge."""

    return pts


def absolute(pts):
    """Calculate gauge."""

    return list([point[0], abs(point[1])] for point in pts)


def derive(pts):
    """calculate derive."""

    result = []

    for i in range(1, len(pts)):
        timestamp, val = pts[i]
        prevts, prevval = pts[i - 1]

        if val > prevval:
            val -= prevval

        interval = abs(timestamp - prevts)
        if interval:
            val = round(float(val) / interval, 3)

        result.append([timestamp, val])

    return result


def counter(pts):
    """Calculate counter."""

    result = []

    val = 0

    for point in pts:
        timestamp, increment = point
        val += increment

        result.append([timestamp, val])

    return result


METHODS = {
    'GAUGE': gauge,
    'ABSOLUTE': absolute,
    'DERIVE': derive,
    'COUNTER': counter
}

def apply_transform(points, method=None):
    """Apply DERIVE, ABSOLUTE, COUNTER, GAUGE transforms to points.

    :param list points: list of points.
    :param str method: method (it's the "type" metadata of perfdata).
    :returns: list of points.
    :rtype: list
    """

    result = points

    if method and method in METHODS:
        result = METHODS[method](points)

    return result
