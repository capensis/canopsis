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

from unittest import main

from canopsis.task.core import get_task

from canopsis.alerts import AlarmField, States
from canopsis.alerts.status import get_previous_step, CANCELED, is_keeped_state

from canopsis.entitylink.manager import Entitylink

from base import BaseTest


class TestTasks(BaseTest):
    def setUp(self):
        super(TestTasks, self).setUp()

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

        self.assertTrue(alarm[AlarmField.ack.value] is not None)
        self.assertEqual(alarm[AlarmField.ack.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.ack.value]['a'], 'testauthor')
        self.assertEqual(alarm[AlarmField.ack.value]['m'], 'test message')
        self.assertTrue(alarm[AlarmField.ack.value] is get_previous_step(alarm, 'ack'))

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

        self.assertTrue(alarm[AlarmField.ack.value] is None)

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
        self.assertTrue(alarm[AlarmField.canceled.value] is not None)
        self.assertEqual(alarm[AlarmField.canceled.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.canceled.value]['a'], 'testauthor')
        self.assertEqual(alarm[AlarmField.canceled.value]['m'], 'test message')
        self.assertTrue(
            alarm[AlarmField.canceled.value] is get_previous_step(alarm, 'cancel')
        )

    def test_comment(self):
        event = {'timestamp': 0}

        task = get_task('alerts.useraction.comment')
        alarm = task(
            self.manager,
            self.alarm,
            'testauthor',
            'test message',
            event
        )

        self.assertFalse(alarm[AlarmField.comment.value] is None)
        self.assertEqual(alarm[AlarmField.comment.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.comment.value]['a'], 'testauthor')
        self.assertEqual(alarm[AlarmField.comment.value]['m'], 'test message')

    def test_restore(self):
        event = {'timestamp': 0}

        task = get_task('alerts.useraction.uncancel')
        self.alarm[AlarmField.canceled.value] = {
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

        self.assertTrue(alarm[AlarmField.canceled.value] is None)

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

        self.assertTrue(alarm[AlarmField.ticket.value] is not None)
        self.assertEqual(alarm[AlarmField.ticket.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.ticket.value]['a'], 'testauthor')
        self.assertEqual(alarm[AlarmField.ticket.value]['m'], 'test message')
        self.assertEqual(alarm[AlarmField.ticket.value]['val'], None)
        self.assertTrue(
            alarm[AlarmField.ticket.value] is get_previous_step(alarm, 'declareticket')
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

        self.assertTrue(alarm[AlarmField.ticket.value] is not None)
        self.assertEqual(alarm[AlarmField.ticket.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.ticket.value]['a'], 'testauthor')
        self.assertEqual(alarm[AlarmField.ticket.value]['m'], 'test message')
        self.assertEqual(alarm[AlarmField.ticket.value]['val'], 1234)
        self.assertTrue(
            alarm[AlarmField.ticket.value] is get_previous_step(alarm, 'assocticket')
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

        self.assertTrue(alarm[AlarmField.state.value] is not None)
        self.assertEqual(alarm[AlarmField.state.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.state.value]['a'], 'testauthor')
        self.assertEqual(alarm[AlarmField.state.value]['m'], 'test message')
        self.assertEqual(alarm[AlarmField.state.value]['val'], 2)
        self.assertTrue(
            alarm[AlarmField.state.value] is get_previous_step(alarm, States.changestate.value)
        )
        self.assertTrue(is_keeped_state(alarm))

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

        self.assertIsNot(alarm[AlarmField.snooze.value], None)
        self.assertEqual(alarm[AlarmField.snooze.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.snooze.value]['a'], 'testauthor')
        self.assertEqual(alarm[AlarmField.snooze.value]['m'], 'test message')
        self.assertEqual(alarm[AlarmField.snooze.value]['val'], 0 + 3600)
        self.assertTrue(
            alarm[AlarmField.snooze.value] is get_previous_step(alarm, 'snooze')
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

        self.assertTrue(alarm[AlarmField.state.value] is not None)
        self.assertEqual(alarm[AlarmField.state.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.state.value]['a'], 'test.test0')
        self.assertEqual(alarm[AlarmField.state.value]['m'], 'test message')
        self.assertEqual(alarm[AlarmField.state.value]['val'], state)
        self.assertTrue(
            alarm[AlarmField.state.value] is get_previous_step(alarm, 'stateinc')
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

        self.assertTrue(alarm[AlarmField.state.value] is not None)
        self.assertEqual(alarm[AlarmField.state.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.state.value]['a'], 'test.test0')
        self.assertEqual(alarm[AlarmField.state.value]['m'], 'test message')
        self.assertEqual(alarm[AlarmField.state.value]['val'], state)
        self.assertTrue(
            alarm[AlarmField.state.value] is get_previous_step(alarm, 'statedec')
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

        self.assertTrue(alarm[AlarmField.status.value] is not None)
        self.assertEqual(alarm[AlarmField.status.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.status.value]['a'], 'test.test0')
        self.assertEqual(alarm[AlarmField.status.value]['m'], 'test message')
        self.assertEqual(alarm[AlarmField.status.value]['val'], statusval)
        self.assertTrue(
            alarm[AlarmField.status.value] is get_previous_step(alarm, 'statusinc')
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

        self.assertTrue(alarm[AlarmField.status.value] is not None)
        self.assertEqual(alarm[AlarmField.status.value]['t'], 0)
        self.assertEqual(alarm[AlarmField.status.value]['a'], 'test.test0')
        self.assertEqual(alarm[AlarmField.status.value]['m'], 'test message')
        self.assertEqual(alarm[AlarmField.status.value]['val'], statusval)
        self.assertTrue(
            alarm[AlarmField.status.value] is get_previous_step(alarm, 'statusdec')
        )

    def test_linklist(self):
        self.llm = Entitylink()

        eid0 = '/entity/id'
        linklist_eid0 = {
            'computed_links': [{'label': 'doc', 'url': 'http://path/to/doc'}],
            'entity_links': [{'label': 'support', 'url': 'http://path/to/sup'}]
        }
        self.llm.put(_id=eid0, document=linklist_eid0)

        task = get_task('alerts.lookup.linklist')

        res = task(self, {'d': eid0})
        self.assertEqual(
            res,
            {
                'd': eid0,
                AlarmField.linklist.value: linklist_eid0
            }
        )

        eid1 = '/no/link/entity'
        res = task(self, {'d': eid1})
        self.assertEqual(res, {'d': eid1, AlarmField.linklist.value: {}})

        del self.llm

if __name__ == '__main__':
    main()
