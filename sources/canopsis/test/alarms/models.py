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
import logging
from unittest import TestCase, main

from canopsis.alarms.models import (
    AlarmStep, AlarmIdentity, Alarm,
    ALARM_STEP_TYPE_STATE_INCREASE, ALARM_STATE_MINOR
)


class AlarmsModelsTest(TestCase):

    @classmethod
    def setUpClass(self):
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
            status=None,
            resolved=None,
            ack=None,
            tags=[],
            creation_date=1385938800,
            canceled=None,
            state=None,
            steps=[self.alarm_step],
            initial_output='Come on Rick',
            last_update_date=1506808800,
            snooze=None,
            ticket=None,
            hard_limit=None,
            extra={},
            alarm_filter=None
        )

    def test_alarm_step(self):
        self.assertEqual(self.alarm_step.author, 'Dan Harmon')

        dico = self.alarm_step.to_dict()
        self.assertEqual(dico['m'], 'Wubbalubbadubdub')

    def test_alarm_identity(self):
        self.assertEqual(self.alarm_identity.connector, 'LawnmowerDog')

        res = self.alarm_identity.get_data_id()
        self.assertEqual(res, 'MeeseeksAndDestroy/ShaymAliens')

    def test_alarm(self):
        self.assertEqual(self.alarm.identity.connector_name, 'AnatomyPark')

        res = self.alarm.to_dict()
        self.assertEqual(res['_id'], 'dim-35-c')
        self.assertEqual(res['v']['initial_output'], 'Come on Rick')
        self.assertEqual(res['v']['steps'][0]['val'], 'Shumshumschilpiddydah')

if __name__ == '__main__':
    main()
