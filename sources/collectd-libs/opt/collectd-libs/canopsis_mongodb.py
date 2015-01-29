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

plugin_name = "canopsis_mongodb"

storage = None

namespaces = ['object', 'cache', 'events', 'events_log', 'ack', 'entities',
    'perfdata3', 'perfdata2', 'perfdata2_daily']


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

    from canopsis.old.storage import get_storage
    from canopsis.old.account import Account

    global storage
    root = Account(user="root", group="root")
    storage = get_storage(account=root, namespace='object')


def config_callback(config):
    log('Config plugin')


def read_callback(data=None):
    for namespace in namespaces:
        put_value(namespace+"_size", storage.get_namespace_size(namespace))

    ## Pyperfstore
    perfdata2_size = storage.get_namespace_size("perfdata2_bin.chunks")
    perfdata2_size += storage.get_namespace_size("perfdata2_bin.files")
    perfdata2_size += storage.get_namespace_size("perfdata2")
    perfdata2_size += storage.get_namespace_size("perfdata2_daily")
    put_value("perfdata2_size", perfdata2_size)

    perfdata3_size = storage.get_namespace_size("perfdata3_bin.chunks")
    perfdata3_size += storage.get_namespace_size("perfdata3_bin.files")
    perfdata3_size += storage.get_namespace_size("perfdata3")
    put_value("perfdata3_size", perfdata3_size)

    ## Briefcase
    files_size = storage.get_namespace_size("binaries.chunks")
    files_size += storage.get_namespace_size("binaries.files")
    files_size += storage.get_namespace_size("files")
    put_value("files_size", files_size)

    size = perfdata2_size + perfdata3_size + files_size

    for namespace in namespaces:
        size += storage.get_namespace_size(namespace)
    put_value("size", size)

### MAIN ###
collectd.register_config(config_callback)
collectd.register_init(init_callback)
collectd.register_read(read_callback)
