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

"""Module in charge of defining main graph task execution.

When an event occured, the related entity is retrieved with its bound
graph nodes in order to execute theirs tasks.
"""

from canopsis.graph.elements import Vertice, Edge, Graph
from canopsis.graph.manager import GraphManager
from canopsis.context.manager import Context
from canopsis.task import register_task, run_task
from canopsis.event import forger, Event
from canopsis.check import Check
from canopsis.common.utils import singleton_per_scope


class BaseTaskedVertice(object):

    TASK = 'task'  #: task field name in info
    ENTITY = Event.ENTITY  #: entity field name in info
    DEFAULT_TASK = 'canopsis.topology.rule.action.change_state'
    NAME = 'name'  #: element name.

    def get_default_task(self):
        """Get default task.
        """

        raise NotImplementedError()

    @property
    def entity(self):
        """Get self entity id.

        :return: self entity id.
        :rtype: str
        """

        return self.info.get(BaseTaskedVertice.ENTITY)

    @entity.setter
    def entity(self, value):
        """Change of entity id and update state.

        :param value: new entity (id) to use.
        :type value: dict or str
        """

        if value is not None:
            if isinstance(value, dict):
                # get entity id
                ctx = singleton_per_scope(Context)
                value = ctx.get_entity_id(value)

            # update entity
            self.info[BaseTaskedVertice.ENTITY] = value
        # call specific set entity
        self.set_entity(value)

    def set_entity(self, entity_id):
        """Specific setting of entity.

        :param str entity_id: new entity id. If None, entity_id is removed.
        """

        pass

    @property
    def task(self):
        """Get self task or default task if task is not setted.
        """

        result = self.info.get(BaseTaskedVertice.TASK)

        if not result:
            result = self.get_default_task()

        return result

    @task.setter
    def task(self, value):
        """Change of task.

        :param value: new task to use.
        """

        if value is not None:
            self.info[BaseTaskedVertice.TASK] = value

    def process(self, event, **kwargs):
        """Process this vertice task in a context of event processing.
        """

        result = None

        task = self.task

        result = run_task(conf=task, vertice=self, event=event, **kwargs)

        return result

    def get_event(self, *args, **kwargs):
        """Get vertice event.

        :param args: event forging args.
        :param kwargs: event forging kwargs.
        """

        result = forger(
            event_type=Check.EVENT_TYPE,
            source_type=self.type,
            component=self.id,
            *args, **kwargs
        )

        return result


class TaskedVertice(Vertice, BaseTaskedVertice):

    def __init__(self, task=None, entity=None, *args, **kwargs):

        super(TaskedVertice, self).__init__(*args, **kwargs)

        if self.info is None:
            self.info = {}

        self.task = task
        self.entity = entity


class TaskedEdge(Edge, BaseTaskedVertice):

    def __init__(self, task=None, entity=None, *args, **kwargs):

        super(TaskedEdge, self).__init__(*args, **kwargs)

        if self.info is None:
            self.info = {}

        self.task = task
        self.entity = entity


class TaskedGraph(Graph, BaseTaskedVertice):

    def __init__(self, task=None, entity=None, *args, **kwargs):

        super(TaskedGraph, self).__init__(*args, **kwargs)

        if self.info is None:
            self.info = {}

        self.task = task
        self.entity = entity


@register_task()
def event_processing(event, ctx=None, cm=None, gm=None, *args, **kwargs):
    """Process input event in getting graph nodes bound to input event entity.

    If at least one graph node is found, execute its tasks.

    :param Context cm:
    :param GraphManager gm:
    """

    if ctx is None:
        ctx = {}

    if cm is None:
        cm = singleton_per_scope(Context)
    if gm is None:
        gm = singleton_per_scope(GraphManager)

    entity = cm.get_entity(event)

    if entity is not None:
        entity_id = cm.get_entity_id(entity)
        vertices = gm.get_elts(
            info={BaseTaskedVertice.ENTITY: entity_id},
            cls=BaseTaskedVertice
        )

        for vertice in vertices:
            vertice.process(event=event, ctx=ctx, *args, **kwargs)

    return event
