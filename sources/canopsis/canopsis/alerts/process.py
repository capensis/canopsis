# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

from canopsis.alarms.services import AlarmService
from canopsis.alarms.adapters import AlarmAdapter
from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader
from canopsis.task.core import register_task
from canopsis.watcher.manager import Watcher

alerts_manager = Alerts(*Alerts.provide_default_basics())
alertsreader_manager = AlertsReader(*AlertsReader.provide_default_basics())

mongo_client = alerts_manager.alerts_storage._backend.database


@register_task
def event_processing(engine, event, alertsmgr=None, **kwargs):
    """
    AMQP Event processing.
    """
    if alertsmgr is None:
        alertsmgr = alerts_manager

    encoded_event = {}

    for key, value in event.items():
        try:
            key = key.encode('utf-8')
        except (UnicodeDecodeError, UnicodeEncodeError):
            pass

        try:
            value = value.encode('utf-8')
        except (UnicodeDecodeError, UnicodeEncodeError, AttributeError):
            pass

        encoded_event[key] = value

    try:
        alertsmgr.archive(encoded_event)
    except ValueError as ex:
        engine.logger.error('cannot store event: {}'.format(ex))


@register_task
def beat_processing(engine, alertsmgr=None, **kwargs):
    """
    Scheduled process.
    """
    alarms_service = AlarmService(
        alarms_adapter=AlarmAdapter(mongo_client),
        watcher_manager=Watcher()
    )

    if alertsmgr is None:
        alertsmgr = alerts_manager

    alertsreader = alertsreader_manager

    # process snoozed alarms first
    snoozed_alarms = alarms_service.find_snoozed_alarms()
    alarms_service.resolve_snoozed_alarms(snoozed_alarms)

    # process all resolution checks on all alarms.
    alarms_service.process_resolution_on_all_alarms()

    alertsmgr.check_alarm_filters()
    alertsreader.clean_fast_count_cache()
