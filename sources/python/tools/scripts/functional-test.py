#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from time import time, sleep

from logging import DEBUG, basicConfig

from canopsis.old.rabbitmq import Amqp
from canopsis.old.storage import Storage
from canopsis.old.account import Account
from canopsis.old.webservices import Webservices
from canopsis.old.tools import parse_perfdata
from canopsis.old.event import get_routingkey, forger

from commands import getstatusoutput


basicConfig(
    level=DEBUG,
    format='%(asctime)s %(name)s %(levelname)s %(message)s')

event = forger(
    connector='canopsis', connector_name='unittest', event_type='check',
    source_type="component", component="test1", state=0, output="Output_1",
    perf_data="mymetric=1s;10;20;0;30",
    tags=['check', 'component', 'test1', 'unittest'])
rk = get_routingkey(event)

myamqp = None
storage = None
event_alert = None
perfstore = None


def on_alert(body, message):
    print("Alert: %s" % body)
    mrk = message.delivery_info['routing_key']
    if mrk == rk:
        global event_alert
        event_alert = body


def clean():
        storage.remove(rk)
        records = storage.find({'rk': rk}, namespace='events_log')
        storage.remove(records, namespace='events_log')

        try:
            perfstore.remove(name='test1mymetric')
        except:
            pass


class KnownValues(TestCase):
    def setUp(self):
        self.rcvmsgbody = None

    def test_1_Init(self):
        global myamqp
        myamqp = Amqp()
        myamqp.add_queue(
            queue_name="unittest_alerts",
            routing_keys="#",
            callback=on_alert,
            exchange_name=myamqp.exchange_name_alerts)
        myamqp.start()
        sleep(1)

        global storage
        storage = Storage(
            Account(user="root", group="root"), namespace='events',
            logging_level=DEBUG)


        clean()

    def test_2_PubState(self):
        myamqp.publish(event, rk, exchange_name=myamqp.exchange_name_events)
        sleep(3)

    def test_3_Check_amqp2engines(self):
        record = storage.get(rk)
        revent = record.data

        if revent['component'] != event['component']:
            raise Exception('Invalid data ...')

        if revent['timestamp'] != event['timestamp']:
            raise Exception('Invalid data ...')

        if revent['state'] != event['state']:
            raise Exception('Invalid data ...')

        del event_alert['_id']

        # remove cps_state

        if 'perf_data_array' in event_alert and len(event_alert['perf_data_array']) >= 2:
            del event_alert['perf_data_array'][1]

        event['perf_data_array'] = parse_perfdata(event['perf_data'])

        try:
            event['rk'] = event_alert['rk']
            event['event_id'] = event_alert['event_id']
        except:
            pass

        del event_alert['last_state_change']


    def test_4_Check_amqp2engines_archiver(self):
        ## change state
        event['state'] = 1
        event['timestamp'] = int(time())
        myamqp.publish(event, rk, exchange_name=myamqp.exchange_name_events)
        sleep(3)

        records = storage.find(
            {'event_id': rk}, sort='timestamp', namespace='events_log')

        if len(records) != 2:
            raise Exception("Archiver don't work ...")

        revent = records[1].data

        if revent['state'] != event['state']:
            raise Exception('Invalid log state')


    def test_6_Check_webserver(self):
        import requests
        r = requests.get('http://localhost:8082')
        self.assertEqual(r.status_code, 200)

    def test_7_Check_collectd2event(self):

        print("Restart collectd ...")
        getstatusoutput("service collectd restart")

        print("Wait collectd events ...")
        i = 0
        while i < 5:
            records = storage.find({'connector': 'collectd'})
            if len(records):
                break
            i += 1
            sleep(5)

        if not len(records):
            raise Exception("Collectd2event don't work ...")



    def test_99_Disconnect(self):
        clean()
        myamqp.stop()
        myamqp.join()

if __name__ == "__main__":
    main(verbosity=2)
