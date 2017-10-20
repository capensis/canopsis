#!/usr/bin/env python
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



class AlarmService(object):

    def __init__(self, alarms_adapter, entities_adapter, logger=None):
        """
        Alarm service constructor.

        :param canopsis.alarms.adapters.Adapter alarms_adapter: Alarms DB adapter
        :param canopsis.entities.adapters.Adapter entities_adapter: Context-graph entities DB adapter
        :param canopsis.logger.Logger: a logger instance
        """
        self.alarms_adapter = alarms_adapter
        self.entities_adapter = entities_adapter

    def find_active_alarms(self):
        """
        finds all active alarms and matches them with their owner entity
        :return: dict
        """
        alarms = self.alarms_adapter.find_unresolved_alarms()
        entities = self.entities_adapter.find_all_enabled()
        alarms_sorted_by_entity = self._build_old_school_format(alarms)
        alarms_with_embedded_entities = self._match_alarms_with_entities(alarms_sorted_by_entity, entities)
        return alarms_with_embedded_entities

    def find_snoozed_alarms(self, use_old__school_format=True):
        alarms = self.alarms_adapter.find_unresolved_snoozed_alarms()
        entities = self.entities_adapter.find_all_enabled()
        if use_old__school_format is True:
            alarms_sorted_by_entity = self._build_old_school_format(alarms)
        else:
            alarms_sorted_by_entity = alarms
        alarms_with_embedded_entities = self._match_alarms_with_entities(alarms_sorted_by_entity, entities, use_old__school_format)
        return alarms_with_embedded_entities



    def resolve_snoozed_alarms(self):
        alarms = self.find_snoozed_alarms(False)
        for alarm in alarms:
            if alarm.resolve_snooze() is True:
                self.alarms_adapter.update(alarm)

    def _build_old_school_format(self, alarms_list):
        """
            Hack to make an alarm list compatible with the current storage system.
            :param []Alarm alarms_list:
            :return dict: a dict fomated agains mongo.periodicalstorage._cursor2periods
        """
        completed_alarms_list = {}
        for a in alarms_list:
            alarm_dict = a.to_dict()
            old_school_alarm = {
                'timestamp': a.creation_date,
                'value': alarm_dict['v']
            }
            entity_id = alarm_dict['d']
            if entity_id not in completed_alarms_list:
                completed_alarms_list[entity_id] = [old_school_alarm]
            else:
                completed_alarms_list[entity_id].append(old_school_alarm)

        return completed_alarms_list


    def _match_alarms_with_entities(self, alarms_list, entities_list, use_old_school_dict=True):

        for entity in entities_list:
            if entity.id_ in alarms_list:
                for entity_id, alarm in alarms_list[entity.id_]:
                    if use_old_school_dict is True:
                        alarm['entity'] = entity.to_dict()
                    else:
                        alarm.entity = entity

        return alarms_list


