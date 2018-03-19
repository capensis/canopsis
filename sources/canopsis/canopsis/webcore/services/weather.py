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

from operator import itemgetter
from bottle import request

from canopsis.alerts.enums import AlarmField, AlarmFilterField
from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader
from canopsis.common.converters import mongo_filter, id_filter
from canopsis.common.utils import get_rrule_freq
from canopsis.stats.manager import Stats
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.tracer.manager import TracerManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_NOT_FOUND

alarm_manager = Alerts(*Alerts.provide_default_basics())
alarmreader_manager = AlertsReader(*AlertsReader.provide_default_basics())
context_manager = alarm_manager.context_manager
tracer_manager = TracerManager(*TracerManager.provide_default_basics())
pbehavior_manager = PBehaviorManager(*PBehaviorManager.provide_default_basics())
stats_manager = Stats()

DEFAULT_LIMIT = '120'
DEFAULT_START = '0'
DEFAULT_SORT = False
DEFAULT_PB_TYPES = []


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


class WatcherFilter(object):
    """
    Filter out active_pb_some and active_pb_all.
    """

    def __init__(self):
        self._all = None
        self._some = None
        self._watcher = None
        self._types = set()

    def all(self):
        """
        :returns: True, False or None
        """
        return self._all

    def some(self):
        """
        :returns: True, False or None
        """
        return self._some

    def watcher(self):
        """
        :returns: True, False or None if the watcher has an active pbehavior.
        """
        return self._watcher

    def types(self):
        """
        :returns: set of pbehavior types to filter on.
        :rtype: set[str]
        """
        return self._types

    def to_bool(self, value):
        if value in ['1', 1, "true", "True"]:
            return True
        return False

    def _filter_dict(self, dictdoc):
        cdoc = copy.deepcopy(dictdoc)
        for k, v in dictdoc.items():
            if k == 'active_pb_all':
                self._all = self.to_bool(v)
                del cdoc['active_pb_all']

            elif k == 'active_pb_some':
                self._some = self.to_bool(v)
                del cdoc['active_pb_some']

            elif k == 'active_pb_type':
                self._types.add(str(v).strip().lower())
                del cdoc['active_pb_type']

            elif k == 'active_pb_watcher':
                self._watcher = self.to_bool(v)
                del cdoc['active_pb_watcher']

            else:
                nv = self._filter(v)
                if nv is not None or v is None:
                    cdoc[k] = nv
                else:
                    del cdoc[k]

        return cdoc

    def _filter_list(self, listdoc):
        cdoc = copy.deepcopy(listdoc)
        j = 0
        for i, item in enumerate(listdoc):
            v = self._filter(item)
            if v is not None:
                cdoc[j] = v
            else:
                del cdoc[j]
                j -= 1
            j += 1

        return cdoc

    def _filter(self, doc):
        if isinstance(doc, dict):
            ndoc = self._filter_dict(doc)
            if len(ndoc) == 0 and len(doc) != 0:
                return None
            return ndoc

        elif isinstance(doc, list):
            ndoc = self._filter_list(doc)
            if len(ndoc) == 0 and len(doc) != 0:
                return None
            return ndoc

        return doc

    def _filter_pb_types(self, pb_types):
        # set flag to true if we are ok to take any pbehavior type
        #   no pb_types   -> True
        #   some pb_types -> False
        if len(pb_types) == 0 or len(self.types()) == 0:
            return True

        for pb_type in pb_types:
            if pb_type.strip().lower() in self.types():
                return True

        return False

    def filter(self, doc):
        res = self._filter(doc)
        if res is None:
            return {}
        return res

    def appendable(self, allstatus, somestatus, watcherstatus, pb_types=None):
        """
        Call WatcherFilter.filter(filter_) first before calling this function.

        :param bool allstatus: watcher has all entities with an active pbehavior
        :param bool somestatus: watcher has some or all entities with an active pbehavior
        :param bool watcherstatus: watcher has a pbehavior or not. If the filter contains a filter on that, it will be
        :param list[str] pb_types: list of pbehavior types to filter on. If None or empty, any types will be ok.
        """
        if not isinstance(allstatus, bool):
            raise ValueError('wrong allstatus value: not a bool')
        if not isinstance(somestatus, bool):
            raise ValueError('wrong somestatus value: not a bool')
        if not isinstance(watcherstatus, bool):
            raise ValueError('wrong watcherstatus value: not a bool')

        if pb_types is None:
            pb_types = list()

        if not (isinstance(pb_types, list) or isinstance(pb_types, set)):
            raise ValueError('wrong pb_types value: not a list')

        logic_some_all = False
        logic_watcher = False

        if self.watcher() is None:
            logic_watcher = True

        elif self.watcher() is True and watcherstatus is True:
            logic_watcher = True

        elif self.watcher() is False and watcherstatus is False:
            logic_watcher = True

        if self.all() is None and self.some() is None:
            logic_some_all = True

        elif self.all() is not None and self.some() is not None:
            # cannot be simplified as only this test.
            if self.all() == allstatus and self.some() == somestatus:
                logic_some_all = True

        elif self.all() is allstatus and self.all() is not None:
            logic_some_all = True

        elif self.some() is somestatus and self.some() is not None:
            logic_some_all = True

        return logic_watcher and logic_some_all and self._filter_pb_types(pb_types)


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


