#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

import os, sys, json, logging

class storage(object):
	def __init__(self, logging_level=logging.INFO):
		self.logger = logging.getLogger('storage')
		self.logger.setLevel(logging_level)

	def encode(self, value):
		try:
			return json.dumps(value)
		except:
			return value

	def decode(self, value):
		try:
			return json.loads(str(value))
		except:
			return value

	def set_raw(self, key, value):
		## Todo
		pass

	def set(self, key, value):
		self.logger.debug("Set '%s'" % key)
		value = self.encode(value)
		self.set_raw(key, value)

	def get_raw(self, key):
		## Todo
		return []
		
	def get(self, key):
		self.logger.debug("Get '%s'" % key)
		try:
			return self.decode(self.get_raw(key))
		except:
			return None

	def rm(self, key):
		self.logger.debug("Remove '%s'" % key)
		## Todo

	def append(self, key, value):
		self.logger.debug("Append data in '%s'" % key)
		self.logger.debug(" + Key: '%s'" % key)
		self.logger.debug(" + Value: '%s'" % value)
		
		## Todo

	def size(self, key):
		## Todo
		return 0
		
	def get_all_nodes(self):
		return []
		
	def get_all_metrics(self):
		return []

	def lock(self, key):
		self.logger.debug("Lock '%s'" % key)
		## Todo

	def wait_lock(self, key):
		self.logger.debug("Wait Lock '%s'" % key)
		## Todo	
