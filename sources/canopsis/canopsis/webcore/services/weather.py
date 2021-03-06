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

from canopsis.common.enumerations import FastEnum
from canopsis.models.entity import Entity
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
from canopsis.watcher.filtering import WatcherFilter
from canopsis.alerts.enums import AlarmField, AlarmFilterField
from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader
from canopsis.common.converters import mongo_filter, id_filter
from canopsis.common.utils import get_rrule_freq
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.stat.manager import StatManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_NOT_FOUND

from pymongo.errors import PyMongoError
from montydb.utils import MontyList


alarm_manager = Alerts(*Alerts.provide_default_basics())
alarmreader_manager = AlertsReader(*AlertsReader.provide_default_basics())
context_manager = alarm_manager.context_manager
pbehavior_manager = PBehaviorManager(*PBehaviorManager.provide_default_basics())

WATCHER_COLLECTION = "default_entities"
mongo = MongoStore.get_default()
collection = mongo.get_collection(WATCHER_COLLECTION)

mongo_collection = MongoCollection(collection)

DEFAULT_LIMIT = 120
DEFAULT_START = '0'
DEFAULT_SORT = False
DEFAULT_PB_TYPES = []
DEFAULT_DIRECTION = "ASC"


class TileColor(FastEnum):
    PAUSE = "pause"
    OK = "ok"
    MINOR = "minor"
    MAJOR = "major"
    CRITICAL = "critical"


TILE_COLOR_SELECTOR = [TileColor.OK,
                       TileColor.MINOR,
                       TileColor.MAJOR,
                       TileColor.CRITICAL]


class TileIcon(FastEnum):
    PAUSE = "pause"
    MAINTENANCE = "maintenance"
    UNMONITORED = "unmonitored"
    OK = "ok"
    MINOR = "minor"
    MAJOR = "major"
    CRITICAL = "critical"


TILE_ICON_SELECTOR = [TileIcon.OK,
                      TileIcon.MINOR,
                      TileIcon.MAJOR,
                      TileIcon.CRITICAL]


class ResultKey(FastEnum):
    """
    Contains the key use to handle the watcher retrieve from the
    watcher pipeline and the rearrange watcher.
    """
    ID = AlarmField._id.value
    INFOS = Entity.INFOS
    LINKS = "links"
    NAME = Entity.NAME
    STATE = "state"
    ALARM = "alarm"
    PBEHAVIORS = "pbehaviors"
    MFILTER = "mfilter"
    WATCHED_ENT_PBH = "watched_entities_pbehaviors"
    WATCHED_ENT_ALRM = "watched_entities_alarm"
    ALRM_VALUE = "v"
    ALRM_STATUS = AlarmField.status.value
    ALRM_STATE = AlarmField.state.value
    ALRM_SNOOZE = AlarmField.snooze.value
    ALRM_ACK = AlarmField.ack.value
    ALARM_FILTER = AlarmField.alarmfilter.value
    ALRM_CONNECTOR = "connector"
    ALRM_CONNECTOR_NAME = "connector_name"
    ALRM_LAST_UPDATE = AlarmField.last_update_date.value
    ALRM_COMPONENT = "component"
    ALRM_RESOURCE = "resource"
    ALRM_OUTPUT = "output"
    ENT = "watched_entities"
    ENT_ID = Entity._ID


