#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

import time

# Alarm statuses
ALARM_STATUS_OFF = 0
ALARM_STATUS_STEALTHY = 2

# Alarm states
ALARM_STATE_OK = 0
ALARM_STATE_MINOR = 1
ALARM_STATE_MAJOR = 2
ALARM_STATE_CRITICAL = 3

# Alarm step types
ALARM_STEP_TYPE_STATE_INCREASE = 'stateinc'
ALARM_STEP_TYPE_STATE_DECREASE = 'statedec'

class AlarmStep:
    """
    An AlarmStep is a Step in the lifecycle of an alarm. They are used as :
    - state
    - status
    - canceled
    - snooze
    attributes
    """
    def __init__(self, author, message, type_, timestamp, value):
        """
        constructor
        :param string author: The author of the Step
        :param string message: The message (displayed in the UI)
        :param string  type_: the Step type (defined in ALARM_STEP_TYPE_* constants)
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

        :return: dict
        """
        return {
            'a': self.author,
            '_t': self.type_,
            'm': self.message,
            't': self.timestamp,
            'val': self.value
        }


class AlarmIdentity:
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
        if self.resource:
            return '{0}/{1}'.format(self.resource, self.component)
        return self.component

    def get_data_id(self):
        return self.display_name()


class Alarm:
    """
    An Alarm representation in Canopsis.
    """
    def __init__(
        self,
        _id,
        identity,
        status,
        resolved,
        ack,
        tags,
        creation_date,
        canceled,
        state,
        steps,
        initial_output,
        last_update_date,
        snooze,
        ticket,
        hard_limit,
        extra={},
        alarm_filter=None

    ):
        self._id = _id
        self.identity = identity
        self.status = status
        self.resolved = resolved
        self.ack = ack
        self.tags = tags
        self.creation_date = creation_date
        self.canceled = canceled
        self.state = state
        self.steps = steps
        self.initial_output = initial_output
        self.last_update_date = last_update_date
        self.snooze = snooze
        self.ticket = ticket
        self.hard_limit = hard_limit
        self.entity = None
        self.extra = extra
        self.alarm_filter = alarm_filter

    def to_dict(self):
        value = {
            'resource': self.identity.resource,
            'tags': self.tags,
            'component': self.identity.component,
            'extra': self.extra,
            'creation_date': self.creation_date,
            'connector': self.identity.connector,
            'connector_name': self.identity.connector_name,
            'initial_output': self.initial_output,
            'last_update_date': self.last_update_date,
            'hard_limit': self.hard_limit,
            'status': None,
            'state': None,
            'ticket': None,
            'snooze': None,
            'canceled': None,
            'resolved': self.resolved,
            'ack': None,
            'steps': []
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

        if self.ack is not None:
            value['ack'] = self.ack.to_dict()

        if self.alarm_filter is not None:
            value['alarm_filter'] = self.alarm_filter
        for s in self.steps:
            value['steps'].append(s.to_dict())

        return {
            '_id': self._id,
            'v': value,
            'd': self.identity.display_name(),
            't': self.creation_date

        }

    def resolve(self, flapping_interval):
        """
        Resolves an alarm if the status is OFF and if the alarm is not flapping.
        :param int flapping_interval: Interval in seconds to determine status flapping
        :return Bool: True if the alarm is resolved, False otherwise
        """
        if self.status is None:
            self.resolved = time.time()
            return True
        else:
            if self.status.value is not ALARM_STATUS_OFF:
                return False
            else:
                if time.time() - self.status.timestamp > flapping_interval:
                    self.resolved = int(self.status.timestamp)
                    return True
                return False



    def resolve_snooze(self):
        """
            Checks if the snooze has expired.

            if yes, this method removes the snooze object and returns True to tell the caller
            that the snooze has been resolved and that the alarm needs to be updated in DB.
            if not, this method returns false, no need to update the Alarm in DB.

            :return Bool: True if snooze is resolved, False otherwise.
        """
        if self.snooze is None:
            return False

        if self.snooze.value < time.time():
            self.snooze = None
            self.last_update_date = time.time()
            return True

    def resolve_cancel(self, cancel_delay):
        """
            Resolve an alarm if it has been canceled, after the cancel_delay
        :param cancel_delay: the delay where the alarm must stay in "canceled" state
        :return Bool: True if the alarm has been canceled, False otherwise
        """
        if self.canceled is not None:
            canceled_date = self.canceled.timestamp
            if (time.time() - canceled_date) >= cancel_delay:
                self.resolved = canceled_date
                return True
        return False

    def get_last_status_value(self):
        """
        Gets the last status of an alarm
        :return int: the last known status (check ALARM_STATUS_* constants)
        """
        if self.status:
            return self.status.value
        return ALARM_STATUS_OFF

    def _is_stealthy(self, stealthy_show_duration, stealthy_interval):
        """
        Checks if an alarm is stealthy.
        :param int stealthy_show_duration:
        :param int stealthy_interval:
        :return Bool: true if the alarm is still stealthy, False otherwise
        """
        last_state_ts = self.state.timestamp
        for step in self.steps:
            delta1 = last_state_ts - step.timestamp
            delta2 = int(time.time()) - step.timestamp
            if delta1 > stealthy_show_duration or \
                    delta1 > stealthy_interval or \
                    delta2 > stealthy_show_duration or \
                    delta2 > stealthy_interval:
                break

            if step.type in [ALARM_STEP_TYPE_STATE_DECREASE, ALARM_STEP_TYPE_STATE_INCREASE]:
                if step.value != ALARM_STATE_OK and self.state.value == ALARM_STATE_OK:
                    return True
        return False

    def resolve_stealthy(self, stealthy_show_duration=0, stealthy_interval=0):
        """
        Resolves alarms that should not be stealthy anymore.
        :param int stealthy_show_duration: duration (in seconds) where the alarm should be shown as stealthy
        :param int stealthy_interval:
        :return Bool: True if the alarm was resolved, False otherwise
        """
        if self.status is None or self.status.value != ALARM_STATUS_STEALTHY:
            return False
        if self._is_stealthy(stealthy_show_duration, stealthy_interval):
            return False

        new_status = AlarmStep(
            '{0}.{1}'.format(self.identity.connector, self.identity.connector_name),
            'automatically resolved after stealthy shown time',
            ALARM_STEP_TYPE_STATE_DECREASE,
            time.time(),
            ALARM_STATUS_OFF)

        self.update_status(new_status)
        return True

    def update_status(self, new_status):
        """
        Updates the alarm status and archives the new Step.
        :param AlarmStep new_status: the new alarm status
        :return : None
        """
        self.status = new_status
        self.steps.append(new_status)




