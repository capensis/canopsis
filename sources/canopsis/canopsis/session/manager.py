#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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
Session manager definition.
"""

from __future__ import unicode_literals

from time import time

from canopsis.confng import Configuration, Ini

DEFAULT_ALIVE_SESSION_DURATION = 300


class Session(object):
    """
    Manage session informations.
    """

    CONF_PATH = 'etc/session/session.conf'

    SESSION_STORAGE_URI = 'mongodb-default-session://'

    SESSION_COLLECTION = 'session_collection'

    def __init__(
            self,
            collection,
    ):
        """
        :param MongoCursor collection: the collection where user sessoins are located
        """

        self.session_collection = collection

        self.config = Configuration.load(self.CONF_PATH, Ini)
        session = self.config.get('SESSION', {})

        self.alive_session_duration = int(session.get('alive_session_duration',
                                                      DEFAULT_ALIVE_SESSION_DURATION))

    def keep_alive(self, username):
        """
        Keep session alive by setting the ``last_check`` field
        to current timestamp.

        :param username: user identifier
        :type username: string
        :returns: last check timestamp
        :rtype: timestamp
        """
        now = time()
        self.session_collection.update({'_id': username},
                                       {'last_check': now})

        return now

    def is_session_active(self, username):
        """
        Check if session is active.
        If the session isn't found, then it is considered inactive.

        :param username: user identifier
        :type username: string
        :returns: check if session is active or not
        :rtype: bool
        """
        session = self.session_collection.find_one({'_id': username})

        if session is None:
            return False

        return session.get('active', False)

    def session_start(self, username):
        """
        Make session active for a user.

        :param username: user identifier
        :type username: string
        :returns: when the has started
        :rtype: timestamp or None
        """

        if not self.is_session_active(username):
            now = time()
            element = {
                'session_start': now,
                'last_check': now,
                'active': True,
                'session_stop': -1
            }
            self.session_collection.update({'_id': username},
                                           element,
                                           upsert=True)

            return now

    def sessions_close(self):
        """
        Close sessions that are expired (last_check + session_duration =< now)

        :returns: Closed sessions
        :rtype: list
        """
        now = time()
        inactive_ts = now - self.alive_session_duration
        query = {
            'last_check': {'$lte': inactive_ts},
            'session_stop': -1,
            'active': True
        }
        inactive_sessions = list(self.session_collection.find(query))

        for session in inactive_sessions:
            session['session_stop'] = now
            session['active'] = False

            self.session_collection.update({'_id': session['_id']},
                                           session,
                                           upsert=True)

        return inactive_sessions