class __TileData(object):
    """
    This object represents one element (a tile) of the array return by the
    weather API. Every __TileData attribute represent a field of the returned
    elements.

    Be careful, this object is used with the vars built-in method. Any new
    instance attribute will be returned by the API.
    """

    def __init__(self, watcher):
        """
        Create a new instance.

        :param dict: a watcher dict from the watcher pipeline.
        """
        self.watcher = copy.deepcopy(watcher)
        self.entity_id = watcher[ResultKey.ID]
        self.infos = watcher[ResultKey.INFOS]
        self.sla_tex = ""
        self.display_name = watcher[ResultKey.NAME]
        self.linklist = []
        self.mfilter = watcher.get(ResultKey.MFILTER, "")
        for key, value in watcher.get(ResultKey.LINKS, {}).items():
            self.linklist.append({'cat_name': key, 'links': value})

        self.watcher_pbehavior = watcher[ResultKey.PBEHAVIORS]

        self.automatic_action_timer = self.__get_next_run(watcher)

        state = watcher.get(ResultKey.STATE, 0)
        if isinstance(state, int):
            self.state = {'val': state}
        else:
            self.state = {'val': 0}

        if not len(watcher[ResultKey.ALARM]) == 0:
            alarm = watcher[ResultKey.ALARM][0]
            alarm = alarm[ResultKey.ALRM_VALUE]
            self.state = alarm[ResultKey.ALRM_STATE]
            self.status = alarm[ResultKey.ALRM_STATUS]
            self.snooze = alarm.get(ResultKey.ALRM_SNOOZE, None)
            self.ack = alarm.get(ResultKey.ALRM_ACK, None)
            self.connector = alarm[ResultKey.ALRM_CONNECTOR]
            self.connector_name = alarm[ResultKey.ALRM_CONNECTOR_NAME]
            self.last_update_date = alarm[ResultKey.ALRM_LAST_UPDATE]
            self.component = alarm[ResultKey.ALRM_COMPONENT]
            self.resource = alarm.get(ResultKey.ALRM_RESOURCE, None)
            self.output = alarm.get(ResultKey.ALRM_OUTPUT, None)

        # properties of the tile
        self.isActionRequired = self.__is_action_required(watcher)
        self.isAllEntitiesPaused = self.__is_all_entities_paused(watcher)
        self.isWatcherPaused = len(watcher[ResultKey.PBEHAVIORS]) != 0
        self.tileColor = self.__get_tile_color(watcher)
        self.tileIcon = self.__get_tile_icon(watcher)
        self.tileSecondaryIcon = self.__get_tile_secondary_icon(watcher)
        self.watcher["tileColor"] = self.tileColor
        self.watcher["tileSecondaryIcon"] = self.tileSecondaryIcon
        self.watcher["tileIcon"] = self.tileIcon

    def remove_tiled_watcher(self):
        del self.watcher

    @classmethod
    def __is_action_required(cls, watcher):
        """
        Return if an action is required on the watcher.

        An action is required if a watched entity have an opened alarm without
        any `ack`.
        :param dict watcher: a watcher with his pbehaviors, watched entities
        and their active pbehaviors(see _rework_watcher_pipeline_element
        function).
        :return boolean: True if an action is required, False otherwise.
        """

        watcher_alarm = watcher.get(ResultKey.ALARM, None)
        if watcher_alarm is None:
            return False

        if len(watcher[ResultKey.PBEHAVIORS]) != 0:
            return False

        for entity in watcher[ResultKey.ENT]:
            if entity[ResultKey.ALARM] is None:
                continue

            alarm_value = entity[ResultKey.ALARM][ResultKey.ALRM_VALUE]
            if alarm_value.get(ResultKey.ALRM_ACK, None) is None:
                if len(entity[ResultKey.PBEHAVIORS]) == 0:
                    return True

        return False

    @classmethod
    def __is_all_entities_paused(cls, watcher):
        """
        Return if every watched entities of a watcher have at least
        one active pbehavior.

        :param dict watcher: a watcher with his pbehaviors, watched entities
        and their active pbehaviors(see _rework_watcher_pipeline_element
        function).
        :return boolean: True if all watched entities have an active
        pbehavior, False otherwise.
        """
        for entity in watcher[ResultKey.ENT]:
            if len(entity[ResultKey.PBEHAVIORS]) == 0:
                return False
        return True

    @classmethod
    def __get_tile_color(cls, watcher):
        """
        Return a string that indicate the color to use to render the watcher
        tile.

        :param dict watcher: a watcher with his pbehaviors, watched entities
        and their active pbehaviors(see _rework_watcher_pipeline_element
        function).
        :return str: 'TileColor.PAUSE' if they are at least one active
        pbehavior on the watcher or an active pbehavior on every watched
        entities.
        'TileColor.OK' if the watcher state is 0.
        'TileColor.MINOR' if the watcher state is 1.
        'TileColor.MAJOR' if the watcher state is 2.
        'TileColor.CRITICAL' if the watcher state is 3.
        """
        watched_ent_paused = 0
        watcher_state = 0

        for ent in watcher[ResultKey.ENT]:
            if len(ent[ResultKey.PBEHAVIORS]) > 0:
                watched_ent_paused += 1

        if len(watcher[ResultKey.ALARM]) > 0:
            alarm = watcher[ResultKey.ALARM][0][ResultKey.ALRM_VALUE]
            watcher_state = alarm[ResultKey.ALRM_STATE]["val"]

        if len(watcher[ResultKey.ENT]) > 0 and \
           len(watcher[ResultKey.ENT]) == watched_ent_paused:
            return TileColor.PAUSE

        if len(watcher[ResultKey.PBEHAVIORS]) > 0:
            return TileColor.PAUSE

        return TILE_COLOR_SELECTOR[watcher_state]

    @classmethod
    def __get_tile_icon(cls, watcher):
        """
        Return a string that indicate the primary tile icon to used to render
        the watcher tile.

        'TileIcon.OK', 'TileIcon.MINOR', 'TileIcon.MAJOR',
        'TileIcon.CRITICAL' are return if they are no active pbehavior on the
        watcher or if not every watched entities have an active pbehavior.

        'TileIcon.PAUSE', 'TileIcon.MAINTENANCE',
        'TileIcon.UNMONITORED' are display if the given watcher is under
        an active pbehavior or if every watched entities have at least one
        active pbehavior. The icon string returned depends on the 'type_' of
        active pbehaviors on the watcher and on the watched entities. If
        at least one 'maintenance' pbehavior is present,
        'TileIcon.MAINTENANCE' will be returned. Then if at least one
        'Hors plage horaire de surveillance' pbehavior is present,
        'TileIcon.UNMONITORED' will be returned. Finally, if no
        'maintenance' pbehavior or 'Hors plage horaire de surveillance'
        pbehavior are present, return 'TileIcon.PAUSE'.
        """
        has_maintenance = False
        has_out_of_surveillance = False
        has_pause = False

        watched_ent_paused = 0

        for ent in watcher[ResultKey.ENT]:
            if len(ent[ResultKey.PBEHAVIORS]) > 0:
                watched_ent_paused += 1

        if watched_ent_paused == len(watcher[ResultKey.ENT]) and\
           len(watcher[ResultKey.PBEHAVIORS]) == 0:

            for ent in watcher[ResultKey.ENT]:
                for pbh in ent[ResultKey.PBEHAVIORS]:
                    if pbh["type_"] == "Hors plage horaire de surveillance":
                        has_out_of_surveillance = True
                    elif pbh["type_"] == "Maintenance":
                        has_maintenance = True
                    elif pbh["type_"] in ["pause", "Pause"]:
                        has_pause = True

        else:
            for pbh in watcher[ResultKey.PBEHAVIORS]:
                if pbh["type_"] == "Hors plage horaire de surveillance":
                    has_out_of_surveillance = True
                elif pbh["type_"] == "Maintenance":
                    has_maintenance = True
                elif pbh["type_"] in ["pause", "Pause"]:
                    has_pause = True

        if has_maintenance:
            return TileIcon.MAINTENANCE
        if has_pause:
            return TileIcon.PAUSE
        if has_out_of_surveillance:
            return TileIcon.UNMONITORED

        watcher_state = 0
        if len(watcher[ResultKey.ALARM]) > 0:
            alarm = watcher[ResultKey.ALARM][0][ResultKey.ALRM_VALUE]
            watcher_state = alarm[ResultKey.ALRM_STATE]["val"]

        return TILE_ICON_SELECTOR[watcher_state]


    @classmethod
    def __get_tile_secondary_icon(cls, watcher):
        """
        Return a string that indicate the secondary tile icon to used to render
        the watcher tile or None

        'TileIcon.PAUSE', 'TileIcon.MAINTENANCE',
        'TileIcon.UNMONITORED' are returned if they are some (not all)
        watched entities with an active pbehavior.

        The icon string returned depends on the 'type_' of active pbehaviors on
        the watche watched entities. If at least one 'maintenance' pbehavior is
        present, 'TileIcon.MAINTENANCE' will be returned. Then if at least one
        'Hors plage horaire de surveillance' pbehavior is present,
        'TileIcon.UNMONITORED' will be returned. Finally, if no
        'maintenance' pbehavior or 'Hors plage horaire de surveillance'
        pbehavior are present, return 'TileIcon.PAUSE'.

        If every watched entities are under an active pbehavior, they are
        no secondary icon displayed on the tile, so None are returned.

        :param dict watcher: a watcher with his pbehaviors, watched entities
        and their active pbehaviors(see _rework_watcher_pipeline_element
        function).
        :return str: 'TileIcon.PAUSE' or 'TileIcon.MAINTENANCE' or
        'TileIcon.UNMONITORED' or None
        """
        has_maintenance = False
        has_out_of_surveillance = False
        has_pause = False
        paused_watched_ent = 0

        for ent in watcher[ResultKey.ENT]:
            if not ent["enabled"]:
                continue
            if len(ent[ResultKey.PBEHAVIORS]) > 0:
                paused_watched_ent += 1
            for pbh in ent[ResultKey.PBEHAVIORS]:
                if pbh["type_"] == "Hors plage horaire de surveillance":
                    has_out_of_surveillance = True
                elif pbh["type_"] == "Maintenance":
                    has_maintenance = True
                elif pbh["type_"] in ["pause", "Pause"]:
                    has_pause = True

        if paused_watched_ent == len(watcher[ResultKey.ENT]):
            return None

        if has_maintenance:
            return TileIcon.MAINTENANCE
        if has_pause:
            return TileIcon.PAUSE
        if has_out_of_surveillance:
            return TileIcon.UNMONITORED

        return None

    @classmethod
    def __get_next_run(cls, watcher):
        """
        Return the smallest next_run field value from all the watched entities
        alarms.

        :param dict watcher: a watcher with his pbehaviors, watched entities
        and their active pbehaviors(see _rework_watcher_pipeline_element
        function).
        :return int: the smallest next_run.
        """
        next_runs = []
        for ent in watcher[ResultKey.ENT]:
            if ent[ResultKey.ALARM] is not None:
                alarm = ent[ResultKey.ALARM]

                alarmfilter = None
                try:
                    alarmfilter = alarm[ResultKey.ALRM_VALUE]
                    alarmfilter = alarmfilter[ResultKey.ALARM_FILTER]
                except KeyError:
                    continue

                if isinstance(alarmfilter, dict) \
                   and AlarmFilterField.next_run.value in alarmfilter:
                    next_runs.append(alarmfilter[AlarmFilterField.next_run.value])

        if len(next_runs) > 0:
            return min(next_runs)
        return None


