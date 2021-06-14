#!/usr/bin/env python2.7
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

import unittest

import canopsis.auth.cas as cas

import canopsis.auth.mock as mock


class TestCASBackend(unittest.TestCase):
    def mock_canopsis(self):
        cas.get_account = mock.auth.mock_get_account
        cas.delete_session = lambda: mock.auth.mock_delete_session(self)

        self.backend.install_account = lambda a: True

    def mock_bottle(self):
        self.session = mock.bottle.MockSession()
        self.request = mock.bottle.MockRequest()
        self.response = mock.bottle.MockResponse()

        cas.request = self.request
        cas.response = self.response
        cas.redirect = mock.bottle.mock_redirect()

    def mock_requests(self):
        self.requests = mock.requests

        cas.requests = self.requests

    def setUp(self):
        self.server = 'http://localhost:5000'
        self.service = 'http://localhost:8082'

        self.backend = cas.get_backend()

        self.mock_canopsis()
        self.mock_bottle()
        self.mock_requests()

    def try_auth(self):
        return self.backend.do_auth(self.session, self.server, self.service)

    def try_unauth(self):
        self.backend.undo_auth(self.session, self.server, self.service)

    def test_cas_auth_step1(self):
        res = self.try_auth()

        self.assertTrue(res is None)

    def test_cas_auth_post_step1(self):
        self.session['username'] = 'canopsis'
        self.session['password'] = 'canopsis'

        url = '{0}/login?service={1}/logged_in'.format(self.server, self.service)

        res = self.try_auth()

        self.assertEqual(res, 'username=canopsis&password=canopsis')
        self.assertEqual(self.response.status, 307)
        self.assertEqual(self.response.header('Location'), url)

    def test_cas_auth_step2(self):
        self.request.params['ticket'] = 'cpsticket'

        mock.requests.Response.status_code = 200
        mock.requests.Response.content = """
<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
    <cas:authenticationSuccess>
        <cas:user>linkdd</cas:user>
    </cas:authenticationSuccess>
</cas:serviceResponse>
"""

        res = self.try_auth()

        self.assertTrue(res)
        self.assertTrue(self.session.get('auth_cas', False))

    def test_cas_auth_step2_http_fail(self):
        self.request.params['ticket'] = 'ticket'

        mock.requests.Response.status_code = 500

        res = self.try_auth()

        self.assertFalse(res)

    def test_cas_auth_step2_validate_fail(self):
        self.request.params['ticket'] = 'cpsticket'

        mock.requests.Response.status_code = 200
        mock.requests.Response.content = """
<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
    <cas:authenticationFailure code="INVALID_TICKET">
        service ticket cpsticket has already been used
    </cas:authenticationFailure>
</cas:serviceResponse>
"""

        res = self.try_auth()

        self.assertFalse(res)

    def test_undo_auth(self):
        self.session['auth_cas'] = True

        self.try_unauth()

        self.assertFalse(self.session.get('auth_on', False))
