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

import unittest
import time, json, sys
import cevent
import uuid
import re
from camqp import camqp
from cstorage import cstorage
from crecord import crecord
from caccount import caccount
from cwebservices import cwebservices
from ctools import parse_perfdata
import pyperfstore2

from subprocess import Popen

event = cevent.forger(connector='canopsis', connector_name='unittest', event_type='check', source_type = "component", component="test1", state=0, output="Output_1", perf_data="mymetric=1s;10;20;0;30", tags = ['check', 'component', 'test1', 'unittest'])
rk = cevent.get_routingkey(event)

myamqp = None
storage = None
event_alert = None
perfstore = None

def on_alert(body, message):
	print "Alert: %s" % body
	mrk = message.delivery_info['routing_key']
	if mrk == rk:
		global event_alert
		event_alert = body
	
def clean():
		storage.remove(rk)
		records = storage.find({'rk': rk}, namespace='events_log')
		storage.remove(records, namespace='events_log')
		
		try:
			perfstore.remove(name='test1mymetric')
		except:
			pass

class KnownValues(unittest.TestCase): 
	def setUp(self):
		self.rcvmsgbody = None

	def test_1_Init(self):
		global myamqp
		myamqp = camqp()
		myamqp.add_queue(	queue_name = "unittest_alerts",
							routing_keys = "#",
							callback = on_alert,
							exchange_name = myamqp.exchange_name_alerts)
		myamqp.start()
		time.sleep(1)
		
		global storage
		storage = cstorage(caccount(user="root", group="root"), namespace='events')
		
		global perfstore
		perfstore = pyperfstore2.manager()
		
		clean()
		
	def test_2_PubState(self):
		myamqp.publish(event, rk, exchange_name=myamqp.exchange_name_events)
		time.sleep(3)
		
	def test_3_Check_amqp2engines(self):
		record = storage.get(rk)
		revent = record.data
		
		if revent['component'] != event['component']:
			raise Exception('Invalid data ...')
			
		if revent['timestamp'] != event['timestamp']:
			raise Exception('Invalid data ...')
			
		if revent['state'] != event['state']:
			raise Exception('Invalid data ...')

		del event_alert['_id']
		
		# remove cps_state
		if len(event_alert['perf_data_array']) >= 2:
			del event_alert['perf_data_array'][1]

		event['perf_data_array'] = parse_perfdata(event['perf_data'])
		
		try:
			event['rk'] = event_alert['rk']
			event['event_id'] = event_alert['event_id']
		except:
			pass
			
		del event_alert['last_state_change']

		if event_alert != event:
			print "event_alert: %s" % event_alert
			print "event: %s" % event
			raise Exception('Invalid alert data ...')
			
		
	def test_4_Check_amqp2engines_archiver(self):
		## change state
		event['state'] = 1
		event['timestamp'] = int(time.time())
		myamqp.publish(event, rk, exchange_name=myamqp.exchange_name_events)
		time.sleep(3)
		
		records = storage.find({'event_id': rk}, sort='timestamp',  namespace='events_log')
		
		if len(records) != 2:
			raise Exception("Archiver don't work ...")
			
		revent = records[1].data
		
		if revent['state'] != event['state']:
			raise Exception('Invalid log state')
			
	def test_5_Check_amqp2engines_perfstore(self):
		values = perfstore.get_points(name='test1mymetric', tstart=int(time.time() - 10), tstop=int(time.time()))
		
		if len(values) != 2:
			raise Exception("Perfsore don't work ...")
			
		if values[1][1] != 1:
			raise Exception("Perfsore don't work ...")		
	
	def test_6_Check_webserver(self):	
		WS = cwebservices()
		WS.login('root', 'root')
		
		data = WS.get('/rest/events/event/%s' % rk)
		data = data[0]
		record = storage.get(rk)
		rdata = record.dump()
		
		WS.logout()
		
		if data['crecord_write_time'] != rdata['crecord_write_time']:
			raise Exception("Webservice don't work ...")

	def test_7_Check_collectd2event(self):

		import commands
		print "Restart collectd ..."
		commands.getstatusoutput("service collectd restart")

		print "Wait collectd events ..."
		i=0
		while i < 5:
			records = storage.find({'connector': 'collectd'})
			if len(records):
				break
			i+=1
			time.sleep(5)
		
		if not len(records):
			raise Exception("Collectd2event don't work ...")

	"""def test_8_Check_amqp2engines_media(self):
		
		## Publish BSON with file
		from bson.binary import Binary
		from bson import BSON
		
		event = cevent.forger(
					connector='test',
					connector_name='test',
					component='EUE',
					resource='Test',
					timestamp=None,
					source_type='resource',
					event_type='check',
					state=0
				)
		rk = cevent.get_routingkey(event)
		
		file_name = "logo_canopsis.png"
		sample_file_path = '/opt/canopsis/var/www/canopsis/themes/canopsis/resources/images/%s' % file_name 
		sample_binary = open(sample_file_path, 'rb').read()
		
		event['media_bin'] = Binary(sample_binary)
		event['media_name'] = file_name
		
		event = BSON.encode(event)

		myamqp.publish(event, rk, myamqp.exchange_name_events, content_encoding="binary", serializer=None)
		
		time.sleep(3)
		
		## Check engine
		from cfile import cfile
		from cfile import get_cfile
		
		record = storage.get(rk, namespace='events')
		
		if not record:
			raise("Impossible to find event in DB")
			
		file_id = record.data['media_id']
		rfile = get_cfile(file_id, storage)
		
		if not rfile.check():
			raise("Impossible to check file")
			
		if sample_binary != rfile.get():
			raise("Binary invalid")
			
		rfile.remove()
		
		from gridfs.errors import NoFile
		
		with self.assertRaises(NoFile):
			rfile.get()
			
		with self.assertRaises(KeyError):
			get_cfile(file_id, storage)
	"""	

	#def test_80_Check_Aps(self):
		#account = caccount(user="root", group="root")
		#storage = cstorage(account=account, namespace="task")

		#task_uuid = str(uuid.uuid4())	

		#data = json.loads('{"name": "%s","interval": {"seconds":1},"args": [],"kwargs":{"task":"task_node","method":"hostname"},"func_ref":"apschedulerlibs.aps_to_celery:launch_celery_task"}' % task_uuid)

		#record = crecord(account=account, storage=storage, data=data)

		#id = storage.put(record)
		#res = Popen(['service', 'apsd', 'restart'])
		#res.wait()

		#time.sleep(1)
		#found = False

		#regexp = re.compile('Job "%s \(trigger: interval\[0:00:01\], next run at: (.*)\)" executed successfully' % task_uuid)

		#for line in open("var/log/apsd.log"):
			#if regexp.search(line):
				#found = True

		#if not found:
			#raise Exception("Task not successfully added or executed")

		#storage.remove(id)		
		#res = Popen(['service', 'apsd', 'restart'])
		#res.wait()

		#time.sleep(1)

	def test_99_Disconnect(self):
		clean()
		myamqp.stop()
		myamqp.join()
		
if __name__ == "__main__":
	unittest.main(verbosity=2)
	
