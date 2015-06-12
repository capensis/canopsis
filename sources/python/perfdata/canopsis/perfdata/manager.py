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

from time import time

from canopsis.monitoring.parser import PerfDataParser
from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths
)
from canopsis.timeserie.timewindow import Period, get_offset_timewindow
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.context.manager import Context
from canopsis.storage.timed import TimedStorage

DEFAULT_PERIOD = Period(**{Period.WEEK: 1})  # save data each week

CONF_PATH = 'perfdata/perfdata.conf'
CATEGORY = 'PERFDATA'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class PerfData(MiddlewareRegistry):
    """Dedicated to access to perfdata (via periodic and timed stores).
    """

    PERFDATA_STORAGE = 'perfdata_storage'
    META_STORAGE = 'meta_storage'

    META_TIMESTAMP = TimedStorage.TIMESTAMP
    META_VALUE = TimedStorage.VALUE
    META_ID = TimedStorage.DATA_ID

    def __init__(
        self, perfdata_storage=None, meta_storage=None, *args, **kwargs
    ):

        super(PerfData, self).__init__(*args, **kwargs)

        self.context = Context()

        if perfdata_storage is not None:
            self[PerfData.PERFDATA_STORAGE] = perfdata_storage
        if meta_storage is not None:
            self[PerfData.META_STORAGE] = meta_storage

    def count(self, metric_id, timewindow=None, period=None):
        """Get number of perfdata identified by metric_id in input timewindow

        :param timewindow: if None, get all perfdata values
        """

        period = self.get_period(metric_id, period=period)

        result = self[PerfData.PERFDATA_STORAGE].count(
            data_id=metric_id, timewindow=timewindow, period=period
        )

        return result

    def get_metrics(self, query=None):
        """Get registered metric ids.

        :return: list of registered metric ids.
        :rtype: list
        """

        result = set()

        documents = self[PerfData.PERFDATA_STORAGE].find_elements(query=query)

        for document in documents:
            result.add(document['i'])

        return list(result)

    def get(
        self, metric_id, period=None, with_meta=True, timewindow=None,
        limit=0, skip=0, timeserie=None
    ):
        """Get a set of data related to input data_id on the timewindow and
        input period.

        If with_meta, result is a couple of (points, list of meta by timestamp)
        """

        period = self.get_period(metric_id, period=period)

        result = self[PerfData.PERFDATA_STORAGE].get(
            data_id=metric_id, timewindow=timewindow, period=period,
            limit=limit, skip=skip, timeserie=timeserie
        )

        if with_meta:

            meta = self[PerfData.META_STORAGE].get(
                data_ids=metric_id, timewindow=timewindow
            )

            result = result, meta

        return result

    def get_point(
        self, metric_id, period=None, with_meta=True, timestamp=None
    ):
        """Get the closest point before input timestamp. Add meta informations
        if with_meta.
        """

        if timestamp is None:
            timestamp = time()

        timewindow = get_offset_timewindow(timestamp)

        period = self.get_period(metric_id, period=period)

        result = self[PerfData.PERFDATA_STORAGE].get(
            data_id=metric_id, timewindow=timewindow, period=period,
            limit=1
        )

        if with_meta is not None:

            meta = self[PerfData.META_STORAGE].get(
                data_id=metric_id, timewindow=timewindow
            )

            result = result, meta

        return result

    def get_meta(
        self, metric_id, timewindow=None, limit=0, sort=None
    ):
        """Get the meta data related to input data_id and timewindow.
        """

        if timewindow is None:
            timewindow = get_offset_timewindow()

        result = self[PerfData.META_STORAGE].get(
            data_id=metric_id,
            timewindow=timewindow,
            limit=limit,
            sort=sort
        )

        return result

    def put(self, metric_id, points, meta=None, period=None, cache=False):
        """Put a (list of) couple (timestamp, value), a meta into
        rated_documents related to input period.

        kwargs will be added to all document in order to extend
        periodic_documents.

        :param iterable points: points to put. One point (timestamp, value) or
            points (timestamp, values)+.
        """

        # do something only if there are points to put
        if points:
            first_point = points[0]
            # if first_point is a timestamp, points is one point
            if isinstance(first_point, (int, float, str)):
                # transform points into a tuple
                points = (points,)

            period = self.get_period(metric_id=metric_id, period=period)
            # update data in a cache (a)synchronous way
            self[PerfData.PERFDATA_STORAGE].put(
                data_id=metric_id, period=period, points=points, cache=cache
            )

            if meta is not None:

                min_timestamp = min(point[0] for point in points)
                # update meta data in a synchronous way
                self[PerfData.META_STORAGE].put(
                    data_id=metric_id,
                    value=meta,
                    timestamp=min_timestamp
                )

    def remove(
        self,
        metric_id, period=None, with_meta=False, timewindow=None, cache=False
    ):
        """Remove values and meta of one metric.

        meta_names is a list of meta_data to remove. An empty list ensure that
        no meta data is removed.
        if meta_names is None, then all meta_names are removed.
        """

        period = self.get_period(metric_id, period=period)

        self[PerfData.PERFDATA_STORAGE].remove(
            data_id=metric_id,
            timewindow=timewindow,
            period=period,
            cache=cache
        )

        if with_meta:
            self[PerfData.META_STORAGE].remove(
                data_ids=metric_id, timewindow=timewindow, cache=cache
            )

    def put_meta(self, metric_id, meta, timestamp=None, cache=False):
        """Update meta information.
        """

        self[PerfData.PERFDATA_STORAGE].put(
            data_id=metric_id, value=meta, timestamp=timestamp, cache=cache)

    def remove_meta(self, metric_id, timewindow=None, cache=False):
        """Remove meta information.
        """

        self[PerfData.PERFDATA_STORAGE].remove(
            data_id=metric_id, timewindow=timewindow, cache=cache)

    def get_period(self, metric_id, period=None):
        """Get default period related to input metric_id.

        DEFAULT_PERIOD if related entity does not exist or does not contain
        a default period.
        """

        result = period

        if result is None:

            result = DEFAULT_PERIOD
            # TODO: restore when the period will be specified by entity
            # entity = self.context.get_entities(ids=metric_id)

            # if entity is not None and 'period' in entity:
            #     result = Period(**entity['period'])

        return result

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
