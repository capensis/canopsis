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

import json

from copy import deepcopy
from json import dumps
from time import time, sleep

from canopsis.event import forger
from canopsis.models.watcher import WatcherModel
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
            'type_': 'pause'
        }

        self.pbehavior_watcher1 = {
            'name': 'buck',
            'author': 'mira',
            'filter_': {'_id': 'watcher_first_gen'},
            'rrule': None,
            'tstart': now,
            'tstop': now + 60 * 60,
            'type_': 'pause'
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
        self.assertTrue(isinstance(json, list))
        self.assertEqual(len(json), 0)

        # !! route weatherwatchers !!

        # With bad id
        r = self._send(url=self.base + '/scotty',
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.NOT_FOUND.value)
        self.assertTrue(isinstance(json, dict))
        self.assertEqual(json['name'], 'resource_not_found')
"""

class TestWeatherFilterAPI(BasicWeatherAPITest):

    def init_tests(self):
        self.evt_ok1 = forger(
            event_type='check',
            connector='gutamaya',
            connector_name='imperial',
            source_type='resource',
            component='test_weatherfilter',
            resource='ok1',
            state=0,
            output='cutter'
        )

        self.evt_ok2 = forger(
            event_type='check',
            connector='gutamaya',
            connector_name='imperial',
            source_type='resource',
            component='test_weatherfilter',
            resource='ok2',
            state=0,
            output='cutter'
        )

        self.evt_paused = forger(
            event_type='check',
            connector='gutamaya',
            connector_name='imperial',
            source_type='resource',
            component='test_weatherfilter',
            resource='paused',
            state=0,
            output='cutter'
        )
        self.evt_paused_id = 'paused/test_weatherfilter'

        self.evt_maintenance = forger(
            event_type='check',
            connector='faulcon',
            connector_name='delacy',
            source_type='resource',
            component='test_weatherfilter',
            resource='maintenance',
            state=0,
            output='cobraMKIII'
        )
        self.evt_maintenance_id = 'maintenance/test_weatherfilter'

        now = int(time())
        self.pb_paused = {
            'name': 'gutamaya_imperial_interdictor',
            'author': 'cmdr',
            'filter_': {'_id': self.evt_paused_id},
            'rrule': None,
            'tstart': now,
            'tstop': now + 60 * 60,
            'type_': 'pause'
        }

        self.pb_maintenance = {
            'name': 'cobramkIII',
            'author': 'cmdr',
            'filter_': {'_id': self.evt_maintenance_id},
            'rrule': None,
            'tstart': now,
            'tstop': now + 60 * 60,
            'type_': 'maintenance'
        }

        self.watcher_ok = WatcherModel(
            "weatherfilter_ok",
            {
                "_id": {
                    "$in": [
                        "ok1/test_weatherfilter",
                        "ok2/test_weatherfilter"
                    ]
                }
            },
            "only_ok"
        )

        self.watcher_partial_pause = WatcherModel(
            "weatherfilter_partial_pause",
            {
                "_id": {
                    "$in": [
                        "paused/test_weatherfilter",
                        "ok1/test_weatherfilter"
                    ]
                }
            },
            "partialy_paused"
        )

        self.watcher_full_pause = WatcherModel(
            "weatherfilter_full_pause",
            {
                "_id": {
                    "$in": [
                        "maintenance/test_weatherfilter",
                        "paused/test_weatherfilter"
                    ]
                }
            },
            "paused_maintenance"
        )

        self.watcher_full_pause_maintenance = WatcherModel(
            "watcherfilter_full_pause_maintenance",
            {
                "_id": {
                    "$in": [
                        "maintenance/test_weatherfilter",
                    ]
                }
            },
            "maintenance"
        )

        self.watcher_mixed = WatcherModel(
            "weatherfilter_mixed",
            {
                "_id": {
                    "$in": [
                        "maintenance/test_weatherfilter",
                        "paused/test_weatherfilter",
                        "ok1/test_weatherfilter"
                    ]
                }
            },
            "mixed"
        )

    def setUp(self):
        self._authenticate()
        self.init_tests()
        self.setup_tests()

    def tearDown(self):
        def delete_watcher(watcher_id):
            self._send(
                url='{}/api/v2/watchers/{}'.format(self.URL_BASE, watcher_id),
                method=Method.delete
            )

        def delete_pbehavior(pb_id):
            self._send(
                url='{}/api/v2/pbehavior/{}'.format(self.URL_BASE, pb_id),
                method=Method.delete
            )

        for wid in self.watcher_ids:
            delete_watcher(wid)

        for pb_id in self.pbehavior_ids:
            delete_pbehavior(pb_id)

    def setup_tests(self):
        self.watcher_ids = []
        self.pbehavior_ids = []

        # send events
        def send_event(evt):
            url = '{}/api/v2/event'.format(self.URL_BASE)
            r = self._send(url=url, method=Method.post, data=dumps(evt))
            self.assertEqual(r.status_code, HTTP.OK.value)

        send_event(self.evt_ok1)
        send_event(self.evt_ok2)
        send_event(self.evt_paused)
        send_event(self.evt_maintenance)

        # ensure context entities
        sleep(1)
        url = '{}/context'.format(self.URL_BASE)
        ctx_filter = {
            "impact":{"$in":["test_weatherfilter"]},
            "type":{"$in":["resource", "component"]}
        }
        data={'_filter': ctx_filter}
        r = self._send(url=url, method=Method.post, data=dumps(data))
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertEqual(r.json()['total'], 4)

        # send pbehaviors and force compute
        r = self._send(
            url='{}/api/v2/pbehavior'.format(self.URL_BASE),
            method=Method.post,
            data=dumps(self.pb_maintenance))
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.pbehavior_ids.append(r.text)

        r = self._send(
            url='{}/api/v2/pbehavior'.format(self.URL_BASE),
            method=Method.post,
            data=dumps(self.pb_paused))
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.pbehavior_ids.append(r.text)

        r = self._send(
            url='{}/api/v2/compute-pbehaviors'.format(self.URL_BASE),
            method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertTrue(r.json())

        # send watchers and force compute
        def send_watcher(watcher):
            r = self._send(
                url='{}/api/v2/watchers'.format(self.URL_BASE),
                method=Method.post,
                data=dumps(watcher.to_dict()))
            self.watcher_ids.append(watcher.id_)
            self.assertEqual(r.status_code, HTTP.OK.value)

        send_watcher(self.watcher_ok)
        send_watcher(self.watcher_partial_pause)
        send_watcher(self.watcher_full_pause)
        send_watcher(self.watcher_full_pause_maintenance)
        send_watcher(self.watcher_mixed)

        url = '{}/api/v2/watchers/compute-watchers-links'.format(self.URL_BASE)
        r = self._send(url=url, method=Method.get)
        self.assertTrue(r.json())

    def test_watcher_filter(self):
        def test_watcher_all():
            url = '{}/api/v2/weather/watchers/{}'.format(self.URL_BASE, '{}')
            r = self._send(url)
            watchers = r.json()
            self.assertEqual(len(watchers), 5)

        def test_watcher_ok():
            filter_ = {
                'active_pb_some': False
            }
            url = '{}/api/v2/weather/watchers/{}'.format(self.URL_BASE, json.dumps(filter_))
            r = self._send(url)

            self.assertEqual(r.status_code, HTTP.OK.value)

            watchers = r.json()
            self.assertEqual(len(watchers), 1)
            watcher = watchers[0]

            self.assertFalse(watcher['active_pb_all'])
            self.assertFalse(watcher['active_pb_some'])
            self.assertEqual(watcher['state']['val'], 0)
            self.assertEqual(watcher['display_name'], 'only_ok')
            self.assertEqual(watcher['entity_id'], 'weatherfilter_ok')
            self.assertEqual(len(watcher['pbehavior']), 0)

        def test_watcher_partially_paused():
            filter_ = {
                'active_pb_some': True,
                'active_pb_all': False
            }
            url = '{}/api/v2/weather/watchers/{}'.format(self.URL_BASE, json.dumps(filter_))
            r = self._send(url)
            self.assertEqual(r.status_code, HTTP.OK.value)
            watchers = r.json()

            self.assertEqual(len(watchers), 2)

            self.assertFalse(watchers[0]['active_pb_all'])
            self.assertTrue(watchers[0]['active_pb_some'])
            self.assertEqual(watchers[0]['entity_id'], 'weatherfilter_mixed')
            self.assertEqual(len(watchers[0]['pbehavior']), 2)
            watcher0_pbehaviors = sorted(watchers[0]['pbehavior'])
            self.assertEqual(watcher0_pbehaviors[0]['name'], 'cobramkIII')
            self.assertEqual(watcher0_pbehaviors[0]['type_'], 'maintenance')
            self.assertEqual(watcher0_pbehaviors[1]['name'], 'gutamaya_imperial_interdictor')
            self.assertEqual(watcher0_pbehaviors[1]['type_'], 'pause')

            self.assertFalse(watchers[1]['active_pb_all'])
            self.assertTrue(watchers[1]['active_pb_some'])
            self.assertEqual(watchers[1]['entity_id'], 'weatherfilter_partial_pause')
            self.assertEqual(len(watchers[1]['pbehavior']), 1)
            self.assertEqual(watchers[1]['pbehavior'][0]['name'], 'gutamaya_imperial_interdictor')
            self.assertEqual(watchers[1]['pbehavior'][0]['type_'], 'pause')

        def test_watcher_full_paused():
            filter_ = {
                'active_pb_some': True,
                'active_pb_all': True,
                'active_pb_type': 'pause'
            }
            url = '{}/api/v2/weather/watchers/{}'.format(self.URL_BASE, json.dumps(filter_))
            r = self._send(url)
            self.assertEqual(r.status_code, HTTP.OK.value)

            watchers = r.json()

            self.assertEqual(len(watchers), 1)
            self.assertEqual(watchers[0]['entity_id'], 'weatherfilter_full_pause')
            self.assertTrue(watchers[0]['active_pb_all'])
            self.assertTrue(watchers[0]['active_pb_some'])
            self.assertEqual(len(watchers[0]['pbehavior']), 2)

            pbehaviors = sorted(watchers[0]['pbehavior'])

            self.assertEqual(pbehaviors[0]['name'], 'cobramkIII')
            self.assertEqual(pbehaviors[0]['type_'], 'maintenance')
            self.assertEqual(pbehaviors[1]['name'], 'gutamaya_imperial_interdictor')
            self.assertEqual(pbehaviors[1]['type_'], 'pause')

        def test_watcher_full_paused_maintenance():
            filter1 = {
                'active_pb_some': True,
                'active_pb_type': 'pause'
            }
            url = '{}/api/v2/weather/watchers/{}'.format(self.URL_BASE, json.dumps(filter1))
            r = self._send(url)
            self.assertEqual(r.status_code, HTTP.OK.value)

            watchers = sorted(r.json())

            self.assertEqual(watchers[0]['entity_id'], 'weatherfilter_mixed')
            self.assertEqual(watchers[1]['entity_id'], 'weatherfilter_partial_pause')
            self.assertEqual(watchers[2]['entity_id'], 'weatherfilter_full_pause')

            self.assertEqual(len(watchers), 3)

            pbehaviors = sorted(watchers[0]['pbehavior'])
            self.assertEqual(len(pbehaviors), 2)
            self.assertEqual(pbehaviors[0]['name'], 'cobramkIII')
            self.assertEqual(pbehaviors[0]['type_'], 'maintenance')
            self.assertEqual(pbehaviors[1]['name'], 'gutamaya_imperial_interdictor')
            self.assertEqual(pbehaviors[1]['type_'], 'pause')

            pbehaviors = sorted(watchers[1]['pbehavior'])
            self.assertEqual(len(pbehaviors), 1)
            self.assertEqual(pbehaviors[0]['name'], 'gutamaya_imperial_interdictor')
            self.assertEqual(pbehaviors[0]['type_'], 'pause')

            pbehaviors = sorted(watchers[2]['pbehavior'])
            self.assertEqual(len(pbehaviors), 2)
            self.assertEqual(pbehaviors[0]['name'], 'cobramkIII')
            self.assertEqual(pbehaviors[0]['type_'], 'maintenance')
            self.assertEqual(pbehaviors[1]['name'], 'gutamaya_imperial_interdictor')
            self.assertEqual(pbehaviors[1]['type_'], 'pause')

            # return two watchers, also one with a "pause" because this watchers also
            # have a maintenance pbehavior.
            filter2 = {
                'active_pb_some': True,
                'active_pb_all': True,
                'active_pb_type': 'maintenance'
            }
            url = '{}/api/v2/weather/watchers/{}'.format(self.URL_BASE, json.dumps(filter2))
            r = self._send(url)
            self.assertEqual(r.status_code, HTTP.OK.value)

            watchers = sorted(r.json())
            self.assertEqual(len(watchers), 2)

            self.assertEqual(watchers[0]['entity_id'], 'watcherfilter_full_pause_maintenance')
            pbehaviors = sorted(watchers[0]['pbehavior'])
            self.assertEqual(len(pbehaviors), 1)
            self.assertEqual(pbehaviors[0]['name'], 'cobramkIII')
            self.assertEqual(pbehaviors[0]['type_'], 'maintenance')

            self.assertEqual(watchers[1]['entity_id'], 'weatherfilter_full_pause')
            pbehaviors = sorted(watchers[1]['pbehavior'])
            self.assertEqual(len(pbehaviors), 2)
            self.assertEqual(pbehaviors[0]['name'], 'cobramkIII')
            self.assertEqual(pbehaviors[0]['type_'], 'maintenance')
            self.assertEqual(pbehaviors[1]['name'], 'gutamaya_imperial_interdictor')
            self.assertEqual(pbehaviors[1]['type_'], 'pause')

            # return no watcher as the only allowed pbehavior type is "YouDied" and we don't have any.
            filter2 = {
                'active_pb_some': True,
                'active_pb_type': 'YouDied'
            }
            url = '{}/api/v2/weather/watchers/{}'.format(self.URL_BASE, json.dumps(filter2))
            r = self._send(url)
            self.assertEqual(r.status_code, HTTP.OK.value)

            watchers = r.json()
            self.assertEqual(len(watchers), 0)

            # will return watchers without pbehavior
            filter3 = {
                'active_pb_type': 'YouDied'
            }
            url = '{}/api/v2/weather/watchers/{}'.format(self.URL_BASE, json.dumps(filter3))
            r = self._send(url)
            self.assertEqual(r.status_code, HTTP.OK.value)

            watchers = r.json()
            self.assertEqual(len(watchers), 1)
            self.assertEqual(len(watchers[0]['pbehavior']), 0)

        def test_watcher_mixed():
            filter_ = {
                'active_pb_some': True,
                'active_pb_all': False
            }
            url = '{}/api/v2/weather/watchers/{}'.format(self.URL_BASE, json.dumps(filter_))
            r = self._send(url)
            self.assertEqual(r.status_code, HTTP.OK.value)

            watchers = sorted(r.json())

            self.assertEqual(len(watchers), 2)

            self.assertEqual(watchers[0]['entity_id'], 'weatherfilter_mixed')
            self.assertEqual(len(watchers[0]['pbehavior']), 2)

            pbehaviors = sorted(watchers[0]['pbehavior'])
            self.assertEqual(pbehaviors[0]['name'], 'cobramkIII')
            self.assertEqual(pbehaviors[0]['type_'], 'maintenance')
            self.assertEqual(pbehaviors[1]['name'], 'gutamaya_imperial_interdictor')
            self.assertEqual(pbehaviors[1]['type_'], 'pause')

            self.assertEqual(watchers[1]['entity_id'], 'weatherfilter_partial_pause')
            self.assertEqual(len(watchers[1]['pbehavior']), 1)

            pbehaviors = sorted(watchers[1]['pbehavior'])
            self.assertEqual(pbehaviors[0]['name'], 'gutamaya_imperial_interdictor')
            self.assertEqual(pbehaviors[0]['type_'], 'pause')

        test_watcher_all()
        test_watcher_ok()
        test_watcher_partially_paused()
        test_watcher_full_paused()
        test_watcher_full_paused_maintenance()
        test_watcher_mixed()

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

        self.pbehavior_url = '{}/api/v2/pbehavior'.format(self.URL_BASE)
        self.pbehavior_ids = []

    def tearDown(self):
        """Deleting contextual informations"""
        self.context_cleanup()

        for pbehavior_id in self.pbehavior_ids:
            self._send(url=self.pbehavior_url + '/' + pbehavior_id,
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
        self.assertTrue(isinstance(r.json(), list))
        self.assertEqual(len(r.json()), 3)

        # Read base watcher 1
        watcher_filter_1 = dumps({'_id': self.watcher_1['_id']})
        r = self._send(url=self.base + '/' + watcher_filter_1,
                       method=Method.get)
        watchers = sorted(r.json())
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertTrue(isinstance(watchers, list))
        self.assertEqual(len(watchers), 1)
        self.assertEqual(watchers[0]['state']['val'], 2)
        self.assertFalse(watchers[0]['active_pb_some'])
        self.assertFalse(watchers[0]['active_pb_all'])

        # Read specific watcher 1
        r = self._send(url=self.base + '/' + self.watcher_1['_id'])
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertTrue(isinstance(r.json(), list))
        self.assertEqual(len(r.json()), 1)
        self.assertEqual(r.json()[0]['state']['val'], 2)
        self.assertIsNone(r.json()[0]['automatic_action_timer'])

        # Adding a pbehavior on event 1
        r = self._send(url=self.pbehavior_url,
                       method=Method.post,
                       data=dumps(self.pbehavior1))
        self.pbehavior_ids.append(r.text)
        self.assertEqual(r.status_code, HTTP.OK.value)

        # Force compute on pbehaviors
        pbehavior_url = '{}/api/v2/compute-pbehaviors'.format(self.URL_BASE)
        r = self._send(url=pbehavior_url,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertTrue(r.json())

        # Verifying watcher state
        r = self.get_watcher(watcher_filter_1)
        self.assertTrue(isinstance(r, list))
        self.assertEqual(len(r), 1)
        self.assertTrue(r[0]['active_pb_some'])
        self.assertTrue(r[0]['active_pb_all'])

        # Verifiyng specific watcher state
        r = self._send(url=self.base + '/' + self.watcher_1['_id'])
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertTrue(isinstance(r.json(), list))
        self.assertEqual(len(r.json()), 1)
        self.assertEqual(r.json()[0]['state']['val'], 2)
        self.assertIsNone(r.json()[0]['automatic_action_timer'])
        pbehavior = r.json()[0]['pbehavior']
        self.assertTrue(isinstance(pbehavior, list))
        self.assertTrue('_id' in pbehavior[0])
        self.assertTrue(pbehavior[0]['enabled'])
        self.assertTrue('type_' in pbehavior[0])
        self.assertTrue('reason' in pbehavior[0])
        self.assertEqual(pbehavior[0]['type_'], self.pbehavior1['type_'])

        # Sending another linked event 2
        r = self._send(url=self.event_url,
                       method=Method.post,
                       data=dumps(self.event2))
        self.assertEqual(r.status_code, HTTP.OK.value)

        sleep(1)

        # Verifying watcher state after event 2
        r = self.get_watcher(watcher_filter_1)
        self.assertTrue(isinstance(r, list))
        self.assertEqual(len(r), 1)
        self.assertTrue(r[0]['active_pb_some'])
        self.assertFalse(r[0]['active_pb_all'])

        # Verifying specific watcher state after event 2
        r = self._send(url=self.base + '/' + self.watcher_1['_id'])
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertTrue(isinstance(r.json(), list))
        self.assertEqual(len(r.json()), 2)
        states = [r.json()[0]['state']['val'], r.json()[1]['state']['val']]
        states.sort()
        self.assertListEqual(states, [2, 3])

        # Verifying watcher state on a wrong pbh type
        pbtyped_wfilter_1 = dumps(
            {
                '_id': self.watcher_1['_id'],
                'active_pb_type': 'wrong_type',
            },
        )
        r = self.get_watcher("{}".format(pbtyped_wfilter_1))
        self.assertTrue(isinstance(r, list))
        self.assertEqual(len(r), 0)

        # Adding a pbehavior on watcher 1
        r = self._send(url=self.pbehavior_url,
                       method=Method.post,
                       data=dumps(self.pbehavior_watcher1))
        self.pbehavior_ids.append(r.text)
        self.assertEqual(r.status_code, HTTP.OK.value)

        # Force compute on pbehaviors
        pbehavior_url = '{}/api/v2/compute-pbehaviors'.format(self.URL_BASE)
        r = self._send(url=pbehavior_url,
                       method=Method.get)
        self.assertEqual(r.status_code, HTTP.OK.value)
        self.assertTrue(r.json())

        # check if watcher1 gets back it's pbehavior
        filter_ = json.dumps({'_id': self.watcher_1['_id']})
        watchers = self.get_watcher(filter_)
        self.assertEqual(len(watchers), 1)
        self.assertEqual(watchers[0]['entity_id'], 'watcher_first_gen')
        pbehaviors = sorted(watchers[0]['watcher_pbehavior'])
        self.assertEqual(len(pbehaviors), 1)
        self.assertEqual(pbehaviors[0]['author'], 'mira')