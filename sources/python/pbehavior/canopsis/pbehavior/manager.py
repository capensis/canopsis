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

from canopsis.common.init import basestring
from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths
)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.storage import Storage

from time import time

from icalendar import Event
from date.rrule import rrulestr
from calendar import timegm

from datetime import datetime

CONF_PATH = 'pbehavior/pbehavior.conf'
CATEGORY = 'PBEHAVIOR'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class PBehaviorManager(MiddlewareRegistry):
    """Dedicated to manage periodic behavior.

    Such period are technically an expression which respects the icalendar
    specification ftp://ftp.rfc-editor.org/in-notes/rfc2445.txt.

    A pbehavior document contains several values. Each value contains
    an icalendar expression (dtstart, rrule, duration) and an array of
    behavior entries:

    {
        id: document_id,
        values: [{period: ..., behaviors: [...]}*]
    }.
    """

    PBEHAVIOR_STORAGE = 'pbehavior_storage'  #: downtime storage name

    ID = Storage.DATA_ID  #: document id
    VALUES = 'values'  #: document values field name
    PERIOD = 'period'  #: period value field name
    BEHAVIORS = 'behaviors'  #: behaviors value field name

    def __init__(self, PBEHAVIOR_storage=None, *args, **kwargs):

        super(PBehaviorManager, self).__init__(*args, **kwargs)

        if PBEHAVIOR_storage is not None:
            self[PBehaviorManager.PBEHAVIOR_STORAGE] = PBEHAVIOR_storage

    def get(self, entity_ids):
        """Get a downtime related to input entity id(s).

        :param entity_ids: entity id(s) bound to downtime.
        :type entity_ids: list or str
        :return: depending on entity_ids type:
            - str: one dictionary with two fields, entity_id which contains
                the entity id, and ``value`` which contains the downtime
                expression.
            - list: list of previous dictionaries.
        :rtype: list or dict
        """

        result = self[PBehaviorManager.PBEHAVIOR_STORAGE].get_elements(
            ids=entity_ids
        )

        return result

    def put(self, entity_id, document, cache=False):
        """Put a downtime related to an entity id.

        :param str entity_id: downtime entity id.
        :param dict document: pbehavior document.
        :param bool cache: if True (False by default), use storage cache.
        :return: entity_id if input downtime has been putted. Otherwise None.
        :rtype: str
        """

        result = self[PBehaviorManager.PBEHAVIOR_STORAGE].put_element(
            _id=entity_id, element=document, cache=cache
        )

        return result

    def remove(self, entity_ids, cache=False):
        """Remove entity downtime(s) related to input entity id(s).

        :param entity_ids: downtime entity id(s).
        :type entity_ids: list or str
        :param bool cache: if True (False by default), use storage cache.
        :return: entity id(s) of removed downtime(s).
        :rtype: list
        """

        result = self.del_elements(ids=entity_ids, cache=cache)

        return result

    def getending(self, entity_ids, behaviors, ts=None):
        """Get end date of corresponding behaviors if a timestamp is in a
        behavior period.

        :param entity_ids: entity id(s).
        :type entity_ids: list or str
        :param behaviors: behavior(s) to check at timestamp.
        :type behaviors: list or str
        :param long ts: timestamp to check. If None, use now.
        :return: depending on entity_ids and behaviors types:
            - behaviors:
                + str: behavior end timestamp.
                + array: dict of end timestamp by behavior.
            - entity_ids:
                + str: one end dates of behaviors if related entity is in
                    behavior at ts. None if entity ids is not in Storage.
                + list: set of (entity id: behaviors). If entity is not
                    registered in behavior state, no entry is added in the
                    result.
        :rtype: dict or long or NoneType
        """

        # initialize result
        result = {}

        # initialize ts
        if ts is None:
            ts = time.time()
        # calculate ts datetime
        dtts = datetime.fromtimestamp(ts)

        # get entity documents(s)
        documents = self.get(entity_ids=entity_ids)
        # check if one entity document is asked
        isunique = isinstance(entity_ids, basestring)

        if isunique:  # ensure documents is a list
            documents = [] if documents is None else [documents]

        # are many behaviors asked ?
        isbunique = isinstance(behaviors, basestring)

        # get sbehaviors such as a behaviors set
        if isbunique:
            sbehaviors = {behaviors}

        else:
            sbehaviors = set(behaviors)

        # check all downtime related to input ts
        for document in documents:
            values = document.get(PBehaviorManager.VALUES)
            behavior_result = self._get_ending(
                behaviors=sbehaviors, values=values, dtts=dtts
            )

            # update result only if behavior result
            if behavior_result:
                document_id = document[PBehaviorManager.ID]

                if isbunique:
                    result[document_id] = behavior_result[behaviors]

                else:
                    result[document_id] = behavior_result

        # update result is isunique
        if isunique:
            # convert result to a bool or None if entity_ids is str
            result = result[entity_ids] if entity_ids in result else None

        return result

    def whois(self, behaviors, ts=None):
        """
        Get entities which currently have specific behaviors.

        :param behaviors: behavior(s) to look for.
        :type behaviors: list or str

        :return: list of entities ids with the specified behaviors
        :rtype: list
        """

        result = []

        if ts is None:
            ts = time()

        # calculate ts datetime
        dtts = datetime.fromtimestamp(ts)

        documents = self[PBehaviorManager.PBEHAVIOR_STORAGE].find_elements()

        if isinstance(behaviors, basestring):
            behaviors = [behaviors]

        len_behaviors = len(behaviors)

        for document in documents:
            values = document.get(PBehaviorManager.VALUES)
            ending = self._get_ending(
                behaviors=behaviors, values=values, dtts=dtts
            )
            if len(ending) != len_behaviors:
                document_id = document[PBehaviorManager.ID]
                result.append(document_id)

        return result

    def _get_ending(self, behaviors, values, dtts):
        """Get ending date of occured behavior(s) at a timestamp among value
        periods.

        :param set behaviors: behavior(s) to check.
        :param value: (list of) period(s) and behaviors.
        :type value: dict or list
        :param datetime dtts: date time moment.
        """

        result = {}

        if isinstance(values, dict):
            values = [values]

        for value in values:
            # get period
            period = value[PBehaviorManager.PERIOD]
            period = Event.from_ical(period)
            # get behaviors intersection
            dbehaviors = set(value[PBehaviorManager.BEHAVIORS])
            behaviors_to_check = behaviors & dbehaviors
            # if intersection contains elements
            for behavior in behaviors_to_check:
                # get duration, dtstart and rrule
                duration = period.get('duration')
                duration = duration.dt
                dtstart = period.get('dtstart')
                dtstart = dtstart.dt
                rrule = period.get('rrule')
                rrule = rrulestr(rrule.to_ical(), cache=True, dtstart=dtstart)
                # calculate first date after dtts including dtts
                before = rrule.before(dt=dtts, inc=True, count=1)
                # if before datetimes exist
                if before:
                    # first date is the first before
                    first = before[1]
                    # add duration
                    end = first + duration
                    # and check if dtstart is in [first; end]
                    if first <= dtts <= end:
                        # update end in the result
                        endts = timegm(end.timetuple())
                        result[behavior] = endts

        return result
