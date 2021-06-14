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

__version__ = "0.1"

from kombu import \
    Connection, Exchange, Queue, Producer, Consumer

from urlparse import urlparse

from canopsis.mom import MOM


class Kombu(MOM):

    __register__ = True  #: register this class to middleware classes
    __protocol__ = 'kombu'  #: register this class to the protocol kombu

    CONF_RESOURCE = 'kombu/kombu.conf'

    class Error(Exception):
        """
        Handle Kombu Exceptions
        """

    def __init__(
        self, port=5672, *args, **kwargs
    ):
        super(Kombu, self).__init__(**kwargs.update({'port': port}))

        self.RUN = True

    def _connect(self):

        result = None

        if not self.connected():

            self.logger.info(
                "Connect to RabbitMQ Broker %s" % (self.uri))

            uri = self.uri

            parsed_url = urlparse(uri)

            uri = 'amqp://%s/%s' % (parsed_url.netloc, parsed_url.path)

            result = Connection(uri)

            try:
                result.connect()
                self.logger.info("Connected to AMQP Broker.")

            except Exception as e:
                result.release()
                self.logger.error(
                    "Connection failure to %s: %s" % (uri, e))

        else:
            self.logger.debug("Allready connected")

        return result

    def connected(self, *args, **kwargs):

        return self.conn is not None and self._conn.connected

    def _get_sender(self, sender, mode):

        exchange = Exchange(name=sender, type=mode, durable=True)

        result = Producer(channel=self.conn, exchange=exchange)

        return result

    def _get_receiver(self, receiver, callback):

        queue = Queue(name=receiver)

        result = Consumer(
            self.conn, queues=queue, no_ack=(not self.safe),
            callbacks=(lambda x, y: callback(x),))

        result.consume()

        return result

    def _send(
        self,
        sender, msg, rk, serializer, compression,
        content_type, content_encoding, out_timeout
    ):

        sender.publish(
            msg, routing_key=rk, serializer=serializer,
            compression=compression, content_type=content_type,
            content_encoding=content_encoding)

    def _receive(self, receiver, callback, in_timeout):

        if in_timeout >= 0:
            self.conn.default_channel.basic_get(
                queue=receiver.name, no_ack=(not self.safe))

        if callback is not None:
            receiver.consume(no_ack=(not self.safe))
