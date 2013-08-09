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
		self.selectors = {}
		self.selectors_events = {}
		
		self.nb_beat_publish = 10
		self.nb_beat = 0

		self.thd_warn_sec_per_evt = 1.5
		self.thd_crit_sec_per_evt = 2
	
	def pre_run(self):
		#load selectors
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
		
		self.unload_all_selectors()
		self.load_selectors()

	def clean_selectors(self):
		## check if selector is already in store
		id_to_clean = []
		ids = [_id for _id in self.selectors]
		
		count = self.storage.count({'_id': {"$in": ids}}, namespace="object")
		if count != len(ids):
			for _id in self.selectors:
				if not self.storage.count({'_id': _id, 'enable': True}, namespace="object"):
					id_to_clean.append(_id)
				
			for _id in id_to_clean:
				self.logger.debug("Clean selector %s: %s" % (_id, self.selectors[_id].name))
				try:
					self.storage.update(_id, {'loaded': False})
				except:
					self.logger.debug(" + Record are deleted.")
					
				del self.selectors[_id]
				del self.selectors_events[_id]
	
	def unload_selectors(self):
		self.clean_selectors()
		
		## Unload selectors
		if self.selectors:
			for _id in self.selectors:
				selector = self.selectors[_id]
				record = self.storage.get(selector._id)
				self.logger.debug("Unload selector %s: %s" % (record._id, record.name))
				self.storage.update(record._id, {'loaded': False})

			del self.selectors
	
		self.selectors = []
		
	def unload_all_selectors(self):
		records = self.storage.find({'crecord_type': 'selector'}, namespace="object")
		
		for record in records:
			self.storage.update(record._id, {'loaded': False})
	
	def load_selectors(self):
		## Load selectors
		self.clean_selectors()
		
		## New selector or modified selector
		records = self.storage.find({'crecord_type': 'selector', 'loaded': False, 'enable': True}, namespace="object")
		
		for record in records:
			self.logger.debug("Load selector %s: %s" % (record._id, record.name))
			_id = record._id
			
			# Set loaded
			self.storage.update(_id, {'loaded': True})
			
			# Del old
			if self.selectors.get(_id, None):
				del self.selectors[_id]
				del self.selectors_events[_id]
				
			## store
			self.selectors[_id] = cselector(storage=self.storage, record=record, logging_level=self.logging_level)
			self.selectors_events[_id] = 0
			
			## Publish state
			if self.selectors[_id].dostate:
				(rk, event) = self.selectors[_id].event()
				if event:
					# Set State
					self.storage.update(_id, {'state': event['state']})
					# Publish
					self.amqp.publish(event, rk, self.amqp.exchange_name_events)
		

	def beat(self):
		self.load_selectors()
		
		publish = False
		if self.nb_beat >= self.nb_beat_publish:
			self.nb_beat = 0
			publish = True
		
		for _id in self.selectors:
			selector = self.selectors[_id]
			
			if not selector.dostate:
				# Dont send state but resolve ids for cache result (tags)
				selector.resolv()
			else:
				if self.selectors_events[_id]:
					publish = True

				if publish:
					(rk, event) = selector.event()
					if event:
						if selector.data.get('sla_rk', None):
							event['sla_rk'] = selector.data.get('sla_rk')

						self.logger.debug("Publish event for '%s' (%s events)" % (selector.name, self.selectors_events[_id]))
						self.storage.update(_id, {'state': event['state']})
						self.amqp.publish(event, rk, self.amqp.exchange_name_events)
						
					self.selectors_events[_id] = 0
						
		self.nb_beat +=1
				
	def work(self, event, *args, **kargs):
		event_id = event["event_id"]
		
		selectors_to_delete = []
		
		## Process selector and prevent Burst
		for sid in self.selectors:
			selector = self.selectors[sid]
			if selector.dostate:
				try:
					if selector.match(event_id):
						self.selectors_events[sid] += 1
				except:
					self.logger.debug("%s wasn't found, it will be erase" % str(sid))
					selectors_to_delete.append(sid)
					
		## delete selectors not found
		for sid in selectors_to_delete:
			del self.selectors[sid]
		
		return event
		
	def post_run(self):
		self.unload_selectors()
