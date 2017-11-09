# -*- coding: utf-8 -*-
import json

from canopsis.backbone.event import Event


class Consumer(object):
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
