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

Three types of vertices exist in topology::

- cluster: operation between vertice states.
- node: vertice bound to an entity, like components, resources, etc.

A topology vertice contains::

- state: vertice state which change at runtime depending on bound entity state
    and event propagation.

A topology edge contains::

- weight: vertice weight in the graph.

"""

from canopsis.graph.elements import Graph, Vertice
from canopsis.task import run_task, TASK


class Topology(Graph):

    TYPE = 'topology'  #: topology type name

    __slots__ = Graph.__slots__


class Node(Vertice):
    """
    Class representation of a topology node.

    Contains:
        - (optionnally) an entity id.
        - a state (integer greater or equal to 0).
        - a weight information.
        - (optionnally) a task information.
    """

    ENTITY = 'entity'  #: entity data name
    STATE = 'state'  #: state data name
    WEIGHT = 'weight'  #: weight data name
    TASK = 'task'  #: task data name

    DEFAULT_WEIGHT = 1  #: default weight value
    DEFAULT_STATE = 0  #: default state value

    PARAM = 'node'  #: parameter name

    __slots__ = Vertice.__slots__

    def __init__(
        self,
        entity=None, state=DEFAULT_STATE, task=None, weight=DEFAULT_WEIGHT,
        *args, **kwargs
    ):
        """
        :param str entity: bound entity id.
        :param int state: state to use.
        :param task: task configuration.
        :type task: dict or str
        :param float weight: node weight.
        """

        super(Node, self).__init__(*args, **kwargs)

        if self.data is None:
            self.data = {}

        if entity is not None:
            self.data[Node.ENTITY] = entity

        if state is not None:
            self.data[Node.STATE] = state

        if task is not None:
            self.data[Node.TASK] = task

        if weight is not None:
            self.data[Node.WEIGHT] = weight

    def process(self, **kwargs):
        """
        Process this node task.
        """

        result = None

        task = self.data[TASK] if TASK in self.data else None

        if task is not None:
            kwargs[Node.PARAM] = self
            result = run_task(conf=task, **kwargs)

        return result
