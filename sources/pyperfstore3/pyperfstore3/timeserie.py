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

# provide only TimeSerie
__all__ = ('TimeSerie')

import logging

logger = logging.getLogger('utils')
logger.setLevel(logging.DEBUG)

from pyperfstore3.timewindow import Period

from .aggregation import get_aggregations
from operator import itemgetter
from time import mktime


class TimeSerie(object):
	"""
	Time serie management. Contain a period and operation of aggregation,
	and a round time and a fill boolean properties.

	- period: interval of steps of aggregated points.
	- aggregation: aggregation operation name.
	- round_time: round_time of input timewindow during calculations.
	- fill: change None values by 0.
	"""

	__slots__ = ('period', 'max_points', 'round_time', 'aggregation', 'fill')

	MAX_POINTS = 500

	def __init__(
		self,
		max_points=MAX_POINTS,
		period=None,
		round_time=True,
		aggregation='MEAN',
		fill=False):

		self.period = period
		self.max_points = max_points
		self.round_time = round_time
		self.aggregation = aggregation
		self.fill = fill

	def __repr__(self):

		message = "TimeSerie(period: {0}, round_time: {1}" + \
			", aggregation: {2}, fill: {3})"
		result = message.format(
			self.period, self.round_time, self.aggregation, self.fill)

		return result

	def _get_period(self, timewindow):
		"""
		Get a period related to input max_points or a period.
		"""

		result = self.period

		if result is None:
			seconds = (timewindow.stop() - timewindow.start()) / self.max_points
			result = Period(second=seconds)

		return result

	def timesteps(self, timewindow):
		"""
		Get a list of same longer intervals inside timewindow.
		The upper bound is timewindow.stop_datetime()
		"""

		# get the right period to apply on timewindow
		period = self._get_period(timewindow=timewindow)

		result = []

		# set start and stop datetime
		start_datetime = timewindow.start_datetime()
		stop_datetime = timewindow.stop_datetime()

		if self.round_time: # normalize if round time is True
			start_datetime = period.round_datetime(
				datetime=start_datetime, normalize=True)

		current_datetime = start_datetime
		delta = period.get_delta()

		while current_datetime < stop_datetime:
			ts = mktime(current_datetime.timetuple())
			result.append(ts)
			current_datetime += delta

		result.append(timewindow.stop())

		return result

	def calculate(self, points, timewindow, fn=None):
		"""
		Do an operation on all points with input timewindow.

		If fn is None, then self aggregation operation is used.

		Return points su as follow:
		Let points of the form: [(T0, V0), ..., (Tn, Vn)]
		and
		[(T0, fn(V0, V1)), (T2, fn(V2, V3), ...]
		"""

		result = list()

		# first, get the right function
		if fn is None:
			fn = get_aggregations()[self.aggregation]

		if len(points) == 1:  # if 1 point, return it
			result = [(timewindow.start(), points[0][1])]

		elif len(points) > 0:
			timesteps = self.timesteps(timewindow)

			#initialize variables for loop
			i = 0
			values_to_aggregate = []
			last_point = None

			# iterate on all timesteps in order to get points
			# between [previous_timestamp, timestamp[
			for index in range(1, len(timesteps)):

				timestamp = timesteps[index]

				previous_timestamp = timesteps[index-1]

				# fill the values_to_aggregate array
				while i < len(points) and points[i][0] < timestamp:

					if points[i][1] is not None:
						values_to_aggregate.append(points[i][1])

					i += 1

				if self.aggregation == "DELTA" and last_point:
					values_to_aggregate.insert(0, last_point)

				_aggregation_value = self._aggregation_value(
					values_to_aggregate, fn)

				aggregation_point = previous_timestamp, _aggregation_value

				result.append(aggregation_point)

				if len(values_to_aggregate) > 0:
					last_point = values_to_aggregate[-1]

				values_to_aggregate = []
			"""
			if i < len(points):

				values_to_aggregate = [point[1] for point in points[i:] if point[1] is not None]

				if self.aggregation == "DELTA" and last_point:
					values_to_aggregate.insert(0, last_point)

				_aggregation_value = self._aggregation_value(
					values_to_aggregate, fn)

				aggregation_point = timesteps[-1], _aggregation_value

				result.append(aggregation_point)
			"""

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
