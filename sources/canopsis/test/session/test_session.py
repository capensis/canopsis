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
from canopsis.common.middleware import Middleware
from canopsis.session.manager import Session
import unittest
from canopsis.common import root_path
import xmlrunner


class SessionManagerTest(TestCase):

    def setUp(self):
        self.storage = Middleware.get_middleware_by_uri(
            'mongodb-default-testsession://'
        )
        self.collection = MongoCollection(self.storage._backend)

        self.manager = Session(collection=self.collection)

        self.user = 'test_user'
        self.id_beaker_session = 'cm9vdF8xNTc2MDY1MzY2'
        self.path = ["view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8","view-tab_edd5855b-54f1-4c51-9550-d88c2da60768"]
        self.path_bis = ["view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8","view-tab_edd5855b-54f1-4c51-azerty"]

    def tearDown(self):
        self.collection.remove()

    def test_keep_alive(self):
        self.manager.session_start(self.id_beaker_session,self.user)
        sleep(1)
        got = self.manager.keep_alive(self.id_beaker_session,self.user,True,self.path)

        session = self.collection.find_one({"id_beaker_session": self.id_beaker_session,'username': self.user})

        self.assertTrue(isinstance(session, dict))
        self.assertEqual(got, session['last_ping'])
        self.assertEqual(self.path,session['last_visible_path'])

        got = self.manager.keep_alive(self.id_beaker_session,self.user,False,self.path_bis)
        session = self.collection.find_one({"id_beaker_session": self.id_beaker_session,'username': self.user})

        self.assertTrue(isinstance(session, dict))
        self.assertEqual(got, session['last_ping'])
        self.assertEqual(self.path,session['last_visible_path'])

        got = self.manager.keep_alive(self.id_beaker_session,self.user,True,self.path_bis)
        session = self.collection.find_one({"id_beaker_session": self.id_beaker_session,'username': self.user})
        self.assertTrue(isinstance(session, dict))
        self.assertEqual(got, session['last_ping'])
        self.assertEqual(self.path_bis,session['last_visible_path'])

    def test_session_start(self):
        got = self.manager.session_start(self.id_beaker_session,self.user)

        session = self.collection.find_one({"id_beaker_session": self.id_beaker_session,'username': self.user})

        self.assertTrue(isinstance(session, dict))
        self.assertTrue(self.manager.is_session_active(self.id_beaker_session))
        self.assertEqual(got, session['start'])

    def test_session_start_already_started(self):
        self.test_session_start()

        got = self.manager.session_start(self.id_beaker_session,self.user)

        self.assertTrue(got is None)

    def test_is_session_active(self):
        self.assertFalse(self.manager.is_session_active(self.id_beaker_session))
        self.manager.session_start(self.id_beaker_session,self.user)
        self.assertTrue(self.manager.is_session_active(self.id_beaker_session))

    def test_session_hide(self):
        self.manager.session_start(self.id_beaker_session,self.user)
        sleep(1)

        got = self.manager.session_hide(self.id_beaker_session,self.user,self.path)

        session = self.collection.find_one({"id_beaker_session": self.id_beaker_session,'username': self.user})

        self.assertTrue(isinstance(session, dict))
        self.assertEqual(got, session['last_ping'])
        self.assertEqual(self.path,session['last_visible_path'])

        got = self.manager.session_hide(self.id_beaker_session,self.user,self.path_bis)
        session = self.collection.find_one({"id_beaker_session": self.id_beaker_session,'username': self.user})

        self.assertTrue(isinstance(session, dict))
        self.assertEqual(got, session['last_ping'])
        self.assertEqual(self.path_bis,session['last_visible_path'])

        got = self.manager.session_hide(self.id_beaker_session,self.user,self.path_bis)
        session = self.collection.find_one({"id_beaker_session": self.id_beaker_session,'username': self.user})
        self.assertTrue(isinstance(session, dict))
        self.assertEqual(got, session['last_ping'])
        self.assertEqual(self.path_bis,session['last_visible_path'])

    def test_sessions_req(self):
        self.manager.session_start(self.id_beaker_session,self.user)
        sleep(1)
        session = self.collection.find_one({"id_beaker_session": self.id_beaker_session,'username': self.user})
        session_req = self.manager.sessions_req(self.id_beaker_session,{"active":"true"})
        self.assertEqual([session],session_req)

        session_req = self.manager.sessions_req(self.id_beaker_session,{"active":"false"})
        self.assertEqual([],session_req)

        self.manager.session_start("azerty","userTest")

        session2 = self.collection.find_one({"id_beaker_session": "azerty",'username': "userTest"})
        session_req = self.manager.sessions_req(self.id_beaker_session,{"usernames[]":[self.user,"userTest"]})
        self.assertEqual([session,session2],session_req)



if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
