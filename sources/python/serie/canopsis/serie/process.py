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

"""Module of serie processing tasks."""

from canopsis.common.utils import singleton_per_scope
from canopsis.task.core import register_task
from canopsis.engines.core import publish

from canopsis.serie.manager import Serie

from time import time


@register_task
def beat_processing(engine, manager=None, logger=None, **_):
    """Engine beat processing task."""

    if manager is None:
        manager = singleton_per_scope(Serie)

    with engine.Lock(engine, 'serie_fetching') as lock:
        if lock.own():
            for serie in manager.get_series(time()):
                publish(
                    publisher=engine.amqp,
                    event=serie,
                    rk=engine.amqp_queue,
                    exchange='amq.direct',
                    logger=logger
                )


@register_task
def serie_processing(engine, event, manager=None, logger=None, **_):
    """Engine work processing task."""

    if manager is None:
        manager = singleton_per_scope(Serie)

    # Generate metric metadata
    metric_tags = {
        tags: event[tags]
        for tags in ['unit', 'min', 'max', 'warn', 'crit']
        if event.get(tags, None) is not None
    }
    metric_tags['type'] = 'GAUGE'

    # Generate metric entity
    entity = {
        'type': 'metric',
        'connector': 'canopsis',
        'connector_name': engine.name,
        'component': event['component'],
        'resource': event['resource'],
        'name': event['crecord_name']
    }

    context = manager[Serie.CONTEXT_MANAGER]
    entity_id = context.get_entity_id(entity)

    # Publish points
    perfdata = manager[Serie.PERFDATA_MANAGER]
    perfdata.put(
        entity_id,
        points=manager.calculate(event),
        tags=metric_tags,
        cache=False
    )
