#!/usr/bin/env python
# -*- coding: utf-8 -*-
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

import logging

from cengine import cengine
import cevent

from caccount import caccount
from cstorage import get_storage
from crecord import crecord
import time

NAME="acknowledgement"


class engine(cengine):
	def __init__(self, name=NAME, acknowledge_on='canopsis.events', *args, **kargs):
		cengine.__init__(self, name=name, *args, **kargs)

		account = caccount(user="root", group="root")

		self.storage = get_storage(namespace='ack', account=account, logging_level=logging.DEBUG)
		self.stbackend = self.storage.get_backend('ack')

		self.acknowledge_on = acknowledge_on
		
	def work(self, event, *args, **kargs):
		logevent = None

		# If event is of type acknowledgement, then acknowledge corresponding event
		if event['event_type'] == 'ack':
			rk = event['referer']

			# add rk to acknowledged rks
			record = self.stbackend.find_and_modify(
				query = {'rk': rk, 'solved': False},
				update = {'$set': {
					'timestamp': event['timestamp'],
					'ackts': int(time.time()),
					'rk': rk,
					'author': event['author'],
					'comment': event['output']
				}},
				upsert = True
			)

			if not record:
				# Emit an event log
				referer_event = self.storage.find_one(mfilter={'rk': rk})

				logevent = cevent.forger(
					connector = "cengine",
					connector_name = NAME,
					event_type = "log",
					source_type = referer_event['source_type'],
					component = referer_event['component'],
					resource = referer_event.get('resource', None),

					state = 0,
					state_type = 1,

					referer = event['rk'],
					output = u'Event {0} acknowledged by {1}'.format(rk, event['author']),
					long_output = event['output'],

					perf_data_array = [
						{
							'metric': 'ack_delay',
							'value': record['ackts'] - record['timestamp'],
							'unit': 's'
						}
					]
				)

		# If event is acknowledged, and went back to normal, remove the ack
		elif event['state'] == 0:
			record = self.stbackend.find_and_modify(
				query = {'rk': event['rk'], 'solved': False},
				update = {'$set': {'solved': True}}
			)

			if record:
				logevent = cevent.forger(
					connector = "cengine",
					connector_name = NAME,
					event_type = "log",
					source_type = event['source_type'],
					component = event['component'],
					resource = event.get('resource', None),

					state = 0,
					state_type = 1,

					referer = event['rk'],
					output = u'Acknowledgement removed for event {0}'.format(event['rk']),
					long_output = u'Everything went back to normal',

					perf_data_array = [
						{
							'metric': 'ack_solved',
							'value': record['ackts'] - int(time.time()),
							'unit': 's'
						}
					]
				)

		# If the event is in problem state, update the solved state of acknowledgement
		else:
			self.stbackend.find_and_modify(
				query = {'rk': event['rk'], 'solved': True},
				update = {'$set': {'solved': False}}
			)

		if logevent:
			self.amqp.publish(logevent, cevent.get_routingkey(logevent), exchange_name=self.acknowledge_on)

		return event
