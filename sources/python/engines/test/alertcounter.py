#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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
import sys
import os
import logging
import time

sys.path.append(os.path.expanduser('~/opt/amqp2engines/engines/'))
sys.path.append(os.path.expanduser('~/lib/canolibs/unittest/'))

import alertcounter
import camqpmock
import managermock

from canopsis.old.storage import get_storage
from canopsis.old.account import Account


#Mocking storage for some tests
class MockStorage(object):
    def find(self, **kwargs):
        return [{'crecord_name': 'name1'}, {'crecord_name': 'name2'}]


class KnownValues(unittest.TestCase):
    def setUp(self):
        self.engine = alertcounter.engine(
            logging_level=logging.WARNING,
        )

        # mocking the manager
        self.engine.amqp = camqpmock.CamqpMock(self.engine)
        self.engine.manager = managermock.ManagerMock(self.engine)
        self.storage = get_storage(namespace='object', account=Account(user="root", group="root"))


    """
    Tests methods has engine method names
    """

    def test_01_Work(self):

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

        # asserts event was not validated
        self.engine.work(event)
        self.assertTrue(self.engine.amqp.events == [])

        # asserts event is being threaten thanks to it s type
        event['event_type'] = 'check'
        self.engine.work(event)
        self.assertTrue(self.engine.amqp.events != [])

    def test_02_load_macros(self):
        macro_name = 'TEST_MACRO'
        self.storage.get_backend('object').update(
            {'crecord_type': 'slamacros'},
            {'$set': {'macro': macro_name }},
            upsert=True
        )

        #effective test
        self.engine.load_macro()
        self.assertTrue(self.engine.MACRO == macro_name)


    def test_04_load_crits(self):
        self.storage.get_backend('object').remove({'crecord_type': 'comment'})
        count = self.storage.get_backend('object').insert({
            'crecord_type': 'comment',
            'referer_event_rks': [{'rk': 'test_rk_1'}]
        })
        self.engine.reload_ack_comments()
        self.assertTrue(self.engine.comments['test_rk_1'] == 1)

    def test_05_perfdata_key(self):
        key = self.engine.perfdata_key({'co': 'co', 're': 're', 'me': 'me'})
        self.assertTrue(key == 'coreme')

        key = self.engine.perfdata_key({'co': 'co', 'me': 'me'})
        self.assertTrue(key == 'come')

        key = self.engine.perfdata_key({})
        self.assertTrue(key == 'missing component or metric key')


    def test_06_increment_counter(self):
        meta = {'co': 'co', 're': 're', 'me': 'me'}
        self.engine.increment_counter(meta, 1)
        self.assertTrue(self.engine.manager.data.pop() == {'meta_data': 'meta_data', 'name':
            u'coreme', 'value': 1})
        self.engine.increment_counter({'co': 'co', 're': 're', 'me': 'me'}, 1)
        del meta['re']
        self.engine.increment_counter(meta, 2)
        self.assertTrue(self.engine.manager.data.pop() == {'meta_data': 'meta_data', 'name': u'come', 'value': 2})

    def test_07_update_global_counter_and_count_alerts(self):
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

            self.assertEqual(event['state'], 0)
            self.assertEqual(event['connector'], 'cengine')
            self.assertEqual(event['connector_name'], self.engine.etype)
            self.assertEqual(event['source_type'], 'resource')
            self.assertEqual(event['component'], alertcounter.INTERNAL_COMPONENT)
            self.assertEqual(event['resource'], None)
            #Let test if they are all generated
            while self.engine.manager.data:
                metric = self.engine.manager.data.pop()
                self.assertTrue(metric['name'] in truth_table)
                #using state as postition in truth table
                self.assertEqual(truth_table[metric['name']][state], metric['value'])

        # all statuses : ok, warning, error, unknown
        for state in xrange(4):
            ugc_each_status(state)

        host_group = 'test_host_group'
        self.engine.update_global_counter({'state': state, 'resource': 'resource', 'hostgroups': [host_group]})

        #8 basic metrics + 8 for hostgroup
        self.assertEqual(len(self.engine.manager.data), 16)
        #reset data
        self.engine.manager.data = []

        while self.engine.amqp.events:
            event = self.engine.amqp.events.pop()
            self.assertTrue(event['resource'] == host_group)



    def test_08_count_sla(self):
        truth_table = {
            "__canopsis__cps_sla_slatype_slaname_out": [1, 0, 0],
            "__canopsis__cps_sla_slatype_slaname_ok": [0, 0, 1],
            "__canopsis__cps_sla_slatype_slaname_nok": [0, 1, 0],
        }
        now = time.time()
        slatype = 'slatype'
        slaname = 'slaname'
        event = {'last_state_change': 1, 'state': 1}

        def test_sla(index,  delay=1):

            self.engine.count_sla(event,slatype, slaname, delay)
            while self.engine.manager.data:
                metric = self.engine.manager.data.pop()
                self.assertTrue(metric['name'] in truth_table)
                #using state as postition in truth table
                self.assertEqual(truth_table[metric['name']][index], metric['value'])

        #Test general cases
        self.storage.get_backend('entities').remove({'type': 'ack'})
        test_sla(0)
        self.storage.get_backend('entities').insert({'type': 'ack', 'timestamp': 2})
        test_sla(1)
        test_sla(2, delay=now + 1)


        #test other stuff like hostgroup care and event state == 0 that publish a new metric
        self.storage.get_backend('entities').remove({'type': 'ack'})

        event ['hostgroups'] = ['hostgroup_test']
        self.engine.count_sla(event,slatype, slaname, 1)
        self.assertEqual(len(self.engine.manager.data), 6)

        event ['hostgroups'] = ['hostgroup_test']
        self.engine.count_sla(event,slatype, slaname, 1)
        self.assertEqual(len(self.engine.manager.data), 12)




    def test_09_count_by_crits(self):
        self.engine.MACRO = 'MACRO'
        self.engine.crits = {'macro_value': 1}
        event = {
            'last_state_change':1,
            'previous_state':1,
            'state' : 0,
            'MACRO': 'macro_value'
        }

        #check metric name is properly built
        self.engine.count_by_crits(event)
        while self.engine.manager.data:
            metric = self.engine.manager.data.pop()
            self.assertTrue('warn' in metric['name'])

        event['previous_state'] = 2
        self.engine.count_by_crits(event)
        while self.engine.manager.data:
            metric = self.engine.manager.data.pop()
            self.assertTrue('crit' in metric['name'])

        #macro is a part of the metric name
        self.engine.crits['mock_test_macro'] = 1
        self.engine.count_by_crits(event)
        #test update other section
        for metric in self.engine.manager.data[3:]:
            metric = self.engine.manager.data.pop()
            self.assertTrue('mock_test_macro' in metric['name'])


    def test_10_count_by_type(self):
        # Simple metrics fetch
        def fetch_metrics():
            result = {}
            for metric in self.engine.manager.data:
                result[metric['name']] = metric['value']
            self.engine.manager.data = []
            return result

        #Test producing metrics when state != 0 and state_type == 1
        self.engine.count_by_type({'state_type':0, 'source_type': 'source', 'component': '__test__', 'state': 1}, 1)
        self.assertEqual(len(self.engine.manager.data), 0)
        self.engine.count_by_type({'state_type':0, 'source_type': 'source', 'component': '__test__', 'state': 0}, 1)
        self.assertEqual(len(self.engine.manager.data), 0)
        self.engine.count_by_type({'state_type':1, 'source_type': 'source', 'component': '__test__', 'state': 1}, 1)
        self.assertNotEqual(len(self.engine.manager.data), 0)

        #Gets metrics as dict
        metrics = fetch_metrics()
        self.assertEqual(metrics['__canopsis__cps_alerts_not_ack'], 1)
        self.assertEqual(metrics['__canopsis__cps_alerts_ack'], 0)
        self.assertEqual(metrics['__canopsis__cps_alerts_ack_by_host'], 0)

        #Test metrics when source type is 'component'
        self.engine.count_by_type({'state_type':0, 'source_type': 'component', 'component': '__test__', 'state': 1}, 1)
        metrics = fetch_metrics()
        self.assertEqual(metrics['__canopsis__cps_statechange_component'], 1)

        self.engine.count_by_type({'state_type':0, 'source_type': 'component', 'component': '__test__', 'state': 0}, 1)
        metrics = fetch_metrics()
        self.assertEqual(metrics['__canopsis__cps_statechange_component'], 0)

        #Test metrics when source type is 'resource'
        def return_false(event):
            return False
        def return_true(event):
            return True
        alertcounter.cevent.is_component_problem = return_false

        self.engine.count_by_type({'state_type':0, 'source_type': 'resource', 'component': '__test__', 'state': 1}, 1)
        metrics = fetch_metrics()
        self.assertEqual(metrics['__canopsis__cps_statechange_resource_by_component'], 0)
        self.assertEqual(metrics['__canopsis__cps_statechange_resource'], 1)

        alertcounter.cevent.is_component_problem = return_true
        self.engine.count_by_type({'state_type':0, 'source_type': 'resource', 'component': '__test__', 'state': 1}, 1)
        metrics = fetch_metrics()
        self.assertEqual(metrics['__canopsis__cps_statechange_resource_by_component'], 1)

        #Test host group metric generation
        self.engine.count_by_type({'state_type':1, 'source_type': 'source', 'component': '__test__', 'state': 1}, 1)
        metrics = fetch_metrics()
        self.assertEqual(len(metrics.keys()), 3)
        self.engine.count_by_type({'hostgroups': ['hg'], 'state_type':1, 'source_type': 'source', 'component': '__test__', 'state': 1}, 1)
        metrics = fetch_metrics()
        self.assertEqual(len(metrics.keys()), 6)


    def test_11_resolve_selectors_name(self):
        self.engine.storage = MockStorage()

        #Starting method test
        start = self.engine.last_resolv
        result = self.engine.resolve_selectors_name()
        #All Record name fetched into one list because refresh time is ok
        self.assertEqual(result[0], 'name1')
        self.assertEqual(result[1], 'name2')
        self.assertTrue(start < self.engine.last_resolv)

        # Refresh time has not append, no selector should be fetched yet
        self.engine.selectors_name = []
        # one minute delay is tested
        self.engine.last_resolv = time.time() - 59
        result = self.engine.resolve_selectors_name()
        self.assertEqual(result, [])

    def test_11_count_by_tags(self):
        self.engine.storage = MockStorage()
        # Selector type, nothing should append
        self.engine.count_by_tags({'event_type': 'selector'},0)
        self.assertEqual(len(self.engine.manager.data), 0)

        self.engine.count_by_tags({'tags':['tag1', 'tag2'], 'event_type': 'notselector'},0)
        self.assertEqual(len(self.engine.manager.data), 0)

        value = 0
        self.engine.count_by_tags({'state': 0, 'tags':['name1'], 'event_type': 'notselector'}, value)
        self.assertEqual(len(self.engine.manager.data), 8)
        result = {}
        for metric in self.engine.manager.data:
            self.assertTrue(metric['name'].startswith('name1'))
            self.assertEqual(metric['value'], value)
        self.engine.manager.data = []

if __name__ == "__main__":
    unittest.main()
