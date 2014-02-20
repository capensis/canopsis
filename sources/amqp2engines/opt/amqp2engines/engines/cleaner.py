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
from bson import BSON
import json
import time

NAME="cleaner"

class engine(cengine):
	def __init__(self, name=NAME, *args, **kargs):
		cengine.__init__(self, name=name, *args, **kargs)
		
	def work(self, body, msg, *args, **kargs):
		## Sanity Checks
		rk = msg.delivery_info['routing_key']
		if not rk:
			raise Exception("Invalid routing-key '%s' (%s)" % (rk, body))
		
		#self.logger.info( body ) 	
		## Try to decode event
		if isinstance(body, dict):
			event = body
		else:
			self.logger.debug(" + Decode JSON")
			try:
				if isinstance(body, str) or isinstance(body, unicode):
					try:
						event = json.loads(body)
						self.logger.debug("   + Ok")
					except Exception, err:
						try:
							self.logger.debug(" + Try hack for windows string")
							# Hack for windows FS -_-
							event = json.loads(body.replace('\\', '\\\\'))
							self.logger.debug("   + Ok")
						except Exception, err :
							try:
								self.logger.debug(" + Decode BSON")
								bson = BSON (body)
								event = bson.decode()
								self.logger.debug("   + Ok")
							except Exception, err:
								raise Exception(err)
			except Exception, err:
				self.logger.error("   + Failed (%s)" % err)
				self.logger.debug("RK: '%s', Body:" % rk)
				self.logger.debug(body)
				raise Exception("Impossible to parse event '%s'" % rk)

		event['rk'] = rk
		
		# Clean tags field
		event['tags'] = event.get('tags', [])
		
		if (isinstance(event['tags'], str) or isinstance(event['tags'], unicode)) and  event['tags'] != "":
			event['tags'] = [ event['tags'] ]
			
		elif not isinstance(event['tags'], list):
			event['tags'] = []

		event["timestamp"] 	= event.get("timestamp", time.time() )
		event["timestamp"] 	= int(event["timestamp"])

		event["state"]		= event.get("state", 0)
		event["state_type"] = event.get("state_type", 1)
		event["event_type"] = event.get("event_type", "check")

		if event['event_type'] == 'check':
			event['component_problem'] = event.get('component_problem', False)

		return event
