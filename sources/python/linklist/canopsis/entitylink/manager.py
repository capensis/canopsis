# -*- coding: utf-8 -*-

from time import time
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
import uuid
from canopsis.context.manager import Context

CONF_PATH = 'linklist/linklist.conf'
CATEGORY = 'LINKLIST'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Entitylink(MiddlewareRegistry):

    ENTITY_STORAGE = 'entitylink_storage'

    """
    Manage entity link information in Canopsis
    """

    def __init__(self, *args, **kwargs):

        super(Entitylink, self).__init__(*args, **kwargs)
        self.context = Context()

    def get_links_from_event(self, event):
        entity = self.context.get_entity(event)
        entity_id = self.context.get_entity_id(entity)
        return self.find(ids=[entity_id])

    def find(
        self,
        limit=None,
        skip=None,
        ids=None,
        sort=None,
        with_count=False,
        _filter={},
    ):

        """
        Retrieve information from data sources

        :param ids: an id list for document to search
        :param limit: maximum record fetched at once
        :param skip: ordinal number where selection should start
        :param with_count: compute selection count when True
        """

        result = self[Entitylink.ENTITY_STORAGE].get_elements(
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
        _id,
        document,
        cache=False
    ):
        """
        Persistance layer for upsert operations

        :param _id: entity id
        :param document: contains link information for entities
        """

        self[Entitylink.ENTITY_STORAGE].put_element(
            _id=_id, element=document, cache=cache
        )

    def remove(
        self,
        ids
    ):
        """
        Remove fields persisted in a default storage.

        :param element_id: identifier for the document to remove
        """

        self[Entitylink.ENTITY_STORAGE].remove_elements(ids=ids)
