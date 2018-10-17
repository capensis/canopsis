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
Event publisher for the statsng engine.
"""

from __future__ import unicode_literals
import os.path

from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.confng.helpers import cfg_to_bool
from canopsis.confng.simpleconf import ConfigurationUnreachable
from canopsis.event import Event, forger
from canopsis.statsng.enums import StatEvents

CONF_PATH = 'etc/statsng/engine.conf'
CONF_SECTION = 'ENGINE'
SEND_EVENTS_CONF_KEY = 'send_events'
LOG_PATH = 'var/log/statsng.log'


class StatEventPublisher(object):
    """
    Event publisher for the statsng engine.
    """
    def __init__(self, logger, amqp_pub):
        self.logger = logger
        self.amqp_pub = amqp_pub

        # Only send events if the configuration file exists and sets
        # send_events to True.
        try:
            cfg = Configuration.load(
                os.path.join(root_path, CONF_PATH), Ini
            ).get(CONF_SECTION, {})

            self.send_events = cfg_to_bool(cfg[SEND_EVENTS_CONF_KEY])
        except ConfigurationUnreachable:
            self.logger.warning(
                'The statsng configuration file does not exist.',
                exc_info=True)
            self.send_events = False
        except KeyError:
            self.logger.warning(
                'The send_event configuration option is not defined.')
            self.send_events = False

        if not self.send_events:
            self.logger.warning('The statistics events are disabled.')

    def publish_statcounterinc_event(self,
                                     timestamp,
                                     counter_name,
                                     entity,
                                     alarm,
                                     author=None):
        """
        Publish a statcounterinc event on amqp.

        :param int timestamp: the time at which the event occurs
        :param str counter_name: the name of the counter to increment
        :param dict entity: the entity
        :param dict alarm: the alarm
        :param Optional[str] author: the author of the event that triggered the
            statistic event.
        """
        if not self.send_events:
            return

        component = alarm.get(Event.COMPONENT)
        resource = alarm.get(Event.RESOURCE)

        # AmqpPublisher.canopsis needs the component and resource of the event
        # to be unicode strings.
        try:
            component = component.decode('utf-8')
        except (UnicodeError, AttributeError):
            pass
        try:
            resource = resource.decode('utf-8')
        except (UnicodeError, AttributeError):
            pass

        # Although only the alarm's value is needed by statsng, the alarm needs
        # to have the same structure as a full alarm for compatibility with the
        # go engines.
        full_alarm = {
            'v': alarm
        }

        event = forger(
            connector="canopsis",
            connector_name="engine",
            event_type=StatEvents.statcounterinc,
            source_type=Event.RESOURCE if resource else Event.COMPONENT,
            component=component,
            resource=resource,
            timestamp=timestamp,
            author=author,
            stat_name=counter_name,
            current_alarm=full_alarm,
            current_entity=entity)

        self.amqp_pub.canopsis_event(event)

    def publish_statduration_event(self,
                                   timestamp,
                                   duration_name,
                                   duration_value,
                                   entity,
                                   alarm,
                                   author=None):
        """
        Publish a statduration event on amqp.

        :param int timestamp: the time at which the event occurs
        :param str duration_name: the name of the duration to add
        :param str duration_value: the value of the duration
        :param dict entity: the entity
        :param dict alarm: the alarm
        :param Optional[str] author: the author of the event that triggered the
            statistic event.
        """
        if not self.send_events:
            return

        component = alarm.get(Event.COMPONENT)
        resource = alarm.get(Event.RESOURCE)

        # AmqpPublisher.canopsis needs the component and resource of the event
        # to be unicode strings.
        try:
            component = component.decode('utf-8')
        except (UnicodeError, AttributeError):
            pass
        try:
            resource = resource.decode('utf-8')
        except (UnicodeError, AttributeError):
            pass

        # Although only the alarm's value is needed by statsng, the alarm needs
        # to have the same structure as a full alarm for compatibility with the
        # go engines.
        full_alarm = {
            'v': alarm
        }

        event = forger(
            connector="canopsis",
            connector_name="engine",
            event_type=StatEvents.statduration,
            source_type=Event.RESOURCE if resource else Event.COMPONENT,
            component=component,
            resource=resource,
            timestamp=timestamp,
            author=author,
            stat_name=duration_name,
            duration=duration_value,
            current_alarm=full_alarm,
            current_entity=entity)

        self.amqp_pub.canopsis_event(event)

    def publish_statstateinterval_event(self,
                                        timestamp,
                                        state_name,
                                        state_duration,
                                        state_value,
                                        entity,
                                        alarm):
        """
        Publish a statstateinterval event on amqp.

        :param int timestamp: the time at which the event occurs
        :param str state_name: the name of the state
        :param str state_duration: the time spent in this state
        :param str state_value: the value of the state
        :param dict entity: the entity
        :param dict alarm: the alarm
        """
        if not self.send_events:
            return

        component = alarm.get(Event.COMPONENT)
        resource = alarm.get(Event.RESOURCE)

        # AmqpPublisher.canopsis needs the component and resource of the event
        # to be unicode strings.
        try:
            component = component.decode('utf-8')
        except (UnicodeError, AttributeError):
            pass
        try:
            resource = resource.decode('utf-8')
        except (UnicodeError, AttributeError):
            pass

        # Although only the alarm's value is needed by statsng, the alarm needs
        # to have the same structure as a full alarm for compatibility with the
        # go engines.
        full_alarm = {
            'v': alarm
        }

        event = forger(
            connector="canopsis",
            connector_name="engine",
            event_type=StatEvents.statstateinterval,
            source_type=Event.RESOURCE if resource else Event.COMPONENT,
            component=component,
            resource=resource,
            timestamp=timestamp,
            stat_name=state_name,
            duration=state_duration,
            state=state_value,
            current_alarm=full_alarm,
            current_entity=entity)

        self.amqp_pub.canopsis_event(event)
