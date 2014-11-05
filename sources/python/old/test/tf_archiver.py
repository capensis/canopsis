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

OFF=0
ONGOING=1
STEALTHY=2
BAGOT=3
CANCEL=4
namespace='events'
logging_level='INFO'

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

        self.amqp = Amqp(
            logging_level=logging_level,
            logging_name='Amqp'
            )


    def change_conf(self, **kwargs):
        conf = {
            '_id': 'statusmanagement',
            'crecord_type': 'statusmanagement',
            'restore_event': True,
            'bagot_time': 3600,
            'bagot_freq': 10,
            'stealthy_time': 300,
            'stealthy_show': 300
            }
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

    def publish_event(self, name, event):
        rk = get_routingkey(event)
        self.logger.info("Sending event {}".format(name))
        self.amqp.publish(
            event,
            rk,
            'canopsis.events'
            )
        time.sleep(1)

    def remove_event(self, **kwargs):
        self.collection.remove(kwargs)
        time.sleep(1)

    def find_event(self, connector):
        cursor = self.collection.find(
            {'crecord_type':'event', 'connector': connector})
        if cursor.count():
            return cursor[0]
        return {'status': -1}

    def test_01_stealthy(self):
        self.logger.info('+ 01 Stealthy')
        self.remove_event(connector='01Stealthy')

        # KO -> OK : STEALTHY
        self.publish_event(*event_ko(connector='01Stealthy'))
        self.publish_event(*event_ok(connector='01Stealthy'))
        self.assertEqual(self.find_event('01Stealthy')['status'], STEALTHY)

        self.change_conf(stealthy_show=3)
        time.sleep(5)
        self.assertEqual(self.find_event('01Stealthy')['status'], OFF)

        # Restore default conf
        self.change_conf()

if __name__ == '__main__':
    unittest.main(verbosity=2)
