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
import time


ctxmgr = Context()  #: default context manager
pbmgr = PBehaviorManager()  #: default pbehavior manager

events_storage = get_storage(
    namespace='events',
    account=Account(user='root', group='root')
).get_backend()

DOWNTIME = 'downtime'  #: downtime pbehavior value
PBEHAVIOR = 'pbehavior'  #: another downtime pbehavior value

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

    if evtype == DOWNTIME or (evtype == PBEHAVIOR and event['pbehavior_name'] == DOWNTIME):
        action = event["action"]
        # listing identical events
        bhvs = manager.get_behaviors(entity_id=eid)
        uids = [b.get('_id', None) for b in bhvs if b['dtstart'] == event['start'] and b['dtend'] == event['end']]

        if action == 'delete':
            # Remove corresponding behaviors (same source/start/end)
            logger.info('Removing behaviors {}'.format(uids))
            manager.remove(uids=uids)

        elif action == 'create' and not event["was_started"]:
            if len(uids) > 0:
                # Behavior already created !
                logger.info('Behavior already created for {}'.format(eid))
                return event

            logger.info('Creating behavior for source {}'.format(eid))
            ev = vEvent()
            ev.add('X-Canopsis-BehaviorType', '["{}"]'.format(DOWNTIME))
            ev.add('summary', event['output'])
            ev.add('dtstart', datetime.fromtimestamp(event['start']))
            ev.add('dtend', datetime.fromtimestamp(event['end']))
            ev.add('dtstamp', datetime.fromtimestamp(event['entry']))
            if not event.get('fixed', True):
                ev.add('duration', timedelta(seconds=event['duration']))
            ev.add('contact', event['author'])

            manager.put(source=eid, vevents=[ev])

        elif event["was_started"]:
            logger.info('Behavior starting event for {}. Ignoring'.format(eid))

    else:
        event[DOWNTIME] = manager.getending(source=eid, behaviors=DOWNTIME) is not None

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
    entities = list(context.get_entities(list(entity_ids)))
    #logger.debug('BEAT: {} //// {}'.format(entity_ids, entities))

    # Generating event filters
    unsetting = {}
    setting = {}
    for entity in entities:
        table = entity['_id'].split('/')[2:]  # striping '' and 'pbehavior'
        for idx, key in enumerate(['connector', 'connector_name', 'component', 'resource']):
            index = key
            if idx == len(table) - 1:
                # The last key of entity is named 'name'
                index = 'name'
            if index not in entity:
                continue

            # List init
            if key not in unsetting:
                unsetting[key] = {'$nin': []}
            if key not in setting:
                setting[key] = {'$in': []}

            unsetting[key]['$nin'].append(entity[index])
            setting[key]['$in'].append(entity[index])

    # (Un)setting passed behaviors
    events_storage.update(unsetting, {'$set': {DOWNTIME: False}})

    # Enabling started behaviors
    events = events_storage.find(setting)
    now = int(time.time())
    for event in events:
        if DOWNTIME in event and event[DOWNTIME]:
            #logger.debug('Downtime already set for {}'.format(event['_id']))
            continue

        entity = context.get_entity(event, from_db=True)
        entity_id = entity['_id'].replace('/resource/', '/pbehavior/')  # Dirty part !

        bhvs = manager.get_behaviors(entity_id=entity_id)[0]
        mge = manager.getending(source=entity_id, behaviors=DOWNTIME)
        if 'dtstart' in bhvs and mge and bhvs['dtstart'] <= now <= mge:
            logger.info('Enabling downtime for event {}'.format(event['_id']))
            events_storage.update({'_id': event['_id']},
                                  {'$set': {DOWNTIME: True}})
