#!/usr/bin/env python
# --------------------------------
# Copyright (c) 2014 'Capensis' [http://www.capensis.com]
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
# ----------------------

import sys
import time
import unittest
from os.path import join

from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.event import get_routingkey
from canopsis.old.rabbitmq import Amqp
from canopsis.old.record import Record

import logging

import ConfigParser
CONFIG = ConfigParser.ConfigParser()

# Default : 60s
BEAT_INTERVAL = 60

# Basic statuses
OFF = 0
ONGOING = 1

# Special Alerts statuses
STEALTHY = 2
BAGOT = 3

# Cancel action from UI
CANCEL = 4

NAMESPACE = 'events'
LOGGING_LEVEL = 'INFO' if not len(sys.argv) >= 2 else sys.argv[1].upper()

# Remove logging_level arg so unittest does not process it
sys.argv = [sys.argv[0]]


def event(name, state, **kwargs):
    event = {
        'connector': 'test',
        'connector_name': 'test',
        'event_type': 'check',
        'source_type': 'resource',
        'component': 'test',
        'resource': 'test',
        'state': state,
        'crecord_type': 'event',
        'state_type': 1,
        'pass_event': 1,
    }
    for key in kwargs:
        event[key] = kwargs[key]
    return (name, event)


def event_ok(**kwargs):
    return event('Event OK', 0, **kwargs)


def event_ko(**kwargs):
    return event('Event KO', 2, **kwargs)


