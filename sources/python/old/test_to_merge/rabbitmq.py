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

import unittest
import threading
import time
import json
import logging

from canopsis.old.rabbitmq import Amqp

logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s %(name)s %(levelname)s %(message)s',
                    )

msgbody = {
    "type": "check",
    "source_name": "Central",
    "source_type": "host",
    "timestamp": "1307518560",
    "host_name": "localhost16",
    "check_type": "0",
    "current_attempt": "1",
    "max_attempts": "10",
    "state_type": "1",
    "state": "0",
    "execution_time": "4.035",
    "latency": "0.218",
    "command_name": "check-host-alive",
    "output": "PING OK -  Paquets perdus = 0%, RTA = 0.04 ms",
    "long_output": "",
    "perf_data": "rta=0.037000ms;3000.000000;5000.000000;0.000000 pl=0%;80;100;0"
}

myamqp = None
rcvmsgbody = None


class KnownValues(unittest.TestCase):
    def setUp(self):
        self.rcvmsgbody = None

    def test_1_Init(self):
        global myamqp
        myamqp = Amqp()

    def test_2_CreateQueue_and_Bind(self):
        global myamqp
        myamqp.add_queue("unit_test", "unit_test.#", self.on_message)

    def test_3_Connect(self):
        global myamqp
        myamqp.start()

    def on_message(self, body, msg):
        rk = msg.delivery_info['routing_key']
        print "Receive message from %s ..." % rk
        if rk == "unit_test.testmessage":
            global rcvmsgbody
            rcvmsgbody = body
            print rcvmsgbody

    def test_4_PublishMessage(self):
        time.sleep(5)
        msg = msgbody
        myamqp.publish(msg, "unit_test.testmessage")

    def test_5_CheckReceiveInQueue(self):
        start = time.time()
        end = start + 20.0
        while not rcvmsgbody:
            time.sleep(0.1)
            if time.time() > end:
                break

        duration = time.time() - start
        print("Receive message in", duration, "ms")
        if rcvmsgbody != msgbody:
            print "msgbody:\t", msgbody
            print "rcvmsgbody:\t", rcvmsgbody
            raise NameError('Received Event is not conform')

    def test_99_Disconnect(self):
        global myamqp
        myamqp.stop()
        myamqp.join()

if __name__ == "__main__":
    unittest.main(verbosity=2)
