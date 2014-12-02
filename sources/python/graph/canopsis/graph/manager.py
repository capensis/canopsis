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
===========
Description
===========

This module defines the GraphManager which interacts between graph elements and
the DB.

Functional
==========

The role of the GraphManager is to ease graph element CRUD operations and
to retrieve graphs, vertices and edges thanks to methods with all element
parameters useful to find them.

This last could be specialized in assigning the GRAPH_TYPE class attribute to a
dedicated Graph type.

All methods may have to be enough generics without the need to override them,
the business code is ensured by graph elements.

Technical
=========

The graph manager permits to get graph elements with any context information.

First, generic methods permit to get/put/delete elements in understanding such
elements such as dictionaries.

Two, it is possible to find graphs, vertices and edges thanks to parameters
which correspond to their properties.
"""

from canopsis.common.init import basestring
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.graph.elements import Graph, Edge, GraphElement, Vertice

CONF_PATH = 'graph/graph.conf'
CATEGORY = 'GRAPH'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class GraphManager(MiddlewareRegistry):
    """
    Manage graph data in using a graph type which choose how to deserialize
    graph elements.
    """

    STORAGE = 'graph_storage'  #: graph storage name

    GRAPH_TYPE = Graph  #: default graph type

    def get_elts(
        self,
        ids=None, types=None, graph_ids=None, data=None, base_type=None,
        query=None
    ):
        """
        Get graph element(s) related to input ids, types and query.

        :param ids: list of ids or id of element to retrieve. If None, get all
            elements. If str, get one element.
        :type ids: list or str
        :param types: graph element types to retrieve.
        :type types: list or str
        :param graph_ids: graph ids from where find elts.
        :type graph_ids: list or str
        :param data: data query
        :param dict query: element search query.
        :param str base_type: elt base type.

        :return: element(s) corresponding to input ids and query.
        :rtype: list or dict
        """

        # init query if None
        if query is None:
            query = {}
        # put base type in query
        if base_type is not None:
            query[GraphElement.BASE_TYPE] = base_type
        # put types in query if not None
        if types is not None:
            if not isinstance(types, basestring):
                types = {'$in': types}
            query[GraphElement.TYPE] = types
        # put data if not None
        if data is not None:
            if isinstance(data, dict):
                for name in data:
                    data_name = 'data.{0}'.format(name)
                    query[data_name] = data[name]
            else:
                query[Vertice.DATA] = data
        # find ids among graphs
        if graph_ids is not None:
            result = []
            graphs = self.get_elts(ids=graph_ids)
            if graphs is not None:
                # all graph elt ids
                elt_ids = set()
                # ensure graphs is a list of graphs
                if isinstance(graphs, dict):
                    graphs = [graphs]
                for graph in graphs:
                    if Graph.ELTS in graph:
                        elts = set(graph[Graph.ELTS])
                        elt_ids |= elts
                # if ids is not given, use elt_ids
                if ids is None:
                    ids = list(elt_ids)
                else:  # else use jonction of elt_ids and ids
                    if isinstance(ids, basestring):
                        ids = [ids]
                    ids = list(elt_ids & set(ids))

        # get elements with ids and query
        result = self[GraphManager.STORAGE].get_elements(ids=ids, query=query)

        return result

    def del_elts(self, ids=None, types=None, query=None):
        """
        Del elements identified by input ids in removing reference before.

        :param ids: list of ids or id elements to delete. If None, delete all
            elements.
        :type ids: list or str
        :param types: element types to delete.
        :type types: list or str
        :param dict query: additional deletion query.
        """

        # initialize query if None
        if query is None:
            query = {}
        # put types in query
        if types is not None:
            query[GraphElement.TYPE] = types
        # remove references in graph
        self.remove_elts(ids=ids)
        # remove edge references
        self.del_edge_refs(vids=ids)
        # remove elements
        self[GraphManager.STORAGE].remove_elements(ids=ids, _filter=query)

    def put_elt(self, elt, graph_ids=None):
        """
        Put an element.

        :param dict elt: element to put.
        :type elt: dict or GraphElement
        :param str graph_ids: element graph id. None if elt is a graph.
        """

        # ensure elt is a dict
        if isinstance(elt, GraphElement):
            elt = elt.to_dict()
        elt_id = elt[GraphElement.ID]

        # put elt value in storage
        self[GraphManager.STORAGE].put_element(_id=elt_id, element=elt)
        # update graph if graph_id is not None
        if graph_ids is not None:
            graphs = self.get_graphs(ids=graph_ids)
            if graphs is not None:
                # ensure graphs is a list of graphs
                if isinstance(graphs, Graph):
                    graphs = [graphs]
                for graph in graphs:
                    # if graph exists and elt_id not already present
                    if elt_id not in graph.elts:
                        # add elt_id in graph elts
                        graph.elts.append(elt_id)
                        # update graph
                        graph_dict = graph.to_dict()
                        self.put_elt(elt=graph_dict)

    def remove_elts(self, ids, graph_ids=None):
        """
        Remove vertices from graphs.

        :param ids: elt ids to remove from graph_ids.
        :type ids: list or str
        :param graph_ids: graph ids from where remove self.
        :type graph_ids: list or str
        """

        # get graphs in order to remove references to self from them
        graphs = self.get_graphs(ids=graph_ids, elts=ids)
        if graphs is not None:
            # ensure graps is a list
            if isinstance(graphs, Graph):
                graphs = [graphs]
            if ids is None:
                ids = []
            elif isinstance(ids, basestring):
                ids = [ids]
            for graph in graphs:
                for _id in ids:
                    if _id in graph.elts:
                        # remove elf from graph.elts
                        graph.remove_elts(_id)
                        # save the graph
                        graph.save(manager=self)

    def del_edge_refs(self, ids=None, vids=None, sources=None, targets=None):
        """
        Delete references of vertices from edges.

        :param ids: edge ids to select for removing vertices links.
        :type ids: list or str
        :param vids: vertice ids to remove from edge links.
        :type vids: list or str
        :param sources: source ids to remove.
        :type sources: list or str
        :param targets: target ids to remove.
        :type targets: list or str
        """

        edges = self.get_edges(ids=ids, sources=sources, targets=targets)

        if edges is not None:
            # ensure edges is a list
            if isinstance(edges, basestring):
                edges = [edges]
            # for all edges
            for edge in edges:
                # del refs
                edge.del_refs(ids=vids, sources=sources, targets=targets)
                # and save them
                edge.save(manager=self)

    def get_graphs(
        self, ids=None, types=None, elts=None, graph_ids=None, query=None,
        add_elts=False
    ):
        """
        Get one or more graphs related to input ids, types and elts.

        :param ids: graph ids to retrieve.
        :type ids: list or str
        :param types: graph types to retrieve.
        :type types: list or str
        :param elts: elts embedded by graphs to retrieve.
        :type elts: list or str
        :param graph_ids: graph ids from where get graphs.
        :type graph_ids: list or str
        :param dict query: additional graph search query. Could help to search
            specific data information.
        :param bool add_elts: (False by default) add elts in the result.

        :return: graph(s) corresponding to input parameters.
        :rtype: list or Graph
        """

        result = []
        # init query
        if query is None:
            query = {}
        # put elts in query
        if elts is not None:
            if not isinstance(elts, basestring):
                elts = {'$in': elts}
            query[Graph.ELTS] = elts
        # get graphs with ids and query
        graphs = self.get_elts(
            ids=ids,
            query=query,
            types=types,
            graph_ids=graph_ids,
            base_type=Graph.BASE_TYPE
        )
        # check if result may one graph or a list
        unique = False
        # if add_elts is asked
        if graphs is not None:
            # if graphs is unique, graphs is a dict
            if isinstance(graphs, dict):
                # set unique value to True
                unique = True
                # put graphs in a list
                graphs = [graphs]
            # save new_element method for quick resolution
            new_element = GraphElement.new_element
            # iterate on graphs
            for graph in graphs:
                try:
                    new_graph = new_element(**graph)
                    # if add elts, find them and add them into the graph
                    if add_elts:
                        # fill the graph with
                        new_graph.update_gelts(manager=self)
                    result.append(new_graph)
                except TypeError as te:
                    self.logger.warning(
                        'element {0} is not a graph. {1}'.format(graph, te)
                    )
        # get the first element if unique
        if unique:
            result = result[0] if result else None

        return result

    def get_vertices(
        self,
        ids=None, graph_ids=None, types=None,
        sources=None, targets=None, source_types=None, target_types=None,
        edge_ids=None, edge_types=None, add_edges=False,
        src_edge_types=None, trgt_edge_types=None,
        union=True,
        query=None
    ):
        """
        Get graph vertices related to some context property.

        :param ids: vertice ids to get.
        :type ids: list or str
        :param graph_ids: vertice graph ids.
        :type graph_ids: list or str
        :param types: vertice type(s).
        :type types: list or str
        :param sources: source edge id(s). If edges exist, add target
            vertices.
        :type sources: list or str
        :param targets: target edge id(s). If edges exist, add source
            vertices.
        :type sources: list or str
        :param edge_ids: edge from where find target/source vertices.
        :type edge_ids: list or str
        :param edge_types: edge types from where find target/source vertices.
        :type edge_types: list or str
        :param bool add_edges: if True (default), add pathed edges in the
            result.
        :param src_edge_types: edge types from where find source vertices.
        :type src_edge_types: list or str
        :param trgt_edge_types: edge types from where find target vertices.
        :type trgt_edge_types: list or str
        :param bool union: if True (default) do an union of all results,
            otherwise, do an intersection.
        :param dict query: additional search query.

        :return: list of vertices if ids is a list, or sources/targets are
            lists or else graph_ids is not None. One vertice if ids is a str
            and other params are None.
        :rtype: list or dict
        """

        result = []

        # init unique result by default
        unique = False

        # elt ids to retrieve
        elt_ids = None

        # boolean value for one value to retrieve.
        one_value = False

        # if ids is not None, get related ids
        if ids is not None:
            elt_ids = set()
            if isinstance(ids, str):
                one_value = True
                elt_ids.add(ids)
            else:
                for _id in ids:
                    elt_ids.add(_id)

        found_edges = []

        # if source vertices are requested
        if (sources, source_types) != (None, None):
            # force edge_types and src_edge_types to be list if basestring
            if isinstance(sources, basestring):
                sources = [sources]
            if isinstance(edge_types, basestring):
                edge_types = [edge_types]
            if isinstance(src_edge_types, basestring):
                src_edge_types = [src_edge_types]
            if src_edge_types is None:
                src_edge_types = edge_types
            elif edge_types is not None:
                src_edge_types = list(set(edge_types + src_edge_types))
            edges = self.get_edges(
                ids=edge_ids, sources=sources, types=src_edge_types
            )
            # if edges exist
            if edges is not None:
                # firce edges to be a list
                if isinstance(edges, Edge):
                    edges = [edges]
                # get target vertices and sources as well if edge is undirected
                for edge in edges:
                    target_ids = edge.targets
                    if not edge.directed:
                        target_ids += edge.sources
                        # remove references to sources if they exist
                        if sources:
                            for source in sources:
                                target_ids.remove(source)
                    elts = self.get_elt(ids=target_ids, types=target_types)
                    target_ids = [elt.id for elt in elts]
                    if elt_ids is None:
                        elt_ids = set(target_ids)
                    elif union:
                        elt_ids |= target_ids
                    else:
                        elt_ids &= target_ids
                # if add_edges, add them to elt_ids
                if add_edges:
                    found_edges += edges

        # if target vertices are requested
        if (targets, target_types) != (None, None):
            # force edge_types and trgt_edge_types to be list if basestring
            if isinstance(targets, basestring):
                targets = [targets]
            if isinstance(edge_types, basestring):
                edge_types = [edge_types]
            if isinstance(trgt_edge_types, basestring):
                trgt_edge_types = [trgt_edge_types]
            if trgt_edge_types is None:
                trgt_edge_types = edge_types
            elif edge_types is not None:
                trgt_edge_types = list(set(edge_types + trgt_edge_types))
            edges = self.get_edges(
                ids=edge_ids, targets=sources, types=trgt_edge_types
            )
            # if edges exist
            if edges is not None:
                # firce edges to be a list
                if isinstance(edges, Edge):
                    edges = [edges]
                # get source vertices and targets as well if edge is undirected
                for edge in edges:
                    source_ids = edge.sources
                    if not edge.directed:
                        source_ids += edge.targets
                        # remove references to targets if they exist
                        if targets:
                            for target in targets:
                                source_ids.remove(target)
                    elts = self.get_elt(ids=source_ids, types=source_types)
                    source_ids = [elt.id for elt in elts]
                    if elt_ids is None:
                        elt_ids = set(source_ids)
                    elif union:
                        elt_ids |= source_ids
                    else:
                        elt_ids &= source_ids
                # if add_edges, add them to elt_ids
                if add_edges:
                    found_edges += edges

        if elt_ids is not None:
            elt_ids = list(elt_ids)
            if one_value and len(elt_ids) == 1:
                elt_ids = elt_ids[0]

        vertices = self.get_elts(
            ids=elt_ids,
            query=query,
            graph_ids=graph_ids,
            types=types,
            base_type=Vertice.BASE_TYPE
        )

        # put right graph elements in result
        if vertices is not None:
            # ensure vertices is a list
            if isinstance(vertices, dict):
                # unique is True if vertices is a dict and found_edges is empty
                unique = not found_edges
                vertices = [vertices]
            graph_type = self.GRAPH_TYPE
            for vertice in vertices:
                elt = graph_type.new_element(**vertice)
                if isinstance(elt, Vertice):
                    result.append(elt)
        else:
            result = None

        # add edges in result
        if found_edges:
            result += found_edges

        # if unique, return first element or None if empty
        if unique:
            result = result[0] if result else None

        return result

    def get_edges(
        self,
        ids=None, types=None, sources=None, targets=None, graph_ids=None,
        query=None
    ):
        """
        Get edges related to input ids, types and source/target ids.

        :param ids: edge ids to find. If ids is a str, result is an Edge or
            None.
        :type ids: list or str
        :param types: edge types to find.
        :type types: list or str
        :param sources: source edge attribute to find.
        :type sources: list or str
        :param targets: target edge attribute to find.
        :type targets: list or str
        :param graph_ids: graph ids from where find edges.
        :type graph_ids: list or str
        :param dict query: additional query.

        :return: Edge(s) corresponding to input parameters.
        :rtype: Edge or list of Edges.
        """

        # by default, result is a list
        result = []

        # init unique by default
        unique = False

        if query is None:
            query = {}

        if sources is not None:
            if not isinstance(sources, basestring):
                sources = {'$in': sources}
            query[Edge.SOURCES] = sources

        if targets is not None:
            if not isinstance(targets, basestring):
                targets = {'$in': targets}
            query[Edge.TARGETS] = targets

        edges = self.get_elts(
            ids=ids,
            types=types,
            query=query,
            graph_ids=graph_ids,
            base_type=Edge.BASE_TYPE
        )

        if edges is not None:
            # ensure edges is a list
            if isinstance(edges, dict):
                unique = True
                edges = [edges]
            graph_type = self.GRAPH_TYPE
            for edge in edges:
                elt = graph_type.new_element(**edge)
                if isinstance(elt, Edge):
                    result.append(elt)
        else:
            result = None

        if unique:
            result = result[0] if result else None

        return result
