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

from __future__ import unicode_literals

"""
Alarm reader manager.

TODO: replace the storage class parameter with a collection (=> rewriting count())
"""

import re

from os.path import join
from time import time

from canopsis.common import root_path

from canopsis.alerts.enums import AlarmField
from canopsis.alerts.manager import Alerts
from canopsis.alerts.search.interpreter import interpret
from canopsis.confng import Configuration, Ini
from canopsis.confng.helpers import cfg_to_bool
from canopsis.entitylink.manager import Entitylink
from canopsis.logger import Logger
from canopsis.middleware.core import Middleware
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.task.core import get_task
from canopsis.timeserie.timewindow import Interval, TimeWindow
from canopsis.tools.schema import get as get_schema

DEFAULT_EXPIRATION = 1800
DEFAULT_OPENED_TRUNC = True
DEFAULT_OPENED_LIMIT = 200000
DEFAULT_RESOLVED_TRUNC = True
DEFAULT_RESOLVED_LIMIT = 1000


class AlertsReader(object):
    """
    Alarm cycle managment.

    Used to retrieve events related to alarms in a TimedStorage.
    """

    LOG_PATH = 'var/log/alertsreader.log'
    CONF_PATH = 'etc/alerts/manager.conf'
    CATEGORY = 'COUNT_CACHE'
    GRAMMAR_FILE = 'etc/alerts/search/grammar.bnf'

    DEFAULT_ACTIVE_COLUMNS = [
        "v.component",
        "v.connector",
        "v.resource",
        "v.connector_name"
    ]

    def __init__(self, logger, config, storage,
                 pbehavior_manager, entitylink_manager):
        """
        :param logger: a logger object
        :param config: a confng instance
        :param storage: a storage instance
        :param pbehavior_manager: a pbehavior manager instance
        :param entitylink_manager: a entitylink manager instance
        """
        self.logger = logger
        self.config = config
        self.alarm_storage = storage
        self.pbehavior_manager = pbehavior_manager
        self.entitylink_manager = entitylink_manager

        category = self.config.get(self.CATEGORY, {})
        self.expiration = int(category.get('expiration', DEFAULT_EXPIRATION))
        self.opened_truncate = cfg_to_bool(category.get('opened_truncate',
                                                        DEFAULT_OPENED_TRUNC))
        self.opened_limit = int(category.get('opened_limit',
                                             DEFAULT_OPENED_LIMIT))
        self.resolved_truncate = cfg_to_bool(category.get('resolved_truncate',
                                                          DEFAULT_RESOLVED_TRUNC))
        self.resolved_limit = int(category.get('resolved_limit',
                                               DEFAULT_RESOLVED_LIMIT))

        _, pb_storage = PBehaviorManager.provide_default_basics()

        self.count_cache = {}

        self.grammar = join(root_path, self.GRAMMAR_FILE)

    @classmethod
    def provide_default_basics(cls):
        """
        Provide logger, config, storages...

        ! Do not use in tests !

        :rtype: Union[logging.Logger,
                      canospis.confng.simpleconf.Configuration,
                      canopsis.storage.core.Storage,
                      canopsis.pbehavior.manager.PBehaviorManager,
                      canopsis.entitylink.manager.EntityLink]
        """
        logger = Logger.get('alertsreader', cls.LOG_PATH)
        conf = Configuration.load(Alerts.CONF_PATH, Ini)
        alerts_storage = Middleware.get_middleware_by_uri(
            Alerts.ALERTS_STORAGE_URI
        )
        pb_storage = Middleware.get_middleware_by_uri(
            PBehaviorManager.PB_STORAGE_URI
        )

        pbm = PBehaviorManager(logger=logger, pb_storage=pb_storage)
        llm = Entitylink(*Entitylink.provide_default_basics())

        return (logger, conf, alerts_storage, pbm, llm)

    @property
    def alarm_fields(self):
        """
        alarm_field parameter
        """
        if not hasattr(self, '_alarm_fields'):
            self._alarm_fields = get_schema('alarm_fields')

        return self._alarm_fields

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

        if isinstance(filter_, list):
            for i, fil in enumerate(filter_):
                filter_[i] = self._translate_filter(fil)

        elif isinstance(filter_, dict):
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

    @staticmethod
    def _get_opened_time_filter(tstart, tstop):
        """
        Get a specific mongo filter.

        :param tstart: Timestamp
        :param tstop: Timestamp
        :type tstart: int or None
        :type tstop: int or None

        :return: Mongo filter
        :rtype: dict
        """

        if tstop is not None and tstart is not None:
            return {
                'v.resolved': None,
                't': {'$lte': tstop, "$gte": tstart}
            }

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

        return {'v.resolved': None}

    @staticmethod
    def _get_resolved_time_filter(tstart, tstop):
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
                't': {'$gte': tstart, '$lte': tstop}
            }

        elif tstart is not None:
            return {'v.resolved': {'$ne': None, '$gte': tstart}}

        elif tstop is not None:
            return {
                'v.resolved': {'$ne': None},
                't': {'$lte': tstop}
            }

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

        for lookup in lookups:
            task = get_task(
                'alerts.lookup.{}'.format(lookup),
                cacheonly=True
            )

            if task is None:
                raise ValueError('Unknown lookup "{}"'.format(lookup))

            for alarm in alarms:
                alarm = task(self, alarm)

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

    def _get_fast_count(self, query, tstart, tstop, opened, resolved, filter_,
                        search):
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

            truncate = self.resolved_truncate
            limit = self.resolved_limit

            count_cache = self.count_cache.get(cache_key, None)
            if count_cache is not None:
                if not now >= count_cache['expiration']:
                    count = count_cache['value']
                    return count, truncate and count == limit

            # No cache entry found (or no up-to-date entry)
            if truncate:
                count = query.limit(limit).count(True)
            else:
                count = query.count(True)

            self.count_cache[cache_key] = {
                'value': count,
                'expiration': now + self.expiration
            }

            return count, truncate and count == limit

        # Opened alarms only
        else:
            if self.opened_truncate:
                limit = self.opened_limit
                count = query.limit(limit).count(True)

                return count, count == limit

            return query.count(True), False

    def _get_final_filter(
        self, view_filter, time_filter, search, active_columns
    ):
        """
        Computes the real filter:

        The view filter and time filter are always part of the final filter,
        if not empty.

        In the search matches the BNF grammar,
        it is appended to the final filter.

        Otherwise, regex on columns is made.


        All filters are aggregated with $and.


        {
            '$and': [
                view_filter,
                time_filter,
                bnf_filter | column_filter
            ]
        }

        :param view_filter dict: the filter given by the canopsis view.
        :param time_filter dict: hehe. dunno.
        :param search str: text to search in columns, or a BNF valid search as
            defined by the grammar in etc/search/grammar.bnf

            The BNF grammar is tried first, if the string does not comply with
            the grammar, column search is used instead.
        :param active_columns list[str]: list of columns to search in.
            in a column ends with '.' it will be ignored.

            The 'd' column is always added.
        """
        final_filter = {'$and': []}

        t_view_filter = self._translate_filter(view_filter)
        # add the view filter if not empty
        if view_filter not in [None, {}]:
            final_filter['$and'].append(t_view_filter)

        if time_filter not in [None, {}]:
            final_filter['$and'].append(time_filter)

        # try grammar search
        try:
            _, bnf_search_filter = self.interpret_search(search)
            bnf_search_filter = self._translate_filter(bnf_search_filter)
        except ValueError:
            bnf_search_filter = None

        if bnf_search_filter is not None:
            final_filter['$and'].append(bnf_search_filter)

        else:
            escaped_search = re.escape(str(search))
            column_filter = {'$or': []}
            for column in active_columns:
                column_filter['$or'].append(
                    {
                        column: {
                            '$regex': '.*{}.*'.format(escaped_search),
                            '$options': 'i'
                        }
                    }
                )
            column_filter['$or'].append(
                {
                    'd': {
                        '$regex': '.*{}.*'.format(escaped_search),
                        '$options': 'i'
                    }
                }
            )

            final_filter['$and'].append(column_filter)

        return final_filter

    def get(
            self,
            tstart=None,
            tstop=None,
            opened=True,
            resolved=False,
            lookups=None,
            filter_=None,
            search='',
            sort_key='opened',
            sort_dir='DESC',
            skip=0,
            limit=None,
            with_steps=False,
            natural_search=False,
            active_columns=None,
            hide_resources=False
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

        :param str sort_key: Name of the column to sort. If the value ends with
            a dot '.', sort_key is replaced with 'v.last_update_date'.
        :param str sort_dir: Either "ASC" or "DESC"

        :param int skip: Number of alarms to skip (pagination)
        :param int limit: Maximum number of alarms to return

        :param bool with_steps: True if you want alarm steps in your alarm.

        :param bool natural_search: True if you want to use a natural search

        :param list active_columns: the list of alarms columns on which to
        apply the natural search filter.

        :param bool hide_resources: hide alarms on resources if an alarm is
        on a parent's component

        :returns: List of sorted alarms + pagination informations
        :rtype: dict
        """
        if hide_resources:
            if resolved:
                self.logger.error('you only can hide resources on current alarms')
            return self.hide_resources(
                tstart,
                tstop,
                filter_,
                sort_key,
                sort_dir,
                skip,
                limit,
                search,
                natural_search,
                active_columns
            )

        if lookups is None:
            lookups = []

        if filter_ is None:
            filter_ = {}

        if active_columns is None:
            active_columns = self.DEFAULT_ACTIVE_COLUMNS

        time_filter = self._get_time_filter(
            opened=opened, resolved=resolved,
            tstart=tstart, tstop=tstop
        )

        if time_filter is None:
            return {'alarms': [], 'total': 0, 'first': 0, 'last': 0}

        result = None
        sort_key, sort_dir = self._translate_sort(sort_key, sort_dir)

        final_filter = self._get_final_filter(
            filter_, time_filter, search, active_columns
        )

        if sort_key[-1] == '.':
            sort_key = 'v.last_update_date'

        pipeline = [
            {
                "$lookup": {
                    "from": "default_entities",
                    "localField": "d",
                    "foreignField": "_id",
                    "as": "entity"
                }
            }, {
                "$unwind": {
                    "path": "$entity",
                    "preserveNullAndEmptyArrays": True,
                }
            }, {
                "$match": final_filter
            }, {
                "$sort": {
                    sort_key: sort_dir
                }
            }, {
                "$skip": skip
            }
        ]

        if limit is not None:
            pipeline.append({"$limit": limit})

        result = self.alarm_storage._backend.aggregate(pipeline, cursor={})

        alarms = list(result)
        limited_total = len(alarms)  # Manual count is much faster than mongo's

        count_query = self.alarm_storage._backend.find(final_filter)
        total, truncated = self._get_fast_count(
            count_query,
            tstart, tstop, opened, resolved,
            final_filter, search
        )

        first = 0 if limited_total == 0 else skip + 1
        last = 0 if limited_total == 0 else skip + limited_total

        alarms = self._lookup(alarms, lookups)

        if not with_steps:
            for alarm in alarms:
                alarm['v'].pop(AlarmField.steps.value)

        res = {
            'alarms': alarms,
            'total': total,
            'truncated': truncated,
            'first': first,
            'last': last
        }

        return res

    def hide_resources(
        self,
        tstart=None,
        tstop=None,
        filter_={},
        sort_key='opened',
        sort_dir='DESC',
        skip=0,
        limit=None,
        search='',
        natural_search=False,
        active_columns=None
    ):
        """
        :rtype: dict
        """
        if filter_ is None:
            filter_ = {}

        if active_columns is None:
            active_columns = self.DEFAULT_ACTIVE_COLUMNS

        time_filter = self._get_time_filter(
            opened=True,
            resolved=False,
            tstart=tstart,
            tstop=tstop
        )

        if time_filter is None:
            return {'alarms': [], 'total': 0, 'first': 0, 'last': 0}

        final_filter = self._get_final_filter(
            filter_,
            time_filter,
            search,
            active_columns
        )

        sort_key, sort_dir = self._translate_sort(sort_key, sort_dir)

        pipeline = [
            {
                "$lookup": {
                    "from": "default_entities",
                    "localField": "d",
                    "foreignField": "_id",
                    "as": "entity"
                }
            }, {
                "$unwind": {
                    "path": "$entity",
                    "preserveNullAndEmptyArrays": True,
                }
            }, {
                "$match": final_filter
            }, {
                "$sort": {
                    sort_key: sort_dir
                }
            }
        ]
        alarms = self.alarm_storage._backend.aggregate(pipeline).get('result')

        alarms_with_resources_hided = self.hide_resources_from_alarm_list(alarms)

        len_before_truncate = len(alarms_with_resources_hided)
        if limit is not None:
            last = limit + skip
            alarms_with_resources_hided = alarms_with_resources_hided[skip:last]
        else:
            alarms_with_resources_hided = alarms_with_resources_hided[skip:]
            last = len(alarms_with_resources_hided)
        len_after_truncate = len(alarms_with_resources_hided)

        #il manque la recherche
        ret_val = {
            'alarms': alarms_with_resources_hided,
            'total': len_after_truncate,
            'truncated': len_after_truncate < len_before_truncate,
            'first': skip,
            'last': last
        }
        return ret_val

    def hide_resources_from_alarm_list(self, alarms):
        """
        :rtype: list
        """
        alarm_dict = {}
        res = []
        for alarm in alarms:
            component = alarm.get('v').get('component')
            if component not in alarm_dict:
                alarm_dict[component] = [alarm]
            else:
                alarm_dict[component].append(alarm)

        for component in alarm_dict:
            for alarm_comp in alarm_dict[component]:
                if 'component' in [alarm_comp.get('entity', {}).get('type', '')]:
                    alarms_with_component = self.check_alarm_list_with_component(alarm_dict[component])
                    res = sum([res, alarms_with_component], [])
                else:
                    res = sum([res, alarm_dict[component]], [])

        return res

    def check_alarm_list_with_component(self, alarms):
        """
        :rtype: list
        """
        state_comp = 0
        states_list = []
        alarm_comp = {}
        for i in alarms:
            val = i.get('v').get('state').get('val')
            states_list.append(val)
            if i.get('entity').get('type') == 'component':
                state_comp = val
                alarm_comp = i

        if state_comp >= max(states_list):
            return [alarm_comp]

        return alarms

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

        :param int start: Beginning timestamp of period
        :param int stop: End timestamp of period
        :param dict subperiod: Cut (stop - start) in ``subperiod`` subperiods
        :param int limit: Counts cannot exceed this value
        :param dict query: Custom mongodb filter for alarms
        :return: List in which each item contains an interval and the
                 related count
        :rtype: list
        """

        intervals = Interval.get_intervals_by_period(start, stop, subperiod)

        results = []
        for date in intervals:
            count = self.alarm_storage.count(
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
