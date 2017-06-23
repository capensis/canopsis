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

from bottle import request

from canopsis.alerts.filter import AlarmFilter
from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader
from canopsis.common.converters import id_filter
from canopsis.common.ws import route
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR


def exports(ws):

    ws.application.router.add_filter('id_filter', id_filter)

    am = Alerts()
    ar = AlertsReader()

    @route(
        ws.application.get,
        name='alerts/get-alarms',
        payload=[
            'tstart',
            'tstop',
            'opened',
            'resolved',
            'lookups',
            'filter',
            'search',
            'sort_key',
            'sort_dir',
            'skip',
            'limit'
        ]
    )
    def get_alarms(
            tstart=None,
            tstop=None,
            opened=True,
            resolved=False,
            lookups=[],
            filter={},
            search='',
            sort_key='opened',
            sort_dir='DESC',
            skip=0,
            limit=50
    ):
        """
        Return filtered, sorted and paginated alarms.

        :param tstart: Beginning timestamp of requested period
        :param tstop: End timestamp of requested period
        :type tstart: int or None
        :type tstop: int or None

        :param bool opened: If True, consider alarms that are currently opened
        :param bool resolved: If True, consider alarms that have been resolved

        :param list lookups: List of extra columns to compute for each
          returned alarm. Extra columns are "pbehaviors" and/or "linklist".

        :param dict filter: Mongo filter. Keys are UI column names.
        :param str search: Search expression in custom DSL

        :param str sort_key: Name of the column to sort
        :param str sort_dir: Either "ASC" or "DESC"

        :param int skip: Number of alarms to skip (pagination)
        :param int limit: Maximum number of alarms to return

        :returns: List of sorted alarms + pagination informations
        :rtype: dict
        """

        return ar.get(
            tstart=tstart,
            tstop=tstop,
            opened=opened,
            resolved=resolved,
            lookups=lookups,
            filter_=filter,
            search=search,
            sort_key=sort_key,
            sort_dir=sort_dir,
            skip=skip,
            limit=limit
        )

    @route(
        ws.application.get,
        name='alerts/search/validate',
        payload=['expression']
    )
    def validate_search(expression):
        """
        Tell if a search expression is valid from a grammatical propespective.

        :param str expression: Search expression

        :returns: True if valid, False otherwise
        :rtype: bool
        """

        try:
            ar.interpret_search(expression)

        except Exception:
            return False

        else:
            return True

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

    @ws.application.get(
        '/api/v2/alerts/filters/<entity_id:id_filter>'
    )
    def get_filter(entity_id):
        """
        Get all filters linked with an alarm.

        :param str entity_id: Entity ID of the alarm-filter

        :returns: a list of <AlarmFilter>
        """
        filters = am.alarm_filters.get_filter(entity_id)
        if filters is None:
            return gen_json_error({'description': 'nothing to return'}, HTTP_ERROR)

        return gen_json([l.serialize() for l in filters])

    @ws.application.post(
        '/api/v2/alerts/filters'
    )
    def create_filter():
        """
        Create a new alarm filter.

        :returns: an <AlarmFilter>
        """
        # element is a full AlarmFilter (dict) to insert
        element = request.json

        if element is None:
            return gen_json_error({'description': 'nothing to insert'}, HTTP_ERROR)

        new = am.alarm_filters.create_filter(element=element)
        new.save()

        return gen_json(new.serialize())

    @ws.application.put(
        '/api/v2/alerts/filters/<entity_id:id_filter>'
    )
    def update_filter(entity_id):
        """
        Update an existing alam filter.

        :param entity_id: Entity ID of the alarm-filter
        :type entity_id: str
        :returns: <AlarmFilter>
        :rtype: dict
        """
        dico = request.json

        if dico is None or not isinstance(dico, dict) or len(dico) <= 0:
            return gen_json_error({'description': 'wrong update dict'}, HTTP_ERROR)

        af = am.alarm_filters.update_filter(filter_id=entity_id, values=dico)
        if not isinstance(af, AlarmFilter):
            return gen_json_error({'description': 'failed to update filter'},
                                  HTTP_ERROR)

        return gen_json(af.serialize())

    @ws.application.delete(
        '/api/v2/alerts/filters/<entity_id:id_filter>'
    )
    def delete_filter(entity_id):
        """
        Delete a filter, based on his id.

        :param entity_id: Entity ID of the alarm-filter
        :type entity_id: str

        :rtype: dict
        """
        ws.logger.info('Delete alarm-filter : {}'.format(entity_id))

        return gen_json(am.alarm_filters.delete_filter(entity_id))
