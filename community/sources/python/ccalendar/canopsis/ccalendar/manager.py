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

from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths
)
from canopsis.vevent.manager import VEventManager

from json import loads


#: calendar manager configuration path
CONF_PATH = 'calendar/calendar.conf'
#: calendar manager configuration category name
CATEGORY = 'CALENDAR'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class CalendarManager(VEventManager):
    """Dedicated to manage calendar event.

    Such period are technically an expression which respects the icalendar
    specification ftp://ftp.rfc-editor.org/in-notes/rfc2445.txt.

    A calendar document contains several values. Each value contains
    an icalendar expression (dtstart, rrule, duration)
    """

    EVENTCATEGORIES = 'EVENTCATEGORIES'
    OUTPUT = 'OUTPUT'

    def _get_vevent_properties(self, vevent, *args, **kwargs):
        """Get information from a vevent.

        :param Event vevent: vevent from where get information
        :return: vevent information in a dictionary
        :rtype: dict
        """

        serialized_eventcategories = vevent[CalendarManager.EVENTCATEGORIES]
        serialized_output = vevent[CalendarManager.OUTPUT]

        result = {
            "eventcategories": serialized_eventcategories,
            "output": serialized_output
        }

        return result

    def _get_document_properties(self, document):
        """Get properties from a document.

        :param dict document: document from where get properties.
        :return: document properties in a dictionary.
        :rtype: dict
        """
        serialized_eventcategories = document.get("eventcategories")
        serialized_output = document.get("output")

        result = {
            "eventcategories": serialized_eventcategories,
            "output": serialized_output
        }

        return result
