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

"""Module of serie tasks."""

from canopsis.timeserie.aggregation import get_aggregations
from canopsis.timeserie.core import TimeSerie
from canopsis.task.core import register_task


def new_operator(opname, manager, period, perfdatas, timewindow):
    """Create a new operator usable in the restricted Python environment.

    This operator will be used in the formula.

    :param opname: Operator name
    :type opname: str

    :param manager: Serie manager
    :type manager: canopsis.serie.manager.Serie

    :param period: Period used for timeserie
    :type period: canopsis.timeserie.timewindow.Period

    :param perfdatas: Perfdata classified by metric id
    :type perfdatas: dict

    :param timewindow: Time window used for consolidation
    :type timewindow: canopsis.timeserie.timewindow.TimeWindow

    :returns: operator as a callable
    """

    def operator(regex):
        """
        Operator returned by ``new_operator()`` function.

        :param regex: Metric filter used to aggregate perfdatas
        :type regex: str

        :returns: consolidated point as float
        """

        points = manager.subset_perfdata_superposed(regex, perfdatas)
        result = float('nan')

        if points:
            timeserie = TimeSerie(
                period=period,
                aggregation=opname,
                round_time=True
            )

            consolidated = timeserie.calculate(points, timewindow)

            if consolidated:
                result = consolidated[0][1]

        return result

    return operator


@register_task('serie.operatorset')
def serie_operatorset(manager, period, perfdatas, timewindow):
    """
    Generate set of operators.

    :param manager: Serie manager
    :type manager: canopsis.serie.manager.Serie

    :param period: Period used for timeserie
    :type period: canopsis.timeserie.timewindow.Period

    :param perfdatas: Perfdata classified by metric id
    :type perfdatas: dict

    :param timewindow: Time window used for consolidation
    :type timewindow: canopsis.timeserie.timewindow.TimeWindow

    :returns: operators classified by name as dict
    """

    operators = {
        key: new_operator(key.lower(), manager, period, perfdatas, timewindow)
        for key in get_aggregations()
    }

    return operators
