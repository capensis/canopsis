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

"""Module in charge of defining topology processing in engines.

When an event occured, the related entity is retrieved and all bound
topological nodes are retrieved as well in order to execute theirs rules.

First, a topology processing is triggered when an event occured.

From this event, bound topology nodes are retrieved in order to apply node
    rules.

A typical topological task condition is an ``canopsis.task.condition.all``
composed of the ``canopsis.topology.rule.condition.new_state`` condition and
``canopsis.topology.rule.action.change_state`` action.
If this condition is checked, then other specific conditions can be applied
such as those defined in the canopsis.topology.rule.action module.
"""
from __future__ import unicode_literals

from canopsis.topology.elements import Topology, TopoNode
from canopsis.topology.manager import TopologyManager
from canopsis.context.manager import Context
from canopsis.task.core import register_task
from canopsis.event import Event
from canopsis.check.manager import CheckManager
from canopsis.common.utils import singleton_per_scope

SOURCE = 'source'
PUBLISHER = 'publisher'


@register_task
def event_processing(
        engine, event, manager=None, logger=None, ctx=None, tm=None, cm=None,
        **kwargs
):
    """Process input event in getting topology nodes bound to input event
    entity.

    One topology nodes are founded, executing related rules.

    :param dict event: event to process.
    :param Engine engine: engine which consumes the event.
    :param TopologyManager manager: topology manager to use.
    :param Logger logger: logger to use in this task.
    :param Context ctx:
    :param TopologManager tm:
    :param CheckManager cm:
    """

    # initialize ctx
    if ctx is None:
        ctx = singleton_per_scope(Context)

    if tm is None:
        tm = singleton_per_scope(TopologyManager)

    if cm is None:
        cm = singleton_per_scope(CheckManager)

    event_type = event[Event.TYPE]

    # apply processing only in case of check event
    if event_type in cm.types:
        # get source type
        source_type = event[Event.SOURCE_TYPE]
        # in case of topology node
        if source_type in [TopoNode.TYPE, Topology.TYPE]:
            # get entity and entity id
            entity = ctx.get_entity(event)
            entity_id = ctx.get_entity_id(entity)
            elt_id = ctx.get_name(entity_id)
            logger.debug(u"elt_id {0}".format(elt_id))
            # process all targets
            elt = tm.get_elts(ids=elt_id)
            if elt is not None:
                targets = tm.get_targets(ids=elt_id)
                logger.debug(u"targets {0}".format(targets))
                # process and save all targets
                for target in targets:
                    target.process(
                        event=event, publisher=engine.amqp,
                        manager=tm, source=elt_id,
                        logger=logger,
                        **kwargs
                    )

        else:  # in case of entity event
            # get entity and entity id
            entity = ctx.get_entity(event)
            if entity is not None:

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

                entity_id = ctx.get_entity_id(encoded_entity)
                logger.debug(u"entity_id {0}".format(entity_id))
                elts = tm.get_elts(info={TopoNode.ENTITY: entity_id})
                logger.debug(u"elts {0}".format(elts))
                # process all elts
                for elt in elts:
                    elt.process(
                        event=event, publisher=engine.amqp,
                        manager=tm, logger=logger,
                        **kwargs
                    )

    return event
