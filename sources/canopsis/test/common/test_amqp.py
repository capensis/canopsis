#!/usr/bin/env python
# -*- coding: utf-8 -*-

from json import dumps
import os
import unittest

from canopsis.backbone.event import Event
from canopsis.common.amqp import (
    AmqpPublisher, AmqpConnection, AmqpConsumer,
    AmqpNotCallableCallback
)

DEFAULT_AMQP_URL = 'amqp://guest:guest@localhost/'
DEFAULT_AMQP_EXCHANGE = 'test'
DEFAULT_AMQP_QUEUE = 'test_queue'


class TestAmqp(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.amqp_url = os.environ.get(
            'TEST_AMQPPUBLISHER_URL', DEFAULT_AMQP_URL)
        cls.amqp_exname = os.environ.get(
            'TEST_AMQPPUBLISHER_EXCHANGE', DEFAULT_AMQP_EXCHANGE)
        cls.amqp_queue = os.environ.get(
            'TEST_AMQPPUBLISHER_QUEUE', DEFAULT_AMQP_QUEUE)


class TestAmqpConn(TestAmqp):

    def test_connection_with_statement(self):
        with AmqpConnection(self.amqp_url) as amqp_conn:
            self.assertIsNotNone(amqp_conn.connection)
            self.assertIsNotNone(amqp_conn.channel)
            self.assertTrue(amqp_conn.connected)

    def test_connection_explicit(self):
        amqp_conn = AmqpConnection(self.amqp_url)
        amqp_conn.connect()

        self.assertIsNotNone(amqp_conn.connection)
        self.assertIsNotNone(amqp_conn.channel)
        self.assertTrue(amqp_conn.connected)

        amqp_conn.disconnect()

        self.assertIsNone(amqp_conn.connection)
        self.assertIsNone(amqp_conn.channel)
        self.assertFalse(amqp_conn.connected)


class TestAmqpPublisher(TestAmqp):

    def test_canopsis_event(self):
        event = {
            'connector': 'test_amqp',
            'connector_name': 'test_amqp',
            'source_type': 'resource',
            'event_type': 'check',
            'component': 'test',
            'resource': 'test'
        }

        with AmqpConnection(self.amqp_url) as ac:
            amqp_pub = AmqpPublisher(ac)
            amqp_pub.canopsis_event(event, self.amqp_exname)

    def test_bad_canopsis_event_raises(self):
        event = {}

        with AmqpConnection(self.amqp_url) as ac:
            amqp_pub = AmqpPublisher(ac)
            with self.assertRaises(KeyError):
                amqp_pub.canopsis_event(event, self.amqp_exname)

    def test_json_document(self):
        jdoc = {'bla': 'bla'}
        with AmqpConnection(self.amqp_url) as ac:
            amqp_pub = AmqpPublisher(ac)
            amqp_pub.json_document(jdoc, self.amqp_exname, '#')


class TestAmqpConsumer(TestAmqp):

    basic_event = {
        "source_type": "resource",
        "event_type": "check",
        "connector": "alien_schemes",
        "connector_name": "mobilize",
        "component": "machines_vs_brains",
        "resource": "galactic_overview",
        "state": 1,
    }

    def test_bad_callback(self):
        with AmqpConnection(self.amqp_url) as ac:
            amqp_cons = AmqpConsumer(ac, self.amqp_queue)
            with self.assertRaises(AmqpNotCallableCallback):
                amqp_cons.on_message()

    def test_consume(self):

        def callback(event):
            self.assertTrue(isinstance(event, Event))
            self.assertEqual(event.resource, self.basic_event['resource'])

        with AmqpConnection(self.amqp_url) as ac:
            amqp_cons = AmqpConsumer(ac, self.amqp_queue)
            amqp_cons._on_message = callback
            amqp_cons._work(None, None, None, dumps(self.basic_event))

if __name__ == '__main__':
    unittest.main()

