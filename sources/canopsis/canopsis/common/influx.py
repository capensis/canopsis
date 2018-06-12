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

from influxdb import InfluxDBClient
from influxdb.exceptions import InfluxDBClientError

from canopsis.common import root_path
from canopsis.confng import Configuration, Ini
from canopsis.confng.helpers import cfg_to_bool

INFLUXDB_CONF_PATH = 'etc/influx/storage.conf'
INFLUXDB_CONF_SECTION = 'DATABASE'

SECONDS = 1000000000


class InfluxDBOptions(object):
    """
    The InfluxDBOptions enumeration contains the names of the options that can
    be defined in the influxdb configuration file.
    """
    host = 'host'
    port = 'port'
    username = 'user'
    password = 'pwd'
    database = 'db'
    ssl = 'ssl'
    verify_ssl = 'verify_ssl'
    timeout = 'timeout'
    retries = 'retries'
    use_udp = 'use_udp'
    udp_port = 'udp_port'


def get_influxdb_client(conf_path=INFLUXDB_CONF_PATH,
                        conf_section=INFLUXDB_CONF_SECTION):
    """
    Read the influxdb database's configuration from conf_path, and return an
    InfluxDBClient for this database.

    If a database name is specified in the configuration file and this database
    does not exist, it will be automatically created.

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
        influxdb_client_args['udp_port'] = int(cfg[InfluxDBOptions.udp_port])

    client = InfluxDBClient(**influxdb_client_args)

    return client


def encode_tags(tags):
    """
    Encode a point's tags in utf-8.

    This is required because of a bug in influxdb-python<=2.12.0.
    """
    encoded_tags = {}
    for key, value in tags.items():
        key = key.encode('utf-8')
        if value:
            value = value.encode('utf-8')
        encoded_tags[key] = value

    return encoded_tags


# The two following functions are defined in the influx.line_protocol module of
# influxdb-python>=4.0.0
def quote_ident(value):
    """
    Quote provided identifier.
    """
    return "\"{}\"".format(value.replace("\\", "\\\\")
                                .replace("\"", "\\\"")
                                .replace("\n", "\\n"))


def quote_literal(value):
    """
    Quote provided literal.
    """
    return "'{}'".format(value.replace("\\", "\\\\")
                              .replace("'", "\\'"))


def quote_regex(value):
    """
    Quote provided regex.
    """
    return "/{}/".format(value.replace("\\", "\\\\")
                              .replace("/", "\\/"))
