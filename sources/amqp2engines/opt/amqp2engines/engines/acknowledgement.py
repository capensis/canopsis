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
			rk = event['referer']

			# add rk to acknowledged rks
			record = self.storage.find_one(mfilter={'rk': rk})

			if not record:
				record = crecord({'rk': rk})
				self.storage.put(record)

			referer_event = self.storage.find_one(mfilter={'rk': rk})

			logevent = cevent.forger(
				connector = "cengine",
				connector_name = "engine",
				event_type = "log",
				source_type = referer_event['source_type'],
				component = referer_event['component'],
				resource = referer_event.get('resource', None),

				state = 0,
				state_type = 1,

				output = u'Event {0} acknowledged by {1}'.format(rk, event['author']),
				long_output = event['output']
			)

			self.amqp.publish(logevent, cevent.get_routingkey(logevent), exchange_name=self.acknowledge_on)

		# If event is acknowledged, and went back to normal, remove the ack
		elif event['state'] == 0:
			record = self.storage.find_one(mfilter={'rk': event['rk']})

			if record:
				self.storage.remove(record._id)

				logevent = cevent.forger(
					connector = "cengine",
					connector_name = "engine",
					event_type = "log",
					source_type = event['source_type'],
					component = event['component'],
					resource = event.get('resource', None),

					state = 0,
					state_type = 1,

					output = u'Acknowledgement removed for event {0}'.format(rk),
					long_output = u'Went back to normal'
				)

				self.amqp.publish(logevent, cevent.get_routingkey(logevent), exchange_name=self.acknowledge_on)

		return event
