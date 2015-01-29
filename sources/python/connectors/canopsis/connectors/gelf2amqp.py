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

from zlib import decompress

from json import loads

from socket import socket, AF_INET, SOCK_DGRAM, IPPROTO_UDP, timeout

from pyparsing import Word, alphas, Suppress, nums, Optional, Regex

from canopsis.old.rabbitmq import Amqp
from canopsis.old.init import Init
from canopsis.old.event import forger, get_routingkey

DAEMON_NAME = 'gelf2amqp'

init = Init()
logger = init.getLogger(DAEMON_NAME)
handler = init.getHandler(logger)

gelf_port = 5555
gelf_interface = "0.0.0.0"

myamqp = None

#sys.path.append(os.path.expanduser("~/opt/event-brokers/nagios/api"))

## Init parser
integer = Word(nums)
serverDateTime = Regex("\S\S\S\s*\d\d?\s*\d\d:\d\d:\d\d")
hostname = Word(alphas + nums + "_" + "-")
daemon = Word(alphas + "/" + "-" + "_") + Optional(
    Suppress("[") + integer + Suppress("]")) + Suppress(":")
output = Regex(".*")
syslog_parser = serverDateTime + hostname + daemon + output


def gelf_uncompress(data):
    logger.debug("Uncompress GELF data ...")

    typeMsg = hex(ord(data[0])) + hex(ord(data[1]))

    if '0x780x9c' in typeMsg:
        compress_type = "ZLIB"
        gelf = loads(str(decompress(data)))

    elif '0x1f0x8b' in typeMsg:
        compress_type = "GZIP"
        gelf = "gzip"

    else:
        GELFCompressType = "NONE"
        gelf = "none"

    logger.debug(" + Gelf: %s" % gelf)
    return gelf


def gelf_level_to_state(gelf_level):
    gelf_level = int(gelf_level)

    #0 Emerg (emergency)
    #1 Alert
    #2 Crit (critical)
    #3 Err (error)
    #4 Warning
    #5 Notice
    #6 Info (informational)
    #7 Debug
    #8 none

    if gelf_level < 4:
        state = 2
    elif gelf_level < 5:
        state = 1
    else:
        state = 0

    return state


def wait_gelf_udp(on_log):
    s = socket(family=AF_INET, type=SOCK_DGRAM, proto=IPPROTO_UDP)

    s.bind((gelf_interface, gelf_port))

    logger.info(
        "Wait GELF data from UDP (%s:%s)" % (gelf_interface, gelf_port))
    try:
        while handler.status():
            try:
                data, peer = s.recvfrom(1024)
            except timeout:
                continue

            try:
                gelf = gelf_uncompress(data)
            except Exception as err:
                logger.error("Data: %s" % data)
                logger.error("Impossible to uncompress gelf data: '%s'" % err)
                continue

            try:
                on_log(gelf)
            except Exception as err:
                logger.error("Impossible to send log to Canopsis: '%s'" % err)
                continue

    except Exception as err:
        logger.error("Exception: '%s'" % err)

    logger.info("Close UDP socket")


def parse_syslog(message):
    logger.debug("Parse message ...")
    logger.debug(" + Raw: %s" % message)
    message = syslog_parser.parseString(message)
    logger.debug(" + Parsed: %s" % message)
    if len(message) < 5:
        result = {
            'timestamp': message[0], 'component': message[1],
            'resource': message[2], 'output': message[3]}
    else:
        result = {
            'timestamp': message[0], 'component': message[1],
            'resource': message[2], 'output': message[4], 'pid': message[3]}

    return result


def on_log(gelf):

    short_message = gelf.get('short_message', '')

    try:
        message = parse_syslog(short_message)
    except Exception as err:
        logger.error('Impossible to parse message ("%s")' % err)
        logger.error('short_message: %s' % short_message)
        return

    long_output = gelf.get('full_message', '')

    state = gelf_level_to_state(gelf.get('level', 6))

    try:
        timestamp = int(gelf['timestamp'])
    except:
        timestamp = None

    output = message['output']
    resource = message['resource']

    #component = str(gelf['host'])
    component = message['component']

    source_type = 'resource'

    event = forger(
            connector='gelf',
            connector_name=DAEMON_NAME,
            component=component,
            resource=resource,
            timestamp=timestamp,
            source_type=source_type,
            event_type='log',
            state=state,
            output=output,
            long_output=long_output)

    event['level'] = gelf['level']
    event['facility'] = gelf['facility']

    logger.debug('Event: %s' % event)

    key = get_routingkey(event)
    myamqp.publish(event, key, myamqp.exchange_name_events)


def main():

    handler.run()

    # global
    global myamqp

    # Connect to amqp bus
    logger.debug("Start AMQP ...")
    myamqp = Amqp()
    myamqp.start()

    wait_gelf_udp(on_log)

    logger.debug("Stop AMQP ...")
    myamqp.stop()
    myamqp.join()
