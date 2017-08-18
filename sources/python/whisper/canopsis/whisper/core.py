# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

from canopsis.common.init import basestring
from canopsis.storage.core import Storage, DataBase, Cursor

from influxdb import InfluxDBClient, InfluxDBClusterClient

from sys import getsizeof


class InfluxDBBase(DataBase):
    """Manage access to influxDB."""

    __protocol__ = 'influxdb'

    DEFAULT_HOST = 'localhost'
    DEFAULT_PORT = 8086
    DEFAULT_USER = 'root'
    DEFAULT_PWD = 'root'
    DEFAULT_DB = None
    DEFAULT_SSL = False
    DEFAULT_TIMEOUT = None
    DEFAULT_PROXIES = None

    def __init__(
            self, host=DEFAULT_HOST, port=DEFAULT_PORT, user=DEFAULT_USER,
            pwd=DEFAULT_PWD, db=DEFAULT_DB, ssl=DEFAULT_SSL,
            conn_timeout=DEFAULT_TIMEOUT, proxies=DEFAULT_PROXIES,
            *args, **kwargs
    ):

        super(InfluxDBBase, self).__init__(
            host=host, port=port, user=user, pwd=pwd, db=db, ssl=ssl,
            conn_timeout=in_timeout, proxies=proxies, *args, **kwargs
        )

    def _connect(self, *args, **kwargs):

        result = None

        connection_args = {}

        conncls = InfluxDBClient

        if self.db:
            connection_args['database'] = self.db

        # if self host is given
        if self.host:
            connection_args['host'] = self.host
        # if self port is given
        if self.port:
            connection_args['port'] = self.port
        # if self replica set is given
        if self.replicaset:
            conncls = InfluxDBClusterClient
            if self.host:
                connection_args['hosts'] = self.host.split(',')
                del connection_args['host']

        if self.ssl:
            connection_args['ssl'] = True

        if self.user:
            connection_args['user'] = self.user

        if self.pwd:
            connection_args['password'] = self.pwd

        if self.conn_timeout:
            connection_args['timeout'] = self.conn_timeout

        if self.proxies:
            connection_args['proxies'] = self.proxies

        self.logger.debug(u'Trying to connect to {0}'.format(connection_args))

        try:
            result = conncls(**connection_args)

        except Exception as cex:
            self.logger.error(
                'Raised {2} during connection attempting to {0}:{1}.'.
                format(self.host, self.port, cex)
            )

        else:
            try:
                result.create_database(self.db)

            except:
                pass

            if self.retention:
                try:
                    result.create_retention_policy(
                        name='{0}'.format(self), duration=self.retention,
                        replication=0 if self.replicaset is None else 1
                    )

                except:
                    pass

            if self.user:
                try:
                    result.create_user(self.user, self.pwd, False)

                except:
                    pass

        return result

    def _disconnect(self, *args, **kwargs):

        if self._conn is not None:
            self._conn = None

    def connected(self, *args, **kwargs):

        result = self._conn is not None

        return result

    def size(self, table=None, criteria=None, *args, **kwargs):

        result = 0

        if table is None:
            table = self.table

        query = paramstoquery(ids='COUNT(*)', query=criteria)

        result = len(self._conn.query(query).get_points()) * getsizeof(0)

        return result

    def drop(self, table=None):

        self._conn.drop_database(table)


class InfluxDBCursor(Cursor):
    """InfluxDB cursor."""

    def __len__(self):

        return len(self._cursor)

    def __iter__(self):

        return iter(self._cursor)

    def __getitem__(self, index):

        return self._cursor[index]


class InfluxDBStorage(InfluxDBDataBase, Storage):
    """Influxdb Storage."""

    def get_elements(
            self,
            ids=None, query=None, limit=0, skip=0, sort=None, with_count=False,
            projection=None, *args, **kwargs
    ):

        _query = paramstoquery(
            ids=ids, query=query, limit=limit, skip=skip, sort=sort,
            projection=projection
        )

        cursor = InfluxDBCursor(self._conn.query(_query))

        one_element = isinstance(ids, basestring)

        if one_element:

            if cursor:
                for key in cursor:
                    result = cursor[key]

            else:
                result = None

        # if with_count, add count to the result
        if with_count:
            # calculate count
            count = len(cursor)
            result = result, count

        return result

    def distinct(self, field, query):

        raise NotImplementedError()

    def find_elements(
            self, query=None, limit=0, skip=0, sort=None, projection=None,
            with_count=False,
            *args, **kwargs
    ):

        return self.get_elements(
            query=query,
            limit=limit,
            skip=skip,
            sort=sort,
            with_count=with_count,
            projection=projection,
            *args, **kwargs
        )

    def put_element(self, element, _id=None, cache=False):

        points = {
            'measurement': _id,
        }

        points.update(element)

        return self._conn.write_points(
            points=points, batch_size=self.cache_size if cache else 0
        )

    def remove_elements(
            self, ids=None, _filter=None, cache=False, *args, **kwargs
    ):

        self._conn.delete_series(measurement=ids)


def paramstoquery(
    ids=None, projection=None, skip=None, sort=None, limit=None, query=None
):
    """transform storage parameters to a influxdb query."""

    #construct from
    if ids is None:
        _from = 'FROM *'

    else:
        one_id = isinstance(ids, basestring)
        if one_id:
            _from = 'FROM {0}'.format(ids)

        else:
            _from = 'FROM {0}'.format(ids[0])
            for _id in ids:
                _from += ',{0}'.format(_id)

    # construct projection
    if projection is None:
        _select = 'SELECT *'

    else:
        one_projection = isinstance(projection, basestring)

        if one_projection:
            _select = projection

        else:
            _select = 'SELECT {0}'.format(projection[0])
            for prj in projection[1:]:
                _select += ',{0}'.format(prj)

    # construct offset
    if skip:
        _offset = 'OFFSET {0}'.format(skip)

    else:
        _offset = ''

    # construct order by
    if sort is None:
        _sort = ''

    else:
        one_sort = isinstance(sort[0], basestring)
        if one_sort:
            _sort = 'SORT BY {0} {1}'.format(sort[0], sort[1])

        else:
            _sort = 'SORT BY {0} {1}'.format(sort[0][0], sort[0][1])
            for _so in sort:
                _sort = 'SORT BY {0} {1}'.format(_so[0], _so[1])

    # construct limit
    _limit = 'LIMIT {0}'.format(limit) if limit else ''

    # construct where
    if query is None:
        _where = ''

    else:
        _where = querytosql(query=query)

    result = '{0} {1} {2} {3} {4} {5}'.format(
        _select, _from, _where, _limit, _sort, _offset
    )

    return result


def querytosql(query):
    """Transform a storage query to a SQL WHERE expression."""

    result = 'WHERE' if query else ''

    for param in query:

        value = query[param]

        if isinstance(value, dict):

            for operator in value:

                if operator in ('$or', '$and'):

                    raise NotImplementedError()

                else:

                    value = value[operator]
                    operator = OPERATORS[operator]

        else:
            operator = '='

        result = '{0} {1} {2} {3},'.format(result, param, operator, value)

    if result:
        result = result[:-1]

    return result

OPERATORS = {
    '$eq': '=',
    '$lt': '<',
    '$lte': '<=',
    '$gt': '>',
    '$gte': '>=',
    '$regex': 'LIKE'
}
