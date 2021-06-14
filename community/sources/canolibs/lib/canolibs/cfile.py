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
from crecord import crecord
import md5
namespace = 'files'

class cfile(crecord):
	def __init__(self, storage=None, *args, **kargs):
		crecord.__init__(self, storage=storage, *args, **kargs)
		self.type = 'bin'
		self.binary = None

	def put_data(self, bin_data, file_name=None, content_type=None):
		self.binary = bin_data
		self.data['file_name'] = file_name
		self.name = file_name
		self.data['content_type'] = content_type 

	def put_file(self, path, file_name=None, content_type=None):
		self.binary = open(path,'r').read()
		self.data['file_name'] = file_name
		self.name = file_name
		self.data['content_type'] = content_type 
		if not self._id:
			self._id = md5.md5(file_name).hexdigest()
	
	def get_binary_id(self):
		bid = self.data.get('binary_id', None)
		if not bid:
			bid = self.data.get('data_id', None)
			
		return bid		
	
	def get(self, storage=None):
		if not storage:
			storage = self.storage
		
		bid = self.get_binary_id()
		
		if storage:
			return storage.get_binary(bid)
		else:
			raise Exception("You must specify storage (GET)")

	def remove(self, storage=None):
		if not storage:
			storage = self.storage
		
		bid = self.get_binary_id()	
		
		if storage:
			storage.remove_binary(bid)
			storage.remove(self._id, namespace=namespace)
		else:
			raise Exception("You must specify storage (REMOVE)")

	def check(self, storage=None):
		if not storage:
			storage = self.storage
		
		bid = self.get_binary_id()	
		
		if storage:
			return storage.check_binary(bid)
		else:
			raise Exception("You must specify storage (CHECK)")


def get_cfile(_id, storage):
	record = storage.get(_id, namespace=namespace)
	rfile = cfile(storage=storage, record=record)
	return rfile
