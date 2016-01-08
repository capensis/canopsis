# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 'Capensis' [http://www.capensis.com]
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
from canopsis.ccalendar.manager import CalendarManager
from json import dumps
cm = CalendarManager()


def exports(ws):

    @route(ws.application.get, name='calendar')
    def get_by_uids(
        uids, limit=0, skip=0, sort=None, projection=None, with_count=False
    ):
        """Get documents by uids.

        :param list uids: list of document uids.
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

        result = cm.get_by_uids(
            uids=uids, limit=limit, skip=skip, sort=sort,
            projection=projection, with_count=with_count
        )

        return result

    @route(
        ws.application.post, name='calendar',
        payload=['eventcategories', 'output', 'dtstart', 'dtend']
    )
    @route(
        ws.application.put, name='calendar',
        payload=['eventcategories', 'output', 'dtstart', 'dtend', 'uid']
    )
    def put(
        eventcategories,
        output,
        dtstart,
        dtend,
        uid=None,
        source=None,
        info=None,
        **kwargs
    ):
        """Add calendar events (and optionally data) related to input source.

        :param str source: calendardata source if not None.
        :param dict info: calendar event info.
        :param str eventcategories: eventcategories of the event.
        :param str output: description of the event.
        :param int dtstart: beginning date.
        :param int dtend: ending date.
        :return: new documents.
        :rtype: list
        """
        ws.logger.debug('eventcategories, {}'.format(
            dumps(eventcategories, indent=2)
        ))
        calendarDocument = cm.get_document(
            uid=uid,
            eventcategories=eventcategories,
            output=output,
            dtstart=dtstart,
            dtend=dtend
        )
        ws.logger.debug('calendarDocument, {}'.format(
            dumps(calendarDocument, indent=2)
        ))
        calendar_vevents = []
        calendar_vevents.append(calendarDocument)

        result = cm.put(source=source, vevents=calendar_vevents)

        return result

    @route(ws.application.delete, payload=['ids'])
    def calendar(ids):
        """Delete a event with the id. We can delete several events with an array of ids

        :param list ids: list of ids to delete
        :return: True if the remove is ok
        :rtype: boolean
        """
        result = cm.remove(uids=ids)
        ws.logger.info('Delete : {}'.format(ids))
        return True

    @route(
        ws.application.get, name='calendar/values',
        payload=['query', 'source', 'dtstart', 'dtend']
    )
    def values(
        query=None, sources=None, dtstart=None, dtend=None
    ):
        """Get source vevent document values.

        :param list sources: sources from where get values. If None, use all
            sources.
        :param int dtstart: vevent dtstart (default 0).
        :param int dtend: vevent dtend (default sys.maxsize).
        :param dict query: vevent information if given.
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
        ws.logger.info("dtend, {}".format(dtend))
        result = cm.values(
            sources=sources, dtstart=dtstart, dtend=dtend, query=query
        )

        return result
