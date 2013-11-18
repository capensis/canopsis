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

import bottle
from bottle import route, get, delete, put, request
from bottle import HTTPError, post, static_file, response

logger = logging.getLogger('calendar')

from libexec.rest import *

namespace = "events"
#########################################################################

@get('/cal/:source/:interval_start/:interval_end')
def cal_get(source, interval_start, interval_end):
	try:
		print "====== cal_get"
		params = request.params

		filter = {
			"$and": [
				{"event_type" : "calendar"},
				{"component" : source},
				{"$or": [
					{"$and": [
								{"start": {"$gt": interval_start}},
								{"start": {"$lt": interval_end}}
							]},
					{"$or": [
								{"end": {"$gt": interval_start}},
								{"end": {"$lt": interval_end}}
							]}
				]}
			]
		}

		params['filter'] = json.dumps(filter)

		regular_events = rest_get("events", params=params)

		filter = {
			"$and": [
				{"event_type" : "calendar"},
				{"component" : source},
				{"rrule" : {"$exists": True}}
			]
		}

		params['filter'] = json.dumps(filter)

		recurrent_events = rest_get("events", params=params)

		return regular_events
	except Exception, e:
		print e
		raise
	else:
		pass
	finally:
		pass

