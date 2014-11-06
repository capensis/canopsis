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

        self.graph = GraphManager()

        self.count = 10

        self.graphs = [None] * 10

        self.types = [i for i in range(self.count)]

        for index in range(len(self.graphs)):
            graph_id = '__test__%s' % index
            # create count-index nodes per graph
            nodes = [
                GraphManager.new_node(
                    graph_id=graph_id,
                    _id='%s%s' % (graph_id, node_index),
                    entity_id='' if (node_index % 2)
                        else '%s%s' % (graph_id, ((index + 1) % self.count)),
                    _type=index,
                    data=None if (node_index % 2) else node_index
                ) for node_index in range(index, self.count)
            ]
            # create count-index edges where sources equal targets
            edges = [
                GraphManager.new_edge(
                    graph_id=graph_id,
                    _id='%s%s-' % (graph_id, edge_index),
                    sources=[
                        '%s%s' % (graph_id, ei)
                        for ei in range(index, edge_index)],
                    targets=[
                        '%s%s' % (graph_id, ei)
                        for ei in range(index, edge_index)],
                    entity_id='' if (edge_index % 2)
                        else '%s%s' % (graph_id, ((index + 1) % self.count)),
                    _type=index,
                    data=None if (edge_index % 2) else edge_index,
                    directed=edge_index % 2
                ) for edge_index in range(index, self.count)
            ]
            # add edges in nodes
            nodes += edges
            # create a graph
            graph = GraphManager.new_graph(_id=graph_id, nodes=nodes)
            # put graph in self.graphs
            self.graphs[index] = graph

    def tearDown(self):
        """
        Del self.graphs
        """

        self.graph.del_graph(
            ids=[g[GraphManager.ID] for g in self.graphs])

    def test_get_graph_which_does_not_exist(self):
        """
        Test to get not existing graphs
        """

        # Test to get graph one by one
        for graph in self.graphs:
            graph = self.graph.get_graph(_id=graph[GraphManager.ID])
            self.assertIsNone(graph)

    def test_del_graph_which_does_not_exist(self):
        """
        Test to delete not existing graphs
        """

        # test elementary calls to del_graph
        for graph in self.graphs:
            self.graph.del_graph(ids=graph[GraphManager.ID])

        # test to del all graphs in one call
        self.graph.del_graph(ids=[g[GraphManager.ID] for g in self.graphs])

    def test_get_node_which_do_not_exists(self):
        """
        Test to get nodes which do not exists.
        """

        nodes = self.graph.get_nodes(ids='')

        self.assertFalse(nodes)

        nodes = self.graph.get_nodes(ids=[''])

        self.assertFalse(nodes)

    def _compare_nodes(self, nodes, _nodes):

        node_ids = [node[GraphManager.ID] for node in nodes]
        _node_ids = [node[GraphManager.ID] for node in _nodes]

        self.assertEqual(len(nodes), len(_nodes))

        for node_id in node_ids:
            self.assertIn(node_id, _node_ids)

    def test_CRUD(self):
        """
        Test to put graph and get them back
        """

        # starts to put graphs
        for graph in self.graphs:
            self.graph.put_graph(graph=graph)

        # check equality
        for index, graph in enumerate(self.graphs):

            # get graph id
            graph_id = graph[GraphManager.ID]
            _graph = self.graph.get_graph(_id=graph_id)
            # assert equality between graph and DB graph
            self.assertEqual(graph[GraphManager.ID], _graph[GraphManager.ID])

            # compare graph nodes
            nodes = graph[GraphManager.NODES]
            _nodes = _graph[GraphManager.NODES]
            self._compare_nodes(nodes, _nodes)

            # get node ids
            node_ids = [node[GraphManager.ID] for node in nodes]

            # ensure get_nodes equals nodes
            _nodes = self.graph.get_nodes(graph_id=graph_id)
            self._compare_nodes(nodes, _nodes)

            # ensure get_nodes equals nodes
            _nodes = self.graph.get_nodes(ids=node_ids)
            self._compare_nodes(nodes, _nodes)

            # delete all nodes
            for node in nodes:
                self.graph.del_nodes(ids=node[GraphManager.ID])

            # ensure get_nodes is empty
            _nodes = self.graph.get_nodes(ids=node_ids)
            self.assertFalse(_nodes)

            # put nodes one by one
            for node in nodes:
                self.graph.put_nodes(node)

            _nodes = self.graph.get_nodes(ids=node_ids)
            self._compare_nodes(nodes, _nodes)

            # delete all nodes at a time
            self.graph.del_nodes(ids=node_ids)
            # ensure get_nodes is empty
            _nodes = self.graph.get_nodes(ids=node_ids)
            self.assertFalse(_nodes)

            self.graph.put_nodes(nodes)
            # ensure get_nodes equals nodes
            _nodes = self.graph.get_nodes(ids=node_ids)
            self._compare_nodes(nodes, _nodes)

            # assert type
            for i in range(self.count):
                _nodes = self.graph.get_nodes(_type=i)
                self.assertEqual(len(_nodes), self.count * 2 - 2 * i)

            # assert entity_id
            for node in nodes:
                entity_id = node[GraphManager.ENTITY_ID]
                _nodes = self.graph.get_nodes(entity_id=entity_id)
                for _node in _nodes:
                    self.assertEqual(_node[GraphManager.ENTITY_ID], entity_id)

            node_ids_set = set(node_ids)  # for sources/targets

            # assert sources
            for node_index, node in enumerate(nodes):
                node_id = node[GraphManager.ID]
                _nodes = self.graph.get_nodes(sources=node_id)
                targets = []
                for _node in _nodes:
                    if self.graph.is_edge(_node):
                        self.assertIn(node_id, _node[GraphManager.SOURCES])
                        targets += _node[GraphManager.TARGETS]
                    else:
                        self.assertIn(_node[GraphManager.ID], targets)

            _nodes = self.graph.get_nodes(sources=node_ids)

            targets = []
            for _node in _nodes:
                if self.graph.is_edge(_node):
                    sources = set(_node[GraphManager.SOURCES])
                    self.assertTrue(sources | node_ids_set)
                    targets += _node[GraphManager.TARGETS]
                else:
                    self.assertIn(_node[GraphManager.ID], targets)

            # assert targets
            for node_index, node in enumerate(nodes):
                node_id = node[GraphManager.ID]
                _nodes = self.graph.get_nodes(targets=node_id)
                targets = []
                for _node in _nodes:
                    if self.graph.is_edge(_node):
                        self.assertIn(node_id, _node[GraphManager.TARGETS])
                        targets += _node[GraphManager.SOURCES]
                    else:
                        self.assertIn(_node[GraphManager.ID], targets)

            _nodes = self.graph.get_nodes(targets=node_ids)

            sources = []
            for _node in _nodes:
                if self.graph.is_edge(_node):
                    targets = set(_node[GraphManager.TARGETS])
                    self.assertTrue(targets | node_ids_set)
                    sources += _node[GraphManager.SOURCES]
                else:
                    self.assertIn(_node[GraphManager.ID], sources)


if __name__ == '__main__':
    main()
