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

"""
Models for Alarms.
"""

from __future__ import unicode_literals

from time import time

from canopsis.common.enumerations import FastEnum

# Alarm step types
ALARM_STEP_TYPE_STATE_INCREASE = 'stateinc'
ALARM_STEP_TYPE_STATE_DECREASE = 'statedec'
ALARM_STEP_TYPE_STATUS_DECREASE = 'statusdec'

ALARM_STEP_AUTHOR = "canopsis.engine"


class AlarmState(FastEnum):
    """Alarm states"""
    OK = 0
    MINOR = 1
    MAJOR = 2
    CRITICAL = 3


class AlarmStatus(FastEnum):
    """Alarm statuses"""
    OFF = 0
    ONGOING = 1
    STEALTHY = 2
    FLAPPING = 3
    CANCELED = 4


class AlarmStep(object):
    """
    An AlarmStep is a Step in the lifecycle of an alarm. They are used as :
    - state
    - status
    - canceled
    - snooze
    attributes
    """

    def __init__(self, author, message, type_, timestamp, value=None):
        """
        :param string author: The author of the Step
        :param string message: The message (displayed in the UI)
        :param string type_: the Step type (defined in ALARM_STEP_TYPE_* constants)
        :param float timestamp: The Step occurrence timestamp
        :param mixed value: a Value for the Step
        """
        self.author = author
        self.message = message
        self.type_ = type_
        self.timestamp = timestamp
        self.value = value

    def to_dict(self):
        """
        :rtype: dict
        """
        return {
            'a': self.author,
            '_t': self.type_,
            'm': self.message,
            't': self.timestamp,
            'val': self.value
        }


class AlarmIdentity(object):
    """
    A value Object that contains the Alarm identity.
    Stores the information about the entity impacted by the alarm.
    """

    def __init__(self, connector, connector_name, component, resource=None):
        self.connector = connector
        self.connector_name = connector_name
        self.component = component
        self.resource = resource

    def display_name(self):
        """
        Printable AlarmIdentity name.

        :rtype: str
        """
        if self.resource:
            return '{}/{}'.format(self.resource, self.component)

        return self.component

    def get_data_id(self):
        """
        :rtype: str
        """
        return self.display_name()


