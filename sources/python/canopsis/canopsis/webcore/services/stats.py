# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

from canopsis.common.ws import route
from canopsis.stats.manager import Stats


def exports(ws):

    sm = Stats()

    @route(
        ws.application.get,
        name='stats/event',
        payload=['tstart', 'tstop', 'tags']
    )
    def event_stats(tstart, tstop, tags={}):
        """
        Get event related stats

        :param int tstart: Start timestamp
        :param int tstop: End timestamp
        :param dict tags: Group stats by values

        :returns: Stats or an empty dict is tags={}
        :rtype: dict

        An example is available in canopsis.stats.manager.Stats.get_event_stats
        """
        return sm.get_event_stats(tstart, tstop, tags=tags)

    @route(
        ws.application.get,
        name='stats/user',
        payload=['tstart', 'tstop', 'users', 'tags']
    )
    def user_stats(tstart, tstop, users=[], tags={}):
        """
        Get user related stats

        :param int tstart: Start timestamp
        :param int tstop: End timestamp
        :param list users: Users whose stats want to be retrieved
        :param dict tags: Groups stats by value

        :returns: Stats or an empty list if users=[]
        :rtype: list

        An example is available in canopsis.stats.manager.Stats.get_user_stats
        """
        return sm.get_user_stats(tstart, tstop, users=users, tags=tags)
