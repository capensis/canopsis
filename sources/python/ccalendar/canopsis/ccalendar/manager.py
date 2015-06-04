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
    conf_paths, add_category
)
from canopsis.middleware.registry import MiddlewareRegistry

CONF_PATH = 'calendar/calendar.conf'
CATEGORY = 'CALENDAR'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Calendar(MiddlewareRegistry):
    """Manage calendar information in Canopsis.
    """

    CALENDAR_STORAGE = 'calendar_storage'

    def find(
        self,
        limit=None,
        skip=None,
        ids=None,
        sort=None,
        with_count=False,
        query={}
    ):
        """Retrieve information from data sources

        :param str ids: an id list for document to search.
        :param int limit: maximum record fetched at once.
        :param int skip: ordinal number where selection should start.
        :param bool with_count: compute selection count when True.
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
        """Persistance layer for upsert operations

        :param _id: entity id
        :param document: contains link information for entities
        """

        self[Calendar.CALENDAR_STORAGE].put_element(
            _id=_id, element=document, cache=cache
        )

    def remove(
        self,
        ids,
        cache=False
    ):
        """Remove fields persisted in a default storage.

        :param element_id: identifier for the document to remove
        """

        self[Calendar.CALENDAR_STORAGE].remove_elements(ids=ids, cache=cache)
