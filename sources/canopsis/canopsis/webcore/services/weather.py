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


"""Weather service routes"""

from __future__ import unicode_literals

import copy
import json
import time

from operator import itemgetter
from bottle import request

from canopsis.common.enumerations import DefaultEnum
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
from canopsis.watcher.filtering import WatcherFilter
from canopsis.alerts.enums import AlarmField, AlarmFilterField
from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader
from canopsis.common.converters import mongo_filter, id_filter
from canopsis.common.utils import get_rrule_freq
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_NOT_FOUND
from canopsis.common.influx import InfluxDBClient

alarm_manager = Alerts(*Alerts.provide_default_basics())
alarmreader_manager = AlertsReader(*AlertsReader.provide_default_basics())
context_manager = alarm_manager.context_manager
pbehavior_manager = PBehaviorManager(*PBehaviorManager.provide_default_basics())

WATCHER_COLLECTION = "default_entities"
mongo = MongoStore.get_default()
collection = mongo.get_collection(WATCHER_COLLECTION)
mongo_collection = MongoCollection(collection)

DEFAULT_LIMIT = '120'
DEFAULT_START = '0'
DEFAULT_SORT = False
DEFAULT_PB_TYPES = []
DEFAULT_DIRECTION = "ASC"

TILE_COLOR_PAUSE = "pause"
TILE_COLOR_OK = "ok"
TILE_COLOR_MINOR = "minor"
TILE_COLOR_MAJOR = "major"
TILE_COLOR_CRITICAL = "critical"
TILE_COLOR_SELECTOR = [TILE_COLOR_OK,
                       TILE_COLOR_MINOR,
                       TILE_COLOR_MAJOR,
                       TILE_COLOR_CRITICAL]

TILE_ICON_PAUSE = "pause"
TILE_ICON_MAINTENANCE = "maintenance"
TILE_ICON_OUT_SURVEILLANCE = "outOfSurveillance"
TILE_ICON_OK = "ok"
TILE_ICON_MINOR = "minor"
TILE_ICON_MAJOR = "major"
TILE_ICON_CRITICAL = "critical"
TILE_ICON_SELECTOR = [TILE_ICON_OK,
                      TILE_ICON_MINOR,
                      TILE_ICON_MAJOR,
                      TILE_ICON_CRITICAL]

class ResultKey(DefaultEnum):
    ID = "_id"
    INFOS = "infos"
    LINKS = "links"
    NAME = "name"
    STATE = "state"
    ALARM = "alarm"
    PBEHAVIORS = "pbehaviors"
    MFILTER = "mfilter"
    WATCHED_ENT_PBH = "watched_entities_pbehaviors"
    WATCHED_ENT_ALRM = "watched_entities_alarm"
    ALRM_VALUE = "v"
    ALRM_STATUS = "status"
    ALRM_STATE = "state"
    ALRM_SNOOZE = "snooze"
    ALRM_ACK = "ack"
    ALRM_CONNECTOR = "connector"
    ALRM_CONNECTOR_NAME = "connector_name"
    ALRM_LAST_UPDATE = "last_update_date"
    ALRM_COMPONENT = "component"
    ALRM_RESOURCE = "resource"
    ENT = "watched_entities"
    ENT_ID = "_id"


