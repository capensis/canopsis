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

from uuid import uuid1

from canopsis.common.collection import CollectionError, MongoCollection
from canopsis.common.init import basestring
from canopsis.common.mongo_store import MongoStore
from canopsis.common.utils import isiterable
from canopsis.storage.core import Cursor, DataBase, Storage
from pymongo import MongoClient
from pymongo.bulk import BulkOperationBuilder
from pymongo.cursor import Cursor as _Cursor
from pymongo.errors import (DuplicateKeyError, NetworkTimeout,
                            OperationFailure, PyMongoError)
from pymongo.read_preferences import ReadPreference


class MongoDataBase(DataBase):
    """Manage access to a mongodb."""

    def __init__(
            self, host=MongoClient.HOST, port=MongoClient.PORT,
            read_preference=ReadPreference.SECONDARY_PREFERRED,
            *args, **kwargs
    ):

        super(MongoDataBase, self).__init__(
            port=port, host=host, *args, **kwargs
        )

    def _connect(self, *args, **kwargs):
        result = MongoStore.get_default()
        self._database = result.client
        self._conn = result
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
            result = MongoStore.hr(self._database.command, "collstats", _backend)['size']

        except Exception as ex:
            self.logger.warning(
                "Impossible to read Collection Size: {0}".format(ex))
            result = None

        return result

    def drop(self, table=None):

        if table is None:
            collections = MongoStore.hr(
                self._database.collection_names,
                include_system_collections=False
            )
            for collection_name in collections:
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
        if backend is None:
            raise Exception('none backend')

        return MongoCollection(self._conn.get_collection(backend))


