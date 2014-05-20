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

import os
import sys
import json
import logging
from time import time

from bson.errors import InvalidStringData
from pymongo import Connection


class Store(object):
	"""
	Manage periodic and timed access to data_name in a mongo collection.
	"""

	class StoreError(Exception):
		def __init__(self):
			super(Store.StoreError, self).__init__(
				"store can not connect")

	def __init__(self,
			data_name="metric",
			mongo_host="127.0.0.1",
			mongo_port=27017,
			mongo_db='canopsis',
			mongo_user=None,
			mongo_pass=None,
			mongo_safe=False,
			logging_level=logging.INFO):

		super(Store, self).__init__()

		self.logger = logging.getLogger(__name__)
		self.logger.setLevel(logging.DEBUG)#logging_level)

		# keep dayly track of ids in cache,
		# TODO must be cleared in beat method
		self.daily_ids = {}

		self.logger.debug(" + Init MongoDB Store")

		# Read db option from conf
		import ConfigParser
		config = ConfigParser.RawConfigParser()
		config.read(os.path.expanduser('~/etc/cstorage.conf'))

		try:
			host = config.get('master', 'host')
			port = config.getint('master', 'port')
			user = config.get('master', 'userid')
			passwd = config.get('master', 'password')

			if host:
				mongo_host = host

			if port:
				mongo_port = port

			if user:
				mongo_user = user

			if passwd:
				mongo_pass = passwd

		except (ConfigParser.NoOptionError, ConfigParser.NoSectionError) as err:
			self.logger.error('Impossible to parse cstorage.conf: {0}'.format(str(err)))

		self.mongo_host = mongo_host
		self.mongo_port = mongo_port
		self.mongo_db = mongo_db
		self.data_name = data_name
		self.mongo_safe = mongo_safe
		self.mongo_user = mongo_user if mongo_user != "" else None
		self.mongo_pass = mongo_pass if mongo_pass != "" else None

		self.connected = False

		if not self.connect():
			raise Store.StoreError()

		indexes = self._get_indexes()
		for index in indexes:
			self.collection.ensure_index(index)

		self.last_rate_time = time()
		self.rate_interval = 10
		self.rate_threshold = 20
		self.last_rate = 0
		self.pushed_values = 0

	def connect(self):
		if self.connected:
			self.logger.debug("Impossible to connect, already connected")
			return True
		else:
			self.logger.debug("Connect to MongoDB ({0}/{1}@{2}:{3})".
				format(
					self.mongo_db, self.get_collection_name(), self.mongo_host, self.mongo_port))

			try:
				self.conn = Connection(
					host=self.mongo_host, port=self.mongo_port, safe=self.mongo_safe)
				self.logger.debug(" + Success")
			except Exception, err:
				self.logger.error(" + %s" % err)
				return False

			self.db = self.conn[self.mongo_db]

			try:
				if self.mongo_user and self.mongo_pass is not None:
						self.logger.debug("Try to auth '{0}'".format(self.mongo_user))
						if not self.db.authenticate(self.mongo_user, self.mongo_pass):
							raise Exception('Invalid user or pass.')
						self.logger.debug(" + Success")
			except Exception, err:
				self.logger.error(" + Impossible to authenticate: {0}".format(err))
				self.disconnect()
				return False

			self.logger.debug("Get collections")

			collection_name = self.get_collection_name()

			self.collection = self.db[collection_name]

			self.connected = True
			self.logger.debug(" + Success")
			return True

	def disconnect(self):
		"""
		Diconnect from mongo data base.
		"""

		if self.connected:
			self.logger.debug("Disconnect from MongoDB")
			self.conn.disconnect()
			self.connected = False
		else:
			self.logger.warning("Impossible to disconnect, you are not connected")

	def check_connection(self):
		if not self.connected:
			if not self.connect():
				raise Exception('Impossible to deal with DB, you are not connected ...')

	def size(self):
		"""
		Get size of collection. None if a connection error occured.
		"""

		self.logger.info("Size of dbs:")
		size = 0
		try:
			size = self.db.command("collstats", self.collection)['size']
		except Exception as e:
			self.logger.warning("Impossible to read Collecion Size: {0}".format(e))
			size = None
			print e
		else:
			self.logger.info(" + Collection: %0.2f MB" % (size/1024.0/1024.0))

		return size

	def _get_collection_prefix(self):
		"""
		Protected method to override in order to get collection prefix.
		"""

		raise NotImplementedError()

	def get_collection_name(self):
		"""
		Get collection name managed by this store related to input data_name.
		"""

		result = "{0}_{1}".format(self._get_collection_prefix(), self.data_name)

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

		self.check_connection()
		self.collection.remove()

	def _get_indexes(self):
		"""
		Get collection indexes. Must be overriden.
		"""

		raise NotImplementedError()

