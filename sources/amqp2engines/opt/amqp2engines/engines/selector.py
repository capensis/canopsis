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
		self.thd_warn_sec_per_evt = 1.5
		self.thd_crit_sec_per_evt = 2
		
	

	
	def pre_run(self):
		#load selectors
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		self.reload_selectors()



	def reload_selectors(self):
		selectors = []
		selectorsjson = self.storage.find({'crecord_type': 'selector', 'enable': True}, namespace="object")
		
		for	selectorjson in selectorsjson:
			selector = cselector(storage=self.storage, record=selectorjson, logging_level=self.logging_level)
			selectors.append(selector)
		#Cache for work method is set in one atomic operation
		self.selectors = selectors



			
	def beat(self):
		self.logger.debug('entered in selector BEAT')
		# Refresh selectors for work method
		self.nb_beat += 1
		if self.nb_beat % 10 == 0:
			self.logger.debug('Refresh selector records cache.')
			self.reload_selectors()
		
	
	def consume_dispatcher(self,  event, *args, **kargs):
		self.logger.debug('entered in selector consume dispatcher')
		# Gets crecord from amqp distribution
		selector = self.get_ready_record(event)
		if selector:
			event_id = event['_id']
			# Loads associated class
			selector = cselector(storage=self.storage, record=selector, logging_level=self.logging_level)
			self.logger.debug('selector found, start processing..')			
			# do I publish a selector event ? Yes if selector have to and it is time or we got to update status 
			if selector.dostate:
				try:
					#TODO improve this full mongo db request
					rk, event = selector.event()
				except Exception as e:
					self.logger.error({'msg': 'unable to select all event matching this selector in order to publish worst state one form them','exception':e})
					event = None
				
					# Publish Sla information when available
					publishSla = selector.data.get('sla_rk', None)
					if publishSla:
						event['sla_rk'] = publishSla
											
					# Ok then i have to update selector statement
					self.storage.update(selector._id, {'state': event['state']})
					self.amqp.publish(event, rk, self.amqp.exchange_name_events)
					self.logger.debug("Published event for selector '%s'" % (selector.name))
					
				self.selector_refresh[selector._id] = False
			else:
				self.logger.debug('Nothing to do with this selector')
		self.nb_beat +=1
		#set record free for dispatcher engine
		self.storage.update(event_id, {'loaded': False})
		
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
