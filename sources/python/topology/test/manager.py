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

from canopsis.topology.manager import TopologyManager


class IsRootTest(TestCase):
    """
    Test TopologyManager.is_root function
    """

    def test_true(self):
        """
        Test if a root is recognized by the TopologyManager.
        """
        root = TopologyManager.new_node(
            graph_id='test', _type=TopologyManager.ROOT)

        self.assertTrue(TopologyManager.is_root(root))

    def test_false(self):
        """
        Test if a common node is not recognized such as a root by the
            TopologyManager.
        """

        root = TopologyManager.new_node(graph_id='test')

        self.assertFalse(TopologyManager.is_root(root))


class IsClusterTest(TestCase):
    """
    Test TopologyManager.is_cluster function
    """

    def test_true(self):
        """
        Test if a cluster is recognized by the TopologyManager.
        """
        cluster = TopologyManager.new_node(
            graph_id='test', _type=TopologyManager.CLUSTER)

        self.assertTrue(TopologyManager.is_cluster(cluster))

    def test_false(self):
        """
        Test if a common node is not recognized such as a cluster by the
            TopologyManager.
        """

        cluster = TopologyManager.new_node(graph_id='test')

        self.assertFalse(TopologyManager.is_cluster(cluster))


class NewGraphTest(TestCase):
    """
    Test Topology.new_graph method.
    """

    def test_without_root_and_nodes(self):
        """
        Test when nodes is not in parameters.
        """

        graph = TopologyManager.new_graph()

        nodes = graph[TopologyManager.NODES]

        self.assertEqual(len(nodes), 1)

        root = nodes[0]

        self.assertTrue(TopologyManager.is_root(root))

    def test_without_root(self):
        """
        Test when root is not in parameter nodes.
        """

        node = TopologyManager.new_node(graph_id='test')

        nodes = [node]

        graph = TopologyManager.new_graph(nodes=nodes)

        nodes = graph[TopologyManager.NODES]

        self.assertEqual(len(nodes), 2)

        root = nodes[0]

        self.assertTrue(TopologyManager.is_root(root))

    def test_with_root(self):
        """
        Test when root is in nodes.
        """

        nodes = [TopologyManager.new_node(graph_id='test') for i in range(10)]

        root = TopologyManager.new_node(
            graph_id='test', _type=TopologyManager.ROOT)

        nodes.insert(5, root)

        graph = TopologyManager.new_graph(nodes=nodes)

        nodes = graph[TopologyManager.NODES]

        self.assertEqual(len(nodes), 11)

        root = nodes[0]

        self.assertTrue(TopologyManager.is_root(root))


class NewNodeTest(TestCase):
    """
    Test TopologyManager.new_node
    """

    def test_without_state_weight(self):

        node = TopologyManager.new_node(graph_id='test')

        self.assertIn(TopologyManager.DATA, node)

        data = node[TopologyManager.DATA]

        node_state = data[TopologyManager.STATE]
        self.assertEqual(node_state, TopologyManager.DEFAULT_STATE)

        node_weight = data[TopologyManager.WEIGHT]
        self.assertEqual(node_weight, TopologyManager.DEFAULT_WEIGHT)

    def test_with_state_weight(self):

        state, weight = 50, 51

        node = TopologyManager.new_node(
            graph_id='test', state=state, weight=weight)

        self.assertIn(TopologyManager.DATA, node)

        data = node[TopologyManager.DATA]

        node_state = data[TopologyManager.STATE]
        self.assertEqual(state, node_state)

        node_weight = data[TopologyManager.WEIGHT]
        self.assertEqual(weight, node_weight)


class PutGraphTest(TestCase):
    """
    Test put_graph method
    """
    def setUp(self):
        """
        Instantiate a new TopologyManager
        """

        self.id = 'test'
        self.tm = TopologyManager()

    def tearDown(self):
        """
        Delete self.id graph.
        """

        self.tm.del_graph(ids=self.id)

    def test_without_root(self):
        """
        Test to put a graph without root.
        """

        graph = TopologyManager.new_graph(_id=self.id)

        graph[TopologyManager.NODES] = []

        self.tm.put_graph(graph=graph)

        graph = self.tm.get_graph(_id=self.id)

        nodes = graph[TopologyManager.NODES]

        self.assertEqual(len(nodes), 1)

        root = nodes[0]

        self.assertTrue(TopologyManager.is_root(root))

    def test_with_root(self):
        """
        Test to put a graph with root
        """
        root = TopologyManager.new_node(
            graph_id=self.id, _type=TopologyManager.ROOT)

        nodes = [TopologyManager.new_node(graph_id=self.id) for i in range(10)]

        nodes.insert(5, root)

        graph = TopologyManager.new_graph(_id=self.id, nodes=nodes)

        self.tm.put_graph(graph=graph)

        graph = self.tm.get_graph(_id=self.id)

        nodes = graph[TopologyManager.NODES]

        self.assertEqual(len(nodes), 11)

        root = nodes[0]

        self.assertTrue(TopologyManager.is_root(root))


if __name__ == '__main__':
    main()
