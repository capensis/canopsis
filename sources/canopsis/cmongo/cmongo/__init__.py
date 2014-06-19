#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

__version__ = '0.1'

from cstorage import Storage, DataBase

# import to remove if mongo
from pymongo import MongoClient, ASCENDING

from pymongo.errors import TimeoutError, OperationFailure, ConnectionFailure


class DataBase(DataBase):
    """
    Manage access to a mongodb.
    """

    def __init__(self, port=MongoClient.PORT, *args, **kwargs):

        super(DataBase, self).__init__(port=port, *args, **kwargs)

    def connect(self, *args, **kwargs):

        if not self.connected():

            self.logger.debug('Trying to connect to {0}:{1}'.format(
                self.host, self.port))

            try:
                args = {
                    'host': self.host,
                    'port': self.port,
                    'j': self.journaling,
                    'w': 1 if self.safe else 0
                }
                if self.ssl:
                    args.update({
                        'ssl': self.ssl,
                        'ssl_keyfile': self.ssl_key,
                        'ssl_certfile': self.ssl_cert
                    })

                self._conn = MongoClient(**args)

            except ConnectionFailure as e:
                self.logger.error(
                    'Raised {2} during connection attempting to {0}:{1}.'.
                    format(self.host, self.port, e))

            else:
                self._database = self._conn[self.db]

                if (self.user, self.pwd) != (None, None):
                    authenticate = self._db.authenticate(self.user, self.pwd)

                    if authenticate:
                        self.logger.debug("Connected on {0}:{1}".format(
                            self.host, self.port))
                        self._connected = True

                    else:
                        self.logger.error(
                            'Impossible to authenticate {0} on {1}:{2}'.format(
                                self.host, self.port))
                        self.disconnect()

                else:
                    self._connected = True
                    self.logger.debug("Connected on {0}:{1}".format(
                        self.host, self.port))

        return self.connected()

    def disconnect(self, *args, **kwargs):

        if getattr(self, '_conn', None)is not None:
            self._conn.close()
            self._conn = None

        self._connected = False

    def connected(self, *args, **kwargs):

        result = getattr(self, '_conn', None) is not None

        return result

    def size(self, table=None, criteria=None, *args, **kwargs):

        result = 0

        _backend = self._get_backend(backend=table)

        try:
            result = self._database.command("collstats", _backend)['size']

        except Exception as e:
            self.logger.warning(
                "Impossible to read Collection Size: {0}".format(e))
            result = None

        return result

    def drop(self, table=None):

        if table is None:
            for collection_name in self._database.collection_names(
                    include_system_collections=False):
                self.drop(table=collection_name)

        else:
            self._get_backend(backend=table).remove()

    def _get_backend(self, backend=None, *args, **kwargs):
        """
        Get a reference to a specific backend where name is input backend.
        If input backend is None, self.backend is used.

        :param backend: backend name. If None, self.backend is used.
        :type backend: basestring

        :returns: backend reference.

        :raises: NotImplementedError
        .. seealso: DataBase.set_backend(self, backend)
        """

        if getattr(self, '_database', None) is None:
            raise DataBase.DataBaseError(
                '{0} is not connected'.format(self))

        result = self._database[backend]

        return result


class Storage(DataBase, Storage):

    def __init__(self, data_type, *args, **kwargs):

        super(Storage, self).__init__(
            data_type=data_type, *args, **kwargs)

    def connect(self, *args, **kwargs):

        result = super(Storage, self).connect(*args, **kwargs)

        if result:
            indexes = self._get_indexes()

            self._backend = self._database[self.get_table()]

            for index in indexes:
                self._backend.ensure_index(index)

        return result

    def drop(self, *args, **kwargs):
        """
        Drop self table.
        """

        table = self.get_table()

        super(Storage, self).drop(table=table, *args, **kwargs)

    def get_elements(
        self, ids=None, limit=0, skip=0, sort=None, *args, **kwargs
    ):

        query = dict()

        if ids is not None:
            query['_id'] = {'$in': ids}

        cursor = self._find(query)

        if limit:
            cursor.limit(limit)
        if skip:
            cursor.skip(skip)
        if sort is not None:
            Storage._update_sort(sort)
            cursor.sort(sort)

        result = list(cursor)

        return result

    def remove_elements(self, ids, *args, **kwargs):

        self._remove({'_id': {'$in': ids}})

    def put_element(self, _id, element, *args, **kwargs):

        self._update(_id={'_id': _id}, document={'$set': element}, multi=False)

    def _element_id(self, element):

        return element['_id']

    def _get_indexes(self):
        """
        Get collection indexes. Must be overriden.
        """

        result = [
            [('_id', ASCENDING)]
        ]

        return result

    def _manage_query_error(self, result_query):
        """
        Manage mongo query error.
        """

        result = None

        if isinstance(result_query, dict):

            error = result_query.get("writeConcernError", None)

            if error is not None:
                self.logger.error(' error in writing document: {0}'.format(
                    error))

            error = result_query.get("writeError")

            if error is not None:
                self.logger.error(' error in writing document: {0}'.format(
                    error))

        else:

            result = result_query

        return result

    def _insert(self, document):

        result = self._run_command(command='insert', doc_or_docs=document)

        return result

    def _update(self, _id, document, multi=True):

        result = self._run_command(
            command='update', spec=_id, document=document,
            upsert=True, multi=multi)

        return result

    def _find(self, document=None, projection=None):

        result = self._run_command(command='find', spec=document,
            projection=projection)

        return result

    def _remove(self, document):

        result = self._run_command(command='remove', spec_or_id=document)

        return result

    def _run_command(self, command, **kwargs):
        """
        Run a specific command on a given backend.

        :param command: command to run
        :type command: str

        :param backend: backend to use
        :type backend: str
        """

        result = None

        try:
            backend_command = getattr(
                self._backend, command)
            result = backend_command(w=1 if self.safe else 0,
                wtimeout=self.wtimeout, **kwargs)

            self._manage_query_error(result)

        except TimeoutError:
            self.logger.warning(
                'Try to run command {0}({1}) on {2} attempts left'
                .format(command, kwargs, self._backend))

        except OperationFailure as of:
            self.logger.error('{0} during running command {1}({2}) of in {3}'
                .format(of, command, kwargs, self._backend))

        return result
