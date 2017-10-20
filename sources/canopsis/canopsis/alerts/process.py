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

from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader
from canopsis.task.core import register_task
from canopsis.alarms.services import AlarmService
from canopsis.alarms.adapters import Adapter as AlarmAdapter
from canopsis.entities.adapters import Adapter as EntityAdapter
from canopsis.logger import Logger
import time

alerts_manager = Alerts(*Alerts.provide_default_basics())
alertsreader_manager = AlertsReader(*AlertsReader.provide_default_basics())

alarms_service = AlarmService(AlarmAdapter(alerts_manager.alerts_storage._backend.database), EntityAdapter(alerts_manager.alerts_storage._backend.database))

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
    logger = Logger.get('alarms', '/opt/canopsis/var/log/alarms.log')

    logger.critical("Starting beat processing.")
    start_time = int(round(time.time() * 1000))
    if alertsmgr is None:
        alertsmgr = alerts_manager

    alertsreader = alertsreader_manager




    #unresolved_alarms = alertsmgr.get_alarms(resolved=False)
    read_time_start = int(round(time.time() * 1000))
    unresolved_alarms = alarms_service.find_active_alarms()
    read_end_time = int(round(time.time() * 1000))

    alertsmgr.logger.critical("DB read time : {} ms".format(read_end_time - read_time_start))

    resolve_start_time = int(round(time.time() * 1000))
    unresolved_alarms = alertsmgr.resolve_alarms(unresolved_alarms)
    resolve_end_time = int(round(time.time() * 1000))

    alertsmgr.logger.critical("DB resolve time : {} ms".format(resolve_end_time - resolve_start_time))

    cancel_start_time = int(round(time.time() * 1000))
    unresolved_alarms = alertsmgr.resolve_cancels(unresolved_alarms)
    cancel_end_time = int(round(time.time() * 1000))
    alertsmgr.logger.critical("DB cancel time : {} ms".format(cancel_end_time - cancel_start_time))

    snooze_start_time = int(round(time.time() * 1000))
    # TODO : vérifier que l'alarme est présente ici et se snooze bien
    snoozed_alarms = alarms_service.find_snoozed_alarms()
    alertsmgr.logger.critical('Found {} alarms'.format(len(snoozed_alarms)))
    alarms_service.resolve_snoozed_alarms()
    #alertsmgr.resolve_snoozes(snoozed_alarms)
    snooze_end_time = int(round(time.time() * 1000))

    alertsmgr.logger.critical("snooze time : {} ms".format(snooze_end_time- snooze_start_time))

    stealthy_start_time = int(round(time.time() * 1000))
    unresolved_alarms = alertsmgr.resolve_stealthy(unresolved_alarms)
    stealthy_end_time = int(round(time.time() * 1000))
    alertsmgr.logger.critical("stealthy time : {} ms".format(stealthy_end_time- stealthy_start_time))

    # unresolved_alarms not used actually but can be used for new actions
    check_start_time = int(round(time.time() * 1000))
    alertsmgr.check_alarm_filters()
    check_end_time = int(round(time.time() * 1000))
    alertsmgr.logger.critical("Check time : {} ms".format(check_end_time- check_start_time))


    cache_start_time = int(round(time.time() * 1000))
    alertsreader.clean_fast_count_cache()
    cache_end_time = int(round(time.time() * 1000))
    alertsmgr.logger.critical("Cache time : {} ms".format(cache_end_time- cache_start_time))


    end_time = int(round(time.time() * 1000))
    logger.critical("End beat processing. Took : {} ms.".format(end_time - start_time))
