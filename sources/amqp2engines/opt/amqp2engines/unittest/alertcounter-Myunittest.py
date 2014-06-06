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
import sys, os
import logging

sys.path.append(os.path.expanduser('~/opt/amqp2engines/engines/'))
sys.path.append(os.path.expanduser('~/lib/canolibs/unittest/'))

import alertcounter
import camqpmock
import managermock

from crecord import crecord
from cstorage import get_storage
from caccount import caccount

class KnownValues(unittest.TestCase):
	def setUp(self):
		self.engine = alertcounter.engine(
			logging_level=logging.INFO,
		)
		# mocking the manager
		self.engine.amqp = camqpmock.CamqpMock(self.engine)
		self.engine.manager = managermock.ManagerMock(self.engine)
		self.storage = get_storage(namespace='object', account=caccount(user="root", group="root"))


	def get_event(self):
		event = {
			'connector': 'test',
			'connector_name': 'test',
			'event_type': 'not a checked type',
			'source_type': 'source',
			'component': 'component',
			'state': 1,
			'state_type': 1
		}
		routing_key = "%s.%s.%s.%s.%s" % (event['connector'], event['connector_name'], event['event_type'], event['source_type'], event['component'])
		event['rk'] = routing_key
		return event


	def test_01_Work(self):

		event = self.get_event()

		# asserts event was not validated
		self.engine.work(event)
		assert(self.engine.amqp.events == [])

		# asserts event is being threaten thanks to it s type
		event['event_type'] = 'check'
		self.engine.work(event)
		assert(self.engine.amqp.events != [])

	def test_02_load_macros(self):
		macro_name = 'TEST_MACRO'
		self.storage.get_backend('object').update(
			{'crecord_type': 'slamacros'},
			{'$set': {'macro': macro_name }},
			upsert=True
		)

		#effective test
		self.engine.load_macro()
		assert(self.engine.MACRO == macro_name)

	def test_03_load_crits(self):
		pass
		"""
		crit_name = 'criticity'
		delay_value = 1

		crit = crecord({
			'_id':'test',
			'crit': crit_name,
			'delay': delay_value,
			'crecord_type': 'slacrit',
		}).dump()

		self.storage.get_backend('object').update(
			{'crecord_type': 'slacrit'},
			{'$set': crit},
			upsert=True
		)

		#effective test
		self.engine.load_crits()

		print self.engine.crits

		assert(self.engine.crits[crit_name] == delay_value)
		"""

	def test_04_load_crits(self):
		self.storage.get_backend('object').remove({'crecord_type': 'comment'})
		count = self.storage.get_backend('object').insert({
			'crecord_type': 'comment',
			'referer_event_rks': [{'rk': 'test_rk_1'}]
		})
		self.engine.reload_ack_comments()
		assert(self.engine.comments['test_rk_1'] == 1)

	def test_05_perfdata_key(self):
		key = self.engine.perfdata_key({'co': 'co', 're': 're', 'me': 'me'})
		assert(key == 'coreme')

		key = self.engine.perfdata_key({'co': 'co', 'me': 'me'})
		assert(key == 'come')

		key = self.engine.perfdata_key({})
		assert(key == 'missing component or metric key')


	def test_06_increment_counter(self):
		meta = {'co': 'co', 're': 're', 'me': 'me'}
		self.engine.increment_counter(meta, 1)
		assert(self.engine.manager.data.pop() == {'meta_data': 'meta_data', 'name':
			u'coreme', 'value': 1})
		self.engine.increment_counter({'co': 'co', 're': 're', 'me': 'me'}, 1)
		del meta['re']
		self.engine.increment_counter(meta, 2)
		assert(self.engine.manager.data.pop() == {'meta_data': 'meta_data', 'name': u'come', 'value': 2})

	def test_07_update_global_counter(self):
		#generated metrics names are listed below.
		truth_table = {
			"__canopsis__cps_statechange": [1,1,1,1],
			"__canopsis__cps_statechange_hard": [0,1,1,1],
			"__canopsis__cps_statechange_soft": [0,0,0,0],
			"__canopsis__cps_statechange_0": [1,0,0,0],
			"__canopsis__cps_statechange_1": [0,1,0,0],
			"__canopsis__cps_statechange_2": [0,0,1,0],
			"__canopsis__cps_statechange_3": [0,0,0,1],
			"__canopsis__cps_statechange_nok": [0,1,1,1]
		}


		#data driven testing
		def ugc_each_status(state):

			self.engine.update_global_counter({'state': state, 'resource': 'resource'})
			event = self.engine.amqp.events.pop()

			assert(event['state'] == 0)
			assert(event['connector'] == 'cengine')
			assert(event['connector_name'] == self.engine.etype)
			assert(event['source_type'] == 'resource')
			assert(event['component'] == alertcounter.INTERNAL_COMPONENT)
			assert(event['resource'] == None)
			#Let test if they are all generated
			while self.engine.manager.data:
				metric = self.engine.manager.data.pop()
				assert(metric['name'] in truth_table)
				#using state as postition in truth table
				assert(truth_table[metric['name']][state] == metric['value'])

		# all statuses : ok, warning, error, unknown
		for state in xrange(4):
			ugc_each_status(state)

		host_group = 'test_host_group'
		self.engine.update_global_counter({'state': state, 'resource': 'resource', 'hostgroups': [host_group]})

		#8 basic metrics + 8 for hostgroup
		assert(len(self.engine.manager.data) == 16)
		self.engine.manager.data = []

		while self.engine.amqp.events:
			event = self.engine.amqp.events.pop()
			assert(event['resource'] == host_group)


if __name__ == "__main__":
	unittest.main()
