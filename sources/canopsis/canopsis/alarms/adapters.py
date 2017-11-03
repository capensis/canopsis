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

from __future__ import unicode_literals

from time import time
from .models import AlarmIdentity, AlarmStep, Alarm


class Adapter(object):

    COLLECTION = 'periodical_alarm'

    def __init__(self, mongo_client):
        self.mongo_client = mongo_client

    def find_unresolved_snoozed_alarms(self):
        """ Returns a list of all unresolved alarms. """

        query = {
            '$and': [
                {
                    'v.resolved': None
                },
                {
                    '$and': [  # include alarms that were never snoozed or alarms for which the snooze time has expired
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
        query = {
            '$and': [
                {
                    'v.resolved': None
                },
                {
                    '$or': [  # include alarms that were never snoozed or alarms for which the snooze time has expired
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
        selector = {
            "_id": alarm._id
        }
        col_adapter = self.mongo_client[self.COLLECTION]
        col_adapter.update(selector, alarm.to_dict())

        return alarm


def make_alarm_from_mongo(alarm_dict):
    al = alarm_dict['v']

    identity = AlarmIdentity(
        al.get('connector'),
        al.get('connector_name'),
        al.get('component'),
        al.get('resource', None)
    )
    status = make_alarm_step_from_mongo(al['status'])
    state = make_alarm_step_from_mongo(al['state'])

    steps = []
    if al.get('steps') is not None:
        for step in al['steps']:
            steps.append(make_alarm_step_from_mongo(step))

    ack = None
    if al.get('ack') is not None:
        ack = make_alarm_step_from_mongo(al['ack'])

    cancel = None
    if al.get('canceled') is not None:
        cancel = make_alarm_step_from_mongo(al['canceled'])

    snooze = None
    if al.get('snooze') is not None:
        snooze = make_alarm_step_from_mongo(al['snooze'])

    ticket = None
    if al.get('ticket') is not None:
        ticket = make_alarm_step_from_mongo(al['ticket'])

    return Alarm(
        _id=alarm_dict['_id'],
        identity=identity,
        ack=ack,
        canceled=cancel,
        creation_date=al.get('creation_date'),
        hard_limit=al.get('hard_limit'),
        initial_output=al.get('initial_output'),
        last_update_date=al.get('last_update_date'),
        resolved=al.get('resolved'),
        snooze=snooze,
        state=state,
        status=status,
        steps=steps,
        tags=al.get('tags'),
        ticket=ticket,
        alarm_filter=al.get('alarm_filter'),
        extra=al.get('extra')
    )


def make_alarm_step_from_mongo(step_dict):
    if not isinstance(step_dict, dict):
        raise TypeError("A dict is required.")
    return AlarmStep(
        author=step_dict.get('a'),
        message=step_dict.get('m'),
        type_=step_dict.get('_t'),
        timestamp=step_dict.get('t'),
        value=step_dict.get('val', None)
    )
