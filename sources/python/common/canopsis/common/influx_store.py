#!/usr/bin/env python
# -*- coding: utf-8  -*-

"""
Manage influxdb connections.
"""

from __future__ import unicode_literals

from influxdb import InfluxDBClient
from influxdb.exceptions import InfluxDBClientError


"""
class GenericSeriesHelper(SeriesHelper):
    # Meta class stores time series helper configuration.
    class Meta:
        # The client should be an instance of InfluxDBClient.
        client = myclient
        # The series name must be a string. Add dependent fields/tags in curly brackets.
        series_name = 'events.stats.{server_name}'
        # Defines all the fields in this time series.
        fields = ['some_stat', 'other_stat']
        # Defines all the tags for the series.
        tags = ['server_name']
        # Defines the number of data points to store prior to writing on the wire.
        bulk_size = 5
        # autocommit must be set to True when using bulk_size
        autocommit = True
"""


class InfluxStore(object):
    """
    Distribute ready-to-use influx series.
    """

    CONF_PATH = 'etc/influx/storage.conf'
    CONF_CAT = 'DATABASE'

    DEFAULT_HOST = 'localhost'
    DEFAULT_PORT = 4444
    DEFAULT_USER = 'admin'
    DEFAULT_PWD = 'admin'
    DEFAULT_DB = None
    DEFAULT_TIMEOUT = None
    #DEFAULT_SSL = False
    #DEFAULT_PROXIES = None

    def __init__(self, config):
        """
        :param config: a configuration object
        """
        self.config = config

        conf = self.config.get(self.CONF_CAT, {})
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
        kwargs = {
            'host': self.host,
            'udp_port': self.port,
            'database': self.database,
            'username': self._user,
            'password': self._pwd,
            'timeout': self.timeout
        }
        try:
            self.client = InfluxDBClient(**kwargs)
        except InfluxDBClientError as ice:
            self.logger.error(
                'Raised {} during connection attempting to {}:{}.'.
                format(ice, self.host, self.port, ice)
            )

        try:
            self.client.create_database(self.database)
        except InfluxDBClientError:
            pass

        #self.client.create_retention_policy('my_retention_policy', '3d', 3, default=True)
