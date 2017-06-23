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

import json

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
        '/api/v2/weather/watchers/<watcher_filter:mongo_filter>'
    )
    def get_watcher(watcher_filter):
        """
        Get a list of watchers from a mongo filter.

        :param dict watcher_filter: a mongo filter to find watchers
        :rtype: dict
        """

        watcher_filter['type'] = 'watcher'
        watcher_list = context_manager.get_entities(query=watcher_filter)

        watchers = []
        for watcher in watcher_list:
            enriched_entity = {}
            tmp_alarm = alarm_manager.get(
                filter_={'d': watcher['_id']}
            )['alarms']

            enriched_entity['entity_id'] = watcher['_id']
            enriched_entity['criticity'] = watcher['infos'].get('criticity', '')
            enriched_entity['org'] = watcher['infos'].get('org', '')
            enriched_entity['sla_text'] = ''  # when sla
            enriched_entity['display_name'] = watcher['name']
            if tmp_alarm != []:
                enriched_entity['state'] = tmp_alarm[0]['v']['state']
                enriched_entity['status'] = tmp_alarm[0]['v']['status']
                enriched_entity['snooze'] = tmp_alarm[0]['v']['snooze']
                enriched_entity['ack'] = tmp_alarm[0]['v']['ack']
            enriched_entity['pbehavior'] = []  # add this when it's ready
            enriched_entity['linklist'] = []  # add this when it's ready

            watchers.append(enriched_entity)

        return gen_json(watchers)

    @ws.application.route("/api/v2/weather/watchers/<watcher_id:id_filter>")
    def weatherwatchers(watcher_id):
        """
        Get watcher and contextual informations.

        :param str watcher_id: the watcher_id to search for
        :return: a list of agglomerated values of entities in the watcher
        :rtype: list
        """
        # Find the selector
        try:
            watcher_entity = context_manager.get_entities(
                query={'_id': watcher_id, 'type': 'selector'})[0]
        except IndexError:
            json_error = {
                "name": "resource_not_found",
                "description": "watcher_id does not match any selector"
            }
            return gen_json_error(json_error, 404)

        # Find entities with the selector filter
        try:
            query = json.loads(watcher_entity['infos']['mfilter'])
        except:
            json_error = {
                "name": "filter_not_found",
                "description": "impossible to load the desired filter"
            }
            ws.logger.error(watcher_entity['infos'])
            return gen_json_error(json_error, 404)

        entities = context_manager.get_entities(query=query)

        entities_list = []
        for entity in entities:
            enriched_entity = {}
            tmp_alarm = alarm_manager.get(
                filter_={'d': entity['_id']}
            )['alarms']

            enriched_entity['entity_id'] = entity['_id']
            enriched_entity['sla_text'] = ''  # TODO when sla, use it
            enriched_entity['org'] = entity['infos'].get('org', '')
            enriched_entity['name'] = entity['name']
            if tmp_alarm != []:
                enriched_entity['state'] = tmp_alarm[0]['v']['state']
                enriched_entity['status'] = tmp_alarm[0]['v']['status']
                enriched_entity['snooze'] = tmp_alarm[0]['v']['snooze']
                enriched_entity['ack'] = tmp_alarm[0]['v']['ack']
            enriched_entity['pbehavior'] = []  # TODO wait for pbehavior
            enriched_entity['linklist'] = []  # TODO wait for linklist

            entities_list.append(enriched_entity)

        return gen_json(entities_list)
