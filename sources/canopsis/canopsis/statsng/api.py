# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

from canopsis.common.influx import SECONDS, quote_ident, quote_literal, \
    get_influxdb_client
from canopsis.statsng.enums import StatRequestFields
from canopsis.statsng.queries import AggregationStatQuery


class StatsAPIError(Exception):
    """
    A StatsAPIError is an Exception that can be raised by a StatsAPI object.

    It should be handled in `webcore/services/statsng.py`, and returned as a
    JSON object in the response.
    """
    def __init__(self, message):
        super(StatsAPIError, self).__init__(message)
        self.message = message


class UnknownStatNameError(StatsAPIError):
    """
    A UnknownStatNameError is an Exception that can be raised by a StatsAPI
    object when requesting an unknown statistic.
    """
    pass


class StatRequest(object):
    """
    A StatRequest is an object containing a request to the statistics API.
    """
    def __init__(self):
        self.stats = None
        self.tstart = None
        self.tstop = None
        self.group_by = []
        self.entity_filter = None

    @staticmethod
    def from_request_body(body):
        """
        Create a StatRequest from the body of a request.

        :param dict body: The parsed body of a request.
        :rtype: StatRequest
        :raises: StatsAPIError
        """
        request = StatRequest()

        request.stats = body.pop(StatRequestFields.stats, None)
        request.tstart = body.pop(StatRequestFields.tstart, None)
        request.tstop = body.pop(StatRequestFields.tstop, None)
        request.group_by = body.pop(StatRequestFields.group_by, [])
        request.entity_filter = body.pop(StatRequestFields.filter, None)

        if body:
            raise StatsAPIError('Unexpected fields: {0}'.format(
                ', '.join(body.keys())))

        return request


class StatsAPIResults(object):
    """
    A StatsAPIResults object stores the results of a request to the statistics
    API.

    :param List[str] group_by: the list of tags used in the GROUP BY statement
    """
    def __init__(self, group_by):
        self.group_by = group_by
        self.data = {}

    def add_stats(self, tags, stats):
        """
        Add statistics to the results.

        :param Dict[str, str] tags: the tags of the statistics
        :param Dict[str, Any] stats: the values of the statistics
        """
        data_key = tuple(
            tags[tag] for tag in self.group_by
        )

        if data_key not in self.data:
            self.data[data_key] = {
                "tags": tags
            }

        self.data[data_key].update(stats)

    def as_list(self):
        """
        Return the results as a list of dictionnaries.

        :rtype: List[Dict[str, Any]]
        """
        return list(self.data.values())


class StatsAPI(object):
    """
    A StatsAPI object handles the computation of statistics.
    """
    def __init__(self, logger):
        self.logger = logger
        self.influxdb_client = get_influxdb_client()

        self.stat_queries = {
            'alarms_canceled': AggregationStatQuery('alarms_canceled', 'sum'),
            'alarms_created': AggregationStatQuery('alarms_created', 'sum'),
            'alarms_resolved': AggregationStatQuery('alarms_resolved', 'sum'),
        }

    def _generate_where_statement(self, request):
        """
        Generate a WHERE statement from a request.

        :param StatRequest request:
        """
        conditions = []

        if request.tstart:
            conditions.append('time >= {}'.format(request.tstart * SECONDS))
        if request.tstop:
            conditions.append('time < {}'.format(request.tstop * SECONDS))

        # TODO: Handle request.filter

        if conditions:
            return 'WHERE {}'.format(
                ' AND '.join(conditions)
            )
        else:
            return ''

    def handle_request(self, request):
        """
        Handle a request to the statistics API.

        :param StatRequest request:
        :rtype dict:
        :raises: StatsAPIError
        """
        results = StatsAPIResults(request.group_by)

        # Generate WHERE statement
        where_statement = self._generate_where_statement(request)

        # Generate GROUP BY statement
        group_by_statement = ''
        if request.group_by:
            group_by_statement = 'GROUP BY {}'.format(
                ', '.join(quote_ident(tag) for tag in request.group_by)
            )

        for stat in request.stats:
            try:
                stat_query = self.stat_queries[stat]
            except KeyError:
                raise UnknownStatNameError('Unknown stat: {0}'.format(stat))

            # Generate SELECT statement
            select_statement = stat_query.get_select_statement()

            # Generate the query
            query = " ".join((
                select_statement,
                where_statement,
                group_by_statement))

            # Run the query
            self.logger.debug("Running the query: {0}".format(query))
            result_set = self.influxdb_client.query(query)

            # Add the stats to results
            for (_, tags), rows in result_set.items():
                results.add_stats(tags, stat_query.get_results(rows))

        return results.as_list()
