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
from time import time

from canopsis.alerts import AlarmField
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_config
from canopsis.configuration.model import Parameter
from canopsis.task.core import get_task

from canopsis.timeserie.timewindow import Interval, TimeWindow

from canopsis.alerts.search.interpreter import interpret

from canopsis.tools.schema import get as get_schema

from canopsis.entitylink.manager import Entitylink
from canopsis.pbehavior.manager import PBehaviorManager


CONF_PATH = 'alerts/manager.conf'

DEFAULT = 'ALERTS'
DEFAULT_CNT = []

CACHE = 'COUNT_CACHE'
CACHE_CNT = [
    Parameter('expiration', int),
    Parameter('opened_truncate', Parameter.bool),
    Parameter('opened_limit', int),
    Parameter('resolved_truncate', Parameter.bool),
    Parameter('resolved_limit', int)
]


@conf_paths(CONF_PATH)
@add_config({DEFAULT: DEFAULT_CNT, CACHE: CACHE_CNT})
class AlertsReader(MiddlewareRegistry):
    """
    Alarm cycle managment.

    Used to retrieve events related to alarms in a TimedStorage.
    """

    CONFIG_STORAGE = 'config_storage'
    ALARM_STORAGE = 'alarm_storage'

    def __init__(self, alarm_storage=None, *args, **kwargs):
        super(AlertsReader, self).__init__(*args, **kwargs)

        if alarm_storage is not None:
            self[AlertsReader.ALARM_STORAGE] = alarm_storage

        self.pbm = PBehaviorManager()
        self.llm = Entitylink()

        self.count_cache = {}

        self.grammar = join(prefix, 'etc/alerts/search/grammar.bnf')

    @property
    def alarm_fields(self):
        if not hasattr(self, '_alarm_fields'):
            self._alarm_fields = get_schema('alarm_fields')

        return self._alarm_fields

    @property
    def cache_config(self):
        if not hasattr(self, '_cache_config'):
            values = self.conf.get(CACHE)

            self._cache_config = {
                'expiration': values.get('expiration').value,
                'resolved_truncate': values.get('resolved_truncate').value,
                'resolved_limit': values.get('resolved_limit').value,
                'opened_truncate': values.get('opened_truncate').value,
                'opened_limit': values.get('opened_limit').value
            }

        return self._cache_config

    @cache_config.setter
    def cache_config(self, value):
        if not hasattr(self, '_cache_config'):
            self._cache_config = value

        else:
            self._cache_config.update(value)

    def _translate_key(self, key):
        if key in self.alarm_fields['properties']:
            return self.alarm_fields['properties'][key]['stored_name']

        return key

    def _translate_filter(self, filter_):
        """
        Translate a mongo filter key names. Input keys are UI column names and
        output keys are corresponding keys in the alarm collection.

        :param dict filter_: Mongo filter written by an user

        :return: Mongo filter usable in a query
        :rtype: dict
        """

        if type(filter_) is list:
            for i, f in enumerate(filter_):
                filter_[i] = self._translate_filter(f)

        elif type(filter_) is dict:
            for key, value in filter_.items():
                new_value = self._translate_filter(value)
                filter_[key] = new_value

                new_key = self._translate_key(key)
                filter_[new_key] = filter_.pop(key)

        return filter_

    def _translate_sort(self, key, dir_):
        """
        Translate sort parameters.

        :param str key: UI column name to sort
        :param str dir_: Direction ('ASC' or 'DESC')

        :return: Key usable in a sort operation and translated direction for
          pymongo
        :rtype: tuple

        :raises ValueError: If dir_ is not 'ASC' nor 'DESC'
        """

        if dir_ not in ['ASC', 'DESC']:
            raise ValueError(
                'Sort direction must be "ASC" or "DESC" (got "{}")'.format(
                    dir_
                )
            )

        tkey = self._translate_key(key)
        tdir = 1 if dir_ == 'ASC' else -1

        return tkey, tdir

    def _get_time_filter(self, opened, resolved, tstart, tstop):
        """
        Transform opened, resolved, tstart and tstop parameters into a mongo
        filter. This filter is specific to alarms collection.

        :param bool opened: If True, select opened alarms
        :param bool resolved: If True, select resolved alarms

        :param tstart: Timestamp
        :param tstop: Timestamp
        :type tstart: int or None
        :type tstop: int or None

        :return: Specific mongo filter or None if opened and resolved are False
        :rtype: dict or None
        """

        if opened and resolved:
            if tstart is None and tstop is None:
                return {}

            else:
                return {
                    '$or': [
                        self._get_opened_time_filter(tstart, tstop),
                        self._get_resolved_time_filter(tstart, tstop)
                    ]
                }

        if opened:
            return self._get_opened_time_filter(tstart, tstop)

        if resolved:
            return self._get_resolved_time_filter(tstart, tstop)

        return None

    def _get_opened_time_filter(self, tstart, tstop):
        """
        Get a specific mongo filter.

        :param tstart: Timestamp
        :param tstop: Timestamp
        :type tstart: int or None
        :type tstop: int or None

        :return: Mongo filter
        :rtype: dict
        """

        if tstop is not None:
            return {
                'v.resolved': None,
                't': {'$lte': tstop}
            }

        elif tstart is not None:
            return {
                'v.resolved': None,
                't': {'$lte': tstart}
            }

        else:
            return {'v.resolved': None}

    def _get_resolved_time_filter(self, tstart, tstop):
        """
        Get a specific mongo filter.

        :param tstart: Timestamp
        :param tstop: Timestamp
        :type tstart: int or None
        :type tstop: int or None

        :return: Specific mongo filter
        :rtype: dict
        """

        if tstart is not None and tstop is not None:
            return {
                'v.resolved': {'$ne': None},
                '$or': [
                    {'t': {'$gte': tstart, '$lte': tstop}},
                    {'v.resolved': {'$gte': tstart, '$lte': tstop}},
                    {'t': {'$lte': tstart}, 'v.resolved': {'$gte': tstop}}
                ]
            }

        elif tstart is not None:
            return {'v.resolved': {'$ne': None, '$gte': tstart}}

        elif tstop is not None:
            return {
                'v.resolved': {'$ne': None},
                't': {'$lte': tstop}
            }

        else:
            return {'v.resolved': {'$ne': None}}

    def interpret_search(self, search):
        """
        Parse a search expression to return a mongo filter and a search scope.

        :param str search: Search expression

        :return: Scope ('this' or 'all') and filter (dict)
        :rtype: tuple

        :raises ValueError: If search is not grammatically correct
        """

        if not search:
            return ('this', {})

        return interpret(search, grammar_file=self.grammar)

    def _lookup(self, alarms, lookups):
        """
        Add extra keys to a list of alarms.

        :param list alarms: List of alarms as dict
        :param list lookups: List of extra keys to add.

        :return: Alarms with extra keys
        :rtype: list
        """

        for l in lookups:
            task = get_task(
                'alerts.lookup.{}'.format(l),
                cacheonly=True
            )

            if task is None:
                raise ValueError('Unknown lookup "{}"'.format(l))

            for a in alarms:
                a = task(self, a)

        return alarms

    def clean_fast_count_cache(self):
        """
        Clean expired entries related to fast count cache.
        """

        now = int(time())

        to_clean = []

        for key, value in self.count_cache.items():
            if now >= value['expiration']:
                to_clean.append(key)

        for key in to_clean:
            self.count_cache.pop(key)

    def _get_fast_count(
        self,
        query,
        tstart, tstop, opened, resolved,
        filter_, search
    ):
        """
        Select the best way to count a query documents (try to avoid as much
        as possible mongo's .count()), and return a count as fast as possible.
        Returned value can be accurate, approximated or truncated.

        :param Cursor query: PyMongo Cursor of query that has to be counted

        :param tstart: Timestamp
        :type tstart: int or None
        :param tstop: Timestamp
        :type tstop: int or None

        :param bool opened: If True, query is about opened alarms
        :param bool resolved: If True, query is about resolved alarms

        :param dict filter_: Potential mongo filter of query
        :param str search: Potential search expression of query

        :return: Tuple with count and a boolean telling if count was truncated
        :rtype: tuple (int, bool)
        """

        if resolved:
            cache_key = '{}-{}-{}-{}-{}-{}'.format(
                tstart, tstop, opened, resolved, filter_, search)

            now = int(time())

            truncate = self.cache_config.get('resolved_truncate')
            limit = self.cache_config.get('resolved_limit')

            count_cache = self.count_cache.get(cache_key, None)
            if count_cache is not None:
                if not now >= count_cache['expiration']:
                    count = count_cache['value']
                    if truncate and count == limit:
                        return count, True

                    else:
                        return count, False

            # No cache entry found (or no up-to-date entry)
            if truncate:
                count = query.limit(limit).count(True)

            else:
                count = query.count(True)

            self.count_cache[cache_key] = {
                'value': count,
                'expiration': now + self.cache_config.get('expiration')
            }

            if truncate and count == limit:
                return count, True

            else:
                return count, False

        # Opened alarms only
        else:
            if self.cache_config.get('opened_truncate'):
                limit = self.cache_config.get('opened_limit')
                count = query.limit(limit).count(True)

                if count == limit:
                    return count, True

                else:
                    return count, False

            else:
                return query.count(True), False

    def get(
            self,
            tstart=None,
            tstop=None,
            opened=True,
            resolved=False,
            lookups=[],
            filter_={},
            search='',
            sort_key='opened',
            sort_dir='DESC',
            skip=0,
            limit=50
    ):
        """
        Return filtered, sorted and paginated alarms.

        :param tstart: Beginning timestamp of requested period
        :param tstop: End timestamp of requested period
        :type tstart: int or None
        :type tstop: int or None

        :param bool opened: If True, consider alarms that are currently opened
        :param bool resolved: If True, consider alarms that have been resolved

        :param list lookups: List of extra columns to compute for each
          returned alarm. Extra columns are "pbehaviors" and/or "linklist".

        :param dict filter_: Mongo filter
        :param str search: Search expression in custom DSL

        :param str sort_key: Name of the column to sort
        :param str sort_dir: Either "ASC" or "DESC"

        :param int skip: Number of alarms to skip (pagination)
        :param int limit: Maximum number of alarms to return

        :returns: List of sorted alarms + pagination informations
        :rtype: dict
        """

        time_filter = self._get_time_filter(
            opened=opened, resolved=resolved,
            tstart=tstart, tstop=tstop
        )

        if time_filter is None:
            return {'alarms': [], 'total': 0, 'first': 0, 'last': 0}

        search_context, search_filter = self.interpret_search(search)
        search_filter = self._translate_filter(search_filter)

        if search_context == 'all':
            filter_ = {'$and': [time_filter, search_filter]}

        else:
            filter_ = self._translate_filter(filter_)

            filter_ = {'$and': [time_filter, filter_]}

            if search_filter:
                filter_ = {'$and': [filter_, search_filter]}

        query = self[AlertsReader.ALARM_STORAGE]._backend.find(filter_)

        sort_key, sort_dir = self._translate_sort(sort_key, sort_dir)
        query = query.sort(sort_key, sort_dir)

        query = query.skip(skip)
        query = query.limit(limit)

        alarms = list(query)
        limited_total = len(alarms)  # Manual count is much faster than mongo's

        count_query = self[AlertsReader.ALARM_STORAGE]._backend.find(filter_)
        total, truncated = self._get_fast_count(
            count_query,
            tstart, tstop, opened, resolved,
            filter_, search
        )

        first = 0 if limited_total == 0 else skip + 1
        last = 0 if limited_total == 0 else skip + limited_total

        alarms = self._lookup(alarms, lookups)

        # Steps are never meant to be used in UI and would just cost
        # unnecessary bandwidth.
        for a in alarms:
            a['v'].pop(AlarmField.steps.value)

        res = {
            'alarms': alarms,
            'total': total,
            'truncated': truncated,
            'first': first,
            'last': last
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
