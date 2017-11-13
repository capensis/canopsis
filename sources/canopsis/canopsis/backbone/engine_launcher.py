import pika
from canopsis.backbone.engine_cleaner import EngineCleaner
from canopsis.common.pika_amqp import Consumer
from canopsis.logger import Logger

logger = Logger.get('cleaner2', 'cleaner2.log')
credentials = pika.PlainCredentials('cpsrabbit', 'canopsis')
connection = pika.BlockingConnection(pika.ConnectionParameters('localhost', 5672, 'canopsis', credentials))
c = EngineCleaner(connection, 'Engine_cleaner_events', logger)
c.consume()

