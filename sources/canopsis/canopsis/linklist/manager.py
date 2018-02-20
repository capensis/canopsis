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

import uuid

from canopsis.logger import Logger
from canopsis.middleware.core import Middleware


class Linklist(object):
    """
    Manage linklist information in Canopsis.
    """

    LOG_PATH = 'var/log/linklist.log'
    LINKLIST_STORAGE_URI = 'mongodb-default-linklist://'

    def __init__(self, logger, storage):
        self.logger = logger
        self.linklist_storage = storage

    @classmethod
    def provide_default_basics(cls):
        """
        Provide logger, config, storages...

        ! Do not use in tests !

        :rtype: Union[logging.Logger,
                      canopsis.storage.core.Storage]
        """
        logger = Logger.get('linklist', cls.LOG_PATH)
        linklist_storage = Middleware.get_middleware_by_uri(
            cls.LINKLIST_STORAGE_URI
        )

        return (logger, linklist_storage)

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

        result = self.linklist_storage.get_elements(
            ids=ids,
            skip=skip,
            sort=sort,
            limit=limit,
            query=_filter,
            with_count=with_count
        )

        return result

    def put(self, document):
        """Persistance layer for upsert operations.

        :param dict document: document to put
        :return: put result
        :rtype: dict
        """
        if not document.get('id'):
            document['_id'] = str(uuid.uuid4())
        else:
            document['_id'] = document.pop('id')

        return self.linklist_storage.put_element(
            _id=document['_id'], element=document
        )

    def remove(self, ids):
        """Remove fields persisted in a default storage.

        :param ids: identifier for documents to remove
        :type: ids: list
        """

        self.linklist_storage.remove_elements(ids=ids)
