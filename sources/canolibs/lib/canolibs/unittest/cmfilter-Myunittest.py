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
import json

import cmfilter

event = {'connector': 'cengine', 'resource': 'Engine_perfstore2_rotate', 'event_type': 'check', 'long_output': None, 'timestamp': 1378713357, 'component': 'wpain-laptop', 'state_type': 1, 'source_type': 'resource', 'state': 0, 'connector_name': 'engine', 'output': '21.10 evt/sec, 0.02050 sec/evt', 'perf_data_array': [{'metric': 'cps_evt_per_sec', 'value': 21.1, 'unit': 'evt', 'retention': 3600}, {'metric': 'cps_sec_per_evt', 'value': 0.0205, 'warn': 0.6, 'crit': 0.9, 'unit': 's', 'retention': 3600}]}

class KnownValues(unittest.TestCase): 
	def setUp(self):
		pass

	def test_01_check(self):
		filter1 = {'connector': 'cengine'}
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match)

		filter1 = {'connector': { '$eq': 'cengine' } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match)

		filter1 = {'connector': 'cengidddddne'}
		match = cmfilter.check(filter1, event)	
		self.assertFalse(match)

		filter1 = {'connector': 'cengine', 'event_type': 'check'}
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match)

		filter1 = {'$and': [ {'connector': 'cengine'}, {'event_type': 'check'}, {'event_type': 'check'} ] }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match)

		filter1 = {'$or': [ {'connector': 'cenginddddde'}, {'event_type': 'check'},  {'event_type': 'checkkkkk'} ] }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match)

		filter1 = {'$or': [ { '$and': [ {'connector': 'cenginddddde'}, {'event_type': 'check'} ] },  {'event_type': 'checkkkkk'} ] }
		match = cmfilter.check(filter1, event)	
		self.assertFalse(match)

if __name__ == "__main__":
	unittest.main(verbosity=2)
	


