#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

from unittest import TestCase, main

from canopsis.old.mfilter import check

event = {
    'connector': 'Engine',
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
    }, {
        'metric': 'cps_sec_per_evt',
        'value': 0.0205,
        'warn': 0.6,
        'crit': 0.9,
        'unit': 's',
        'retention': 3600
    }]
}


class KnownValues(TestCase):

    def test_01_simple(self):
        filter1 = {'connector': 'Engine'}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'connector': 'cengidddddne'}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

    def test_02_exists(self):
        filter1 = {'timestamp': {'$exists': True}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'timestamp': {'$exists': False}}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

    def test_03_eq(self):
        filter1 = {'connector': {'$eq': 'Engine'}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'$or': [{'state': {'$eq': 0}}]}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'connector': {'$eq': 'Enginessssss'}}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

        filter1 = {'timestamp': {'$eq': 1378713357}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

    def test_04_gt_gte(self):
        filter1 = {'timestamp': {'$gt': 1378713357}}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

        filter1 = {'timestamp': {'$gte': 1378713357}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'timestamp': {'$gt': 137871335}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

    def test_05_in_nin(self):
        filter1 = {'timestamp': {'$in': [0, 5, 6, 1378713357]}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'timestamp': {'$nin': [0, 5, 6]}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

    def test_06_complex(self):
        filter1 = {'timestamp': {'$gt': 0, '$lt': 2378713357}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'$and': [{'timestamp': {'$gt': 0}}, {'timestamp': {'$lt': 2378713357}}]}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'connector': {'$eq': 'Engine'}, 'timestamp': {'$gt': 137871335}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'connector': {'$not': {'$eq': 'cccenngine'}}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'connector': {'$not': {'$eq': 'Engine'}}}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

        filter1 = {'$nor': [{'connector': {'$eq': 'cEngine'}}, {'connector': {'$eq': 'ccEngine'}}]}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'$nor': [{'connector': {'$eq': 'Engine'}}, {'connector': {'$eq': 'ccEngine'}}]}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

        filter1 = {'connector': 'Engine', 'event_type': 'check'}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'$and': [{'connector': 'Engine'}, {'event_type': 'check'}, {'event_type': 'check'}]}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'$or': [{'connector': 'cenginddddde'}, {'event_type': 'check'},  {'event_type': 'checkkkkk'}]}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'$or': [{'$and': [{'connector': 'cenginddddde'}, {'event_type': 'check'}]},  {'event_type': 'checkkkkk'}]}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

    def test_07_all(self):
        filter1 = {'connector': {'$all': ['Engine']}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'connector': {'$all': ['Engine', 'cEngine']}}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

    def test_08_regex(self):
        filter1 = {'connector': {'$regex': '.ngInE'}}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

        filter1 = {'connector': {'$regex': '.ngInE', '$options': 'i'}}
        match = check(filter1, event)
        self.assertTrue(match, msg='Filter: %s' % filter1)

        filter1 = {'connector': {'$regex': '..ngine', '$options': 'i'}}
        match = check(filter1, event)
        self.assertFalse(match, msg='Filter: %s' % filter1)

if __name__ == "__main__":
    main(verbosity=2)
