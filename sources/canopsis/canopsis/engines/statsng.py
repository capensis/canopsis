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
import os
try:
    from threading import Lock
except ImportError:
    from dummy_threading import Lock

from influxdb import InfluxDBClient
from influxdb.exceptions import InfluxDBClientError

from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.engines.core import Engine

SECONDS = 1000000000


class engine(Engine):
    etype = "statsng"

    CONF_PATH = "etc/statsng/engine.conf"
    DEFAULT_DATABASE = 'statsng'
    DEFAULT_MAX_BATCH_SIZE = 100

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.batch_lock = Lock()
        self.batch = []

        cfg = Configuration.load(os.path.join(root_path, self.CONF_PATH), Ini)
        batch_cfg = cfg.get('BATCH', {})
        influxdb_cfg = cfg.get('DATABASE', {})

        self.max_batch_size = int(
            batch_cfg.get('max_batch_size', self.DEFAULT_MAX_BATCH_SIZE))

        # Set the default database name (use InfluxDBClient's default for the
        # other values)
        if 'database' not in influxdb_cfg:
            influxdb_cfg['database'] = self.DEFAULT_DATABASE

        self.influx_client = InfluxDBClient(**influxdb_cfg)
        try:
            self.influx_client.create_database(influxdb_cfg['database'])
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
