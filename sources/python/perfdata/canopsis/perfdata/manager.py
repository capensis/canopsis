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

from __future__ import unicode_literals
from time import time

from canopsis.common.init import basestring
from canopsis.monitoring.parser import PerfDataParser
from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths
)
from canopsis.timeserie.timewindow import get_offset_timewindow, TimeWindow
from canopsis.middleware.registry import MiddlewareRegistry

from numbers import Number

CONF_PATH = 'perfdata/perfdata.conf'
CATEGORY = 'PERFDATA'

SLIDING_TIME = 'sliding_time'
SLIDING_TIME_UPPER_BOUND = 365 * 86400 * 3


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class PerfData(MiddlewareRegistry):
    """Dedicated to access to perfdata (via periodical and timed stores)."""

    PERFDATA_STORAGE = 'perfdata_storage'
    CONTEXT_MANAGER = 'context'

    @property
    def context(self):
        return self[PerfData.CONTEXT_MANAGER]

    def __init__(
            self, perfdata_storage=None, context=None,
            *args, **kwargs
    ):

        super(PerfData, self).__init__(*args, **kwargs)

        if perfdata_storage is not None:
            self[PerfData.PERFDATA_STORAGE] = perfdata_storage

        if context is not None:
            self[PerfData.CONTEXT_MANAGER] = context

    def _get_tags_metric_id(self, metric_id, meta=None, event={}):

        tags = {} if meta is None else meta.copy()

        # entity = self[PerfData.CONTEXT_MANAGER].get_entities_by_id(metric_id)[0]

        entity = {}

        eid = '/metric/{0}/{1}/{2}/{3}/{4}'.format(
            event['connector'],
            event['connector_name'],
            event['component'],
            event['resource'],
            event['perf_metric']
        )

        entity = {
            'connector':  event['connector'],
            'connector_name': event['connector_name'],
            'component': event['component'],
            'resource': event['resource'],
            'eid': eid,
            'type': 'metric',
            # 'retention': meta['retention'],
            # 'unit': meta['unit']

        }
        tags.update(entity)
        tags[eid] = metric_id

        data_id = tags.pop(eid)

        return data_id, tags

    def count(self, metric_id, timewindow=None, meta=None):
        """Get number of perfdata identified by metric_id in input timewindow

        :param timewindow: if None, get all perfdata values
        """

        data_id, tags = self._get_tags_metric_id(metric_id, meta)

        result = self[PerfData.PERFDATA_STORAGE].count(
            data_id=data_id, timewindow=timewindow, tags=tags
        )

        return result

    def get_metrics(self, query=None):
        """Get registered metric ids.

        :return: list of registered metric ids.
        :rtype: list
        """

        return self[PerfData.CONTEXT_MANAGER].find(
            _type='metric', _filter=query
        )

    def get(
            self, metric_id, meta=None, with_meta=True, timewindow=None,
            limit=0, skip=0, sort=None, timeserie=None, sliding_time=False
    ):
        """Get a set of data related to input data_id on the timewindow and
        input period.

        If with_meta, result is a couple of (points, list of tags by timestamp)
        """

        data_id, tags = self._get_tags_metric_id(metric_id, meta)

        if sliding_time:  # apply sliding_time

            if timewindow is None:

                timewindow = TimeWindow()

            _timewindow = timewindow

            timewindow = TimeWindow(
                start=timewindow.start(),
                stop=timewindow.stop() + SLIDING_TIME_UPPER_BOUND
            )

        result = self[PerfData.PERFDATA_STORAGE].get(
            data_id=data_id, timewindow=timewindow, limit=limit,
            skip=skip, timeserie=timeserie, tags=tags, with_tags=with_meta,
            sort=sort
        )

        if sliding_time:

            if with_meta:
                points = result[0]

            else:
                points = result

            points = [
                (min(ts, _timewindow.stop()), val)
                for (ts, val) in points
            ]

            if with_meta:
                result = points, result[1]

            else:
                result = points

        return result

    def get_point(self, metric_id, with_meta=True, timestamp=None, meta=None):
        """Get the closest point before input timestamp. Add tags informations
        if with_tags.
        """

        data_id, tags = self._get_tags_metric_id(metric_id, meta)

        if timestamp is None:
            timestamp = time()

        timewindow = get_offset_timewindow(timestamp)

        result = self[PerfData.PERFDATA_STORAGE].get(
            data_id=data_id, timewindow=timewindow,
            limit=1, tags=tags, with_tags=with_meta
        )

        return result

    def put(self, metric_id, points, meta=None, cache=False, event={}):
        """Put a (list of) couple (timestamp, value), a tags into
        rated_documents.

        kwargs will be added to all document in order to extend timed
        documents.

        :param iterable points: points to put. One point (timestamp, value) or
            points (timestamp, values)+.
        """

        # do something only if there are points to put
        if points:
            first_point = points[0]
            # if first_point is a timestamp, points is one point
            if isinstance(first_point, (Number, basestring)):
                # transform points into a tuple
                points = (points,)

            data_id, tags = self._get_tags_metric_id(metric_id, meta, event=event)

            # update data in a cache (a)synchronous way
            self[PerfData.PERFDATA_STORAGE].put(
                data_id=data_id, points=points, tags=tags, cache=cache
            )

    def remove(
            self, metric_id, timewindow=None, meta=None, cache=False
    ):
        """Remove values and tags of one metric."""

        data_id, tags = self._get_tags_metric_id(metric_id, meta)

        self[PerfData.PERFDATA_STORAGE].remove(
            data_id=data_id,
            timewindow=timewindow,
            cache=cache,
            tags=tags
        )

    def parse_perfdata(self, perf_data_raw):
        """Try to get a perf data array from input perf_data_raw.

        :param str perf_data_raw: perf_data_raw to parse.
        :return: array of perfdata.
        :rtype: list
        :raises: parsing error if perf_data_raw is not in an understood format.
        """

        self.logger.debug("Parse: {0}".format(perf_data_raw))

        parser = PerfDataParser(perf_data_raw)
        result = parser.perf_data_array

        return result

    def is_internal(self, metric):

        return metric['metric'].startswith('cps_')

    def get_metric_infos(self, limit, start):
        """
        Retreive metrics informations from influx.

        :param int limit: how many records to retreive
        :param int start: skip n first elements
        :rtype: list(dict)
        """
        data = self[PerfData.PERFDATA_STORAGE]._conn.query('SHOW SERIES;').raw

        if "series" not in data:
            return []

        metric_infos = []
        for serie in data['series']:
            try:
                index_ = serie['columns'].index('eid')
            except:
                self.logger.debug("Could not find eid in columns")
                continue

            for value in serie['values']:
                split = value[index_].split('/')
                dict_ = {
                    '_id': value[index_],
                    'type': 'metric',
                    'connector': split[2],
                    'connector_name': split[3],
                    'component': split[4],
                    'resource': split[5],
                    'name': split[6],
                    'internal': False  # TODO: calculate that
                }
                metric_infos.append(dict_)

        if limit == None:         # No limit
            end = len(metric_infos)
        else:
            end = start + limit

        return metric_infos[start:end]
