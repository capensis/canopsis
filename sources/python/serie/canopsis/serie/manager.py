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

"""Module of the serie manager."""

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.configurable.decorator import conf_paths

from canopsis.serie.utils import build_filter_from_regex
from canopsis.timeserie.timewindow import Period, Interval, TimeWindow
from canopsis.timeserie.core import TimeSerie
from canopsis.task.core import get_task

from canopsis.old.mfilter import check

from RestrictedPython import compile_restricted
from RestrictedPython.Guards import safe_builtins

from operator import itemgetter

from math import isnan

from time import time


CONF_PATH = 'serie/manager.conf'
CATEGORY = 'SERIE'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Serie(MiddlewareRegistry):
    """Serie manager."""

    SERIE_STORAGE = 'serie_storage'  #: serie storage name.
    CONTEXT_MANAGER = 'context'  #: serie context manager name.
    PERFDATA_MANAGER = 'perfdata'  #: serie perfdata manager name.

    def __init__(
            self,
            serie_storage=None,
            context=None,
            perfdata=None,
            *args, **kwargs
    ):
        super(Serie, self).__init__(*args, **kwargs)

        if serie_storage is not None:
            self[Serie.SERIE_STORAGE] = serie_storage

        if context is not None:
            self[Serie.CONTEXT_MANAGER] = context

        if perfdata is not None:
            self[Serie.PERFDATA_MANAGER] = perfdata

    def get_metrics(self, regex, metrics=None):
        """
        Get metrics with regular expression from storage or existing set.

        :param regex: "co:regex re:regex me:regex" filter
        :type regex: str

        :param metrics: If not specified, get from storage (optional)
        :type metrics: list

        :returns: list of matching metrics
        """

        mfilter = build_filter_from_regex(regex)

        if metrics is None:
            return self[Serie.CONTEXT_MANAGER].find(
                _type='metric',
                _filter=mfilter
            )

        else:
            result = [
                metric
                for metric in metrics
                if check(mfilter, metric)
            ]

            return result

    def get_perfdata(self, metrics, timewindow=None):
        """
        Internal method to fetch perfdata from metrics.

        :returns: perfdata per metric id as dict
        """

        result = {}

        for metric in metrics:
            mid = self[Serie.CONTEXT_MANAGER].get_entity_id(metric)
            perfdata = self[Serie.PERFDATA_MANAGER].get(
                mid,
                timewindow=timewindow,
                with_meta=False
            )

            result[mid] = {
                'entity': metric,
                'points': perfdata
            }

        return result

    def subset_perfdata_superposed(self, regex, perfdatas):
        """
        Get superposed points of metric matching filter.

        :param regex: filter for ``get_metric()`` method
        :type regex: str

        :param perfdatas: perfdata fetched with ``get_perfdata()`` method
        :type perfdatas: dict

        :returns: superposed points as list
        """

        selected_metrics = [perfdatas[key]['entity'] for key in perfdatas]

        metrics = self.get_metrics(regex, selected_metrics)

        metric_ids = [
            self[Serie.CONTEXT_MANAGER].get_entity_id(metric)
            for metric in metrics
        ]

        points = []

        for metric_id in metric_ids:
            points += perfdatas[metric_id]['aggregated']

        points = sorted(points, key=itemgetter(0))

        return points

    def aggregation(
            self, serieconf, timewindow, period=None, usenan=True, fixed=True
    ):
        """
        Get aggregated perfdata from serie.

        :param serieconf: Serie used for aggregation
        :type serieconf: dict

        :param timewindow: Time window used for perfdata fetching
        :type timewindow: canopsis.timeserie.timewindow.TimeWindow

        :returns: aggregated perfdata classified by metric id as dict
        """

        tw, period, usenan, fixed = self.get_timewindow_period_usenan_fixed(
            serieconf, timewindow, period, usenan, fixed
        )

        timeserie = TimeSerie(
            period=period,
            aggregation=serieconf.get(
                'aggregation_method',
                TimeSerie.VDEFAULT_AGGREGATION
            ),
            round_time=fixed
        )

        metrics = self.get_metrics(serieconf['metric_filter'])
        perfdatas = self.get_perfdata(metrics, timewindow=tw)

        for key in perfdatas:
            perfdatas[key]['aggregated'] = timeserie.calculate(
                points=perfdatas[key]['points'],
                timewindow=tw,
                usenan=usenan
            )

        return perfdatas

    def consolidation(
            self, serieconf, perfdatas, timewindow,
            period=None, usenan=True, fixed=True
    ):
        """
        Get consolidated point from serie.

        :param serieconf: Serie used for consolidation
        :type serieconf: dict

        :param perfdatas: Aggregated perfdatas from ``aggregation()`` method
        :type perfdatas: dict

        :param timewindow: Time window used for consolidation
        :type timewindow: canopsis.timeserie.timewindow.TimeWindow

        :returns: Consolidated points
        """

        # configure consolidation period (same as aggregation period)
        tw, period, usenan, fixed = self.get_timewindow_period_usenan_fixed(
            serieconf, timewindow, period, usenan, fixed
        )

        intervals = Interval.get_intervals_by_period(
            tw.start(),
            tw.stop(),
            period
        )

        points = []

        # generator consolidation operators
        operatorset = get_task('serie.operatorset')

        # generate one point per aggregation interval in timewindow
        for interval in intervals:
            tw = TimeWindow(
                start=interval['begin'],
                stop=interval['end'] - 1
            )

            # operators are acting on a specific timewindow
            operators = operatorset(self, period, perfdatas, tw, usenan)

            # execute formula in sand-boxed environment
            restricted_globals = {
                '__builtins__': safe_builtins,
            }

            restricted_globals.update(operators)

            formula = serieconf['formula']
            code = compile_restricted(formula, '<string>', 'eval')

            try:
                val = eval(code, restricted_globals)

            except Exception as ex:
                self.logger.warning(
                    'Wrong serie formula: {0}/{1} ({2})'.format(
                        serieconf['crecord_name'], formula, ex
                    )
                )
                val = float('nan')

            else:
                if isnan(val):
                    self.logger.warning(
                        'Formula result is nan: {0}/{1}.'.format(
                            serieconf['crecord_name'], formula
                        )
                    )

            # result contains consolidated value
            # point is computed at the start of interval
            points.append((interval['begin'], val))

        return points

    def calculate(self, serieconf, lastts=None):
        """
        Compute serie point.

        :param serieconf: Serie to compute
        :type serieconf: dict

        :param lastts: timestamp corresponding to the last point to calculate.
            Default is now.
        :type lastts: float

        :returns: Computed points
        """

        if lastts is None:
            lastts = time()

        timewindow = TimeWindow(
            start=serieconf.get('last_computation', lastts),
            stop=lastts
        )

        tw, period, usenan, fixed = self.get_timewindow_period_usenan_fixed(
            serieconf, timewindow
        )

        perfdatas = self.aggregation(
            serieconf, tw, period=period, usenan=usenan, fixed=fixed
        )

        points = self.consolidation(
            serieconf, perfdatas, tw,
            period=period, usenan=usenan, fixed=fixed
        )

        points = [point for point in points if point[0] <= lastts]

        serieconf['last_computation'] = lastts
        self[Serie.SERIE_STORAGE].put_element(element=serieconf)

        return points

    def get_series(self, timestamp):
        """
        Get series that need to be computed at specified timestamp.

        :param timestamp: Timestamp used to determine if a serie needs
                          to be computed
        :type timestamp: int

        :returns: list of serie
        """

        storage = self[Serie.SERIE_STORAGE]

        javascript_condition = '({0} - {1}) >= ({2} / {3})'.format(
            timestamp,
            'this.last_computation',
            'this.aggregation_interval',
            'this.computations_per_interval'
        )

        return storage.find_elements(query={'$where': javascript_condition})

    @staticmethod
    def get_timewindow_period_usenan_fixed(
            serieconf, timewindow, period=None, usenan=None, fixed=None
    ):
        """Get the right timewindow, period and usenan."""

        if fixed is None:
            fixed = serieconf.get('round_time_interval', TimeSerie.VROUND_TIME)

        if period is None:
            interval = serieconf.get('aggregation_interval', TimeSerie.VPERIOD)
            period = Period.new(interval)

            if fixed:
                timewindow = timewindow.get_round_timewindow(period=period)

        if usenan is None:
            usenan = serieconf.get('usenan', True)

        result = timewindow, period, usenan, fixed

        return result
