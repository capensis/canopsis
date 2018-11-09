#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Healthcheck manager.
"""

from __future__ import unicode_literals

from canopsis.common.amqp import AmqpConnection
from canopsis.common.collection import MongoCollection
from canopsis.common.mongo_store import MongoStore
from canopsis.common.redis_store import RedisStore
from canopsis.logger import Logger
from canopsis.models.healthcheck import Healthcheck, ServiceState


class HealthcheckManager(object):
    """
    Action managment.
    """
    LOG_PATH = 'var/log/healthcheck.log'

    CHECK_COLLECTIONS = ['default_entities', 'periodical_alarm']

    def __init__(self, logger):
        self.logger = logger

        self.db_store = MongoStore.get_default()
        self.cache_store = RedisStore.get_default()
        self.amqp_url, self.amqp_exchange = AmqpConnection.parse_conf()

    @classmethod
    def provide_default_basics(cls):
        """
        Provide logger, config, storages...

        ! Do not use in tests !

        :rtype: (logging.Logger)
        """
        logger = Logger.get('healthcheck', cls.LOG_PATH)

        return (logger,)

    def check_amqp(self):
        """
        Check if amqp service is available.

        :rtype: ServiceState
        """
        try:
            with AmqpConnection(self.amqp_url) as amqp_conn:
                channel = amqp_conn._channel
                try:
                    channel.basic_publish(self.amqp_exchange, '', 'test')
                except Exception as exc:
                    msg = 'Failed to publish: {}'.format(exc)
                    return ServiceState(message=msg)

                if not amqp_conn._connection.is_open:
                    return ServiceState(message='Connection is not opened')
                if not channel.is_open:
                    return ServiceState(message='Channel is not opened')
        except Exception as exc:
            return ServiceState(message='Failed to connect: {}'.format(exc))

        return ServiceState()

    def check_cache(self):
        """
        Check if cache service is available.

        :rtype: ServiceState
        """
        message = "import this"
        try:
            response = self.cache_store.echo(message)
        except Exception as exc:
            msg = 'Cache service crash on echo: {}'.format(exc)
            return ServiceState(message=msg)

        if response != message:
            return ServiceState(message='Failed to validate echo')

        return ServiceState()

    def check_db(self):
        """
        Check if database service is available.

        :rtype: ServiceState
        """
        existing_cols = self.db_store.client.collection_names()
        for collection_name in self.CHECK_COLLECTIONS:
            # Existence test
            if collection_name not in existing_cols:
                msg = 'Missing collection {}'.format(collection_name)
                return ServiceState(message=msg)

            # Read test
            collection = self.db_store.get_collection(name=collection_name)
            mongo_collection = MongoCollection(collection)
            try:
                mongo_collection.find({}, limit=1)
            except Exception as exc:
                return ServiceState(message='Find error: {}'.format(exc))

        return ServiceState()

    def check_engines(self):
        """
        Check if engines are available.

        :rtype: ServiceState
        """
        return ServiceState()

    def check_time_series(self):
        """
        Check if time_series service is available.

        :rtype: ServiceState
        """
        return ServiceState()

    def check(self, criticals=None):
        """
        Check all services.

        :param list criticals: service names considered as critical (for overall)
        :rtype: dict
        """
        check = Healthcheck(
            amqp=self.check_amqp(),
            cache=self.check_cache(),
            database=self.check_db(),
            engines=self.check_engines(),
            time_series=self.check_time_series(),
            criticals=criticals
        )
        return check.to_dict()
