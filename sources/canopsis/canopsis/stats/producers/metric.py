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

from hashlib import sha1

from canopsis.configuration.model import Parameter
from canopsis.influxdb.core import InfluxDBStorage
from canopsis.common.middleware import Middleware
from canopsis.old.mfilter import check
from canopsis.timeserie.core import DEFAULT_ROUND_TIME, DEFAULT_PERIOD

CONF_PATH = 'stats/producers/metric.conf'
CATEGORY = 'METRIC_PRODUCER'


class MetricProducer(Middleware):
    """
    Base Metric producer.

    This object is used to generate events containing statistics metrics.
    """

    FILTER_STORAGE = 'filter_storage'
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
            value = DEFAULT_PERIOD.total_seconds()

        self._default_aggregation_interval = value

    @property
    def round_time_interval(self):
        if not hasattr(self, '_round_time_interval'):
            self.round_time_interval = None

        return self._round_time_interval

    @round_time_interval.setter
    def round_time_interval(self, value):
        if value is None:
            value = DEFAULT_ROUND_TIME

        self._round_time_interval = value

    def __init__(
        self,
        default_aggregation_interval=None,
        round_time_interval=None,
        filter_storage=None,
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

        if context is not None:
            self[MetricProducer.CONTEXT_MANAGER] = context

        if perfdata is not None:
            self[MetricProducer.PERFDATA_MANAGER] = perfdata

        self.idbs = InfluxDBStorage()

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

    def _count(
        self,
        name,
        author='__canopsis__',
        extra_fields={}
    ):
        tags = extra_fields
        tags['component'] = author

        point = {
            'measurement': name,
            'tags': tags,
            'fields': {'value': 1}
        }

        self.idbs._conn.write_points([point])

    def _delay(
        self,
        name,
        value,
        author='__canopsis__', extra_fields={}
    ):
        tags = extra_fields
        tags['component'] = author

        point = {
            'measurement': name,
            'tags': tags,
            'fields': {'value': value}
        }

        self.idbs._conn.write_points([point])
