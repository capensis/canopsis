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
from canopsis.common.utils import ensure_iterable

from json import loads

from time import time

from datetime import datetime

#: pbehavior manager configuration path
CONF_PATH = 'pbehavior/pbehavior.conf'
#: pbehavior manager configuration category name
CATEGORY = 'PBEHAVIOR'


def get_query(behaviors):
    """Get a query related to input behaviors.

    :param behaviors: behaviors to find.
    :type behaviors: str or list
    :return: query.
    :rtype: dict
    """
    result = {}

    if behaviors is not None:
        result[PBehaviorManager.BEHAVIORS] = behaviors

    return result


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

    def _get_document_properties(self, document, *args, **kwargs):

        behaviors = document[PBehaviorManager.BEHAVIORS]

        result = {
            PBehaviorManager.BEHAVIORS: behaviors
        }

        return result

    def _get_vevent_properties(self, vevent, *args, **kwargs):

        serialized_behaviors = vevent.get(PBehaviorManager.BEHAVIOR, "[]")
        behaviors = loads(serialized_behaviors)

        result = {
            PBehaviorManager.BEHAVIORS: behaviors
        }

        return result

    def getending(self, source, behaviors, ts=None):
        """Get end date of corresponding behaviors if a timestamp is in a
        behavior period.

        :param str source: source id.
        :param behaviors: behavior(s) to check at timestamp.
        :type behaviors: list or str
        :param long ts: timestamp to check. If None, use now.
        :return: dict of end timestamp by behavior.
        :rtype: dict
        """

        if ts is None:
            ts = time()

        behaviors = ensure_iterable(behaviors)
        vevents = self.values(
            sources=[source],
            dtstart=ts,
            query=self.get_query(behaviors)
        )

        result = {}

        for vevent in vevents:
            vbehaviors = vevent['behaviors']
            dtend = vevent['dtend']

            for behavior in vbehaviors:
                if behavior in behaviors and dtend > result.get(behavior, -1):
                    result[behavior] = dtend

        return result
