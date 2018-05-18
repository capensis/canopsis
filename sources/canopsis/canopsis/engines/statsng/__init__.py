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

from __future__ import unicode_literals
try:
    from threading import Lock
except ImportError:
    from dummy_threading import Lock

from influxdb import InfluxDBClient
from influxdb.exceptions import InfluxDBClientError

from canopsis.engines.core import Engine

SECONDS = 1000000000


class engine(Engine):
    etype = "statsng"

    def pre_run(self):
        self.batch_lock = Lock()
        self.batch = []

        self.max_batch_size = 5000

        self.influx_client = InfluxDBClient('192.168.0.62', 8086, 'root', 'root', 'statsng-test')
        try:
            self.influx_client.create_database('statsng-test')
        except InfluxDBClientError:
            pass

    def beat(self):
        with self.batch_lock:
            self.flush()

    def work(self, event, *args, **kargs):
        """
        AMQP event processing.

        :param dict event: event to process.
        """

        if event['event_type'] == 'statcounterinc':
            self.handle_statcounterinc_event(event)
        elif event['event_type'] == 'statduration':
            self.handle_statduration_event(event)

    def handle_statcounterinc_event(self, event):
        self.logger.info('received statcounterinc event')

        alarm = event['alarm']

        self.add_point({
            'measurement': 'statcounters',
            'time': event['timestamp'] * SECONDS,
            'tags': {
                'connector': alarm['connector'],
                'connector_name': alarm['connector_name'],
                'component': alarm['component'],
                'resource': alarm['resource']
            },
            'fields': {
                event['counter_name']: 1
            }
        })

    def handle_statduration_event(self, event):
        self.logger.info('received statduration event')

        alarm = event['alarm']

        self.add_point({
            'measurement': 'statdurations',
            'time': event['timestamp'] * SECONDS,
            'tags': {
                'connector': alarm['connector'],
                'connector_name': alarm['connector_name'],
                'component': alarm['component'],
                'resource': alarm['resource']
            },
            'fields': {
                event['duration_name']: event['duration']
            }
        })

    def add_point(self, point):
        with self.batch_lock:
            self.batch.append(point)

            if len(self.batch) >= self.max_batch_size:
                self.flush()

    def flush(self):
        if self.batch:
            self.influx_client.write_points(self.batch)
            self.batch = []
