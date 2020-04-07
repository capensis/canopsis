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

import json
from bottle import request
from pymongo.errors import OperationFailure
from time import time

from canopsis.alerts.filter import AlarmFilter
from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader
from canopsis.alerts.utils import compat_go_crop_states
from canopsis.check import Check
from canopsis.common.converters import id_filter
from canopsis.common.ws import route, WebServiceError
from canopsis.context_graph.manager import ContextGraph
from canopsis.event import forger
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR
from canopsis.pbehavior.manager import PBehaviorManager,PBehavior


def exports(ws):

    ws.application.router.add_filter('id_filter', id_filter)

    context_manager = ContextGraph(ws.logger)
    am = Alerts(*Alerts.provide_default_basics())
    ar = AlertsReader(*AlertsReader.provide_default_basics())
    pbm = PBehaviorManager(*PBehaviorManager.provide_default_basics())

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
            'limit',
            'with_steps',
            'natural_search',
            'active_columns',
            'hide_resources'
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
            limit=None,
            with_steps=False,
            natural_search=False,
            active_columns=None,
            hide_resources=False
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
          returned alarm. Extra columns are "pbehaviors".

        :param dict filter: Mongo filter. Keys are UI column names.
        :param str search: Search expression in custom DSL

        :param str sort_key: Name of the column to sort
        :param str sort_dir: Either "ASC" or "DESC"

        :param int skip: Number of alarms to skip (pagination)
        :param int limit: Maximum number of alarms to return

        :param list active_columns: list of active columns on the brick
        listalarm .

        :param bool hide_resources: hide_resources if component has an alarm

        :returns: List of sorted alarms + pagination informations
        :rtype: dict
        """
        if isinstance(search, int):
            search = str(search)

        try:
            alarms = ar.get(
                tstart=tstart,
                tstop=tstop,
                opened=opened,
                resolved=resolved,
                lookups=lookups,
                filter_=filter,
                search=search.strip(),
                sort_key=sort_key,
                sort_dir=sort_dir,
                skip=skip,
                limit=limit,
                with_steps=with_steps,
                natural_search=natural_search,
                active_columns=active_columns,
                hide_resources=hide_resources
            )
        except OperationFailure as of_err:
            message = 'Operation failure on get-alarms: {}'.format(of_err)
            raise WebServiceError(message)

        alarms_ids = []
        for alarm in alarms['alarms']:
            tmp_id = alarm.get('d')
            if tmp_id:
                alarms_ids.append(tmp_id)
        entities = context_manager.get_entities_by_id(alarms_ids, with_links=True)
        entity_dict = {}
        for entity in entities:
            entity_dict[entity.get('_id')] = entity

        list_alarm = []
        for alarm in alarms['alarms']:
            now = int(time())

            alarm_end = alarm.get('v', {}).get('resolved')
            if not alarm_end:
                alarm_end = now
            alarm["v"]['duration'] = (
                alarm_end - alarm.get('v', {}).get('creation_date', alarm_end))

            state_time = alarm.get('v', {}).get('state', {}).get('t', now)
            alarm["v"]['current_state_duration'] = now - state_time
            tmp_entity_id = alarm['d']

            if alarm['d'] in entity_dict:
                alarm['links'] = entity_dict[alarm['d']]['links']

                # TODO: 'infos' is already present in entity.
                # Remove this one if unused.
                if tmp_entity_id in entity_dict:
                    data = entity_dict[alarm['d']]['infos']
                    if alarm.get('infos'):
                        alarm['infos'].update(data)
                    else:
                        alarm['infos'] = data

            alarm = compat_go_crop_states(alarm)

            list_alarm.append(alarm)

        alarms['alarms'] = list_alarm

        return alarms

    @route(
        ws.application.get,
        name='alerts/get-counters',
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
            'limit',
            'with_steps',
            'natural_search',
            'active_columns',
            'hide_resources'
        ]
    )
    def get_counters(
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
        limit=None,
        with_steps=False,
        natural_search=False,
        active_columns=None,
        hide_resources=False
    ):

        if isinstance(search, int):
            search = str(search)

        try:
            alarms = ar.get(
                tstart=tstart,
                tstop=tstop,
                opened=opened,
                resolved=resolved,
                lookups=lookups,
                filter_=filter,
                search=search.strip(),
                sort_key=sort_key,
                sort_dir=sort_dir,
                skip=skip,
                limit=limit,
                with_steps=with_steps,
                natural_search=natural_search,
                active_columns=active_columns,
                hide_resources=hide_resources,
                add_pbh_filter=False
            )
        except OperationFailure as of_err:
            message = 'Operation failure on get-alarms: {}'.format(of_err)
            raise WebServiceError(message)

        counters = {
            "total": len(alarms['alarms']),
            "total_active": 0,
            "snooze": 0,
            "ack": 0,
            "ticket": 0,
            "pbehavior_active": 0
        }

        alarms_ids = []
        for alarm in alarms['alarms']:
            tmp_id = alarm.get('d')
            if tmp_id:
                alarms_ids.append(tmp_id)
        entities = context_manager.get_entities_by_id(alarms_ids, with_links=True)
        entity_id = []
        for entity in entities:
            _id = entity.get('_id')
            if _id:
                entity_id.append(_id)

        active_pbh = pbm.get_active_pbehaviors_on_entities(entity_id)
        enabled_pbh_entity_dict = set()
        for pbh in active_pbh:
            if pbh[PBehavior.ENABLED]:
                for eid in pbh.get(PBehavior.EIDS, []):
                    if eid in entity_id:
                        enabled_pbh_entity_dict.add(eid)

        for alarm in alarms['alarms']:
            v = alarm.get('v')
            if isinstance(v, dict):
                if v.get('ack', {}).get('_t') == 'ack':
                    counters['ack'] += 1
                if v.get('snooze', {}).get('_t') == 'snooze':
                    counters['snooze'] += 1
                if v.get('ticket', {}).get('_t') in ['declareticket', 'assocticket']:
                    counters['ticket'] += 1
            d = alarm.get('d')
            if d in enabled_pbh_entity_dict:
                counters['pbehavior_active'] += 1

        counters['total_active'] = counters['total'] - counters['pbehavior_active'] - counters['snooze']
        return counters

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
            return gen_json_error({'description': 'nothing to return'},
                                  HTTP_ERROR)

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
            return gen_json_error(
                {'description': 'nothing to insert'}, HTTP_ERROR)

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
            return gen_json_error(
                {'description': 'wrong update dict'}, HTTP_ERROR)

        af = am.alarm_filters.update_filter(filter_id=entity_id, values=dico)
        if not isinstance(af, AlarmFilter):
            return gen_json_error({'description': 'failed to update filter'},
                                  HTTP_ERROR)

        return gen_json(af.serialize())

    @ws.application.delete(
        '/api/v2/alerts/filters/<entity_id:id_filter>'
    )
    def delete_id(entity_id):
        """
        Delete a filter, based on his id.

        :param entity_id: Entity ID of the alarm-filter
        :type entity_id: str

        :rtype: dict
        """
        ws.logger.info('Delete alarm-filter : {}'.format(entity_id))

        return gen_json(am.alarm_filters.delete_filter(entity_id))

    @ws.application.delete(
        '/api/v2/alerts/<mfilter>'
    )
    def delete_filter(mfilter):
        """
        :param str mfilter: mongo filter
        :rtype: dict
        """
        return gen_json(ar.alarm_storage._backend.remove(json.loads(mfilter)))

    @ws.application.post(
        '/api/v2/alerts/done'
    )
    def done_action():
        """
        Trigger done action.

        For json payload, see doc/docs/fr/guide_developpeur/apis/v2/alerts.md

        :rtype: dict
        """
        dico = request.json

        if dico is None or not isinstance(dico, dict) or len(dico) <= 0:
            return gen_json_error(
                {'description': 'wrong done dict'}, HTTP_ERROR)

        author = dico.get(am.AUTHOR)
        event = forger(
            event_type=Check.EVENT_TYPE,
            author=author,
            connector=dico.get('connector'),
            connector_name=dico.get('connector_name'),
            component=dico.get('component'),
            output=dico.get('comment')
        )
        if dico.get('source_type', None) == 'resource':
            event['resource'] = dico['resource']
            event['source_type'] = 'resource'
        ws.logger.debug('Received done action: {}'.format(event))

        entity_id = am.context_manager.get_id(event)
        retour = am.execute_task(
            'alerts.useraction.done',
            event=event,
            author=author,
            entity_id=entity_id
        )
        return gen_json(retour)
