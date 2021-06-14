#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
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

import collectd
import ConfigParser
import json
import urllib2
import os
import sys
import httplib
import base64
import urlparse
import socket

plugin_name = "canopsis_rabbitmq"

url = "None"

canopsis_exchanges = ['canopsis.events', 'canopsis.alerts']

#opener = None

filename = '~/etc/amqp.conf'
filename = os.path.expanduser(filename)

config = ConfigParser.RawConfigParser()

try:
    config.read(filename)
    section = 'master'
    amqp_host = config.get(section, "host")
    amqp_userid = config.get(section, "userid")
    amqp_password = config.get(section, "password")

    url = "http://%s:55672/api" % amqp_host
except Exception as err:
    log("Impossible to load configurations (%s) !" % err)


### Functions
def put_value(metric, value, type='gauge'):
    metric = collectd.Values(
        plugin=plugin_name,
        type=type,
        values=[value],
        type_instance=metric
    )
    metric.dispatch()


def log(msg):
    collectd.info("%s: %s" % (plugin_name, msg))


### Callbacks
def init_callback():
    log('Init plugin')

    if not url:
        log(' + url is not defined !')
        return


def config_callback(config):
    log('Config plugin')


def read_callback(data=None):
    #f = opener.open(url+"/exchanges")
    exchanges = api(
        "GET", "/api/exchanges",
        "", 'localhost', '15672',
        amqp_userid, amqp_password
    )

    try:
        exchanges = json.loads(exchanges)
        for exchange in exchanges:
            name = exchange['name']
            if name in canopsis_exchanges:
                name = name.split('.')[1]

                try:
                    message_stats_out = exchange['message_stats_out']
                    put_value(
                        '%s_msg_out' % name,
                        message_stats_out['publish'],
                        type='derive'
                    )

                except Exception as err:
                    #log("Impossible to put OUT values of %s (%s)" % (name, err))
                    pass

                try:
                    message_stats_in = exchange['message_stats_in']
                    put_value(
                        '%s_msg_in' % name,
                        message_stats_in['publish'],
                        type='derive'
                    )

                except Exception, err:
                    #log("Impossible to put IN values of %s (%s)" % (name, err))
                    pass

    except Exception as err:
        log("Impossible to read json data (%s)" % err)
        pass


def api(method, path, body, hostname, port, username, password):

        #TODO SSL
        ssl = False
        if ssl:
            conn = httplib.HTTPSConnection(hostname, port, ssl_key_file, ssl_cert_file)
        else:
            conn = httplib.HTTPConnection(hostname, port)
        headers = {"Authorization": "Basic " + base64.b64encode(username + ":" + password)}

        if body != "":
            headers["Content-Type"] = "application/json"
        try:
            conn.request(method, path, body, headers)
        except socket.error, e:
            log("Could not connect: {0}".format(e))
            return []
        resp = conn.getresponse()
        if resp.status == 400:
            log(json.loads(resp.read())['reason'])
            return []
        if resp.status == 401:
            log("Access refused: {0}".format(path))
            return []
        if resp.status == 404:
            log("Not found: {0}".format(path))
            return []
        if resp.status == 301:
            url = urlparse.urlparse(resp.getheader('location'))
            [host, port] = url.netloc.split(':')
            raise Exception("Error while retreiving rabbit information")
        if resp.status < 200 or resp.status > 400:
            raise Exception("Received %d %s for path %s\n%s"
                            % (resp.status, resp.reason, path, resp.read()))
        return resp.read()

### MAIN ###
collectd.register_config(config_callback)
collectd.register_init(init_callback)
collectd.register_read(read_callback)
