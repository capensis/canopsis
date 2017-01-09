#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from unittest import main
from time import sleep

from canopsis.alerts.manager import Alerts
from canopsis.alerts.reader import AlertsReader

from base import BaseTest


class TestReader(BaseTest):
    def setUp(self):
        super(TestReader, self).setUp()

        self.reader = AlertsReader()
        self.reader[AlertsReader.ALARM_STORAGE] = \
            self.manager[Alerts.ALARM_STORAGE]

        self.reader._alarm_fields = {
            'properties': {
                'connector': {'stored_name': 'v.ctr'},
                'component': {'stored_name': 'v.cpt'},
                'entity_id': {'stored_name': 'd'}
            }
        }

    def test__translate_key(self):
        cases = [
            {
                'key': 'untranslated_key',
                'tkey': 'untranslated_key'
            },
            {
                'key': 'connector',
                'tkey': 'v.ctr'
            },
            {
                'key': 'entity_id',
                'tkey': 'd'
            }
        ]

        for case in cases:
            tkey = self.reader._translate_key(case['key'])
            self.assertEqual(tkey, case['tkey'])

    def test__translate_filter(self):
        cases = [
            {
                'filter': {},
                'tfilter': {}
            },
            {
                'filter': {'connector': 'c'},
                'tfilter': {'v.ctr': 'c'}
            },
            {
                'filter': {'$or': [{'connector': 'c1'}, {'component': 'c2'}]},
                'tfilter': {'$or': [{'v.ctr': 'c1'}, {'v.cpt': 'c2'}]}
            },
            {
                'filter': {
                    '$or': [
                        {'entity_id': {'$gte': 12}, 'untranslated': 'val'},
                        {'connector': 'c1'},
                        {'$or': [{'component': 'c2'}, {'untranslated': 'val'}]}
                    ]
                },
                'tfilter': {
                    '$or': [
                        {'d': {'$gte': 12}, 'untranslated': 'val'},
                        {'v.ctr': 'c1'},
                        {'$or': [{'v.cpt': 'c2'}, {'untranslated': 'val'}]}
                    ]
                }
            }
        ]

        for case in cases:
            tfilter = self.reader._translate_filter(case['filter'])
            self.assertEqual(tfilter, case['tfilter'])

    def test__get_time_filter(self):
        # opened=False, resolved=False
        self.assertIs(
            self.reader._get_time_filter(
                opened=False, resolved=False, tstart=0, tstop=0),
            None
        )

        # opened=True, resolved=False
        expected_opened = {'v.resolved': None, 't': {'$lte': 2}}
        self.assertEqual(
            self.reader._get_time_filter(
                opened=True, resolved=False, tstart=1, tstop=2),
            expected_opened
        )

        # opened=False, resolved=True
        expected_resolved = {
            'v.resolved': {'$ne': None},
            '$or': [
                {'t': {'$gte': 1, '$lte': 2}},
                {'v.resolved': {'$gte': 1, '$lte': 2}},
                {'t': {'$lte': 1}, 'v.resolved': {'$gte': 2}}
            ]
        }
        self.assertEqual(
            self.reader._get_time_filter(
                opened=False, resolved=True, tstart=1, tstop=2),
            expected_resolved
        )

        # opened=True, resolved=True
        expected_both = {'$or': [expected_opened, expected_resolved]}
        self.assertEqual(
            self.reader._get_time_filter(
                opened=True, resolved=True, tstart=1, tstop=2),
            expected_both
        )

        # opened=True, resolved=True, tstart=tstop=None
        self.assertEqual(
            self.reader._get_time_filter(
                opened=True, resolved=True,
                tstart=None, tstop=None
            ),
            {}
        )

    def test__get_opened_time_filter(self):
        cases = [
            {
                'tstart': None,
                'tstop': None,
                'expected': {'v.resolved': None}
            },
            {
                'tstart': None,
                'tstop': 0,
                'expected': {'v.resolved': None, 't': {'$lte': 0}}
            },
            {
                'tstart': None,
                'tstop': 42,
                'expected': {'v.resolved': None, 't': {'$lte': 42}}
            },
            {
                'tstart': 13,
                'tstop': None,
                'expected': {'v.resolved': None, 't': {'$lte': 13}}
            },
            {
                'tstart': 13,
                'tstop': 42,
                'expected': {'v.resolved': None, 't': {'$lte': 42}}
            }
        ]

        for case in cases:
            time_filter = self.reader._get_opened_time_filter(
                case['tstart'],
                case['tstop']
            )
            self.assertEqual(time_filter, case['expected'])

    def test__get_resolved_time_filter(self):
        cases = [
            {
                'tstart': None,
                'tstop': None,
                'expected': {'v.resolved': {'$ne': None}}
            },
            {
                'tstart': 13,
                'tstop': None,
                'expected': {
                    'v.resolved': {'$ne': None, '$gte': 13}
                }
            },
            {
                'tstart': None,
                'tstop': 42,
                'expected': {
                    'v.resolved': {'$ne': None},
                    't': {'$lte': 42}
                }
            },
            {
                'tstart': 0,
                'tstop': 0,
                'expected': {
                    'v.resolved': {'$ne': None},
                    '$or': [
                        {'t': {'$gte': 0, '$lte': 0}},
                        {'v.resolved': {'$gte': 0, '$lte': 0}},
                        {'t': {'$lte': 0}, 'v.resolved': {'$gte': 0}}
                    ]
                }
            },
            {
                'tstart': 1,
                'tstop': 2,
                'expected': {
                    'v.resolved': {'$ne': None},
                    '$or': [
                        {'t': {'$gte': 1, '$lte': 2}},
                        {'v.resolved': {'$gte': 1, '$lte': 2}},
                        {'t': {'$lte': 1}, 'v.resolved': {'$gte': 2}}
                    ]
                }
            }
        ]

        for case in cases:
            time_filter = self.reader._get_resolved_time_filter(
                case['tstart'],
                case['tstop']
            )
            self.assertEqual(time_filter, case['expected'])

    def test__translate_sort(self):
        cases = [
            {
                'sort_key': 'untranslated',
                'sort_dir': 'DESC',
                'tkey': 'untranslated',
                'tdir': -1
            },
            {
                'sort_key': 'untranslated',
                'sort_dir': 'ASC',
                'tkey': 'untranslated',
                'tdir': 1
            },
            {
                'sort_key': 'component',
                'sort_dir': 'DESC',
                'tkey': 'v.cpt',
                'tdir': -1
            }
        ]

        for case in cases:
            tkey, tdir = self.reader._translate_sort(
                case['sort_key'],
                case['sort_dir']
            )

            self.assertEqual(tkey, case['tkey'])
            self.assertEqual(tdir, case['tdir'])

    def test_clean_fast_count_cache(self):
        class MockQuery(object):
            def __init__(self, count):
                self.c = count

            def count(self):
                return self.c

            def limit(self, _):
                return self

        self.reader.cache_config = {'expiration': 1}

        query = MockQuery(count=11)
        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=1, tstop=1,
            opened=True, resolved=True,
            filter_=1, search=1
        )

        self.assertEqual(count, 11)
        self.assertEqual(truncated, False)

        query = MockQuery(count=12)
        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=1, tstop=1,
            opened=True, resolved=True,
            filter_=1, search=1
        )

        self.assertEqual(count, 11)
        self.assertEqual(truncated, False)

        sleep(1)
        self.reader.clean_fast_count_cache()

        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=1, tstop=1,
            opened=True, resolved=True,
            filter_=1, search=1
        )

        self.assertEqual(count, 12)
        self.assertEqual(truncated, False)

    def test__get_fast_count(self):
        class MockQuery(object):
            def __init__(self, count, limit=None):
                self.c = count
                self.l = limit

            def count(self):
                if self.l is None:
                    return self.c

                return self.c if self.c <= self.l else self.l

            def limit(self, limit):
                self.l = limit
                return self

        # 1) resolved alarms, no cache, not truncated
        query = MockQuery(count=1)
        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=1, tstop=1,
            opened=True, resolved=True,
            filter_=1, search=1
        )

        self.assertEqual(count, 1)
        self.assertEqual(truncated, False)

        # 2) resolved alarms, no cache, truncated
        query = MockQuery(count=200002)

        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=2, tstop=2,
            opened=True, resolved=True,
            filter_=2, search=2
        )

        self.assertEqual(count, 200000)
        self.assertEqual(truncated, True)

        # 3) resolved alarms, cache, not truncated
        query = MockQuery(count=3)

        # caching value
        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=3, tstop=3,
            opened=True, resolved=True,
            filter_=3, search=3
        )

        self.assertEqual(count, 3)
        self.assertEqual(truncated, False)

        query = MockQuery(count='random value')

        # retrieving cached value with same cache_key but not same query
        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=3, tstop=3,
            opened=True, resolved=True,
            filter_=3, search=3
        )

        self.assertEqual(count, 3)
        self.assertEqual(truncated, False)

        # 4) resolved alarms, cache, truncated
        query = MockQuery(count=200004)

        # caching value
        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=4, tstop=4,
            opened=True, resolved=True,
            filter_=4, search=4
        )

        self.assertEqual(count, 200000)
        self.assertEqual(truncated, True)

        query = MockQuery(count='random value')

        # retrieving cached value with same cache_key but not same query
        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=4, tstop=4,
            opened=True, resolved=True,
            filter_=4, search=4
        )

        self.assertEqual(count, 200000)
        self.assertEqual(truncated, True)

        # 5) opened alarms, not truncated
        query = MockQuery(count=5)

        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=5, tstop=5,
            opened=True, resolved=False,
            filter_=5, search=5
        )

        self.assertEqual(count, 5)
        self.assertEqual(truncated, False)

        # 6) opened alarms, truncated
        query = MockQuery(count=1006)

        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=6, tstop=6,
            opened=True, resolved=False,
            filter_=6, search=6
        )

        self.assertEqual(count, 1000)
        self.assertEqual(truncated, True)

        # 7) change configuration
        self.reader.cache_config = {
            'resolved_truncate': False,
            'opened_truncate': False
        }

        query = MockQuery(count=200071)

        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=71, tstop=71,
            opened=True, resolved=True,
            filter_=71, search=71
        )

        self.assertEqual(count, 200071)
        self.assertEqual(truncated, False)

        query = MockQuery(count=1072)

        count, truncated = self.reader._get_fast_count(
            query=query,
            tstart=72, tstop=72,
            opened=True, resolved=False,
            filter_=72, search=72
        )

        self.assertEqual(count, 1072)
        self.assertEqual(truncated, False)

    def test_count_alarms_by_period(self):
        day = 24 * 3600

        alarm0_id = '/fake/alarm/id0'
        event0 = {
            'connector': 'ut',
            'connector_name': 'ut0',
            'component': 'c',
            'output': '...',
            'timestamp': day / 2
        }
        alarm0 = self.manager.make_alarm(
            alarm0_id,
            event0
        )
        alarm0 = self.manager.update_state(alarm0, 1, event0)
        new_value0 = alarm0[self.manager[Alerts.ALARM_STORAGE].VALUE]
        self.manager.update_current_alarm(alarm0, new_value0)

        alarm1_id = '/fake/alarm/id1'
        event1 = {
            'connector': 'ut',
            'connector_name': 'ut0',
            'component': 'c',
            'output': '...',
            'timestamp': 3 * day / 2
        }
        alarm1 = self.manager.make_alarm(
            alarm1_id,
            event1
        )
        alarm1 = self.manager.update_state(alarm1, 1, event1)
        new_value1 = alarm1[self.manager[Alerts.ALARM_STORAGE].VALUE]
        self.manager.update_current_alarm(alarm1, new_value1)

        # Are subperiods well cut ?
        count = self.reader.count_alarms_by_period(0, day)
        self.assertEqual(len(count), 1)

        count = self.reader.count_alarms_by_period(0, day * 3)
        self.assertEqual(len(count), 3)

        count = self.reader.count_alarms_by_period(day, day * 10)
        self.assertEqual(len(count), 9)

        count = self.reader.count_alarms_by_period(
            0, day,
            subperiod={'hour': 1},
        )
        self.assertEqual(len(count), 24)

        # Are counts by period correct ?
        count = self.reader.count_alarms_by_period(0, day / 4)
        self.assertEqual(count[0]['count'], 0)

        count = self.reader.count_alarms_by_period(0, day)
        self.assertEqual(count[0]['count'], 1)

        count = self.reader.count_alarms_by_period(day / 2, 3 * day / 2)
        self.assertEqual(count[0]['count'], 2)

        # Does limit limits count ?
        count = self.reader.count_alarms_by_period(0, day, limit=100)
        self.assertEqual(count[0]['count'], 1)

        count = self.reader.count_alarms_by_period(day / 2, 3 * day / 2,
                                                   limit=1)
        self.assertEqual(count[0]['count'], 1)


if __name__ == '__main__':
    main()
