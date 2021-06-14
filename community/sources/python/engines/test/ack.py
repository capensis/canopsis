#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from unittest import TestCase, main


from os.path import expanduser

from sys import path

from time import time, sleep

from canopsis.engines.core import publish as _publish
from canopsis.old.rabbitmq import Amqp

from canopsis.old.storage import get_storage
from canopsis.old.account import Account

from logging import INFO

path.append(expanduser('~/opt/amqp2engines/engines/'))

before_ack_sent = time()


# Mock ready function for amqp
def ready():
    pass

amqp = Amqp(
    logging_level=INFO,
    logging_name='ack-amqp',
    on_ready=ready
)


def get_rk(event):
    routing_key = "%s.%s.%s.%s.%s" % (
        event['connector'], event['connector_name'], event['event_type'],
        event['source_type'], event['component'])
    if event['source_type'] == "resource":
        routing_key += ".%s" % event['resource']
    return routing_key


def publish(event):

    _publish(publisher=amqp, event=event)


def log(message):
    print(' + ' + message)


class KnownValues(TestCase):
    def setUp(self):
        self.storage = get_storage(
            namespace='object', account=Account(user="root", group="root"))
        self.ack = self.storage.get_backend('ack')
        self.event = {
            "connector": "unit-test",
            "connector_name": "canopsis",
            "event_type": "check",
            "source_type": "resource",
            "component": "ack-test-event",
            "resource": "error-test-event",
            "state": 1,
            "state_type": 1,
            "output": "ERROR-UNITTEST",
        }
        self.rk = get_rk(self.event)

    def test_01_send_error_event_to_ack(self):

        log('Send an erroneous event to be ack')

        # Clean event
        self.ack.remove({'rk': self.rk})

        publish(self.event)

        log('Wait until event is inserted before acquiting')
        sleep(3)

    def test_02_send_ack_to_error_event(self):

        event = self.event.copy()
        event['referer'] = self.rk
        event['event_type'] = 'ack'
        event['author'] = 'canopsis_unit_test'
        #publishing ack event
        publish(event)

        log('Wait for ack event beeing threaten')
        sleep(3)

        log(
            'Now entities collection should contain an ack event for this scenario'
        )

        ack_insert = self.ack.find_one({'rk': self.rk})
        self.assertTrue(ack_insert)
        self.assertFalse(ack_insert['solved'])

    def test_03_send_ok_event_ack_solved(self):

        # Tests solved ack by setting event to ok state
        self.event['state'] = 0
        publish(self.event)
        log('Wait for valid event to solve ack')
        sleep(3)

        ack_insert = self.ack.find_one({'rk': self.rk})
        self.assertTrue(ack_insert)
        self.assertTrue(ack_insert['solved'])
        self.assertTrue(ack_insert['solvedts'] > before_ack_sent)

    def test_04_send_again_error_event_ack_reset(self):

        # Tests solved ack by setting event to ok state
        self.event['state'] = 1
        publish(self.event)
        log('Wait for invalid event to reset ack')
        sleep(3)

        ack_insert = self.ack.find_one({'rk': self.rk})
        self.assertTrue(ack_insert)
        self.assertFalse(ack_insert['solved'])
        self.assertTrue(ack_insert['solvedts'] == -1)
        self.assertTrue(ack_insert['timestamp'] == -1)

if __name__ == "__main__":
    main()
