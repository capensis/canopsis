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

from __future__ import unicode_literals

from time import sleep
from unittest import TestCase, main

from canopsis.common.collection import MongoCollection
from canopsis.middleware.core import Middleware
from canopsis.session.manager import Session


class SessionManagerTest(TestCase):

    def setUp(self):
        self.storage = Middleware.get_middleware_by_uri(
            'mongodb-default-testsession://'
        )
        self.collection = MongoCollection(self.storage._backend)

        self.manager = Session(collection=self.collection)

        self.user = 'test_user'

    def tearDown(self):
        self.collection.remove()

    def test_keep_alive(self):
        self.manager.session_start(self.user)
        sleep(1)
        got = self.manager.keep_alive(self.user)

        session = self.collection.find_one({'_id': self.user})

        self.assertTrue(session is not None)
        self.assertEqual(got, session['last_check'])

    def test_session_start(self):
        got = self.manager.session_start(self.user)

        session = self.collection.find_one({'_id': self.user})

        self.assertTrue(session is not None)
        self.assertTrue(session['active'])
        self.assertEqual(got, session['session_start'])

    def test_session_start_already_started(self):
        self.test_session_start()

        got = self.manager.session_start(self.user)

        self.assertTrue(got is None)

    def test_is_session_active(self):
        self.assertFalse(self.manager.is_session_active(self.user))
        self.manager.session_start(self.user)
        self.assertTrue(self.manager.is_session_active(self.user))

    def test_sessions_close(self):
        got = self.manager.session_start(self.user)

        self.manager.alive_session_duration = 0
        self.assertTrue(got is not None)

        sessions = self.manager.sessions_close()
        self.assertTrue(len(sessions) > 0)
        self.assertEqual(got, sessions[0]['last_check'])

if __name__ == '__main__':
    main()
