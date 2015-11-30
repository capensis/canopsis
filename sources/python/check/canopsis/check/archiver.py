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

from time import time

from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.old.record import Record
from canopsis.old.rabbitmq import Amqp
from canopsis.event import get_routingkey

from canopsis.engines.core import publish
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths
)

import pprint
pp = pprint.PrettyPrinter(indent=2)

from pymongo.errors import BulkWriteError

legend_type = ['soft', 'hard']
OFF = 0
ONGOING = 1
STEALTHY = 2
BAGOT = 3
CANCELED = 4

CONF_PATH = 'check/archiver.conf'
CATEGORY = 'ARCHIVER'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Archiver(Configurable):

    def __init__(
        self, namespace, confnamespace='object', storage=None, autolog=False,
        *args, **kwargs
    ):

        super(Archiver, self).__init__(*args, **kwargs)
        self.namespace = namespace
        self.namespace_log = namespace + '_log'

        # Bulk operation configuration
        self.last_bulk_insert_date = time()
        self.bulk_ids = []
        # How many events can be buffered
        self.bulk_amount = 500
        # What is the maximum duration until bulk insert
        self.bulk_delay = 3
        self.incoming_events = {}

        self.autolog = autolog

        self.logger.debug(u"Init Archiver on %s" % namespace)

        self.account = Account(user="root", group="root")

        if not storage:
            self.logger.debug(" + Get storage")
            self.storage = get_storage(
                namespace=namespace,
                logging_level=self.log_lvl
            )
        else:
            self.storage = storage

        self.conf_storage = get_storage(
            namespace=confnamespace,
            logging_level=self.log_lvl
        )
        self.conf_collection = self.conf_storage.get_backend(confnamespace)
        self.collection = self.storage.get_backend(namespace)

        self.amqp = Amqp(
            logging_level=self.log_lvl,
            logging_name='archiver-amqp'
        )

        self.reset_stealthy_event_duration = time()
        self.reset_stats()

    def reset_stats(self):
        self.stats = {
            'update': 0,
            'insert ' + self.namespace: 0,
            'insert ' + self.namespace_log: 0
        }

    def beat(self):

        self.logger.info(
            (
                'DB documents stats : ' +
                'update: {} in events, ' +
                'insert: {} in events, ' +
                'insert: {} in events_log').format(
                self.stats['update'],
                self.stats['insert ' + self.namespace],
                self.stats['insert ' + self.namespace_log]
            )
        )
        self.reset_stats()

    def process_insert_operations_collection(self, operations, collection):

        self.stats['insert ' + collection] += len(operations)

        if operations:
            # is there any event to process ?
            backend = self.storage.get_backend(collection)
            bulk = backend.initialize_unordered_bulk_op()
            for operation in operations:
                record = Record(operation['event'])
                record.type = "event"
                event = record.dump()
                bulk.insert(event)
            try:
                bulk.execute({'w': 0})
            except BulkWriteError as bwe:
                self.logger.warning(pp.pformat(bwe.details))
            self.logger.info(u'inserted log events {}'.format(len(operations)))

    def process_update_operations(self, operations):

        self.stats['update'] += len(operations)

        if operations:
            # is there any event to process ?
            backend = self.storage.get_backend('events')
            bulk = backend.initialize_unordered_bulk_op()
            for operation in operations:
                bulk.find(operation['query']).update(operation['update'])
            bulk.execute({'w': 0})

    def process_insert_operations(self, insert_operations):

        events = {}
        events_log = {}
        # Avoid same RK insert
        for operation in insert_operations:
            if '_id' not in operation['event']:
                self.logger.error(
                    u'Unable to find _id value in event {}'.format(
                        operation['event']
                    )
                )
            else:
                _id = operation['event']['_id']

                if operation['collection'] == self.namespace:
                    events[_id] = operation
                elif operation['collection'] == self.namespace_log:
                    _id = u'{}.{}'.format(_id, time())
                    operation['event']['_id'] = _id
                    events_log[_id] = operation
                else:
                    self.logger.critical(
                        'Wrong operation type {}'.format(
                            operation['collection']
                        )
                    )

        self.process_insert_operations_collection(
            events.values(),
            'events'
        )
        self.process_insert_operations_collection(
            events_log.values(),
            'events_log'
        )

    def reload_configuration(self):
                # Default values
        self.restore_event = True
        self.bagot_freq = 10
        self.bagot_time = 3600
        self.stealthy_time = 360
        self.stealthy_show = 360

        state_config = self.conf_collection.find_one(
            {'crecord_type': 'statusmanagement'}
        )

        if state_config is not None:
            self.bagot_freq = state_config.get('bagot_freq', self.bagot_freq)
            self.bagot_time = state_config.get('bagot_time', self.bagot_time)
            self.stealthy_time = state_config.get(
                'stealthy_time',
                self.stealthy_time
            )
            self.stealthy_show = state_config.get(
                'stealthy_show',
                self.stealthy_show
            )
            self.restore_event = state_config.get(
                'restore_event',
                self.restore_event
            )

        self.logger.debug(
            (
                'Archiver configuration ->\n' +
                ' bagot_freq {}\n' +
                ' bagot_time {}\n' +
                ' stealthy_show {}\n' +
                ' stealthy_time {}'
            ).format(
                self.bagot_freq,
                self.bagot_time,
                self.stealthy_show,
                self.stealthy_time
            )
        )

    def reset_status_event(self, reset_type):

        """Trigger event status reset to off/on going status if event are in
        BAGOT or STEALTHY status.

        :param reset_type: event status to consider and change.
        :type int: This is en enum, can be either BAGOT or STEALTHY
        """

        def _publish_event(event):
            rk = event.get('rk', get_routingkey(event))
            self.logger.info(u"Sending event {}".format(rk))
            self.logger.debug(event)
            publish(
                event=event, rk=rk, publisher=self.amqp
            )

        if reset_type not in [BAGOT, STEALTHY]:
            self.logger.info('wrong reset type given, will not process')
            return

        # Dynamic method parameter depends on reset type input
        compare_property = {
            BAGOT: 'last_state_change',
            STEALTHY: 'ts_first_stealthy'
        }[reset_type]

        configuration_delay = {
            BAGOT: self.bagot_time,
            STEALTHY: self.stealthy_show
        }[reset_type]

        event_cursor = self.collection.find(
            {
                'crecord_type': 'event',
                'status': reset_type
            }
        )

        # Change all potention reset type events
        for event in event_cursor:
            # This is a bagot event.
            is_show_delay_passed = \
                time() - event[compare_property] >= configuration_delay

            # Check the stealthy intervals
            if is_show_delay_passed:

                self.logger.info(
                    u'Event {} no longer in status {}'.format(
                        event['rk'],
                        reset_type
                    )
                )

                new_status = ONGOING if event['state'] else OFF
                self.set_status(event, new_status)
                event['pass_status'] = 1
                _publish_event(event)

    def is_bagot(self, event):
        """
        Args:
            event map of the current evet
        Returns:
            ``True`` if the event is bagot
            ``False`` otherwise
        """

        ts_curr = event['timestamp']
        ts_first_bagot = event.get('ts_first_bagot', 0)
        ts_diff_bagot = (ts_curr - ts_first_bagot)
        freq = event.get('bagot_freq', -1)

        result = ts_diff_bagot <= self.bagot_time and freq >= self.bagot_freq

        return result

    def is_stealthy(self, event, d_status):
        """
        Args:
            event map of the current evet
            d_status status of the previous event
        Returns:
            ``True`` if the event is stealthy
            ``False`` otherwise
        """

        ts_diff = event['timestamp'] - event['ts_first_stealthy']
        result = ts_diff <= self.stealthy_time and d_status != STEALTHY
        return result

    def set_status(self, event, status, devent=None):
        """
        Args:
            event map of the current event
            status status of the current event
        """

        log = u'Status is set to {} for event {}'.format(status, event['rk'])
        bagot_freq = event.get('bagot_freq', 0)
        values = {
            OFF: {
                'freq': bagot_freq,
                'name': 'Off'
            },
            ONGOING: {
                'freq': bagot_freq,
                'name': 'On going'
            },
            STEALTHY: {
                'freq': bagot_freq,
                'name': 'Stealthy'
            },
            BAGOT: {
                'freq': bagot_freq + 1,
                'name': 'Bagot'
            },
            CANCELED: {
                'freq': bagot_freq,
                'name': 'Cancelled'
            }
        }

        self.logger.debug(log.format(values[status]['name']))

        # This is an additional check as stealthy
        # status is not properly managed until now
        if status != STEALTHY:
            event['status'] = status
        elif devent['state'] != 0 and event['state'] == 0:
            delta = time() - event['last_state_change']
            if delta < self.stealthy_time:
                event['status'] = status

        event['bagot_freq'] = values[status]['freq']

        if status not in [STEALTHY, BAGOT]:
            event['ts_first_stealthy'] = 0

    def check_stealthy(self, devent, ts):
        """
        Args:
            devent map of the previous event
            ts timestamp of the current event
        Returns:
            ``True`` if the event should stay stealthy
            ``False`` otherwise
        """
        result = False

        if devent['status'] == STEALTHY:
            result = (ts - devent['ts_first_stealthy']) <= self.stealthy_show

        return result

    def check_statuses(self, event, devent):
        """
        Args:
            event map of the current event
            devent map of the previous evet
        """

        if event.get('pass_status', 0):
            event['pass_status'] = 0
            return

        event_ts = event['timestamp']

        event['bagot_freq'] = devent.get('bagot_freq', 0)
        event['ts_first_stealthy'] = devent.get('ts_first_stealthy', 0)
        event['ts_first_bagot'] = devent.get('ts_first_bagot', 0)
        dstate = devent['state']
        # Increment frequency if state changed and set first occurences
        if ((not dstate and event['state']) or
                dstate and not event['state']):

            if event['state']:
                event['ts_first_stealthy'] = event_ts
            else:
                event['ts_first_stealthy'] = event_ts

            event['bagot_freq'] += 1

            if not event['ts_first_bagot']:
                event['ts_first_bagot'] = event_ts

        # Out of bagot interval, reset variables
        if event['ts_first_bagot'] - event_ts > self.bagot_time:
            event['ts_first_bagot'] = 0
            event['bagot_freq'] = 0

        # If not canceled, proceed to check the status
        if (devent.get('status', ONGOING) != CANCELED
            or (dstate != event['state']
                and (self.restore_event
                or event['state'] == OFF
                or dstate == OFF))):
            # Check the stealthy intervals
            if self.check_stealthy(devent, event_ts):
                if self.is_bagot(event):
                    self.set_status(event, BAGOT)
                else:
                    self.set_status(event, STEALTHY, devent=devent)
            # Else proceed normally
            else:
                if (event['state'] == OFF):
                    # If still non-alert, can only be OFF
                    if (not self.is_bagot(event)
                            and not self.is_stealthy(event, devent['status'])):
                        self.set_status(event, OFF)
                    elif self.is_bagot(event):
                        self.set_status(event, BAGOT)
                    elif self.is_stealthy(event, devent['status']):
                        self.set_status(event, STEALTHY, devent=devent)
                else:
                    # If not bagot/stealthy, can only be ONGOING
                    if (not self.is_bagot(event)
                            and not self.is_stealthy(event, devent['status'])):
                        self.set_status(event, ONGOING)
                    elif self.is_bagot(event):
                        self.set_status(event, BAGOT)
                    elif self.is_stealthy(event, devent['status']):
                        if devent['status'] == OFF:
                            self.set_status(event, ONGOING)
                        else:
                            self.set_status(event, STEALTHY, devent=devent)
        else:
            self.set_status(event, CANCELED)

    def check_event(self, _id, event):
        """
            This method aims to buffer and process incoming events.
            Processing is done on buffer to reduce database operations.
        """

        # As this was not done until now... setting event primary key
        event['_id'] = _id

        # Buffering event informations
        self.bulk_ids.append(_id)
        self.incoming_events[_id] = event

        # Processing many events condition computation
        bulk_modulo = len(self.bulk_ids) % self.bulk_amount
        elapsed_time = time() - self.last_bulk_insert_date

        # When enough event buffered/time elapsed
        # processing events buffers
        if bulk_modulo == 0 or elapsed_time > self.bulk_delay:

            insert_operations = []
            update_operations = []

            query = {'_id': {'$in': self.bulk_ids}}

            devents = {}

            # Put previous events in pretty data structure
            backend = self.storage.get_backend(self.namespace)
            for devent in backend.find(query):
                devents[devent['_id']] = devent

            # Try to match previous and new incoming event
            for _id in self.incoming_events:
                event = self.incoming_events[_id]
                devent = None
                if _id in devents:
                    devent = devents[_id]
                else:
                    self.logger.info(
                        u'Previous event for rk {} not found'.format(_id))

                # Effective archiver processing call
                operations = self.process_an_event(_id, event, devent)
                for operation in operations:
                    if operation['type'] == 'insert':
                        insert_operations.append(operation)
                    else:
                        update_operations.append(operation)

            self.process_insert_operations(insert_operations)
            self.process_update_operations(update_operations)

            # Buld processing done, reseting informations
            self.bulk_ids = []
            self.incoming_events = {}
            self.last_bulk_insert_date = time()

        # Half useless retro compatibility
        if 'state' in event and event['state']:
            return _id

    def process_an_event(self, _id, event, devent):

        operations = []

        changed = False
        new_event = False

        state = event['state']
        state_type = event['state_type']

        now = int(time())

        event['timestamp'] = event.get('timestamp', now)
        try:
            # Get old record
            exclusion_fields = {
                'perf_data_array',
                'processing'
            }

            if not devent:
                new_event = True
                # may have side effects on acks/cancels
                devent = {}

            old_state = devent['state']
            old_state_type = devent['state_type']
            event['last_state_change'] = devent.get(
                'last_state_change',
                event['timestamp']
            )

            if state != old_state:
                event['previous_state'] = old_state

            if state != old_state or state_type != old_state_type:
                changed = True

            self.check_statuses(event, devent)

        except:
            # No old record
            event['ts_first_stealthy'] = 0
            changed = True
            old_state = state

        if changed:
            # Tests if change is from alert to non alert
            if ('last_state_change' in event
                    and (state == 0 or (state > 0 and old_state == 0))):
                event['previous_state_change_ts'] = event['last_state_change']
            event['last_state_change'] = event.get('timestamp', now)

        if new_event:
            # copy avoid side effects
            operations.append({
                'type': 'insert',
                'event': event.copy(),
                'collection': 'events'
            })
            self.logger.info(u' + New event, have to log {}'.format(_id))

        else:
            change = {}

            # Clear ack related information when event goes ok state
            # Ack is kept in the current event below
            if event['state'] == 0 and devent.get('state') != 0:
                for key in (
                    'ticket_declared_author',
                    'ticket_declared_date',
                    'ticket_date',
                    'ticket'
                ):
                    change[key] = None

            # keep ack information if status does not reset event
            if 'ack' in devent:
                if event['status'] == 0:
                    was_ack = devent.get('ack', {}).get('isAck', False)
                    # save was ack for stats purposes
                    change['ack'] = {
                        'wasAck': was_ack
                    }

                else:
                    change['ack'] = devent['ack']
                    # remove was ack for accurate stats
                    # when event change to any alert state
                    if 'wasAck' in change['ack']:
                        del change['ack']['wasAck']

            # keep cancel information if status does not reset event
            if 'cancel' in devent:
                if event['status'] not in [0, 1]:
                    change['cancel'] = devent['cancel']
                else:
                    change['cancel'] = {}

            # Remove ticket information in case state is back to normal
            # (both ack and ticket declaration case)
            if 'ticket_declared_author' in devent and event['status'] == 0:
                change['ticket_declared_author'] = None
                change['ticket_declared_date'] = None

            # Remove ticket information in case state is back to normal
            # (ticket number declaration only case)
            if 'ticket' in devent and event['status'] == 0:
                del devent['ticket']
                if 'ticket_date' in devent:
                    del devent['ticket_date']

            # Generate diff change from old event to new event
            for key in event:
                if key not in exclusion_fields:
                    if (key in event and
                        key in devent and
                            devent[key] != event[key]):
                        change[key] = event[key]
                    elif key in event and key not in devent:
                        change[key] = event[key]

            # Manage keep state key that allow
            # from UI to keep the choosen state
            # into until next ok state
            event_reset = False

            # When a event is ok again, dismiss keep_state statement
            if devent.get('keep_state') and event['state'] == 0:
                change['keep_state'] = False
                event_reset = True

            # assume we do not just received a keep state and
            # if keep state was sent previously
            # then override state of new event
            if 'keep_state' not in event:
                if not event_reset and devent.get('keep_state'):
                    change['state'] = devent['state']

            # Keep previous output
            if 'keep_state' in event:
                change['change_state_output'] = event['output']
                change['output'] = devent.get('output', '')

            if change:
                operations.append(
                    {
                        'type': 'update',
                        'update': {'$set': change},
                        'query': {'_id': _id},
                        'collection': 'events'
                    }
                )

        # I think that is the right condition to log
        have_to_log = event.get('previous_state', state) != state
        if have_to_log:

            # store ack information to log collection
            if 'ack' in devent:
                event['ack'] = devent['ack']

            self.logger.info(u' + State changed, have to log {}'.format(_id))

            # copy avoid side effects
            operations.append(
                {
                    'type': 'insert',
                    'event': event.copy(),
                    'collection': 'events_log'
                }
            )

        return operations
