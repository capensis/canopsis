#!/usr/bin/env python2
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

from __future__ import unicode_literals

from enum import Enum
import re
import requests
import unittest


class HTTP(Enum):
    OK = 200
    ERROR = 400
    NOT_FOUND = 404
    NOT_ALLOWED = 405


class Method(Enum):
    """
    List of accepted HTTP methods
    """
    get = 'GET'
    post = 'POST'
    put = 'PUT'
    patch = 'PATCH'
    delete = 'DELETE'
    # And CONNECT, OPTIONS


class BaseApiTest(unittest.TestCase):

    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    }

    #WEB_HOST = "localhost"
    WEB_HOST = "172.16.0.20"

    URL_AUTHKEY = "{}/?authkey={}"
    URL_PLAIN = "{}/auth"

    """
    def _authent_with_authkey(self):
        # Send authentification with a authkey by reading it in the database.
        from canopsis.middleware.core import Middleware
        user_storage = Middleware.get_middleware_by_uri(
            'storage-default-rights://'
        )
        authkey = user_storage.find_elements(query={'_id': 'root'})
        authkey = list(authkey)[0]['authkey']
        url_auth = self.URL_AUTH.format(self.URL_BASE, authkey)

        return self._send(url_auth)
    """

    def _authent_plain(self):
        """
        Send authentification with login/passwd.
        """
        form = {
            'username': "root",
            'password': "root"
        }
        headers = {
            'Content-type': "application/x-www-form-urlencoded"
        }
        url_auth = self.URL_PLAIN.format(self.URL_BASE)

        return self._send(url_auth, data=form, headers=headers, method=Method.post)

    def _authenticate(self):
        self.URL_BASE = "http://{}:8082".format(self.WEB_HOST)

        self.session = requests.Session()
        response = self._authent_plain()
        #print("Login on {}".format(response.request.path_url))

        # Auth validation
        self.assertEqual(response.status_code, 200)
        if re.search("<title>Canopsis | Login</title>", response.text)\
           is not None:
            self.fail("Authentication error.")

        self.cookies = response.cookies

    def _send(self, url, data=None, headers=None, method=Method.get):
        """Send a http request.

        :param str url: the url to access
        :param str data: data to with with the request
        :param dict headers: change headers on the request
        :param Method method: which method to use
        :rtype: <Response>
        """
        kwargs = {
            'method': method.value,
            'headers': self.headers if headers is None else headers
        }
        if hasattr(self, 'cookies'):
            kwargs['cookies'] = self.cookies
        if data is not None:
            kwargs['data'] = data

        response = self.session.request(url=url, **kwargs)

        return response

    def setUp(self):
        self._authenticate()
