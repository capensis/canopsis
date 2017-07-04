# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from json import dumps

from canopsis.context_graph.manager import ContextGraph
from canopsis.common.utils import singleton_per_scope
from canopsis.task.core import register_task
from canopsis.pbehavior.manager import PBehaviorManager, PBehavior
from canopsis.watcher.manager import Watcher


EVENT_TYPE = 'pbehavior'
PBEHAVIOR_CREATE = 'create'
PBEHAVIOR_DELETE = 'delete'
ERROR_MSG = 'Failed to perform the action {} for the event: {}'

COMPONENT = 'component'
RESOURCE = 'resource'
SELECTOR = 'selector'
CONNECTOR = 'connector'
CONNECTOR_NAME = 'connector_name'

TEMPLATE = '/{}/{}/{}/{}'
TEMPLATE_RESOURCE = '/{}/{}/{}/{}/{}'

watcher_manager = Watcher()


def get_entity_id(event):
    source_type = event['source_type']

    if source_type == COMPONENT:
        first_word = COMPONENT
    elif source_type == RESOURCE:
        first_word = RESOURCE
    else:
        first_word = SELECTOR

    args = [first_word, event.get(CONNECTOR, ''), event.get(CONNECTOR_NAME, '')]

    template = TEMPLATE
    if first_word == RESOURCE:
        template = TEMPLATE_RESOURCE
        args.append(event.get(COMPONENT, ''))

    args.append(event.get(first_word, ''))

    entity_id = template.format(*args)

    return entity_id


@register_task
def event_processing(engine, event, pbm=None, logger=None, **kwargs):

    if pbm is None:
        pbm = singleton_per_scope(PBehaviorManager)

    if event.get('event_type') == EVENT_TYPE:
        entity_id = ContextGraph.get_id(event)
        logger.debug("Start processing event {}".format(event))

        logger.debug("entity_id: {}\naction: {}".format(
            entity_id, event.get('action'))
        )

        filter = {'_id': entity_id}
        if event.get('action') == PBEHAVIOR_CREATE:
            result = pbm.create(
                event['pbehavior_name'], filter, event['author'],
                event['start'], event['end'],
                connector=event['connector'],
                comments=event.get('comments', None),
                connector_name=event['connector_name'],
                rrule=event['rrule']
            )
            if not result:
                logger.error(ERROR_MSG.format(event['action'], event))

        elif event.get('action') == PBEHAVIOR_DELETE:
            result = pbm.delete(_filter={
                PBehavior.FILTER: dumps(filter),
                PBehavior.NAME: event['pbehavior_name'],
                PBehavior.TSTART: event['start'],
                PBehavior.TSTOP: event['end'],
                PBehavior.RRULE: event['rrule'],
                PBehavior.CONNECTOR: event['connector'],
                PBehavior.CONNECTOR_NAME: event['connector_name'],
            })
            if not result:
                logger.error(ERROR_MSG.format(event['action'], event))
        else:
            logger.error(ERROR_MSG.format(event['action'], event))

    watcher_manager.compute_watchers()

    return event


@register_task
def beat_processing(engine, pbm=None, logger=None, **kwargs):
    logger.debug("Start beat processing")

    if pbm is None:
        pbm = singleton_per_scope(PBehaviorManager)
    try:
        pbm.compute_pbehaviors_filters()
    except Exception as e:
        logger.error('Processing error {}'.format(str(e)))
