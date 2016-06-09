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
from __future__ import unicode_literals

from canopsis.context.manager import Context
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.task.core import register_task
from canopsis.event import Event

from canopsis.old.account import Account
from canopsis.old.storage import get_storage

from datetime import datetime, timedelta
from icalendar import Event as vEvent


ctxmgr = Context()  #: default context manager
pbmgr = PBehaviorManager()  #: default pbehavior manager

events = get_storage(
    namespace='events',
    account=Account(user='root', group='root')
).get_backend()

DOWNTIME = 'downtime'  #: downtime pbehavior value

DOWNTIME_QUERY = PBehaviorManager.get_query(behaviors=DOWNTIME)


@register_task
def event_processing(
        engine, event, context=None, manager=None, logger=None, **kwargs
):
    """Process input event.

    :param dict event: event to process.
    :param Engine engine: engine which consumes the event.
    :param Context manager: context manager to use. Default is shared ctxmgr.
    :param PBehaviorManager manager: pbehavior manager to use. Default is
        pbmgr.
    :param Logger logger: logger to use in this task.
    """

    if context is None:
        context = ctxmgr

    if manager is None:
        manager = pbmgr

    evtype = event[Event.TYPE]
    entity = context.get_entity(event)

    encoded_entity = {}
    for k, v in entity.items():
        try:
            k = k.encode('utf-8')
        except:
            pass
        try:
            v = v.encode('utf-8')
        except:
            pass
        encoded_entity[k] = v

    eid = context.get_entity_id(encoded_entity)

    if evtype == DOWNTIME:
        ev = vEvent()
        ev.add('X-Canopsis-BehaviorType', '["{0}"]'.format(DOWNTIME))
        ev.add('summary', event['output'])
        ev.add('dtstart', datetime.fromtimestamp(event['start']))
        ev.add('dtend', datetime.fromtimestamp(event['end']))
        ev.add('dtstamp', datetime.fromtimestamp(event['entry']))
        ev.add('duration', timedelta(event['duration']))
        ev.add('contact', event['author'])

        manager.put(source=eid, vevents=[ev])

        if manager.getending(
            source=eid, behaviors=DOWNTIME, ts=event['timestamp']
        ):
            events.update(
                {
                    'connector': event['connector'],
                    'connector_name': event['connector_name'],
                    'component': event['component'],
                    'resource': event.get('resource', None)
                },
                {
                    '$set': {
                        DOWNTIME: True
                    }
                }
            )

    else:
        event[DOWNTIME] = manager.getending(
            source=eid, behaviors=DOWNTIME
        ) is not None

    return event


@register_task
def beat_processing(engine, context=None, manager=None, logger=None, **kwargs):
    """Process periodic task.

    :param Engine engine: engine which consumes the event.
    :param Context manager: context manager to use. Default is shared ctxmgr.
    :param PBehaviorManager manager: pbehavior manager to use. Default is
        pbmgr.
    :param Logger logger: logger to use in this task.
    """

    if context is None:
        context = ctxmgr

    if manager is None:
        manager = pbmgr

    entity_ids = manager.whois(query=DOWNTIME_QUERY)
    entities = context.get_entities(list(entity_ids))

    spec = {}

    for key in ['connector', 'connector_name', 'component', 'resource']:
        spec[key] = {
            '$nin': [
                e.get(key, None)
                for e in entities
            ]
        }

    events.update(spec, {'$set': {DOWNTIME: False}})
