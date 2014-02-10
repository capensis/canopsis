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
import pyperfstore2

from cstorage import get_storage
from caccount import caccount
import time

NAME="alertcounter"

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		
	def pre_run(self):
		self.listened_event_type = ['check','selector','eue','sla', 'log']
		self.manager = pyperfstore2.manager()	
		
		self.selectors_name = []
		self.last_resolv = 0
		
	def resolve_selectors_name(self):
		if int(time.time()) > (self.last_resolv + 60):
			storage = get_storage(namespace='object', account=caccount(user="root", group="root"))
			records = storage.find(mfilter={'crecord_type': 'selector'}, mfields=['crecord_name'])
			
			self.selectors_name = [record['crecord_name'] for record in records]
			
			self.last_resolv = int(time.time())
			
			del storage
		
	
	def count_alert(self, component, state, value, resource=None, tags=[]):
		
		if resource:
			meta_data = {'type': 'COUNTER', 'co': component, 're': resource }
			name = "%s%s" % (meta_data['co'], meta_data['re'])
		else:
			meta_data = {'type': 'COUNTER', 'co': component }
			name = meta_data['co']

		if tags:
			meta_data['tg'] = tags
		
		metric = "cps_statechange"
		meta_data['me'] = metric
				
		self.logger.debug("Increment %s: %s: %s" % (name, metric, value))
		self.manager.push(name="%s%s" % (name, metric), value=value, meta_data=meta_data)

		metric = "cps_statechange_nok"
		meta_data['me'] = metric
		
		cvalue = 0
		if state != 0:
			cvalue = value
			
		self.logger.debug("Increment %s: %s: %s" % (name, metric, cvalue))
		self.manager.push(name="%s%s" % (name, metric), value=cvalue, meta_data=meta_data)
		
		for cstate in [ 0, 1, 2, 3 ]:
			cvalue = 0
			if cstate == state:
				cvalue = value
				
			metric = "cps_statechange_%s" % cstate
			meta_data['me'] = metric
			meta_data['type'] = 'COUNTER'
			
			self.logger.debug(" + Increment '%s': %s" % (metric, cvalue))
			self.manager.push(name="%s%s" % (name, metric), value=cvalue, meta_data=meta_data)
		
	
	def work(self, event, *args, **kargs):
		if event['event_type'] in self.listened_event_type:
							
			# Hard state
			if event.get('state_type', 1) == 1 and event['component'] != 'derogation':
				
				# By Tags (Selector)
				if event['event_type'] != 'selector':
					tags = event.get('tags', [])
					if tags:
						self.resolve_selectors_name()
						for tag in tags:
							if tag in self.selectors_name:
								self.logger.debug("Increment Tag: '%s'" % tag)
								self.count_alert(
									component	= tag,
									resource	= 'selector',
									state		= event['state'],
									value		= 1,
									tags		= event['tags']
							)
				
				# By name
				self.count_alert(
					component	= event['component'],
					resource	= event.get('resource', None),
					state		= event['state'],
					value		= 1,
					tags		= event['tags']
				)
	
		return event
