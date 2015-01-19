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
Description
===========

Functional
----------

A topology is a graph dedicated to enriches status model of entities with state
dependency between entities.

Topological tasks consist to update status vertice information and to propagate
the change of state in sending check events.

vertices could be finally connected to the topology in order to propagate all
 change of state to the topology itelf.

An example of application is root cause analysis where a topology may react
when an entity change of state and can propagate over topology nodes the change
of state in some propagation conditions such as operations like ``worst state``

Topology tasks are commonly rules of (condition, actions). A condition takes in
parameter the execution context of an event, the engine and the vertice which
embeds the rule.

Technical
---------

A topology node can be bound to an entity or/and contain a task.

Both permits to update the node state. The first one will update its state
related to the bound entity state, while the task can update the state
independently to the entity state.

A topology node contains::

- state: node state which change at runtime depending on bound entity state
    and event propagation.

A topology edge contains::

- weight: node weight in the graph related to edge targets.

A topology contains::

- state: current topology state.
- data.task: the change_state operator with a propagation to the entity.

"""

__all__ = ['Topology', 'TopoEdge', 'TopoNode']

from canopsis.graph.elements import Graph, Vertice, Edge
from canopsis.task import run_task, TASK
from canopsis.check import Check
from canopsis.check.manager import CheckManager
from canopsis.context.manager import Context

_context = Context()
_check = CheckManager()


class Topology(Graph):

    TYPE = 'topology'  #: topology type name

    __slots__ = Graph.__slots__

    CONNECTOR = 'canopsis'
    CONNECTOR_NAME = 'canopsis'
    COMPONENT = 'canopsis'

    def __init__(self, *args, **kwargs):

        super(Topology, self).__init__(*args, **kwargs)

        if self.data is None:
            self.data = {}
        # ensure entity id is in topology
        if TopoNode.ENTITY not in self.data:
            entity_id = self.entity_id()
            self.data[TopoNode.ENTITY] = entity_id
        # ensure state is in data
        if Check.STATE not in self.data:
            entity_id = self.data[TopoNode.ENTITY]
            self.data[TASK] = _check.state(ids=entity_id)

    def entity_id(self):
        """
        Get self entity id.
        """

        ctx, entity = self.get_context_w_entity()
        entity.update(ctx)
        result = _context.get_entity_id(entity)
        return result

    def get_context_w_entity(self):
        """
        Get self entity structure and its context.

        :return: tuple of self context and entity.
        :rtype: tuple
        """

        context = {
            'connector': Topology.CONNECTOR,
            'connector_name': Topology.CONNECTOR_NAME,
            'component': Topology.COMPONENT
        }

        entity = {
            Context.NAME: self.id
        }

        return context, entity

    def save(self, context=None, *args, **kwargs):

        super(Topology, self).save(*args, **kwargs)

        # use global context if input context is None
        if context is None:
            context = _context
        # get self entity
        ctx, entity = self.get_context_w_entity()
        # put the topology in the context by default
        context.put(_type=Topology.TYPE, entity=entity, context=ctx)


class TopoNode(Vertice):
    """
    Class representation of a topology node.

    Contains:
        - (optionnally) an entity id.
        - an entity state.
        - a weight information.
        - (optionnally) a task information.
    """

    TYPE = 'topovertice'  #: node type name

    ENTITY = 'entity'  #: entity data name

    DEFAULT_STATE = Check.OK  #: default state value

    PARAM = 'toponode'  #: parameter name

    __slots__ = Vertice.__slots__

    def __init__(
        self,
        entity=None, state=DEFAULT_STATE, task=None,
        *args, **kwargs
    ):
        """
        :param str entity: bound entity id.
        :param int state: state to use.
        :param task: task configuration.
        :type task: dict or str
        :param float weight: node weight.
        """

        super(TopoNode, self).__init__(*args, **kwargs)

        if self.data is None:
            self.data = {}

        if entity is not None:
            self.data[TopoNode.ENTITY] = entity

        if state is not None:
            self.data[Check.STATE] = state

        if task is not None:
            self.data[TASK] = task

    @staticmethod
    def task(elt, value=None):
        """
        Get task value. If value is not None, update task value.

        :param GraphElement elt: from where get task.
        :param value: new task to use.
        :type value: str or dict
        :return: elt task.
        :rtype: str or dict
        """

        if value is not None:
            elt.data[TASK] = value
        return elt.data[TASK] if TASK in elt.data else None

    @staticmethod
    def entity(elt, value=None):
        """
        Get task entity id. If value is not None, update entity value.

        :param GraphElement elt: from where get entity.
        :param str value: new entity to use.
        :return: elt entity.
        :rtype: str
        """

        if value is not None:
            elt.data[TopoNode.ENTITY] = value

        result = None

        if TopoNode.ENTITY in elt.data:
            result = elt.data[TopoNode.ENTITY]

        return result

    @staticmethod
    def state(elt, value=None):
        """
        Get state value. If value is not None, update state value.

        :param GraphElement elt: from where get state.
        :param float value: new state to use.
        :return: elt state.
        :rtype: float
        """

        if value is not None:
            elt.data[Check.STATE] = value
        return elt.data[Check.STATE] if Check.STATE in elt.data else Check.OK

    @staticmethod
    def process(toponode, event, **kwargs):
        """
        Process this node task.
        """

        result = None

        task = TopoNode.task(toponode)

        if task is not None:
            result = run_task(
                conf=task, toponode=toponode, event=event, **kwargs
            )

        return result


class TopoEdge(Edge):
    """
    Topology edge.
    """

    __slots__ = Edge.__slots__

    TYPE = 'topoedge'  #: topology edge type name

    def __init__(self, data=None, *args, **kwargs):

        super(TopoEdge, self).__init__(
            data={} if data is None else data,
            *args, **kwargs
        )
