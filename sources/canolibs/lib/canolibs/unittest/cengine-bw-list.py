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

class KnownValues(unittest.TestCase): 
	def setUp(self):
		pass

	def test_01_Init(self):
		module = __import__('bwlist')		
		engine = module.engine(logging_level=logging.DEBUG)
		
		events = [
			{'result': True, 'connector': 'nagios'},
			{'result': False, 'connector': 'collectd'},
			{'result': True, 'connector': 'lapin'}
		]
		engine.configuration = {'rules': [
			{'filter':'connector', 'pattern': 'nagios', 'action': 'pass'},
			{'filter':'connector', 'pattern': 'collectd', 'action': 'drop'},
		]}

		for event in events:
			
			if not event['result'] or engine.work(event) == event:
				result = '[OK]'
			else:
				result = '[KO]'
				raise Exception('Rule not work')
				
			print result + ' TEST EVENT ', event  
	

if __name__ == "__main__":

	unittest.main()
	


