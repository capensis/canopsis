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


@register_task
def event_processing(
    engine,
    event,
    usermgr=None,
    eventmgr=None,
    logger=None,
    **kwargs
):
    if usermgr is None:
        usermgr = singleton_per_scope(UserMetricProducer)

    if eventmgr is None:
        eventmgr = singleton_per_scope(EventMetricProducer)

    if event['type'] in ['ack', 'check']:
        for manager in [usermgr, eventmgr]:
            manager.cache(event)

    events = []

    if event['type'] == 'ack':
        events.append(usermgr.alarm_ack(event, event['ack']['author']))
        events.append(eventmgr.alarm_ack(event))

    elif event['type'] == 'check':
        if event['state'] == 0:
            events.append(eventmgr.alarm_solved(event))

            involved_events = eventmgr.get_cache(event)

            for involved_event in involved_events:
                if involved_event['type'] == 'ack':
                    events.append(eventmgr.alarm_ack_solved(event))
                    break

            eventmgr.clear_cache(event)

        elif event.get('ack', {}).get('isAck', False):
            events.append(eventmgr.alarm_ack(event))

        else:
            events.append(eventmgr.alarm(event))

    for event in events:
        publish(publisher=engine.amp, event=event, logger=logger)