class __TileData:

    def __init__(self, watcher):
        self.entity_id = watcher[ResultKey.ID.value]
        self.infos = watcher[ResultKey.INFOS.value]
        self.sla_tex = ""
        self.display_name = watcher[ResultKey.NAME.value]
        self.linklist = []
        self.mfilter = watcher[ResultKey.MFILTER.value]
        for key, value in watcher[ResultKey.LINKS.value].items():
            self.linklist.append({'cat_name': key, 'links': value})

        state = watcher.get(ResultKey.STATE.value, 0)
        if isinstance(state, int):
            self.state = {'val': state}
        else:
            self.state = {'val': 0}

        if not len(watcher[ResultKey.ALARM.value]) == 0:
            alarm = watcher[ResultKey.ALARM.value][0]
            alarm = alarm[ResultKey.ALRM_VALUE.value]
            self.state = alarm[ResultKey.ALRM_STATE.value]
            self.status = alarm[ResultKey.ALRM_STATUS.value]
            self.snooze = alarm[ResultKey.ALRM_SNOOZE.value]
            self.ack = alarm[ResultKey.ALRM_ACK.value]
            self.connector = alarm[ResultKey.ALRM_CONNECTOR.value]
            self.connector_name = alarm[ResultKey.ALRM_CONNECTOR_NAME.value]
            self.last_update_date = alarm[ResultKey.ALRM_LAST_UPDATE.value]
            self.component = alarm[ResultKey.ALRM_COMPONENT.value]
            self.resource = alarm[ResultKey.ALRM_RESOURCE.value]

        # properties of the tile
        self.isActionRequired = self.__is_action_required(watcher)
        self.isAllEntitiesPaused = self.__is_all_entities_paused(watcher)
        self.isWatcherPaused = len(watcher[ResultKey.PBEHAVIORS.value]) != 0
        self.tileColor = self.__get_tile_color(watcher)
        self.tileIcon = self.__get_tile_icon(watcher)
        self.TileSecondaryIcon = self.__get_tile_secondary_icon(watcher)

    @classmethod
    def __is_action_required(cls, watcher):

        watcher_alarm = watcher.get(ResultKey.ALARM.value, None)
        if watcher_alarm is None:
            return False

        if len(watcher[ResultKey.PBEHAVIORS.value]) != 0:
            return False

        for entity in watcher[ResultKey.ENT.value]:
            if entity[ResultKey.ALARM.value] is None:
                continue

            if entity[ResultKey.ALARM.value]["v"].get("ack", None) is None:
                if len(entity[ResultKey.PBEHAVIORS.value]) == 0:
                    return True

        return False

    @classmethod
    def __is_all_entities_paused(cls, watcher):
        for entity in watcher[ResultKey.ENT.value]:
            if len(entity[ResultKey.PBEHAVIORS.value]) == 0:
                return False
        return True

    @classmethod
    def __get_tile_color(cls, watcher):
        watched_ent_paused = 0
        for ent in watcher[ResultKey.ENT.value]:
            if len(ent[ResultKey.PBEHAVIORS.value]) != 0:
                watched_ent_paused += 1

        if watched_ent_paused == len(watcher[ResultKey.ENT.value]) or \
           len(watcher[ResultKey.PBEHAVIORS.value]) != 0:
            return TILE_COLOR_PAUSE

        return TILE_COLOR_SELECTOR[watcher[ResultKey.STATE.value]]

    @classmethod
    def __get_tile_icon(cls, watcher):
        has_maintenance = False
        has_out_of_surveillance = False
        has_pause = False
        for pbh in watcher[ResultKey.PBEHAVIORS.value]:
            if pbh["type_"] == "Hors plage horaire de surveillance":
                has_out_of_surveillance = True
            elif pbh["type_"] == "Maintenance":
                has_maintenance = True
            elif pbh["type_"] in ["pause", "Pause"]:
                has_maintenance = True

        if has_maintenance:
            return TILE_ICON_MAINTENANCE
        if has_pause:
            return TILE_ICON_PAUSE
        if has_out_of_surveillance:
            return TILE_ICON_OUT_SURVEILLANCE

        return TILE_ICON_SELECTOR[watcher[ResultKey.STATE.value]]

    @classmethod
    def __get_tile_secondary_icon(cls, watcher):
        has_maintenance = False
        has_out_of_surveillance = False
        has_pause = False
        for ent in watcher[ResultKey.ENT.value]:
            for pbh in ent[ResultKey.PBEHAVIORS.value]:
                if pbh["type_"] == "Hors plage horaire de surveillance":
                    has_out_of_surveillance = True
                elif pbh["type_"] == "Maintenance":
                    has_maintenance = True
                elif pbh["type_"] in ["pause", "Pause"]:
                    has_maintenance = True

        if has_maintenance:
            return TILE_ICON_MAINTENANCE
        if has_pause:
            return TILE_ICON_PAUSE
        if has_out_of_surveillance:
            return TILE_ICON_OUT_SURVEILLANCE

        return None


