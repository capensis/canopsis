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

from canopsis.timeserie.timewindow import get_offset_timewindow
from canopsis.common.utils import ensure_iterable
from canopsis.task.core import get_task

from canopsis.alerts.status import get_last_state
from canopsis.event.manager import Event
from canopsis.check import Check


CONF_PATH = 'alerts/manager.conf'
CATEGORY = 'ALERTS'
CONTENT = []


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class Alerts(MiddlewareRegistry):

    CONFIG_STORAGE = 'config_storage'
    ALARM_STORAGE = 'alarm_storage'
    CONTEXT_MANAGER = 'context'
    CHECK_MANAGER = 'check'

    @property
    def config(self):
        if not hasattr(self, '_config'):
            self.config = None

        return self._config

    @config.setter
    def config(self, value):
        if value is None:
            value = self[Alerts.CONFIG_STORAGE].get_elements(
                query={'crecord_type': 'statusmanagement'}
            )
            value = {} if not value else value[0]

        self._config = value

    @property
    def flapping_interval(self):
        return self.config.get('bagot_time', 0)

    @property
    def flapping_freq(self):
        return self.config.get('bagot_freq', 0)

    @property
    def stealthy_interval(self):
        return self.config.get('stealthy_time', 0)

    @property
    def stealthy_show_duration(self):
        return self.config.get('stealthy_show', 0)

    @property
    def restore_event(self):
        return self.config.get('restore_event', False)

    def __init__(
        self,
        config_storage=None,
        alarm_storage=None,
        context=None,
        check=None,
        *args, **kwargs
    ):
        super(Alerts, self).__init__(*args, **kwargs)

        if config_storage is not None:
            self[Alerts.CONFIG_STORAGE] = config_storage

        if alarm_storage is not None:
            self[Alerts.ALARM_STORAGE] = alarm_storage

        if context is not None:
            self[Alerts.CONTEXT_MANAGER] = context

        if check is not None:
            self[Alerts.CHECK_MANAGER] = check

    def get_alarms(self, resolved=True, tags=None, exclude_tags=None):
        query = {}

        if resolved:
            query['resolved'] = {'$ne': None}

        else:
            query['resolved'] = None

        tags_cond = None

        if tags is not None:
            tags_cond = {'$in': ensure_iterable(tags)}

        notags_cond = None

        if exclude_tags is not None:
            notags_cond = {'$nin': ensure_iterable(tags)}

        if tags_cond is None and notags_cond is not None:
            query['tags'] = notags_cond

        elif tags_cond is not None and notags_cond is None:
            query['tags'] = tags_cond

        elif tags_cond is not None and notags_cond is not None:
            query = {'$and': [query, tags_cond, notags_cond]}

        return self[Alerts.ALARM_STORAGE].find(query=query)

    def get_current_alarm(self, alarm_id):

        timewindow = get_offset_timewindow()

        return self[Alerts.ALARM_STORAGE].get(
            alarm_id,
            timewindow=timewindow,
            _filter={
                'resolved': None
            },
            limit=1
        )

    def update_current_alarm(self, alarm, new_value, tags=None):
        storage = self[Alerts.ALARM_STORAGE]

        alarm_id = alarm[storage.DATA_ID]
        alarm_ts = alarm[storage.TIMESTAMP]

        if tags is not None:
            for tag in ensure_iterable(tags):
                if tag not in new_value['tags']:
                    new_value['tags'].append(tag)

        storage.put(alarm_id, new_value, alarm_ts)

    def get_events(self, alarm):
        storage = self[Alerts.ALARM_STORAGE]
        alarm_id = alarm[storage.DATA_ID]
        alarm = alarm[storage.VALUE]

        entity = self[Alerts.CONTEXT_MANAGER].get_entity_by_id(alarm_id)

        no_author_tupes = ['stateinc', 'statedec', 'statusinc', 'statusdec']
        check_referer_types = [
            'ack', 'ackremove',
            'cancel', 'uncancel',
            'declareticket',
            'assocticket',
            'changestate'
        ]

        typemap = {
            'stateinc': Check.EVENT_TYPE,
            'statedec': Check.EVENT_TYPE,
            'statusinc': Check.EVENT_TYPE,
            'statusdec': Check.EVENT_TYPE,
            'ack': 'ack',
            'ackremove': 'ackremove',
            'cancel': 'cancel',
            'uncancel': 'uncancel',
            'declareticket': 'declareticket',
            'assocticket': 'assocticket',
            'changestate': 'changestate'
        }
        valmap = {
            'stateinc': Check.STATE,
            'statedec': Check.STATE,
            'changestate': Check.STATE,
            'statusinc': Check.STATUS,
            'statusdec': Check.STATUS,
            'assocticket': 'ticket'
        }

        events = []
        eventmodel = self[Alerts.CONTEXT_MANAGER].get_event(entity)

        for step in alarm['steps']:
            event = eventmodel.copy()
            event['timestamp'] = step['t']
            event['output'] = step['m']

            if step['_t'] in valmap:
                field = valmap[step['_t']]
                event[field] = step['val']

            if step['_t'] not in no_author_tupes:
                event['author'] = step['a']

            if step['_t'] in check_referer_types:
                event['event_type'] = 'check'
                event['ref_rk'] = Event.get_rk(event)

            if Check.STATE not in event:
                event[Check.STATE] = get_last_state(alarm)

            event['event_type'] = typemap[step['_t']]

            events.append(event)

        return events

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
                value = alarm.get(self[Alerts.ALARM_STORAGE].VALUE)
                new_value = task(self, value, author, message, event)
                self.update_current_alarm(alarm, new_value)

    def change_of_state(self, entity, old_state, state, event):
        entity_id = self[Alerts.CONTEXT_MANAGER].get_entity_id(entity)

        if state > old_state:
            task = get_task('alerts.systemaction.state_increase')

        elif state < old_state:
            task = get_task('alerts.systemaction.state_decrease')

        alarm = self.get_current_alarm(entity_id)
        value = alarm.get(self[Alerts.ALARM_STORAGE].VALUE)
        new_value, status = task(self, value, state, event)
        self.update_current_alarm(alarm, new_value)

        # TODO: implementation needed in check manager
        # old_status = self[Alerts.CHECK_MANAGER].status(ids=entity_id)
        # status = self[Alerts.CHECK_MANAGER].status(
        #     ids=entity_id,
        #     status=status
        # )
        #
        # if status != old_status:
        #     self.change_of_status(
        #         alarm,
        #         old_status,
        #         status,
        #         event
        #     )

    def change_of_status(self, alarm, old_status, status, event):
        value = alarm.get(self[Alerts.ALARM_STORAGE].VALUE)

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
                'tags': []
            }

            self[Alerts.ALARM_STORAGE].put(alarm_id, value, timestamp)
