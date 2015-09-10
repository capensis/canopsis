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
from canopsis.graph.elements import GraphElement, Vertice, Edge, Graph


class GraphElementTest(TestCase):
    """
    Test graph element.
    """

    __COUNT__ = 10  #: number of graph elements to create

    def get_type(self):
        """
        Get graph element type.
        """

        return GraphElement

    def get_params(self):
        """
        Get graph element kwargs.
        """

        return {}

    def setUp(self):
        """
        Create count graph elements and a manager
        """

        self.count = self.__COUNT__
        self.elts = [None] * self.count
        self.manager = GraphManager(data_scope='test_graph')
        self.manager.del_elts()
        for i in range(self.count):
            params = self.get_params()
            self.elts[i] = self.get_type()(**params)
        self.elt_ids = [elt.id for elt in self.elts]

    def tearDown(self):
        """
        Delete all managers.
        """

        for elt in self.elts:
            elt.delete(manager=self.manager)
        self.manager.del_elts()

    def test_init(self):
        """
        Test initialization
        """

        # test unique ids
        self.assertEqual(len(set(self.elt_ids)), self.count)

        # test same types
        self.assertEqual(len(set(elt.type for elt in self.elts)), 1)

        # change ids and types
        for index, elt in enumerate(self.elts):
            elt.id, elt.type = '0', index

        # test same ids
        self.assertEqual(len(set(elt.id for elt in self.elts)), 1)
        # test different types
        self.assertEqual(len(set(elt.type for elt in self.elts)), self.count)

    def test_save(self):
        """
        Test to save elements.
        """

        # ensure generated elements are not in the manager
        elts = self.manager.get_elts(ids=self.elt_ids)
        self.assertFalse(elts)

        # save all managers
        for elt in self.elts:
            elt.save(manager=self.manager)

        # ensure all can be retrieved
        elts = self.manager.get_elts(ids=self.elt_ids)
        self.assertEqual(len(elts), self.count)

    def test_delete(self):

        # ensure generated elements are not in the manager
        elts = self.manager.get_elts(ids=self.elt_ids)
        self.assertFalse(elts)

        # save all managers
        for elt in self.elts:
            elt.save(manager=self.manager)

        # ensure generated elements are in the manager
        elts = self.manager.get_elts(ids=self.elt_ids)
        self.assertEqual(len(elts), len(self.elts))

        for elt in self.elts:
            elt.delete(manager=self.manager)

        # ensure generated elements are not in the manager
        elts = self.manager.get_elts(ids=self.elt_ids)
        self.assertFalse(elts)

    def test_resolve_refs(self):
        """
        Test resolve_refs.
        """

    def test_serialization(self):
        """
        Test serializations
        """

        # assert serialization are ok
        for elt in self.elts:
            elt_dict = elt.to_dict()
            new_elt = self.get_type()(**elt_dict)
            self.assertEqual(new_elt, elt)
            new_elt_dict = new_elt.to_dict()
            self.assertEqual(elt_dict, new_elt_dict)


class VerticeTest(GraphElementTest):
    """
    Test vertice.
    """

    def get_type(self):

        return Vertice


class EdgeTest(VerticeTest):
    """
    Test edge.
    """

    def get_type(self):

        return Edge


class GraphTest(VerticeTest):
    """
    Test graph.
    """

    def get_type(self):

        return Graph

    def test_delete_orphans(self):
        """
        Test if graph vertices exist after deleting a graph.
        """

        graph = Graph()
        vertice = Vertice()
        graph.add_elts(vertice)

        graph.save(manager=self.manager)

        vertice = self.manager.get_elts(vertice.id)
        self.assertIsNotNone(vertice)

        graph.delete(manager=self.manager)
        vertice = self.manager.get_elts(vertice.id)
        self.assertIsNone(vertice)

if __name__ == '__main__':
    main()
