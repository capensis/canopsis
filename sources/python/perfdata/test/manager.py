#!/usr/bin/env python
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

from unittest import TestCase, main

from canopsis.perfdata.manager import PerfData
from canopsis.timeserie.timewindow import TimeWindow, get_offset_timewindow


class PerfDataTest(TestCase):

    def setUp(self):
        self.perfdata = PerfData(data_scope='test')

    def test_put_get_data(self):

        timewindow = TimeWindow()

        metric_id = 'test_manager'

        self.perfdata.remove(metric_id=metric_id, with_meta=True)

        count = self.perfdata.count(metric_id=metric_id)
        self.assertEqual(count, 0)

        tv0 = (int(timewindow.start()), None)
        tv1 = (int(timewindow.start() + 1), 0)
        tv2 = (int(timewindow.stop()), 2)
        tv3 = (int(timewindow.stop() + 1000000), 3)

        # set values with timestamp without order
        points = [tv0, tv2, tv1, tv3]

        meta = {'plop': None}

        self.perfdata.put(
            metric_id=metric_id,
            points=points,
            meta=meta)

        data, _meta = self.perfdata.get(
            metric_id=metric_id,
            timewindow=timewindow,
            with_meta=True)

        self.assertEqual(meta, _meta[0][PerfData.META_VALUE])

        self.assertEqual([tv0, tv1, tv2], data)

        # remove 1 data at stop point
        _timewindow = get_offset_timewindow(timewindow.stop())

        self.perfdata.remove(
            metric_id=metric_id,
            timewindow=_timewindow)

        data, _meta = self.perfdata.get(
            metric_id=metric_id,
            timewindow=timewindow,
            with_meta=True)

        self.assertEqual(meta, _meta[0][PerfData.META_VALUE])

        self.assertEqual(data, [tv0, tv1])

        # get data on timewindow
        data, _meta = self.perfdata.get(
            metric_id=metric_id,
            timewindow=timewindow,
            with_meta=True)

        self.assertEqual(meta, _meta[0][PerfData.META_VALUE])

        # get all data
        data, _meta = self.perfdata.get(
            metric_id=metric_id,
            with_meta=True)

        self.assertEqual(meta, _meta[0][PerfData.META_VALUE])

        self.assertEqual(len(data), 3)

        # remove all data
        self.perfdata.remove(
            metric_id=metric_id,
            with_meta=True)

        data, _meta = self.perfdata.get(
            metric_id=metric_id,
            with_meta=True)

        self.assertEqual(len(_meta), 0)

        self.assertEqual(len(data), 0)

if __name__ == '__main__':
    main()
