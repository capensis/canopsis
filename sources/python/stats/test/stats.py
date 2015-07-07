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

        sevent = self.stats_manager.compute_ack_alerts(event, devent)
        self.assertEqual(len(sevent['perf_data_array']), 2)

        # New alert metric only
        event['state'] = 1
        devent['state'] = 0

        sevent = self.stats_manager.compute_ack_alerts(event, devent)
        self.assertEqual(len(sevent['perf_data_array']), 1)

        # No metric as no state change
        event['state'] = 0
        devent['state'] = 0

        sevent = self.stats_manager.compute_ack_alerts(event, devent)
        self.assertIsNone(sevent)

        # No metric as no state change
        event['state'] = 1
        devent['state'] = 1

        sevent = self.stats_manager.compute_ack_alerts(event, devent)
        self.assertIsNone(sevent)

    def test_add_metric(self):
        self.stats_manager.add_metric('name', 5)
        self.assertEqual(
            self.stats_manager.perf_data_array[0],
            {
                'metric': 'name',
                'value': 5,
                'type': 'COUNTER'
            }
        )

        self.stats_manager.add_metric('name1', 6, mtype='MYTYPE')
        self.assertEqual(len(self.stats_manager.perf_data_array), 2)

        self.assertEqual(
            self.stats_manager.perf_data_array[1],
            {
                'metric': 'name1',
                'value': 6,
                'type': 'MYTYPE'
            }
        )

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

    def test_event_count_by_source(self):
        def mockfind(query='Q', with_count=True):
            return 'fakecursor', 5

        self.stats_manager.event_manager.find = mockfind

        self.stats_manager.event_count_by_source()
        perfs = self.stats_manager.perf_data_array
        self.assertEqual(len(perfs), 2)
        self.assertEqual([x['metric'] for x in perfs], [
            'cps_count_resource',
            'cps_count_component',
        ])
        self.assertEqual([x['value'] for x in perfs], [5, 5])

    def test_event_count_by_source_and_state(self):
        def mockfind(query='Q', with_count=True):
            return 'fakecursor', 5

        self.stats_manager.event_manager.find = mockfind

        self.stats_manager.event_count_by_source_and_state()
        perfs = self.stats_manager.perf_data_array
        self.assertEqual(len(perfs), 8)
        self.assertEqual([x['metric'] for x in perfs], [
            'cps_states_resource_info',
            'cps_states_resource_minor',
            'cps_states_resource_major',
            'cps_states_resource_critical',
            'cps_states_component_info',
            'cps_states_component_minor',
            'cps_states_component_major',
            'cps_states_component_critical'
        ])

    def test_event_count_by_state(self):
        def mockfind(query='Q', with_count=True):
            return 'fakecursor', 5

        self.stats_manager.event_manager.find = mockfind

        self.stats_manager.event_count_by_state()
        perfs = self.stats_manager.perf_data_array
        self.assertEqual(len(perfs), 4)
        self.assertEqual([x['metric'] for x in perfs], [
            'cps_states_info',
            'cps_states_minor',
            'cps_states_major',
            'cps_states_critical'
        ])


if __name__ == '__main__':
    main()
