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

from canopsis.alerts.reader import AlertsReader
from canopsis.context_graph.manager import ContextGraph
from canopsis.webcore.utils import gen_json_error, HTTP_ERROR


def exports(ws):
    context_manager = ContextGraph(ws.logger)
    ar = AlertsReader(*AlertsReader.provide_default_basics())

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
