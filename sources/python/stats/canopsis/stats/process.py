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

from canopsis.common.utils import singleton_per_scope
from canopsis.task.core import register_task
from canopsis.engines.core import publish

from canopsis.stats.producers.user import UserMetricProducer
from canopsis.stats.producers.event import EventMetricProducer

from time import time


@register_task
def event_processing(
    engine, event,
    usermgr=None, eventmgr=None,
    logger=None, **kwargs
):
    if usermgr is None:
        usermgr = singleton_per_scope(UserMetricProducer)

    if eventmgr is None:
        eventmgr = singleton_per_scope(EventMetricProducer)

    events = []

    for manager in [usermgr, eventmgr]:
        for counter in manager.counters(event):
            metric = {
                'metric': counter['filter'],
                'value': 1,
                'type': 'COUNTER'
            }

            for meta in ['unit', 'min', 'max', 'warn', 'crit']:
                if counter.get(meta, None) is not None:
                    metric[meta] = counter[meta]

            events.append({
                'timestamp': int(time()),
                'connector': 'canopsis',
                'connector_name': engine.name,
                'event_type': 'perf',
                'source_type': 'resource',
                'component': counter['component'],
                'resource': counter['name'],
                'perf_data_array': [metric]
            })

        for gauge in manager.gauges(event):
            metric = {
                'metric': 'last',
                'value': gauge['value'],
                'type': 'GAUGE'
            }

            for meta in ['unit', 'min', 'max', 'warn', 'crit']:
                if gauge.get(meta, None) is not None:
                    metric[meta] = gauge[meta]

            events.append({
                'timestamp': int(time()),
                'connector': 'canopsis',
                'connector_name': engine.name,
                'event_type': 'perf',
                'source_type': 'resource',
                'component': counter['component'],
                'resource': counter['name'],
                'perf_data_array': [metric]
            })

    for event in events:
        publish(publisher=engine.amp, event=event, logger=logger)
