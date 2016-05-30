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
from canopsis.configuration.model import Parameter

from canopsis.timeserie.timewindow import get_offset_timewindow
from canopsis.common.utils import ensure_iterable
from canopsis.task.core import get_task

from canopsis.event.manager import Event
from canopsis.check import Check

from canopsis.alerts.status import get_last_state, get_last_status, OFF

from time import time


CONF_PATH = 'alerts/manager.conf'
CATEGORY = 'ALERTS'
CONTENT = [
    Parameter('extra_fields', Parameter.array())
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class Alerts(MiddlewareRegistry):
    """
    Alarm cycle managment.

    Used to archive events related to alarms in a TimedStorage.
    """

    CONFIG_STORAGE = 'config_storage'
    ALARM_STORAGE = 'alarm_storage'
    CONTEXT_MANAGER = 'context'

    @property
    def config(self):
        """
        Property computed from configuration storage.
        """

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
        """
        Interval used to check for flapping alarm status.
        """

        return self.config.get('bagot_time', 0)

    @property
    def flapping_freq(self):
        """
        Number of alarm oscillation during flapping interval
        to consider an alarm as flapping.
        """

        return self.config.get('bagot_freq', 0)

    @property
    def stealthy_interval(self):
        """
        Interval used to check for stealthy alarm status.
        """

        return self.config.get('stealthy_time', 0)

    @property
    def stealthy_show_duration(self):
        """
        Interval used to check if alarm is still in stealthy status.
        """

        return self.config.get('stealthy_show', 0)

    @property
    def restore_event(self):
        """
        When alarm is restored, reset the previous status if ``True``,
        recompute status with alarm history if ``False``.
        """

        return self.config.get('restore_event', False)

    @property
    def extra_fields(self):
        """
        Array of fields to save from event in alarm.
        """

        if not hasattr(self, '_extra_fields'):
            self.extra_fields = None

        return self._extra_fields

    @extra_fields.setter
    def extra_fields(self, value):
        if value is None:
            value = []

        self._extra_fields = value

    def __init__(
        self,
        extra_fields=None,
        config_storage=None,
        alarm_storage=None,
        context=None,
        *args, **kwargs
    ):
        super(Alerts, self).__init__(*args, **kwargs)

        if extra_fields is not None:
            self.extra_fields = extra_fields

        if config_storage is not None:
            self[Alerts.CONFIG_STORAGE] = config_storage

        if alarm_storage is not None:
            self[Alerts.ALARM_STORAGE] = alarm_storage

        if context is not None:
            self[Alerts.CONTEXT_MANAGER] = context

    def get_alarms(
            self,
            resolved=True,
            tags=None,
            exclude_tags=None,
            timewindow=None
    ):
        """
        Get alarms from TimedStorage.

        :param resolved: If ``True``, returns only resolved alarms, else
                         returns only unresolved alarms (default: ``True``).
        :type resolved: bool
        :param tags: Tags which must be set on alarm (optional)
        :type tags: str or list

        :param exclude_tags: Tags which must not be set on alarm (optional)
        :type tags: str or list

        :param timewindow: Time Window used for fetching (optional)
        :type timewindow: canopsis.timeserie.timewindow.TimeWindow

        :returns: Iterable of alarms matching
        """

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
            notags_cond = {'$nin': ensure_iterable(exclude_tags)}

        if tags_cond is None and notags_cond is not None:
            query['tags'] = notags_cond

        elif tags_cond is not None and notags_cond is None:
            query['tags'] = tags_cond

        elif tags_cond is not None and notags_cond is not None:
            query = {'$and': [query, tags_cond, notags_cond]}

        return self[Alerts.ALARM_STORAGE].find(
            _filter=query,
            timewindow=timewindow
        )

    def get_current_alarm(self, alarm_id):
        """
        Get current unresolved alarm.

        :param alarm_id: Alarm entity ID
        :type alarm_id: str

        :returns: Alarm as dict if found, else None
        """

        storage = self[Alerts.ALARM_STORAGE]

        result = storage.get(
            alarm_id,
            timewindow=get_offset_timewindow(),
            _filter={
                'resolved': None
            },
            limit=1
        )

        if result is not None:
            result = result[0]
            result[storage.DATA_ID] = alarm_id

        return result

    def update_current_alarm(self, alarm, new_value, tags=None):
        """
        Update alarm's history and tags.

        :param alarm: Alarm to update
        :type alarm: dict

        :param new_value: New history to set on alarm
        :type new_value: dict

        :param tags: Tags to add on alarm (optional)
        :type tags: str or list
        """

        storage = self[Alerts.ALARM_STORAGE]

        alarm_id = alarm[storage.DATA_ID]
        alarm_ts = alarm[storage.TIMESTAMP]

        if tags is not None:
            for tag in ensure_iterable(tags):
                if tag not in new_value['tags']:
                    new_value['tags'].append(tag)

        storage.put(alarm_id, new_value, alarm_ts)

    def get_events(self, alarm):
        """
        Rebuild events from alarm history.

        :param alarm: Alarm to use for events reconstruction
        :type alarm: dict

        :returns: Array of events
        """

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

            for field in self.extra_fields:
                if field in alarm['extra']:
                    event[field] = alarm['extra'][field]

            events.append(event)

        return events

    def archive(self, event):
        """
        Archive event in corresponding alarm history.

        :param event: Event to archive
        :type event: dict
        """

        entity = self[Alerts.CONTEXT_MANAGER].get_entity(event)
        entity_id = self[Alerts.CONTEXT_MANAGER].get_entity_id(entity)

        author = event.get('author', None)
        message = event.get('output', None)

        if event['event_type'] == Check.EVENT_TYPE:
            if event[Check.STATE] != Check.OK:
                self.make_alarm(entity_id, event)

            alarm = self.get_current_alarm(entity_id)

            if alarm is not None:
                self.archive_state(alarm, event[Check.STATE], event)

        else:
            try:
                task = get_task('alerts.useraction.{0}'.format(
                    event['event_type']
                ), cacheonly=True)

            except ImportError:
                task = None

            if task is not None:
                alarm = self.get_current_alarm(entity_id)
                value = alarm.get(self[Alerts.ALARM_STORAGE].VALUE)

                new_value = task(self, value, author, message, event)
                status = None

                if isinstance(new_value, tuple):
                    new_value, status = new_value

                self.update_current_alarm(alarm, new_value)

                if status is not None:
                    self.archive_status(alarm, status, event)

    def archive_state(self, alarm, state, event):
        """
        Archive state if needed.

        :param alarm: Alarm associated to state change event
        :type alarm: dict

        :param state: New state to archive
        :type state: int

        :param event: Associated event
        :type event: dict
        """

        value = alarm.get(self[Alerts.ALARM_STORAGE].VALUE)

        old_state = get_last_state(value, ts=event['timestamp'])

        if state != old_state:
            self.change_of_state(alarm, old_state, state, event)

    def archive_status(self, alarm, status, event):
        """
        Archive status if needed.

        :param alarm: Alarm associated to status change event
        :type alarm: dict

        :param status: New status to archive
        :type status: int

        :param event: Associated event
        :type event: dict
        """

        value = alarm.get(self[Alerts.ALARM_STORAGE].VALUE)

        old_status = get_last_status(value, ts=event['timestamp'])

        if status != old_status:
            self.change_of_status(
                alarm,
                old_status,
                status,
                event
            )

    def change_of_state(self, alarm, old_state, state, event):
        """
        Archive state change when ``archive_state()`` detected a state change.

        :param alarm: Associated alarm to state change event
        :type alarm: dict

        :param old_state: Previous state
        :type old_state: int

        :param state: New state
        :type state: int

        :param event: Associated event
        :type event: dict
        """

        if state > old_state:
            task = get_task(
                'alerts.systemaction.state_increase', cacheonly=True
            )

        elif state < old_state:
            task = get_task(
                'alerts.systemaction.state_decrease', cacheonly=True
            )

        value = alarm.get(self[Alerts.ALARM_STORAGE].VALUE)
        new_value, status = task(self, value, state, event)
        self.update_current_alarm(alarm, new_value)

        self.archive_status(alarm, status, event)

    def change_of_status(self, alarm, old_status, status, event):
        """
        Archive status change when ``archive_status()`` detected a status
        change.

        :param alarm: Associated alarm to status change event
        :type alarm: dict

        :param old_status: Previous status
        :type old_status: int

        :param status: New status
        :type status: int

        :param event: Associated event
        :type event: dict
        """

        if status > old_status:
            task = get_task(
                'alerts.systemaction.status_increase', cacheonly=True
            )

        elif status < old_status:
            task = get_task(
                'alerts.systemaction.status_decrease', cacheonly=True
            )

        value = alarm.get(self[Alerts.ALARM_STORAGE].VALUE)
        new_value = task(self, value, status, event)
        self.update_current_alarm(alarm, new_value)

    def make_alarm(self, alarm_id, event):
        """
        Create a new alarm from event if not already existing.

        :param alarm_id: Alarm entity ID
        :type alarm_id: str

        :param event: Associated event
        :type event: dict
        """

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
                'tags': [],
                'extra': {
                    field: event[field]
                    for field in self.extra_fields
                    if field in event
                }
            }

            self[Alerts.ALARM_STORAGE].put(alarm_id, value, event['timestamp'])

    def resolve_alarms(self):
        """
        Loop over all unresolved alarms, and check if it can be resolved.
        """

        storage = self[Alerts.ALARM_STORAGE]
        result = self.get_alarms(resolved=False)

        for data_id in result:
            for docalarm in result[data_id]:
                docalarm[storage.DATA_ID] = data_id
                alarm = docalarm.get(storage.VALUE)

                if get_last_status(alarm) == OFF:
                    t = alarm['status']['t']
                    now = int(time())

                    if (now - t) > self.flapping_interval:
                        alarm['resolved'] = t
                        self.update_current_alarm(docalarm, alarm)