def _pbehavior_types(watcher):
    """
    Return a set containing all type_ found in pbehaviors.

    :param dict watcher: one element from the query
    :return set: a set of string.
    """
    pb_types = set()

    pbehaviors = watcher[ResultKey.PBEHAVIORS][:]  # create a new list
    for ent in watcher[ResultKey.ENT]:
        pbehaviors += ent[ResultKey.PBEHAVIORS]

    for pbh in pbehaviors:
        pb_type = pbh.get('type_', None)
        if pb_type is not None:
            pb_types.add(pb_type)

    return pb_types


def _watcher_status(watcher):

    ent_with_active_pbh = set()

    for ent in watcher[ResultKey.ENT]:
        ent_with_active_pbh.add(len(ent[ResultKey.PBEHAVIORS]) > 0)

    at_least_one = True in ent_with_active_pbh
    if at_least_one and False in ent_with_active_pbh:
        # has_active_pbh
        return True, False
    elif at_least_one:
        # has_all_active_pbh
        return True, True

    return False, False


def _remove_inactive_pbh(pbehaviors):
    """
    Return a list without the inactive pbehavior of at the time of call.

    :param pbehavior: a list of pbehavior.
    :return list: a list without any inactive pbehavior
    """
    now = time.time()

    active_pbh = []
    for pbh in pbehaviors:
        if pbehavior_manager.check_active_pbehavior(now, pbh):
            active_pbh.append(pbh)

    return active_pbh


