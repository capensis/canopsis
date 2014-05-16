#!/usr/bin/env python
# -*- coding: utf-8 -*-
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

from time import mktime, time

from md5 import new as md5

from operator import itemgetter

from datetime import datetime

from pyperfstore3.timewindow import TimeWindow, Period
from pyperfstore3.store import Store

from collections import Iterable

from sys import maxint

import logging


class Manager(object):

	class ManagerError(Exception):
		pass

	__slots__ = ['perfdata3', 'logger']

	def __init__(self, perfdata3=None, logger=None):

		super(Manager, self).__init__()

		if logger is None:
			logger = logging.getLogger(Manager.__name__)
		self.logger = logger

		if perfdata3 is None:

			store = Store(logging_level=self.logger.level)
			perfdata3 = store.collection

		self.perfdata3 = perfdata3

	@staticmethod
	def get_metric_id(component, resource, metric):

		metric_id_md5 = md5()
		# add co in id
		metric_id_md5.update(component.encode('ascii', 'ignore'))

		if resource is not None:
			# add re in id
			metric_id_md5.update(resource.encode('ascii', 'ignore'))
		# add me in id
		metric_id_md5.update(metric.encode('ascii', 'ignore'))

		result = metric_id_md5.hexdigest()

		return result

	@staticmethod
	def get_document_id(metric_id, timestamp, period):

		md5_result = md5()

		# add metric_id in id
		md5_result.update(metric_id.encode('ascii', 'ignore'))

		# add id_timestamp in id
		md5_result.update(str(timestamp).encode('ascii', 'ignore'))

		# add period in id
		unit_with_value = period.get_max_unit()
		if unit_with_value is None:
			raise Manager.ManagerError(
				"period {0} must contain at least one valid unit among {1}".
				format(period, Period.UNITS))
		md5_result.update(unit_with_value[Period.UNIT].encode('ascii', 'ignore'))

		# resolve md5
		result = md5_result.hexdigest()

		return result

	@staticmethod
	def get_document(timewindow, period):
		"""
		Get id timestamps related to input timewindow and period.
		"""

		# get minimal timestamp
		start_timestamp = int(period.sliding_timestamp(timewindow.start(), normalize=True))
		# and maximal timestamp
		stop_timestamp = int(period.sliding_timestamp(timewindow.stop(), normalize=True))
		stop_datetime = datetime.fromtimestamp(stop_timestamp)
		delta = period.get_delta()
		stop_datetime += delta
		stop_timestamp = mktime(stop_datetime.timetuple())

		result = start_timestamp, stop_timestamp

		return result

	@staticmethod
	def get_documents_query(metric_id, timewindow, period):
		"""
		Get mongo documents query about metric_id, timewindow and period.

		If period is None and timewindow is not None, period takes default period value for metric_id.
		"""

		query = {"metric_id": metric_id}

		if period is not None:  # manage specific period
			query['period'] = period.unit_values

		if timewindow is not None:  # manage specific timewindow
			if period is None:
				period = Manager.get_default_period(metric_id)  # TODO : remove all documents whatever timestamps if period is None
			start_timestamp, stop_timestamp = Manager.get_document(timewindow, period)
			query['timestamp'] = {'$gte': start_timestamp, '$lte': stop_timestamp}

		return query

	@staticmethod
	def get_default_period(metric_id):
		"""
		Get default period related to input metric_id.
		"""

		result = Period(minute=60)

		return result

	def put_data(self, metric_id, points, meta=dict(), period=None, **kwargs):
		"""
		Put a (list of) couple (timestamp, value), a meta into perfdata3 related to input period.
		kwargs will be added to all document in order to extend perfdata3 documents.
		"""

		if period is None:
			period = Manager.get_default_period(metric_id)

		# if points is a couple, transform it into a tuple of couple
		if len(points) > 0:
			if not isinstance(points[0], Iterable):
				points = (points,)

		# initialize a dictionary of perfdata value by value field and id_timestamp
		document_properties_by_id_timestamp = dict()
		# previous variable contains a dictionary of entries to put in the related document

		# prepare data to insert/update in db
		for timestamp, value in points:

			timestamp = int(timestamp)
			id_timestamp = int(period.sliding_timestamp(timestamp, normalize=True))
			document_properties = document_properties_by_id_timestamp.setdefault(
				id_timestamp, kwargs.copy())

			self.logger.debug(' + id_timestamp: {0}'.format(id_timestamp))

			if '_id' not in document_properties:
				document_properties['_id'] = Manager.get_document_id(
					metric_id, id_timestamp, period)
				document_properties['last_update'] = timestamp

			else:
				if document_properties['last_update'] < timestamp:
					document_properties['last_update'] = timestamp

			field_name = "values.{0}".format(timestamp - id_timestamp)

			document_properties[field_name] = value

		for id_timestamp, document_properties in document_properties_by_id_timestamp.iteritems():

			# remove _id and last_update
			_id = document_properties.pop('_id')

			_set = {
				'period': period.unit_values,
				'metric_id': metric_id,
				'timestamp': id_timestamp,
				'meta': meta  # TODO : add meta in the dedicated collection perfdata3_meta
			}
			_set.update(document_properties)

			document_properties['_id'] = _id

			result = self.perfdata3.update(
				{
					'_id': _id
				},
				{
					'$set': _set
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

	def get_data(self, metric_id, timewindow=None, period=None, return_meta=False):
		"""
		Get a set of data related to input metric_id on the timewindow and input period.
		If return_meta, result is the couple (set of data, meta)
		"""

		if period is None:
			period = Manager.get_default_period(metric_id)

		query = Manager.get_documents_query(
			metric_id=metric_id, timewindow=timewindow, period=period)

		projection = {
			'timestamp': 1,
			'values': 1
		}
		if return_meta:
			projection['meta'] = 1

		cursor = self.perfdata3.find(query,	projection=projection)

		cursor.hint([('metric_id', 1), ('period', 1), ('timestamp', 1)])

		meta = (0, None)

		result = list()

		for document in cursor:

			timestamp = int(document.get('timestamp'))

			# get last meta information
			if return_meta and timestamp > meta[0]:
					meta = (timestamp, document.get('meta'))

			values = document.get('values')

			for t, value in values.iteritems():
				value_timestamp = timestamp + int(t)

				if timewindow is None or value_timestamp in timewindow:
					result.append( (value_timestamp, value) )

		result.sort(key=itemgetter(0))

		result = tuple(result)

		if return_meta:
			result = (result, meta[1])

		return result

	def remove(self, metric_id, timewindow=None, period=None):
		"""
		Remove values and meta of one metric.
		"""

		if period is None:
			period = Manager.get_default_period(metric_id)

		query = Manager.get_documents_query(metric_id, timewindow, period)

		# TODO : only remove values where timestamps are in document which contain more values than in timewindow
		# for example, remove values in the set ([start_timestamp; start_timestamp + period.value] - [start_timestamp; start])

		self.perfdata3.remove(query)

	def update_meta(self, metric_id, meta, timewindow=None, period=None):
		"""
		"""

		raise NotImplemented()
