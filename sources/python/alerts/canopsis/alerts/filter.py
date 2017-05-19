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

from datetime import timedelta
import json
import operator
from uuid import uuid4 as uuid


class AlarmFilters(object):
    """
        Access to a set of alarm filters.
    """

    def __init__(self, storage, alarm_storage):
        self.storage = storage  # A alarmfilter storage
        self.alarm_storage = alarm_storage  # An alarm storage

    def create_filter(self, element):
        """
        Create the selected alarm-filter.

        :param element: the filter informations
        :type element: dict
        :rtype: <AlarmFilter>
        """
        # Validating element minimal structure
        for key in [AlarmFilter.LIMIT, AlarmFilter.KEY, AlarmFilter.OPERATOR,
                    AlarmFilter.VALUE, AlarmFilter.TASKS, AlarmFilter.FILTER]:
            if key not in element:
                return None

        af = AlarmFilter(element=element, storage=self.storage)
        af.save()

        return af

    def delete_filter(self, entity_id):
        """
        Delete the selected alarm filter.

        :param entity_id: the desired Entity ID
        :type entity_id: str
        :rtype: dict
        """
        return self.storage.remove_elements(ids=[entity_id])

    def update_filter(self, alarm_id, key, value):
        """
        Retreive the list of filters linked to a specific alarm.

        :param alarm_id: the desried alarm_id
        :type alarm_id: str
        :param key: the key to update
        :type key: str
        :param value: the value to put
        :type value: str
        :rtype: <AlarmFilter>
        """
        query = {'_id': alarm_id}
        element = list(self.storage.get_elements(query=query))
        if element is None or len(element) <= 0:
            return None

        lifter = self.create_filter(element.pop())
        lifter[key] = value
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

        return all_filters[0]

    def get_filters(self):
        """
        Retreive the list of all filters.

        :rtype: [(AlarmFilter, Alarm)]
        """
        results = []

        all_filters = list(self.storage.get_elements())
        for yummy in all_filters:
            mfilter = yummy[AlarmFilter.FILTER]

            # Instanciate each AlarmFilter on this alarm
            new_filter = self.create_filter(yummy)
            try:
                query = json.loads(mfilter)
            except:
                # Cannot parse mfilter
                continue

            for alarm in list(self.alarm_storage.get_elements(query=query)):
                if new_filter is not None:
                    results.append((new_filter, alarm))

        return results

    def __repr__(self):
        return "AlarmFilters of {}".format(self.storage)


class AlarmFilter(object):
    """
        An alarm filter.

        filter = {
            '_id': 'deadbeef',
            'entity_filter': '{\"$or\":[{\"connector\":{\"$eq\":\"connector\"}}]}'
            'limit': timedelta(seconds=30),
            'condition_key': 'connector',
            'condition_operator': operator.eq,
            'condition_value': 'connector_value',
            'tasks': ['alerts.systemaction.status_increase'],
        }
    """
    UID = '_id'
    LIMIT = 'limit'
    KEY = 'condition_key'
    OPERATOR = 'condition_operator'
    VALUE = 'condition_value'
    FILTER = 'entity_filter'
    TASKS = 'tasks'
    #FORMAT = 'output_format'

    def __init__(self, element, storage=None):
        self.element = element  # has persisted in the db
        self.storage = storage

        if not element.get(self.UID, False):
            element[self.UID] = str(uuid())

        # Map and converter element parts as attribute
        for k, v in self.element.items():
            self[k] = v

    def __setitem__(self, key, item):
        value = item
        # Limit conversion
        if key == self.LIMIT and isinstance(item, (int, float)):
            value = timedelta(seconds=item)
        # Operator conversion
        elif key == self.OPERATOR and hasattr(operator, item):
            value = getattr(operator, item)

        setattr(self, key, value)
        self.element[key] = item

    def __getitem__(self, key):
        if hasattr(self, key):
            return getattr(self, key)

    def check_alarm(self, alarm):
        """
        Check if a filter is valide for a specified alarm.

        :param alarm: An alarm
        :type alarm: dict
        :rtype: bool
        """
        # Unstack the targeted value
        for mckey in self[self.KEY].split('.'):
            alarm = alarm.get(mckey, None)
            if alarm is None:
                # Cannot find the value
                return False

        # Try to evaluate the filter condition
        try:
            return self[self.OPERATOR](alarm, self[self.VALUE])
        except:
            return False

    def save(self):
        """
        Save this filter into the db.

        :raises Exception: when storage is not avalaible
        """
        if self.storage is not None:
            return self.storage.put_element(element=self.element)

        raise Exception("No storage available to save into !")

    def serialize(self):
        """
        Return a printable element (especially for json serialization)
        """
        return self.element

    def __repr__(self):
        return "AlarmFilter: {}".format(self.element)
