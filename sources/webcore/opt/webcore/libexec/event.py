#!/usr/bin/env python
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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
import json
import time

from bottle import put, request, HTTPError, post

## Canopsis
from canopsis.old.rabbitmq import Amqp
import canopsis.old.event as cevent
import requests

logger = logging.getLogger('Event')
logger.setLevel(logging.INFO)

group_managing_access = 'group.CPS_event_admin'

enable_crossdomain_send_events = False

try:
	enable_crossdomain_send_events = config.getboolean('server', "enable_crossdomain_send_events")
except Exception as e:
	logger.info("Setting \"enable_crossdomain_send_events\" does not seems to appear on webserver configuration. Using False as default value")
	logger.debug(str(e))

print("enable_crossdomain_send_events")
print(enable_crossdomain_send_events)


## load functions
def load_amqp(host="localhost"):
	amqp = Amqp(logging_name="Event-amqp", host=host)
	amqp.start()
	return amqp


def unload_amqp(amqp):
	if amqp:
		amqp.stop()
		amqp.join()


## Handlers
@post('/event')
@put('/event')
def send_event_route():

	items = request.params.get('event', '[]')
	if isinstance(items, str):
		try:
			items = json.loads(items)

			if not isinstance(items, list):
				items = [items]

		except Exception as err:
			logger.error("PUT: Impossible to parse data ({})".format(err))
			return HTTPError(500, "Impossible to parse data")

	host = request.params.get("host", "localhost")

	messages = []
	success = True
	total = 0
	for event in items:
		result = send_event(event, host)
		total += result['total']
		success = success and result['success']
		messages.append({
			'event': cevent.get_routingkey(event),
			'message': result['data']['message']
		})

	return {'total': total, 'success': success, 'messages': messages}


def send_event(event, host):

	timestamp = event.get('timestamp', time.time())

	if timestamp is None:
		timestamp = int(time.time())

	if not isinstance(timestamp, int):
		event['timestamp'] = int(timestamp)

	mandatory_fields = ['connector', 'connector_name', 'event_type', 'source_type', 'component', 'state']
	missing_fields = []

	for field in mandatory_fields:
		if event.get(field, None) is None:
			missing_fields.append(field)

	if missing_fields:
		message = 'Missing {} argument in payload'.format(missing_fields)
		logger.error(message)
		return {'total': 0, 'success': False, 'data': {'message': message}}

	if event.get('state_type', None) is None:
		event['state_type'] = 1

	def json2py(data, key):
		value = data.get(key, None)
		if value:
			try:
				value = json.loads(value)
			except Exception as err:
				logger.error("Impossible to parse '{0}' ({1})".format(key, err))

		if not isinstance(value, list):
			value = []
		data[key] = value

	json2py(event, 'tags')
	json2py(event, 'perf_data_array')

	#------------------------------forging event----------------------------------
	forged_event = cevent.forger(
		connector=event.get('connector', None),
		connector_name=event.get('connector_name', None),
		event_type=event.get('event_type', None),
		source_type=event.get('source_type', None),
		component=event.get('component', None),
		resource=event.get('resource', None),
		state=int(event.get('state', None)),
		state_type=int(event.get('state_type', None)),
		output=event.get('output', None),
		long_output=event.get('long_output', None),
		perf_data=event.get('perf_data', None),
		perf_data_array=event.get('perf_data_array', None),
		timestamp=event.get('timestamp', None),
		display_name=event.get('display_name', None),
		tags=event.get('tags', None),
		ticket=event.get('ticket', None),
		ref_rk=event.get('ref_rk', None),
		cancel=event.get('cancel', None),
		author=event.get('author', None),
	)

	logger.debug('Event crafted {}'.format(forged_event))

	if host == "localhost":

		key = cevent.get_routingkey(forged_event)

		logger.info("now send event to amqp")

		amqp = load_amqp(host)
		amqp.publish(forged_event, key, amqp.exchange_name_events)
		unload_amqp(amqp)

		logger.debug('Amqp event published')

		return {'total': 1, 'success': True, 'data': {'message': 'Event send completed successfully'}}

	else:
		global enable_crossdomain_send_events

		logger.info("sending event to host: " + host)
		if enable_crossdomain_send_events is True:

			payload = {'event': json.dumps(forged_event)}

			response = requests.post(host, data=payload)
			if response.status_code == 200:

				try:
					response = json.loads(response.text)
					if response['success']:
						response = 'Event send to remote host {} completed successfully'.format(host)
				except Exception:
					response = 'error while parsing return valure from remote host {} event sent is uncertain'.format(host)
				return {'total': 1, 'success': True, 'data': {'message': response}}

			else:
				return {'total': 0, 'success': False, 'data': {'message': 'was unable to follow event @ {}'.format(host)}}

		else:
			logger.info("Cross domain send_events are not authorized on the webserver, please check your webserver config file")
			return HTTPError(403, "Cross domain send_events are not authorized on the webserver")