class Alarm(object):
    """
    An Alarm representation in Canopsis.
    """

    def __init__(
            self,
            _id,
            identity,
            ack,
            canceled,
            creation_date,
            display_name,
            hard_limit,
            initial_output,
            last_update_date,
            resolved,
            snooze,
            state,
            status,
            steps,
            tags,
            ticket,
            alarm_filter=None,
            done=None,
            extra={},
            parents=[],
            children=[]
    ):
        """
        :param str _id: db ID of the alarm
        :param AlarmIdentity identity: data to identify the alarm (resource, component...)
        :param AlarmStep ack: acknoledgment step
        :param dict alarm_filter: alarm filters informations
        :param AlarmStep canceled: canceled step
        :param int creation_date: alarm creation timestamp
        :param str display_name: displayed name of the alarm
        :param dict extra: extra informations (domain, perimeter...)
        :param bool hard_limit: hardlimit reached
        :param str initial_output: first output message
        :param int last_update_date: last update timestamp
        :param timestamp resolved: the alarm resolve timestamp
        :param AlarmStep snooze: the alarm snoozed step
        :param AlarmStep state: state step
        :param AlarmStep status: status step
        :param list steps: all alarm steps
        :param list tags: list of associated tags
        :param AlarmStep ticket: declareticket step
        :param AlarmStep done: done step
        :param list parents: list of metaalarm parents
        :param list children: list of metaalarm children
        """
        self._id = _id
        self.identity = identity
        self.ack = ack
        self.alarm_filter = alarm_filter
        self.canceled = canceled
        self.creation_date = creation_date
        self.display_name = display_name
        self.done = done
        self.entity = None
        self.extra = extra
        self.hard_limit = hard_limit
        self.initial_output = initial_output
        self.last_update_date = last_update_date
        self.resolved = resolved
        self.snooze = snooze
        self.state = state
        self.status = status
        self.steps = steps
        self.tags = tags
        self.ticket = ticket
        self.parents = parents
        self.children = children

    def to_dict(self):
        """
        Give a dict representation of the Alarm.

        :rtype: dict
        """
        value = {
            'resource': self.identity.resource,
            'tags': self.tags,
            'component': self.identity.component,
            'extra': self.extra,
            'creation_date': self.creation_date,
            'connector': self.identity.connector,
            'connector_name': self.identity.connector_name,
            'display_name': self.display_name,
            'initial_output': self.initial_output,
            'last_update_date': self.last_update_date,
            'hard_limit': self.hard_limit,
            'done': None,
            'status': None,
            'state': None,
            'ticket': None,
            'snooze': None,
            'canceled': None,
            'resolved': self.resolved,
            'ack': None,
            'steps': [],
            'parents': self.parents,
            'children': self.children
        }
        if self.status is not None:
            value['status'] = self.status.to_dict()

        if self.state is not None:
            value['state'] = self.state.to_dict()

        if self.ticket is not None:
            value['ticket'] = self.ticket.to_dict()

        if self.snooze is not None:
            value['snooze'] = self.snooze.to_dict()

        if self.canceled is not None:
            value['canceled'] = self.canceled.to_dict()

        if self.done is not None:
            value['done'] = self.done.to_dict()

        if self.ack is not None:
            value['ack'] = self.ack.to_dict()

        if self.alarm_filter is not None:
            value['alarm_filter'] = self.alarm_filter

        for ste in self.steps:
            value['steps'].append(ste.to_dict())

        return {
            '_id': self._id,
            'v': value,
            'd': self.identity.display_name(),
            't': self.creation_date
        }

    def get_last_status_value(self):
        """
        Gets the last status of an alarm.

        :returns: the last known status (check the AlarmStatus enum)
        :rtype: int
        """
        if self.status:
            return self.status.value

        return AlarmStatus.OFF

    def resolve(self, flapping_interval):
        """
        Resolves an alarm if the status is OFF and if the alarm is not
        flapping.

        :param int flapping_interval: Interval in seconds to determine status flapping
        :returns: True if the alarm is resolved, False otherwise
        :rtype: bool
        """
        now = int(time())
        if self.status is None:
            self.resolved = now
            return True
        else:
            if self.status.value is AlarmStatus.OFF \
                    and now - self.status.timestamp > flapping_interval:
                self.resolved = int(self.status.timestamp)
                return True

        return False

    def resolve_flapping(self, flapping_interval):
        """
        Resolve an alarm if it has a FLAPPING status, an OK state and a
        last state change > flapping_interval

        :param int flapping_interval: the considered flapping interval, in seconds
        :returns: True if the alarm has been resolved, False otherwise
        :rtype: bool
        """

        if self.status is None or self.status.value is not AlarmStatus.FLAPPING:
            return False

        now = int(time())
        if self.state.value == AlarmState.OK and (now - self.state.timestamp) > flapping_interval:
            self.resolved = int(self.status.timestamp)
            self.status.value = AlarmStatus.OFF
            return True

        return False

    def resolve_cancel(self, cancel_delay):
        """
        Resolve an alarm if it has been canceled, after the cancel_delay

        :param int cancel_delay: the delay where the alarm must stay in "canceled" state
        :returns: True if the alarm has been canceled, False otherwise
        :rtype: bool
        """
        if self.canceled is not None:
            canceled_date = self.canceled.timestamp
            if (int(time()) - canceled_date) >= cancel_delay:
                self.resolved = canceled_date
                return True

        return False

    def resolve_done(self, done_delay):
        """
        Resolve an alarm if it has been done, after the done_delay

        :param int done_delay: the delay where the alarm must stay in "done" state
        :returns: True if the alarm has been done, False otherwise
        :rtype: bool
        """
        if self.done is not None:
            done_date = self.done.timestamp
            if (int(time()) - done_date) >= done_delay:
                self.resolved = done_date
                return True

        return False

    def resolve_snooze(self):
        """
        Checks if the snooze has expired.

        If yes, this method removes the snooze object and returns True to
        tell the caller that the snooze has been resolved and that the
        alarm needs to be updated in DB.
        If not, this method returns false, no need to update the Alarm in DB.

        :returns: True if snooze is resolved, False otherwise.
        :rtype: bool
        """
        if self.snooze is None:
            return False

        now = int(time())
        if self.snooze.value < now:
            self.snooze = None
            self.last_update_date = now
            return True

        return False

    def resolve_stealthy(self, stealthy_interval=0):
        """
        Resolves alarms that should not be stealthy anymore.

        :param int stealthy_interval:
        :returns: True if the alarm was resolved, False otherwise
        :rtype: bool
        """
        if self.status is None \
                or self.status.value != AlarmStatus.STEALTHY \
                or self._is_stealthy(stealthy_interval):
            return False

        new_status = AlarmStep(
            author='{}.{}'.format(self.identity.connector,
                                  self.identity.connector_name),
            message='automatically resolved after stealthy shown time',
            type_=ALARM_STEP_TYPE_STATE_DECREASE,
            timestamp=int(time()),
            value=AlarmStatus.OFF
        )
        self.update_status(new_status)

        return True

    def _is_stealthy(self, stealthy_interval):
        """
        Checks if an alarm is stealthy.

        :param int stealthy_interval:
        :returns: true if the alarm is still stealthy, False otherwise
        :rtype: bool
        """
        if self.state.value != AlarmState.OK:
            return False

        for step in self.steps:
            delta = int(time()) - step.timestamp
            if delta > stealthy_interval:
                break

            if step.type_ in [ALARM_STEP_TYPE_STATE_DECREASE,
                              ALARM_STEP_TYPE_STATE_INCREASE] \
               and step.value != AlarmState.OK:
                return True

        return False

    def update_status(self, new_status):
        """
        Updates the alarm status and archives the new Step.

        :param AlarmStep new_status: the new alarm status
        """
        self.status = new_status
        self.steps.append(new_status)

    def __repr__(self):
        return '<Alarm {}>'.format(self._id)
