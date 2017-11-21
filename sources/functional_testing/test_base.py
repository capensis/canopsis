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
from urllib import quote
import requests
import unittest


class Event:

    def __init__(self, connector, connector_name, component, state, state_type,
                 status, source_type, resource=None):
        self.event_type = "check"
        self.connector = connector
        self.connector_name = connector_name
        self.component = component
        self.resource = resource
        self.state_type = state_type
        self.status = status
        self.source_type = source_type

class HTTP(Enum):
    """
    HTTP codes transcription.
    """
    OK = 200
    ERROR = 400
    NOT_FOUND = 404
    NOT_ALLOWED = 405


class Method(Enum):
    """
    List of accepted HTTP methods.
    """
    get = 'GET'
    post = 'POST'
    put = 'PUT'
    patch = 'PATCH'
    delete = 'DELETE'
    # And CONNECT, OPTIONS


class BaseApiTest(unittest.TestCase):

    """
    Generic class to instanciate an API test.
    """

    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    }

    WEB_HOST = "192.168.0.74"
    WEB_PORT = "8082"

    URL_BASE = "http://{}:{}".format(WEB_HOST, WEB_PORT)
    URL_PLAIN = "{}/auth".format(URL_BASE)

    # URL_AUTHKEY = "{}/?authkey={}"
    # def _authent_with_authkey(self):
    #     # Send authentification with a authkey by reading it in the database.
    #     from canopsis.middleware.core import Middleware

    #     user_storage = Middleware.get_middleware_by_uri(
    #         'storage-default-rights://'
    #     )
    #     authkey = user_storage.find_elements(query={'_id': 'root'})
    #     authkey = list(authkey)[0]['authkey']
    #     url_auth = self.URL_AUTH.format(self.URL_BASE, authkey)

    #     return self._send(url_auth)

    def _build_url(self, url, params):
        """add every key/value in at the end of the url.
        url utr : the url
        param dict : the parameters to add at the end of the url
        return str : an url with every parameters concatenated
        """

        querystring  = ""

        keys = params.keys()
        for i in range(len(keys)):
            key = keys[i]
            if i != len(keys) - 1:
                sep = "&"
            else:
                sep = ""

            querystring += "{0}={1}{2}".format(key, params[key], sep)
        return self.URL_BASE + url + "?" + quote(querystring)

    def _authent_plain(self):
        """
        Send authentification through clear login/passwd.
        """
        form = {
            'username': "root",
            'password': "root"
        }
        headers = {
            'Content-type': "application/x-www-form-urlencoded"
        }

        return self._send(self.URL_PLAIN,
                          data=form,
                          headers=headers,
                          method=Method.post,
                          allow_redirects=False)

    def _authenticate(self):
        """
        Do the authentification.
        """
        self.session = requests.Session()
        response = self._authent_plain()
        # print("Login on {}".format(response.request.path_url))

        # Auth validation
        self.assertEqual(response.status_code, 303)
        #search_title = re.search("<title>Canopsis | Login</title>",
        #                         response.text)
        #if search_title is not None:
        #    self.fail("Authentication error.")

        self.cookies = response.cookies

    def _send(self,
              url,
              data=None,
              headers=None,
              method=Method.get,
              params=None,
              allow_redirects=True
              ):
        """Send a http request.

        :param str url: the url to access
        :param str data: data to with with the request
        :param dict headers: change headers on the request
        :param Method method: which method to use
        :param dict params: querystring parameters
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
        if params is not None:
            kwargs['params'] = params

        response = self.session.request(url=url,
                                        allow_redirects=allow_redirects,
                                        **kwargs)

        return response

    def _send_event(self, event):
        return self._send(url=self.URL_BASE + "/event",
                   method=Method.post,
                   params={"event": event.__dict__},
                   allow_redirects=False)

    def setUp(self):
        self._authenticate()
