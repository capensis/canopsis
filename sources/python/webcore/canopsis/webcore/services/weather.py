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

from ast import literal_eval

LOGGER = None

context_manager = ContextGraph()
alarm_manager = AlertsReader()
pbehavior_manager = PBehaviorManager()


def __format_pbehavior(pbehavior):
    """Rewrite en pbehavior from db format to front format.

    :param dict pbehavior: a pbehavior dict
    :return: a formatted pbehavior
    """
    EVERY = "Every {0}"
    to_delete = [
        "_id", "connector", "author", "comments", "filter", "connector_name",
        "eids"
    ]

    pbehavior["behavior"] = pbehavior.pop("name")
    pbehavior["dtstart"] = pbehavior.pop("tstart")
    pbehavior["dtend"] = pbehavior.pop("tstop")

    # parse the rrule to get is "text"
    rrule = {}
    rrule["rrule"] = pbehavior["rrule"]

    if pbehavior["rrule"] is not None:
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

    return pbehavior


def add_pbehavior_info(enriched_entity):
    """Add pbehavior related field to selectors. This function will add
    the related pbehavior in 'pbehavior'.

    :param dict enriched_entity: the entity to enrich
    """

    enriched_entity["pbehavior"] = pbehavior_manager.get_pbehaviors_by_eid(
        enriched_entity['entity_id'])

    pmcp = pbehavior_manager._check_pbehavior
    for pbehavior in enriched_entity["pbehavior"]:
        pbehavior = __format_pbehavior(pbehavior)

        pbehavior["isActive"] = pmcp(entity_id=enriched_entity['entity_id'],
                                     pb_names=[pbehavior['behavior']])


def add_pbehavior_status(watchers):
    """Add "haspbehaviorinentities" and "hasallactivepbehaviorinentities" fields
    on every dict in data. Data must be a list of dict that contains a key
    "pbehavior" in order to work properly.

    If the field "mfilter" is present in the element of data, ignore
    the pbehavior present in the element en retreive them directly from
    database. Then remove the field "mfilter".

    :param list watchers: the watchers to parse
    """
    for entity in watchers:

        has_active_pbh = False
        has_all_active_pbh = False
        act_eids = []
        if "mfilter" in entity:  # retreive pbehavior using the filter
            entities = context_manager.get_entities(
                literal_eval(entity["mfilter"]),
                {"_id": 1}
            )

            eids = [ent["_id"] for ent in entities]

            pbh_active_list = pbehavior_manager.get_active_pbehaviors(eids)

            has_active_pbh = len(pbh_active_list) > 0

            for p_eid in [x['eids'] for x in pbh_active_list]:
                act_eids = act_eids + p_eid

            # as many active entity as all entities and at least one pbehavior
            has_all_active_pbh = set(eids) == set(act_eids) and len(act_eids) > 0

        # has_active and has_all_active are exclude each one anothers
        has_active_pbh = has_active_pbh and not has_all_active_pbh

        entity["hasallactivepbehaviorinentities"] = has_all_active_pbh
        entity["hasactivepbehaviorinentities"] = has_active_pbh

        # cleaning entity
        if "mfilter" in entity:
            del entity["mfilter"]


def exports(ws):
    global LOGGER
    LOGGER = ws.logger

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
                filter_={'d': '{0}/{1}'.format(
                    watcher['_id'],
                    watcher['name']
                )}
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
                enriched_entity['connector'] = tmp_alarm[0]['v']['connector']
                enriched_entity['connector_name'] = (
                    tmp_alarm[0]['v']['connector_name']
                )
                enriched_entity['component'] = tmp_alarm[0]['v']['component']
                enriched_entity['component'] = tmp_alarm[0]['v']['component']
                if 'resource' in tmp_alarm[0]['v'].keys():
                    enriched_entity['resource'] = tmp_alarm[0]['v']['resource']
            else:
                enriched_entity['state'] = {'val': 0}

            enriched_entity['linklist'] = []  # add this when it's ready
            enriched_entity["mfilter"] = watcher["infos"]["mfilter"]
            add_pbehavior_info(enriched_entity)
            watchers.append(enriched_entity)

        add_pbehavior_status(watchers)

        return gen_json(watchers)

    @ws.application.route("/api/v2/weather/watchers/<watcher_id:id_filter>")
    def weatherwatchers(watcher_id):
        """
        Get watcher and contextual informations.

        :param str watcher_id: the watcher_id to search for
        :return: a list of agglomerated values of entities in the watcher
        :rtype: list
        """
        try:
            watcher_entity = context_manager.get_entities(
                query={'_id': watcher_id, 'type': 'watcher'})[0]
        except IndexError:
            json_error = {
                "name": "resource_not_found",
                "description": "the watcher_id does not match any watcher"
            }
            return gen_json_error(json_error, HTTP_NOT_FOUND)

        # Find entities with the watcher filter
        try:
            query = json.loads(watcher_entity['infos']['mfilter'])
        except:
            json_error = {
                "name": "filter_not_found",
                "description": "impossible to load the desired filter"
            }
            ws.logger.error(watcher_entity['infos'])
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
            enriched_entity['source_type'] = entity['type']
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
                enriched_entity['state'] = {'val': 0}
            enriched_entity['linklist'] = []  # TODO wait for linklist

            add_pbehavior_info(enriched_entity)

            entities_list.append(enriched_entity)

        return gen_json(entities_list)
