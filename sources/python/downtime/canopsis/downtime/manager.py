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

from time import time

from date.rrule import rrulestr

from datetime import datetime

CONF_PATH = 'downtime/downtime.conf'
CATEGORY = 'DOWNTIME'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class DowntimeManager(MiddlewareRegistry):
    """Dedicated to manage entity downtimes.

    A downtime is an expression which respects the icalendar specification
    ftp://ftp.rfc-editor.org/in-notes/rfc2445.txt.
    """

    DOWNTIME_STORAGE = 'downtime_storage'  #: downtime storage name

    VALUE = 'value'  #: document value field name

    def __init__(self, downtime_storage=None, *args, **kwargs):

        super(DowntimeManager, self).__init__(*args, **kwargs)

        if downtime_storage is not None:
            self[DowntimeManager.DOWNTIME_STORAGE] = downtime_storage

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

        result = self[DowntimeManager.DOWNTIME_STORAGE].get_elements(
            ids=entity_ids
        )

        return result

    def put(self, entity_id, downtime, cache=False):
        """Put a downtime related to an entity id.

        :param str entity_id: downtime entity id.
        :param str downtime: downtime value respecting icalendar format.
        :param bool cache: if True (False by default), use storage cache.
        :return: entity_id if input downtime has been putted. Otherwise None.
        :rtype: str
        """

        element = {DowntimeManager.VALUE: downtime}

        result = self[DowntimeManager.DOWNTIME_STORAGE].put_element(
            _id=entity_id, element=element, cache=cache
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

    def isdown(self, entity_ids, ts=None):
        """Check if entitie(s) is/are down ts.

        :param entity_ids: entity id(s).
        :type entity_ids: list or str
        :param long ts: timestamp to check. If None, use now.
        :return: depending on entity_ids type:
            - str: True/False if related entity is down at ts. None if entity
                ids is not in Storage.
            - list: set of (entity id: down status at ts). If entity is not
                registered in downtime, no entry is added in the result.
        :rtype: dict or bool or NoneType
        """

        # initialize result
        result = {}
        # initialize ts
        if ts is None:
            ts = time.time()
        # calculate ts datetime
        dtts = datetime.fromtimestamp(ts)
        # get entity downtime(s)
        downtimes = self.get(entity_ids=entity_ids)
        # check if one entity downtime is asked
        isunique = isinstance(entity_ids, basestring)
        if isunique:  # ensure downtimes is a list
            downtimes = [downtimes]
        # check all downtime related to input ts
        for downtime in downtimes:
            # check downtime if not None
            if downtime is not None:
                # get rrule object
                rrule = rrulestr(downtimes, cache=True, dtstart=dtts)
                # calculate first date after dtts including dtts
                after = rrule.after(ts, inc=True)[:1]
                # check if after equals dtts
                isdown = len(after) and after == dtts
                # update isdown flag in the result
                result[downtimes['id']] = isdown

        if isunique:
            # convert result to a bool or None if entity_ids is str
            result = result[entity_ids] if entity_ids in result else None

        return result
