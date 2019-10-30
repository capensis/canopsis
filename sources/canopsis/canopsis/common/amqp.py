# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

import configparser
import json
import os
import time

import pika

from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.event import get_routingkey

DIRECT_EXCHANGE_NAME = 'amq.direct'
DEFAULT_CONF_FILE = "etc/amqp.conf"


class AmqpPublishError(Exception):
    pass


class AmqpConnection(object):

    def __init__(self, url):
        """
        :param url: url of the form: amqp://[<user>:<pass>]@host:port/vhost
        :type url: str
        """

        self._url = url
        self._connection = None
        self._channel = None

    def __enter__(self):
        self.connect()
        return self

    def __exit__(self, type_, value, traceback):
        self.disconnect()

    @property
    def channel(self):
        """
        If no channel is declared, try to reconnect to the bus.
        """
        if self._channel is None:
            self.connect()

        return self._channel

    @property
    def connection(self):
        if self._connection is None:
            self.connect()

        return self._connection

    @classmethod
    def parse_conf(self):
        """
        Read config file and return parsed informations.

        :returns: amqp url and the exchange name
        :rtype: string, string
        """
        config = configparser.RawConfigParser()
        config.read(os.path.join(root_path, DEFAULT_CONF_FILE))

        url = "amqp://{0}:{1}@{2}:{3}/{4}".format(
            config["master"]["userid"],
            config["master"]["password"],
            config["master"]["host"],
            config["master"]["port"],
            config["master"]["virtual_host"]
        )
        exname = config["master"]["exchange_name"]

        return url, exname

    def connect(self):
        """
        If connection is already made, disconnect then connect.

        You don't need te connect yourself if you use the channel or connection
        properties, is if they are None, AmqpConnection will
        handle (re)connection for you.

        :raises pika.exceptions.ConnectionClosed:
        """
        self.disconnect()
        parameters = pika.URLParameters(self._url)
        parameters.heartbeat = 3600
        parameters.retry_delay = 3
        parameters.connection_attempts = 3
        self._connection = pika.BlockingConnection(parameters)
        self._channel = self._connection.channel()

    def disconnect(self):
        """
        Close current connection, if connected, and resets
        self.connection and self.channel to None.
        """
        if self._channel is not None:
            try:
                self._channel.close()
            except (
                pika.exceptions.ChannelClosed,
                pika.exceptions.ConnectionClosed
            ):
                pass

            self._channel = None

        if self._connection is not None:
            try:
                self._connection.close()
            except pika.exceptions.ConnectionClosed:
                pass

            self._connection = None


class AmqpPublisher(object):
    """
    Easy to use synchronous AMQP publisher.

    Example:

    url = 'amqp://cpsrabbit:canopsis@localhost/canopsis'

    evt = {...}
    with AmqpConnection(url) as apc:
        pub = AmqpPublisher(apc, logger)
        pub.canopsis_event(evt)

    or:

    apc = AmqpConnection(url)
    apc.connect()

    pub = AmqpPublisher(apc, logger)
    pub.canopsis_event(evt)

    apc.disconnect()

    """

    def __init__(self, connection, logger):
        """
        :param connection AmqpConnection:
        :param logger Logger:
        """
        self.connection = connection
        self.logger = logger
        self._json_props = pika.BasicProperties(
            content_type='application/json')

    def json_document(self,
                      document,
                      exchange_name,
                      routing_key,
                      retries=3,
                      wait=1):
        """
        Sends a JSON document with AMQP content_type application/json

        :param document Any: a JSON serializable object
        :param exchange_name str: the name of the exchange to publish to.
        :param routing_key str: the event's routing key
        :param retries int: the number of times the publication should be
            retried in case of failure.
        :param wait float: the number of seconds to wait before retrying to
            publish the event.
        :raises AmqpPublishError: when all retries failed, raise this error.
        :raises TypeError: when the document cannot be serialized
        """
        # just ensure the connection is alive, if not, reconnect
        jdoc = json.dumps(document)

        retry = 0
        while retry <= retries:

            try:
                return self.connection.channel.basic_publish(
                    exchange_name, routing_key, jdoc, self._json_props
                )

            except (
                pika.exceptions.ConnectionClosed,
                pika.exceptions.ChannelClosed
            ):
                self.logger.warning(
                    "Failed to publish the following event ({}/{} retries)\n"
                    "{}".format(retry, retries, jdoc))
                try:
                    self.connection.connect()
                except pika.exceptions.ConnectionClosed:
                    if retry < retries:
                        time.sleep(wait)

            retry += 1

        raise AmqpPublishError(
            'cannot publish ({} times): cannot connect'.format(retry))

    def canopsis_event(self,
                       event,
                       exchange_name='canopsis.events',
                       retries=3,
                       wait=1):
        """
        Send an event to canopsis.

        :param event dict: a canopsis event (as a dictionnary).
        :param exchange_name str: the name of the exchange to publish to.
        :param retries int: the number of times the publication should be
            retried in case of failure.
        :param wait float: the number of seconds to wait before retrying to
        :raises KeyError: on invalid event, if routing key cannot be built.
        :raises AmqpPublishError: when all retries failed, raise this error.
        :raises TypeError: when the document cannot be serialized
        """
        return self.json_document(
            event, exchange_name, get_routingkey(event),
            retries=retries, wait=wait
        )

    def direct_event(self,
                     event,
                     queue_name,
                     exchange_name=DIRECT_EXCHANGE_NAME,
                     retries=3,
                     wait=1):
        """
        Send an event directly to a queue.

        :param event dict: a canopsis event (as a dictionnary).
        :param queue_name str: the name of the queue to publish to.
        :param exchange_name str: the name of the exchange to publish to.
        :param retries int: the number of times the publication should be
            retried in case of failure.
        :param wait float: the number of seconds to wait before retrying to
        :raises AmqpPublishError: when all retries failed, raise this error.
        :raises TypeError: when the document cannot be serialized
        """
        return self.json_document(
            event, exchange_name, queue_name,
            retries=retries, wait=wait
        )


def get_default_connection():
    """
    Provide default connection with parameters from etc/amqp.conf.
    """
    amqp_conf = Configuration.load(os.path.join('etc', 'amqp.conf'), Ini)
    amqp_url = 'amqp://{}:{}@{}:{}/{}'.format(
        amqp_conf['master']['userid'],
        amqp_conf['master']['password'],
        amqp_conf['master']['host'],
        amqp_conf['master']['port'],
        amqp_conf['master']['virtual_host']
    )

    return AmqpConnection(amqp_url)
