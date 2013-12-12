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
from cengine import DROP

from caccount import caccount
from cstorage import get_storage
import cevent
import logging
import cmfilter
import ast

NAME='event_filter'

class engine(cengine):

	def __init__(self, *args, **kargs):
		cengine.__init__(self, name=NAME, *args, **kargs)
		account = caccount(user="root", group="root")
		self.storage = get_storage(logging_level=logging.DEBUG, account=account)
		
		
	def pre_run(self):
		self.drop_event_count = 0
		self.pass_event_count = 0
		self.beat()


	def work(self, event, *xargs, **kwargs):		

		rk = cevent.get_routingkey(event)

		default_action = self.configuration.get('default_action', 'pass')

		#When list configuration then check black and white lists depending on json configuration
		for filterItem in self.configuration.get('rules', []):

			action = filterItem.get('action')

			name = filterItem.get('name', 'no_name')
		
			# Try filter rules on current event
			if cmfilter.check(filterItem['mfilter'], event):
				if action == 'pass':
					self.logger.debug("Event passed by rule '%s'" % name)
					self.pass_event_count += 1
					return event

				elif action == 'drop':
					self.logger.debug("Event dropped by rule '%s'" % name)
					self.drop_event_count += 1
					return DROP

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

		self.configuration = { 'rules' : [], 'default_action': self.find_default_action()}

		try:
			records = self.storage.find({'crecord_type':'rule'}, sort="priority")

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
			source_type="resource",
			resource=self.amqp_queue + '_data',
			state=0,
			state_type=1,
			output="%s event dropped since %s" % (self.drop_event_count, self.beat_interval),
			perf_data_array=[
								{'metric': 'pass_event' , 'value': self.pass_event_count, 'type': 'GAUGE' },
								{'metric': 'drop_event' , 'value': self.drop_event_count, 'type': 'GAUGE' }
							]
		)

		self.logger.debug("%s event dropped since %s" % (self.drop_event_count, self.beat_interval))
		self.logger.debug("%s event passed since %s" % (self.pass_event_count, self.beat_interval))


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
