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

from canopsis.context.manager import Context
from canopsis.common.utils import singleton_per_scope
from canopsis.task.core import register_task
from canopsis.pbehavior.manager import PBehaviorManager, PBehavior


EVENT_TYPE = 'pbehavior'
PBEHAVIOR_CREATE = 'create'
PBEHAVIOR_DELETE = 'delete'


def logger_error(logger, *args):
    logger.error('Failed to perform the action {} '
                 'for the event: {}'.format(*args))


@register_task
def event_processing(engine, event, pbm=None, logger=None, **kwargs):
    if pbm is None:
        pbm = singleton_per_scope(PBehaviorManager)

    if event.get('event_type') == EVENT_TYPE:
        cm = singleton_per_scope(Context)
        entity = cm.get_entity(event)
        entity_id = cm.get_entity_id(entity)

        filter = {'entity_id': entity_id}
        if event.get('action') == PBEHAVIOR_CREATE:
            result = pbm.create(
                event['pbehavior_name'], filter, event['author'],
                event['start'], event['end'],
                connector=event['connector'],
                connector_name=event['connector_name']
            )
            if not result:
                logger_error(logger, event['action'], event)

        elif event.get('action') == PBEHAVIOR_DELETE:
            result = pbm.delete(_filter={
                PBehavior.FILTER: dumps(filter),
                PBehavior.NAME: event['pbehavior_name'],
                PBehavior.TSTART: event['start'],
                PBehavior.TSTOP: event['end'],
                PBehavior.RRULE: '',
                PBehavior.CONNECTOR: event['connector'],
                PBehavior.CONNECTOR_NAME: event['connector_name'],
            })
            if not result:
                logger_error(logger, event['action'], event)
        else:
            logger_error(logger, event['action'], event)


@register_task
def beat_processing(engine, pbm=None, logger=None, **kwargs):
    if pbm is None:
        pbm = singleton_per_scope(PBehaviorManager)
    try:
        pbm.compute_pbehaviors_filters()
    except Exception as e:
        logger.error('Processing error {}'.format(str(e)))

