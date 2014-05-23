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

from cmongo import MongoDB

import logging

# import to remove if mongo
try
	from pymongo import MongoClient as Connection
except ImportError:
	from pymongo import Connection

CONFIG = ConfigParser.RawConfigParser()
CONFIG.read(expanduser('~/etc/referential.conf'))


class Referential(MongoDB):
	"""
	Manage access to a referential.
	"""

	DEFAULT_CONFIGURATION_FILE = "~/etc/referential.conf"

	def __init__(self, configuration_file=DEFAULT_CONFIGURATION_FILE,
		**kwargs):

		super(Referential, self).__init__(
			configuration_file=configuration_file, **kwargs)

	def get_entity(self, id):

		raise NotImplementedError()

	def get_entities(self):

		raise NotImplementedError()

	def get_connectors(self):

		raise NotImplementedError()

	def get_components(self):

		raise NotImplementedError()

	def get_resources(self):

		raise NotImplementedError()

	def delete_entities(self):

		raise NotImplementedError()
