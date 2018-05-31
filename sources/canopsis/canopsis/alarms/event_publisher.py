#!/usr/binenv python
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

"""
Event publisher for alarms.
"""

from __future__ import unicode_literals

from canopsis.alerts.enums import AlarmField
from canopsis.common.amqp import AmqpPublisher
from canopsis.common.amqp import get_default_connection as \
    get_default_amqp_conn
from canopsis.event import Event, forger
from canopsis.logger import Logger
from canopsis.statsng.enums import StatEvents

LOG_PATH = 'var/log/alarms.log'


class AlarmEventPublisher(object):
    """
    Event publisher for alarms.
    """
    def __init__(self, amqp_pub):
        self.logger = Logger.get('alarms', LOG_PATH)
        self.amqp_pub = amqp_pub

    def publish_statcounterinc_event(self, counter_name, entity, alarm):
        """
        Publish a statcounterinc event on amqp.

        :param str counter_name: the name of the counter to increment
        :param dict entity: the entity
        :param dict alarm: the alarm
        """
        creation_date = alarm.get(AlarmField.creation_date.value)

        if not creation_date:
            self.logger.warning(
                "The alarm does not have a creation date. Ignoring it.")
            return

        component = alarm.get(Event.COMPONENT)
        resource = alarm.get(Event.RESOURCE)

        event = forger(
            connector="canopsis",
            connector_name="engine",
            event_type=StatEvents.statcounterinc,
            source_type=Event.RESOURCE if resource else Event.COMPONENT,
            component=component,
            resource=resource,
            timestamp=creation_date,
            counter_name=counter_name,
            alarm=alarm,
            entity=entity)

        self.amqp_pub.canopsis_event(event)

    def publish_statduration_event(self,
                                   timestamp,
                                   duration_name,
                                   entity,
                                   alarm):
        """
        Publish a statduration event on amqp.

        :param int timestamp: the time at which the event occurs
        :param str duration_name: the name of the duration to add
        :param dict entity: the entity
        :param dict alarm: the alarm
        """
        creation_date = alarm.get(AlarmField.creation_date.value)

        if not creation_date:
            self.logger.warning(
                "The alarm does not have a creation date. Ignoring it.")
            return

        duration = timestamp - creation_date

        component = alarm.get(Event.COMPONENT)
        resource = alarm.get(Event.RESOURCE)

        event = forger(
            connector="canopsis",
            connector_name="engine",
            event_type=StatEvents.statduration,
            source_type=Event.RESOURCE if resource else Event.COMPONENT,
            component=component,
            resource=resource,
            timestamp=timestamp,
            duration_name=duration_name,
            duration=duration,
            alarm=alarm,
            entity=entity)

        self.amqp_pub.canopsis_event(event)
