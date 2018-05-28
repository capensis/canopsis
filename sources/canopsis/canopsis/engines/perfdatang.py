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

from canopsis.common.influx import SECONDS, get_influxdb_client
from canopsis.engines.core import Engine
from canopsis.event import Event
from canopsis.monitoring.parser import PerfDataParser

MEASUREMENT = 'perfdata'


class engine(Engine):
    etype = "perfdatang"

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.influxdb_client = get_influxdb_client()

    def work(self, event, *args, **kwargs):
        """
        AMQP event processing.

        :param dict event: event to process.
        """
        # If the event does not have a resource, no perfdata can be created.
        if "resource" not in event:
            return

        # Get perfdata
        perf_data = event.get('perf_data')
        perf_data_array = event.get('perf_data_array', [])

        if perf_data_array is None:
            perf_data_array = []

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

        tags = {
            'connector': event[Event.CONNECTOR],
            'connector_name': event[Event.CONNECTOR_NAME],
            'component': event[Event.COMPONENT],
            'resource': event[Event.RESOURCE]
        }

        fields = {}
        for data in perf_data_array:
            metric = data.get('metric')
            value = data.get('value')

            if value is not None and metric:
                fields[metric] = value

        self.influxdb_client.write_points([{
            'measurement': MEASUREMENT,
            'time': event['timestamp'] * SECONDS,
            'tags': tags,
            'fields': fields,
        }])
