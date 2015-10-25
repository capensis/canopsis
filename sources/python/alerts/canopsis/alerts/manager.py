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

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category

from canopsis.task.core import get_task
from canopsis.check import Check


CONF_PATH = 'alerts/manager.conf'
CATEGORY = 'ALERTS'
CONTENT = []


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class Alerts(MiddlewareRegistry):

    ALARM_STORAGE = 'alarm_storage'
    CONTEXT_MANAGER = 'context'
    CHECK_MANAGER = 'check'

    def __init__(
        self,
        alarm_storage=None,
        context=None,
        check=None,
        *args, **kwargs
    ):
        super(Alerts, self).__init__(*args, **kwargs)

        if alarm_storage is not None:
            self[Alerts.ALARM_STORAGE] = alarm_storage

        if context is not None:
            self[Alerts.CONTEXT_MANAGER] = context

        if check is not None:
            self[Alerts.CHECK_MANAGER] = check

    def get_current_alarm(self, alarm_id):
        return self[Alerts.ALARM_STORAGE].get(alarm_id, limit=1)

    def update_current_alarm(self, alarm, new_value):
        storage = self[Alerts.ALARM_STORAGE]

        alarm_id = alarm[storage.Key.DATA_ID]
        alarm_ts = alarm[storage.Key.TIMESTAMP]

        storage.put(alarm_id, new_value, alarm_ts)

    def archive(self, event):
        entity = self[Alerts.CONTEXT_MANAGER].get_entity(event)
        entity_id = self[Alerts.CONTEXT_MANAGER].get_entity_id(entity)

        author = event.get('author', None)
        message = event.get('output', None)

        if event['type'] == Check.EVENT_TYPE:
            if event[Check.STATE] != Check.OK:
                self.make_alarm(entity_id, event['timestamp'])

            old_state = self[Alerts.CHECK_MANAGER].state(ids=entity_id)
            state = self[Alerts.CHECK_MANAGER].state(
                ids=entity_id,
                state=event[Check.STATE]
            )

            if state != old_state:
                self.change_of_state(entity, old_state, state, event)

        else:
            task = get_task('alerts.useraction.{0}'.format(event['type']))

            if task is not None:
                alarm = self.get_current_alarm(entity_id)
                value = alarm.get(self[Alerts.ALARM_STORAGE].Key.VALUE)
                new_value = task(self, value, author, message, event)
                self.update_current_alarm(alarm, new_value)

    def change_of_state(self, entity, old_state, state, event):
        entity_id = self[Alerts.CONTEXT_MANAGER].get_entity_id(entity)

        if state > old_state:
            task = get_task('alerts.systemaction.state_increase')

        elif state < old_state:
            task = get_task('alerts.systemaction.state_decrease')

        alarm = self.get_current_alarm(entity_id)
        value = alarm.get(self[Alerts.ALARM_STORAGE].Key.VALUE)
        new_value = task(self, value, state, event)
        self.update_current_alarm(alarm, new_value)

        # TODO: implementation needed in check manager
        # old_status = self[Alerts.CHECK_MANAGER].status(ids=entity_id)
        # status = self[Alerts.CHECK_MANAGER].status(
        #     ids=entity_id,
        #     status=new_value['status']['val']
        # )
        #
        # if status != old_status:
        #     self.change_of_status(
        #         alarm,
        #         old_status,
        #         new_value['status'],
        #         event
        #     )

    def change_of_status(self, alarm, old_status, status, event):
        value = alarm.get(self[Alerts.ALARM_STORAGE].Key.VALUE)

        if status > old_status:
            task = get_task('alerts.systemaction.status_increase')

        elif status < old_status:
            task = get_task('alerts.systemaction.status_decrease')

        new_value = task(self, value, status, event)
        self.update_current_alarm(alarm, new_value)

    def make_alarm(self, alarm_id, timestamp):
        alarm = self.get_current_alarm(alarm_id)

        if alarm is None:
            value = {
                'state': None,
                'status': None,
                'ack': None,
                'canceled': None,
                'ticket': None,
                'resolved': None,
                'steps': [],
            }

            self[Alerts.ALARM_STORAGE].put(alarm_id, value, timestamp)
