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

from canopsis.common.influx import SECONDS, quote_ident, quote_regex,\
    escape_regex


class SelectColumn(object):
    """
    A SelectColumn is an object encoding a column in a SELECT statement.

    ```
    >>> str(SelectColumn("value", "sum", "total"))
    sum('value') AS 'total'
    ```
    """
    def __init__(self, key, function=None, alias=None):
        self.key = key
        self.function = function
        self.alias = alias or key

    def __str__(self):
        if self.function is None:
            return "{key} AS {alias}".format(
                key=quote_ident(self.key),
                alias=quote_ident(self.alias))

        return "{function}({key}) AS {alias}".format(
            function=self.function,
            key=quote_ident(self.key),
            alias=quote_ident(self.alias))


class SelectQuery(object):
    """
    A SelectQuery is an object encoding an InfluxQL SELECT query.

    ```
    >>> query = (SelectQuery('measurement')
    ...     .select('value')
    ...     .where('type = "watcher"')
    ...     .group_by('id'))
    >>> print(query.build())
    SELECT 'value' FROM 'measurement' WHERE type = "watcher" GROUP_BY 'id'
    ```

    :param Union[str, SelectQuery] from_series: The name of a measurement or a
        subquery.
    """
    def __init__(self, from_series):
        self.from_series = from_series
        self.into_measurement = None
        self.select_columns = []
        self.group_by_tags = []
        self.group_by_time_interval = None
        self.group_by_time_offset = None
        self.group_by_time_fill = None
        self.where_conditions = []

    def build(self):
        """
        Returns the query as a string.

        :rtype: str
        """
        if not self.select_columns:
            raise ValueError(
                "The select statement should contain at least one value.")

        # Generate select statement
        select_columns = ', '.join(
            str(select_column)
            for select_column in self.select_columns)

        # Generate from
        if isinstance(self.from_series, SelectQuery):
            from_series = '({})'.format(self.from_series.build())
        else:
            from_series = quote_ident(self.from_series)

        # Generate into statement
        into_statement = ''
        if self.into_measurement is not None:
            into_statement = 'INTO {}'.format(
                quote_ident(self.into_measurement))

        # Generate where statement
        where_statement = ''
        if self.where_conditions:
            where_statement = 'WHERE {}'.format(
                ' AND '.join(
                    '({})'.format(condition)
                    for condition in self.where_conditions))

        # Generate group by statement
        group_by_tags = []
        if self.group_by_time_interval:
            if self.group_by_time_offset:
                group_by_tags.append('time({}, {})'.format(
                    self.group_by_time_interval,
                    self.group_by_time_offset))
            else:
                group_by_tags.append('time({})'.format(
                    self.group_by_time_interval))

        for tag in self.group_by_tags:
            group_by_tags.append(quote_ident(tag))

        group_by_statement = ''
        if group_by_tags:
            group_by_statement = 'GROUP BY {}'.format(
                ', '.join(group_by_tags))

        fill_statement = ''
        if not (self.group_by_time_interval is None
                or self.group_by_time_fill is None):
            fill_statement = 'fill({})'.format(self.group_by_time_fill)

        return (
            "SELECT {select_columns} "
            "{into_statement} "
            "FROM {from_series} "
            "{where_statement} "
            "{group_by_statement}"
            "{fill_statement}"
        ).format(
            select_columns=select_columns,
            into_statement=into_statement,
            from_series=from_series,
            where_statement=where_statement,
            group_by_statement=group_by_statement,
            fill_statement=fill_statement)

    def select(self, key, function=None, alias=None):
        """
        Add a column to the select statement.

        :param str key: The name of the tag or field to query.
        :param Optional[str] function: An optional InfluxQL function (see
        https://docs.influxdata.com/influxdb/v1.5/query_language/functions/).
        :param Optional[str] alias: An optional alias for the column.
        """
        select_column = SelectColumn(key, function, alias)

        for other_column in self.select_columns:
            if other_column.alias == select_column.alias:
                raise ValueError((
                    "There is already a select column with the alias {}"
                ).format(select_column.alias))

        self.select_columns.append(select_column)
        return self

    def select_all(self):
        """
        Select all the columns.
        """
        self.select_columns = ['*']
        return self

    def into(self, measurement):
        """
        Add an INTO statement.

        :param str measurement: The name of the measurement
        """
        self.into_measurement = measurement
        return self

    def group_by_time(self, interval, offset=None, fill=None):
        """
        Add a GROUP BY time(...) statement.

        interval and offset should be influxdb duration literals, e.g. "50s" or
        "10d". See the documentation for more details:
        https://docs.influxdata.com/influxdb/v1.6/query_language/spec/#durations

        :param str interval: An influxdb duration literal indicating the
            duration of the interval
        :param Optional[str] interval: An influxdb duration literal that shifts
            influxdb's time boundaries
        :param Optional[str] fill: An influxdb fill option (linear, none, null,
            previous or a number)
        """
        self.group_by_time_interval = interval
        self.group_by_time_offset = offset
        self.group_by_time_fill = fill
        return self

    def group_by(self, *tags):
        """
        Add columns to the GROUP BY statement.

        :param List[str] *tags: A list of tag names.
        """
        self.group_by_tags.extend(tags)
        return self

    def where(self, *conditions):
        """
        Add conditions to the WHERE statement.

        :param List[str] *conditions: A list of InfluxQL expressions.
        """
        for condition in conditions:
            if condition:
                self.where_conditions.append(condition)
        return self

    def where_in(self, tag, values):
        """
        Add a condition to the WHERE statement checking that the value of a tag
        is in a list of values.

        :param str tag: The name of the tag
        :param List[str] values: A list of values the tag should be in.
        """
        if not values:
            # The list of values is empty. No value should match this
            # condition.
            return self.where('false')

        # Build a regular expression checking that the tag matches one of the
        # values.
        regex = '^({})$'.format('|'.join(
            escape_regex(value)
            for value in values
        ))
        condition = '{} =~ {}'.format(
            quote_ident(tag),
            quote_regex(regex))
        return self.where(condition)

    def after(self, timestamp):
        """
        Add a condition `time >= timestamp` to the WHERE statement.

        :param Union[int, float] timestamp:
        """
        if timestamp is not None:
            return self.where('time >= {:.0f}ns'.format(timestamp * SECONDS))

        return self

    def before(self, timestamp):
        """
        Add a condition `time < timestamp` to the WHERE statement.

        :param Union[int, float] timestamp:
        """
        if timestamp is not None:
            return self.where('time < {:.0f}ns'.format(timestamp * SECONDS))

        return self

    def copy(self):
        """
        Return a copy of the SelectQuery.

        :rtype: SelectQuery
        """
        from_series = self.from_series
        if isinstance(from_series, SelectQuery):
            from_series = from_series.copy()

        copy = SelectQuery(from_series)
        copy.into_measurement = self.into_measurement
        copy.select_columns = self.select_columns[:]
        copy.group_by_tags = self.group_by_tags[:]
        copy.group_by_time_interval = self.group_by_time_interval
        copy.group_by_time_offset = self.group_by_time_offset
        copy.group_by_time_fill = self.group_by_time_fill
        copy.where_conditions = self.where_conditions[:]

        return copy
