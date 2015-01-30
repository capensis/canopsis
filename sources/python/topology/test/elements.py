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

from canopsis.task import register_task
from canopsis.topology.elements import TopoNode, Topology
from canopsis.topology.manager import TopologyManager
from canopsis.context.manager import Context


class TopoNodeTest(TestCase):
    """
    Test event processing function.
    """

    def test_empty(self):
        """
        Test to process a toponode without task.
        """

        toponode = TopoNode()
        result = toponode.process(event=None)
        self.assertIsNone(result)

    def test_process_task(self):
        """
        Process a task which returns all toponode data.
        """

        @register_task('process')
        def process_node(toponode, ctx, event=None, **kwargs):

            return toponode, ctx, kwargs

        ctx, entity, state, operator = {'b': 1}, 'e', 0, 'process'

        toponode = TopoNode(state=state, entity=entity, operator=operator)

        _node, _ctx, _kwargs = toponode.process(ctx=ctx, event=None)

        self.assertIs(_node, toponode)
        self.assertIs(_ctx, ctx)
        self.assertFalse(_kwargs)


class TopologyGraphTest(TestCase):
    """
    Test topology element.
    """

    def setUp(self):

        self.context = Context(data_scope='test')
        self.manager = TopologyManager(data_scope='test')

    def tearDown(self):

        self.context.remove()
        self.manager.del_elts()

    def test_save(self):
        """
        Test if an entity exists after saving a topology.
        """
        _id = 'test'

        topology = Topology(_id=_id)
        topology.save(manager=self.manager, context=self.context)

        topology = self.context.get(_type=topology.type, names=_id)

        self.assertEqual(topology[Context.NAME], _id)

    def test_delete(self):
        """
        Test if topology nodes exist after deleting a topology.
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
