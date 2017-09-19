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

from canopsis.common.utils import singleton_per_scope
from canopsis.task.core import register_task

from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader


@register_task
def event_processing(engine, event, alertsmgr=None, **kwargs):
    """
    AMQP Event processing.
    """
    if alertsmgr is None:
        alertsmgr = singleton_per_scope(Alerts)

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
    if alertsmgr is None:
        alertsmgr = singleton_per_scope(Alerts)

    alertsreader = singleton_per_scope(AlertsReader)

    alertsmgr.config = alertsmgr.load_config()

    unresolved_alerts = alertsmgr.get_alarms(resolved=False)

    alertsmgr.resolve_alarms(unresolved_alerts)

    alertsmgr.resolve_cancels(unresolved_alerts)

    alertsmgr.resolve_snoozes()

    alertsmgr.resolve_stealthy(unresolved_alerts)

    alertsmgr.check_alarm_filters()

    alertsreader.clean_fast_count_cache()