def _parse_direction(direction):
    """
    Parse the sort direction retrieved from the request.

    If direction is `ASC`, retrun 1. If direction is `DESC` return -1.
    If the value does not match `ASC` or `DESC` raise a ValueError exception.

    :param int direction: 1 or -1
    :return str: `ASC` or `DESC`
    """
    if direction == "ASC":
        return 1
    elif direction == "DESC":
        return -1
    else:
        raise ValueError("Direction must be 'ASC' or 'DESC' not {}.".format(
            direction))


def _generate_tile_pipeline(watcher_filter, limit, start, orderby, direction):
    """
    Return the aggregation pipeline use to retrieve every watcher, their
    alarm and pbehavior and their
    watched entities and their respective alarm and pbehavior.

    :param watcher_filter:
    :param limit: the number of watcher (tile) to return
    :param start: the number of watcher to skip
    :param orderby: the watcher field use the sort the result
    :param direction: the direction of the sort 'ASC' or 'DESC'
    :return list: return a list of mongodb aggregation stage
    """
    # Select the watchers
    select_watcher_stage = {"$match": watcher_filter}

    # Pagination
    skip = {"$skip": start}

    # Retrieve opened alarm for the watchers
    # I use the `$graphLookup` stage in order to retrieve only the opened alarms
    # with the `restrictSearchWithMatch` option.
    alarms = {"$graphLookup":
              {"from": "periodical_alarm",
               "startWith": "$_id",
               "connectFromField": "_id",
               "connectToField": "d",
               "restrictSearchWithMatch": {'v.resolved': None},
               "as": "alarm",
               "maxDepth": 0}}

    # Retrieve every pbehaviors on the watcher
    pbehaviors = {"$lookup":
                  {"from": "default_pbehavior",
                   "localField": "_id",
                   "foreignField": "eids",
                   "as": "pbehaviors"}}

    # Retrieve watched entities
    entities = {"$lookup":
                {"from": "default_entities",
                 "localField": "depends",
                 "foreignField": "_id",
                 "as": "watched_entities"}}

    # Retrieve every pbehaviors on the watched entities
    pbehaviors_watched_ent = {"$graphLookup":
                              {"from": "default_pbehavior",
                               "startWith": "$watched_entities._id",
                               "connectFromField": "watched_entities._id",
                               "connectToField": "eids",
                               "maxDepth": 0,
                               "as": "watched_entities_pbehaviors"}}

    # Retrieve every opened alarm on the watched entities
    # I use the `$graphLookup` stage in order to retrieve only the opened
    # alarms with the `restrictSearchWithMatch` option.
    alarm_watched_ent = {"$graphLookup":
                         {"from": "periodical_alarm",
                          "startWith": "$watched_entities._id",
                          "connectFromField": "_id",
                          "connectToField": "d",
                          "as": "watched_entities_alarm",
                          "restrictSearchWithMatch": {'v.resolved': None},
                          "maxDepth": 0}}

    # Hide watched_entities fields to fit document into 16M limit
    hide_fields = {"$project": {"watched_entities.impact": 0,
                                "watched_entities.depends": 0, "watched_entities.infos": 0}}

    # Genenate the pipeline
    pipeline = [select_watcher_stage,
                skip,
                alarms,
                pbehaviors,
                entities,
                pbehaviors_watched_ent,
                alarm_watched_ent,
                hide_fields]

    # Insert optionnal stage limit
    if limit is not None:
        pipeline.insert(2, {"$limit": limit})

    # Insert optionnal stage orderby
    if orderby is not None:
        direction = _parse_direction(direction)
        pipeline.insert(1, {"$sort": {orderby: direction}})

    return pipeline


