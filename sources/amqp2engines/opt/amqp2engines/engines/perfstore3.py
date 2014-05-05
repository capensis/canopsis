#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

from ctools import parse_perfdata
from datetime import timedelta

from cstorage import get_storage

from cengine import cengine

from ctools import internal_metrics

import calendar

from datetime import datetime
import time

from md5 import new as md5

INTERNAL_QUEUE = "beat_perfdata3"


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
			result = Period(unit=Period.UNITS[index+step])  # TODO: set right value

		except ValueError:
			pass
		except IndexError:
			pass

		return result

	def sliding_timestamp(self, timestamp, timezone=0):

		datetime.fromutctimestamp(timestamp)

		utcdatetime = datetime.utcfromtimestamp(timestamp)

		utcdatetime = self.sliding_datetime(utcdate=utcdatetime, timezone=timezone)

		result = int(time.mktime(utcdatetime.timetuple()))

		return result


	def sliding_datetime(self, utcdate, timezone=0):
		"""
		Calculate roudtime relative to an UTC date, a period time/type
		and a timezone.
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

	__slots__ = ['start', 'stop', 'exclusion_intervals', '_delta']

	def __init__(self, start, stop=time.time(), exclusion_intervals=[], timezone=0):
		self.start = start if start else stop - 60 * 60 * 24
		self.stop = stop
		self.exclusion_intervals = exclusion_intervals
		self._get_exclusion_intervals(exclusion_intervals)
		self._delta = None

	def __repr__(self):
		message = "start = %s, stop = %s, exclusion_intervals = %s"
		result = message % (self.start, self.stop, self.exclusion_intervals)
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


class engine(cengine):

	etype = 'perfstore3'

	def __init__(self, *args, **kargs):

		super(engine, self).__init__(*args, **kargs)

	"""
	def create_amqp_queue(self):

		super(engine, self).create_amqp_queue()
	"""

	def pre_run(self):

		storage = get_storage(logging_level=self.logger.level)
		perfdata3 = storage.get_backend('perfdata3')

		self.manager = Manager(perfdata3, self.logger)

		"""
		self.internal_amqp = camqp(logging_level=logging.INFO, logging_name="%s-internal-amqp" % self.name)
		self.internal_amqp.add_queue(
			queue_name=INTERNAL_QUEUE,
			routing_keys=["#"],
			callback=self.on_internal_event,
			no_ack=True,
			exclusive=False,
			auto_delete=False
		)

		self.internal_amqp.start()
		"""

	"""
	def post_run(self):

		self.internal_amqp.cancel_queues()
		self.internal_amqp.stop()
		self.internal_amqp.join()
	"""

	def work(self, event, *args, **kargs):

		## Get perfdata
		perf_data = event.get('perf_data', None)
		perf_data_array = event.get('perf_data_array', list())

		if perf_data_array is None:
			perf_data_array = list()

		### Parse perfdata
		if perf_data:

			self.logger.debug(' + perf_data: {0}'.format(perf_data))

			try:
				perf_data_array = parse_perfdata(perf_data)

			except Exception as err:
				self.logger.error("Impossible to parse: {0} ('{1}')".format(perf_data, err))

		self.logger.debug(' + perf_data_array: {0}'.format(perf_data_array))

		### Add status informations
		event_type = event.get('event_type', None)

		if event_type is not None and event_type in ['check', 'selector', 'sla']:

			self.logger.debug('Add status informations')

			state = int(event.get('state', 0))
			state_type = int(event.get('state_type', 0))
			state_extra = 0

			# Multiplex state
			cps_state = state * 100 + state_type * 10 + state_extra

			perf_data_array.append(
				{
					"metric": "cps_state",
					"value": cps_state
				})

		event['perf_data_array'] = perf_data_array

		# remove perf_data_akeys where values are None
		for index, perf_data in enumerate(event['perf_data_array']):

			event['perf_data_array'][index] = \
				dict(((key, value) for key, value in perf_data.iteritems() if value is not None))

		self.logger.debug('perf_data_array: {0}'.format(perf_data_array))

		#self.internal_amqp.publish(event, INTERNAL_QUEUE)
		self.on_internal_event(event)

		return event

	def on_internal_event(self, event, msg=None):
		## Metrology
		timestamp = event.get('timestamp', None)

		if timestamp is not None:

			component = event.get('component', None)

			if component is not None:

				resource = event.get('resource', None)

				perf_data_array = event.get('perf_data_array')

				for perf_data in perf_data_array:

					name = perf_data.get('metric', None)

					if name is not None:

						# get node id
						md5_id = md5()

						# add component in id
						md5_id.update(component.encode('ascii', 'ignore'))
						# add resource in id
						if resource is not None:
							md5_id.update(resource.encode('ascii', 'ignore'))
						# add metric name in id
						md5_id.update(name)

						# get metric id
						metric_id = md5_id.hexdigest()

						value = perf_data.pop('value', None)

						self.manager.put_data(metric_id, timestamp, value, perf_data)

					else:
						self.logger.warning('metric name does not exist: {0}'.format(event))

	def beat(self):
		# Counts metric not in internal metrics for webserver cache purposes
		self.logger.info('Computing cache value for perfdata3 metric count')

		metrics_cursor = self.entities.find(
			{
				'type': 'metric',
				'name': {'$nin': internal_metrics}
			})
		count = len(metrics_cursor)

		self.object.update(
			{'crecord_name': 'perfdata2_count_no_internal'},
			{'$set':
				{'count': count}
			},
			upsert=True,
			w=1
		)
		self.logger.info('Cache value for perfdata3 metric count computed > {0}'.format(count))


class Manager(object):

	__slots__ = ['perfdata3', 'logger']

	def __init__(self, perfdata3, logger):

		super(Manager).__init__()
		self.perfdata3 = perfdata3
		self.logger = logger

	@staticmethod
	def get_document_id(metric_id, timestamp, period):

		md5_result = md5()

		# add id_timestamp in id
		md5_result.update(str(timestamp).encode('ascii', 'ignore'))
		# add period in id
		md5_result.update(period.unit.encode('ascii', 'ignore'))

		result = metric_id.join(md5_result.hexdigest())

		return result

	def put_data(self, metric_id, timestamp, value, meta=dict(), period=Period(unit=Period.MINUTE, value=60)):

		sliding_period = period.next_period()

		timestamp = int(timestamp)

		id_timestamp = sliding_period.sliding_timestamp(timestamp)

		self.logger.debug(' + id_timestamp: {0}'.format(id_timestamp))

		_id = self.get_document_id(metric_id, id_timestamp, period)

		field_name = "values.{0}".format(timestamp - id_timestamp)

		result = self.perfdata3.update(
			{
				'_id': _id,
			},
			{
				'$set': {
					'period': period.to_dict(),
					'metric_id': metric_id,
					'timestamp': id_timestamp,
					'last_update': timestamp,
					field_name: value,
					'meta': meta
				}
			},
			upsert=True,
			w=1)

		error = result.get("writeConcernError", None)

		if error is not None:
			self.logger.error(' error in updating document: {0}'.format(error))

		error = result.get("writeError")

		if error is not None:
			self.logger.error(' error in updating document: {0}'.format(error))

		self.logger.debug(' + metric updated: {0}'.format(meta))

	def get_data(metric_id, interval, period=Period(unit=Period.MINUTE, value=60)):

		timewindow = TimeWindow(**interval)

		sliding_period = period.next_period()

		start = sliding_period.sliding_timestamp(interval.start)

		stop = sliding_period.sliding_timestamp(interval.stop)

		dt = datetime.fromutctimestamp(start)

		dtstop = datetime.fromutctimestamp(stop)

		ids = list()

		start_id = Manager.get_document_id(metric_id, start, period)
		ids.append(start_id)

		stop_id = Manager.get_document_id(metric_id, stop, period)
		ids.append(stop_id)

		while dt <= dtstop:

			dt = timewindow.get_next_date(dt, period)

			t = time.mktime(dt)

			t_id = Manager.get_document_id(metric_id, t, period)

			ids.append(t_id)

		request = {"_id": {'$in': ids}}

		result = self.perfdata3.find(
			request,
			projection={'timestamp': 1, 'values': 1, 'meta': 1})

		return result