def __format_pbehavior(pbehavior):
    """
    Rewrite a pbehavior from db format to front format.

    :param dict pbehavior: a pbehavior dict
    :return: a formatted pbehavior
    """
    EVERY = "Every {}"
    to_delete = [
        "connector", "filter", "connector_name", "eids"
    ]

    pbehavior["behavior"] = pbehavior.pop("name")
    pbehavior["dtstart"] = pbehavior.pop("tstart")
    pbehavior["dtend"] = pbehavior.pop("tstop")

    # parse the rrule to get is "text"
    rrule = {}
    rrule["rrule"] = pbehavior["rrule"]

    if pbehavior["rrule"] is not None:

        rrule_str = pbehavior["rrule"]
        if rrule_str[0:6] == "RRULE:":
            rrule_str = rrule_str[6:]

        freq = get_rrule_freq(rrule_str)

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


def get_ok_ko(influx_client, entity_id):
    """
    For an entity defined by its id, return the number of OK check and KO
    check.

    :param InfluxDBClient influx_client:
    :param str entity_id: the id of the entity
    :return: a dict with two key ok and ko or none if no data are found for
    the given entity.
    """
    query_sum = "SELECT SUM(ok) as ok, SUM(ko) as ko FROM " \
                "event_state_history WHERE \"eid\"='{}'"
    query_last_event = "SELECT LAST(\"ko\") FROM event_state_history WHERE " \
                       "\"eid\"='{}'"
    query_last_ko = query_last_event + " and \"ko\"=1"


    # Why did I use a double '\' ? It's simple, for some mystical reason,
    # somewhere between the call of influxdbstg.raw_query and the HTTP
    # request is sent, the escaped simple quote are deescaped. So like the
    # song says "you can't touch this".
    entity_id = entity_id.replace("'", "\\'")
    entity_id = entity_id.replace('"', '\\"')

    result = influx_client.query(query_sum.format(entity_id))

    stats = {}
    data = list(result.get_points())
    if len(data) > 0:
        data = data[0]
        stats["ok"] = data["ok"]
        stats["ko"] = data["ko"]

    result = influx_client.query(query_last_event.format(entity_id))
    data = list(result.get_points())
    if len(data) > 0:
        data = data[0]
        time = data["time"]
        time = time.replace("T", " ")
        time = time.replace("Z", "")
        stats["last_event"] = time

    result = influx_client.query(query_last_ko.format(entity_id))
    data = list(result.get_points())
    if len(data) > 0:
        data = data[0]
        time = data["time"]
        time = time.replace("T", " ")
        time = time.replace("Z", "")
        stats["last_ko"] = time

    if len(stats) > 0:
        return stats

    return None


def pbehavior_types(pbehaviors):
    """
    Return a set containing all type_ found in pbehaviors.
    :param pbehaviors
    """
    pb_types = set()

    for pb in pbehaviors:
        pb_type = pb.get('type_', None)
        if pb_type is not None:
            pb_types.add(pb_type)

    return pb_types


def watcher_status(watcher, pbehavior_eids_merged):
    """
    watcher_status

    :param dict watcher: watcher entity document
    :param set pbehavior_eids_merged: set with eids
    :returns: has active pb status (or all), has all active pb status
    :rtype: (bool, bool)
    """
    bool_set = set([e in pbehavior_eids_merged for e in watcher['depends']])

    at_least_one = True in bool_set
    if at_least_one and False in bool_set:
        # has_active_pbh
        return True, False
    elif at_least_one:
        # has_all_active_pbh
        return True, True

    return False, False


