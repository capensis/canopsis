from __future__ import unicode_literals

import pika

from canopsis.backbone.engine_cleaner import EngineCleaner
from canopsis.logger import Logger

logger = Logger.get('cleaner2', 'var/log/engines/cleaner2.log')
credentials = pika.PlainCredentials('cpsrabbit', 'canopsis')
parameters = pika.ConnectionParameters('localhost',
                                       5672,
                                       'canopsis',
                                       credentials)
connection = pika.BlockingConnection(parameters)
c = EngineCleaner(connection, 'Engine_cleaner_events', logger)
c.consume()
