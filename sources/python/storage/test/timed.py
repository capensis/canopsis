#!/usr/bin/env python
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

from unittest import main

from canopsis.timeserie.timewindow import TimeWindow

from .base import BaseTestConfiguration, BaseStorageTest


@conf_paths('storage/test-timed.conf')
class TestConfiguration(BaseTestConfiguration):
    """Default test configuration."""


class PeriodicStoreTest(BaseStorageTest):

    def _testconfcls(self):

        return TestConfiguration

    def _test(self, storage):

        self._test_CRUD(storage)

    def _test_CRUD(self, storage):

        # let's play with different data_names
        data_ids = ['m0', 'm1']

        timewindow = TimeWindow()

        points = [
            (timewindow.start(), None),  # lower bound
            (timewindow.stop(), 0),  # upper bound
            (timewindow.start() - 1, 1),  # outside timewindow (minus 1)
            (timewindow.start() + 1, 2),  # inside timewindow (plus 1)
            (timewindow.stop() + 1, 3)  # outside timewindow (plus 1)
        ]

        sorted_points = sorted(points, key=lambda point: point[0])

        inserted_points = {}

        # starts to put points for every aggregations and periods
        for data_id in data_ids:

            inserted_points[data_id] = points
            # add documents
            storage.put(data_id=data_id, points=points)

        points_count_in_timewindow = len(
            [point for point in points if point[0] in timewindow])

        # check for reading methods
        for data_id in data_ids:
            # iterate on data_ids
            count = storage.count(data_id=data_id)
            self.assertEquals(count, len(points))

            count = storage.count(data_id=data_id, timewindow=timewindow)
            self.assertEquals(count, points_count_in_timewindow)

            data = storage.get(data_id=data_id)
            self.assertEquals(len(data), len(points))
            self.assertEquals(data, sorted_points)

            data = storage.get(data_id=data_id, timewindow=timewindow)
            self.assertEquals(len(data), points_count_in_timewindow)
            self.assertEquals(data,
                [point for point in sorted_points
                if point[0] in timewindow])

            storage.remove(data_id=data_id, timewindow=timewindow)

            # check for count equals 1
            count = storage.count(data_id=data_id, timewindow=timewindow)
            self.assertEquals(count, 0)
            count = storage.count(data_id=data_id)
            self.assertEquals(
                count, len(points) - points_count_in_timewindow)

            storage.remove(data_id=data_id)
            # check for count equals 0
            count = storage.count(data_id=data_id)
            self.assertEquals(count, 0)


if __name__ == '__main__':
    main()
