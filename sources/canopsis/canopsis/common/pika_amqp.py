# -*- coding: utf-8 -*-
import pika
#from concurrent.futures import ThreadPoolExecutor
import threading
from canopsis.backbone.event import Event
import json

class Consumer(object):
    def __init__(self, connection, queue):
        self.connection = connection
        self.queue = queue
        self.channel = self.connection.channel()
        self.channel.connection.process_data_events(time_limit=1)
        self.callback = None

    def consume(self):
        self.channel.basic_consume(self._work, queue=self.queue, no_ack=True)
        self.channel.start_consuming()

    def _work(self, unused_channel, basic_deliver, properties, body):
        truc = json.loads(body)
        event = Event(**truc)
        self.callback(event)
