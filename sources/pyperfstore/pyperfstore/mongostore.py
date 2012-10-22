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

from pymongo import ASCENDING
from pymongo import Connection
from bson.errors import InvalidStringData
from gridfs import GridFS
from bson import BSON

class mongostore(storage):
	def __init__(self, mongo_host="127.0.0.1", mongo_port=27017, mongo_db='canopsis', mongo_collection='perfdata', mongo_safe=False):
		storage.__init__(self)
		self.logger.debug(" + Init MongoDB Store")

		self.mongo_host = mongo_host
		self.mongo_port = mongo_port
		self.mongo_db = mongo_db
		self.mongo_collection = mongo_collection
		self.mongo_safe = mongo_safe

		self.logger.debug(" + Connect to MongoDB (%s/%s@%s:%s)" % (mongo_db, mongo_collection, mongo_host, mongo_port))
		self.conn=Connection(self.mongo_host, self.mongo_port)
		self.db=self.conn[self.mongo_db]
		self.collection = self.db[self.mongo_collection]

		self.grid = GridFS(self.db, self.mongo_collection+".fs")

	def drop(self):
		self.db.drop_collection(self.mongo_collection)
		self.db.drop_collection(self.mongo_collection+".fs.chunks")
		self.db.drop_collection(self.mongo_collection+".fs.files")

	def set_raw(self, key, value):
		try:
			self.collection.update({'_id': key}, {"$set": { 'd': value } }, upsert=True, safe=self.mongo_safe)
		except InvalidStringData:
			self.rm(key)
			self.grid.put(value, _id=key)
		except Exception, err:
			self.logger.error(err)
			self.logger.error(self.db.error())
			

	def set(self, key, value):
		self.logger.debug("Set '%s'" % key)
		self.set_raw(key, value)

	def get_raw(self, key):
		record = self.collection.find_one({'_id': key}, safe=self.mongo_safe)
		if record:
			return record['d']
		else:
			if self.grid.exists(key):
				return self.grid.get(key).read()
			else:
				return None
		
	def get(self, key):
		self.logger.debug("Get '%s'" % key)
		try:
			return self.get_raw(key)
		except:
			return None

	def rm(self, key):
		self.logger.debug("Remove '%s'" % key)
		if self.grid.exists(key):
			self.grid.delete(key)
		else:
			self.collection.remove(key, safe=self.mongo_safe)

	def append(self, key, value):
		self.logger.debug("Append data in '%s'" % key)
		self.logger.debug(" + Key: '%s'" % key)
		self.logger.debug(" + Value: '%s'" % value)
		
		try:
			self.collection.update({'_id': key}, { "$push": { 'd': value } }, upsert=True, safe=self.mongo_safe)
		except:
			self.set(key, [ value ])

	def size(self, key=None):
		size = 0
		if key:
			#TODO: Value is strange ...
			data = self.get(key)
			if type(data) == list or type(data) == dict:
				size = sys.getsizeof(BSON.encode({'_id': key, 'd': data }))
			else:
				size = sys.getsizeof(data)
			pass
			
		else:
			size = 0
			try:
				size = self.db.command("collstats", self.mongo_collection)['size']
				#print 'Col Size: %s' % size
			except:
				self.logger.warning("Impossible to read Collecion Size")
				
			try:
				chunks_size = self.db.command("collstats", self.mongo_collection+".fs.chunks")['size']
				#print 'Chunks Size: %s' % chunks_size
				bin_size = self.db.command("collstats", self.mongo_collection+".fs.files")['size']
				#print 'Bin Size: %s' % bin_size
				size += chunks_size + bin_size
			except:
				self.logger.warning("Impossible to read GridFS Size")
				pass
							
		return size

	def get_all_nodes(self,limit=None,offset=None,search=None):
		nodes = []
		
		filter = { 'd.metrics' : {'$exists' : True}}
		
		if search:
			filter = {'$and':[
								filter,
								{ 'd.dn': { '$regex' : '.*'+search+'.*', '$options': 'i' }}
							]}
		
		raw_output = self.collection.find(filter,sort =[('_id',ASCENDING)])
		total = raw_output.count()
		if raw_output and limit:
			raw_output = raw_output.limit(int(limit))
		if raw_output and offset:
			raw_output = raw_output.skip(int(offset))
		
		for record in raw_output:
			nodes.append({'node':record['_id'],'dn':record['d']['dn'],'metrics':record['d']['metrics']})

		return {'total':total,'data':nodes}
		
	def get_all_metrics(self,limit=None,offset=None,search=None):
		nodes = []
		
		filter = {'$and' : [
								{'d.dn':{'$exists' : True}},
								{ 'd.metrics' : {'$exists' : False}}
							]}
							
		if search:
			filter['$and'].append({ 'd.dn': { '$regex' : '.*'+search+'.*', '$options': 'i' }})
		
		raw_output = self.collection.find(filter, sort =[('_id',ASCENDING)])
		total = raw_output.count()
		if raw_output and limit:
			raw_output = raw_output.limit(int(limit))
		if raw_output and offset:
			raw_output = raw_output.skip(int(offset))
			
		for record in raw_output:
			nodes.append({'node':record['d']['node_id'],'metric':record['d']['dn']})
		
		return {'total':total,'data':nodes}
		
	def lock(self, key):
		self.logger.debug("Lock '%s'" % key)
		## Todo

	def wait_lock(self, key):
		self.logger.debug("Wait Lock '%s'" % key)
		## Todo	
