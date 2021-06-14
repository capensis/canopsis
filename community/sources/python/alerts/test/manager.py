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

from canopsis.middleware.core import Middleware
from canopsis.task.core import get_task

from canopsis.timeserie.timewindow import get_offset_timewindow
from canopsis.alerts.manager import Alerts
from canopsis.alerts.status import get_previous_step, OFF, CANCELED


class BaseTest(TestCase):
    def setUp(self):
        self.alarm_storage = Middleware.get_middleware_by_uri(
            'storage-periodical-testalarm://'
        )
        self.config_storage = Middleware.get_middleware_by_uri(
            'storage-default-testconfig://'
        )

        self.manager = Alerts()
        self.manager[Alerts.ALARM_STORAGE] = self.alarm_storage
        self.manager[Alerts.CONFIG_STORAGE] = self.config_storage

        self.config_storage.put_element(
            element={
                '_id': 'test_config',
                'crecord_type': 'statusmanagement',
                'bagot_time': 3600,
                'bagot_freq': 10,
                'flapping_peristant_steps': 10,
                'hard_limit': 100,
                'stealthy_time': 300,
                'stealthy_show': 600,
                'restore_event': True,
                'auto_snooze': False,
                'snooze_default_time': 300,
            },
            _id='test_config'
        )

    def tearDown(self):
        self.alarm_storage.remove_elements()


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
            {'connector': 'ut-connector', 'timestamp': 0}
        )
        self.assertTrue(alarm is not None)

    def test_get_alarms(self):
        storage = self.manager[Alerts.ALARM_STORAGE]

        alarm0_id = '/fake/alarm/id0'
        event0 = {
            'connector': 'ut',
            'connector_name': 'ut0',
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

    def test_count_alarms_by_period(self):
        day = 24 * 3600

        alarm0_id = '/fake/alarm/id0'
        event0 = {
            'connector': 'ut',
            'connector_name': 'ut0',
            'output': '...',
            'timestamp': day / 2
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
            'output': '...',
            'timestamp': 3 * day / 2
        }
        alarm1 = self.manager.make_alarm(
            alarm1_id,
            event1
        )
        alarm1 = self.manager.update_state(alarm1, 1, event1)
        new_value1 = alarm1[self.manager[Alerts.ALARM_STORAGE].VALUE]
        self.manager.update_current_alarm(alarm1, new_value1)

        # Are subperiods well cut ?
        count = self.manager.count_alarms_by_period(0, day)
        self.assertEqual(len(count), 1)

        count = self.manager.count_alarms_by_period(0, day * 3)
        self.assertEqual(len(count), 3)

        count = self.manager.count_alarms_by_period(day, day * 10)
        self.assertEqual(len(count), 9)

        count = self.manager.count_alarms_by_period(
            0, day,
            subperiod={'hour': 1},
        )
        self.assertEqual(len(count), 24)

        # Are counts by period correct ?
        count = self.manager.count_alarms_by_period(0, day / 4)
        self.assertEqual(count[0]['count'], 0)

        count = self.manager.count_alarms_by_period(0, day)
        self.assertEqual(count[0]['count'], 1)

        count = self.manager.count_alarms_by_period(day / 2, 3 * day / 2)
        self.assertEqual(count[0]['count'], 2)

        # Does limit limits count ?
        count = self.manager.count_alarms_by_period(0, day, limit=100)
        self.assertEqual(count[0]['count'], 1)

        count = self.manager.count_alarms_by_period(day / 2, 3 * day / 2,
                                                    limit=1)
        self.assertEqual(count[0]['count'], 1)

    def test_get_current_alarm(self):
        alarm_id = '/fake/alarm/id'

        got = self.manager.get_current_alarm(alarm_id)
        self.assertTrue(got is None)

        event = {
            'connector': 'ut',
            'connector_name': 'ut0',
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
            {'connector': 'ut-connector', 'timestamp': 0}
        )

        value = alarm[storage.VALUE]

        value['state'] = {'val': 0}

        self.manager.update_current_alarm(alarm, value, tags='test')

        alarm = self.manager.get_current_alarm(alarm_id)
        value = alarm[storage.VALUE]

        self.assertTrue(value['state'] is not None)
        self.assertTrue('test' in value['tags'])

    def test_resolve_alarms(self):
        storage = self.manager[Alerts.ALARM_STORAGE]

        alarm_id = '/fake/alarm/id'
        alarm = self.manager.make_alarm(
            alarm_id,
            {'connector': 'ut-connector', 'timestamp': 0}
        )

        self.assertIsNotNone(alarm)

        value = alarm[storage.VALUE]
        value['status'] = {
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
                'resolved': {'$exists': True}
            },
            limit=1
        )
        self.assertTrue(alarm)
        alarm = alarm[0]
        value = alarm[storage.VALUE]

        self.assertEqual(value['resolved'], value['status']['t'])

    def test_change_of_state(self):
        alarm_id = '/fake/alarm/id'

        event = {
            'timestamp': 0,
            'connector': 'ut-connector',
            'connector_name': 'ut-connector0',
            'output': 'UT message',
        }

        alarm = self.manager.make_alarm(
            alarm_id,
            {'connector': 'ut-connector', 'timestamp': 0}
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
        self.assertEqual(len(alarm['value']['steps']), 2)

        self.assertEqual(alarm['value']['state'], expected_state)
        self.assertEqual(alarm['value']['steps'][0], expected_state)
        self.assertEqual(alarm['value']['status'], expected_status)
        self.assertEqual(alarm['value']['steps'][1], expected_status)

        alarm = self.manager.change_of_state(alarm, 2, 1, event)

        expected_state = {
            'a': 'ut-connector.ut-connector0',
            '_t': 'statedec',
            'm': 'UT message',
            't': 0,
            'val': 1,
        }

        # Make sure no more steps are added
        self.assertEqual(len(alarm['value']['steps']), 3)

        self.assertEqual(alarm['value']['state'], expected_state)
        self.assertEqual(alarm['value']['steps'][2], expected_state)

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
            {'connector': 'ut-connector', 'timestamp': 0}
        )

        alarm = self.manager.change_of_status(alarm, 0, 1, event)

        expected_status = {
            'a': 'ut-connector.ut-connector0',
            '_t': 'statusinc',
            'm': 'UT message',
            't': 0,
            'val': 1,
        }

        self.assertEqual(alarm['value']['status'], expected_status)

        self.assertEqual(len(alarm['value']['steps']), 1)
        self.assertEqual(alarm['value']['steps'][0], expected_status)

    def test_archive_state_nochange(self):
        alarm_id = '/component/test/test0/ut-comp'

        event0 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': 1,
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

        self.assertEqual(len(alarm['value']['steps']), 2)
        self.assertEqual(alarm['value']['steps'][0], expected_state)
        self.assertEqual(alarm['value']['state'], expected_state)

        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': 1,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        self.assertEqual(len(alarm['value']['steps']), 2)
        self.assertEqual(alarm['value']['steps'][0], expected_state)
        self.assertEqual(alarm['value']['state'], expected_state)

    def test_archive_state_changed(self):
        alarm_id = '/component/test/test0/ut-comp'

        event0 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': 1,
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

        self.assertEqual(len(alarm['value']['steps']), 2)
        self.assertEqual(alarm['value']['steps'][0], expected_state)
        self.assertEqual(alarm['value']['state'], expected_state)

        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': 2,
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

        self.assertEqual(len(alarm['value']['steps']), 3)
        self.assertEqual(alarm['value']['steps'][2], expected_state)
        self.assertEqual(alarm['value']['state'], expected_state)

    def test_archive_status_nochange(self):
        alarm_id = '/component/test/test0/ut-comp'

        event0 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': 1,
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

        self.assertEqual(len(alarm['value']['steps']), 2)
        self.assertEqual(alarm['value']['steps'][1], expected_status)
        self.assertEqual(alarm['value']['status'], expected_status)

        # Force status to stealthy
        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': 2,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        self.assertEqual(len(alarm['value']['steps']), 3)
        self.assertEqual(alarm['value']['steps'][1], expected_status)
        self.assertEqual(alarm['value']['status'], expected_status)

    def test_archive_status_changed(self):
        alarm_id = '/component/test/test0/ut-comp'

        event0 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': 1,
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

        self.assertEqual(len(alarm['value']['steps']), 2)
        self.assertEqual(alarm['value']['steps'][1], expected_status)
        self.assertEqual(alarm['value']['status'], expected_status)

        # Force status to stealthy
        event1 = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': 0,
        }
        self.manager.archive(event1)

        alarm = self.manager.get_current_alarm(alarm_id)

        expected_status = {
            'a': 'test.test0',
            '_t': 'statusinc',
            'm': 'test message',
            't': 0,
            'val': 2,
        }

        self.assertEqual(len(alarm['value']['steps']), 4)
        self.assertEqual(alarm['value']['steps'][3], expected_status)
        self.assertEqual(alarm['value']['status'], expected_status)

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
            'state': 0,
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

        alarm_id = '/component/test/test0/ut-comp1'
        docalarm = self.manager.get_current_alarm(alarm_id)

        self.assertIsNot(docalarm, None)

        alarm = docalarm[self.manager[Alerts.ALARM_STORAGE].VALUE]

        last_status_i = alarm['steps'].index(alarm['status'])

        state_steps = filter(
            lambda step: step['_t'] in ['stateinc', 'statedec'],
            alarm['steps'][last_status_i + 1:]
        )
        self.assertEqual(len(state_steps), 9)

        # 4 KO + 4 OK + 5 assocticket + 1 KO = 14 steps
        all_steps = alarm['steps'][last_status_i + 1:]
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

        alarm_id = '/component/test/test0/ut-comp2'
        docalarm = self.manager.get_current_alarm(alarm_id)

        self.assertIsNot(docalarm, None)

        alarm = docalarm[self.manager[Alerts.ALARM_STORAGE].VALUE]

        last_status_i = alarm['steps'].index(alarm['status'])

        state_steps = filter(
            lambda step: step['_t'] in ['stateinc', 'statedec'],
            alarm['steps'][last_status_i + 1:]
        )
        self.assertEqual(len(state_steps), 10)

        # 10 remaining state changes + 6 assocticket + 1 statecounter
        all_steps = alarm['steps'][last_status_i + 1:]
        self.assertEqual(len(all_steps), 17)

        expected_counter = {
            'stateinc': 1,
            'state:1': 1
        }
        counter = alarm['steps'][last_status_i + 1]
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

        alarm_id = '/component/test/test0/ut-comp3'
        docalarm = self.manager.get_current_alarm(alarm_id)

        self.assertIsNot(docalarm, None)

        alarm = docalarm[self.manager[Alerts.ALARM_STORAGE].VALUE]

        last_status_i = alarm['steps'].index(alarm['status'])

        state_steps = filter(
            lambda step: step['_t'] in ['stateinc', 'statedec'],
            alarm['steps'][last_status_i + 1:]
        )
        self.assertEqual(len(state_steps), 10)

        # 10 remaining state changes + 36 assocticket + 1 statecounter
        all_steps = alarm['steps'][last_status_i + 1:]
        self.assertEqual(len(all_steps), 47)

        expected_counter = {
            'statedec': 30,
            'stateinc': 30,
            'state:0': 30,
            'state:1': 30
        }
        counter = alarm['steps'][last_status_i + 1]
        self.assertEqual(counter['val'], expected_counter)

    def test_is_hard_limit_reached(self):
        cases = [
            {
                'alarm': {'hard_limit': None},
                'expected': False
            },
            {
                'alarm': {'hard_limit': {'val': 99}},
                'expected': False
            },
            {
                'alarm': {'hard_limit': {'val': 100}},
                'expected': True
            },
            {
                'alarm': {'hard_limit': {'val': 101}},
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
                    'hard_limit': None,
                    'steps': []
                },
                'expected': {
                    'type_hard_limit': NoneType,
                    'len_steps': 0
                }
            },
            {
                'alarm': {
                    'hard_limit': None,
                    'steps': [i for i in range(99)]
                },
                'expected': {
                    'type_hard_limit': NoneType,
                    'len_steps': 99
                }
            },
            {
                'alarm': {
                    'hard_limit': None,
                    'steps': [i for i in range(100)]
                },
                'expected': {
                    'type_hard_limit': dict,
                    'len_steps': 101
                }
            },
            {
                'alarm': {
                    'hard_limit': {'val': 101},
                    'steps': [i for i in range(200)]
                },
                'expected': {
                    'type_hard_limit': dict,
                    'len_steps': 200
                }
            },
            {
                'alarm': {
                    'hard_limit': {'val': 99},
                    'steps': [i for i in range(100)]
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
                type(alarm['hard_limit']),
                case['expected']['type_hard_limit']
            )
            self.assertEqual(
                len(alarm['steps']),
                case['expected']['len_steps']
            )

    def test_get_events(self):
        # Empty alarm ; no events sent
        alarm0_id = '/fake/alarm/id0'

        alarm0 = self.manager.make_alarm(
            alarm0_id,
            {'connector': 'ut-connector', 'timestamp': 0}
        )

        events = self.manager.get_events(alarm0)
        self.assertEqual(events, [])

        # Only a check OK
        alarm1_id = '/component/test/test0/ut-comp'

        event = {
            'source_type': 'component',
            'connector': 'test',
            'connector_name': 'test0',
            'component': 'ut-comp',
            'timestamp': 0,
            'output': 'test message',
            'event_type': 'check',
            'state': 0,
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
            'state': 1,
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
            'state': 1,
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
            'state': 1,
            'state_type': 1,
            'timestamp': 0,
            'type': 'component',
        }

        expected_event1 = {
            'component': 'ut-comp',
            'connector': 'test',
            'connector_name': 'test0',
            'event_type': 'check',
            'long_output': None,
            'output': u'test message',
            'source_type': 'component',
            'state': 0,
            'state_type': 1,
            'status': 1,
            'timestamp': 0,
            'type': 'component',
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
            'state': 0,
            'timestamp': 0,
            'type': 'component',
        }

        self.assertEqual(len(events), 3)
        self.assertEqual(events[0], expected_event0)
        self.assertEqual(events[1], expected_event1)
        self.assertEqual(events[2], expected_event2)


class TestTasks(BaseTest):
    def setUp(self):
        super(TestTasks, self).setUp()

        self.alarm = {
            'state': None,
            'status': None,
            'ack': None,
            'canceled': None,
            'ticket': None,
            'resolved': None,
            'steps': [],
            'tags': []
        }

    def test_acknowledge(self):
        event = {'timestamp': 0}

        task = get_task('alerts.useraction.ack')
        alarm = task(
            self.manager,
            self.alarm,
            'testauthor',
            'test message',
            event
        )

        self.assertTrue(alarm['ack'] is not None)
        self.assertEqual(alarm['ack']['t'], 0)
        self.assertEqual(alarm['ack']['a'], 'testauthor')
        self.assertEqual(alarm['ack']['m'], 'test message')
        self.assertTrue(alarm['ack'] is get_previous_step(alarm, 'ack'))

    def test_unacknowledge(self):
        event = {'timestamp': 0}

        task = get_task('alerts.useraction.ackremove')
        alarm = task(
            self.manager,
            self.alarm,
            'testauthor',
            'test message',
            event
        )

        self.assertTrue(alarm['ack'] is None)

        unack = get_previous_step(alarm, 'ackremove')
        self.assertEqual(unack['t'], 0)
        self.assertEqual(unack['a'], 'testauthor')
        self.assertEqual(unack['m'], 'test message')

    def test_cancel(self):
        event = {'timestamp': 0}

        task = get_task('alerts.useraction.cancel')
        alarm, statusval = task(
            self.manager,
            self.alarm,
            'testauthor',
            'test message',
            event
        )

        self.assertEqual(statusval, CANCELED)
        self.assertTrue(alarm['canceled'] is not None)
        self.assertEqual(alarm['canceled']['t'], 0)
        self.assertEqual(alarm['canceled']['a'], 'testauthor')
        self.assertEqual(alarm['canceled']['m'], 'test message')
        self.assertTrue(
            alarm['canceled'] is get_previous_step(alarm, 'cancel')
        )

    def test_restore(self):
        event = {'timestamp': 0}

        task = get_task('alerts.useraction.uncancel')
        self.alarm['canceled'] = {
            '_t': 'cancel',
            't': 0,
            'a': 'testauthor',
            'm': 'test message'
        }

        alarm, _ = task(
            self.manager,
            self.alarm,
            'testauthor',
            'test message',
            event
        )

        self.assertTrue(alarm['canceled'] is None)

        uncancel = get_previous_step(alarm, 'uncancel')
        self.assertFalse(uncancel is None)
        self.assertEqual(uncancel['t'], 0)
        self.assertEqual(uncancel['a'], 'testauthor')
        self.assertEqual(uncancel['m'], 'test message')

    def test_declare_ticket(self):
        event = {'timestamp': 0}

        task = get_task('alerts.useraction.declareticket')
        alarm = task(
            self.manager,
            self.alarm,
            'testauthor',
            'test message',
            event
        )

        self.assertTrue(alarm['ticket'] is not None)
        self.assertEqual(alarm['ticket']['t'], 0)
        self.assertEqual(alarm['ticket']['a'], 'testauthor')
        self.assertEqual(alarm['ticket']['m'], 'test message')
        self.assertEqual(alarm['ticket']['val'], None)
        self.assertTrue(
            alarm['ticket'] is get_previous_step(alarm, 'declareticket')
        )

    def test_assoc_ticket(self):
        event = {
            'timestamp': 0,
            'ticket': 1234
        }

        task = get_task('alerts.useraction.assocticket')
        alarm = task(
            self.manager,
            self.alarm,
            'testauthor',
            'test message',
            event
        )

        self.assertTrue(alarm['ticket'] is not None)
        self.assertEqual(alarm['ticket']['t'], 0)
        self.assertEqual(alarm['ticket']['a'], 'testauthor')
        self.assertEqual(alarm['ticket']['m'], 'test message')
        self.assertEqual(alarm['ticket']['val'], 1234)
        self.assertTrue(
            alarm['ticket'] is get_previous_step(alarm, 'assocticket')
        )

    def test_change_state(self):
        event = {
            'timestamp': 0,
            'state': 2
        }

        task = get_task('alerts.useraction.changestate')
        alarm = task(
            self.manager,
            self.alarm,
            'testauthor',
            'test message',
            event
        )

        self.assertTrue(alarm['state'] is not None)
        self.assertEqual(alarm['state']['t'], 0)
        self.assertEqual(alarm['state']['a'], 'testauthor')
        self.assertEqual(alarm['state']['m'], 'test message')
        self.assertEqual(alarm['state']['val'], 2)
        self.assertTrue(
            alarm['state'] is get_previous_step(alarm, 'changestate')
        )

    def test_snooze(self):
        event = {
            'connector': 'test',
            'connector_name': 'test0',
            'timestamp': 0,
            'output': 'test message',
            'duration': 3600,
        }

        task = get_task('alerts.useraction.snooze')
        alarm = task(
            self.manager,
            self.alarm,
            'testauthor',
            'test message',
            event,
        )

        self.assertIsNot(alarm['snooze'], None)
        self.assertEqual(alarm['snooze']['t'], 0)
        self.assertEqual(alarm['snooze']['a'], 'testauthor')
        self.assertEqual(alarm['snooze']['m'], 'test message')
        self.assertEqual(alarm['snooze']['val'], 0 + 3600)
        self.assertTrue(
            alarm['snooze'] is get_previous_step(alarm, 'snooze')
        )

    def test_state_increase(self):
        event = {
            'connector': 'test',
            'connector_name': 'test0',
            'timestamp': 0,
            'output': 'test message'
        }
        state = 2

        task = get_task('alerts.systemaction.state_increase')
        alarm, _ = task(self.manager, self.alarm, state, event)

        self.assertTrue(alarm['state'] is not None)
        self.assertEqual(alarm['state']['t'], 0)
        self.assertEqual(alarm['state']['a'], 'test.test0')
        self.assertEqual(alarm['state']['m'], 'test message')
        self.assertEqual(alarm['state']['val'], state)
        self.assertTrue(
            alarm['state'] is get_previous_step(alarm, 'stateinc')
        )

    def test_state_decrease(self):
        event = {
            'connector': 'test',
            'connector_name': 'test0',
            'timestamp': 0,
            'output': 'test message'
        }
        state = 0

        task = get_task('alerts.systemaction.state_decrease')
        alarm, _ = task(self.manager, self.alarm, state, event)

        self.assertTrue(alarm['state'] is not None)
        self.assertEqual(alarm['state']['t'], 0)
        self.assertEqual(alarm['state']['a'], 'test.test0')
        self.assertEqual(alarm['state']['m'], 'test message')
        self.assertEqual(alarm['state']['val'], state)
        self.assertTrue(
            alarm['state'] is get_previous_step(alarm, 'statedec')
        )

    def test_status_increase(self):
        event = {
            'connector': 'test',
            'connector_name': 'test0',
            'timestamp': 0,
            'output': 'test message'
        }
        statusval = 2

        task = get_task('alerts.systemaction.status_increase')
        alarm = task(self.manager, self.alarm, statusval, event)

        self.assertTrue(alarm['status'] is not None)
        self.assertEqual(alarm['status']['t'], 0)
        self.assertEqual(alarm['status']['a'], 'test.test0')
        self.assertEqual(alarm['status']['m'], 'test message')
        self.assertEqual(alarm['status']['val'], statusval)
        self.assertTrue(
            alarm['status'] is get_previous_step(alarm, 'statusinc')
        )

    def test_status_decrease(self):
        event = {
            'connector': 'test',
            'connector_name': 'test0',
            'timestamp': 0,
            'output': 'test message'
        }
        statusval = 0

        task = get_task('alerts.systemaction.status_decrease')
        alarm = task(self.manager, self.alarm, statusval, event)

        self.assertTrue(alarm['status'] is not None)
        self.assertEqual(alarm['status']['t'], 0)
        self.assertEqual(alarm['status']['a'], 'test.test0')
        self.assertEqual(alarm['status']['m'], 'test message')
        self.assertEqual(alarm['status']['val'], statusval)
        self.assertTrue(
            alarm['status'] is get_previous_step(alarm, 'statusdec')
        )

    def test_update_state_counter(self):
        cases = [
            {
                'alarm': {
                    'status': {
                        'a': 'ut',
                        't': 0
                    },
                    'steps': [{'a': 'ut', 't': 0}]
                },
                'diff_counter': {},
                'expected_steps': [
                    {
                        'a': 'ut',
                        't': 0
                    },
                    {
                        '_t': 'statecounter',
                        'a': 'ut',
                        't': 0,
                        'm': '',
                        'val': {}
                    }
                ]
            },
            {
                'alarm': {
                    'status': {
                        'a': 'ut',
                        't': 0
                    },
                    'steps': [{'a': 'ut', 't': 0}]
                },
                'diff_counter': {'item': 0},
                'expected_steps': [
                    {
                        'a': 'ut',
                        't': 0
                    },
                    {
                        '_t': 'statecounter',
                        'a': 'ut',
                        't': 0,
                        'm': '',
                        'val': {'item': 0}
                    }
                ]
            },
            {
                'alarm': {
                    'status': {
                        'a': 'ut',
                        't': 0
                    },
                    'steps': [{'a': 'ut', 't': 0}, {'_t': 'customstep'}]
                },
                'diff_counter': {'item1': 10, 'item2': 15},
                'expected_steps': [
                    {
                        'a': 'ut',
                        't': 0
                    },
                    {
                        '_t': 'statecounter',
                        'a': 'ut',
                        't': 0,
                        'm': '',
                        'val': {'item1': 10, 'item2': 15}
                    },
                    {
                        '_t': 'customstep'
                    }
                ]
            },
            {
                'alarm': {
                    'status': {
                        'a': 'ut',
                        't': 0
                    },
                    'steps': [
                        {
                            'a': 'ut',
                            't': 0
                        },
                        {
                            '_t': 'statecounter',
                            'a': 'ut',
                            't': 0,
                            'm': '',
                            'val': {'item1': 3}
                        },
                        {
                            '_t': 'customstep'
                        }
                    ]
                },
                'diff_counter': {'item1': 2, 'item2': 4},
                'expected_steps': [
                    {
                        'a': 'ut',
                        't': 0
                    },
                    {
                        '_t': 'statecounter',
                        'a': 'ut',
                        't': 0,
                        'm': '',
                        'val': {'item1': 5, 'item2': 4}
                    },
                    {
                        '_t': 'customstep'
                    }
                ]
            },
            {
                'alarm': {
                    'status': {
                        'a': 'ut',
                        't': 10
                    },
                    'steps': [
                        {
                            'a': 'ut',
                            't': 0
                        },
                        {
                            '_t': 'statecounter',
                            'a': 'ut',
                            't': 0,
                            'm': '',
                            'val': {'item1': 3}
                        },
                        {
                            '_t': 'customstep'
                        },
                        {
                            'a': 'ut',
                            't': 10
                        }
                    ]
                },
                'diff_counter': {'item1': 2, 'item2': 4},
                'expected_steps': [
                    {
                        'a': 'ut',
                        't': 0
                    },
                    {
                        '_t': 'statecounter',
                        'a': 'ut',
                        't': 0,
                        'm': '',
                        'val': {'item1': 3}
                    },
                    {
                        '_t': 'customstep'
                    },
                    {
                        'a': 'ut',
                        't': 10
                    },
                    {
                        '_t': 'statecounter',
                        'a': 'ut',
                        't': 10,
                        'm': '',
                        'val': {'item1': 2, 'item2': 4}
                    }
                ]
            }
        ]

        task = get_task('alerts.systemaction.update_state_counter')
        for case in cases:
            alarm = task(case['alarm'], case['diff_counter'])

            self.assertEqual(alarm['steps'], case['expected_steps'])

    def test_hard_limit(self):
        class Manager(object):
            hard_limit = 100

        mgr = Manager()
        alarm = {'hard_limit': None, 'steps': []}

        task = get_task('alerts.systemaction.hard_limit')

        alarm = task(mgr, alarm)

        self.assertIsNot(alarm['hard_limit'], None)
        self.assertEqual(len(alarm['steps']), 1)
        self.assertEqual(alarm['steps'][0], alarm['hard_limit'])

        self.assertEqual(alarm['hard_limit']['_t'], 'hardlimit')
        self.assertIs(type(alarm['hard_limit']['t']), int)
        self.assertEqual(alarm['hard_limit']['a'], '__canopsis__')
        self.assertEqual(
            alarm['hard_limit']['m'],
            (
                'This alarm has reached an hard limit (100 steps recorded) : '
                'no more steps will be appended. Please cancel this alarm or '
                'extend hard limit value.'
            )
        )
        self.assertEqual(alarm['hard_limit']['val'], 100)


if __name__ == '__main__':
    main()
