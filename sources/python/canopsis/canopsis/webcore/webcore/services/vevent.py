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

from canopsis.common.ws import route
from canopsis.vevent.manager import VEventManager

vem = VEventManager()

DEFAULT_ROUTE = 'vevent'  #: route specifics to vevents document


def exports(ws):

    @route(ws.application.get, name=DEFAULT_ROUTE)
    def get_by_uids(
        ids, limit=0, skip=0, sort=None, projection=None, with_count=False
    ):
        """Get documents by uids.

        :param list ids: list of document ids.
        :param int limit: max number of elements to get.
        :param int skip: first element index among searched list.
        :param sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order.
        :type sort: list of {(str, {ASC, DESC}}), or str}
        :param dict projection: key names to keep from elements.
        :param bool with_count: If True (False by default), add count to the
            result.
        :return: documents where uids are in uids.
        :rtype: list
        """

        result = vem.get_by_uids(
            uids=ids, limit=limit, skip=skip, sort=sort,
            projection=projection, with_count=with_count
        )

        return result

    @route(
        ws.application.post, name=DEFAULT_ROUTE,
        payload=['source', 'vevents']
    )
    @route(
        ws.application.put, name=DEFAULT_ROUTE,
        payload=['source', 'vevents']
    )
    def put(vevents, source=None):
        """Add vevents (and optionally data) related to input source.

        :param str source: vevent source if not None.
        :param list vevents: vevents (document, str or ical vevent).
        :return: new documents.
        :rtype: list
        """

        result = vem.put(source=source, vevents=vevents)

        return result

    @route(
        ws.application.delete, name=DEFAULT_ROUTE,
        payload=['ids']
    )
    def remove(ids=None):
        """Remove elements from storage where uids are given.

        :param list ids: list of document uids to remove from storage
            (default all empty storage documents).
        """

        result = vem.remove(uids=ids)

        return result

    @route(
        ws.application.delete, name='vevent/source',
        payload=['sources']
    )
    def remove_by_source(sources=None):
        """Remove vevent documents related to input sources.

        :param list sources: sources from where remove related vevent
            documents.
        """

        result = vem.remove_by_source(sources=sources)

        return result

    @route(ws.application.get, name='vevent/values')
    def values(
        sources=None, dtstart=None, dtend=None, query=None,
        limit=0, skip=0, sort=None, projection=None, with_count=False
    ):
        """Get source vevent document values.

        :param list sources: sources from where get values. If None, use all
            sources.
        :param int dtstart: vevent dtstart (default 0).
        :param int dtend: vevent dtend (default sys.maxsize).
        :param dict query: additional filtering query to apply in the search.
        :param int limit: max number of elements to get.
        :param int skip: first element index among searched list.
        :param sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order.
        :type sort: list of {(str, {ASC, DESC}}), or str}
        :param dict projection: key names to keep from elements.
        :param bool with_count: If True (False by default), add count to the
            result.
        :return: matchable documents.
        :rtype: list
        """

        result = vem.values(
            sources=sources, dtstart=dtstart, dtend=dtend, query=query,
            limit=limit, skip=skip, sort=sort, projection=projection,
            with_count=with_count
        )

        return result

    @route(ws.application.get, name='vevent/whois')
    def whois(sources=None, dtstart=None, dtend=None, query=None):
        """Get a set of sources which match with timed condition and query.

        :param list sources: sources from where get values. If None, use all
            sources.
        :param int dtstart: vevent dtstart (default 0).
        :param int dtend: vevent dtend (default sys.maxsize).
        :param dict query: additional filtering query to apply in the search.
        :return: sources.
        :rtype: set
        """

        result = vem.whois(
            sources=sources, dtstart=dtstart, dtend=dtend, query=query
        )

        return result
