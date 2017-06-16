# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from canopsis.context_graph.manager import ContextGraph
import json

from canopsis.alerts.reader import AlertsReader
from canopsis.common.converters import mongo_filter, id_filter
from canopsis.context_graph.manager import ContextGraph

context_manager = ContextGraph()
alarm_manager = AlertsReader()

from bottle import abort, response

def exports(ws):

    ws.application.router.add_filter('mongo_filter', mongo_filter)
    ws.application.router.add_filter('id_filter', id_filter)

    @ws.application.route('/weather/selectors/<selector_filter:mongo_filter>')
    def get_selector(selector_filter):
        """
        Get a list of selectors from a mongo filter.

        :param dict selector_filter: a mongo filter to find selectors
        :rtype: dict
        """
        selector_filter['type'] = 'selector'
        selector_list = context_manager.get_entities(query=selector_filter)

        ret_val = []
        for i in selector_list:
            tmp = {}
            tmp_alarm = alarm_manager.get(filter_={'d': i['_id']})['alarms']

            tmp['entity_id'] = i['_id']
            tmp['criticity'] = i['infos'].get('criticity', '')
            tmp['org'] = i['infos'].get('org', '')
            tmp['sla_text'] = ''  # when sla
            tmp['display_name'] = i['name']
            if tmp_alarm != []:
                tmp['state'] = tmp_alarm[0]['v']['state']
                tmp['status'] = tmp_alarm[0]['v']['status']
                tmp['snooze'] = tmp_alarm[0]['v']['snooze']
                tmp['ack'] = tmp_alarm[0]['v']['ack']
            tmp['pbehavior'] = []  # add this when it's ready
            tmp['linklist'] = []  # add this when it's ready
            ret_val.append(tmp)

        return ret_val

    @ws.application.route("/weather/selectors/<selector_id:id_filter>")
    def weatherselectors(selector_id):
        """
        Get selector and contextual informations.

        :param str selector_id: the selector_id to search for
        :return: a list of agglomerated values of entities in the selector
        :rtype: list
        """
        try:
            selector_entity = context_manager.get_entities(
                query={'_id': selector_id})[0]
        except IndexError:
            abort(404, "No selfuch selector")

        entities = context_manager.get_entities(
            query=json.loads(selector_entity['infos']['mfilter']))

        ret_val = []
        for i in entities:
            tmp_val = {}
            tmp_alarm = alarm_manager.get(filter_={'d': i['_id']})['alarms']

            tmp_val['entity_id'] = i['_id']
            tmp_val['sla_text'] = ''  # TODO when sla, use it
            tmp_val['org'] = i['infos'].get('org', '')
            tmp_val['display_name'] = selector_entity[
                'name']  # check if we need selector here
            tmp_val['name'] = i['name']
            if tmp_alarm != []:
                tmp_val['state'] = tmp_alarm[0]['v']['state']
                tmp_val['status'] = tmp_alarm[0]['v']['status']
                tmp_val['snooze'] = tmp_alarm[0]['v']['snooze']
                tmp_val['ack'] = tmp_alarm[0]['v']['ack']
            tmp_val['pbehavior'] = []  # TODO wait for pbehavior
            tmp_val['linklist'] = []  # TODO wait for linklist
            ret_val.append(tmp_val)

        response.content_type = "application/json"
        return json.dumps(ret_val)
