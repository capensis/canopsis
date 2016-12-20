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

from sys import prefix
from os.path import join
from json import load

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.model import Parameter

from canopsis.timeserie.timewindow import Interval, TimeWindow


CONF_PATH = 'alerts/manager.conf'
CATEGORY = 'ALERTS'
CONTENT = [
    Parameter('extra_fields', Parameter.array())
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class AlertsReader(MiddlewareRegistry):
    """
    Alarm cycle managment.

    Used to retrieve events related to alarms in a TimedStorage.
    """

    CONFIG_STORAGE = 'config_storage'
    ALARM_STORAGE = 'alarm_storage'
    ALARM_FIELDS_STORAGE = 'alarm_fields_storage'
    CONTEXT_MANAGER = 'context'

    def __init__(
            self,
            extra_fields=None,
            alarm_storage=None,
            alarm_fields_storage=None,
            context=None,
            *args, **kwargs
    ):
        super(AlertsReader, self).__init__(*args, **kwargs)

        if extra_fields is not None:
            self.extra_fields = extra_fields

        if alarm_storage is not None:
            self[AlertsReader.ALARM_STORAGE] = alarm_storage

        if alarm_fields_storage is not None:
            self[AlertsReader.ALARM_FIELDS_STORAGE] = alarm_fields_storage

        if context is not None:
            self[AlertsReader.CONTEXT_MANAGER] = context

    @property
    def alarm_fields(self):
        if not hasattr(self, '_alarms'):
            self.load_config()

        return self._alarms

    def load_config(self):
        with open(join(prefix, 'etc/schema.d/alarm_fields.json')) as fh:
            self._alarms = load(fh)

    def translate_key(self, key):
        if key in self.alarm_fields['properties']:
            return self.alarm_fields['properties'][key]['stored_name']

        return key

    def translate_filter(self, filter_):
        if type(filter_) is list:
            for i, f in enumerate(filter_):
                filter_[i] = self.translate_filter(f)

        elif type(filter_) is dict:
            for key, value in filter_.items():
                new_value = self.translate_filter(value)
                filter_[key] = new_value

                new_key = self.translate_key(key)
                filter_[new_key] = filter_.pop(key)

        return filter_

    def translate_sort(self, sort):
        for key, value in sort.items():
            new_key = self.translate_key(key)
            sort[new_key] = sort.pop(key)

        return sort

    def interpret_search(self, search):
        return {}, 'current'

    def get_time_filter(self, opened, resolved, tstart, tstop):
        if opened and resolved:
            return {
                '$or': [
                    self.get_opened_time_filter(tstop),
                    self.get_resolved_time_filter(tstart, tstop)
                ]
            }

        if opened and not resolved:
            return self.get_opened_time_filter(tstop)

        if not opened and resolved:
            return self.get_resolved_time_filter(tstart, tstop)

        return None

    def get_opened_time_filter(self, tstop):
        return {
            'v.resolved': None,
            't': {'$gte': tstop}
        }

    def get_resolved_time_filter(self, tstart, tstop):
        return {
            'v.resolved': {'$neq': None},
            '$or': [
                {'t': {'$gte': tstart, '$lte': tstop}},
                {'v.resolved': {'$gte': tstart, '$lte': tstop}},
                {'t': {'$lte': tstart}, 'v.resolved': {'$gte': tstop}}
            ]
        }

    def get(
            self,
            tstart,
            tstop,
            opened=True,
            resolved=False,
            consolidations=[],
            filter_={},
            search='',
            sort={'opened': 'DESC'},
            limit=50,
            offset=0
    ):
        """
        Return filtered, sorted and paginated alarms.

        :param int tstart: Beginning timestamp of requested period
        :param int tstop: End timestamp of requested period

        :param bool opened: If True, consider alarms that are currently opened
        :param bool resolved: If False, consider alarms that have been resolved

        :param list consolidations: List of extra columns to compute for each
          returned result.

        :param dict filter_: Mongo filter.
        :param str search: Search expression in custom DSL.

        :param dict sort: Dict with only one key. Key is the name of the column
          to sort, and value is either "ASC" or "DESC".

        :param int limit: Maximum number of alarms to return.
        :param int offset: Number of alarms to skip (pagination)

        :returns: List of sorted alarms + pagination informations
        :rtype: dict
        """

        search_filter, search_context = self.interpret_search(search)
        search_filter = self.translate_filter(search_filter)

        if search_context == 'all':
            # Use only this filter to search
            alarms = 'a'

        else:
            time_filter = self.get_time_filter(
                opened=opened, resolved=resolved,
                tstart=tstart, tstop=tstop
            )

            if time_filter is None:
                return {'alarms': [], 'total': 0, 'first': 0, 'last': 0}

            filter_ = {
                '$and': [
                    time_filter,
                    self.translate_filter(filter_)
                ]
            }

            # Use filter to get results
            alarms = 'b'

            # Use search_filter to get results from previous results
            if search_filter:
                alarms = alarms + 'c'

        sort = self.translate_sort(sort)
        alarms = alarms.sort(sort)

        alarms = alarms.offset(offset).limit(limit)

        res = {
            'alarms': alarms,
            'total': self.get_total(),
            'first': 0,
            'last': 0
        }

        return res

    def count_alarms_by_period(
            self,
            start,
            stop,
            subperiod={'day': 1},
            limit=100,
            query={},
    ):
        """
        Count alarms that have been opened during (stop - start) period.

        :param start: Beginning timestamp of period
        :type start: int

        :param stop: End timestamp of period
        :type stop: int

        :param subperiod: Cut (stop - start) in ``subperiod`` subperiods
        :type subperiod: dict

        :param limit: Counts cannot exceed this value
        :type limit: int

        :param query: Custom mongodb filter for alarms
        :type query: dict

        :return: List in which each item contains an interval and the
                 related count
        :rtype: list
        """

        intervals = Interval.get_intervals_by_period(start, stop, subperiod)

        results = []
        for date in intervals:
            count = self[AlertsReader.ALARM_STORAGE].count(
                data_ids=None,
                timewindow=TimeWindow(start=date['begin'], stop=date['end']),
                window_start_bind=True,
                _filter=query,
            )

            results.append(
                {
                    'date': date,
                    'count': limit if count > limit else count,
                }
            )

        return results
