# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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
Module in charge of defining main topological rules.

When an event occured, the related entity is retrieved and all bound
topological nodes are retrieved as well in order to execute theirs rules.

First, a topology processing is triggered when an event occured.

From this event, bound topology nodes are got in order to apply node rules.

A typical topological rule condition is an ``canopsis.rule.condition.all``
composed of the ``canopsis.topology.process.new_state`` condition.
If this condition is checked, then other specific conditions can be applied
such as::
    - ``canopsis.topology.process.condition``
    - ``canopsis.topology.process.all``
    - ``canopsis.topology.process.any``

A topology rule use the condition ``new_state``
"""

from canopsis.topology.manager import Topology
from canopsis.context.manager import Context
from canopsis.task import process_rule, RULE, register_task
from canopsis.event import Event
from canopsis.check import Check


context = Context()
topology = Topology()

SOURCE = 'source'
SOURCES = 'sources'
NODE = 'node'
PUBLISHER = 'publisher'
WEIGHT = 'weight'


@register_task(name='topology.event_processing')
def event_processing(event, ctx=None, **params):
    """
    Process input event in getting topology nodes bound to input event entity.

    One topology nodes are founded, executing related rules.
    """

    event_type = event[Event.TYPE]

    nodes = []

    if ctx is None:
        ctx = {}

    # TODO: remove when Check event will be used
    # apply processing only in case of check event
    if event_type == Check.get_type():

        source_type = event[Event.SOURCE_TYPE]

        # in case of topology node
        if source_type == Topology.TOPOLOGY_NODE_TYPE:
             # get nodes from the event topology node id
            nodes = [topology.get_nodes(event[Topology.ID])]

        else:  # in case of entity event
            # get nodes from entity
            entity = context.get_entity(event)
            if entity is not None:
                entity_id = context.get_entity_id(entity)
                nodes = topology.find_bound_nodes(entity_id=entity_id)

        # iterate on nodes
        for node in nodes:

            # add node in the ctx
            ctx[NODE] = node

            # save old state in order to check for its modification
            if Check.STATE in node:
                old_state = node[Check.STATE]
            else:
                old_state = event[Check.STATE]

            # process rule
            if RULE in node:
                rule = node[RULE]
                process_rule(event=event, rule=rule, ctx=ctx)
                # propagate the change of state in case of new state
                if old_state != node[Check.STATE]:

                    # get next nodes
                    next_nodes = topology.get_next_nodes(node)

                    # send the event_to_propagate to all next_nodes
                    for next_node in next_nodes:

                        # create event to propagate with source and node ids
                        event_to_propagate = {
                            Event.TYPE: Check.get_type(),
                            Event.SOURCE_TYPE: Topology.TOPOLOGY_NODE_TYPE,
                            Topology.ID: next_node[Topology.ID],
                            SOURCE: node[Topology.ID]
                        }

                        ctx[PUBLISHER].publish(event_to_propagate)

    return event
