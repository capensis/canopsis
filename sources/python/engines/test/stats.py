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
from logging import ERROR
from canopsis.engines.stats import engine


# Simple mock for canopsis amqp class
class CamqpMock(object):

    exchange_name_events = 'testexchange'

    def publish(self, event, rk, exchange_name):
        self.event = event
        self.rk = rk
        self.exchange_name = exchange_name


class StatsTest(TestCase):

    def setUp(self):
        self.engine = engine(logging_level=ERROR)

    def test_init(self):
        self.assertIsNotNone(getattr(self.engine, 'storage', None))
        states_str = {
            0: 'info',
            1: 'minor',
            2: 'major',
            3: 'critical'
        }
        self.assertEqual(getattr(self.engine, 'states_str', None), states_str)

    def test_compute_states(self):
        # Test all following metrics are produced and their values are positive
        metric_list = [
            'cps_states_info',
            'cps_states_minor',
            'cps_states_major',
            'cps_states_critical',
            'cps_count_resource',
            'cps_states_resource_info',
            'cps_states_resource_minor',
            'cps_states_resource_major',
            'cps_states_resource_critical',
            'cps_count_component',
            'cps_states_component_info',
            'cps_states_component_minor',
            'cps_states_component_major',
            'cps_states_component_critical'
        ]

        self.engine.perf_data_array = []
        self.engine.compute_states()
        for metric in self.engine.perf_data_array:
            self.assertIn(metric['metric'], metric_list)
            self.assertTrue(metric['value'] >= 0)
            metric_list.remove(metric['metric'])
        self.assertEqual(len(metric_list), 0)

    def test_publish_states(self):
        self.engine.perf_data_array = []
        self.engine.compute_states()
        self.engine.amqp = CamqpMock()
        self.engine.publish_states()
        self.assertEqual(
            len(self.engine.amqp.event['perf_data_array']),
            14
        )
        self.assertIn(
            'engine.engine.perf.resource',
            self.engine.amqp.rk
        )
        self.assertIn(
            'Engine_stats',
            self.engine.amqp.rk
        )



if __name__ == "__main__":
    main()
