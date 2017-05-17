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
import operator
from uuid import uuid4 as uuid


class AlarmFilters(object):
    """
        Access to a set of alarm filters.
    """

    def __init__(self, storage):
        self.storage = storage
        self.filters = None  # All alarm filters, as dict, grouped by alarm_id

    def create_filter(self, element):
        """
        Create the selected alarm-filter.

        :param element: the filter informations
        :type element: dict
        :rtype: <AlarmFilter>
        """
        # Validating element minimal structure
        for key in ['limit', 'key', 'operator', 'value', 'tasks', 'alarms']:
            if key not in element:
                return None

        self.filters = None  # Invalidate filter list cache
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
        self.filters = None  # Invalidate filter list cache
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
        self.filters = None  # Invalidate filter list cache
        query = {'_id': alarm_id}
        element = list(self.storage.find_elements(query=query))
        if element is None or len(element) <= 0:
            return None

        lifter = self.create_filter(element.pop())
        lifter[key] = value
        lifter.save()

        return lifter

    def get_filter(self, alarm_id):
        """
        Retreive the list of filters linked to a specific alarm.

        :param alarm_id: the desired alarm_id
        :type alarm_id: str
        :rtype: list or None
        """
        if self.filters is None:
            self.filters = self._get_filters()

        return self.filters.get(alarm_id, None)

    def get_filters(self):
        """
        Retreive the list of all filters.

        :rtype: list
        """
        if self.filters is None:
            self.filters = self._get_filters()

        return self.filters

    def _get_filters(self):
        """
        List and generate all <AlarmFilter> grouped by alarm id.

        :rtype: dict
        """
        result = {}
        known_filters = {}

        # Get all alarms with at least one filter
        alarms = list(self.storage._backend.distinct('alarms'))
        for alarm_id in alarms:
            query = {
                'alarms': {
                    '$in': [alarm_id]
                }
            }

            # Get filters associate with this alarm
            for element in self.storage.find_elements(query=query):
                if alarm_id not in result:
                    result[alarm_id] = []

                # Deduplicate AlarmFilter objects
                element_id = element.get('_id', None)
                if element_id in known_filters:
                    result[alarm_id].append(known_filters[element_id])
                    continue

                # Instanciate each AlarmFilter on this alarm
                new_filter = self.create_filter(element)
                result[alarm_id].append(new_filter)
                known_filters[element_id] = new_filter

        return result

    def __repr__(self):
        return "AlarmFilters of {}".format(self.storage)


class AlarmFilter(object):
    """
        An alarm filter.

        filter = {
            '_id': 'deadbeef',
            'alarms': ['/id/of/linked/alarm']
            'limit': timedelta(minutes=30),
            'key': 'connector',
            'operator': operator.eq,
            'value': 'connector_value',
            'tasks': ['alerts.systemaction.status_increase'],
        }
    """
    UID = '_id'
    LIMIT = 'limit'
    OPERATOR = 'operator'

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
            value = timedelta(minutes=item)
        # Operator conversion
        elif key == self.OPERATOR and hasattr(operator, item):
            value = getattr(operator, item)

        setattr(self, key, value)
        self.element[key] = item

    def check_alarm(self, alarm_value):
        """
        Check if a filter is valide for a specified alarm value.

        :param alarm_value: An alarm value
        :type alarm_value: dict
        :rtype: bool
        """
        # Find the targeted value
        val = alarm_value
        for mckey in self.key.split('.'):
            val = val.get(mckey)

        # Try to evaluate the filter condition
        try:
            return self.operator(val, self.value)
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
