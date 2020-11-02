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
from canopsis.metaalarmrule.manager import MetaAlarmRuleManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR
from canopsis.pbehavior.manager import PBehaviorManager,PBehavior


def exports(ws):

    ws.application.router.add_filter('id_filter', id_filter)

    context_manager = ContextGraph(ws.logger)
    am = Alerts(*Alerts.provide_default_basics())
    ar = AlertsReader(*AlertsReader.provide_default_basics())
    ma_rule_manager = MetaAlarmRuleManager(
        *MetaAlarmRuleManager.provide_default_basics())
    pbm = PBehaviorManager(*PBehaviorManager.provide_default_basics())

    @route(
        ws.application.get,
        name='alerts/get-alarms',
        payload=[
            'authkey',
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
            'hide_resources',
            'with_consequences',
            'with_causes',
            'correlation',
            'manual_only'
        ]
    )
    def get_alarms(
            authkey=None,
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
            hide_resources=False,
            with_consequences=False,
            with_causes=False,
            correlation=False,
            manual_only=False
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

        :param bool manual_only: manual meta-alarms only with open status (resolved=None)
        :returns: List of sorted alarms + pagination informations
        :rtype: dict
        """
        if isinstance(search, int):
            search = str(search)

        if manual_only:
            return ar.manual_only()

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
                with_consequences=with_consequences,
                correlation=correlation
            )
        except OperationFailure as of_err:
            message = 'Operation failure on get-alarms: {}'.format(of_err)
            raise WebServiceError(message)

        alarms_ids, consequences_children = [], []
        alarm_children = {'alarms': [], 'total': 0}
        for alarm in alarms['alarms']:
            if with_consequences:
                consequences_children.extend(alarm.get('consequences', {}).get('data', []))
            elif with_causes and alarm.get('v') and alarm['v'].get('parents'):
                consequences_children.extend(alarm['v']['parents'])
            tmp_id = alarm.get('d')
            if tmp_id:
                alarms_ids.append(tmp_id)
        entities = context_manager.get_entities_by_id(alarms_ids, with_links=False)

        entity_dict = {}
        for entity in entities:
            entity_dict[entity.get('_id')] = entity

        if consequences_children:
            alarm_children = ar.get(
                tstart=tstart,
                tstop=tstop,
                opened=True,
                resolved=True,
                lookups=lookups,
                filter_={'d': {'$in': consequences_children}},
                sort_key=sort_key,
                sort_dir=sort_dir,
                skip=skip,
                limit=None,
                natural_search=natural_search,
                active_columns=active_columns,
                hide_resources=hide_resources,
                correlation=correlation,
                consequences_children=True
            )

        list_alarm = []
        if ('rules' in alarms) or not correlation:
            if not correlation:
                parent_eids = set()
                alarms['rules'] = dict()
                for alarm in alarms['alarms']:
                    if 'd' in alarm and 'v' in alarm and alarm['v'].get('parents'):
                        for v in alarm['v']['parents']:
                            parent_eids.add(v)
                named_rules = ar.meta_parents_with_rules(list(parent_eids))
                for alarm in alarms['alarms']:
                    if 'd' in alarm and 'v' in alarm and alarm['v'].get('parents'):
                        alarm_named_rules = dict()
                        for p in alarm['v']['parents']:
                            for r in named_rules[p]:
                                alarm_named_rules[r['id']] = r
                        alarms['rules'][alarm['d']] = alarm_named_rules.values()
            else:
                rule_ids = set()
                for alarm_rules in alarms['rules'].values():
                    for v in alarm_rules:
                        rule_ids.add(v)

                named_rules = ma_rule_manager.read_rules_with_names(list(rule_ids))
                for d, alarm_rules in alarms['rules'].items():
                    alarm_named_rules = []
                    for v in alarm_rules:
                        alarm_named_rules.append({'id': v, 'name': named_rules.get(v, "")})
                    alarms['rules'][d] = alarm_named_rules
        else:
            alarms['rules'] = dict()

        children_ent_ids = set()
        for alarm in alarms['alarms']:
            rules = alarms['rules'].get(alarm['d'], []) if 'd' in alarm and 'v' in alarm and \
                alarm['v'].get('parents') else None
            if rules:
                if with_causes:
                    alarm['causes'] = {
                        'total': len(alarm_children['alarms']),
                        'data': alarm_children['alarms'],
                    }
                    for al_child in alarm_children['alarms']:
                        children_ent_ids.add(al_child['d'])
                else:
                    alarm['causes'] = {
                        'total': len(rules),
                        'rules': rules,
                    }

            if alarm.get('v') is None:
                alarm['v'] = dict()
            if alarm.get('v').get('meta'):
                del alarm['v']['meta']

            if isinstance(alarm.get('rule'), basestring) and alarm['rule'] != "":
                alarm['rule'] = {'id': alarm['rule'], 'name': named_rules.get(alarm['rule'], alarm['rule'])}

            now = int(time())

            alarm_end = alarm.get('v', {}).get('resolved')
            if not alarm_end:
                alarm_end = now
            alarm_value = alarm.get('v', {})
            activation_date = alarm_value.get('activation_date')
            alarm["v"]['duration'] = alarm_end - (alarm_value.get('creation_date', alarm_end) \
                if activation_date is None else activation_date)

            state_time = alarm.get('v', {}).get('state', {}).get('t', now)
            alarm["v"]['current_state_duration'] = now - state_time
            tmp_entity_id = alarm['d']

            if alarm['d'] in entity_dict:
                alarm['links'] = context_manager.enrich_links_to_entity_with_alarm(entity_dict[alarm['d']], alarm)

                # TODO: 'infos' is already present in entity.
                # Remove this one if unused.
                if tmp_entity_id in entity_dict:
                    data = entity_dict[alarm['d']]['infos']
                    if alarm.get('infos'):
                        alarm['infos'].update(data)
                    else:
                        alarm['infos'] = data

            alarm = compat_go_crop_states(alarm)

            if with_consequences and isinstance(alarm.get('consequences'), dict) and alarm_children['total'] > 0:
                map(lambda al_ch: al_ch.update({'causes': {'rules': [alarm['rule']], 'total': 1}}),  alarm_children['alarms'])
                alarm['consequences']['data'] = alarm_children['alarms']
                alarm['consequences']['total'] = alarm_children['total']
                for al_child in alarm_children['alarms']:
                    children_ent_ids.add(al_child['d'])

            list_alarm.append(alarm)

        if children_ent_ids:
            children_entities = context_manager.get_entities_by_id(
                list(children_ent_ids), with_links=False)
            for entity in children_entities:
                entity_dict[entity.get('_id')] = entity

            for alarm in alarms['alarms']:
                for cat in ('causes', 'consequences'):
                    if cat in alarm and alarm[cat].get('data'):
                        for child in alarm[cat]['data']:
                            if child['d'] in entity_dict:
                                child['links'] = context_manager.enrich_links_to_entity_with_alarm(
                                    entity_dict[child['d']], child)

        del alarms['rules']
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
            'hide_resources',
            'correlation'
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
        hide_resources=False,
        correlation=False
    ):

        if isinstance(search, int):
            search = str(search)

        try:
            alrs = ar.get(
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
                add_pbh_filter=False,
                correlation=False
            )
        except OperationFailure as of_err:
            message = 'Operation failure on get-alarms: {}'.format(of_err)
            raise WebServiceError(message)

        def count(alarms):
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

            pbehavior_active_snooze = 0

            for alarm in alarms['alarms']:
                v = alarm.get('v')
                snoozed = False
                if isinstance(v, dict):
                    if v.get('ack', {}).get('_t') == 'ack':
                        counters['ack'] += 1
                    snoozed = v.get('snooze', {}).get('_t') == 'snooze'
                    if snoozed:
                        counters['snooze'] += 1
                    if v.get('ticket', {}).get('_t') in ['declareticket', 'assocticket']:
                        counters['ticket'] += 1
                d = alarm.get('d')
                if d in enabled_pbh_entity_dict:
                    counters['pbehavior_active'] += 1
                    if snoozed:
                        pbehavior_active_snooze += 1

            counters['total_active'] = counters['total'] - counters['pbehavior_active'] - counters['snooze'] + \
                pbehavior_active_snooze
            return counters

        if not correlation:
            return count(alrs)

        try:
            alrs_with_meta_group = ar.get(
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
                add_pbh_filter=False,
                correlation=True
            )

            regular_counters = count(alrs)
            meta_alarm_counters = count(alrs_with_meta_group)
            counters = dict(regular_counters.items(
            ) + {k + "_correlation": v for k, v in meta_alarm_counters.items()}.items())
            return counters
        except OperationFailure as of_err:
            message = 'Operation failure on get-alarms: {}'.format(of_err)
            raise WebServiceError(message)

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

    @ws.application.post(
        '/api/v2/links',
    )
    def links():
        r = request.json
        if not r or not isinstance(r, dict) or not isinstance(r.get('entities'), list):
            return gen_json_error(
                {'description': 'wrong entities payload'}, HTTP_ERROR)

        links, alarm_ids, entity_ids = [], [], []

        for en in r['entities']:
            if isinstance(en,  dict):
                if en.get('alarm'):
                    alarm_ids.append(en['alarm'])
                if en.get('entity'):
                    entity_ids.append(en['entity'])

        if alarm_ids:
            alarms = ar.alarm_collection.find({'_id': {'$in': alarm_ids}})
            entities = context_manager.get_entities_by_id(
                entity_ids, with_links=False)

            entity_dict = {}
            for entity in entities:
                entity_dict[entity.get('_id')] = entity
            for alarm in alarms:
                if alarm['d'] in entity_dict:
                    links.append({
                        'entity': alarm['d'],
                        'alarm': alarm['_id'],
                        'links': context_manager.enrich_links_to_entity_with_alarm(entity_dict[alarm['d']], alarm)
                    })
        elif entity_ids:
            entities = context_manager.get_entities(
                query={"_id": {"$in": entity_ids}},
                with_links=True
            )
            for entity in entities:
                if isinstance(entity, dict) and '_id' in entity and entity.get('links'):
                    links.append({
                        'entity': entity['_id'],
                        'links': entity['links']
                    })
        return {
            'data': links
        }
