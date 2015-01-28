#!/usr/bin/env python
# --------------------------------
# Copyright (c) 2014 'Capensis' [http://www.capensis.com]
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
# ----------------------

from time import time, sleep
import unittest
from canopsis.old.selector import Selector
from canopsis.old.record import Record
from canopsis.old.sla import Sla


class MockStorage(object):
    def update(self, id, document):
        pass


class KnownValues(unittest.TestCase):

    def setUp(self):
        self.selector_record = {
            "dostate": False,
            "sla_timewindow": None,
            "state_when_all_ack": "worststate",
            "display_name": "test",
            "state": None,
            "crecord_type": "selector",
            "sla_state": None,
            "sla_timewindow": None,
            "dosla": True,
            "output_tpl": '',
            "sla_output_tpl": ''
        }

        self.sla_information = {
            '0': [],
            '1': [],
            '2': [],
            '3': [],
        }

    def get_sla(self):
        self.record = Record(data=self.selector_record)
        self.selector = Selector(None, record=self.record)
        self.selector_event = {
            'state': 0
        }
        self.sla = Sla(
            self.selector,
            self.selector_event,
            MockStorage(),
            '123'
        )
        return self.sla

    def test_01_init(self):
        self.get_sla()

    def test_02_udpate_sla_info_nochange(self):
        sla = self.get_sla()
        sla.update_sla_information(
            100,
            0,
            0,
            self.sla_information
        )
        self.assertEqual(
            self.sla_information,
            {'0': [], '1': [], '2': [], '3': []}
        )

    def test_02_udpate_sla_info(self):
        # Apply time window rules to remove information from sla dict
        sla = self.get_sla()
        self.sla_information['0'].append({
            'start': time() - 101
        })
        self.sla_information['1'].append({
            'start': time() - 99
        })
        self.sla_information['2'].append({
            'start': time() - 101,
            'stop': time() - 99
        })
        self.sla_information['2'].append({
            'start': time() - 102,
            'stop': time() - 103
        })

        sla.update_sla_information(
            100,
            0,
            0,
            self.sla_information
        )
        self.assertEqual(len(self.sla_information['0']), 0)
        self.assertEqual(len(self.sla_information['1']), 1)
        self.assertEqual(len(self.sla_information['2']), 1)
        self.assertEqual(len(self.sla_information['3']), 0)

    def test_03_udpate_sla_info_many_entries(self):
        # It should behave properly when many slice for each state
        sla = self.get_sla()
        self.sla_information['0'].append({
            'start': time() - 101
        })
        self.sla_information['0'].append({
            'start': time() - 99
        })
        self.sla_information['0'].append({
            'start': time() - 98
        })

        sla.update_sla_information(
            100,
            0,
            0,
            self.sla_information
        )
        self.assertEqual(len(self.sla_information['0']), 2)

    def test_03_udpate_sla_info_many_entries(self):
        # It should add entries when a change state appends
        sla = self.get_sla()
        info = self.sla_information.copy()

        # Changestate first time, should create an entry for state 0
        sla.update_sla_information(100, 0, 1, info)
        self.assertEqual(len(info['0']), 1)
        self.assertIn('start', info['0'][0])
        self.assertNotIn('stop', info['0'][0])

        # Changestate second time, should create an entry for state 1
        sla.update_sla_information(100, 1, 0, info)
        self.assertEqual(len(info['0']), 1)
        self.assertEqual(len(info['1']), 1)
        self.assertIn('start', info['0'][0])
        self.assertIn('stop', info['0'][0])

        # Changestate third time, should create a second entry for state 0
        sla.update_sla_information(100, 0, 1, info)
        self.assertEqual(len(info['0']), 2)
        self.assertIn('start', info['1'][0])
        self.assertIn('stop', info['1'][0])

    def test_04_compute_sla(self):
        # It should add entries when a change state appends
        sla = self.get_sla()
        info = self.sla_information.copy()
        sla.update_sla_information(100, 0, 1, info)
        # Sleep allow better testing
        sleep(1)
        sla.update_sla_information(100, 1, 0, info)
        sleep(1)

        result = sla.compute_sla(2, info, 0)
        result1 = result[0] - 0.5
        result2 = result[1] - 0.5

        # sla measure should be about 50% each (approx 5%)
        self.assertTrue(result1 < 0.05)
        self.assertTrue(result2 < 0.05)

    def test_04_compute_sla_missing_time(self):
        # It should compute sla with specific values and expected results
        # Missing time at between last state change and now
        sla = self.get_sla()
        now = time()
        info = {
            u'0': [{u'start': now - 80, u'stop': now - 60}],
            u'1': [],
            u'2': [
                {u'start': now - 90, u'stop': now - 80},
                {u'start': now - 60}],  # missing time for stop here
            u'3': []
        }

        result = sla.compute_sla(100, info, 0)

        # result 1 is 30% because of start missing time and result 2 is 70%
        # we got a total of about 100%
        result1 = result[0] - 0.3
        result2 = result[2] - 0.7

        # sla measure (approx 1%)
        self.assertTrue(abs(result1) < 0.01)
        self.assertTrue(abs(result2) < 0.01)

        # same test but with missing time attribution to the minor state
        minor = 1
        sla = self.get_sla()
        now = time()
        info = {
            u'0': [{u'start': now - 80, u'stop': now - 60}],
            u'1': [],
            u'2': [
                {u'start': now - 90, u'stop': now - 80},
                {u'start': now - 60}],  # missing time for stop here
            u'3': []
        }

        # Set missing time on minor state
        result = sla.compute_sla(100, info, minor)
        # result 1 is 30% because of start missing time and result 2 is 70%
        # we got a total of about 100%
        result1 = result[0] - 0.2
        result2 = result[2] - 0.7
        result3 = result[1] - 0.1

        # sla measure (approx 1%)
        self.assertTrue(abs(result1) < 0.01)
        self.assertTrue(abs(result2) < 0.01)
        self.assertTrue(abs(result3) < 0.01)

    def test_04_compute_sla_output(self):

        # It should generate a specific string with given data structure
        sla = self.get_sla()
        template = '[OFF],[MINOR],[MAJOR],[CRITICAL],[ALERTS]'
        output = sla.compute_output(template, {0: 0, 1: 1.2, 2: 2.23333, 3: 3})
        self.assertEqual(output, '0,1.2,2.23333,3,6.43333')

        # User may not fill all fields, alert should always compute properly
        template = '[OFF] - [MAJOR] - [ALERTS]'
        output = sla.compute_output(template, {0: 11, 1: 1.2, 2: 11, 3: 3})
        self.assertEqual(output, '11 - 11 - 15.2')

    def test_04_prepare_event(self):

        class MockSelector(object):
            display_name = 'test_display'

        sla = self.get_sla()
        measures = {0: 0, 1: 1, 2: 2, 3: 3}
        event = sla.prepare_event(MockSelector(), measures, 'output')
        self.assertEqual(event['event_type'], 'sla')
        self.assertEqual(event['component'], 'sla')
        self.assertEqual(event['source_type'], 'resource')
        self.assertEqual(event['source_type'], 'resource')
        self.assertEqual(
            event['perf_data_array'],
            [
                {'metric': 'cps_sla_off', 'value': 0},
                {'metric': 'cps_sla_minor', 'value': 1},
                {'metric': 'cps_sla_major', 'value': 2},
                {'metric': 'cps_sla_critical', 'value': 3}
            ]
        )
        self.assertEqual(event['display_name'], 'test_display')
        self.assertEqual(event['connector'], 'sla')
        self.assertEqual(event['state'], 0)
        self.assertEqual(event['output'], 'output')
        self.assertEqual(event['connector_name'], 'engine')


if __name__ == "__main__":
    unittest.main(verbosity=2)
