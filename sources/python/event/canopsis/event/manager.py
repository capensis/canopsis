# -*- coding: utf-8 -*-

from time import time
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
import uuid

CONF_PATH = 'event/event.conf'
CATEGORY = 'EVENT'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Event(MiddlewareRegistry):

    EVENT_STORAGE = 'event_storage'
    """
    Manage events in Canopsis
    """

    def __init__(self, *args, **kwargs):

        super(Event, self).__init__(*args, **kwargs)

    def find(
        self,
        limit=None,
        skip=None,
        ids=None,
        sort=None,
        with_count=False,
        query={},
        projection=None
    ):

        """
        Retrieve information from data sources

        :param ids: an id list for document to search
        :param limit: maximum record fetched at once
        :param skip: ordinal number where selection should start
        :param with_count: compute selection count when True
        """

        result = self[Event.EVENT_STORAGE].get_elements(
            ids=ids,
            skip=skip,
            sort=sort,
            limit=limit,
            query=query,
            with_count=with_count,
            projection=projection
        )
        return result
