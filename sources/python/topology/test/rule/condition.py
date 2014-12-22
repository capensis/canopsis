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

from canopsis.check import Check
from canopsis.topology.manager import TopologyManager
from canopsis.topology.elements import Node, Edge
from canopsis.topology.rule.condition import new_state, at_least, _all, nok


class NewStateTest(TestCase):
    """
    Test new_state function.
    """

    def setUp(self):

        self.node = Node()

    def test_new_state_from_event(self):
        """
        Test if state is different from an event.
        """

        result = new_state(node=self.node, event={'state': 1})

        self.assertTrue(result)

    def test_same_state_from_event(self):
        """
        Test if state is the same than the one from an event.
        """

        result = new_state(node=self.node, event={'state': 0})

        self.assertFalse(result)

    def test_new_state_from_state(self):
        """
        Test if state is different than the input state.
        """

        result = new_state(node=self.node, event={}, state=1)

        self.assertTrue(result)

    def test_same_state_from_state(self):
        """
        Test if state is the same than the input state.
        """

        result = new_state(node=self.node, event={}, state=0)

        self.assertFalse(result)


class _AtLeastTest(TestCase):
    """
    Default test class for at least function.
    """

    def setUp(self):

        self.manager = TopologyManager()

    def tearDown(self):
        """
        Drop nodes and edges.
        """

        self.manager.del_elts()


class AtLeastTest(_AtLeastTest):
    """
    Test at least test.
    """

    def test_empty(self):
        """
        Test to process at least with no children.
        """

        target = Node()
        target.save(self.manager)
        check = at_least(event={}, ctx={}, node=target)

        self.assertFalse(check)

    def test_default(self):
        """
        Test to check default condition.
        """

        source = Node()
        source.save(self.manager)
        target = Node()
        target.save(self.manager)
        edge = Edge(sources=source.id, targets=target.id)
        edge.save(self.manager)

        check = at_least(event={}, ctx={}, node=target)

        self.assertTrue(check)

    def test_false(self):
        """
        Test to check if there are at least one
        """

        source = Node()
        source.save(self.manager)
        target = Node()
        target.save(self.manager)
        edge = Edge(sources=source.id, targets=target.id)
        edge.save(self.manager)

        check = at_least(event={}, ctx={}, state=1, node=target)
        self.assertFalse(check)

        edge.data = {Node.WEIGHT: 0.5}
        edge.save(self.manager)
        source.data[Check.STATE] = 1
        source.save(self.manager)

        check = at_least(event={}, ctx={}, state=1, node=target)
        self.assertFalse(check)

        edge.data[Node.WEIGHT] = 1.5
        edge.save(self.manager)

        check = at_least(event={}, ctx={}, state=1, node=target)
        self.assertTrue(check)


class AllTest(_AtLeastTest):
    """
    Test _all function.
    """

    def test_one(self):
        """
        Test one source.
        """

        source = Node()
        source.save(self.manager)
        target = Node()
        target.save(self.manager)
        edge = Edge(sources=source.id, targets=target.id)
        edge.save(self.manager)

        check = _all(event={}, ctx={}, node=target)
        self.assertTrue(check)

        source.data[Check.STATE] = 1
        source.save(self.manager)

        check = _all(event={}, ctx={}, node=target)
        self.assertFalse(check)

    def test_many(self):

        target = Node()
        target.save(self.manager)
        count = 5
        sources = [Node() for i in range(count)]
        for source in sources:
            source.save(self.manager)
        edge = Edge(
            sources=[source.id for source in sources], targets=target.id
        )
        edge.save(self.manager)

        check = _all(event={}, ctx={}, node=target)
        self.assertTrue(check)

        sources[0].data[Check.STATE] = 1
        sources[0].save(self.manager)

        check = _all(event={}, ctx={}, node=target)
        self.assertFalse(check)


class NOKTest(_AtLeastTest):
    """
    Test nok test.
    """

    def test_one(self):
        """
        Test one source.
        """

        source = Node()
        source.save(self.manager)
        target = Node()
        target.save(self.manager)
        edge = Edge(sources=source.id, targets=target.id)
        edge.save(self.manager)

        check = nok(event={}, ctx={}, node=target)
        self.assertFalse(check)

        source.data[Check.STATE] = 1
        source.save(self.manager)

        check = nok(event={}, ctx={}, node=target)
        self.assertTrue(check)

    def test_many(self):

        target = Node()
        target.save(self.manager)
        count = 5
        sources = [Node() for i in range(count)]
        for source in sources:
            target.save(self.manager)
        edge = Edge(
            sources=[source.id for source in sources], targets=target.id
        )
        edge.save(self.manager)

        check = nok(event={}, ctx={}, node=target, min_weight=count)
        self.assertFalse(check)

        sources[0].data[Check.STATE] = 1
        sources[0].save(self.manager)

        check = nok(event={}, ctx={}, node=target)
        self.assertTrue(check)
        check = nok(event={}, ctx={}, node=target, min_weight=count)
        self.assertFalse(check)

        for source in sources:
            source.data[Check.STATE] = 1
            source.save(self.manager)

        check = nok(event={}, ctx={}, node=target, min_weight=count)
        self.assertTrue(check)


if __name__ == '__main__':
    main()
