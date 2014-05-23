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

__all__ = ('MongoDB')

from cdatabase import DataBase

# import to remove if mongo
try:
	from pymongo import MongoClient as Connection
except ImportError:
	from pymongo import Connection


class MongoDB(DataBase):
	"""
	Manage access to a mongodb.
	"""

	MONGO_SECTION = 'MONGO'

	PARSERS_BY_OPTIONS_BY_SECTIONS = {
		MONGO_SECTION: {
			'backend': str
		}
	}

	def __init__(self, backend=None, host='localhost', port=27017, *args, **kwargs):
		"""
		:param backend: default backend to use.

		:type backend: str.
		"""

		super(MongoDB, self).__init__(host=host, port=port, *args, **kwargs)

		self.backend = backend

		self._connected = False

		self._conn = None

	def get_parsers_by_option_by_section(self, *args, **kwargs):

		result = super(MongoDB, self).get_parsers_by_option_by_section()

		result.update({
			MongoDB.PARSERS_BY_OPTIONS_BY_SECTIONS
			})

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
					w=1 if self.safe else 0)

			except Exception as e:
				self.logger.error('Raised Exception during connection attempting \
					to {0}:{1}: {0}'.format(self.host, self.port, e))

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

		_backend = self._get_backend(backend)

		try:
			result = self._db.command("collstats", _backend)['size']

		except Exception as e:
			self.logger.warning(
				"Impossible to read Collection Size: {0}".format(e))
			result = None

		return result

	def count(self, criteria=None, hint=None, backend=None, *args, **kwargs):

		_backend = self._get_backend(backend)

		if criteria is not None:

			cursor = _backend.find(criteria)
			if hint is not None:
				cursor.hint(hint)
			result = cursor.count()

		else:
			result = _backend.count()

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
