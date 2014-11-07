#!/usr/bin/env python
#--------------------------------
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
import imp
import time
import unittest
from os.path import join

from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.old.event import get_routingkey, forger
from canopsis.old.rabbitmq import Amqp
from canopsis.old.record import Record

import logging

# Basic statuses
OFF=0
ONGOING=1

# Special Alerts statuses
STEALTHY=2
BAGOT=3

# Cancel action from UI
CANCEL=4

namespace='events'
logging_level= 'INFO' if not len(sys.argv) >= 2 else sys.argv[1].upper()

# Remove logging_level arg so unittest does not process it
sys.argv = [sys.argv[0]]

def event(name, state, **kwargs):
    event = {
        'connector'     : 'test',
        'connector_name': 'test',
        'event_type'    : 'check',
        'source_type'   : 'resource',
        'component'     : 'test',
        'resource'      : 'test',
        'state'         : state,
        'crecord_type'  : 'event',
        'state_type'    : 1
        }
    for key in kwargs:
        event[key] = kwargs[key]
    return (name, event)

def event_ok(**kwargs):
    return event('Event OK', 0, **kwargs)

def event_ko(**kwargs):
    return event('Event KO', 2, **kwargs)

class KnownValues(unittest.TestCase):

    def setUp(self):
        self.logger = logging.getLogger('TF_Archiver')
        self.logger.setLevel(logging_level)

        stdout_handler = logging.StreamHandler(sys.stdout)
        stdout_handler.setLevel(logging_level)
        stdout_handler.setFormatter(logging.Formatter(
                '%(asctime)s [%(name)s] [%(levelname)s] %(message)s'
            ))
        self.logger.addHandler(stdout_handler)

        self.logger.debug(' + Init TF_Archiver on %s' % namespace)
        self.account = Account(user='root', group='root')

        self.logger.debug(' + Get storage')
        self.storage = get_storage(namespace=namespace,
                                   logging_level=logging_level)
        self.collection = self.storage.get_backend('events')

        self.default_conf = self.collection.find(
            {'crecord_type': 'statusmanagement'},
            namespace='object'
            )

        if self.default_conf.count():
            self.default_conf = self.default_conf[0]
        else:
            self.default_conf =  {
                '_id': 'statusmanagement',
                'crecord_type': 'statusmanagement',
                'restore_event': True,
                'bagot_time': 3600,
                'bagot_freq': 10,
                'stealthy_time': 300,
                'stealthy_show': 300
            }

        self.amqp = Amqp(
            logging_level=logging_level,
            logging_name='Amqp'
            )

    def change_conf(self, sleep=1, **kwargs):
        conf = self.default_conf.copy()
        for key in kwargs:
            conf[key] = kwargs[key]
        record = Record(
            data=conf,
            name="event state specifications",
            _type='statusmanagement'
            )
        record.chmod("g+w")
        record.chmod("o+r")
        record.chgrp('group.CPS_root')
        self.storage.put(record,
                         namespace='object',
                         account=self.account)
        time.sleep(sleep)

    def publish_event(self, name, event, sleep=0):
        rk = get_routingkey(event)
        self.logger.debug("Sending event {}".format(name))
        self.amqp.publish(
            event,
            rk,
            'canopsis.events'
            )
        time.sleep(1+sleep)

    def remove_event(self, **kwargs):
        self.collection.remove(kwargs)
        time.sleep(1)

    def find_event(self, connector):
        cursor = self.collection.find(
            {'crecord_type':'event', 'connector': connector})
        if cursor.count():
            return cursor[0]
        return {'status': -1}

    def test_01_off_basic_ok(self):
        self.logger.debug('+ 01 Off')
        self.remove_event(connector='01Off')

        # OK : OFF
        self.publish_event(*event_ok(connector='01Off'))
        self.assertEqual(self.find_event('01Off')['status'], OFF)

    def test_02_off_basic_okokokokok(self):
        self.logger.debug('+ 02 Off')
        self.remove_event(connector='02Off')

        # OK -> OK -> OK -> OK -> OK : OFF
        self.publish_event(*event_ok(connector='02Off'))
        self.publish_event(*event_ok(connector='02Off'))
        self.publish_event(*event_ok(connector='02Off'))
        self.publish_event(*event_ok(connector='02Off'))
        self.publish_event(*event_ok(connector='02Off'))
        self.assertEqual(self.find_event('02Off')['status'], OFF)

    def test_03_off_basic_kook_stealthytime(self):
        self.logger.debug('+ 03 Off')
        self.remove_event(connector='03Off')

        # Reduce time of stealthy time so the switch from KO to OK
        # does not trigger the Stealthy status
        self.change_conf(sleep=5, stealthy_time=1)

        # KO -> OK : OFF
        self.publish_event(*event_ko(connector='03Off'), sleep=2)
        self.publish_event(*event_ok(connector='03Off'))
        self.assertEqual(self.find_event('03Off')['status'], OFF)

        # Restore default conf
        self.change_conf()

    def test_04_off_basic_okkook_stealthytime(self):
        self.logger.debug('+ 04 Off')
        self.remove_event(connector='04Off')

        # Reduce time of stealthy time so the switch from KO to OK
        # does not trigger the Stealthy status
        self.change_conf(sleep=5, stealthy_time=1)

        # KO -> OK : OFF
        self.publish_event(*event_ok(connector='04Off'), sleep=2)
        self.publish_event(*event_ko(connector='04Off'), sleep=2)
        self.publish_event(*event_ok(connector='04Off'))
        self.assertEqual(self.find_event('04Off')['status'], OFF)

        # Restore default conf
        self.change_conf()

    def test_05_off_ko_ok_stealthyshow(self):
        self.logger.debug('+ 05 Off')
        self.remove_event(connector='05Off')

        # KO -> OK : STEALTHY [5s] OFF
        self.publish_event(*event_ko(connector='05Off'))
        self.publish_event(*event_ok(connector='05Off'))
        self.assertEqual(self.find_event('05Off')['status'], STEALTHY)

        # Reduce the time of stealthy show so the event goes from STEALTHY
        # to basic state avec 2 sec
        self.change_conf(sleep=5, stealthy_show=2)
        self.assertEqual(self.find_event('05Off')['status'], OFF)

        # Restore default conf
        self.change_conf()

    def test_06_ongoing_basic_okko(self):
        self.logger.debug('+ 06 OnGoing')
        self.remove_event(connector='06OnGoing')

        # OK -> KO : ONGOING
        self.publish_event(*event_ok(connector='06OnGoing'))
        self.publish_event(*event_ko(connector='06OnGoing'))
        self.assertEqual(self.find_event('06OnGoing')['status'], ONGOING)

    def test_07_ongoing_okkokook_stealthytime(self):
        self.logger.debug('+ 07 OnGoing')
        self.remove_event(connector='07OnGoing')

        self.change_conf(sleep=5, stealthy_time=1)

        # OK -> KO : ONGOING
        self.publish_event(*event_ok(connector='07OnGoing'), sleep=2)
        self.publish_event(*event_ko(connector='07OnGoing'), sleep=2)
        self.publish_event(*event_ok(connector='07OnGoing'), sleep=2)
        self.publish_event(*event_ko(connector='07OnGoing'), sleep=2)
        self.assertEqual(self.find_event('07OnGoing')['status'], ONGOING)

        # Restore default conf
        self.change_conf()

    def test_08_ongoing_okkookko_stealthyshow(self):
        self.logger.debug('+ 08 OnGoing')
        self.remove_event(connector='08OnGoing')

        # OK -> KO : ONGOING
        self.publish_event(*event_ok(connector='08OnGoing'))
        self.publish_event(*event_ko(connector='08OnGoing'))
        self.publish_event(*event_ok(connector='08OnGoing'))
        self.publish_event(*event_ko(connector='08OnGoing'))
        self.assertEqual(self.find_event('08OnGoing')['status'], STEALTHY)

        # Reduce the time of stealthy show so the event goes from STEALTHY
        # to basic state avec 2 sec
        self.change_conf(sleep=5, stealthy_show=2)
        self.assertEqual(self.find_event('08OnGoing')['status'], ONGOING)

        # Restore default conf
        self.change_conf()

    def test_09_stealthy_basic_kook(self):
        self.logger.debug('+ 09 Stealthy')
        self.remove_event(connector='09Stealthy')

        # KO -> OK : STEALTHY
        self.publish_event(*event_ko(connector='09Stealthy'))
        self.publish_event(*event_ok(connector='09Stealthy'))
        self.assertEqual(self.find_event('09Stealthy')['status'], STEALTHY)

    def test_10_stealthy_basic_okkook(self):
        self.logger.debug('+ 10 Stealthy')
        self.remove_event(connector='10Stealthy')

        # OK -> KO -> OK : STEALTHY
        self.publish_event(*event_ok(connector='10Stealthy'))
        self.publish_event(*event_ko(connector='10Stealthy'))
        self.publish_event(*event_ok(connector='10Stealthy'))
        self.assertEqual(self.find_event('10Stealthy')['status'], STEALTHY)

    def test_11_stealthy_basic_kookokko(self):
        self.logger.debug('+ 11 Stealthy')
        self.remove_event(connector='11Stealthy')

        # OK -> KO -> OK : STEALTHY
        self.publish_event(*event_ko(connector='11Stealthy'))
        self.publish_event(*event_ok(connector='11Stealthy'))
        self.publish_event(*event_ok(connector='11Stealthy'))
        self.publish_event(*event_ko(connector='11Stealthy'))
        self.assertEqual(self.find_event('11Stealthy')['status'], STEALTHY)

    def test_12_stealthy_notbagot(self):
        self.logger.debug('+ 12 Bagot')
        self.remove_event(connector='12Bagot')

        # KO -> OK (x3) : STEALTHY
        for i in xrange(3):
            self.publish_event(*event_ko(connector='12Bagot'))
            self.publish_event(*event_ok(connector='12Bagot'))
        self.assertEqual(self.find_event('12Bagot')['status'], STEALTHY)

    def test_13_bagot_basic(self):
        self.logger.debug('+ 13 Bagot')
        self.remove_event(connector='13Bagot')

        # KO -> OK (x10) : BAGOT
        for i in xrange(10):
            self.publish_event(*event_ko(connector='13Bagot'))
            self.publish_event(*event_ok(connector='13Bagot'))
        self.assertEqual(self.find_event('13Bagot')['status'], BAGOT)

    def test_13_bagot_basic(self):
        self.logger.debug('+ 13 Bagot')
        self.remove_event(connector='13Bagot')

        self.change_conf(sleep=1, bagot_freq=3)
        # KO -> OK (x3) : BAGOT
        for i in xrange(3):
            self.publish_event(*event_ko(connector='13Bagot'))
            self.publish_event(*event_ok(connector='13Bagot'))
        self.assertEqual(self.find_event('13Bagot')['status'], BAGOT)

        # Restore conf
        self.change_conf()

if __name__ == '__main__':
    unittest.main(verbosity=2)

    # functest=TestCase()

    # func_tests = [getattr(functest, func)
    #               for func in sorted(dir(functest))
    #               if len(func) >= 4 and func[:4] == 'test']
    # for test in func_tests:
    #     test()
