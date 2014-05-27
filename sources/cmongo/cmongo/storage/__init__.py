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

from cstorage import Storage

from cdatabase import DataBase

# import to remove if mongo
try:
	from pymongo import MongoClient as Connection
except ImportError:
	from pymongo import Connection

from pymongo.errors import TimeoutError, OperationFailure, ConnectionFailure


class Storage(Storage):
	"""
	Manage access to a mongodb.
	"""

	BACKEND = 'backend'

	def __init__(self, data_type, backend=None, host='localhost', port=27017,
		*args, **kwargs):
		"""
		:param backend: default backend to use.
		:type backend: str.
		"""

		super(Storage, self).__init__(
			_ready_to_conf=False, data_type=data_type, host=host, port=port,
			*args, **kwargs)

		self.backend = backend

		self._connected = False

		self._conn = None

		if self.auto_conf:
			self.apply_configuration()

	def get_parsing_rules(self, *args, **kwargs):

		result = super(Storage, self).get_parsing_rules(*args, **kwargs)

		storage_parsing_rule = result[-1]

		section = self._get_section()

		storage_parsing_rule[section][Storage.BACKEND] = str

		return result

	def connect(self, *args, **kwargs):
		"""
		Connect this referential to the backend.
		"""

		if not self._connected:

			self.logger.debug('Trying to connect to {0}:{1}'.format(
				self.host, self.port))

			try:
				self._conn = \
					Connection(host=self.host, port=self.port,
					j=self.journaling, ssl=self.ssl,
					# TODO: remove comment with mongo version > 2.5
					# ssl_keyfile=self.key, ssl_certfile=self.cert,
					w=1 if self.safe else 0)

			except ConnectionFailure as e:
				self.logger.error('Raised Exception during connection attempting \
					to {0}:{1}. {0}'.format(self.host, self.port, e))

			else:
				self._db = self._conn[self.db]

				if (self.user, self.pwd) != (None, None):
					authenticate = self._db.authenticate(self.user, self.pwd)

					if authenticate:
						self.logger.debug("Connected on {0}:{1}".format(
							self.host, self.port))
						self._connected = True

					else:
						self.logger.error(
							'Impossible to authenticate {0} on {1}:{2}'.format(
								self.host, self.port))
						self.disconnect()

				else:
					self._connected = True
					self.logger.debug("Connected on {0}:{1}".format(
						self.host, self.port))

			if self._connected:
				indexes = self._get_indexes()

				for index in indexes:
					self._get_backend().ensure_index(index)

		result = self._connected

		return result

	def disconnect(self, *args, **kwargs):
		"""
		"""

		if self._conn is not None:
			self._conn.disconnect()

		self._connected = False

	def connected(self, *args, **kwargs):
		"""
		:returns: True if this is connected.
		"""

		result = self._connected

		return result

	def size(self, criteria=None, hint=None, backend=None, *args, **kwargs):
		"""
		Get size of collection. None if a connection error occured.
		"""

		result = 0

		_backend = self._get_backend(backend=backend)

		try:
			result = self._db.command("collstats", _backend)['size']

		except Exception as e:
			self.logger.warning(
				"Impossible to read Collection Size: {0}".format(e))
			result = None

		return result

	def count(self, criteria=None, hint=None, backend=None, *args, **kwargs):

		_backend = self._get_backend(backend=backend)

		if criteria is not None:
			cursor = _backend.find(criteria)

			if hint is not None:
				cursor.hint(hint)
			result = cursor.count()

		else:
			result = _backend.count()

		return result

	def drop(self, backend=None):
		"""
		Drop all data.
		"""

		self._get_backend(backend=backend).remove()

	def get_collection_name(self, data_type=None):
		"""
		Get collection name managed by this store related to input data_type.
		"""

		if data_type is None:
			data_type = self._data_type

		result = "{0}_{1}".format(self._get_collection_prefix(), data_type)

		return result

	def _get_collection_prefix(self):
		result = type(self).__name__

		return result

	def _get_backend(self, backend=None, *args, **kwargs):
		"""
		Get a reference to a specific backend where name is input backend.
		If input backend is None, self.backend is used.

		:param backend: backend name. If None, self.backend is used.
		:type backend: basestring

		:returns: backend reference.

		:raises: NotImplementedError
		.. seealso: DataBase.set_backend(self, backend)
		"""

		if backend is None:
			backend = self.backend

		result = self._db[backend]

		return result

	def _get_indexes(self):
		"""
		Get collection indexes. Must be overriden.
		"""

		result = [
			[('_id', 1)]
		]

		return result

	def _manage_query_error(self, result_query):
		"""
		Manage mongo query error.
		"""

		result = None

		if isinstance(result_query, dict):

			error = result_query.get("writeConcernError", None)

			if error is not None:
				self.logger.error(' error in writing document: {0}'.format(error))

			error = result_query.get("writeError")

			if error is not None:
				self.logger.error(' error in writing document: {0}'.format(error))

		else:

			result = result_query

		return result

	def _insert(self, document, backend=None):

		result = self._run_command(
			backend=backend, command='insert', doc_or_docs=document)

		return result

	def _update(self, _id, document, backend=None):

		result = self._run_command(
			backend=backend, command='update', spec=_id, document=document,
			upsert=True, multi=True)

		return result

	def _find(self, document, projection=None, backend=None):

		result = self._run_command(
			backend=backend, command='find', spec=document, projection=projection)

		return result

	def _remove(self, document, backend=None):

		result = self._run_command(
			backend=backend, command='remove', spec_or_id=document)

		return result

	def _run_command(self, command, backend=None, **kwargs):

		result = None

		attempts = DataBase.MAX_ATTEMPTS

		while attempts > 0:
			try:
				backend_command = getattr(self._get_backend(backend=backend), command)
				result = backend_command(w=1 if self.safe else 0, wtimeout=self.wtimeout,
					**kwargs)

				self._manage_query_error(result)

				break

			except TimeoutError:
				attempts -= 1
				self.logger.warning('Try to run command {0}({1}) on {2}, {3} attempts left'
					.format(command, kwargs, backend, attempts))

			except OperationFailure as of:
				self.logger.error('{0} during running command {1}({2}) of in {3}'
					.format(of, command, kwargs, backend))
				break

		return result


class Store(Storage):
	"""
	Manage periodic and timed access to data_type in a mongo collection.
	"""

	def __init__(self, data_type, *args, **kwargs):

		super(Store, self).__init__(
			backend=self.get_collection_name(data_type=data_type),
			_ready_to_conf=False, *args, **kwargs)

		self._data_type = data_type

		if self.auto_conf:
			self.apply_configuration()

	def get_parsing_rules(self, *args, **kwargs):

		result = dict()

		mongodb_parsers = super(Store, self).get_parsing_rules()

		for section in list(mongodb_parsers.keys()):
			new_section_name = "{0}_{1}_{2}".format(
				section, type(self), self._data_type)
			result[new_section_name] = mongodb_parsers[section]

		return result

	def _get_data_id(self, data_id):
		"""
		Get data id if data_id is an entity _id.
		"""

		result = data_id

		if data_id.startswith(self._data_type):
			result = data_id[len(self._data_type) + 1:]

		return result