def get_active_pbehaviors_on_watchers(watchers,
                                      active_pb_dict,
                                      active_pb_dict_full):
    """
    get_active_pbehaviors_on_watchers.

    :param list watchers:
    :param list active_pb_dict:
    :param list active_pb_dict_full: list of pbehavior dict
    :returns: dict of watcher with list of active pbehavior
    """
    active_pb_on_watchers = {}
    active_watcher_pbehaviors = {}

    for watcher in watchers:
        tmp_pbh = []
        tmp_wpbh = []
        watcher_depends = set(watcher.get('depends', []))

        for pb_id, eids in active_pb_dict.items():
            # add pbehaviors linked to this watcher's entities
            for eid in eids:
                if eid in watcher_depends:
                    tmp_pbh.append(active_pb_dict_full[pb_id])

            # add pbehaviors linked to this watcher
            if watcher['_id'] in active_pb_dict[pb_id]:
                wpb = active_pb_dict_full[pb_id]
                wpb['isActive'] = True
                tmp_wpbh.append(wpb)

        for pbh in tmp_pbh:
            pbh['isActive'] = True

        active_pb_on_watchers[watcher['_id']] = tmp_pbh
        active_watcher_pbehaviors[watcher['_id']] = tmp_wpbh

    return active_pb_on_watchers, active_watcher_pbehaviors


def is_action_required(watcher, alarm_dict, active_pbehaviors, active_watchers_pbehaviors):

    watcher_alarm = alarm_dict.get(watcher["_id"], None)
    if watcher_alarm is None:
        return False

    entities_alarm = {}
    entities_pbh = {}
    for key in watcher["depends"]:
        entities_alarm[key] = alarm_dict.get(key, None)
        entities_pbh[key] = active_pbehaviors.get(key, None)

    w_pbh = active_watchers_pbehaviors[watcher["_id"]]
    if len(w_pbh) != 0:
        return False

    for entity in entities_alarm:
        if entities_alarm[entity] is None:
            continue

        if entities_alarm[entity].get("ack", None) is None:
            if entities_pbh[entity] is None:
                return True

    return False

def remove_inactive_pbh(pbehaviors):
    now = time.time()

    active_pbh = []
    for pbh in pbehaviors:
        if pbehavior_manager.check_active_pbehavior(now, pbh):
            active_pbh.append(pbh)

    return active_pbh

def get_next_run_alert(watcher_depends, alert_next_run_dict):
    """
    get the next run of alarm filter

    :param watcher_depends: list of eids
    :param alert_next_run_dict: dict with next run infos for alarm filter
    :returns: a timestamp with next alarm filter information or None
    """
    list_next_run = []
    for depend in watcher_depends:
        tmp_next_run = alert_next_run_dict.get(depend, None)
        if tmp_next_run:
            list_next_run.append(tmp_next_run)
    if list_next_run:
        return min(list_next_run)

    return None


def alert_not_ack_in_watcher(watcher_depends, alarm_dict):
    """
    alert_not_ack_in_watcher check if an alert is not ack in watcher depends

    :param watcher_depends: list of depends
    :param alarm_dict: alarm dict
    :rtype: bool
    """
    for depend in watcher_depends:
        tmp_alarm = alarm_dict.get(depend, {})
        if (tmp_alarm != {}
                and tmp_alarm.get('ack', None) is None
                and tmp_alarm.get('state', {}).get('val', 0) != 0):
            return True

    return False


def _parse_direction(direction):
    if direction == "ASC":
        return 1
    elif direction == "DESC":
        return -1
    else:
        ValueError("Direction must be 'ASC' or 'DESC' not {}.".format(
            direction))


