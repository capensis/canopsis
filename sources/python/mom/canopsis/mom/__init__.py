# -*- coding: utf-8 -*-
#--------------------------------
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

__version__ = "0.1"


from canopsis.middleware import Middleware
from canopsis.configuration import Category, Parameter

from sys import stdout

from kombu import Connection, Exchange, Queue, pools, __version__

from socket import timeout

from time import sleep
from traceback import print_exc


class MOM(Middleware):

    CONF_FILE = 'mom/mom.conf'

    CATEGORY = 'PUBSUB'

    SENDERS = 'senders'
    RECEIVERS = 'receivers'

    def _get_conf_files(self, *args, **kwargs):

        result = super(MOM, self)._get_conf_files(*args, **kwargs)

        result.append(MOM.CONF_FILE)

        return result

    def _conf(self, *args, **kwargs):

        result = super(MOM, self)._conf(*args, **kwargs)

        result += Category(
            Parameter(MOM.SENDERS, self.senders),
            Parameter(MOM.RECEIVERS, self.receivers))

        return result

    def _configure(self, unified_conf, *args, **kwargs):

        super(MOM, self)._configure(unified_conf)

        self._update_property(
            unified_conf=unified_conf, param_name=MOM.SENDERS)
        self._update_property(
            unified_conf=unified_conf, param_name=MOM.RECEIVERS)

    def __init__(
        self, senders=None, receivers=None
    ):
        super(MOM, self).__init__()

        self.receivers = receivers
        self.senders = senders

    def run(self):
        self.logger.debug("Start thread ...")
        reconnect = False

        while self.RUN:

            self.connect()

            #self.wait_connection()

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
                        self.logger.error("Connection error ! (%s)" % err)
                        break

                    except Exception as err:
                        self.logger.error(
                            "Unknown error: %s (%s)" % (err, type(err)))
                        print_exc(file=stdout)
                        break

                self.disconnect()

            if self.RUN:
                self.logger.error(
                    "Connection loss, try to reconnect in few seconds ...")
                reconnect = True
                self.wait_connection(timeout=5)

        self.logger.debug("End of thread ...")

    def stop(self):

        raise NotImplementedError()

    def connect(self):
        if not self.connected:
            self.logger.info(
                "Connect to AMQP Broker (%s:%s)" % (self.host, self.port))

            self.conn = Connection(self.amqp_uri)

            try:
                self.logger.debug(" + Connect")
                self.conn.connect()
                self.logger.info("Connected to AMQP Broker.")
                self.senders = pools.Producers(limit=10)
                self.connected = True
            except Exception as err:
                self.conn.release()
                self.logger.error("Impossible to connect (%s)" % err)

            if self.connected:
                self.logger.debug(" + Open channel")
                try:
                    self.chan = self.conn.channel()

                    self.logger.debug(
                        "Channel openned. Ready to send messages")

                    try:
                        ## declare exchange
                        self.logger.debug("Declare exchanges")
                        for exchange_name in self.exchanges:
                            self.logger.debug(" + %s" % exchange_name)
                            self.exchanges[exchange_name](self.chan).declare()
                    except Exception as err:
                        self.logger.error(
                            "Impossible to declare exchange (%s)" % err)

                except Exception as err:
                    self.logger.error(err)
        else:
            self.logger.debug("Allready connected")

    def get_sender(self, name):
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
                self.logger.debug(" + %s" % queue_name)
                qsettings = self.queues[queue_name]

                if not qsettings['queue']:
                    self.logger.debug("   + Create queue")

                    # copy list
                    routing_keys = list(qsettings['routing_keys'])
                    routing_key = None

                    if len(routing_keys):
                        routing_key = routing_keys[0]
                        routing_keys = routing_keys[1:]

                    exchange = self.get_sender(qsettings['exchange_name'])

                    if qsettings['exchange_name'] == "amq.direct" \
                            and not routing_key:
                        routing_key = queue_name

                    #self.logger.debug("   + exchange: '%s', routing_key: '%s', exclusive: %s, auto_delete: %s, no_ack: %s" % (qsettings['exchange_name'], routing_key, qsettings['exclusive'], qsettings['auto_delete'], qsettings['no_ack']))
                    self.logger.debug(
                        "   + exchange: '%s', exclusive: %s, auto_delete: %s, no_ack: %s" % (qsettings['exchange_name'], qsettings['exclusive'], qsettings['auto_delete'], qsettings['no_ack']))
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
                                " + routing_key: '%s'" % routing_key)
                            try:
                                qsettings['queue'].bind_to(
                                    exchange=exchange, routing_key=routing_key)
                            except:
                                self.logger.error(
                                    "You need upgrade your Kombu version (%s)" % __version__)

                if not qsettings['consumer']:
                    self.logger.debug("   + Create Consumer")
                    qsettings['consumer'] = self.conn.Consumer(
                        qsettings['queue'], callbacks=[qsettings['callback']])

                self.logger.debug("   + Consume queue")
                qsettings['consumer'].consume()

    def add_queue(
        self, queue_name, routing_keys, callback, exchange_name=None,
        no_ack=True, exclusive=False, auto_delete=True
    ):
        #if exchange_name == "amq.direct":
        #   routing_keys = queue_name

        c_routing_keys = []

        if not isinstance(routing_keys, list):
            if isinstance(routing_keys, str):
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

    def _publish(
        self, sender, msg, rk, serializer='json', compression=None,
        content_type=None, content_encoding=None
    ):

        raise NotImplementedError()

    def publish(
        self, msg, rk, sender=None, serializer="json",
        compression=None, content_type=None, content_encoding=None
    ):
        if self.connected():

            senders = self.senders if sender is None else (sender)

            for sender in senders:
                _sender = self.get_sender(sender)
                self.logger.debug(
                    "Send message to %s in %s" % (rk, sender))
                try:
                    self._publish(
                        sender=_sender, rk=rk, serializer=serializer,
                        compression=compression, content_type=content_type,
                        content_encoding=content_encoding)
                except Exception as err:
                    self.logger.error("Impossible to send message (%s)" % err)
                else:
                    self.logger.debug(
                        "Message sent message to %s in %s" % (rk, sender))

        else:
            self.logger.error("%s is not connected" % self)

    def cancel_senders(self, sender=None):

        raise NotImplementedError()

    def cancel_receivers(self):

        if self.connected:
            for queue_name in self.queues.keys():
                if self.queues[queue_name]['consumer']:
                    self.logger.debug(" + Cancel consumer on %s" % queue_name)
                    try:
                        self.queues[queue_name]['consumer'].cancel()
                    except:
                        pass

                    del(self.queues[queue_name]['consumer'])
                    self.queues[queue_name]['consumer'] = False
                    del(self.queues[queue_name]['queue'])
                    self.queues[queue_name]['queue'] = False
