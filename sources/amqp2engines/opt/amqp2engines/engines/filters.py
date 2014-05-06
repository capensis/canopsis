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

from cengine import cengine, DROP
from cstorage import get_storage
from cstatemap import cstatemap
from caccount import caccount

import cmfilter
import logging
import cevent
import time
import json
import ast

NAME='filters'

class engine(cengine):

	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		account = caccount(user="root", group="root")
		self.storage = get_storage(logging_level=logging.DEBUG, account=account)
		self.derogations = []

	def pre_run(self):
		self.drop_event_count = 0
		self.pass_event_count = 0
		self.beat()

	def time_conditions(self, derogation):
		conditions = derogation.get('time_conditions', None)

		if not isinstance(conditions, list):
			self.logger.error(("Invalid time conditions field in '%s': %s"
					   % (derogation['_id'], conditions)))
			self.logger.debug(derogation)
			return False

		result = False

		now = time.time()
		for condition in conditions:
			if (condition['type'] == 'time_interval'
			    and condition['startTs']
			    and condition['stopTs']):
				always = condition.get('always', False)

				if always:
					self.logger.debug(" + 'time_interval' is 'always'")
					result = True

				elif (now >= condition['startTs']
				      and now < condition['stopTs']):
					self.logger.debug(" + 'time_interval' Match")
					result = True

		return result



	def override(self, event, derogation, action):
		name = derogation.get('name', None)
		description = derogation.get('description', None)
		_id = derogation.get('_id', None)

		if not _id or not name or not description:
			self.logger.error("Malformed derogation: %s" % derogation)
			return event

		# If _id is ObjectId(), transform it to str()
		if not isinstance(_id, basestring):
			_id = str(_id)

		derogated = False
		atype = action.get('type', None)

		if atype == "override":
			afield = action.get('field', None)
			avalue = action.get('value', None)

			if afield and avalue:
				event[afield] = avalue
				derogated = True
				self.logger.debug(("    + %s: Override: '%s' -> '%s'"
						   % (event['rk'], afield, avalue)))

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



	def work(self, event, *xargs, **kwargs):
		if 'downtime' in event and event['downtime']:
			event['state'] = 0
			self.logger.debug('derogation to apply on event')
		else:
			self.logger.debug('no derogation to apply on event %s ',
					  (str(event)))

		rk = cevent.get_routingkey(event)
		default_action = self.configuration.get('default_action', 'pass')

		# When list configuration then check black and
		# white lists depending on json configuration
		for filterItem in self.configuration.get('rules', []):
			actions = filterItem.get('actions')
			name = filterItem.get('name', 'no_name')

			# Try filter rules on current event
			if cmfilter.check(filterItem['mfilter'], event):
				for action in actions:
					if action['type'] == 'pass':
						self.logger.debug("Event passed by rule '%s'" % name)
						self.pass_event_count += 1
						return event

					elif action['type'] == 'drop':
						self.logger.debug("Event dropped by rule '%s'" % name)
						self.drop_event_count += 1
						return DROP

					elif (action['type'] == 'override'
					      or action['type'] == 'requalificate'):
						#if self.time_conditions(filterItem):
						self.logger.debug("Event changed by rule '%s'" % name)
						return self.override(event, filterItem, action)

					else:
						self.logger.warning("Unknown action '%s'" % action)

		# No rules matched
		if default_action == 'drop':
			self.logger.debug("Event '%s' dropped by default action" % (rk))
			self.drop_event_count += 1
			return DROP

		self.logger.debug("Event '%s' passed by default action" % (rk))
		self.pass_event_count += 1

		return event



	def beat(self, *args, **kargs):
		""" Configuration reload for realtime ui changes handling """
		self.derogations = []
		self.configuration = {'rules': [],
				      'default_action': self.find_default_action()}

		try:
			records = self.storage.find({'crecord_type':'rule'},
						    sort='priority')

			for record in records:
				record_dump = record.dump()
				record_dump["mfilter"] = ast.literal_eval(record_dump["mfilter"])
				self.configuration['rules'].append(record_dump)

			self.send_stat_event()

		except Exception, e:
			self.logger.warning(str(e))



	def send_stat_event(self):
		""" Send AMQP Event for drop and pass metrics """

		event = cevent.forger(
			connector = "cengine",
			connector_name = "engine",
			event_type = "check",
			source_type = "resource",
			resource = self.amqp_queue + '_data',
			state = 0,
			state_type = 1,
			output = ("%s event dropped since %s"
				  % (self.drop_event_count, self.beat_interval)),
			perf_data_array = [
				{'metric': 'pass_event',
				 'value': self.pass_event_count,
				 'type': 'GAUGE'},
				{'metric': 'drop_event',
				 'value': self.drop_event_count,
				 'type': 'GAUGE' }])

		self.logger.debug(("%s event dropped since %s"
				   % (self.drop_event_count, self.beat_interval)))
		self.logger.debug(("%s event passed since %s"
				   % (self.pass_event_count, self.beat_interval)))

		rk = cevent.get_routingkey(event)
		self.amqp.publish(event, rk, self.amqp.exchange_name_events)
		self.drop_event_count = 0
		self.pass_event_count = 0



	def find_default_action(self):
		""" Find the default action stored and returns it, else assume it default action is pass """

		records = self.storage.find({'crecord_type':'defaultrule'})
		if records:
			return records[0].dump()["action"]

		self.logger.debug("No default action found. Assuming default action is pass")
		return 'pass'



	def set_derogation_state(self, derogation, active):
		dactive = derogation.get('active', False)
		name = derogation.get('crecord_name', None)
		notify = False
		state = 0

		if active and not dactive:
			self.logger.info(("%s (%s) is now active"
					  % (derogation['crecord_name'],
					     derogation['_id'])))
			self.storage.update(derogation['_id'], {'active': True})
			notify = True

		elif not active and dactive:
			self.logger.info(("%s (%s) is now inactive"
					  % (derogation['crecord_name'],
					     derogation['_id'])))
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
				source_type = "component",
				component = NAME,
				state = state,
				output = output,
				long_output = derogation.get('description', None)
			)
			rk = cevent.get_routingkey(event)

			self.amqp.publish(event, rk, self.amqp.exchange_name_events)



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
			self.logger.debug(("Load %s derogations."
					   % len(self.derogations)))


