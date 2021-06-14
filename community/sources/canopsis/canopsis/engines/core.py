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

from traceback import format_exc, print_exc

from itertools import cycle

from logging import INFO, DEBUG, FileHandler, Formatter

from time import time, sleep
from json import loads
from os import getpid
from os.path import join

from canopsis.common import root_path
from canopsis.common.amqp import AmqpPublisher, DIRECT_EXCHANGE_NAME
from canopsis.common.amqp import get_default_connection as \
    get_default_amqp_connection
from canopsis.common.init import Init
from canopsis.old.rabbitmq import Amqp
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.event import forger, get_routingkey
from canopsis.task.core import register_task
from canopsis.tools import schema as cschema

DROP = -1


class Engine(object):
    etype = 'Engine'

    def __init__(self,
                 next_amqp_queues=[],
                 next_balanced=False,
                 name="worker1",
                 beat_interval=60,
                 logging_level=INFO,
                 exchange_name=DIRECT_EXCHANGE_NAME,
                 routing_keys=[],
                 camqp_custom=None,
                 max_retries=5,
                 *args, **kwargs):
        super(Engine, self).__init__()

        self.logging_level = logging_level
        self.debug = logging_level == DEBUG

        init = Init()

        self.logger = init.getLogger(name, logging_level=self.logging_level)

        log_handler = FileHandler(
            filename=join(
                root_path, 'var', 'log', 'engines', '{0}.log'.format(name)
            )
        )

        log_handler.setFormatter(
            Formatter(
                "%(asctime)s %(levelname)s %(name)s %(message)s"
            )
        )

        # Log in file
        self.logger.addHandler(log_handler)

        self.RUN = True

        self.name = name

        # Set parametrized Amqp for testing purposes
        if camqp_custom is None:
            self.Amqp = Amqp
        else:
            self.Amqp = camqp_custom

        # self.amqp handles the consumption of events from rabbitmq. The
        # publication of events from self.amqp is deprecated.
        self.amqp = None
        # self.beat_amqp_publisher and self.work_amqp_publisher handle the
        # publication of events (they are separated to prevent sharing a
        # channel between two threads).
        self.beat_amqp_publisher = AmqpPublisher(
            get_default_amqp_connection(), self.logger)
        self.work_amqp_publisher = AmqpPublisher(
            get_default_amqp_connection(), self.logger)

        self.amqp_queue = "Engine_{0}".format(self.name)
        self.routing_keys = routing_keys
        self.exchange_name = exchange_name

        self.perfdata_retention = 3600

        self.next_amqp_queues = next_amqp_queues
        self.get_amqp_queue = cycle(self.next_amqp_queues)

        # Get from internal or external queue
        self.next_balanced = next_balanced

        self.max_retries = max_retries

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

        self.logger.info("Engine initialized")

    def new_amqp_queue(self,
                       amqp_queue,
                       routing_keys,
                       on_amqp_event,
                       exchange_name):
        self.amqp.add_queue(
            queue_name=amqp_queue,
            routing_keys=routing_keys,
            callback=on_amqp_event,
            exchange_name=exchange_name,
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

        self.logger.info("Start Engine with pid {0}".format(getpid()))

        self.amqp = self.Amqp(
            logging_level=self.logging_level,
            logging_name="{0}-amqp".format(self.name),
            on_ready=ready,
            max_retries=self.max_retries
        )

        if self.create_queue:
            self.new_amqp_queue(
                self.amqp_queue,
                self.routing_keys,
                self.on_amqp_event,
                self.exchange_name
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
                self.logger.error("Error in break time: {0}".format(err))
                self.RUN = False

            except KeyboardInterrupt:
                self.logger.info(u'Stop request')
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
                self.logger.error("Impossible to deal with: {0}".format(event))
                self.rk_on_error.append(event['rk'])

            self.next_queue(event)

    def _work(self, event, msg=None, *args, **kargs):
        start = time()
        error = False

        try:
            if self.debug:
                if 'processing' not in event:
                    event['processing'] = {}

                event['processing'][self.etype] = start

            wevent = self.work(event, msg, *args, **kargs)

            if wevent != DROP:
                if isinstance(wevent, dict):
                    event = wevent

                self.next_queue(event)

        except Exception as err:
            error = True
            self.logger.error("Worker raise exception: {0}".format(err))
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
                try:
                    self.work_amqp_publisher.direct_event(event, queue_name)
                except Exception as e:
                    self.logger.exception("Unable to send event to next queue")

        else:
            for queue_name in self.next_amqp_queues:
                try:
                    self.work_amqp_publisher.direct_event(event, queue_name)
                except Exception as e:
                    self.logger.exception("Unable to send event to next queue")

    def _beat(self):
        now = int(time())

        if self.last_stat + 60 <= now:
            self.logger.debug(" + Send stats")
            self.last_stat = now

            evt_per_sec = 0
            sec_per_evt = 0

            if self.counter_event:
                evt_per_sec = float(self.counter_event) / self.beat_interval
                self.logger.debug(" + %0.2f event(s)/seconds", evt_per_sec)

            if self.counter_worktime and self.counter_event:
                sec_per_evt = self.counter_worktime / self.counter_event
                self.logger.debug(" + %0.5f seconds/event", sec_per_evt)

            # Submit event
            if self.send_stats_event and self.counter_event != 0:
                state = 0

                if sec_per_evt > self.thd_warn_sec_per_evt:
                    state = 1

                if sec_per_evt > self.thd_crit_sec_per_evt:
                    state = 2

                perf_data_array = [
                    {
                        'retention': self.perfdata_retention,
                        'metric': 'cps_evt_per_sec',
                        'value': round(evt_per_sec, 2),
                        'unit': 'evt'
                    }, {
                        'retention': self.perfdata_retention,
                        'metric': 'cps_sec_per_evt',
                        'value': round(sec_per_evt, 5),
                        'unit': 's',
                        'warn': self.thd_warn_sec_per_evt,
                        'crit': self.thd_crit_sec_per_evt
                    }
                ]

                self.logger.debug(" + State: {0}".format(state))

                event = forger(
                    connector="Engine",
                    connector_name="engine",
                    event_type="check",
                    source_type="resource",
                    resource=self.amqp_queue,
                    state=state,
                    state_type=1,
                    output="%0.2f evt/sec, %0.5f sec/evt" % (
                        evt_per_sec, sec_per_evt),
                    perf_data_array=perf_data_array
                )

                try:
                    self.beat_amqp_publisher.canopsis_event(event)
                except Exception as e:
                    self.logger.exception("Unable to send perfdata")

            self.counter_error = 0
            self.counter_event = 0
            self.counter_worktime = 0

        try:
            self.beat()

        except Exception as err:
            self.logger.error("Beat raise exception: {0}".format(err))
            self.logger.error(print_exc())

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
            super(Engine.Lock, self).__init__()

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

            if self.lock.get('l', False) \
                    and (now - last) < self.engine.beat_interval:
                self.engine.logger.debug(
                    'Another engine {0} is already holding the lock {1}'.
                    format(self.engine.etype, self.name))

                return False

            else:
                self.engine.logger.debug(
                    'Lock {1} on engine {0}, processing...'.format(
                        self.engine.etype, self.name))
                return True

            return False

        def __enter__(self):
            lock = self.storage.find_and_modify(
                query={'_id': self.lock_id},
                update={'$set': {'l': True}},
                upsert=True
            )

            if lock is not None:
                self.lock = lock

            if 't' not in self.lock:
                self.lock['t'] = 0

            return self

        def __exit__(self, type, value, tb):
            if self.own():
                self.engine.logger.debug(
                    'Release lock {1} on engine {0}'.format(
                        self.engine.etype, self.name))

                self.storage.save({
                    '_id': self.lock_id,
                    'l': False,
                    't': time()
                })

        @classmethod
        def release(cls, lock_id, object_storage):
            object_storage.update(
                {'_id': lock_id},
                {'$set': {
                    'l': False,
                    't': time()
                }},
                upsert=True
            )


class TaskHandler(Engine):
    etype = 'Task'

    def __init__(self, *args, **kwargs):
        super(TaskHandler, self).__init__(*args, **kwargs)
        self.amqp_queue = self.name

    def work(self, msg, *args, **kwargs):
        self.logger.info(u'Received job: {0}'.format(msg))

        start = int(time())

        job = None
        output = None
        state = 3

        try:
            if not isinstance(msg, dict):
                job = loads(msg)

            else:
                job = msg

        except ValueError as err:
            output = 'Impossible to decode message: {0}'.format(err)
            self.logger.warning(output)
            state = 2

        else:
            if not cschema.validate(job, 'task.{0}'.format(self.etype)):
                output = 'Invalid job'
                self.logger.warning(output)
                state = 2

            else:
                try:
                    state, output = self.handle_task(job)

                except NotImplementedError:
                    state = 1
                    output = 'Not implemented'

        end = int(time())

        event = {
            'timestamp': end,
            'connector': 'taskhandler',
            'connector_name': self.name,
            'event_type': 'check',
            'source_type': 'resource',
            'component': 'job',
            'resource': job['jobid'],
            'state': state,
            'state_type': 1,
            'output': output,
            'execution_time': end - start
        }

        try:
            self.work_amqp_publisher.canopsis_event(event)
        except Exception as e:
            self.logger.exception("Unable to send event to next queue")

    def handle_task(self, job):
        """
            :param job: Job's informations
            :type job: dict

            :returns: (<state>, <output>)
        """

        raise NotImplementedError()


@register_task
def publish(event, publisher, rk=None, exchange=None, logger=None, **kwargs):
    """Task dedicated to publish an event from an engine.

    :param dict event to send.
    :param publisher: object in charge of publishing the event. Its method
        ``publish`` takes three parameters, the ``event``, the ``rk`` related
        to the event and an ``exchange name``.

    :param str rk: routing key to use. If None, use get_routingkey(event).
    :param str exchange: exchange name. If None, use
        ``publisher.exchange_name_events``.
    """

    if exchange is None:
        exchange = publisher.exchange_name_events

    if rk is None:
        rk = get_routingkey(event)

    publisher.publish(
        event, rk, exchange
    )
