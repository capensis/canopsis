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


#: pbehavior manager configuration path
CONF_PATH = 'pbehavior/pbehavior.conf'
#: pbehavior manager configuration category name
CATEGORY = 'PBEHAVIOR'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class PBehaviorManager(VEventManager):
    """Dedicated to manage periodic behavior.

    Such period are technically an expression which respects the icalendar
    specification ftp://ftp.rfc-editor.org/in-notes/rfc2445.txt.

    A pbehavior document contains several values. Each value contains
    an icalendar expression (dtstart, rrule, duration) and an array of
    behavior entries:

    {
        id: document_id,
        entity_id: entity id,
        period: period,
        behaviors: behavior ids
    }.
    """

    BEHAVIOR = 'X-Canopsis-BehaviorType'  #: behavior type key in period

    BEHAVIORS = 'behaviors'  #: behaviors value field name

    def _get_info(self, vevent, *args, **kwargs):

        serialized_behaviors = vevent[PBehaviorManager.BEHAVIOR]

        result = loads(serialized_behaviors)

        return result
