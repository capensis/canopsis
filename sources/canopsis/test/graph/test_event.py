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
from canopsis.graph.event import TaskedVertice, TaskedEdge, TaskedGraph


class TaskedVerticeTest(TestCase):
    """
    Test event processing function.
    """

    def test_process_vertice(self):

        cls = TaskedVertice

        self._test_process_cls(cls)

    def test_process_edge(self):

        cls = TaskedEdge

        self._test_process_cls(cls)

    def test_process_graph(self):

        cls = TaskedGraph

        self._test_process_cls(cls)

    def _test_process_cls(self, cls):
        """
        Process a task which returns all elt data.
        """

        task_id = 'process'

        @register_task(task_id, force=True)
        def process_node(vertice, ctx, event=None, **kwargs):

            return vertice, ctx, kwargs

        ctx, entity, operator = {'b': 1}, 'e', task_id

        elt = cls(
            entity=entity, task=operator
        )

        _node, _ctx, _kwargs = elt.process(ctx=ctx, event=None)

        self.assertIs(_node, elt)
        self.assertIs(_ctx, ctx)
        self.assertFalse(_kwargs)


if __name__ == '__main__':
    main()
