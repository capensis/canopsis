#!/usr/bin/env python
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

from unittest import TestCase, main

from canopsis.graph.manager import GraphManager
from canopsis.graph.elements import Graph, Vertice, Edge, GraphElement


class GraphTest(TestCase):

    def setUp(self):
        """
        Create self.count=10 and self.graphs
        """

        self.manager = GraphManager(data_scope='test_graph')

        self.count = 10
        self.graphs = [None] * 10
        self.types = (i for i in range(self.count))

        # create self.count vertices
        self.vertices = [None] * self.count
        for index in range(self.count):
            self.vertices[index] = Vertice(
                _id='vertice-{0}'.format(index),
                data=None if index % 2 else {'index': index},
                _type='tvertice-{0}'.format(index),
            )
        self.vertice_ids = [vertice.id for vertice in self.vertices]
        self.vertice_types = [vertice.type for vertice in self.vertices]

        # create self.count edges
        self.edges = [None] * self.count
        for index in range(self.count):
            self.edges[index] = Edge(
                _id='edge-{0}'.format(index),
                data=None if index % 2 else {'index': index},
                _type='tedge-{0}'.format(index),
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
                _id='graph-{0}'.format(index),
                data=None if index % 2 else {'index': index},
                _type='tgraph-{0}'.format(index),
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
        elts = self.manager.get_elts()
        self.assertFalse(elts)

        # put vertices in DB
        for vertice in self.vertices:
            self.manager.put_elt(elt=vertice)

        # put edges in DB
        for edge in self.edges:
            self.manager.put_elt(elt=edge)

        # put graphs in DB
        for graph in self.graphs:
            self.manager.put_elt(elt=graph)

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
            self.manager.put_elt(elt=vertice)
        for edge in self.edges:
            self.manager.put_elt(elt=edge)
        for graph in self.graphs:
            self.manager.put_elt(elt=graph)

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
        elts = self.manager.get_elts(
            ids=last_vertice_id, graph_ids=self.graph_ids)
        self.assertEqual(len(elts), 1)

        # check get data
        for elt in self.elts:
            _elt = self.manager.get_elts(ids=elt.id, data=elt.data)
            self.assertIsNotNone(_elt)
            self.assertEqual(_elt[GraphElement.ID], elt.id)
            if elt.data is not None:
                elts = self.manager.get_elts(data=elt.data)
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

        # check get edges by targets
        edges = set(self.manager.get_edges(targets=self.elt_ids))
        self.assertEqual(edges, set(self.edges))

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


if __name__ == '__main__':
    main()
