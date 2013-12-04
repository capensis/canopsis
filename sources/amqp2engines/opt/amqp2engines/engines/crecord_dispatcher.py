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

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		self.crecords = []
		self.delays = {}
		self.beat_interval = 5
		self.nb_beat = 0
		self.crecords_types = {
			'selector'		: {'ttl': 15, 'elapsed_time' : 0},
			'topology'		: {'ttl': 10, 'elapsed_time' : 0},
			'derogation'	: {'ttl': 5, 'elapsed_time' : 0},
			'consolidation'	: {'ttl': 20, 'elapsed_time' : 0}}
		#self.init_amqp()
		
	def init_amqp(self):	
		self.logger.debug('Initiating amqp connection to dispatchers')
		# Connection
		with Connection(hostname='127.0.0.1', userid='guest', virtual_host='canopsis') as conn:
			# Get one producer
			with producers[conn].acquire(block=True) as producer:
				self.amqp_producer = producer

	def pre_run(self):
		#load crecords from database
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.crecords = []


	def load_crecords(self):
		crecords_json = self.storage.find({'crecord_type': {'$in': self.crecords_types.keys() }, 'enable': True, 'loaded': False}, namespace="object")
		for	crecord_json in crecords_json:
			# let say selector is loaded
			self.storage.update(crecord_json._id, {'loaded': True})
			crecord = cselector(storage=self.storage, record=crecord_json, logging_level=self.logging_level)
			self.crecords.append(crecord)

			
			
	def beat(self):

		""" Reinitialize crecords and may publish event related credort targeted to other engines crecord queues"""
		self.load_crecords()
					
		# list of not yet sent crecords
		crecords_not_ready = []
		
		for crecord_type in self.crecords_types: 
			self.crecords_types[crecord_type]['elapsed_time'] += self.beat_interval
	
		# Loop until list is empty
		self.logger.debug(' + %s beat, %s crecords queued to publish' % (self.name, len(self.crecords)))
		while self.crecords:
			crecord = self.crecords.pop()		
			# Every crecord is sent to rabbit mq queues for each listening engines
			dump = crecord.dump()
			record_id = dump['_id']
			#Is this crecord ready to be threaten by engines ? if yes then it is sent to amqp and removed from local records
			if self.crecords_types[dump['crecord_type']]['elapsed_time'] > self.crecords_types[dump['crecord_type']]['ttl']:
				#crecord is sent to other engines and is not kept anymore
				try:
					self.storage.update(record_id, {'last_dispatch': time.time()})
					if '_id' in dump:
						# just sending key and type to build back object from dedicated engines
						dump['_id'] = str(dump['_id'])
						
						with Connection('amqp://guest:guest@localhost/canopsis') as conn:
							# Get one producer
							with producers[conn].acquire(block=True) as producer:
								# Publish crecord to dispatcher amqp queues

								rk = 'dispatcher.' + dump['crecord_type']
								exchange = Exchange('media', 'direct', durable=True)
								queue = Queue('Dispatcher_' + dump['crecord_type'],exchange=exchange, routing_key=rk)
							
								producer.publish(
									dump, 
									serializer='json', 
									exchange=exchange, 
									routing_key=rk, 
									declare=[queue])					
					
					self.logger.debug('publishing on queue : Dispatcher_' + dump['crecord_type'])
					#self.amqp.publish(dump, 'Dispatcher_' + self.name, "amq.direct")

					self.logger.debug("Published event for engine '%s'" % (dump['crecord_type']))
				except Exception, e:
					#Crecord gets out of queue and will be reloaded on next beat
					self.logger.error('Dispatcher was unable to send crecord_type error : %s' % (e))
					self.storage.update(record_id, {'loaded': False})



			else:
				crecords_not_ready.append(crecord)

		# New list with left items is kept
		self.crecords = crecords_not_ready

		for crecord_type in self.crecords_types: 
			if self.crecords_types[crecord_type]['elapsed_time'] > self.crecords_types[crecord_type]['ttl']:
				self.crecords_types[crecord_type]['elapsed_time'] = 0


		self.nb_beat +=1
