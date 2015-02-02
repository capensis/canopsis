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


from time import time, sleep
import unittest
from canopsis.old.selector import Selector
from canopsis.old.record import Record
from canopsis.sla import Sla


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

        tw_prev_state = 0
        result = sla.compute_sla(2, info, 0, tw_prev_state)
        result1 = result[0] - 0.5
        result2 = result[1] - 0.5

        # sla measure should be about 50% each (approx 5%)
        self.assertTrue(result1 < 0.05)
        self.assertTrue(result2 < 0.05)

    def test_05_compute_sla_missing_time(self):
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

        prev_state_tw_start = 0
        result = sla.compute_sla(100, info, 0, prev_state_tw_start)

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
        prev_state_tw_start = 1
        result = sla.compute_sla(100, info, minor, prev_state_tw_start)
        # result 1 is 30% because of start missing time and result 2 is 70%
        # we got a total of about 100%
        result1 = result[0] - 0.2
        result2 = result[2] - 0.7
        result3 = result[1] - 0.1

        # sla measure (approx 1%)
        self.assertTrue(abs(result1) < 0.01)
        self.assertTrue(abs(result2) < 0.01)
        self.assertTrue(abs(result3) < 0.01)

    def test_06_compute_timewindow_start_state(self):
        # It should compute state at start of timewindow on update
        # It should add entries when a change state appends
        now = time()
        sla = self.get_sla()

        # It should produce the 1 state as sla information is out of timewindow
        info = self.sla_information.copy()
        info['1'] = [{'start': now - 15, 'stop': now - 14}]
        tw_start_state = sla.update_sla_information(10, 0, 0, info)
        self.assertEqual(tw_start_state, 1)

        # It should produce the 2 state as sla information is out of timewindow
        info = self.sla_information.copy()
        info['2'] = [{'start': now - 15, 'stop': now - 14}]
        tw_start_state = sla.update_sla_information(10, 0, 0, info)
        self.assertEqual(tw_start_state, 2)

        # It should produce the None state as sla information is in timewindow
        timewindow = 20
        info = self.sla_information.copy()
        info['2'] = [{'start': now - 15, 'stop': now - 14}]
        info['3'] = [{'start': now - 14, 'stop': now - 10}]
        tw_start_state = sla.update_sla_information(timewindow, 0, 0, info)
        self.assertEqual(tw_start_state, None)

        # It should produce the 3 state as sla information are in timewindow
        # and state 1 is the latest to be excluded
        info = self.sla_information.copy()
        info['2'] = [{'start': now - 15, 'stop': now - 14}]
        info['3'] = [{'start': now - 14, 'stop': now - 10}]
        tw_start_state = sla.update_sla_information(10, 0, 0, info)
        self.assertEqual(tw_start_state, 3)

    def test_07_compute_sla_output(self):

        # It should generate a specific string with given data structure
        sla = self.get_sla()
        template = '[OFF],[MINOR],[MAJOR],[CRITICAL],[ALERTS]'
        output = sla.compute_output(template, {
            0: 0,
            1: 0.012,
            2: 0.0223333,
            3: 0.03
        })
        self.assertEqual(output, '0.00,1.20,2.23,3.00,6.43')

        # User may not fill all fields, alert should always compute properly
        template = '[OFF] - [MAJOR] - [ALERTS]'
        output = sla.compute_output(template, {
            0: 0.11,
            1: 0.012,
            2: 0.11,
            3: 0.03
        })
        self.assertEqual(output, '11.00 - 11.00 - 15.20')

    def test_08_compute_sla_output_accuracy(self):
        # It should make sla repartition to 25% on each state
        sla = self.get_sla()
        now = time()
        timewindow = 40
        info = {
            u'0': [{u'start': now - 40, u'stop': now - 30}],
            u'1': [{u'start': now - 30, u'stop': now - 20}],
            u'2': [{u'start': now - 20, u'stop': now - 10}],
            u'3': [{u'start': now - 10, u'stop': now}]
        }

        result = sla.compute_sla(timewindow, info, 0, None)

        # we got a total of about 100%
        result0 = result[0] - 0.25
        result1 = result[1] - 0.25
        result2 = result[2] - 0.25
        result3 = result[3] - 0.25

        # sla measure (approx 1%)
        self.assertTrue(abs(result0) < 0.01)
        self.assertTrue(abs(result1) < 0.01)
        self.assertTrue(abs(result2) < 0.01)
        self.assertTrue(abs(result3) < 0.01)

        # It should make sla repartition to 25% on each state
        # Becaule missing time is attribued to the previous state
        sla = self.get_sla()
        now = time()
        timewindow = 40
        info = {
            u'0': [],
            u'1': [{u'start': now - 30, u'stop': now - 20}],
            u'2': [{u'start': now - 20, u'stop': now - 10}],
            u'3': [{u'start': now - 10, u'stop': now}]
        }

        result = sla.compute_sla(timewindow, info, 0, None)

        # we got a total of about 100%
        result0 = result[0] - 0.25
        result1 = result[1] - 0.25
        result2 = result[2] - 0.25
        result3 = result[3] - 0.25

        # sla measure (approx 1%)
        self.assertTrue(abs(result0) < 0.01)
        self.assertTrue(abs(result1) < 0.01)
        self.assertTrue(abs(result2) < 0.01)
        self.assertTrue(abs(result3) < 0.01)

    def test_09_prepare_event(self):

        class MockSelector(object):
            display_name = 'test_display'

        sla = self.get_sla()
        measures = {0: 0, 1: 1, 2: 2, 3: 3}
        event = sla.prepare_event(MockSelector(), measures, 'output', 0)
        self.assertEqual(event['event_type'], 'sla')
        self.assertEqual(event['component'], 'test_display')
        self.assertEqual(event['source_type'], 'resource')
        self.assertEqual(event['resource'], 'sla')
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

    def test_10_compute_sla_state(self):

        # Create and test mock
        class MockSelector(object):
            def __init__(self, warning, critical):

                def warn():
                    return warning
                self.get_sla_warning = warn

                def crit():
                    return critical
                self.get_sla_critical = crit

        mock = MockSelector(2, 3)
        self.assertEqual(mock.get_sla_warning(), 2)
        self.assertEqual(mock.get_sla_critical(), 3)

        # It compute state for sla depending on given thresholds
        sla = self.get_sla()
        now = time()

        # Alerts mesaure is sum(M[1],M[2],M[3]) witch is tested below
        sla_measures = {
            0: 0.25,
            1: 0.25,
            2: 0.25,
            3: 0.25
        }

        # It should compute a state equal to 0 because no limit passed
        state = sla.compute_state(
            sla_measures,
            MockSelector(80, 90)
        )
        self.assertEqual(state, 0)

        # It should compute a state equal to 1
        state = sla.compute_state(
            sla_measures,
            MockSelector(70, 90)
        )
        self.assertEqual(state, 1)

        # It should compute a state equal to 3
        state = sla.compute_state(
            sla_measures,
            MockSelector(70, 71)
        )
        self.assertEqual(state, 3)


if __name__ == "__main__":
    unittest.main(verbosity=2)


