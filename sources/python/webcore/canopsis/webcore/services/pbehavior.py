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

from canopsis.common.utils import ensure_iterable
from canopsis.common.ws import route
from canopsis.pbehavior.manager import PBehaviorManager


DEFAULT_ROUTE = 'pbehavior'  #: route specifics to pbehavior document


def exports(ws):
    pbm = PBehaviorManager()

    @route(
        ws.application.post,
        payload=['entity_ids', 'behaviors', 'start', 'end'],
        name=DEFAULT_ROUTE
    )
    def find(entity_ids=None, behaviors=None, start=None, end=None):
        """Find documents related to input entity id(s) and behavior(s).

        :param entity_ids:
        :type entity_ids: list or str
        :param behaviors:
        :type behaviors: list or str
        :param int start: start timestamp.
        :param int end: end timestamp.
        :return: entity documents with input behaviors.
        :rtype: list
        """

        query = PBehaviorManager.get_query(behaviors)

        entity_ids = ensure_iterable(entity_ids)

        result = pbm.values(
            sources=entity_ids, query=query, dtstart=start, dtend=end
        )

        return result

    @route(
        ws.application.put, name=DEFAULT_ROUTE,
        payload=['document']
    )
    def put(document):
        """Put a pbehavior document.

        :param str _id: document entity id.
        :param dict document: pbehavior document.
        :param bool cache: if True (False by default), use storage cache.
        :return: _id if input pbehavior document has been putted. Otherwise
            None.
        :rtype: str
        """

        result = pbm.put(vevents=[document])

        return result

    @route(
        ws.application.delete, payload=['ids'],
        name=DEFAULT_ROUTE
    )
    def remove(ids=None):
        """Remove document(s) by id.

        :param ids: pbehavior document id(s). If None, remove all documents.
        :type ids: list or str
        :param bool cache: if True (False by default), use storage cache.
        :return: removed document id(s).
        :rtype: list
        """

        result = pbm.remove(uids=ids)

        return result

    @route(
        ws.application.get,
        payload=['behaviors', 'start', 'end'],
        name='pbehavior/calendar'
    )
    def find_pbehavior(behaviors=None, start=None, end=None):
        """Get pbehavior which are between a starting and a ending date. They are filtered by behavior

        :param string behaviors: behavior to filter the query (optionnal)
        :param timestamp start: begin of the filtered period
        :param timestamp end: end of the filtered period
        :return: matchable events
        :rtype: list
        """

        storage = pbm[PBehaviorManager.STORAGE]

        if behaviors is None:
            pbehavior_list = storage.find_elements(query={
                '$or': [
                    {
                        'dtstart': {'$gte': start},
                        'dtstart': {'$lte': end}
                    },
                    {
                        'dtend': {'$gte': start},
                        'dtend': {'$lte': end}
                    },
                    {
                        'dtstart': {'$lte': start},
                        'dtend': {'gte': end}
                    }
                ]
            })
        else:
            pbehavior_list = storage.find_elements(query={
                '$and': [
                    {'behaviors': behaviors},
                    {
                        '$or': [
                            {
                                'dtstart': {'$gte': start},
                                'dtstart': {'$lte': end}
                            },
                            {
                                'dtend': {'$gte': start},
                                'dtend': {'$lte': end}
                            },
                            {
                                'dtstart': {'$lte': start},
                                'dtend': {'gte': end}
                            }
                        ]
                    }
                ]
            })

        return list(pbehavior_list)
