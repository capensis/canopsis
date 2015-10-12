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
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.parameters import Parameter
from time import time

CONF_PATH = 'session/session.conf'
CATEGORY = 'SESSION'

CONFIG = [
    Parameter('alive_session_duration', int),
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONFIG)
class Session(MiddlewareRegistry):

    """
    Manage session information in Canopsis
    """

    ENTITY_STORAGE = 'session_storage'

    @property
    def alive_session_duration(self):
        if not hasattr(self, '_alive_session_duration'):
            self.alive_session_duration = 0

        return self._alive_session_duration

    @alive_session_duration.setter
    def alive_session_duration(self, value):
        self._alive_session_duration = value

    def keep_alive(self, username):
        self[Session.ENTITY_STORAGE].put_element(
            _id=username, element={'last_check': time()}
        )

    def is_user_session_active(self, username):

        """
        Returns wether or not the user session is active
        :param: username string
        """

        session = self.get_user_info(username)
        if not session:
            return False
        else:
            return session['active']

    def session_start(self, username):

        """
        Starts a session for a given username
        :param: username string
        """

        # active session test avoid start session
        # reset if session is still active
        if not self.is_user_session_active(username):

            self[Session.ENTITY_STORAGE].put_element(
                _id=username,
                element={
                    'session_start': time(),
                    'last_check': time(),
                    'active': True,
                    'session_stop': -1
                }
            )

    def get_new_inactive_sessions(self):

        """
        Retrieve a user session list that are newly inactive.
        Sets the session stop time if session is found as inactive
        since at least
        :returns: a list of user session newly inactive
        """

        active_limit_date = time() - self.alive_session_duration

        sessions = self[Session.ENTITY_STORAGE].get_elements(
            query={
                'last_check': {'$lte': active_limit_date}
            }
        )

        new_inactive_sessions = []
        # Upsert end date if not already set
        now = time()
        for session in sessions:
            if session['session_stop'] == -1:
                session['session_stop'] = now

                self[Session.ENTITY_STORAGE].put_element(
                    _id=session['_id'],
                    element={
                        'session_stop': now,
                        'active': False
                    }
                )
                new_inactive_sessions.append(session)

        return new_inactive_sessions

    def get_delta_session_time_metrics(self, sessions):

        """
        Compute metrics from a user session list
        :param: sessions is a list of user session (session records)
        :returns: a list of metrics formated as an event perfdata array
        """

        metrics = []
        for session in sessions:
            delta = session['session_stop'] - session['session_start']
            metrics.append({
                'type': 'COUNTER',
                'value': delta,
                'metric': 'cps_session_delay_user_{}'.format(session['_id'])
            })
        return metrics

    def get_user_info(self, username):
        """
        Retrieve database user session information
        :param: username
        :returns: a dict with user's session information
        """

        return self[Session.ENTITY_STORAGE].get_elements(ids=username)
