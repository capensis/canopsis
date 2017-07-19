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

from ast import literal_eval
import json
import copy

from canopsis.alerts.enums import AlarmField, AlarmFilterField
from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader
from canopsis.common.converters import mongo_filter, id_filter
from canopsis.common.utils import get_rrule_freq
from canopsis.context_graph.manager import ContextGraph
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_NOT_FOUND

context_manager = ContextGraph()
alarm_manager = Alerts()
alarmreader_manager = AlertsReader()
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
        active_eids = []
        next_action_timers = []

        if "mfilter" in entity:  # retreive pbehavior using the filter
            entities = context_manager.get_entities(
                literal_eval(entity["mfilter"])
            )

            eids = [ent["_id"] for ent in entities]

            pbh_active_list = pbehavior_manager.get_active_pbehaviors(eids)

            has_active_pbh = len(pbh_active_list) > 0

            for p_eid in [x['eids'] for x in pbh_active_list]:
                active_eids = active_eids + p_eid

            # as many active entity as all entities and at least one pbehavior
            has_all_active_pbh = set(eids) == set(active_eids) and len(active_eids) > 0

            # Enumerate all next_run timers
            for ent in entities:
                alarms = alarmreader_manager.get(
                    filter_={'d': '{}'.format(ent['name'])}
                )['alarms']
                if len(alarms) == 0:
                    continue

                alarmfilter = alarms[0]['v'].get(AlarmField.alarmfilter.value, {})
                action_timer = alarmfilter.get(AlarmFilterField.next_run.value, None)
                if action_timer is not None:
                    next_action_timers.append(action_timer)

        # has_active and has_all_active are exclude each one anothers
        has_active_pbh = has_active_pbh and not has_all_active_pbh

        entity["hasallactivepbehaviorinentities"] = has_all_active_pbh
        entity["hasactivepbehaviorinentities"] = has_active_pbh

        if len(next_action_timers) > 0:
            entity["automatic_action_timer"] = min(next_action_timers)

        # cleaning entity
        if "mfilter" in entity:
            del entity["mfilter"]

    return watchers


def watcher_status(watcher, pbehavior_eids_merged):
    """
        watcher_status

        :param dict watcher: watcher entity document
        :param set pbehavior_eids_merged: set with eids
        :return dict pbhevahior: dict with pbehavior infos has active
    """
    ret_dict = {
        'has_active_pbh': False,
        'has_all_active_pbh': False
    }
    bool_set = set([])
    for entity_id in watcher['depends']:
        bool_set.add(entity_id in pbehavior_eids_merged)

    if True in bool_set and False in bool_set:
        ret_dict['has_active_pbh'] = True
        return ret_dict
    elif True in bool_set:
        ret_dict['has_all_active_pbh'] = True
        return ret_dict
    return ret_dict


def get_active_pbehaviors_on_watchers(
    watchers_ids,
    active_pb_dict,
    active_pb_dict_full
):
    """
        get_active_pbehaviors_on_watchers.

        :param list watchers_ids:
        :param list active_pb_dict: list of dict with key: pbheavior_id value: set of eids
        :param list active_pb_dict_full: list of pbehavior dict
        :return dict: dict of watcher with list of active pbehavior

    """

    active_pb_on_watchers = {}
    for watcher_id in watchers_ids:
        tmp_pb = []
        for key, eids in active_pb_dict.items():
            if watcher_id in eids:
                tmp_pb.append(active_pb_dict_full[key])
        for pb in tmp_pb:
            pb['isActive'] = True
        active_pb_on_watchers[watcher_id] = tmp_pb

    return active_pb_on_watchers


def get_next_run_alert(watcher_depends, alert_next_run_dict):
    """
    get the next run of alarm filter

    :param watcher_depends: list of eids
    :param alert_next_run_dict: dict with next run infos for alarm filter
    :return: a timestamp with next alarm filter information or None
    """
    list_next_run = []
    for depend in watcher_depends:
        tmp_next_run = alert_next_run_dict.get(depend, None)
        if tmp_next_run:
            list_next_run.append(tmp_next_run)
    if list_next_run:
        return min(list_next_run)
    else:
        return None

