#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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


"""
Alerts manager (also known as Alarm Manager).

Warning : this class is deprecated and will be replaced in a near future.
Please see the canopsis.alarms package to add/edit features.
"""

from __future__ import unicode_literals
from datetime import datetime
from operator import itemgetter
from time import time, mktime

from canopsis.alerts import DEFAULT_AUTHOR
from canopsis.alarms.models import AlarmState
from canopsis.alerts.enums import AlarmField, States, AlarmFilterField
from canopsis.alerts.filter import AlarmFilters
from canopsis.alerts.status import (
    get_last_state, get_last_status, OFF, STEALTHY,
    CANCELED, is_stealthy, is_keeped_state
)
from canopsis.check import Check
from canopsis.common.amqp import AmqpPublisher
from canopsis.common.amqp import get_default_connection as \
    get_default_amqp_conn
from canopsis.common.collection import MongoCollection
from canopsis.common.ethereal_data import EtherealData
from canopsis.common.mongo_store import MongoStore
from canopsis.common.utils import ensure_iterable, gen_id
from canopsis.confng import Configuration, Ini
from canopsis.confng.helpers import cfg_to_array, cfg_to_bool
from canopsis.context_graph.manager import ContextGraph
from canopsis.event import get_routingkey
from canopsis.lock.manager import AlertLockRedis
from canopsis.logger import Logger
from canopsis.common.middleware import Middleware
from canopsis.models.entity import Entity
from canopsis.task.core import get_task
from canopsis.timeserie.timewindow import get_offset_timewindow
from canopsis.statsng.enums import StatCounters, StatStateIntervals
from canopsis.statsng.event_publisher import StatEventPublisher
from canopsis.watcher.manager import Watcher

# register tasks manually
import canopsis.alerts.tasks as __alerts_tasks

# Extra fields from the event that should be stored in the alarm
DEFAULT_EXTRA_FIELDS = 'domain,perimeter'

# if set to True, the last_event_date will be updated on each event that
# triggers the alarm
DEFAULT_RECORD_LAST_EVENT_DATE = False

DEFAULT_FLAPPING_INTERVAL = 0
DEFAULT_FLAPPING_FREQ = 0
DEFAULT_PERSISTANT_STEPS = 10
DEFAULT_HARD_LIMIT = 2000
DEFAULT_STEALTHY_INTERVAL = 0
DEFAULT_STEALTHY_SHOW_DURATION = 0
DEFAULT_RESTORE_EVENT = False
DEFAULT_CANCEL_AUTOSOLVE_DELAY = 3600
DEFAULT_DONE_AUTOSOLVE_DELAY = 900


