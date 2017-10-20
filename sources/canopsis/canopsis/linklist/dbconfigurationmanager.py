# -*- coding: utf-8 -*-

from canopsis.linklist.manager import Linklist
from canopsis.logger import Logger
from canopsis.middleware.core import Middleware


class DBConfiguration(object):

    """Manage Canopsis database configuration information."""

    DBCONFIGURATION_STORAGE_URI = 'mongodb-default-dbconfiguration://'

    def __init__(self):
        self.logger = Logger.get('linklist', Linklist.LOG_PATH)
        self.dbconfiguration_storage = Middleware.get_middleware_by_uri(
            Linklist.LINKLIST_STORAGE_URI
        )

    def get(self, _id, default=None):
        return self.find_one(query={'_id': _id}) or default

    def find_one(
        self,
        skip=None,
        ids=None,
        sort=None,
        with_count=False,
        query={},
    ):

        results = list(self.find(
            limit=1,
            skip=skip,
            ids=ids,
            sort=sort,
            with_count=with_count,
            query=query,
        ))

        if len(results):
            return results[0]
        else:
            return None

    def find(
        self,
        limit=None,
        skip=None,
        ids=None,
        sort=None,
        with_count=False,
        query={},
    ):
        """
        Retrieve information from data sources

        :param ids: an id list for document to search
        :param limit: maximum record fetched at once
        :param skip: ordinal number where selection should start
        :param with_count: compute selection count when True
        """

        result = self.dbconfiguration_storage.get_elements(
            ids=ids,
            skip=skip,
            sort=sort,
            limit=limit,
            query=query,
            with_count=with_count
        )
        return result

    # Not tested with old object collection system: TODO
    def put(self, _id, document):
        """
        Persistance layer for upsert operations

        :param _id: entity id
        :param document: contains link information for entities
        """

        self.dbconfiguration_storage.put_element(
            _id=_id, element=document
        )

    # Not tested with old object collection system: TODO
    def remove(self, ids):
        """
        Remove fields persisted in a default storage.

        :param element_id: identifier for the document to remove
        """

        self.dbconfiguration_storage.remove_elements(ids=ids)
