#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.configuration import Parameter
from canopsis.timeserie.timewindow import Period, get_offset_timewindow
from canopsis.storage.manager import Manager
from canopsis.context.manager import Context

from collections import Iterable

from time import time

DEFAULT_PERIOD = Period(**{Period.HOUR: 24})
DEFAULT_AGGREGATION = 'MEAN'


class PerfData(Manager):
    """
    Dedicated to access to perfdata (via periodic and timed stores).
    """

    CONF_FILE = '~/etc/perfdata.conf'

    CATEGORY = 'PERFDATA'

    CONTEXT = 'context'
    PERFDATA_STORAGE = 'perfdata_storage'

    DATA_TYPE = 'metric'

    def __init__(
        self,
        context=None,
        data_type=DATA_TYPE, perfdata_storage=None, meta_storage=None,
        *args, **kwargs
    ):

        super(PerfData, self).__init__(*args, **kwargs)

        self.data_type = data_type
        self.perfdata_storage = perfdata_storage
        self.meta_storage = meta_storage

        self.context = Context() if context is None else context

    def count(self, metric_id, period=None, timewindow=None):

        period = self.get_period(data_id=metric_id, period=period)

        result = self.perfdata_storage.count(
            data_id=metric_id, period=period, timewindow=timewindow)

        return result

    def get(
        self, metric_id, with_meta=True, period=None, timewindow=None,
        limit=0, skip=0, *args, **kwargs
    ):
        """
        Get a set of data related to input data_id on the timewindow \
        and input period.
        If with_meta, result is a couple of (points, list of meta by timestamp)
        """

        period = self.get_period(data_id=metric_id, period=period)

        result = self.periodic_store.get(
            data_id=metric_id, period=period, timewindow=timewindow,
            limit=limit, skip=skip)

        if with_meta is not None:

            meta = self.meta_storage.get(
                data_id=metric_id, timewindow=timewindow)

            result = result, meta

        return result

    def get_point(
        self, metric_id, with_meta=True, period=None, timestamp=time(),
        *args, **kwargs
    ):
        """
        Get the closest point before input timestamp. Add meta informations \
        if with_meta.
        """

        period = self.get_period(data_id=metric_id, period=period)

        timewindow = get_offset_timewindow(timestamp)

        result = self.perfdata_storage.get(
            data_id=metric_id, period=period, timewindow=timewindow,
            limit=1)

        if with_meta is not None:

            meta = self.meta_storage.get(
                data_id=metric_id, timewindow=timewindow)

            result = result, meta

        return result

    def get_meta(
        self, metric_id, timewindow=None, limit=0, sort=None, *args, **kwargs
    ):
        """
        Get the meta data related to input data_id and timewindow.
        """

        if timewindow is None:
            timewindow = get_offset_timewindow()

        result = self.meta_storage.get(
            data_id=metric_id, timewindow=timewindow, limit=limit, sort=sort)

        return result

    def put(
        self, metric_id, points_or_point, meta=None, period=None,
        *args, **kwargs
    ):
        """
        Put a (list of) couple (timestamp, value), a meta into rated_documents
        related to input period.
        kwargs will be added to all document in order to extend
        periodic_documents
        """

        # if points_or_point is a point, transform it into a tuple of couple
        if len(points_or_point) > 0:
            if not isinstance(points_or_point[0], Iterable):
                points_or_point = (points_or_point,)

        period = self.get_period(data_id=metric_id, period=period)

        self.perfdata_storage.put(
            data_id=metric_id, period=period, points=points_or_point)

        if meta is not None:

            min_timestamp = min([point[0] for point in points_or_point])

            self.meta_storage.put(
                data_id=metric_id, value=meta, timestamp=min_timestamp)

    def remove(
        self, metric_id, with_meta=False, period=None, timewindow=None,
        *args, **kwargs
    ):
        """
        Remove values and meta of one metric.
        meta_names is a list of meta_data to remove. An empty list ensure that
        no meta data is removed.
        if meta_names is None, then all meta_names are removed.
        """

        aggregation, period = self.get_period(
            metric_id=metric_id, period=period)

        self.perfdata_storage.remove(
            data_id=metric_id, period=period, timewindow=timewindow)

        if with_meta:
            self.perfdata_storage.remove(
                data_id=metric_id, timewindow=timewindow)

    def update_meta(self, metric_id, meta, timestamp=None, *args, **kwargs):
        """
        Update meta information.
        """

        self.perfdata_storage.put(
            data_id=metric_id, value=meta, timestamp=timestamp)

    def remove_meta(self, metric_id, timewindow=None, *args, **kwargs):
        """
        Remove meta information.
        """

        self.perfdata_storage.remove(data_id=metric_id, timewindow=timewindow)

    def get_period(self, metric_id, aggregation=None, period=None):
        """
        Get default period related to input metric_id.
        DEFAULT_PERIOD if related entity does not exist or does not contain
        a default period.
        """

        result = period

        if result is None:

            result = DEFAULT_PERIOD

            entity = self.context.get(
                element_id=metric_id, element_type=PerfData.DATA_TYPE)

            result = DEFAULT_PERIOD if entity is None else entity.get(
                'period', DEFAULT_PERIOD)

        return result

    def _conf(self, *args, **kwargs):

        result = super(PerfData, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=PerfData.CATEGORY,
            new_content=(
                Parameter(PerfData.CONTEXT),
                Parameter(PerfData.PERFDATA_STORAGE)))

        return result

    def _get_conf_files(self, *args, **kwargs):

        result = super(PerfData, self)._get_conf_files(*args, **kwargs)

        result.append(PerfData.CONF_FILE)

        return result
