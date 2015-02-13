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
from canopsis.sla.sla import Sla


class MongoQuery(object):
    def find(self, query, projection):
        return []

    def find_one(self, query, projection, sort=[]):
        return None


class MockStorage(object):
    def update(self, id, document):
        pass

    def get_backend(self, collection):
        return MongoQuery()


class KnownValues(unittest.TestCase):

    def setUp(self):
        pass

    def get_sla(
        self,
        rk='mock.rk',
        template='template sla',
        timewindow=60,
        warning=80,
        critical=60,
        alert_level='minor',
        display_name='mysla'
    ):
        sla = Sla(
            MockStorage(),
            rk,
            template,
            timewindow,
            warning,
            critical,
            alert_level,
            display_name,
        )
        return sla

    def test_01_init(self):
        self.get_sla()

    def test_compute_sla_output(self):

        # It should generate a specific string with given data structure
        sla = self.get_sla()
        template = '[OFF],[MINOR],[MAJOR],[CRITICAL],[ALERTS],[TSTART]'
        output = sla.compute_output(
            template,
            {
                0: 0,
                1: 0.012,
                2: 0.0223333,
                3: 0.03
            },
            0.98,
            1423753091
        )
        self.assertEqual(output, '0.00,1.20,2.23,3.00,98.00,2015-02-12 15:58:11')

        # User may not fill all fields, alert should always compute properly
        template = '[OFF] - [MAJOR] - [ALERTS]'
        output = sla.compute_output(
            template,
            {
                0: 0.11,
                1: 0.012,
                2: 0.11,
                3: 0.03
            },
            0.97,
            1423753091
        )
        self.assertEqual(output, '11.00 - 11.00 - 97.00')

    def test_prepare_event(self):

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


if __name__ == "__main__":
    unittest.main(verbosity=2)


