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
from json import dumps
from time import time, sleep

from canopsis.event import forger
from test_base import BaseApiTest, Method, HTTP


class BasicWeatherAPITest(BaseApiTest):

    """
    Helper class wich initialize events, watchers and pbehavior for all tests.
    """

    def init_tests(self):
        """
        Create basic objects that will be manipulated.
        """
        self.event1 = forger(
            event_type='check',
            connector='cap_kirk',
            connector_name='spock',
            source_type='resource',
            component='mc_coy',
            resource='uhura',
            state=2,
            output='NCC_1701'
        )

        self.event2 = forger(
            event_type='check',
            connector='cap_kirk',
            connector_name='spock',
            source_type='resource',
            component='zulu',
            resource='chekov',
            state=3,
            output='NCC_1701-B'
        )

        # Retrieve futur event id
        get_entity_id = '{}/api/v2/context_graph/get_id/'.format(self.URL_BASE)
        self.event1_id = self._send(url=get_entity_id,
                                    method=Method.post,
                                    data=dumps(self.event1)).json()
        self.event2_id = self._send(url=get_entity_id,
                                    method=Method.post,
                                    data=dumps(self.event2)).json()

        # Simple watcher (to insert)
        self.watcher_dict = {
            "description": "a_description",
            "display_name": "a_displayed_name",
            "enable": True,
            "mfilter": dumps({'_id': {'$in': [self.event1_id, self.event2_id]}}),
            "_id": "watcher_id"
        }
        self.watcher_id = self.watcher_dict['_id']

        now = int(time())
        self.pbehavior1 = {
            'name': 'imagine',
            'author': 'lennon',
            'filter_': {'_id': self.event1_id},
            'rrule': None,
            'tstart': now,
            'tstop': now + 60 * 60,
        }


class TestWeatherAPI_Empty(BasicWeatherAPITest):

    def setUp(self):
        self._authenticate()  # default setup
        self.init_tests()

        self.base = '{}/{}'.format(self.URL_BASE, '/api/v2/weather/watchers')

    def test_weather_get_watcher_empty(self):
        # Cleanup
        watcher_url = '{}/api/v2/watchers/{}'.format(self.URL_BASE,
                                                     self.watcher_id)
        self._send(url=watcher_url,
                   method=Method.delete)

        # Without mongo filter
        r = self._send(url=self.base,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.NOT_FOUND.value)

        # With a mongo filter
        r = self._send(url=self.base + '/' + '{}',
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 0)


class TestWeatherAPI(BasicWeatherAPITest):

    """
    NB : Because of api/v2/compute-pbehaviors, tests can be runned once per
         minute
    """

    def setUp(self):
        self._authenticate()  # default setup
        self.init_tests()

        self.base = '{}/{}'.format(self.URL_BASE, '/api/v2/weather/watchers')

    def init_tests(self):
        super(TestWeatherAPI, self).init_tests()

        # Adding watcher and alarms to watch upon
        watcher_url = '{}/api/v2/watchers'.format(self.URL_BASE)
        r = self._send(url=watcher_url,
                       method=Method.post,
                       data=dumps(self.watcher_dict))
        self.assertEqual(r.status_code, HTTP.OK.value)

        self.event_url = '{}/api/v2/event'.format(self.URL_BASE)
        r = self._send(url=self.event_url,
                       method=Method.post,
                       data=dumps(self.event1))
        self.assertEqual(r.status_code, HTTP.OK.value)

        self.pbheavior_url = '{}/api/v2/pbehavior'.format(self.URL_BASE)
        self.pbehavior_ids = []

        self.context_url = '{}/api/v2/context'.format(self.URL_BASE)

    def tearDown(self):
        """Deleting contextual informations"""
        watcher_url = '{}/api/v2/watchers/{}'.format(self.URL_BASE,
                                                     self.watcher_id)
        self._send(url=watcher_url,
                   method=Method.delete)

        for pbehavior_id in self.pbehavior_ids:
            self._send(url=self.pbheavior_url + '/' + pbehavior_id,
                       method=Method.delete)

        for event_id in [self.event1_id, self.event2_id]:
            self._send(url=self.context_url + '/' + event_id,
                       method=Method.delete)

    def get_watcher(self, watcher_filter):
        """
        Helper to read watcher informations (the tested interface).
        """
        r = self._send(url=self.base + '/' + watcher_filter,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)

        return r.json()

    def test_weather_get_watcher(self):

        watcher_filter = dumps({'_id': self.watcher_id})
        r = self._send(url=self.base + '/' + watcher_filter,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 1)
        self.assertEqual(json[0]['state']['val'], 2)
        self.assertFalse(json[0]['hasactivepbehaviorinentities'])
        self.assertFalse(json[0]['hasallactivepbehaviorinentities'])

        # Adding a new pbehavior
        r = self._send(url=self.pbheavior_url,
                       method=Method.post,
                       data=dumps(self.pbehavior1))
        self.pbehavior_ids.append(r.json())
        self.assertEqual(r.status_code, HTTP.OK.value)

        # Force compute on pbehaviors
        pbheavior_url = '{}/api/v2/compute-pbehaviors'.format(self.URL_BASE)
        r = self._send(url=pbheavior_url,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(json)

        # Verifying watcher state
        json = self.get_watcher(watcher_filter)
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 1)
        self.assertFalse(json[0]['hasactivepbehaviorinentities'])
        self.assertTrue(json[0]['hasallactivepbehaviorinentities'])

        # Sending another event
        r = self._send(url=self.event_url,
                       method=Method.post,
                       data=dumps(self.event2))
        self.assertEqual(r.status_code, HTTP.OK.value)

        sleep(1)

        # Verifying watcher state
        json = self.get_watcher(watcher_filter)
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 1)
        self.assertTrue(json[0]['hasactivepbehaviorinentities'])
        self.assertFalse(json[0]['hasallactivepbehaviorinentities'])

    # TODO: def test_weather_weatherwatchers(self):
