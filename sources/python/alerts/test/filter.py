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


class TestFilter(BaseTest):
    def test_get_filter_s(self):
        alarm, value = self.gen_fake_alarm()
        self.alarm_storage.put_elements([alarm])

        lifter = self.gen_alarm_filter({
            AlarmFilter.FILTER: '{"$or":[{"value.connector":{"$eq":"wrong-connector"}}]}'
        }, storage=self.filter_storage)
        lifter.save()
        lifter = self.gen_alarm_filter({
            AlarmFilter.FILTER: '{"$or":[{"value.connector":{"$eq":"fake-connector"}}]}'
        }, storage=self.filter_storage)
        lifter.save()

        # get_filters()
        alarm_filters = AlarmFilters(storage=self.filter_storage,
                                     alarm_storage=self.alarm_storage)
        result = alarm_filters.get_filters()
        _id = result[0][0]._id
        self.assertTrue(isinstance(result, list))
        self.assertEqual(len(result), 1)

        # get_filter()
        result = alarm_filters.get_filter(_id)
        self.assertEqual(result[AlarmFilter.UID], _id)

    def test_crud(self):
        alarm_filters = AlarmFilters(storage=self.filter_storage,
                                     alarm_storage=self.alarm_storage)

        result = alarm_filters.get_filters()
        self.assertEqual(result, [])

        # Init
        alarm, value = self.gen_fake_alarm()
        self.alarm_storage.put_elements([alarm])

        element = {
            AlarmFilter.LIMIT: 30.0,
            AlarmFilter.KEY: 'key',
            AlarmFilter.OPERATOR: 'neq',
            AlarmFilter.VALUE: 'value',
            AlarmFilter.TASKS: ['alerts.useraction.comment'],
            AlarmFilter.FILTER: '{"data_id":{"$eq":"/fake/alarm/id"}}'
        }

        # CREATE
        result = alarm_filters.create_filter(element)
        element['_id'] = result._id
        self.assertEqual(result[AlarmFilter.VALUE], 'value')

        # GET
        result = alarm_filters.get_filters()
        self.assertEqual(result[0][0].element, element)

        result = alarm_filters.update_filter('/not/an/alarm',
                                             key=AlarmFilter.KEY,
                                             value='another')
        self.assertTrue(result is None)

        # UPDATE
        result = alarm_filters.update_filter(element['_id'],
                                             key=AlarmFilter.KEY,
                                             value='another')
        self.assertEqual(result[AlarmFilter.KEY], 'another')

        # GET
        result = alarm_filters.get_filters()
        self.assertEqual(result[0][0][AlarmFilter.KEY], 'another')

        # DELETE
        result = alarm_filters.delete_filter(element['_id'])
        self.assertEqual(result['ok'], 1.0)

        # GET
        result = alarm_filters.get_filters()
        self.assertEqual(result, [])

    def test_check_alarm(self):
        alarm, value = self.gen_fake_alarm()

        lifter = AlarmFilter({
            AlarmFilter.OPERATOR: 'eq',
            AlarmFilter.KEY: 'cacao',
            AlarmFilter.VALUE: 'maigre'
        })
        self.assertFalse(lifter.check_alarm(alarm))

        lifter = AlarmFilter({
            AlarmFilter.OPERATOR: 'eq',
            AlarmFilter.KEY: 'value.component',
            AlarmFilter.VALUE: 'bbb'
        })
        self.assertFalse(lifter.check_alarm(alarm))

        lifter = AlarmFilter({
            AlarmFilter.OPERATOR: 'eq',
            AlarmFilter.KEY: 'value.component',
            AlarmFilter.VALUE: 'c'
        })
        self.assertTrue(lifter.check_alarm(alarm))

        lifter = AlarmFilter({
            AlarmFilter.OPERATOR: 'ge',
            AlarmFilter.KEY: 'value.state.val',
            AlarmFilter.VALUE: 1
        })
        self.assertTrue(lifter.check_alarm(alarm))

if __name__ == '__main__':
    main()
