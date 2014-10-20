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

from canopsis.topology.manager import Topology


class TopologyTest(TestCase):

    def setUp(self):
        self.topology = Topology(data_scope='test')

        # set of default topologies
        self.count = 10
        self.topologies = ({Topology.ID: str(i)} for i in range(self.count))
        # delete all topologies
        self.topology.delete()
        # check if all topologies are gone
        topologies = self.topology.get_topologies()
        self.assertFalse(topologies)

    def test_CRUD(self):
        # start to push topologies
        for topology in self.topologies:
            self.topology.push(topology=topology)

        # ensure topology count is equal to pushed topologies
        topologies = self.topology.get_topologies()
        self.assertEqual(len(topologies), self.count)

        # delete the first topology
        self.topology.delete(ids='0')
        topologies = self.topology.get_topologies()
        self.assertEqual(len(topologies), self.count - 1)

        # delete third more
        self.topology.delete(
            ids=('-1', '2', '3'))
        topologies = self.topology.get_topologies()
        self.assertEqual(len(topologies), self.count - 3)

    def test_get_topologies(self):
        topologies = self.topology.get_topologies()
        self.assertFalse(topologies)

    def test_find(self):
        for topology in self.topologies:
            self.topology.push(topology=topology)

        # find the topology where id is 3
        topologies = self.topology.find(regex='3')

        self.assertEqual(len(topologies), 1)

    def tearDown(self):
        # remove topologies from database
        self.topology.delete()


class TopologyNodeTest(TestCase):

    def setUp(self):
        self.topology = Topology(data_scope='test')
        self.topology.delete_nodes()
        self.count = 10
        self.nodes = (
            {
                Topology.ID: str(i),
                Topology.ENTITY_ID: str(i % 2),
                Topology.NEXT: str(i % 3)
            }
            for i in range(self.count)
        )
        for node in self.nodes:
            self.topology.push_node(node)

    def test_CRUD(self):
        nodes = self.topology.get_nodes()
        self.assertEqual(len(nodes), self.count)

        self.topology.delete_nodes(ids='0')
        nodes = self.topology.get_nodes()
        self.assertEqual(len(nodes), self.count - 1)

        self.topology.delete_nodes(ids=('-1', '1', '2'))
        nodes = self.topology.get_nodes()
        self.assertEqual(len(nodes), self.count - 3)

    def test_bounds(self):
        nodes = self.topology.find_bound_nodes(entity_id='1')
        self.assertEqual(len(nodes), self.count / 2)

    def test_nexts(self):
        for node in self.nodes:
            next_node = self.topology.get_next_nodes(node=node)
            self.assertEqual(
                int(node[Topology.ID]), int(next_node[Topology.ID] % 3))

    def test_sources(self):
        for node in self.nodes:
            source_node = self.topology.find_source_nodes(node=node)
            self.assertEqual(
                int(node[Topology.ID]), int(source_node[Topology.ID]) + 1 % 3)

    def tearDown(self):
        self.topology.delete_nodes()


if __name__ == '__main__':
    main()
