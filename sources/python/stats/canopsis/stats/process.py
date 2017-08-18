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

from copy import copy

from canopsis.common.utils import singleton_per_scope
from canopsis.task.core import register_task

from canopsis.stats.producers.user import UserMetricProducer
from canopsis.stats.producers.event import EventMetricProducer

from canopsis.session.manager import Session

from canopsis.alerts.status import get_previous_step
from canopsis.alerts.manager import Alerts


def session_stats(usermgr, sessionmgr, logger):
    for expired in sessionmgr.sessions_close():
        duration = expired['session_stop'] - expired['session_start']

        usermgr.session_duration(
            user=expired['_id'],
            duration=duration
        )


def opened_alarm_stats(eventmgr, alertsmgr, storage, logger):
    for resolved in [False, True]:
        alarms = alertsmgr.get_alarms(
            resolved=resolved,
            exclude_tags='stats-opened'
        )

        for entity_id in alarms:
            for docalarm in alarms[entity_id]:

                docalarm[storage.DATA_ID] = entity_id

                alarm = docalarm[storage.VALUE]
                extra = copy(alarm['extra'])

                eventmgr.alarm_opened(
                    extra_fields=extra
                )

                alertsmgr.update_current_alarm(
                    docalarm,
                    alarm,
                    tags='stats-opened'
                )


def resolved_alarm_stats(eventmgr, usermgr, alertsmgr, storage, logger):
    resolved_alarms = alertsmgr.get_alarms(
        resolved=True,
        exclude_tags='stats-resolved'
    )

    for entity_id in resolved_alarms:
        for docalarm in resolved_alarms[entity_id]:
            docalarm[storage.DATA_ID] = entity_id
            alarm = docalarm[storage.VALUE]
            alarm_ts = docalarm[storage.TIMESTAMP]

            extra = copy(alarm['extra'])

            solved_delay = alarm['resolved'] - alarm_ts
            eventmgr.alarm_solved_delay(
                delay=solved_delay,
                extra_fields=extra
            )

            if alarm['ack'] is not None:
                ack_ts = alarm['ack']['t']

                eventmgr.alarm_ack_solved_delay(
                    delay=alarm['resolved'] - ack_ts,
                    extra_fields=extra
                )

            alarm_events = alertsmgr.get_events(docalarm)
            for event in alarm_events:
                if event['event_type'] == 'ack':
                    ack_ts = event['timestamp']

                    ackremove = get_previous_step(
                        alarm,
                        'ackremove',
                        ts=ack_ts
                    )

                    ref_ts = alarm_ts if ackremove is None else ackremove['t']
                    ack_delay = ack_ts - ref_ts

                    usermgr.alarm_ack_delay(
                        user=event['author'],
                        delay=ack_delay,
                        extra_fields=extra
                    )

            alertsmgr.update_current_alarm(
                docalarm,
                alarm,
                tags='stats-resolved'
            )


@register_task
def beat_processing(
    engine,
    sessionmgr=None,
    eventmgr=None,
    usermgr=None,
    alertsmgr=None,
    logger=None,
    **kwargs
):
    if sessionmgr is None:
        sessionmgr = singleton_per_scope(Session)

    if eventmgr is None:
        eventmgr = singleton_per_scope(EventMetricProducer)

    if usermgr is None:
        usermgr = singleton_per_scope(UserMetricProducer)

    if alertsmgr is None:
        alertsmgr = singleton_per_scope(Alerts)

    storage = alertsmgr[alertsmgr.ALARM_STORAGE]

    session_stats(usermgr, sessionmgr, logger)

    with engine.Lock(engine, 'alarm_stats_computation') as l:
        if l.own():
            opened_alarm_stats(
                eventmgr,
                alertsmgr,
                storage,
                logger
            )

            resolved_alarm_stats(
                eventmgr,
                usermgr,
                alertsmgr,
                storage,
                logger
            )
