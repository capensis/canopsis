#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Healthcheck manager.
"""

from __future__ import unicode_literals

from canopsis.common.collection import MongoCollection
from canopsis.common.mongo_store import MongoStore
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

    @classmethod
    def provide_default_basics(cls):
        """
        Provide logger, config, storages...

        ! Do not use in tests !

        :rtype: (logging.Logger)
        """
        logger = Logger.get('healthcheck', cls.LOG_PATH)

        return (logger,)

    def check_db(self):
        """
        :returns: a working/not working status and a message string
        :rtype: bool, string
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

    def check(self, criticals):
        """
        Check all services.

        :param list criticals:
        :rtype: dict
        """
        # TODO: implements criticals
        check = Healthcheck(
            amqp=self.check_db(),
            cache=ServiceState(),  # TODO
            database=ServiceState(),  # TODO
            engines=ServiceState(),  # TODO
            time_series=ServiceState()  # TODO
        )
        return check.to_dict()
