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


class TimeSerie(object):
	"""
	Time serie management. Contain a TimeWindow, a period, a sliding_time,
	an operation and fill property.
	"""

	__slots__ = ('period', 'sliding_time', 'operation', 'fill')

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
		sliding_time=True,
		operation=MEAN,
		fill=False):

		self.period = self._get_period(max_points=max_points, period=period)
		self.sliding_time = sliding_time
		self.operation = operation
		self.fill = fill

	def __repr__(self):

		message = "period: {0}, sliding_time: {1}" + \
			", operation: {2}, fill: {3}"
		result = message.format(
			self.period, self.sliding_time, self.operation, self.fill)

		return result

	def _get_period(self, max_points, period):

		result = Period(value=period.value, unit=period.unit) \
			if period is not None else None

		if result is None and max_points:
			interval = self.timewindow.stop - self.timewindow.start
			# may reduce value with timewindow.exclusion_intervals
			seconds = interval.total_seconds() / max_points
			result = Period(value=seconds)

		return result

	def getTimeSteps(self, timewindow):

		logger.debug('getTimeSteps:')

		timeSteps = []

		logger.debug(' + TimeSerie: {0}'.format(self))

		start_datetime = datetime.utcfromtimestamp(timewindow.start)
		stop_datetime = datetime.utcfromtimestamp(timewindow.stop)

		if self.sliding_time:
			stop_datetime = self.period.sliding_time(
				stop_datetime,
				self.timezone)

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
		Get self operation.
		"""

		if self.operation == TimeSerie.MEAN or not self.operation:
			result = lambda x: sum(x) / float(len(x))
		elif self.operation == TimeSerie.MIN:
			result = lambda x: min(x)
		elif self.operation == TimeSerie.MAX:
			result = lambda x: max(x)
		elif self.operation == TimeSerie.FIRST:
			result = lambda x: x[0]
		elif self.operation == TimeSerie.LAST:
			result = lambda x: x[-1]
		elif self.operation == TimeSerie.DELTA:
			result = lambda x: (max(x) - min(x)) / 2.

		return result

	def calculate_points(
		self,
		points,
		timewindow):
		"""
		Calculate an array of points.
		"""

		result = []

		logger.debug("Calculate {0} points (timeserie: {1}, timewindow: {2})".format(
			len(points),
			self,
			timewindow))

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

			if self.operation == TimeSerie.DELTA and last_point:
				points_to_aggregate.insert(0, last_point)

			aggregation_point = self.get_aggregation_point(
				points_to_aggregate,
				agfn,
				previous_timestamp)

			result.append(aggregation_point)

			if points_to_aggregate:
				last_point = points_to_aggregate[-1]

			points_to_aggregate = []

		if i < len(points):

			points_to_aggregate = [point[1] for point in points[i:]]

			if self.operation == TimeSerie.DELTA and last_point:
				points_to_aggregate.insert(0, last_point)

			aggregation_point = self.get_aggregation_point(
				points_to_aggregate,
				agfn,
				timeSteps[-1])

			result.append(aggregation_point)

		logger.debug(" + Nb points: {0}".format(len(result)))

		return result

	def get_aggregation_point(self, points_to_aggregate, fn, timestamp):

		if points_to_aggregate:

			logger.debug(
				" + {0} points to aggregate".format(
					(len(points_to_aggregate))))

			_points_to_aggregate = \
				[point for point in points_to_aggregate if point is None]

			agvalue = round(fn(_points_to_aggregate), 2)

			result = [timestamp, agvalue]

		else:
			logger.debug(" + No points")

			result = [timestamp, 0 if self.fill else None]

		logger.debug(" + Point : {0} ".format(result))

		return result

	def crush_series(self, series, timewindow):

		points = list()

		for serie in series:
			points += serie['values']

		points.sort(key=lambda x: x[0])

		result = self.calculate_points(points, timewindow)

		return result
