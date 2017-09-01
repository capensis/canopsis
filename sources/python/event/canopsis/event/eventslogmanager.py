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

from canopsis.middleware.core import Middleware
from canopsis.timeserie.timewindow import Interval


class EventsLog(object):
    """
    Manage events log in Canopsis
    """

    DEFAULT_EL_STORAGE_URI = 'mongodb-default-eventslog://'

    @classmethod
    def provide_default_basics(cls):
        """
        Provide default storage for EventsLog
        """
        return Middleware.get_middleware_by_uri(
            cls.DEFAULT_EL_STORAGE_URI, table='events_log')

    def __init__(self, el_storage, *args, **kwargs):
        super(EventsLog, self).__init__(*args, **kwargs)
        self.el_storage = el_storage

    def get_eventlog_count_by_period(
            self, tstart, tstop, limit=100, query={}.copy()):
        """Get an eventlog count for each interval found in the given period and
           with a given filter.
           This period is given by tstart and tstop.
           Counts can be limited thanks to the 'limit' parameter.

           :param start: begin interval timestamp
           :param stop: end interval timestamp
           :param limit: number that limits the max count returned
           :param query: filter for events_log collection
           :return: list in which each item contains an interval and the
           related count
           :rtype: list
        """
        period = {'day': 1}
        intervals = Interval.get_intervals_by_period(tstart, tstop, period)
        results = []

        for date in intervals:
            eventfilter = {
                '$and': [
                    query,
                    {
                        'timestamp': {
                            '$gte': date['begin'],
                            '$lte': date['end']
                        }
                    }
                ]
            }

            _, count = self.el_storage.find_elements(
                query=eventfilter,
                limit=limit,
                with_count=True
            )

            results.append({
                'date': date,
                'count': count
            })

        return results
