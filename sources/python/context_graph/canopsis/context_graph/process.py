# -*- coding: utf-8 -*-

"""Module in chage of defining the graph context and updating
 it when it's needed."""

from __future__ import unicode_literals

from canopsis.task.core import register_task
from canopsis.context_graph.manager import ContextGraph
from canopsis.common.utils import singleton_per_scope

from canopsis.context_graph.manager import ContextGraph


context_graph_manager = ContextGraph()

@register_task
def event_processing(
        engine, event, manager=None, logger=None, ctx=None, tm=None, cm=None,
        **kwargs
):
    logger.error(event)
