#!/usr/bin/env python
# --------------------------------
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

import logging, json

import bottle
from bottle import route, get, put, delete, request, HTTPError, post, response

## Canopsis
from caccount import caccount
from cstorage import cstorage
from cstorage import get_storage
from crecord import crecord
from camqp import camqp
import cevent

#import protection function
from libexec.auth import check_auth, get_account, check_group_rights

logger = logging.getLogger('Event')

amqp = None
group_managing_access = 'group.CPS_event_admin'

##################################################################################

## load functions
def load():
	global amqp
	amqp = camqp(logging_name="Event-amqp")
	amqp.start()	
	
def unload():
	global amqp
	if amqp:
		amqp.stop()
		amqp.join()
	
##################################################################################

## Handlers

@post('/event/',checkAuthPlugin={'authorized_grp':group_managing_access})
@post('/event/:routing_key',checkAuthPlugin={'authorized_grp':group_managing_access})
def send_event(	routing_key=None):
	account = get_account()
	
	connector = None
	connector_name = None
	event_type = None
	source_type = None
	component = None
	resource = None
	state = None
	state_type = None
	perf_data = None
	perf_data_array = None
	output = None
	long_output = None
	timestamp = None
	display_name = None
				
	#--------------------explode routing key----------
	if routing_key :
		logger.debug('The routing key is : %s' % str(routing_key))
		
		routing_key = routing_key.split('.')
		if len(routing_key) > 6 or len(routing_key) < 5:
			logger.error('Bad routing key')
			return HTTPError(400, 'Bad routing key')
			
		connector = routing_key[0]
		connector_name = routing_key[1]
		event_type = routing_key[2]
		source_type = routing_key[3]
		component = routing_key[4]
		if routing_key[5]:
			resource = routing_key[5]
	
	
	#-----------------------get params-------------------
	if not timestamp:
		timestamp = request.params.get('timestamp', default=None)
	
	#fix timestamp type
	if timestamp and not isinstance(timestamp, int):
		timestamp = int(timestamp)
		
	if not display_name:
		display_name = request.params.get('display_name', default=None)
	
	if not connector:
		connector = request.params.get('connector', default=None)
		if not connector :
			logger.error('No connector argument')
			return HTTPError(400, 'Missing connector argument')
			
	if not connector_name:
		connector_name = request.params.get('connector_name', default=None)
		if not connector_name:
			logger.error('No connector name argument')
			return HTTPError(400, 'Missing connector name argument')
			
	if not event_type:
		event_type = request.params.get('event_type', default=None)
		if not event_type:
			logger.error('No event_type argument')
			return HTTPError(400, 'Missing event type argument')
		
	if not source_type:
		source_type = request.params.get('source_type', default=None)
		if not source_type:
			logger.error('No source_type argument')
			return HTTPError(400, 'Missing source type argument')
	
	if not component:
		component = request.params.get('component', default=None)
		if not component:
			logger.error('No component argument')
			return HTTPError(400, 'Missing component argument')
	
	if not resource:
		resource = request.params.get('resource', default=None)
		if not resource:
			logger.error('No resource argument')
			return HTTPError(400, 'Missing resource argument')
		
	if not state:
		state = request.params.get('state', default=None)
		if not state:
			logger.error('No state argument')
			return HTTPError(400, 'Missing state argument')
		
	if not state_type:
		state_type = request.params.get('state_type', default=1)
		
	if not output:
		output = request.params.get('output', default=None)
		
	if not long_output:
		long_output = request.params.get('long_output', default=None)
		
	if not perf_data:
		perf_data = request.params.get('perf_data', default=None)
		
	if not perf_data_array:
		perf_data_array = request.params.get('perf_data_array', default=None)
		if perf_data_array:
			try:
				perf_data_array = json.loads(perf_data_array)
			except Exception, err:
				logger.error("Impossible to parse 'perf_data_array': %s (%s)" % (perf_data_array, err))
		
	#------------------------------forging event----------------------------------

	event = cevent.forger(
				connector = connector,
				connector_name = connector_name,
				event_type = event_type,
				source_type = source_type,
				component = component,
				resource= resource,
				state = int(state),
				state_type = int(state_type),
				output = output,
				long_output = long_output,
				perf_data = perf_data,
				perf_data_array = perf_data_array,
				timestamp = timestamp,
				display_name = display_name
			)
	
	logger.debug(type(perf_data_array))
	logger.debug(perf_data_array)
	logger.debug('The forged event is : ')
	logger.debug(str(event))
	
	#------------------------------AMQP Part--------------------------------------
	
	key = cevent.get_routingkey(event)
	
	global amqp
	amqp.publish(event, key, amqp.exchange_name_events)
		
	logger.debug('Amqp event published')
	
	return {'total':1,'success':True,'data':{'event':event}}
