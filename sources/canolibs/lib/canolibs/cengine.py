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

import multiprocessing
import time
import Queue
import logging
import os, sys
from cinit import cinit
import traceback
import cevent

class cengine(multiprocessing.Process):

	def __init__(self, next_engines=[], name="worker1", beat_interval=60, use_internal_queue=True, queue_maxsize=1000, logging_level=logging.INFO):
		
		multiprocessing.Process.__init__(self)
		
		self.logging_level = logging_level
	
		self.signal_queue = multiprocessing.Queue(maxsize=5)
		self.input_queue = multiprocessing.Queue(maxsize=queue_maxsize)
		self.RUN = True
		
		self.name = name
		
		self.amqp_queue = "Engine_%s" % name
		
		self.perfdata_retention = 3600
		
		## Get from internal or external queue
		self.next_engines = next_engines
		
		init 	= cinit()
		
		self.logger = init.getLogger(name, logging_level=self.logging_level)
		
		# Log in file
		self.logger.addHandler(logging.FileHandler(filename=os.path.expanduser("~/var/log/engines/%s.log" % name)))	
		
		self.counter_error = 0
		self.counter_event = 0
		self.counter_worktime = 0
		
		self.thd_warn_sec_per_evt = 0.6
		self.thd_crit_sec_per_evt = 0.9
		
		self.beat_interval = beat_interval
		self.beat_last = time.time()
		
		self.create_queue =  True
		
		# If use_internal_queue is false, use only AMQP Queue
		self.use_internal_queue = use_internal_queue
		
		self.amqp_flow = True
		
		self.send_stats_event = True
				
		self.logger.info("Engine initialised")
		
	def create_amqp_queue(self):
		self.amqp.add_queue(self.amqp_queue, None, self.on_amqp_event, "amq.direct", no_ack=False, exclusive=False, auto_delete=False)
	
	def pre_run(self):
		pass
		
	def post_run(self):
		pass
	
	def run(self):
		def ready():
			self.logger.info(" + Ready!")
			
		self.logger.info("Start Engine with pid %s" % (os.getpid()))
		
		from camqp import camqp
		
		self.amqp = camqp(logging_level=logging.INFO, logging_name="%s-amqp" % self.name, on_ready=ready)
		
		if self.create_queue:
			self.create_amqp_queue()
		
		self.amqp.start()
		
		self.pre_run()
		
		while self.RUN:
			# Internal signals
			try:
				signal = self.signal_queue.get_nowait()
				self.logger.debug("Signal: %s" % signal)
				if signal == "STOP":
					self.RUN = False
			except Queue.Empty:
				pass
				
			# Beat
			if self.beat_interval:
				now = time.time()
				if now > (self.beat_last + self.beat_interval):
					self._beat()						
					self.beat_last = now
			
			# Input Queue
			if not self.input_queue.empty():
				while self.RUN :
					try:
						event = self.input_queue.get_nowait()
						self._work(event)
					except Queue.Empty:
						if self.amqp.paused and self.RUN:
							self.logger.info("Re-start AMQP Flow")
							self.amqp.paused = False
							
						break
			
			time.sleep(0.5)
		
		self.post_run()
		
		self.logger.info("Stop Engine")
		self.stop()
		self.logger.info("End of Engine")
		
	def on_amqp_event(self, event, msg):
		if self.use_internal_queue:
			try:
				if not self.input_queue.full():
					self.input_queue.put(event)
					msg.ack()
				else:
					if self.amqp_flow:
						self.logger.warning("Stop AMQP Flow")
						self.amqp.paused = True
					msg.requeue()
					
			except Exception, err:
				self.logger.error("Impossible to put event on internal queue (%s), requeue event." % err)
				msg.requeue()
				
		else:
			self._work(event)
			msg.ack()
				
	
	def _work(self, event, *args, **kargs):
		start = time.time()
		error = False
		try:
			wevent = self.work(event, *args, **kargs)
			# Forward event to next queue
			if self.next_engines:
				if wevent:
					#self.logger.debug("Forward event '%s' to next engines" % wevent['rk'])
					self.next_queue(wevent)
				else:
					#self.logger.debug("Forward original event '%s' to next engines" % event['rk'])
					self.next_queue(event)
					
		except Exception, err:
			error = True
			self.logger.error("Worker raise exception: %s" % err)
			traceback.print_exc(file=sys.stdout)
	
		if error:
			self.counter_error +=1
			
		self.counter_event += 1
		self.counter_worktime += time.time() - start
		
	def work(self, event, amqp_msg):
		return event
		
	def next_queue(self, event):
		for engine in self.next_engines:
			if not engine.input_queue.full() and self.use_internal_queue:
				#self.logger.debug(" + Forward via internal Q to '%s'" % engine.name)
				engine.input_queue.put(event)
			else:
				self.logger.debug(" + Forward via amqp to '%s'" % engine.amqp_queue)
				self.amqp.publish(event, engine.amqp_queue, "amq.direct")
		
	def _beat(self):
		self.logger.debug("Beat: %s event(s), %s error" % (self.counter_event, self.counter_error))
		
		if not self.input_queue.empty():
			size = self.input_queue.qsize()
			if size > 110:
				self.logger.info("%s event(s) in internal queue" % size)
			
		evt_per_sec = 0
		sec_per_evt = 0
		
		if self.counter_event:
			evt_per_sec = float(self.counter_event) / self.beat_interval
			self.logger.debug(" + %0.2f event(s)/seconds" % evt_per_sec)
		
		if self.counter_worktime and self.counter_event:
			sec_per_evt = self.counter_worktime / self.counter_event
			self.logger.debug(" + %0.5f seconds/event" % sec_per_evt)
		
		## Submit event
		if self.send_stats_event and self.counter_event != 0:
			state = 0
			
			if sec_per_evt > self.thd_warn_sec_per_evt:
				state = 1
				
			if sec_per_evt > self.thd_crit_sec_per_evt:
				state = 2
			
			perf_data_array = [
				{'retention': self.perfdata_retention, 'metric': 'cps_queue_size', 'value': self.input_queue.qsize(), 'unit': 'evt' },
				{'retention': self.perfdata_retention, 'metric': 'cps_evt_per_sec', 'value': round(evt_per_sec,2), 'unit': 'evt' },
				{'retention': self.perfdata_retention, 'metric': 'cps_sec_per_evt', 'value': round(sec_per_evt,5), 'unit': 's',
					'warn': self.thd_warn_sec_per_evt,
					'crit': self.thd_crit_sec_per_evt
				},
			]

			self.logger.debug(" + State: %s" % state)
			
			event = cevent.forger(
				connector = "cengine",
				connector_name = "engine",
				event_type = "check",
				source_type="resource",
				resource=self.amqp_queue,
				state=state,
				state_type=1,
				output="%0.2f evt/sec, %0.5f sec/evt" % (evt_per_sec, sec_per_evt),
				perf_data_array=perf_data_array
			)
			
			rk = cevent.get_routingkey(event)
			self.amqp.publish(event, rk, self.amqp.exchange_name_events)
		

		self.counter_error = 0
		self.counter_event = 0
		self.counter_worktime = 0

		try:
			self.beat()
		except Exception, err:
			self.logger.error("Beat raise exception: %s" % err)
			traceback.print_exc(file=sys.stdout)
				
	def beat(self):
		pass
			
	def stop(self):
		self.RUN = False
		
		# cancel self consumer
		self.amqp.cancel_queues()
		
		# transfer internal queue to AMQP queue
		if not self.input_queue.empty():
			self.logger.info("Transfer internal queue to AMQP queue")
			try:
				i =0
				while True:
					event = self.input_queue.get_nowait()
					if event:
						i+=1
						self.amqp.publish(event, self.amqp_queue, "amq.direct")
						
			except Queue.Empty:
				self.logger.info(" + %s event(s) transfered to AMQP" % i)
					
		self.amqp.stop()
		self.amqp.join()
		self.signal_queue.empty()
		del self.signal_queue
		self.logger.debug(" + Stopped")
