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
from canopsis.stats.manager import Stats


class StatsManagerTest(TestCase):

    def setUp(self):
        self.stats_manager = Stats()
        self.fake_event = {
            'connector': 'fake_connector',
            'connector_name': 'fake_connector_name',
            'event_type': 'fake_event_type',
            'source_type': 'fake_source_type',
            'component': 'fake_component',
        }


class StatsTest(StatsManagerTest):
    def test_new_alert_event_count(self):

        event = self.fake_event.copy()
        devent = self.fake_event.copy()

        event['state'] = 0
        devent['state'] = 1

        metric = self.stats_manager.new_alert_event_count(event, devent)
        self.assertEqual(metric['value'], -1)

        event['state'] = 1
        devent['state'] = 0

        metric = self.stats_manager.new_alert_event_count(event, devent)
        self.assertEqual(metric['value'], 1)

        # When no metrics, event is not generated
        event['state'] = 0
        devent['state'] = 0

        metric = self.stats_manager.new_alert_event_count(event, devent)
        self.assertIsNone(metric)

        event['state'] = 1
        devent['state'] = 1

        metric = self.stats_manager.new_alert_event_count(event, devent)
        self.assertIsNone(metric)

    def test_solved_alarm_ack(self):

        devent = {}

        # devent is not ack, so incremented metric name contains not
        metric = self.stats_manager.solved_alarm_ack(devent)
        self.assertEqual(
            metric['metric'],
            'cps_solved_not_ack_alarms'
        )

        # Simulate ack in devent
        devent['ack'] = {'isAck': True}
        metric = self.stats_manager.solved_alarm_ack(devent)
        self.assertEqual(
            metric['metric'],
            'cps_solved_ack_alarms'
        )

    def test_compute_ack_alerts(self):
        event = self.fake_event.copy()
        devent = self.fake_event.copy()

        # Solve alert produce two metrics
        event['state'] = 0
        devent['state'] = 1

        metrics = self.stats_manager.compute_ack_alerts(event, devent)
        self.assertEqual(len(metrics), 2)

        # New alert metric only
        event['state'] = 1
        devent['state'] = 0

        metrics = self.stats_manager.compute_ack_alerts(event, devent)
        self.assertEqual(len(metrics), 1)

        # No metric as no state change
        event['state'] = 0
        devent['state'] = 0

        metrics = self.stats_manager.compute_ack_alerts(event, devent)
        self.assertEqual(len(metrics), 0)

        # No metric as no state change
        event['state'] = 1
        devent['state'] = 1

        metrics = self.stats_manager.compute_ack_alerts(event, devent)
        self.assertEqual(len(metrics), 0)

    def test_users_session_duration(self):
        # below methods are tested in the session manager
        def mock_gnis():
            return ['testvalue']

        def mock_gdstm(sessions):
            return sessions
        sm = self.stats_manager.session_manager

        sm.get_new_inactive_sessions = mock_gnis
        sm.get_delta_session_time_metrics = mock_gdstm

        # Just test the method acts as expected
        self.stats_manager.users_session_duration()
        self.assertEqual(self.stats_manager.perf_data_array, ['testvalue'])

    def test_event_add_by_source(self):
        for source in ('resource', 'component', 'test'):
            metrics = self.stats_manager.event_add_by_source(
                {'source_type': source}, True
            )
            self.assertEqual(len(metrics), 1)
            self.assertEqual(metrics[0]['value'], 1)
            self.assertEqual(
                metrics[0]['metric'],
                'cps_count_{}'.format(source)
            )

        # Not new event, no metric to produce
        metrics = self.stats_manager.event_add_by_source(
            {'source_type': 'resource'}, False
        )
        self.assertEqual(len(metrics), 0)

    def test_count_by_source_and_state(self):

        # New event, no metric to produce
        metrics = self.stats_manager.event_count_by_source_and_state(
            {}, {}, True
        )
        self.assertEqual(len(metrics), 0)

        # No state change, no metrics
        metrics = self.stats_manager.event_count_by_source_and_state(
            {'state': 1}, {'state': 1}, False
        )
        self.assertEqual(len(metrics), 0)

        # Test works
        metrics = self.stats_manager.event_count_by_source_and_state(
            {'state': 0, 'source_type': 'testsource'}, {'state': 1}, False
        )
        self.assertEqual(len(metrics), 2)
        self.assertEqual(metrics[0]['metric'], 'cps_states_testsource_minor')
        self.assertEqual(metrics[0]['value'], -1)
        self.assertEqual(metrics[1]['metric'], 'cps_states_testsource_info')
        self.assertEqual(metrics[1]['value'], 1)

        metrics = self.stats_manager.event_count_by_source_and_state(
            {'state': 3, 'source_type': 'source1'}, {'state': 2}, False
        )
        self.assertEqual(len(metrics), 2)
        self.assertEqual(metrics[0]['metric'], 'cps_states_source1_major')
        self.assertEqual(metrics[0]['value'], -1)
        self.assertEqual(metrics[1]['metric'], 'cps_states_source1_critical')
        self.assertEqual(metrics[1]['value'], 1)

    def test_event_add_by_state(self):
        for x, state in enumerate(['info', 'minor', 'major', 'critical']):
            metrics = self.stats_manager.event_add_by_state(
                {'state': x}, True
            )
            self.assertEqual(len(metrics), 1)
            self.assertEqual(metrics[0]['value'], 1)
            self.assertEqual(
                metrics[0]['metric'],
                'cps_states_{}'.format(state)
            )

        # Not new event, no metric to produce
        metrics = self.stats_manager.event_add_by_source(
            {'source_type': 'resource'}, False
        )
        self.assertEqual(len(metrics), 0)


if __name__ == '__main__':
    main()
