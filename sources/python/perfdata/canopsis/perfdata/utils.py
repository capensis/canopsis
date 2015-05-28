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

from canopsis.perfdata.manager import PerfData
from canopsis.timeserie.timewindow import TimeWindow, Period
from canopsis.timeserie import TimeSerie


class PerfDataInterface(object):
    """
    Implement common interactions with perfdata.
    """

    def __init__(self, manager=None, *args, **kwargs):
        super(PerfDataInterface, self).__init__(*args, **kwargs)

        self.manager = manager if manager is not None else PerfData()

    def count(self, metric_id, timewindow=None):
        if timewindow is not None:
            timewindow = TimeWindow(**timewindow)

        return self.manager.count(
            metric_id=metric_id, timewindow=timewindow
        )

    def get(
        self, metric_id, timewindow=None, period=None, with_meta=True,
        limit=0, skip=0, timeserie=None
    ):
        if timewindow is not None:
            timewindow = TimeWindow(**timewindow)

        if timeserie is not None:
            if period is None:
                period = timeserie.pop('period', None)

            timeserie = TimeSerie(**timeserie)

            if period is not None:
                timeserie.period = Period(**period)

        if not isinstance(metric_id, list):
            metrics = [metric_id]

        else:
            metrics = metric_id

        result = []

        for metric_id in metrics:
            pts, meta = self.manager.get(
                metric_id=metric_id, with_meta=True,
                timewindow=timewindow, limit=limit, skip=skip
            )

            meta = meta[0]

            if timeserie is not None:
                pts = timeserie.calculate(pts, timewindow, meta=meta)

            if with_meta:
                result.append({
                    "points": pts,
                    "meta": meta
                })

            else:
                result.append({
                    "points": pts
                })

        return (result, len(result))

    def meta(self, metric_id, timewindow=None, limit=0, sort=None):
        if timewindow is not None:
            timewindow = TimeWindow(**timewindow)

        return self.manager.get_meta(
            metric_id=metric_id, timewindow=timewindow, limit=limit, sort=sort
        )

    def put(self, metric_id, points, meta=None):
        self.manager.put(
            metric_id=metric_id, points=points, meta=meta
        )

        return points

    def remove(self, metric_id, with_meta=False, timewindow=None):
        if timewindow is not None:
            timewindow = TimeWindow(**timewindow)

        self.manager.remove(
            metric_id=metric_id, with_meta=with_meta,
            timewindow=timewindow
        )
