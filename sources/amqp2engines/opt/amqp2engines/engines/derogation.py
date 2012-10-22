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
import cevent
import time

NAME="derogation"

class engine(cengine):
	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		
		self.derogations = []
		
	def pre_run(self):
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))		
		self.beat()
		
	def time_conditions(self, derogation):
		conditions = derogation.get('time_conditions', None)
		
		if not isinstance(conditions, list):
			self.logger.error("Invalid time conditions field in '%s': %s" % (derogation['_id'], conditions))
			self.logger.debug(derogation)
			return False
		
		result = False
		
		now = time.time()
		for condition in conditions:
			if condition['type'] == 'time_interval':
				if now >= condition['startTs'] and now < condition['stopTs']:
					self.logger.debug(" + 'time_interval' Match")
					result = True
					
		return result		
		
	
	def conditions(self, event, derogation):
		return True
	
	def actions(self, event, derogation):
		actions = derogation.get('actions', None)
		name = derogation.get('crecord_name', None)
		
		if not isinstance(actions, list):
			self.logger.error("Invalid actions field in '%s': %s" % (derogation['_id'], actions))
			return event
		
		for action in derogation['actions']:
			if action['type'] == "override":
				self.logger.debug("    + %s: Override: '%s' -> '%s'" % (event['rk'], action['field'], action['value']))
				event[action['field']] = action['value']
				event["derogation_id"] = derogation['_id']
				event["derogation_description"] = derogation['description']
				event["derogation_name"] = derogation['crecord_name']
				event["tags"].append(name)
				event["tags"].append("derogated")
			else:
				self.logger.warning("Unknown action '%s'" % action['type'])
				
		return event
	
	def work(self, event, *args, **kargs):
		for derogation in self.derogations:
			# Check scope
			if event['rk'] in derogation['ids']:
				self.logger.debug("%s is in %s (%s)" % (event['rk'], derogation['crecord_name'], derogation['_id']))
				# Check Time
				if self.time_conditions(derogation):
					# Check conditions
					if self.conditions(event, derogation):
						# Actions
						return self.actions(event, derogation)
		
		return event
		
	def set_derogation_state(self, derogation, active):
		dactive = derogation.get('active', False)
		name = derogation.get('crecord_name', None)
		notify = False
		state = 0
		
		if active:
			if not dactive:
				self.logger.info("%s (%s) is now active" % (derogation['crecord_name'], derogation['_id']))
				self.storage.update(derogation['_id'], {'active': True})
				notify = True
		else:
			if dactive:
				self.logger.info("%s (%s) is now inactive" % (derogation['crecord_name'], derogation['_id']))
				self.storage.update(derogation['_id'], {'active': False})
				notify = True
				
		if notify:
			if active:
				output = "Derogation '%s' is now active" % name
				state = 1
			else:
				output = "Derogation '%s' is now inactive" % name
			
			
			tags = derogation.get('tags', None)
			self.logger.debug(" + Tags: '%s' (%s)" % (tags, type(tags)))
			
			if isinstance(tags, str) or isinstance(tags, unicode):
				tags = [ tags ]
			
			if not isinstance(tags, list) or tags == "":
				tags = None
				
			event = cevent.forger(
				connector = "cengine",
				connector_name = "engine",
				event_type = "log",
				source_type="component",
				component=NAME,
				state=state,
				output=output,
				long_output=derogation.get('description', None),
				tags=tags
			)
			rk = cevent.get_routingkey(event)
			
			self.amqp.publish(event, rk, self.amqp.exchange_name_events)

	def beat(self):
		self.derogations = []
		
		## Extract ids
		records = self.storage.find( {	'crecord_type': 'derogation',
										'enable': True,
										'ids': { '$exists' : True },
										'actions': { '$exists' : True },
										'conditions': { '$exists' : True } },
										namespace="object")
		
		for record in records:
			if record.data.get('ids'):
				derogation = record.dump()
				if self.time_conditions(derogation):
					self.set_derogation_state(derogation, True)	
				else:
					self.set_derogation_state(derogation, False)
					
				self.derogations.append(derogation)
				
		self.logger.debug("Load %s derogations." % len(self.derogations))
