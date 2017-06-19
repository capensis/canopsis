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

from bottle import abort, response
import json

from canopsis.context_graph.manager import ContextGraph
from canopsis.alerts.reader import AlertsReader
from canopsis.common.converters import mongo_filter, id_filter
from canopsis.context_graph.manager import ContextGraph
from canopsis.webcore.utils import gen_json, gen_json_error


context_manager = ContextGraph()
alarm_manager = AlertsReader()


def exports(ws):

    ws.application.router.add_filter('mongo_filter', mongo_filter)
    ws.application.router.add_filter('id_filter', id_filter)

    @ws.application.route(
        '/api/v2/weather/selectors/<selector_filter:mongo_filter>'
    )
    def get_selector(selector_filter):
        """
        Get a list of selectors from a mongo filter.

        :param dict selector_filter: a mongo filter to find selectors
        :rtype: dict
        """

        selector_filter['type'] = 'selector'
        selector_list = context_manager.get_entities(query=selector_filter)

        selectors = []
        for selector in selector_list:
            enriched_entity = {}
            tmp_alarm = alarm_manager.get(
                filter_={'d': selector['_id']}
            )['alarms']

            enriched_entity['entity_id'] = selector['_id']
            enriched_entity['criticity'] = selector['infos'].get(
                'criticity',
                ''
            )
            enriched_entity['org'] = selector['infos'].get('org', '')
            enriched_entity['sla_text'] = ''  # when sla
            enriched_entity['display_name'] = selector['name']
            if tmp_alarm != []:
                enriched_entity['state'] = tmp_alarm[0]['v']['state']
                enriched_entity['status'] = tmp_alarm[0]['v']['status']
                enriched_entity['snooze'] = tmp_alarm[0]['v']['snooze']
                enriched_entity['ack'] = tmp_alarm[0]['v']['ack']
            enriched_entity['pbehavior'] = []  # add this when it's ready
            enriched_entity['linklist'] = []  # add this when it's ready
            selectors.append(enriched_entity)

        return gen_json(response, enriched_entity)

    @ws.application.route("/api/v2/weather/selectors/<selector_id:id_filter>")
    def weatherselectors(selector_id):
        """
        Get selector and contextual informations.

        :param str selector_id: the selector_id to search for
        :return: a list of agglomerated values of entities in the selector
        :rtype: list
        """
        context_manager.logger.critical(selector_id)
        try:
            selector_entity = context_manager.get_entities(
                query={'_id': selector_id})[0]
        except IndexError:
            json_error = {"name" : "resource_not_found",
                      "description": "the selector_id does not match"
                      " any selector"}
            return gen_json_error(response, json_error, 404)

        entities = context_manager.get_entities(
            query=json.loads(selector_entity['infos']['mfilter']))

        entities_list = []
        for entity in entities:
            enriched_entity = {}
            tmp_alarm = alarm_manager.get(
                filter_={'d': entity['_id']}
            )['alarms']

            enriched_entity['entity_id'] = entity['_id']
            enriched_entity['sla_text'] = ''  # TODO when sla, use it
            enriched_entity['org'] = entity['infos'].get('org', '')
            enriched_entity['display_name'] = selector_entity[
                'name']  # check if we need selector here
            enriched_entity['name'] = entity['name']
            if tmp_alarm != []:
                enriched_entity['state'] = tmp_alarm[0]['v']['state']
                enriched_entity['status'] = tmp_alarm[0]['v']['status']
                enriched_entity['snooze'] = tmp_alarm[0]['v']['snooze']
                enriched_entity['ack'] = tmp_alarm[0]['v']['ack']
            enriched_entity['pbehavior'] = []  # TODO wait for pbehavior
            enriched_entity['linklist'] = []  # TODO wait for linklist
            entities_list.append(enriched_entity)

        return gen_json(response, enriched_entity)
