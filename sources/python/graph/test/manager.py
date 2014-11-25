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


class GraphTest(TestCase):

    def setUp(self):
        """
        Create self.count=10 and self.graphs
        """

        self.manager = GraphManager()

        self.count = 10

        self.graphs = [None] * 10

        self.types = (i for i in range(self.count))

        # create self.count vertices
        self.vertices = [None] * self.count
        for index in range(self.count):
            self.vertices[index] = {
                'id': 'vertice-{0}'.format(index),
                'data': None if index % 2 else {'index': index},
                'type': 'tvertice-{0}'.format(index),
            }
        self.vertice_ids = [vertice['id'] for vertice in self.vertices]

        # create self.count edges
        self.edges = [None] * self.count
        for index in range(self.count):
            self.edges[index] = {
                'id': 'edge-{0}'.format(index),
                'data': None if index % 2 else {'index': index},
                'type': 'tedge-{0}'.format(index),
                'sources': self.vertice_ids[index:],
                'targets': self.vertice_ids[:-index],
                'directed': index % 2 == 0
            }
        self.edge_ids = [edge['id'] for edge in self.edges]

        # create self.count graphs
        self.graphs = [None] * self.count
        for index in range(self.count):
            self.graphs[index] = {
                'id': 'graph-{0}'.format(index),
                'data': None if index % 2 else {'index': index},
                'type': 'tgraph-{0}'.format(index),
                'elts': self.edge_ids[index:] + self.vertice_ids[index:]
            }
        self.graph_ids = [graph['id'] for graph in self.graphs]

        # get all elt ids
        self.elt_ids = self.graph_ids + self.vertice_ids + self.edge_ids

        # add edges and graphs in edges
        for index in range(self.count):
            edge = self.edges[index]
            edge['sources'] += self.edge_ids[index:] + self.graph_ids[index:]
            edge['targets'] += self.edge_ids[:-index] + self.graph_ids[:-index]

        # add graph_ids in graph
        for index in range(self.count):
            self.graphs[index]['elts'] += self.graph_ids[index:]

    def tearDown(self):
        """
        Del all elements
        """

        for graph in self.graphs:
            self.manager.del_elts(ids=self.elt_ids)

    def test_CRUD(self):
        """
        Test to put graph and get them back
        """

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
        graphs = self.maanger.get_elts(ids=self.graph_ids)
        self.assertEqual(len(graphs), self.count)

        # check for get all elts
        elts = self.manager.get_elts()
        self.assertEqual(len(elts), 3 * self.count)


if __name__ == '__main__':
    main()
