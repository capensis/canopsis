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

#from __future__ import unicode_literals

from datetime import datetime
from operator import itemgetter
from time import time, mktime, sleep
from unittest import main

from canopsis.alerts import AlarmField, States
from canopsis.alerts.filter import AlarmFilter
from canopsis.alerts.manager import Alerts
from canopsis.alerts.status import OFF, STEALTHY, is_keeped_state
from canopsis.alerts.tasks import snooze
from canopsis.check import Check
from canopsis.timeserie.timewindow import get_offset_timewindow

from base import BaseTest


class TestManager(BaseTest):
    def test_config(self):
        self.assertEqual(self.manager.flapping_interval, 3600)
        self.assertEqual(self.manager.flapping_freq, 10)
        self.assertEqual(self.manager.flapping_persistant_steps, 10)
        self.assertEqual(self.manager.hard_limit, 100)
        self.assertEqual(self.manager.stealthy_interval, 300)
        self.assertEqual(self.manager.stealthy_show_duration, 600)
        self.assertEqual(self.manager.restore_event, True)

    def test_make_alarm(self):
        alarm_id = '/fake/alarm/id'

        alarm = self.manager.make_alarm(
            alarm_id,
            {
                'connector': 'ut-connector',
                'connector_name': 'ut-connector0',
                'component': 'c',
                'timestamp': 0
            }
        )
        self.assertTrue(alarm is not None)

    def test_get_alarms(self):
        storage = self.manager[Alerts.ALARM_STORAGE]

        alarm0_id = '/fake/alarm/id0'
        event0 = {
            'connector': 'ut',
            'connector_name': 'ut0',
            'component': 'c',
            'output': '...',
            'timestamp': 0
        }
        alarm0 = self.manager.make_alarm(
            alarm0_id,
            event0
        )
        alarm0 = self.manager.update_state(alarm0, 1, event0)
        new_value0 = alarm0[self.manager[Alerts.ALARM_STORAGE].VALUE]
        self.manager.update_current_alarm(alarm0, new_value0)

        alarm1_id = '/fake/alarm/id1'
        event1 = {
            'connector': 'ut',
            'connector_name': 'ut0',
            'component': 'c',
            'output': '...',
            'timestamp': 0
        }
        alarm1 = self.manager.make_alarm(
            alarm1_id,
            event1
        )
        alarm1 = self.manager.update_state(alarm1, 1, event1)
        new_value1 = alarm1[self.manager[Alerts.ALARM_STORAGE].VALUE]
        self.manager.update_current_alarm(alarm1, new_value1)

        # Case 1: unresolved alarms
        got = self.manager.get_alarms(resolved=False)
        ids = [a for a in got]

        self.assertTrue(alarm0_id in ids)
        self.assertTrue(alarm1_id in ids)

        alarm0 = got[alarm0_id][0]
        alarm0[storage.DATA_ID] = alarm0_id

        # Case 2: resolved alarms
        got = self.manager.get_alarms(resolved=True)
        ids = [a for a in got]

        self.assertFalse(alarm0_id in ids)
        self.assertFalse(alarm1_id in ids)

        # Case 3: with tags
        self.manager.update_current_alarm(
            alarm0,
            alarm0[storage.VALUE],
            tags='test'
        )

        got = self.manager.get_alarms(tags='test', resolved=False)
        ids = [a for a in got]

        self.assertTrue(alarm0_id in ids)
        self.assertFalse(alarm1_id in ids)

        # Case 4: without tags
        got = self.manager.get_alarms(
            exclude_tags='test',
            resolved=False
        )
        ids = [a for a in got]

        self.assertFalse(alarm0_id in ids)
        self.assertTrue(alarm1_id in ids)

    def test_get_current_alarm(self):
        alarm_id = '/fake/alarm/id'

        got = self.manager.get_current_alarm(alarm_id)
        self.assertTrue(got is None)

        event = {
            'connector': 'ut',
            'connector_name': 'ut0',
            'component': 'c',
            'output': '...',
            'timestamp': 0
        }
        alarm = self.manager.make_alarm(
            alarm_id,
            event
        )
        alarm = self.manager.update_state(alarm, 1, event)
        new_value = alarm[self.manager[Alerts.ALARM_STORAGE].VALUE]
        self.manager.update_current_alarm(alarm, new_value)

        got = self.manager.get_current_alarm(alarm_id)
        self.assertTrue(got is not None)

    def test_update_current_alarm(self):
        storage = self.manager[Alerts.ALARM_STORAGE]

        alarm_id = '/fake/alarm/id'
        alarm = self.manager.make_alarm(
            alarm_id,
            {
                'connector': 'ut-connector',
                'connector_name': 'ut-connector0',
                'component': 'c',
                'timestamp': 0
            }
        )

        value = alarm[storage.VALUE]

        value[AlarmField.state.value] = {'val': 0}

        self.manager.update_current_alarm(alarm, value, tags='test')

        alarm = self.manager.get_current_alarm(alarm_id)
        value = alarm[storage.VALUE]

        self.assertTrue(value[AlarmField.state.value] is not None)
        self.assertTrue('test' in value[AlarmField.tags.value])

    def test_resolve_alarms(self):
        storage = self.manager[Alerts.ALARM_STORAGE]

        alarm_id = '/fake/alarm/id'
        alarm = self.manager.make_alarm(
            alarm_id,
            {
                'connector': 'ut-connector',
                'connector_name': 'ut-connector0',
                'component': 'c',
                'timestamp': 0
            }
        )

        self.assertIsNotNone(alarm)

        value = alarm[storage.VALUE]
        value[AlarmField.status.value] = {
            't': 0,
            'val': OFF
        }

        self.manager.update_current_alarm(alarm, value)
        self.manager.resolve_alarms()

        alarm = self.manager.get_current_alarm(alarm_id)
        self.assertIsNone(alarm)

        alarm = storage.get(
            alarm_id,
            timewindow=get_offset_timewindow(),
            _filter={
                AlarmField.resolved.value: {'$exists': True}
            },
            limit=1
        )
        self.assertTrue(alarm)
        alarm = alarm[0]
        value = alarm[storage.VALUE]

        self.assertEqual(value[AlarmField.resolved.value],
                         value[AlarmField.status.value]['t'])

    def test_resolve_snoozes(self):
        storage = self.manager[Alerts.ALARM_STORAGE]

        event = {
            'connector': 'ut-connector',
            'connector_name': 'ut-connector0',
            'component': 'c',
            'duration': 0,
            'timestamp': 0
        }
        alarm_id = '/fake/alarm/id'
        alarm = self.manager.make_alarm(alarm_id, event)
        self.assertIsNotNone(alarm)
        alarm[AlarmField.steps.value] = []

        # Execute a snooze on this alarm !!!
        value = snooze(self.manager, alarm, 'author', 'message', event)
        self.manager.update_current_alarm(alarm, value)

        alarm = self.manager.get_current_alarm(alarm_id)
        self.assertIsNotNone(alarm)
        value = alarm[storage.VALUE]
        self.assertIsNotNone(value[AlarmField.snooze.value])

        self.manager.resolve_snoozes()

        alarm = self.manager.get_current_alarm(alarm_id)
        self.assertIsNotNone(alarm)
        value = alarm[storage.VALUE]
        self.assertIsNone(value[AlarmField.snooze.value])

    def test_resolve_stealthy(self):
        storage = self.manager[Alerts.ALARM_STORAGE]
        now = int(time()) - self.manager.stealthy_show_duration - 1

        alarm_id = '/fake/alarm/id'
        alarm = self.manager.make_alarm(
            alarm_id,
            {
                'connector': 'ut-connector',
                'connector_name': 'ut-connector0',
                'component': 'c',
                'timestamp': now
            }
        )
        self.assertIsNotNone(alarm)

        # Init stealthy state
        value = alarm[storage.VALUE]
        value[AlarmField.status.value] = {
            't': now,
            'val': STEALTHY
        }
        value[AlarmField.state.value] = {
            't': now,
            'val': Check.OK
        }
        value[AlarmField.steps.value] = [
            {
                '_t': 'stateinc',
                't': now - 1,
                'a': 'test',
                'm': 'test',
                'val': Check.CRITICAL
            },
            {
                '_t': 'statedec',
                't': now,
                'a': 'test',
                'm': 'test',
                'val': Check.OK
            }
        ]
        self.manager.update_current_alarm(alarm, value)

        self.manager.resolve_stealthy()

        alarm = storage.get(
            alarm_id,
            timewindow=get_offset_timewindow(),
            _filter={
                AlarmField.resolved.value: {'$exists': True}
            },
            limit=1
        )
        self.assertTrue(alarm)
        alarm = alarm[0]
        value = alarm[storage.VALUE]

        self.assertEqual(value[AlarmField.status.value]['val'], OFF)

    def test_change_of_state(self):
        alarm_id = '/fake/alarm/id'

        event = {
            'timestamp': 0,
            'connector': 'ut-connector',
            'connector_name': 'ut-connector0',
            'output': 'UT message'
        }

        alarm = self.manager.make_alarm(
            alarm_id,
            {
                'connector': 'ut-connector',
                'connector_name': 'ut-connector0',
                'component': 'c',
                'timestamp': 0
            }
        )
        alarm = self.manager.change_of_state(alarm, 0, 2, event)

        expected_state = {
            'a': 'ut-connector.ut-connector0',
            '_t': 'stateinc',
            'm': 'UT message',
            't': 0,
            'val': 2,
        }

        expected_status = {
            'a': 'ut-connector.ut-connector0',
            '_t': 'statusinc',
            'm': 'UT message',
            't': 0,
            'val': 1,
        }

        # Make sure no more steps are added
        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 2)

        self.assertEqual(alarm['value'][AlarmField.state.value], expected_state)
        self.assertEqual(alarm['value'][AlarmField.steps.value][0], expected_state)
        self.assertEqual(alarm['value'][AlarmField.status.value], expected_status)
        self.assertEqual(alarm['value'][AlarmField.steps.value][1], expected_status)

        alarm = self.manager.change_of_state(alarm, 2, 1, event)

        expected_state = {
            'a': 'ut-connector.ut-connector0',
            '_t': 'statedec',
            'm': 'UT message',
            't': 0,
            'val': 1,
        }

        # Make sure no more steps are added
        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 3)

        self.assertEqual(alarm['value'][AlarmField.state.value], expected_state)
        self.assertEqual(alarm['value'][AlarmField.steps.value][2], expected_state)

    def test_change_of_status(self):
        alarm_id = '/fake/alarm/id'

        event = {
            'timestamp': 0,
            'connector': 'ut-connector',
            'connector_name': 'ut-connector0',
            'output': 'UT message',
        }

        alarm = self.manager.make_alarm(
            alarm_id,
            {
                'connector': 'ut-connector',
                'connector_name': 'ut-connector0',
                'component': 'c',
                'timestamp': 0
            }
        )

        alarm = self.manager.change_of_status(alarm, 0, 1, event)

        expected_status = {
            'a': 'ut-connector.ut-connector0',
            '_t': 'statusinc',
            'm': 'UT message',
            't': 0,
            'val': 1,
        }

        self.assertEqual(alarm['value'][AlarmField.status.value], expected_status)

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 1)
        self.assertEqual(alarm['value'][AlarmField.steps.value][0], expected_status)

    def test_archive_state_nochange(self):
        alarm_id = 'ut-comp'

        event0 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.MINOR,
        }
        self.manager.archive(event0)

        alarm = self.manager.get_current_alarm(alarm_id)

        expected_state = {
            'a': 'test.test0',
            '_t': 'stateinc',
            'm': 'test message',
            't': 0,
            'val': 1,
        }

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 2)
        self.assertEqual(alarm['value'][AlarmField.steps.value][0], expected_state)
        self.assertEqual(alarm['value'][AlarmField.state.value], expected_state)

        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.MINOR,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 2)
        self.assertEqual(alarm['value'][AlarmField.steps.value][0], expected_state)
        self.assertEqual(alarm['value'][AlarmField.state.value], expected_state)

    def test_archive_state_changed(self):
        alarm_id = 'ut-comp'

        # Testing state creation
        event0 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.MINOR,
        }
        self.manager.archive(event0)

        alarm = self.manager.get_current_alarm(alarm_id)

        expected_state = {
            'a': 'test.test0',
            '_t': 'stateinc',
            'm': 'test message',
            't': 0,
            'val': 1,
        }

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 2)
        self.assertEqual(alarm['value'][AlarmField.steps.value][0], expected_state)
        self.assertEqual(alarm['value'][AlarmField.state.value], expected_state)

        # Testing state increase
        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.MAJOR,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        expected_state = {
            'a': 'test.test0',
            '_t': 'stateinc',
            'm': 'test message',
            't': 0,
            'val': 2,
        }

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 3)
        self.assertEqual(alarm['value'][AlarmField.steps.value][2], expected_state)
        self.assertEqual(alarm['value'][AlarmField.state.value], expected_state)

        # Testing keeped state
        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': States.changestate.value,
            'state': Check.MINOR,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 4)
        self.assertEqual(alarm['value'][AlarmField.state.value]['val'], 1)
        self.assertTrue(is_keeped_state(alarm['value']))

        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.CRITICAL,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 4)
        self.assertEqual(alarm['value'][AlarmField.state.value]['val'], 1)
        self.assertTrue(is_keeped_state(alarm['value']))

        # Disengaging keepstate
        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.OK,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 6)
        self.assertEqual(alarm['value'][AlarmField.state.value]['val'], 0)
        self.assertFalse(is_keeped_state(alarm['value']))

    def test_archive_status_nochange(self):
        alarm_id = 'ut-comp'

        event0 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.MINOR,
        }
        self.manager.archive(event0)

        alarm = self.manager.get_current_alarm(alarm_id)

        expected_status = {
            'a': 'test.test0',
            '_t': 'statusinc',
            'm': 'test message',
            't': 0,
            'val': 1,
        }

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 2)
        self.assertEqual(alarm['value'][AlarmField.steps.value][1], expected_status)
        self.assertEqual(alarm['value'][AlarmField.status.value], expected_status)

        # Force status to stealthy
        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.MAJOR,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 3)
        self.assertEqual(alarm['value'][AlarmField.steps.value][1], expected_status)
        self.assertEqual(alarm['value'][AlarmField.status.value], expected_status)

    def test_archive_status_changed(self):
        alarm_id = 'ut-comp'

        # Force status to minor
        event0 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.MINOR,
        }
        self.manager.archive(event0)

        alarm = self.manager.get_current_alarm(alarm_id)

        expected_status = {
            'a': 'test.test0',
            '_t': 'statusinc',
            'm': 'test message',
            't': 0,
            'val': Check.MINOR,
        }

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 2)
        self.assertEqual(alarm['value'][AlarmField.steps.value][1], expected_status)
        self.assertEqual(alarm['value'][AlarmField.status.value], expected_status)

        # Force status to stealthy
        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.OK,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        expected_status = {
            'a': 'test.test0',
            '_t': 'statusdec',
            'm': 'test message',
            't': 0,
            'val': Check.OK,
        }

        self.assertEqual(len(alarm['value'][AlarmField.steps.value]), 4)
        self.assertEqual(alarm['value'][AlarmField.steps.value][3], expected_status)
        self.assertEqual(alarm['value'][AlarmField.status.value], expected_status)

    def test_crop_flapping_steps(self):
        # Creating alarm /component/test/test0/ut-comp1
        KO = {
            'connector': 'test',
            'connector_name': 'test0',
            'source_type': 'component',
            'component': 'ut-comp1',
            'event_type': 'check',
            'state': '1',
            'output': '...'
        }

        OK = {
            'connector': 'test',
            'connector_name': 'test0',
            'source_type': 'component',
            'component': 'ut-comp1',
            'event_type': 'check',
            'state': Check.OK,
            'output': '...'
        }

        assocticket = {
            'connector': 'test',
            'connector_name': 'test0',
            'source_type': 'component',
            'component': 'ut-comp1',
            'author': 'user0',
            'event_type': 'declareticket',
            'ticket': 'ticket',
            'output': '...'
        }

        for i in range(9):
            KO['timestamp'] = 2 * i * 300
            self.manager.archive(KO)

            OK['timestamp'] = (2 * i + 1) * 300
            self.manager.archive(OK)

            assocticket['timestamp'] = (2 * i + 1) * 300 + 10
            self.manager.archive(assocticket)

        KO['timestamp'] = 2 * (i + 1) * 300
        self.manager.archive(KO)

        # At this point we have inserted 19 check events. The alarm has
        # changed its status to flapping after the 10th. There are only 9
        # state changes after the last status change. It means we should not
        # have any state crop.

        alarm_id = 'ut-comp1'
        docalarm = self.manager.get_current_alarm(alarm_id)

        self.assertIsNot(docalarm, None)

        alarm = docalarm[self.manager[Alerts.ALARM_STORAGE].VALUE]

        last_status_i = alarm[AlarmField.steps.value].index(alarm[AlarmField.status.value])

        state_steps = filter(
            lambda step: step['_t'] in ['stateinc', 'statedec'],
            alarm[AlarmField.steps.value][last_status_i + 1:]
        )
        self.assertEqual(len(state_steps), 9)

        # 4 KO + 4 OK + 5 assocticket + 1 KO = 14 steps
        all_steps = alarm[AlarmField.steps.value][last_status_i + 1:]
        self.assertEqual(len(all_steps), 14)

        # Creating alarm /component/test/test0/ut-comp2
        KO['component'] = 'ut-comp2'
        OK['component'] = 'ut-comp2'
        assocticket['component'] = 'ut-comp2'

        for i in range(10):
            KO['timestamp'] = 2 * i * 300
            self.manager.archive(KO)

            OK['timestamp'] = (2 * i + 1) * 300
            self.manager.archive(OK)

            assocticket['timestamp'] = (2 * i + 1) * 300 + 10
            self.manager.archive(assocticket)

        KO['timestamp'] = 2 * (i + 1) * 300
        self.manager.archive(KO)

        # 21 flapping checks inserted. 10 checks to trigger flapping status.
        # 11 state changes after this change of status. Expecting 1 state to
        # be cropped.

        alarm_id = 'ut-comp2'
        docalarm = self.manager.get_current_alarm(alarm_id)

        self.assertIsNot(docalarm, None)

        alarm = docalarm[self.manager[Alerts.ALARM_STORAGE].VALUE]

        last_status_i = alarm[AlarmField.steps.value].index(alarm[AlarmField.status.value])

        state_steps = filter(
            lambda step: step['_t'] in ['stateinc', 'statedec'],
            alarm[AlarmField.steps.value][last_status_i + 1:]
        )
        self.assertEqual(len(state_steps), 10)

        # 10 remaining state changes + 6 assocticket + 1 statecounter
        all_steps = alarm[AlarmField.steps.value][last_status_i + 1:]
        self.assertEqual(len(all_steps), 17)

        expected_counter = {
            'stateinc': 1,
            'state:1': 1
        }
        counter = alarm[AlarmField.steps.value][last_status_i + 1]
        self.assertEqual(counter['val'], expected_counter)

        # Creating alarm /component/test/test0/ut-comp3
        KO['component'] = 'ut-comp3'
        OK['component'] = 'ut-comp3'
        assocticket['component'] = 'ut-comp3'

        for i in range(40):
            KO['timestamp'] = 2 * i * 300
            self.manager.archive(KO)

            OK['timestamp'] = (2 * i + 1) * 300
            self.manager.archive(OK)

            assocticket['timestamp'] = (2 * i + 1) * 300 + 10
            self.manager.archive(assocticket)

        # 80 flapping checks inserted. 10 checks to trigger flapping status.
        # 70 state changes after this change of status. Expecting 60 state to
        # be cropped.

        alarm_id = 'ut-comp3'
        docalarm = self.manager.get_current_alarm(alarm_id)

        self.assertIsNot(docalarm, None)

        alarm = docalarm[self.manager[Alerts.ALARM_STORAGE].VALUE]

        last_status_i = alarm[AlarmField.steps.value].index(alarm[AlarmField.status.value])

        state_steps = filter(
            lambda step: step['_t'] in ['stateinc', 'statedec'],
            alarm[AlarmField.steps.value][last_status_i + 1:]
        )
        self.assertEqual(len(state_steps), 10)

        # 10 remaining state changes + 36 assocticket + 1 statecounter
        all_steps = alarm[AlarmField.steps.value][last_status_i + 1:]
        self.assertEqual(len(all_steps), 47)

        expected_counter = {
            'statedec': 30,
            'stateinc': 30,
            'state:0': 30,
            'state:1': 30
        }
        counter = alarm[AlarmField.steps.value][last_status_i + 1]
        self.assertEqual(counter['val'], expected_counter)

    def test_is_hard_limit_reached(self):
        cases = [
            {
                'alarm': {AlarmField.hard_limit.value: None},
                'expected': False
            },
            {
                'alarm': {AlarmField.hard_limit.value: {'val': 99}},
                'expected': False
            },
            {
                'alarm': {AlarmField.hard_limit.value: {'val': 100}},
                'expected': True
            },
            {
                'alarm': {AlarmField.hard_limit.value: {'val': 101}},
                'expected': True
            }
        ]

        for case in cases:
            res = self.manager.is_hard_limit_reached(case['alarm'])

            self.assertIs(res, case['expected'])

    def test_check_hard_limit(self):
        from types import NoneType

        cases = [
            {
                'alarm': {
                    AlarmField.hard_limit.value: None,
                    AlarmField.steps.value: []
                },
                'expected': {
                    'type_hard_limit': NoneType,
                    'len_steps': 0
                }
            },
            {
                'alarm': {
                    AlarmField.hard_limit.value: None,
                    AlarmField.steps.value: [i for i in range(99)]
                },
                'expected': {
                    'type_hard_limit': NoneType,
                    'len_steps': 99
                }
            },
            {
                'alarm': {
                    AlarmField.hard_limit.value: None,
                    AlarmField.steps.value: [i for i in range(100)]
                },
                'expected': {
                    'type_hard_limit': dict,
                    'len_steps': 101
                }
            },
            {
                'alarm': {
                    AlarmField.hard_limit.value: {'val': 101},
                    AlarmField.steps.value: [i for i in range(200)]
                },
                'expected': {
                    'type_hard_limit': dict,
                    'len_steps': 200
                }
            },
            {
                'alarm': {
                    AlarmField.hard_limit.value: {'val': 99},
                    AlarmField.steps.value: [i for i in range(100)]
                },
                'expected': {
                    'type_hard_limit': dict,
                    'len_steps': 101
                }
            }
        ]

        for case in cases:
            alarm = self.manager.check_hard_limit(case['alarm'])

            self.assertIs(
                type(alarm[AlarmField.hard_limit.value]),
                case['expected']['type_hard_limit']
            )
            self.assertEqual(
                len(alarm[AlarmField.steps.value]),
                case['expected']['len_steps']
            )

    def test_get_events(self):
        # Empty alarm ; no events sent
        alarm0_id = '/fake/alarm/id0'

        alarm0 = self.manager.make_alarm(
            alarm0_id,
            {
                'connector': 'ut-connector',
                'connector_name': 'ut-connector0',
                'component': 'c',
                'timestamp': 0
            }
        )

        events = self.manager.get_events(alarm0)
        self.assertEqual(events, [])

        component = {
            "_id": "ut-comp",
            "impact": [],
            "name": "ut-comp",
            "measurements": [],
            "depends": [],
            "infos": {},
            "type": "component",
            "connector": "test",
            "connector_name": "test0"
        }

        self.manager.context_manager._put_entities(component)

        # Only a check OK
        alarm1_id = 'ut-comp'

        event = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.OK,
        }
        self.manager.archive(event)

        alarm1 = self.manager.get_current_alarm(alarm1_id)
        self.assertIs(alarm1, None)

        # Check KO
        event = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': Check.MINOR,
        }
        self.manager.archive(event)

        # Ack
        event = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'ack',
            'state_type': 1,
            'state': Check.MINOR,
        }
        self.manager.archive(event)

        alarm1 = self.manager.get_current_alarm(alarm1_id)
        events = self.manager.get_events(alarm1)

        expected_event0 = {
            'component': 'ut-comp',
            'connector': 'test',
            'connector_name': 'test0',
            'event_type': 'check',
            'long_output': None,
            'output': 'test message',
            'source_type': 'component',
            'state': Check.MINOR,
            'state_type': 1,
            'timestamp': 0,
        }

        expected_event1 = {
            'component': 'ut-comp',
            'connector': 'test',
            'connector_name': 'test0',
            'event_type': 'check',
            'long_output': None,
            'output': u'test message',
            'source_type': 'component',
            'state': Check.OK,
            'state_type': 1,
            'status': 1,
            'timestamp': 0,
        }

        expected_event2 = {
            'author': None,
            'component': 'ut-comp',
            'connector': 'test',
            'connector_name': 'test0',
            'event_type': 'ack',
            'ref_rk': 'test.test0.check.component.ut-comp',
            'long_output': None,
            'output': 'test message',
            'source_type': 'component',
            'state_type': 1,
            'state': Check.OK,
            'timestamp': 0,
        }

        self.assertEqual(len(events), 3)

        self.assertDictEqual(events[0], expected_event0)
        self.assertDictEqual(events[1], expected_event1)
        self.assertDictEqual(events[2], expected_event2)

        self.manager.context_manager.delete_entity(component["_id"])

    def test_check_alarm_filters(self):
        # Apply a filter on a new alarm
        now_stamp = int(mktime(datetime.now().timetuple()))
        alarm, value = self.gen_fake_alarm(moment=now_stamp)
        alarm_id = alarm[self.manager[Alerts.ALARM_STORAGE].DATA_ID]
        did = self.manager[Alerts.ALARM_STORAGE].Key.DATA_ID

        lifter = self.gen_alarm_filter({
            AlarmFilter.FILTER: {did: {"$eq": alarm_id}},
            AlarmFilter.LIMIT: -1,
            AlarmFilter.CONDITION: '{"v.connector": {"$eq": "fake-connector"}}',
            AlarmFilter.TASKS: ['alerts.systemaction.state_increase'],
            AlarmFilter.FORMAT: '>> foo',
        }, storage=self.manager[Alerts.FILTER_STORAGE])
        lifter.save()

        self.manager.update_current_alarm(alarm, value)

        self.manager.check_alarm_filters()

        result = self.manager.get_alarms(resolved=False)
        self.assertTrue(alarm_id in result)
        self.assertEqual(len(result[alarm_id]), 1)
        res_alarm = result[alarm_id][0]
        self.assertEqual(res_alarm['value']['state']['val'],
                         Check.MAJOR)
        self.assertTrue(AlarmField.filter_runs.value in res_alarm['value'])
        alarm_filters1 = res_alarm['value'][AlarmField.filter_runs.value]
        self.assertTrue(isinstance(alarm_filters1, dict))

        # Output transcription validation
        steps = result[alarm_id][0]['value'][AlarmField.steps.value]
        message = sorted(steps, key=itemgetter('t'))[-1]['m']
        self.assertEqual(message, '>> foo')

        # The filter has already been applied => alarm must not change
        now_stamp = int(mktime(datetime.now().timetuple()))
        alarm, value = self.gen_fake_alarm(moment=now_stamp)
        alarm_id2 = alarm[self.manager[Alerts.ALARM_STORAGE].DATA_ID]

        self.manager.check_alarm_filters()
        result = self.manager.get_alarms(resolved=False)
        alarm_filters2 = result[alarm_id2][0]['value'][AlarmField.filter_runs.value]
        for key in alarm_filters1.keys():
            self.assertEqual(alarm_filters1[key], alarm_filters2[key])

        # Verify that the state has correctly been increased
        self.assertEqual(result[alarm_id][0]['value']['state']['val'],
                         Check.MAJOR)

    def test_check_alarm_filters_keepstate(self):
        # Testing keepstate flag
        now_stamp = int(mktime(datetime.now().timetuple()))
        alarm, value = self.gen_fake_alarm(moment=now_stamp)
        alarm_id = alarm[self.manager[Alerts.ALARM_STORAGE].DATA_ID]
        did = self.manager[Alerts.ALARM_STORAGE].Key.DATA_ID

        lifter = self.gen_alarm_filter({
            AlarmFilter.FILTER: {did: {"$eq": alarm_id}},
            AlarmFilter.LIMIT: -1,
            AlarmFilter.CONDITION: {},
            AlarmFilter.TASKS: ['alerts.systemaction.state_increase',
                                'alerts.useraction.keepstate']
        }, storage=self.manager[Alerts.FILTER_STORAGE])
        lifter.save()

        self.manager.update_current_alarm(alarm, value)

        self.manager.check_alarm_filters()

        result = self.manager.get_alarms(resolved=False)
        state = result[alarm_id][0]['value']['state']
        self.assertEqual(state['_t'], States.changestate.value)
        self.assertEqual(state['val'], Check.MAJOR)

    def test_check_alarm_filters_repeat(self):
        # Testing repeat flag
        now_stamp = int(mktime(datetime.now().timetuple()))
        alarm, value = self.gen_fake_alarm(moment=now_stamp)
        alarm_id = alarm[self.manager[Alerts.ALARM_STORAGE].DATA_ID]
        did = self.manager[Alerts.ALARM_STORAGE].Key.DATA_ID

        lifter = self.gen_alarm_filter({
            AlarmFilter.FILTER: {did: {"$eq": alarm_id}},
            AlarmFilter.LIMIT: 1,
            AlarmFilter.CONDITION: {},
            AlarmFilter.TASKS: ['alerts.systemaction.state_increase'],
            AlarmFilter.REPEAT: 2
        }, storage=self.manager[Alerts.FILTER_STORAGE])
        lifter.save()

        self.manager.update_current_alarm(alarm, value)

        sleep(1.1)
        self.manager.check_alarm_filters()

        sleep(1.1)
        self.manager.check_alarm_filters()

        result = self.manager.get_alarms(resolved=False)
        state = result[alarm_id][0]['value']['state']
        self.assertEqual(state['_t'], 'stateinc')
        self.assertEqual(state['val'], Check.CRITICAL)

        # This one should not do anything
        sleep(1.1)
        self.manager.check_alarm_filters()

        result = self.manager.get_alarms(resolved=False)
        state = result[alarm_id][0]['value']['state']
        self.assertEqual(state['_t'], 'stateinc')
        self.assertEqual(state['val'], Check.CRITICAL)

if __name__ == '__main__':
    main()
