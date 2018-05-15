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

from __future__ import unicode_literals

import json
import time

from test_base import BaseApiTest, Method, HTTP


class TestAlertsAPI(BaseApiTest):

    def setUp(self):

        self._authenticate()  # default setup

        self.base_filter = '{}/{}'.format(self.URL_BASE,
                                          'api/v2/alerts/filters')

    def test_alert_filters(self):
        alarm_id = 'fake_alarm_id'
        filter_id = 'fake_filter_id'

        # GET
        r = self._send(url=self.base_filter)
        self.assertEqual(r.status_code, HTTP.NOT_ALLOWED.value)

        r = self._send(url=self.base_filter + '/' + filter_id,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertEqual(r.json(), [])

        # POST
        r = self._send(url=self.base_filter,
                       method=Method.post)
        self.assertEqual(r.status_code, HTTP.ERROR.value)

        params = {
            'condition': {},
            'entity_filter': {"d": {"$eq": alarm_id}},
            'limit': 30.0,
            'tasks': ['alerts.systemaction.status_increase',
                      'alerts.systemaction.status_decrease'],
        }
        r = self._send(url=self.base_filter,
                       data=json.dumps(params),
                       method=Method.post)
        self.assertEqual(r.status_code, HTTP.OK.value)
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        filter_id = data['_id']  # Get the real filter_id
        self.assertEqual(data['limit'], 30.0)

        # PUT
        r = self._send(url=self.base_filter,
                       method=Method.put)
        self.assertEqual(r.status_code, HTTP.NOT_ALLOWED.value)

        params = {
            'condition': {'key': {'$neq': 'value'}},
            'repeat': 6
        }
        r = self.session.put(self.base_filter + '/' + filter_id,
                             data=json.dumps(params),
                             headers=self.headers, cookies=self.cookies)
        self.assertEqual(r.status_code, HTTP.OK.value)
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        self.assertEqual(data['_id'], filter_id)
        self.assertTrue('condition' in data)
        self.assertEqual(data['repeat'], 6)

        # GET (again)
        r = self._send(url=self.base_filter + '/' + filter_id)
        self.assertEqual(r.status_code, HTTP.OK.value)
        data = r.json()
        self.assertTrue(isinstance(data, list))
        self.assertEqual(data[0]['_id'], filter_id)
        self.assertTrue('condition' in data[0])

        # DELETE
        r = self._send(url=self.base_filter,
                       method=Method.delete)
        self.assertEqual(r.status_code, HTTP.NOT_ALLOWED.value)

        r = self._send(url=self.base_filter + '/' + filter_id,
                       method=Method.delete)
        self.assertEqual(r.status_code, HTTP.OK.value)
        data = r.json()
        self.assertTrue(isinstance(data, dict))
        self.assertTrue('ok' in data)

    def test_hide_resources(self):
        delete_all_alarms = self._send(
            url=self.URL_BASE + '/api/v2/alerts/{}',
            method=Method.delete
        )
        events = [
            {
                "event_type": "check",
                "connector": "connector",
                "connector_name": "connector_name",
                "component": "component",
                "source_type": "component",
                "state": 2
            },
            {
                "event_type": "check",
                "connector": "connector",
                "connector_name": "connector_name",
                "component": "component",
                "resource": "resource1",
                "source_type": "resource",
                "state": 2
            },
            {
                "event_type": "check",
                "connector": "connector",
                "connector_name": "connector_name",
                "component": "component",
                "resource": "resource2",
                "source_type": "resource",
                "state": 2
            },
            {
                "event_type": "check",
                "connector": "connector",
                "connector_name": "connector_name",
                "component": "component",
                "resource": "resource3",
                "source_type": "resource",
                "state": 2
            }
        ]
        r = self._send(
            url=self.URL_BASE + '/api/v2/event',
            data=json.dumps(events),
            method=Method.post
        )
        time.sleep(0.2)
        alarms = self._send(
            url=self.URL_BASE + '/alerts/get-alarms?hide_resources=false',
            method=Method.get
        )
        self.assertEqual(alarms.json().get('data')[0].get('total'), 4)
        alarms = self._send(
            url=self.URL_BASE + '/alerts/get-alarms?hide_resources=true',
            method=Method.get
        )
        self.assertEqual(alarms.json().get('data')[0].get('total'), 1)
        self.assertEqual(alarms.json().get('data')[0].get('alarms')[0].get('d'), 'component')
        event_component_minor = {
                "event_type": "check",
                "connector": "connector",
                "connector_name": "connector_name",
                "component": "component",
                "source_type": "component",
                "state": 1
            }
        r = self._send(
            url=self.URL_BASE + '/api/v2/event',
            data=json.dumps(event_component_minor),
            method=Method.post
        )
        time.sleep(0.2)
        alarms = self._send(
            url=self.URL_BASE + '/alerts/get-alarms?hide_resources=true',
            method=Method.get
        )
        self.assertEqual(alarms.json().get('data')[0].get('total'), 4)
