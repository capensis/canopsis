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

"""
Module in charge of defining main topological rule actions.

A topological vertice has at least one of four actions in charge of changing
of state::

    - ``change_state``: change of state related to an input or event state.
    - ``state_from_sources``: change of state related to source nodes.
    - ``best_state``: change of state related to the best source node state.
    - ``worst_state``: change of state related to the worst source node state.
"""

from canopsis.task.core import register_task
from canopsis.common.init import basestring
from canopsis.common.utils import lookup
from canopsis.topology.manager import TopologyManager
from canopsis.topology.rule.condition import SOURCES_BY_EDGES
from canopsis.check import Check
from canopsis.check.manager import CheckManager

#: default topology manager
tm = TopologyManager()
#: default check manager
check = CheckManager()


@register_task
def change_state(
    event, vertice,
    state=None, update_entity=False, criticity=CheckManager.HARD,
    check_manager=None,
    **kwargs
):
    """
    Change of state on node and propagate the change of state on bound entity
        if necessary.

    :param event: event to process in order to change of state.
    :param vertice: vertice to change of state.
    :param state: new state to apply on input vertice. If None, get state from
        input event.
    :param bool update_entity: update entity state if True (False by default).
        The topology graph may have this flag to True.
    :param int criticity: criticity level. Default HARD.
    """

    # if state is None, use event state
    if state is None:
        state = event.get(Check.STATE, Check.OK)
    # update vertice state from ctx
    vertice.state = state

    # update entity if necessary
    if update_entity:
        entity = vertice.entity
        if entity is not None:
            # init check manager
            if check_manager is None:
                check_manager = check
            check_manager.state(ids=entity, state=state, criticity=criticity)


@register_task
def state_from_sources(event, vertice, ctx, f, manager=None, *args, **kwargs):
    """
    Change ctx vertice state which equals to f result on source nodes.
    """

    # get function f
    if isinstance(f, basestring):
        f = lookup(f)
    # init manager
    if manager is None:
        manager = tm

    # if sources are in ctx, get them
    if SOURCES_BY_EDGES in ctx:
        sources_by_edges = ctx[SOURCES_BY_EDGES]
    else:  # else get them with the topology object
        sources_by_edges = manager.get_sources(ids=vertice.id, add_edges=True)

    if sources_by_edges:  # do something only if sources exist
        # calculate the state
        sources = []
        for edge_id in sources_by_edges:
            _, edge_sources = sources_by_edges[edge_id]
            sources += edge_sources

        if sources:  # if sources exist, check state
            state = f(
                source_node.state
                for source_node in sources
            )
        else:  # else get OK
            state = Check.OK

        # change state
        change_state(
            state=state, event=event, vertice=vertice, ctx=ctx,
            *args, **kwargs
        )


@register_task
def worst_state(**kwargs):
    """
    Check the worst state among source nodes.
    """

    state_from_sources(f=max, **kwargs)


@register_task
def best_state(**kwargs):
    """
    Get the best state among source nodes.
    """

    state_from_sources(f=min, **kwargs)
