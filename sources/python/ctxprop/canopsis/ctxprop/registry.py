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

"""This module is dedicated to CTXPropRegistry interface.

It aims to execute ctxprop manager read and delete methods related to its own
property managers such as the perfdata manager for example.
"""

from canopsis.common.init import basestring
from canopsis.middleware.core import Middleware


class CTXPropRegistry(Middleware):
    """dedicated to ctx events.
    """

    class Error(Exception):
        """Handle CTXPropRegistry errors.
        """

    CTX_ID = 'ctx_id'  #: default ctx id field name

    __register__ = True  #: register this class such as a middleware
    #: protocol registration name if __register__
    __protocol__ = 'ctxpropreg'

    def get(self, ids=None, query=None, *args, **kwargs):
        """Get ctx property related to input ``ctx id(s)`` and
        additional specific filter ``query`` per ctx id.

        :param str(s) ids: ctx id(s) from where find property. If None
            (default), get all property of all ctxs.
        :param dict query: specific query to apply on the method execution.
        :return: couples of (ctx id, [property]).
        :rtype: dict
        """

        result = self._process(
            cmd=self._get, ids=ids, query=query, *args, **kwargs
        )

        return result

    def _get(self, *args, **kwargs):
        """Protected get method implementation.
        """

        raise NotImplementedError()

    def count(self, ids=None, query=None, *args, **kwargs):
        """Get ctx property count related to input ``ctx id(s)`` and
        additional specific filter ``query`` per ctx id.

        :param ids: ctx id(s) from where find property. If None
            (default), get all property of all ctxs.
        :type ids: str or list
        :param dict query: specific query to apply on the method execution.
        :return: couples of (ctx id, count).
        :rtype: dict
        """

        result = self._process(
            cmd=self._count, ids=ids, query=query,
            *args, **kwargs
        )

        return result

    def _count(self, *args, **kwargs):
        """Protected count method implementation.

        By default execute the _get method and generate documents containing
        both self.ctx_id_field() and count keys.
        """

        documents = self._get(*args, **kwargs)

        if isinstance(documents, dict):
            documents = [documents]

        counts = {}

        ctx_id_field = self._ctx_id_field()

        for document in documents:
            ctx_id = document[ctx_id_field]
            if ctx_id in counts:
                counts[ctx_id] += 1
            else:
                counts[ctx_id] = 1

        result = [
            {ctx_id: count, 'count': counts[count]} for count in counts
        ]

        return result

    def delete(
            self, ids=None, query=None, cache=False, force=False,
            *args, **kwargs
    ):
        """Get ctx property count related to input ``ctx id(s)`` and
        additional specific filter ``query`` per ctx id.

        :param str(s) ids: ctx id(s) from where find property. If None
            (default), get all property of all ctxs.
        :param dict query: specific query to apply on the method execution.
        :param bool cache: storage cache property.
        :return: couples of (ctx id, count).
        :rtype: dict
        """

        # check if force is True if registries is None
        if ids is None and not force:
            raise CTXPropRegistry.Error(
                "Impossible to remove all existing ctx properties. Use force."
            )

        result = self._process(
            cmd=self._delete, ids=ids, query=query, cache=cache,
            *args, **kwargs
        )

        return result

    def _delete(self, *args, **kwargs):
        """Protected delete method implementation.
        """

        raise NotImplementedError()

    def ids(self, query=None):
        """Get ctx id(s) related to input query.

        :param dict query: query to apply on ctxpropregistry documents in order
            to retrieve related ctx id. If None, get all self ctx ids.
        :return: list of ctx ids.
        :rtype: list
        """

        raise NotImplementedError()

    def _process(self, cmd, ids, *args, **kwargs):

        ids, unique = CTXPropRegistry._ctx_ids_unique(ids)

        result = cmd(
            ids=ids, *args, **kwargs
        )

        result = self._final_result(result, unique)

        return result

    def _ctx_id_field(self):
        """Get ctx id document field name.

        :return: ctx id document field name. ``CTXPropRegistry.CTX_ID`` by
            default.
        :rtype: str
        """

        return CTXPropRegistry.CTX_ID

    @staticmethod
    def _ctx_ids_unique(ids):
        """Return a couple of (list of ctx ids, ids is unique).

        :param str(s) ids: ctx id(s).
        :return: a couple of (list of ctx ids, ids is unique) where
            unique is True iif ids is a string and, the list of ctx
            ids depending on ids type:

        - None: empty list.
        - str: [ids].
        - iterable: ids.
        :rtype: tuple
        """

        unique = isinstance(ids, basestring)

        if unique:
            ids = [ids]

        result = ids, unique

        return result

    def _final_result(self, queryresult, unique):
        """Transform a query result to a final result which is a dictionary of
        document by ctx id.

        :param dict(s) queryresult: query result to transform into a dict of
            documents by ctx id.
        :param bool unique: True if related ids where a string.
        :return: dict of documents by ctx id or documents if unique.
        :rtype: dict or list
        """

        # get default ctx id field name
        ctx_id_field = self._ctx_id_field()

        # initialiaze the result with an empty dict
        result = {}

        # ensure queryresult such as a list of documents
        if isinstance(queryresult, dict):
            queryresult = [queryresult]
        elif queryresult is None:
            queryresult = []

        # iterate on query result
        for r in queryresult:
            # get document ctx id
            ctx_id = r[ctx_id_field]
            # get result documents
            documents = result.setdefault(ctx_id, [])
            # append a document in result documents
            documents.append(r)

        if unique:  # apply unique rule
            result = iter(result).next() if result else None

        return result
