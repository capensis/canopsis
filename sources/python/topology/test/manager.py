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

from canopsis.topology.elements import Topology, Node
from canopsis.topology.manager import TopologyManager


class GetNodesTest(TestCase):
    """
    Test TopologyManager.get_nodes method.
    """

    def setUp(self):
        """
        Initialize a manager.
        """

        self.manager = TopologyManager(data_scope='test_topology')
        self.manager.del_elts()

    def tearDown(self):
        """
        Empty DB.
        """

        self.manager.del_elts()

    def test_empty(self):
        """
        Test to get nodes from an not existing entity.
        """

        nodes = self.manager.get_nodes(entity='')
        self.assertFalse(nodes)

    def test_one_node(self):
        """
        Test to get one node from an existing entity.
        """

        node = Node(entity='test')
        node.save(manager=self.manager)
        nodes = self.manager.get_nodes(entity='test')
        self.assertEqual(len(nodes), 1)
        node.delete(manager=self.manager)

    def test_nodes(self):
        """
        Test to get nodes from an existing entity.
        """

        count = 10
        nodes = [None] * count
        for i in range(count):
            node = Node(entity='test')
            node.save(manager=self.manager)
            nodes[i] = node
        _nodes = self.manager.get_nodes(entity='test')
        self.assertEqual(len(_nodes), count)
        for node in nodes:
            node.delete(manager=self.manager)

if __name__ == '__main__':
    main()
