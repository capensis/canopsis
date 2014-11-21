#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
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

from canopsis.old.record import Record
from canopsis.old.storage import Storage
from canopsis.old.account import Account

from time import time


class Cache(object):
    def __init__(self, storage, namespace=None):
        super(Cache, self).__init__()

        self.storage = storage

        if not namespace:
            namespace = storage.namespace
        self.namespace = namespace

    def make_record(self, _id):
        record = Record()
        record.type = "cache"
        #record._id = 'cache.'+_id
        record._id = _id
        record.access_owner = ['r', 'w']
        record.access_group = []
        record.access_other = []
        record.access_unauth = []
        return record

    def deco(self, _id='', freshness=10, args=[], account=None):
        def dec(func):
            def gave(*k, **a):
                try:
                    logger = k[0].logger
                except:
                    logger = None

                cid = _id
                for arg in k:
                    if isinstance(arg, int) or isinstance(arg, basestring):
                        cid += "." + str(arg)

                for arg in args:
                    try:
                        if isinstance(arg, int) or isinstance(arg, basestring):
                            cid += "." + str(a[arg])
                    except:
                        pass

                data = self.get(cid, freshness, account)
                if data:
                    if logger:
                        logger.debug('   + Get from cache (%s)...' % cid)
                    return data

                data = func(*k, **a)
                self.put(cid, data, account)
                if logger:
                    logger.debug('   + Update cache (%s)...' % cid)
                return data
            return gave
        return dec

    def remove(self, _id, account=None):
        self.storage.remove(
            'cache.' + _id, namespace=self.namespace, account=account)

    def put(self, _id, data, account=None):
        self.remove(_id)
        record = self.make_record('cache.' + _id)
        record.data = {'d': data}
        self.storage.put(record, namespace=self.namespace, account=account)

    def get(self, _id, freshness=10, account=None):
        if freshness == -1:
            return None

        try:
            record = self.storage.get(
                'cache.' + _id, namespace=self.namespace, account=account)

            if record.write_time < (time() - freshness) and freshness != 0:
                self.remove(_id)
                return None
            else:
                return record.data['d']
        except:
            return None

## Cache Cache
CACHE = None


def get_cache(storage=None):
    global CACHE

    if not storage:
        storage = Storage(Account(), namespace='cache')

    if CACHE:
        return CACHE
    else:
        #global CACHE
        CACHE = Cache(storage, namespace='cache')
        return CACHE
