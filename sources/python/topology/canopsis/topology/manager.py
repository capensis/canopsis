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

from canopsis.common.utils import force_iterable
from canopsis.configuration import conf_paths, add_category
from canopsis.middleware.manager import Manager
from canopsis.storage.composite import CompositeStorage
from canopsis.storage.filter import Filter

CONF_PATH = 'topology/topology.conf'
CATEGORY = 'TOPOLOGY'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class Topology(Manager):
    """
    Manage access to topologies
    """

    STORAGE = 'topology_storage'

    ENTITY_ID = 'entity'
    NEXT = 'next'
    NODES = 'nodes'
    TYPE = 'type'

    def __init__(self, *args, **kwargs):

        super(Topology, self).__init__(*args, **kwargs)

    @staticmethod
    def _get_topology_scope():
        """
        Get topology scope
        """
        return {Topology.TYPE: Topology.DATA_SCOPE}

    @staticmethod
    def _get_topology_node_scope(topology_id=None):
        """
        Get topology node scope related to a topology_id
        """

        result = {'type': 'topology_node'}
        if topology_id is not None:
            result[Topology.STORAGE] = topology_id
        return result

    def _add_nodes(self, topologies):
        """
        Add nodes into input topologies.
        """

        if topologies:
            topologies = force_iterable(topologies)
            for topology in topologies:
                nodes = self.get_nodes(ids=topology[CompositeStorage.ID])
                topology['nodes'] = nodes

    def get_topologies(self, ids=None, add_nodes=False):
        """
        Get one or more topologies.

        :return: depending on ids:

            - an id: one topology or None if corresponding topology does not
                exist.
            - list of id: list of topologies or empty list if ids do not
                correspond to existing topologies.
            - None: all existing topologies.
        """
        path = self._get_topology_scope()
        result = self[Topology.STORAGE].get(path=path, ids=ids)
        if add_nodes:
            self._add_nodes(result)
        return result

    def find(self, regex=None, add_nodes=False):
        """
        Find topologies where id match with inpur regex_id
        """

        path = self._get_topology_scope()
        _filter = Filter()
        _filter.add_regex(
            name=CompositeStorage.ID, value=regex, case_sensitive=True)
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

        path = self._get_topology_scope()
        self[Topology.STORAGE].remove(path=path, ids=ids)

    def push(self, _id, topology):
        """
        Push one topology.
        """

        path = self._get_topology_scope()
        self[Topology.STORAGE].put(path=path, _id=_id, data=topology)

    def push_node(self, _id, node):
        """
        Push a node.
        """

        path = self._get_topology_node_scope(node['id'])
        self[Topology.STORAGE].put(path=path, _id=_id, data=node)

    def get_nodes(self, ids=None):
        """
        Get nodes
        """

        path = self._get_topology_node_scope()
        result = self[Topology.STORAGE].get(path=path, ids=ids)

        return result

    def find_nodes_by_entity_id(self, entity_id):
        """
        Find all nodes related to input entity_id
        """
        path = self._get_topology_node_scope()
        _filter = Filter()
        _filter[Topology.ENTITY_ID] = entity_id
        result = self[Topology.STORAGE].find(path=path, _filter=_filter)

        return result

    def find_nodes_by_next(self, next=None):
        """
        Find source nodes from input next node. If next is None, get all root
        nodes
        """
        path = self._get_topology_node_scope()
        _filter = Filter()
        _filter[Topology.NEXT] = next
        result = self[Topology.STORAGE].find(path=path, _filter=_filter)

        return result
