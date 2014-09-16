# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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
from canopsis.event import Event
from canopsis.rule import process_rule

context = Context()
topology = Topology()

SOURCE = 'source'
SOURCES = 'sources'
NODE = 'node'
STATE = 'state'
RULE = 'rule'


def event_processing(event, ctx=None, **params):
    """
    Process input event in getting topology nodes bound to input event entity.

    One topology nodes are founded, executing related rules.
    """

    # TODO: remove when Check event will be used
    if event['event_type'] == Event.CHECK:

        # get nodes
        entity = context.get_entity(event)
        entity_id = context.get_entity_id(entity)
        nodes = topology.find_bound_nodes(entity_id=entity_id)

        # iterate on bound nodes
        for node in nodes:

            rule = node[RULE]
            ctx = {NODE: node}

            process_rule(event=event, rule=rule, ctx=ctx)

    return event
