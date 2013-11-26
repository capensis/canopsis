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

event = {
	'connector': 'cengine',
	'resource': 'Engine_perfstore2_rotate',
	'event_type': 'check',
	'long_output': None,
	'timestamp': 1378713357,
	'component': 'wpain-laptop',
	'state_type': 1,
	'source_type': 'resource',
	'state': 0,
	'connector_name': 'engine',
	'output': '21.10 evt/sec, 0.02050 sec/evt',
	'perf_data_array': [{
		'metric': 'cps_evt_per_sec',
		'value': 21.1,
		'unit': 'evt',
		'retention': 3600
	},{
		'metric': 'cps_sec_per_evt',
		'value': 0.0205,
		'warn': 0.6,
		'crit': 0.9,
		'unit': 's',
		'retention': 3600
	}]
}

class KnownValues(unittest.TestCase): 
	def setUp(self):
		pass

	def test_01_simple(self):
		filter1 = {'connector': 'cengine'}
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = {'connector': 'cengidddddne'}
		match = cmfilter.check(filter1, event)	
		self.assertFalse(match, msg='Filter: %s' % filter1)

	def test_02_exists(self):
		filter1 = {'timestamp': { '$exists': True } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = {'timestamp': { '$exists': False } }
		match = cmfilter.check(filter1, event)	
		self.assertFalse(match, msg='Filter: %s' % filter1)

	def test_03_eq(self):
		filter1 = {'connector': { '$eq': 'cengine' } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = {'connector': { '$eq': 'cenginessssss' } }
		match = cmfilter.check(filter1, event)	
		self.assertFalse(match, msg='Filter: %s' % filter1)

		filter1 = {'timestamp': { '$eq': 1378713357 } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

	def test_04_gt_gte(self):
		filter1 = {'timestamp': { '$gt': 1378713357 } }
		match = cmfilter.check(filter1, event)	
		self.assertFalse(match, msg='Filter: %s' % filter1)

		filter1 = {'timestamp': { '$gte': 1378713357 } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = {'timestamp': { '$gt': 137871335 } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

	def test_05_in_nin(self):
		filter1 = {'timestamp': { '$in': [ 0, 5, 6, 1378713357 ] } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = {'timestamp': { '$nin': [ 0, 5, 6 ] } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

	def test_06_complex(self):
		filter1 = {'timestamp': { '$gt': 0, '$lt': 2378713357 } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = { '$and': [ {'timestamp': {'$gt': 0} } , {'timestamp': {'$lt': 2378713357} }] }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = { 'connector': { '$eq': 'cengine' },  'timestamp': { '$gt': 137871335 }}
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = { 'connector': { '$not': { '$eq': 'cccenngine' } } }
		match = cmfilter.check(filter1, event)
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = { 'connector': { '$not': { '$eq': 'cengine' } } }
		match = cmfilter.check(filter1, event)
		self.assertFalse(match, msg='Filter: %s' % filter1)

		filter1 = { '$nor': [ { 'connector': { '$eq': 'ccengine' } }, {'connector': { '$eq': 'cccengine' } } ] }
		match = cmfilter.check(filter1, event)
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = { '$nor': [ { 'connector': { '$eq': 'cengine' } }, {'connector': { '$eq': 'cccengine' } } ] }
		match = cmfilter.check(filter1, event)
		self.assertFalse(match, msg='Filter: %s' % filter1)
		
		filter1 = {'connector': 'cengine', 'event_type': 'check'}
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = {'$and': [ {'connector': 'cengine'}, {'event_type': 'check'}, {'event_type': 'check'} ] }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = {'$or': [ {'connector': 'cenginddddde'}, {'event_type': 'check'},  {'event_type': 'checkkkkk'} ] }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = {'$or': [ { '$and': [ {'connector': 'cenginddddde'}, {'event_type': 'check'} ] },  {'event_type': 'checkkkkk'} ] }
		match = cmfilter.check(filter1, event)	
		self.assertFalse(match, msg='Filter: %s' % filter1)

	def test_07_all(self):
		filter1 = { 'connector': { '$all': [ 'cengine' ] } }
		match = cmfilter.check(filter1, event)
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = { 'connector': { '$all': [ 'cengine', 'ccengine' ] } }
		match = cmfilter.check(filter1, event)
		self.assertFalse(match, msg='Filter: %s' % filter1)

	def test_08_regex(self):
		filter1 = { 'connector': { '$regex': 'c.ngInE' } }
		match = cmfilter.check(filter1, event)	
		self.assertFalse(match, msg='Filter: %s' % filter1)

		filter1 = { 'connector': { '$regex': 'c.ngInE', '$options': 'i' } }
		match = cmfilter.check(filter1, event)	
		self.assertTrue(match, msg='Filter: %s' % filter1)

		filter1 = { 'connector': { '$regex': 'c..ngine', '$options': 'i' } }
		match = cmfilter.check(filter1, event)	
		self.assertFalse(match, msg='Filter: %s' % filter1)

if __name__ == "__main__":
	unittest.main(verbosity=2)

