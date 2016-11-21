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

from re import escape

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.influxdb.core import InfluxDBStorage
from canopsis.common.utils import ensure_iterable

CONF_PATH = 'stats/manager.conf'
CATEGORY = 'STATS'
CONTENT = []


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class Stats(MiddlewareRegistry):
    """
    Stats management
    """

    def __init__(
        self,
        influxdbstg=None,
        *args,
        **kwargs
    ):
        super(Stats, self).__init__(*args, **kwargs)

        if influxdbstg is not None:
            self.influxdbstg = influxdbstg
        else:
            self.influxdbstg = InfluxDBStorage()

    def get_event_stats(self, tstart, tstop, tags={}):
        """
        Get event related stats

        :param int tstart: Start timestamp
        :param int tstop: End timestamp
        :param dict tags: Query formatting tags

        :returns: Stats or an empty dict if tags={}
        :rtype: dict

        Example:

        tags = {'domain': ['d1', 'd2'], 'perimeter': ['p1']}
        return = {
          'domain': [
            {
              'name': 'd1',
              'stats_count': {
                'alarms_new': 12,
                'alarms_ack': 3,
                'alarms_solved_ack': 10,
                'alarms_solved_without_ack': 10
              },
              'stats_delay': {
                'ack_delay_min' : 3,
                'ack_delay_max' : 10,
                'ack_delay_avg' : 5,
                'ack_solved_delay_min' : 3,
                'ack_solved_delay_max' : 10,
                'ack_solved_delay_avg' : 5,
                'alarm_solved_delay_min' : 3,
                'alarm_solved_delay_max' : 10,
                'alarm_solved_delay_avg' : 5
              }
            },
            {
              'name': 'd1',
              [...]
            },
            {
              'name': '__total__',
              [...]
            }
          ],
          'perimeter': [
            {
              'name': 'p1',
              [...]
            },
            {
              'name': '__total__',
              [...]
          ],
          '__total__': {
            'stats_count': {...},
            'stats_delay': {...}
          }
        }
        """
        self.logger.info("Retrieving event based stats")

        result = {}

        time_where = 'time >= {}s AND time <= {}s'.format(tstart, tstop)

        tags_where_total = ''
        for tag_group, tag_names in tags.items():
            tags_where_total += '{} =~ {} OR '.format(
                tag_group,
                self._influx_or_regex(tag_names)
            )
        else:
            # Remove trailing ' OR '
            tags_where_total = tags_where_total[:-4]

        where_total = '{} AND ({})'.format(time_where, tags_where_total)

        g_alarm_new_total = self._influx_query(
            metric_id='alarm_opened_count',
            aggregations='count',
            condition=where_total
        )

        g_ack_new_total = self._influx_query(
            metric_id='alarm_ack_delay',
            aggregations=['count', 'min', 'mean', 'max'],
            condition=where_total
        )

        g_ack_solved_total = self._influx_query(
            metric_id='alarm_ack_solved_delay',
            aggregations=['count', 'min', 'mean', 'max'],
            condition=where_total
        )

        g_alarm_solved_total = self._influx_query(
            metric_id='alarm_solved_delay',
            aggregations=['count', 'min', 'mean', 'max'],
            condition=where_total
        )

        al_new_cnt_tot = self._get_stats(
            result_set=g_alarm_new_total, column='count')

        ack_new_cnt_tot = self._get_stats(
            result_set=g_ack_new_total, column='count')
        ack_new_min_tot = self._get_stats(
            result_set=g_ack_new_total, column='min')
        ack_new_avg_tot = self._get_stats(
            result_set=g_ack_new_total, column='mean')
        ack_new_max_tot = self._get_stats(
            result_set=g_ack_new_total, column='max')

        ack_sol_cnt_tot = self._get_stats(
            result_set=g_ack_solved_total, column='count')
        ack_sol_min_tot = self._get_stats(
            result_set=g_ack_solved_total, column='min')
        ack_sol_avg_tot = self._get_stats(
            result_set=g_ack_solved_total, column='mean')
        ack_sol_max_tot = self._get_stats(
            result_set=g_ack_solved_total, column='max')

        al_sol_cnt_tot = self._get_stats(
            result_set=g_alarm_solved_total, column='count')
        al_sol_min_tot = self._get_stats(
            result_set=g_alarm_solved_total, column='min')
        al_sol_avg_tot = self._get_stats(
            result_set=g_alarm_solved_total, column='mean')
        al_sol_max_tot = self._get_stats(
            result_set=g_alarm_solved_total, column='max')

        no_ack_sol_tot = al_sol_cnt_tot - ack_sol_cnt_tot

        result['__total__'] = {
            'stats_count': {
                'alarms_new': al_new_cnt_tot,
                'alarms_ack': ack_new_cnt_tot,
                'alarms_solved_ack': ack_sol_cnt_tot,
                'alarms_solved_without_ack': no_ack_sol_tot
            },
            'stats_delay': {
                'ack_delay_min': ack_new_min_tot,
                'ack_delay_avg': int(ack_new_avg_tot),
                'ack_delay_max': ack_new_max_tot,
                'ack_solved_delay_min': ack_sol_min_tot,
                'ack_solved_delay_avg': int(ack_sol_avg_tot),
                'ack_solved_delay_max': ack_sol_max_tot,
                'alarm_solved_delay_min': al_sol_min_tot,
                'alarm_solved_delay_avg': int(al_sol_avg_tot),
                'alarm_solved_delay_max': al_sol_max_tot
            }
        }

        for tag_group, tag_names in tags.items():
            self.logger.debug(
                "Retrieving stats for tag id '{}'".format(tag_group))

            result[tag_group] = []

            tag_regex = self._influx_or_regex(tag_names)
            where = '{} AND {} =~ {}'.format(time_where, tag_group, tag_regex)

            g_alarm_new = self._influx_query(
                metric_id='alarm_opened_count',
                aggregations='count',
                condition=where,
                groupby=tag_group
            )

            g_ack_new = self._influx_query(
                metric_id='alarm_ack_delay',
                aggregations=['count', 'min', 'mean', 'max'],
                condition=where,
                groupby=tag_group
            )

            g_ack_solved = self._influx_query(
                metric_id='alarm_ack_solved_delay',
                aggregations=['count', 'min', 'mean', 'max'],
                condition=where,
                groupby=tag_group
            )

            g_alarm_solved = self._influx_query(
                metric_id='alarm_solved_delay',
                aggregations=['count', 'min', 'mean', 'max'],
                condition=where,
                groupby=tag_group
            )

            for tag_name in tag_names:
                self.logger.debug(
                    "Retrieving stats for tag value '{}'".format(tag_name))

                al_new_cnt = self._get_stats(
                    result_set=g_alarm_new,
                    tags={tag_group: tag_name},
                    column='count'
                )

                ack_new_cnt = self._get_stats(
                    result_set=g_ack_new,
                    tags={tag_group: tag_name},
                    column='count'
                )
                ack_new_min = self._get_stats(
                    result_set=g_ack_new,
                    tags={tag_group: tag_name},
                    column='min'
                )
                ack_new_avg = self._get_stats(
                    result_set=g_ack_new,
                    tags={tag_group: tag_name},
                    column='mean'
                )
                ack_new_max = self._get_stats(
                    result_set=g_ack_new,
                    tags={tag_group: tag_name},
                    column='max'
                )

                ack_sol_cnt = self._get_stats(
                    result_set=g_ack_solved,
                    tags={tag_group: tag_name},
                    column='count'
                )
                ack_sol_min = self._get_stats(
                    result_set=g_ack_solved,
                    tags={tag_group: tag_name},
                    column='min'
                )
                ack_sol_avg = self._get_stats(
                    result_set=g_ack_solved,
                    tags={tag_group: tag_name},
                    column='mean'
                )
                ack_sol_max = self._get_stats(
                    result_set=g_ack_solved,
                    tags={tag_group: tag_name},
                    column='max'
                )

                al_sol_cnt = self._get_stats(
                    result_set=g_alarm_solved,
                    tags={tag_group: tag_name},
                    column='count'
                )
                al_sol_min = self._get_stats(
                    result_set=g_alarm_solved,
                    tags={tag_group: tag_name},
                    column='min'
                )
                al_sol_avg = self._get_stats(
                    result_set=g_alarm_solved,
                    tags={tag_group: tag_name},
                    column='mean'
                )
                al_sol_max = self._get_stats(
                    result_set=g_alarm_solved,
                    tags={tag_group: tag_name},
                    column='max'
                )

                no_ack_sol = al_sol_cnt - ack_sol_cnt

                stats = {
                    'name': tag_name,
                    'stats_count': {
                        'alarms_new': al_new_cnt,
                        'alarms_ack': ack_new_cnt,
                        'alarms_solved_ack': ack_sol_cnt,
                        'alarms_solved_without_ack': no_ack_sol
                    },
                    'stats_delay': {
                        'ack_delay_min': ack_new_min,
                        'ack_delay_avg': int(ack_new_avg),
                        'ack_delay_max': ack_new_max,
                        'ack_solved_delay_min': ack_sol_min,
                        'ack_solved_delay_avg': int(ack_sol_avg),
                        'ack_solved_delay_max': ack_sol_max,
                        'alarm_solved_delay_min': al_sol_min,
                        'alarm_solved_delay_avg': int(al_sol_avg),
                        'alarm_solved_delay_max': al_sol_max
                    }
                }

                result[tag_group].append(stats)

            # Insert total stats for tag_group at the end of the list
            g_alarm_new_total = self._influx_query(
                metric_id='alarm_opened_count',
                aggregations='count',
                condition=where
            )

            g_ack_new_total = self._influx_query(
                metric_id='alarm_ack_delay',
                aggregations=['count', 'min', 'mean', 'max'],
                condition=where
            )

            g_ack_solved_total = self._influx_query(
                metric_id='alarm_ack_solved_delay',
                aggregations=['count', 'min', 'mean', 'max'],
                condition=where
            )

            g_alarm_solved_total = self._influx_query(
                metric_id='alarm_solved_delay',
                aggregations=['count', 'min', 'mean', 'max'],
                condition=where
            )

            al_new_cnt_tot = self._get_stats(
                result_set=g_alarm_new_total, column='count')

            ack_new_cnt_tot = self._get_stats(
                result_set=g_ack_new_total, column='count')
            ack_new_min_tot = self._get_stats(
                result_set=g_ack_new_total, column='min')
            ack_new_avg_tot = self._get_stats(
                result_set=g_ack_new_total, column='mean')
            ack_new_max_tot = self._get_stats(
                result_set=g_ack_new_total, column='max')

            ack_sol_cnt_tot = self._get_stats(
                result_set=g_ack_solved_total, column='count')
            ack_sol_min_tot = self._get_stats(
                result_set=g_ack_solved_total, column='min')
            ack_sol_avg_tot = self._get_stats(
                result_set=g_ack_solved_total, column='mean')
            ack_sol_max_tot = self._get_stats(
                result_set=g_ack_solved_total, column='max')

            al_sol_cnt_tot = self._get_stats(
                result_set=g_alarm_solved_total, column='count')
            al_sol_min_tot = self._get_stats(
                result_set=g_alarm_solved_total, column='min')
            al_sol_avg_tot = self._get_stats(
                result_set=g_alarm_solved_total, column='mean')
            al_sol_max_tot = self._get_stats(
                result_set=g_alarm_solved_total, column='max')

            no_ack_sol_tot = al_sol_cnt_tot - ack_sol_cnt_tot

            stats = {
                'name': '__total__',
                'stats_count': {
                    'alarms_new': al_new_cnt_tot,
                    'alarms_ack': ack_new_cnt_tot,
                    'alarms_solved_ack': ack_sol_cnt_tot,
                    'alarms_solved_without_ack': no_ack_sol_tot
                },
                'stats_delay': {
                    'ack_delay_min': ack_new_min_tot,
                    'ack_delay_avg': int(ack_new_avg_tot),
                    'ack_delay_max': ack_new_max_tot,
                    'ack_solved_delay_min': ack_sol_min_tot,
                    'ack_solved_delay_avg': int(ack_sol_avg_tot),
                    'ack_solved_delay_max': ack_sol_max_tot,
                    'alarm_solved_delay_min': al_sol_min_tot,
                    'alarm_solved_delay_avg': int(al_sol_avg_tot),
                    'alarm_solved_delay_max': al_sol_max_tot
                }
            }

            result[tag_group].append(stats)

        self.logger.debug("Event based stats have succesfully been retrieved")

        return result

    def get_user_stats(self, tstart, tstop, users=[], tags={}):
        """
        Get user related stats

        :param int tstart: Start timestamp
        :param int tstop: End timestamp
        :param list users: Users whose stats want to be retrieved
        :param dict tags: Groups stats by value

        :returns: Stats or an empty list if users=[]
        :rtype: list

        Example:

        users = ['u1', 'u2']
        tags = {'domain': ['d1', 'd2'], 'perimeter': ['p1']}
        result = [
          {
            'author': 'MMA',
            'ack': {
              'total': 12,
              'delay_min': 3,
              'delay_avg': 5,
              'delay_max': 10,
            },
            'session': {
              'duration_min': 12,
              'duration_avg': 50,
              'duration_max': 300,
            },
            'tags': {
              'domain': [
                {
                  'name': 'd1',
                  'ack_total': 3
                },
                {
                  'name': 'd2',
                  'ack_total': 4
                }
              ],
              'perimeter': [
                {
                  'name': 'p1',
                  'ack_total': 7
                }
              ]
            }
          }
        ]
        """
        self.logger.info("Retrieving user based stats")

        result = []

        if not users:
            return result

        time_where = 'time >= {}s AND time <= {}s'.format(tstart, tstop)

        user_regex = self._influx_or_regex(users)
        user_where = 'component =~ {}'.format(user_regex)

        tags_where = ''
        for tag_group, tag_names in tags.items():
            tags_where += '{} =~ {} OR '.format(
                tag_group,
                self._influx_or_regex(tag_names)
            )
        else:
            # Remove trailing ' OR '
            tags_where = tags_where[:-4]

        ack_new_where = '{} AND {} AND ({})'.format(
            time_where,
            user_where,
            tags_where
        )

        g_ack_new = self._influx_query(
            metric_id='alarm_ack_delay',
            aggregations=['count', 'min', 'mean', 'max'],
            condition=ack_new_where,
            groupby='component'
        )

        duration_where = '{} AND {}'.format(time_where, user_where)

        g_duration = self._influx_query(
            metric_id='session_duration',
            aggregations=['min', 'mean', 'max'],
            condition=duration_where,
            groupby='component'
        )

        for user in users:
            self.logger.debug("Retrieving stats for user '{}'".format(user))

            ack_new_cnt = self._get_stats(
                result_set=g_ack_new, tags={'component': user}, column='count')
            ack_new_min = self._get_stats(
                result_set=g_ack_new, tags={'component': user}, column='min')
            ack_new_avg = self._get_stats(
                result_set=g_ack_new, tags={'component': user}, column='mean')
            ack_new_max = self._get_stats(
                result_set=g_ack_new, tags={'component': user}, column='max')

            duration_min = self._get_stats(
                result_set=g_duration, tags={'component': user}, column='min')
            duration_avg = self._get_stats(
                result_set=g_duration, tags={'component': user}, column='mean')
            duration_max = self._get_stats(
                result_set=g_duration, tags={'component': user}, column='max')

            stats = {
                'author': user,
                'ack': {
                    'total': ack_new_cnt,
                    'delay_min': ack_new_min,
                    'delay_avg': int(ack_new_avg),
                    'delay_max': ack_new_max
                },
                'session': {
                    'duration_min': duration_min,
                    'duration_avg': duration_avg,
                    'duration_max': duration_max
                },
                'tags': {}
            }

            result.append(stats)

        for tag_group, tag_names in tags.items():
            self.logger.debug(
                "Retrieving stats for tag id '{}'".format(tag_group)
            )

            tag_regex = self._influx_or_regex(tag_names)
            tag_where = '{} AND {} =~ {}'.format(
                user_where, tag_group, tag_regex)

            g_ack_new = self._influx_query(
                metric_id='alarm_ack_delay',
                aggregations='count',
                condition=tag_where,
                groupby=['component', tag_group]
            )

            for tag_name in tag_names:
                self.logger.debug(
                    "Retrieving stats for tag value '{}'".format(tag_name)
                )

                for i in range(len(users)):
                    ack_new = self._get_stats(
                        result_set=g_ack_new,
                        tags={'component': users[i], tag_group: tag_name},
                        column='count'
                    )

                    result[i]['tags'].setdefault(tag_group, [])
                    result[i]['tags'][tag_group].append(
                        {
                            'name': tag_name,
                            'ack_total': ack_new
                        }
                    )

        self.logger.debug("User based stats have succesfully been retrieved")

        return result

    def _influx_or_regex(self, items):
        """
        Transform a list of strings in an influx regex with an or filter
        between each element.

        :param list items: List of strings

        :return: Ready to use regex
        :rtype: str
        """
        regex = '/^('

        for item in items:
            regex = '{}{}|'.format(regex, escape(item))

        # Truncate regex for the last '|'
        regex = '{})$/'.format(regex[:-1])

        return regex

    def _influx_query(
            self,
            metric_id,
            aggregations=[],
            condition=None,
            groupby=[]
    ):
        query = 'select'

        if not aggregations:
            query = '{} value'.format(query)

        else:
            for aggr in ensure_iterable(aggregations):
                query = '{} {}(value),'.format(query, aggr)

            else:
                # Remove last ','
                query = query[:-1]

        query = '{} from {}'.format(query, metric_id)

        if condition is not None:
            query = '{} where {}'.format(query, condition)

        if groupby:
            query = '{} group by'.format(query)
            for gb in ensure_iterable(groupby):
                query = '{} {},'.format(query, gb)

            # Remove last ','
            else:
                query = query[:-1]

        self.logger.debug('Processing query `{}`'.format(query))
        result = self.influxdbstg.raw_query(query)

        # If something went wrong (bad metric_id, bad aggregation...),
        # InfluxDBStorage will just return None
        if result is None:
            raise ValueError(
                'Query `{}` failed : unable to retrieve stats'.format(query)
            )

        return result

    def _get_stats(self, result_set, tags={}, column=None):
        self.logger.debug(
            "Retrieving stats on column {} for tags '{}'".format(
                column,
                tags
            )
        )

        try:
            stats = next(
                result_set.get_points(tags=tags)
            )

        # If there is no entry for tag name, stats cannot be iterated on
        except StopIteration:
            return 0

        if column is None:
            return stats

        else:
            # stats[column] might be None if requested period does not have
            # any data.
            if stats[column] is not None:
                return stats[column]

            else:
                return 0
