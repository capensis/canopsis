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

		self.storage = get_storage(namespace='ack', account=account)
		self.stbackend = self.storage.get_backend('ack')

		self.acknowledge_on = acknowledge_on
		
	def work(self, event, *args, **kargs):
		logevent = None

		# If event is of type acknowledgement, then acknowledge corresponding event
		if event['event_type'] == 'ack':
			rk = event.get('referer', event.get('ref_rk', None))

			if not rk:
				self.logger.error("Cannot get acknowledged event, missing key referer or ref_rk")
				return event

			# add rk to acknowledged rks
			response = self.stbackend.find_and_modify(
				query = {'rk': rk, 'solved': False},
				update = {'$set': {
					'timestamp': event['timestamp'],
					'ackts': int(time.time()),
					'rk': rk,
					'author': event['author'],
					'comment': event['output']
				}},
				upsert = True,
				full_response = True,
				new = True
			)

			self.logger.error(str(response))

			if not response['lastErrorObject']['updatedExisting']:
				record = response['value']

				# Emit an event log
				referer_event = self.storage.find_one(mfilter={'rk': rk}, namespace='events').dump()

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

			# Now update counters
			alerts_event = cevent.forger(
				connector = "cengine",
				connector_name = NAME,
				event_type = "perf",
				source_type = "component",
				component = "__canopsis__",

				perf_data_array = [
					{'metric': 'cps_alerts_ack', 'value': 1, 'type': 'COUNTER'},
					{'metric': 'cps_alerts_ack_by_host', 'value': 0, 'type': 'COUNTER'},
					{'metric': 'cps_alerts_not_ack', 'value': -1, 'type': 'COUNTER'}
				]
			)

			self.amqp.publish(alerts_event, cevent.get_routingkey(alerts_event), self.amqp.exchange_name_events)

			for hostgroup in event.get('hostgroups', []):
				alerts_event = cevent.forger(
					connector = "cengine",
					connector_name = NAME,
					event_type = "perf",
					source_type = "resource",
					component = "__canopsis__",
					resource = hostgroup,

					perf_data_array = [
						{'metric': 'cps_alerts_ack', 'value': 1, 'type': 'COUNTER'},
						{'metric': 'cps_alerts_ack_by_host', 'value': 0, 'type': 'COUNTER'},
						{'metric': 'cps_alerts_not_ack', 'value': -1, 'type': 'COUNTER'}
					]
				)

				self.amqp.publish(alerts_event, cevent.get_routingkey(alerts_event), self.amqp.exchange_name_events)


		# If event is acknowledged, and went back to normal, remove the ack
		elif event['state'] == 0 and event.get('state_type', 1) == 1:
			solvedts = int(time.time())

			response = self.stbackend.find_and_modify(
				query = {
					'rk': event['rk'],
					'solved': False,
					'ackts': {'$gt': -1}
				},
				update = {'$set': {
					'solved': True,
					'solvedts': solvedts
				}},
				full_response = True,
				new = True
			)

			if response and response['value']:
				record = response['value']

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
							'metric': 'ack_solved_delay',
							'value': solvedts - record['ackts'],
							'unit': 's'
						}
					]
				)

				logevent['acknowledged_connector'] = event['connector']
				logevent['acknowledged_source'] = event['connector_name']
				logevent['acknowledged_at'] = record['ackts']
				logevent['solved_at'] = solvedts

		# If the event is in problem state, update the solved state of acknowledgement
		elif event['state'] != 0 and event.get('state_type', 1) == 1:
			self.stbackend.find_and_modify(
				query = {'rk': event['rk'], 'solved': True},
				update = {'$set': {
					'solved': False,
					'solvedts': -1,
					'ackts': -1,
					'timestamp': -1,
					'author': '',
					'comment': ''
				}}
			)

		if logevent:
			self.amqp.publish(logevent, cevent.get_routingkey(logevent), exchange_name=self.acknowledge_on)

		return event