def alert_not_ack_in_watcher(watcher_depends, alarm_dict):
    """
        alert_not_ack_in_watcher check if an alert is not ack in watcher depends

        :param watcher_depends: list of depends
        :param alarm_dict: alarm dict

    """
    for depend in watcher_depends:
        tmp_alarm = alarm_dict.get(depend, {})
        if tmp_alarm != {} and tmp_alarm.get('ack', None) is None:
            return True
    return False       


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

        depends_merged = set([])
        active_pb_dict = {}
        active_pb_dict_full = {}
        alarm_watchers_ids = []
        entity_watchers_ids = []
        alarm_dict = {}
        merged_pbehaviors_eids = set([])
        next_run_dict = {}
        watchers = []

        actives_pb = pbehavior_manager.get_all_active_pbehaviors()
        for pb in actives_pb:

            active_pb_dict[pb['_id']] = set(pb.get('eids', []))
            active_pb_dict_full[pb['_id']] = pb

        for watcher in watcher_list:
            for depends_id in watcher['depends']:
                depends_merged.add(depends_id)
            entity_watchers_ids.append(watcher['_id'])
            alarm_watchers_ids.append(
                '{0}/{1}'.format(watcher['_id'], watcher['name'])
            )
        active_pbehaviors = get_active_pbehaviors_on_watchers(
            entity_watchers_ids,
            active_pb_dict,
            active_pb_dict_full
        )
        for eids_tab in active_pb_dict.values():
            for eid in eids_tab:
                merged_pbehaviors_eids.add(eid)

        alarm_list = alarmreader_manager.get(
            filter_={}
        )['alarms']

        for alarm in alarm_list:
            alarm_dict[alarm['d']] = alarm['v']

        alerts_list_on_depends = alarmreader_manager.get(
            filter_={'d': {'$in': list(depends_merged)}}
        )['alarms']

        for alert in alerts_list_on_depends:
            if 'alarmfilter' in alert['v']:
                next_run_dict[alert['d']] = alert['v']['alarmfilter']['next_run']

        for watcher in watcher_list:
            enriched_entity = {}
            tmp_alarm = alarm_dict.get(
                '{0}/{1}'.format(watcher['_id'], watcher['name']),
                []
            )

            enriched_entity['entity_id'] = watcher['_id']
            enriched_entity['criticity'] = watcher['infos'].get('criticity', '')
            enriched_entity['org'] = watcher['infos'].get('org', '')
            enriched_entity['sla_text'] = ''  # when sla
            enriched_entity['display_name'] = watcher['name']
            enriched_entity['state'] = {'val': 0}
            enriched_entity['linklist'] = watcher['links']
            if tmp_alarm != []:
                enriched_entity['state'] = tmp_alarm['state']
                enriched_entity['status'] = tmp_alarm['status']
                enriched_entity['snooze'] = tmp_alarm['snooze']
                enriched_entity['ack'] = tmp_alarm['ack']
                enriched_entity['connector'] = tmp_alarm['connector']
                enriched_entity['connector_name'] = (
                    tmp_alarm['connector_name']
                )
                enriched_entity['component'] = tmp_alarm['component']
                if 'resource' in tmp_alarm.keys():
                    enriched_entity['resource'] = tmp_alarm['resource']

            enriched_entity["mfilter"] = watcher["mfilter"]

            enriched_entity['pbehavior'] = active_pbehaviors.get(
                watcher['_id'],
                []
            )
            enriched_entity['alerts_not_ack'] = alert_not_ack_in_watcher(watcher['depends'], alarm_dict)
            truc = watcher_status(watcher, merged_pbehaviors_eids)
            enriched_entity["hasallactivepbehaviorinentities"] = truc['has_all_active_pbh']
            enriched_entity["hasactivepbehaviorinentities"] = truc['has_active_pbh']
            tmp_next_run = get_next_run_alert(watcher.get('depends', []), next_run_dict)
            if tmp_next_run:
                enriched_entity['automatic_action_timer'] = tmp_next_run

            watchers.append(enriched_entity)

        return gen_json(watchers)

    @ws.application.route("/api/v2/weather/watchers/<watcher_id:id_filter>")
    def weatherwatchers(watcher_id):
        """
        Get a watcher and his contextual informations.

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
            query = json.loads(watcher_entity['mfilter'])
        except:
            json_error = {
                "name": "filter_not_found",
                "description": "impossible to load the desired filter"
            }
            return gen_json_error(json_error, HTTP_NOT_FOUND)

        raw_entities = context_manager.get_entities(query=query)
        entity_ids = [entity['_id'] for entity in raw_entities]
        enriched_entities = []

        entities = {}
        for raw_entity in raw_entities:
            reid = raw_entity['_id']
            entities[reid] = {
                'entity': raw_entity,
                'cur_alarm': None,
                'pbehaviors': []
            }

        tmp_alarms = alarmreader_manager.get(filter_={'d': {'$in': entity_ids}})
        alarms = tmp_alarms['alarms']
        for alarm in alarms:
            eid = alarm['d']
            if entities[eid]['cur_alarm'] is None:
                entities[eid]['cur_alarm'] = alarm['v']

        active_pbs = pbehavior_manager.get_all_active_pbehaviors()
        pmcp = pbehavior_manager._check_pbehavior

        for active_pb in active_pbs:
            active_pb_eids = set(active_pb['eids'])
            active_pb_dirty = copy.deepcopy(active_pb)
            active_pb_cleaned = __format_pbehavior(active_pb_dirty)

            for eid in active_pb_eids:
                active_pb_cleaned["isActive"] = pmcp(
                    entity_id=eid,
                    pb_names=[active_pb_cleaned['behavior']]
                )

                entities[eid]['pbehaviors'].append(active_pb_cleaned)

        for entity_id, entity in entities.iteritems():
            enriched_entity = {}

            current_alarm = entity['cur_alarm']
            raw_entity = entity['entity']
            enriched_entity['pbehavior'] = entity['pbehaviors']

            enriched_entity['entity_id'] = entity_id
            enriched_entity['sla_text'] = ''  # TODO when sla, use it
            enriched_entity['org'] = raw_entity['infos'].get('org', '')
            enriched_entity['name'] = raw_entity['name']
            enriched_entity['source_type'] = raw_entity['type']
            enriched_entity['state'] = {'val': 0}
            enriched_entity['linklist'] = raw_entity['links']

            if current_alarm:
                enriched_entity['state'] = current_alarm['state']
                enriched_entity['status'] = current_alarm['status']
                enriched_entity['snooze'] = current_alarm['snooze']
                enriched_entity['ack'] = current_alarm['ack']
                enriched_entity['connector'] = current_alarm['connector']
                enriched_entity['connector_name'] = (
                    current_alarm['connector_name']
                )
                enriched_entity['component'] = current_alarm['component']
                next_run = (current_alarm.get(AlarmField.alarmfilter.value, {})
                            .get(AlarmFilterField.next_run.value, None))
                enriched_entity['automatic_action_timer'] = next_run
                if 'resource' in current_alarm:
                    enriched_entity['resource'] = current_alarm['resource']

            enriched_entities.append(enriched_entity)

        return gen_json(enriched_entities)
