#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

from unittest import TestCase, main
from mock import Mock

from canopsis.alerts.enums import AlarmField
from canopsis.event import Event
from canopsis.statsng.enums import StatEvents, StatEventFields
from canopsis.statsng.event_publisher import StatEventPublisher


class StatEventPublisherTest(TestCase):

    def setUp(self):
        logger = Mock()
        self.amqp_pub = Mock()
        self.event_publisher = StatEventPublisher(logger, self.amqp_pub)

    def test_no_events(self):
        self.event_publisher.send_events = False

        timestamp = 3
        self.event_publisher.publish_statcounterinc_event(
            1, 'counter_name', {}, {})
        self.event_publisher.publish_statduration_event(
            timestamp, 'duration_name', 2, {}, {})

        self.assertEqual(self.amqp_pub.canopsis_event.call_count, 0)

    def test_publish_statcounterinc_event(self):
        self.event_publisher.send_events = True

        self.event_publisher.publish_statcounterinc_event(
            1, 'counter_name', {}, {})

        self.assertEqual(self.amqp_pub.canopsis_event.call_count, 1)

        event = self.amqp_pub.canopsis_event.call_args[0][0]
        self.assertEqual(event[Event.EVENT_TYPE], StatEvents.statcounterinc)
        self.assertEqual(event[StatEventFields.stat_name], 'counter_name')
        self.assertEqual(event['timestamp'], 1)

    def test_publish_statduration_event(self):
        self.event_publisher.send_events = True

        timestamp = 3
        self.event_publisher.publish_statduration_event(
            timestamp, 'duration_name', 2, {}, {})

        self.assertEqual(self.amqp_pub.canopsis_event.call_count, 1)

        event = self.amqp_pub.canopsis_event.call_args[0][0]
        self.assertEqual(event[Event.EVENT_TYPE], StatEvents.statduration)
        self.assertEqual(event[StatEventFields.stat_name], 'duration_name')
        self.assertEqual(event['timestamp'], 3)
        self.assertEqual(event[StatEventFields.duration], 2)


if __name__ == '__main__':
    main()
