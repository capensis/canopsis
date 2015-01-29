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

from unittest import TestCase, main

from canopsis.mongo.periodic import MongoPeriodicStorage
from canopsis.timeserie.timewindow import Period, TimeWindow


class PeriodicStoreTest(TestCase):

    def setUp(self):
        # create a storage on test_store collections
        self.storage = MongoPeriodicStorage(data_scope="test_store")

    def test_connect(self):
        self.assertTrue(self.storage.connected())

        self.storage.disconnect()

        self.assertFalse(self.storage.connected())

        self.storage.connect()

        self.assertTrue(self.storage.connected())

    def test_CRUD(self):
        # start in droping data
        self.storage.drop()

        # let's play with different data_names
        data_ids = ['m0', 'm1']
        periods = [
            Period(**{Period.MINUTE: 60}),
            Period(**{Period.HOUR: 24})
        ]

        timewindow = TimeWindow()

        points = [
            (timewindow.start(), None),  # lower bound
            (timewindow.stop(), 0),  # upper bound
            (timewindow.start() - 1, 1),  # outside timewindow (minus 1)
            (timewindow.start() + 1, 2),  # inside timewindow (plus 1)
            (timewindow.stop() + 1, 3)  # outside timewindow (plus 1)
        ]

        sorted_points = sorted(points, key=lambda point: point[0])

        inserted_points = dict()

        # starts to put points for every aggregations and periods
        for data_id in data_ids:
            inserted_points[data_id] = dict()
            for period in periods:
                inserted_points[data_id][period] = points
                # add documents
                self.storage.put(data_id=data_id, period=period, points=points)

        points_count_in_timewindow = len(
            [point for point in points if point[0] in timewindow])

        # check for reading methods
        for data_id in data_ids:
            # iterate on data_ids

            for period in periods:

                count = self.storage.count(data_id=data_id, period=period)
                self.assertEquals(count, len(points))

                count = self.storage.count(
                    data_id=data_id,
                    period=period, timewindow=timewindow)
                self.assertEquals(count, points_count_in_timewindow)

                data = self.storage.get(data_id=data_id, period=period)
                self.assertEquals(len(data), len(points))
                self.assertEquals(data, sorted_points)

                data = self.storage.get(
                    data_id=data_id,
                    period=period, timewindow=timewindow)
                self.assertEquals(len(data), points_count_in_timewindow)
                self.assertEquals(data,
                    [point for point in sorted_points
                    if point[0] in timewindow])

                self.storage.remove(
                    data_id=data_id,
                    period=period, timewindow=timewindow
                )

                # check for count equals 1
                count = self.storage.count(
                    data_id=data_id,
                    period=period, timewindow=timewindow)
                self.assertEquals(count, 0)
                count = self.storage.count(data_id=data_id, period=period)
                self.assertEquals(
                    count, len(points) - points_count_in_timewindow)

                self.storage.remove(data_id=data_id, period=period)
                # check for count equals 0
                count = self.storage.count(data_id=data_id, period=period)
                self.assertEquals(count, 0)

if __name__ == '__main__':
    main()
