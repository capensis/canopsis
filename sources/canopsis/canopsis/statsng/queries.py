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

from canopsis.common.influx import quote_ident, quote_literal


class StatQuery(object):
    """
    A StatQuery is an object that is used by the statistics API to compute
    statistics from the influxdb database.
    """
    def get_select_statement(self):
        """
        Generate a SELECT statement (`SELECT ... FROM ...`).

        This statement will be followed by a WHERE statement in the query.
        """
        raise NotImplementedError()

    def get_results(self, rows):
        """
        Given the rows of an influxdb queries, return the corresponding
        statistic(s) in a dictionnary.

        :params generator rows: a generator of dictionnaries
        :rtype: dict
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
    def __init__(self, measurement, aggregation, name=None):
        self.measurement = measurement
        self.aggregation = aggregation
        self.name = name or measurement

    def get_select_statement(self):
        return """
            SELECT {aggregation}(value) AS value
            FROM {measurement}
        """.format(
            aggregation=self.aggregation,
            measurement=quote_ident(self.measurement))

    def get_results(self, rows):
        # Get first and only row
        try:
            row = next(rows)
        except StopIteration:
            return {}

        return {
            self.name: row['value']
        }
