# -*- coding: utf-8 -*-

from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry

CONF_PATH = 'configuration/dbconfiguration.conf'
CATEGORY = 'DBCONFIGURATION'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class DBConfiguration(MiddlewareRegistry):

    DBCONFIGURATION_STORAGE = 'dbconfiguration_storage'

    """Manage Canopsis database configuration information."""

    def __init__(self, dbconfiguration_storage=None, *args, **kwargs):

        super(DBConfiguration, self).__init__(*args, **kwargs)

        if dbconfiguration_storage is not None:
            self[DBConfiguration.DBCONFIGURATION_STORAGE] = dbconfiguration_storage

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

        result = self[DBConfiguration.DBCONFIGURATION_STORAGE].get_elements(
            ids=ids,
            skip=skip,
            sort=sort,
            limit=limit,
            query=query,
            with_count=with_count
        )
        return result

    # Not tested with old object collection system: TODO
    def put(
        self,
        _id,
        document
    ):
        """
        Persistance layer for upsert operations

        :param _id: entity id
        :param document: contains link information for entities
        """

        self[DBConfiguration.DBCONFIGURATION_STORAGE].put_element(
            _id=_id, element=document
        )

    # Not tested with old object collection system: TODO
    def remove(
        self,
        ids
    ):
        """
        Remove fields persisted in a default storage.

        :param element_id: identifier for the document to remove
        """

        self[DBConfiguration.DBCONFIGURATION_STORAGE].remove_elements(ids=ids)