class MongoStorage(MongoDataBase, Storage):

    __register__ = True  #: register this class to middleware
    __protocol__ = 'mongodb'  #: register this class to the protocol mongodb

    ID = '_id'  #: ID mongo
    TAGS = 'tags'  #: tags field name

    def _connect(self, *args, **kwargs):
        result = super(MongoStorage, self)._connect(*args, **kwargs)

        # initialize cache
        if not hasattr(self, '_cache'):
            self._cache = None

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


            if self.all_indexes() is not None:

                for index in self.all_indexes():
                    try:
                        self._backend.ensure_index(index)

                    except Exception as ex:
                        self.logger.error(ex)

        return result

    def _disconnect(self, *args, **kwargs):
        super(MongoStorage, self)._disconnect(*args, **kwargs)

        self.halt_cache_thread()

    def _new_cache(self, *args, **kwargs):
        backend = self._get_backend(self.get_table())
        result = BulkOperationBuilder(backend, self._cache_ordered)

        return result

    def _execute_cache(self, *args, **kwargs):

        return self._cache.execute()

    def drop(self, *args, **kwargs):
        super(MongoStorage, self).drop(table=self.get_table(), *args, **kwargs)

    def get_elements(
            self,
            ids=None, query=None, limit=0, skip=0, sort=None, with_count=False,
            hint=None, projection=None, tags=None,
            *args, **kwargs
    ):

        _query = {} if query is None else query.copy()

        one_element = isinstance(ids, basestring)

        if ids is not None:
            if one_element:
                _query[MongoStorage.ID] = ids

            else:
                _query[MongoStorage.ID] = {'$in': ids}

        if tags:
            _query[MongoStorage.TAGS] = tags

        cursor = self._find(_query, projection)

        # set limit, skip and sort properties
        if limit:
            cursor.limit(limit)

        if skip:
            cursor.skip(skip)

        if sort is not None:
            sort = Storage._resolve_sort(sort)
            if sort:
                cursor.sort(sort)

        if hint is None:
            hint = self._get_hint(query=_query, cursor=cursor)

        if hint is not None:
            cursor.hint(hint)

        result = MongoCursor(cursor)

        if one_element:
            result = result[0] if result else None

        # if with_count, add count to the result
        if with_count:
            # calculate count
            count = cursor.count()
            result = result, count

        return result

    def _get_hint(self, query, cursor):
        """Get the best hint on input cursor for input query and returns it."""

        result = None

        # search for the best hint if cursor is not None and query exists
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

    def distinct(self, field, query):

        return self.find_elements(query=query)._cursor.distinct(field)

    def find_elements(
            self, query=None, limit=0, skip=0, sort=None, projection=None,
            tags=None, with_count=False,
            *args, **kwargs
    ):

        return self.get_elements(
            query=query, limit=limit, skip=skip, sort=sort, tags=tags,
            projection=projection, with_count=with_count,
            *args, **kwargs
        )

    def remove_elements(
            self, ids=None, _filter=None, tags=None, cache=False,
            *args, **kwargs
    ):

        query = {}

        if ids is not None:
            if isiterable(ids, is_str=False):
                query[MongoStorage.ID] = {'$in': ids}

            else:
                query[MongoStorage.ID] = ids

        if tags:
            query[MongoStorage.TAGS] = tags

        if _filter is not None:
            query.update(_filter)

        return self._remove(query, cache=cache)

    def put_element(
        self, element, _id=None, tags=None, cache=False, *args, **kwargs
    ):

        if tags is not None:
            element.update(tags)

        if _id is None:
            _id = self._element_id(element)

        if _id is None:
            return self._insert(document=element, cache=cache)

        else:
            return self._update(
                spec={MongoStorage.ID: _id}, document={'$set': element},
                multi=False, cache=cache
            )

    def put_elements(self, elements, tags=None, *args, **kwargs):

        for element in elements:
            self.put_element(element=element, tags=tags)

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

    def aggregate(self, query, table=None, **kwargs):
        query = query or {}
        table = table or self.get_table()

        result = self._get_backend(backend=table).aggregate(query, **kwargs)

        return result

    def _element_id(self, element):

        return element.get(MongoStorage.ID, None)

    def all_indexes(self, *args, **kwargs):

        result = super(MongoStorage, self).all_indexes(*args, **kwargs)

        result.append([(MongoStorage.ID, 1)])

        return result

    def _insert(self, document=None, cache=False, **kwargs):

        if cache and self._cache is None:
            self._init_cache()

        cache_op = self._cache.insert if cache else None

        if isinstance(document, dict):
            query_kwargs = {'command': 'insert_one', 'document': document}
            if '_id' not in document:
                document['_id'] = str(uuid1())
        else:
            for i, doc in enumerate(document):
                if '_id' not in doc:
                    doc['_id'] = str(uuid1())
                    document[i] = doc

            query_kwargs = {'command': 'insert_many', 'documents': document}

        result = self._process_query(
            query_op=self._run_command,
            cache_op=cache_op,
            cache_kwargs={'document': document},
            query_kwargs=query_kwargs,
            cache=cache,
            **kwargs
        )

        return result.inserted_id

    def _update(
            self, spec, document, cache=False, multi=True, upsert=True,
            **kwargs
    ):

        if cache and self._cache is None:
            self._init_cache()

        if cache:
            cache_op = self._cache.find(selector=spec)
            if upsert:
                cache_op = cache_op.upsert()
            if multi:
                cache_op = cache_op.update
            else:
                cache_op = cache_op.update_one
        else:
            cache_op = None

        # catch any mongo command like $set or $addToSet so we do not
        # override them
        if len(document.keys()) > 1 or document.keys()[0][0] != '$':
            document = {'$set': document}

        result = self._process_query(
            cache_op=cache_op,
            query_op=self._run_command,
            cache_kwargs={'update': document},
            query_kwargs={
                'filter': spec,
                'command': 'update_many',
                'update': document,
                'upsert': upsert,
            },
            cache=cache,
            **kwargs
        )

        return {
            'updatedExisting': True if result is not None else False,
            'n': result.modified_count if result is not None else 0,
        }

    def _find(self, document=None, projection=None, **kwargs):

        result = self._run_command(
            'find', self.get_table(), document, projection, **kwargs
        )

        return result

    def _remove(self, document, cache=False, **kwargs):

        if cache and self._cache is None:
            self._init_cache()

        cache_op = self._cache.find(selector=document).remove if cache \
            else None

        result = self._process_query(
            query_op=self._run_command,
            cache_op=cache_op,
            cache_kwargs={},
            query_kwargs={'command': 'remove', 'spec_or_id': document},
            cache=cache,
            **kwargs
        )

        return result

    def _count(self, document=None, **kwargs):

        cursor = self._find(document=document, **kwargs)

        result = cursor.count(False)

        return result

    def _process_query(self, *args, **kwargs):
        result = super(MongoStorage, self)._process_query(*args, **kwargs)

        result = self._manage_query_error(result)

        return result

    def _manage_query_error(self, result_query):
        """Manage mongo query error.

        Returns result_query if no error encountered. Else None.
        """

        result = result_query

        if isinstance(result_query, dict):

            error = result_query.get("writeConcernError", None)

            if error is not None:
                self.logger.error(u' error in writing document: {0}'.format(
                    error))
                result = None

            error = result_query.get("writeError")

            if error is not None:
                self.logger.error(u' error in writing document: {0}'.format(
                    error))
                result = None

        return result

    def _run_command(self, command, table=None, *args, **kwargs):
        """Run a specific command on a given backend.

        :param str command: command to run.
        :param str table: table to use. self table if None.
        """

        result = None

        try:
            if table is None:
                table = self.get_table()
            # get pymongo raw collection
            backend = self._get_backend(backend=table).collection
            backend_command = getattr(backend, command)
            result = MongoStore.hr(backend_command, *args, **kwargs)

        except NetworkTimeout:
            self.logger.warning(
                'Try to run command {0}({1}) on {2} attempts left'
                .format(command, kwargs, backend))

        except (OperationFailure, CollectionError) as of:
            self.logger.error(u'{0} during running command {1}({2}) of in {3}'
                .format(of, command, kwargs, backend))

        return result


class MongoCursor(Cursor):
    """In charge of handle cursors with MongoDB."""

    __slots__ = ('_len', ) + Cursor.__slots__

    def __init__(self, *args, **kwargs):

        super(MongoCursor, self).__init__(*args, **kwargs)

        self._len = None

    def __getitem__(self, index):

        result = self._cursor.__getitem__(index)

        if isinstance(result, _Cursor):  # in case of slice
            result = MongoCursor(result)

        return result

    def __iter__(self):

        return self._cursor.__iter__()

    def __len__(self):

        if self._len is None:
            self._len = self._cursor.count(True)

        return self._len
