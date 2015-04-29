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
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.task import register_task
from canopsis.event import Event

from canopsis.old.account import Account
from canopsis.old.storage import get_storage

from datetime import datetime
from icalendar import Event as icalEvent


context = Context()
dm = PBehaviorManager()

events = get_storage(
    namespace='events',
    account=Account(user='root', group='root')
).get_backend()


@register_task
def event_processing(engine, event, manager=None, logger=None, **kwargs):
    """Process input event.

    :param dict event: event to process.
    :param Engine engine: engine which consumes the event.
    :param PBehaviorManager manager: pbehavior manager to use.
    :param Logger logger: logger to use in this task.
    """

    if manager is None:
        manager = dm

    evtype = event[Event.TYPE]
    entity = context.get_entity(event)
    entity_id = context.get_entity_id(entity)

    if evtype == 'downtime':
        ev = icalEvent()
        ev.add('summary', event['output'])
        ev.add('dtstart', datetime.fromtimestamp(event['start']))
        ev.add('dtend', datetime.fromtimestamp(event['end']))
        ev.add('dtstamp', datetime.fromtimestamp(event['entry']))
        ev.add('duration', event['duration'])
        ev.add('contact', event['author'])

        manager.put(entity_id, ev.to_ical())

        if manager.until(entity_id, 'downtime', event['timestamp']):
            events.update(
                {
                    'connector': event['connector'],
                    'connector_name': event['connector_name'],
                    'component': event['component'],
                    'resource': event.get('resource', None)
                },
                {
                    '$set': {
                        'downtime': True
                    }
                }
            )

    else:
        event['downtime'] = manager.until(entity_id, 'downtime') is not None

    return event
