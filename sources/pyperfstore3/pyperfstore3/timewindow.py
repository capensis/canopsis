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

import time
from datetime import datetime, timedelta
import calendar


class Period(object):
	"""
	Period management with a value and an unitself.
	"""

	__slots__ = ['value', 'unit']

	SECOND = 'second'
	MINUTE = 'minute'
	HOUR = 'hour'
	DAY = 'day'
	WEEK = 'week'
	MONTH = 'month'
	YEAR = 'year'

	UNITS = [SECOND, MINUTE, HOUR, DAY, WEEK, MONTH, YEAR]

	def __init__(self, value=1, unit=SECOND):

		super(Period, self).__init__()

		self.value = value
		self.unit = unit

	def to_dict(self):

		return {'unit': self.unit, 'value': self.value}

	def next_period(self, step=1):
		"""
		Get next period with input step or none if next period can't be associated to a specific unit.
		"""
		result = None

		try:
			index = Period.UNITS.index(self.unit)
			next_unit = Period.UNITS[index+step]
			result = Period(unit=next_unit)  # TODO: set right value

		except ValueError:
			pass
		except IndexError:
			pass

		return result

	def sliding_timestamp(self, timestamp, timezone=0):
		"""
		Get round timestamp relative to an input timestamp and a timezone.
		"""
		utcdatetime = datetime.utcfromtimestamp(timestamp)

		utcdatetime = self.sliding_datetime(utcdate=utcdatetime, timezone=timezone)

		result = int(time.mktime(utcdatetime.timetuple()))

		return result

	def sliding_datetime(self, utcdate, timezone=0):
		"""
		Calculate roudtime relative to an UTC date and a timezone.
		"""

		result = utcdate

		dt = timedelta(seconds=timezone)
		result -= dt

		# assume result does not contain seconds and microseconds in this case
		result = result.replace(second=0, microsecond=0)

		if self.unit == Period.SECOND:
			seconds = (result.second * 60 / self.value) * self.value / 60
			result = result.replace(second=seconds)

		elif self.unit == Period.MINUTE:
			minutes = (result.minute * 60 / self.value) * self.value / 60
			result = result.replace(minute=minutes)

		else:
			result = result.replace(minute=0)

			index_unit = Period.UNITS.index(self.unit)

			if index_unit >= Period.UNITS.index(Period.DAY):
				result = result.replace(hour=0)

				if index_unit >= Period.UNITS.index(Period.WEEK):

					weeks = calendar.monthcalendar(result.year, result.month)
					for week in weeks:
						if result.day in week:
							result = result.replace(
								day=week[0] if week[0] != 0 else 1)
						break

				if index_unit >= Period.UNITS.index(Period.MONTH):
					result = result.replace(day=1)

					if index_unit >= Period.UNITS.index(Period.YEAR):
						result = result.replace(month=1)

			result += dt

		return result


class TimeWindow(object):
	"""
	Time window management with exclusion intervals.
	Contains a start/stop dates and an array of exclusion intervals
	(couple of START, STOP timestamp).
	"""

	__slots__ = ['start', 'stop', 'exclusion_intervals', 'timezone']

	def __init__(self, start, stop=time.time(), exclusion_intervals=[], timezone=0):

		self.start = start if start else stop - 60 * 60 * 24
		self.stop = stop
		self.exclusion_intervals = exclusion_intervals
		self._get_exclusion_intervals(exclusion_intervals)
		self.timezone = timezone

	def __repr__(self):

		message = "start = {0}, stop = {1}, exclusion_intervals = {2}"
		result = message.format(self.start, self.stop, self.exclusion_intervals)

		return result

	def _get_exclusion_intervals(self, exclusion_intervals):

		# should sort and simplify exclusion intervals
		return exclusion_intervals

	def total_seconds(self):
		"""
		Returns seconds inside this timewindow.
		"""

		delta = self.stop - self.start
		# remove exclusion_intervals
		result = delta.total_seconds()

		return result

	def get_next_date(self, date, period, delta=None):
		"""
		Get next date of input date with timewindow parameters,
		period and optionaly a previous calculated delta.
		"""

		delta = period.get_delta(date, delta)
		# check if next date is in exclusion dates of the input timewindow
		result = date + delta

		return result

	def get_previous_date(self, date, delta=None):
		"""
		Get previous date of input date with timewindow parameters,
		period and optionaly a previous calculated delta.
		"""

		delta = self.period.get_delta(date, delta)
		# check if next date is in exclusion dates of the input timewindow
		result = date - delta

		return result

	@staticmethod
	def get_datetime(timestamp, timezone=0):
		"""
		Get the date time corresponding to both input timestamp and timezone.
		"""

		dt = timedelta(seconds=timezone)
		result = datetime.utcfromtimestamp(timestamp) - dt

		return result
