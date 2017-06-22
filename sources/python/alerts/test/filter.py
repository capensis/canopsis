#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from unittest import main

from canopsis.alerts.filter import AlarmFilters, AlarmFilter

from base import BaseTest
from datetime import timedelta
import logging


class TestFilter(BaseTest):
    def test_get_filter_s(self):
        alarm, value = self.gen_fake_alarm()
        self.manager.update_current_alarm(alarm, value)

        lifter = self.gen_alarm_filter({
            AlarmFilter.FILTER: '{"v.connector":{"$eq":"wrong-connector"}}'
        }, storage=self.filter_storage)
        lifter.save()
        lifter = self.gen_alarm_filter({
            AlarmFilter.FILTER: '{"v.connector":{"$eq":"fake-connector"}}'
        }, storage=self.filter_storage)
        lifter.save()

        # get_filters()
        alarm_filters = AlarmFilters(storage=self.filter_storage,
                                     alarm_storage=self.alarm_storage,
                                     logger=logging.getLogger('alerts'))
        result = alarm_filters.get_filters()
        _id = result[0][0]._id
        self.assertTrue(isinstance(result, list))
        self.assertEqual(len(result), 1)

        # get_filter()
        result = alarm_filters.get_filter(_id)
        self.assertEqual(result[0][AlarmFilter.UID], _id)

    def test_get_filters_empty_filter(self):
        alarm, value = self.gen_fake_alarm()
        self.manager.update_current_alarm(alarm, value)
        lifter = self.gen_alarm_filter({
            AlarmFilter.FILTER: ''
        }, storage=self.filter_storage)
        lifter.save()

        alarm_filters = AlarmFilters(storage=self.filter_storage,
                                     alarm_storage=self.alarm_storage,
                                     logger=logging.getLogger('alerts'))
        result = alarm_filters.get_filters()
        self.assertTrue(isinstance(result, list))
        self.assertTrue(len(result) >= 1)

    def test_crud(self):
        alarm_filters = AlarmFilters(storage=self.filter_storage,
                                     alarm_storage=self.alarm_storage,
                                     logger=logging.getLogger('alerts'))

        result = alarm_filters.get_filters()
        self.assertEqual(result, [])

        # Init
        alarm, value = self.gen_fake_alarm()
        self.manager.update_current_alarm(alarm, value)

        element = {
            AlarmFilter.LIMIT: 30.0,
            AlarmFilter.CONDITION: {},
            AlarmFilter.TASKS: ['alerts.useraction.comment'],
            AlarmFilter.FILTER: {"d": {"$eq": alarm[self.alarm_storage.DATA_ID]}}
        }

        # CREATE
        result = alarm_filters.create_filter(element)
        result.save()
        self.assertEqual(result[AlarmFilter.CONDITION], {})
        element['_id'] = result._id

        # GET
        result = alarm_filters.get_filters()
        self.assertEqual(result[0][0].element, element)

        another_cond = '{"key": {"$eq": "another"}}'
        update = {AlarmFilter.CONDITION: another_cond}
        result = alarm_filters.update_filter('/not/an/alarm',
                                             values=update)
        self.assertTrue(result is None)

        # UPDATE
        result = alarm_filters.update_filter(element['_id'], values=update)
        self.assertEqual(result[AlarmFilter.CONDITION]['key']['$eq'], 'another')

        update = {AlarmFilter.LIMIT: 666, AlarmFilter.REPEAT: 3}
        result = alarm_filters.update_filter(element['_id'], values=update)
        self.assertEqual(result[AlarmFilter.LIMIT], timedelta(seconds=666))
        self.assertEqual(result[AlarmFilter.REPEAT], 3)

        # GET
        result = alarm_filters.get_filters()
        self.assertEqual(result[0][0][AlarmFilter.CONDITION]['key']['$eq'], 'another')

        # DELETE
        result = alarm_filters.delete_filter(element['_id'])
        self.assertEqual(result['ok'], 1.0)

        # GET
        result = alarm_filters.get_filters()
        self.assertEqual(result, [])

    def test_check_alarm(self):
        alarm, value = self.gen_fake_alarm()
        self.manager.update_current_alarm(alarm, value)

        lifter = AlarmFilter(element={
                AlarmFilter.CONDITION: {"cacao": {"$eq": 'maigre'}},
            },
            storage=self.filter_storage, alarm_storage=self.alarm_storage, logger=self.logger)
        self.assertFalse(lifter.check_alarm(alarm))

        lifter = AlarmFilter(element={
                AlarmFilter.CONDITION: {"v.component": {"$eq": 'bb'}},
            },
            storage=self.filter_storage, alarm_storage=self.alarm_storage, logger=self.logger)
        self.assertFalse(lifter.check_alarm(alarm))

        lifter = AlarmFilter(element={
                AlarmFilter.CONDITION: {"v.component": {"$eq": 'c'}},
            },
            storage=self.filter_storage, alarm_storage=self.alarm_storage, logger=self.logger)
        self.assertTrue(lifter.check_alarm(alarm))

        lifter = AlarmFilter(element={
                AlarmFilter.CONDITION: {"v.state.val": {"$gte": 1}},
            },
            storage=self.filter_storage, alarm_storage=self.alarm_storage, logger=self.logger)
        self.assertTrue(lifter.check_alarm(alarm))

    def test_output(self):
        alarm, value = self.gen_fake_alarm()
        self.manager.update_current_alarm(alarm, value)

        lifter = AlarmFilter({
            AlarmFilter.FORMAT: "",
        }, storage=self.filter_storage, alarm_storage=self.alarm_storage, logger=self.logger)

        self.assertEqual(lifter.output(''), "")

        lifter = AlarmFilter({
            AlarmFilter.FORMAT: "{old} -- foo",
        }, storage=self.filter_storage, alarm_storage=self.alarm_storage, logger=self.logger)

        self.assertEqual(lifter.output('toto'), "toto -- foo")

if __name__ == '__main__':
    main()
