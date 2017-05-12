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
        Access to a set of alarm filters.
    """

    def __init__(self, storage):
        self.storage = storage

    def get_filters(self):
        """
        List and generate all <AlarmFilter> grouped by alarm id.

        :rtype: dict
        """
        result = {}

        # Get all alarms with at least one filter
        alarms = self.storage._backend.distinct('alarms')
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
                # Instanciate each AlarmFilter on this alarm
                result[alarm_id].append(AlarmFilter(element=element,
                                                    storage=self.storage))

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

    def __init__(self, element, storage=None):
        self.element = element  # has persisted in the db
        self.storage = storage

        # Map and converter element parts as attribute
        for k, v in self.element.items():
            self[k] = v

    def __setitem__(self, key, item):
        value = item
        if key == 'limit' and isinstance(item, int):
            # Limit conversion
            value = timedelta(minutes=item)
        elif key == 'operator' and hasattr(operator, item):
            # Operator conversion
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
        # Find the target value
        val = alarm_value
        for mckey in self.key.split('.'):
            val = val.get(mckey)

        # Try to evaluate the filter
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

    def __repr__(self):
        if hasattr(self, 'limit') and hasattr(self, 'key') \
           and hasattr(self, 'operator') and hasattr(self, 'value') \
           and hasattr(self, 'tasks'):

            return ("AlarmFilter: {(after {} ; {} {} {} ; {})}"
                    .format(self.limit, self.key, self.operator,
                            self.value, self.tasks))

        return "AlarmFilter: {}".format(self.element)
