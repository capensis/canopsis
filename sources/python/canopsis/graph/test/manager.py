#!/usr/bin/env python
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

from unittest import TestCase, main

from canopsis.graph.manager import GraphManager
from canopsis.graph.elements import Graph, Vertice, Edge


class GraphTest(TestCase):

    def setUp(self):
        """
        Create self.count=10 and self.graphs
        """

        self.manager = GraphManager(data_scope='test_graph')

        self.manager.del_elts()

        self.count = 10
        self.graphs = [None] * 10
        self.types = (i for i in range(self.count))

        # create self.count vertices
        self.vertices = [None] * self.count
        for index in range(self.count):
            self.vertices[index] = Vertice(
                id='vertice-{0}'.format(index),
                info=None if index % 2 else {'index': index},
                type='tvertice-{0}'.format(index),
            )
        self.vertice_ids = [vertice.id for vertice in self.vertices]
        self.vertice_types = [vertice.type for vertice in self.vertices]

        # create self.count edges
        self.edges = [None] * self.count
        for index in range(self.count):
            self.edges[index] = Edge(
                id='edge-{0}'.format(index),
                info=None if index % 2 else {'index': index},
                type='tedge-{0}'.format(index),
                sources=self.vertice_ids[index:],
                targets=self.vertice_ids[:-index],
                directed=index % 2 == 0
            )
        self.edge_ids = [edge.id for edge in self.edges]
        self.edge_types = [edge.type for edge in self.edges]

        # create self.count graphs
        self.graphs = [None] * self.count
        for index in range(self.count):
            self.graphs[index] = Graph(
                id='graph-{0}'.format(index),
                info=None if index % 2 else {'index': index},
                type='tgraph-{0}'.format(index),
                elts=self.edge_ids[index:] + self.vertice_ids[index:]
            )
        self.graph_ids = [graph.id for graph in self.graphs]
        self.graph_types = [graph.type for graph in self.graphs]

        # get all elt ids
        self.elt_ids = self.graph_ids + self.vertice_ids + self.edge_ids
        # get all elts
        self.elts = self.vertices + self.edges + self.graphs

        # add edges and graphs in edges
        for index in range(self.count):
            edge = self.edges[index]
            edge.sources += self.edge_ids[index:] + self.graph_ids[index:]
            edge.targets += self.edge_ids[index:] + self.graph_ids[index:]

        # add graph_ids in graph
        for index in range(self.count):
            self.graphs[index].elts += self.graph_ids[index:]

        # ensure DB is empty
        self.manager.del_elts()
        elts = self.manager.get_elts()
        self.assertFalse(elts)

        # put vertices in DB
        for vertice in self.vertices:
            vertice.save(manager=self.manager)

        # put edges in DB
        for edge in self.edges:
            edge.save(manager=self.manager)

        # put graphs in DB
        for graph in self.graphs:
            graph.save(manager=self.manager)

    def tearDown(self):
        """
        Del all elements
        """

        self.manager.del_elts()

    def test_get_elts(self):
        """
        Test to put graph and get them back
        """

        self.manager.del_elts()
        # ensure no element exist in storage
        vertices = self.manager.get_elts(ids=self.vertice_ids)
        self.assertFalse(vertices)
        edges = self.manager.get_elts(ids=self.edge_ids)
        self.assertFalse(edges)
        graphs = self.manager.get_elts(ids=self.graph_ids)
        self.assertFalse(graphs)

        # put all elements
        for vertice in self.vertices:
            vertice.save(manager=self.manager)
        for edge in self.edges:
            edge.save(manager=self.manager)
        for graph in self.graphs:
            graph.save(manager=self.manager)

        # ensure to get all elements
        vertices = self.manager.get_elts(ids=self.vertice_ids)
        self.assertEqual(len(vertices), self.count)
        edges = self.manager.get_elts(ids=self.edge_ids)
        self.assertEqual(len(edges), self.count)
        graphs = self.manager.get_elts(ids=self.graph_ids)
        self.assertEqual(len(graphs), self.count)

        # check for get all elts
        elts = self.manager.get_elts()
        self.assertEqual(len(elts), 3 * self.count)

        # check to get elements by graph ids
        elts = self.manager.get_elts(graph_ids=self.graph_ids)
        self.assertEqual(len(elts), 3 * self.count)

        # check to get no elements by graph ids where ids does not exist
        elts = self.manager.get_elts(ids='', graph_ids=self.graph_ids)
        self.assertFalse(elts)

        # check to get elements by ids and graph ids
        for graph in self.graphs:
            elts = self.manager.get_elts(graph_ids=graph.id)
            self.assertEqual(len(elts), len(graph.elts))

        # check to get one element from graph_ids
        last_vertice_id = self.vertice_ids[-1]
        vertice = self.manager.get_elts(
            ids=last_vertice_id, graph_ids=self.graph_ids
        )
        self.assertTrue(isinstance(vertice, Vertice))

        # check get info
        for elt in self.elts:
            _elt = self.manager.get_elts(ids=elt.id, info=elt.info)
            self.assertIsNotNone(_elt)
            self.assertEqual(_elt.id, elt.id)
            if elt.info is not None:
                elts = self.manager.get_elts(info=elt.info)
                self.assertEqual(len(elts), 3)

    def test_get_edges(self):
        """
        Test GraphManager.get_edges.
        """

        # check get all edges
        edges = set(self.manager.get_edges())
        self.assertEqual(edges, set(self.edges))

        # check get all self edges
        edges = set(self.manager.get_edges(ids=self.edge_ids))
        self.assertEqual(edges, set(self.edges))

        # check get one edge
        for edge_id in self.edge_ids:
            edge = self.manager.get_edges(ids=edge_id)
            self.assertIn(edge, set(self.edges))

        # check get edge types
        edges = set(self.manager.get_edges(types=self.edge_types))
        self.assertEqual(edges, set(self.edges))

        # check to get one edge type
        for edge_type in self.edge_types:
            edges = self.manager.get_edges(types=edge_type)
            self.assertEqual(len(edges), 1)

        # check get edges by sources
        edges = set(self.manager.get_edges(sources=self.elt_ids))
        self.assertEqual(edges, set(self.edges))

        # check get edges by source
        for edge in self.edges:
            sources = set(edge.sources)
            edges = set()
            for _edge in self.edges:
                _sources = set(_edge.sources)
                if sources & _sources:
                    edges.add(_edge)
            _edges = set(self.manager.get_edges(sources=list(sources)))
            self.assertEqual(_edges, edges)

        # check get edges by targets
        edges = set(self.manager.get_edges(targets=self.elt_ids))
        self.assertEqual(edges, set(self.edges))

        # check get edges by target
        for edge in self.edges:
            targets = set(edge.targets)
            edges = set()
            for _edge in self.edges:
                _targets = set(_edge.targets)
                if targets & _targets:
                    edges.add(_edge)
            _edges = set(self.manager.get_edges(targets=list(targets)))
            self.assertEqual(_edges, edges)

    def test_get_vertices(self):
        """
        Test GraphManager.get_vertices.
        """

        # check get all edges
        vertices = set(self.manager.get_vertices())
        self.assertEqual(vertices, set(self.vertices))

        # check get all self vertices

        vertices = set(self.manager.get_vertices(ids=self.vertice_ids))
        self.assertEqual(vertices, set(self.vertices))

        # check get one vertice
        for vertice_id in self.vertice_ids:
            vertice = self.manager.get_vertices(ids=vertice_id)
            self.assertIn(vertice, set(self.vertices))

        # check get vertice types
        vertices = set(self.manager.get_vertices(types=self.vertice_types))
        self.assertEqual(vertices, set(self.vertices))

        # check to get one vertice type
        for vertice_type in self.vertice_types:
            vertices = self.manager.get_vertices(types=vertice_type)
            self.assertEqual(len(vertices), 1)

    def test_get_graphs(self):
        """
        Test GraphManager.get_graphs
        """

        # check get all graphs
        graphs = set(self.manager.get_graphs())
        self.assertEqual(graphs, set(self.graphs))

        # check get all self graphs
        graphs = set(self.manager.get_graphs(ids=self.graph_ids))
        self.assertEqual(graphs, set(self.graphs))

        # check get one graph
        for graph_id in self.graph_ids:
            graph = self.manager.get_graphs(ids=graph_id)
            self.assertIn(graph, self.graphs)

        # check gelts
        for graph in self.graphs:
            graph_elts = graph.elts
            graph = self.manager.get_graphs(ids=graph.id, add_elts=True)
            self.assertEqual(len(graph_elts), len(graph._gelts))

        # check get graph types
        graphs = set(self.manager.get_graphs(types=self.graph_types))
        self.assertEqual(graphs, set(self.graphs))

        # check to get one graph type
        for graph_type in self.graph_types:
            graphs = self.manager.get_graphs(types=graph_type)
            self.assertEqual(len(graphs), 1)

        # check get graphs by elts
        graphs = set(self.manager.get_graphs(elts=self.elt_ids))
        self.assertEqual(graphs, set(self.graphs))

        # check get graphs by one vertice
        for index, vertice_id in enumerate(self.vertice_ids):
            graphs = self.manager.get_graphs(elts=vertice_id)
            self.assertEqual(len(graphs), index + 1)

        # check get graphs by one edge
        for index, edge_id in enumerate(self.edge_ids):
            graphs = self.manager.get_graphs(elts=edge_id)
            self.assertEqual(len(graphs), index + 1)

        # check get graphs by one graph
        for index, graph_id in enumerate(self.graph_ids):
            graphs = self.manager.get_graphs(elts=graph_id)
            self.assertEqual(len(graphs), index + 1)

        # check remove one elt
        for index, graph_id in enumerate(self.graph_ids):
            graph = self.manager.get_graphs(ids=graph_id)
            graph_elts = graph.elts
            if graph_elts:
                first_elt = graph.elts[0]
                self.manager.remove_elts(ids=first_elt, graph_ids=graph.id)
                graph = self.manager.get_graphs(ids=graph_id)
                self.assertEqual(len(graph.elts), len(graph_elts) - 1)
                if len(graph.elts) > 1:
                    # check remove two elts
                    self.manager.remove_elts(ids=graph.elts[:2])
                    graph = self.manager.get_graphs(ids=graph_id)
                    self.assertEqual(len(graph.elts), len(graph_elts) - 3)

        # check remove one elt from multiple graphs
        for index in range(len(self.graph_ids)):
            # get first "index" graph
            graph_ids = self.graph_ids[:index]
            graphs = self.manager.get_graphs(ids=graph_ids)
            graph_elts = graph.elts
            elt_id = self.vertice_ids[-1]
            self.manager.remove_elts(ids=elt_id, graph_ids=graph_ids)
            graphs = self.manager.get_graphs(ids=graph_ids)
            for graph in graphs:
                self.assertNotIn(elt_id, graph.elts)

    # TODO 4-01-2017
    #def test_get_neighbourhood(self):
    #    """
    #    Test get_neighbourhood method.
    #    """

    #    # test empty result
    #    neighbourhood = self.manager.get_neighbourhood(ids='')
    #    self.assertFalse(neighbourhood)

    #    # test all vertices
    #    neighbourhood = set(self.manager.get_neighbourhood())
    #    self.assertEqual(len(neighbourhood), len(self.elts))

    #    # let's play with few vertices, edges and graphs
    #    v0, v1, v2 = Vertice(type='0'), Vertice(type='1'), Vertice(type='2')
    #    e0, e1, e2 = Edge(type='0'), Edge(type='1'), Edge(type='2')
    #    g0, g1, g2 = Graph(type='0'), Graph(type='1'), Graph(type='2')

    #    # connect v0 to v1
    #    e0.directed = True
    #    e0.sources = [v0.id, v1.id, v2.id]
    #    e0.targets = [g0.id, g1.id, g2.id]

    #    # save all elements
    #    v0.save(self.manager), v1.save(self.manager), v2.save(self.manager)
    #    e0.save(self.manager), e1.save(self.manager), e2.save(self.manager)
    #    g0.save(self.manager), g1.save(self.manager), g2.save(self.manager)

    #    # test ids
    #    neighbourhood = set(self.manager.get_neighbourhood(ids=v0.id))
    #    self.assertEqual(neighbourhood, {g0, g1, g2})

    #    # test sources
    #    neighbourhood = set(self.manager.get_neighbourhood(
    #        ids=v0.id, sources=True))
    #    self.assertEqual(neighbourhood, {g0, g1, g2, v0, v1, v2})

    #    # test not targets
    #    neighbourhood = set(self.manager.get_neighbourhood(
    #        ids=v0.id, targets=False))
    #    self.assertEqual(neighbourhood, set())

    #    # test not directed
    #    e0.directed = False
    #    e0.save(self.manager)
    #    neighbourhood = set(self.manager.get_neighbourhood(
    #        ids=v0.id))
    #    self.assertEqual(neighbourhood, {g0, g1, g2, v0, v1, v2})

    #    # test info

    #    # test source_data

    #    # test target_data

    #    # test types

    #    # test source_types

    #    # test target_types

    #    # test edge_ids

    #    # test edge_types

    #    # test add_edges

    #    # test source_edge_types

    #    # test target_edge_types

    #    # test query

    #    # test edge_query

    #    # test source_query

    #    # test target_query

    def test_orphans(self):
        """
        Test get orphans method.
        """
        # check if no orphans exist
        orphans = self.manager.get_orphans()
        self.assertFalse(orphans)
        # generate self.count vertices and edges
        [Vertice().save(self.manager) for i in range(self.count)]
        [Edge().save(self.manager) for i in range(self.count)]
        # check if previous vertices and edges are orphans
        orphans = self.manager.get_orphans()
        self.assertEqual(len(orphans), 2 * self.count)
        # create a graph and add orphans to the graph
        graph = Graph()
        graph.add_elts(orphans)
        graph.save(manager=self.manager)
        # check if only the graph is an orphan
        orphans = self.manager.get_orphans()
        self.assertEqual(len(orphans), 1)
        # delete the graph and check if vertices and edges became orphans
        graph.delete(manager=self.manager, del_orphans=False)
        orphans = self.manager.get_orphans()
        self.assertEqual(len(orphans), 2 * self.count)


class PutEltsTest(TestCase):

    def setUp(self):

        self.manager = GraphManager(data_scope='graph_test')

    def tearDown(self):

        self.manager.del_elts()

    # TODO 4-01-2017
    #def test_put_elts(self):
    #    """
    #    Test put elts method.
    #    """

    #    vertice = Vertice()

    #    def assertVertice():
    #        """
    #        Assert vertice exists in DB
    #        """
    #        # get vertice
    #        _vertice = self.manager.get_elts(ids=vertice.id)
    #        # check ids
    #        self.assertEqual(_vertice.id, vertice.id)
    #        # delete vertice
    #        self.manager.del_elts()

    #    self.manager.put_elts(elts=vertice)
    #    assertVertice()

    #    self.manager.put_elts(elts=vertice.to_dict())
    #    assertVertice()

    #    self.manager.put_elts(elts=[vertice])
    #    assertVertice()

    #    self.manager.put_elts(elts=[vertice.to_dict()])
    #    assertVertice()


if __name__ == '__main__':
    main()
