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
from canopsis.basecrud.manager import BaseCrud
from time import time

CONF_PATH = 'session/session.conf'
CATEGORY = 'SESSION'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class Session(BaseCrud):

    """
    Manage session information in Canopsis
    """
    alive_session_duration = 60 * 5
    ENTITY_STORAGE = 'session_storage'

    def keep_alive(self, username):
        self.put(username, {'last_check': time()})

    def is_user_session_active(self, username):
        session = self.get_user_info(username)
        if not session:
            return False
        else:
            return session['active']

    def session_start(self, username):
        # active session test avoid start session
        # reset if session is still active
        if not self.is_user_session_active(username):
            self.put(username, {
                'session_start': time(),
                'active': True
            })

    def get_inactive_sessions(self):
        active_limit_date = time() - self.alive_session_duration
        sessions = list(self.find(query={
            'last_check': {'$lte': active_limit_date}
        }))

        # Upsert end date if not already set
        now = time()
        for session in sessions:
            if 'session_stop' not in session:
                session['session_stop'] = now
                self.put(session['_id'], {'session_stop': now})

        return sessions

    def get_delta_session_time_metrics(self, sessions):
        metrics = []
        for session in sessions:
            delta = session['session_stop'] - session['session_start']
            metrics.append({
                'type': 'COUNTER',
                'value': delta,
                'metric': 'cps_session_delta_user_{}'.format(session['_id'])
            })
        return metrics

    def get_user_info(self, username):
        return self.find(ids=username)
