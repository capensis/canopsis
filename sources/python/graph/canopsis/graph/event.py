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

from canopsis.graph.elements import Vertice
from canopsis.graph.manager import Graph
from canopsis.context.manager import Context
from canopsis.task import register_task, run_task
from canopsis.event import forger, Event

_context = Context()
graph = Graph()


class TaskedVertice(object):

    TASK = 'task'  #: task field name in data
    ENTITY = 'entity'  #: entity field name in data
    DEFAULT_TASK = 'canopsis.topology.rule.action.change_state'
    NAME = 'name'  #: element name.

    def get_default_task(self):
        """Get default task.
        """

        return {}

    @property
    def name(self):
        return self.data.get(TaskedVertice.NAME, self.id)

    @name.setter
    def name(self, value):
        self.data[TaskedVertice.NAME] = value

    @property
    def entity(self):
        """Get self entity id.

        :return: self entity id.
        :rtype: str
        """

        return self.data[TaskedVertice.ENTITY]

    @entity.setter
    def entity(self, value):
        """Change of entity id and update state.

        :param value: new entity (id) to use.
        :type value: dict or str
        """

        if isinstance(value, dict):
            # get entity id
            entity_id = _context.get_entity_id(value)
        else:
            entity_id = value

        # update entity
        self.data[TaskedVertice.ENTITY] = entity_id
        # call specific set entity
        self.set_entity(entity_id)

    def set_entity(self, entity_id):
        """Specific setting of entity.
        """

        pass

    def get_context_w_entity(self):
        """Get self entity structure and its context.

        :return: tuple of self context and entity.
        :rtype: tuple
        """

        context = {
            'connector': Event.CONNECTOR,
            'connector_name': Event.CONNECTOR_NAME,
            'component': self.id
        }

        entity = {
            Context.NAME: self.name
        }

        return context, entity

    @property
    def task(self):
        """Get self task or default task if task is not setted.
        """
        result = self.data.get(TaskedVertice.TASK, self.get_default_task())

        return result

    @task.setter
    def task(self, value):
        """Change of task.

        :param value: new task to use.
        """

        self.data[TaskedVertice.TASK] = value

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
            event_type=self.type,
            component=self.id if self.type == Vertice.TYPE else None,
            resource=self.id if self.type == Vertice.TYPE else None,
            id=self.id,
            *args, **kwargs
        )
        return result


@register_task('graph.event_processing')
def event_processing(event, ctx=None, *args, **kwargs):
    """Process input event in getting graph nodes bound to input event entity.

    If at least one graph node is found, execute its tasks.
    """

    if ctx is None:
        ctx = {}

    entity = _context.get_entity(event)

    if entity is not None:
        entity_id = _context.get_entity_id(entity)
        vertices = graph.get_elts(
            data={TaskedVertice.ENTITY: entity_id},
            cls=TaskedVertice
        )

        for vertice in vertices:
            vertice.process(event=event, ctx=ctx, *args, **kwargs)

    return event