class Alerts(object):
    """
    Alarm cycle managment.

    Used to archive events related to alarms in a TimedStorage.
    """
    LOG_PATH = 'var/log/alerts.log'
    CONF_PATH = 'etc/alerts/manager.conf'
    ALERTS_CAT = 'ALERTS'
    FILTER_CAT = 'FILTER'

    ALERTS_STORAGE_URI = 'mongodb-periodical-alarm://'
    CONFIG_STORAGE_URI = 'mongodb-object://'
    CONFIG_COLLECTION = 'object'
    FILTER_STORAGE_URI = 'mongodb-default-alarmfilter://'

    AUTHOR = 'author'

    @property
    def alarm_filters(self):
        """
        Automatic filters and actions for alarms.
        """
        return AlarmFilters(storage=self.filter_storage,
                            alarm_storage=self.alerts_storage,
                            logger=self.logger)

    def __init__(
            self,
            config,
            logger,
            alerts_storage,
            config_data,
            filter_storage,
            context_graph,
            watcher,
            event_publisher
    ):
        self.config = config
        self.logger = logger
        self.alerts_storage = alerts_storage
        self.config_data = config_data
        self.filter_storage = filter_storage
        self.context_manager = context_graph
        self.watcher_manager = watcher

        self.event_publisher = event_publisher

        alerts_ = self.config.get(self.ALERTS_CAT, {})
        self.extra_fields = cfg_to_array(alerts_.get('extra_fields',
                                                     DEFAULT_EXTRA_FIELDS))
        self.record_last_event_date = cfg_to_bool(alerts_.get('record_last_event_date',
                                                              DEFAULT_RECORD_LAST_EVENT_DATE))

        self.update_longoutput_fields = alerts_.get("update_long_output",
                                                          False)
        filter_ = self.config.get(self.FILTER_CAT, {})
        self.filter_author = filter_.get('author', DEFAULT_AUTHOR)
        self.lock_manager = AlertLockRedis(*AlertLockRedis.provide_default_basics())

    @classmethod
    def provide_default_basics(cls):
        """
        Provide logger, config, storages...

        ! Do not use in tests !

        :rtype: Union[canopsis.confng.simpleconf.Configuration
                      logging.Logger,
                      canopsis.storage.core.Storage,
                      canopsis.common.ethereal_data.EtherealData,
                      canopsis.storage.core.Storage,
                      canopsis.context_graph.manager.ContextGraph,
                      canopsis.watcher.manager.Watcher]
        """
        config = Configuration.load(Alerts.CONF_PATH, Ini)
        conf_store = Configuration.load(MongoStore.CONF_PATH, Ini)

        mongo = MongoStore(config=conf_store)
        config_collection = MongoCollection(
            mongo.get_collection(name=cls.CONFIG_COLLECTION))
        filter_ = {'crecord_type': 'statusmanagement'}
        config_data = EtherealData(collection=config_collection,
                                   filter_=filter_)

        logger = Logger.get('alerts', cls.LOG_PATH)
        alerts_storage = Middleware.get_middleware_by_uri(
            cls.ALERTS_STORAGE_URI
        )
        filter_storage = Middleware.get_middleware_by_uri(
            cls.FILTER_STORAGE_URI
        )
        context_manager = ContextGraph(logger)
        watcher_manager = Watcher()

        amqp_pub = AmqpPublisher(get_default_amqp_conn(), logger)
        event_publisher = StatEventPublisher(logger, amqp_pub)

        return (config, logger, alerts_storage, config_data,
                filter_storage, context_manager, watcher_manager,
                event_publisher)

    @property
    def cancel_autosolve_delay(self):
        """
        Once a canceled alarm is resolved, it cannot be uncanceled. This delay
        allows users to uncancel an alarm if they made a mistake.
        """
        return self.config_data.get('cancel_autosolve_delay',
                                    DEFAULT_CANCEL_AUTOSOLVE_DELAY)

    @property
    def done_autosolve_delay(self):
        """
        Automatically close done alarms after a delay.
        """
        return self.config_data.get('done_autosolve_delay',
                                    DEFAULT_DONE_AUTOSOLVE_DELAY)

    @property
    def flapping_freq(self):
        """
        Number of alarm oscillation during flapping interval
        to consider an alarm as flapping.
        """
        #  The minimum accepted frequency is 3 changes, otherwise all alarms will bagot
        freq = self.config_data.get('bagot_freq', DEFAULT_FLAPPING_FREQ)
        if freq < 3:
            return 3

        return freq

    @property
    def flapping_interval(self):
        """
        Interval used to check for flapping alarm status.
        """
        return self.config_data.get('bagot_time', DEFAULT_FLAPPING_INTERVAL)

    @property
    def flapping_persistant_steps(self):
        """
        Number of state change steps to keep in case of long term flapping.
        Most recent steps are kept.
        """
        return self.config_data.get('persistant_steps',
                                    DEFAULT_PERSISTANT_STEPS)

    @property
    def hard_limit(self):
        """
        Maximum number of steps an alarm can have. Only alarm cancelation or
        hard limit extension are possible ways to interact with an alarm that
        has reached this point.
        """
        return self.config_data.get('hard_limit', DEFAULT_HARD_LIMIT)

    @property
    def restore_event(self):
        """
        When alarm is restored, reset the previous status if ``True``,
        recompute status with alarm history if ``False``.
        """
        return self.config_data.get('restore_event', DEFAULT_RESTORE_EVENT)

    @property
    def stealthy_interval(self):
        """
        Interval used to check for stealthy alarm status.
        """
        return self.config_data.get('stealthy_time', DEFAULT_STEALTHY_INTERVAL)

    def get_alarms(
            self,
            resolved=True,
            tags=None,
            exclude_tags=None,
            timewindow=None,
            snoozed=False
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

        :param snoozed: If ``False``, return all non-snoozed alarms, else
                        returns alarms even if they are snoozed.
        :type snoozed: bool

        :returns: Iterable of alarms matching: {alarm_id: [alarm_dict]}
        """

        query = {}

        if resolved:
            query['resolved'] = {'$ne': None}

        else:
            query['resolved'] = None

        tags_cond = None

        if tags is not None:
            tags_cond = {'$all': ensure_iterable(tags)}

        notags_cond = None

        if exclude_tags is not None:
            notags_cond = {'$not': {'$all': ensure_iterable(exclude_tags)}}

        if tags_cond is None and notags_cond is not None:
            query['tags'] = notags_cond

        elif tags_cond is not None and notags_cond is None:
            query['tags'] = tags_cond

        elif tags_cond is not None and notags_cond is not None:
            query = {'$and': [
                query,
                {'tags': tags_cond},
                {'tags': notags_cond}
            ]}

        # used to fetch alarms that were never snoozed OR alarms for which the snooze has expired
        if not snoozed:
            no_snooze_cond = {
                '$or': [
                    {AlarmField.snooze.value: None},
                    {'snooze.val': {'$lte': int(time())}}
                ]
            }
            query = {'$and': [query, no_snooze_cond]}

        alarms_by_entity = self.alerts_storage.find(
            _filter=query,
            timewindow=timewindow
        )

        for entity_id, alarms in alarms_by_entity.items():
            entity = self.context_manager.get_entities_by_id(entity_id)
            try:
                entity = entity[0]
            except IndexError:
                entity = {}

            entity['entity_id'] = entity_id
            for alarm in alarms:
                alarm['entity'] = entity

        return alarms_by_entity

    def get_alarm_with_eid(self, eid, resolved=False):
        """
        Get alarms from an eid.

        :param eid: The desired entity_id
        :type eid: str
        :param resolved: Only see resolved, or unresolved alarms
        :type resolved: bool
        """
        query = {'d': eid}
        if resolved:
            query['resolved'] = {'$ne': None}
        else:
            query['resolved'] = None

        return list(self.alerts_storage.get_elements(query=query))

    def get_current_alarm(self, alarm_entity_id):
        """
        Get current unresolved alarm.

        :param alarm_id: Alarm entity ID
        :type alarm_id: str

        :returns: Alarm as dict if found, else None
        """
        storage = self.alerts_storage

        result = storage.get(
            alarm_entity_id,
            timewindow=get_offset_timewindow(),
            _filter={
                AlarmField.resolved.value: None
            },
            limit=1
        )

        if result is not None:
            result = result[0]
            result[storage.DATA_ID] = alarm_entity_id

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
        storage = self.alerts_storage

        alarm_id = alarm[storage.DATA_ID]
        alarm_ts = alarm[storage.TIMESTAMP]

        if AlarmField.display_name.value not in new_value:
            display_name = gen_id()
            while self.check_if_display_name_exists(display_name):
                display_name = gen_id()
            new_value[AlarmField.display_name.value] = display_name

        if tags is not None:
            for tag in ensure_iterable(tags):
                if tag not in new_value[AlarmField.tags.value]:
                    new_value[AlarmField.tags.value].append(tag)

        storage.put(alarm_id, new_value, alarm_ts)

        self.watcher_manager.alarm_changed(alarm['data_id'])

    def get_events(self, alarm):
        """
        Rebuild events from alarm history.

        :param alarm: Alarm to use for events reconstruction
        :type alarm: dict

        :returns: Array of events
        """

        storage = self.alerts_storage
        alarm_id = alarm[storage.DATA_ID]
        alarm = alarm[storage.VALUE]

        entity = self.context_manager.get_entities_by_id(alarm_id)
        try:
            entity = entity[0]
        except IndexError:
            entity = {}

        no_author_types = ['stateinc', 'statedec', 'statusinc', 'statusdec']
        check_referer_types = [
            'ack',
            'ackremove',
            'cancel',
            'uncancel',
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
            'changestate': States.changestate.value,
            'snooze': 'snooze'
        }
        valmap = {
            'stateinc': Check.STATE,
            'statedec': Check.STATE,
            'changestate': Check.STATE,
            'statusinc': Check.STATUS,
            'statusdec': Check.STATUS,
            'assocticket': 'ticket',
            'snooze': 'duration'
        }
        events = []
        eventmodel = self.context_manager.get_event(entity)
        try:
            eventmodel.pop("_id")
            eventmodel.pop("depends")
            eventmodel.pop("impact")
            eventmodel.pop("infos")
            eventmodel.pop("measurements")
            eventmodel.pop("type")
        except KeyError:
            # FIXME : A logger would be nice
            pass

        for step in alarm[AlarmField.steps.value]:
            event = eventmodel.copy()
            event['timestamp'] = step['t']
            event['output'] = step['m']

            if step['_t'] in valmap:
                field = valmap[step['_t']]
                event[field] = step['val']

            if step['_t'] not in no_author_types:
                event[self.AUTHOR] = step['a']

            if step['_t'] in check_referer_types:
                event['ref_rk'] = get_routingkey(event)

            if Check.STATE not in event:
                event[Check.STATE] = get_last_state(alarm)

            # set event_type to step['_t'] if we don't have any valid mapping.
            event['event_type'] = typemap.get(step['_t'], step['_t'])

            for field in self.extra_fields:
                if field in alarm[AlarmField.extra.value]:
                    event[field] = alarm[AlarmField.extra.value][field]

            events.append(event)

        return events

    def check_if_the_entity_is_enabled(self, entity_id):
        """
        check if the entity is enabled

        :param str entity_id: entity_id
        :return bool is active: return if the entity is active true by default
        """
        entity = self.context_manager.get_entities_by_id(entity_id)

        if entity != []:
            try:
                return entity[0]['enabled']
            except Exception:
                self.logger.warning('entity not in context')

        return True

    def update_output_fields(self, value, event, state_updated):
        """
        Update the field output, long_output, long_output_history.
        :param value: the alarm.value field of an alarm
        :param event: the event used to update the alarm
        :param state_updated: if the state of the alarm change.
        """

        if not self.update_longoutput_fields and not state_updated:
            return value

        value[AlarmField.output.value] = event["output"]

        if value.get(AlarmField.long_output.value, "") != event["long_output"]:

            if AlarmField.long_output_history.value not in value:
                value[AlarmField.long_output_history.value] = []

            if len(value[AlarmField.long_output_history.value]) == 0:
                message = "Initial long_output set to \"{}\".".format(
                    event["long_output"])
            else:
                message = "Update long_output from \"{0}\" to \"{1}\".".format(
                    value[AlarmField.long_output.value],
                    event["long_output"])

            value[AlarmField.long_output.value] = event["long_output"]

            value[AlarmField.long_output_history.value].append(
                event[AlarmField.long_output.value]
            )

            if len(value[AlarmField.long_output_history.value]) > 100:
                new_hist = value[AlarmField.long_output_history.value][0:99]
                value[AlarmField.long_output_history.value] = new_hist

            value[AlarmField.steps.value].append({
                "a": event.get(self.AUTHOR, self.filter_author),
                "_t": "long_output",
                "m": message,
                "t": int(time()),
                "val": value["state"]["val"]
            })

        return value

    def archive(self, event):
        """
        Archive event in corresponding alarm history.

        :param dict event: Event to archive
        """
        entity_id = self.context_manager.get_id(event)
        event_type = event['event_type']
        initial_state = None

        lock_id = self.lock_manager.lock(entity_id)
        if event_type in [Check.EVENT_TYPE, 'watcher']:
            initial_state = event["state"]
            alarm = self.get_current_alarm(entity_id)

            is_new_alarm = alarm is None

            if is_new_alarm:
                if event[Check.STATE] == Check.OK:
                    # If a check event with an OK state concerns an entity for
                    # which no alarm is opened, there is no point continuing
                    self.lock_manager.unlock(lock_id)
                    return
                if not self.check_if_the_entity_is_enabled(entity_id):
                    self.lock_manager.unlock(lock_id)
                    return
                # Check is not OK
                alarm = self.make_alarm(entity_id, event)
                alarm = self.update_state(alarm, event[Check.STATE], event)

            else:  # Alarm is already opened
                initial_state = alarm["value"]["state"]["val"]
                value = alarm.get(self.alerts_storage.VALUE)
                if self.is_hard_limit_reached(value):
                    self.lock_manager.unlock(lock_id)
                    return
                if not self.check_if_the_entity_is_enabled(entity_id):
                    self.lock_manager.unlock(lock_id)
                    return

                alarm = self.update_state(alarm, event[Check.STATE], event)

            # set default value to event["long_output"] and event["output"]
            if "long_output" not in event:
                event["long_output"] = alarm.get(AlarmField.long_output.value,
                                                 "")

            if "output" not in event:
                event["output"] = alarm.get(AlarmField.output.value, "")

            state_updated = not initial_state == alarm["value"]["state"]["val"]

            value = alarm.get(self.alerts_storage.VALUE)

            value = self.update_output_fields(value, event, state_updated)

            value = self.crop_flapping_steps(value)

            value = self.check_hard_limit(value)

            self.update_current_alarm(alarm, value)

            if is_new_alarm:
                self.check_alarm_filters()
                self.publish_new_alarm_stats(alarm, event.get(self.AUTHOR))

        else:
            self.execute_task('alerts.useraction.{}'.format(event_type),
                              event=event,
                              author=event.get(self.AUTHOR, self.filter_author),
                              entity_id=entity_id)
        self.lock_manager.unlock(lock_id)

    def execute_task(self, name, event, entity_id,
                     author=None, new_state=None, diff_counter=None):
        """
        Find and execute a task.

        :param str name: Name of the task to execute
        :param dict event: Event to archive
        :param str entity_id: Id of the alarm
        :param str author: If needed, the author of the event
        :param int new_state: If needed, the new state in the event
        :param int diff_counter: For crop events, the new value of the counter
        """
        # Find the corresponding task
        try:
            task = get_task(name, cacheonly=True)
            # FIXIT: https://git.canopsis.net/canopsis/canopsis/issues/298
            if not callable(task):
                raise ImportError('cannot import task "{}"'.format(name))

        except ImportError:
            self.logger.debug('Unknown task {}'.format(name))
            return

        # Find the corresponding alarm
        alarm = self.get_current_alarm(entity_id)
        if alarm is None:
            self.logger.debug(
                'Entity {} has no current alarm : ignoring'.format(entity_id)
            )
            return

        value = alarm.get(self.alerts_storage.VALUE)

        if self.is_hard_limit_reached(value):
            # Only cancel and ack are allowed when hard limit has been reached
            if event['event_type'] != 'cancel' and event['event_type'] != 'ack':
                self.logger.debug('Hard limit reached. Cancelling')

                return

        # Execute the desired task
        if '.systemaction' in name:
            new_value = task(self, value, new_state, event)
        elif '.useraction' in name:
            message = event.get('output', None)
            new_value = task(self, value, author, message, event)
        elif '.lookup' in name or '.check' in name:
            new_value = task(self, value)
        elif '.crop' in name:
            new_value = task(self, value, diff_counter)
        else:
            self.logger.warning('Unknown task type for {}'.format(name))
            return

        # Some tasks return two values (a value and a status)
        status = None
        if isinstance(new_value, tuple):
            new_value, status = new_value

        if event['event_type'] != 'cancel' and event['event_type'] != 'ack':
            new_value = self.check_hard_limit(new_value)

        self.update_current_alarm(alarm, new_value)

        # If needed, update status
        if status is not None:
            alarm = self.update_status(alarm, status, event)
            new_value = alarm[self.alerts_storage.VALUE]

            self.update_current_alarm(alarm, new_value)

        return new_value

    def update_state(self, alarm, state, event):
        """
        Update alarm state if needed.

        :param dict alarm: Alarm associated to state change event
        :param int state: New state to archive
        :param dict event: Associated event
        :return: updated alarm
        :rtype: dict
        """

        value = alarm.get(self.alerts_storage.VALUE)

        if self.record_last_event_date:
            value[AlarmField.last_event_date.value] = int(time())

        old_state = get_last_state(value, ts=event['timestamp'])

        if state != old_state:
            return self.change_of_state(alarm, old_state, state, event)

        return alarm

    def update_status(self, alarm, status, event):
        """
        Update alarm status if needed.

        :param dict alarm: Alarm associated to status change event
        :param int status: New status to archive
        :param dict event: Associated event
        :return: updated alarm
        :rtype: dict
        """

        value = alarm.get(self.alerts_storage.VALUE)

        old_status = get_last_status(value, ts=event['timestamp'])

        if status != old_status:
            return self.change_of_status(
                alarm,
                old_status,
                status,
                event
            )

        return alarm

    def change_of_state(self, alarm, old_state, state, event):
        """
        Change state when ``update_state()`` detected a state change.

        :param dict alarm: Associated alarm to state change event
        :param int old_state: Previous state
        :param int state: New state
        :param dict event: Associated event
        :return: alarm with changed state
        :rtype: dict
        """

        storage_value = self.alerts_storage.VALUE
        # Check for a forced state on this alarm
        if is_keeped_state(alarm['value']):
            if state == Check.OK:
                # Disengaging 'keepstate' flag
                alarm[storage_value][AlarmField.state.value]['_t'] = None
            else:
                self.logger.info('Entity {} not allowed to change state: '
                                 'ignoring'.format(alarm['data_id']))
                return alarm

        # Escalation
        if state > old_state:
            task = get_task(
                'alerts.systemaction.state_increase', cacheonly=True
            )

        elif state < old_state:
            task = get_task(
                'alerts.systemaction.state_decrease', cacheonly=True
            )

        # Executing task
        now = int(time())
        value = alarm.get(self.alerts_storage.VALUE)
        new_value, status = task(self, value, state, event)
        new_value[AlarmField.last_update_date.value] = now

        entity_id = alarm[self.alerts_storage.DATA_ID]
        try:
            entity = self.context_manager.get_entities_by_id(entity_id)[0]
        except IndexError:
            entity = {}

        # Send statistics event
        last_state_change = entity.get(Entity.LAST_STATE_CHANGE)
        if last_state_change:
            self.event_publisher.publish_statstateinterval_event(
                now,
                StatStateIntervals.time_in_state,
                now - last_state_change,
                old_state,
                entity,
                new_value)

        if state == AlarmState.CRITICAL:
            self.event_publisher.publish_statcounterinc_event(
                now,
                StatCounters.downtimes,
                entity,
                new_value,
                event.get(self.AUTHOR))

        # Update entity's last_state_change
        if entity:
            entity[Entity.LAST_STATE_CHANGE] = now
            self.context_manager.update_entity_body(entity)

        alarm[storage_value] = new_value

        return self.update_status(alarm, status, event)

    def change_of_status(self, alarm, old_status, status, event):
        """
        Change status when ``update_status()`` detected a status
        change.

        :param dict alarm: Associated alarm to status change event
        :param int old_status: Previous status
        :param int status: New status
        :param dict event: Associated event
        :return: alarm with changed status
        :rtype: dict
        """

        if status > old_status:
            task = get_task(
                'alerts.systemaction.status_increase', cacheonly=True
            )

        elif status < old_status:
            task = get_task(
                'alerts.systemaction.status_decrease', cacheonly=True
            )

        value = alarm.get(self.alerts_storage.VALUE)
        new_value = task(self, value, status, event)
        new_value[AlarmField.last_update_date.value] = int(time())

        alarm[self.alerts_storage.VALUE] = new_value

        entity_id = alarm[self.alerts_storage.DATA_ID]

        if status == CANCELED:
            entity = self.context_manager.get_entities_by_id(entity_id)
            try:
                entity = entity[0]
            except IndexError:
                entity = {}
            self.event_publisher.publish_statcounterinc_event(
                new_value[AlarmField.last_update_date.value],
                StatCounters.alarms_canceled,
                entity,
                new_value,
                event.get(self.AUTHOR))

        return alarm

    def make_alarm(self, alarm_id, event):
        """
        Create a new alarm from event if not already existing.

        :param str alarm_id: Alarm entity ID
        :param dict event: Associated event
        :return alarm document:
        :rtype: dict
        """
        display_name = gen_id()
        while self.check_if_display_name_exists(display_name):
            display_name = gen_id()

        return {
            self.alerts_storage.DATA_ID: alarm_id,
            self.alerts_storage.TIMESTAMP: event['timestamp'],
            self.alerts_storage.VALUE: {
                AlarmField.display_name.value: display_name,
                'connector': event['connector'],
                'connector_name': event['connector_name'],
                'component': event['component'],
                'resource': event.get('resource', None),
                AlarmField.initial_output.value: event.get('output', ''),
                AlarmField.creation_date.value: int(time()),
                AlarmField.last_update_date.value: int(time()),
                AlarmField.state.value: None,
                AlarmField.status.value: None,
                AlarmField.ack.value: None,
                AlarmField.canceled.value: None,
                AlarmField.snooze.value: None,
                AlarmField.hard_limit.value: None,
                AlarmField.ticket.value: None,
                AlarmField.resolved.value: None,
                AlarmField.steps.value: [],
                AlarmField.tags.value: [],
                AlarmField.extra.value: {
                    field: event[field]
                    for field in self.extra_fields
                    if field in event
                },
                AlarmField.initial_long_output.value:
                event.get("long_output", "")
            }
        }

    def check_if_display_name_exists(self, display_name):
        """
        Check if a display_name is already associated.

        :param str display_name: the name to check
        :rtype: bool
        """
        tmp_alarms = self.alerts_storage.get_elements(
            query={'v.display_name': display_name}
        )
        if len(tmp_alarms) == 0:
            return False

        return True

    def crop_flapping_steps(self, alarm):
        """
        Remove old state changes for alarms that are flapping over long periods
        of time.

        :param dict alarm: Alarm value

        :return: Alarm with cropped steps or alarm if nothing to remove
        :rtype: dict
        """

        p_steps = self.flapping_persistant_steps

        if p_steps < 0:
            self.logger.warning(
                "Peristant steps is {} (< 0) : aborting flapping steps crop "
                "operation".format(p_steps)
            )
            return alarm

        last_status_i = alarm[AlarmField.steps.value].index(
            alarm[AlarmField.status.value])

        state_changes = filter(
            lambda step: step['_t'] in ['stateinc', 'statedec'],
            alarm[AlarmField.steps.value][last_status_i + 1:]
        )

        number_to_crop = len(state_changes) - p_steps

        if not number_to_crop > 0:
            return alarm

        crop_counter = {}

        # Removed steps are supposed unique due to their timestamps, so as
        # `remove` method does not cause any collisions.
        for i in range(number_to_crop):
            # Increase statedec or stateinc counter
            t = state_changes[i]['_t']
            crop_counter[t] = crop_counter.get(t, 0) + 1

            # Increase {0,1,2,3} counter
            s = 'state:{}'.format(state_changes[i]['val'])
            crop_counter[s] = crop_counter.get(s, 0) + 1

            alarm[AlarmField.steps.value].remove(state_changes[i])

        task = get_task('alerts.crop.update_state_counter')
        alarm = task(alarm, crop_counter)

        return alarm

    def is_hard_limit_reached(self, alarm):
        """
        Check if an hard limit is on going.

        :param dict alarm: Alarm value

        :return: False if hard_limit property is None or if configured value is
          greater than recorded value, True otherwise
        :rtype: boolean
        """

        limit = alarm.get(AlarmField.hard_limit.value, None)

        if limit is None:
            return False

        if limit['val'] < self.hard_limit:
            return False

        return True

    def check_hard_limit(self, alarm):
        """
        Update hard limit informations if number of steps has exceeded this
        limit.

        :param dict alarm: Alarm value

        :return: Alarm with hard limit informations or alarm if nothing to do
        :rtype: dict
        """

        limit = alarm.get(AlarmField.hard_limit.value, None)

        if limit is not None and limit['val'] >= self.hard_limit:
            return alarm

        if len(alarm[AlarmField.steps.value]) >= self.hard_limit:
            task = get_task('alerts.check.hard_limit')
            return task(self, alarm)

        return alarm

    def resolve_alarms(self, alarms):
        """
        Loop over all unresolved alarms, and check if it can be resolved.

        :param alarms: a list of unresolved alarms
        :return: a list of unresolved alarms (excluding locally processed alarms)
        :deprecated: see canopsis.alarms
        """
        self.logger.info("DEPRECATED: see the canopsis.alarms package instead.")

        for data_id in alarms:
            for docalarm in alarms[data_id]:
                docalarm[self.alerts_storage.DATA_ID] = data_id
                alarm = docalarm.get(self.alerts_storage.VALUE)

                if get_last_status(alarm) == OFF:
                    t = alarm[AlarmField.status.value]['t']
                    now = int(time())

                    if (now - t) > self.flapping_interval:
                        alarm[AlarmField.resolved.value] = t
                        self.update_current_alarm(docalarm, alarm)
                        alarms[data_id].remove(docalarm)

        return alarms

    def resolve_cancels(self, alarms):
        """
        Loop over all canceled alarms, and resolve the ones that are in this
        status for too long.

        :param alarms: a list of unresolved alarms
        :return: a list of unresolved alarms (excluding locally processed alarms)
        :deprecated: see canopsis.alarms
        """
        self.logger.info("DEPRECATED: see the canopsis.alarms package instead.")

        now = int(time())

        for data_id in alarms:
            for docalarm in alarms[data_id]:
                docalarm[self.alerts_storage.DATA_ID] = data_id
                alarm = docalarm.get(self.alerts_storage.VALUE)

                if alarm[AlarmField.canceled.value] is not None:
                    canceled_ts = alarm[AlarmField.canceled.value]['t']

                    if (now - canceled_ts) >= self.cancel_autosolve_delay:
                        alarm[AlarmField.resolved.value] = canceled_ts
                        self.update_current_alarm(docalarm, alarm)
                        alarms[data_id].remove(docalarm)

        return alarms

    def resolve_snoozes(self, alarms=None):
        """
        Loop over all snoozed alarms, and restore them if needed.
        :param list alarms: existing alarms (hack to bypass the self.getAlarms that is deprecated)
        :deprecated: see canopsis.alarms
        """

        self.logger.info("DEPRECATED: see the canopsis.alarms package instead.")
        now = int(time())
        if alarms is None:
            result = self.get_alarms(resolved=False, snoozed=True)
        else:
            result = alarms

        for data_id in result:
            for docalarm in result[data_id]:

                docalarm[self.alerts_storage.DATA_ID] = data_id
                alarm = docalarm.get(self.alerts_storage.VALUE)

                # if the alarm is snoozed...
                if alarm is None or \
                   AlarmField.snooze.value not in alarm or \
                   not isinstance(alarm[AlarmField.snooze.value], dict):
                    continue

                # ... and snooze is over ...
                if now > alarm[AlarmField.snooze.value]['val']:
                    # ... remove the 'snooze' key in alarm
                    alarm[AlarmField.snooze.value] = None
                    self.logger.info('Clear snooze value on alarm {}'
                                     .format(data_id))
                    self.update_current_alarm(docalarm, alarm)

    def resolve_stealthy(self, alarms):
        """
        Loop over all stealthy alarms, and check if it can be return to off
        status.

        :param alarms: a list of unresolved alarms
        :return: a list of unresolved alarms (excluding locally processed alarms)
        :deprecated: see canopsis.alarms
        """
        self.logger.info("DEPRECATED: see the canopsis.alarms package instead.")

        for data_id in alarms:
            for docalarm in alarms[data_id]:
                docalarm[self.alerts_storage.DATA_ID] = data_id
                alarm = docalarm.get(self.alerts_storage.VALUE)

                # Only look at stealthy status
                if get_last_status(alarm) != STEALTHY:
                    continue

                # Is the state stealthy still valid ?
                if is_stealthy(self, alarm):
                    continue

                event = {
                    'timestamp': int(time()),  # now
                    'output': 'automaticly resolved after stealthy shown time',
                    'connector': alarm['connector'],
                    'connector_name': alarm['connector_name']
                }
                # Updating the alarm state
                alarm_new = self.change_of_status(
                    docalarm,
                    STEALTHY,
                    OFF,
                    event
                )
                self.update_current_alarm(docalarm, alarm_new['value'])
                alarms[data_id].remove(docalarm)

        return alarms

    def check_alarm_filters(self):
        """
        Do actions on alarms based on certain conditions/filters.

        This method can alter alarm[AlarmField.alarm_filter]
        """
        now = datetime.now()
        now_stamp = int(mktime(now.timetuple()))
        RUNS = AlarmFilterField.runs.value
        NEXT = AlarmFilterField.next_run.value

        storage = self.alerts_storage

        for lifter, docalarm in self.alarm_filters.get_filters():
            # Thanks to get_alarms(), we must renaming keys
            # (... as shittily as MongoPeriodicalStorage)
            docalarm[storage.DATA_ID] = docalarm.pop(storage.Key.DATA_ID)
            docalarm[storage.TIMESTAMP] = docalarm.pop(storage.Key.TIMESTAMP)
            docalarm[storage.VALUE] = docalarm.pop(storage.Key.VALUE)
            # TODO: fix MongoPeriodicalStorage and go back remove that

            alarm_id = docalarm[storage.DATA_ID]
            self.logger.debug('Checking alarmfilter {}'.format(lifter))

            value = docalarm[storage.VALUE]
            if AlarmField.alarmfilter.value not in value:
                value[AlarmField.alarmfilter.value] = {}
            # Updating next_run timestamp
            next_run = lifter.next_run(docalarm)
            old_next_run = value[AlarmField.alarmfilter.value].get(NEXT, None)
            if old_next_run != next_run:
                value[AlarmField.alarmfilter.value][NEXT] = next_run
                self.update_current_alarm(docalarm, value)

            date = datetime.fromtimestamp(docalarm[storage.TIMESTAMP])
            # Continue only if the limit condition is valid
            if date + lifter.limit > now:
                self.logger.debug('AlarmFilter {}: Limit condition is invalid'
                                  .format(lifter._id))
                continue

            # Continue only if the filter condition is valid
            if not lifter.check_alarm(docalarm):
                self.logger.debug('AlarmFilter {}: Filter condition is invalid'
                                  .format(lifter._id))
                continue

            alarmfilter = value.get(AlarmField.alarmfilter.value, {})
            # Only execute the filter once per reached limit
            if len(alarmfilter) > 0 and RUNS in alarmfilter \
               and lifter._id in alarmfilter[RUNS]:
                executions = alarmfilter[RUNS][lifter._id]
                if len(executions) >= lifter.repeat:
                    # Already repeated enough times
                    continue

                last = datetime.fromtimestamp(max(executions))
                if last + lifter.limit > now:
                    # Too soon to execute one more time all tasks
                    continue
                self.logger.info('Rerunning tasks on {} after {} hours'
                                 .format(alarm_id, lifter.limit))

            # Getting most recent step message
            steps = docalarm[storage.VALUE][AlarmField.steps.value]
            message = sorted(steps, key=itemgetter('t'))[-1]['m']

            event = {
                'timestamp': now_stamp,
                'connector': value['connector'],
                'connector_name': value['connector_name'],
                'output': lifter.output(message),
                'event_type': Check.EVENT_TYPE,
                'component': value["component"]
            }

            if value["resource"] is not None:
                event["source_type"] = "resource"
                event["resource"] = value["resource"]
            else:
                event["source_type"] = "component"

            vstate = AlarmField.state.value

            # Execute each defined action
            updated_once = False
            new_value = self.get_current_alarm(alarm_id)[storage.VALUE]
            for task in lifter.tasks:

                if vstate in new_value:
                    event[vstate] = new_value[vstate]['val']  # for changestate

                if 'systemaction.state_increase' in task:
                    event[vstate] = event[vstate] + 1
                elif 'systemaction.state_decrease' in task:
                    event[vstate] = event[vstate] - 1

                self.logger.info('Automatically execute {} on {}'
                                 .format(task, alarm_id))

                updated_alarm_value = self.execute_task(
                    name=task,
                    event=event,
                    entity_id=alarm_id,
                    author=self.filter_author,
                    new_state=event[vstate]
                )
                if updated_alarm_value is not None:
                    new_value = updated_alarm_value
                    updated_once = True
                    self.update_current_alarm(docalarm, updated_alarm_value)

            if not updated_once:
                continue

            # Mark the alarm that this filter has been applied
            new_value = self.get_current_alarm(alarm_id)[storage.VALUE]
            alarmfilter = new_value.get(AlarmField.alarmfilter.value, {})
            if RUNS not in alarmfilter:
                alarmfilter[RUNS] = {}
            if lifter._id not in alarmfilter[RUNS]:
                alarmfilter[RUNS][lifter._id] = []

            alarmfilter[RUNS][lifter._id].append(now_stamp)
            new_value[AlarmField.alarmfilter.value] = alarmfilter

            self.update_current_alarm(docalarm, new_value)

    def publish_new_alarm_stats(self, alarm, author):
        """
        Publish statistics events for a new alarm.

        :param Dict[str, Any] alarm:
        """
        entity_id = alarm[self.alerts_storage.DATA_ID]
        entity = self.context_manager.get_entities_by_id(entity_id)
        try:
            entity = entity[0]
        except IndexError:
            entity = {}

        alarm_value = alarm[self.alerts_storage.VALUE]
        creation_date = alarm_value[AlarmField.creation_date.value]

        # Increment alarms_created counter
        self.event_publisher.publish_statcounterinc_event(
            creation_date,
            StatCounters.alarms_created,
            entity,
            alarm_value,
            author)
