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

import kombu
from kombu import Connection, Exchange, Queue
import kombu.pools

try:
	from amqplib.client_0_8.exceptions import AMQPConnectionException as ConnectionError
except ImportError as IE:
	from amqp.exceptions import ConnectionError

import socket

import time, logging, threading, os, traceback
import sys
#from kombu.pools import producers


class camqp(threading.Thread):
	def __init__(self, host="localhost", port=5672, userid="guest", password="guest", virtual_host="canopsis", exchange_name="canopsis", logging_name="camqp", logging_level=logging.INFO, read_config_file=True, auto_connect=True, on_ready=None):
		threading.Thread.__init__(self)
		
		self.logger = logging.getLogger(logging_name)
		
		
		self.host=host
		self.port=port
		self.userid=userid
		self.password=password
		self.virtual_host=virtual_host
		self.exchange_name=exchange_name
		self.logging_level = logging_level
		
		if (read_config_file):
			self.read_config("amqp")
		
		self.amqp_uri = "amqp://%s:%s@%s:%s/%s" % (self.userid, self.password, self.host, self.port, self.virtual_host)
		
		self.logger.setLevel(logging_level)
		
		self.exchange_name_events=exchange_name+".events"
		self.exchange_name_alerts=exchange_name+".alerts"
		self.exchange_name_incidents=exchange_name+".incidents"
		
		self.chan = None
		self.conn = None
		self.connected = False
		self.on_ready = on_ready
		
		self.RUN = True
	
		self.exchanges = {}	
		
		self.queues = {}
		
		self.paused = False
		
		self.connection_errors = (	 ConnectionError,
									 socket.error,
									 IOError,
									 OSError,)
									 #AttributeError)
									 						 
		## create exchange
		self.logger.debug("Create exchanges object")
		for exchange_name in [self.exchange_name, self.exchange_name_events, self.exchange_name_alerts, self.exchange_name_incidents]:
			self.logger.debug(" + %s" % exchange_name)
			self.get_exchange(exchange_name)
		
		if auto_connect:
			self.connect()
		
		self.logger.debug("Object canamqp initialized")

	def run(self):
		self.logger.debug("Start thread ...")
		reconnect = False
		
		while self.RUN:
			
			self.connect()
			
			#self.wait_connection()
			
			if self.connected:
				self.init_queue(reconnect=reconnect)
				
				self.logger.debug("Drain events ...")
				while self.RUN:
					try:
						if not self.paused:
							self.conn.drain_events(timeout=0.5)
						else:
							time.sleep(0.5)

					except socket.timeout:
						pass

					except self.connection_errors as err:
						self.logger.error("Connection error ! (%s)" % err)
						break

					except Exception as err:
						self.logger.error("Unknown error: %s (%s)" % (err, type(err)))
						traceback.print_exc(file=sys.stdout)
						break
					
				self.disconnect()
		
			if self.RUN:
				self.logger.error("Connection loss, try to reconnect in few seconds ...")
				reconnect = True
				self.wait_connection(timeout=5)
			
		self.logger.debug("End of thread ...")
		
	def stop(self):
		self.logger.debug("Stop thread ...")
		self.RUN = False	

	def connect(self):
		if not self.connected:
			self.logger.info("Connect to AMQP Broker (%s:%s)" % (self.host, self.port))
			
			self.conn = Connection(self.amqp_uri)

			try:
				self.logger.debug(" + Connect")
				self.conn.connect()
				self.logger.info("Connected to AMQP Broker.")
				self.producers = kombu.pools.Producers(limit=10)
				self.connected = True
			except Exception as err:
				self.conn.release()
				self.logger.error("Impossible to connect (%s)" % err)
			
			if self.connected:
				self.logger.debug(" + Open channel")
				try:
					self.chan = self.conn.channel()
					
					self.logger.debug("Channel openned. Ready to send messages")
					
					try:
						## declare exchange
						self.logger.debug("Declare exchanges")
						for exchange_name in self.exchanges:
							self.logger.debug(" + %s" % exchange_name)
							self.exchanges[exchange_name](self.chan).declare()
					except Exception as err:
						self.logger.error("Impossible to declare exchange (%s)" % err)
					
				except Exception as err:
					self.logger.error(err)
		else:
			self.logger.debug("Allready connected")
	
	def get_exchange(self, name):
		if name:
			try:
				return self.exchanges[name]
			except:
				if name == "amq.direct":
					self.exchanges[name] = Exchange(name, "direct", durable=True)
				else:
					self.exchanges[name] =  Exchange(name , "topic", durable=True, auto_delete=False)
				return self.exchanges[name]
		else:
			return None
		
	def init_queue(self, reconnect=False):
		if self.queues:
			self.logger.debug("Init queues")
			for queue_name in self.queues.keys():
				self.logger.debug(" + %s" % queue_name)
				qsettings = self.queues[queue_name]
				
				if not qsettings['queue']:
					self.logger.debug("   + Create queue")

					# copy list
					routing_keys = list(qsettings['routing_keys'])
					routing_key = None

					if len(routing_keys):
						routing_key = routing_keys[0]
						routing_keys = routing_keys[1:]

					exchange = self.get_exchange(qsettings['exchange_name'])

					if (qsettings['exchange_name'] == "amq.direct" and not routing_key):
						routing_key = queue_name

					#self.logger.debug("   + exchange: '%s', routing_key: '%s', exclusive: %s, auto_delete: %s, no_ack: %s" % (qsettings['exchange_name'], routing_key, qsettings['exclusive'], qsettings['auto_delete'], qsettings['no_ack']))
					self.logger.debug("   + exchange: '%s', exclusive: %s, auto_delete: %s, no_ack: %s" % (qsettings['exchange_name'], qsettings['exclusive'], qsettings['auto_delete'], qsettings['no_ack']))
					qsettings['queue'] = Queue(queue_name,
											exchange = exchange,
											routing_key = routing_key,
											exclusive = qsettings['exclusive'],
											auto_delete = qsettings['auto_delete'],
											no_ack = qsettings['no_ack'],
											channel=self.conn.channel())

					qsettings['queue'].declare()

					if len(routing_keys):
						self.logger.debug("   + Bind on all routing keys")
						for routing_key in routing_keys:
							self.logger.debug("     + routing_key: '%s'" % routing_key)
							try:
								qsettings['queue'].bind_to(exchange=exchange, routing_key=routing_key)
							except:
								self.logger.error("You need upgrade your Kombu version (%s)" % kombu.__version__)

				if not qsettings['consumer']:
					self.logger.debug("   + Create Consumer")
					qsettings['consumer'] = self.conn.Consumer(qsettings['queue'], callbacks=[ qsettings['callback'] ])
				
				self.logger.debug("   + Consume queue")
				qsettings['consumer'].consume()
			
			if self.on_ready:
				self.on_ready()

	
	def add_queue(self, queue_name, routing_keys, callback, exchange_name=None, no_ack=True, exclusive=False, auto_delete=True):
		#if exchange_name == "amq.direct":
		#	routing_keys = queue_name

		c_routing_keys = []

		if not isinstance(routing_keys, list):
			if isinstance(routing_keys, str):
				c_routing_keys = [ routing_keys ]
		else:
			c_routing_keys = routing_keys
		
		if not exchange_name:
			exchange_name = self.exchange_name		
		
		self.queues[queue_name]={	'queue': False,
									'consumer': False,
									'queue_name': queue_name,
									'routing_keys': c_routing_keys,
									'callback': callback,
									'exchange_name': exchange_name,
									'no_ack': no_ack,
									'exclusive': exclusive,
									'auto_delete': auto_delete
							}


	def publish(self, msg, routing_key, exchange_name=None, serializer="json", compression=None, content_type=None, content_encoding=None):
		self.wait_connection()
		if self.connected:
			if not exchange_name:
				exchange_name = self.exchange_name
				
			self.logger.debug("Send message to %s in %s" % (routing_key, exchange_name))
			with self.producers[self.conn].acquire(block=True) as producer:
				try:
					_msg = msg.copy()
					camqp._clean_msg_for_serialization(_msg)
					producer.publish(_msg, serializer=serializer, compression=compression, routing_key=routing_key, exchange=self.get_exchange(exchange_name))
					self.logger.debug(" + Sended")
				except Exception, err:
					self.logger.error(" + Impossible to send (%s)" % err)
		else:
			self.logger.error("You are not connected ...")

	@staticmethod
	def _clean_msg_for_serialization(msg):
		from bson import objectid
		for key in msg:
			if isinstance(msg[key], objectid.ObjectId):
				msg[key] = str(msg[key])

	def cancel_queues(self):
		if self.connected:
			for queue_name in self.queues.keys():
				if self.queues[queue_name]['consumer']:
					self.logger.debug(" + Cancel consumer on %s" % queue_name)
					try:
						self.queues[queue_name]['consumer'].cancel()
					except:
						pass
						
					del(self.queues[queue_name]['consumer'])
					self.queues[queue_name]['consumer'] = False
					del(self.queues[queue_name]['queue'])
					self.queues[queue_name]['queue'] = False		
	
	def disconnect(self):
 		if self.connected:
			self.logger.info("Disconnect from AMQP Broker")
	
			self.cancel_queues()

			for exchange in self.exchanges:
				del exchange
			self.exchanges = {}

			try:
				kombu.pools.reset()
			except Exception as err:
				self.logger.error("Impossible to reset kombu pools: %s (%s)" % (err, type(err)))

			try:
				self.conn.release()
				del self.conn
			except Exception as err:
				self.logger.error("Impossible to release connection: %s (%s)" % (err, type(err)))

			self.connected = False

	def wait_connection(self, timeout=5):
		i=0
		while self.RUN and not self.connected and i < (timeout*2):
			try:
				time.sleep(0.5)
			except:
				pass
			i+=1

	def read_config(self, name):

		filename = '~/etc/' + name + '.conf'
		filename = os.path.expanduser(filename)

		import ConfigParser
		self.config = ConfigParser.RawConfigParser()

		try:
			self.config.read(filename)

			section = 'master'

			self.host = self.config.get(section, "host")
			self.port = self.config.getint(section, "port")
			self.userid = self.config.get(section, "userid")
			self.password = self.config.get(section, "password")
			self.virtual_host = self.config.get(section, "virtual_host")
			self.exchange_name = self.config.get(section, "exchange_name")

		except Exception, err:
			self.logger.error("Impossible to load configurations (%s), use default ..." % err)
			
	def __del__(self):
		self.stop()
