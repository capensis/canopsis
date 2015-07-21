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

__version__ = '0.1'

from sqlalchemy import create_engine, Table, Column, Integer, String, MetaData
from sqlalchemy.sql import select

from canopsis.common.init import basestring
from canopsis.storage.core import Storage, DataBase
from canopsis.common.utils import isiterable


class SQLAlchemyDataBase(DataBase):
    """
    Manage access to a mongodb.
    """

    def _connect(self, *args, **kwargs):

        result = None

        try:
            engine = create_engine(
                self.db,
                url=self.url,
                logging_name=self.log_name,
                echo=True
            )
        except Exception as e:
            self.logger.warning(e)
        else:
            connection_args = {}

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
                result = engine.connect(**connection_args)
            except Exception as e:
                self.logger.error(
                    'Raised {2} during connection attempting to {0}:{1}.'.
                    format(self.host, self.port, e))

        return result

    def _disconnect(self, *args, **kwargs):

        if self._conn is not None:
            self._conn.close()
            self._conn = None

    def connected(self, *args, **kwargs):

        result = self._conn is not None and not self._conn.closed

        return result

    def size(self, table=None, criteria=None, *args, **kwargs):

        raise NotImplementedError()

    def drop(self, table=None):

        self.conn.execute(self._backend.delete())


class SQLAlchemyStorage(SQLAlchemyDataBase, Storage):

    __register__ = True  #: register this class to middleware
    __protocol__ = 'sqlalchemy'  #: register this class to the protocol sqlalchemy

    def _connect(self, *args, **kwargs):

        result = super(SQLAlchemyStorage, self)._connect(*args, **kwargs)

        if result:
            metadata = MetaData()
            table = self.get_table()

            self._backend = Table(
                table,
                metadata
            )

            for index in self.all_indexes():
                try:
                    self._backend.ensure_index(index)
                except Exception as e:
                    self.logger.error(e)

            metadata.create_all(self._conn)

        return result

    def get_elements(
        self,
        ids=None, query=None, limit=0, skip=0, sort=None, with_count=False,
        *args, **kwargs
    ):

        sel = select([self._get_backend()], limit=limit, skip=skip, sort=sort)

        if ids is not None:
            sel = sel.where(ids)

        if query is not None:
            sel = sel.where(query)

        result = self._conn.execute(sel)

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

        del_ = self._get_backend().delete()

        if ids is not None:
            del_.where(ids=ids)

        if _filter is not None:
            del_.where(_filter)

    def put_element(self, _id, element, *args, **kwargs):

        pass

    def bool_compare_and_swap(self, _id, oldvalue, newvalue):

        return self.val_compare_and_swap(_id, oldvalue, newvalue) == newvalue

    def val_compare_and_swap(self, _id, oldvalue, newvalue):

        try:
            result = self._run_command(
                'find_and_modify',
                query={SQLAlchemyStorage.ID: _id, 'value': oldvalue},
                update={'$set': {'value': newvalue}},
                upsert=True)

        except DuplicateKeyError:
            result = None

        return oldvalue if not result else newvalue

    def all_indexes(self, *args, **kwargs):

        result = super(SQLAlchemyStorage, self).all_indexes(*args, **kwargs)

        result.append([(SQLAlchemyStorage.ID, 1)])

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
