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

"""
"""

from canopsis.common.init import basestring
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.middleware.registry import MiddlewareRegistry

from time import time

from icalendar import Event

from dateutil.rrule import rrulestr

from calendar import timegm

from datetime import datetime, time as datetime_time

from uuid import uuid4 as uuid

from sys import maxsize

CONF_PATH = 'vevent/vevent.conf'
CATEGORY = 'VEVENT'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class VEventManager(MiddlewareRegistry):
    """Manage virtual event data.

    Such vevent are technically an expression which respects the icalendar
    specification ftp://ftp.rfc-editor.org/in-notes/rfc2445.txt.

    A vevent document contains several values. Each value contains
    an icalendar expression (dtstart, rrule, duration) and an array of
    behavior entries:

    {
        id: document_id,
        source: source element id,
        dtstart: datetime start,
        dtend: datetime end,
        duration: vevent duration,
        freq: vevent freq,
        vevent: vevent ical format value,
        info: data information
    }.
    """

    STORAGE = 'vevent_storage'  #: vevent storage name

    UID = 'uid'  #: document id
    SOURCE = 'source'  #: source field name
    DTSTART = 'dtstart'  #: dtstart field name
    DTEND = 'dtend'  #: dtend field name
    DURATION = 'duration'  #: duration field name
    FREQ = 'freq'  #: freq field name
    VEVENT = 'vevent'  #: vevent value field name

    def __init__(self, vevent_storage=None, *args, **kwargs):
        """
        :param Storage vevent_storage: vevent storage.
        """

        super(VEventManager, self).__init__(*args, **kwargs)
        # set storage if given
        if vevent_storage is not None:
            self[VEventManager.STORAGE] = vevent_storage

    def _get_info(self, vevent):
        """Get information from an ical Event.

        :param Event vevent: vevent from where get information.
        :return: vevent information in a dictionary.
        :rtype: dict
        """

        return None

    def get_by_uids(
        self, uids,
        limit=0, skip=0, sort=None, projection=None, with_count=False
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

        documents = self[VEventManager.STORAGE].get_element(
            ids=uids,
            limit=limit, skip=skip, sort=sort, projection=projection,
            with_count=with_count
        )

        if with_count:
            result = list(documents[0]), documents[1]
        else:
            result = list(documents[0])

        return result

    def values(
        self, sources=None, dtstart=None, dtend=None, query=None,
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

        # initialize query
        if query is None:
            query = {}

        # put sources in query if necessary
        if sources is not None:
            query[VEventManager.SOURCE] = {'$in': sources}
        # put dtstart and dtend in query
        if dtstart is None:
            dtstart = 0
        if dtend is None:
            dtend = maxsize

        query['$and'] = [
            {
                '$or': [
                    {VEventManager.DTSTART: {'$geq': dtstart}},
                    {VEventManager.DTSTART: {'$leq': dtend}}
                ]
            },
            {
                '$or': [
                    {VEventManager.DTEND: {'$geq': dtstart}},
                    {VEventManager.DTEND: {'$leq': dtend}}
                ]
            }
        ]

        documents = self[VEventManager.STORAGE].find_elements(
            query=query,
            limit=limit, skip=skip, sort=sort, projection=projection,
            with_count=with_count
        )

        if with_count:
            result = list(documents[0]), documents[1]
        else:
            result = list(documents)

        return result

    def whois(self, sources=None, dtstart=None, dtend=None, query=None):
        """Get a set of sources which match with timed condition and query.

        :param list sources: sources from where get values. If None, use all
            sources.
        :param int dtstart: vevent dtstart (default 0).
        :param int dtend: vevent dtend (default sys.maxsize).
        :param dict query: additional filtering query to apply in the search.
        :return: sources.
        :rtype: set
        """

        values = self.values(
            sources=sources, dtstart=dtstart, dtend=dtend, query=query
        )

        result = set([value[VEventManager.SOURCE] for value in values])

        return result

    def put(self, vevents, source=None, cache=False):
        """Add vevents (and optionally data) related to input source.

        :param str source: vevent source if not None.
        :param list vevents: vevents (document, str or ical vevent).
        :param dict info: vevent info.
        :param bool cache: if True (default False), use storage cache.
        :return: new documents.
        :rtype: list
        """

        result = []

        for vevent in vevents:

            document = vevent if isinstance(vevent, dict) else None

            # if document has to be generated ...
            if document is None:
                # ensure vevent is an ical format
                if isinstance(vevent, basestring):
                    vevent = Event.from_ical(vevent)
                # get dtstart
                dtstart = vevent.get(VEventManager.DTSTART, 0)
                if isinstance(dtstart, datetime):
                    dtstart = timegm(dtstart.timetuple())
                # get dtend
                dtend = vevent.get(VEventManager.DTEND, 0)
                if isinstance(dtend, datetime):
                    dtend = timegm(dtend.timetuple())
                # get duration
                duration = vevent.get([VEventManager.DURATION])
                # prepare the document
                document = {
                    VEventManager.SOURCE: source,
                    VEventManager.DTSTART: dtstart,
                    VEventManager.DTEND: dtend,
                    VEventManager.DURATION: duration,
                    VEventManager.VEVENT: vevent.to_ical()
                }
                # get info
                document_info = self._get_info(vevent)
                if document_info is not None:
                    document.update(document_info)

            # get document uid
            if VEventManager.UID in document:
                uid = document[VEventManager.UID]
            else:
                uid = str(uuid())
                # put it in document if not already present
                document[VEventManager.UID] = uid

            result.append(document)

            self[VEventManager.STORAGE].put_element(
                _id=uid, document=document
            )

        return result

    def remove(self, uids=None, cache=False):
        """Remove elements from storage where uids are given.

        :param list uids: list of document uids to remove from storage
            (default all empty storage documents).
        """

        result = self[VEventManager.STORAGE].remove_elements(
            ids=uids, cache=cache
        )

        return result

    def remove_by_source(self, sources=None, cache=False):
        """Remove vevent documents related to input sources.

        :param list sources: sources from where remove related vevent
            documents.
        """
        _filter = {}

        if sources is not None:
            _filter[VEventManager.SOURCE] = {'$in': sources}

        result = self[VEventManager.STORAGE].remove_elements(
            _filter=_filter, cache=cache
        )

        return result

    def get_before(
        self, sources=None, ts=None, dtstart=None, dtstop=None, query=None
    ):
        """Get before date related to one timestamp and additional parameters.

        :param list sources: sources from where parse vevent documents.
        :param int ts: timestamp from when find vevent documents.
        :param int dtstart: vevent dtstart.
        :param int dtend: vevent dtend.
        :param dict query: additional filtering query to apply in the search.
        :return: list of before date per source.
        """

        result = self._get_around(
            sources=sources, ts=ts, dtstart=dtstart, dtstop=dtstop,
            query=query, after=False
        )

        return result

    def get_after(
        self, sources=None, ts=None, dtstart=None, dtstop=None, query=None
    ):
        """Get afer date related to one timestamp and additional parameters.

        :param list sources: sources from where parse vevent documents.
        :param int ts: timestamp from when find vevent documents.
        :param int dtstart: vevent dtstart.
        :param int dtend: vevent dtend.
        :param dict query: additional filtering query to apply in the search.
        :return: list of after date per source.
        """

        result = self._get_around(
            sources=sources, ts=ts, dtstart=dtstart, dtstop=dtstop,
            query=query, after=True
        )

        return result

    def _get_around(
        self,
        sources=None, ts=None, dtstart=None, dtend=None, query=None,
        after=True
    ):
        """Get around date related to one timestamp and additional parameters.

        :param list sources: sources from where parse vevent documents.
        :param int ts: timestamp from when find vevent documents.
        :param int dtstart: vevent dtstart.
        :param int dtend: vevent dtend.
        :param dict query: additional filtering query to apply in the search.
        :param bool after: if True (default), get period after ts (included),
            otherwise, before ts.
        :return: list of around date per source.
        """

        result = {}

        # initialize ts
        if ts is None:
            ts = time()
        # calculate ts datetime
        dtts = datetime.fromtimestamp(ts)

        # get entity documents(s)
        documents = self.values(
            sources=sources, dtstart=dtstart, dtend=dtend, query=query
        )

        for document in documents:
            # get event
            vevent = document[VEventManager.VEVENT]
            event = Event.from_ical(vevent)

            # get duration, dtstart and rrule
            duration = event.get('duration')
            duration = duration.dt
            dtstart = event.get('dtstart')
            dtstart = dtstart.dt

            if isinstance(dtstart, datetime_time):
                dtstart = datetime.now().replace(
                    hour=dtstart.hour, minute=dtstart.minute,
                    second=dtstart.second, tzinfo=dtstart.tzinfo
                )

            rrule = event.get('rrule')
            rrule = rrulestr(rrule.to_ical(), cache=True, dtstart=dtstart)
            # calculate first date after dtts including dtts
            if after:
                around = rrule.after(dt=dtts, inc=True)
            else:
                around = rrule.before(dt=dtts, inc=True)

            # if around datetime exist
            if around is not None:
                source = document[VEventManager.SOURCE]

                if source not in result:
                    result[source] = {}

                # add duration
                end = around + duration

                # and check if dtstart is in [first; end]
                if around <= dtts <= end:
                    # update end in the result
                    endts = timegm(end.timetuple())

                    result[source] = endts

        return result
