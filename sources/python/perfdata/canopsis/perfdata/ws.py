#!/usr/bin/env python
# --------------------------------
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

from bottle import get, delete, put, post

from canopsis.common.ws import response, route
from canopsis.perfdata.manager import PerfData

manager = PerfData()


@route(get)
def perfdata_count(metric_id, timewindow=None, period=None):

    count = manager.count(
        metric_id=metric_id, timewindow=timewindow, period=period)

    result = response(count)

    return result


@route(get)
def perfdata(
    metric_id, period=None, with_meta=True, timewindow=None,
    limit=0, skip=0, timeserie=None
):

    points = manager.get(
        metric_id=metric_id, period=period, with_meta=with_meta,
        timewindow=timewindow, limit=limit, skip=skip, timeserie=timeserie)

    if timeserie is not None:
        _points = points[0] if with_meta else points
        points = timeserie.calculate(_points, timewindow=timewindow)

    result = response(perfdata)

    return result


@route(get)
def perfdata_meta(metric_id, timewindow=None, limit=0, sort=None):

    meta = manager.get_meta(
        metric_id=metric_id, timewindow=timewindow, limit=limit, sort=sort)

    result = response(meta)

    return result


@route(put)
def perfdata(metric_id, points, meta=None, period=None):

    manager.put(metric_id=metric_id, points=points, meta=meta, period=period)

    result = response(points)

    return result


@route(delete)
def perfdata(metric_id, period=None, with_meta=False, timewindow=None):

    manager.remove(
        metric_id=metric_id, period=period, with_meta=with_meta,
        timewindow=timewindow)

    result = response(None)

    return result


@route(post)
def perfdata_meta(metric_id, meta, timestamp=None):

    nodes = manager.put_meta(
        metric_id=metric_id, meta=meta, timestamp=timestamp)

    result = response(nodes)

    return result


@route(get)
def perfdata_period(metric_id):

    period = manager.get_period(metric_id)

    result = response(period)

    return result


@route(get)
def perfdata_internal(metric):

    internal = manager.is_internal(metric)

    result = response(internal)

    return result
