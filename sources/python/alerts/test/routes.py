#!/usr/bin/env python
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

import json
import requests
from unittest import main, TestCase

from canopsis.middleware.core import Middleware

# TODO: replace session.post(data=...) with 'json=' param, on requests>=2.4.2 only


def is_successful(r):
    # requests built-in exception handler. Is None if okay
    r.raise_for_status()
    # additional response validation:
    try:
        assert r.headers['content-type'] == "application/json", \
            "Reponse is not JSON format ({})".format(r.headers.get('content-type', 'UNKNOWN'))
    except AssertionError as e:
        print(e)
        return False
    else:
        return True


class TestRoutes(TestCase):
    def setUp(self):
        self.headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        }
        self.url = 'http://localhost:8082/'
        self.session = requests.Session()

        user_storage = Middleware.get_middleware_by_uri(
            'storage-default-rights://'
        )
        authkey = user_storage.find_elements(query={'_id': 'root'})
        authkey = list(authkey)[0]['authkey']

        link = "{}?authkey={}".format(self.url, authkey)
        login = self.session.get(link)

        self.cookies = login.cookies
        self.assertEqual(login.status_code, 200)

    def test_alert_filter_routes(self):
        alarm_id = 'fake_alarm_id'
        filter_id = 'fake_filter_id'
        base = '{}{}'.format(self.url, 'api/v2/alerts/filters')

        # GET
        r = self.session.get(base, cookies=self.cookies)
        self.assertEqual(r.status_code, 405)

        r = self.session.get(base + '/' + filter_id,
                             cookies=self.cookies)
        self.assertTrue(is_successful(r))
        self.assertEqual(r.json(), [])

        # POST
        r = self.session.post(base, cookies=self.cookies, headers=self.headers)
        self.assertEqual(r.status_code, 400)

        params = {
            'condition': {},
            'entity_filter': {"d": {"$eq": alarm_id}},
            'limit': 30.0,
            'tasks': ['alerts.systemaction.status_increase',
                      'alerts.systemaction.status_decrease'],
        }
        r = self.session.post(base,
                              data=json.dumps(params),
                              cookies=self.cookies, headers=self.headers)
        self.assertTrue(is_successful(r))
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        filter_id = data['_id']  # Get the real filter_id
        self.assertEqual(data['limit'], 30.0)

        # PUT
        r = self.session.put(base, cookies=self.cookies, headers=self.headers)
        self.assertEqual(r.status_code, 405)

        params = {
            'condition': {'key': {'$neq': 'value'}},
            'repeat': 6
        }
        r = self.session.put(base + '/' + filter_id,
                             data=json.dumps(params),
                             headers=self.headers, cookies=self.cookies)
        self.assertTrue(is_successful(r))
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        self.assertEqual(data['_id'], filter_id)
        self.assertTrue('condition' in data)
        self.assertEqual(data['repeat'], 6)

        # GET (again)
        r = self.session.get(base + '/' + filter_id,
                             cookies=self.cookies)
        self.assertTrue(is_successful(r))
        data = r.json()
        self.assertTrue(isinstance(data, list))
        self.assertEqual(data[0]['_id'], filter_id)
        self.assertTrue('condition' in data[0])

        # DELETE
        r = self.session.delete(base, cookies=self.cookies)
        self.assertEqual(r.status_code, 405)

        r = self.session.delete(base + '/' + filter_id,
                                cookies=self.cookies)
        self.assertTrue(is_successful(r))
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        self.assertTrue('ok' in data)


if __name__ == '__main__':
    # TODO: run webserver with testing storage instead of real ones
    # actually, the database can become unclean !
    #main()
    pass
