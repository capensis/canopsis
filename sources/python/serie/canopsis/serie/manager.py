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

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.model import Parameter

from canopsis.serie.utils import build_filter_from_regex
from canopsis.timeserie.timewindow import Period
from canopsis.timeserie.core import TimeSerie
from canopsis.task.core import get_task

from canopsis.old.mfilter import check

from RestrictedPython import compile_restricted
from RestrictedPython.Guards import safe_builtins
from time import time


CONF_PATH = 'serie/manager.conf'
CATEGORY = 'SERIE'
CONTENT = [
    Parameter('points_per_interval', int)
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class Serie(MiddlewareRegistry):
    """
    Consolidation manager.
    """

    SERIE_STORAGE = 'serie_storage'
    CONTEXT_MANAGER = 'context'
    PERFDATA_MANAGER = 'perfdata'

    @property
    def points_per_interval(self):
        """
        Maximum number of points in a consolidation interval.
        """

        if not hasattr(self, '_points_per_interval'):
            self.points_per_interval = None

        return self._points_per_interval

    @points_per_interval.setter
    def points_per_interval(self, value):
        if value is None:
            value = 10

        self._points_per_interval = value

    def __init__(
        self,
        points_per_interval=None,
        serie_storage=None,
        context=None,
        perfdata=None,
        *args, **kwargs
    ):
        super(Serie, self).__init__(*args, **kwargs)

        if points_per_interval is not None:
            self.points_per_interval = points_per_interval

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
                _filter=mfilter
            )

        else:
            result = [
                metric
                for metric in metrics
                if check(mfilter, metric)
            ]

            return result

    def get_perfdata(self, metrics, period=None, timewindow=None):
        """
        Internal method to fetch perfdata from metrics.

        :returns: perfdata per metric id as dict
        """

        result = {}

        for metric in metrics:
            mid = self[Serie.CONTEXT_MANAGER].get_entity_id(metric)
            perfdata = self[Serie.PERFDATA_MANAGER].get(
                mid,
                period=period,
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

        selected_metrics = [
            perfdatas[key]['entity']
            for key in perfdatas.keys()
        ]

        metrics = self.get_metrics(regex, selected_metrics)
        metric_ids = [
            self[Serie.CONTEXT_MANAGER].get_entity_id(metric)
            for metric in metrics
        ]

        points = []

        for metric_id in metric_ids:
            points += perfdatas[metric_id]['aggregated']

        points = sorted(points, key=lambda point: point[0])

        return points

    def aggregation(self, serieconf, timewindow=None):
        """
        Get aggregated perfdata from serie.

        :param serieconf: Serie used for aggregation
        :type serieconf: dict

        :param timewindow: Time window used for perfdata fetching (optional)
        :type timewindow: canopsis.timeserie.timewindow.TimeWindow

        :returns: aggregated perfdata classified by metric id as dict
        """

        interval = serieconf.get('aggregation_interval', None)

        if interval is None:
            period = TimeSerie.VPERIOD

        else:
            period = Period(second=interval)

        ts = TimeSerie(
            period=period,
            aggregation=serieconf.get(
                'aggregation_method',
                TimeSerie.VDEFAULT_AGGREGATION
            ),
            round_time=serieconf.get(
                'round_time_interval',
                TimeSerie.VROUND_TIME
            )
        )

        metrics = self.get_metrics(serieconf['metric_filter'])
        perfdatas = self.get_perfdata(metrics, timewindow=timewindow)

        for key in perfdatas:
            perfdatas[key]['aggregated'] = ts.calculate(
                perfdatas[key]['points'],
                timewindow
            )

        return perfdatas

    def consolidation(self, serieconf, perfdatas):
        """
        Get consolidated point from serie.

        :param serieconf: Serie used for consolidation
        :type serieconf: dict

        :param perfdatas: Aggregated perfdatas from ``aggregation()`` method
        :type perfdatas: dict

        :returns: Consolidated point as float
        """

        # configure consolidation period (same as aggregation period)
        interval = serieconf.get('aggregation_interval', None)

        if interval is None:
            period = TimeSerie.VPERIOD

        else:
            period = Period(second=interval)

        # generator consolidation operators
        operatorset = get_task('serie.operatorset')
        operators = operatorset(self, period, perfdatas)

        # execute formula in sand-boxed environment
        restricted_globals = {
            '__builtins__': safe_builtins,
        }

        restricted_globals.update(operators)

        expression = 'result = {0}'.format(serieconf['formula'])
        code = compile_restricted(expression, '<string>', 'exec')

        exec(code) in restricted_globals

        # result contains consolidated point
        return restricted_globals['result']

    def calculate(self, serieconf, timewindow=None):
        """
        Compute serie point.

        :param serieconf: Serie to compute
        :type serieconf: dict

        :param timewindow: Time Window to use for aggregation (optional)
        :type timewindow: canopsis.timeserie.timewindow.TimeWindow

        :returns: Computed point as int
        """

        perfdatas = self.aggregation(serieconf, timewindow)
        point = self.consolidation(serieconf, perfdatas)

        serieconf['last_computation'] = int(time())
        self[Serie.SERIE_STORAGE].put_element(element=serieconf)

        return point

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
            self.points_per_interval
        )

        return storage.find_elements(query={'$where': javascript_condition})
