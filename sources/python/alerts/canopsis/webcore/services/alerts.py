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
from canopsis.alerts.reader import AlertsReader


def exports(ws):

    am = Alerts()
    ar = AlertsReader()

    @route(
        ws.application.get,
        name='alerts/alarms2',
        payload=[
            'tstart',
            'tstop',
            'opened',
            'closed',
            'filter',
            'search',
            'sort',
            'limit',
            'offset'
        ]
    )
    def get_alarms2(
            tstart,
            tstop,
            opened=True,
            closed=False,
            consolidations=[],
            filter={},
            search='',
            sort={'opened': 'DESC'},
            limit=50,
            offset=0
    ):
        """
        Return filtered, sorted and paginated alarms.

        :param int tstart: Beginning timestamp of requested period
        :param int tstop: End timestamp of requested period

        :param bool opened: If True, consider alarms that are currently opened
        :param bool closed: If False, consider alarms that have been closed

        :param list consolidations: List of extra columns to compute for each
          returned result.

        :param dict filter: Mongo filter. Keys are UI column names.
        :param str search: Search expression in custom DSL.

        :param dict sort: Dict with only one key. Key is the name of the column
          to sort, and value is either "ASC" or "DESC".

        :param int limit: Maximum number of alarms to return.
        :param int offset: Number of alarms to skip (pagination)

        :returns: List of sorted alarms + pagination informations
        :rtype: dict
        """

        return ar.get_alarms2(
            tstart=tstart,
            tstop=tstop,
            opened=opened,
            closed=closed,
            consolidations=[],
            filter_=filter,
            search=search,
            sort=sort,
            limit=limit,
            offset=offset
        )

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

        return ar.count_alarms_by_period(
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
