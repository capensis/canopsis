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


ALARM_STATUS_OFF = "OFF"


class AlarmStep:
    def __init__(self, author, message, type_, timestamp, value):
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
    def __init__(self, connector, connector_name, component, resource=None):
        self.connector = connector
        self.connector_name = connector_name
        self.component = component
        self.resource = resource

    def display_name(self):
        if self.resource:
            return '{0}/{1}'.format(self.component, self.resource)
        return self.component


class Alarm:
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
        extra={}

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
            'resolved': None,
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

        for s in self.steps:
            value['steps'].append(s.to_dict())

        return {
            '_id': self._id,
            'v': value,
            'd': self.identity.display_name(),
            't': self.creation_date

        }

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

    def get_last_status_value(self):
        if self.status:
            return self.status.value
        return ALARM_STATUS_OFF
