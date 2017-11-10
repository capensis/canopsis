# -*- coding: utf-8 -*-
import json
import pika

from canopsis.backbone.event import Event
from canopsis.event import get_routingkey


class AmqpConnection(object):

    def __init__(self, url):
        """
        :param url: url of the form: amqp://[<user>:<pass>]@host:port/vhost
        :type url: str
        """

        self._url = url
        self.connection = None
        self.channel = None

    def __enter__(self):
        self.connect()
        return self

    def __exit__(self, type_, value, traceback):
        self.disconnect()

    @property
    def connected(self):
        """
        Property checking for connection state.
        """
        if self.connection is None:
            return False

        return self.connection.is_open

    def connect(self):
        """
        If connection is already made, disconnect then connect.
        """
        if self.connected:
            self.disconnect()

        parameters = pika.URLParameters(self._url)
        self.connection = pika.BlockingConnection(parameters)
        self.channel = self.connection.channel()

    def disconnect(self):
        """
        Close current connection, if connected, and resets
        self.connection and self.channel to None.
        """
        if self.connected:
            self.connection.close()

        self.connection = None
        self.channel = None


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

    def json_document(self, document, exchange_name, routing_key):
        """
        Sends a JSON document with AMQP content_type application/json

        :param document: valid JSON document
        :type document: dict
        :param exchange_name: exchange to publish to
        :type exchange_name: str
        :param routing_key: event's routing key
        :type routing_key: str
        """
        jdoc = json.dumps(document)
        props = pika.BasicProperties(content_type='application/json')
        return self.connection.channel.basic_publish(
            exchange_name, routing_key, jdoc, props
        )

    def canopsis_event(self, event, exchange_name):
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
        return self.json_document(event, exchange_name, get_routingkey(event))


class AmqpConsumer(object):
    def __init__(self, connection, queue):
        """
        :param Connection connection:
        :param str queue:
        """
        self.connection = connection
        self.queue = queue
        self.channel = self.connection.channel()
        self.channel.connection.process_data_events(time_limit=1)
        self.callback = None

    def consume(self):
        """
        Start consuming events.
        """
        self.channel.basic_consume(self._work, queue=self.queue, no_ack=True)
        self.channel.start_consuming()

    def _work(self, unused_channel, basic_deliver, properties, body):
        dico = json.loads(body)
        event = Event(**dico)
        self.callback(event)
