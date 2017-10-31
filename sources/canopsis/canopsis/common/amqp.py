import pika
import json

from canopsis.event import get_routingkey


class AmqpPublisher(object):
    """
    Easy to use synchronous AMQP publisher.

    Example:

    url = 'amqp://cpsrabbit:canopsis@localhost/canopsis'

    evt = {...}
    with AmqpPublisher(url) as ap:
        ap.event_publish(evt, 'canopsis.events')

    or:

    ap = AmqpPublisher(url)
    ap.connect()
    ap.publish(evt)
    ap.disconnect()

    """

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

    def __exit__(self, type, value, traceback):
        self.disconnect()

    @property
    def connected(self):
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

    def publish_json(self, document, exchange_name, routing_key):
        """
        :param document: valid JSON document
        :type document: dict
        :param exchange_name: exchange to publish to
        :type exchange_name: str
        :param routing_key: event's routing key
        :type routing_key: str
        """
        jdoc = json.dumps(document)
        props = pika.BasicProperties(content_type='application/json')
        return self.channel.basic_publish(
            exchange_name, routing_key, jdoc, props
        )

    def publish_event(self, event, exchange_name):
        """
        Shortcut to publish_json, builds the routing key
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
        return self.publish_json(event, exchange_name, get_routingkey(event))

    def disconnect(self):
        if self.connected:
            self.connection.close()

        self.connection = None
        self.channel = None
