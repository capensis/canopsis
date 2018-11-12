#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Healthcheck and ServiceState objects.
"""

from __future__ import unicode_literals

import time

OK_MSG = ''


class ServiceState(object):

    """
    State of a single service and description of the error.
    """

    def __init__(self, message=OK_MSG):
        """
        :param string message: error message associated (if state is down)
        """
        self.message = message
        # state is the up/down state of the service
        self.state = message == OK_MSG

    @property
    def value(self):
        """
        Return a message for the service.
        '' if it works, the error message if not.

        :rtype: string
        """
        return self.message


class Healthcheck(object):

    """
    Representation of a health check over services.
    """

    AMQP = 'amqp'
    CACHE = 'cache'
    DATABASE = 'database'
    ENGINES = 'engines'
    TIME_SERIES = 'time_series'

    SERVICES = [AMQP, CACHE, DATABASE, ENGINES, TIME_SERIES]

    OVERALL = 'overall'
    TIME = 'timestamp'

    def __init__(self, amqp, cache, database, engines, time_series,
                 criticals=None):
        """
        :param ServiceState amqp: state of amqp service
        :param ServiceState cache: state of cache service
        :param ServiceState db: state of db service
        :param ServiceState engines: state of engines
        :param ServiceState time_series: state of time_series service
        :param list criticals: a list of service names considered as critical
        :raises: TypeError if parameters are in a wrong type
        """
        self.amqp = amqp
        self.cache = cache
        self.database = database
        self.engines = engines
        self.time_series = time_series

        if criticals is None:
            self.criticals = self.SERVICES  # default all
        else:
            if not isinstance(criticals, list):
                raise TypeError('criticals must be a list: {}'
                                .format(criticals))
            self.criticals = criticals

        for service in self.SERVICES:
            if not isinstance(getattr(self, service), ServiceState):
                raise TypeError('{} is not a ServiceState'.format(service))

    def __str__(self):
        return '{}'.format(self._id)

    def __repr__(self):
        return '<Healthcheck {}>'.format(self.__str__())

    @property
    def overall(self):
        """
        Check critical service states (all by default).

        :rtype: bool
        """
        services = []
        for service in self.criticals:
            attr = getattr(self, service)
            if isinstance(attr, ServiceState):
                services.append(attr.state)

        return all(services)

    def to_dict(self):
        """
        Give a dict representation of the healthcheck object.

        :rtype: dict
        """
        dictionnary = {
            self.AMQP: self.amqp.value,
            self.ENGINES: self.engines.value,
            self.CACHE: self.cache.value,
            self.DATABASE: self.database.value,
            self.TIME_SERIES: self.time_series.value,
            self.OVERALL: self.overall,
            self.TIME: int(time.time())
        }

        return dictionnary
