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

"""Module in charge of defining downtime processing in engines."""

from canopsis.context.manager import Context
from canopsis.downtime.manager import DowntimeManager
from canopsis.task import register_task
from canopsis.event import Event


context = Context()
dm = DowntimeManager()


@register_task
def event_processing(engine, event, manager=None, logger=None, **kwargs):
    """Process input event.

    :param dict event: event to process.
    :param Engine engine: engine which consumes the event.
    :param DowntimeManager manager: downtime manager to use.
    :param Logger logger: logger to use in this task.
    """

    if manager is None:
        manager = dm

    evtype = event[Event.TYPE]
    entity = context.get_entity(event)
    entity_id = context.get_entity_id(entity)

    if evtype == 'downtime':
        # manager.put(entity_id, ical_downtime)
        pass

    else:
        event['downtime'] = manager.isdown(entity_id)

    return event
