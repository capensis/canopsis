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

from cstorage.timed import TimedStorage

from ctimeserie.timewindow import get_offset_timewindow


class TimedStorage(Storage, TimedStorage):
	"""
	Storage dedicated to manage timed data.
	"""

	class Index:

		DATA_ID = 'd'
		VALUE = 'v'
		TIMESTAMP = 't'

		TIMESTAMP_INDEX = 0
		VALUE_INDEX = 1

		TIMESTAMP_BY_ID = [(DATA_ID, 1), (TIMESTAMP, -1)]

	def get(self, data_id, timewindow=None, limit=0, sort=None):
		"""
		Get a sorted list of triplet of dictionaries such as :
		tuple(timestamp, dict(data_type, data_value), dict(meta_name, meta_value)).

		If timewindow is None, result is all timed document.
		"""

		result = list()

		data_id = self._get_data_id(data_id)

		# set a where clause for the search
		where = {
					TimedStorage.Index.DATA_ID: data_id
				}

		# if timewindow is not None, get latest timestamp before timewindow.stop()
		if timewindow is not None:
			timestamp = timewindow.stop()
			where[TimedStorage.Index.TIMESTAMP] = {'$lte': timewindow.stop()}

		# do the query
		cursor = self._find(document=where)

		# if timewindow is None or contains only one point, get only last document
		# respectively before now or before the one point
		if limit != 0:
			cursor.limit(limit)

		# apply a specific index
		cursor.hint(TimedStorage.Index.TIMESTAMP_BY_ID)

		if sort is not None:
			_sort = [item if isinstance(item, tuple) else (item, Storage.ASC)
					for item in sort]
			cursor.sort(_sort)

		# iterate on all documents
		for document in cursor:
			timestamp = document.pop(TimedStorage.Index.TIMESTAMP)
			value = document.pop(TimedStorage.Index.VALUE)

			value_to_append = (timestamp, value, document)

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
			TimedStorage.Index.DATA_ID: data_id
		}

		cursor = self._find(document=query)
		cursor.hint(TimedStorage.Index.TIMESTAMP_BY_ID)
		result = cursor.count()

		return result

	def put(self, data_id, value, timestamp):
		"""
		Put a dictionary of value by name in collection.
		"""

		timewindow = get_offset_timewindow(offset=timestamp)

		data_id = self._get_data_id(data_id)

		data = self.get(
			data_id=data_id, timewindow=timewindow, limit=1)

		data_value = None

		if data:
			data_value = data[0][TimedStorage.VALUE_INDEX]

		if value != data_value:  # new entry to insert

			values_to_insert = {
					TimedStorage.Index.DATA_ID: data_id,
					TimedStorage.Index.TIMESTAMP: timestamp,
					TimedStorage.Index.VALUE: value
				}
			self._insert(document=values_to_insert)

	def remove(self, data_id, timewindow=None):
		"""
		Remove timed_data existing on input timewindow.
		"""

		data_id = self._get_data_id(data_id)

		where = {
			TimedStorage.Index.DATA_ID: data_id
		}

		if timewindow is not None:
			where[TimedStorage.Index.TIMESTAMP] = \
				{'$gte': timewindow.start(), '$lte': timewindow.stop()}

		self._remove(document=where)

	def size(self, data_id=None):
		"""
		Get documents size for data if data_id else for the entire collection.
		"""

		raise NotImplementedError()

	def _get_indexes(self):

		result = super(TimedStorage, self)._get_indexes()

		result.append(TimedStorage.Index.TIMESTAMP_BY_ID)

		return result

	def _get_collection_prefix(self):
		"""
		Get collection prefix.
		"""

		return "timed"
