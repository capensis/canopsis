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

from canopsis.task import register_task
from canopsis.topology.elements import Node


class NodeTest(TestCase):
    """
    Test event processing function.
    """

    def test_empty(self):
        """
        Test to process a node without task.
        """

        node = Node()
        result = node.process()
        self.assertIsNone(result)

    def test_process_task(self):
        """
        Process a task which returns all node data.
        """

        @register_task('process')
        def process_node(node, ctx, **kwargs):

            return node, ctx, kwargs

        ctx, entity, state, task, weight = {'b': 1}, 'e', 0, 'process', 1

        node = Node(task=task, entity=entity, state=state, weight=weight)

        _node, _ctx, _kwargs = node.process(ctx=ctx)

        self.assertIs(_node, node)
        self.assertIs(_ctx, ctx)
        self.assertFalse(_kwargs)

if __name__ == '__main__':
    main()
