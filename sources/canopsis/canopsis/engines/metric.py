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
from canopsis.common.influx import SECONDS, InfluxDBClient
from canopsis.confng import Configuration, Ini
from canopsis.confng.helpers import cfg_to_array
from canopsis.context_graph.manager import ContextGraph
from canopsis.engines.core import Engine
from canopsis.event import Event
from canopsis.models.entity import Entity
from canopsis.monitoring.parser import PerfDataParser


class engine(Engine):
    etype = "metric"

    CONF_PATH = "etc/metric/engine.conf"
    CONF_SECTION = 'ENGINE'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.context_manager = ContextGraph(self.logger)
        self.influxdb_client = InfluxDBClient.from_configuration(self.logger)

        cfg = Configuration.load(
            os.path.join(root_path, self.CONF_PATH), Ini
        ).get(self.CONF_SECTION, {})
        self.tags = cfg_to_array(cfg.get('tags', ''))

    def work(self, event, *args, **kwargs):
        """
        AMQP event processing.

        :param dict event: event to process.
        """
        # Get perfdata
        perf_data = event.get('perf_data')
        perf_data_array = event.get('perf_data_array', [])

        if perf_data_array is None:
            perf_data_array = []

        # If the event does not have a resource, no perfdata can be created.
        # Ignore events without perf_data.
        if "resource" not in event or (not perf_data and not perf_data_array):
            return

        # Parse perfdata
        if perf_data:
            self.logger.debug(u' + perf_data: {0}'.format(perf_data))

            try:
                parser = PerfDataParser(perf_data)
                perf_data_array += parser.perf_data_array
            except Exception as err:
                self.logger.error(
                    "Impossible to parse perfdata from: {0} ({1})".format(
                        event, err
                    )
                )

        self.logger.debug(u'perf_data_array: {0}'.format(perf_data_array))

        # Write perfdata to influx
        timestamp = event['timestamp'] * SECONDS
        tags = self.get_tags(event)

        points = []
        for data in perf_data_array:
            metric = data.get('metric')
            value = data.get('value')
            warn = data.get('warn')
            crit = data.get('crit')

            if value is not None and metric:
                point = {
                    'measurement': metric,
                    'time': timestamp,
                    'tags': tags,
                    'fields': {
                        'value': value
                    }
                }

                if warn is not None:
                    point['fields']['warn'] = warn
                if crit is not None:
                    point['fields']['crit'] = crit

                points.append(point)

        self.influxdb_client.write_points(points)

    def get_tags(self, event):
        """
        Returns the tags corresponding to an event, to be used in
        `InfluxDBClient.write_points`.

        :param dict event:
        :rtype dict:
        """
        tags = {
            'connector': event[Event.CONNECTOR],
            'connector_name': event[Event.CONNECTOR_NAME],
            'component': event[Event.COMPONENT],
            'resource': event[Event.RESOURCE]
        }

        entity = self.context_manager.get_entities_by_id(event['_id'])
        try:
            entity = entity[0]
        except IndexError:
            entity = {}

        infos = entity.get(Entity.INFOS, {})
        for tag in self.tags:
            tags[tag] = infos.get(tag, {}).get('value', '')

        return tags
