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

from re import compile as re_compile

from time import time

from canopsis.common.utils import force_iterable
from canopsis.configuration import add_category, conf_paths
from canopsis.timeserie.timewindow import Period, get_offset_timewindow
from canopsis.middleware.manager import Manager
from canopsis.context.manager import Context

DEFAULT_PERIOD = Period(**{Period.WEEK: 1})  # save data each week

CONF_PATH = 'perfdata/perfdata.conf'
CATEGORY = 'PERFDATA'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class PerfData(Manager):
    """
    Dedicated to access to perfdata (via periodic and timed stores).
    """

    PERFDATA_STORAGE = 'perfdata_storage'
    META_STORAGE = 'meta_storage'

    DATA_SCOPE = 'metric'

    def __init__(
        self,
        data_scope=DATA_SCOPE, perfdata_storage=None, meta_storage=None,
        *args, **kwargs
    ):

        super(PerfData, self).__init__(
            data_scope=data_scope, *args, **kwargs)

        self.context = Context()

        if perfdata_storage is not None:
            self[PerfData.PERFDATA_STORAGE] = perfdata_storage
        if meta_storage is not None:
            self[PerfData.META_STORAGE] = meta_storage

    def count(self, metric_id, timewindow=None, period=None):
        """
        Get number of perfdata identified by metric_id in input timewindow

        :param timewindow: if None, get all perfdata values
        """

        period = self.get_period(metric_id, period=period)

        result = self[PerfData.PERFDATA_STORAGE].count(
            data_id=metric_id, timewindow=timewindow, period=period)

        return result

    def get(
        self, metric_id, period=None, with_meta=True, timewindow=None,
        limit=0, skip=0
    ):
        """
        Get a set of data related to input data_id on the timewindow \
        and input period.
        If with_meta, result is a couple of (points, list of meta by timestamp)
        """

        period = self.get_period(metric_id, period=period)

        result = self[PerfData.PERFDATA_STORAGE].get(
            data_id=metric_id, timewindow=timewindow, period=period,
            limit=limit, skip=skip)

        if with_meta is not None:

            meta = self[PerfData.META_STORAGE].get(
                data_ids=metric_id, timewindow=timewindow)

            result = result, meta

        return result

    def get_point(
        self, metric_id, period=None, with_meta=True, timestamp=None
    ):
        """
        Get the closest point before input timestamp. Add meta informations \
        if with_meta.
        """

        if timestamp is None:
            timestamp = time()

        timewindow = get_offset_timewindow(timestamp)

        period = self.get_period(metric_id, period=period)

        result = self[PerfData.PERFDATA_STORAGE].get(
            data_id=metric_id, timewindow=timewindow, period=period,
            limit=1)

        if with_meta is not None:

            meta = self[PerfData.META_STORAGE].get(
                data_id=metric_id, timewindow=timewindow)

            result = result, meta

        return result

    def get_meta(self, metric_id, timewindow=None, limit=0, sort=None):
        """
        Get the meta data related to input data_id and timewindow.
        """

        if timewindow is None:
            timewindow = get_offset_timewindow()

        result = self[PerfData.META_STORAGE].get(
            data_id=metric_id, timewindow=timewindow, limit=limit, sort=sort)

        return result

    def put(self, metric_id, points, meta=None, period=None):
        """
        Put a (list of) couple (timestamp, value), a meta into rated_documents
        related to input period.
        kwargs will be added to all document in order to extend
        periodic_documents
        """

        # if points is a point, transform it into a tuple of couple
        points = force_iterable(points, iterable=tuple)

        period = self.get_period(metric_id=metric_id, period=period)

        self[PerfData.PERFDATA_STORAGE].put(
            data_id=metric_id, period=period, points=points)

        if meta is not None:

            min_timestamp = min(point[0] for point in points)

            self[PerfData.META_STORAGE].put(
                data_id=metric_id, value=meta, timestamp=min_timestamp)

    def remove(self, metric_id, period=None, with_meta=False, timewindow=None):
        """
        Remove values and meta of one metric.
        meta_names is a list of meta_data to remove. An empty list ensure that
        no meta data is removed.
        if meta_names is None, then all meta_names are removed.
        """

        period = self.get_period(metric_id, period=period)

        self[PerfData.PERFDATA_STORAGE].remove(
            data_id=metric_id, timewindow=timewindow, period=period)

        if with_meta:
            self[PerfData.META_STORAGE].remove(
                data_ids=metric_id, timewindow=timewindow)

    def update_meta(self, metric_id, meta, timestamp=None,):
        """
        Update meta information.
        """

        self[PerfData.PERFDATA_STORAGE].put(
            data_id=metric_id, value=meta, timestamp=timestamp)

    def remove_meta(self, metric_id, timewindow=None):
        """
        Remove meta information.
        """

        self[PerfData.PERFDATA_STORAGE].remove(
            data_id=metric_id, timewindow=timewindow)

    def get_period(self, metric_id, period=None):
        """
        Get default period related to input metric_id.
        DEFAULT_PERIOD if related entity does not exist or does not contain
        a default period.
        """

        result = period

        if result is None:

            result = DEFAULT_PERIOD

            entity = self.context.get_entities(ids=metric_id)

            if entity is not None and 'period' in entity:
                result = Period(**entity['period'])

        return result

    RE_PERF_DATA = re_compile(
        "('?([0-9A-Za-z/\\\:\.%%\-{}\?\[\]_ ]*)'?=(\-?[0-9.,]*)(([A-Za-z%%/]*))(;@?(\-?[0-9.,]*):?)?(;@?(\-?[0-9.,]*):?)?(;@?(\-?[0-9.,]*):?)?(;@?(\-?[0-9.,]*):?)?(;? ?))")

    def parse_perfdata(self, perf_data_raw):
        # 'label'=value[UOM];[warn];[crit];[min];[max]
        #   load1=0.440     ;5.000 ;10.000;0    ;

        self.logger.debug("Parse: %s" % perf_data_raw)

        perfs = PerfData.RE_PERF_DATA.split(perf_data_raw)

        perf_data_array = []
        perf_data = {}
        i = 0
        for info in perfs:
            if info == '':
                info = None

            #self.logger.debug(" + %s '%s'" % (i, info))
            try:
                if info and i == 2:
                    perf_data['metric'] = info
                elif info and i == 3:
                    perf_data['value'] = info.replace(',', '.')
                elif info and i == 4:
                    perf_data['unit'] = info
                elif info and i == 7:
                    perf_data['warn'] = info.replace(',', '.')
                elif info and i == 9:
                    perf_data['crit'] = info.replace(',', '.')
                elif info and i == 11:
                    perf_data['min'] = info.replace(',', '.')
                elif info and i == 13:
                    perf_data['max'] = info.replace(',', '.')

                i += 1
                if i is 15:
                    try:
                        perf_data_clean = {}
                        for key in perf_data.keys():
                            if perf_data[key]:
                                try:
                                    perf_data_clean[key] = float(
                                        perf_data[key])
                                except:
                                    if key == 'metric' or key == 'unit':
                                        perf_data_clean[key] = perf_data[key]
                                    else:
                                        self.logger.debug(
                                            "Invalid value, '%s' = '%s'" % (
                                                key, perf_data[key]))

                        if 'value' in perf_data_clean \
                                and 'metric' in perf_data_clean:
                            perf_data_array.append(perf_data_clean)

                        if not perf_data_clean.get('unit', None):
                            # split: g[in_bps]= ...
                            metric_ori = perf_data_clean['metric']
                            if metric_ori[len(metric_ori) - 1] == ']':
                                metric_ori = metric_ori[:len(metric_ori) - 1]
                                metric = metric_ori.split('[', 1)
                                if len(metric) == 2:
                                    perf_data_clean['metric'] = metric[0]
                                    perf_data_clean['unit'] = metric[1]

                        self.logger.debug(" + %s" % perf_data_clean)

                    except Exception as err:

                        self.logger.error("perf_data: Raw: %s" % perf_data_raw)
                        self.logger.error(
                            "perf_data: Impossible to clean '%s': %s" % (
                                perf_data, err))

                    perf_data = {}
                    i = 0

            except Exception as err:
                self.logger.error(
                    "perf_data: Invalid metric %s: %s (%s)" % (i, info, err))

        return perf_data_array

    def is_internal(self, metric):

        return metric['metric'].startswith('cps_')
