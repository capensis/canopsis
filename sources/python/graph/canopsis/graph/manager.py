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

Technical
=========

The graph manager permits to get graph elements with any context information.

First, generic methods permit to get/put/delete elements in understanding such
elements such as dictionaries or GraphElement depending on serialize parameter
value.

Two, it is possible to find graphs, vertices and edges thanks to parameters
which correspond to their properties.
"""

from canopsis.common.init import basestring
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.graph.elements import Graph, Edge, GraphElement, Vertice

from itertools import chain

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

    SOURCES = 1 << 0  #: source orientation
    TARGETS = 1 << 1  #: target orientation
    ALL = SOURCES | TARGETS  #: source and target orientation

    def get_elts(
        self,
        ids=None, types=None, graph_ids=None, info=None, base_type=None,
        query=None, serialize=True, cls=None
    ):
        """Get graph element(s) related to input ids, types and query.

        :param ids: list of ids or id of element to retrieve. If None, get all
            elements. If str, get one element.
        :type ids: list or str
        :param types: graph element types to retrieve.
        :type types: list or str
        :param graph_ids: graph ids from where find elts.
        :type graph_ids: list or str
        :param info: info query.
        :param dict query: element search query.
        :param str base_type: elt base type.
        :param bool serialize: serialize result to GraphElements if True
            (by default).
        :param type cls: GraphElement type to retrieve if not None.

        :return: element(s) corresponding to input ids and query.
        :rtype: list or dict
        """

        # check if only one element is asked
        unique = isinstance(ids, basestring)
        # init query
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
        # put info if not None
        if info is not None:
            if isinstance(info, dict):
                for name in info:
                    data_name = 'info.{0}'.format(name)
                    query[data_name] = info[name]
            else:
                query[Vertice.INFO] = info
        # find ids among graphs
        if graph_ids is not None:
            result = []
            graphs = self.get_elts(ids=graph_ids, serialize=False)
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
        if result is not None and serialize:
            if isinstance(result, dict):
                result = GraphElement.new_element(**result)
                # ensure cls is respected
                if cls is not None and not isinstance(result, cls):
                    result = None
            else:
                # save reference to new_element in order to ease its use
                new_element = GraphElement.new_element
                result = list(
                    new_element(**elt) for elt in result
                )
                # ensure cls is respected
                if cls is not None:
                    result = [elt for elt in result if isinstance(elt, cls)]

        if unique and isinstance(result, list):
            result = result[0] if result else None

        return result

    def del_elts(
        self, ids=None, types=None, query=None, info=None, cache=False
    ):
        """Del elements identified by input ids in removing reference before.

        :param ids: list of ids or id elements to delete. If None, delete all
            elements.
        :type ids: list or str
        :param types: element types to delete.
        :type types: list or str
        :param dict query: additional deletion query.
        :param info: info query.
        :param bool cache: use query cache if True (False by default).
        """

        # initialize query if None
        if query is None:
            query = {}
        # put types in query
        if types is not None:
            query[GraphElement.TYPE] = types
        # put info if not None
        if info is not None:
            if isinstance(info, dict):
                for name in info:
                    data_name = 'info.{0}'.format(name)
                    query[data_name] = info[name]
            else:
                query[Vertice.INFO] = info
        # remove references in graph
        self.remove_elts(ids=ids, cache=cache, del_orphans=False)
        # remove edge references
        self.del_edge_refs(vids=ids, cache=cache)
        # remove elements
        self[GraphManager.STORAGE].remove_elements(
            ids=ids, _filter=query, cache=cache
        )

    def put_elt(self, elt, graph_ids=None, cache=False):
        """
        Put an element.

        :param dict elt: element to put.
        :type elt: dict or GraphElement
        :param str graph_ids: element graph id. None if elt is a graph.
        :param bool cache: use query cache if True (False by default).
        """

        # ensure elt is a dict
        if isinstance(elt, GraphElement):
            elt = elt.to_dict()
        # get elt uuid
        if GraphElement.ID not in elt:
            elt[GraphElement.ID] = GraphElement.new_id()
        elt_id = elt[GraphElement.ID]

        # put elt value in storage
        self[GraphManager.STORAGE].put_element(
            _id=elt_id, element=elt, cache=cache
        )
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
                        graph.save(self, cache=cache)

    def put_elts(self, elts, graph_ids=None, cache=False):
        """
        Put graph elements in DB.

        :param elts: graph elements to put in DB.
        :type elts: dict, GraphElement or list of dict/GraphElement.
        :param str graph_ids: element graph id. None if elt is a graph.
        :param bool cache: use query cache if True (False by default).
        """

        # ensure elts is a list
        if isinstance(elts, (dict, GraphElement)):
            elts = [elts]

        for elt in elts:
            gelt = elt
            if isinstance(gelt, dict):
                if not gelt.get(GraphElement.ID):
                    gelt[GraphElement.ID] = GraphElement.new_id()
                gelt = GraphElement.new_element(**gelt)
            # save elt
            gelt.save(manager=self, cache=cache, graph_ids=graph_ids)

        return elts

    def remove_elts(
        self, ids=None, graph_ids=None, cache=False, del_orphans=True
    ):
        """
        Remove vertices from graphs.

        :param ids: elt ids to remove from graph_ids.
        :type ids: list or str
        :param graph_ids: graph ids from where remove self.
        :type graph_ids: list or str
        :param bool cache: use query cache if True (False by default).
        :param bool del_orphans: delete vertices which are orphans.
        """

        # get graphs in order to remove references to self from them
        graphs = self.get_graphs(ids=graph_ids, elts=ids)
        if graphs is not None:
            # ensure graps is a list
            if isinstance(graphs, Graph):
                graphs = [graphs]
            # ensure ids is a list
            if isinstance(ids, basestring):
                ids = [ids]
            # if del orphans
            if del_orphans:
                # save graph ids such as a set
                graph_ids = set(graph.id for graph in graphs)
                # save elt_ids to remove in a list
                elt_ids = [] if ids is None else ids
            for graph in graphs:
                if ids is None:  # if ids is None, remove all graph elts
                    if del_orphans:
                        elt_ids += graph.elts
                    graph.remove_elts(graph.elts)
                else:  # else remove specific elt ids
                    graph.remove_elts(ids)
                    graph.save(manager=self, cache=cache)
            # delete orphans
            if del_orphans:
                # save elts to delete
                elts_to_delete = []
                # for all removed elt ids
                elts = self.get_elts(ids=elt_ids)
                for elt in elts:
                    # delete elts which are not graphs
                    if isinstance(elt, Graph):
                        continue
                    # find a graph
                    graphs = self.get_graphs(elts=elt.id)
                    # keep graphs which are not in graph_ids
                    elt_graphs = [
                        graph for graph in graphs
                        if graph.id not in graph_ids
                    ]
                    if not elt_graphs:  # if elt graphs do not exist
                        elts_to_delete.append(elt)
                # if there are elts to delete
                for elt in elts_to_delete:
                    elt.delete(manager=self, cache=cache)

    def del_edge_refs(
        self, ids=None, vids=None, sources=None, targets=None, del_empty=False,
        cache=False
    ):
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
        :param bool del_empty: if True and edges are not connected, delete
            them.
        :param bool cache: use query cache if True (False by default).
        """

        edges = self.get_edges(ids=ids, sources=sources, targets=targets)

        if edges is not None:
            # ensure edges is a list of edges
            if isinstance(edges, Edge):
                edges = [edges]
            # for all edges
            for edge in edges:
                # del refs
                edge.del_refs(ids=vids, sources=sources, targets=targets)

                # if del_empty and sources and targets are empty
                if del_empty and not (edge.sources or edge.targets):
                    # delete the edge
                    edge.delete(manager=self, cache=cache)
                else:  # save the edge
                    edge.save(manager=self, cache=cache)

    def get_graphs(
        self, ids=None, types=None, elts=None, graph_ids=None, info=None,
        query=None, add_elts=False, serialize=True
    ):
        """
        Get one or more graphs related to input ids, types and elts.

        :param ids: graph ids to retrieve.
        :type ids: list or str
        :param types: graph types to retrieve.
        :type types: list or str
        :param elts: graph elt ids.
        :type elts: basestring or list
        :param graph_ids: graph ids from where get graphs.
        :type graph_ids: list or str
        :param info: info to find among graphs.
        :param dict query: additional graph search query. Could help to search
            specific info information.
        :param bool add_elts: (False by default) add elts in the result. Works
            only if serialize is True.
        :param bool serialize: serialize result in GraphElements if True
            (by default).

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
        result = self.get_elts(
            ids=ids,
            query=query,
            types=types,
            graph_ids=graph_ids,
            info=info,
            base_type=Graph.BASE_TYPE,
            serialize=serialize
        )
        # if add_elts is asked
        if result is not None and add_elts:
            if serialize:  # add elts in _gelts
                if isinstance(result, Graph):
                    result.update_gelts(manager=self)
                else:
                    for graph in result:
                        graph.update_gelts(manager=self)
            else:  # add elts in _delts
                if isinstance(result, dict):
                    _delts = {}
                    elts = self.get_elts(
                        ids=result[Graph.ELTS],
                        serialize=False
                    )
                    for elt in elts:
                        elt_id = elt[GraphElement.ID]
                        _delts[elt_id] = elt
                    result[Graph._DELTS] = _delts
                else:
                    result = list(result)
                    for graph in result:
                        _delts = {}
                        elts = self.get_elts(
                            ids=graph[Graph.ELTS],
                            serialize=False
                        )
                        for elt in elts:
                            elt_id = elt[GraphElement.ID]
                            _delts[elt_id] = elt
                        graph[Graph._DELTS] = _delts

        return result

    def get_targets(
        self,
        ids=None, graph_ids=None,
        info=None, query=None,
        types=None, edge_ids=None, add_edges=False, edge_types=None,
        edge_data=None, edge_query=None, serialize=True
    ):
        """
        Ease the use of get_neighbourhood method in order to get targets
            vertices.

        :param ids: graph ids to retrieve.
        :type ids: list or str
        :param types: graph types to retrieve.
        :type types: list or str
        :param graph_ids: graph ids from where get graphs.
        :type graph_ids: list or str
        :param info: info to find among graphs.
        :param dict query: additional graph search query. Could help to search
            specific info information.
        :param bool serialize: serialize result in GraphElements if True
            (by default).

        :return: graph(s) corresponding to input parameters.
        :rtype: list or Graph
        """

        return self.get_neighbourhood(
            ids=ids, graph_ids=graph_ids, orientation=GraphManager.TARGETS,
            target_data=info, target_query=query, target_types=types,
            edge_ids=edge_ids, add_edges=add_edges,
            target_edge_types=edge_types, target_edge_data=edge_data,
            edge_query=edge_query, serialize=serialize
        )

    def get_sources(
        self,
        ids=None, graph_ids=None,
        info=None, query=None,
        types=None, edge_ids=None, add_edges=False, edge_types=None,
        edge_data=None, edge_query=None, serialize=True
    ):
        """
        Ease the use of get_neighbourhood method in order to get sources
            vertices.

        :param ids: graph ids to retrieve.
        :type ids: list or str
        :param types: graph types to retrieve.
        :type types: list or str
        :param graph_ids: graph ids from where get graphs.
        :type graph_ids: list or str
        :param info: info to find among graphs.
        :param dict query: additional graph search query. Could help to search
            specific info information.
        :param bool serialize: serialize result in GraphElements if True
            (by default).

        :return: graph(s) corresponding to input parameters.
        :rtype: list or Graph
        """

        return self.get_neighbourhood(
            ids=ids, graph_ids=graph_ids, orientation=GraphManager.SOURCES,
            source_data=info, source_query=query, source_types=types,
            edge_ids=edge_ids, add_edges=add_edges,
            source_edge_types=edge_types, source_edge_data=edge_data,
            edge_query=edge_query, serialize=serialize
        )

    def get_neighbourhood(
            self,
            ids=None, orientation=TARGETS,
            graph_ids=None,
            info=None, source_data=None, target_data=None,
            types=None, source_types=None, target_types=None,
            edge_ids=None, edge_types=None, add_edges=False,
            source_edge_types=None, target_edge_types=None,
            edge_data=None, source_edge_data=None, target_edge_data=None,
            query=None, edge_query=None, source_query=None, target_query=None,
            serialize=True, depth=None
    ):
        """
        Get neighbour vertices identified by context parameters.

        Sources and targets are handled in not directed edges.

        :param ids: vertice ids from where get neighbours.
        :type ids: list or str
        :param int orientation: edge orientation to use, among GRAPH.SOURCES,
            GRAPH.TARGETS (default) and GRAPH.ALL.
        :param bool sources: if True (False by default) add source vertices.
        :param bool targets: if True (default) add target vertices.
        :param graph_ids: vertice graph ids.
        :type graph_ids: list or str
        :param dict info: neighbourhood info to find.
        :param dict source_data: source neighbourhood info to find.
        :param dict target_data: target neighbourhood info to find.
        :param types: vertice type(s).
        :type types: list or str
        :param types: neighbourhood types to retrieve.
        :type types: list or str
        :param source_types: neighbourhood source types to retrieve.
        :type source_types: list or str
        :param target_types: neighbourhood target types to retrieve.
        :type target_types: list or str
        :param edge_ids: edge from where find target/source vertices.
        :type edge_ids: list or str
        :param edge_types: edge types from where find target/source vertices.
        :type edge_types: list or str
        :param bool add_edges: if True (False by default), add pathed edges in
            the result such as {edge_id: (edge, list(vertices))}.
        :param source_edge_types: edge types from where find source vertices.
        :type source_edge_types: list or str
        :param target_edge_types: edge types from where find target vertices.
        :type target_edge_types: list or str
        :param dict edge_data: edge info to find.
        :param dict source_edge_data: source edge info to find.
        :param dict target_edge_data: target edge info to find.
        :param dict query: additional search query.
        :param dict edge_query: additional edge query.
        :param dict source_query: additional source query.
        :param dict target_query: additional target query.
        :param bool serialize: serialize result in GraphElements if True
            (by default).
        :param int depth: if not None (default), repeat recursively the depth
            search and sort results by depth in ensuring a minimal depth for
            found neighbourhoods.

        :return: list of neighbour vertices designed by ids, or dict of
            {edge_id: (edge, list(vertices))} if add_edges. If depth is greater
            than 1 or negative, result a set of (search depth, previous result
            structure).
        :rtype: list or dict
        """

        result = {} if add_edges else []

        # init types
        if isinstance(types, basestring):
            types = [types]

        # init edges
        edges = dict()

        # search among source edges even if not sources because
        # edges can be not directed
        # init source_types
        if source_types is not None:
            if isinstance(source_types, basestring):
                source_types = [source_types]
        if types is not None:
            source_types += types
        # init source query
        if source_query is not None:
            if query is not None:
                source_query.update(query)
        else:
            source_query = query
        # init source edge types
        if source_edge_types is not None:
            if isinstance(source_edge_types, basestring):
                source_edge_types = [source_edge_types]
            if edge_types is not None:
                if isinstance(edge_types, basestring):
                    source_edge_types.append(edge_types)
                else:
                    source_edge_types += edge_types
        else:
            source_edge_types = edge_types
        # init source edge info
        if source_edge_data is not None:
            if edge_data is not None:
                source_edge_data.update(edge_data)
            else:
                source_edge_data = edge_data
        # get all source edges
        source_edges = self.get_edges(
            ids=edge_ids,
            graph_ids=graph_ids,
            types=source_edge_types,
            targets=ids,
            info=source_edge_data,
            query=source_query,
            serialize=False
        )

        # fill edges
        if source_edges is not None:
            sources = orientation & GraphManager.SOURCES
            # if source_edges is an edge
            if isinstance(source_edges, Edge):
                # and sources or source_edges is not directed
                if sources or not source_edges[Edge.DIRECTED]:
                    edges[source_edges[GraphElement.ID]] = source_edges
            elif sources:  # if sources
                for source_edge in source_edges:
                    source_edge_id = source_edge[GraphElement.ID]
                    if source_edge_id not in edges:
                        edges[source_edge_id] = source_edge
            else:
                for source_edge in source_edges:
                    # add not directed edges
                    if not source_edge[Edge.DIRECTED]:
                        edges[source_edge[GraphElement.ID]] = source_edge

        # search among target edges
        # init target types
        if target_types is not None:
            if isinstance(target_types, basestring):
                target_types = [target_types]
        if types is not None:
            target_types += types
        # init target query
        if target_query is not None:
            if query is not None:
                target_query.update(query)
        else:
            target_query = query
        # init target edge types
        if target_edge_types is not None:
            if isinstance(target_edge_types, basestring):
                target_edge_types = [target_edge_types]
            if edge_types is not None:
                if isinstance(edge_types, basestring):
                    target_edge_types.append(edge_types)
                else:
                    target_edge_types += edge_types
        else:
            target_edge_types = edge_types
        # init target edge info
        if target_edge_data is not None:
            if edge_data is not None:
                target_edge_data.update(edge_data)
            else:
                target_edge_data = edge_data
        # get all target edges
        target_edges = self.get_edges(
            ids=edge_ids,
            graph_ids=graph_ids,
            types=target_edge_types,
            sources=ids,
            info=target_edge_data,
            query=target_query,
            serialize=False
        )
        # fill edges
        if target_edges is not None:
            targets = orientation & GraphManager.TARGETS
            # if target_edges is an edge
            if isinstance(target_edges, Edge):
                # and targets or target_edges is not directed
                if targets or not target_edges[Edge.DIRECTED]:
                    if target_edges[GraphElement.ID] not in edges:
                        edges[target_edges[GraphElement.ID]] = target_edges
            elif targets:  # if targets
                for target_edge in target_edges:
                    if target_edge[GraphElement.ID] not in edges:
                        edges[target_edge[GraphElement.ID]] = target_edge
            else:
                for target_edge in target_edges:
                    # add not directed edges
                    if not target_edge[Edge.DIRECTED]:
                        if target_edge[GraphElement.ID] not in edges:
                            edges[target_edge[GraphElement.ID]] = target_edge

        # store edge sources and targets ids before get them at a time
        if not add_edges:
            edge_sources = []
            edge_targets = []

        if serialize:  # save new_element method in memory for quicker access
            new_element = GraphElement.new_element

        # add sources and targets from all edges
        for edge_id in edges:
            edge = edges[edge_id]
            if sources or not edge[Edge.DIRECTED]:
                if add_edges:
                    elts = self.get_elts(
                        ids=edge[Edge.SOURCES],
                        graph_ids=graph_ids,
                        info=source_data,
                        types=source_types,
                        query=source_query,
                        serialize=serialize
                    )
                    # serialize edge if required
                    _edge = new_element(**edge) if serialize else edge
                    if edge_id in result:
                        # TODO: check if this case can happen
                        result[edge_id][1] += elts
                    else:
                        result[edge_id] = (_edge, elts)
                else:
                    edge_sources += edge[Edge.SOURCES]
            if targets or not edge[Edge.DIRECTED]:
                if add_edges:
                    elts = self.get_elts(
                        ids=edge[Edge.TARGETS],
                        graph_ids=graph_ids,
                        info=target_data,
                        types=target_types,
                        query=target_query,
                        serialize=serialize
                    )
                    # serialize edge if required
                    _edge = new_element(**edge) if serialize else edge
                    if edge_id in result:
                        # TODO: check if this case can happen
                        result[edge_id][1] += elts
                    else:
                        result[edge_id] = (_edge, elts)
                else:
                    edge_targets += edge[Edge.TARGETS]

        # improve complexity if not add_edges
        if not add_edges:
            # get source graph elements
            if edge_sources:
                elts = self.get_elts(
                    ids=edge_sources,
                    graph_ids=graph_ids,
                    info=source_data,
                    types=source_types,
                    query=source_query,
                    serialize=serialize
                )
                result += elts

            # get target graph elements
            if edge_targets:
                elts = self.get_elts(
                    ids=edge_targets,
                    graph_ids=graph_ids,
                    info=target_data,
                    types=target_types,
                    query=target_query,
                    serialize=serialize
                )
                result += elts

        # if depth search is asked
        if depth is not None:
            # initialize the new result
            result = {0: result}
            foundvertices = []
            # initialize query
            if query is None:
                depth_query = {'$id': {'$nin': foundvertices}}
            else:
                depth_query = {
                    '$and': [{'$id': {'$nin': foundvertices}}, query]
                }

            # initialize the new query
            def getvertices(res):
                """Get found vertice from parent function result and
                fill query.

                :param res: neighbourhood result to parse.
                :return: found vertices.
                :rtype: set
                """

                result = set()
                # if res is a set of vertices by edges.
                if isinstance(res, dict):
                    for edge in res:
                        for vertice in res[edge]:
                            result.add(vertice)
                # if res is a list of vertices
                elif isinstance(res, list):
                    result = set(res)
                # if res is one vertice
                elif res is not None:
                    result.add(res)

                if result:  # update foundvertices if necessary
                    for vertice in result:
                        verticeid = vertice.id if serialize else vertice[
                            GraphElement.ID
                        ]
                        if verticeid not in foundvertices:
                            foundvertices.append(verticeid)

                return result

            verticeids = ids

            while depth != 0:
                depth -= 1
                new_result = self.get_neighbourhood(
                    ids=verticeids, orientation=orientation,
                    graph_ids=graph_ids, info=info, types=types,
                    source_data=source_data, target_data=target_data,
                    source_types=source_types, target_types=target_types,
                    edge_ids=edge_ids, edge_types=edge_types,
                    add_edges=add_edges, serialize=serialize,
                    source_edge_types=source_edge_types,
                    target_edge_types=target_edge_types,
                    edge_data=edge_data,
                    source_edge_data=source_edge_data,
                    target_edge_data=target_edge_data,
                    query=depth_query, edge_query=edge_query,
                    source_query=source_query, target_query=target_query,
                )
                if new_result:  # if new vertices are founded
                    # update vertice ids
                    verticeids = getvertices(new_result)
                    result[len(result)] = verticeids
                else:  # stop to search vertices
                    break

        return result

    def get_vertices(
        self,
        ids=None, graph_ids=None, types=None, info=None, query=None,
        serialize=True
    ):
        """
        Get graph vertices related to some context property.

        :param ids: vertice ids to get.
        :type ids: list or str
        :param graph_ids: vertice graph ids.
        :type graph_ids: list or str
        :param types: vertice type(s).
        :type types: list or str
        :param info: info to find among vertices.
        :param dict query: additional search query.
        :param bool serialize: serialize result in GraphElements if True
            (by default).

        :return: list of vertices if ids is a list. One vertice if ids is a
            str.
        :rtype: list or dict
        """

        result = self.get_elts(
            ids=ids,
            graph_ids=graph_ids,
            types=types,
            info=info,
            base_type=Vertice.BASE_TYPE,
            query=query,
            serialize=serialize
        )

        return result

    def get_edges(
        self,
        ids=None, types=None, sources=None, targets=None, graph_ids=None,
        info=None, query=None, serialize=True
    ):
        """Get edges related to input ids, types and source/target ids.

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
        :param bool serialize: serialize result in GraphElements if True
            (by default).

        :return: Edge(s) corresponding to input parameters.
        :rtype: Edge or list of Edges.
        """

        # by default, result is a list
        result = []

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

        result = self.get_elts(
            ids=ids,
            types=types,
            query=query,
            graph_ids=graph_ids,
            info=info,
            base_type=Edge.BASE_TYPE,
            serialize=serialize
        )

        return result

    def get_orphans(self, serialize=True, query=None, info=None, types=None):
        """Get all elements which are not associated to graphs.

        :param bool serialize: serialize result in GraphElements if True
            (by default).
        :param dict query: additional query.
        :param info: info to find among vertices.
        :param types: edge types to find.
        :type types: list or str
        :return: element(s) corresponding to input ids and query.
        :rtype: list or dict
        """

        # get all elt ids
        graphs = self.get_graphs()
        elt_ids = list(chain(*(graph.elts for graph in graphs)))
        # init query
        nin = {'$nin': elt_ids}
        if query is None:
            query = {}
        # get graph element id for quick access
        elt_id = GraphElement.ID
        # if elt_id is already in query, add 'nin' in an '$and' query
        if elt_id in query:
            query[elt_id] = {'$and': [query['id'], nin]}
        else:  # else use nin such as the elt_id query
            query[elt_id] = nin
        # find elts
        result = self.get_elts(
            types=types,
            query=query,
            serialize=serialize,
            info=info
        )
        return result
