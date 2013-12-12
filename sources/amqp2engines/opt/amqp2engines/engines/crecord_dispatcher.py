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

from cengine import cengine
from cstorage import get_storage
from caccount import caccount
from cselector import cselector
from kombu.pools import producers
from kombu import Connection, Queue, Exchange
import logging, time
		
NAME="crecord_dispatcher"
# Delay since the lock document is released in any cases
UNLOCK_DELAY = 60

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		self.crecords = []
		self.delays = {}
		self.beat_interval = 5
		self.nb_beat = 0
		self.crecords_types = ['selector', 'topology', 'derogation', 'consolidation']
		
	def init_amqp(self):	
		self.logger.debug('Initiating amqp connection to dispatchers')
		# Connection
		with Connection('amqp://guest:guest@localhost/canopsis') as conn:
			# Get one producer
			with producers[conn].acquire(block=True) as producer:
				self.amqp_producer = producer

	def pre_run(self):
		#load crecords from database
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))

		self.backend = self.storage.get_backend('object')
		

	def load_crecords(self):
		
		crecords = []		
		now = int(time.time())
		
		lock = self.backend.find_and_modify(
			query  = {'crecord_type': 'dispatcher_lock'}, 
			update = {'$set': {'lock': True}}, 
		)
				
		if lock:

			#init cases
			if 'lock' not in lock:
				lock['lock'] = False
			if 'last_update' not in lock: 
				lock['last_update'] = time.time()

			if lock['lock'] and now - lock['last_update'] < UNLOCK_DELAY:
				self.logger.debug('Information beeing threan by another dispatcher engine. Nothing has to be done. delay : %s' % (now - lock['last_update']))
				return
			
			elif now - lock['last_update'] > UNLOCK_DELAY:
				#Information message
				self.logger.debug('Dispatcher lock duration exeeded lock limit (%s seconds) , record selection will be performed' % (UNLOCK_DELAY))
				
			# Crecord selection can be performed now as lock is ready. Crecord are loaded again if last dispatch update is not set or < now - 60 seconds
			crecords_json = self.storage.find({
				'crecord_type': {'$in': self.crecords_types }, 
				'enable': True,
				'$or':[
					{'next_ready_time'		: {'$exists': False}},# crecord is new
					{'last_dispatch_update'	: {'$exists': False}},# crecord is new
					{'last_dispatch_update' : {'$lte': now - 60}},# unlock case
					{'$and': [					
						{'next_ready_time' 	: {'$lte': now}}, # record is ready
						{'loaded'			: False}
					]}
				]
			}, namespace="object")
			for	crecord_json in crecords_json:
				# let say selector is loaded
				self.storage.update(crecord_json._id, {'loaded': True, 'last_dispatch_update': now})
				crecord = cselector(storage=self.storage, record=crecord_json, logging_level=self.logging_level)
				crecords.append(crecord)

			#Updating lock status
			self.backend.update(
				{'crecord_type': 'dispatcher_lock'}, 
				{'$set': {'lock': False, 'last_update': time.time()}}
			)
		else:
			#New insert, no information being threaten this time
			self.backend.insert({'crecord_type': 'dispatcher_lock', 'last_update': time.time(), 'lock': False})
	
		return crecords

		
	#Factorised code method
	def publish_record(self, event, crecord_type, producer):

		rk = 'dispatcher.' + crecord_type
		exchange = Exchange('media', 'direct', durable=True)
		queue = Queue('Dispatcher_' + crecord_type, exchange=exchange, routing_key=rk)

		producer.publish(
			event, 
			serializer='json', 
			exchange=exchange, 
			routing_key=rk, 
			declare=[queue])					
		self.logger.debug('publishing on queue : Dispatcher_' + crecord_type)
						
			
			
	def beat(self):		

		""" Reinitialize crecords and may publish event related credort targeted to other engines crecord queues"""
		crecords = self.load_crecords()
						
		# Loop until list is empty
		self.logger.debug(' + %s beat, %s crecords queued to publish @ %s' % (self.name, len(crecords), int(time.time())))
				# Connection
		with Connection('amqp://guest:guest@localhost/canopsis') as conn:
			# Get one producer
			with producers[conn].acquire(block=True) as producer:
		
				for crecord in crecords:

					# Every crecord is sent to rabbit mq queues for each listening engines
					dump = crecord.dump()
					record_id = dump['_id']

					#crecord is sent to other engines and is not kept anymore
					if '_id' in dump:
						dump['_id'] = str(dump['_id'])
						try:
							# just sending key and type to build back object from dedicated engines
							self.publish_record(dump, dump['crecord_type'], producer)			
						
							#Special case: selector crecords targeted to SLA
							if dump['crecord_type'] == 'selector' and 'rk' in dump and dump['rk'] and 'dosla' in dump and dump['dosla'] in [ True, 'on'] and 'dostate' in dump and dump['dostate'] in [ True, 'on']:
								self.publish_record(dump, 'sla', producer)


						except Exception, e:
							#Crecord gets out of queue and will be reloaded on next beat
							self.logger.error('Dispatcher was unable to send crecord_type error : %s' % (e))
							self.storage.update(record_id, {'loaded': False})
		
						
		
		self.nb_beat +=1
