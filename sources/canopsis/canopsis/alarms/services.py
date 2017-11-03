#!/usr/bin/env python
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

import logging

DEFAULT_FLAPPING_INTERVAL = 0
DEFAULT_CANCEL_AUTOSOLVE_DELAY = 3600
DEFAULT_STEALTHY_SHOW_DURATION = 0
DEFAULT_STEALTHY_INTERVAL = 0

class AlarmService(object):

    def __init__(self, alarms_adapter, entities_adapter, watcher_manager,
                 logger=None, config=None):
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
        if config is None:
            self.config = {}
        else:
            if not isinstance(config, dict):
                raise ValueError("config must be a dict")

            self.config = config

    def _log(self, level, message):
        """
        Logs 'message' according to 'level' if the logger is present. Does nothing otherwise.

        :param int level: a level from logging package
        :param string message: the log message
        """
        if self.logger is not None:
            self.logger.log(level, message)

    def update_alarm(self, alarm):
        self.alarms_adapter.update(alarm)
        self.watcher_manager.alarm_changed(alarm.identity.get_data_id())

    def find_snoozed_alarms(self):
        return self.alarms_adapter.find_unresolved_snoozed_alarms()

    def resolve_snoozed_alarms(self, alarms=None):
        if alarms is None:
            alarms = self.find_snoozed_alarms()
        for alarm in alarms:
            if alarm.resolve_snooze() is True:
                self._log(logging.DEBUG, "alarm : {} has been unsnoozed"
                                         .format(alarm._id))
                self.update_alarm(alarm)
                alarms.remove(alarm)

        return alarms

    def process_resolution_on_all_alarms(self):
        """
        This method processes all open alarms to check if they need to be resolved;

        This method is meant to be used in the Alarm Engine's beat processing.
        :return:
        """
        alarm_counter = 0
        updated_alarm_counter = 0
        for alarm in self.alarms_adapter.stream_unresolved_alarms():
            alarm_needs_update = False

            if alarm.resolve(self.config.get('bagot_time', DEFAULT_FLAPPING_INTERVAL)) is True:
                alarm_needs_update = True

            if alarm.resolve_cancel(self.config.get('cancel_autosolve_delay', DEFAULT_CANCEL_AUTOSOLVE_DELAY)) is True:
                alarm_needs_update = True

            stealthy_show_duration = self.config.get('stealthy_show', DEFAULT_STEALTHY_SHOW_DURATION)
            stealthy_interval = self.config.get('stealthy_time', DEFAULT_STEALTHY_INTERVAL)
            if alarm.resolve_stealthy(stealthy_show_duration, stealthy_interval) is True:
                alarm_needs_update = True

            alarm_counter += 1
            if alarm_needs_update is True:
                self.update_alarm(alarm)
                updated_alarm_counter += 1

        self._log(
            logging.DEBUG,
            "alarms resolution processing : {0} alarms processed. {1} updates. ".format(
                alarm_counter,
                updated_alarm_counter
            )
        )
