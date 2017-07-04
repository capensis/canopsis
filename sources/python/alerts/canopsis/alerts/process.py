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
def event_processing(engine, event, alertsmgr=None, logger=None, **kwargs):
    if alertsmgr is None:
        alertsmgr = singleton_per_scope(Alerts)

    encoded_event = {}

    for k, v in event.items():
        try:
            k = k.encode('utf-8')
        except:
            pass
        try:
            v = v.encode('utf-8')
        except:
            pass
        encoded_event[k] = v

    alertsmgr.archive(encoded_event)


@register_task
def beat_processing(engine, alertsmgr=None, logger=None, **kwargs):
    if alertsmgr is None:
        alertsmgr = singleton_per_scope(Alerts)

    alertsreader = singleton_per_scope(AlertsReader)

    alertsmgr.config = alertsmgr.load_config()

    alertsmgr.resolve_alarms()

    alertsmgr.resolve_cancels()

    alertsmgr.resolve_snoozes()

    alertsmgr.resolve_stealthy()

    alertsmgr.check_alarm_filters()

    alertsreader.clean_fast_count_cache()
