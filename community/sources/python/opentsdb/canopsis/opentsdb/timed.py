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


class InfluxDBTimedStorage(InfluxDBStorage, TimedStorage):
    """InfluxDB storage dedicated to manage timed data."""

    __register__ = True  #: register this class to middleware.

    def count(self, data_id, timewindow=None, *args, **kwargs):

        ids = 'COUNT({0})'.format(data_id)

        query = self._timewindowtowhere(timewindow=timewindow)

        result = self.get_elements(ids=ids, query=query)[data_id]

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
            result = {
                'time': {'$gte': timewindow.start()},
                'time': {'$lte': timewindow.stop()}
            }

        else:
            result = None

        return result

    def get(self, data_id, timewindow=None, limit=0, *args, **kwargs):

        query = self._timewindowtowhere(timewindow=timewindow)

        result = self.get_elements(ids=data_id, query=query, limit=limit)[data_id]

        return result

    def put(self, data_id, points, cache=False, *args, **kwargs):

        pointstoput = []

        for point in points:
            pointstoput.append(
                {
                    'measurement': data_id,
                    'time': point[0],
                    'fields': {'value': point[1]}
                }
            )

        return self._conn.write_points(
            points=pointstoput, time_precision='s',
            batch_size=self.cache_size if cache else 0
        )

    def remove(self, data_id, timewindow=None, **_):

        if timewindow is not None:
            raise ValueError(
                'This storage can not delete points in a specific timewindow'
            )

        return self._conn.delete_series(measurement=data_id)
