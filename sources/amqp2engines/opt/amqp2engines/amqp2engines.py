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

from cinit import cinit

## Engines path
import sys, os
sys.path.append(os.path.expanduser('~/opt/amqp2engines/engines/'))

## Configurations
DAEMON_NAME="amqp2engines"

init 	= cinit()
logger 	= init.getLogger(DAEMON_NAME, level="INFO")
handler = init.getHandler(logger)

engines=[]

## Very Dirty HACK !
## Remove old queues (temporary workaround)
import ConfigParser
amqp_config = ConfigParser.RawConfigParser()
section = 'master'
amqp_config.read(os.path.expanduser("~/etc/amqp.conf"))
amqp_host = amqp_config.get(section, "host")
amqp_port = amqp_config.getint(section, "port")
amqp_userid = amqp_config.get(section, "userid")
amqp_password = amqp_config.get(section, "password")
amqp_virtual_host = amqp_config.get(section, "virtual_host")

import subprocess
subprocess.call('rabbitmqadmin -H %s --vhost=%s --username=%s --password=%s delete queue name="amqp2engines"'
	% (amqp_host, amqp_virtual_host, amqp_userid, amqp_password), shell=True)
subprocess.call('rabbitmqadmin -H %s --vhost=%s --username=%s --password=%s delete queue name="amqp2engines_alerts"'
	% (amqp_host, amqp_virtual_host, amqp_userid, amqp_password), shell=True)

###### END of HACK ####


def start_engines():
	global engines

	##################
	# Events
	##################

	# Engine_cleaner
	import cleaner
	engines.append( cleaner.engine(			next_amqp_queues=['Engine_derogation'], routing_keys=["#"], exchange_name="canopsis.events", name='cleaner_events'))

	# Engine_derogation
	import derogation
	engines.append( derogation.engine(		next_amqp_queues=['Engine_tag']) )

	# Engine_tag
	import tag
	engines.append( tag.engine(				next_amqp_queues=['Engine_perfstore2']) )

	# Engine_perfstore2
	import perfstore2
	engines.append( perfstore2.engine(		next_amqp_queues=['Engine_eventstore']) )

	# Engine_eventstore
	import eventstore
	engines.append( eventstore.engine() )


	##################
	# Alerts
	##################

	# Engine_cleaner
	import cleaner
	engines.append( cleaner.engine(			next_amqp_queues=['Engine_alertcounter'], routing_keys=["#"], exchange_name="canopsis.alerts", name='cleaner_alerts'))

	# Engine_alertcounter
	import alertcounter
	engines.append( alertcounter.engine(	next_amqp_queues=['Engine_topology']),)

	# Engine_topology
	import topology
	engines.append( topology.engine(		next_amqp_queues=['Engine_selector']) )

	# Engine_selector
	import selector
	engines.append( selector.engine() )


	##################
	# Autres
	##################

	# Engine_collectdgw
	import collectdgw
	engines.append( collectdgw.engine() )

	# Engine_sla (no queue)
	import sla
	engines.append( sla.engine() )

	# Engine_consolidation
	import consolidation
	engines.append( consolidation.engine(logging_level=logging.DEBUG) )	


	##################
	# Start engines
	##################

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
	
	logger.info("Waitting ...")
	handler.wait()
	
	stop_engines()

	logger.info("Process finished")
	
if __name__ == "__main__":
	main()
