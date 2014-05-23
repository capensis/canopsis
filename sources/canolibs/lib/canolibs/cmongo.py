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

import logging

# import to remove if mongo
try
	from pymongo import MongoClient as Connection
except ImportError:
	from pymongo import Connection


class MongoDB(DataBase):
	"""
	Manage access to a mongodb.
	"""

	DEFAULT_CONFIGURATION_FILE = "~/etc/mongodb.conf"

	def __init__(self,
		configuration_file=DEFAULT_CONFIGURATION_FILE, backend=None, **kwargs):
		"""
		:param backend: default backend to use.

		:type backend: str
		"""
		super(MongoDB, self).__init__(configuration_file=configuration_file,
			**kwargs)

		self.backend = backend

	PARSERS_BY_OPTIONS_BY_SECTIONS = {
		'OPTION': {
			'backend': str
		}
	}

	def get_parsers_by_option_by_section(self, **kwargs):

		result = super(MongoDB, self).get_parsers_by_option_by_section()

		result.update({
			MongoDB.PARSERS_BY_OPTIONS_BY_SECTIONS
			})

		return result

	def connect(self):
		"""
		Connect this referential to the backend.
		"""

		self.logger.debug('Trying to connect to {0}:{1}'.format(
			self.host, self.port))

		try:
			self.conn = \
				Connection(host=self.mongo_host, port=self.mongo_port,
				j=self.journaling, ssl=self.ssl,
				ssl_keyfile=self.ssl_keyfile, ssl_certfile=self.ssl_certfile)

		self.logger.debug("Connected on {0}:{1}".format(self.host, self.port))

		self.db = self.conn[self.db_name]

		self.backend = self.db[self.backend_name]

	def disconnect(self):
		self.conn.disconnect()

	def connected(self):
		"""
		:returns: True if this is connected.
		"""

		return self.conn.alive()

	def size(self, criteria=None, hint=None, backend=None):

		raise NotImplementedError()

	def count(self, criteria=None, hint=None, backend=None):

		raise NotImplementedError()

	def copy_to_backend(self, criteria, backends=None):

		raise NotImplementedError()
