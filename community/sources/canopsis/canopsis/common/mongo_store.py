#!/usr/bin/env python
# -*- coding: utf-8  -*-

"""
Manage mongodb connections.
"""

from __future__ import unicode_literals

import hashlib
import time

from canopsis.confng import Configuration, Ini
from pymongo import MongoClient, ReadPreference
from pymongo.errors import AutoReconnect

DEFAULT_HOST = 'localhost'
DEFAULT_PORT = 27017
DEFAULT_DB_NAME = 'canopsis'

singletons_cache = {}


class MongoStore(object):
    """
    Distribute ready-to-use mongo collections.
    """

    CONF_PATH = 'etc/common/mongo_store.conf'
    CONF_CAT = 'DATABASE'

    @staticmethod
    def get_default(from_singleton=True):
        """
        :returns: a defautl connection to Mongo using etc/common/mongo_store.conf
        :rtype: MongoStore
        """
        global singletons_cache

        cfg = Configuration.load(MongoStore.CONF_PATH, Ini)

        if from_singleton:
            cfg_values = cfg.get(MongoStore.CONF_CAT, {}).values()
            fingerprint = hashlib.md5('.'.join(sorted(cfg_values))).hexdigest()

            if fingerprint not in singletons_cache:
                singletons_cache[fingerprint] = MongoStore(cfg)

            return singletons_cache.get(fingerprint)

        return MongoStore(cfg)

    def __init__(self, config):
        """
        To use a replicaset, just use a list of hosts in the configuration.

        Example:

        host = host1:27017,host2:27017

        :param config dict: a configuration object
        """
        self.config = config
        conf = self.config.get(self.CONF_CAT, {})
        self.db_name = conf.get('db', DEFAULT_DB_NAME)
        self.host = conf.get('host', DEFAULT_HOST)
        try:
            self.port = int(conf.get('port', DEFAULT_PORT))
        except ValueError:
            self.port = DEFAULT_PORT

        self._db_uri = conf.get('db_uri')

        self.read_preference = getattr(
            ReadPreference,
            conf.get('read_preference', 'PRIMARY_PREFERRED'),
            ReadPreference.PRIMARY_PREFERRED
        )

        # missing from storage: journaling, sharding, retention ;;
        # cache_size, cache_autocommit, cache_ordered

        # missing from middleware: uri, protocol, data_type, data_scope, path,
        # auto_connect=true, safe=true, conn_timeout=20000, in_timeout=20000,
        # out_timeout=100, ssl=false, ssl_key, ssl_cert, user, pwd

        self._user = conf.get('user')
        self._pwd = conf.get('pwd')

        self._authenticated = False
        self._connect()

    def _connect(self):
        """
        Connect to the desired database.
        """
        self._authenticated = False
        if self._db_uri:
            self.conn = MongoClient(
                self._db_uri,
                w=1, j=True, read_preference=self.read_preference
            )

        else:
            db_uri = 'mongodb://{}:{}@{}:{}/{}'.format(
                self._user, self._pwd, self.host, self.port, self.db_name)
            self.conn = MongoClient(db_uri, w=1, j=True)

        self.client = self.get_database()

    def get_collection(self, name):
        """
        Return the desired collection.

        This function returns the raw pymongo Collection object.

        You must wrap it with MongoCollection if you want automatic AutoReconnect handling.

        :param name: the name of the collection
        :rtype: pymongo.collection.Collection
        """
        return MongoStore.hr(getattr, self.client, name)

    def get_database(self):
        """
        Returns a raw pymongo Database object.

        :rtype: pymongo.database.Database
        """
        return MongoStore.hr(self.conn.get_database)

    def alive(self):
        return self.conn is not None

    def authenticate(self):
        """
        Authenticate against the requested database.

        This method used to use MongoClient.authenticate, which is now
        deprecated. Some parts of the code still need authenticate to raise an
        exception when the authentication fails, so it now calls get_database,
        which also raises an exception when the authentication fails.

        .. deprecated:: 3.3.1
           The authentication check is already done in the get_database method.
        """
        self.get_database()
        self._authenticated = True

    @property
    def authenticated(self):
        """
        :rtype: bool

        .. deprecated:: 3.3.1
           This is set by the authenticate method, which is deprecated.
        """
        return self._authenticated

    def close(self):
        return MongoStore.hr(self.conn.close)

    def fsync(self, **kwargs):
        return MongoStore.hr(self.conn.fsync, **kwargs)

    @staticmethod
    def hr(func, *args, **kwargs):
        """
        hr means "Handle Reconnect". This function will loop forever until
        the pymongo driver has succeeded to reconnect to the database.

        This works only when using a replicaset or sharding setup.

        Reconnections are tried every second.
        """

        # try to work
        try:
            return func(*args, **kwargs)
        except AutoReconnect:
            pass

        # fast retry
        try:
            return func(*args, **kwargs)
        except AutoReconnect:
            pass

        # slow retries
        retries = 0
        allowed_retries = int(Configuration.load(
            MongoStore.CONF_PATH, Ini
        ).get(
            MongoStore.CONF_CAT, {}
        ).get('autoreconnect_retries', 20))

        while retries < allowed_retries:
            try:
                return func(*args, **kwargs)
            except AutoReconnect:
                pass

            time.sleep(1)
            retries += 1

        raise AutoReconnect(
            'failed to reconnect to MongoDB after {} tries, giving up.'.format(retries)
        )
