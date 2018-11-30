#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from time import time, sleep
from pymongo import errors
from redlock import Redlock

from canopsis.common.mongo_store import MongoStore
from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.logger import Logger
import os


class AlertLock(object):

    LOG_PATH = 'var/log/alert_lock.log'
    LOCK_COLLECTION = 'lock'

    @classmethod
    def provide_default_basics(cls):
        """
            provide default basics
        """
        conf_store = Configuration.load(MongoStore.CONF_PATH, Ini)

        mongo = MongoStore(config=conf_store)
        lock_collection = mongo.get_collection(name=cls.LOCK_COLLECTION)

        logger = Logger.get('lock', cls.LOG_PATH)

        return (logger, lock_collection)

    def __init__(self, logger, lock_collection):
        """
            AlertLock constructor
        """
        self.logger = logger
        self.lock_collection = lock_collection

    def lock(self, entity_id):
        """
            create a document in lock collection
        """
        document = {
            "_id": entity_id,
            "timestamp": time()
        }
        try:
            self.lock_collection.insert(document)
        except errors.DuplicateKeyError:
            sleep(0.2)
            self.clear_old_locks()
            self.lock(entity_id)

    def unlock(self, entity_id):
        """
            remove lock documentt
        """
        self.lock_collection.remove({'_id': entity_id})

    def clear_old_locks(self):
        """
            remove locks older than 13 seconds
        """
        self.lock_collection.remove({'timestamp': {'$lt': time() - 13}})


class AlertLockRedis(object):

    LOG_PATH = 'var/log/alert_lock.log'
    LOCK_COLLECTION = 'lock'
    CONF_PATH = 'etc/common/redis_store.conf'
    CONF_SECTION = 'DATABASE'
    DEFAULT_DB_HOST = 'localhost'
    DEFAULT_DB_PORT = '6379'
    DEFAULT_DB_NUM = '0'

    @classmethod
    def provide_default_basics(cls):
        """
            provide default basics
        """
        config = Configuration.load(
            os.path.join(root_path, cls.CONF_PATH), Ini).get(cls.CONF_SECTION)
        redis_host = config.get('host', cls.DEFAULT_DB_HOST)
        redis_port = int(config.get('port', cls.DEFAULT_DB_PORT))
        redis_db_num = int(config.get('dbnum', cls.DEFAULT_DB_NUM))
        redlock = Redlock(
            [{'host': redis_host, 'port': redis_port, 'db': redis_db_num}])

        logger = Logger.get('lock', cls.LOG_PATH)

        return (logger, redlock)

    def __init__(self, logger, redlock):
        """
            AlertLock constructor
        """
        self.logger = logger
        self.redlock = redlock

    def lock(self, entity_id):
        """
            create a document in lock collection
        """
        lock_id = 'redlock_{0}'.format(entity_id)
        stop = self.redlock.lock(lock_id, 12000)
        while type(stop) == bool:
            sleep(0.2)
            stop = self.redlock.lock(lock_id, 12000)
        return stop

    def unlock(self, lock):
        """
            remove lock documentt
        """
        return self.redlock.unlock(lock)
