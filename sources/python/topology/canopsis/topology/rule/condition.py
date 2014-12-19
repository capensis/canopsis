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

A topological node doesn't require a condition.

In addition to this condition, it is possible to test another condition which
refers to source nodes if they exist.

Such conditions are::
    - ``new_state``: test if state (or event state) is not equal to node state
    - ``condition``: test if source node states match with an input state.
    - ``all``: test if all source node states match with an input state.
    - ``any``: test if one of source node states match with an input state.

The ``new_state`` condition may be used by nodes bound to entities in order to
update such nodes when the entity change of state.

Related rule actions are defined in ``canopsis.topology.rule.action`` module.
"""

from canopsis.topology.manager import TopologyManager
from canopsis.topology.process import SOURCES, WEIGHT
from canopsis.check import Check
from canopsis.task import register_task

from sys import maxint

tm = TopologyManager()


@register_task
def new_state(event, node, state=None, **kwargs):
    """
    Condition triggered when state is different than node state.

    :param dict event: event from where get state if input state is None.
    :param Node node: node from where check if node state != input state.
    :param int state: state to compare with input node state.
    """

    # if state is None, use event state
    if state is None:
        state = event[Check.STATE]

    # True if node state is different than state
    result = node.data[Check.STATE] != state

    return result


@register_task
def at_least(
    event, ctx, node, min_weight=1, state=None, rrule=None, f=None, **kwargs
):
    """
    Generic condition applied on sources of node which check if at least source
        nodes check a condition.

    :param dict event: processed event.
    :param dict ctx: rule context which must contain rule node.
    :param Node node: node to check.
    :param int min_weight: minimal node weight.
    :param int state: state to check among sources nodes.
    :param rrule rrule: rrule to consider in order to check condition in time.
    :param f: function to apply on source node state. If None, use equality
        between input state and source node state.

    :return: True if condition is checked among source nodes.
    :rtype: bool
    """

    source_nodes = tm.get_sources(ids=node.id)

    weights = list(source_node.data[WEIGHT] for source_node in source_nodes)
    if weights:
        weights.append(min_weight)
        min_weight = min(weights)

    result = False

    for source_node in source_nodes:

        source_node_state = source_node.data[Check.STATE]

        if source_node_state == state if f is None else f(source_node_state):
            min_weight -= source_node.data[WEIGHT]

            if min_weight <= 0:
                result = True
                break

    # if result, save source_nodes in ctx in order to save read data from db
    if result:
        ctx[SOURCES] = source_nodes

    return result


@register_task
def _all(event, ctx, node, min_weight=1, state=None, rrule=None, **kwargs):
    """
    Check if all source nodes match with input check_state.

    :param dict event: processed event.
    :param dict ctx: rule context which must contain rule node.
    :param Node node: node to check.
    :param int min_weight: minimal node weight.
    :param int state: state to check among sources nodes.
    :param rrule rrule: rrule to consider in order to check condition in time.
    :param f: function to apply on source node state. If None, use equality
        between input state and source node state.

    :return: True if condition is checked among source nodes.
    :rtype: bool
    """

    result = at_least(
        event=event,
        ctx=ctx,
        node=node,
        min_weight=maxint,
        state=state,
        rrule=rrule,
        **kwargs
    )

    return result


@register_task
def nok(event, ctx, node, min_weight=1, rrule=None, **kwargs):
    """
    Condition which check if source nodes are not ok.

    :param dict event: processed event.
    :param dict ctx: rule context which must contain rule node.
    :param Node node: node to check.
    :param int min_weight: minimal node weight.
    :param int state: state to check among sources nodes.
    :param rrule rrule: rrule to consider in order to check condition in time.
    :param f: function to apply on source node state. If None, use equality
        between input state and source node state.

    :return: True if condition is checked among source nodes.
    :rtype: bool
    """

    return at_least(
        event=event,
        ctx=ctx,
        node=node,
        min_weight=min_weight,
        rrule=rrule,
        f=lambda x: x != 0,
        **kwargs
    )
