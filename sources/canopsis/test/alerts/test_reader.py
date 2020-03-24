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

from __future__ import unicode_literals

import time
import unittest

import xmlrunner
from canopsis.alerts.reader import AlertsReader
from canopsis.common import root_path
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
from canopsis.confng import Configuration, Ini
from canopsis.context_graph.manager import ContextGraph
from canopsis.logger import Logger
from canopsis.pbehavior.manager import PBehaviorManager

from base import BaseTest


class LoggerMock():
    def debug(self, *args, **kwargs):
        pass

    def info(self, *args, **kwargs):
        pass

    def warning(self, *args, **kwargs):
        pass

    def critical(self, *args, **kwargs):
        pass


class TestReader(BaseTest):
    def setUp(self):
        super(TestReader, self).setUp()

        mongo = MongoStore.get_default()
        collection = mongo.get_collection("default_testpbehavior")
        pb_coll = MongoCollection(collection)

        self.logger = Logger.get('alertsreader', '/tmp/null')
        conf = Configuration.load(PBehaviorManager.CONF_PATH, Ini)
        self.pbehavior_manager = PBehaviorManager(config=conf,
                                                  logger=self.logger,
                                                  pb_collection=pb_coll)

        self.reader = AlertsReader(config=conf,
                                   logger=self.logger,
                                   storage=self.manager.alerts_storage,
                                   pbehavior_manager=self.pbehavior_manager)

        self.reader._alarm_fields = {
            'properties': {
                'connector': {'stored_name': 'v.ctr'},
                'component': {'stored_name': 'v.cpt'},
                'entity_id': {'stored_name': 'd'}
            }
        }

    def tearDown(self):
        """Teardown"""
        super(TestReader, self).setUp()
        self.pbehavior_manager.delete(_filter={})

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
        expected_opened = {'v.resolved': None, 't': {'$lte': 2, "$gte": 1}}
        self.assertEqual(
            self.reader._get_time_filter(
                opened=True, resolved=False, tstart=1, tstop=2),
            expected_opened
        )

        # opened=False, resolved=True
        expected_resolved = {
            'v.resolved': {'$ne': None},
            't': {'$gte': 1, '$lte': 2}
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
                'expected': {'v.resolved': None, 't': {'$lte': 42, "$gte": 13}}
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
                    't': {'$gte': 0, '$lte': 0}
                }
            },
            {
                'tstart': 1,
                'tstop': 2,
                'expected': {
                    'v.resolved': {'$ne': None},
                    't': {'$gte': 1, '$lte': 2}
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

    def test__get_final_filter_bnf(self):
        view_filter = {'$and': [{'resource': 'companion cube'}]}
        time_filter = {'glados': 'shell'}
        bnf_search = 'NOT resource="turret"'
        active_columns = ['resource', 'component']

        filter_ = self.reader._get_final_filter(
            view_filter, time_filter, bnf_search, active_columns
        )

        ref_filter = {
            '$and': [
                view_filter,
                time_filter,
                {'resource': {'$not': {'$eq': 'turret'}}}
            ]
        }
        self.assertEqual(ref_filter, filter_)

    def test__get_final_filter_natural(self):
        view_filter = {'$and': [{'resource': 'companion cube'}]}
        time_filter = {'glados': 'shell'}
        search = 'turret'
        active_columns = ['resource', 'component']

        filter_ = self.reader._get_final_filter(
            view_filter, time_filter, search, active_columns
        )

        self.maxDiff = None
        ref_filter = {
            '$and': [
                view_filter,
                time_filter,
                {
                    '$or': [
                        {'resource': {
                            '$regex': '.*turret.*', '$options': 'i'}},
                        {'component': {
                            '$regex': '.*turret.*', '$options': 'i'}},
                        {'d': {
                            '$regex': '.*turret.*', '$options': 'i'}}
                    ]
                }
            ]
        }
        self.assertEqual(ref_filter, filter_)

    def test__get_final_filter_natural_numonly(self):
        view_filter = {}
        time_filter = {}
        search = 11111
        active_columns = ['resource']

        filter_ = self.reader._get_final_filter(
            view_filter, time_filter, search, active_columns
        )

        self.maxDiff = None
        res_filter = {
            '$and': [
                {'$or': [
                    {'resource': {'$options': 'i', '$regex': '.*11111.*'}},
                    {'d': {'$options': 'i', '$regex': '.*11111.*'}}
                ]}
            ]
        }
        self.assertEqual(res_filter, filter_)

    def test_contains_wildcard_dynamic_filter(self):
        # not contains dynamic wildcard filter
        view_filter = {}
        time_filter = {}
        search = 11111
        active_columns = ['resource']

        filter_ = self.reader._get_final_filter(
            view_filter, time_filter, search, active_columns
        )

        self.maxDiff = None
        res_filter = {
            '$and': [
                {'$or': [
                    {'resource': {'$options': 'i', '$regex': '.*11111.*'}},
                    {'d': {'$options': 'i', '$regex': '.*11111.*'}}
                ]}
            ]
        }
        t = self.reader.contains_wildcard_dynamic_filter(filter_)
        self.assertFalse(t)
        self.assertEqual(res_filter, filter_)

        # contains dynamic wildcard filter
        view_filter = {}
        time_filter = {}
        search = 11111
        active_columns = ['v.infos.*.type']

        filter_ = self.reader._get_final_filter(
            view_filter, time_filter, search, active_columns
        )

        t = self.reader.contains_wildcard_dynamic_filter(filter_)
        self.maxDiff = None
        res_filter = {
            '$and': [
                {'$or': [
                    {'infos_array.v.type': {'$options': 'i', '$regex': '.*11111.*'}},
                    {'d': {'$options': 'i', '$regex': '.*11111.*'}}
                ]}
            ]
        }
        self.assertTrue(t)
        self.assertEqual(res_filter, filter_)

        # contains dynamic wildcard filter
        view_filter = {'$and': [{'v.infos.*.tt': 'companion cube'}]}
        time_filter = {'glados': 'shell'}
        bnf_search = 'NOT resource="turret"'
        active_columns = ['resource', 'component']

        filter_ = self.reader._get_final_filter(
            view_filter, time_filter, bnf_search, active_columns
        )

        ref_filter = {
            '$and': [
                {'$and': [{'infos_array.v.tt': 'companion cube'}]},
                time_filter,
                {'resource': {'$not': {'$eq': 'turret'}}}
            ]
        }
        t = self.reader.contains_wildcard_dynamic_filter(filter_)
        self.assertTrue(t)
        self.assertEqual(ref_filter, filter_)

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
        new_value0 = alarm0[self.manager.alerts_storage.VALUE]
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
        new_value1 = alarm1[self.manager.alerts_storage.VALUE]
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

    def test__get_disable_entity(self):
        event = {
            'connector': '03-K64_Firefly',
            'connector_name': 'serenity',
            'component': 'Malcolm_Reynolds',
            'output': 'the big red recall button',
            'timestamp': int(time.time()) - 100,
            "source_type": "component"
        }
        alarm_id = '/strawberry'
        alarm = self.manager.make_alarm(
            alarm_id,
            event
        )

        context_manager = ContextGraph(logger=LoggerMock())
        ent_id = context_manager.get_id(event)

        entity = context_manager.create_entity_dict(ent_id,
                                                    "inara",
                                                    "component")
        entity["enabled"] = False
        context_manager._put_entities(entity)

        alarms = self.reader.get(opened=True)
        print(alarms)
        self.assertEqual(len(alarms["alarms"]), 0)


if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
