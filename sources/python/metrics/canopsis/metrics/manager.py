# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

"""
Metrics manager.

Handle basic operations on series and metrics data.
"""

from __future__ import unicode_literals

#from canopsis.monitoring.parser import PerfDataParser
#from canopsis.timeserie.timewindow import get_offset_timewindow, TimeWindow

#SLIDING_TIME = 'sliding_time'
#SLIDING_TIME_UPPER_BOUND = 365 * 86400 * 3


class MetricsManager(object):
    """Access to metrics."""

    LOG_PATH = 'var/log/metrics.log'

    def __init__(self, logger, store):
        """
        :param store: an <InfluxStore> object
        """
        super(MetricsManager, self).__init__()
        self.logger = logger
        self.store = store

    def count(self, metric_id, timewindow=None, meta=None):
        """Get number of perfdata identified by metric_id in input timewindow

        :param timewindow: if None, get all perfdata values
        """
        pass

    def get_metrics(self, **_):
        return self.get_all_metrics()

    def get_all_metrics(self):
        """Get registered metric ids.

        :return: list of registered metric ids.
        :rtype: list
        """
        self.store.get_list_measurements()
        # TODO: upgrade python-influx to 4.1.1

    def get_metric_infos(self, limit, start, filter_):
        pass

    def get_period(self, metric_id):
        pass

    def is_internal(self, metric):
        pass

    def put(self, metric_id, points, meta):
        pass

    def remove(self, metric_id, timewindow, meta):
        pass
