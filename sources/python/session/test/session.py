#!/usr/bin/env python
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

from unittest import TestCase, main
from canopsis.session.manager import Session
from time import time


class SessionManagerTest(TestCase):

    def setUp(self):
        self.session_manager = Session()
        self.user = 'test_user'
        self.session_manager[Session.ENTITY_STORAGE].remove_elements()


class SessionTest(SessionManagerTest):

    def test_get_user_info(self):
        self.assertIsNone(
            self.session_manager.get_user_info(self.user)
        )
        self.session_manager.keep_alive(self.user)
        session = self.session_manager.get_user_info(self.user)
        self.assertEqual(session['_id'], self.user)

    def test_keep_alive(self):
        self.session_manager.keep_alive(self.user)
        session = self.session_manager.get_user_info(self.user)
        self.assertLessEqual(session['last_check'], time())

    def test_session_start(self):
        self.session_manager.session_start(self.user)
        session = self.session_manager.get_user_info(self.user)
        self.assertLessEqual(session['session_start'], time())
        self.assertTrue(session['active'])

        first_start = session['session_start']
        self.session_manager.session_start(self.user)
        session = self.session_manager.get_user_info(self.user)
        self.assertEqual(first_start, session['session_start'])

    def test_is_user_session_active(self):
        self.assertFalse(
            self.session_manager.is_user_session_active(self.user)
        )
        self.session_manager.session_start(self.user)
        self.assertTrue(
            self.session_manager.is_user_session_active(self.user)
        )

    def test_check_inactive_sessions(self):
        # initilize session
        self.session_manager.session_start(self.user)

        sessions = list(self.session_manager.get_new_inactive_sessions())

        self.assertEqual(len(sessions), 0)

        delta = self.session_manager.alive_session_duration + 1
        self.session_manager[Session.ENTITY_STORAGE].put_element(
            _id=self.user,
            element={'last_check': time() - delta}
        )
        sessions = list(self.session_manager.get_new_inactive_sessions())

        self.assertEqual(len(sessions), 1)

        self.session_manager.session_start(self.user + '1')

        self.session_manager[Session.ENTITY_STORAGE].put_element(
            _id=self.user + '1',
            element={'last_check': time() - delta}
        )
        sessions = list(self.session_manager.get_new_inactive_sessions())
        self.assertEqual(len(sessions), 1)
        for session in sessions:
            self.assertIn('session_stop', session)

    def test_get_delta_session_time_metrics(self):
        # This should produce events metrics
        sessions = [
            {
                '_id': self.user,
                'session_start': 500,
                'session_stop': 1000
            }
        ]
        metrics = self.session_manager.get_delta_session_time_metrics(sessions)
        self.assertEqual(
            metrics[0],
            {
                'metric': 'cps_session_delay_user_test_user',
                'type': 'COUNTER',
                'value': 500
            }
        )
        # We have two session so two metrics
        sessions.append(sessions[0].copy())
        metrics = self.session_manager.get_delta_session_time_metrics(sessions)
        self.assertEqual(len(metrics), 2)
if __name__ == '__main__':
    main()
