# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

_AGGREGATIONS = {}


class AggregationError(Exception):
    """Handle Aggregation errors."""


def get_aggregations():
    """Get aggregation functions by name."""

    result = _AGGREGATIONS.copy()

    return result


def get_aggregation(name):
    """Get registered aggregation.

    :param str name: aggregation name to retrieve.
    :return: aggregation function registered with input name.
    """

    return _AGGREGATIONS.get(name.lower())


def get_aggregation_value(name, points):
    """Get aggregation value where with related name and points.

    Points must have only real values.
    """

    aggregation = get_aggregation(name)

    if aggregation is None:
        raise NotImplementedError("No aggregation {0} exists".format(name))

    result = aggregation(points)

    return result


def add_aggregation(name, function=None, push=False):
    """Set an aggregation function to this AGGREGATIONS module variable.

    - push : if False, raise an AggregationError if an aggregation has already
      been added with the same name.
    - push : change of aggregation if name already exists.

    Added aggregations are available through module properties.
    """

    if push is False and name.lower() in _AGGREGATIONS:
        raise AggregationError("name {0} already exists".format(name))

    def _setfunc(function):

        _AGGREGATIONS[name.lower()] = function

        return function

    if function is None:
        return _setfunc

    else:
        _setfunc(function)

    return function


NONE = 'NONE'
add_aggregation(NONE, None)


MEAN = 'MEAN'
AVERAGE = 'AVERAGE'


@add_aggregation(MEAN)
def _mean(points):
    """Calculate mean of points."""
    return sum(points) / len(points)
add_aggregation(AVERAGE, _mean)


LAST = 'LAST'


@add_aggregation(LAST)
def _last(points):
    """Get the last point."""
    return points[-1]


FIRST = 'FIRST'


@add_aggregation(FIRST)
def _first(points):
    """Get the first point."""
    return points[0]


DELTA = 'DELTA'


@add_aggregation(DELTA)
def _delta(points):
    """Get the delta value."""
    return (max(points) - min(points)) / 2


SUM = 'SUM'


@add_aggregation(SUM)
def _sum(points):
    """Get the sum."""
    return sum(points)


MAX = 'MAX'


@add_aggregation(MAX)
def _max(points):
    """Get the max."""
    return max(points)


MIN = 'MIN'


@add_aggregation(MIN)
def _min(points):
    """Get the min."""
    return min(points)
