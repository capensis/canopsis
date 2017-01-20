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

"""
Functional tests for webservice requests on alerts/get-alarms
"""

from unittest import TestCase, main
from requests import Session
from json import dumps

from canopsis.middleware.core import Middleware


BASE_URL = 'http://localhost:8082'
USERNAME = 'root'
PASSWORD = 'root'


class WebserviceTest(TestCase):
    def setUp(self):
        self.f = Feeder()
        self.f.feed()

        auth_url = '{}/auth'.format(BASE_URL)
        self.s = Session()
        self.s.post(
            auth_url,
            params={'username': USERNAME, 'password': PASSWORD}
        )

        self.alarm_url = '{}/alerts/get-alarms'.format(BASE_URL)

    def tearDown(self):
        self.f.clean()

    def test_01(self):
        expected_order = [
            '/resource/ctr1/prefix_name1/c1/r7',
            '/resource/ctr0/name_suffix0/c1/r8',
            '/resource/ctr1/prefix_name0/c1/r6',
            '/resource/ctr0/name_suffix0/c0/r4',
            '/resource/ctr0/name_suffix1/c1/r9',
            '/resource/ctr0/name_suffix1/c1/r5'
        ]

        r = self.s.get(self.alarm_url).json()

        self.assertIs(r['success'], True)

        data = r['data'][0]

        self.assertEqual(data['first'], 1)
        self.assertEqual(data['last'], 6)
        self.assertEqual(data['total'], 6)
        self.assertIs(data['truncated'], False)

        for i, expected_eid in enumerate(expected_order):
            self.assertEqual(data['alarms'][i]['d'], expected_eid)

    def test_02(self):
        expected_order = [
            '/resource/ctr0/name_suffix1/c1/r9',
            '/resource/ctr0/name_suffix0/c1/r8',
            '/resource/ctr1/prefix_name1/c1/r7',
            '/resource/ctr1/prefix_name0/c1/r6',
            '/resource/ctr0/name_suffix1/c1/r5',
            '/resource/ctr0/name_suffix0/c0/r4',
            '/resource/ctr0/name_suffix1/c0/r1'
        ]

        r = self.s.get(
            self.alarm_url,
            params={'tstart': 3000, 'resolved': 'true', 'sort_key': 'resource'}
        ).json()

        self.assertIs(r['success'], True)

        data = r['data'][0]

        self.assertEqual(data['first'], 1)
        self.assertEqual(data['last'], 7)
        self.assertEqual(data['total'], 7)
        self.assertIs(data['truncated'], False)

        for i, expected_eid in enumerate(expected_order):
            self.assertEqual(data['alarms'][i]['d'], expected_eid)

    def test_03(self):
        expected_order = [
            '/resource/ctr1/prefix_name0/c1/r6',
            '/resource/ctr0/name_suffix0/c0/r4',
            '/resource/ctr0/name_suffix1/c1/r9',
            '/resource/ctr0/name_suffix1/c1/r5'
        ]

        r = self.s.get(
            self.alarm_url,
            params={'tstart': 1200, 'tstop': 1800, 'opened': 'true'}
        ).json()

        self.assertIs(r['success'], True)

        data = r['data'][0]

        self.assertEqual(data['first'], 1)
        self.assertEqual(data['last'], 4)
        self.assertEqual(data['total'], 4)
        self.assertIs(data['truncated'], False)

        for i, expected_eid in enumerate(expected_order):
            self.assertEqual(data['alarms'][i]['d'], expected_eid)

    def test_04(self):
        expected_order = [
            '/resource/ctr0/name_suffix0/c0/r4',
            '/resource/ctr0/name_suffix1/c1/r9',
            '/resource/ctr0/name_suffix1/c1/r5'
        ]

        r = self.s.get(
            self.alarm_url,
            params={'tstart': 0, 'tstop': 1500, 'limit': 3}
        ).json()

        self.assertIs(r['success'], True)

        data = r['data'][0]

        self.assertEqual(data['first'], 1)
        self.assertEqual(data['last'], 3)
        self.assertEqual(data['total'], 3)
        self.assertIs(data['truncated'], False)

        for i, expected_eid in enumerate(expected_order):
            self.assertEqual(data['alarms'][i]['d'], expected_eid)

    def test_05(self):
        expected_order = [
            '/resource/ctr0/name_suffix1/c0/r1',
            '/resource/ctr1/prefix_name0/c0/r2',
            '/resource/ctr1/prefix_name1/c0/r3',
            '/resource/ctr0/name_suffix0/c0/r0'
        ]

        r = self.s.get(
            self.alarm_url,
            params={'opened': 'false', 'resolved': 'true'}
        ).json()

        self.assertIs(r['success'], True)

        data = r['data'][0]

        self.assertEqual(data['first'], 1)
        self.assertEqual(data['last'], 4)
        self.assertEqual(data['total'], 4)
        self.assertIs(data['truncated'], False)

        for i, expected_eid in enumerate(expected_order):
            self.assertEqual(data['alarms'][i]['d'], expected_eid)

    def test_06(self):
        expected_order = []

        r = self.s.get(
            self.alarm_url,
            params={'filter': dumps({"domain": "d0"})}
        ).json()

        self.assertIs(r['success'], True)

        data = r['data'][0]

        self.assertEqual(data['first'], 0)
        self.assertEqual(data['last'], 0)
        self.assertEqual(data['total'], 0)
        self.assertIs(data['truncated'], False)

        for i, expected_eid in enumerate(expected_order):
            self.assertEqual(data['alarms'][i]['d'], expected_eid)

    def test_07(self):
        expected_order = [
            '/resource/ctr0/name_suffix1/c1/r5',
            '/resource/ctr0/name_suffix0/c0/r4',
            '/resource/ctr1/prefix_name0/c1/r6',
            '/resource/ctr1/prefix_name1/c1/r7'
        ]

        filter_ = dumps(
            {
                '$or': [
                    {'domain': 'd1'},
                    {'connector': 'ctr1'}
                ]
            }
        )

        r = self.s.get(
            self.alarm_url,
            params={'filter': filter_, 'sort_dir': 'ASC'}
        ).json()

        self.assertIs(r['success'], True)

        data = r['data'][0]

        self.assertEqual(data['first'], 1)
        self.assertEqual(data['last'], 4)
        self.assertEqual(data['total'], 4)
        self.assertIs(data['truncated'], False)

        for i, expected_eid in enumerate(expected_order):
            self.assertEqual(data['alarms'][i]['d'], expected_eid)

    def test_08(self):
        expected_order = [
            '/resource/ctr0/name_suffix1/c1/r5',
            '/resource/ctr0/name_suffix0/c0/r4',
            '/resource/ctr1/prefix_name0/c1/r6',
            '/resource/ctr1/prefix_name1/c1/r7'
        ]

        filter_ = dumps(
            {
                '$or': [
                    {'domain': 'd1'},
                    {'connector': 'ctr1'}
                ]
            }
        )

        r = self.s.get(
            self.alarm_url,
            params={'filter': filter_, 'sort_dir': 'ASC', 'limit': 4}
        ).json()

        self.assertIs(r['success'], True)

        data = r['data'][0]

        self.assertEqual(data['first'], 1)
        self.assertEqual(data['last'], 4)
        self.assertEqual(data['total'], 4)
        self.assertIs(data['truncated'], False)

        for i, expected_eid in enumerate(expected_order):
            self.assertEqual(data['alarms'][i]['d'], expected_eid)

    def test_09(self):
        expected_order = []

        filter_ = dumps(
            {
                '$or': [
                    {'domain': 'd1'},
                    {'connector': 'ctr1'}
                ]
            }
        )

        search = "ALL perimeter = 'p0' AND connector_name = 'prefix_name1'"

        r = self.s.get(
            self.alarm_url,
            params={'filter': filter_, 'search': search}
        ).json()

        self.assertIs(r['success'], True)

        data = r['data'][0]

        self.assertEqual(data['first'], 0)
        self.assertEqual(data['last'], 0)
        self.assertEqual(data['total'], 0)
        self.assertIs(data['truncated'], False)

        for i, expected_eid in enumerate(expected_order):
            self.assertEqual(data['alarms'][i]['d'], expected_eid)


