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

from canopsis.common.influx import quote_ident
from canopsis.statsng.errors import StatsAPIError


class StatQuery(object):
    """
    A StatQuery is an object that is used by the statistics API to compute
    statistics from the influxdb database.
    """
    def __init__(self, logger, influxdb_client):
        self.logger = logger
        self.influxdb_client = influxdb_client

    def _run_query(self, select_statement, where='', group_by=''):
        """
        Runs an influxdb query.

        Runs the following query:
        ```
        {select_statement}
        WHERE {where}
        GROUP BY {group_by}
        ```

        :param str select_statement:
        :param str where:
        :param str group_by:
        :rtype ResultSet:
        """
        # Generate WHERE statement
        where_statement = ''
        if where:
            where_statement = 'WHERE {}'.format(where)

        # Generate GROUP BY statement
        group_by_statement = ''
        if group_by:
            group_by_statement = 'GROUP BY {}'.format(group_by)

        # Generate the query
        query = ' '.join((
            select_statement,
            where_statement,
            group_by_statement))

        # Run the query
        self.logger.info("Running the query: {0}".format(query))
        return self.influxdb_client.query(query)

    def run(self, where, group_by, parameters):
        """
        Run the StatsQuery

        This is an iterator yielding tuples `(tags, stats)` where `tags` is a
        dictionary containing the values of the tags of `group_by` and `stats`
        is a dictionary containing the values of the statistics for this group.

        :param str where: a condition to be used in a WHERE statement, used to
        set the time interval and to filter the entities
        :param str group_by: a list of comma separated tags to be used in a
        GROUP BY statement
        :param Dict[str, Any] parameters: a dictionary containing additional
        parameters
        :rtype: Iterator[Tuple[Dict[str, str], Dict[str, Any]]]
        """
        raise NotImplementedError()


class AggregationStatQuery(StatQuery):
    """
    An AggregationStatQuery is a StatQuery that aggregates the `value` field in
    a measurement, and returns it as a dictionnary.

    :param str measurement: the name of the measurement
    :param str aggregation: the aggregation function
    :param str name: the name of the statistic
    """
    def __init__(self, logger, influxdb_client, measurement, aggregation):
        super(AggregationStatQuery, self).__init__(logger, influxdb_client)
        self.measurement = measurement
        self.aggregation = aggregation

    def run(self, where, group_by, parameters):
        # Run the query
        select_statement = """
            SELECT {aggregation}(value) AS value
            FROM {measurement}
        """.format(
            aggregation=self.aggregation,
            measurement=quote_ident(self.measurement))

        result_set = self._run_query(select_statement, where, group_by)

        # Yield the results
        for (_, tags), rows in result_set.items():
            # Get first and only row
            try:
                row = next(rows)
            except StopIteration:
                continue

            yield tags, row['value']


class SLAStatQuery(StatQuery):
    """
    An SLAStatQuery is a StatQuery that given an SLA value, returns a
    dictionary containing :
     - the number of values above and below the SLA
     - the percentage of values above and below the SLA

    :param str measurement: the name of the measurement
    :param str name: the name of the statistic
    """
    def __init__(self,
                 logger,
                 influxdb_client,
                 measurement,
                 sla_field):
        super(SLAStatQuery, self).__init__(logger, influxdb_client)
        self.measurement = measurement
        self.sla_field = sla_field

    def run(self, where, group_by, parameters):
        try:
            sla = float(parameters[self.sla_field])
        except KeyError:
            raise StatsAPIError('Missing field: {0}'.format(self.sla_field))
        except (TypeError, ValueError):
            raise StatsAPIError(
                'The {0} field should be a number'.format(self.sla_field))

        below_where = 'value <= {}'.format(sla)
        if where:
            below_where = '({}) AND {}'.format(where, below_where)

        # Run the query
        select_statement = """
            SELECT count(value) AS value
            FROM {measurement}
        """.format(measurement=quote_ident(self.measurement))

        total_result_set = self._run_query(select_statement, where, group_by)
        below_result_set = self._run_query(select_statement,
                                           below_where,
                                           group_by)

        # Yield the results
        for (measurement, tags), total_rows in total_result_set.items():
            below_rows = below_result_set.get_points(measurement, tags)

            # Get first and only row
            try:
                total = next(total_rows)['value']
            except StopIteration:
                continue

            try:
                below = next(below_rows)['value']
            except StopIteration:
                below = 0

            above = total - below

            below_rate = -1
            above_rate = -1
            if total > 0:
                below_rate = below / float(total)
                above_rate = above / float(total)

            results = {
                'below': below,
                'above': above,
                'below_rate': below_rate,
                'above_rate': above_rate,
            }
            yield tags, results
