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
from canopsis.topology.manager import TopologyManager
from canopsis.event import forger, Event

_context = Context()
_check = CheckManager()
_topology = TopologyManager()


class TopoVertice(object):

    STATE = Check.STATE  #: state field name in data
    ENTITY = 'entity'  #: entity field name in data
    OPERATOR = 'operator'  #: operator field name in data
    NAME = 'name'  #: element name.

    DEFAULT_STATE = Check.OK  #: default state value

    @property
    def entity(self):
        """
        Get self entity id.
        :return: self entity id.
        :rtype: str
        """
        return self.data[TopoVertice.ENTITY]

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
        self.data[TopoVertice.ENTITY] = entity_id
        # update entity state
        self.data[TopoVertice.STATE] = _check.state(ids=entity_id)

    def get_context_w_entity(self):
        """
        Get self entity structure and its context.

        :return: tuple of self context and entity.
        :rtype: tuple
        """

        context = {
            'connector': Event.CONNECTOR,
            'connector_name': Event.CONNECTOR_NAME,
            'component': self.id
        }

        entity = {
            Context.NAME: self.name
        }

        return context, entity

    @property
    def name(self):
        return self.data.get(TopoVertice.NAME, self.id)

    @name.setter
    def name(self, value):
        self.data[TopoVertice.NAME] = value

    @property
    def state(self):
        result = self.data.get(TopoVertice.STATE)
        return result

    @state.setter
    def state(self, value):
        self.data[TopoVertice.STATE] = value

    @property
    def operator(self):
        result = self.data.get(TopoVertice.OPERATOR)
        return result

    @operator.setter
    def operator(self, value):
        self.data[TopoVertice.OPERATOR] = value

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

    def get_event(self, state, source):
        """
        Get topo element event.
        :param int state: new state to apply.
        """

        result = forger(
            event_type=self.type,
            component=self.id if self.type == Topology.TYPE else None,
            resource=self.id if self.type == TopoNode.TYPE else None,
            state=state,
            state_type=1,
            id=self.id,
            source=source
        )
        return result


class Topology(Graph, TopoVertice):

    TYPE = 'topo'  #: topology type name

    __slots__ = Graph.__slots__

    def __init__(
        self,
        operator=None, state=TopoVertice.DEFAULT_STATE, _type=TYPE,
        *args, **kwargs
    ):

        super(Topology, self).__init__(_type=_type, *args, **kwargs)

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
        # set operator
        if operator is not None:
            self.operator = operator
        # set state
        self.state = state

    @property
    def entity(self):
        """
        Get self entity id.
        :return: self entity id.
        :rtype: str
        """

        return self.data[TopoVertice.ENTITY]

    def save(self, context=None, *args, **kwargs):

        super(Topology, self).save(*args, **kwargs)

        # use global context if input context is None
        if context is None:
            context = _context
        # get self entity
        ctx, entity = self.get_context_w_entity()
        # put the topology in the context by default
        context.put(_type=self.type, entity=entity, context=ctx)


class TopoNode(Vertice, TopoVertice):
    """
    Class representation of a topology node.

    Contains:
        - (optionnally) an entity id.
        - an entity state.
        - a weight information.
        - (optionnally) a task information.
    """

    TYPE = 'toponode'  #: node type name

    PARAM = 'toponode'  #: parameter name

    __slots__ = Vertice.__slots__

    def __init__(
        self,
        entity=None, state=TopoVertice.DEFAULT_STATE, operator=None,
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
        # init data
        if self.data is None:
            self.data = {}
        # set entity
        if entity is not None:
            self.entity = entity
        # set state
        self.state = state
        # set operator
        if operator is not None:
            self.operator = operator

    def get_event(self, state, source, *args, **kwargs):
        """
        Get topo element event.
        :param int state: new state to apply.
        """

        result = super(TopoNode, self).get_event(
            state, source, *args, **kwargs
        )

        # get topology id
        topologies = _topology.get_graphs(elts=self.id)
        if topologies:
            topology = topologies[0]
            # update component field
            result['component'] = topology.id

        return result

    def get_context_w_entity(self, *args, **kwargs):

        ctx, entity = super(TopoNode, self).get_context_w_entity(
            *args, **kwargs
        )

        ctx['resource'] = self.id
        ctx['component'] = None

        # get topology id
        topologies = _topology.get_graphs(elts=self.id)
        if topologies:
            topology = topologies[0]
            # update component field
            ctx['component'] = topology.id

        return ctx, entity


class TopoEdge(Edge, TopoVertice):
    """
    Topology edge.
    """

    __slots__ = Edge.__slots__

    TYPE = 'topoedge'  #: topology edge type name

    def __init__(self, *args, **kwargs):

        super(TopoEdge, self).__init__(*args, **kwargs)

        if self.data is None:
            self.data = {}
