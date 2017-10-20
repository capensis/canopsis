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

from __future__ import print_function

from unittest import main

from canopsis.timeserie.timewindow import TimeWindow
from canopsis.configuration.model import Parameter
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category

from base import BaseTestConfiguration, BaseStorageTest

from time import time


@conf_paths('storage/timedbench.conf')
@add_category(
    'TEST',
    content=[
        Parameter('batch_sizes', Parameter.array),
        Parameter('counts', Parameter.array),
        Parameter('defaultcount', int)
    ]
)
class TestConfiguration(BaseTestConfiguration):
    """Default test configuration."""


class TimedBench(BaseStorageTest):

    def _testconfcls(self):

        return TestConfiguration

    def _test(self, storage):

        def measure(msg, count, batch, func, **kwargs):
            data_id = '{0}-{1}'.format(count, batch)
            print(msg)
            lastts = time()
            func(data_id=data_id, **kwargs)
            lastts = time() - lastts
            print('Duration: {0}'.format(lastts))

        print('Starting Bench of {0}:'.format(storage))

        print('Starting to measure CRUD operations')

        for count in self.testconf.counts:

            count = int(count)

            print('Number of perfdata: {0}'.format(count))

            for batch_size in self.testconf.batch_sizes:

                print('Batch size: {0}'.format(batch_size))

                points = [(i, i) for i in range(count)]

                measure(
                    msg='Writing operation in process...',
                    func=storage.put, points=points,
                    count=count, batch=batch_size
                )

                measure(
                    msg='Read operation in process...',
                    func=storage.get,
                    timewindow=TimeWindow(
                        start=points[0][0], stop=points[-1][0]
                    ),
                    count=count, batch=batch_size
                )

                measure(
                    msg='Deleting operation in process...',
                    func=storage.remove, count=count, batch=batch_size
                )


if __name__ == '__main__':
    main()
