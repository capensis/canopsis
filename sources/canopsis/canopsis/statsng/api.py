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
    quote_regex, get_influxdb_client
from canopsis.statsng.enums import StatRequestFields
from canopsis.statsng.errors import StatsAPIError, UnknownStatNameError
from canopsis.statsng.queries import AggregationStatQuery, SLAStatQuery


class StatRequest(object):
    """
    A StatRequest is an object containing a request to the statistics API.
    """
    def __init__(self):
        self.stats = None
        self.tstart = None
        self.tstop = None
        self.group_by = []
        self.entity_filter = []
        self.parameters = {}

    @staticmethod
    def from_request_body(body):
        """
        Create a StatRequest from the body of a request.

        :param dict body: The parsed body of a request.
        :rtype: StatRequest
        :raises: StatsAPIError
        """
        request = StatRequest()

        try:
            body = body.copy()
            request.stats = body.pop(StatRequestFields.stats, [])
            request.tstart = body.pop(StatRequestFields.tstart, None)
            request.tstop = body.pop(StatRequestFields.tstop, None)
            request.group_by = body.pop(StatRequestFields.group_by, [])
            request.entity_filter = body.pop(StatRequestFields.filter, [])
        except AttributeError:
            raise StatsAPIError("The request's body should be a JSON object")

        if not isinstance(request.stats, list):
            raise StatsAPIError("The stats field should be a JSON array")
        if not isinstance(request.group_by, list):
            raise StatsAPIError("The group_by field should be a JSON array")
        if not isinstance(request.entity_filter, list):
            raise StatsAPIError("The filter field should be a JSON array")

        if request.tstart is not None:
            try:
                request.tstart = float(request.tstart)
            except (TypeError, ValueError):
                raise StatsAPIError("The tstart field should be a number")
        if request.tstop is not None:
            try:
                request.tstop = float(request.tstop)
            except (TypeError, ValueError):
                raise StatsAPIError("The tstop field should be a number")

        request.parameters = body

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

    def add_stat(self, tags, stat_name, stat_value):
        """
        Add statistics to the results.

        :param Dict[str, str] tags: the tags of the statistics
        :param str stat_name: the name of the statistic
        :param Any stat_value: the value of the statistic
        """
        data_key = tuple(
            tags[tag] for tag in self.group_by
        )

        if data_key not in self.data:
            self.data[data_key] = {
                "tags": tags
            }

        self.data[data_key][stat_name] = stat_value

    def as_list(self):
        """
        Return the results as a list of dictionnaries.

        :rtype: List[Dict[str, Any]]
        """
        return list(self.data.values())


class EntityFilterParser(object):
    """
    The EntityFilterParser class converts entity filters into influxdb
    conditions.
    """
    def parse_tag_filter(self, tag_name, tag_filter):
        """
        Convert a tag filter `tag_name: tag_filter` into an influxdb condition.

        :param str tag_name:
        :param Union[Dict[str, str], List[str], str] tag_name:
        :rtype: str
        """
        if isinstance(tag_filter, dict) and \
           StatRequestFields.matches in tag_filter:
            regex = tag_filter[StatRequestFields.matches]
            return "{} =~ {}".format(quote_ident(tag_name),
                                     quote_regex(regex))

        elif isinstance(tag_filter, list):
            return "({})".format(" OR ".join(
                "{} = {}".format(quote_ident(tag_name),
                                 quote_literal(tag_value))
                for tag_value in tag_filter
            ))

        elif isinstance(tag_filter, basestring):
            return "{} = {}".format(quote_ident(tag_name),
                                    quote_literal(tag_filter))

        raise StatsAPIError('Invalid tag filter : {}'.format(tag_filter))

    def parse_entity_group(self, entity_group):
        """
        Convert an entity group into an influxdb condition.

        :param Dict[str, Any] entity_group:
        :rtype: str
        """
        return " AND ".join(
            self.parse_tag_filter(tag_name, tag_filter)
            for tag_name, tag_filter in entity_group.items()
        )

    def parse(self, entity_filter):
        """
        Convert an entity filter into an influxdb condition.

        :param List[Dict[str, Any]] entity_filter:
        :rtype: str
        """
        return " OR ".join(
            "({})".format(self.parse_entity_group(entity_group))
            for entity_group in entity_filter
        )


class StatsAPI(object):
    """
    The StatsAPI class handles the requests to the statistics API.
    """
    def __init__(self, logger):
        self.logger = logger
        self.entity_filter_parser = EntityFilterParser()
        influxdb_client = get_influxdb_client()

        self.stat_queries = {
            'alarms_canceled': AggregationStatQuery(
                logger, influxdb_client, 'alarms_canceled', 'sum'),
            'alarms_created': AggregationStatQuery(
                logger, influxdb_client, 'alarms_created', 'sum'),
            'alarms_resolved': AggregationStatQuery(
                logger, influxdb_client, 'alarms_resolved', 'sum'),
            'mean_ack_time': AggregationStatQuery(
                logger, influxdb_client, 'ack_time', 'mean'),
            'mean_resolve_time': AggregationStatQuery(
                logger, influxdb_client, 'resolve_time', 'mean'),
            'ack_time_sla': SLAStatQuery(
                logger, influxdb_client, 'ack_time', 'ack_time_sla'),
            'resolve_time_sla': SLAStatQuery(
                logger, influxdb_client, 'resolve_time', 'resolve_time_sla'),
        }

    def _generate_where_statement(self, request):
        """
        Generate a WHERE statement from a request.

        :param StatRequest request:
        """
        conditions = []

        if request.tstart is not None:
            conditions.append(
                'time >= {:.0f}'.format(request.tstart * SECONDS))
        if request.tstop is not None:
            conditions.append(
                'time < {:.0f}'.format(request.tstop * SECONDS))

        if request.entity_filter:
            conditions.append(
                self.entity_filter_parser.parse(request.entity_filter))

        return ' AND '.join(conditions)

    def handle_request(self, request):
        """
        Handle a request to the statistics API.

        :param StatRequest request:
        :rtype dict:
        :raises: StatsAPIError
        """
        results = StatsAPIResults(request.group_by)

        # Generate WHERE statement
        where = self._generate_where_statement(request)

        # Generate GROUP BY statement
        group_by = ', '.join(quote_ident(tag) for tag in request.group_by)

        for stat_name in request.stats:
            try:
                stat_query = self.stat_queries[stat_name]
            except KeyError:
                raise UnknownStatNameError(
                    'Unknown stat: {0}'.format(stat_name))

            # Add the stats to results
            for tags, stats in stat_query.run(where, group_by,
                                              request.parameters):
                results.add_stat(tags, stat_name, stats)

        return results.as_list()
