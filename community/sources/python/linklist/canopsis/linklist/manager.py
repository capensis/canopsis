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

from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
import uuid

CONF_PATH = 'linklist/linklist.conf'
CATEGORY = 'LINKLIST'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Linklist(MiddlewareRegistry):
    """Manage linklist information in Canopsis.
    """

    LINKLIST_STORAGE = 'linklist_storage'  #: linklist storage name

    def find(
        self,
        limit=None,
        skip=None,
        ids=None,
        sort=None,
        with_count=False,
        _filter={},
    ):
        """Retrieve information from data sources.

        :param ids: id(s) for document to search. All documents by default.
        :type ids: str or list
        :param int limit: maximum record fetched at once.
        :param int skip: ordinal number where selection should start.
        :param bool with_count: compute selection count when True.
        :return: depending on ids type:

            - str: one document or None if no related document exists.
            - list or None: storage Cursor.
        :rtype: Cursor or dict
        """

        result = self[Linklist.LINKLIST_STORAGE].get_elements(
            ids=ids,
            skip=skip,
            sort=sort,
            limit=limit,
            query=_filter,
            with_count=with_count
        )

        return result

    def put(
        self,
        document,
        cache=False
    ):
        """Persistance layer for upsert operations.

        :param dict document: document to put.
        :param bool cache: if True (default false), use storage cache.
        :return: put result.
        :rtype: dict
        """
        if not document.get('id'):
            document['_id'] = str(uuid.uuid4())
        else:
            document['_id'] = document.pop('id')

        return self[Linklist.LINKLIST_STORAGE].put_element(
            _id=document['_id'], element=document, cache=cache
        )

    def remove(
        self,
        ids
    ):
        """Remove fields persisted in a default storage.

        :param element_id: identifier for the document to remove
        """

        self[Linklist.LINKLIST_STORAGE].remove_elements(ids=ids)
