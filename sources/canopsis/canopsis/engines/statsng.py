# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.confng.helpers import cfg_to_array
from canopsis.engines.core import Engine
from canopsis.statsng.enums import StatEvents, StatMeasurements
from canopsis.statsng.influx import get_influxdb_client

SECONDS = 1000000000


class engine(Engine):
    etype = "statsng"

    CONF_PATH = "etc/statsng/engine.conf"
    CONF_SECTION = 'ENGINE'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        cfg = Configuration.load(
            os.path.join(root_path, self.CONF_PATH), Ini
        ).get(self.CONF_SECTION, {})

        self.influxdb_client = get_influxdb_client()

        self.tags = cfg_to_array(cfg.get('tags', ''))

    def work(self, event, *args, **kargs):
        """
        AMQP event processing.

        :param dict event: event to process.
        """

        if event['event_type'] == StatEvents.statcounterinc:
            self.handle_statcounterinc_event(event)
        elif event['event_type'] == StatEvents.statduration:
            self.handle_statduration_event(event)

    def get_tags(self, event):
        """
        Returns the tags corresponding to an event, to be used in
        `InfluxDBClient.write_points`.

        :param dict event:
        :rtype dict:
        """
        alarm = event['alarm']

        tags = {
            'connector': alarm['connector'],
            'connector_name': alarm['connector_name'],
            'component': alarm['component'],
            'resource': alarm['resource']
        }

        infos = event.get('entity', {}).get('infos', {})
        for tag in self.tags:
            tags[tag] = infos.get(tag, {}).get('value', '')

        return tags

    def handle_statcounterinc_event(self, event):
        """
        Process a statcounterinc event.

        :param dict event:
        """
        self.logger.info('received statcounterinc event')

        self.influxdb_client.write_points([{
            'measurement': StatMeasurements.counters,
            'time': event['timestamp'] * SECONDS,
            'tags': self.get_tags(event),
            'fields': {
                event['counter_name']: 1
            }
        }])

    def handle_statduration_event(self, event):
        """
        Process a statduration event.

        :param dict event:
        """
        self.logger.info('received statduration event')

        self.influxdb_client.write_points([{
            'measurement': StatMeasurements.durations,
            'time': event['timestamp'] * SECONDS,
            'tags': self.get_tags(event),
            'fields': {
                event['duration_name']: event['duration']
            }
        }])
