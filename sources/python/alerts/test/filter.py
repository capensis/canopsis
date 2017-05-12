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
    def test_get_filters(self):
        alarm, value = self.gen_fake_alarm()
        alarm_id = '/fake/alarm/id'

        lifter = AlarmFilter({
            'alarms': [alarm_id, '/not/a/real/alarm/id']
        }, storage=self.filter_storage)
        lifter.save()
        lifter = AlarmFilter({
            'alarms': [alarm_id]
        }, storage=self.filter_storage)
        lifter.save()

        alarm_filters = AlarmFilters(storage=self.filter_storage).get_filters()
        self.assertTrue(isinstance(alarm_filters, dict))
        self.assertEqual(len(alarm_filters), 2)
        self.assertTrue(alarm_id in alarm_filters)

    def test_check_alarm(self):
        alarm, value = self.gen_fake_alarm()

        lifter = AlarmFilter({
            'operator': 'eq',
            'key': 'cacao',
            'value': 'maigre'
        })
        self.assertFalse(lifter.check_alarm(value))

        lifter = AlarmFilter({
            'operator': 'eq',
            'key': 'component',
            'value': 'bbb'
        })
        self.assertFalse(lifter.check_alarm(value))

        lifter = AlarmFilter({
            'operator': 'eq',
            'key': 'component',
            'value': 'c'
        })
        self.assertTrue(lifter.check_alarm(value))

        lifter = AlarmFilter({
            'operator': 'ge',
            'key': 'state.val',
            'value': 1
        })
        self.assertTrue(lifter.check_alarm(value))


if __name__ == '__main__':
    main()
