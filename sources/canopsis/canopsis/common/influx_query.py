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

from canopsis.common.influx import SECONDS, quote_ident


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
    """
    def __init__(self, measurement):
        self.measurement = measurement
        self.select_columns = []
        self.group_by_tags = []
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

        # Generate where statement
        where_statement = ''
        if self.where_conditions:
            where_statement = 'WHERE {}'.format(
                ' AND '.join(
                    '({})'.format(condition)
                    for condition in self.where_conditions))

        # Generate group by statement
        group_by_statement = ''
        if self.group_by_tags:
            group_by_statement = 'GROUP BY {}'.format(
                ', '.join(
                    quote_ident(tag)
                    for tag in self.group_by_tags))

        return (
            "SELECT {select_columns} "
            "FROM {measurement} "
            "{where_statement} "
            "{group_by_statement}"
        ).format(
            select_columns=select_columns,
            measurement=quote_ident(self.measurement),
            where_statement=where_statement,
            group_by_statement=group_by_statement)

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

    def after(self, timestamp):
        """
        Add a condition `time >= timestamp` to the WHERE statement.

        :param Union[int, float] timestamp:
        """
        if timestamp is not None:
            return self.where('time >= {:.0f}'.format(timestamp * SECONDS))

        return self

    def before(self, timestamp):
        """
        Add a condition `time < timestamp` to the WHERE statement.

        :param Union[int, float] timestamp:
        """
        if timestamp is not None:
            return self.where('time < {:.0f}'.format(timestamp * SECONDS))

        return self
