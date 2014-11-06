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

A graph is an entity which permits to define relationships of entities in using
composed of nodes and edges each one bound to entities.

A node contains information about itself and associated entity.

An edge is a node which binds one or more nodes with direction property which
    specifies target and source nodes if true.

Technical
---------

As a graph is composed of nodes and edges, an edge inherits from a node.

A graph node contains::

    - graph: required field which binds a graph element to a graph id.
    - entity: optional field which binds a graph element to an entity id.
    - id: unique id among elements of the same graph_id and the same element
        kind.
    - type: optional field for graph element typing.
    - data: optional field which contains element data.
    - task: optional field which contains task information.

A graph edge contains::

    - sources: graph element ids.
    - targets: graph element ids.
    - directed: graph edge boolean direction information.
        If false, sources and targets are same.
"""

from uuid import uuid4 as uuid

from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.context.manager import Context
from canopsis.storage import Storage
from canopsis.task.manager import TaskManager

CONF_PATH = 'graph/graph.conf'
CATEGORY = 'GRAPH'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class GraphManager(MiddlewareRegistry):
    """
    Manage graph data.
    """

    STORAGE = 'graph_storage'  #: graph storage name

    ID = Storage.DATA_ID  #: graph element ID

    # graph specific fields
    NODES = 'nodes'  #: graph nodes attribute name
    EDGES = 'edges'  #: graph edges attribute name
    ENTITY = 'entity'  #: graph node entity attribute name

    # node specific fields
    GRAPH_ID = 'graph_id'  #: graph node graph id
    ENTITY_ID = 'entity_id'  #: graph node entity id field name
    TYPE = 'type'  #: graph node type field name
    DATA = 'data'  #: graph node data field name
    TASK = TaskManager.NAME  #: graph node task field name

    # edge specific fields
    SOURCES = 'sources'  #: graph node next field name
    TARGETS = 'targets'  #: graph targets field name
    DIRECTED = 'directed'  #: graph edge direction field name

    # element type specific field
    GRAPH_TYPE = 'graph'  #: graph type name
    GRAPH_NODE_TYPE = 'node'  #: graph node type name
    GRAPH_EDGE_TYPE = 'edge'  #: graph node type name

    def __init__(self, *args, **kwargs):

        super(GraphManager, self).__init__(*args, **kwargs)

        self.context = Context()

    def get_graph(self, graph_id, _type=None):
        """
        Get one graph and optionally its nodes.

        :param str graph_id: get a graph and its nodes related to input
            graph_id
        :param str _type: graph type to retrieve (default
            GraphManager.GRAPH_TYPE)
        """

        # get the right path
        if _type is None:
            _type = GraphManager.GRAPH_TYPE

        result = self.context.get(_type=_type, names=graph_id)

        if result is not None:
            nodes = self.get_nodes(graph_id=graph_id)
            result[GraphManager.NODES] = nodes

        return result

    def del_graph(self, graph_ids=None):
        """
        Get graph elements related to input graph_id.

        :param graph_ids: graph id from where get elements.
        :type graph_ids: None, str or list
        """

        # if graph ids is a str
        if isinstance(graph_ids, str):
            graph_ids = [graph_ids]

        # if graph ids is None, remove all
        if graph_ids is None:
            self[GraphManager.STORAGE].remove({})

        else:
            for graph_id in graph_ids:
                query = {GraphManager.GRAPH_ID: graph_id}
                self[GraphManager.STORAGE].remove_elements(query=query)
                self.context.remove(
                    _type=GraphManager.GRAPH_TYPE, ids=graph_id)

    def put_graph(self, graph):
        """
        Put graph in DB.

        :param dict graph: graph to put in DB.
        :return: putted graph.
        """

        # get graph id
        if GraphManager.ID not in graph:
            # set it in graph if it does not exist
            graph[GraphManager.ID] = str(uuid())
        graph_id = graph[GraphManager.ID]

        # create an entity sucha as a copy of input graph
        entity = graph.copy()
        # add type if not specified in the entity
        if GraphManager.TYPE in graph:
            _type = graph[GraphManager.TYPE]
        else:
            _type = GraphManager.GRAPH_TYPE

        # put nodes if they exist
        if GraphManager.NODES in graph:
            # get nodes
            nodes = graph[GraphManager.NODES]
            # put nodes in storage
            self.put_nodes(nodes=nodes, graph_id=graph_id)
            # delete nodes from entity
            del entity[GraphManager.NODES]

        # put graph in context
        self.context.put(_type=_type, entity=entity)

        return graph

    def put_nodes(self, nodes, graph_id=None):
        """
        Put nodes in DB and add node ids in nodes if they do not exist.

        :param nodes: node(s) to put in DB.
        :type nodes: dict or list
        :param str graph_id: graph id of nodes. If None, let the one in nodes.
        :return: putted nodes.
        """

        if isinstance(nodes, dict):
            nodes = [nodes]

        for node in nodes:
            # set node id
            if GraphManager.ID not in node:
                node[GraphManager.ID] = str(uuid())
            # set graph_id
            if graph_id is not None:
                node[GraphManager.GRAPH_ID] = graph_id
            _id = node[GraphManager.ID]
            # put node in storage
            self[GraphManager.STORAGE].put_element(_id=_id, element=node)

        return nodes

    def get_nodes(
        self,
        graph_id=None,
        ids=None, sources=None, targets=None,
        _type=None,
        entity_id=None,
    ):
        """
        Get graph nodes related to some context.

        :param str graph_id: graph id from where get nodes.
        :param ids: node ids.
        :type ids: list or str.
        :param list sources: source edge ids. If edges exist, add target nodes.
        :param list targets: target edge ids. If edges exist, add source nodes.
        :param str _type: graph type (default GraphManager.GRAPH_NODE_TYPE)
        :param str entity_id: related entity id.
        """

        result = []

        query = {}

        # add graph id in query
        if graph_id is not None:
            query[GraphManager.GRAPH_ID] = graph_id

        # initialize ids
        if ids is not None:
            if isinstance(ids, str):
                ids = [ids]
        else:
            ids = []

        # if target nodes are requested
        if sources is not None:
            # transform sources into a list
            if isinstance(sources, str):
                sources = [sources]
            # get all edges which have sources
            query[GraphManager.SOURCES] = {'$in': sources}
            edges = self[GraphManager.STORAGE].find_elements(query=query)
            # remove sources from query
            del query[GraphManager.SOURCES]
            # put edge targets in ids
            ids = [] if edges else [None]  # cancel future search if not edges
            for edge in edges:
                # put other sources in case of undirected edge
                if GraphManager.DIRECTED in edge \
                        and not edge[GraphManager.DIRECTED]:
                    ids += edge[GraphManager.SOURCES]
                    # remove one occurence of sources from edge sources
                    for source in sources:
                        if source in ids:
                            ids.remove(source)
                # put targets
                ids += edge[GraphManager.TARGETS]
            # add edges in the result
            result += edges

        # if source nodes are requested
        if targets is not None:
            # transform targets into a list
            if isinstance(targets, str):
                targets = [targets]
            # get all edges which have targets
            query[GraphManager.TARGETS] = {'$in': targets}
            edges = self[GraphManager.STORAGE].find_elements(query=query)
            # remove sources from query
            del query[GraphManager.TARGETS]
            # put edge sources in ids
            ids = [] if edges else [None]  # cancel future search if not edges
            targets = []
            for edge in edges:
                # put other targets in case of undirected edge
                if GraphManager.DIRECTED in edge \
                        and not edge[GraphManager.DIRECTED]:
                    ids += edge[GraphManager.TARGETS]
                    # remove one occurence of targets from edge targets
                    for target in targets:
                        if target in ids:
                            ids.remove(target)
                # put sources
                ids += edge[GraphManager.SOURCES]
            # add edges in the result
            result += edges

        # set entity_id into the query
        if entity_id is not None:
            query[GraphManager.ENTITY_ID] = entity_id

        # set type if not None
        if _type is not None:
            query[GraphManager.TYPE] = _type

        # get nodes related to ids
        if ids:
            query = {GraphManager.ID: {'$in': ids}}

        nodes = self[GraphManager.STORAGE].find_elements(query=query)

        # add nodes in result
        result += nodes

        # add entity in nodes
        for node in nodes:
            if GraphManager.ENTITY_ID in node:
                entity_id = node[GraphManager.ENTITY_ID]
                entities = self.context.get_entities(ids=entity_id)
                node[GraphManager.ENTITY] = entities

        return result

    def del_nodes(self, ids=None):
        """
        Delete one or more nodes depending on input ids:

            - an id: try to delete nodes where id correspond to id
            - list of ids: try to delete nodes where id are in input ids
            - None: delete all nodes
        """

        _filter = {}

        if ids is not None:
            if isinstance(ids, str):
                ids = [ids]
            _filter = {GraphManager.ID: {'$in': ids}}

        self[GraphManager.STORAGE].remove_elements(_filter=_filter)

    def is_edge(self, node):
        """
        True if node is an edge.

        :param dict node: node to compare to an edge.
        """

        # a node is an edge only if it has sources and targets
        result = GraphManager.SOURCES in node and GraphManager.TARGETS in node

        return result

    @classmethod
    def new_edge(
        cls, graph_id,
        _id=None, entity_id=None, _type=None, data=None, task=None,
        sources=None, targets=None, directed=True
    ):
        """
        Create a new edge with parameters

        :param str graph_id: graph id.
        :param str _id: edge id.
        :param str entity_id: bound entity id.
        :param str _type: edge type. If None, equals cls.GRAPH_EDGE_TYPE.
        :param data: edge data.
        :param task: edge task.
        :param list sources: source node ids.
        :param list targets: target node ids.
        :param bool directed: edge directed property (True by default).
        """

        # get default _type value if None
        if _type is None:
            _type = cls.GRAPH_EDGE_TYPE

        result = cls.new_node(
            graph_id=graph_id,
            _id=_id,
            entity_id=entity_id,
            _type=_type,
            data=data,
            task=task)
        # set sources, targets and directed
        result[GraphManager.SOURCES] = sources
        result[GraphManager.TARGETS] = targets
        result[GraphManager.DIRECTED] = directed

        return result

    @classmethod
    def new_node(
        cls,
        graph_id, _id=None, entity_id=None, _type=None, data=None, task=None
    ):
        """
        Create a new node with parameters.

        :param str graph_id: graph id.
        :param str _id: node id.
        :param str entity_id: bound entity id.
        :param str _type: node type. If None, equals cls.GRAPH_NODE_TYPE.
        :param data: node data.
        :param task: task information
        """

        # create a new node with graph_id
        result = {GraphManager.GRAPH_ID: graph_id}

        # set id
        if _id is None:
            _id = cls.new_id()
        result[GraphManager.ID] = _id
        # set entity_id
        result[GraphManager.ENTITY_ID] = entity_id
        # set _type
        if _type is None:
            _type = cls.GRAPH_NODE_TYPE
        result[GraphManager.TYPE] = _type
        # set data
        result[GraphManager.DATA] = data
        # set task
        result[GraphManager.TASK] = task

        return result

    @classmethod
    def new_graph(cls, _id=None, nodes=None, _type=None):
        """
        Create a graph related to an id, nodes and a type.

        :param str _id: graph id. Generated if None.
        :param list nodes: graph nodes to put in the graph.
        :param str _type: graph type. cls.GRAPH_TYPE if None.
        """

        result = {}

        result[GraphManager.TYPE] = cls.GRAPH_TYPE if _type is None else _type

        if _id is None:
            _id = cls.new_id()

        result[GraphManager.ID] = _id

        result[GraphManager.NODES] = nodes

        return result

    @classmethod
    def new_id(cls):
        """
        Generate a new id
        """

        return str(uuid())
