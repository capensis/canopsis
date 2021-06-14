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

from canopsis.engines.core import Engine

from time import time


class engine(Engine):
    etype = "eventduration"

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        self.durations = []

    def work(self, event, *args, **kargs):
        now = time()
        ts = [now] + event['processing'].values()

        duration = max(ts) - min(ts)
        self.durations.append(duration)

    def beat(self):
        durations = []

        while True:
            try:
                durations.append(self.durations.pop())

            except IndexError:
                break

        if durations:
            durmin = min(durations)
            durmax = max(durations)
            duravg = sum(durations) / len(durations)

            event = {
                'connector': 'Engine',
                'connector_name': self.etype,
                'event_type': 'perf',
                'source_type': 'component',
                'component': '__canopsis__',
                'perf_data_array': [
                    {
                        'metric': 'cps_evt_duration_min',
                        'value': durmin, 'unit': 's', 'type': 'GAUGE'
                    },
                    {
                        'metric': 'cps_evt_duration_max',
                        'value': durmax, 'unit': 's', 'type': 'GAUGE'
                    },
                    {
                        'metric': 'cps_evt_duration_avg', 'value': duravg,
                        'unit': 's', 'type': 'GAUGE'
                    }
                ]
            }

            try:
                self.beat_amqp_publisher.canopsis_event(event)
            except Exception as e:
                self.logger.exception("Unable to send event")
