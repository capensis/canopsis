# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals
import os

from influxdb import InfluxDBClient as Client
from influxdb.exceptions import InfluxDBClientError

from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.confng.helpers import cfg_to_bool

INFLUXDB_CONF_PATH = 'etc/influx/storage.conf'
INFLUXDB_CONF_SECTION = 'DATABASE'

SECONDS = 1000000000

# This string contains the special characters that need to be escaped in a go
# regexp:
# https://github.com/golang/go/blob/c5d38b896df504e3354d7a27f7ad86fa9661ce6b/src/regexp/regexp.go#L628
REGEX_SPECIAL_CHARS = '\.+*?()|[]{}^$'


class InfluxDBOptions(object):
    """
    The InfluxDBOptions enumeration contains the names of the options that can
    be defined in the influxdb configuration file.

    The names of the options are set to be compatible with the
    canopsis.influxdb module, in order to be able to use the same configuration
    file.
    """
    host = 'host'
    port = 'tcp_port'
    username = 'user'
    password = 'pwd'
    database = 'db'
    ssl = 'ssl'
    verify_ssl = 'verify_ssl'
    timeout = 'timeout'
    retries = 'retries'
    use_udp = 'use_udp'
    udp_port = 'port'


class InfluxDBClient(Client):
    """
    This class is a wrapper around influxdb.InfluxDBClient that provides
    additional functionalities (initialization from configuration, continuous
    queries).
    """
    def __init__(self, logger, **kwargs):
        self.logger = logger
        self.database = kwargs.get('database')
        super(InfluxDBClient, self).__init__(**kwargs)

    @staticmethod
    def from_configuration(logger,
                           conf_path=INFLUXDB_CONF_PATH,
                           conf_section=INFLUXDB_CONF_SECTION):
        """
        Read the influxdb database's configuration from conf_path, and return
        an InfluxDBClient for this database.

        If a database name is specified in the configuration file and this
        database does not exist, it will be automatically created.

        :param str conf_path: the path of the file containing the database
            configuration.
        :param str conf_section: the section of the ini file containing the
            database configuration.
        :rtype: InfluxDBClient
        """
        influxdb_client_args = {}

        cfg = Configuration.load(
            os.path.join(root_path, conf_path), Ini
        ).get(conf_section, {})

        if InfluxDBOptions.host in cfg:
            influxdb_client_args['host'] = cfg[InfluxDBOptions.host]

        if InfluxDBOptions.port in cfg:
            influxdb_client_args['port'] = int(cfg[InfluxDBOptions.port])

        if InfluxDBOptions.username in cfg:
            influxdb_client_args['username'] = cfg[InfluxDBOptions.username]

        if InfluxDBOptions.password in cfg:
            influxdb_client_args['password'] = cfg[InfluxDBOptions.password]

        if InfluxDBOptions.database in cfg:
            influxdb_client_args['database'] = cfg[InfluxDBOptions.database]
        else:
            raise RuntimeError(
                "The {} option is required.".format(InfluxDBOptions.database))

        if InfluxDBOptions.ssl in cfg:
            influxdb_client_args['ssl'] = cfg_to_bool(cfg[InfluxDBOptions.ssl])

        if InfluxDBOptions.verify_ssl in cfg:
            influxdb_client_args['verify_ssl'] = cfg_to_bool(
                cfg[InfluxDBOptions.verify_ssl])

        if InfluxDBOptions.timeout in cfg:
            influxdb_client_args['timeout'] = int(cfg[InfluxDBOptions.timeout])

        if InfluxDBOptions.retries in cfg:
            influxdb_client_args['retries'] = int(cfg[InfluxDBOptions.retries])

        if InfluxDBOptions.use_udp in cfg:
            influxdb_client_args['use_udp'] = cfg_to_bool(
                cfg[InfluxDBOptions.use_udp])

        if InfluxDBOptions.udp_port in cfg:
            influxdb_client_args['udp_port'] = int(cfg[
                InfluxDBOptions.udp_port])

        return InfluxDBClient(logger, **influxdb_client_args)

    def create_continuous_query(self,
                                name,
                                query,
                                resample_every=None,
                                resample_for=None,
                                overwrite=False):
        """
        Create a continuous query.

        See https://docs.influxdata.com/influxdb/v1.6/query_language/continuous_queries
        for details on continuous queries.

        :param str name: The name of the continous query
        :param str query: The InfluxQL query
        :param str resample_every:
        :param str resample_for:
        :param bool overwrite: True to overwrite the query if it already exists
        :rtype: ResultSet
        """
        resample_statement = ''
        if resample_every is not None and resample_for is not None:
            resample_statement = 'RESAMPLE EVERY {} FOR {}'.format(
                resample_every, resample_for)
        elif resample_every is not None or resample_for is not None:
            raise ValueError(
                'resample_every and resample_for should either be both defined '
                'or both undefined.')

        creation_query = (
            "CREATE CONTINUOUS QUERY {name} ON {database} "
            "{resample_statement} "
            "BEGIN "
            "{query} "
            "END"
        ).format(
            name=quote_ident(name),
            database=quote_ident(self.database),
            resample_statement=resample_statement,
            query=query)

        try:
            return self.query(creation_query, epoch='s')
        except InfluxDBClientError as error:
            # I could not find a better way to catch this specific error
            if (error.content == 'continuous query already exists'
                and overwrite):
                # The continuous query already exists, recreate it.
                self.logger.info(
                    'A different continuous query already exists with this '
                    'name, overwriting it.')
                self.drop_continuous_query(name)
                return self.query(creation_query, epoch='s')
            else:
                # This is a different error, raise it.
                raise

    def drop_continuous_query(self, name):
        """
        Drop a continuous query.

        :param str name: The name of the continuous query.
        :rtype: ResultSet
        """
        query = 'DROP CONTINUOUS QUERY {name} ON {database}'.format(
            name=quote_ident(name),
            database=quote_ident(self.database))
        return self.query(query)

    def write_points(self, points):
        """Write points into measurements.

        The original InfluxDBClient.write_points method fails to write points
        with tags that contain a newline character (\n) [1]. This is a wrapper
        arround the InfluxDBClient.write_points that escapes newline characters
        from the tags of the point to prevent this error from happening.

        WARNING: when this issue is fixed upstream (with this PR [2] for
        example), this wrapper should be remove before upgrading
        influxdb-python. If not, this will prevent the stats of entities with a
        newline from being computed.

        [1]: https://github.com/influxdata/influxdb-python/issues/632
        [2]: https://github.com/influxdata/influxdb-python/pull/716
        """
        for point in points:
            tags = point.get('tags')
            if tags:
                for key, value in tags.items():
                    tags[key] = unicode(value).replace('\n', '\\n')

        super(InfluxDBClient, self).write_points(points)


# The two following functions are defined in the influx.line_protocol module of
# influxdb-python>=4.0.0
def quote_ident(value):
    """
    Quote provided identifier.

    :param str value: An influxdb identifier (e.g. a tag, field or measurement
    name).
    :rtype: str
    """
    return "\"{}\"".format(value.replace("\\", "\\\\")
                                .replace("\"", "\\\"")
                                .replace("\n", "\\n"))


def quote_literal(value):
    """
    Quote provided literal.

    :param str value: An influxdb literal (e.g. a tag value or a field value).
    :rtype: str
    """
    return "'{}'".format(value.replace("\\", "\\\\")
                              .replace("'", "\\'"))


def quote_regex(value):
    """
    Quote provided regex.

    :param str value: An influxdb regex.
    :rtype: str
    """
    return "/{}/".format(value.replace("/", "\\/"))


def escape_regex(pattern):
    """
    Escape the special characters in the pattern.

    This function does the same thing as re.escape, but for go regexp.
    """
    s = list(pattern)
    for i, c in enumerate(pattern):
        if c in REGEX_SPECIAL_CHARS:
            s[i] = "\\" + c
    return pattern[:0].join(s)
