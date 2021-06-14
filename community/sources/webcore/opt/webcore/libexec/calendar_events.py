#!/usr/bin/env python
# -*- coding: utf-8 -*-
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

import sys
import os
import logging
import json
import gevent
from datetime import datetime
from dateutil.rrule import *
from time import mktime as mktime

import bottle
from bottle import route, get, delete, put, request
from bottle import HTTPError, post, static_file, response

logger = logging.getLogger('calendar')

from libexec.rest import *

namespace = "events"
#########################################################################

@get('/cal/:source/:interval_start/:interval_end')
def cal_get(source, interval_start, interval_end):
	params = request.params

	filter = {
		"$and": [
			{"event_type" : "calendar"},
			{"component" : source},
			{"rrule" : {"$exists": False}},
			{"$or": [
				{"$and": [
							{"start": {"$gt": int(interval_start)}},
							{"start": {"$lt": int(interval_end)}}
						]},
				{"$and": [
							{"end": {"$gt": int(interval_start)}},
							{"end": {"$lt": int(interval_end)}}
						]}
			]}
		]
	}

	params['filter'] = json.dumps(filter)

	events = rest_get("events", params=params)

	filter = {
		"$and": [
			{"event_type" : "calendar"},
			{"component" : source},
			{"rrule" : {"$exists": True}}
		]
	}

	params['filter'] = json.dumps(filter)

	recurrent_events = rest_get("events", params=params)
	#TODO install dateutil

	for event in recurrent_events["data"]:
		try:
			dtstart = datetime.fromtimestamp(float(interval_start))
			dtend = datetime.fromtimestamp(float(interval_end))

			eventStart = datetime.fromtimestamp(float(event["start"]))
			eventEnd = datetime.fromtimestamp(float(event["end"]))

			occurences = list(rrulestr(event["rrule"], dtstart=eventStart).between(dtstart, dtend))

			eventDuration = eventEnd - eventStart

			occurenceCount = 0
			#instantiate an event occurence for each found date
			for occurence in occurences:
				occurenceCount += 1
				newEvent = event.copy()
				occurenceStart = mktime(occurence.timetuple())
				occurenceEnd = occurence + eventDuration
				occurenceEnd = mktime(occurenceEnd.timetuple())
				newEvent["start"] = int(occurenceStart)
				newEvent["end"] = int(occurenceEnd)
				events["data"].append(newEvent)
		except Exception, e:
			print "Error parsing rrule for an event : %s" % e


	return events
