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

from canopsis.common.collection import MongoCollection


class AlarmAdapter(object):

    """
    Adapter for Alarm collection.
    """

    COLLECTION = 'periodical_alarm'

    def __init__(self, mongo_store):
        """
        :param canopsis.common.mongo_store.MongoStore mongo_store: optional MongoStore that handle HA
        """
        super(AlarmAdapter, self).__init__()

        self.mongo_store = mongo_store

        self.collection = self.mongo_store.get_collection(self.COLLECTION)

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

        for alarm in self.collection.find(query):
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

        for alarm in self.collection.find(query):
            yield make_alarm_from_mongo(alarm)

    def find_unresolved_alarms(self, entity_ids=None):
        """
        Yields all the alarms that are unresolved, optionally filtered by
        entity id.

        :param Optional[List[str]] entity_ids: A list of entity ids.
        :rtype: Iterator[Dict[str, Any]]
        """
        query = {
            "$or": [
                {"v.resolved": None},
                {"v.resolved": {"$exists": False}},
            ],
        }
        if entity_ids is not None:
            query = {
                '$and': [
                    {'d': {'$in': entity_ids}},
                    query
                ]
            }

        for alarm in self.collection.find(query):
            yield make_alarm_from_mongo(alarm)

    def get_current_alarm(self, eid, connector_name=None):
        """
        Returns exactly one alarm or None.

        This function uses the mongo_store attribute.

        :param str eid: entity ID
        :param str connector_name: add filter on connector_name if not None
        :returns: alarm document or None if no alarm found
        """

        filter_ = {
            "d": eid,
            "$or": [
                {"v.resolved": None},
                {"v.resolved": {"$exists": False}},
            ],
        }

        if connector_name is not None:
            filter_['v.connector_name'] = connector_name

        return self.collection.find_one(filter_)

    def find_last_alarms(self):
        """
        Returns the last alarm for each entity.

        This is a generator that yields the last alarm that was opened for each
        entity. The alarms may or may not have been resolved.

        :rtype: Iterator[Alarm]
        """
        pipeline = [{
            "$group": {
                "_id": "$d",
                "last_alarm": { "$last": "$$ROOT" }
            }
        }]

        for document in self.collection.aggregate(pipeline):
            yield make_alarm_from_mongo(document['last_alarm'])

    def update(self, alarm):
        """
        Update an alarm in db.

        :param Alarm alarm: an alarm object
        :rtype: Alarm
        """
        selector = {
            "_id": alarm._id
        }

        self.collection.update(selector, alarm.to_dict())

        return alarm


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

    done = None
    if ald.get('done') is not None:
        done = make_alarm_step_from_mongo(ald['done'])

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
        done=done,
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
