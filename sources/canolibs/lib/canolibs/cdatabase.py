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
			'MASTER': {
				'HOST': str,
				'PORT': int,
				'DB': str,
				'AUTO_CONNECT': bool,
				'BACKEND': str,
				'JOURNALING': str,
				'SSL': bool,
				'SSL_KEYFILE': str,
				'SSL_CERTFILE': str
			}
		}

	def __init__(self, host=None, port=0, db=None, auto_connect=False,
		journaling=False, ssl=False, ssl_keyfile=None, ssl_certfile=None,
		**kwargs):
		"""
		:param host: db host name
		:param port: db port
		:param db: db name
		:param auto_connect: auto connect to database when initialised.
		:param backend: default backend to use.
		:param journaling: journaling mode enabling.
		:param ssl: ssl mode
		:param ssl_keyfile: ssl keys file.
		:param ssl_certfile: ssl certification file.

		:type host: str
		:type port: int
		:type db: str
		:type auto_connect: bool
		:type backend: str
		:type journaling: bool
		:type ssl: bool
		:type ssl_keyfile: str
		:type ssl_certfile: str
		"""

		super(DataBase, self).__init__(**kwargs)

		# initialize instance properties with default values
		self.host = None
		self.port = 0
		self.db = None
		self.auto_connect = False
		self.backend = None
		self.journaling = False
		self.ssl = False
		self.ssl_keyfile = None
		self.ssl_certfile = None

	def get_parsers_by_option_by_section(self, **kwargs):

		result = super(DataBase, self).get_parsers_by_option_by_section()

		result.update({
			DataBase.PARSERS_BY_OPTIONS_BY_SECTIONS
			})

		return result

	def apply_configuration(self, parsers_by_option_by_section=None,
		configuration_file=None, naming_rule=None, **kwargs):
		"""
		Load configuration file and connect if self.auto_connect.
		"""

		super(DataBase, self).apply_configuration(
			parsers_by_option_by_section=parsers_by_option_by_section,
			configuration_file=configuration_file, naming_rule=naming_rule)

		if self.auto_connect:
			self.connect()

	def get_backend(self, backend=None):
		"""
		Get a reference to a specific backend where name is input backend.
		If input backend is None, self.backend is used.

		:param backend: backend name. If None, self.backend is used.
		:type backend: basestring

		:returns: backend reference.

		:raises: NotImplementedError
		.. seealso: DataBase.set_backend(self, backend)
		"""

		raise NotImplementedError()

	def set_backend(self, backend):
		"""
		Change of backend with input backend.

		:note:

		:raises: NotImplementedError
		.. seealso: DataBase.get_backend(self, backend=None)
		"""

		raise NotImplementedError()

	def connect(self):
		"""
		Connect this database.

		:raises: NotImplementedError

		.. seealso:: disconnect(self), connected(self)
		"""

		raise NotImplementedError()

	def disconnect(self):
		"""
		Disconnect this database.

		:raises: NotImplementedError
		"""

		raise NotImplementedError()

	def connected(self):
		"""
		:returns: True if this is connected.
		"""

		raise NotImplementedError()

	def get_element(self, id, backend=None):
		"""
		:param id: id of the element to get.
		:type id: basestring

		:return: input id element.
		:rtype: dict

		:raises NotImplementedError:
		"""

		raise NotImplementedError()

	def delete_element(self, id, backend=None):

		raise NotImplementedError()

	def update_element(self, id, element, backend=None):

		raise NotImplementedError()

	def drop(self, backend=None):

		raise NotImplementedError()
