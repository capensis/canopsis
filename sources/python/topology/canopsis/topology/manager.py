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

A topology enrichs a model of entities in adding dependency among them and in
executing rules which can modify the model, send events, etc.

A topology permits to run automatized and distributed operations over a context

An example of application is root cause analysis where a topology may react
when an entity change of state and can propagate over topology nodes the change
of state in some propagation conditions such as operations like ``worst state``

As a topology and topology nodes are entities, it is possible to bind a node to
a topology, therefore, a topology is bound to a root node.

A rule is a couple of (condition, actions). A condition takes in parameter the
execution context of an event, the engine and the node which embeds the rule.

Technical
---------

A topology is an entity which contains following fields::

    - id: unique id among topologies.
    - nodes: list of topology node ids.
    - root: root node id.

A topology node contains following fields::

    - id: unique topology node id.
    - entity_id: bound entity id. May be self id if no entity is bound.
    - next: list of topology-node ids which depends on this.
    - rule: rule to apply.
"""

from canopsis.common.utils import ensure_iterable

from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)

from canopsis.storage import Storage
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.storage.filter import Filter

CONF_PATH = 'topology/topology.conf'
CATEGORY = 'TOPOLOGY'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class Topology(MiddlewareRegistry):
    """
    Manage topological data
    """

    STORAGE = 'topology_storage'  #: topology storage name

    ENTITY_ID = 'entity_id'  #: topology node entity id field name
    NEXT = 'next'  #: topology node next field name
    NODES = 'nodes'  #: topology nodes field name
    TYPE = 'type'  #: topology/topology node type field name
    ROOT = 'root'  #: topology root node
    TOPOLOGY_ID = 'pid'  #: topology node topology id

    TOPOLOGY_TYPE = 'topology'  #: topology type name
    TOPOLOGY_NODE_TYPE = 'topology-node'  #: topology node type name

    ID = Storage.DATA_ID  #: topology and topology node id field name

    def _add_nodes(self, topologies):
        """
        Add nodes into input topologies.
        """

        if topologies:
            topologies = ensure_iterable(topologies)

            for topology in topologies:
                nodes = self.get_nodes(ids=topology[Topology.ID])
                topology[Topology.NODES] = nodes

    def get(self, ids=None, add_nodes=False):
        """
        Get one or more topologies.

        :return: depending on ids:

            - an id: one topology or None if corresponding topology does not
                exist.
            - list of id: list of topologies or empty list if ids do not
                correspond to existing topologies.
            - None: all existing topologies.
        """

        # get the right path
        path = {
            Topology.TYPE: Topology.TOPOLOGY_TYPE
        }

        result = self[Topology.STORAGE].get(path=path, data_ids=ids)

        if add_nodes:
            self._add_nodes(result)

        return result

    def find(self, regex=None, add_nodes=False):
        """
        Find topologies where id match with inpur regex_id
        """

        # get the right filter
        _filter = Filter()
        _filter.add_regex(
            name=Storage.DATA_ID, value=regex, case_sensitive=True
        )

        # get the right path
        path = {
            Topology.TYPE: Topology.TOPOLOGY_TYPE
        }

        result = self[Topology.STORAGE].find(path=path, _filter=_filter)

        if add_nodes:
            self._add_nodes(result)

        return result

    def delete(self, ids=None):
        """
        Delete one or more topologies depending on input ids:

            - an id: try to delete topology where id correspond to id
            - list of ids: try to delete topologies where id are in input ids
            - None: delete all topologies
        """

        path = {
            Topology.TYPE: Topology.TOPOLOGY_TYPE
        }

        self[Topology.STORAGE].remove(path=path, data_ids=ids)

    def delete_nodes(self, ids=None):
        """
        Delete one or more topology nodes depending on input ids:

            - an id: delete the topology node where id corresponds to input id
            - list of ids: delete topology nodes where ids equal input ids
            - None: delete all topology nodes
        """

        path = {
            Topology.TYPE: Topology.TOPOLOGY_NODE_TYPE
        }

        self[Topology.STORAGE].remove(path=path, data_ids=ids)

    def push(self, topology):
        """
        Push one topology.
        """

        path = {
            Topology.TYPE: Topology.TOPOLOGY_TYPE
        }

        _id = topology[Topology.ID]

        # if topology contains nodes
        if Topology.NODES in topology:
            # get nodes
            nodes = topology[Topology.NODES]
            # in case of nodes are dictionaries
            if nodes and isinstance(nodes[0], dict):
                # add nodes
                for node in nodes:
                    self.push_node(node=node)
                # and transform the content of topology nodes into node ids
                nodes = [node[Topology.ID] for node in nodes]
                topology[Topology.NODES] = nodes

        # finally, put the topology in storage
        self[Topology.STORAGE].put(path=path, data_id=_id, data=topology)

    def push_node(self, node):
        """
        Push a node.
        """

        path = {
            Topology.TYPE: Topology.TOPOLOGY_NODE_TYPE
        }

        _id = node.pop(Topology.ID)
        self[Topology.STORAGE].put(path=path, data_id=_id, data=node)

    def get_nodes(self, topology_id=None, ids=None):
        """
        Get topology nodes

        :param str topology_id: topology id
        :param ids: topology node id or list of ids
        :type ids: str or list(str)

        :return: depending on type of ids::
            - str: node or None
            - list: list of nodes
        """

        path = {
            Topology.TYPE: Topology.TOPOLOGY_NODE_TYPE
        }

        _filter = None if topology_id is None else {
            Topology.TOPOLOGY_ID: topology_id
        }

        result = self[Topology.STORAGE].get(
            path=path, data_ids=ids, _filter=_filter)

        return result

    def find_bound_nodes(self, entity_id):
        """
        Find all nodes related to input entity_id
        """

        _filter = Filter()
        _filter[Topology.ENTITY_ID] = entity_id

        path = {
            Topology.TYPE: Topology.TOPOLOGY_NODE_TYPE
        }

        result = self[Topology.STORAGE].find(path=path, _filter=_filter)

        return result

    def get_next_nodes(self, node):
        """
        Get next nodes from input source node.
        """

        result = []

        if Topology.NEXT in node:
            next_ids = node[Topology.NEXT]

            path = {
                Topology.TYPE: Topology.TOPOLOGY_NODE_TYPE
            }

            result = self[Topology.STORAGE].get(path=path, data_ids=next_ids)

        return result

    def find_source_nodes(self, node=None):
        """
        Get previous nodes from input next node. If node is None, get all root
        nodes.
        """

        # get the right path
        path = {
            Topology.TYPE: Topology.TOPOLOGY_NODE_TYPE
        }

        # get the right filter
        _filter = Filter()

        if node is None:
            # get root nodes which don't have next field or next field is empty
            _filter[Topology.NEXT] = {
                '$or': [
                    {'$empty'},
                    {'$not': {'$exists'}}
                ]
            }

        else:
            node_id = node[Topology.ID]
            _filter[Topology.NEXT] = node_id

        result = self[Topology.STORAGE].get(path=path, _filter=_filter)

        return result

    def get_id(self, element):
        """
        Get input element id

        :param node: topology or topology node
        """

        result = element[Topology.ID]

        return result
