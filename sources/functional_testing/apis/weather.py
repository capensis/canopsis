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
        self.context_url = '{}/api/v2/context'.format(self.URL_BASE)

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

        # Unlinked event
        self.event3 = forger(
            event_type='check',
            connector='picard',
            connector_name='ricker',
            source_type='resource',
            component='laforge',
            resource='worf',
            state=1,
            output='NCC_1701-D'
        )

        # Retrieve futur event id
        get_entity_id = '{}/api/v2/context_graph/get_id/'.format(self.URL_BASE)
        self.event1_id = self._send(url=get_entity_id,
                                    method=Method.post,
                                    data=dumps(self.event1)).json()
        self.event2_id = self._send(url=get_entity_id,
                                    method=Method.post,
                                    data=dumps(self.event2)).json()
        self.event3_id = self._send(url=get_entity_id,
                                    method=Method.post,
                                    data=dumps(self.event3)).json()

        # Simple watcher (to insert)
        self.watcher_1 = {
            "description": "first gen",
            "display_name": "st 1",
            "enable": True,
            "mfilter": dumps({'_id': {'$in': [self.event1_id, self.event2_id]}}),
            "_id": "watcher_first_gen"
        }

        # Simple watcher (to insert)
        self.watcher_2 = {
            "description": "pikes crew",
            "display_name": "st 1 - pilot",
            "enable": True,
            "mfilter": dumps({'_id': self.event1_id}),
            "_id": "watcher_pikes_crew"
        }

        self.watcher_3 = {
            "description": "next gen",
            "display_name": "st 2",
            "enable": True,
            "mfilter": dumps({'_id': {'$in': [self.event3_id]}}),
            "_id": "watcher_next_gen"
        }

        now = int(time())
        self.pbehavior1 = {
            'name': 'imagine',
            'author': 'lennon',
            'filter_': {'_id': self.event1_id},
            'rrule': None,
            'tstart': now,
            'tstop': now + 60 * 60,
        }

    def context_cleanup(self):
        # Cleanup existing watchers
        watcher_url = '{}/api/v2/watchers'.format(self.URL_BASE)
        self._send(url=watcher_url + '/' + self.watcher_1['_id'],
                   method=Method.delete)
        self._send(url=watcher_url + '/' + self.watcher_2['_id'],
                   method=Method.delete)
        self._send(url=watcher_url + '/' + self.watcher_3['_id'],
                   method=Method.delete)

        # Cleanup whole entity graph
        entity_ids = [
            self.event1_id, self.event2_id, self.event3_id,
            self.event1['component'], self.event2['component'], self.event3['component'],
            '{}/{}'.format(self.event1['connector'],
                           self.event1['connector_name']),
            '{}/{}'.format(self.event3['connector'],
                           self.event3['connector_name'])
        ]
        for entity_id in entity_ids:
            self._send(url=self.context_url + '/' + entity_id,
                       method=Method.delete)

"""
class TestWeatherAPI_Empty(BasicWeatherAPITest):

    def setUp(self):
        self._authenticate()  # default setup
        self.init_tests()

        self.base = '{}/{}'.format(self.URL_BASE, '/api/v2/weather/watchers')

    def test_weather_service_empty(self):
        # !! route get_watcher !!

        # Safety cleanup
        self.context_cleanup()

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

        # !! route weatherwatchers !!

        # With bad id
        r = self._send(url=self.base + '/scotty',
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.NOT_FOUND.value)
        json = r.json()
        self.assertTrue(isinstance(json, dict))
        self.assertEqual(json['name'], 'resource_not_found')
"""


