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

"""
Alarm reader manager.

TODO: replace the storage class parameter with a collection (=> rewriting count())
"""

from __future__ import unicode_literals

import re
from os.path import join
from time import time

from canopsis.alerts.manager import Alerts
from canopsis.alerts.search.interpreter import interpret
from canopsis.common import root_path
from canopsis.common.collection import MongoCollection
from canopsis.common.redis_store import RedisStore
from canopsis.common.utils import get_sub_key
from canopsis.confng import Configuration, Ini
from canopsis.confng.helpers import cfg_to_bool
from canopsis.logger import Logger
from canopsis.common.middleware import Middleware
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.task.core import get_task
from canopsis.timeserie.timewindow import Interval, TimeWindow
from canopsis.tools.schema import get as get_schema

rconn = RedisStore.get_default()

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

    DEFAULT_ACTIVE_COLUMNS = ["v.component",
                              "v.connector",
                              "v.resource",
                              "v.connector_name"]

    def __init__(self, logger, config, storage, pbehavior_manager):
        """
        :param logger: a logger object
        :param config: a confng instance
        :param storage: a storage instance
        :param pbehavior_manager: a pbehavior manager instance
        """
        self.logger = logger
        self.config = config
        self.alarm_storage = storage
        self.alarm_collection = MongoCollection(self.alarm_storage._backend)
        self.pbehavior_manager = pbehavior_manager
        self.pbh_filter = None

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

        self.count_cache = {}

        self.grammar = join(root_path, self.GRAMMAR_FILE)
        self.has_active_pbh = None

    @classmethod
    def provide_default_basics(cls):
        """
        Provide logger, config, storages...

        ! Do not use in tests !

        :rtype: Union[logging.Logger,
                      canospis.confng.simpleconf.Configuration,
                      canopsis.storage.core.Storage,
                      canopsis.pbehavior.manager.PBehaviorManager]
        """
        logger = Logger.get('alertsreader', cls.LOG_PATH)
        conf = Configuration.load(Alerts.CONF_PATH, Ini)
        alerts_storage = Middleware.get_middleware_by_uri(
            Alerts.ALERTS_STORAGE_URI
        )

        pbm = PBehaviorManager(*PBehaviorManager.provide_default_basics())

        return (logger, conf, alerts_storage, pbm)

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

    @classmethod
    def __convert_to_bool(cls, value):
        """Take a string and return the corresponding boolean. This method is
        case insensitive. Raise a a ValueError if the string can not be parsed.
        :param str value: a string containing the following value true, false
        : return bool: True or false"""
        if isinstance(value, bool):
            return value
        if value.lower() == "true":
            return True
        if value.lower() == "false":
            return False
        msg_err = "Can not convert {} to a boolean. true or false (case insensitive)"
        raise ValueError(msg_err.format(value))

    def _filter_list(self, filter_):
        for item in filter_:
            self._filter(item)

    def _filter_dict(self, filter_):
        for key in filter_:
            if key == "has_active_pb":
                self.has_active_pbh = self.__convert_to_bool(filter_[key])
                del filter_[key]
                return
            else:
                self._filter(filter_[key])

    def _filter(self, filter_):
        if isinstance(filter_, dict):
            self._filter_dict(filter_)

        elif isinstance(filter_, list):
            self._filter_list(filter_)

    def parse_filter(self, filter_):
        """Set self.has_active_pbh true if the filter contain a active_pb key
        set to true or false if it set to false. If the key is not present or
        set to None, None. This method store the first occurrence.
        :param dict alarms: a filter from the brick listalarm
        """

        if filter_ is not None:
            self._filter(filter_)

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

    def add_pbh_filter(self, pipeline, filter_):
        """Add to the aggregation pipeline the stages to filter the alarm
        with their pbehavior.
        :param list pipeline: the aggregation pipeline
        :param filter_ the filter received from the front."""
        self.parse_filter(filter_)
        pipeline.append({"$lookup": {
            "from": "default_pbehavior",
            "localField": "d",
            "foreignField": "eids",
            "as": "pbehaviors"}})

        if self.has_active_pbh is not None:
            tnow = int(time())
            stage = {
                "$project": {
                    "pbehaviors": {
                        "$filter": {
                            "as": "pbh",
                            "input": "$pbehaviors",
                            "cond":
                            {
                                "$and":
                                [
                                    {"$lte": ["$$pbh.tstart", tnow]},
                                    {"$gte": ["$$pbh.tstop", tnow]}
                                ]
                            }
                        }
                    },
                    "_id": 1,
                    "v": 1,
                    "d": 1,
                    "t": 1,
                    "entity": 1
                }
            }
            pipeline.append(stage)

            pbh_filter = {"$match": {"pbehaviors": None}}

            if self.has_active_pbh is True:
                pbh_filter["$match"]["pbehaviors"] = {"$ne": []}
            if self.has_active_pbh is False:
                pbh_filter["$match"]["pbehaviors"] = {"$eq": []}

            pipeline.append(pbh_filter)
        self.has_active_pbh = None

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
          returned alarm. Extra columns are "pbehaviors".

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

        :param bool hide_resources: hide resources' alarms if the component has
        an alarm

        :returns: List of sorted alarms + pagination informations
        :rtype: dict
        """
        if sort_key == 'v.duration':
            sort_key = 'v.creation_date'
        elif sort_key == 'v.current_state_duration':
            sort_key = 'v.state.t'
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
        sort_key, sort_dir = self._translate_sort(sort_key, sort_dir)

        final_filter = self._get_final_filter(
            filter_, time_filter, search, active_columns
        )

        if sort_key[-1] == '.':
            sort_key = 'v.last_update_date'


        total = self.alarm_collection.find(final_filter).count()

        if limit is None:
            limit = total

        # truncate results if more than required
        api_limit = limit

        # get a little bit more results so we may avoid querying the database
        # more than once.
        if hide_resources:
            limit = limit * 2
            hide_resources &= rconn.exists('featureflag:hide_resources')


        def search_aggregate(skip, limit):
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
                    "$match": {"$or": [
                        {"entity.enabled": True}, {
                            "entity": {"$exists": False}}
                    ]}
                }, {
                    "$match": final_filter
                }, {
                    "$sort": {
                        sort_key: sort_dir
                    }
                }
            ]

            if not with_steps:
                pipeline.insert(0, {"$project": {"v.steps": False}})

            self.add_pbh_filter(pipeline, filter_)

            pipeline.append({
                "$skip": skip
            })

            if limit is not None:
                pipeline.append({"$limit": limit})

            result = self.alarm_collection.aggregate(
                pipeline, allowDiskUse=True, cursor={}
            )

            alarms = list(result)
            # Manual count is much faster than mongo's
            truncated = len(alarms)

            res = {
                'alarms': alarms,
                'truncated': truncated,
            }

            return res

        def offset_aggregate(results, skip, limit, filters):
            """
            :param dict results:
            :param int skip:
            :param int limit:
            :param list filters: list of functions to apply on alarms
            """
            tmp_res = search_aggregate(skip, limit)
            pre_filter_len = len(tmp_res['alarms'])

            # no results, all good
            if tmp_res['alarms']:
                results['truncated'] |= tmp_res['truncated']

                # filter useless data
                for filter_ in filters:
                    tmp_res['alarms'] = filter_(tmp_res['alarms'])

                results['alarms'].extend(tmp_res['alarms'])

                if skip < total:
                    results['first'] = 1+skip
                    results['last'] = skip+min(len(tmp_res['alarms']), limit)
                else:
                    results['first'] = skip
                    results['last'] = skip

                skip += limit

            truncated_by = pre_filter_len - len(tmp_res['alarms'])

            return results, skip, truncated_by

        def loop_aggregate(skip, limit, filters, post_sort):
            """
            :param int skip:
            :param int limit:
            :param list filters: list of functions to apply on alarms. Called
                only in offset_aggregate
            :param bool post_sort: post filtering sort with sort_key
                and sort_dir on alarms.
            """
            len_alarms = 0
            results = {
                'alarms': [],
                'total': total,
                'truncated': False,
                'first': 1+skip,
                'last': skip+min(api_limit, total)
            }
            if skip > total:
                results['first'] = skip
                results['last'] = skip

            while len(results['alarms']) < api_limit:
                results, skip, truncated_by = offset_aggregate(
                    results,
                    skip,
                    limit,
                    filters,
                )

                len_alarms = len(results['alarms'])

                # premature break in case we do not have any filter that could
                # modify the real count.
                # this condition cannot be embedded in while <cond> because the
                # loop needs to be ran at least one time.
                if not filters:
                    break

                # filters did not filtered any thing: we don't need to loop
                # again, even if we don't have enough results.
                elif filters and truncated_by == 0:
                    break

            if post_sort:
                results['alarms'] = self._aggregate_post_sort(
                    results['alarms'], sort_key, sort_dir
                )

            if len_alarms > api_limit:
                results['alarms'] = results['alarms'][0:api_limit]

            if limit >= total:
                results['total'] = len(results['alarms'])

            return results

        filters = []
        post_sort = False
        if hide_resources:
            post_sort = True
            filters.append(self._hide_resources)

        return loop_aggregate(skip, limit, filters, post_sort=post_sort)

    @staticmethod
    def _aggregate_post_sort(alarms, sort_key, sort_dir):
        return sorted(
            alarms,
            key=lambda k: get_sub_key(k, sort_key),
            reverse=(sort_dir == -1)
        )

    @staticmethod
    def _hide_resources(alarms):
        """
        Reads alarm_hideresources_resource:<connector>/<connector_name>/<resource>/<component>:drop
        key from redis. if such key exists then the alarm is removed
        from the result set.
        """
        filtered_alarms = []
        for alarm in alarms:
            if alarm['v'].get('resource', '') == '':
                filtered_alarms.append(alarm)
                continue

            drop_id = 'alarm_hideresources_resource:{}/{}/{}/{}:drop'.format(
                alarm['v'].get('connector'),
                alarm['v'].get('connector_name'),
                alarm['v'].get('resource'),
                alarm['v'].get('component'),
            )

            drop_value = rconn.get(drop_id)
            to_drop = False
            try:
                to_drop = drop_value is not None
            except (TypeError, ValueError):
                pass

            if not to_drop:
                filtered_alarms.append(alarm)

        return filtered_alarms

    def count_alarms_by_period(
            self,
            start,
            stop,
            subperiod=None,
            limit=100,
            query=None,
    ):
        """
        Count alarms that have been opened during (stop - start) period.

        :param start: Beginning timestamp of period
        :type start: int

        :param stop: End timestamp of period
        :type stop: int

        :param subperiod: Cut (stop - start) in ``subperiod`` subperiods.
        :type subperiod: dict

        :param limit: Counts cannot exceed this value
        :type limit: int

        :param query: Custom mongodb filter for alarms
        :type query: dict

        :return: List in which each item contains an interval and the
                 related count
        :rtype: list
        """

        if subperiod is None:
            subperiod = {'day': 1}

        if query is None:
            query = {}

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


def remove_resources_alarms(alarms):
    """
    take a list of alarms on a component and remove resources' alarms if needed
    :param list alarms: list of alarms
    :rtype: list
    """
    state_comp = 0
    states_list = []
    alarm_comp = {}
    for i in alarms:
        val = i.get('v').get('state').get('val')
        states_list.append(val)
        if i.get('entity', {}).get('type') == 'component':
            state_comp = val
            alarm_comp = i

    if state_comp >= max(states_list):
        return [alarm_comp]

    ret_val = [alarm_comp]
    for alarm in alarms:
        if alarm.get('v', {}).get('state', {}).get('val') > state_comp:
            ret_val.append(alarm)

    return ret_val
