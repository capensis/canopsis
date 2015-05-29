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

from time import time

from datetime import datetime

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

    def _get_document_properties(self, document, *args, **kwargs):

        behaviors = document[PBehaviorManager.BEHAVIORS]

        result = {
            PBehaviorManager.BEHAVIORS: behaviors
        }

        return result

    def _get_vevent_properties(self, vevent, *args, **kwargs):

        serialized_behaviors = vevent[PBehaviorManager.BEHAVIOR]
        behaviors = loads(serialized_behaviors)

        result = {
            PBehaviorManager.BEHAVIORS: behaviors
        }

        return result

    def getending(self, source, behaviors, ts=None, dtstart=None, dtend=None):
        """Get end date of corresponding behaviors if a timestamp is in a
        behavior period.

        :param str source: source id.
        :param behaviors: behavior(s) to check at timestamp.
        :type behaviors: list or str
        :param long ts: timestamp to check. If None, use now.
        :param int start: start timestamp.
        :param int end: end timestamp.
        :return: depending on behaviors types:
            - behaviors:
                + str: behavior end timestamp.
                + array: dict of end timestamp by behavior.
        :rtype: dict or long or NoneType
        """

        # initialize ts
        if ts is None:
            ts = time()
        # calculate ts datetime
        dtts = datetime.fromtimestamp(ts)

        query = get_query(behaviors)

        # get entity documents(s)
        documents = self.values(
            sources=source, query=query, dtstart=dtstart, dtend=dtend
        )
        # check if one entity document is asked
        isunique = isinstance(behaviors, basestring)

        # get sbehaviors such as a behaviors set
        sbehaviors = {behaviors} if isunique else set(behaviors)

        # check all pbehavior related to input ts
        result = self._get_ending(
            behaviors=sbehaviors, documents=documents, dtts=dtts
        )

        # keep only entity_id ending dates
        if entity_id in result:
            result = result[entity_id]

        # update result is isunique
        if isunique:
            # convert result to a bool or None if behaviors is str
            result = result[behaviors] if behaviors in result else None

        return result

    def _get_ending(self, behaviors, documents, dtts):
        """Get ending date of occured behavior(s) at a timestamp among value
        periods per entity id.

        :param set behaviors: behavior(s) to check.
        :param list documents: document(s).
        :param datetime dtts: date time moment.
        :return: dict of ending date per entity and behavior.
        """

        result = {}

        for document in documents:
            # get event
            vevent = document[PBehaviorManager.VEVENT]
            event = Event.from_ical(vevent)
            # get behaviors intersection
            dbehaviors = set(document[PBehaviorManager.BEHAVIORS])
            behaviors_to_check = behaviors & dbehaviors

            # if intersection contains elements
            if behaviors_to_check:
                # get duration, dtstart and rrule
                duration = event.get('duration')
                duration = duration.dt
                dtstart = event.get('dtstart')
                dtstart = dtstart.dt

                if isinstance(dtstart, datetime_time):
                    dtstart = datetime.now().replace(
                        hour=dtstart.hour, minute=dtstart.minute,
                        second=dtstart.second, tzinfo=dtstart.tzinfo
                    )

                rrule = event.get('rrule')
                rrule = rrulestr(rrule.to_ical(), cache=True, dtstart=dtstart)
                # calculate first date after dtts including dtts
                before = rrule.before(dt=dtts, inc=True)

                # if before datetime exist
                if before is not None:
                    entity_id = document[PBehaviorManager.ENTITY]

                    if entity_id not in result:
                        result[entity_id] = {}

                    # add duration
                    end = before + duration

                    # and check if dtstart is in [first; end]
                    if before <= dtts <= end:
                        # update end in the result
                        endts = timegm(end.timetuple())

                        for behavior in behaviors_to_check:
                            result[entity_id][behavior] = endts

        return result

    def getending(
        self,
        sources=None, ts=None, dtstart=None, dtend=None, behaviors=None
    ):
        """Get around period dates related to one timestamp and additional
        parameters.

        :param list sources: sources from where parse vevent documents.
        :param int ts: timestamp from when find vevent documents.
        :param int dtstart: vevent dtstart.
        :param int dtend: vevent dtend.
        :param dict query: additional filtering query to apply in the search.
        :param bool after: if True (default), get period after ts (included),
            otherwise, before ts.
        :return: list of around date per source.
        :rtype: dict
        """

        result = {}

        # initialize ts
        if ts is None:
            ts = time()
        # calculate ts datetime
        dtts = datetime.fromtimestamp(ts)

        # get entity documents(s)
        documents = self.values(
            sources=sources, dtstart=dtstart, dtend=dtend, query=query
        )

        for document in documents:
            # get event
            vevent = document[VEventManager.VEVENT]
            event = Event.from_ical(vevent)

            # if rrule/duration are in event
            if (
                VEventManager.DURATION in event
                and VEventManager.RRULE in event
            ):


                # get duration, dtstart and rrule
                duration = event.get(VEventManager.DURATION)
                duration = duration.dt
                dtstart = event.get(VEventManager.DTSTART)
                dtstart = dtstart.dt

                if isinstance(dtstart, datetime_time):
                    dtstart = datetime.now().replace(
                        hour=dtstart.hour, minute=dtstart.minute,
                        second=dtstart.second, tzinfo=dtstart.tzinfo
                    )

                rrule = event.get('rrule')
                rrule = rrulestr(rrule.to_ical(), cache=True, dtstart=dtstart)
                # calculate first date after dtts including dtts
                if after:
                    around = rrule.after(dt=dtts, inc=True)
                else:
                    around = rrule.before(dt=dtts, inc=True)

                # if around datetime exist
                if around is not None:
                    source = document[VEventManager.SOURCE]

                    if source not in result:
                        result[source] = {}

                    # add duration
                    end = around + duration

                    # and check if dtstart is in [first; end]
                    if around <= dtts <= end:
                        # update end in the result
                        endts = timegm(end.timetuple())

                        result[source] = endts

        return result


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
