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
from time import time

from canopsis.task.core import get_task

from canopsis.alerts.status import get_previous_step, CANCELED

from canopsis.entitylink.manager import Entitylink
from canopsis.pbehavior.manager import PBehaviorManager

from base import BaseTest


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

        self.assertFalse(alarm['steps'] is None)
        self.assertEqual(alarm['steps'][0]['t'], 0)
        self.assertEqual(alarm['steps'][0]['_t'], 'comment')
        self.assertEqual(alarm['steps'][0]['a'], 'testauthor')
        self.assertEqual(alarm['steps'][0]['m'], 'test message')

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
                'linklist': linklist_eid0
            }
        )

        eid1 = '/no/link/entity'
        res = task(self, {'d': eid1})
        self.assertEqual(res, {'d': eid1, 'linklist': {}})

        del self.llm

    def test_pbehaviors(self):
        class TestManager(object):
            pbm = PBehaviorManager()

        tm = TestManager()
        eid = '/eid'
        # Ugly, but required if last execution failed
        tm.pbm['vevent_storage']._backend.remove({'source': eid})

        task = get_task('alerts.lookup.pbehaviors')

        pb1 = {
            'behaviors': ['downtime'],
            'crecord_write_time': None,
            'enable': True,
            'id': None,
            'source': eid,
            'crecord_type': 'pbehavior',
            'rrule': None,
            'duration': None,
            'dtend': 0,
            'dtstart': 0,
            'crecord_creation_time': None,
            'crecord_name': None
        }
        tm.pbm.put([pb1])

        alarm1 = task(tm, {'d': eid})
        expected_pb1 = {
            u'enabled': True,
            u'name': [u'downtime'],
            u'dtstart': 0,
            u'dtend': 0,
            u'rrule': None
        }

        self.assertEqual(alarm1, {'d': eid, 'pbehaviors': [expected_pb1]})

        now = int(time())
        pb2 = {
            'behaviors': ['downtime'],
            'crecord_write_time': None,
            'enable': True,
            'id': None,
            'source': eid,
            'crecord_type': 'pbehavior',
            'rrule': None,
            'duration': None,
            'dtend': now + 1000,
            'dtstart': now,
            'crecord_creation_time': None,
            'crecord_name': None
        }
        tm.pbm.put([pb2])

        alarm2 = task(tm, {'d': eid})
        expected_pb2 = {
            u'enabled': True,
            u'name': [u'downtime'],
            u'dtstart': now,
            u'dtend': now + 1000,
            u'rrule': None
        }

        self.assertEqual(
            alarm2,
            {'d': eid, 'pbehaviors': [expected_pb2, expected_pb1]}
        )

        alarm3 = task(tm, {'d': '/unexisting/eid'})
        self.assertEqual(alarm3, {'d': '/unexisting/eid', 'pbehaviors': []})

        tm.pbm['vevent_storage']._backend.remove({'source': eid})


if __name__ == '__main__':
    main()
