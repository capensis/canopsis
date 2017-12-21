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
Adapters for alarm object.
"""

from __future__ import unicode_literals

from time import time
from .models import AlarmIdentity, AlarmStep, Alarm
from canopsis.common.utils import gen_id


class AlarmAdapter(object):

    """
    Adapter for Alarm collection.
    """

    COLLECTION = 'periodical_alarm'

    def __init__(self, mongo_client):
        self.mongo_client = mongo_client

    def find_unresolved_snoozed_alarms(self):
        """
        Returns a list of all unresolved alarms.

        :rtype: [Alarm]
        """
        query = {
            '$and': [
                {
                    'v.resolved': None
                },
                {
                    '$and': [
                        # include alarms that were never snoozed or alarms
                        # for which the snooze time has expired
                        {'v.snooze.val': {'$ne': None}},
                        {'v.snooze.val': {'$lte': int(time())}}
                    ]
                },
            ]
        }

        alarms = []

        col_adapter = self.mongo_client[self.COLLECTION]
        for alarm in col_adapter.find(query):
            alarms.append(make_alarm_from_mongo(alarm))

        return alarms

    def stream_unresolved_alarms(self):
        """
        Yield unresolved alarms.

        :rtype: Alarm
        """
        query = {
            '$and': [
                {
                    'v.resolved': None
                },
                {
                    '$or': [
                        # include alarms that were never snoozed or alarms
                        # for which the snooze time has expired
                        {'v.snooze.val': None},
                        {'v.snooze.val': {'$lte': int(time())}}
                    ]
                },
            ]
        }

        col_adapter = self.mongo_client[self.COLLECTION]

        for alarm in col_adapter.find(query):
            yield make_alarm_from_mongo(alarm)

    def update(self, alarm):
        """
        Update an alarm in db.

        :param Alarm alarm: an alarm object
        :rtype: Alarm
        """
        selector = {
            "_id": alarm._id
        }
        col_adapter = self.mongo_client[self.COLLECTION]

        alarm_dict = alarm.to_dict()

        # Enforce display_name calculation
        if alarm_dict.get('display_name') in [None, '']:
            display_name = gen_id()
            while self.check_if_display_name_exists(display_name):
                display_name = gen_id()
            alarm_dict['display_name'] = display_name

        col_adapter.update(selector, alarm_dict)

        return alarm

    def check_if_display_name_exists(self, display_name):
        """
        Check if a display_name is already associated.

        :param str display_name: the name to check
        :rtype: bool
        """
        alarms = self.mongo_client[self.COLLECTION].find(
            {'v.display_name': display_name}
        )

        return alarms.count() != 0


def make_alarm_from_mongo(alarm_dict):
    """
    Build an alarm object from a mongo dict.

    :param dict alarm_dict: an alarm document
    :rtype: Alarm
    """
    ald = alarm_dict['v']

    identity = AlarmIdentity(
        ald.get('connector'),
        ald.get('connector_name'),
        ald.get('component'),
        ald.get('resource', None)
    )
    status = make_alarm_step_from_mongo(ald['status'])
    state = make_alarm_step_from_mongo(ald['state'])

    steps = []
    if ald.get('steps') is not None:
        for step in ald['steps']:
            steps.append(make_alarm_step_from_mongo(step))

    ack = None
    if ald.get('ack') is not None:
        ack = make_alarm_step_from_mongo(ald['ack'])

    cancel = None
    if ald.get('canceled') is not None:
        cancel = make_alarm_step_from_mongo(ald['canceled'])

    snooze = None
    if ald.get('snooze') is not None:
        snooze = make_alarm_step_from_mongo(ald['snooze'])

    ticket = None
    if ald.get('ticket') is not None:
        ticket = make_alarm_step_from_mongo(ald['ticket'])

    return Alarm(
        _id=alarm_dict['_id'],
        identity=identity,
        ack=ack,
        canceled=cancel,
        creation_date=ald.get('creation_date'),
        display_name=ald.get('display_name', None),
        hard_limit=ald.get('hard_limit'),
        initial_output=ald.get('initial_output'),
        last_update_date=ald.get('last_update_date'),
        resolved=ald.get('resolved'),
        snooze=snooze,
        state=state,
        status=status,
        steps=steps,
        tags=ald.get('tags'),
        ticket=ticket,
        alarm_filter=ald.get('alarm_filter'),
        extra=ald.get('extra')
    )


def make_alarm_step_from_mongo(step_dict):
    """
    Build an AlarmStep from a mongo dict.

    :param dict step_dict: an alarmstep document
    :rtype: AlarmStep
    """
    if not isinstance(step_dict, dict):
        raise TypeError("A dict is required.")

    return AlarmStep(
        author=step_dict.get('a'),
        message=step_dict.get('m'),
        type_=step_dict.get('_t'),
        timestamp=step_dict.get('t'),
        value=step_dict.get('val', None)
    )
