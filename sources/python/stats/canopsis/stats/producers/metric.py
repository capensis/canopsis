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

from canopsis.timeserie.core import TimeSerie
from canopsis.old.mfilter import check

from hashlib import sha1


CONF_PATH = 'stats/producers/metric.conf'
CATEGORY = 'METRIC_PRODUCER'
CONTENT = [
    Parameter('default_aggregation_interval', int),
    Parameter('round_time_interval', Parameter.bool)
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class MetricProducer(MiddlewareRegistry):
    """
    Base Metric producer.

    This object is used to generate events containing statistics metrics.
    """

    FILTER_STORAGE = 'filter_storage'
    SERIE_STORAGE = 'serie_storage'
    CONTEXT_MANAGER = 'context'
    PERFDATA_MANAGER = 'perfdata'

    @property
    def default_aggregation_interval(self):
        if not hasattr(self, '_default_aggregation_interval'):
            self.default_aggregation_interval = None

        return self._default_aggregation_interval

    @default_aggregation_interval.setter
    def default_aggregation_interval(self, value):
        if value is None:
            value = TimeSerie.VPERIOD.total_seconds()

        self._default_aggregation_interval = value

    @property
    def round_time_interval(self):
        if not hasattr(self, '_round_time_interval'):
            self.round_time_interval = None

        return self._round_time_interval

    @round_time_interval.setter
    def round_time_interval(self, value):
        if value is None:
            value = TimeSerie.VROUND_TIME

        self._round_time_interval = value

    def __init__(
        self,
        default_aggregation_interval=None,
        round_time_interval=None,
        filter_storage=None,
        serie_storage=None,
        context=None,
        perfdata=None,
        *args, **kwargs
    ):
        super(MetricProducer, self).__init__(*args, **kwargs)

        if default_aggregation_interval is not None:
            self.default_aggregation_interval = default_aggregation_interval

        if round_time_interval is not None:
            self.round_time_interval = round_time_interval

        if filter_storage is not None:
            self[MetricProducer.FILTER_STORAGE] = filter_storage

        if serie_storage is not None:
            self[MetricProducer.SERIE_STORAGE] = serie_storage

        if context is not None:
            self[MetricProducer.CONTEXT_MANAGER] = context

        if perfdata is not None:
            self[MetricProducer.PERFDATA_MANAGER] = perfdata

    def match(self, event):
        """
        Get filters names which match the event.

        :param event: Event to check
        :type event: dict

        :returns: filters names as list
        """

        storage = self[MetricProducer.FILTER_STORAGE]
        matches = [
            evfilter['crecord_name']
            for evfilter in storage.find_elements()
            if check(evfilter.get('filter', None) or {}, event)
        ]

        return matches

    def get_stats_serie_id(self, metric_id, operator):
        """
        Get serie name from metric id and operator.

        :returns: sha1(metric_id, operator) as str
        """

        serie_id = sha1()
        serie_id.update(metric_id)
        serie_id.update(operator)
        return serie_id.hexdigest()

    def may_create_stats_serie(self, metric, operator):
        """Create serie for metric and operator if not existing yet.

        :param metric: Metric entity
        :type metric: dict

        :param operator: Aggregation operator used on metric
        :type operator: str

        :returns: the newly created, or existing, serie
        """

        storage = self[MetricProducer.SERIE_STORAGE]

        metric_id = self[MetricProducer.CONTEXT_MANAGER].get_entity_id(metric)
        serie_id = self.get_stats_serie_id(metric_id, operator)

        result = storage.get_elements(ids=serie_id)

        if result is None:
            serie = {
                'crecord_name': operator,
                'component': metric['component'],
                'resource': metric['resource'],
                'metric_filter': 'co:{0} re:{1} me:{2}'.format(
                    metric['component'],
                    metric['resource'],
                    metric['name']
                ),
                'aggregation_method': operator,
                'aggregation_interval': self.default_aggregation_interval,
                'round_time_interval': self.round_time_interval,
                # only one metric selected, so SUM is the identity
                'formula': 'SUM("me:.*")',
                'last_computation': 0
            }

            _, tags = self[MetricProducer.PERFDATA_MANAGER].get(
                metric_id, with_tags=True
            )

            if tags is not None:
                serie.update(tags)

            storage.put_element(serie, _id=serie_id)

            serie[storage.DATA_ID] = serie_id
            result = serie

        return result

    def _counter(self, name, event, author='__canopsis__'):
        """
        Generate counters for each matching filter.

        :param name: Counter's name
        :type name: str

        :param event: Event to check against filters
        :type event: dict

        :param author: Statistic author (default: __canopsis__)
        :type author: str

        :returns: perf event as dict
        """

        event = {
            'connector': 'canopsis',
            'connector_name': 'stats',
            'event_type': 'perf',
            'source_type': 'resource',
            'component': author,
            'resource': name,
            'perf_data_array': [
                {
                    'metric': filtername,
                    'value': 1,
                    'type': 'COUNTER'
                }
                for filtername in self.match(event)
            ]
        }

        return event

    def _delay(self, name, value, author='__canopsis__'):
        """
        Generate gauge and counter for delay statistic.

        :param name: Delay's name
        :type name: str

        :param value: Delay
        :type value: float

        :param author: Statistic author (default: __canopsis__)
        :type author: str

        :returns: perf event as dict
        """

        event = {
            'connector': 'canopsis',
            'connector_name': 'stats',
            'event_type': 'perf',
            'source_type': 'resource',
            'component': author,
            'resource': name,
            'perf_data_array': [
                {
                    'metric': 'sum',
                    'value': value,
                    'type': 'COUNTER'
                },
                {
                    'metric': 'last',
                    'value': value,
                    'type': 'GAUGE'
                }
            ]
        }

        for operator in ['min', 'max', 'average']:
            entity = self[MetricProducer.PERFDATA_MANAGER].get_metric_entity(
                'last', event
            )

            self.may_create_stats_serie(entity, operator)

        return event
