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
from canopsis.middleware.core import Middleware
from canopsis.session.manager import Session


class SessionManagerTest(TestCase):

    def setUp(self):
        self.storage = Middleware.get_middleware_by_uri(
            'mongodb-default-testsession://'
        )
        self.storage.connect()

        self.manager = Session()
        self.manager[Session.SESSION_STORAGE] = self.storage

        self.user = 'test_user'

    def tearDown(self):
        self.storage.remove_elements()
        self.storage.disconnect()

    def test_keep_alive(self):
        got = self.manager.keep_alive(self.user)

        session = self.storage.get_elements(ids=self.user)

        self.assertTrue(session is not None)
        self.assertEqual(got, session['last_check'])

    def test_session_start(self):
        got = self.manager.session_start(self.user)

        session = self.storage.get_elements(ids=self.user)

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

#TODO4-01-2017
#    def test_duration(self):
#        raise NotImplementedError('missing test')


if __name__ == '__main__':
    main()
