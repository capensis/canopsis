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

from canopsis.task.core import register_task
from canopsis.topology.elements import TopoNode, Topology, TopoEdge
from canopsis.topology.manager import TopologyManager
from canopsis.context.manager import Context


class TopoNodeTest(TestCase):
    """Test event processing function.
    """

    def setUp(self):

        self.manager = TopologyManager(data_scope='test')

    def tearDown(self):

        self.manager.del_elts()

    def test_default(self):
        """Test to process a toponode without default task.
        """

        toponode = TopoNode()
        self.assertEqual(toponode.state, 0)
        event = {'state': 1}
        result = toponode.process(event=event, manager=self.manager)
        self.assertIsNone(result)
        self.assertEqual(toponode.state, 1)

    def test_process_task(self):
        """Process a task which returns all toponode data.
        """

        @register_task('process')
        def process_node(
            vertice, ctx,
            event=None, publisher=None, source=None, manager=None, logger=None,
            **kwargs
        ):

            return vertice, ctx, kwargs

        ctx, entity, state, operation = {'b': 1}, 'e', 0, 'process'

        toponode = TopoNode(state=state, entity=entity, operation=operation)

        _node, _ctx, _kwargs = toponode.process(
            ctx=ctx, event=None, manager=self.manager
        )

        self.assertIs(_node, toponode)
        self.assertIs(_ctx, ctx)
        self.assertFalse(_kwargs)

    def test_proccess_task_with_propagation(self):
        """Process a node and check if node and edge states have changed.
        """

        state = 0
        new_state = state + 1
        operation = {
            'id': 'canopsis.topology.rule.action.change_state',
            'params': {'state': new_state}
        }

        toponode = TopoNode(state=state, operation=operation)
        toponode.save(manager=self.manager)
        edge = TopoEdge(sources=toponode)
        edge.save(manager=self.manager)

        self.assertEqual(toponode.state, state)
        self.assertEqual(edge.state, state)

        toponode.process(event={}, manager=self.manager)
        elts = self.manager.get_elts(ids=[toponode.id, edge.id])
        toponode = elts[0]
        edge = elts[1]
        self.assertEqual(toponode.state, new_state)
        self.assertEqual(edge.state, new_state)


class TopologyGraphTest(TestCase):
    """Test topology element.
    """

    def setUp(self):

        self.context = Context(data_scope='test')
        self.manager = TopologyManager(data_scope='test')

    def tearDown(self):

        self.context.remove()
        self.manager.del_elts()

    def test_save(self):
        """Test if an entity exists after saving a topology.
        """
        id = 'test'

        topology = Topology(id=id)
        topology.save(manager=self.manager, context=self.context)

        topology = self.context.get(_type=topology.type, names=id)

        self.assertEqual(topology[Context.NAME], id)

    def test_delete(self):
        """Test if topology nodes exist after deleting a topology.
        """

        topology = Topology()
        node = TopoNode()
        topology.add_elts(node)

        topology.save(manager=self.manager)

        node = self.manager.get_elts(node.id)
        self.assertIsNotNone(node)

        topology.delete(manager=self.manager)
        node = self.manager.get_elts(node.id)
        self.assertIsNone(node)


if __name__ == '__main__':
    main()
