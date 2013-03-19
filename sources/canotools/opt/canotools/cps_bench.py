#!/usr/bin/env python
# --------------------------------
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

import time
import logging
import sys

from camqp import camqp
import cevent
from cstorage import get_storage
from caccount import caccount
import traceback

########################################################
#
#   Configuration
#
########################################################

AMQP_HOST = "localhost"

logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s %(name)s %(levelname)s %(message)s',
                    )

logger = logging.getLogger("bench")
amqp = camqp()

storage = get_storage(namespace='events', account=caccount(user="root", group="root"))

base_component_event = cevent.forger(
					connector =			'bench',
					connector_name =	"engine",
					event_type =		"check",
					source_type =		"component",
					component =			"component-",
					state =				0,
					state_type =		1,
					output =			"Output",
					long_output =		"",
					#perf_data =			None,
					#perf_data_array =	[
					#	{'metric': 'shortterm', 'value': 0.25, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' },
					#	{'metric': 'midterm',   'value': 0.16, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' },
					#	{'metric': 'longterm',  'value': 0.12, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' }
                    #]
					#display_name =		""
				)

base_resource_event = cevent.forger(
					connector =			'bench',
					connector_name =	"engine",
					event_type =		"check",
					source_type =		"resource",
					component =			"component-",
					resource =			"resource-",
					state =				0,
					state_type =		1,
					output =			"Output",
					long_output =		"",
					#perf_data =			None,
					#perf_data_array =	[],
					#display_name =		""
				)
				

########################################################
#
#   Functions
#
########################################################

#### Connect signals
RUN = 1
import signal
def signal_handler(signum, frame):
	logger.warning("Receive signal to stop daemon...")
	global RUN
	RUN = 0
 
signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)

def send_events(n, rate=0, burst=10):
	i = 0
	
	logger.info("Send %s events" % n)
	if (rate):
		#1/10 ou 1/15
		logger.info(" + @ %s events/second (%s events/5min)" % (rate, (rate*300)))
		#time_target = time.time() + float(n)/rate
		time_break = ((float(n) / rate) / n) * burst
		#logger.info(" + @ %s events / %s seconds" % (rate, burst))
		#logger.info(" + sleep %s seconds / %s events" % (time_break, burst))

	time_start_burst = time.time()
	start_time = time.time()
	while RUN and i < n:
		event = base_component_event.copy()

		event['component'] += str(i)
		event['bench_timestamp'] = time.time()

		rk = cevent.get_routingkey(event)
		amqp.publish(event, rk, amqp.exchange_name_events)

		if (rate and (i % burst == 0)):
			elapsed = time.time() - time_start_burst
			if (time_break > elapsed):
				time.sleep(time_break - elapsed)
			time_start_burst = time.time()

		i+=1

	duration = time.time() - start_time
	logger.info(" + Done, elapsed: %.3f ms (%s events/second)" % ((duration*1000), int(n/duration)) )

	# Get last event
	record = None
	elapsed = None
	logger.info("Wait last record ...")
	timeout = time.time() + 300
	while RUN:
		raw = storage.find_one({'_id': rk}, mfields={'bench_timestamp': 1})
		if raw:
			elapsed = time.time() - float(raw['bench_timestamp'])
			storage.get_backend('events').remove({'_id': rk}, safe=True)
			logger.info(" + Done, Delta: %.3f s" % elapsed )
			break
		
		if time.time() > timeout:
			logger.info(" + Fail, timeout")
			break

		time.sleep(0.001)

	return elapsed

def clean_db():
	# Clean DB
	logger.info("Remove old data")
	storage.get_backend('perfdata2').remove({'co': {'$regex': 'component-.*'}}, safe=True)
	storage.get_backend('events').remove({'connector': 'bench'}, safe=True)
	storage.get_backend('events_log').remove({'connector': 'bench'}, safe=True)
	time.sleep(1)
	if (storage.get_backend('events').find({'connector': 'bench'}).count()):
		logger.error(" + All data are not removed ...")
	else:	
		logger.info(" + Done")

def mean(values) :
	return sum(values) / float(len(values))

def stat_variance( echantillon ) :
	n = len( echantillon ) # taille
	mq = mean( echantillon )**2
	s = sum( [ x**2 for x in echantillon ] )
	variance = s / n - mq
	return variance

########################################################
#
#   Main
#
########################################################


amqp.start()

clean_db()
stats = []

try:
	# Send n events and check lattency
	#nbs = [ 500, 1000, 1250, 1500, 1750, 2000, 2250, 2500, 2750, 3000, 3250, 3500, 3750, 4000 ]
	#rates =  [100, 150, 200, 250, 300, 350, 400, 450, 500]
	nbs = [ 4000 ]
	rates = [ 200 ]
	
	for rate in rates:
		for nb in nbs:
			result = send_events(nb, rate)
			stats.append((nb, rate, result))
			time.sleep(5)
			#clean_db()

	clean_db()

except Exception as err:
	logger.error('Bench Failed !')
	logger.error(err)
	traceback.print_exc(file=sys.stdout)

times = [stat[2] for stat in stats]
print "Stats:"
print " + Mean: 	%.2f" % mean(times)
print " + Variance:	%.2f" % stat_variance(times)

amqp.stop()
amqp.join()