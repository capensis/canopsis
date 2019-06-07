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

import unittest
from canopsis.common import root_path

from canopsis.alerts.enums import AlarmField, AlarmFilterField
from canopsis.alerts.filter import AlarmFilters, AlarmFilter

from base import BaseTest
from datetime import timedelta
import logging
import xmlrunner


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
                                     alarm_storage=self.alerts_storage,
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
                                     alarm_storage=self.alerts_storage,
                                     logger=logging.getLogger('alerts'))
        result = alarm_filters.get_filters()
        self.assertTrue(isinstance(result, list))
        self.assertTrue(len(result) >= 1)

    def test_crud(self):
        alarm_filters = AlarmFilters(storage=self.filter_storage,
                                     alarm_storage=self.alerts_storage,
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
            AlarmFilter.FILTER: {
                "d": alarm[self.alerts_storage.DATA_ID]
            }
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
        self.assertEqual(
            result[AlarmFilter.CONDITION]['key']['$eq'], 'another')

        update = {AlarmFilter.LIMIT: 666, AlarmFilter.REPEAT: 3}
        result = alarm_filters.update_filter(element['_id'], values=update)
        self.assertEqual(result[AlarmFilter.LIMIT], timedelta(seconds=666))
        self.assertEqual(result[AlarmFilter.REPEAT], 3)

        # GET
        result = alarm_filters.get_filters()
        self.assertEqual(
            result[0][0][AlarmFilter.CONDITION]['key']['$eq'], 'another')

        # DELETE
        result = alarm_filters.delete_filter(element['_id'])
        self.assertEqual(result['ok'], 1.0)

        # GET
        result = alarm_filters.get_filters()
        self.assertEqual(result, [])

    def test_get_and_check_alarm(self):
        alarm, value = self.gen_fake_alarm()
        self.manager.update_current_alarm(alarm, value)

        # get back the alert's MongoDB ID for AlarmFilter
        doc = self.manager.alerts_storage._backend.find({})
        alarm['_id'] = list(doc)[0]['_id']

        lifter = AlarmFilter(element={},
                             storage=self.filter_storage,
                             alarm_storage=self.alerts_storage,
                             logger=self.logger)
        self.assertIsNotNone(lifter.get_and_check_alarm(alarm))

        lifter = AlarmFilter(
            element={
                AlarmFilter.CONDITION: {"cacao": {"$eq": 'maigre'}},
            },
            storage=self.filter_storage,
            alarm_storage=self.alerts_storage,
            logger=self.logger)
        self.assertIsNone(lifter.get_and_check_alarm(alarm))

        lifter = AlarmFilter(
            element={
                AlarmFilter.CONDITION: {"v.component": {"$eq": 'bb'}},
            },
            storage=self.filter_storage,
            alarm_storage=self.alerts_storage,
            logger=self.logger)
        self.assertIsNone(lifter.get_and_check_alarm(alarm))

        lifter = AlarmFilter(
            element={
                AlarmFilter.CONDITION: {"v.component": {"$eq": 'c'}},
            },
            storage=self.filter_storage,
            alarm_storage=self.alerts_storage,
            logger=self.logger)
        self.assertIsNotNone(lifter.get_and_check_alarm(alarm))

        lifter = AlarmFilter(
            element={
                AlarmFilter.CONDITION: {"v.state.val": {"$gte": 1}},
            },
            storage=self.filter_storage,
            alarm_storage=self.alerts_storage,
            logger=self.logger)
        self.assertIsNotNone(lifter.get_and_check_alarm(alarm))

    def test_next_run(self):
        delta = 100
        alarm, value = self.gen_fake_alarm(moment=0)
        self.manager.update_current_alarm(alarm, value)
        doc_id = list(self.manager.alerts_storage._backend.find({}))[0]['_id']
        alarm['_id'] = doc_id

        # Check no repeat
        lifter = AlarmFilter({AlarmFilter.REPEAT: 0},
                             storage=self.filter_storage,
                             alarm_storage=self.alerts_storage,
                             logger=self.logger)
        self.assertIsNone(lifter.next_run(alarm))

        # Check simple next run
        lifter = AlarmFilter({AlarmFilter.LIMIT: delta},
                             storage=self.filter_storage,
                             alarm_storage=self.alerts_storage,
                             logger=self.logger)
        self.assertTrue(lifter.next_run(alarm) >= delta)

        # Check next next run date
        value[AlarmField.alarmfilter.value] = {}
        value[AlarmField.alarmfilter.value][AlarmFilterField.runs.value] = {
            alarm['_id']: [666]
        }
        self.manager.update_current_alarm(alarm, value)
        lifter = AlarmFilter(
            {
                AlarmFilter.LIMIT: delta,
                AlarmFilter.REPEAT: 2
            },
            storage=self.filter_storage,
            alarm_storage=self.alerts_storage,
            logger=self.logger
        )
        self.assertEqual(lifter.next_run(alarm), 666 + delta)

    def test_output(self):
        alarm, value = self.gen_fake_alarm()
        self.manager.update_current_alarm(alarm, value)

        lifter = AlarmFilter({AlarmFilter.FORMAT: ""},
                             storage=self.filter_storage,
                             alarm_storage=self.alerts_storage,
                             logger=self.logger)

        self.assertEqual(lifter.output(''), "")

        lifter = AlarmFilter({AlarmFilter.FORMAT: "{old} -- foo"},
                             storage=self.filter_storage,
                             alarm_storage=self.alerts_storage,
                             logger=self.logger)

        self.assertEqual(lifter.output('toto'), "toto -- foo")

    def test_update_correct_alarm(self):
        """
        This test ensures the AlarmFilter updates the right alarm and not old
        ones that are already resolved.

        For that we create two alarms with the same data, but the first one
        will get resolved at a given timestamp.

        The alarmfilter will gather every single alarm, not matter the state
        or the resolved field, then apply the condition on them.
        """
        alarm, value = self.gen_fake_alarm(moment=42)
        self.manager.update_current_alarm(alarm, value)

        coll_alerts = self.manager.alerts_storage._backend

        alarm_doc = list(coll_alerts.find({}))[0]
        alarm['_id'] = alarm_doc['_id']
        # set arbitrary resolution time
        coll_alerts.update(
            {'_id': alarm['_id']},
            {'$set': {'v.resolved': 42, 'v.state.val': 0}}
        )

        alarm2, value2 = self.gen_fake_alarm(moment=4242)
        self.manager.update_current_alarm(alarm2, value2)

        all_alarms = list(coll_alerts.find({}))
        self.assertEqual(all_alarms[0]['v']['resolved'], 42)
        self.assertEqual(all_alarms[1]['v']['resolved'], None)
        self.assertEqual(all_alarms[0]['v']['state']['val'], 0)
        self.assertEqual(all_alarms[1]['v']['state']['val'], 1)

        filter_ = AlarmFilter(
            {
                AlarmFilter.FILTER: {},
                AlarmFilter.FORMAT: "{old} -- 2424",
                AlarmFilter.LIMIT: 100,
                AlarmFilter.REPEAT: 100,
                AlarmFilter.CONDITION: {'v.state.val': 1},
                AlarmFilter.TASKS: [
                    'alerts.systemaction.state_increase',
                    'alerts.useraction.keepstate'
                ]
            },
            storage=self.filter_storage,
            alarm_storage=self.alerts_storage,
            logger=self.logger
        )
        filter_.save()

        self.manager.check_alarm_filters()
        all_alarms = list(coll_alerts.find({}))

        self.assertEqual(all_alarms[0]['v']['resolved'], 42)
        self.assertEqual(all_alarms[1]['v']['resolved'], None)
        self.assertEqual(all_alarms[0]['v']['state']['val'], 0)
        self.assertEqual(all_alarms[1]['v']['state']['val'], 2)


if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
