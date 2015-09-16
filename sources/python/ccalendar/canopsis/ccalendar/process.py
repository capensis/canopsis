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

"""Module in charge of defining topology processing in engines.

When an event occured, the related entity is retrieved and all bound
topological nodes are retrieved as well in order to execute theirs rules.

First, a topology processing is triggered when an event occured.

From this event, bound topology nodes are retrieved in order to apply node
    rules.

A typical topological task condition is an ``canopsis.task.condition.all``
composed of the ``canopsis.topology.rule.condition.new_state`` condition and
``canopsis.topology.rule.action.change_state`` action.
If this condition is checked, then other specific conditions can be applied
such as those defined in the canopsis.topology.rule.action module.
"""

from canopsis.topology.elements import Topology, TopoNode
from canopsis.topology.manager import TopologyManager
from canopsis.context.manager import Context
from canopsis.task.core import register_task
from canopsis.event import Event
from canopsis.check.manager import CheckManager

context = Context()
tm = TopologyManager()
_check = CheckManager()

SOURCE = 'source'
PUBLISHER = 'publisher'


@register_task
def event_processing(event, logger=None, **kwargs):
    """Process input event in getting topology nodes bound to input event
    entity.

    One topology nodes are founded, executing related rules.

    :param dict event: event to process.
    :param Engine engine: engine which consumes the event.
    :param TopologyManager manager: topology manager to use.
    :param Logger logger: logger to use in this task.
    """

    return event
