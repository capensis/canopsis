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

from time import time

from canopsis.confng import Configuration, Ini
from canopsis.middleware.core import Middleware

DEFAULT_SESSION_STORAGE_URI = 'mongodb-default-session://'
DEFAULT_METRIC_PRODUCER_VALUE = 'canopsis.stats.producers.metric.MetricProducer'
DEFAULT_PERFDATA_MANAGER_VALUE = 'canopsis.perfdata.manager.PerfData'
DEFAULT_ALIVE_SESSION_DURATION = 300


class Session():
    """
    Manage session informations.
    """

    CONF_PATH = 'etc/session/session.conf'

    SESSION_STORAGE = 'session_storage'
    METRIC_PRODUCER = 'metric_producer'
    PERFDATA_MANAGER = 'perfdata_manager'

    def __init__(
        self,
        session_storage=None,
        metric_producer=None,
        perfdata_manager=None,
        *args, **kwargs
    ):
        self.session_storage = session_storage
        self.metric_producer = metric_producer
        self.perfdata_manager = perfdata_manager

        self.config = Configuration.load(self.CONF_PATH, Ini)
        session = self.config.get('SESSION', {})
        if session_storage is None:
            self.session_storage_uri = session.get('session_storage_uri',
                                                   DEFAULT_SESSION_STORAGE_URI)
            self.session_storage = Middleware.get_middleware_by_uri(
                self.session_storage_uri
            )

        self.metric_producer_value = session.get('metric_producer_value',
                                                 DEFAULT_METRIC_PRODUCER_VALUE)
        self.perfdata_manager_value = session.get('perfdata_manager_value',
                                                  DEFAULT_PERFDATA_MANAGER_VALUE)
        self.alive_session_duration = int(session.get('alive_session_duration',
                                                      DEFAULT_ALIVE_SESSION_DURATION))

    def keep_alive(self, username):
        """
        Keep session alive by setting the ``last_check`` field
        to current timestamp.

        :param username: user identifier
        :type username: string

        :returns: check timestamp
        """

        now = time()
        self.session_storage.put_element(
            _id=username, element={'last_check': now}
        )
        return now

    def is_session_active(self, username):
        """
        Check if session is active.
        If the session isn't found, then it is considered inactive.

        :param username: user identifier
        :type username: string

        :returns: True if session is active, False otherwise
        """

        session = self.session_storage.get_elements(ids=username)

        if session is None:
            return False

        else:
            return session['active']

    def session_start(self, username):
        """
        Make session active for a user.

        :param username: user identifier
        :type username: string

        :returns: Start timestamp, or None if already started
        """

        if not self.is_session_active(username):
            now = time()

            self.session_storage.put_element(
                _id=username,
                element={
                    'session_start': now,
                    'last_check': now,
                    'active': True,
                    'session_stop': -1
                }
            )

            return now

    def sessions_close(self):
        """
        Close sessions that are expired (last_check + session_duration =< now)

        :returns: Closed sessions
        :rtype: list
        """
        now = time()
        inactive_ts = now - self.alive_session_duration
        inactive_sessions = list(self.session_storage.get_elements(
            query={
                'last_check': {'$lte': inactive_ts},
                'session_stop': -1,
                'active': True
            }
        ))

        # Update sessions in storage
        for session in inactive_sessions:
            session['session_stop'] = now
            session['active'] = False

            self.session_storage.put_element(element=session)

        return inactive_sessions
