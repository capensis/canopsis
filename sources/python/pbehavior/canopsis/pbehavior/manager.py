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
from canopsis.common.init import basestring

from json import loads

from time import time

from datetime import datetime, timedelta

from dateutil import rrulestr

from calendar import timegm

#: pbehavior manager configuration path
CONF_PATH = 'pbehavior/pbehavior.conf'
#: pbehavior manager configuration category name
CATEGORY = 'PBEHAVIOR'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class PBehaviorManager(VEventManager):
    """Dedicated to manage periodic behaviors documents which inherits from
    the vevent documents.
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

    def getending(self, source, behaviors, ts=None):
        """Get end date of corresponding behaviors if a timestamp is in a
        behavior period.

        :param str source: source id.
        :param behaviors: behavior(s) to check at timestamp.
        :type behaviors: list or str
        :param float ts: timestamp to check. If None, use now.
        :return: dict of end timestamp by behavior.
        :rtype: dict
        """

        result = {}
        # get the right ts datetime
        if ts is None:
            ts = time()
        dtts = datetime.fromtimestamp(ts)

        # check if behaviors is unique and ensure it is a set
        isunique = isinstance(behaviors, basestring)
        if isunique:
            behaviors = [behaviors]
        # prepare query
        query = self.get_query(behaviors)
        behaviors = set(behaviors)
        # get documents
        documents = self.values(
            sources=[source],
            dtstart=ts,
            query=query
        )
        # prepare CONSTS
        DURATION = PBehaviorManager.DURATION
        FREQ = PBehaviorManager.FREQ
        DTEND = PBehaviorManager.DTEND
        DTSTART = PBehaviorManager.DTSTART
        BEHAVIORS = PBehaviorManager.BEHAVIORS
        # iterate on documents in order to update result with end ts
        for document in documents:
            # prepare end ts to update in result
            endts = None
            # prepare doc_behaviors such as a conjuguaison with behaviors
            doc_behaviors = set(document[BEHAVIORS]) & behaviors
            # get the right end ts
            if DURATION in document:
                dtstart = document[DTSTART]
                duration = document[DURATION]
                duration = timedelta(seconds=duration)
                if FREQ in document:
                    freq = document[FREQ]
                    dtts = datetime.fromtimestamp(dtstart)
                    rrule = rrulestr(freq, dtts=dtts)
                    before = rrule.before(dtts=ts, inc=True)
                    if before:
                        endbefore = before + duration
                        if endbefore >= dtts:
                            endts = timegm(endbefore.timetuple())
            elif FREQ in document:  # check if ts in freq
                freq = document[FREQ]
                dtts = datetime.fromtimestamp(dtstart)
                rrule = rrulestr(freq, dtts=ts)
                if rrule[0] == dtts:
                    endts = ts
            else:  # get simply dtend
                endts = document[DTEND]

            # update result with upper values
            for behavior in doc_behaviors:
                if behavior not in result or result[behavior] < endts:
                    result[behavior] = endts

        return result
