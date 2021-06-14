#!/usr/bin/env python
# -*- coding: utf-8  -*-

"""
Access data from mongodb through a timed cache.
"""

from __future__ import unicode_literals

from time import time

DEFAULT_TIMEOUT = 60


class EtherealData(object):
    """
    Representation of a data in mongodb, as seen recently.

    A deeper solution should be to use event queues to inform managers that
    the data has been updated.
    """

    _data = {}
    _last_update = 0

    def __init__(self, collection, filter_, timeout=DEFAULT_TIMEOUT):
        """
        :param Collection collection: a mongo collection object
        :param dict filter_: subset of the collection to access to
        :param int timeout: time before invalidating data
        """
        self.collection = collection
        self.filter_ = filter_
        self.timeout = timeout
        self._update_cache()

    def get(self, value, default=None):
        """
        Update cache if needed and retreive a value.

        :param value: th value to insert
        :param default: default value to return if not founded
        :rtype: builtin type|default
        """
        if int(time()) - self._last_update > self.timeout:
            self._update_cache()

        return self._data.get(value, default)

    def set(self, key, value):
        """
        Insert in db the key/value pair.

        :param key: the key
        :param value: the value
        """
        self._data[key] = value
        self.collection.update(self.filter_, {'$set': {key: value}}, upsert=True)

    def _update_cache(self):
        """
        Update local cache from db.
        """
        values = self.collection.find_one(self.filter_)
        if values is None:
            values = {}
        self._data = values
        self._last_update = int(time())
