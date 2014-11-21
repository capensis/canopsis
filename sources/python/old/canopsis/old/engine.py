#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.old.init import Init
from canopsis.old.rabbitmq import Amqp
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.old.event import forger, get_routingkey

from traceback import format_exc, print_exc
from itertools import cycle
from logging import INFO, DEBUG, FileHandler, Formatter
from time import time, sleep

from os import getpid
from os.path import join
from sys import prefix as sys_prefix


DROP = -1
DISPATCHER_READY_TIME = 5


class Engine(object):
    etype = 'Engine'

    def __init__(
        self,
        next_amqp_queues=[],
        next_balanced=False,
        name="worker1",
        beat_interval=60,
        logging_level=INFO,
        exchange_name='amq.direct',
        routing_keys=[],
        *args, **kwargs
    ):
        super(Engine, self).__init__(*args, **kwargs)

        self.logging_level = logging_level

        self.debug = logging_level == DEBUG

        self.RUN = True

        self.name = name

        self.Amqp = Amqp

        self.amqp_queue = "Engine_{0}".format(self.name)
        self.routing_keys = routing_keys
        self.exchange_name = exchange_name

        self.perfdata_retention = 3600

        self.next_amqp_queues = next_amqp_queues
        self.get_amqp_queue = cycle(self.next_amqp_queues)
        self.original_next_amqp_queues = list(self.next_amqp_queues)

        ## Get from internal or external queue
        self.next_balanced = next_balanced

        init = Init()

        self.logger = init.getLogger(name, logging_level=logging_level)

        logHandler = FileHandler(
            filename=join(
                sys_prefix,
                'var', 'log', 'engines', '%s.log' % name
            )
        )

        logHandler.setFormatter(
            Formatter("%(asctime)s %(levelname)s %(name)s %(message)s")
        )

        # Log in file
        self.logger.addHandler(logHandler)

        self.counter_error = 0
        self.counter_event = 0
        self.counter_worktime = 0

        self.thd_warn_sec_per_evt = 0.6
        self.thd_crit_sec_per_evt = 0.9

        self.beat_interval = beat_interval
        self.beat_last = time()

        self.create_queue = True

        self.send_stats_event = True

        self.rk_on_error = []

        self.last_stat = int(time())

        self.logger.info("Engine initialised")

        self.dispatcher_crecords = [
            'selector',
            'topology',
            'derogation',
            'consolidation',
            'sla',
            'downtime'
        ]

    def crecord_task_complete(self, crecord_id):
        next_ready = time() + DISPATCHER_READY_TIME

        self.storage.update(crecord_id, {
            'loaded': False,
            'next_ready_time': next_ready
        })

        self.logger.debug('next ready time for crecord %s : %s' % (
            crecord_id,
            next_ready
        ))

    def get_ready_record(self, event):
        """
        crecord dispatcher sent an event with type and crecord id.
        So I load a crecord object instance (dynamically) from these infos
        """

        ctype = event.get('crecord_type', None)

        if '_id' in event and ctype and ctype in self.dispatcher_crecords:
            self.logger.warning('Record type %s not found' % ctype)

            record_object = None

            try:
                record_object = self.storage.get(
                    event['_id'],
                    account=Account(user="root", group="root")
                )

            except Exception as e:
                self.logger.critical(
                    'Record <%s> not found in %s : %s' % (
                        ctype,
                        str(self.dispatcher_crecords),
                        e
                    )
                )

            return record_object

    def new_amqp_queue(self, amqp_queue, routing_keys, on_amqp_event, ename):
        self.amqp.add_queue(
            queue_name=amqp_queue,
            routing_keys=routing_keys,
            callback=on_amqp_event,
            exchange_name=ename,
            no_ack=True,
            exclusive=False,
            auto_delete=False
        )

    def pre_run(self):
        pass

    def post_run(self):
        pass

    def run(self):
        def ready():
            self.logger.info(" + Ready!")

        self.logger.info("Start Engine with pid %s" % (getpid()))

        self.amqp = self.Amqp(
            logging_level=self.logging_level,
            logging_name="%s-amqp" % self.name,
            on_ready=ready
        )

        if self.create_queue:
            self.new_amqp_queue(
                self.amqp_queue,
                self.routing_keys,
                self.on_amqp_event,
                self.exchange_name
            )

            # This is an async engine and it needs engine dispatcher
            # bindinds to be feed properly

        if self.etype in self.dispatcher_crecords:
            rk = 'dispatcher.{0}'.format(self.etype)

            self.logger.debug('Creating dispatcher queue : %s' % self.etype)

            self.amqp.get_exchange('media')
            self.new_amqp_queue(
                'Dispatcher_{0}'.format(self.etype),
                rk,
                self.consume_dispatcher,
                'media'
            )

        self.amqp.start()

        self.pre_run()

        while self.RUN:
            # Beat
            if self.beat_interval:
                now = time()

                if now > (self.beat_last + self.beat_interval):
                    self._beat()
                    self.beat_last = now

            try:
                sleep(1)

            except Exception as err:
                self.logger.error("Error in break time: %s" % err)
                self.RUN = False

            except KeyboardInterrupt:
                self.logger.info('Stop request')
                self.RUN = False

        self.post_run()

        self.logger.info("Stop Engine")
        self.stop()
        self.logger.info("End of Engine")

    def on_amqp_event(self, event, msg):
        try:
            self._work(event, msg)

        except Exception as err:
            if event['rk'] not in self.rk_on_error:
                self.logger.error(err)
                self.logger.error("Impossible to deal with: %s" % event)
                self.rk_on_error.append(event['rk'])

            self.next_queue(event)

    def _work(self, event, msg=None, *args, **kargs):
        start = time()
        error = False

        try:
            wevent = self.work(event, msg, *args, **kargs)

            if 'rk' in wevent:
                self.logger.debug('Event rk >' + wevent['rk'])

            if wevent != DROP:
                if isinstance(wevent, dict):
                    event = wevent

                if self.debug:
                    if 'processing' not in event:
                        event['processing'] = {}

                    event['processing'][self.etype] = start

                self.next_queue(event)

        except Exception as err:
            error = True
            self.logger.error("Worker raise exception: %s" % err)
            self.logger.error(format_exc())

        if error:
            self.counter_error += 1

        elapsed = time() - start

        if elapsed > 3:
            self.logger.warning("Elapsed time %.2f seconds" % elapsed)

        self.counter_event += 1
        self.counter_worktime += elapsed

    def work(self, event, amqp_msg):
        return event

    def next_queue(self, event):
        if self.next_balanced:
            queue_name = self.get_amqp_queue.next()
            if queue_name:
                self.amqp.publish(event, queue_name, "amq.direct")

        else:
            for queue_name in self.next_amqp_queues:
                self.amqp.publish(event, queue_name, "amq.direct")

                queues = set(self.original_next_amqp_queues)
                queues = queues.intersection(self.next_amqp_queues)
                lqueues = len(queues)
                lorig = len(self.original_next_amqp_queues)

                if lqueues != lorig:
                    self.next_amqp_queues = [
                        q
                        for q in self.original_next_amqp_queues
                    ]

    def _beat(self):
        now = int(time())

        if self.last_stat + 60 <= now:
            self.logger.debug(" + Send stats")
            self.last_stat = now

            evt_per_sec = 0
            sec_per_evt = 0

            if self.counter_event:
                evt_per_sec = float(self.counter_event) / self.beat_interval
                self.logger.debug(" + %0.2f event(s)/seconds" % evt_per_sec)

            if self.counter_worktime and self.counter_event:
                sec_per_evt = self.counter_worktime / self.counter_event
                self.logger.debug(" + %0.5f seconds/event" % sec_per_evt)

            ## Submit event
            if self.send_stats_event and self.counter_event != 0:
                state = 0

                if sec_per_evt > self.thd_warn_sec_per_evt:
                    state = 1

                if sec_per_evt > self.thd_crit_sec_per_evt:
                    state = 2

                perf_data_array = [{
                    'retention': self.perfdata_retention,
                    'metric': 'cps_evt_per_sec',
                    'value': round(evt_per_sec, 2),
                    'unit': 'evt'
                },
                {
                    'retention': self.perfdata_retention,
                    'metric': 'cps_sec_per_evt',
                    'value': round(sec_per_evt, 5),
                    'unit': 's',
                    'warn': self.thd_warn_sec_per_evt,
                    'crit': self.thd_crit_sec_per_evt
                }]

                self.logger.debug(" + State: %s" % state)

                msg = "%0.2f evt/sec, %0.5f sec/evt" % (
                    evt_per_sec,
                    sec_per_evt
                )

                event = forger(
                    connector="Engine",
                    connector_name="engine",
                    event_type="check",
                    source_type="resource",
                    resource=self.amqp_queue,
                    state=state,
                    state_type=1,
                    output=msg,
                    perf_data_array=perf_data_array
                )

                rk = get_routingkey(event)
                self.amqp.publish(event, rk, self.amqp.exchange_name_events)

            self.counter_error = 0
            self.counter_event = 0
            self.counter_worktime = 0

        try:
            self.beat()

        except Exception as err:
            self.logger.error("Beat raise exception: %s" % err)
            self.logger.error(print_exc())

        finally:
            self.beat_lock = False

    def beat(self):
        pass

    def stop(self):
        self.RUN = False

        # cancel self consumer
        self.amqp.cancel_queues()

        self.amqp.stop()
        self.amqp.join()
        self.logger.debug(" + Stopped")

    class Lock(object):
        def __init__(self, engine, name, *args, **kwargs):
            super(Engine.Lock, self).__init__(*args, **kwargs)

            self.name = name
            self.lock_id = '{0}.{1}'.format(engine.etype, name)

            self.storage = get_storage(
                namespace='lock',
                logging_level=engine.logging_level,
                account=Account(user='root', group='root')
            ).get_backend()

            self.engine = engine
            self.lock = {}

        def own(self):
            now = time()
            last = self.lock.get('t', now)
            beat = self.engine.beat_interval

            if self.lock.get('l', False) and (now - last) < beat:
                self.engine.logger.debug(
                    'Another engine {} is already holding the lock {}'.format(
                        self.engine.etype,
                        self.name
                    )
                )

                return False

            else:
                self.engine.logger.debug(
                    'Lock {} on engine {}, processing...'.format(
                        self.name,
                        self.engine.etype
                    )
                )

                return True

            return False

        def __enter__(self):
            self.lock = self.storage.find_and_modify(
                query={'_id': self.lock_id},
                update={'$set': {'l': True}},
                upsert=True
            )

            if 't' not in self.lock:
                self.lock['t'] = 0

            return self

        def __exit__(self, type, value, tb):
            if self.own():
                self.engine.logger.debug('Release lock {} on engine {}'.format(
                    self.name,
                    self.engine.etype
                ))

                self.storage.save({
                    '_id': self.lock_id,
                    'l': False,
                    't': time()
                })
