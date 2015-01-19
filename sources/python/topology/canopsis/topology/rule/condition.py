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
Module in charge of defining main topological rule conditions.

A topological toponode doesn't require a condition.

In addition to this condition, it is possible to test another condition which
refers to source nodes if they exist.

Such conditions are::
    - ``new_state``: test if state (or event state) is not equal to node state
    - ``at_least``: test if source node states match with an input state.
    - ``_all``: test if all source node states match with an input state.
    - ``nok``: test if source node states are not OK.

The ``new_state`` condition may be used by nodes bound to entities in order to
update such nodes when the entity change of state.

Related rule actions are defined in ``canopsis.topology.rule.action`` module.
"""

from canopsis.common.init import basestring
from canopsis.topology.manager import TopologyManager
from canopsis.topology.elements import TopoNode
from canopsis.check import Check
from canopsis.task import register_task

tm = TopologyManager()

#: parameter name which contain sources by edges
SOURCES_BY_EDGES = 'sources_by_edges'


@register_task
def new_state(event, toponode, state=None, manager=None, **kwargs):
    """
    Condition triggered when state is different than toponode state.

    :param dict event: event from where get state if input state is None.
    :param TopoNode toponode: toponode from where check if toponode state != input state.
    :param int state: state to compare with input toponode state.
    """

    if manager is None:
        manager = tm

    if isinstance(toponode, basestring):
        toponode = manager.get_elts(ids=toponode)

    # if state is None, use event state
    if state is None:
        state = event[Check.STATE]

    # True if toponode state is different than state
    result = toponode.state != state

    return result


@register_task
def at_least(
    event, ctx, toponode, state=Check.OK, min_weight=1, rrule=None, f=None,
    manager=None, **kwargs
):
    """
    Generic condition applied on sources of toponode which check if at least source
        nodes check a condition.

    :param dict event: processed event.
    :param dict ctx: rule context which must contain rule toponode.
    :param TopoNode toponode: toponode to check.
    :param int state: state to check among sources nodes.
    :param float min_weight: minimal weight (default 1) to reach in order to
        validate this condition. If None, condition results in checking all
            sources.
    :param rrule rrule: rrule to consider in order to check condition in time.
    :param f: function to apply on source toponode state. If None, use equality
        between input state and source toponode state.

    :return: True if condition is checked among source nodes.
    :rtype: bool
    """

    result = False

    if manager is None:
        manager = tm

    if isinstance(toponode, basestring):
        toponode = manager.get_elts(ids=toponode)

    sources_by_edges = manager.get_sources(ids=toponode.id, add_edges=True)

    if sources_by_edges and min_weight is None:
        # if edges & checking all nodes is required, result is True by default
        result = True

    # for all edges
    for edge_id in sources_by_edges:
        # get edge and sources
        edge, sources = sources_by_edges[edge_id]
        # get edge_weight which is 1 by default
        for source in sources:
            source_state = source.state
            if (source_state == state if f is None else f(source_state)):
                if min_weight is not None:  # if min_weight is not None
                    min_weight -= edge.weight  # remove edge_weight from result
                    if min_weight <= 0:  # if min_weight is negative, ends loop
                        result = True
                        break
            elif min_weight is None:
                # stop if condition is not checked and min_weight is None
                result = False
                break

    # if result, save source_nodes in ctx in order to save read data from db
    if result:
        ctx[SOURCES_BY_EDGES] = sources_by_edges

    return result


@register_task
def _all(**kwargs):
    """
    Check if all source nodes match with input check_state.

    :param dict event: processed event.
    :param dict ctx: rule context which must contain rule toponode.
    :param TopoNode toponode: toponode to check.
    :param int min_weight: minimal toponode weight to check.
    :param int state: state to check among sources nodes.
    :param rrule rrule: rrule to consider in order to check condition in time.
    :param f: function to apply on source toponode state. If None, use equality
        between input state and source toponode state.

    :return: True if condition is checked among source nodes.
    :rtype: bool
    """

    result = at_least(
        min_weight=None,
        **kwargs
    )

    return result


@register_task
def nok(**kwargs):
    """
    Condition which check if source nodes are not ok.

    :param dict event: processed event.
    :param dict ctx: rule context which must contain rule toponode.
    :param TopoNode toponode: toponode to check.
    :param int min_weight: minimal toponode weight to check.
    :param int state: state to check among sources nodes.
    :param rrule rrule: rrule to consider in order to check condition in time.
    :param f: function to apply on source toponode state. If None, use equality
        between input state and source toponode state.

    :return: True if condition is checked among source nodes.
    :rtype: bool
    """

    return at_least(
        f=lambda x: x != Check.OK,
        **kwargs
    )
