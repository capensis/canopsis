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
Module in charge of defining main topological rule actions.

A topological node has at least one of three rule action in charge of changing
of state::

    - ``change_state``: change of state related to an input or event state.
    - ``best_state``: change of state related to the best source node state.
    - ``worst_state``: change of state related to the worst source node state.

For logical reasons, the propagate action runned such as the last action.
"""

from time import sleep

from canopsis.rule import register_task
from canopsis.common.utils import lookup
from canopsis.topology.manager import Topology
from canopsis.check import Check
from canopsis.topology.process import SOURCES, NODE

topology = Topology()


@register_task(name='topology.change_state')
def change_state(event, ctx, state=None, **kwargs):
    """
    Change of state for node ctx.
    """

    # if state is None, use event state
    if state is None:
        state = event[Check.STATE]

    # update node state from ctx
    node = ctx[NODE]
    node[Check.STATE] = state
    topology.push_node(node)


@register_task(name='topology.worst_state')
def worst_state(event, ctx, **kwargs):
    """
    Check the worst state among source nodes.
    """

    change_state_from_source_nodes(event=event, ctx=ctx, f=max)


@register_task(name='topology.best_state')
def best_state(event, ctx, **kwargs):
    """
    Get the best state among source nodes.
    """

    change_state_from_source_nodes(event=event, ctx=ctx, f=min)


@register_task(name='topology.change_state_from_source_nodes')
def change_state_from_source_nodes(event, ctx, f, **kwargs):
    """
    Change ctx node state which equals to f result on source nodes.
    """

    # get function f
    if isinstance(f, basestring):
        f = lookup(f)

    # retrieve the node from where find source nodes
    node = ctx[NODE]

    # if sources are in ctx, get them
    if SOURCES in ctx:
        sources = ctx[SOURCES]
    else:  # else get them with the topology object
        sources = topology.find_source_nodes(node=node)

    # calculate the state
    state = f(source_node[Check.STATE] for source_node in sources)

    # update the node state
    node[Check.STATE] = state
    topology.push_node(node=node)
