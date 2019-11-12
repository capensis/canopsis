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
To apply filters and do automated actions on alerts.
"""

from __future__ import unicode_literals

from datetime import timedelta
import json
from uuid import uuid4 as uuid
from six import string_types

from canopsis.alerts.enums import AlarmField, AlarmFilterField
from canopsis.storage.exceptions import StorageUnavailable


class AlarmFilters(object):
    """
        Easy access to a set of alarm filters.
    """

    def __init__(self, storage, alarm_storage, logger):
        """AlarmFilters Constructor !

        :param storage storage: where to find alarm filters definitions
        :param storage alarm_storage: where to find all alarms
        :param logger logger: where to log
        """
        self.storage = storage  # A alarmfilter storage
        self.alarm_storage = alarm_storage  # An alarm storage
        self.logger = logger

    def create_filter(self, element):
        """
        Create the selected alarm-filter.

        :param element: the filter informations
        :type element: dict
        :rtype: <AlarmFilter>
        """
        # Validating element minimal structure
        for key in [AlarmFilter.LIMIT, AlarmFilter.CONDITION,
                    AlarmFilter.TASKS, AlarmFilter.FILTER]:
            if key not in element:
                self.logger.error('Missing key "{}" to properlly create the filter'
                                  .format(key))
                return None

        alarmfilter = AlarmFilter(element=element,
                                  storage=self.storage,
                                  alarm_storage=self.alarm_storage,
                                  logger=self.logger)

        return alarmfilter

    def delete_filter(self, entity_id):
        """
        Delete the selected alarm filter.

        :param entity_id: the desired Entity ID
        :type entity_id: str
        :rtype: dict
        """
        return self.storage.remove_elements(ids=[entity_id])

    def update_filter(self, filter_id, values):
        """
        Update an alarm filter value.

        :param filter_id: the desired filter_id
        :type filter_id: str
        :param values: an update dict to merge with current filter
        :type values: dict
        :rtype: <AlarmFilter>
        """
        query = {'_id': filter_id}
        element = list(self.storage.get_elements(query=query))
        if element is None or len(element) <= 0:
            self.logger.debug('No alarm filter to update')
            return None

        lifter = self.create_filter(element.pop())
        for key in values.keys():
            lifter[key] = values[key]
        lifter.save()

        return lifter

    def get_filter(self, alarmfilter_id):
        """
        Retreive the list of filters linked to a specific alarm.

        :param alarmfilter_id: the desired alarmfilter_id
        :type alarmfilter_id: str
        :rtype: list or None
        """
        query = {
            AlarmFilter.UID: alarmfilter_id
        }
        all_filters = list(self.storage.get_elements(query=query))

        return [AlarmFilter(fil, logger=self.logger) for fil in all_filters]

    def get_filters(self):
        """
        Retreive the list of all filters with their alarm.

        :rtype: [(AlarmFilter, Alarm)]
        """
        results = []

        all_filters = list(self.storage.get_elements())
        for yummy in all_filters:
            mfilter = yummy[AlarmFilter.FILTER]

            # Instanciate each AlarmFilter on this alarm
            new_filter = self.create_filter(yummy)
            query = None
            if isinstance(mfilter, string_types) and mfilter != '':
                try:
                    query = json.loads(mfilter)
                except ValueError as exc:
                    self.logger.warning('Cannot parse mfilter "{}": {}'
                                        .format(mfilter, exc))
                    continue

            # Associate a filter with his matching alarm
            for alarm in self.alarm_storage._backend.find(query):
                if new_filter is not None:
                    results.append((new_filter, alarm))

        return results

    def __repr__(self):
        return "AlarmFilters of {}".format(self.storage)


class AlarmFilter(object):
    """
        An alarm filter object.

        filter = {
            "entity_filter": {"d": {"$eq": "/fake/alarm/id"}},
            "limit": timedelta(seconds=30),
            "condition": {"v.state.val": {"$eq": 1}},
            "tasks": ["alerts.systemaction.state_increase"],
            "output_format": "{old} -- message",
            "repeat": 1
        }
    """
    UID = '_id'
    LIMIT = 'limit'
    CONDITION = 'condition'
    FILTER = 'entity_filter'
    TASKS = 'tasks'
    FORMAT = 'output_format'
    REPEAT = 'repeat'

    DEFAULT_REPEAT_NUMBER = 1

    def __init__(self, element, logger, storage=None, alarm_storage=None):
        self.element = element  # has persisted in the db
        self.logger = logger
        self.storage = storage
        self.alarm_storage = alarm_storage

        if not element.get(self.UID, False):
            element[self.UID] = str(uuid())

        # Map and converter element parts as attribute
        if self.REPEAT not in self.element:
            self[self.REPEAT] = self.DEFAULT_REPEAT_NUMBER
        for k, value in self.element.items():
            self[k] = value

    def __setitem__(self, key, item):
        value = item
        # Limit conversion
        if key == self.LIMIT and isinstance(item, (int, float)):
            value = timedelta(seconds=item)
        # Condition conversion
        elif key == self.CONDITION and isinstance(item, string_types):
            try:
                value = json.loads(item)
            except Exception:
                self.logger.error('Cannot parse condition item "{}"'
                                  .format(item))
                return

        # Dict serialization conversion
        if key in [self.CONDITION, self.FILTER] and isinstance(item, dict):
            item = json.dumps(item)

        setattr(self, key, value)
        self.element[key] = item

    def __getitem__(self, key):
        if hasattr(self, key):
            return getattr(self, key)

    def get_and_check_alarm(self, alarm):
        """
        Get an up to date alarm from the database, and check that the alarm
        filter should still be applied to it.

        If the alarm filter should not be applied to the alarm, None is
        returned.

        The alarm document MUST contain _id key.

        :param alarm: An alarm
        :type alarm: dict
        :rtype: Union[dict, None]
        """
        and_ = [{'_id': alarm['_id']}]
        if self[self.CONDITION] is not None:
            and_.append(self[self.CONDITION])

        mfilter = self[self.FILTER]
        if isinstance(mfilter, string_types) and mfilter != '':
            try:
                and_.append(json.loads(mfilter))
            except ValueError as exc:
                self.logger.warning('Cannot parse mfilter "{}": {}'
                                    .format(mfilter, exc))
                return None

        query = {'$and': and_}
        alarms = list(self.alarm_storage._backend.find(query))

        if not alarms:
            return None
        return alarms[0]

    def next_run(self, alarm):
        """
        Compute the next run timestamp on the desired alarm

        :param dict alarm: the arlarm dict to test
        :returns: a timestamp
        :rtype: int or None
        """
        value = alarm[self.alarm_storage.VALUE]
        alarmfilter = value.get(AlarmField.alarmfilter.value, {})
        runs = alarmfilter.get(AlarmFilterField.runs.value, {})
        executions = runs.get(alarm[AlarmField._id.value], [])

        if self.get_and_check_alarm(alarm) and len(executions) < self.repeat:
            limit = self.element[self.LIMIT]
            if len(executions) > 0:

                return executions[-1] + limit

            return value[AlarmField.state.value]['t'] + limit

        return None

    def output(self, old=''):
        """
        Modifiy an output message according to the filter parameters.

        :param old: the old output message, if needed in the format
        :type old': str
        :rtype: str
        """
        if self[self.FORMAT] is not None:
            return self[self.FORMAT].format(old=old)

        return old

    def save(self):
        """
        Save this filter into the db.

        :raises StorageException: when storage is not avalaible
        """
        if self.storage is not None:
            return self.storage._backend.save(self.element)

        raise StorageUnavailable()

    def serialize(self):
        """
        Return a printable element (especially for json serialization).

        :rtype: dict
        """
        return self.element

    def __repr__(self):
        return "AlarmFilter: {}".format(self.element)
