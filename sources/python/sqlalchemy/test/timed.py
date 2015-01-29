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

from canopsis.mongo.timed import MongoTimedStorage
from canopsis.timeserie.timewindow import TimeWindow


class TimedStorageTest(TestCase):
    """
    MongoTimedStorage UT on data_type = "test_store"
    """

    def setUp(self):
        # create a storage on test_store collections
        self.storage = MongoTimedStorage(data_scope="test_store", safe=True)

    def test_connect(self):
        self.assertTrue(self.storage.connected())

        self.storage.disconnect()

        self.assertFalse(self.storage.connected())

        self.storage.connect()

        self.assertTrue(self.storage.connected())

    def test_CRUD(self):

        data_id = 'test_store_id'

        # start in droping data
        self.storage.drop()

        # ensure count is 0
        count = self.storage.count(data_id=data_id)
        self.assertEquals(count, 0)

        # let's play with different data_names
        meta = {'min': None, 'max': 0}

        timewindow = TimeWindow()

        before_timewindow = [timewindow.start() - 1000]
        in_timewindow = [
            timewindow.start(),
            timewindow.start() + 5,
            timewindow.stop() - 5,
            timewindow.stop()
        ]
        after_timewindow = [timewindow.stop() + 1000]

        # set timestamps without sort
        timestamps = after_timewindow + before_timewindow + in_timewindow

        for timestamp in timestamps:
            # add a document at starting time window
            self.storage.put(data_id=data_id, value=meta, timestamp=timestamp)

        # check for count equals 5
        count = self.storage.count(data_id=data_id)
        self.assertEquals(count, 2)

        # check for_data before now
        data = self.storage.get(data_ids=[data_id])
        self.assertEquals(len(data[data_id]), 2)

        # check for data inside timewindow and just before
        data = self.storage.get(data_ids=[data_id], timewindow=timewindow)
        self.assertEquals(len(data), 1)

        # remove data inside timewindow
        self.storage.remove(data_ids=[data_id], timewindow=timewindow)
        # check for data outside timewindow
        count = self.storage.count(data_id=data_id)
        self.assertEquals(
            count, len(before_timewindow) + len(after_timewindow))

        # remove all data
        self.storage.remove(data_ids=[data_id])
        # check for count equals 0
        count = self.storage.count(data_id=data_id)
        self.assertEquals(count, 0)

if __name__ == '__main__':
    main()
