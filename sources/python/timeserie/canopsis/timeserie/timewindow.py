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

__all__ = ['Period', 'Interval', 'TimeWindow', 'get_offset_timewindow']

from time import time

from dateutil.relativedelta import relativedelta as rd
from dateutil.tz import tzoffset

from calendar import timegm, monthcalendar

from collections import Iterable

from operator import itemgetter

from datetime import datetime as dt


class Period(object):
    """Period management with a value and an unitself.
    """

    __slots__ = ['unit_values']

    MICROSECOND = 'microsecond'
    SECOND = 'second'
    MINUTE = 'minute'
    HOUR = 'hour'
    DAY = 'day'
    WEEK = 'week'
    MONTH = 'month'
    YEAR = 'year'

    UNITS = (MICROSECOND, SECOND, MINUTE, HOUR, DAY, WEEK, MONTH, YEAR)

    MAX_UNIT_VALUES = (10000000, 60, 60, 24, 7, 4, 31, 100)

    UNIT = 'unit'
    VALUE = 'value'

    def __init__(self, **unit_values):

        super(Period, self).__init__()

        self.unit_values = unit_values

    def __repr__(self):

        return "Period{0}".format(self.unit_values)

    def __eq__(self, other):

        result = isinstance(other, Period) and repr(self) == repr(other)

        return result

    def __getitem__(self, name):

        return self.unit_values[name]

    def __delitem__(self, name):

        del self.unit_values[name]

    def __setitem__(self, name, value):

        self.unit_values[name] = value

    def __iter__(self):

        result = iter(
            [unit for unit in Period.UNITS if unit in self.unit_values]
        )

        return result

    def __len__(self):
        """Get number of seconds.

        :return: this period in seconds. Approximation if period is in months
            or years.
        :rtype: int
        """

        return self.total_seconds()

    def __mul__(self, value):
        """Get a new period which is a factor of value by self.

        :param float value: new multiplicative factor.
        :return: new period equals to value * self.
        :rtype: Period
        """

        result_unit_values = {
            k: self.unit_values[k] * value for k in self.unit_values
        }

        result = Period(**result_unit_values)

        return result

    def __imul__(self, value):

        for unit in self.unit_values.copy():
            self.unit_values[unit] *= value

        return self

    def total_seconds(self):
        """Get number of seconds.

        :return: this period in seconds. Approximation if this period has
            months.
        :rtype: int
        """

        result = 0

        if Period.MICROSECOND in self:
            result += self[Period.MICROSECOND] * 10 ** -9

        if Period.SECOND in self:
            result += self[Period.SECOND]

        if Period.MINUTE in self:
            result += self[Period.MINUTE] * 60

        if Period.HOUR in self:
            result += self[Period.HOUR] * 3600

        if Period.DAY in self:
            result += self[Period.DAY] * 86400

        if Period.WEEK in self:
            result += self[Period.WEEK] * 604800

        if Period.MONTH in self:
            result += self[Period.MONTH] * 2592000

        if Period.YEAR in self:
            result += self[Period.YEAR] * 31536000

        return result

    def get_delta(self):
        """Get a delta object in order to add/remove a period on a datetime.

        :return: delta object in order to add/remove a period on a datetime.
        :rtype: relativedelta
        """

        unit_values = self.unit_values

        parameters = {
            (name + 's'): unit_values[name] for name in unit_values
        }

        result = rd(**parameters)

        return result

    def next_period(self):
        """Get next period with input step or none if next period can't be
        associated to a specific unit.

        Example: Period(minute=60).next_period() == Period(hour=1)
        """

        result = Period()

        counts_with_unit = zip(Period.UNITS, Period.MAX_UNIT_VALUES)
        previous_unit, previous_count = counts_with_unit.pop(-1)
        counts_with_unit.reverse()

        for unit, count in counts_with_unit:
            value = self.unit_values.get(unit)

            if value is not None:
                next_value = value * count
                result[previous_unit] = next_value

            previous_unit = unit

        return result

    def round_timestamp(self, timestamp, normalize=False, next_period=False):
        """Get round timestamp relative to an input timestamp.

        :param long timestamp: timestamp to round.
        :param bool normalize: normalization property.
        :param bool next_period: computes current period next timestamp.


        Example: Let a timestamp ``t`` related to the date: 2015/03/04 15:05.
        r = Period(week=1).round_timestamp(timestamp=t)
        In this case, r corresponds to "2015/03/01 15:05".
        If normalize equals True, r corresponds to "2015/03/01 00:00"
        """

        datetime = dt.utcfromtimestamp(float(timestamp))
        datetime = self.round_datetime(datetime=datetime, normalize=normalize)

        utctimetuple = datetime.utctimetuple()
        result = timegm(utctimetuple)

        # restore microsecond because utctimetuple() does not
        microseconds = datetime.microsecond * 0.000001
        result += microseconds

        if next_period:
            result += len(self)

        return result

    def round_datetime(self, datetime, normalize=False):
        """Calculate roudtime relative to an UTC date.
        normalize unsure to set to 0 for not given units under the minimal unit

        :param datetime datetime: datetime to round.
        :param bool normalize: normalize.

        Example: Let a datetime ``d`` related to the date: 2015/03/04 15:05.
        r = Period(week=1).round_datetime(datetime=d)
        In this case, r corresponds to "2015/m/d 15:05" where m/d is first
        monday before 2015/03/01.
        If normalize equals True, r corresponds to "2015/m/d 00:00"
        """

        result = None

        parameters = {}

        unit_values = self.unit_values

        for unit in unit_values:
            value = max(1, unit_values[unit])
            if unit == Period.WEEK:
                _monthcalendar = monthcalendar(
                    datetime.year, datetime.month
                )
                for week_index, week in enumerate(_monthcalendar):
                    if datetime.day in week:
                        datetime_value = week_index
                        break
            else:
                datetime_value = getattr(datetime, unit)
            rounding_period_value = datetime_value % value
            parameters[unit] = rounding_period_value

        rounding_period = Period(**parameters)

        delta = rounding_period.get_delta()

        result = datetime - delta

        if normalize:  # set to minimal value for all units before minimal unit
            parameters = {}
            if Period.MICROSECOND not in self:
                parameters[Period.MICROSECOND] = 0
                if Period.SECOND not in self:
                    parameters[Period.SECOND] = 0
                    if Period.MINUTE not in self:
                        parameters[Period.MINUTE] = 0
                        if Period.HOUR not in self:
                            parameters[Period.HOUR] = 0
                            # check week have to be normalized
                            if Period.WEEK in self:
                                # get the right monday
                                # save day, month and year which can change
                                day = datetime.day
                                month = datetime.month
                                year = datetime.year
                                # get week value
                                v = self[Period.WEEK]
                                # find the right week which corresponds to day
                                for week_index, week in enumerate(
                                    _monthcalendar
                                ):
                                    if day in week:
                                        # get the right normalized  week index
                                        norm_idx = (week_index // (v + 1)) * v
                                        # get the right normalized week
                                        norm_week = _monthcalendar[norm_idx]
                                        # get last monday
                                        day = norm_week[0]
                                        # if monday appeared previous month
                                        if day == 0:
                                            month -= 1
                                            if month == 0:
                                                month = 12
                                                year -= 1
                                                # update year
                                                parameters[Period.YEAR] = year
                                            # update month
                                            parameters[Period.MONTH] = month
                                            # get previous month calendar
                                            mc = monthcalendar(year, month)
                                            # get last old monday
                                            day = mc[-1][0]
                                        break
                                parameters[Period.DAY] = day
                            elif Period.DAY not in self:
                                parameters[Period.DAY] = 1
                                if Period.MONTH not in self:
                                    parameters[Period.MONTH] = 1
                                    if Period.YEAR not in self:
                                        parameters[Period.YEAR] = 0
            result = result.replace(**parameters)

        return result

    def get_max_unit(self):
        """Get a dictionary which contains a unit and a value
        where unit is the last among Period.UNITS.

        Example: period=Period(minute=10, hour=13)
        period.get_max_unit()  # equals {'hour': 13}
        """

        result = None

        units = list(Period.UNITS)
        units.reverse()

        for unit in units:
            if unit in self:
                result = {Period.UNIT: unit, Period.VALUE: self[unit]}

        return result

    def copy(self):
        """Get a period which is a copy of self.
        """

        result = Period(**self.unit_values)

        return result

    @staticmethod
    def from_str(serialized):
        """Get a Period from a string of shape "(unit=value,)+".

        :param str serialized: serialized period of shape "(unit=value,)+"
        :return: period from a str.
        :rtype: Period
        """
        params = {}

        splitted = serialized.split(',')

        for s in splitted:
            args = s.split('=')
            if len(args) == 2:
                params[args[0]] = float(args[1])
            else:
                # TODO: display an error
                pass

        result = Period(**params)

        return result


class Interval(object):
    """Manage points interval with sub intervals
    which are (lower value, upper value).
    """

    class IntervalError(Exception):
        pass

    __slots__ = ['sub_intervals']

    _NUMBER = (float, int, complex)

    def __init__(self, *intervals):

        super(Interval, self).__init__()

        self.sub_intervals = Interval.sort_and_join_intersections(*intervals)

    def __eq__(self, other):

        result = isinstance(other, Interval) and repr(self) == repr(other)

        return result

    def __repr__(self):

        result = "Interval{0}".format(self.sub_intervals)

        return result

    def __contains__(self, numbers_or_intervals):
        """True iif input values or intervals are in this interval.
        values_or_interval must be numbers or Intervals.
        """

        # return False by default.
        result = False

        def check_number_or_interval(number_or_interval, pos=None):
            """Check if input number_or_interval is in self.sub_intervals.
            """

            result = False

            if isinstance(number_or_interval, Iterable) \
                    and len(number_or_interval) == 2:

                result = number_or_interval[0] in self \
                    and number_or_interval[1] in self

            elif isinstance(number_or_interval, Interval):
                result = True

                for sub_interval in number_or_interval:
                    if sub_interval[0] not in self \
                            or sub_interval[1] not in self:
                        result = False
                        break

            elif isinstance(number_or_interval, Interval._NUMBER):

                for sub_interval in self:
                    if number_or_interval >= sub_interval[0] \
                            and number_or_interval <= sub_interval[1]:
                        result = True
                        break

            else:
                raise Interval.IntervalError(
                    "Wrong input parameter {0}({1}){2}."
                    .format(
                        number_or_interval,
                        type(number_or_interval),
                        "" if pos is None else "at pos {0}".format(pos)))

            return result

        if isinstance(numbers_or_intervals, Iterable):
            result = len(numbers_or_intervals) > 0

            for index, number_or_interval in enumerate(numbers_or_intervals):
                if not check_number_or_interval(number_or_interval, index):
                    result = False
                    break

        else:
            result = check_number_or_interval(numbers_or_intervals)

        return result

    def __len__(self):
        """Get number of values between all sub intervals.
        """
        result = 0

        for sub_interval in self:
            result += sub_interval[1] - sub_interval[0]

        return result

    def __or__(self, interval):

        result = None

        if isinstance(interval, Interval):
            result = Interval(self.sub_intervals + interval.sub_intervals)

        else:
            raise NotImplementedError()

        return result

    def __ior__(self, interval):

        if isinstance(interval, Interval):
            self.sub_intervals = (self | interval).sub_intervals

    def __and__(self, interval):

        raise NotImplementedError()

    def __sub__(self, interval):

        raise NotImplementedError()

    def __iter__(self):
        """Get self sub_intervals iterator.
        """

        return iter(self.sub_intervals)

    def __getitem__(self, key):
        """Get the right sub interval.
        """

        return self.sub_intervals[key]

    def min(self):
        """Get minimal point or None if no sub intervals.
        """

        return self.sub_intervals[0][0] if self.sub_intervals else None

    def max(self):
        """Get maximal point or None if no sub intervals.
        """

        return self.sub_intervals[-1][1] if self.sub_intervals else None

    def is_empty(self):
        """True iif this interval does not contain sub intervals.
        """

        result = len(self.sub_intervals) == 0

        return result

    @staticmethod
    def sort_and_join_intersections(*intervals):
        """Get intervals which are the result of a clean, sort and a join
        intersection operation on input intervals.
        Get an interval which is a cleanable version of all input intervals.

        Input intervals can be empty or contains Intervals or Iterable of
        two floats.
        """

        result = []

        for index, interval in enumerate(intervals):

            if isinstance(interval, Interval):
                result += tuple(interval.sub_intervals)

            elif isinstance(interval, Iterable):
                if len(interval) != 2:
                    raise Interval.IntervalError(
                        "Iterable interval {0} at pos {1} must contain only \
                        two elements"
                        .format(interval, index))

                if isinstance(interval[0], Interval._NUMBER) \
                        and isinstance(interval[1], Interval._NUMBER):
                    result.append(tuple(interval))

                else:
                    raise Interval.IntervalError(
                        "Wrong input interval {0} at pos {1}"
                        .format(interval, index))

            elif isinstance(interval, Interval._NUMBER):
                sub_interval = (0, interval) if interval > 0 else (interval, 0)
                result.append(sub_interval)

            else:
                raise Interval.IntervalError(
                    "Wrong input interval {0} at pos {1}"
                    .format(interval, index))

        # sort intervals
        result, _result = [], sorted(result, key=itemgetter(0))

        index = 0

        while index < len(_result):
            interval = _result[index]

            index += 1

            for _index in range(index, len(_result)):
                _interval = _result[_index]

                index = _index

                if _interval[0] >= interval[0] and _interval[0] <= interval[1]:
                    interval = (interval[0], max(interval[1], _interval[1]))
                    index += 1

                else:
                    break

            result.append(interval)

        result = tuple(result)

        return result

    def reduce(self, lower, upper):
        """Returns an interval reduced with input lower and upper bounds.

        :param float lower: lower bound.
        :param float upper: upper bound.
        :return: intersection of this and [lower; upper].
        :rtype: Interval
        """
        # list which will be used to return an interval
        intervals = []

        for sub_interval in self:
            sub_lower, sub_upper = sub_interval
            # do something until sub interval & [lower, upper] != 0
            if sub_upper < lower:
                continue
            if upper < sub_lower:
                break
            # update sub interval if inside [lower, upper]
            if sub_lower < lower < sub_upper:
                sub_lower = lower
            if sub_upper > upper > sub_lower:
                sub_upper = upper
            # add sub interval in final intervals
            intervals.append((sub_lower, sub_upper))

        return Interval(*intervals)


class TimeWindow(object):
    """Manage second intervals with a timezone.
    """

    class TimeWindowError(Exception):
        pass

    DEFAULT_DURATION = 60 * 60 * 24  # one day

    __slots__ = ['interval', 'timezone']

    def __init__(self, start=None, stop=None, intervals=None, timezone=0):
        """This interval is created from:

        - an interval with stop, start :
            - stop is now if None,
            - start is stop - TimeWindow.DEFAULT_DURATION if None,
            - if intervals is not empty and start and stop equal None
            then they are not calculated.
        - intervals is a list of (lower timestamp, upper timestamp) or
            Interval.
        """

        super(TimeWindow, self).__init__()

        # initialize intervals
        if intervals is None:
            intervals = []
        elif isinstance(intervals, Interval):
            intervals = [intervals]
        # initialize start/stop related to intervals
        if len(intervals):
            if start is None:
                start = intervals[0][0]
            if stop is None:
                stop = intervals[-1][0]

        # if stop is None, stop = now
        if stop is None:
            stop = round(time())

        # if start is None, start is stop - TimeWindow.DEFAULT_DURATION
        if start is None:
            start = stop - TimeWindow.DEFAULT_DURATION

        # if no interval is proposed, initialize it with start and stop
        if not intervals:
            intervals.append((start, stop))

        interval = Interval(*intervals)

        if interval.is_empty():
            raise TimeWindow.TimeWindowError("Interval can not be empty")

        self.interval = TimeWindow.convert_to_seconds_interval(interval)

        self.timezone = timezone

    def __eq__(self, other):

        result = isinstance(other, TimeWindow) and repr(self) == repr(other)

        return result

    def __repr__(self):

        message = "TimeWindow(tz:{0}):{1}"
        result = message.format(self.timezone, self.interval)

        return result

    def __contains__(self, *timestamps):
        """True if input timestamps are in this timewindow.
        """

        result = timestamps in self.interval

        return result

    def copy(self):
        """Get a copy of self.
        """

        result = TimeWindow(
            intervals=[Interval(self.interval)], timezone=self.timezone
        )

        return result

    def start(self):
        """Get first timestamp.
        """

        result = float(self.interval.min())

        return result

    def start_datetime(self, utc=False):
        """Get start datetime.
        """

        result = TimeWindow.get_datetime(
            self.start(), - self.timezone if utc else 0
        )

        return result

    def stop(self):
        """Get last timestamp.
        """

        result = float(self.interval.max())

        return result

    def stop_datetime(self, utc=False):
        """Get stop datetime.
        """

        result = TimeWindow.get_datetime(
            self.stop(), self.timezone if utc else 0
        )

        return result

    def total_seconds(self):
        """Returns seconds inside this timewindow.
        """

        result = len(self.interval)

        return result

    @staticmethod
    def get_timestamp(datetime):
        """Get the timestamp corresponding to input datetime.
        """

        result = timegm(datetime.timetuple())

        return result

    @staticmethod
    def get_datetime(timestamp, timezone=0):
        """Get the datetime corresponding to both input timestamp and timezone.
        """

        tz = tzoffset(None, timezone)
        result = dt.fromtimestamp(timestamp, tz)

        return result

    @staticmethod
    def convert_to_seconds_interval(interval):
        """Get interval in seconds from an interval.
        """

        sub_intervals = [
            (int(sub_interval[0]), int(sub_interval[1]))
            for sub_interval in interval
        ]

        result = Interval(*sub_intervals)

        return result

    def reduce(self, start, stop):
        """Returns a timewindow where start and stop are redefined.

        :param float start: new start time.
        :param float stop: new stop time.
        :returns: new timewindow with start/stop such as lower/upper bounds.
        :rtype: TimeWindow
        """

        interval = self.interval.reduce(start, stop)
        result = TimeWindow(intervals=interval, timezone=self.timezone)
        return result


def get_offset_timewindow(offset=time()):
    """Get a timewindow with one point.
    """

    return TimeWindow(start=offset, stop=offset)
