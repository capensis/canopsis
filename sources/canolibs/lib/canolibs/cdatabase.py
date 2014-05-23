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

__all__ = ('DataBase')

from cconfiguration import Configurable


class DataBase(Configurable):
	"""
	Abstract class which aims to manage access to a data base.

	Related to a configuration file, it can connects to a database
	depending on several parameters like.

	:param host: db host name
	:type host: basestring
	:param port: db port
	:type port: int
	:param db: db name
	:type db: basestring
	:param auto_connect: auto connect to database when initialised.
	:type auto_connect: bool
	:param backend: default backend to use.

	It provides a DataBaseError for internal errors
	"""

	class DataBaseError(Exception):
		"""
		Errors raised by the DataBase class.
		"""

		pass

	PARSERS_BY_OPTIONS_BY_SECTIONS = {
			'DATABASE': {
				'host': str,
				'port': int,
				'db': str,
				'auto_connect': bool,
				'journaling': str,
				'safe': bool,
				'ssl': bool,
				'ssl_keyfile': str,
				'ssl_certfile': str,
				'user': str,
				'pwd': str
			}
		}

	ASC = 1  # ASC order
	DESC = -1  # DESC order

	def __init__(self, host='localhost', port=0, db='canopsis', auto_connect=True,
		journaling=False, safe=False, ssl=False, ssl_keyfile=None, ssl_certfile=None,
		user=None, pwd=None, _ready_to_connect=True, *args, **kwargs):
		"""
		:param host: db host name
		:param port: db port
		:param db: db name
		:param auto_connect: auto connect to database when initialised.
		:param backend: default backend to use.
		:param journaling: journaling mode enabling.
		:param safe: ensure writing data.
		:param ssl: ssl mode
		:param ssl_keyfile: ssl keys file.
		:param ssl_certfile: ssl certification file.
		:param user: user
		:param pwd: password

		:type host: str
		:type port: int
		:type db: str
		:type auto_connect: bool
		:type backend: str
		:type journaling: bool
		:type safe: bool
		:type ssl: bool
		:type ssl_keyfile: str
		:type ssl_certfile: str
		:type user: str
		:type pwd: str
		"""

		super(DataBase, self).__init__(*args, **kwargs)

		# initialize instance properties with default values
		self.host = host
		self.port = port
		self.db = db
		self.auto_connect = auto_connect
		self.journaling = journaling
		self.safe = safe
		self.ssl = ssl
		self.ssl_keyfile = ssl_keyfile
		self.ssl_certfile = ssl_certfile
		self.user = user
		self.pwd = pwd

		if _ready_to_connect and self.auto_connect:
			self.connect()

	def get_parsers_by_option_by_section(self, *args, **kwargs):

		result = super(DataBase, self).get_parsers_by_option_by_section()

		result.update(
			DataBase.PARSERS_BY_OPTIONS_BY_SECTIONS
			)

		return result

	def apply_configuration(self, parsers_by_option_by_section=None,
		configuration_file=None, naming_rule=None, *args, **kwargs):
		"""
		Load configuration file and connect if self.auto_connect.
		"""

		super(DataBase, self).apply_configuration(
			parsers_by_option_by_section=parsers_by_option_by_section,
			configuration_file=configuration_file, naming_rule=naming_rule)

		if self.auto_connect:
			self.disconnect()
			self.connect()

	def connect(self, *args, **kwargs):
		"""
		Connect this database.

		:raises: NotImplementedError

		.. seealso:: disconnect(self), connected(self)
		"""

		raise NotImplementedError()

	def disconnect(self, *args, **kwargs):
		"""
		Disconnect this database.

		:raises: NotImplementedError
		"""

		raise NotImplementedError()

	def connected(self, *args, **kwargs):
		"""
		:returns: True if this is connected.
		"""

		raise NotImplementedError()

	def get_element(self, id, backend=None, *args, **kwargs):
		"""
		:param id: id of the element to get.
		:type id: basestring

		:return: input id element.
		:rtype: dict

		:raises NotImplementedError:
		"""

		raise NotImplementedError()

	def delete_element(self, id, backend=None, *args, **kwargs):

		raise NotImplementedError()

	def update_element(self, id, element, backend=None, *args, **kwargs):

		raise NotImplementedError()

	def drop(self, backend=None, *args, **kwargs):

		raise NotImplementedError()

	def size(self, criteria=None, backend=None, *args, **kwargs):

		raise NotImplementedError()
