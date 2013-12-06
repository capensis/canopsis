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

import logging
		
NAME="selector"

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		self.selectors = []
		self.selector_refresh = {}
		self.beat_interval = 5
		self.nb_beat_publish = 10
		self.nb_beat = 0
		#self.beat_interval = 5
		self.thd_warn_sec_per_evt = 1.5
		self.thd_crit_sec_per_evt = 2
		
	
	def pre_run(self):
		self.beat_lock = False
		#load selectors
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.reload_selectors()



	def reload_selectors(self):
		self.selectors = []
		selectorsjson = self.storage.find({'crecord_type': 'selector', 'enable': True}, namespace="object")
		
		for	selectorjson in selectorsjson:
			# let say selector is loaded
			self.storage.update(selectorjson._id, {'loaded': True})
			selector = cselector(storage=self.storage, record=selectorjson, logging_level=self.logging_level)
			self.selectors.append(selector)

			
			
	def beat(self):

		if self.beat_lock:
			return 
			
		self.beat_lock = True
		""" Reinitialize selectors and may publish event if they have to"""
		self.reload_selectors()

		publish = False
		if self.nb_beat >= self.nb_beat_publish:
			self.nb_beat = 0
			publish = True
		
		for selector in self.selectors:
			
			# do I publish a selector event ? Yes if selector have to and it is time or we got to update status 
			if selector.dostate and (publish or (selector._id in self.selector_refresh and self.selector_refresh[selector._id])):
				try:
					#TODO improve this full mongo db request
					rk, event = selector.event()
				except Exception as e:
					self.logger.error({'msg': 'unable to select all event matching this selector in order to publish worst state one form them','exception':e})
					event = None
				

					
				if event:
					# Publish Sla information when available
					publishSla = selector.data.get('sla_rk', None)
					if publishSla:
						event['sla_rk'] = publishSla
											
					# Ok then i have to update selector statement
					self.storage.update(selector._id, {'state': event['state']})
					self.amqp.publish(event, rk, self.amqp.exchange_name_events)
					self.logger.debug("Published event for selector '%s'" % (selector.name))
					
				self.selector_refresh[selector._id] = False

		self.nb_beat +=1
		self.beat_lock = False
		
	def work(self, event, *args, **kargs):
						
		## Process selector and prevent Burst
		for selector in self.selectors:
			if selector.dostate and selector.match(event):
				self.selector_refresh[selector._id] = True
							
		return event
		
	def post_run(self):
		for selector in self.selectors:
			self.storage.update(selector._id, {'loaded': False})
		self.selector = []
