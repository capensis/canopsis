#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals
import time
import unittest
from canopsis.alarms.adapters import (
    AlarmAdapter, make_alarm_from_mongo, make_alarm_step_from_mongo
)
from canopsis.alarms.models import AlarmStep, AlarmIdentity
from canopsis.common import root_path
import xmlrunner


class AlarmAdapterTest(unittest.TestCase):

    @classmethod
    def setUp(self):

        self.adapter = AlarmAdapter(mongo_client={})

    def test_make_alarm_from_mongo(self):
        now = time.time()
        state_dict = {
            'a': 'Arthur Dent',
            'm': 'The answer to life, the universe and everything',
            '_t': 'statedec',
            't': time.time(),
            'val': 42
        }

        status_dict = {
            'a': 'Arthur Dent',
            'm': 'The answer to life, the universe and everything',
            '_t': 'stateinc',
            't': now,
            'val': 2
        }

        alarm_dict = {
            'v': {
                'state': state_dict,
                'status': status_dict,
                'connector': 'snmp',
                'connector_name': 'snmp',
                'component': 'comp1',
                'resource': 'res1'
            },
            '_id': 'abc',
            'd': now,

        }

        alarm = make_alarm_from_mongo(alarm_dict)

        self.assertIsInstance(alarm.status, AlarmStep)
        self.assertIsInstance(alarm.state, AlarmStep)
        self.assertIsInstance(alarm.identity, AlarmIdentity)
        self.assertIsNone(alarm.ack)
        self.assertIsNone(alarm.resolved)
        self.assertIsNone(alarm.ticket)
        self.assertIsNone(alarm.snooze)
        self.assertIsNone(alarm.canceled)
        self.assertIsNone(alarm.alarm_filter)
        self.assertEquals(0, len(alarm.steps))
        self.assertEquals("{}/{}".format(alarm_dict['v']['resource'],
                                         alarm_dict['v']['component']),
                          alarm.identity.display_name())

        step_dict = {
            'a': 'Arthur Dent',
            'm': 'The answer to life, the universe and everything',
            '_t': 'ack',
            't': now,
            'val': 2
        }
        alarm_dict['v']['ack'] = step_dict
        alarm_dict['v']['snooze'] = step_dict
        alarm_dict['v']['ticket'] = step_dict
        alarm_dict['v']['canceled'] = step_dict
        alarm_dict['v']['steps'] = [step_dict]
        al2 = make_alarm_from_mongo(alarm_dict)
        self.assertIsInstance(al2.ack, AlarmStep)
        al2 = make_alarm_from_mongo(alarm_dict)
        self.assertIsInstance(al2.snooze, AlarmStep)
        self.assertIsInstance(al2.ticket, AlarmStep)
        self.assertIsInstance(al2.canceled, AlarmStep)
        self.assertIsInstance(al2.steps, list)

    def test_make_alarm_step_from_mongo(self):

        test_dict = {
            'a': 'Arthur Dent',
            'm': 'The answer to life, the universe and everything',
            '_t': 'statedec',
            't': time.time(),
            'val': 42
        }
        step = make_alarm_step_from_mongo(test_dict)
        self.assertIsInstance(step, AlarmStep)
        self.assertEquals(step.author, test_dict['a'])
        self.assertEquals(step.message, test_dict['m'])
        self.assertEquals(step.type_, test_dict['_t'])
        self.assertEquals(step.timestamp, test_dict['t'])
        self.assertEquals(step.value, test_dict['val'])

    def test_make_alarm_step_from_mongo_without_dict(self):
        test = 1

        with self.assertRaises(TypeError) as context:
            make_alarm_step_from_mongo(test)
        self.assertTrue("A dict is required." in context.exception)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
