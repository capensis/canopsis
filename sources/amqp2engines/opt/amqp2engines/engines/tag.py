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
import cmfilter
import json

NAME="tag"

class engine(cengine):

	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		self.nb_beat = 0		
			
	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))		
		self.reload_selectors()
		self.beat()

	def reload_selectors(self):
		
		"""
		Loads selectors crecords to find out witch event in work method match selector.
		When an event match, selector name and dispay name fields are added to the event tag list
		"""
		self.selectors = []
		self.selByRk = {}
		selectors_json = self.storage.find({'crecord_type': 'selector', 'enable': True }, namespace="object")
		
		for selector_json in selectors_json:
			selector_dump = selector_json.dump()
			# add selector name to tag when 
			if 'rk' in selector_dump:
				## Extract ids resolved by selectors
				sel = selector_dump.get('crecord_name')
				ids = selector_dump.get('ids', [])
				
				if isinstance(ids, list):
					for rk in ids:
						if rk in self.selByRk:
							self.selByRk[rk].append(sel)
						else:
							self.selByRk[rk] = [sel]
												
		
			# Put selector witch can tag event in cache
			if 'dostate' in selector_dump and selector_dump['dostate']:
				
				selector = cselector(storage=self.storage, record=selector_json)

				# tag field is defined here only
				selector.tags = []

				for selector_tag in ['crecord_name', 'display_name']:
					if  selector_tag in selector_dump and selector_dump[selector_tag]:
						selector.tags.append(selector_dump[selector_tag])
				self.selectors.append(selector)

	
		self.logger.debug('Reloaded %s selectors' % (len(self.selectors)))


	def add_tag(self, event, field=None, value=None):
		"""
		Adds a tag to event depending on values
		"""
		
		if not value and not field:
			return event
			
		if not value and field:
			value = event.get(field, None)
			
		if value and value not in event['tags']:
			event['tags'].append(value)
			
		return event
		
	def work(self, event, *args, **kargs):
	
		"""
		Each event comming to tag work method is beeing tagged with many informations.
		Those tags aim to display tags into UI to enhence information search.
		"""
		
		event['tags'] = event.get('tags', [])
		
		event = self.add_tag(event, 'connector_name')
		event = self.add_tag(event, 'event_type')
		event = self.add_tag(event, 'source_type')
		event = self.add_tag(event, 'component')
		event = self.add_tag(event, 'resource')

		# Adds tag to event if selector crecord matches current event.
		self.logger.debug('Will process selector tag on event %s ' % (event['rk']))		

		for selector in self.selectors:	

			add_tag = False
			cfilter = False
			self.logger.debug('Filter %s: type %s' % (selector.mfilter, type(selector.mfilter)) )
			if selector.mfilter:
				cfilter = cmfilter.check(selector.mfilter, event)					
				
			if 'rk' in event:
				if event['rk'] not in selector.exclude_ids and (event['rk'] in selector.include_ids or cfilter):
					add_tag = True
			elif cfilter:
				add_tag = True
			
			if add_tag:
				self.logger.debug('Will write tag to event: %s' % (selector.tags))
				for tag in selector.tags:
					if tag not in event['tags']:
						event['tags'].append(tag)

		### Tag with dynamic tags
		sels = self.selByRk.get(event['rk'], [])

		for sel in sels:
			event = self.add_tag(event, value=sel)
			
		return event

	def beat(self):

		"""
		Reload inforamtion allowing event tag in work method. 
		selectors are implied in tag definition
		"""

		self.nb_beat += 1

		self.logger.debug('Refresh selector records cache for event tag by selector purposes.')
		self.reload_selectors()