FEED_CONF = [
    {
        "connector": "ctr0",
        "connector_name": "name_suffix0",
        "component": "c0",
        "resource": "r0",
        "domain": "d0",
        "perimeter": "p0",
        "state": 1,
        "opened": 0,
        "resolved": 1200
    },
    {
        "connector": "ctr0",
        "connector_name": "name_suffix1",
        "component": "c0",
        "resource": "r1",
        "domain": "d0",
        "perimeter": "p1",
        "state": 1,
        "opened": 1800,
        "resolved": 3000
    },
    {
        "connector": "ctr1",
        "connector_name": "prefix_name0",
        "component": "c0",
        "resource": "r2",
        "domain": "d0",
        "perimeter": "p2",
        "state": 1,
        "opened": 300,
        "resolved": 1500
    },
    {
        "connector": "ctr1",
        "connector_name": "prefix_name1",
        "component": "c0",
        "resource": "r3",
        "domain": "d0",
        "perimeter": "p2",
        "state": 1,
        "opened": 0,
        "resolved": 1200
    },
    {
        "connector": "ctr0",
        "connector_name": "name_suffix0",
        "component": "c0",
        "resource": "r4",
        "domain": "d1",
        "perimeter": "p1",
        "state": 1,
        "opened": 1500
    },
    {
        "connector": "ctr0",
        "connector_name": "name_suffix1",
        "component": "c1",
        "resource": "r5",
        "domain": "d1",
        "perimeter": "p1",
        "state": 1,
        "opened": 0
    },
    {
        "connector": "ctr1",
        "connector_name": "prefix_name0",
        "component": "c1",
        "resource": "r6",
        "domain": "d2",
        "perimeter": "p0",
        "state": 1,
        "opened": 1800
    },
    {
        "connector": "ctr1",
        "connector_name": "prefix_name1",
        "component": "c1",
        "resource": "r7",
        "domain": "d2",
        "perimeter": "p1",
        "state": 1,
        "opened": 3000
    },
    {
        "connector": "ctr0",
        "connector_name": "name_suffix0",
        "component": "c1",
        "resource": "r8",
        "domain": "d2",
        "perimeter": "p2",
        "state": 1,
        "opened": 2400
    },
    {
        "connector": "ctr0",
        "connector_name": "name_suffix1",
        "component": "c1",
        "resource": "r9",
        "domain": "d2",
        "perimeter": "p2",
        "state": 1,
        "opened": 900
    }
]


