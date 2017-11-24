import json
import os
import time

import pika

from canopsis.confng import Configuration, Ini
from canopsis.event import get_routingkey


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
        pub = AmqpPublisher(apc)
        pub.canopsis_event(evt, 'canopsis.events')

    or:

    apc = AmqpConnection(url)
    apc.connect()

    pub = AmqpPublisher(apc)
    pub.canopsis_event(evt)

    apc.disconnect()

    """

    def __init__(self, connection):
        """
        :type connection: AmqpConnection
        """
        self.connection = connection
        self._json_props = pika.BasicProperties(
            content_type='application/json')

    def json_document(
        self, document, exchange_name,
        routing_key, retries=3, wait=1
    ):
        """
        Sends a JSON document with AMQP content_type application/json

        :param retries: if the first try doesn't suceed, retry X times.
        :param document: valid JSON document
        :type document: dict
        :param exchange_name: exchange to publish to
        :type exchange_name: str
        :param routing_key: event's routing key
        :type routing_key: str
        :raises AmqpPublishError: when all retries failed, raise this error.
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
                try:
                    self.connection.connect()
                except pika.exceptions.ConnectionClosed:
                    if retry < retries:
                        time.sleep(wait)

            retry += 1

        raise AmqpPublishError(
            'cannot publish ({} times): cannot connect'.format(retry))

    def canopsis_event(self, event, exchange_name, retries=3, wait=1):
        """
        Shortcut to self.json_document, builds the routing key
        for you from the event.

        Event required keys:

            connector
            connector_name
            event_type
            source_type
            component
            resource if source_type == 'resource'

        :param event: valid Canopsis event.
        :raises KeyError: on invalid event, if routing key cannot be built.
        :param canopsis_exchange: exchange to publish to
        """
        return self.json_document(
            event, exchange_name, get_routingkey(event),
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
