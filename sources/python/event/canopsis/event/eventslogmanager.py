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
from canopsis.middleware.core import Middleware
from canopsis.timeserie.timewindow import Interval

from canopsis.confng import Configuration, Ini

DEFAULT_EL_STORAGE_URI = 'mongodb-default-eventslog://'


class EventsLog(object):
    """
    Manage events log in Canopsis
    """

    CONF_PATH = 'etc/event/eventlog.conf'
    CONF_CATEGORY = 'EVENTSLOG'

    def __init__(self, config=None, el_storage=None, *args, **kwargs):
        super(EventsLog, self).__init__(*args, **kwargs)

        if config is None:
            self.config = Configuration.load(self.CONF_PATH, Ini)
        else:
            self.config = config

        self.config_el = self.config.get(self.CONF_CATEGORY, {})

        el_storage_uri = self.config_el.get('eventslog_storage_uri', DEFAULT_EL_STORAGE_URI)
        if el_storage is None:
            self.el_storage = Middleware.get_middleware_by_uri(el_storage_uri)


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
