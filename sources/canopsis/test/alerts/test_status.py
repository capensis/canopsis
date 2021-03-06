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

from time import time

from canopsis.alerts.enums import AlarmField, States
from canopsis.alerts.status import (
    ONGOING, CANCELED, OFF,
    is_flapping, is_stealthy, compute_status, get_last_state, get_last_status,
    get_previous_step, is_keeped_state
)
from canopsis.check import Check

from base import BaseTest
import unittest
from canopsis.common import root_path
import xmlrunner


class TestStatus(BaseTest):
    def setUp(self):
        super(TestStatus, self).setUp()

        self.config_storage.put_element(
            element={
                '_id': 'test_config',
                'crecord_type': 'statusmanagement',
                'bagot_time': 3600,
                'bagot_freq': 10,
                'stealthy_time': 300,
                'restore_event': True
            },
            _id='test_config'
        )

        self.alarm = {
            AlarmField.state.value: None,
            AlarmField.status.value: None,
            AlarmField.ack.value: None,
            AlarmField.canceled.value: None,
            AlarmField.ticket.value: None,
            AlarmField.resolved.value: None,
            AlarmField.steps.value: [],
            AlarmField.tags.value: []
        }

    def test_is_flapping(self):
        self.alarm[AlarmField.steps.value] = [
            {
                '_t': 'stateinc',
                't': 0,
                'a': 'test',
                'm': 'test',
                'val': Check.CRITICAL
            },
            {
                '_t': 'statedec',
                't': 1,
                'a': 'test',
                'm': 'test',
                'val': Check.OK
            },
            {
                '_t': 'stateinc',
                't': 2,
                'a': 'test',
                'm': 'test',
                'val': Check.CRITICAL
            },
            {
                '_t': 'statedec',
                't': 3,
                'a': 'test',
                'm': 'test',
                'val': Check.OK
            },
            {
                '_t': 'stateinc',
                't': 4,
                'a': 'test',
                'm': 'test',
                'val': Check.CRITICAL
            },
            {
                '_t': 'statedec',
                't': 5,
                'a': 'test',
                'm': 'test',
                'val': Check.OK
            },
            {
                '_t': 'stateinc',
                't': 6,
                'a': 'test',
                'm': 'test',
                'val': Check.CRITICAL
            },
            {
                '_t': 'statedec',
                't': 7,
                'a': 'test',
                'm': 'test',
                'val': Check.OK
            },
            {
                '_t': 'stateinc',
                't': 8,
                'a': 'test',
                'm': 'test',
                'val': Check.CRITICAL
            },
            {
                '_t': 'statedec',
                't': 9,
                'a': 'test',
                'm': 'test',
                'val': Check.OK
            },
            {
                '_t': 'stateinc',
                't': 10,
                'a': 'test',
                'm': 'test',
                'val': Check.CRITICAL
            },
        ]

        self.alarm[AlarmField.state.value] = self.alarm[AlarmField.steps.value][-1]

        got = is_flapping(self.manager, self.alarm)
        self.assertTrue(got)

    def test_isnot_flapping(self):
        self.alarm[AlarmField.steps.value] = [
            {
                '_t': 'stateinc',
                't': 0,
                'a': 'test',
                'm': 'test',
                'val': Check.CRITICAL
            },
            {
                '_t': 'statedec',
                't': 1,
                'a': 'test',
                'm': 'test',
                'val': Check.OK
            }
        ]

        self.alarm[AlarmField.state.value] = self.alarm[AlarmField.steps.value][-1]

        got = is_flapping(self.manager, self.alarm)
        self.assertFalse(got)

    def test_is_stealthy(self):
        now = int(time())
        self.alarm[AlarmField.steps.value].append({
            '_t': 'stateinc',
            't': now - 1,
            'a': 'test',
            'm': 'test',
            'val': Check.CRITICAL
        })
        self.alarm[AlarmField.state.value] = {
            '_t': 'statedec',
            't': now,
            'a': 'test',
            'm': 'test',
            'val': Check.OK
        }

        got = is_stealthy(self.manager, self.alarm)

        self.assertTrue(got)

    def test_isnot_stealthy(self):
        self.alarm[AlarmField.steps.value].append({
            '_t': 'stateinc',
            't': 0,
            'a': 'test',
            'm': 'test',
            'val': Check.CRITICAL
        })
        self.alarm[AlarmField.state.value] = {
            '_t': 'statedec',
            't': 601,
            'a': 'test',
            'm': 'test',
            'val': Check.OK
        }

        got = is_stealthy(self.manager, self.alarm)

        self.assertFalse(got)

    def test_is_keeped_state(self):
        self.alarm[AlarmField.state.value] = {}
        self.alarm[AlarmField.state.value]['_t'] = States.changestate.value

        self.assertTrue(is_keeped_state(self.alarm))

    def test_isnot_keeped_state(self):
        self.alarm[AlarmField.state.value] = {}
        self.alarm[AlarmField.state.value]['_t'] = None

        self.assertFalse(is_keeped_state(self.alarm))

    def test_is_ongoing(self):
        self.alarm[AlarmField.state.value] = {
            '_t': 'stateinc',
            't': 0,
            'a': 'test',
            'm': 'test',
            'val': Check.CRITICAL
        }
        got = compute_status(self.manager, self.alarm)
        self.assertEqual(got, ONGOING)

    def test_is_canceled(self):
        self.alarm[AlarmField.canceled.value] = {
            '_t': 'cancel',
            't': 0,
            'a': 'test',
            'm': 'test'
        }
        got = compute_status(self.manager, self.alarm)
        self.assertEqual(got, CANCELED)

    def test_is_off(self):
        self.alarm[AlarmField.state.value] = {
            '_t': 'statedec',
            't': 0,
            'a': 'test',
            'm': 'test',
            'val': Check.OK
        }
        got = compute_status(self.manager, self.alarm)
        self.assertEqual(got, OFF)

    def test_get_last_state(self):
        got = get_last_state(self.alarm)
        self.assertEqual(got, Check.OK)

        self.alarm[AlarmField.state.value] = {
            '_t': 'stateinc',
            't': 0,
            'a': 'test',
            'm': 'test',
            'val': Check.CRITICAL
        }

        got = get_last_state(self.alarm)
        self.assertEqual(got, Check.CRITICAL)

    def test_get_last_status(self):
        got = get_last_status(self.alarm)
        self.assertEqual(got, OFF)

        self.alarm[AlarmField.status.value] = {
            '_t': 'statusinc',
            't': 0,
            'a': 'test',
            'm': 'test',
            'val': ONGOING
        }

        got = get_last_status(self.alarm)
        self.assertEqual(got, ONGOING)

    def test_get_previous_step(self):
        expected = {
            '_t': 'teststep',
            't': 0,
            'a': 'test',
            'm': 'test'
        }
        self.alarm[AlarmField.steps.value].append(expected)

        step = get_previous_step(self.alarm, 'teststep')

        self.assertTrue(expected is step)

    def test_bagot_issue_482(self):
        """
        With this issue, if the bagot_freq and bagot_time are both set to 1, the issue bagots
        :return:
        """
        self.manager.config_data.set('bagot_freq', 1)
        self.manager.config_data.set('bagot_time', 1)
        self.alarm[AlarmField.steps.value] = [
            {
                '_t': 'stateinc',
                't': 0,
                'a': 'test',
                'm': 'test',
                'val': Check.CRITICAL
            },
            {
                '_t': 'statedec',
                't': 1,
                'a': 'test',
                'm': 'test',
                'val': Check.OK
            }]

        self.alarm[AlarmField.state.value] = {
                '_t': 'statedec',
                't': 1,
                'a': 'test',
                'm': 'test',
                'val': Check.OK
            }

        got = is_flapping(self.manager, self.alarm)
        self.assertFalse(got)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
