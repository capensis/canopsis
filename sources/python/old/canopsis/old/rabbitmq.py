#!/usr/bin/env python
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
from __future__ import unicode_literals

from sys import stdout
from sys import prefix as sys_prefix

from kombu import Connection, Exchange, Queue, pools, __version__

import traceback

try:
    from amqplib.client_0_8.exceptions import AMQPConnectionException \
        as ConnectionError
except ImportError as IE:
    from amqp.exceptions import ConnectionError

from socket import error, timeout

from time import sleep
from logging import INFO, getLogger
from threading import Thread
from os.path import join
from traceback import print_exc

# Number of tries to re-publish an event before it is lost
# when connection problems


class Amqp(Thread):
    def __init__(
        self,
        host="localhost",
        port=5672,
        userid="guest",
        password="guest",
        virtual_host="canopsis",
        exchange_name="canopsis",
        logging_name="Amqp",
        logging_level=INFO,
        read_config_file=True,
        auto_connect=True,
        on_ready=None,
        max_retries=5
    ):
        super(Amqp, self).__init__()

        self.logger = getLogger(logging_name)

        self.host = host
        self.port = port
        self.userid = userid
        self.password = password
        self.virtual_host = virtual_host
        self.exchange_name = exchange_name
        self.logging_level = logging_level

        if (read_config_file):
            self.read_config("amqp")

        self.amqp_uri = "amqp://{0}:{1}@{2}:{3}/{4}".format(
            self.userid, self.password, self.host, self.port,
            self.virtual_host)

        self.logger.setLevel(logging_level)

        # Event sent try count before event drop in case of connection problem
        if max_retries != 5:
            self.logger.info(u'Custom retries value : {} {}'.format(
                max_retries,
                type(max_retries)
            ))
        self.max_retries = max_retries

        self.exchange_name_events = exchange_name + ".events"
        self.exchange_name_alerts = exchange_name + ".alerts"
        self.exchange_name_incidents = exchange_name + ".incidents"

        self.chan = None
        self.conn = None
        self.connected = False
        self.on_ready = on_ready

        self.RUN = True

        self.exchanges = {}

        self.queues = {}

        self.paused = False

        self.connection_errors = (
            ConnectionError,
            error,
            IOError,
            OSError)

        # Create exchange
        self.logger.debug("Create exchanges object")
        for exchange_name in [
                self.exchange_name, self.exchange_name_events,
                self.exchange_name_alerts, self.exchange_name_incidents]:
            self.logger.debug(u' + {0}'.format(exchange_name))
            self.get_exchange(exchange_name)

        if auto_connect:
            self.connect()

        self.logger.debug("Object canamqp initialized")

    def run(self):
        self.logger.debug("Start thread ...")
        reconnect = False

        while self.RUN:

            self.connect()

            if self.connected:
                self.init_queue(reconnect=reconnect)

                self.logger.debug("Drain events ...")
                while self.RUN:
                    try:
                        if not self.paused:
                            self.conn.drain_events(timeout=0.5)
                        else:
                            sleep(0.5)

                    except timeout:
                        pass

                    except self.connection_errors as err:
                        self.logger.error(
                            u"Connection error ! ({})".format(err)
                        )
                        break

                    except Exception as err:
                        self.logger.error(u"Error: {} ({})".format(
                            err,
                            type(err)
                        ))
                        print_exc(file=stdout)
                        break

                self.disconnect()

            if self.RUN:
                self.logger.error(
                    'Connection lost, try to reconnect in few seconds ...'
                )
                reconnect = True
                self.wait_connection(timeout=5)

        self.logger.debug("End of thread ...")

    def stop(self):
        self.logger.debug("Stop thread ...")
        self.RUN = False

    def connect(self):
        if not self.connected:
            self.logger.info(
                "Connect to AMQP Broker (%s:%s)" % (self.host, self.port))

            self.conn = Connection(self.amqp_uri)

            try:
                self.logger.debug(" + Connect")
                self.conn.connect()
                self.logger.info("Connected to AMQP Broker.")
                self.producers = pools.Producers(limit=10)
                self.connected = True
            except Exception as err:
                self.conn.release()
                self.logger.error(u"Impossible to connect ({})".format(err))

            if self.connected:
                self.logger.debug(" + Open channel")
                try:
                    self.chan = self.conn.channel()

                    self.logger.debug(
                        "Channel openned. Ready to send messages")

                    try:
                        # Declare exchange
                        self.logger.debug("Declare exchanges")
                        for exchange_name in self.exchanges:
                            self.logger.debug(u" + {}".format(exchange_name))
                            self.exchanges[exchange_name](self.chan).declare()
                    except Exception as err:
                        self.logger.error(
                            u"Impossible to declare exchange ({})".format(err))

                except Exception as err:
                    self.logger.error(err)
        else:
            self.logger.debug("Already connected")

    def get_exchange(self, name):
        if name:
            try:
                return self.exchanges[name]
            except:
                if name == "amq.direct":
                    self.exchanges[name] = Exchange(
                        name, "direct", durable=True)
                else:
                    self.exchanges[name] = Exchange(
                        name, "topic", durable=True, auto_delete=False)
                return self.exchanges[name]
        else:
            return None

    def init_queue(self, reconnect=False):
        if self.queues:
            self.logger.debug("Init queues")
            for queue_name in self.queues.keys():
                self.logger.debug(u" + {}".format(queue_name))
                qsettings = self.queues[queue_name]

                if not qsettings['queue']:
                    self.logger.debug("   + Create queue")

                    # copy list
                    routing_keys = list(qsettings['routing_keys'])
                    routing_key = None

                    if len(routing_keys):
                        routing_key = routing_keys[0]
                        routing_keys = routing_keys[1:]

                    exchange = self.get_exchange(qsettings['exchange_name'])

                    if qsettings['exchange_name'] == "amq.direct" \
                            and not routing_key:
                        routing_key = queue_name

                    self.logger.debug(
                        (u"exchange: '{}', exclusive: {}," +
                            " auto_delete: {},no_ack: {}")
                        .format(
                            qsettings['exchange_name'],
                            qsettings['exclusive'],
                            qsettings['auto_delete'],
                            qsettings['no_ack'])
                    )
                    qsettings['queue'] = Queue(
                        queue_name,
                        exchange=exchange,
                        routing_key=routing_key,
                        exclusive=qsettings['exclusive'],
                        auto_delete=qsettings['auto_delete'],
                        no_ack=qsettings['no_ack'],
                        channel=self.conn.channel())

                    qsettings['queue'].declare()

                    if len(routing_keys):
                        self.logger.debug(" + Bind on all routing keys")
                        for routing_key in routing_keys:
                            self.logger.debug(
                                " + routing_key: '{}'".format(routing_key)
                            )
                            try:
                                qsettings['queue'].bind_to(
                                    exchange=exchange,
                                    routing_key=routing_key
                                )
                            except:
                                self.logger.error(
                                    u"You need upgrade your Kombu version ({})"
                                    .format(__version__)
                                )

                if not qsettings['consumer'] or reconnect:
                    self.logger.debug("   + Create Consumer")
                    qsettings['consumer'] = self.conn.Consumer(
                        qsettings['queue'], callbacks=[qsettings['callback']])

                self.logger.debug("   + Consume queue")
                qsettings['consumer'].consume()

            if self.on_ready:
                self.on_ready()

    def add_queue(
        self,
        queue_name,
        routing_keys,
        callback,
        exchange_name=None,
        no_ack=True,
        exclusive=False,
        auto_delete=True
    ):

        c_routing_keys = []

        if not isinstance(routing_keys, list):
            if isinstance(routing_keys, basestring):
                c_routing_keys = [routing_keys]
        else:
            c_routing_keys = routing_keys

        if not exchange_name:
            exchange_name = self.exchange_name

        self.queues[queue_name] = {
            'queue': False,
            'consumer': False,
            'queue_name': queue_name,
            'routing_keys': c_routing_keys,
            'callback': callback,
            'exchange_name': exchange_name,
            'no_ack': no_ack,
            'exclusive': exclusive,
            'auto_delete': auto_delete
        }

    def publish(
        self,
        msg,
        routing_key,
        exchange_name=None,
        serializer="json",
        compression=None,
        content_type=None,
        content_encoding=None
    ):

        operation_success = False
        retries = 0

        while not operation_success and retries < self.max_retries:
            retries += 1

            if self.connected:

                if not exchange_name:
                    exchange_name = self.exchange_name

                with self.producers[self.conn].acquire(block=True) as producer:
                    try:
                        _msg = msg.copy()

                        Amqp._clean_msg_for_serialization(_msg)

                        producer.publish(
                            _msg,
                            serializer=serializer,
                            compression=compression,
                            routing_key=routing_key,
                            exchange=self.get_exchange(exchange_name.encode('utf-8'))
                        )

                        self.logger.debug('publish {} in exchange {}'.format(
                            routing_key,
                            exchange_name
                        ))

                        operation_success = True

                    except Exception as e:
                        self.logger.error(
                            u' + Impossible to send {}'.format(traceback.format_exc(e))
                        )
                        self.disconnect()
                        self.connect()
            else:
                self.logger.error('Not connected ... try reconnecting')
                self.connect()

            if not operation_success:
                # Event and it's information are buffered until next send retry
                self.logger.info(u'Retry count {}'.format(
                    retries
                ))

        if not operation_success:
            # Event and it's information are buffered until next send retry
            self.logger.error(u'Too much retries for event {}, give up'.format(
                routing_key
            ))

    @staticmethod
    def _clean_msg_for_serialization(msg):
        from bson import objectid
        for key in msg:
            if isinstance(msg[key], objectid.ObjectId):
                msg[key] = str(msg[key])

    def cancel_queues(self):
        if self.connected:
            for queue_name in self.queues.keys():
                if self.queues[queue_name]['consumer']:
                    self.logger.debug(
                        u" + Cancel consumer on {}".format(queue_name)
                    )
                    try:
                        self.queues[queue_name]['consumer'].cancel()
                    except:
                        pass

                    del(self.queues[queue_name]['consumer'])
                    self.queues[queue_name]['consumer'] = False
                    del(self.queues[queue_name]['queue'])
                    self.queues[queue_name]['queue'] = False

    def disconnect(self):
        if self.connected:
            self.logger.info("Disconnect from AMQP Broker")

            self.cancel_queues()

            for exchange in self.exchanges:
                del exchange
            self.exchanges = {}

            try:
                pools.reset()
            except Exception as err:
                self.logger.error(
                    u"Impossible to reset kombu pools: {} ({})".format(
                        err, type(err)))

            try:
                self.conn.release()
                del self.conn
            except Exception as err:
                self.logger.error(
                    u"Impossible to release connection: {} ({})".format(
                        err, type(err)))

            self.connected = False

    def wait_connection(self, timeout=5):
        i = 0
        while self.RUN and not self.connected and i < (timeout * 2):
            try:
                sleep(0.5)
            except:
                pass
            i += 1

    def read_config(self, name):

        filename = join(sys_prefix, 'etc', u'{0}.conf'.format(name))

        import ConfigParser
        self.config = ConfigParser.RawConfigParser()

        try:
            self.config.read(filename)

            section = 'master'

            self.host = self.config.get(section, "host")
            self.port = self.config.getint(section, "port")
            self.userid = self.config.get(section, "userid")
            self.password = self.config.get(section, "password")
            self.virtual_host = self.config.get(section, "virtual_host")
            self.exchange_name = self.config.get(section, "exchange_name")

        except Exception as err:
            self.logger.error(
                u"Can't to load configurations ({}), use default ...".format(
                    err
                ))

    def __del__(self):
        self.stop()
