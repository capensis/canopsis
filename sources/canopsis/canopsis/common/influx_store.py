#!/usr/bin/env python
# -*- coding: utf-8  -*-

"""
Manage influxdb connections.
"""

from __future__ import unicode_literals

from influxdb import InfluxDBClient
from influxdb.exceptions import InfluxDBClientError

# TODO: passer en tcp au lieu d'udp


class InfluxStore(object):
    """
    Distribute ready-to-use influx store.
    """

    CONF_PATH = 'etc/common/influx_store.conf'
    CONF_CAT = 'DATABASE'

    DEFAULT_HOST = 'localhost'
    DEFAULT_PORT = 4444
    DEFAULT_USER = 'admin'
    DEFAULT_PWD = 'admin'
    DEFAULT_DB = 'canopsis'
    DEFAULT_TIMEOUT = None

    def __init__(self, logger, config):
        """
        :param logger: a logger object
        :param config: a configuration object
        """
        self.logger = logger
        self.config = config

        conf = config.get(self.CONF_CAT, {})
        self.database = conf.get('database', self.DEFAULT_DB)
        self.host = conf.get('host', self.DEFAULT_HOST)
        self.port = int(conf.get('port', self.DEFAULT_PORT))
        self.timeout = conf.get('timeout', self.DEFAULT_TIMEOUT)

        self._user = conf.get('user', self.DEFAULT_USER)
        self._pwd = conf.get('pwd', self.DEFAULT_PWD)

        self._connect()

    def _connect(self):
        """
        Connect to the server and create the database if needed.
        """
        try:
            self.client = InfluxDBClient(
                host=self.host,
                udp_port=self.port,
                database=self.database,
                username=self._user,
                password=self._pwd,
                timeout=self.timeout
            )
        except InfluxDBClientError as ice:
            self.logger.error(
                'Raised {} during connection attempting to {}:{}.'.
                format(ice, self.host, self.port, ice)
            )

        #self.client.create_retention_policy('my_retention_policy', '3d', 3, default=True)
