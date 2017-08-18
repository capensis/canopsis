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
DEFAULT_AUTHOR = 'Default Author'

TEMPLATE = '/{}/{}/{}/{}'
TEMPLATE_RESOURCE = '/{}/{}/{}/{}/{}'

watcher_manager = Watcher()


def get_entity_id(event):
    return ContextGraph.get_id(event)


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
                event['pbehavior_name'], filter, event.get('author', DEFAULT_AUTHOR),
                event['start'], event['end'],
                connector=event['connector'],
                comments=event.get('comments', None),
                connector_name=event['connector_name'],
                rrule=event.get('rrule', None)
            )
            if not result:
                logger.error(ERROR_MSG.format(event['action'], event))

            else:
                watcher_manager.compute_watchers()


        elif event.get('action') == PBEHAVIOR_DELETE:
            result = pbm.delete(_filter={
                PBehavior.FILTER: dumps(filter),
                PBehavior.NAME: event['pbehavior_name'],
                PBehavior.TSTART: event['start'],
                PBehavior.TSTOP: event['end'],
                PBehavior.RRULE: event.get('rrule', None),
                PBehavior.CONNECTOR: event['connector'],
                PBehavior.CONNECTOR_NAME: event['connector_name'],
            })
            if not result:
                logger.error(ERROR_MSG.format(event['action'], event))
            else:
                watcher_manager.compute_watchers()

        else:
            logger.error(ERROR_MSG.format(event['action'], event))


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
