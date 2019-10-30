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
import os

from canopsis.common.utils import singleton_per_scope
from canopsis.context_graph.manager import ContextGraph
from canopsis.task.core import register_task
from canopsis.models.pbehavior import PBehavior as PBehaviorModel
from canopsis.pbehavior.manager import PBehaviorManager, PBehavior
from canopsis.watcher.manager import Watcher

EVENT_TYPE = 'pbehavior'
PBEHAVIOR_CREATE = 'create'
PBEHAVIOR_DELETE = 'delete'
ERROR_MSG = 'Failed to perform the action {} for the event: {}'

DOWNTIME_ID = 'downtime_id'
COMPONENT = 'component'
RESOURCE = 'resource'
SELECTOR = 'selector'
CONNECTOR = 'connector'
CONNECTOR_NAME = 'connector_name'
DEFAULT_AUTHOR = 'Default Author'

ENV_RECOMPUTE_ON_NEW_ENTITY = 'CPS_RECOMPUTE_PBEHAVIOR_ON_NEW_ENTITY'

watcher_manager = Watcher()

# known_entities is a set containing the ids of the entities corresponding to
# the events that have passed through the pbehavior engine. It is used when the
# CPS_RECOMPUTE_PBEHAVIORS_ON_NEW_ENTITY option is enabled, and only works with
# che.
known_entities = set()


def init_managers():
    """
    Init managers [sic].
    """
    config, pb_logger, pb_collection = PBehaviorManager.provide_default_basics()
    pb_kwargs = {'config': config,
                 'logger': pb_logger,
                 'pb_collection': pb_collection}
    pb_manager = singleton_per_scope(PBehaviorManager, kwargs=pb_kwargs)

    return pb_manager

_pb_manager = init_managers()


def get_entity_id(event):
    """
    get entity id from event.
    """
    return ContextGraph.get_id(event)


def pb_id(event):
    """
    Build a pbehavior ID from event if applicable.

    :returns: (None, None) if not applicable, (id, source) if applicable
    """
    did = event.get(DOWNTIME_ID)
    connector = event.get(CONNECTOR)
    connector_name = event.get(CONNECTOR_NAME)

    if did is not None and connector is not None and connector_name is not None:
        return 'pb_downtime_{}-{}_{}'.format(connector, connector_name, did), 'nagioslike'

    return None, None


@register_task
def event_processing(engine, event, pbm=_pb_manager, logger=None, **kwargs):
    """
    Event processing.
    """
    # This is a hack to make sure that new entities created by the go engine
    # che are immediately added to the pbehaviors. It is required to know
    # immediately if an entity is in maintenance, and to prevent tickets from
    # being declared in this case.
    # The pbehavior engine needs to receive events from che and publish them to
    # axe for this to work.
    if os.environ.get(ENV_RECOMPUTE_ON_NEW_ENTITY) == '1':
        entity_id = event.get('current_entity', {}).get('_id')
        if entity_id and entity_id not in known_entities:
            try:
                pbm.compute_pbehaviors_filters()
            except Exception as ex:
                engine.logger.exception('Processing error {}'.format(str(ex)))

            known_entities.add(entity_id)

    if event.get('event_type') == EVENT_TYPE:
        entity_id = ContextGraph.get_id(event)
        engine.logger.debug("Start processing event {}".format(event))

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
            pbehavior_id, pb_source = pb_id(event)
            if event.get('action') == PBEHAVIOR_CREATE and pbehavior_id is None and pb_source is None:
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

            elif event.get('action') == PBEHAVIOR_CREATE and pbehavior_id is not None and pb_source is not None:
                pbehavior = PBehaviorModel(
                    pbehavior_id, pb_name, filter_, pb_start, pb_end, pb_rrule, pb_author,
                    connector=pb_connector,
                    connector_name=pb_connector_name,
                    source=pb_source
                )
                success, result = pbm.upsert(pbehavior)
                if not success:
                    logger.critical('pbehavior upsert: {}'.format(result))

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
                logger.error(ERROR_MSG.format(event.get('action', 'no_action'), event))

        except ValueError as err:
            logger.error('cannot handle event: {}'.format(err))

    return event


@register_task
def beat_processing(engine, pbm=_pb_manager, **kwargs):
    """
    Beat processing.
    """
    engine.logger.debug("Start beat processing")

    try:
        pbm.compute_pbehaviors_filters()
        pbm.launch_update_watcher(watcher_manager)
    except Exception as ex:
        engine.logger.exception('Processing error {}'.format(str(ex)))
