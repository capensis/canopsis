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

NAME="tag"

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		
		self.tags_ids = []
		
	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))		
		self.beat()
				
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
		
		### Tag with dynamic tags
		if self.tags_ids:
			for data in self.tags_ids:
				if event['rk'] in data[1]:
					event = self.add_tag(event, value=data[0])
			
		return event

	def beat(self):
		self.tags_ids = []
		
		## Extract ids resolved by selectors
		datas = self.storage.find({ 'crecord_type': 'selector', 'enable': True, 'rk': { '$exists' : True } }, mfields=['_id', 'crecord_name', 'ids'], namespace="object")
		for data in datas:
			_id = data.get('_id')
			ids = data.get('ids')
			if ids:
				self.tags_ids.append( (data.get('crecord_name'), ids) )
