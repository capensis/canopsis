#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import os
import unittest
from kombu import pools, Connection
from contextlib import contextmanager

from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.old.rabbitmq import Amqp

import xmlrunner


class TestAmqp(unittest.TestCase):

    @classmethod
    def setUpClass(cls):

        amqp_conf = Configuration.load(os.path.join('etc', 'amqp.conf'), Ini)
        cls.amqp_uri = 'amqp://{}:{}@{}:{}/{}'.format(
            amqp_conf['master']['userid'],
            amqp_conf['master']['password'],
            amqp_conf['master']['host'],
            amqp_conf['master']['port'],
            amqp_conf['master']['virtual_host']
        )

        cls.conn = Connection(cls.amqp_uri)
        cls.producers = pools.Producers(limit=1)
        cls.exchange_name = "canopsis"

        cls.amqp = Amqp(logging_level='INFO',
                        logging_name='Amqp')
        cls.amqp.producers = cls.producers
        cls.amqp.conn = cls.conn

        cls.event = {
            'connector': 'test_amqp',
            'connector_name': 'test_amqp',
            'source_type': 'resource',
            'event_type': 'check',
            'component': 'test',
            'resource': 'test'
        }

    @contextmanager
    def assertNotRaises(self, exc_type):
        """
        Opposite of assertRaises()
        # https://gist.github.com/AntoineCezar/9086f5374888f24eb315
        """
        try:
            yield None
        except exc_type:
            raise self.failureException('{} raised'.format(exc_type.__name__))


class TestKombuAmqpPublisher(TestAmqp):

    def test_too_long_rk(self):
        msg = 'coucou'
        rk = 'Johann.Gambolputty.de.von.Ausfern-schplenden-schiltter-crasscrenbon-fried-digger-dingle-dangle-dongle-dungle-burstein-von-knacker-thrasher-apple-banger-horowitz-ticolensic-grander-knotty-spelltinkle-grandlich-grumblemeyer-spelterwasser-kurstlich-himbleeisen-bahnwagen-gutenabend-bitte-ein-nürnburger-bratwustle-gerspurten-mitz-weimache-luber-hundsfut-gumeraber-shönendanker-kalbsfleisch-mittler-aucher.von.Hautkopft.of.Ulm'

        # Isolated test
        import struct
        with self.producers[self.conn].acquire(block=True) as producer:
            with self.assertRaises(struct.error):
                producer.publish(
                    msg,
                    serializer="json",
                    routing_key=rk,
                    exchange=self.exchange_name
                )
            with self.assertNotRaises(struct.error):
                try:
                    producer.publish(
                        msg,
                        serializer="json",
                        routing_key=rk,
                        exchange=self.exchange_name
                    )
                except struct.error:
                    print('RK is too long !')

        # Inplace test (should not publish)
        ret = self.amqp.publish(
            msg,
            rk,
            exchange_name=self.exchange_name,
            serializer="json"
        )
        self.assertFalse(ret)


if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