class TestWeatherAPI(BasicWeatherAPITest):

    """
    NB : Because of api/v2/compute-pbehaviors, tests can be runned once per 10s
    """

    def setUp(self):
        self._authenticate()  # default setup
        self.init_tests()

        self.base = '{}/{}'.format(self.URL_BASE, '/api/v2/weather/watchers')

    def init_tests(self):
        super(TestWeatherAPI, self).init_tests()

        self.context_cleanup()

        # Adding watcher and alarm to watch upon
        self.event_url = '{}/api/v2/event'.format(self.URL_BASE)
        r = self._send(url=self.event_url,
                       method=Method.post,
                       data=dumps(self.event1))
        self.assertEqual(r.status_code, HTTP.OK.value)

        sleep(2)

        watcher_url = '{}/api/v2/watchers'.format(self.URL_BASE)
        r = self._send(url=watcher_url,
                       method=Method.post,
                       data=dumps(self.watcher_1))
        self.assertEqual(r.status_code, HTTP.OK.value)

        r = self._send(url=watcher_url,
                       method=Method.post,
                       data=dumps(self.watcher_2))
        self.assertEqual(r.status_code, HTTP.OK.value)

        # Sending unlinked event 3
        r = self._send(url=self.event_url,
                       method=Method.post,
                       data=dumps(self.event3))
        self.assertEqual(r.status_code, HTTP.OK.value)
        # Adding outside watcher and alarm (should not affect computations)
        r = self._send(url=watcher_url,
                       method=Method.post,
                       data=dumps(self.watcher_3))
        self.assertEqual(r.status_code, HTTP.OK.value)

        self.pbheavior_url = '{}/api/v2/pbehavior'.format(self.URL_BASE)
        self.pbehavior_ids = []

    def tearDown(self):
        """Deleting contextual informations"""
        self.context_cleanup()

        for pbehavior_id in self.pbehavior_ids:
            self._send(url=self.pbheavior_url + '/' + pbehavior_id,
                       method=Method.delete)

    def get_watcher(self, watcher_filter):
        """
        Helper to read watcher informations (the tested interface).
        """
        r = self._send(url=self.base + '/' + watcher_filter,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)

        return r.json()

    def test_weather_all_routes(self):
        """
        Generic scenario:
        To begin, we have 2 watchers and 2 events isolated.
        After verifying base state, we verify that a pbehavior and a new
        event correctly influence watcher state.
        """

        sleep(1)

        # Read all watchers
        r = self._send(url=self.base + '/' + '{}',
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 3)

        # Read base watcher 1
        watcher_filter_1 = dumps({'_id': self.watcher_1['_id']})
        r = self._send(url=self.base + '/' + watcher_filter_1,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 1)
        self.assertEqual(json[0]['state']['val'], 2)
        self.assertFalse(json[0]['hasactivepbehaviorinentities'])
        self.assertFalse(json[0]['hasallactivepbehaviorinentities'])

        # Read specific watcher 1
        r = self._send(url=self.base + '/' + self.watcher_1['_id'])
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 1)
        self.assertEqual(json[0]['state']['val'], 2)
        self.assertIsNone(json[0]['automatic_action_timer'])

        # Adding a pbehavior on event 1
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
        json = self.get_watcher(watcher_filter_1)
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 1)
        self.assertFalse(json[0]['hasactivepbehaviorinentities'])
        self.assertTrue(json[0]['hasallactivepbehaviorinentities'])

        # Verifiyng specific watcher state
        r = self._send(url=self.base + '/' + self.watcher_1['_id'])
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 1)
        self.assertEqual(json[0]['state']['val'], 2)
        self.assertIsNone(json[0]['automatic_action_timer'])
        pbehavior = json[0]['pbehavior']
        self.assertTrue(isinstance(pbehavior, list))
        self.assertTrue('_id' in pbehavior[0])
        self.assertTrue(pbehavior[0]['enabled'])

        # Sending another linked event 2
        r = self._send(url=self.event_url,
                       method=Method.post,
                       data=dumps(self.event2))
        self.assertEqual(r.status_code, HTTP.OK.value)

        sleep(1)

        # Verifying watcher state after event 2
        json = self.get_watcher(watcher_filter_1)
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 1)
        self.assertTrue(json[0]['hasactivepbehaviorinentities'])
        self.assertFalse(json[0]['hasallactivepbehaviorinentities'])

        # Verifying specific watcher state after event 2
        r = self._send(url=self.base + '/' + self.watcher_1['_id'])
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 2)
        states = [json[0]['state']['val'], json[1]['state']['val']]
        self.assertListEqual(states, [2, 3])
