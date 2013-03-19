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

import os, sys, json, logging, time

from bson.errors import InvalidStringData
from pymongo import Connection
from gridfs import GridFS

class store(object):
	def __init__(self, mongo_host="127.0.0.1", mongo_port=27017, mongo_db='canopsis', mongo_collection='perfdata2', mongo_user=None, mongo_pass=None, mongo_safe=False, logging_level=logging.INFO):
		self.logger = logging.getLogger('store')
		self.logger.setLevel(logging_level)
		
		self.logger.debug(" + Init MongoDB Store")

		# Read db option from conf
		import ConfigParser
		config = ConfigParser.RawConfigParser()
		config.read(os.path.expanduser('~/etc/cstorage.conf'))

		try:
			mongo_host = config.get('master', 'host')		if config.get('master', 'host') != "" and not mongo_host else mongo_host
			mongo_port = config.getint('master', 'port')	if config.get('master', 'port') != "" and not mongo_port else mongo_port
			mongo_user = config.get('master', 'userid')		if config.get('master', 'userid') != "" and not mongo_user else mongo_user
			mongo_pass = config.get('master', 'password')	if config.get('master', 'password') != "" and not mongo_pass else mongo_pass
		except:
			pass

		self.mongo_host = mongo_host
		self.mongo_port = mongo_port
		self.mongo_db = mongo_db
		self.mongo_collection = mongo_collection
		self.mongo_safe = mongo_safe
		self.mongo_user = mongo_user if mongo_user != "" else None
		self.mongo_pass = mongo_pass if mongo_pass != "" else None

		
		self.connected = False
		
		self.connect()

	def connect(self):
		if self.connected:
			self.logger.warning("Impossible to connect, already connected")
			return True
		else:
			self.logger.debug("Connect to MongoDB (%s/%s@%s:%s)" % (self.mongo_db, self.mongo_collection, self.mongo_host, self.mongo_port))
			
			try:
				self.conn=Connection(host=self.mongo_host, port=self.mongo_port, safe=self.mongo_safe)
				self.logger.debug(" + Success")
			except Exception, err:
				self.logger.error(" + %s" % err)
				return False
				
			self.db=self.conn[self.mongo_db]

			try:
				self.logger.debug("Try to auth '%s'" % self.mongo_user)
				if self.mongo_user and self.mongo_pass != None:
						if not self.db.authenticate(self.mongo_user, self.mongo_pass):
							raise Exception('Invalid user or pass.')
						self.logger.debug(" + Success")
			except Exception, err:
				self.logger.error(" + Impossible to authenticate: %s" % err)
				self.disconnect()
				return False

			self.logger.debug("Get collections")
			self.collection = self.db[self.mongo_collection]
			self.grid = GridFS(self.db, self.mongo_collection+"_bin")
			self.connected = True
			self.logger.debug(" + Success")
			return True
			
	def check_connection(self):
		if not self.connected:
			if not self.connect():
				raise Exception('Impossible to deal with DB, you are not connected ...')
						
	def count(self, _id):
		return self.collection.find({'_id': _id}).count()
		
	def update(self, _id, mset=None, munset=None, mpush=None, mpush_all=None, mpop=None, upsert=True):
		self.check_connection()
		data = {}
		if mset:
			data['$set'] = mset
		if munset:
			data['$unset'] = munset
		if mpush:
			data['$push'] = mpush
		if mpush_all:
			data['$pushAll'] = mpush_all
		if mpop:
			data['$pop'] = mpop
		
		if data:
			return self.collection.update({'_id': _id}, data, upsert=upsert)
	
	def push(self, _id, point, meta_data={}):
		self.check_connection()
		self.logger.debug("Push point '%s' in '%s'" % (point, _id))
		
		meta_data['lts'] = point[0]
		meta_data['lv'] = point[1]
		
		return self.update(_id=_id, mset=meta_data, mpush={'d': point})

	def create(self, _id, data):
		self.check_connection()
		data['_id'] = _id
		self.logger.debug("Create record '%s'" % _id)
		return self.collection.insert(data)

	def create_bin(self, _id, data):
		self.check_connection()
		self.logger.debug("Create bin record '%s'" % _id)
		return self.grid.put(data, _id=_id)
			
	def remove(self, _id=None, mfilter=None):
		self.check_connection()
		if mfilter:
			return self.collection.remove(mfilter)
		elif _id:
			return self.collection.remove({'_id': _id})

	def size(self):
		self.logger.info("Size of dbs:")
		size = 0
		try:
			size = self.db.command("collstats", self.mongo_collection)['size']
		except:
			self.logger.warning("Impossible to read Collecion Size")
		
		self.logger.info(" + Collection:    %0.2f MB" % (size/1024.0/1024.0))
		try:			
			bin_size = self.db.command("collstats", self.mongo_collection+"_bin.files")['size']
			self.logger.info(" + Binaries Meta: %0.2f MB" % (bin_size /1024.0/1024.0))
			
			chunks_size = self.db.command("collstats", self.mongo_collection+"_bin.chunks")['size']
			self.logger.info(" + Binaries:      %0.2f MB" % (chunks_size /1024.0/1024.0))
			
			size += chunks_size + bin_size
		except:
			self.logger.warning("Impossible to read GridFS Size")
			pass

		
		return size
	
	def get(self, _id, mfields=None):
		self.check_connection()
		return self.collection.find_one({'_id': _id}, fields=mfields)
	
	def get_bin(self, _id):
		self.check_connection()
		return self.grid.get(_id).read()

	def find(self, limit=0, skip=0, mfilter={}, mfields=None, sort=None):
		self.check_connection()
		if limit == 1:
			return self.collection.find_one(mfilter, limit=limit, fields=mfields, sort=sort)
		else:		
			return self.collection.find(mfilter, limit=limit, skip=skip, fields=mfields, sort=sort)
							
	def drop(self):
		self.check_connection()
		self.db.drop_collection(self.mongo_collection)
		self.db.drop_collection(self.mongo_collection+"_bin.chunks")
		self.db.drop_collection(self.mongo_collection+"_bin.files")
		
	def disconnect(self):
		if self.connected:
			self.logger.debug("Disconnect from MongoDB")
			self.conn.disconnect()
		else:
			self.logger.warning("Impossible to disconnect, you are not connected")

	def __del__(self):
		self.disconnect()
