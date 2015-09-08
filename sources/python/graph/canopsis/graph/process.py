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

"""
Module in charge of defining main graph task execution.

When an event occured, the related entity is retrieved with its bound
graph nodes in order to execute theirs tasks.
"""

from canopsis.graph.manager import Graph
from canopsis.context.manager import Context
from canopsis.task.core import run_task, register_task

context = Context()
graph = Graph()


@register_task('graph.event_processing')
def event_processing(event, ctx=None, **params):
    """Process input event in getting graph nodes bound to input event entity.

    If at least one graph nodes is found, execute its tasks.
    """

    nodes = []

    if ctx is None:
        ctx = {}

    entity = context.get_entity(event)

    if entity is not None:
        entity_id = context.get_entity_id(entity)
        nodes = graph.get_nodes(entity_id=entity_id)

        for node in nodes:
            task = node[Graph.TASK]
            if task is not None:
                run_task(
                    event=event,
                    task=task,
                    ctx=ctx,
                    node=node
                )

    return event
