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
    pass


def get_aggregations():
    """
    Get aggregation functions by name.
    """

    result = _AGGREGATIONS.copy()

    return result


def get_aggregation_value(name, points):
    """
    Get aggregation value where with related name and points.

    Points must have only real values.
    """

    aggregation = _AGGREGATIONS.get(name, None)
    if aggregation is None:
        raise NotImplementedError("No aggregation {0} exists".format(name))

    result = aggregation(points)

    return result


def add_aggregation(name, function, push=False):
    """
    Set an aggregation function to this AGGREGATIONS module variable.

    - push : if False, raise an AggregationError if an aggregation has already
      been added with the same name.
    - push : change of aggregation if name already exists.

    Added aggregations are available through module properties.
    """

    if push is False and name in _AGGREGATIONS:
        raise AggregationError("name {0} already exists".format(name))

    _AGGREGATIONS[name] = function


add_aggregation('NONE', None)


def _mean(points):
    return sum(points) / len(points)
add_aggregation('MEAN', _mean)
add_aggregation('AVERAGE', _mean)


def _last(points):
    return points[-1]
add_aggregation('LAST', _last)


def _first(points):
    return points[0]
add_aggregation('FIRST', _first)


def _delta(points):
    return (max(points) - min(points)) / 2
add_aggregation('DELTA', _delta)


def _sum(points):
    return sum(points)
add_aggregation('SUM', _sum)


def _max(points):
    return max(points)
add_aggregation('MAX', _max)


def _min(points):
    return min(points)
add_aggregation('MIN', _min)
