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


class AlarmFilters(object):
    """
        A set of alarm filters, linked to a specific alarm.
    """

    def __init__(self, storage, alarm_id):
        self.storage = storage
        self.alarm_id = alarm_id

        query = {
            'alarms': {
                '$in': [self.alarm_id]
            }
        }
        self.elements = self.storage.find_elements(query=query)

        self.filters = []
        for element in self.elements:
            self.filters.append(AlarmFilter(element=element))

    def __iter__(self):
        return iter(self.filters)

    def __repr__(self):
        return "AlarmFilters of {}: {}".format(self.alarm_id, self.elements)


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

    def __init__(self, element, storage=None):
        self.element = element  # has persisted in the db
        self.storage = storage

        for k, v in self.element.items():
            self[k] = v

    def __setitem__(self, key, item):
        value = item
        if key == 'operator' and hasattr(operator, item):
            # Operator conversion
            value = getattr(operator, item)
        elif key == 'limit' and isinstance(item, int):
            # Limit conversion
            value = timedelta(minutes=item)

        setattr(self, key, value)
        self.element[key] = item

    def check_alarm(self, alarm_value):
        """
        Check if a filter is valide for a specified alarm value.

        :param alarm_value: An alarm value
        :type alarm_value: dict
        :rtype: bool
        """
        if not isinstance(alarm_value, dict):
            return False

        if self.key in alarm_value:
            return self.operator(alarm_value[self.key], self.value)

        return False

    def save(self):
        """
        Save this filter into the db.
        """
        if self.storage is not None:
            return self.storage.put_element(element=self.element)

        # TODO: use the local logger instead
        print("No storage available to save into !")

    def __repr__(self):
        #return "AlarmFilter: {}".format(dir(self))
        return ("AlarmFilter: (after {} ; {} {} {} ; {})"
                .format(self.limit, self.key, self.operator,
                        self.value, self.tasks))
