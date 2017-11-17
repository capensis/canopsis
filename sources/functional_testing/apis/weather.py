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

from canopsis.context_graph.manager import ContextGraph
from canopsis.event import forger
from test_base import BaseApiTest, Method, HTTP

event1 = forger(
    event_type='check',
    connector='monconnectorlol',
    connector_name='tonconnectorlol',
    source_type='resource',
    component='ceciestuncomposant',
    resource='ceciestuneresource',
    state=2,
    output='onsenfout'
)

event1_id = ContextGraph.get_id(event1)  # TODO foutre ça en API

# Sample watcher (to insert)
watcher_dict = {
    "description": "a_description",
    "display_name": "a_displayed_name",
    "enable": True,
    "mfilter": dumps({'_id': event1_id}),
    "_id": "watcher_id"
}
watcher_id = watcher_dict['_id']

pbehavior1 = {
    'name': 'imagine',
    'author': 'lennon',
    'filter_': dumps({'_id': event1_id}),
    'rrule': None,
    'tstart': 0,
    'tstop': 1,
}


class TestWeatherAPI(BaseApiTest):

    def setUp(self):
        self._authenticate()  # default setup

        self.base = '{}/{}'.format(self.URL_BASE, '/api/v2/weather/watchers')

    def init_watchers(self):
        """Adding watcher and alarms to watch upon"""
        watcher_url = '{}/api/v2/watchers'.format(self.URL_BASE)
        r = self._send(url=watcher_url,
                       method=Method.post,
                       data=dumps(watcher_dict))
        self.assertEqual(r.status_code, HTTP.OK.value)

        event_url = '{}/api/v2/event'.format(self.URL_BASE)
        r = self._send(url=event_url,
                       method=Method.post,
                       data=dumps(event1))
        self.assertEqual(r.status_code, HTTP.OK.value)

        # TODO: webservice pour forcer le recalcul des watchers

    def delete_watchers(self):
        """Deleting watcher"""
        watcher_url = '{}/api/v2/watchers/{}'.format(self.URL_BASE, watcher_id)
        self._send(url=watcher_url,
                   method=Method.delete)

    def test_weather_get_watcher_empty(self):
        self.delete_watchers()

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

    def test_weather_get_watcher(self):

        self.delete_watchers()
        self.init_watchers()

        watcher_filter = dumps({'_id': watcher_id})
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
        pbheavior_url = '{}/api/v2/pbehavior'.format(self.URL_BASE)
        r = self._send(url=pbheavior_url,
                       method=Method.post,
                       data=dumps(pbehavior1))
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertEqual(len(json), 1)
        print(json)

        pbheavior_url = '{}/api/v2/compute-pbehaviors'.format(self.URL_BASE)
        r = self._send(url=pbheavior_url,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        print(r.json())
        self.assertTrue(r.json())

        # Verifying watcher state
        r = self._send(url=self.base + '/' + watcher_filter,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        json = r.json()
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 1)
        self.assertFalse(json['hasactivepbehaviorinentities'])
        self.assertFalse(json['hasallactivepbehaviorinentities'])

        self.delete_watchers()

    #def test_weather_weatherwatchers(self):
