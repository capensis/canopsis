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

from canopsis.common.utils import singleton_per_scope
from canopsis.task.core import register_task
from canopsis.engines.core import publish

from canopsis.stats.producers.user import UserMetricProducer
from canopsis.stats.producers.event import EventMetricProducer

from canopsis.session.manager import Session

from canopsis.alerts.status import get_previous_step
from canopsis.alerts.manager import Alerts


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
    events = sessionmgr.duration()

    with engine.Lock(engine, 'alarm_stats_computation') as l:
        if l.own():
            resolved_alarms = alertsmgr.get_alarms(
                resolved=True,
                exclude_tags='stats'
            )

            for data_id in resolved_alarms:
                for docalarm in resolved_alarms[data_id]:
                    docalarm[storage.DATA_ID] = data_id
                    alarm = docalarm[storage.VALUE]
                    alarm_ts = docalarm[storage.TIMESTAMP]
                    alarm_events = alertsmgr.get_events(docalarm)

                    solved_delay = alarm['resolved'] - alarm_ts
                    events.append(eventmgr.alarm_solved_delay(solved_delay))

                    if alarm['ack'] is not None:
                        ack_ts = alarm['ack']['t']
                        ackremove = get_previous_step(
                            alarm,
                            'ackremove',
                            ts=ack_ts
                        )
                        ts = alarm_ts if ackremove is None else ackremove['t']
                        ack_delay = ack_ts - ts

                        events.append(eventmgr.alarm_ack_delay(ack_delay))
                        events.append(
                            eventmgr.alarm_ack_solved_delay(
                                solved_delay - ack_delay
                            )
                        )

                        events.append(usermgr.alarm_ack_delay(
                            alarm['ack']['a'],
                            ack_delay
                        ))

                    if len(alarm_events) > 0:
                        events.append(eventmgr.alarm(alarm_events[0]))

                    for event in alarm_events:
                        if event['event_type'] == 'ack':
                            events.append(eventmgr.alarm_ack(event))
                            events.append(
                                usermgr.alarm_ack(event, event['author'])
                            )

                        elif event['timestamp'] == alarm['resolved']:
                            events.append(eventmgr.alarm_solved(event))

                            if alarm['ack'] is not None:
                                events.append(eventmgr.alarm_ack_solved(event))

                                events.append(
                                    usermgr.alarm_ack_solved(
                                        alarm['ack']['a'],
                                        alarm['resolved'] - alarm['ack']['t']
                                    )
                                )

                                events.append(
                                    usermgr.alarm_solved(
                                        alarm['ack']['a'],
                                        alarm['resolved'] - alarm_ts
                                    )
                                )

                    alertsmgr.update_current_alarm(
                        docalarm,
                        alarm,
                        tags='stats'
                    )

    for event in events:
        publish(publisher=engine.amqp, event=event, logger=logger)