def _rework_watcher_pipeline_element(watcher, logger):
    """Return a rearrange element from the watcher pipeline.

    This function will remove every inactive pbehaviors from the fields
    `ResultKey.PBEHAVIORS` and `ResultKey.WATCHED_ENT_PBH`.

    Then create a dict with every watched entities with the respective alarm
    and pbehaviors. It will be store under the ResultKey.ENT field.

    :param dict watcher: an element from the watcher pipeline
    :return dict: a rearrange watcher that respect the following pattern

    {
        "ResultKey.ID": str,
        "ResultKey.ALARM": list of alarm,
        "depends": list of str,
        "enabled_history": list of int,
        "enabled": boolean,
        "impact": list of str,
        "ResultKey.INFOS": dict,
        "ResultKey.LINKS": dict,
        "measurements": dict,
        "ResultKey.MFILTER": str,
        "ResultKey.NAME": str,
        "ResultKey.PBEHAVIORS": list of pbehavior,
        "ResultKey.STATE": int,
        "type": str,
        "ResultKey.ENT": list of entities, their respective pbehaviors
            and alarms
    }
    """
    # remove the inactive pbehaviors from the pipeline result
    pbhs = watcher[ResultKey.PBEHAVIORS]
    watcher[ResultKey.PBEHAVIORS] = _remove_inactive_pbh(pbhs)
    pbhs = watcher[ResultKey.WATCHED_ENT_PBH]
    watcher[ResultKey.WATCHED_ENT_PBH] = _remove_inactive_pbh(pbhs)

    # assign watched entities pbehaviors to the correct entities
    entities = {}
    for entity in watcher[ResultKey.ENT]:
        entity[ResultKey.PBEHAVIORS] = []
        entity[ResultKey.ALARM] = None
        entities[entity[ResultKey.ENT_ID]] = entity

    for pbh in watcher[ResultKey.WATCHED_ENT_PBH]:
        for ent_id in pbh["eids"]:
            try:
                entities[ent_id][ResultKey.PBEHAVIORS].append(pbh)
            except KeyError:
                logger.error("Can not find entities {} in the"
                             "pipeline result".format(ent_id))

    # assign watched entities alarms to the correct entities
    for alarm in watcher[ResultKey.WATCHED_ENT_ALRM]:
        try:
            entities[alarm["d"]][ResultKey.ALARM] = alarm
        except KeyError:
            logger.error("Can not find entities {} in the"
                         "pipeline result".format(alarm["d"]))

    watcher[ResultKey.ENT] = entities.values()
    del watcher[ResultKey.WATCHED_ENT_PBH]
    del watcher[ResultKey.WATCHED_ENT_ALRM]

    return watcher


