# -*- coding: utf-8 -*-
# --------------------------------
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
Module in charge of defining main topological rule actions.

A topological node has at least one of three rule action in charge of changing
of state::

    - ``change_state``: change of state related to an input or event state.
    - ``best_state``: change of state related to the best source node state.
    - ``worst_state``: change of state related to the worst source node state.

For logical reasons, the propagate action runned such as the last action.
"""

from canopsis.task import register_task
from canopsis.common.init import basestring
from canopsis.common.utils import lookup
from canopsis.topology.manager import TopologyManager
from canopsis.check import Check
from canopsis.topology.process import SOURCES

tm = TopologyManager()


@register_task
def change_state(event, node, state=None, manager=None, **kwargs):
    """
    Change of state for node.

    :param event: event to process in order to change of state.
    :param node: node to change of state.
    :param state: new state to apply on input node. If None, get state from
        input event.
    """

    # if state is None, use event state
    if state is None:
        state = event[Check.STATE]

    if manager is None:
        manager = tm

    # update node state from ctx
    node.data[Check.STATE] = state
    node.save(manager=manager)


@register_task
def state_from_sources(event, node, ctx, f, manager=None, **kwargs):
    """
    Change ctx node state which equals to f result on source nodes.
    """

    # get function f
    if isinstance(f, basestring):
        f = lookup(f)

    if manager is None:
        manager = tm

    # if sources are in ctx, get them
    if SOURCES in ctx:
        sources = ctx[SOURCES]
    else:  # else get them with the topology object
        sources = manager.get_sources(ids=node.id)

    if sources:  # do something only if sources exist
        # calculate the state
        state = f(source_node.data[Check.STATE] for source_node in sources)

        # update the node state
        node.data[Check.STATE] = state
        node.save(manager=manager)


@register_task
def worst_state(event, ctx, manager=None, **kwargs):
    """
    Check the worst state among source nodes.
    """

    state_from_sources(event=event, ctx=ctx, f=max, manager=manager, **kwargs)


@register_task
def best_state(event, ctx, manager=None, **kwargs):
    """
    Get the best state among source nodes.
    """

    state_from_sources(event=event, ctx=ctx, f=min, manager=manager, **kwargs)
