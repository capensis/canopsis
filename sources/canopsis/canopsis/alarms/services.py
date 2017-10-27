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

import logging

class AlarmService(object):

    def __init__(self, alarms_adapter, entities_adapter, watcher_manager, logger=None):
        """
        Alarm service constructor.

        :param canopsis.alarms.adapters.Adapter alarms_adapter: Alarms DB adapter
        :param canopsis.entities.adapters.Adapter entities_adapter: Context-graph entities DB adapter
        :param canopsis.logger.Logger: a logger instance
        """
        self.alarms_adapter = alarms_adapter
        self.entities_adapter = entities_adapter
        self.watcher_manager = watcher_manager
        self.logger = logger

    def _log(self, level, message):
        """
        Logs 'message' according to 'level' if the logger is present. Does nothing otherwise.
        :param int level: a level from logging package
        :param string message: the log message
        :return: None
        """
        if self.logger is not None:
            self.logger.log(level, message)

    def update_alarm(self, alarm):
        self.alarms_adapter.update(alarm)
        self.watcher_manager.alarm_changed(alarm.identity.get_data_id())

    def find_active_alarms(self):
        """
        finds all active alarms and matches them with their owner entity
        :return: dict
        """
        alarms = self.alarms_adapter.find_unresolved_alarms()
        self._log(logging.DEBUG, 'found {} active alarms'.format(len(alarms)))
        return alarms
        #entities = self.entities_adapter.find_all_enabled()
        #self._log(logging.DEBUG, 'found {} enabled entities'.format(len(entities)))
        #alarms_with_embedded_entities = self._match_alarms_with_entities(alarms, entities)
        #return alarms_with_embedded_entities

    def find_snoozed_alarms(self):
        alarms = self.alarms_adapter.find_unresolved_snoozed_alarms()
        return alarms
        #entities = self.entities_adapter.find_all_enabled()
        #alarms_with_embedded_entities = self._match_alarms_with_entities(alarms, entities)
        #return alarms_with_embedded_entities

    def resolve_snoozed_alarms(self, alarms=None):
        if alarms is None:
            alarms = self.find_snoozed_alarms(False)
        for alarm in alarms:
            if alarm.resolve_snooze() is True:
                self._log(logging.DEBUG, "alarm : {} has been unsnoozed".format(alarm._id))
                self.update_alarm(alarm)
                alarms.remove(alarm)

        return alarms



    def resolved_canceled_alarms(self, alarms, cancel_delay=3600):
        for alarm in alarms:
            if alarm.resolve_cancel(cancel_delay) is True:
                self._log(logging.DEBUG, "alarm : {0} was cancelled on {1} and will now be resolved".format(alarm._id, alarm.canceled.timestamp))
                self.update_alarm(alarm)
                alarms.remove(alarm)
        return alarms

    def resolve_alarms(self, alarms, flapping_interval=60):
        for alarm in alarms:
            if alarm.resolve(flapping_interval) is True:
                self._log(logging.DEBUG, "alarm : {} has been resolved and will now be resolved".format(alarm._id))
                self.update_alarm(alarm)
                alarms.remove(alarm)
        return alarms

    def resolve_stealthy_alarms(self, alarms, stealthy_show_duration=120, stealthy_interval=0):
        for alarm in alarms:
            if alarm.resolve_stealthy(stealthy_show_duration, stealthy_interval) is True:
                self._log(logging.DEBUG, "alarm : {} is not stealthy anymore and will now be resolved".format(alarm._id))
                self.update_alarm(alarm)
                alarms.remove(alarm)
        return alarms

    def _match_alarms_with_entities(self, alarms_list, entities_list):


        for entity in entities_list:
            if entity.id_ in alarms_list:
                for alarm in alarms_list[entity.id_]:
                    alarm.entity = entity
        return alarms_list


