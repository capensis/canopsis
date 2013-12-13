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
import cevent

from caccount import caccount
from cstorage import get_storage
from crecord import crecord

NAME="acknowledgement"


class engine(cengine):
	def __init__(self, name=NAME, acknowledge_on='canopsis.events', *args, **kargs):
		cengine.__init__(self, name=name, *args, **kargs)

		account = caccount(user="root", group="root")
		self.storage = get_storage(namespace='ack', account=account, logging_level=logging.DEBUG)

		self.acknowledge_on = acknowledge_on
		
	def work(self, event, *args, **kargs):
		# If event is of type acknowledgement, then acknowledge corresponding event
		if event['event_type'] == 'ack':
			if event['source_type'] == 'resource':
				rk = event['referer'].split('.', 5)

			elif event['source_type'] == 'component':
				rk = event['referer'].split('.', 4)

			new_event = cevent.forger(
				connector = rk[0],
				connector_name = rk[1],
				event_type = rk[2],
				source_type = rk[3],
				component = rk[4],
				state = event['state'],
				output = u'Event acknowledged by {0}'.format(event['author']),
				long_output = event['output']
			)

			if event['source_type'] == 'resource':
				new_event['resource'] = rk[5]

			rk = rk.join('.')
			self.amqp.publish(new_event, rk, exchange_name=self.acknowledge_on)

			# add rk to acknowledged rks
			record = self.storage.find_one(mfilter={'rk': rk})

			if not record:
				record = crecord({'rk': rk})
				self.storage.put(record)

		# If event is acknowledged, and went back to normal, remove the ack
		elif event['state'] == 0:
			record = self.storage.find_one(mfilter={'rk': event['rk']})

			if record:
				self.storage.remove(record._id)

		return event
