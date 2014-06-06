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

from crecord import crecord
from cstorage import get_storage
from caccount import caccount

class KnownValues(unittest.TestCase):
	def setUp(self):
		self.engine = alertcounter.engine(
			logging_level=logging.INFO,
			camqp_custom=camqpmock.CamqpMock
		)
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



if __name__ == "__main__":
	unittest.main()
