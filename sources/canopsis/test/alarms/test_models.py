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

from copy import copy
import logging
from time import time
from unittest import TestCase, main

from canopsis.alarms.models import (
    AlarmStep, AlarmIdentity, Alarm, AlarmStatus, AlarmState,
    ALARM_STEP_TYPE_STATE_INCREASE
)
import unittest
from canopsis.common import root_path
import xmlrunner


class AlarmsModelsTest(TestCase):

    def setUp(self):
        self.logger = logging.getLogger('alarms')

        self.alarm_step = AlarmStep(
            author='Dan Harmon',
            message='Wubbalubbadubdub',
            type_=ALARM_STEP_TYPE_STATE_INCREASE,
            timestamp=1385938800,
            value='Shumshumschilpiddydah'
        )

        self.alarm_identity = AlarmIdentity(
            connector='LawnmowerDog',
            connector_name='AnatomyPark',
            component='ShaymAliens',
            resource='MeeseeksAndDestroy'
        )

        self.alarm = Alarm(
            _id='dim-35-c',
            identity=self.alarm_identity,
            ack=None,
            canceled=None,
            creation_date=1385938800,
            hard_limit=None,
            initial_output='Come on Rick',
            last_update_date=1506808800,
            resolved=None,
            snooze=None,
            state=None,
            status=None,
            steps=[self.alarm_step],
            tags=[],
            ticket=None,
            alarm_filter=None,
            extra={}
        )

    def test_alarmstep(self):
        self.assertEqual(self.alarm_step.author, 'Dan Harmon')

        dico = self.alarm_step.to_dict()
        self.assertEqual(dico['m'], 'Wubbalubbadubdub')

    def test_alarmidentity(self):
        self.assertEqual(self.alarm_identity.connector, 'LawnmowerDog')

        res = self.alarm_identity.get_data_id()
        self.assertEqual(res, 'MeeseeksAndDestroy/ShaymAliens')

    def test_alarm(self):
        self.assertEqual(self.alarm.identity.connector_name, 'AnatomyPark')

        res = self.alarm.to_dict()
        self.assertEqual(res['_id'], 'dim-35-c')
        self.assertEqual(res['v']['initial_output'], 'Come on Rick')
        self.assertEqual(res['v']['steps'][0]['val'], 'Shumshumschilpiddydah')

    def test_alarm_get_last_status_value(self):
        self.assertEqual(self.alarm.get_last_status_value(),
                         AlarmStatus.OFF)

        self.alarm.status = AlarmStep(
            author='Morty',
            message='Smith',
            type_=ALARM_STEP_TYPE_STATE_INCREASE,
            timestamp=1506808800
        )
        self.assertEqual(self.alarm.get_last_status_value(),
                         self.alarm.status.value)

    def test_alarm_resolve(self):
        self.assertTrue(self.alarm.resolve(0))

        self.alarm.status = AlarmStep(
            author='Jerry',
            message='Smith',
            type_=ALARM_STEP_TYPE_STATE_INCREASE,
            timestamp=1506808800,
            value=AlarmStatus.ONGOING
        )
        self.assertFalse(self.alarm.resolve(0))

        self.alarm_step.value = AlarmStatus.OFF
        self.alarm.status = self.alarm_step
        self.assertTrue(self.alarm.resolve(0))
        self.assertNotEqual(self.alarm.resolved, self.alarm_step.value)

    def test_alarm_resolve_cancel(self):
        self.assertFalse(self.alarm.resolve_cancel(0))

        ts = 1506808800
        self.alarm.canceled = AlarmStep(
            author='Beth',
            message='Smith',
            type_=ALARM_STEP_TYPE_STATE_INCREASE,
            timestamp=ts
        )
        self.assertTrue(self.alarm.resolve_cancel(0))
        self.assertEqual(self.alarm.resolved, ts)

    def test_alarm_resolve_snooze(self):
        self.assertFalse(self.alarm.resolve_snooze())

        self.alarm.snooze = AlarmStep(
            author='Summer',
            message='Smith',
            type_=ALARM_STEP_TYPE_STATE_INCREASE,
            timestamp=1506808800,
            value=AlarmStatus.ONGOING
        )
        last = self.alarm.last_update_date
        self.assertTrue(self.alarm.resolve_snooze())
        self.assertTrue(self.alarm.snooze is None)
        self.assertNotEqual(self.alarm.last_update_date, last)

        self.alarm.snooze = AlarmStep(
            author='Summer',
            message='Smith',
            type_=ALARM_STEP_TYPE_STATE_INCREASE,
            timestamp=int(time()) + 10000,
            value=AlarmStatus.ONGOING
        )
        self.alarm.snooze.value = int(time()) + 100
        self.assertFalse(self.alarm.resolve_snooze())

    def test_alarm_is_stealthy(self):
        self.alarm.state = AlarmStep(
            author='Coach',
            message='Feratu',
            type_=ALARM_STEP_TYPE_STATE_INCREASE,
            timestamp=int(time()) - 100,
            value=AlarmState.OK
        )
        step = copy(self.alarm.state)
        step.value = AlarmState.MAJOR
        self.alarm.steps = [step]

        self.assertFalse(self.alarm._is_stealthy(0, 0))

        self.assertTrue(self.alarm._is_stealthy(9999, 9999))

    def test_alarm_resolve_stealthy(self):
        self.assertFalse(self.alarm.resolve_stealthy())

        self.alarm.status = AlarmStep(
            author='Bird',
            message='person',
            type_=ALARM_STEP_TYPE_STATE_INCREASE,
            timestamp=int(time()) - 100,
            value=AlarmStatus.STEALTHY
        )
        self.alarm.state = AlarmStep(
            author='Xenon',
            message='Bloom',
            type_=ALARM_STEP_TYPE_STATE_INCREASE,
            timestamp=int(time()) - 100,
            value=AlarmState.OK
        )
        last = len(self.alarm.steps)
        self.assertTrue(self.alarm.resolve_stealthy(9999, 9999))
        self.assertEqual(self.alarm.status.author, 'LawnmowerDog.AnatomyPark')
        self.assertEqual(len(self.alarm.steps), last + 1)

if __name__ == '__main__':
    output = root_path + "tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
