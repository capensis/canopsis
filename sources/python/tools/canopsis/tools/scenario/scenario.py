#!/usr/bin/python

import time, sys
from kombu import Connection
from kombu.pools import producers
from random import random, randint

user 		= "canopsis"
password 	= "canopsis"
vhost 		= "canopsis"
exchange 	= "canopsis.events"
port		= '5672'
host 		= 'localhost'

try:
	print 'Try to load : scenarii ' + sys.argv[1]
	module = __import__('scenarii.' + sys.argv[1])
	scenario = getattr(module, sys.argv[1]).scenario
except Exception ,e:
	print 'unable to load scenario. ', e
	sys.exit(1)

def get_rk(event):
	return "%s.%s.%s.%s.%s" % (event['connector'], event['connector_name'], event['event_type'], event['source_type'], event['component'])

with Connection(port=port, hostname=host, userid=user, password=password, virtual_host=vhost) as conn:
	with producers[conn].acquire(block=True) as producer:
		event = {}
		for part in scenario:
			title, event = part(event)
			if event:
				print 'Processing... %s' % title
				event["timestamp"] = time.time()

				producer.publish(
				event,
				serializer='json',
				exchange=exchange,
				routing_key=get_rk(event))
				print 'sent event %s' % get_rk(event)
				time.sleep(1)
		producer.close()