class Feeder(object):
    def __init__(self):
        self.storage = Middleware.get_middleware_by_uri(
            'mongodb-periodical-alarm://')

        self.conf = FEED_CONF

    def feed(self):
        for alarm in self.get_records():
            self.storage._backend.insert(alarm)

    def clean(self):
        for eid in self.get_records_eids():
            self.storage._backend.remove({'d': eid})

    def get_records(self):
        for a in self.conf:
            alarm = {
                'v': {
                    'connector': a['connector'],
                    'connector_name': a['connector_name'],
                    'component': a['component'],
                    'resource': a['resource'],
                    'extra': {
                        'domain': a['domain'],
                        'perimeter': a['perimeter']
                    },
                    # Prevent stats to be computed for those alarms
                    'tags': [
                        'stats-opened',
                        'stats-resolved'
                    ],
                    # Useless for this test
                    'steps': {},
                    'ack': None,
                    'canceled': None,
                    'snooze': None,
                    'ticket': None,
                    'hard_limit': None
                },
                'd': '/resource/{}/{}/{}/{}'.format(
                    a['connector'],
                    a['connector_name'],
                    a['component'],
                    a['resource']
                ),
                't': a['opened']
            }

            if 'resolved' in a:
                alarm['v']['status'] = {
                    'a': '{}.{}'.format(a['connector'], a['connector_name']),
                    '_t': 'statusdec',
                    'm': '',
                    't': a['resolved'],
                    'val': 0
                }

                alarm['v']['state'] = {
                    'a': '{}.{}'.format(a['connector'], a['connector_name']),
                    '_t': 'statedec',
                    'm': '',
                    't': a['resolved'],
                    'val': 0
                }

                alarm['v']['resolved'] = a['resolved']

            else:
                alarm['v']['status'] = {
                    'a': '{}.{}'.format(a['connector'], a['connector_name']),
                    '_t': 'statusinc',
                    'm': '',
                    't': a['opened'],
                    'val': a['state']
                }

                alarm['v']['state'] = {
                    'a': '{}.{}'.format(a['connector'], a['connector_name']),
                    '_t': 'stateinc',
                    'm': '',
                    't': a['opened'],
                    'val': a['state']
                }

                alarm['v']['resolved'] = None

            yield alarm

    def get_records_eids(self):
        for a in self.conf:
            eid = '/resource/{}/{}/{}/{}'.format(
                a['connector'],
                a['connector_name'],
                a['component'],
                a['resource']
            )

            yield eid


if __name__ == '__main__':
    main()
