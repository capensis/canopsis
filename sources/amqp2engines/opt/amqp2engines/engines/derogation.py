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
from cstatemap import cstatemap
import cmfilter
import time
import json

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
				always = condition.get('always', False)

				if always:
					self.logger.debug(" + 'time_interval' is 'always'")
					result = True

				elif now >= condition['startTs'] and now < condition['stopTs']:
					self.logger.debug(" + 'time_interval' Match")
					result = True
					
		return result		
		
	
	def conditions(self, event, derogation):
		conditions_json = derogation.get('conditions', None)

		try:
			conditions = json.loads(conditions_json)

		except ValueError:
			self.logger.error("Invalid conditions field in '%s': %s" % (derogation['_id'], conditions_json))
			self.logger.debug(derogation)

			return False
		
		if conditions:
			check = cmfilter.check(conditions, event)

			self.logger.debug(" + 'conditions' check is %s" % check)

			return check

		else:
			return True
	
	def actions(self, event, derogation):
		name = derogation.get('crecord_name', None)
		description = derogation.get('description', None)
		actions = derogation.get('actions', None)
		_id = derogation.get('_id', None)

		if not _id or not name or not description or not actions:
			self.logger.error("Malformed derogation: %s" % derogation)
			return event

		# If _id is ObjectId(), transform it to str()
		if not isinstance(_id, basestring):
			_id = str(_id)

		if not isinstance(actions, list):
			self.logger.error("Invalid actions field in '%s': %s" % (_id, actions))
			return event

		derogated = False
		
		for action in actions:
			atype = action.get('type', None)

			if atype == "override":
				afield = action.get('field', None)
				avalue = action.get('value', None)

				self.logger.debug("    + %s: Override: '%s' -> '%s'" % (event['rk'], afield, avalue))

				if afield and avalue:
					event[afield] = avalue

					derogated = True

				else:
					self.logger.error("Action malformed (needs 'field' and 'value'): %s" % action)

			elif atype == "requalificate":
				statemap_id = action.get('statemap', None)

				self.logger.debug("    + %s: Requalificate using statemap '%s'" % (event['rk'], statemap_id))

				if statemap_id:
					record = self.storage.find_one(mfilter={'crecord_type': 'statemap', '_id': statemap_id})

					if not record:
						self.logger.error("Statemap '%s' not found" % statemap_id)

					statemap = cstatemap(record=record)

					event['real_state'] = event['state']
					event['state'] = statemap.get_mapped_state(event['real_state'])

					derogated = True

				else:
					self.logger.error("Action malformed (needs 'statemap'): %s" % action)

			else:
				self.logger.warning("Unknown action '%s'" % atype)
				
		# If the event was derogated, fill some informations
		if derogated:
			event["derogation_id"] = _id
			event["derogation_description"] = description
			event["derogation_name"] = name
			event["tags"].append(name)
			event["tags"].append("derogated")

		return event
	
	def work(self, event, *args, **kargs):
		for derogation in self.derogations:
			# Check Time
			if self.time_conditions(derogation):
				# Check conditions
				if self.conditions(event, derogation):
					self.logger.debug("%s is in %s (%s)" % (event['rk'], derogation['crecord_name'], derogation['_id']))

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
				
			event = cevent.forger(
				connector = "cengine",
				connector_name = "engine",
				event_type = "log",
				source_type="component",
				component=NAME,
				state=state,
				output=output,
				long_output=derogation.get('description', None)
			)
			rk = cevent.get_routingkey(event)
			
			self.amqp.publish(event, rk, self.amqp.exchange_name_events)

	def beat(self):
		self.derogations = []
		
		## Extract ids
		records = self.storage.find( {	'crecord_type': 'derogation',
										'enable': True,
										'actions': { '$exists' : True },
										'conditions': { '$exists' : True } },
										namespace="object")
	def beat(self):
		self.logger.debug('Derogation BEAT')
		
		

	def consume_dispatcher(self,  event, *args, **kargs):
		self.logger.debug("Consolidate metrics:")

		now = time.time()
		beat_elapsed = 0

		record = self.get_ready_record(event)
		if record:	

			if record.data.get('ids'):
				derogation = record.dump()
				if self.time_conditions(derogation):
					self.set_derogation_state(derogation, True)	
				else:
					self.set_derogation_state(derogation, False)
					
				self.derogations.append(derogation)

			self.crecord_task_complete(event['_id'])

			self.logger.debug("Load %s derogations." % len(self.derogations))

