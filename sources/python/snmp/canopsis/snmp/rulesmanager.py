# -*- coding: utf-8 -*-

from time import time
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
import uuid

CONF_PATH = 'snmp/snmprules.conf'
CATEGORY = 'SNMPRULES'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class RulesManager(MiddlewareRegistry):

    STORAGE = 'snmprules_storage'

    """
    Manage snmp rules information in Canopsis
    """

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

        result = self[RulesManager.STORAGE].get_elements(
            ids=ids,
            skip=skip,
            sort=sort,
            limit=limit,
            query=query,
            with_count=with_count
        )
        return result

    def put(
        self,
        document,
        cache=False
    ):
        """
        Persistance layer for upsert operations

        :param document: contains link information for entities
        """
        if 'id' not in document or not document['id']:
            document['_id'] = str(uuid.uuid4())
        else:
            document['_id'] = document['id']
            del document['id']

        self[RulesManager.STORAGE].put_element(
            _id=document['_id'], element=document, cache=cache
        )

    def remove(
        self,
        ids
    ):
        """
        Remove fields persisted in a default storage.

        :param element_id: identifier for the document to remove
        """

        self[RulesManager.STORAGE].remove_elements(ids=ids)