def check_baseline(merged_eids_tracer, watcher_depends):
    """
    check if the watcher has an entity with a baseline active

    :param set merged_eids_tracer: all entities withan active baseline
    :param list watcher_depends: watcher entities
    :returns: true if the watcher has an entity with an active active_baseline
    :rtype: bool
    """
    for entity_id in watcher_depends:
        if entity_id in merged_eids_tracer:
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
        limit = request.query.limit or DEFAULT_LIMIT
        start = request.query.start or DEFAULT_START
        sort = request.query.sort or DEFAULT_SORT

        # FIXIT: service weather has no pagination capability at all.
        try:
            #start = int(start)
            start = 0
        except ValueError:
            start = int(DEFAULT_START)
        try:
            #limit = int(limit)
            limit = None
        except ValueError:
            limit = int(DEFAULT_LIMIT)

        wf = WatcherFilter()
        watcher_filter['type'] = 'watcher'
        watcher_filter = wf.filter(watcher_filter)
        pb_types = wf.types()

        watcher_list = context_manager.get_entities(
            query=watcher_filter,
            limit=limit,
            start=start,
            sort=sort
        )

        depends_merged = set([])
        active_pb_dict = {}
        active_pb_dict_full = {}
        alarm_watchers_ids = []
        entity_watchers_ids = []
        alarm_dict = {}
        merged_pbehaviors_eids = set([])
        next_run_dict = {}
        watchers = []
        merged_eids_tracer = []

        # Find all entites with an activated baseline
        active_baseline_tracer = tracer_manager.get(
            {
                'triggered_by': 'baseline',
                'extra.active': True
            }
        )
        for tracer in active_baseline_tracer:
            merged_eids_tracer = merged_eids_tracer + tracer['impact_entities']
        merged_eids_tracer = set(merged_eids_tracer)

        # List all activated pbh eids, ordered by pbh id
        actives_pb = pbehavior_manager.get_all_active_pbehaviors()

        for pbh in actives_pb:
            active_pb_dict[pbh['_id']] = set(pbh.get('eids', []))
            active_pb_dict_full[pbh['_id']] = pbh

        # List all watcher ids on entities and alarms
        for watcher in watcher_list:
            for depends_id in watcher['depends']:
                depends_merged.add(depends_id)
            entity_watchers_ids.append(watcher['_id'])
            alarm_watchers_ids.append(watcher['_id'])

        active_pbehaviors, active_watchers_pbehaviors = get_active_pbehaviors_on_watchers(
            watcher_list,
            active_pb_dict,
            active_pb_dict_full
        )
        # List all actived pbh eids
        for eids_tab in active_pb_dict.values():
            for eid in eids_tab:
                merged_pbehaviors_eids.add(eid)

        # List alarm values has a dict
        alarm_list = alarmreader_manager.get(
            filter_={'v.resolved': None}
        )['alarms']

        for alarm in alarm_list:
            alarm_dict[alarm['d']] = alarm['v']

        # List all next_run timers, grouped by alarm
        alerts_list_on_depends = alarmreader_manager.get(
            filter_={'d': {'$in': list(depends_merged)}}
        )['alarms']
        for alert in alerts_list_on_depends:
            if 'alarmfilter' in alert['v']:
                alarmfilter = alert['v']['alarmfilter']
                if isinstance(alarmfilter, dict) and "next_run" in alarmfilter:
                    next_run_dict[alert['d']] = alarmfilter['next_run']

        for watcher in watcher_list:
            enriched_entity = {}
            tmp_alarm = alarm_dict.get(
                '{}'.format(watcher['_id']),
                []
            )
            tmp_linklist = []
            for k, val in watcher['links'].items():
                tmp_linklist.append({'cat_name': k, 'links': val})

            enriched_entity['entity_id'] = watcher['_id']
            enriched_entity['infos'] = watcher['infos']
            enriched_entity['criticity'] = watcher['infos'].get('criticity', '')
            enriched_entity['org'] = watcher['infos'].get('org', '')
            enriched_entity['sla_text'] = ''  # when sla
            enriched_entity['display_name'] = watcher['name']
            enriched_entity['linklist'] = tmp_linklist
            enriched_entity['state'] = {'val': watcher.get('state', 0)}

            if tmp_alarm != []:
                enriched_entity['state'] = tmp_alarm['state']
                enriched_entity['status'] = tmp_alarm['status']
                enriched_entity['snooze'] = tmp_alarm.get('snooze')
                enriched_entity['ack'] = tmp_alarm.get('ack')
                enriched_entity['connector'] = tmp_alarm['connector']
                enriched_entity['connector_name'] = (
                    tmp_alarm['connector_name']
                )
                enriched_entity['last_update_date'] = tmp_alarm.get(
                    'last_update_date', None
                )
                enriched_entity['component'] = tmp_alarm['component']
                if 'resource' in tmp_alarm.keys():
                    enriched_entity['resource'] = tmp_alarm['resource']

            enriched_entity['pbehavior'] = active_pbehaviors.get(watcher['_id'], [])
            enriched_entity['watcher_pbehavior'] = active_watchers_pbehaviors.get(watcher['_id'], [])
            enriched_entity["mfilter"] = watcher["mfilter"]
            enriched_entity['alerts_not_ack'] = alert_not_ack_in_watcher(
                watcher['depends'],
                alarm_dict
            )
            wstatus = watcher_status(watcher, merged_pbehaviors_eids)
            enriched_entity["active_pb_some"] = wstatus[0]
            enriched_entity["active_pb_all"] = wstatus[1]
            enriched_entity['active_pb_watcher'] = len(enriched_entity['watcher_pbehavior']) > 0
            enriched_entity['has_baseline'] = check_baseline(
                merged_eids_tracer,
                watcher['depends']
            )
            tmp_next_run = get_next_run_alert(
                watcher.get('depends', []),
                next_run_dict
            )
            if tmp_next_run is not None:
                enriched_entity['automatic_action_timer'] = tmp_next_run

            watcher_pb_types = pbehavior_types(enriched_entity['pbehavior'])
            watcher_pb_types |= pbehavior_types(enriched_entity['watcher_pbehavior'])
            if wf.appendable(wstatus[1], wstatus[0], enriched_entity['active_pb_watcher'], pb_types=watcher_pb_types) is True:
                watchers.append(enriched_entity)

        watchers = sorted(watchers, key=itemgetter("display_name"))

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
        except (ValueError, KeyError, TypeError):
            json_error = {
                "name": "filter_not_found",
                "description": "impossible to load the desired filter"
            }
            return gen_json_error(json_error, HTTP_NOT_FOUND)

        query["enabled"] = True

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

            tmp_linklist = []
            for k, val in raw_entity['links'].items():
                tmp_linklist.append({'cat_name': k, 'links': val})

            enriched_entity['pbehavior'] = entity['pbehaviors']
            enriched_entity['entity_id'] = entity_id
            enriched_entity['linklist'] = tmp_linklist
            enriched_entity['infos'] = raw_entity['infos']
            enriched_entity['sla_text'] = ''  # TODO when sla, use it
            enriched_entity['org'] = raw_entity['infos'].get('org', '')
            enriched_entity['name'] = raw_entity['name']
            enriched_entity['source_type'] = raw_entity['type']
            enriched_entity['state'] = {'val': 0}
            # enriched_entity['stats'] = stats_manager.get_ok_ko(entity_id)
            # FIXME: canopsis/canopsis#650 UTF-8 error
            if current_alarm is not None:
                enriched_entity['ticket'] = current_alarm['ticket']
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
                if 'resource' in current_alarm:
                    enriched_entity['resource'] = current_alarm['resource']

            enriched_entities.append(enriched_entity)

        enriched_entities = sorted(enriched_entities, key=itemgetter("name"))

        return gen_json(enriched_entities)
