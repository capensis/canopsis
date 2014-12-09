# -*- coding: utf-8 -*-
# --------------------------------
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

from canopsis.common.init import basestring
from canopsis.storage import Storage, DataBase
from canopsis.common.utils import isiterable

from pymongo import MongoClient

from pymongo.errors import TimeoutError, OperationFailure, ConnectionFailure,\
    DuplicateKeyError


class MongoDataBase(DataBase):
    """
    Manage access to a mongodb.
    """

    def __init__(
        self, host=MongoClient.HOST, port=MongoClient.PORT, *args, **kwargs
    ):

        super(MongoDataBase, self).__init__(
            port=port, host=host, *args, **kwargs
        )

    def _connect(self, *args, **kwargs):

        result = None

        connection_args = {}

        # if self host is given
        if self.host:
            connection_args['host'] = self.host

        if self.port:
            connection_args['port'] = self.port

        self.logger.debug('Trying to connect to %s' % (connection_args))

        connection_args['j'] = self.journaling
        connection_args['w'] = 1 if self.safe else 0

        if self.ssl:
            connection_args.update({
                'ssl': self.ssl,
                'ssl_keyfile': self.ssl_key,
                'ssl_certfile': self.ssl_cert
            })

        try:
            result = MongoClient(**connection_args)
        except ConnectionFailure as e:
            self.logger.error(
                'Raised {2} during connection attempting to {0}:{1}.'.
                format(self.host, self.port, e))
        else:
            self._database = result[self.db]

            if (self.user, self.pwd) != (None, None):

                authenticate = self._database.authenticate(
                    self.user, self.pwd)

                if authenticate:
                    self.logger.debug("Connected on {0}:{1}".format(
                        self.host, self.port))

                else:
                    self.logger.error(
                        'Impossible to authenticate {0} on {1}:{2}'.format(
                            self.host, self.port))
                    self.disconnect()
                    result = None

            else:
                self.logger.debug("Connected on {0}:{1}".format(
                    self.host, self.port))

        return result

    def _disconnect(self, *args, **kwargs):

        if self._conn is not None:
            self._conn.close()
            self._conn = None

    def connected(self, *args, **kwargs):

        result = self._conn is not None and self._conn.alive()

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


