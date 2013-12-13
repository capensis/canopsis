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
		self.selByRk = {}
			
	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))		
		self.reload_selectors()
		self.beat()

	def reload_selectors(self):

		self.selectors = []
		selectorsjson = self.storage.find({'crecord_type': 'selector', 'enable': True, 'dostate': True}, namespace="object")

		for selectorjson in selectorsjson:
			selector = cselector(storage=self.storage, record=selectorjson, logging_level=self.logging_level)

			# Defined here only
			selector.tags = []
			dump = selectorjson.dump()
			for selector_tag in ['crecord_name', 'display_name']:
				if  selector_tag in dump and dump[selector_tag]:
					selector.tags.append(dump[selector_tag])
			self.selectors.append(selector)

			#Cache for work method is set in one atomic operation

		self.logger.debug('Reloaded %s selectors' % (len(self.selectors)))


	def add_tag(self, event, field=None, value=None):
		if not value and not field:
			return event
			
		if not value and field:
			value = event.get(field, None)
			
		if value and value not in event['tags']:
			event['tags'].append(value)
			
		return event
		
	def work(self, event, *args, **kargs):
	
		event['tags'] = event.get('tags', [])
		
		event = self.add_tag(event, 'connector_name')
		event = self.add_tag(event, 'event_type')
		event = self.add_tag(event, 'source_type')
		event = self.add_tag(event, 'component')
		event = self.add_tag(event, 'resource')

		self.logger.debug('Will process selector tag on event %s ' % (event['rk']))		
		for selector in self.selectors:	
			add_tag = False
			cfilter = False
			self.logger.debug('Super Filter %s: type %s' % (selector.mfilter, type(selector.mfilter)) )
			if selector.mfilter:
				cfilter = cmfilter.check(selector.mfilter, event)					
				self.logger.debug('cfilter result %s' % (cfilter) )
				
			if 'rk' in event:
				if event['rk'] not in selector.exclude_ids and (event['rk'] in selector.include_ids or cfilter):
					self.logger.debug('swag')
					add_tag = True
			elif cfilter:
				self.logger.debug('beauf')
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

		self.nb_beat += 1

		self.logger.debug('Refresh selector records cache for event tag by selector purposes.')
		self.reload_selectors()

	
		self.selByRk = {}
		
		## Extract ids resolved by selectors
		datas = self.storage.find({ 'crecord_type': 'selector', 'enable': True, 'rk': { '$exists' : True } }, mfields=['_id', 'crecord_name', 'ids'], namespace="object")
		for data in datas:
			_id = data.get('_id')
			sel = data.get('crecord_name')
			ids = data.get('ids', [])
				
			if isinstance(ids, list):
				for rk in ids:
					try:
						self.selByRk[rk].append(sel)
					except:
						self.selByRk[rk] = [ sel ]
												

