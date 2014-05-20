#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

import logging

logger = logging.getLogger('utils')
logger.setLevel(logging.DEBUG)

from pyperfstore3.timewindow import Period

from datetime import datetime
import calendar


# TODO: add a dynamic aggregation operator manager

class TimeSerie(object):
	"""
	Time serie management. Contain a period and operation of aggregation,
	and a round time and a fill boolean properties.

	- period: interval of steps of aggregated points.
	- aggregation: aggregation operation name.
	- round_time: round_time of input timewindow during calculations.
	- fill: change None values by 0.
	"""

	__slots__ = ('period', 'round_time', 'aggregation', 'fill')

	MAX_POINTS = 500

	MEAN = 'MEAN'
	LAST = 'LAST'
	FIRST = 'FIRST'
	DELTA = 'DELTA'
	SUM = 'SUM'
	MAX = 'MAX'
	MIN = 'MIN'

	def __init__(
		self,
		max_points=MAX_POINTS,
		period=None,
		round_time=True,
		aggregation=MEAN,
		fill=False):

		self.period = self._get_period(max_points=max_points, period=period)
		self.round_time = round_time
		self.aggregation = aggregation
		self.fill = fill

	def __repr__(self):

		message = "TimeSerie(period: {0}, round_time: {1}" + \
			", aggregation: {2}, fill: {3})"
		result = message.format(
			self.period, self.round_time, self.aggregation, self.fill)

		return result

	def _get_period(self, max_points, period):
		"""
		Get a period related to input max_points or a period.
		"""

		result = period.copy() if period is not None else None

		if result is None and max_points:
			interval = self.timewindow.stop - self.timewindow.start
			# may reduce value with timewindow.exclusion_intervals
			seconds = interval.total_seconds() / max_points
			result = Period(value=seconds)

		return result

	def getTimeSteps(self, timewindow):
		"""
		Get a list of same longer intervals inside timewindow.
		"""

		logger.debug('getTimeSteps:')

		timeSteps = []

		logger.debug(' + TimeSerie: {0}'.format(self))

		start_datetime = datetime.utcfromtimestamp(timewindow.start())
		stop_datetime = datetime.utcfromtimestamp(timewindow.stop())

		if self.round_time:
			stop_datetime = self.period.round_time(stop_datetime, self.timezone)

		date = stop_datetime
		start_delta = self.period.get_delta(start_datetime)
		start_datetime_minus_delta = start_datetime - start_delta

		date_delta = self.period.get_delta(date)

		while date > start_datetime_minus_delta:
			ts = calendar.timegm(date.timetuple())
			timeSteps.append(ts)
			date_delta = self.period.get_delta(date, date_delta)
			date -= date_delta

		timeSteps.reverse()

		logger.debug(' + timeSteps: {0}'.format(timeSteps))

		return timeSteps

	def get_operation(self):
		"""
		Get self aggregation.
		"""

		if self.aggregation == TimeSerie.MEAN or not self.aggregation:
			result = lambda x: sum(x) / float(len(x))
		elif self.aggregation == TimeSerie.MIN:
			result = lambda x: min(x)
		elif self.aggregation == TimeSerie.MAX:
			result = lambda x: max(x)
		elif self.aggregation == TimeSerie.FIRST:
			result = lambda x: x[0]
		elif self.aggregation == TimeSerie.LAST:
			result = lambda x: x[-1]
		elif self.aggregation == TimeSerie.DELTA:
			result = lambda x: (max(x) - min(x)) / 2.

		return result

	def calculate_points(self, points, timewindow):
		"""
		Do self aggregation on input points with input timewindow.
		"""

		result = []

		logger.debug("Calculate {0} point(s) (timeserie: {1}, timewindow: {2})"\
			.format(len(points), self, timewindow))

		agfn = self.get_operation()

		if len(points) == 1:
			return [[timewindow.start, points[0][1]]]

		timeSteps = self.getTimeSteps(timewindow)

		#initialize variables for loop
		i = 0
		points_to_aggregate = []
		last_point = None

		for index in range(1, len(timeSteps)):

			timestamp = timeSteps[index]

			previous_timestamp = timeSteps[index-1]

			logger.debug(
				" + Interval {0} -> {1}".format(
					(previous_timestamp, timestamp)))

			while i < len(points) and points[i][0] < timestamp:

				points_to_aggregate.append(points[i][1])

				i += 1

			if self.aggregation == TimeSerie.DELTA and last_point:
				points_to_aggregate.insert(0, last_point)

			aggregation_value = self.get_aggregation_value(
				points_to_aggregate, agfn)

			aggregation_point = previous_timestamp, aggregation_value

			result.append(aggregation_point)

			if points_to_aggregate:
				last_point = points_to_aggregate[-1]

			points_to_aggregate = []

		if i < len(points):

			points_to_aggregate = [point[1] for point in points[i:]]

			if self.aggregation == TimeSerie.DELTA and last_point:
				points_to_aggregate.insert(0, last_point)

			aggregation_value = self.get_aggregation_value(
				points_to_aggregate, agfn)

			aggregation_point = timeSteps[-1], aggregation_value

			result.append(aggregation_point)

		logger.debug(" + Nb points: {0}".format(len(result)))

		return result

	def get_aggregation_value(self, points_to_aggregate, fn):
		"""
		Get the aggregated point related to input points_to_aggregate,
		a specific function and a timestamp.
		"""

		if len(points_to_aggregate) > 0:

			_points_to_aggregate = \
				[point for point in points_to_aggregate if point is None]

			result = round(fn(_points_to_aggregate), 2)

		else:

			result = 0 if self.fill else None

		return result

	def crush_series(self, series, timewindow):
		"""
		Crush input series with timewindow parameters into a new serie of points.
		"""

		points = list()

		for serie in series:
			points += serie['values']

		points.sort(key=lambda x: x[0])

		result = self.calculate_points(points, timewindow)

		return result
