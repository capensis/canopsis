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
        alarm_id = '/fake/alarm/id'
        base = '{}{}'.format(self.url, 'alerts-filter')

        # GET
        r = self.session.get(base, cookies=self.cookies)
        self.assertTrue(is_successful(r))
        self.assertFalse(r.json()['success'])

        r = self.session.get(base, params={'entity_id': alarm_id},
                             cookies=self.cookies)
        self.assertTrue(is_successful(r))
        rjson = r.json()
        self.assertTrue('data' in rjson)
        self.assertEqual(rjson['data'], [])

        # PUT
        r = self.session.put(base, cookies=self.cookies)
        self.assertTrue(is_successful(r))
        self.assertFalse(r.json()['success'])

        params = {
            "element": {
                'limit': 30.0,
                'condition': {},
                'tasks': ['alerts.systemaction.status_increase',
                          'alerts.systemaction.status_decrease'],
                'entity_filter': {"d": {"$eq": alarm_id}},
            }
        }
        r = self.session.put(base,
                             data=json.dumps(params),
                             cookies=self.cookies, headers=self.headers)
        self.assertTrue(is_successful(r))
        data = r.json()['data']
        self.assertEqual(len(data), 1)
        filter_id = data[0]['_id']
        self.assertEqual(data[0]['limit'], 30.0)

        # POST
        r = self.session.post(base, cookies=self.cookies)
        self.assertTrue(is_successful(r))
        self.assertFalse(r.json()['success'])

        params = {
            'entity_id': filter_id,
            'key': 'condition',
            'value': {'key': {'$neq': 'value'}},
        }
        r = self.session.post(base, data=json.dumps(params),
                              headers=self.headers, cookies=self.cookies)
        self.assertTrue(is_successful(r))
        data = r.json()['data']
        self.assertEqual(len(data), 1)
        self.assertEqual(data[0]['_id'], filter_id)
        self.assertTrue('condition' in data[0])

        # GET (again)
        r = self.session.get(base,
                             params={'entity_id': filter_id},
                             cookies=self.cookies)
        self.assertTrue(is_successful(r))
        data = r.json()['data']
        self.assertEqual(len(data), 1)
        self.assertEqual(data[0]['_id'], filter_id)
        self.assertTrue('condition' in data[0])

        # DELETE
        r = self.session.delete(base, cookies=self.cookies)
        self.assertTrue(is_successful(r))
        self.assertFalse(r.json()['success'])

        r = self.session.delete(base,
                                params={'entity_id': filter_id},
                                cookies=self.cookies)
        self.assertTrue(is_successful(r))
        data = r.json()['data']
        self.assertEqual(len(data), 1)
        self.assertTrue('ok' in data[0])


if __name__ == '__main__':
    # TODO: run webserver with testing storage instead of real ones
    # actually, the database can become unclean !
    #main()
    pass
