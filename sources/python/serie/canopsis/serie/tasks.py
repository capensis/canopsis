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

from canopsis.task.core import register_task, get_task


_OPERATORS_CACHE = set()


@register_task
def register_operator(operator):
    _OPERATORS_CACHE.add(operator.upper())

    return register_task('serie.op.{0}'.format(operator))


@register_operator('first')
def op_first(values):
    return values[0]


@register_operator('last')
def op_last(values):
    return values[-1]


@register_operator('average')
def op_average(values):
    return sum(values) / len(values)


@register_operator('min')
def op_min(values):
    return min(values)


@register_operator('max')
def op_max(values):
    return max(values)


@register_operator('sum')
def op_sum(values):
    return sum(values)


@register_operator('sub')
def op_sub(values):
    result = values[0]

    for value in values[1:]:
        result -= value

    return result


@register_operator('mul')
def op_mul(values):
    result = values[0]

    for value in values[1:]:
        result *= value

    return result


@register_operator('div')
def op_div(values):
    result = values[0]

    for value in values[1:]:
        result /= value

    return result


@register_task('operatorset')
def serie_operatorset(manager, perfdatas):
    def call_operator(op):
        task = get_task('serie.op.{}'.format(op))

        return lambda regex: task(
            manager.subset_perfdata_values_at_x(
                regex,
                globals()['x'],  # defined in sand-boxed environment
                perfdatas
            )
        )

    operators = {
        key: call_operator(key.lower())
        for key in _OPERATORS_CACHE
    }

    return operators
