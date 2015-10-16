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

from canopsis.task.core import register_task


@register_task('serie.operator.first')
def serie_operator_first(values):
    return values[0]


@register_task('serie.operator.last')
def serie_operator_last(values):
    return values[-1]


@register_task('serie.operator.average')
def serie_operator_average(values):
    return sum(values) / len(values)


@register_task('serie.operator.min')
def serie_operator_min(values):
    return min(values)


@register_task('serie.operator.max')
def serie_operator_max(values):
    return max(values)


@register_task('serie.operator.sum')
def serie_operator_sum(values):
    return sum(values)


@register_task('serie.operator.sub')
def serie_operator_sub(values):
    result = values[0]

    for value in values[1:]:
        result -= value

    return result


@register_task('serie.operator.mul')
def serie_operator_mul(values):
    result = values[0]

    for value in values[1:]:
        result *= value

    return result


@register_task('serie.operator.div')
def serie_operator_div(values):
    result = values[0]

    for value in values[1:]:
        result /= value

    return result


@register_task('operatorset')
def serie_operatorset(manager, serieconf, perfdatas):
    def perfdata(regex):
        selected_metrics = [
            perfdatas[key]['metric']
            for key in perfdatas.keys()
        ]

        metrics = manager.get_metrics(regex, selected_metrics)
        metric_ids = [
            manager[manager.CONTEXT_MANAGER].get_entity_id(metric)
            for metric in metrics
        ]

        # all perfdata are aggregated with the same period
        # so all x values are the same
        mid = metric_ids[0]
        i = 0

        for point in perfdatas[mid]['aggregated']:
            # x is defined by the sand-boxed environment
            if point[0] == globals()['x']:
                break

            i += 1

        selected_points = [
            perfdatas[key]['aggregated'][i]
            for key in metric_ids
        ]

        return selected_points

    operators = {
        'FIRST': lambda regex: serie_operator_first(perfdata(regex)),
        'LAST': lambda regex: serie_operator_last(perfdata(regex)),
        'AVERAGE': lambda regex: serie_operator_average(perfdata(regex)),
        'MIN': lambda regex: serie_operator_min(perfdata(regex)),
        'MAX': lambda regex: serie_operator_max(perfdata(regex)),
        'SUM': lambda regex: serie_operator_sum(perfdata(regex)),
        'SUB': lambda regex: serie_operator_sub(perfdata(regex)),
        'MUL': lambda regex: serie_operator_mul(perfdata(regex)),
        'DIV': lambda regex: serie_operator_div(perfdata(regex))
    }

    return operators
