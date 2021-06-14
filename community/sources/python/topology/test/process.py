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

from canopsis.context.manager import Context
from canopsis.topology.elements import TopoNode, TopoEdge
from canopsis.topology.manager import TopologyManager
from canopsis.topology.process import event_processing
from canopsis.topology.rule.action import change_state
from canopsis.task.core import new_conf
from canopsis.check import Check

#TODO4-01-2017
#class ProcessingTest(TestCase):
#    """
#    Test event processing function.
#    """
#
#    class _Amqp(object):
#        """
#        In charge of processing publishing of test.
#        """
#
#        def __init__(self, processingTest):
#
#            self.exchange_name_events = None
#            self.processingTest = processingTest
#            self.event = None
#
#        def publish(self, event, rk, exchange):
#            """
#            Called when an event process publishes an event.
#            """
#
#            self.event = event
#            self.processingTest.count += 1
#
#            event_processing(
#                event=event,
#                engine=self.processingTest,
#                manager=self.processingTest.manager
#            )
#
#    def setUp(self):
#
#        self.context = Context(data_scope='test_context')
#        self.manager = TopologyManager(data_scope='test_topology')
#        self.check = {
#            'type': 'check',
#            'event_type': 'check',
#            'connector': 'c',
#            'connector_name': 'c',
#            'component': 'c',
#            'source_type': 'component',
#            'state': Check.OK
#        }
#        entity = self.context.get_entity(self.check)
#        entity_id = self.context.get_entity_id(entity)
#        self.node = TopoNode(entity=entity_id)
#        self.node.save(self.manager)
#        self.count = 0
#        self.amqp = ProcessingTest._Amqp(self)
#
#    def tearDown(self):
#
#        self.manager.del_elts()
#
#    def test_no_bound(self):
#        """
#        Test in case of not bound nodes.
#        """
#
#        event_processing(event=self.check, engine=self, manager=self.manager)
#        self.assertEqual(self.count, 0)
#
#    def test_one_node(self):
#        """
#        Test in case of one bound node
#        """
#
#        source = TopoNode()
#        source.save(self.manager)
#        edge = TopoEdge(sources=source.id, targets=self.node.id)
#        edge.save(self.manager)
#
#        event_processing(event=self.check, engine=self, manager=self.manager)
#        self.assertEqual(self.count, 0)
#
#    def test_change_state(self):
#        """
#        Test in case of change state.
#        """
#        # create a change state operation with minor state
#        change_state_conf = new_conf(
#            change_state,
#            state=Check.MINOR
#        )
#        self.node.operation = change_state_conf
#        self.node.save(self.manager)
#
#        self.node.process(event=self.check, manager=self.manager)
#        event_processing(event=self.check, engine=self, manager=self.manager)
#
#        target = self.manager.get_elts(ids=self.node.id)
#        self.assertEqual(target.state, Check.MINOR)
#
#    def test_chain_change_state(self):
#        """
#        Test to change of state in a chain of nodes.
#
#        This test consists to link three node in such way:
#        self.node(state=0) -> node(state=0) -> root(state=0)
#        And to propagate the change state task with state = 1 in order to check
#        if root state equals 1.
#        """
#
#        # create a simple task which consists to change of state
#        change_state_conf = new_conf(
#            change_state,
#            state=Check.MINOR
#        )
#
#        # create a root node with the change state task
#        root = TopoNode(operator=change_state_conf)
#        root.save(self.manager)
#        # create a node with the change state task
#        node = TopoNode(operator=change_state_conf)
#        node.save(self.manager)
#        # create a leaf with the change state task
#        self.node.operation = change_state_conf
#        self.node.save(self.manager)
#        # link node to root
#        rootnode = TopoEdge(targets=root.id, sources=node.id)
#        rootnode.save(self.manager)
#        # link self.node to node
#        self_node = TopoEdge(targets=node.id, sources=self.node.id)
#        self_node.save(self.manager)
#
#        event_processing(event=self.check, engine=self, manager=self.manager)
#        self.assertEqual(self.count, 3)
#
#        self.node = self.manager.get_elts(ids=self.node.id)
#        self.assertEqual(self.node.state, Check.MINOR)
#
if __name__ == '__main__':
    main()
