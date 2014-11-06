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

Topological tasks consist to update status node information and propagate the
change of state in sending check events.

A topology has at least a root node which is bound to a topology entity.

This root is the only node where the state is the same as the entity state.

An example of application is root cause analysis where a topology may react
when an entity change of state and can propagate over topology nodes the change
of state in some propagation conditions such as operations like ``worst state``

Topology tasks are commonly rules of (condition, actions). A condition takes in
parameter the execution context of an event, the engine and the node which
embeds the rule.

Technical
---------

Three types of nodes exist in topology::

- cluster: operation between node states.
- node: node bound to an entity, like components, resources, and the root.
- root: root node.

A topology data contains several thinks among::

- weight: node weight in the graph.
- state: node state which change at runtime depending on bound entity state and
    event propagation.
"""

from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.graph.manager import GraphManager

CONF_PATH = 'topology/topology.conf'
CATEGORY = 'TOPOLOGY'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class TopologyManager(GraphManager):
    """
    Manage topological graph.
    """

    GRAPH_TYPE = 'topology'  #: topology graph type name

    WEIGHT = 'weight'  #: weight node field name
    STATE = 'state'  #: state node field name

    DEFAULT_WEIGHT = 1  #: default weight
    DEFAULT_STATE = 0  #: default state

    ROOT = 'root'  #: root node type name
    CLUSTER = 'cluster'  #: cluster node type name

    def put_graph(self, graph, *args, **kwargs):
        """
        Put graph in DB and ensuring a root is bound to graph.

        :param dict graph: graph to put in DB.
        :return: putted graph.
        """

        result = super(TopologyManager, self).put_graph(graph, *args, **kwargs)

        graph_id = result[GraphManager.ID]

        nodes = result[GraphManager.NODES]

        # get is_root class method for cache reasons
        is_root = TopologyManager.is_root
        # find a root in nodes
        for node in nodes:
            # if a root is found
            if is_root(node):
                # compare root entity_id with graph_id
                entity_id = node[GraphManager.ENTITY_ID]
                if entity_id != graph_id:
                    node[GraphManager.ENTITY_ID] = graph_id
                    # in case of a difference, put the root with graph_id
                    self.put_nodes(nodes=node)
                break
        else:  # if no root exist in topology
            # create it
            root = TopologyManager.new_node(
                graph_id=graph_id,
                entity_id=graph_id,
                _type=TopologyManager.ROOT
            )
            # and put it in DB
            self.put_nodes(nodes=root)
            nodes.insert(0, root)

        return result

    @classmethod
    def new_node(
        cls,
        data=None, weight=None, state=None, *args, **kwargs
    ):
        """
        Apply GraphManager.new_node in adding weight and state in node data.

        :param dict data: node data.
        :param float weight: node weight (default TopologyManager.WEIGHT = 1).
        :param int state: node state (default TopologyManager.STATE = 0).
        :param args: GraphManager.new_node varargs
        :param kwargs: GraphManager.new_node kwargs
        :return: GraphManager.new_node result with weight and state in node
            data.
        :rtype: dict
        """

        # initialize weight and state
        if weight is None:
            weight = TopologyManager.DEFAULT_WEIGHT
        if state is None:
            state = TopologyManager.DEFAULT_STATE

        if data is None:
            data = {
                TopologyManager.WEIGHT: weight,
                TopologyManager.STATE: state
            }
        else:
            data[TopologyManager.WEIGHT] = weight
            data[TopologyManager.STATE] = state

        result = GraphManager.new_node(data=data, *args, **kwargs)

        return result

    @classmethod
    def new_graph(cls, nodes=None, _type=None, *args, **kwargs):
        """
        Apply GraphManager.new_graph in ensuring a root exists in nodes.

        :param str _type: graph type.
        :param args: GraphManager.new_graph varargs
        :param kwargs: GraphManager.new_graph kwargs
        :return: GraphManager.new_graph result with a root in nodes.
        :rtype: dict
        """

        root = None

        if _type is None:
            _type = cls.GRAPH_TYPE

        if nodes is None:
            nodes = []
        else:
            if isinstance(nodes, dict):
                nodes = [nodes]
            # try to find a root in nodes
            for node in nodes:
                if TopologyManager.is_root(node):
                    root = node
                    nodes.remove(root)
                    break

        if root is None:
            root = TopologyManager.new_node(
                graph_id='', _type=TopologyManager.ROOT)

        # move the root at the beginning of nodes
        nodes.insert(0, root)

        result = GraphManager.new_graph(
            nodes=nodes, _type=_type, *args, **kwargs)

        return result

    @staticmethod
    def is_root(node):
        """
        Check if node is a root.

        :return: True if node is a root.
        :rtype: bool
        """
        return node[GraphManager.TYPE] == TopologyManager.ROOT

    @staticmethod
    def is_cluster(node):
        """
        Check if node is a cluster.

        :return: True if node is a cluster.
        :rtype: bool
        """

        return node[GraphManager.TYPE] == TopologyManager.CLUSTER
