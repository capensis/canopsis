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

from canopsis.engines.core import Engine, publish
from canopsis.check.archiver import Archiver, BAGOT, STEALTHY
from canopsis.context_graph.manager import ContextGraph
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.old.storage import CONFIG
from copy import deepcopy
from csv import reader
from time import time


class engine(Engine):
    etype = 'eventstore'

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        self.archiver = Archiver(
            namespace='events', confnamespace='object',
            autolog=False, log_lvl=self.logging_level
        )

        self.event_types = reader([CONFIG.get('events', 'types')]).next()
        self.check_types = reader([CONFIG.get('events', 'checks')]).next()
        self.log_types = reader([CONFIG.get('events', 'logs')]).next()
        self.comment_types = reader([CONFIG.get('events', 'comments')]).next()

        self.context = ContextGraph()
        self.pbehavior = PBehaviorManager()
        self.beat()

        self.log_bulk_amount = 100
        self.log_bulk_delay = 3
        self.last_bulk_insert_date = time()
        self.events_log_buffer = []

    def beat(self):
        self.archiver.beat()

        with self.Lock(self, 'eventstore_reset_status') as l:
            if l.own():
                self.reset_stealthy_event_duration = time()
                self.archiver.reload_configuration()
                self.archiver.reset_status_event(BAGOT)
                self.archiver.reset_status_event(STEALTHY)

    def store_check(self, event):
        _id = self.archiver.check_event(event['rk'], event)

        if _id:
            event['_id'] = _id
            event['event_id'] = event['rk']
            # Event to Alert
            publish(
                publisher=self.amqp, event=event, rk=event['rk'],
                exchange=self.amqp.exchange_name_alerts
            )

    def store_log(self, event, store_new_event=True):

        """
            Stores events in events_log collection
            Logged events are no more in event collection at the moment
        """

        # Ensure event Id exists from rk key
        event['_id'] = event['rk']

        # Prepare log event collection async insert
        log_event = deepcopy(event)
        self.events_log_buffer.append({
            'event': log_event,
            'collection': 'events_log'
        })

        bulk_modulo = len(self.events_log_buffer) % self.log_bulk_amount
        elapsed_time = time() - self.last_bulk_insert_date

        if bulk_modulo == 0 or elapsed_time > self.log_bulk_delay:
            self.archiver.process_insert_operations(
                self.events_log_buffer
            )
            self.events_log_buffer = []
            self.last_bulk_insert_date = time()

        # Event to Alert
        event['event_id'] = event['rk']
        publish(
            publisher=self.amqp, event=event, rk=event['rk'],
            exchange=self.amqp.exchange_name_alerts
        )

    def work(self, event, *args, **kargs):

        if 'exchange' in event:
            del event['exchange']

        event_type = event['event_type']

        if event_type not in self.event_types:
            self.logger.warning(
                "Unknown event type '{}', id: '{}', event:\n{}".format(
                    event_type,
                    event['rk'],
                    event
                ))
            return event

        elif event_type in self.check_types:
            self.store_check(event)

        elif event_type in self.log_types:
            self.store_log(event)

        elif event_type in self.comment_types:
            self.store_log(event, store_new_event=False)

        return event