def exports(ws):
    ws.application.router.add_filter('mongo_filter', mongo_filter)
    ws.application.router.add_filter('id_filter', id_filter)

    influx_client = InfluxDBClient.from_configuration(ws.logger)

    @ws.application.route(
        '/api/v2/weather/watchers/<watcher_filter:mongo_filter>'
    )
    def get_watcher(watcher_filter):
        """
        Get a list of watchers from a mongo filter.

        :param dict watcher_filter: a mongo filter to find watchers
        :rtype: dict
        """
        limit = request.query.limit or DEFAULT_LIMIT
        start = request.query.start or DEFAULT_START
        orderby = request.query.orderby or None
        direction = request.query.direction or None


        # FIXIT: service weather has no pagination capability at all.
        try:
            start = int(start)
        except ValueError:
            start = int(DEFAULT_START)
        try:
            limit = int(limit)
        except ValueError:
            limit = int(DEFAULT_LIMIT)

        wf = WatcherFilter()
        watcher_filter['type'] = 'watcher'
        watcher_filter = wf.filter(watcher_filter)

        # select the watchers
        select_watcher_stage = {"$match": watcher_filter}

        # pagination
        skip = {"$skip": start}
        limit = {"$limit": limit}

        # retreive opened alarm for the watchers
        alarms = {"$graphLookup":
                  {"from": "periodical_alarm",
                   "startWith": "$_id",
                   "connectFromField": "_id",
                   "connectToField": "d",
                   "restrictSearchWithMatch": {'v.resolved': None},
                   "as": "alarm",
                   "maxDepth": 0
                  }
        }

        # retreive every pbehaviors on the watcher
        pbehaviors = {"$lookup":
                      {"from": "default_pbehavior",
                       "localField": "_id",
                       "foreignField": "eids",
                       "as": "pbehaviors"
                      }
        }

        # retrieve watched entities
        entities = {"$lookup": {
            "from": "default_entities",
            "localField": "depends",
            "foreignField": "_id",
            "as": "watched_entities",
            }
        }

        # retreive every pbehaviors on the watched entities
        pbehaviors_watched_ent = {"$graphLookup":
                                  {"from": "default_pbehavior",
                                   "startWith": "$watched_entities._id",
                                   "connectFromField": "watched_entities._id",
                                   "connectToField": "eids",
                                   "maxDepth": 0,
                                   "as": "watched_entities_pbehaviors",
                                  }
        }

        # retreive every opened alarm on the watched entities
        alarm_watched_ent = {"$graphLookup":
                             {"from": "periodical_alarm",
                              "startWith": "$watched_entities._id",
                              "connectFromField": "_id",
                              "connectToField": "d",
                              "restrictSearchWithMatch": {'v.resolved': None},
                              "as": "watched_entities_alarm",
                              "maxDepth": 0
                             }
        }

        pipeline = [select_watcher_stage,
                    skip,
                    limit,
                    alarms,
                    pbehaviors,
                    entities,
                    pbehaviors_watched_ent,
                    alarm_watched_ent]

        # retreive
        if orderby is not None:
            # TODO if needed, set the correction direction value
            pipeline.insert(1, {"$sort": {orderby: direction}})

        pipeline_result = mongo_collection.aggregate(pipeline)

        result = []

        for watcher in pipeline_result:
            # remove the inactive pbehaviors from the pipeline result
            pbhs = watcher[ResultKey.PBEHAVIORS.value]
            watcher[ResultKey.PBEHAVIORS.value] = remove_inactive_pbh(pbhs)
            pbhs = watcher[ResultKey.WATCHED_ENT_PBH.value]
            watcher[ResultKey.WATCHED_ENT_PBH.value] = remove_inactive_pbh(pbhs)

            # assign entities pbehaviors to the correct entities
            entities = {}
            for entity in watcher[ResultKey.ENT.value]:
                entity[ResultKey.PBEHAVIORS.value] = []
                entity[ResultKey.ALARM.value] = None
                entities[entity[ResultKey.ENT_ID.value]] = entity

            for pbh in watcher[ResultKey.WATCHED_ENT_PBH.value]:
                for ent_id in pbh["eids"]:
                    try:
                        entities[ent_id][ResultKey.PBEHAVIORS.value].append(pbh)
                    except KeyError:
                        ws.logger.error("Can not find entities {} in the"
                                        "pipeline result".format(ent_id))

            for alarm in watcher[ResultKey.WATCHED_ENT_ALRM.value]:
                try:
                    entities[alarm["d"]][ResultKey.ALARM.value] = alarm
                except KeyError:
                    ws.logger.error("Can not find entities {} in the"
                                    "pipeline result".format(alarm["d"]))

            watcher[ResultKey.ENT.value] = entities.values()
            del watcher[ResultKey.WATCHED_ENT_PBH.value]
            del watcher[ResultKey.WATCHED_ENT_ALRM.value]

            tileData = __TileData(watcher)
            result.append(vars(tileData))

        return gen_json(result)

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
        # when entities is in watcher_entity, the watcher is handled in go engines
        # thus the mfilter query is not needed and not present as well
        if "entities" in watcher_entity:
            query = {"_id":{"$in": watcher_entity.get('depends', [])}}
        else:
            try:
                query = json.loads(watcher_entity['mfilter'])
            except (ValueError, KeyError, TypeError):
                json_error = {
                    "name": "filter_not_found",
                    "description": "impossible to load the desired filter"
                }
                return gen_json_error(json_error, HTTP_NOT_FOUND)

        query["enabled"] = True

        raw_entities = context_manager.get_entities(
            query=query,
            with_links=True
        )
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

        for active_pb in active_pbs:
            active_pb_eids = set(active_pb['eids'])
            active_pb_dirty = copy.deepcopy(active_pb)
            active_pb_cleaned = __format_pbehavior(active_pb_dirty)

            for eid in active_pb_eids:
                active_pb_cleaned['isActive'] = True
                if eid in entities:
                    entities[eid]['pbehaviors'].append(active_pb_cleaned)

        for entity_id, entity in entities.iteritems():
            enriched_entity = {}

            current_alarm = entity['cur_alarm']
            raw_entity = entity['entity']

            tmp_links = []
            for k, val in raw_entity['links'].items():
                tmp_links.append({'cat_name': k, 'links': val})

            enriched_entity['pbehavior'] = entity['pbehaviors']
            enriched_entity['entity_id'] = entity_id
            enriched_entity['linklist'] = tmp_links
            enriched_entity['infos'] = raw_entity['infos']
            enriched_entity['sla_text'] = ''  # TODO when sla, use it
            enriched_entity['org'] = raw_entity['infos'].get('org', '')
            enriched_entity['name'] = raw_entity['name']
            enriched_entity['source_type'] = raw_entity['type']
            enriched_entity['state'] = {'val': 0}
            enriched_entity['stats'] = get_ok_ko(influx_client, entity_id)
            if current_alarm is not None:
                enriched_entity['alarm_creation_date'] = current_alarm.get("creation_date")
                enriched_entity['alarm_display_name'] = current_alarm.get("display_name")
                enriched_entity['ticket'] = current_alarm.get('ticket')
                enriched_entity['state'] = current_alarm['state']
                enriched_entity['status'] = current_alarm['status']
                enriched_entity['snooze'] = current_alarm.get('snooze')
                enriched_entity['ack'] = current_alarm.get('ack')
                enriched_entity['connector'] = current_alarm['connector']
                enriched_entity['connector_name'] = (
                    current_alarm['connector_name']
                )
                enriched_entity['last_update_date'] = current_alarm.get(
                    'last_update_date', None
                )
                enriched_entity['component'] = current_alarm['component']
                next_run = (current_alarm.get(AlarmField.alarmfilter.value, {})
                            .get(AlarmFilterField.next_run.value, None))
                enriched_entity['automatic_action_timer'] = next_run
                if current_alarm.get('resource', ''):
                    enriched_entity['resource'] = current_alarm['resource']

            enriched_entities.append(enriched_entity)

        enriched_entities = sorted(enriched_entities, key=itemgetter("name"))

        return gen_json(enriched_entities)
