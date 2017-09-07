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

"""
pbehavior process
"""

from __future__ import unicode_literals

from json import dumps

from canopsis.context_graph.manager import ContextGraph
from canopsis.common.utils import singleton_per_scope
from canopsis.task.core import register_task
from canopsis.pbehavior.manager import PBehaviorManager, PBehavior
from canopsis.watcher.manager import Watcher as WatcherManager


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

WATCHER_MANAGER = WatcherManager()


@register_task
def event_processing(engine, event, pbm=None, logger=None,
                     watcher_manager=WATCHER_MANAGER, **kwargs):
    """
    Event processing.
    """
    if pbm is None:
        pbm = singleton_per_scope(PBehaviorManager)

    if event.get('event_type') == EVENT_TYPE:
        entity_id = ContextGraph.get_id(event)
        logger.debug("Start processing event {}".format(event))

        logger.debug("entity_id: {}\naction: {}".format(
            entity_id, event.get('action')))

        try:
            pb_start = event.get('start')
            pb_end = event.get('end')
            pb_connector = event.get('connector')
            pb_name = event.get('pbehavior_name')
            pb_connector_name = event.get('connector_name')
        except KeyError as ex:
            logger.error('missing key in event: {}'.format(ex))
            return event

        pb_rrule = event.get('rrule', None)
        pb_comments = event.get('comments', None)
        pb_author = event.get('author', DEFAULT_AUTHOR)

        try:
            filter_ = {'_id': entity_id}
            if event.get('action') == PBEHAVIOR_CREATE:
                result = pbm.create(
                    pb_name, filter_, pb_author,
                    pb_start, pb_end,
                    connector=pb_connector,
                    comments=pb_comments,
                    connector_name=pb_connector_name,
                    rrule=pb_rrule
                )
                if not result:
                    logger.error(ERROR_MSG.format(event['action'], event))

                else:
                    watcher_manager.compute_watchers()

            elif event.get('action') == PBEHAVIOR_DELETE:
                result = pbm.delete(_filter={
                    PBehavior.FILTER: dumps(filter_),
                    PBehavior.NAME: pb_name,
                    PBehavior.TSTART: pb_start,
                    PBehavior.TSTOP: pb_end,
                    PBehavior.RRULE: pb_rrule,
                    PBehavior.CONNECTOR: pb_connector,
                    PBehavior.CONNECTOR_NAME: pb_connector_name,
                })
                if not result:
                    logger.error(ERROR_MSG.format(event['action'], event))
                else:
                    watcher_manager.compute_watchers()

            else:
                logger.error(ERROR_MSG.format(event['action'], event))

        except ValueError as err:
            logger.error('cannot handle event: {}'.format(err))

    return event


@register_task
def beat_processing(engine, pbm=None, logger=None, **kwargs):
    """
    Beat processing.
    """
    logger.debug("Start beat processing")

    if pbm is None:
        pbm = singleton_per_scope(PBehaviorManager)
    try:
        pbm.compute_pbehaviors_filters()
    except Exception as err:
        logger.error('Processing error {}'.format(str(err)))
    pbm.launch_update_watcher(WATCHER_MANAGER)