class KnownValues(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.logger = logging.getLogger('TF_Archiver')
        cls.logger.setLevel(LOGGING_LEVEL)

        stdout_handler = logging.StreamHandler(sys.stdout)
        stdout_handler.setLevel(LOGGING_LEVEL)
        stdout_handler.setFormatter(
            logging.Formatter(
                '%(asctime)s [%(name)s] [%(levelname)s] %(message)s'
            )
        )
        cls.logger.addHandler(stdout_handler)

        cls.logger.debug(' + Init TF_Archiver on {}'.format(NAMESPACE))
        cls.account = Account(user='root', group='root')

        cls.logger.debug(' + Get storage')
        cls.storage = get_storage(namespace=NAMESPACE,
                                  logging_level=LOGGING_LEVEL)
        cls.collection = cls.storage.get_backend('events')

        cls.default_conf = cls.collection.find(
            {'crecord_type': 'statusmanagement'},
            namespace='object'
        )

        if cls.default_conf.count():
            cls.default_conf = cls.default_conf[0]
        else:
            cls.default_conf = {
                '_id': 'statusmanagement',
                'crecord_type': 'statusmanagement',
                'restore_event': True,
                'bagot_time': 3600,
                'bagot_freq': 10,
                'stealthy_time': 300,
                'stealthy_show': 300
            }

        cls.amqp = Amqp(logging_level=LOGGING_LEVEL,
                        logging_name='Amqp')

    def setUp(self):
        # Restore default conf
        self.change_conf()
        # Get current test name
        test_func = self.id().split('.')[-1].split('_')
        self.current_test = test_func[1] + ' ' + test_func[2].capitalize()
        self.logger.debug('+ {0}'.format(self.current_test))
        self.remove_event(connector=self.current_test)

    def tearDown(self):
        # Remove events
        self.remove_event(connector=self.current_test)

    def change_conf(self, sleep=1, **kwargs):
        conf = self.default_conf.copy()
        for key in kwargs:
            conf[key] = kwargs[key]
        record = Record(data=conf,
                        name="event state specifications",
                        _type='statusmanagement')
        self.storage.put(record,
                         namespace='object',
                         account=self.account)
        time.sleep(BEAT_INTERVAL + sleep)

    def publish_event(self, name, event, sleep=2):
        rk = get_routingkey(event)
        self.logger.debug("Sending event {}".format(name))
        self.amqp.publish(event,
                          rk,
                          'canopsis.events')
        time.sleep(sleep)

    def remove_event(self, **kwargs):
        self.collection.remove(**kwargs)
        time.sleep(1)

    def find_event(self, connector):
        cursor = self.collection.find(
            {'crecord_type': 'event', 'connector': connector})
        if cursor.count():
            return cursor[0]
        return {'status': -1}

    def test_01_off_basic_ok(self):
        self.logger.debug('OK : OFF')
        self.publish_event(*event_ok(connector='01 Off'))
        self.assertEqual(self.find_event('01 Off')['status'], OFF)

    def test_02_off_basic_okokokokok(self):
        self.logger.debug('OK -> OK -> OK -> OK -> OK : OFF')
        self.publish_event(*event_ok(connector='02 Off'))
        self.publish_event(*event_ok(connector='02 Off'))
        self.publish_event(*event_ok(connector='02 Off'))
        self.publish_event(*event_ok(connector='02 Off'))
        self.publish_event(*event_ok(connector='02 Off'))
        self.assertEqual(self.find_event('02 Off')['status'], OFF)

    def test_03_off_basic_kook_stealthytime(self):
        self.logger.debug(
            'Reduce time of stealthy time so the switch from KO to OK')
        self.logger.debug('does not trigger the Stealthy status')
        self.change_conf(sleep=5, stealthy_time=1)

        self.logger.debug('KO -> OK : OFF')
        self.publish_event(*event_ko(connector='03 Off'))
        self.publish_event(*event_ok(connector='03 Off'))
        self.assertEqual(self.find_event('03 Off')['status'], OFF)

    def test_04_off_basic_okkook_stealthytime(self):
        self.logger.debug(
            'Reduce time of stealthy time so the switch from KO to OK')
        self.logger.debug('does not trigger the Stealthy status')
        self.change_conf(sleep=5, stealthy_time=1)

        self.logger.debug('KO -> OK : OFF')
        self.publish_event(*event_ok(connector='04 Off'), sleep=5)
        self.publish_event(*event_ko(connector='04 Off'), sleep=5)
        self.publish_event(*event_ok(connector='04 Off'), sleep=5)
        self.assertEqual(self.find_event('04 Off')['status'], OFF)

    def test_05_off_ko_ok_stealthyshow(self):
        self.logger.debug('KO -> OK : STEALTHY [5s] OFF')
        self.publish_event(*event_ko(connector='05 Off'))
        self.publish_event(*event_ok(connector='05 Off'))
        self.assertEqual(self.find_event('05 Off')['status'], STEALTHY)

        self.logger.debug(
            'Reduce the time of stealthy show so the event goes from STEALTHY')
        self.logger.debug('to basic state avec 2 sec')
        self.change_conf(sleep=5, stealthy_show=2)
        self.assertEqual(self.find_event('05 Off')['status'], OFF)

    def test_06_ongoing_basic_okko(self):
        self.logger.debug('OK -> KO : ONGOING')
        self.publish_event(*event_ok(connector='06 OnGoing'))
        self.publish_event(*event_ko(connector='06 OnGoing'))
        self.assertEqual(self.find_event('06 OnGoing')['status'], ONGOING)

    def test_07_ongoing_okkokook_stealthytime(self):
        self.change_conf(sleep=5, stealthy_time=1)

        self.logger.debug('OK -> KO : ONGOING')
        self.publish_event(*event_ok(connector='07 OnGoing'), sleep=4)
        self.publish_event(*event_ko(connector='07 OnGoing'), sleep=4)
        self.publish_event(*event_ok(connector='07 OnGoing'), sleep=4)
        self.publish_event(*event_ko(connector='07 OnGoing'), sleep=4)
        self.assertEqual(self.find_event('07 OnGoing')['status'], ONGOING)

    def test_08_ongoing_okkookko_stealthyshow(self):
        self.logger.debug('OK -> KO : ONGOING')
        self.publish_event(*event_ok(connector='08 OnGoing'))
        self.publish_event(*event_ko(connector='08 OnGoing'))
        self.publish_event(*event_ok(connector='08 OnGoing'))
        self.publish_event(*event_ko(connector='08 OnGoing'))
        self.assertEqual(self.find_event('08 OnGoing')['status'], STEALTHY)

        self.logger.debug(
            'Reduce the time of stealthy show so the event goes from STEALTHY')
        self.logger.debug('to basic state avec 2 sec')
        self.change_conf(sleep=5, stealthy_show=2)
        self.assertEqual(self.find_event('08 OnGoing')['status'], ONGOING)

    def test_09_stealthy_basic_kook(self):
        self.logger.debug('KO -> OK : STEALTHY')
        self.publish_event(*event_ko(connector='09 Stealthy'))
        self.publish_event(*event_ok(connector='09 Stealthy'))
        self.assertEqual(self.find_event('09 Stealthy')['status'], STEALTHY)

    def test_10_stealthy_basic_okkook(self):
        self.logger.debug('OK -> KO -> OK : STEALTHY')
        self.publish_event(*event_ok(connector='10 Stealthy'))
        self.publish_event(*event_ko(connector='10 Stealthy'))
        self.publish_event(*event_ok(connector='10 Stealthy'))
        self.assertEqual(self.find_event('10 Stealthy')['status'], STEALTHY)

    def test_11_stealthy_basic_kookokko(self):
        self.logger.debug('OK -> KO -> OK : STEALTHY')
        self.publish_event(*event_ko(connector='11 Stealthy'))
        self.publish_event(*event_ok(connector='11 Stealthy'))
        self.publish_event(*event_ok(connector='11 Stealthy'))
        self.publish_event(*event_ko(connector='11 Stealthy'))
        self.assertEqual(self.find_event('11 Stealthy')['status'], STEALTHY)

    def test_12_stealthy_notbagot(self):
        self.logger.debug('KO -> OK (x3) : STEALTHY')
        for i in xrange(3):
            self.publish_event(*event_ko(connector='12 Bagot'))
            self.publish_event(*event_ok(connector='12 Bagot'))
        self.assertEqual(self.find_event('12 Bagot')['status'], STEALTHY)

    def test_13_bagot_basic(self):
        self.logger.debug('KO -> OK (x10) : BAGOT')
        for i in xrange(10):
            self.publish_event(*event_ko(connector='13 Bagot'))
            self.publish_event(*event_ok(connector='13 Bagot'))
        self.assertEqual(self.find_event('13 Bagot')['status'], BAGOT)

    def test_14_bagot_basic(self):
        self.change_conf(sleep=5, bagot_freq=3)
        self.logger.debug('KO -> OK (x3) : BAGOT')
        for i in xrange(3):
            self.publish_event(*event_ko(connector='14 Bagot'))
            self.publish_event(*event_ok(connector='14 Bagot'))
        self.assertEqual(self.find_event('14 Bagot')['status'], BAGOT)


if __name__ == '__main__':

    CONFIG.read(join(sys.prefix, 'etc', 'amqp2engines.conf'))
    try:
        BEAT_INTERVAL = int(CONFIG.get('engine:eventstore', 'beat_interval'))
    except:
        BEAT_INTERVAL = 60

    unittest.main(verbosity=2)

    #functest=TestCase()

    #func_tests = [getattr(functest, func)
    #              for func in sorted(dir(functest))
    #              if len(func) >= 4 and func[:4] == 'test']
    #for test in func_tests:
    #    test()
