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
import cmfilter
import time

import acknowledgement as engine_ack

NAME="downtime"


class engine(cengine):
	def __init__(self, name=NAME, *args, **kwargs):
		cengine.__init__(self, name=name, *args, **kwargs)

		account = caccount(user="root", group="root")

		self.storage = get_storage(namespace='downtime', account=account, logging_level=logging.DEBUG)
		self.dt_backend = self.storage.get_backend('downtime')
		self.evt_backend = self.storage.get_backend('events')

	def beat(self):
		self.logger.debug('Removing expired downtime entries')

		now = time.time()

		# Remove downtime that are expired
		records = self.storage.find({
			'_expire': {
				'$lt': now
			}
		})

		# No downtime found
		if not records:
			return

		self.storage.remove([r._id for r in records])

		# Build query
		matching = []

		for record in records:
			record = record.dump()

			matching.append({
				'connector': record['connector'],
				'connector_name': record['source'],
				'component': record['component'],
				'resource': record['resource'],
				'downtime': True
			})

		# Now, update all matching events unset the downtime information
		records = self.evt_backend.update(
			{ '$or': matching },
			{
				'$set': {
					'downtime': False
				}
			},
			multi = True
		)

	def work(self, event, *args, **kwargs):
		now = time.time()

		# If the event is a downtime event, add entry to the downtime collection
		if event['event_type'] == 'downtime':
			self.logger.debug('Event downtime received: {0}'.format(event['rk']))

			# Build entry, so we know there is a downtime on the component
			record = crecord({
				'_expire': event['start'] + event['duration'],

				'connector': event['connector'],
				'source': event['connector_name'],
				'component': event['component'],
				'resource': event.get('resource', None),

				'start': event['start'],
				'end': event['end'],
				'fixed': event['fixed'],
				'timestamp': event['entry'],

				'author': event['author'],
				'comment': event['output']
			})

			# Save record, and log the action
			record.save(self.storage)

			logevent = cevent.forger(
				connector = "cengine",
				connector_name = NAME,
				event_type = "log",
				source_type = event['source_type'],
				component = event['component'],
				resource = event.get('resource', None),

				state = 0,
				state_type = 1,

				output = u'Downtime scheduled by {0} from {1} to {2}'.format(
					event['author'],
					event['start'],
					event['end']
				),

				long_output = event['output']
			)

			logevent['downtime_connector'] = event['connector']
			logevent['downtime_source'] = event['connector_name']

			self.amqp.publish(logevent, cevent.get_routingkey(logevent), exchange_name='canopsis.events')

			# Set downtime for events already in database
			self.evt_backend.update(
				{
					'connector': event['connector'],
					'connector_name': event['connector_name'],
					'component': event['component'],
					'resource': event.get('resource', None)
				},
				{
					'$set': {
						'downtime': True
					}
				},
				multi = True
			)

		# For every other case, check if the event is in downtime
		else:
			self.logger.debug('Received event: {0}'.format(event['rk']))

			event['downtime'] = False

			records = self.storage.find({
				'connector': event['connector'],
				'source': event['connector_name'],
				'component': event['component'],
				'resource': event.get('resource', None)
			})

			for record in records:
				downtime_info = record.dump()

				# If the downtime is active
				if (downtime_info['start'] <= now <= downtime_info['end']):
					# Set the event as 'in downtime'
					event['downtime'] = True
					break

		return event