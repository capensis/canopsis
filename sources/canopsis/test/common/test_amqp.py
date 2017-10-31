import os
import unittest

from unittest import TestCase

from canopsis.common.amqp import AmqpPublisher


DEFAULT_AMQP_URL = 'amqp://guest:guest@localhost/'
DEFAULT_AMQP_EXCHANGE = 'test'


class TestAmqpPublisher(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.amqp_url = os.environ.get(
            'TEST_AMQPPUBLISHER_URL', DEFAULT_AMQP_URL)
        cls.amqp_exname = os.environ.get(
            'TEST_AMQPPUBLISHER_EXCHANGE', DEFAULT_AMQP_EXCHANGE)

    def test_connection_with_statement(self):

        with AmqpPublisher(self.amqp_url) as amqp_pub:
            self.assertIsNotNone(amqp_pub.connection)
            self.assertIsNotNone(amqp_pub.channel)

    def test_connection_explicit(self):

        amqp_pub = AmqpPublisher(self.amqp_url)
        amqp_pub.connect()

        self.assertIsNotNone(amqp_pub.connection)
        self.assertIsNotNone(amqp_pub.channel)

        amqp_pub.disconnect()

        self.assertIsNone(amqp_pub.connection)
        self.assertIsNone(amqp_pub.channel)

    def test_publish_event(self):

        event = {
            'connector': 'test_amqp',
            'connector_name': 'test_amqp',
            'source_type': 'resource',
            'event_type': 'check',
            'component': 'test',
            'resource': 'test'
        }

        with AmqpPublisher(self.amqp_url) as amqp_pub:
            amqp_pub.publish_event(event, self.amqp_exname)

    def test_publish_badevent_raises(self):

        event = {}

        with AmqpPublisher(self.amqp_url) as amqp_pub:
            with self.assertRaises(KeyError):
                amqp_pub.publish_event(event, self.amqp_exname)

if __name__ == '__main__':
    unittest.main()