from .timewindow import TimeWindow


class TimedStore(Store):
	"""
	Store dedicated to manage timed data.
	"""

	DATA_ID = 'd'
	VALUE = 'v'
	TIMESTAMP = 't'

	def _get_indexes(self):

		result = [
			'_id',
			TimedStore._get_query_hint()
		]

		return result

	def _get_collection_prefix(self):
		"""
		Get collection prefix.
		"""

		return "timed"

	def get(self, data_id, timewindow=None, with_id=False):
		"""
		Get a sorted list of couple of (timestamp, dict(data_name, data_value)).
		If timewindow is None, the list contains at most one element where timestamp is now.
		"""

		result = list()

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
		cursor = self.collection.find(where, projection=projection)

		# if timewindow is None, get only the last document before now
		if timewindow is None:
			cursor.limit(1)

		# apply a specific index
		cursor.hint(TimedStore._get_query_hint())

		# iterate on all documents
		for document in cursor:
			timestamp = document[TimedStore.TIMESTAMP]
			value = document[TimedStore.VALUE]

			value_to_append = (timestamp, value,)
			if with_id:
				value_to_append += (document['_id'],)

			result.append(value_to_append)

			if timewindow is not None and timestamp not in timewindow:
				# stop when a document is just before the start timewindow
				break

		return result

	def count(self, data_id):
		"""
		Get number of timed documents for input data_id.
		"""
		self.check_connection()
		query = {
			TimedStore.DATA_ID: data_id
		}

		cursor = self.collection.find(query)
		cursor.hint(TimedStore._get_query_hint())
		result = cursor.count()

		return result

	def put(self, data_id, value, timestamp):
		"""
		Put a dictionary of value by name in collection.
		"""

		timewindow = TimeWindow(interval=(timestamp, timestamp))

		data = self.get(data_id=data_id, timewindow=timewindow, with_id=True)

		_set = {
					TimedStore.DATA_ID: data_id,
					TimedStore.TIMESTAMP: timestamp,
					TimedStore.VALUE: value
				}

		if len(data) == 0 or data[0][0] != timestamp:
			values_to_insert = {
					TimedStore.DATA_ID: data_id,
					TimedStore.TIMESTAMP: timestamp,
					TimedStore.VALUE: value
				}
			self.collection.insert(values_to_insert, w=1)

		else:
			_id = data[0][2]
			value_to_update = {TimedStore.VALUE: value}
			self.collection.update(
				{'_id': _id},
				{
					'$set': value_to_update
				},
				w=1)

	def remove(self, data_id, timewindow=None):
		"""
		Remove timed_data existing on input timewindow.
		"""

		where = {
			TimedStore.DATA_ID: data_id
		}

		if timewindow is not None:
			where[TimedStore.TIMESTAMP] = \
				{'$gte': timewindow.start(), '$lte': timewindow.stop()}

		self.collection.remove(where, w=1)

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

	def _get_indexes(self):
		"""
		Get indexes.
		"""

		result = [
			'_id',
			PeriodicStore._get_query_hint()
		]

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

		self.check_connection()

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

		cursor = self.collection.find(where)
		cursor.hint(PeriodicStore._get_query_hint())

		result = cursor.count()

		return result

	def put(self, data_id, aggregation, period, points):
		"""
		Put periodic points in periodic collection with specific aggregation and
		period values.

		points is an iterable of (timestamp, value)
		"""

		# initialize a dictionary of perfdata value by value field and id_timestamp
		document_properties_by_id_timestamp = dict()
		# previous variable contains a dictionary of entries to put in the related document

		# prepare data to insert/update in db
		for timestamp, value in points:

			timestamp = int(timestamp)
			id_timestamp = int(period.round_timestamp(timestamp, normalize=True))
			document_properties = \
				document_properties_by_id_timestamp.setdefault(id_timestamp, dict())

			self.logger.debug(' + id_timestamp: {0}'.format(id_timestamp))

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

			result = self.collection.update(
				{
					'_id': _id
				},
				{
					'$set': _set
				},
				upsert=True,
				w=1)

			self._manage_query_error(result)

	def get(self, data_id, aggregation, period, timewindow=None):
		"""
		Get a list of points.
		"""

		query = PeriodicStore._get_documents_query(data_id=data_id,
			aggregation=aggregation, timewindow=timewindow, period=period)

		projection = {
			PeriodicStore.TIMESTAMP: 1,
			PeriodicStore.VALUES: 1
		}

		if aggregation is None:
			projection[PeriodicStore.AGGREGATION] = 1

		if period is None:
			projection[PeriodicStore.PERIOD] = 1

		cursor = self.collection.find(query, projection=projection)

		cursor.hint(PeriodicStore._get_query_hint())

		result = list()

		doc_aggregation = aggregation
		doc_period = period

		for document in cursor:

			doc_aggregation = document.get(PeriodicStore.AGGREGATION, aggregation)
			doc_period = document.get(PeriodicStore.PERIOD, period)
			timestamp = int(document[PeriodicStore.TIMESTAMP])

			values = document[PeriodicStore.VALUES]

			for t, value in values.iteritems():
				value_timestamp = timestamp + int(t)

				if timewindow is None or value_timestamp in timewindow:
					result.append( (value_timestamp, value) )

		result.sort(key=itemgetter(0))

		return result

	def remove(self, data_id, aggregation=None, period=None, timewindow=None):
		"""
		Remove periodic data related to data_id, timewindow and period.
		If timewindow is None, remove all periodic_data with input period.
		If period is None
		"""

		query = PeriodicStore._get_documents_query(data_id=data_id,
			aggregation=aggregation, timewindow=timewindow, period=period)

		if timewindow is not None:

			projection = {
				PeriodicStore.TIMESTAMP: 1,
				PeriodicStore.VALUES: 1
			}

			documents = self.collection.find(query, projection=projection)

			for document in documents:
				timestamp = document.get(PeriodicStore.TIMESTAMP)
				values = document.get(PeriodicStore.VALUES)
				values_to_save = \
					{t: v for t, v in values.iteritems() if (timestamp + int(t)) not in timewindow}
				_id = document.get('_id')

				if len(values_to_save) > 0:
					self.collection.update(
						{
							'_id': _id
						},
						{
							'$set': {PeriodicStore.VALUES: values_to_save}
						},
						w=1)
				else:
					self.collection.remove(_id, w=1)

		else:
			self.collection.remove(query)

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


	@staticmethod
	def _get_documents_query(data_id, aggregation, timewindow, period):
		"""
		Get mongo documents query about data_id, timewindow and period.

		If period is None and timewindow is not None, period takes default period value for data_id.
		"""

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
			result[PeriodicStore.TIMESTAMP] = \
				{'$gte': start_timestamp, '$lte': stop_timestamp}

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
			raise Manager.ManagerError(
				"period {0} must contain at least one valid unit among {1}".
				format(period, Period.UNITS))
		md5_result.update(unit_with_value[Period.UNIT].encode('ascii', 'ignore'))

		# resolve md5
		result = md5_result.hexdigest()

		return result
