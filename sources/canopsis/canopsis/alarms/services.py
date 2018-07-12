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

"""
Service for Alarm.
"""

from __future__ import unicode_literals

import logging

from canopsis.statsng.enums import StatCounters, StatDurations

DEFAULT_CANCEL_AUTOSOLVE_DELAY = 3600
DEFAULT_FLAPPING_INTERVAL = 0
DEFAULT_STEALTHY_SHOW_DURATION = 0
DEFAULT_STEALTHY_INTERVAL = 0


class AlarmService(object):

    """
    Alarm Service class.
    """

    def __init__(self,
                 alarms_adapter,
                 context_manager,
                 event_publisher,
                 watcher_manager,
                 bagot_time=DEFAULT_FLAPPING_INTERVAL,
                 cancel_autosolve_delay=DEFAULT_CANCEL_AUTOSOLVE_DELAY,
                 stealthy_duration=DEFAULT_STEALTHY_SHOW_DURATION,
                 stealthy_interval=DEFAULT_STEALTHY_INTERVAL,
                 logger=None):
        """
        Alarm service constructor.

        :param AlarmAdapter alarms_adapter: Alarms DB adapter
        :param ContextGraph context_manager: Context graph
        :param StatEventPublisher event_publisher: Event publisher
        :param WatcherManager watcher_manager: ref to a WatcherManager object
        :param int bagot_time: period to consider status oscilations
        :param int cancel_autosolve_delay: delay before validating a cancel
        :param int stealthy_duration: period to consider an alarm as stealthy
        :param int stealthy_interval: period to show an alarm as stealthy
        :param Logger logger: a logger instance
        """
        self.alarms_adapter = alarms_adapter
        self.context_manager = context_manager
        self.event_publisher = event_publisher
        self.watcher_manager = watcher_manager
        self.bagot_time = bagot_time
        self.cancel_delay = cancel_autosolve_delay
        self.stealthy_duration = stealthy_duration
        self.stealthy_interval = stealthy_interval
        self.logger = logger

    def _log(self, level, message):
        """
        Logs 'message' according to 'level' if the logger is present.
        Does nothing otherwise.

        :param int level: a level from logging package
        :param string message: the log message
        """
        if self.logger is not None:
            self.logger.log(level, message)

    def update_alarm(self, alarm):
        """
        Update an alarm.

        :param Alarm alarm: an alarm
        """
        self.alarms_adapter.update(alarm)
        self.watcher_manager.alarm_changed(alarm.identity.get_data_id())

    def find_snoozed_alarms(self):
        """
        List all snoozed alarms.

        :rtype: [Alarm]
        """
        return self.alarms_adapter.find_unresolved_snoozed_alarms()

    def resolve_snoozed_alarms(self, alarms=None):
        """
        Resolve snooze on a list of alarms.

        :param [Alarm] alarms: a list of Alarm
        :rtype: [Alarm]
        """
        if alarms is None:
            alarms = self.find_snoozed_alarms()
        for alarm in alarms:
            if alarm.resolve_snooze():
                self._log(logging.DEBUG,
                          'alarm : {} has been unsnoozed'.format(alarm._id))
                self.update_alarm(alarm)
                alarms.remove(alarm)

        return alarms

    def process_resolution_on_all_alarms(self):
        """
        This method processes all open alarms to check if they need to be
        resolved;

        This method is meant to be used in the Alarm Engine's beat processing.
        """
        alarm_counter = 0
        updated_alarm_counter = 0

        for alarm in self.alarms_adapter.stream_unresolved_alarms():
            alarm_counter += 1

            resolved = alarm.resolve(self.bagot_time)
            resolved_flapping = alarm.resolve_flapping(self.bagot_time)
            resolved_cancel = alarm.resolve_cancel(self.cancel_delay)
            resolved_stealthy = alarm.resolve_stealthy(self.stealthy_interval)
            if resolved or resolved_cancel or resolved_stealthy or resolved_flapping:
                self.update_alarm(alarm)
                updated_alarm_counter += 1

                entity_id = alarm.identity.get_data_id()
                entity = self.context_manager.get_entities_by_id(entity_id)
                try:
                    entity = entity[0]
                except IndexError:
                    entity = {}
                alarm_dict = alarm.to_dict()
                self.event_publisher.publish_statcounterinc_event(
                    alarm.last_update_date,
                    StatCounters.alarms_resolved,
                    entity,
                    alarm_dict['v'])
                self.event_publisher.publish_statduration_event(
                    alarm.last_update_date,
                    StatDurations.resolve_time,
                    alarm.last_update_date - alarm.creation_date,
                    entity,
                    alarm_dict['v'])

        self._log(
            logging.DEBUG,
            "alarms resolution processing : {} alarms processed. {} updates."
            .format(alarm_counter, updated_alarm_counter)
        )
