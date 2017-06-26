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
from canopsis.common.utils import get_rrule_freq
from canopsis.context_graph.manager import ContextGraph
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_NOT_FOUND

LOGGER = None

context_manager = ContextGraph()
alarm_manager = AlertsReader()
pbehavior_manager = PBehaviorManager()


def __format_pbehavior(pbehavior):
    """Rewrite en pbehavior from db format to front format.

    :param dict pbehavior: a pbehavior dict
    """
    EVERY = "Every {0}"
    to_delete = [
        "_id", "connector", "author", "comments", "filter", "connector_name",
        "eids"
    ]

    pbehavior["behavior"] = pbehavior.pop("name")
    pbehavior["dtstart"] = pbehavior.pop("tstart")
    pbehavior["dtend"] = pbehavior.pop("tstop")
    pbehavior["isActive"] = pbehavior.pop("enabled")

    # parse the rrule to get is "text"
    rrule = {}
    rrule["rrule"] = pbehavior["rrule"]

    freq = get_rrule_freq(pbehavior["rrule"])

    if freq == "SECONDLY":
        rrule["text"] = EVERY.format("second")
    elif freq == "MINUTELY":
        rrule["text"] = EVERY.format("minute")
    elif freq == "HOURLY":
        rrule["text"] = EVERY.format("hour")
    elif freq == "DAILY":
        rrule["text"] = EVERY.format("day")
    elif freq == "WEEKLY":
        rrule["text"] = EVERY.format("week")
    elif freq == "MONTHLY":
        rrule["text"] = EVERY.format("month")
    elif freq == "YEARLY":
        rrule["text"] = EVERY.format("year")

    pbehavior["rrule"] = rrule

    for key in to_delete:
        try:
            pbehavior.pop(key)
        except KeyError:
            pass


def add_pbehavior_info(enriched_entity):
    """Add pbehavior related field to selectors. This function will add
    the related pbehavior in 'pbehavior'.

    :param dict enriched_entity: the entity to enrich
    """

    enriched_entity["pbehavior"] = pbehavior_manager.get_pbehaviors_by_eid(
        enriched_entity['entity_id'])

    LOGGER.debug("Pbehavior list : {0}".format(enriched_entity["pbehavior"]))

    for pbehavior in enriched_entity["pbehavior"]:
        __format_pbehavior(pbehavior)


def add_pbehavior_status(data):
    """Add "haspbehaviorinentities" and "hasallactivepbehaviorinentities" fields
    on every dict in data. Data must be a list of dict that contains a key
    "pbehavior" in order to work properly

    :param list data: the data to parse
    """
    for entity in data:
        enabled_list = [pbh["isActive"] for pbh in entity["pbehavior"]]

        if len(enabled_list) == 0:
            all_active = False
        else:
            all_active = all(enabled_list)

        one_more_active = any(enabled_list)

        entity["hasactivepbehaviorinentities"] = one_more_active
        entity["hasallactivepbehaviorinentities"] = all_active


def exports(ws):
    global LOGGER
    LOGGER = ws.logger

    ws.application.router.add_filter('mongo_filter', mongo_filter)
    ws.application.router.add_filter('id_filter', id_filter)

    @ws.application.route(
        '/api/v2/weather/selectors/<selector_filter:mongo_filter>')
    def get_selector(selector_filter):
        """
        Get a list of selectors from a mongo filter.

        :param dict selector_filter: a mongo filter to find selectors
        :rtype: dict
        """

        selector_filter['type'] = 'selector'
        selector_list = context_manager.get_entities(query=selector_filter)

        ws.logger.debug("Selector list: {}".format(selector_list))

        selectors = []
        for selector in selector_list:
            enriched_entity = {}
            tmp_alarm = alarm_manager.get(filter_={'d':
                                                   selector['name']})['alarms']

            enriched_entity['entity_id'] = selector['_id']
            enriched_entity['criticity'] = selector['infos'].get(
                'criticity', '')
            enriched_entity['org'] = selector['infos'].get('org', '')
            enriched_entity['sla_text'] = ''  # when sla
            enriched_entity['display_name'] = selector['name']
            if tmp_alarm != []:
                enriched_entity['state'] = tmp_alarm[0]['v']['state']
                enriched_entity['status'] = tmp_alarm[0]['v']['status']
                enriched_entity['snooze'] = tmp_alarm[0]['v']['snooze']
                enriched_entity['ack'] = tmp_alarm[0]['v']['ack']
                enriched_entity['connector'] = tmp_alarm[0]['v']['connector']
                enriched_entity['connector_name'] = (
                    tmp_alarm[0]['v']['connector_name']
                )
                enriched_entity['component'] = tmp_alarm[0]['v']['component']
                enriched_entity['component'] = tmp_alarm[0]['v']['component']
                if 'resource' in tmp_alarm[0]['v'].keys():
                    enriched_entity['resource'] = tmp_alarm[0]['v']['resource']
            else:
                enriched_entity['state'] = 0

            enriched_entity['linklist'] = []  # add this when it's ready
            add_pbehavior_info(enriched_entity)
            selectors.append(enriched_entity)

        add_pbehavior_status(selectors)

        return gen_json(selectors)

    @ws.application.route("/api/v2/weather/selectors/<selector_id:id_filter>")
    def weatherselectors(selector_id):
        """
        Get selector and contextual informations.

        :param str selector_id: the selector_id to search for
        :return: a list of agglomerated values of entities in the selector
        :rtype: list
        """
        try:
            selector_entity = context_manager.get_entities(
                query={'_id': selector_id, 'type': 'selector'})[0]
        except IndexError:
            json_error = {
                "name": "resource_not_found",
                "description": "the selector_id does not match any selector"
            }
            return gen_json_error(json_error, HTTP_NOT_FOUND)

        # Find entities with the selector filter
        try:
            query = json.loads(selector_entity['infos']['mfilter'])
        except:
            json_error = {
                "name": "filter_not_found",
                "description": "impossible to load the desired filter"
            }
            ws.logger.error(selector_entity['infos'])
            return gen_json_error(json_error, HTTP_NOT_FOUND)

        entities = context_manager.get_entities(query=query)

        entities_list = []
        for entity in entities:
            enriched_entity = {}
            tmp_alarm = alarm_manager.get(filter_={'d':
                                                   entity['_id']})['alarms']

            enriched_entity['entity_id'] = entity['_id']
            enriched_entity['sla_text'] = ''  # TODO when sla, use it
            enriched_entity['org'] = entity['infos'].get('org', '')
            enriched_entity['name'] = entity['name']
            enriched_entitity['source_type'] = entity['type']
            if tmp_alarm != []:
                enriched_entity['state'] = tmp_alarm[0]['v']['state']
                enriched_entity['status'] = tmp_alarm[0]['v']['status']
                enriched_entity['snooze'] = tmp_alarm[0]['v']['snooze']
                enriched_entity['ack'] = tmp_alarm[0]['v']['ack']
                enriched_entity['connector'] = tmp_alarm[0]['v']['connector']
                enriched_entity['connector_name'] = (
                    tmp_alarm[0]['v']['connector_name']
                )
                enriched_entity['component'] = tmp_alarm[0]['v']['component']
                if 'resource' in tmp_alarm[0]['v'].keys():
                    enriched_entity['resource'] = tmp_alarm[0]['v']['resource']
            else:
                enriched_entity['state'] = 0
            enriched_entity['linklist'] = []  # TODO wait for linklist

            add_pbehavior_info(enriched_entity)

            entities_list.append(enriched_entity)

        add_pbehavior_status(entities_list)

        return gen_json(entities_list)
