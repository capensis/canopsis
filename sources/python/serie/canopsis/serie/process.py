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

from canopsis.timeserie.timewindow import TimeWindow
from canopsis.serie.manager import Serie

from time import time


@register_task
def beat_processing(engine, manager=None, logger=None, **kwargs):
    if manager is None:
        manager = singleton_per_scope(Serie)

    for serie in manager.get_series(time()):
        publish(
            publisher=engine.amqp,
            event=serie,
            rk=engine.amqp_queue,
            exchange='amq.direct',
            logger=logger
        )


@register_task
def serie_processing(engine, event, manager=None, logger=None, **kwargs):
    if manager is None:
        manager = singleton_per_scope(Serie)

    now = time()

    timewin = TimeWindow(
        start=event['last_computation'],
        stop=now
    )

    point = manager.calculate(event, timewin)

    metric = {
        'metric': event['crecord_name'],
        'value': point,
        'type': 'GAUGE'
    }

    for meta in ['unit', 'min', 'max', 'warn', 'crit']:
        if event.get(meta, None) is not None:
            metric[meta] = event[meta]

    event = {
        'timestamp': int(now),
        'connector': 'canopsis',
        'connector_name': engine.name,
        'event_type': 'perf',
        'source_type': 'resource',
        'component': event['component'],
        'resource': event['resource'],
        'perf_data_array': [metric]
    }

    publish(publisher=engine.amqp, event=event, logger=logger)
