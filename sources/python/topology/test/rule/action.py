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

from canopsis.topology.elements import Node
from canopsis.topology.manager import TopologyManager
from canopsis.topology.rule.action import (
    change_state, state_from_sources, worst_state, best_state
)
from canopsis.graph.elements import Edge


tm = TopologyManager(data_scope='test_topology')


class ChangeStateTest(TestCase):
    """
    Test change state function.
    """

    def setUp(self):

        self.node = Node()
        self.assertEqual(Node.state(self.node), 0)
        self.new_state = 1

    def test_state(self):
        """
        Test to change of state from a state.
        """

        change_state(
            node=self.node, event={}, state=self.new_state, manager=tm
        )

    def test_event(self):
        """
        Test to change of state from an event.
        """

        event = {'state': self.new_state}
        change_state(node=self.node, event=event, manager=tm)

    def tearDown(self):

        node = tm.get_vertices(ids=self.node.id)
        self.assertEqual(Node.state(node), self.new_state)
        tm.del_elts()


class StateFromSourcesTest(TestCase):
    """
    Test to change state from sources function.
    """

    def setUp(self):

        # empty DB
        tm.del_elts()

        self.count = 5

    def tearDown(self):

        # empty DB
        tm.del_elts()

    def get_function(self):
        """
        Get change of state function
        """

        return state_from_sources

    def get_kwargs(self):

        return {'f': max}

    def get_new_state(self):
        """
        Get new state.
        """

        return self.count - 1

    def test_no_sources(self):
        """
        Test to change of state without sources.
        """

        node = Node()
        event = {}

        self.get_function()(
            node=node, event=event, manager=tm, ctx={}, **self.get_kwargs()
        )
        self.assertEqual(node.data['state'], 0)

    def test_sources(self):
        """
        Test to change of state with sources.
        """

        node = Node()
        event = {}

        count = 5

        sources = [Node(state=i) for i in range(count)]
        for source in sources:
            source.save(manager=tm)
        edge = Edge(
            targets=[node.id], sources=list(source.id for source in sources)
        )
        edge.save(manager=tm)
        self.get_function()(
            node=node, event=event, manager=tm, ctx={}, **self.get_kwargs()
        )
        self.assertEqual(node.data['state'], self.get_new_state())
        edge.delete(manager=tm)
        for source in sources:
            source.delete(manager=tm)


class WorstStateTest(StateFromSourcesTest):
    """
    Test to change state depending on worst source state.
    """

    def get_kwargs(self):

        return {}

    def get_function(self):

        return worst_state


class BestStateTest(StateFromSourcesTest):

    def get_kwargs(self):

        return {}

    def get_function(self):

        return best_state

    def get_new_state(self):

        return 0

if __name__ == '__main__':
    main()
