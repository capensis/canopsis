# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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
from __future__ import unicode_literals

from canopsis.common.ws import route
from canopsis.alerts.manager import Alerts


def exports(ws):

    am = Alerts()

    @route(
            ws.application.get,
            name='alerts/alarms',
            payload=['resolved', 'tags', 'exclude_tags'],
    )
    def get_alarms(
            resolved=False,
            tags=None,
            exclude_tags=None,
            snoozed=False,
            limit=5,
            start=0,
            filter={},
            sort=[{"direction": "ASC"}],
    ):
        """
        Get alarms

        :param bool resolved: If ``True``, returns only resolved alarms, else
          returns only unresolved alarms (default: ``False``).

        :param tags: Tags which must be set on alarm (optional)
        :type tags: str or list

        :param exclude_tags: Tags which must not be set on alarm (optional)
        :type tags: str or list

        :param int limit: Number of entries returned (TODO)

        :param int start: Pagination index (TODO)

        :param dict filter: TODO

        :param list sort: TODO

        :returns: Iterable of alarms matching
        """

        alarms = am.get_alarms(
            resolved=resolved,
            tags=tags,
            exclude_tags=exclude_tags,
            snoozed=snoozed,
        )

        return alarms

    @route(
            ws.application.get,
            name='alerts/count',
            payload=['start', 'stop', 'limit', 'select'],
    )
    def count_by_period(
            start,
            stop,
            limit=100,
            select=None,
    ):
        """
        Count alarms that have been opened during (stop - start) period.

        :param start: Beginning timestamp of period
        :type start: int

        :param stop: End timestamp of period
        :type stop: int

        :param limit: Counts cannot exceed this value
        :type limit: int

        :param query: Custom mongodb filter for alarms
        :type query: dict

        :return: List in which each item contains a time interval and the
                 related count
        :rtype: list
        """

        return am.count_alarms_by_period(
            start,
            stop,
            limit=limit,
            query=select,
        )

    @route(
            ws.application.get,
            name='alerts/get-current-alarm',
            payload=['entity_id'],
    )
    def get_current_alarm(entity_id):
        """
        Get current unresolved alarm for a entity.

        :param str entity_id: Entity ID of the alarm

        :returns: Alarm as dict if something is opened, else None
        """

        return am.get_current_alarm(entity_id)
