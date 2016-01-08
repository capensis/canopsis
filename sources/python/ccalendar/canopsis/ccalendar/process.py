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

"""Process input calendarevent

A calendarevent could be created by the user using the widgetCalendar in the
UI or by adding a ical Calendar.

This kind of event is processing by this process.
"""

from canopsis.ccalendar.manager import CalendarManager
from canopsis.topology.elements import Topology, TopoNode
from canopsis.topology.manager import TopologyManager
from canopsis.context.manager import Context
from canopsis.task.core import register_task
from canopsis.event import Event
from canopsis.check.manager import CheckManager

cm = CalendarManager()
_check = CheckManager()

SOURCE = 'source'
PUBLISHER = 'publisher'


@register_task
def calendar_event_processing(event, **kwargs):
    """Process input event in getting calendarevent to input event
    entity.

    :param dict event: event to process.
    """

    return event
