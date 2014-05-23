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

from time import time

from cmongo import MongoDB


class Store(MongoDB):
	"""
	Manage periodic and timed access to data_name in a mongo collection.
	"""

	DEFAULT_CONFIGURATION_FILE = '~/etc/store.conf'

	def __init__(self, data_name="metric", *args, **kwargs):

		super(Store, self).__init__(
			backend=self.get_collection_name(data_name=data_name),
			_ready_to_conf=False, _ready_to_connect=False, *args, **kwargs)

		self._data_name = data_name

		if self.auto_conf:
			self.apply_configuration()

	def get_parsers_by_option_by_section(self, *args, **kwargs):

		result = dict()

		mongodb_parsers = super(MongoDB, self).get_parsers_by_option_by_section()

		for section in list(mongodb_parsers.keys()):
			new_section_name = "{0}_{1}_{2}".format(
				section, type(self), self._data_name)
			result[new_section_name] = mongodb_parsers[section]

		return result

	def _get_collection_prefix(self):
		"""
		Protected method to override in order to get collection prefix.
		"""

		raise NotImplementedError()

	def get_collection_name(self, data_name=None):
		"""
		Get collection name managed by this store related to input data_name.
		"""

		if data_name is None:
			data_name = self._data_name

		result = "{0}_{1}".format(
			self._get_collection_prefix(), data_name)

		return result

	def _manage_query_error(self, result_query):
		"""
		Manage mongo query error.
		"""

		error = result_query.get("writeConcernError", None)

		if error is not None:
			self.logger.error(' error in updating document: {0}'.format(error))

		error = result_query.get("writeError")

		if error is not None:
			self.logger.error(' error in updating document: {0}'.format(error))

	def drop(self):
		"""
		Drop all data.
		"""

		self._get_backend().remove()

	def _get_data_id(self, data_id):
		"""
		Get data id if data_id is an entity _id
		"""

		result = data_id

		if data_id.startswith(self._data_name):
			result = data_id[len(self._data_name) + 1:]

		return result

from .timewindow import get_offset_timewindow


class TimedStore(Store):
	"""
	Store dedicated to manage timed data.
	"""

	DATA_ID = 'd'
	VALUE = 'v'
	TIMESTAMP = 't'

	def _get_indexes(self):

		result = super(TimedStore, self)._get_indexes()

		result.append(TimedStore._get_query_hint())

		return result

	def _get_collection_prefix(self):
		"""
		Get collection prefix.
		"""

		return "timed"

	def get(self, data_id, timewindow=None, with_meta=False, limit=0, sort=None):
		"""
		Get a sorted list of couple of dictionary
		(timestamp, dict(data_name, data_value)).
		If timewindow is None, the list contains at most one element where timestamp
		is now.
		If with_meta is True, the third field contains meta information.
		"""

		result = list()

		data_id = self._get_data_id(data_id)

		# set a where clause for the search
		where = {
					TimedStore.DATA_ID: data_id
				}

		# if timewindow is None, max timestamp is now
		timestamp = timewindow.stop() if timewindow is not None else time()
		where[TimedStore.TIMESTAMP] = {'$lte': timestamp}

		# set a projection
		projection = {
			TimedStore.VALUE: 1,
			TimedStore.TIMESTAMP: 1
		}

		# do the query
		cursor = self._get_backend().find(where, projection=projection)

		# if timewindow is None or contains only one point, get only last document
		# respectively before now or before the one point
		if limit != 0:
			cursor.limit(limit)

		# apply a specific index
		cursor.hint(TimedStore._get_query_hint())

		if sort is not None:
			_sort = [item if isinstance(item, tuple) else (item, Store.ASC)
					for item in sort]
			cursor.sort(_sort)

		# iterate on all documents
		for document in cursor:
			timestamp = document.pop(TimedStore.TIMESTAMP)
			value = document.pop(TimedStore.VALUE)

			value_to_append = (timestamp, value,)
			if with_meta:
				value_to_append += (document,)

			result.append(value_to_append)

			if timewindow is not None and timestamp not in timewindow:
				# stop when a document is just before the start timewindow
				break

		return result

	def count(self, data_id):
		"""
		Get number of timed documents for input data_id.
		"""

		data_id = self._get_data_id(data_id)

		query = {
			TimedStore.DATA_ID: data_id
		}

		cursor = self._get_backend().find(query)
		cursor.hint(TimedStore._get_query_hint())
		result = cursor.count()

		return result

	def put(self, data_id, value, timestamp, fields_to_override=dict()):
		"""
		Put a dictionary of value by name in collection.

		fields_to_override are added after the search operation.
		"""

		timewindow = get_offset_timewindow(offset=timestamp)

		data_id = self._get_data_id(data_id)

		data = self.get(
			data_id=data_id, timewindow=timewindow, with_meta=True, limit=1)

		# if we have to update previous value
		if data and data[0][1] == value:
			_id = data[0][2].pop('_id')
			if fields_to_override and data[0][2] != fields_to_override:
				# impossible to replace native fields by fields_to_override
				_set = fields_to_override.copy()
				_set.update(
					{
						TimedStore.DATA_ID: data_id,
						TimedStore.TIMESTAMP: timestamp,
						TimedStore.VALUE: value
					})
				self._get_backend().update(
					{'_id': _id},
					{
						'$set': _set
					},
					w=1)

		else:  # let's insert a document
			values_to_insert = fields_to_override.copy()
			values_to_insert.update({
					TimedStore.DATA_ID: data_id,
					TimedStore.TIMESTAMP: timestamp,
					TimedStore.VALUE: value
				})
			self._get_backend().insert(values_to_insert, w=1)

	def remove(self, data_id, timewindow=None):
		"""
		Remove timed_data existing on input timewindow.
		"""

		data_id = self._get_data_id(data_id)

		where = {
			TimedStore.DATA_ID: data_id
		}

		if timewindow is not None:
			where[TimedStore.TIMESTAMP] = \
				{'$gte': timewindow.start(), '$lte': timewindow.stop()}

		self._get_backend().remove(where, w=1)

	@staticmethod
	def _get_query_hint():
		"""
		Get query hint.
		"""

		result = [
				(TimedStore.DATA_ID, 1),
				(TimedStore.TIMESTAMP, -1)
			]

		return result

	def size(self, data_id=None):
		"""
		Get documents size for data if data_id else for the entire collection.
		"""

		raise NotImplementedError()

