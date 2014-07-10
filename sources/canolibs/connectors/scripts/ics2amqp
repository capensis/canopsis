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

from pyparsing import Word, alphas, Suppress, nums, Optional, Regex
import os

from camqp import camqp

from cinit import cinit
import urllib
from icalendar import Calendar
import pytz
import ConfigParser
import pickle

import time
from datetime import datetime
import calendar

import cevent

DAEMON_NAME='ics2amqp'
config_file=os.path.expanduser('~/etc/ics2amqp.conf')
timestamps_file=os.path.expanduser('~/var/lib/ics2amqp/ics2amqp.timestamps')


init 	= cinit()
logger 	= init.getLogger(DAEMON_NAME)
handler = init.getHandler(logger)

sleep_time = 1

ics_streams = []
myamqp = None

## Init parser
integer = Word(nums)
serverDateTime = Regex("\S\S\S\s*\d\d?\s*\d\d:\d\d:\d\d")
hostname = Word(alphas + nums + "_" + "-")
daemon = Word(alphas + "/" + "-" + "_") + Optional(Suppress("[") + integer + Suppress("]")) + Suppress(":")
output = Regex(".*")
syslog_parser = serverDateTime + hostname + daemon + output

########################################################
#
#   Functions
#
########################################################

sources = []


utc=pytz.UTC

lastUpdate = {}

count = 0

def read_config():
	global config_file
	try:
		logger.info(" + reading config")
		config = ConfigParser.RawConfigParser()
		#TODO add this file in canopsis install
		config.read(config_file)

		for source_name, source_url in config.items('sources'):
			sources.append({"name" : source_name, "url" : source_url})

	except Exception, e:
		logger.error(e)


def init_timestamps():
	global lastUpdate
	global timestamps_file
	try:
		timestamps_input = open(timestamps_file, 'r')
		lastUpdate = pickle.load(timestamps_input)
		timestamps_input.close()
	except:
		logger.warn("No timestamps file found.")
		lastUpdate = {}


def update_timestamps():
	global timestamps_file
	global lastUpdate

	if not os.path.exists(os.path.dirname(timestamps_file)):
		os.makedirs(os.path.dirname(timestamps_file))

	output = open(timestamps_file, 'w')
	p = pickle.Pickler(output)
	p.dump(lastUpdate)
	output.close()

def send_event(source, event):
	#state should be info
	state = 0

	try:
		timestamp = calendar.timegm(event.get("last-modified").dt)
	except:
		timestamp = None

	component = source["name"]
	resource = event.get('uid')

	output = event.get('summary')
	long_output = event.get('description')

	start = str(event.get("dtstart").dt)
	end = str(event.get("dtend").dt)

	if type(start) is datetime:
		all_day = False
	else:
		all_day = True

	# start = time.mktime(start.timetuple())
	# end = time.mktime(end.timetuple())

	source_type='resource'

	event = cevent.forger(
					connector='ics',
					connector_name=DAEMON_NAME,
					component=component,
					resource=resource,
					timestamp=timestamp,
					source_type=source_type,
					event_type='calendar',
					state=state,
					output=output,
					long_output=long_output)

	event["start"] = start
	event["end"] = end

	event["all_day"] = all_day

	logger.debug('Event: %s' % event)

	key = cevent.get_routingkey(event)
	myamqp.publish(event, key, myamqp.exchange_name_events)

def parse_ics(source):
	global count
	global lastUpdate
	global timestamps_file

	ics = urllib.urlopen(source["url"]).read()
			# events = []

	cal = Calendar.from_ical(ics)

	if not(source["name"] in lastUpdate):
		lastUpdate[source["name"]] = utc.localize(datetime(2000, 1, 1))

	for event in cal.walk('vevent'):
		if event.get('LAST-MODIFIED') and event.get('LAST-MODIFIED').dt > lastUpdate[source["name"]]:
			send_event(source, event)

			count = count + 1

	lastUpdate[source["name"]] = event.get('dtstamp').dt

	update_timestamps()

	if count > 0 : print "send %i events from source %s" % (count, source["name"])

	count = 0

def wait_ics():
	global sources
	# try:
	logger.info("starting to poll ics files ")
	while handler.status():
		for source in sources:
			# logger.info("parsing source : " + str(source))
			parse_ics(source)

		time.sleep(sleep_time)


	# except Exception, err:
	# 	logger.error("Exception: '%s'" % err)

	logger.info("Close ICS handler")


########################################################
#
#   Main
#
########################################################

def main():

	handler.run()

	# global
	global myamqp
	read_config()
	init_timestamps()

	# Connect to amqp bus
	logger.info("Starting AMQP ...")
	myamqp = camqp()
	myamqp.start()

	logger.info("Started AMQP ...")

	wait_ics()

	logger.debug("Stop AMQP ...")
	myamqp.stop()
	myamqp.join()

if __name__ == "__main__":
	main()
