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

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.configuration.model import Parameter

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
    Manage session informations.
    """

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
        super(Session, self).__init__(*args, **kwargs)

        if session_storage is not None:
            self[Session.SESSION_STORAGE] = session_storage

        if metric_producer is not None:
            self[Session.METRIC_PRODUCER] = metric_producer

        if perfdata_manager is not None:
            self[Session.PERFDATA_MANAGER] = perfdata_manager

    @property
    def alive_session_duration(self):
        if not hasattr(self, '_alive_session_duration'):
            self.alive_session_duration = None

        return self._alive_session_duration

    @alive_session_duration.setter
    def alive_session_duration(self, value):
        if value is None:
            value = 300

        self._alive_session_duration = value

    def keep_alive(self, username):
        """
        Keep session alive by setting the ``last_check`` field
        to current timestamp.

        :param username: user identifier
        :type username: string

        :returns: check timestamp
        """

        now = time()
        self[Session.SESSION_STORAGE].put_element(
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

        session = self[Session.SESSION_STORAGE].get_elements(ids=username)

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

            self[Session.SESSION_STORAGE].put_element(
                _id=username,
                element={
                    'session_start': now,
                    'last_check': now,
                    'active': True,
                    'session_stop': -1
                }
            )

            return now

    def duration(self):
        """
        Return event, for each user, containing the session_duration metric.

        :returns: list of events
        """

        storage = self[Session.SESSION_STORAGE]

        now = time()
        inactive_ts = now - self.alive_session_duration
        inactive_sessions = list(storage.get_elements(
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

            storage.put_element(element=session)

        # Generate events
        events = []

        for session in inactive_sessions:
            duration = session['session_stop'] - session['session_start']
            event = {
                'timestamp': now,
                'connector': 'canopsis',
                'connector_name': 'session',
                'event_type': 'perf',
                'source_type': 'resource',
                'component': session[storage.ID],
                'resource': 'session_duration',
                'perf_data_array': [
                    {
                        'metric': 'last',
                        'value': duration,
                        'type': 'GAUGE',
                        'unit': 's'
                    },
                    {
                        'metric': 'sum',
                        'value': duration,
                        'type': 'COUNTER',
                        'unit': 's'
                    }
                ]
            }

            events.append(event)

            for operator in ['min', 'max', 'average']:
                perfdatamgr = self[Session.PERFDATA_MANAGER]
                producer = self[Session.METRIC_PRODUCER]

                entity = perfdatamgr.get_metric_entity('last', event)
                producer.may_create_stats_serie(entity, operator)

        return events
