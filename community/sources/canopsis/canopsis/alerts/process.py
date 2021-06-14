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

from canopsis.alarms.adapters import AlarmAdapter
from canopsis.alarms.services import AlarmService
from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader
from canopsis.task.core import register_task
from canopsis.watcher.manager import Watcher
from canopsis.common.mongo_store import MongoStore

# import this module so tasks are registered in canopsis.tasks.core
import canopsis.alerts.tasks as __alerts_tasks

work_alerts_manager = Alerts(*Alerts.provide_default_basics())
beat_alerts_manager = Alerts(*Alerts.provide_default_basics())
alertsreader_manager = AlertsReader(*AlertsReader.provide_default_basics())

mongo_store = MongoStore.get_default()

@register_task
def event_processing(engine, event, alertsmgr=None, **kwargs):
    """
    AMQP Event processing.
    """
    if alertsmgr is None:
        alertsmgr = work_alerts_manager

    encoded_event = {}

    for key, value in event.items():
        try:
            key = key.encode('utf-8')
        except UnicodeError:
            pass

        try:
            value = value.encode('utf-8')
        except (UnicodeError, TypeError, AttributeError):
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
    if alertsmgr is None:
        alertsmgr = beat_alerts_manager

    alarms_service = AlarmService(
        alarms_adapter=AlarmAdapter(mongo_store),
        context_manager=alertsmgr.context_manager,
        event_publisher=alertsmgr.event_publisher,
        watcher_manager=Watcher(),
        bagot_time=alertsmgr.flapping_interval,
        cancel_autosolve_delay=alertsmgr.cancel_autosolve_delay,
        done_autosolve_delay=alertsmgr.done_autosolve_delay,
        stealthy_interval=alertsmgr.stealthy_interval,
        logger=alertsmgr.logger
    )

    # process snoozed alarms first
    snoozed_alarms = alarms_service.find_snoozed_alarms()
    alarms_service.resolve_snoozed_alarms(snoozed_alarms)

    # process all resolution checks on all alarms.
    alarms_service.process_resolution_on_all_alarms()

    alertsmgr.check_alarm_filters()

    # Recompute watcher states
    alertsmgr.watcher_manager.compute_watchers()
