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

import pyperfstore
from pyperfstore.storage import storage

class memstore(storage):
	def __init__(self):
		storage.__init__(self)
		self.logger.debug(" + Init Mem Store")

		self.data = {}

	def set_raw(self, key, value):
		self.data[key] = value

	def set(self, key, value):
		self.logger.debug("Set '%s'" % key)
		self.set_raw(key, value)

	def get_raw(self, key):
		return self.data[key]
		
	def get(self, key):
		self.logger.debug("Get '%s'" % key)
		try:
			return self.get_raw(key)
		except:
			return None

	def rm(self, key):
		self.logger.debug("Remove '%s'" % key)
		del self.data[key]

	def append(self, key, value):
		self.logger.debug("Append data in '%s'" % key)
		self.logger.debug(" + Key: '%s'" % key)
		self.logger.debug(" + Value: '%s'" % value)
		
		try:
			self.data[key].append(value)
		except:
			self.data[key] = [ value ]

	def size(self, key):
		return int(sys.getsizeof(self.data[key]))		
		
	def get_all_nodes(self):
		index = []
		for key,value in self.data.items():
			if isinstance(value,dict):
				if 'metrics' in value:
					index.append(key)
		return index
		
	def get_all_metrics(self):
		index = []
		for key,value in self.data.items():
			if isinstance(value,dict):
				if 'dn' in value and not 'metrics' in value:
					index.append({'node':value['node_id'],'metric':value['dn']})
		return index
		
		

	def lock(self, key):
		self.logger.debug("Lock '%s'" % key)
		## Todo

	def wait_lock(self, key):
		self.logger.debug("Wait Lock '%s'" % key)
		## Todo	
