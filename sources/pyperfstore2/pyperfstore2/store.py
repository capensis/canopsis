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
from gridfs import GridFS, errors
import redis

import threading

class store(object):
	def __init__(self,
			mongo_host="127.0.0.1",
			mongo_port=27017,
			mongo_db='canopsis',
			mongo_collection='perfdata2',
			mongo_user=None,
			mongo_pass=None,
			mongo_safe=False,
			redis_host=None,
			redis_port=6379,
			redis_db=0,
			redis_sync_interval=10,
			logging_level=logging.INFO):

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

		self.redis_sync_interval = redis_sync_interval
		self.redis_db = redis_db
		self.redis_port = redis_port
		self.redis_host = redis_host

		if not redis_host:
			self.redis_host = self.mongo_host

		self.connected = False
		
		self.connect()

		self.last_sync = time.time()

		self.last_rate_time = time.time()
		self.rate_interval = 10
		self.rate_threshold = 20
		self.last_rate = 0
		self.pipe_size = 0
		self.pushed_values = 0

	def connect(self):
		if self.connected:
			self.logger.debug("Impossible to connect, already connected")
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

			self.redis = redis.StrictRedis(host=self.redis_host, port=self.redis_port, db=self.redis_db)
			self.redis_pipe = self.redis.pipeline()

			try:
				if self.mongo_user and self.mongo_pass != None:
						self.logger.debug("Try to auth '%s'" % self.mongo_user)
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
		if not self.connected or not self.connect():
			raise Exception('Impossible to deal with DB, you are not connected ...')
						
	def count(self, _id):
		self.check_connection()
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
	
	def sync(self):
		if self.connected:
			self.logger.debug("Sync pipeline to Redis")
			self.redis_pipe.execute()
			self.last_sync = time.time()
			self.pipe_size = 0

	def push(self, _id, point, meta_data={}, bulk=True):
		self.check_connection()
		self.logger.debug("Push point '%s' in '%s'" % (point, _id))
		
		meta_data['lts'] = point[0]
		meta_data['lv'] = point[1]

		# Update meta data on mongo
		if not self.redis.exists(_id):
			self.update(_id=_id, mset=meta_data)

		now = time.time()

		# Calcul push rate
		if self.pushed_values and (self.last_rate_time + self.rate_interval) < now:
			elapsed = now - self.last_rate_time
			self.last_rate =  self.pushed_values // elapsed
			self.pushed_values = 0
			self.last_rate_time = now

		# Disable bulk mode in lower rate
		if bulk and self.last_rate < self.rate_threshold:
			bulk = False
		
		if bulk and self.pipe_size == 0:
			self.logger.debug("Bulk mode is enabled (rate: %s push/sec)" % self.last_rate)

		# Push perfdata to db
		if not bulk and self.pipe_size == 0:
			self.redis.rpush(_id, '%s|%s' % (point[0], point[1]))
		else:
			self.redis_pipe.rpush(_id, '%s|%s' % (point[0], point[1]))
			self.pipe_size += 1

		# Sync DB if need
		if self.pipe_size and self.redis_sync_interval and (self.last_sync + self.redis_sync_interval) < now:
			self.sync()

		self.pushed_values += 1

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
		result = None
		self.check_connection()
		try:
			document = self.grid.get(_id)
			result = document.read()
		except errors.NoFile as nf:
			self.logger.error(nf)
		return result

	def find(self, limit=0, skip=0, mfilter={}, mfields=None, sort=None):
		self.check_connection()
		if limit == 1:
			result = self.collection.find_one(mfilter, limit=limit, fields=mfields, sort=sort)
		else:
			result = self.collection.find(mfilter, limit=limit, skip=skip, fields=mfields, sort=sort)
		
		return result

	def drop(self):
		self.check_connection()
		self.db.drop_collection(self.mongo_collection)
		self.db.drop_collection(self.mongo_collection+"_bin.chunks")
		self.db.drop_collection(self.mongo_collection+"_bin.files")
		self.redis.flushdb()
		
	def disconnect(self):
		# Sync redis
		self.sync()

		if self.connected:
			self.logger.debug("Disconnect from MongoDB")
			self.conn.fsync()
			del self.conn
			self.connected = False
		else:
			self.logger.warning("Impossible to disconnect, you are not connected")

