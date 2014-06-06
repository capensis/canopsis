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

from cmongo.storage import Storage

from cstorage.periodic import PeriodicStore

from md5 import new as md5
from ctimeserie.timewindow import Period
from operator import itemgetter
from datetime import datetime
from time import mktime


class PeriodicStore(Storage, PeriodicStore):
	"""
	Storage dedicated to manage periodic data.
	"""

	class Index:

		DATA_ID = 'i'
		TIMESTAMP = 't'
		AGGREGATION = 'a'
		VALUES = 'v'
		PERIOD = 'p'
		LAST_UPDATE = 'l'

		QUERY = [(DATA_ID, 1), (AGGREGATION, 1), (PERIOD, 1), (TIMESTAMP, 1)]

	def _get_indexes(self):
		"""
		Get indexes.
		"""

		result = super(PeriodicStore, self)._get_indexes()

		result.append(PeriodicStore.Index.QUERY)

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
			PeriodicStore.Index.DATA_ID: data_id
		}

		if aggregation is not None:
			where[PeriodicStore.Index.AGGREGATION] = aggregation

		if timewindow is not None:
			where[PeriodicStore.Index.TIMESTAMP] = {
				'$gte': timewindow.start(),
				'$lte': timewindow.stop()
			}

		if period is not None:
			where[PeriodicStore.Index.PERIOD] = period

		cursor = self._find(document=where)
		cursor.hint(PeriodicStore.Index.QUERY)

		result = cursor.count()

		return result

	def get(self, data_id, aggregation, period, timewindow=None, limit=0):
		"""
		Get a list of points.
		"""

		query = self._get_documents_query(data_id=data_id,
			aggregation=aggregation, timewindow=timewindow, period=period)

		projection = {
			PeriodicStore.Index.TIMESTAMP: 1,
			PeriodicStore.Index.VALUES: 1
		}

		if aggregation is None:
			projection[PeriodicStore.Index.AGGREGATION] = 1

		if period is None:
			projection[PeriodicStore.Index.PERIOD] = 1

		cursor = self._find(document=query, projection=projection)

		cursor.hint(PeriodicStore.Index.QUERY)

		result = list()

		if limit != 0:
			cursor = cursor[:limit]

		for document in cursor:

			timestamp = int(document[PeriodicStore.Index.TIMESTAMP])

			values = document[PeriodicStore.Index.VALUES]

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
				document_properties[PeriodicStore.Index.LAST_UPDATE] = timestamp

			else:
				if document_properties[PeriodicStore.Index.LAST_UPDATE] < timestamp:
					document_properties[PeriodicStore.Index.LAST_UPDATE] = timestamp

			field_name = "{0}.{1}".format(
				PeriodicStore.Index.VALUES, timestamp - id_timestamp)

			document_properties[field_name] = value

		for id_timestamp, document_properties in \
			document_properties_by_id_timestamp.iteritems():

			# remove _id and last_update
			_id = document_properties.pop('_id')

			_set = {
				PeriodicStore.Index.DATA_ID: data_id,
				PeriodicStore.Index.AGGREGATION: aggregation,
				PeriodicStore.Index.PERIOD: period.unit_values,
				PeriodicStore.Index.TIMESTAMP: id_timestamp
			}
			_set.update(document_properties)

			document_properties['_id'] = _id

			result = self._update(_id={'_id': _id}, document={'$set': _set})

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
				PeriodicStore.Index.TIMESTAMP: 1,
				PeriodicStore.Index.VALUES: 1
			}

			documents = self._find(document=query, projection=projection)

			for document in documents:
				timestamp = document.get(PeriodicStore.Index.TIMESTAMP)
				values = document.get(PeriodicStore.Index.VALUES)
				values_to_save = {t: v for t, v in values.iteritems()
					if (timestamp + int(t)) not in timewindow}
				_id = document.get('_id')

				if len(values_to_save) > 0:
					self._update(
						_id={'_id': _id},
						document={
							'$set': {PeriodicStore.Index.VALUES: values_to_save}
						})
				else:
					self._remove(document=_id)

		else:
			self._remove(document=query)

	def _get_documents_query(self, data_id, aggregation, timewindow, period):
		"""
		Get mongo documents query about data_id, timewindow and period.

		If period is None and timewindow is not None, period takes default period
		value for data_id.
		"""

		data_id = self._get_data_id(data_id)

		result = {
			PeriodicStore.Index.DATA_ID: data_id
		}

		if aggregation is not None:
			result[PeriodicStore.Index.AGGREGATION] = aggregation

		if period is not None:  # manage specific period
			result[PeriodicStore.Index.PERIOD] = period.unit_values

		if timewindow is not None:  # manage specific timewindow
			start_timestamp, stop_timestamp = \
				PeriodicStore._get_id_timestamps(
					timewindow=timewindow, period=period)
			result[PeriodicStore.Index.TIMESTAMP] = {
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
