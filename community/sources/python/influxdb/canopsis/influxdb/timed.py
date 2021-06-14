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

from .core import InfluxDBStorage

from canopsis.storage.timed import TimedStorage

from sys import getsizeof

from dateutil.parser import parse

from calendar import timegm

from numbers import Number


class InfluxDBTimedStorage(InfluxDBStorage, TimedStorage):
    """InfluxDB storage dedicated to manage timed data."""

    __register__ = True  #: register this class to middleware.

    def count(self, data_id, timewindow=None, tags=None, *args, **kwargs):

        result = 0

        query = self._timewindowtowhere(timewindow=timewindow)

        points = self.get_elements(
            projection='COUNT(value)', ids=data_id, query=query
        )

        if points:
            result = next(points)['count']

        return result

    def size(self, data_id=None, timewindow=None, *args, **kwargs):

        return (
            getsizeof(0) *
            self.count(data_id=data_id, timewindow=timewindow, *args, **kwargs)
        )

    @staticmethod
    def _timewindowtowhere(timewindow):
        """Transform a timewindow into a WHERE query."""

        if timewindow is not None:
            factor = 1e9
            result = {
                'time': {
                    '$gte': int(timewindow.start() * factor),
                    '$lte': int(timewindow.stop() * factor)
                }
            }

        else:
            result = None

        return result

    def get(
        self,
        data_id,
        timewindow=None,
        limit=0,
        with_tags=False,
        tags=None,
        **_
    ):

        query = self._timewindowtowhere(timewindow=timewindow)

        points = self.get_elements(
            projection=None if with_tags else 'value',
            ids=data_id,
            query=query,
            limit=limit,
            tags=tags
        )

        _points, tags = [], {}

        if points:
            for point in points:
                timestamp = timegm(parse(point['time']).timetuple())
                _points.append((timestamp, point.pop('value')))
                if with_tags:
                    tags.update(point)
                    tags['timestamp'] = point['time']

        result = (_points, tags) if with_tags else _points

        return result

    def put(self, data_id, points, tags=None, cache=False, *args, **kwargs):

        pointstoput = []

        factor = 1e9

        for point in points:
            value = point[1]
            if isinstance(value, Number):
                value = float(value)
            pointstoput.append(
                {
                    'measurement': data_id,
                    'time': int(point[0] * factor),
                    'fields': {'value': value}
                }
            )

        return self.put_elements(elements=pointstoput, cache=cache, tags=tags)

    def remove(self, data_id, timewindow=None, tags=None, **_):

        if timewindow is not None:
            raise ValueError(
                'This storage can not delete points in a specific timewindow'
            )

        return self.remove_elements(ids=data_id, tags=tags)
