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

A topology operation (to) contains:

- info.state: to state which change at runtime depending on bound entity
    state and event propagation.
- info.entity: to entity.
- info.operation: to operation.

A topology node inherits from both vertice and to.

A topology edge contains:

- weight: node weight in the graph related to edge targets.

A topology inherits from both graph and to and contains.
"""

__all__ = ['Topology', 'TopoEdge', 'TopoNode', 'TopoVertice']

from canopsis.graph.elements import Graph, Vertice, Edge
from canopsis.task.core import new_conf
from canopsis.check import Check
from canopsis.check.manager import CheckManager
from canopsis.context.manager import Context
from canopsis.topology.manager import TopologyManager
from canopsis.graph.event import BaseTaskedVertice
from canopsis.engines.core import publish
from canopsis.common.utils import singleton_per_scope


class TopoVertice(BaseTaskedVertice):

    STATE = Check.STATE  #: state field name in info.
    ENTITY = 'entity'  #: entity field name in info.
    OPERATION = 'operation'  #: operation field name in info.
    NAME = 'name'  #: element name.

    DEFAULT_STATE = Check.OK  #: default state value
    #: default task
    DEFAULT_TASK = 'canopsis.topology.rule.action.change_state'

    def get_default_task(self):
        """Get default task.
        """

        return new_conf(self.DEFAULT_TASK)

    def set_entity(self, entity_id, *args, **kwargs):

        super(TopoVertice, self).set_entity(
            entity_id=entity_id, *args, **kwargs
        )
        # update entity state
        if entity_id is not None:
            cm = singleton_per_scope(CheckManager)
            state = cm.state(ids=entity_id)
            if state is None:
                state = TopoVertice.DEFAULT_STATE
            self.info[TopoVertice.STATE] = state

    @property
    def state(self):

        result = self.info.get(
            TopoVertice.STATE, TopoVertice.DEFAULT_STATE
        )
        return result

    @state.setter
    def state(self, value):

        if value is not None:
            self.info[TopoVertice.STATE] = value

    @property
    def operation(self):

        result = self.task
        return result

    @operation.setter
    def operation(self, value):

        self.task = value

    def get_event(self, state=DEFAULT_STATE, source=None, *args, **kwargs):

        result = super(TopoVertice, self).get_event(
            state=state, source=source,
            *args, **kwargs
        )

        return result

    def process(
            self, event, publisher=None, manager=None, source=None,
            logger=None,
            **kwargs
    ):

        """
        :param TopologyManager manager:
        """

        if manager is None:
            manager = singleton_per_scope(TopologyManager)

        # save old state
        old_state = self.state
        # process task
        result = super(TopoVertice, self).process(
            event=event, publisher=publisher, manager=manager, source=source,
            logger=logger,
            **kwargs
        )
        # compare old state and new state
        if self.state != old_state:
            # update edges
            targets_by_edge = manager.get_targets(
                ids=self.id, add_edges=True
            )
            for edge_id in targets_by_edge:
                edge, _ = targets_by_edge[edge_id]
                # update edge state
                edge.state = self.state
                edge.save(manager=manager)
            # if not equal
            new_event = self.get_event(state=self.state, source=source)
            # publish a new event
            if publisher is not None:
                publish(event=new_event, publisher=publisher)
            # save self
            self.save(manager=manager)

        return result


class Topology(Graph, TopoVertice):

    TYPE = 'topo'  #: topology type name

    DEFAULT_TASK = 'canopsis.topology.rule.action.worst_state'

    __slots__ = Graph.__slots__

    def __init__(
            self,
            operation=None, state=None, type=TYPE,
            entity=None, *args, **kwargs
    ):

        super(Topology, self).__init__(type=type, *args, **kwargs)

        # set info
        if self.info is None:
            self.info = {}
        # set operation
        if operation is not None:
            self.operation = operation
        # set state
        self.state = state
        # set entity
        self.entity = entity

    def set_entity(self, entity_id, *args, **kwargs):

        super(Topology, self).set_entity(entity_id=entity_id, *args, **kwargs)

        # set default entity if entity_id is None
        if entity_id is None and self.entity is None:
            # set entity
            ctxm = singleton_per_scope(Context)
            event = self.get_event(source=0, state=0)
            entity = ctxm.get_entity(event)
            entity_id = ctxm.get_entity_id(entity)
            self.entity = entity_id

    def save(self, context=None, *args, **kwargs):

        super(Topology, self).save(*args, **kwargs)

        # use global context if input context is None
        if context is None:
            context = singleton_per_scope(Context)
        # get self entity
        event = self.get_event()
        entity = context.get_entity(event)
        ctx, _id = context.get_entity_context_and_name(entity=entity)
        entity = {Context.NAME: _id}
        # put the topology in the context by default
        context.put(_type=self.type, entity=entity, context=ctx)


class TopoNode(Vertice, TopoVertice):
    """Class representation of a topology node.

    Contains:
        - (optionnally) an entity id.
        - an entity state.
        - a weight information.
        - (optionnally) a task information.
    """

    TYPE = 'toponode'  #: node type name.

    PARAM = 'toponode'  #: parameter name.

    __slots__ = Vertice.__slots__

    def __init__(
            self,
            entity=None, state=None, operation=None,
            *args, **kwargs
    ):
        """
        :param int state: state to use.
        :param float weight: node weight.
        """

        super(TopoNode, self).__init__(*args, **kwargs)

        # set info
        if self.info is None:
            self.info = {}
        # set state
        self.state = state
        # set entity
        self.entity = entity
        # set operation
        if operation is not None:
            self.operation = operation

    def get_event(self, *args, **kwargs):

        result = super(TopoNode, self).get_event(*args, **kwargs)

        tm = singleton_per_scope(TopologyManager)
        graphs = tm.get_graphs(elts=self.id)
        # iterate on existing graphs
        for graph in graphs:
            # update result as soon as a graph has been founded
            result['component'] = graph.id
            break
        result['resource'] = self.id

        return result

    def delete(self, manager, cache=False, *args, **kwargs):

        super(TopoNode, self).delete(
            manager=manager, cache=cache, *args, **kwargs
        )

        # delete edges where source is self
        edges = manager.get_edges(sources=self.id)
        for edge in edges:
            edge.delete(manager=manager, cache=cache)

        # delete reference of self in edges
        manager.del_edge_refs(targets=self.id, del_empty=True, cache=cache)


class TopoEdge(Edge, TopoVertice):
    """Topology edge.
    """

    __slots__ = Edge.__slots__

    TYPE = 'topoedge'  #: topology edge type name

    def __init__(self, state=None, *args, **kwargs):

        super(TopoEdge, self).__init__(*args, **kwargs)
        # set info
        if self.info is None:
            self.info = {}
        # set state
        self.state = state
