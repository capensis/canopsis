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

A topology operator (to) contains::

- data.state: to state which change at runtime depending on bound entity
    state and event propagation.
- data.entity: to entity.
- data.operator: to operator.

A topology node inherits from both vertice and to.

A topology edge contains::

- weight: node weight in the graph related to edge targets.

A topology inherits from both grapg and to and contains.
"""

__all__ = ['Topology', 'TopoEdge', 'TopoNode']

from canopsis.graph.elements import Graph, Vertice, Edge
from canopsis.task import run_task, TASK, new_conf
from canopsis.check import Check
from canopsis.check.manager import CheckManager
from canopsis.context.manager import Context

_context = Context()
_check = CheckManager()


class TopoOp(object):

    STATE = Check.STATE  #: state field name in data
    ENTITY = 'entity'  #: entity field name in data
    OPERATOR = 'operator'  #: operator field name in data

    @property
    def entity(self):
        """
        Get self entity id.
        :return: self entity id.
        :rtype: str
        """
        return self.data[TopoOp.ENTITY]

    @entity.setter
    def entity(self, value):
        """
        Change of entity id and update state.

        :param value: new entity (id) to use.
        :type value: dict or str
        """

        if isinstance(value, dict):
            # get entity id
            entity_id = _context.get_entity_id(value)
        else:
            entity_id = value

        # update entity
        self.data[TopoOp.ENTITY] = entity_id
        # update entity state
        self.data[TopoOp.STATE] = _check.state(ids=entity_id)

    @property
    def state(self):
        result = self.data.get(TopoOp.STATE)
        return result

    @state.setter
    def state(self, value):
        self.data[TopoOp.STATE] = value

    @property
    def operator(self):
        result = self.data.get(TopoOp.OPERATOR)
        return result

    @operator.setter
    def operator(self, value):
        self.data[TopoOp.OPERATOR] = value

    def process(self, event, **kwargs):
        """
        Process this node task.
        """

        result = None

        task = self.operator

        if task is not None:
            result = run_task(
                conf=task, toponode=self, event=event, **kwargs
            )

        return result


class Topology(Graph, TopoOp):

    TYPE = 'topology'  #: topology type name

    __slots__ = Graph.__slots__

    CONNECTOR = 'canopsis'
    CONNECTOR_NAME = 'canopsis'
    COMPONENT = 'canopsis'

    def __init__(self, *args, **kwargs):

        super(Topology, self).__init__(*args, **kwargs)

        if self.data is None:
            self.data = {}

        # ensure change state is the default task
        if TASK not in self.data:
            self.data[TASK] = new_conf(
                'canopsis.topology.rule.action.change_state',
                **{
                    'update_entity': True
                }
            )

    @property
    def entity(self):
        """
        Get self entity id.
        :return: self entity id.
        :rtype: str
        """

        return self.data[TopoOp.ENTITY]

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
            Context.NAME: self.name
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
        context.put(_type=self.type, entity=entity, context=ctx)


class TopoNode(Vertice, TopoOp):
    """
    Class representation of a topology node.

    Contains:
        - (optionnally) an entity id.
        - an entity state.
        - a weight information.
        - (optionnally) a task information.
    """

    TYPE = 'toponode'  #: node type name

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


class TopoEdge(Edge, TopoOp):
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
