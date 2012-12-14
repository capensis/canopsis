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


import unittest
import time, json, logging
from bson import BSON

from camqp import camqp
from cinit import cinit
## Engines path
import sys, os
sys.path.append(os.path.expanduser('~/opt/amqp2engines/engines/'))

## Configurations

DAEMON_NAME="amqp2engines"

init 	= cinit()
logger 	= init.getLogger(DAEMON_NAME, level="INFO")
#logger 	= init.getLogger(DAEMON_NAME)
handler = init.getHandler(logger)

engines=[]
amqp = None
next_event_engines = []
next_alert_engines = []

def clean_message(body, msg):
	## Sanity Checks
	rk = msg.delivery_info['routing_key']
	if not rk:
		raise Exception("Invalid routing-key '%s' (%s)" % (rk, body))
		
	#logger.debug("Event: %s" % rk)
	#logger.info( body ) 	
	## Try to decode event
	if isinstance(body, dict):
		event = body
	else:
		logger.debug(" + Decode JSON")
		try:
			if isinstance(body, str) or isinstance(body, unicode):
				try:
					event = json.loads(body)
					logger.debug("   + Ok")
				except Exception, err:
					try:
						logger.debug(" + Try hack for windows string")
						# Hack for windows FS -_-
						event = json.loads(body.replace('\\', '\\\\'))
						logger.debug("   + Ok")
					except Exception, err :
						try:
							logger.debug(" + Decode BSON")
							bson = BSON (body)
							event = bson.decode()
							logger.debug("   + Ok")
						except Exception, err:
							raise Exception(err)
		except Exception, err:
			logger.error("   + Failed (%s)" % err)
			logger.debug("RK: '%s', Body:" % rk)
			logger.debug(body)
			raise Exception("Impossible to parse event '%s'" % rk)

	event['rk'] = rk
	
	# Clean tags field
	event['tags'] = event.get('tags', [])
	
	if (isinstance(event['tags'], str) or isinstance(event['tags'], unicode)) and  event['tags'] != "":
		event['tags'] = [ event['tags'] ]
		
	elif not isinstance(event['tags'], list):
		event['tags'] = []
		
	return event

def on_event(body, msg):
	## Clean message	
	event = clean_message(body, msg)
	
	event['exchange'] = amqp.exchange_name_events
	
	## Forward to engines
	for engine in next_event_engines:
		try:
			engine.input_queue.put(event)
		except Queue.Full:
			logger.warngin("Internal queue of '%s' is full, forward event to AMQP queue." % engine.name)
			amqp.publish(event, engine.amqp_queue, "amq.direct")

	
def on_alert(body, msg):	
	## Clean message	
	event = clean_message(body, msg)
	
	event['exchange'] = amqp.exchange_name_alerts
	
	## Forward to engines
	for engine in next_alert_engines:
		try:
			engine.input_queue.put(event)
		except Queue.Full:
			logger.warngin("Internal queue of '%s' is full, forward event to AMQP queue." % engine.name)
			amqp.publish(event, engine.amqp_queue, "amq.direct")

def start_engines():
	global engines
	# Init Engines
	## TODO: Use routing table for dynamic routing
	### Route:
	
	# Events:
	### Nagios/Icinga/Shinken... ----------------------------> canopsis.events -> tag -> perfstore -> eventstore
	### collectd ------------------> amq.topic -> collectdgw |
	
	# Alerts:
	### canopsis.alerts -> selector -> eventstore
	
	import perfstore2
	import eventstore
	import collectdgw
	import tag
	import selector
	import sla
	import alertcounter
	import derogation
	import topology
	import consolidation
	

	engine_selector		= selector.engine(logging_level=logging.INFO)
	engines.append(engine_selector)
	
	engine_topology		= topology.engine(next_engines=[engine_selector], logging_level=logging.INFO)
	engines.append(engine_topology)

	engine_alertcounter	= alertcounter.engine(next_engines=[engine_topology], logging_level=logging.INFO)
	engines.append(engine_alertcounter)
	
	engine_collectdgw	= collectdgw.engine()
	engines.append(engine_collectdgw)
	
	engine_eventstore	= eventstore.engine( logging_level=logging.INFO)
	engines.append(engine_eventstore)

	engine_perfstore	= perfstore2.engine(next_engines=[engine_eventstore])
	engines.append(engine_perfstore)
	
	engine_tag			= tag.engine(		next_engines=[engine_perfstore])
	engines.append(engine_tag)
	
	engine_derogation 	= derogation.engine( next_engines=[engine_tag], logging_level=logging.INFO)
	engines.append(engine_derogation)

	engine_sla			= sla.engine(logging_level=logging.INFO)
	engines.append(engine_sla)
	
	engine_consolidation		= consolidation.engine(logging_level=logging.DEBUG)
	engines.append(engine_consolidation)
	
	# Set Next queue
	## Events
	next_event_engines.append(engine_derogation)
	#next_event_engines.append(engine_tag)
	## Alerts
	next_alert_engines.append(engine_alertcounter)
	
	logger.info("Start engines")
	for engine in engines:
		engine.start()
	
def stop_engines():
	logger.info("Stop engines")
	for engine in engines:
		engine.signal_queue.put("STOP")
	
	logger.info("Join engines")
	for engine in engines:
		engine.join()
		while engine.is_alive():
			time.sleep(0.1)
			
	time.sleep(0.5)

def main():
	global amqp
		
	logger.info("Initialyze process")
	handler.run()

	logger.info("Start Engines")
	start_engines()

	# Safety wait
	time.sleep(3)
	
	# Init AMQP
	amqp = camqp(logging_name="%s-amqp" % DAEMON_NAME)
	amqp.add_queue(DAEMON_NAME, ['#'], on_event, amqp.exchange_name_events, auto_delete=False)
	amqp.add_queue("%s_alerts" % DAEMON_NAME, ['#'], on_alert, amqp.exchange_name_alerts, auto_delete=False)
	
	# Start AMQP
	amqp.start()
	
	logger.info("Wait")
	handler.wait()
	
	# Stop AMQP
	amqp.stop()
	amqp.join()
	
	stop_engines()

	logger.info("Process finished")
	
if __name__ == "__main__":
	main()
