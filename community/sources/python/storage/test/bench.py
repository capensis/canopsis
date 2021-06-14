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

from base import BaseStorageTest

# from time import time

# try:
#     from threading import Thread

# except ImportError:
#     from dummy_threading import Thread


class Bench(BaseStorageTest):

    #: only one variable to set in order to instanciate the right storage.
    __protocol__ = None

    def setUp(self):
        """initalize storages."""

        super(Bench, self).setUp()

        # initialize numerical values
        self.iteration = 5
        self.count = 10000
        # generate count documents
        self.documents = list(
            {
                "_id": str(i),
                "data": {"index": i}
            } for i in range(self.count)
        )
        self.documents_to_update = list(
            {
                "$set": {
                    'data.index': None
                }
            } for i in range(self.count)
        )
        # initialize commands
        self.commands = (
            ("insert", lambda spec, document: self.storage._insert(
                document=document, cache=True)),
            ("update", lambda spec, document: self.storage._update(
                spec=spec, document=document, multi=False, cache=True)),
            ("find", lambda spec, document: self.storage._find(
                document=document, cache=True)),
            ("remove", lambda spec, document: self.storage._remove(
                document=document, multi=True, cache=True))
        )

        self.max_connections = 1

    # TODO 4-01-2017
    # def tearDown(self):
    #     """
    #     End the test
    #     """

    #     # remove all data from collection
    #     self.storage.drop()

    # TODO 4-01-2017
    # def test(self):
    #     """
    #     Run tests.
    #     """

    #     threads = [
    #         Thread(target=self._test_CRUD)
    #         for i in range(self.max_connections)
    #     ]

    #     for thread in threads:
    #         thread.start()

    #     for thread in threads:
    #         thread.join()

    # TODO 4-01-2017
    # def _test_CRUD(self):
    #     """
    #     Bench CRUD commands
    #     """

    #     print(
    #         'Starting bench on %s with %s documents'
    #         %
    #         (self.commands, self.count)
    #     )

    #     stats_per_command = {
    #         command[0]: {'min': None, 'max': None, 'mean': 0}
    #         for command in self.commands
    #     }

    #     min_t, max_t, mean_t = None, None, 0

    #     for i in range(self.iteration):

    #         total = 0

    #         for index, command in enumerate(self.commands):

    #             fn = command[1]
    #             command = command[0]

    #             start = time()

    #             if command == 'update':
    #                 documents = self.documents_to_update
    #             else:
    #                 documents = self.documents

    #             for j in range(self.count):

    #                 fn(
    #                     spec={'_id': str(self.count)},
    #                     document=documents[j]
    #                 )

    #             stop = time()

    #             duration = stop - start
    #             stats = stats_per_command[command]
    #             if stats['min'] is None or stats['min'] > duration:
    #                 stats['min'] = duration
    #             if stats['max'] is None or stats['max'] < duration:
    #                 stats['max'] = duration
    #             stats['mean'] += duration

    #             bench_msg = 'command: %s, iteration: %s, time: %s'

    #             print(
    #                 bench_msg % (command, i, duration)
    #             )

    #             total += duration

    #         if min_t is None or min_t > total:
    #             min_t = total
    #         if max_t is None or max_t < total:
    #             max_t = total
    #         mean_t += total

    #         bench_msg = 'CRUD: %s' % total

    #         print(bench_msg)

    #     # update mean per command
    #     for command in self.commands:
    #         stats = stats_per_command[command[0]]
    #         stats['mean'] = stats['mean'] / self.iteration
    #         bench_msg = 'command: %s, min: %s, max: %s, mean: %s'
    #         print(
    #             bench_msg
    #             %
    #             (command[0], stats['min'], stats['max'], stats['mean'])
    #         )

    #     # update mean for all CRUD operations
    #     mean_t = mean_t / self.iteration

    #     bench_msg = "all command, min: %s, max: %s, mean: %s"
    #     print(
    #         bench_msg % (min_t, max_t, mean_t)
    #     )


if __name__ == '__main__':
    main()
