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

from canopsis.common.ws import route
from canopsis.perfdata.manager import PerfData
from canopsis.timeserie.timewindow import TimeWindow, Period
from canopsis.timeserie.core import TimeSerie

manager = PerfData()


def exports(ws):

    @route(ws.application.post, payload=['metric_id', 'timewindow'])
    def perfdata_count(metric_id, timewindow=None):
        if timewindow is not None:
            timewindow = TimeWindow(**timewindow)

        result = manager.count(
            metric_id=metric_id, timewindow=timewindow
        )

        return result

    @route(
        ws.application.post,
        payload=[
            'metric_id', 'with_meta',
            'limit', 'skip', 'period',
            'timewindow', 'period', 'timeserie'
        ])
    def perfdata(
        metric_id, timewindow=None, period=None, with_meta=True,
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
            # meta -> _meta
            pts, _meta = manager.get(
                metric_id=metric_id, with_meta=True,
                timewindow=timewindow, limit=limit, skip=skip
            )

            meta = _meta[0] if _meta is not None else {}
            meta[manager[PerfData.META_STORAGE].DATA_ID] = metric_id

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

    @route(ws.application.post, payload=['timewindow', 'limit', 'sort'])
    def perfdata_meta(metric_id, timewindow=None, limit=0, sort=None):
        if timewindow is not None:
            timewindow = TimeWindow(**timewindow)

        result = manager.get_meta(
            metric_id=metric_id, timewindow=timewindow, limit=limit, sort=sort
        )

        return result

    @route(ws.application.put, payload=[
        'metric_id', 'points', 'meta'
    ])
    def perfdata(metric_id, points, meta=None):
        manager.put(
            metric_id=metric_id, points=points, meta=meta
        )

        result = points

        return result

    @route(ws.application.delete, payload=[
        'metric_id', 'with_meta', 'timewindow'
    ])
    def perfdata(metric_id, with_meta=False, timewindow=None):
        if timewindow is not None:
            timewindow = TimeWindow(**timewindow)

        manager.remove(
            metric_id=metric_id, with_meta=with_meta,
            timewindow=timewindow
        )

        result = None

        return result

    @route(ws.application.put, payload=['metric_id', 'meta', 'timestamp'])
    def perfdata_meta(metric_id, meta, timestamp=None):
        result = manager.put_meta(
            metric_id=metric_id, meta=meta, timestamp=timestamp
        )

        return result

    @route(ws.application.get)
    def perfdata_period(metric_id):
        result = manager.get_period(metric_id)

        return result

    @route(ws.application.get)
    def perfdata_internal(metric):
        result = manager.is_internal(metric)

        return result