from md5 import new as md5
from .timewindow import Period
from operator import itemgetter
from datetime import datetime
from time import mktime


class PeriodicStore(Store):
	"""
	Store dedicated to manage periodic data.
	"""

	DATA_ID = 'i'
	TIMESTAMP = 't'
	AGGREGATION = 'a'
	VALUES = 'v'
	PERIOD = 'p'
	LAST_UPDATE = 'l'

	class PeriodicStoreError(Exception):
		pass

	def _get_indexes(self):
		"""
		Get indexes.
		"""

		result = super(PeriodicStore, self)._get_indexes()

		result.append(PeriodicStore._get_query_hint())

		return result

	def _get_collection_prefix(self):
		"""
		Get collection prefix.
		"""

		return "periodic"

	def count(
		self, data_id, aggregation, period, timewindow=None
	):
		"""
		Get number of periodic documents for input data_id.
		"""

		data = self.get(data_id=data_id, aggregation=aggregation,
			timewindow=timewindow, period=period)

		result = len(data)

		return result

	def size(self, data_id=None, aggregation=None, period=None, timewindow=None):
		"""
		Get size occupied by research filter data_id
		"""

		where = {
			PeriodicStore.DATA_ID: data_id
		}

		if aggregation is not None:
			where[PeriodicStore.AGGREGATION] = aggregation

		if timewindow is not None:
			where[PeriodicStore.TIMESTAMP] = {
				'$gte': timewindow.start(),
				'$lte': timewindow.stop()
			}

		if period is not None:
			where[PeriodicStore.PERIOD] = period

		cursor = self._get_backend().find(where)
		cursor.hint(PeriodicStore._get_query_hint())

		result = cursor.count()

		return result

	def get(self, data_id, aggregation, period, timewindow=None, limit=0):
		"""
		Get a list of points.
		"""

		query = self._get_documents_query(data_id=data_id,
			aggregation=aggregation, timewindow=timewindow, period=period)

		projection = {
			PeriodicStore.TIMESTAMP: 1,
			PeriodicStore.VALUES: 1
		}

		if aggregation is None:
			projection[PeriodicStore.AGGREGATION] = 1

		if period is None:
			projection[PeriodicStore.PERIOD] = 1

		cursor = self._get_backend().find(query, projection=projection)

		cursor.hint(PeriodicStore._get_query_hint())

		result = list()

		if limit != 0:
			cursor = cursor[:limit]

		for document in cursor:

			timestamp = int(document[PeriodicStore.TIMESTAMP])

			values = document[PeriodicStore.VALUES]

			for t, value in values.iteritems():
				value_timestamp = timestamp + int(t)

				if timewindow is None or value_timestamp in timewindow:
					result.append((value_timestamp, value))

		result.sort(key=itemgetter(0))

		return result

	def put(self, data_id, aggregation, period, points):
		"""
		Put periodic points in periodic collection with specific aggregation and
		period values.

		points is an iterable of (timestamp, value)
		"""

		# initialize a dictionary of perfdata value by value field and id_timestamp
		document_properties_by_id_timestamp = dict()
		# previous variable contains a dict of entries to put in the related document

		# prepare data to insert/update in db
		for timestamp, value in points:

			timestamp = int(timestamp)
			id_timestamp = int(period.round_timestamp(timestamp, normalize=True))
			document_properties = document_properties_by_id_timestamp.setdefault(
				id_timestamp, dict())

			if '_id' not in document_properties:
				document_properties['_id'] = PeriodicStore._get_document_id(
					data_id=data_id, aggregation=aggregation,
					timestamp=id_timestamp, period=period)
				document_properties[PeriodicStore.LAST_UPDATE] = timestamp

			else:
				if document_properties[PeriodicStore.LAST_UPDATE] < timestamp:
					document_properties[PeriodicStore.LAST_UPDATE] = timestamp

			field_name = "{0}.{1}".format(
				PeriodicStore.VALUES, timestamp - id_timestamp)

			document_properties[field_name] = value

		for id_timestamp, document_properties in \
			document_properties_by_id_timestamp.iteritems():

			# remove _id and last_update
			_id = document_properties.pop('_id')

			_set = {
				PeriodicStore.DATA_ID: data_id,
				PeriodicStore.AGGREGATION: aggregation,
				PeriodicStore.PERIOD: period.unit_values,
				PeriodicStore.TIMESTAMP: id_timestamp
			}
			_set.update(document_properties)

			document_properties['_id'] = _id

			result = self._get_backend().update(
				{
					'_id': _id
				},
				{
					'$set': _set
				},
				upsert=True,
				w=1)

			self._manage_query_error(result)

	def remove(self, data_id, aggregation=None, period=None, timewindow=None):
		"""
		Remove periodic data related to data_id, timewindow and period.
		If timewindow is None, remove all periodic_data with input period.
		If period is None
		"""

		query = self._get_documents_query(data_id=data_id,
			aggregation=aggregation, timewindow=timewindow, period=period)

		if timewindow is not None:

			projection = {
				PeriodicStore.TIMESTAMP: 1,
				PeriodicStore.VALUES: 1
			}

			documents = self._get_backend().find(query, projection=projection)

			for document in documents:
				timestamp = document.get(PeriodicStore.TIMESTAMP)
				values = document.get(PeriodicStore.VALUES)
				values_to_save = {t: v for t, v in values.iteritems()
					if (timestamp + int(t)) not in timewindow}
				_id = document.get('_id')

				if len(values_to_save) > 0:
					self._get_backend().update(
						{
							'_id': _id
						},
						{
							'$set': {PeriodicStore.VALUES: values_to_save}
						},
						w=1)
				else:
					self._get_backend().remove(_id, w=1)

		else:
			self._get_backend().remove(query)

	@staticmethod
	def _get_query_hint():
		"""
		Get query hint.
		"""

		result = [
				(PeriodicStore.DATA_ID, 1),
				(PeriodicStore.AGGREGATION, 1),
				(PeriodicStore.PERIOD, 1),
				(PeriodicStore.TIMESTAMP, 1)
			]

		return result

	def _get_documents_query(self, data_id, aggregation, timewindow, period):
		"""
		Get mongo documents query about data_id, timewindow and period.

		If period is None and timewindow is not None, period takes default period
		value for data_id.
		"""

		data_id = self._get_data_id(data_id)

		result = {
			PeriodicStore.DATA_ID: data_id
		}

		if aggregation is not None:
			result[PeriodicStore.AGGREGATION] = aggregation

		if period is not None:  # manage specific period
			result[PeriodicStore.PERIOD] = period.unit_values

		if timewindow is not None:  # manage specific timewindow
			start_timestamp, stop_timestamp = \
				PeriodicStore._get_id_timestamps(timewindow=timewindow, period=period)
			result[PeriodicStore.TIMESTAMP] = {
				'$gte': start_timestamp,
				'$lte': stop_timestamp}

		return result

	@staticmethod
	def _get_id_timestamps(timewindow, period):
		"""
		Get id timestamps related to input timewindow and period.
		"""

		# get minimal timestamp
		start_timestamp = int(
			period.round_timestamp(timewindow.start(), normalize=True))
		# and maximal timestamp
		stop_timestamp = int(
			period.round_timestamp(timewindow.stop(), normalize=True))
		stop_datetime = datetime.fromtimestamp(stop_timestamp)
		delta = period.get_delta()
		stop_datetime += delta
		stop_timestamp = mktime(stop_datetime.timetuple())

		result = start_timestamp, stop_timestamp

		return result

	@staticmethod
	def _get_document_id(data_id, aggregation, timestamp, period):
		"""
		Get periodic document id related to input data_id, timestamp aggregation and
		period.
		"""

		md5_result = md5()

		# add data_id in id
		md5_result.update(data_id.encode('ascii', 'ignore'))

		# add aggregation
		md5_result.update(aggregation.encode('ascii', 'ignore'))

		# add id_timestamp in id
		md5_result.update(str(timestamp).encode('ascii', 'ignore'))

		# add period in id
		unit_with_value = period.get_max_unit()
		if unit_with_value is None:
			raise PeriodicStore.PeriodicStoreError(
				"period {0} must contain at least one valid unit among {1}".
				format(period, Period.UNITS))

		md5_result.update(unit_with_value[Period.UNIT].encode('ascii', 'ignore'))

		# resolve md5
		result = md5_result.hexdigest()

		return result