def exports(ws):
    ws.application.router.add_filter('mongo_filter', mongo_filter)
    ws.application.router.add_filter('id_filter', id_filter)

    stat_manager = StatManager(*StatManager.provide_default_basics(ws.logger))

    @ws.application.route(
        '/api/v2/weather/watchers/<watcher_filter:mongo_filter>'
    )
    def get_watcher(watcher_filter):
        """
        Return a list of tile ready to be displayed by the front-end.

        For more informations, see the __TileData object.

        :param dict watcher_filter: a mongo filter to find watchers
        :rtype: list of __TileData as a JSON
        """
        limit = request.query.limit or None
        start = request.query.start or DEFAULT_START
        orderby = request.query.orderby or None
        direction = request.query.direction or None

        try:
            start = int(start)
        except ValueError:
            start = int(DEFAULT_START)
        if limit is not None:
            try:
                limit = int(limit)
            except ValueError:
                limit = DEFAULT_LIMIT

        wf = WatcherFilter()
        original_filter = copy.deepcopy(watcher_filter)
        watcher_filter['type'] = 'watcher'
        watcher_filter = wf.filter(watcher_filter, True)
        try:
            pipeline = _generate_tile_pipeline(watcher_filter,
                                               limit,
                                               start,
                                               orderby,
                                               direction)
        except ValueError as error:
            return gen_json_error({"name": "Can not parse sort direction.",
                                   "description": str(error)}, 400)

        try:
            pipeline_result = mongo_collection.aggregate(pipeline)
        except Exception as error:
            ws.logger.error('Watcher aggregation {} error {}'.format(pipeline, str(error)))
            return gen_json_error({"name": "Query error",
                                   "description": str(error)}, 500)

        result = []
        tile_enriched_watchers = []
        try:
            for watcher in pipeline_result:

                try:
                    watcher = _rework_watcher_pipeline_element(watcher, ws.logger)
                except Exception as error:
                    ws.logger.error('_rework_watcher_pipeline_element {} error {}'.format(watcher, str(error)))
                    return gen_json_error({"name": "Query error",
                                        "description": str(error)}, 500)

                # This part should not exist and must be considered deprecated.
                # This filter has to be done inside the aggregation pipeline but
                # currently it is impossible as there is no way to check if a
                # pbehavior is active directly inside the database.

                some_watched_ent_paused, all_watched_ent_paused = _watcher_status(
                    watcher
                )

                if wf.match(all_watched_ent_paused,
                            some_watched_ent_paused,
                            len(watcher[ResultKey.PBEHAVIORS]) > 0,
                            _pbehavior_types(watcher)):
                    tileData = __TileData(watcher)
                    tiled_watcher = tileData.watcher
                    tile_enriched_watchers.append(tiled_watcher)
                    tileData.remove_tiled_watcher()
                    tiled_watcher["weather_format"] = vars(tileData)

        except PyMongoError as error:
            ws.logger.warning('get_watcher {} {} {}'.format(pipeline, type(error).__name__, str(error)))
        except Exception as error:
            ws.logger.error('get_watcher iterate result {} {} {}'.format(pipeline, type(error).__name__, str(error)))
            return gen_json_error({"name": "Query error",
                                "description": str(error)}, 500)

        tile_filter = wf.filter(original_filter)
        mtl = MontyList(tile_enriched_watchers).find(tile_filter)
        mtl = [m["weather_format"] for m in mtl]
        return gen_json(list(mtl))

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

        active_pbehaviors = pbehavior_manager.get_active_pbehaviors_on_entities(
            entity_ids)
        for pbehavior in active_pbehaviors:
            pbehavior['isActive'] = True
            for eid in pbehavior.get('eids', []):
                if eid in entities:
                    entities[eid]['pbehaviors'].append(pbehavior)

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
            enriched_entity['stats'] = stat_manager.get_stats(entity_id).as_dict()
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