class MongoStorage(MongoDataBase, Storage):

    __register__ = True  #: register this class to middleware
    __protocol__ = 'mongodb'  #: register this class to the protocol mongodb

    ID = '_id'  #: ID mongo

    INSERT_COUNT = 'insert_count'  #: insert count

    def _connect(self, *args, **kwargs):

        result = super(MongoStorage, self)._connect(*args, **kwargs)

        if result:
            table = self.get_table()
            self._backend = self._database[table]
            # enable sharding
            if self.sharding:
                # on db
                self._database.command(enableSharding=self.db)
                # and on collection
                collection_full_name = '{0}.{1}'.format(self.db, table)
                self._database.command(
                    shardCollection=collection_full_name, key={'_id': 1}
                )

            # initialize self.to_insert
            self.to_insert = []
            self.insert_count = 30

            for index in self.all_indexes():
                try:
                    self._backend.ensure_index(index)
                except Exception as e:
                    self.logger.error(e)

        return result

    def __del__(self):
        """
        Insert last documents if there are such documents.
        """

        if self.to_insert:
            self.insert_count = 1
            self._insert(document=self.to_insert)

    def drop(self, *args, **kwargs):
        """
        Drop self table.
        """

        super(MongoStorage, self).drop(table=self.get_table(), *args, **kwargs)

    def get_elements(
        self,
        ids=None, query=None, limit=0, skip=0, sort=None, with_count=False,
        *args, **kwargs
    ):

        _query = {} if query is None else query.copy()

        one_element = False

        if ids is not None:
            if isiterable(ids, is_str=False):
                _query[MongoStorage.ID] = {'$in': ids}

            else:
                one_element = True
                _query[MongoStorage.ID] = ids

        cursor = self._find(_query)

        # set limit, skip and sort properties
        if limit:
            cursor.limit(limit)
        if skip:
            cursor.skip(skip)
        if sort is not None:
            MongoStorage._update_sort(sort)
            cursor.sort(sort)

        # calculate count
        count = cursor.count() if with_count else 0

        hint = self._get_hint(query=_query, cursor=cursor)

        if hint is not None:
            cursor.hint(hint)

        # TODO: enrich a cursor with methods to use it such as a tuple
        result = list(cursor)

        if one_element:
            if result:
                result = result[0]
            else:
                result = None

        # if with_count, add count to the result
        if with_count:
            result = result, count

        return result

    def _get_hint(self, query, cursor):
        """
        Get the best hint on input cursor for input query and returns it.
        """

        result = None

        # search for the best hint
        if cursor is not None and query:
            index = None
            # maximize the best hint related to query size
            query_len = len(query)
            # find the right index
            max_correspondance = 0
            # iterate on all indexes
            for self_index in self.all_indexes():
                # find the higher correspondance score
                correspondance = 0
                for index_value in self_index:
                    if index_value[0] in query:
                        correspondance += 1
                    else:
                        break
                # increment an old score if a new one is better
                if correspondance > max_correspondance:
                    index = self_index
                    max_correspondance = correspondance

                # ends if correspondance equals len(query)
                if correspondance == query_len:
                    break

            # if index has been founded
            if index is not None:
                # construct the right hint
                result = []
                for item in index:
                    if isinstance(item, basestring):
                        result.append((item, 1))
                    elif isinstance(item, tuple):
                        result.append(item)

        return result

    def find_elements(
        self, query, limit=0, skip=0, sort=None, with_count=False,
        *args, **kwargs
    ):

        return self.get_elements(
            query=query,
            limit=limit,
            skip=skip,
            sort=sort,
            with_count=with_count)

    def remove_elements(self, ids=None, _filter=None, *args, **kwargs):

        query = {}

        if ids is not None:
            if isiterable(ids, is_str=False):
                query[MongoStorage.ID] = {'$in': ids}
            else:
                query[MongoStorage.ID] = ids

        if _filter is not None:
            query.update(_filter)

        self._remove(query)

    def put_element(self, _id, element, *args, **kwargs):

        return self._update(
            spec={MongoStorage.ID: _id}, document={'$set': element},
            multi=False)

    def bool_compare_and_swap(self, _id, oldvalue, newvalue):

        return self.val_compare_and_swap(_id, oldvalue, newvalue) == newvalue

    def val_compare_and_swap(self, _id, oldvalue, newvalue):

        try:
            result = self._run_command(
                'find_and_modify',
                query={MongoStorage.ID: _id, 'value': oldvalue},
                update={'$set': {'value': newvalue}},
                upsert=True)

        except DuplicateKeyError:
            result = None

        return oldvalue if not result else newvalue

    def _element_id(self, element):

        return element[MongoStorage.ID]

    def all_indexes(self, *args, **kwargs):

        result = super(MongoStorage, self).all_indexes(*args, **kwargs)

        result.append([(MongoStorage.ID, 1)])

        return result

    def _manage_query_error(self, result_query):
        """
        Manage mongo query error.

        Returns result_query if no error encountered. Else None.
        """

        result = result_query

        if isinstance(result_query, dict):

            error = result_query.get("writeConcernError", None)

            if error is not None:
                self.logger.error(' error in writing document: {0}'.format(
                    error))
                result = None

            error = result_query.get("writeError")

            if error is not None:
                self.logger.error(' error in writing document: {0}'.format(
                    error))
                result = None

        return result

    def _insert(self, document=None, **kwargs):
        result = None
        # if insertion count is required
        if self.insert_count >= 0:
            # test if count is greater than expected
            if len(self.to_insert) > self.insert_count:
                if document is not None:
                    self.to_insert.append(document)
                result = self._run_command(
                    command='insert', doc_or_docs=self.to_insert, **kwargs
                )
                self.to_insert = []
            else:
                # insert a new document in the bulk
                self.to_insert.append(document)

        else:  # else, run directly the command
            result = self._run_command(
                command='insert', doc_or_docs=document, **kwargs
            )

        return result

    def _update(self, spec, document, multi=True, upsert=True, **kwargs):

        result = self._run_command(
            command='update', spec=spec, document=document,
            upsert=upsert, multi=multi, **kwargs
        )

        return result

    def _find(self, document=None, projection=None, **kwargs):

        result = self._run_command(command='find', spec=document,
            projection=projection, **kwargs)

        return result

    def _remove(self, document, **kwargs):

        result = self._run_command(
            command='remove', spec_or_id=document, **kwargs)

        return result

    def _count(self, document=None, **kwargs):

        cursor = self._find(document=document, **kwargs)

        result = cursor.count(False)

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
            backend = self._get_backend(backend=self.get_table())
            backend_command = getattr(backend, command)
            w = 1 if self.safe else 0
            result = backend_command(w=w, wtimeout=self.out_timeout, **kwargs)

            result = self._manage_query_error(result)

        except TimeoutError:
            self.logger.warning(
                'Try to run command {0}({1}) on {2} attempts left'
                .format(command, kwargs, backend))

        except OperationFailure as of:
            self.logger.error('{0} during running command {1}({2}) of in {3}'
                .format(of, command, kwargs, backend))

        return result
