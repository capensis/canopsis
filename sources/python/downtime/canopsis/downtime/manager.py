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

from date.rrule import rrulestr

from datetime import datetime

CONF_PATH = 'downtime/downtime.conf'
CATEGORY = 'DOWNTIME'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class DowntimeManager(MiddlewareRegistry):
    """Dedicated to manage downtime.
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
        """

        element = {DowntimeManager.VALUE: downtime}

        self[DowntimeManager.DOWNTIME_STORAGE].put_element(
            _id=entity_id, element=element, cache=cache
        )

    def remove(self, entity_ids, cache=False):
        """Remove entity downtime(s) related to input entity id(s).

        :param entity_ids: downtime entity id(s).
        :type entity_ids: list or str
        :param bool cache: if True (False by default), use storage cache.
        """

        self.del_elements(ids=entity_ids, cache=cache)

    def isdown(self, entity_ids, ts):
        """Check if entitie(s) is/are down ts.

        :param entity_ids: entity id(s).
        :type entity_ids: list or str
        :param long ts: timestamp to check.
        :return: depending on entity_ids type:
            - str: True/False if related entity is down at ts. None if entity
                ids is not in Storage.
            - list: set of (entity id: down status at ts). If entity is not
                registered in downtime, no entry is added in the result.
        :rtype: dict or bool or NoneType
        """

        downtimes = self.get(entity_ids=entity_ids)

        isunique = isinstance(entity_ids, basestring)

        result = {}

        if isunique:
            downtimes = [downtimes]

        for downtime in downtimes:
            # check downtime if not None
            if downtime is not None:
                dtts = datetime.fromtimestamp(ts)
                rrule = rrulestr(downtimes, cache=True, dtstart=dtts)
                isdown = rrule.after(ts, inc=True)
                result[downtimes['id']] = isdown

        if isunique:
            # convert result to a bool or None if entity_ids is str
            result = result[entity_ids] if entity_ids in result else None

        return result
