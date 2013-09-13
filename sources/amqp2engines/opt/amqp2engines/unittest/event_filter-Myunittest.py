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

import event_filter
from cengine import DROP

class KnownValues(unittest.TestCase): 
	def setUp(self):
		pass

	def test_01_Init(self):	
		engine = event_filter.engine(logging_level=logging.DEBUG)
		engine.drop_event_count = 0
		engine.configurations = [{
			'rules': [
				{'filter': {'connector': 'nagios'}	, 'action': 'pass'},
				{'filter': {'connector': 'collectd'}, 'action': 'drop'},
				{'filter': {'connector': 'priority'}, 'action': 'pass'},
				{'filter': {'test_field': { '$eq': 'cengine' } }, 'action': 'pass'},
				{'filter': {'test_field': { '$gt': 1378713357 } },'action': 'drop'},						
			], 
			'priority' : 2,
			'default_action': 'drop',
			'configuration': 'white_black_lists',
		},{
			'rules': [
				{'filter': {'connector': 'nagios'}	, 'action': 'pass'},
				{'filter': {'connector': 'second_rule'}, 'action': 'pass'},
				{'filter': {'connector': 'priority'}, 'action': 'drop'},
				{'filter': {'test_field': { '$eq': 'cengine' } }, 'action': 'pass'},
				{'filter': {'test_field': { '$gt': 1378713357 } },'action': 'drop'},						
			], 
			'priority' : 1,
			'default_action': 'drop',
			'configuration': 'white_black_lists',
		}
		]

		# Test normal behaviors
		event = {'connector': 'nagios'}


		self.assertTrue(engine.work(event) == event)

	
		event['connector'] = 'collectd'
		self.assertTrue(engine.work(event) == DROP)

		# second rule matched
		event['connector'] = 'second_rule'
		self.assertTrue(engine.work(event) == event)


		# Test default actions
		event['connector'] = 'default_drop'
		self.assertTrue(engine.work(event) == DROP)

		# Change default action
		for configuration in engine.configurations:
			configuration['default_action'] = 'pass'
		event['connector'] = 'default_pass'
		self.assertTrue(engine.work(event) == event)

		# rule priority validation sorted is the same used in beat method in the engine
		event['connector'] = 'priority'
		self.assertTrue(engine.work(event) == event)
		engine.configurations = sorted(engine.configurations, key=lambda x: x['priority'])
		self.assertTrue(engine.work(event) == DROP)


		# No configuration, default configuration is loaded
		engine.configurations = []
		self.assertTrue(engine.work(event) == event)



if __name__ == "__main__":

	unittest.main()
	


