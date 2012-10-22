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

from pyperfstore.storage import storage

class filestore(storage):
	def __init__(self, base_path):
		storage.__init__(self)

		self.logger.debug(" + Init Files Store")
		self.base_path = base_path

		try:
			os.mkdir(self.base_path)
		except:
			pass

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
		filename = "%s/%s" % (self.base_path, key)
		file = open(filename, "w")
		file.write(value)
		file.close()

	def set(self, key, value):
		value = self.encode(value)

		#self.logger.debug("Set '%s'" % key)
		self.set_raw(key, value)

	def get_raw(self, key):
		try:
			filename = "%s/%s" % (self.base_path, key)
			file = open(filename, "r")
			data = file.read()
			file.close()
			return data
		except:
			return None
		
	def get(self, key):
		self.logger.debug("Get '%s'" % key)
		try:
			return self.decode(self.get_raw(key))
		except:
			return None

	def rm(self, key):
		self.logger.debug("Remove '%s'" % key)
		try:
			filename = "%s/%s" % (self.base_path, key)
			os.remove(filename)
		except Exception, err:
			self.logger.error("Impossible to remove '%s' (%s)" % (filename, err))

	def append(self, key, value):
		self.logger.debug("Append data in '%s'" % key)
		self.logger.debug(" + Key: '%s'" % key)
		self.logger.debug(" + Value: '%s'" % value)

		raw = self.get_raw(key)
		#self.logger.debug(" + Raw: '%s'" % raw)
		
		if raw:
			raw = raw[:-1] + "," + self.encode(value) + "]"
			self.set_raw(key, raw)
		else:
			self.logger.debug("   + Init key")
			self.set(key, [value])

	def size(self, key):
		#self.logger.debug("Get size of '%s'" % key)
		return int(sys.getsizeof(self.get_raw(key)))
		

	def lock(self, key):
		self.logger.debug("Lock '%s'" % key)

	def wait_lock(self, key):
		self.logger.debug("Wait Lock '%s'" % key)


	def pretty_print(self):
		for key in self.store:
			self.logger.debug("%s: %s" % (key, self.store[key]))	
