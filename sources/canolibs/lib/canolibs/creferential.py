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

__all__ = ('Referential')

import ConfigParser
from os.path import expanduser

import logging

# import to remove if mongo
try
	from pymongo import MongoClient as Connection
except ImportError:
	from pymongo import Connection

CONFIG = ConfigParser.RawConfigParser()
CONFIG.read(expanduser('~/etc/referential.conf'))


class Referential(object):
	"""
	Manage access to a referential.
	"""

	DEFAULT_CONFIGURATION_FILE = "~/etc/referential.conf"

	def __init__(self, logging_level=logging.INFO,
		configuration_file=DEFAULT_CONFIGURATION_FILE):

		super(Referential, self).__init__()

		self.logger = logging.getLogger('cstorage')
		self.logger.setLevel(logging_level)

	def load_configuration_file(self):

		self.host = CONFIG.get("master", "host")

		self.port = CONFIG.getint("master", "port")

		self.db_name = CONFIG.get("master", "db")

		self.auto_connect = CONFIG.getbool('master', 'auto_connect')

		self.backend_name = CONFIG.get('master', 'backend')

		self.journaling = CONFIG.getbool('master', 'journaling')

		self.ssl = CONFIG.getbool('master', 'ssl')

		self.ssl_keyfile = CONFIG.get('master', 'ssl_keyfile')

		self.ssl_certfile = CONFIG.get('master', 'ssl_certfile')

		if self.auto_connect:
			self.connect()

	def connect(self):
		"""
		Connect this referential to the backend.
		"""

		self.logger.debug('Trying to connect to {0}:{1}'.format(self.host, self.port))

		try:
			self.conn = Connection(
				host=self.mongo_host, port=self.mongo_port,
				j=self.journaling, ssl=self.ssl,
				ssl_keyfile=self.ssl_keyfile, ssl_certfile=self.ssl_certfile)


		self.logger.debug("Connected on {0}:{1}".format(self.host, self.port))

		self.db = self.conn[self.db_name]

		self.backend = self.db[self.backend_name]

		self.connected = True

	def disconnect(self):
		self.conn.disconnect()

	def connected(self):
		"""
		:returns: True if this is connected.
		"""

		return self.conn.alive()

	def get_entity(self):
		pass

	def get_entities(self):
		pass

	def get_connectors(self):
		pass

	def get_components(self):
		pass

	def get_resources(self):
		pass


	def delete_entities(self):
		pass



def get_entity(id):
	pass


def update_entity(entity):
	pass