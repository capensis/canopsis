# -*- coding: utf-8 -*-

from time import time
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
import uuid

CONF_PATH = 'calendar/calendar.conf'
CATEGORY = 'CALENDAR'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Calendar(MiddlewareRegistry):

    CALENDAR_STORAGE = 'calendar_storage'

    """
    Manage calendar information in Canopsis
    """

    def __init__(self, *args, **kwargs):

        super(Calendar, self).__init__(*args, **kwargs)

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

        result = self[Calendar.CALENDAR_STORAGE].get_elements(
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
        _id,
        document,
        cache=False
    ):
        """
        Persistance layer for upsert operations

        :param _id: entity id
        :param document: contains link information for entities
        """

        self[Calendar.CALENDAR_STORAGE].put_element(
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

        self[Calendar.CALENDAR_STORAGE].remove_elements(ids=ids)
