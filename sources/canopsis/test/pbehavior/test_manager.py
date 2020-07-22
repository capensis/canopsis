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

import unittest
import time
import datetime
from calendar import timegm
from copy import deepcopy
from datetime import datetime, timedelta
from json import dumps
from uuid import uuid4

import xmlrunner
from canopsis.common import root_path
from canopsis.models.pbehavior import PBehavior as PBModel
from canopsis.pbehavior.manager import PBehavior
from canopsis.pbehavior.manager import PBehaviorManager

from test_base import BaseTest
import re


class TestManager(BaseTest):
    def setUp(self):
        super(TestManager, self).setUp()

        self.pbehavior_id = str(uuid4())
        self.comment_id = str(uuid4())
        now = datetime.utcnow()
        self.pbehavior = {
            'name': 'test',
            'filter': dumps({'connector': 'nagios-test-connector'}),
            'comments': [{
                '_id': self.comment_id,
                'author': 'test_author_comment',
                'ts': timegm(now.timetuple()),
                'message': 'test_message'
            }],
            'tstart': timegm(now.timetuple()),
            'tstop': timegm((now + timedelta(days=1)).timetuple()),
            'rrule': 'FREQ=DAILY;INTERVAL=2;COUNT=3',
            'enabled': True,
            'connector': 'test_connector',
            'connector_name': 'test_connector_name',
            'author': 'test_author',
            PBehavior.TYPE: 'pause',
            'reason': 'reason is treason'
        }

        data = deepcopy(self.pbehavior)
        data.update({'_id': self.pbehavior_id})
        self.pbm.collection.insert(data)

        self.entity_id_1 = '/component/collectd/pbehavior/test1/'
        self.entity_id_2 = '/component/collectd/pbehavior/test2/'
        self.entity_id_3 = '/component/collectd/pbehavior/test3/'

        self.entities = [{
            '_id': self.entity_id_1,
            'name': 'engine-test1',
            'type': 'metric-test',
            PBehavior.FILTER: {},
            'infos': {}
        }, {
            '_id': self.entity_id_2,
            'name': 'big-engine-test2',
            'type': 'metric-test',
            PBehavior.FILTER: {},
            'infos': {}
        }, {
            '_id': self.entity_id_3,
            'name': 'test_context3',
            'type': 'resource-test',
            PBehavior.FILTER: {},
            'infos': {}
        }]

    def test_create(self):
        pb = self.pbm.create(**self.pbehavior)
        self.assertTrue(pb is not None)

        pb_new = deepcopy(self.pbehavior)
        pb_new[PBehavior.RRULE] = ''
        pb_new.update({'filter': {'author': {'$in': ['author1, author2']}}})
        pb_new[PBehavior.TSTOP] = int(time.time()) - 24 * 3600
        pb = self.pbm.create(**pb_new)
        self.assertTrue(pb is not None)

        # create with expired
        pb_new_2 = deepcopy(self.pbehavior)
        pb_new_2['pbh_id'] = pb
        pb_new_2[PBehavior.TSTOP] = int(time.time()) + 24 * 3600
        e = None
        try:
            self.pbm.create(**pb_new_2)
        except Exception as e:
            pass
        finally:
            self.assertTrue(e is not None)
            self.assertTrue(isinstance(e, ValueError))
            e = None

        pb_new_2['replace_expired'] = True
        new_pb = self.pbm.create(**pb_new_2)
        pbh = self.pbm.get(_id=None)
        self.assertGreaterEqual(pbh.get('count'), 1)
        regex = r"EXP\d+-{}".format(pb)
        self.assertEqual(len(filter(lambda x: re.match(regex, x['_id']), pbh.get('data'))), 1)
        self.assertEqual(pbh.get('data')[0].get('name'), pb_new_2[PBehavior.NAME])
        self.assertTrue(new_pb, pb)


    def test_read(self):
        pb = self.pbm.read(_id=self.pbehavior_id).get('data')[0]
        pbs = self.pbm.read().get('data')
        self.assertTrue(pb is not None)
        self.assertTrue(isinstance(pbs, list))
        self.assertEqual(len(pbs), 1)
        self.assertEqual(pbs[0][PBehavior.TYPE],
                         self.pbehavior[PBehavior.TYPE])
        self.assertEqual(pbs[0]['reason'], self.pbehavior['reason'])

    def test_update(self):
        self.pbm.update(self.pbehavior_id, name='test_name2',
                        connector=None, connector_name=None)
        pb = self.pbm.get(self.pbehavior_id).get('data')[0]
        self.assertTrue(pb is not None)
        self.assertEqual(pb['name'], 'test_name2')
        self.assertEqual(pb['connector'], 'test_connector')
        self.assertEqual(pb['connector_name'], 'test_connector_name')

    def test_delete(self):
        self.pbm.delete(_id=self.pbehavior_id)
        pb = self.pbm.get(self.pbehavior_id).get('data')
        self.assertTrue(len(pb) == 0)

    def test_create_pbehavior_comment(self):
        self.pbm.create_pbehavior_comment(self.pbehavior_id, 'author', 'msg')
        pb = self.pbm.get(self.pbehavior_id).get('data')[0]
        self.assertTrue('comments' in pb)
        self.assertTrue(isinstance(pb['comments'], list))
        self.assertEqual(len(pb['comments']), 2)

        self.pbm._update_pbehavior(
            self.pbehavior_id, {'$set': {'comments': []}})
        self.pbm.create_pbehavior_comment(self.pbehavior_id, 'author', 'msg')
        pb = self.pbm.get(self.pbehavior_id).get('data')[0]
        self.assertTrue('comments' in pb)
        self.assertTrue(isinstance(pb['comments'], list))
        self.assertEqual(len(pb['comments']), 1)

    def test_update_pbehavior_comment(self):
        new_author, new_message = 'new_author', 'new_message'
        result = self.pbm.update_pbehavior_comment(
            self.pbehavior_id, self.comment_id,
            author=new_author,
            message=new_message
        )
        self.assertIsNotNone(result)
        self.assertTrue(isinstance(result, dict))
        self.assertEqual(result['message'], new_message)
        self.assertEqual(result['author'], new_author)

        pb = self.pbm.get(self.pbehavior_id).get('data')[0]
        self.assertTrue(isinstance(pb['comments'], list))
        self.assertEqual(pb['comments'][0]['author'], new_author)
        self.assertEqual(pb['comments'][0]['message'], new_message)

        pb2 = self.pbm.get('id_does_not_exist').get('data')
        self.assertEqual(pb2, [])

    def test_delete_pbehavior_comment(self):
        self.pbm.create_pbehavior_comment(self.pbehavior_id, 'author', 'msg')
        pb = self.pbm.get(self.pbehavior_id).get('data')[0]
        self.assertEqual(len(pb['comments']), 2)

        self.pbm.delete_pbehavior_comment(self.pbehavior_id, self.comment_id)
        pb = self.pbm.get(self.pbehavior_id).get('data')[0]
        self.assertEqual(len(pb['comments']), 1)

    def test_get_pbehaviors(self):
        pbehavior_1 = deepcopy(self.pbehavior)
        pbehavior_2 = deepcopy(self.pbehavior)
        self.pbehavior.update({'eids': [1, 2, 3],
                               'tstart': timegm((datetime.utcnow() + timedelta(days=3)).timetuple())})
        pbehavior_1.update({'eids': [2, 4, 5],
                            'tstart': timegm((datetime.utcnow() + timedelta(days=2)).timetuple())})
        pbehavior_2.update({'eids': [5, 6],
                            'tstart': timegm((datetime.utcnow() + timedelta(days=1)).timetuple())})

        self.pbm.collection.insert([self.pbehavior, pbehavior_1, pbehavior_2])
        pbs = self.pbm.get_pbehaviors(2)

        self.assertTrue(isinstance(pbs, list))
        self.assertEqual(len(pbs), 2)

        pbs_2 = sorted(pbs, key=lambda el: el['tstart'], reverse=True)
        self.assertEqual(pbs, pbs_2)

    def test_compute_pbehaviors_filters(self):
        self.pbm.context._put_entities(self.entities)
        self.pbm.compute_pbehaviors_filters()
        pb = self.pbm.get(self.pbehavior_id).get('data')[0]

        self.assertTrue(pb is not None)
        self.assertTrue('eids' in pb)
        self.assertTrue(isinstance(pb['eids'], list))

        pb = deepcopy(self.pbehavior)
        pb.update({
            'filter': {
                '$or': [
                    {'type': 'resource-test'},
                    {'name': {'$in': ['engine-test1', 'big-engine-test2']}}
                ]
            }
        })
        pb_id = self.pbm.create(**pb)
        self.pbm.compute_pbehaviors_filters()
        pb = self.pbm.get(pb_id).get('data')[0]

        self.assertTrue(pb is not None)
        self.assertTrue('eids' in pb)
        self.assertTrue(isinstance(pb['eids'], list))

        # A bad pbehavior filter does not crash compute()
        pb = deepcopy(self.pbehavior)
        pb.update({
            'filter': "\"{\\\"_id\\\": \\\"Sc_aude_eid_02/scenario\\\"}\""
        })
        pb_id = self.pbm.create(**pb)
        self.pbm.compute_pbehaviors_filters()
        pb = self.pbm.get(pb_id)

    def test_check_pbehaviors(self):
        now = datetime.utcnow()
        pbehavior_1 = deepcopy(self.pbehavior)
        pbehavior_2 = deepcopy(self.pbehavior)
        pbehavior_3 = deepcopy(self.pbehavior)
        pbehavior_4 = deepcopy(self.pbehavior)

        pb_name1, pb_name2, pb_name3, pb_name4 = 'name1', 'name2', 'name3', 'name4'

        pbehavior_1.update(
            {
                'name': pb_name1,
                'eids': [self.entity_id_1, self.entity_id_2],
                'tstart': timegm(now.timetuple()),
                'tstop': timegm((now + timedelta(days=8)).timetuple())
            }
        )

        pbehavior_2.update({'name': pb_name2})

        pbehavior_3.update(
            {
                'name': pb_name3,
                'eids': [self.entity_id_2, self.entity_id_3],
                'tstart': timegm(now.timetuple()),
                'tstop': timegm((now + timedelta(days=8)).timetuple())
            }
        )

        pbehavior_4.update({'name': pb_name4})

        self.pbm.collection.insert(
            [pbehavior_1, pbehavior_2, pbehavior_3, pbehavior_4])

        self.entities[0]['timestamp'] = timegm(
            (now + timedelta(days=2)).timetuple())
        self.entities[1]['timestamp'] = timegm(now.timetuple())
        self.entities[2]['timestamp'] = timegm(
            (now - timedelta(days=2)).timetuple())

        self.pbm.context._put_entities(self.entities)

        result = self.pbm.check_pbehaviors(
            self.entity_id_1, [pb_name1, pb_name2], [pb_name3, pb_name4]
        )
        self.assertTrue(result)

        result = self.pbm.check_pbehaviors(
            self.entity_id_3, [pb_name3, pb_name4], [pb_name1, pb_name2]
        )
        self.assertFalse(result)

    def test_check_pbehavior(self):
        now = datetime.utcnow()
        pbehavior_1 = deepcopy(self.pbehavior)
        pb_name1 = 'name1'
        pbehavior_1.update(
            {
                'name': pb_name1,
                'eids': [self.entity_id_1, self.entity_id_2],
                'tstart': timegm(now.timetuple()),
                'tstop': timegm((now + timedelta(days=8)).timetuple())
            }
        )
        self.pbm.collection.insert(pbehavior_1)

        self.entities[0]['timestamp'] = timegm(
            (now - timedelta(days=2)).timetuple())
        self.entities[1]['timestamp'] = timegm(
            (now - timedelta(seconds=2)).timetuple())
        self.pbm.context._put_entities(self.entities)

        # Check is a passed pbehavior is detected as not triggered
        result = self.pbm._check_pbehavior(
            self.entity_id_1, [pb_name1]
        )
        self.assertFalse(result)

        # Check for bad tstart/stop values
        pbehavior_2 = deepcopy(self.pbehavior)
        pbehavior_2.update(
            {
                'name': pb_name1,
                'eids': [self.entity_id_1, self.entity_id_2],
                'tstart': None,
                'tstop': None
            }
        )
        self.pbm.collection.insert(pbehavior_2)
        result = self.pbm._check_pbehavior(self.entity_id_1, [pb_name1])
        self.assertFalse(result)

        pbehavior_3 = deepcopy(self.pbehavior)
        pbehavior_3.update(
            {
                'name': pb_name1,
                'eids': [self.entity_id_1, self.entity_id_2],
                'tstart': 'han',
                'tstop': 'solo'
            }
        )
        self.pbm.collection.insert(pbehavior_3)
        result = self.pbm._check_pbehavior(self.entity_id_1, [pb_name1])
        self.assertFalse(result)

    def test_get_active_pbheviors(self):
        now = datetime.utcnow()
        pbehavior_1 = deepcopy(self.pbehavior)
        pbehavior_2 = deepcopy(self.pbehavior)
        pbehavior_1.update({
            'eids': [self.entity_id_1],
            'name': 'pb1',
            'tstart': timegm(now.timetuple()),
            'tstop': timegm((now + timedelta(days=2)).timetuple()),
            'rrule': None
        })
        pbehavior_2.update({'eids': [self.entity_id_3],
                            'name': 'pb2',
                            'tstart': timegm(now.timetuple())})

        self.pbm.collection.insert([pbehavior_1, pbehavior_2])

        self.pbm.context._put_entities(self.entities)

        tab = self.pbm.get_active_pbehaviors([self.entity_id_1,
                                              self.entity_id_2])
        names = [x['name'] for x in tab]
        self.assertEqual(names, ['pb1'])

    def test_get_active_pbehavior_from_type(self):
        pbehavior_1 = deepcopy(self.pbehavior)
        pbehavior_2 = deepcopy(self.pbehavior)
        pb_name1, pb_name2, = 'cheerfull', 'blue'

        now = int(time.time())
        hour = 3600

        pbehavior_1.update(
            {
                'name': pb_name1,
                'eids': [self.entity_id_1, self.entity_id_2],
                PBehavior.TYPE: 'maintenance',
                "rrule": 'FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR,SA,SU',
                'tstart': now - hour,
                'tstop': now + hour,
                'timezone': "UTC"
            }
        )

        pbehavior_2.update({'name': pb_name2})

        self.pbm.collection.insert([pbehavior_1, pbehavior_2])

        self.pbm.context._put_entities(self.entities)

        result = self.pbm.get_active_pbehaviors_from_type(['maintenance'])
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0][PBehavior.TYPE], 'maintenance')

    def test_check_active_pbehavior(self):
        now = int(time.mktime(datetime.utcnow().timetuple()))
        hour = 3600

        # tstart < now < tstop
        pb_w_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now - 1,
            now + 1,
            'FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR,SA,SU',
            'test'
        ).to_dict()

        self.assertTrue(self.pbm.check_active_pbehavior(now, pb_w_rrule))

        # tstart is one hour ahead from now
        pb_w_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now + hour * 1,
            now + hour * 2,
            'FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR,SA,SU',
            'test'
        ).to_dict()

        self.assertFalse(self.pbm.check_active_pbehavior(now, pb_w_rrule))

        # tstart is two hour behind from now, tstop one hour
        pb_w_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now - hour * 2,
            now - hour * 1,
            'FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR,SA,SU',
            'test'
        ).to_dict()

        self.assertFalse(self.pbm.check_active_pbehavior(now, pb_w_rrule))

        # no rrule, tstart and tstop in the past
        pb_n_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now - hour * 2,
            now - hour * 1,
            '',
            'test'
        ).to_dict()

        self.assertFalse(self.pbm.check_active_pbehavior(now, pb_n_rrule))

        # no rrule, now between tstart and tstop
        pb_n_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now - hour * 1,
            now + hour * 1,
            '',
            'test'
        ).to_dict()

        self.assertTrue(self.pbm.check_active_pbehavior(now, pb_n_rrule))

        # no rrule, tstart and tstop in the future
        pb_n_rrule = PBModel(
            'w_rrule',
            'w_rrule',
            {},
            now + hour * 1,
            now + hour * 2,
            '',
            'test'
        ).to_dict()

        self.assertFalse(self.pbm.check_active_pbehavior(now, pb_n_rrule))

    def test_check_active_pbehavior_2(self):
        timestamps = []

        # Vendredi 15 Juin 2018 15h13
        timestamps.append((False, 1529154801 - 24 * 3600))
        timestamps.append((True, 1529154801 - 24 * 3600 + 5 * 3600)
                          )  # Vendredi 15 Juin 2018 20h13
        timestamps.append((True, 1529154801))  # Samedi 16 Juin 2018 15h13
        timestamps.append((True, 1529290800))  # Lundi 18 Juin 2018 05h00

        timestamps.append((False, 1529308800))  # Lundi 18 Juin 2018 10h00
        timestamps.append((False, 1529308800 + 7 * 24 * 3600))
        timestamps.append((False, 1529308800 + 7 * 24 * 3600 * 2))
        timestamps.append((False, 1529308800 + 7 * 24 * 3600 * 3))
        timestamps.append((False, 1529308800 + 7 * 24 * 3600 * 4))
        timestamps.append((False, 1529308800 + 7 * 24 * 3600 * 5))

        timestamps.append((True, 1529740800))  # Samedi 23 Juin 2018 10h00
        timestamps.append((True, 1529740800 + 7 * 24 * 3600))  # +7j
        timestamps.append((True, 1529740800 + 7 * 24 * 3600 * 2))  # ...
        timestamps.append((True, 1529740800 + 7 * 24 * 3600 * 3))
        timestamps.append((True, 1529740800 + 7 * 24 * 3600 * 4))
        timestamps.append((True, 1529740800 + 7 * 24 * 3600 * 5))

        pbehavior = {
            "rrule": "FREQ=WEEKLY;BYDAY=FR",
            "tstart": 1529085600,
            "tstop": 1529294400,
            "timezone": "UTC"
        }

        for i, ts in enumerate(timestamps):
            res = self.pbm.check_active_pbehavior(ts[1], pbehavior)
            self.assertEqual(res, ts[0])

    def test_get_active_intervals(self):
        day = 24 * 3600
        tstart = 1530288000  # 2018/06/29 18:00:00
        tstop = tstart + 3600

        pbehavior = {
            'rrule': 'FREQ=DAILY',
            'tstart': tstart,
            'tstop': tstop
        }

        # after = tstart
        expected_intervals = [
            (tstart, tstop),
            (tstart + day, tstop + day),
            (tstart + 2 * day, tstop + 2 * day),
            (tstart + 3 * day, tstop + 3 * day),
            (tstart + 4 * day, tstop + 4 * day),
        ]
        intervals = list(PBehaviorManager.get_active_intervals(
            tstart, tstart + 5 * day, pbehavior))
        self.assertEqual(intervals, expected_intervals)

        # after < tstart
        intervals = list(PBehaviorManager.get_active_intervals(
            tstart - 3 * day, tstart + 5 * day, pbehavior))
        self.assertEqual(intervals, expected_intervals)

        # after > tstart
        intervals = list(PBehaviorManager.get_active_intervals(
            tstart + 2 * day, tstart + 5 * day, pbehavior))
        expected_intervals = [
            (tstart + 2 * day, tstop + 2 * day),
            (tstart + 3 * day, tstop + 3 * day),
            (tstart + 4 * day, tstop + 4 * day),
        ]
        self.assertEqual(intervals, expected_intervals)

        intervals = list(PBehaviorManager.get_active_intervals(
            tstart + 2 * day + 1800, tstart + 5 * day, pbehavior))
        expected_intervals = [
            (tstart + 2 * day + 1800, tstop + 2 * day),
            (tstart + 3 * day, tstop + 3 * day),
            (tstart + 4 * day, tstop + 4 * day),
        ]
        self.assertEqual(intervals, expected_intervals)

        # before < tstart
        intervals = list(PBehaviorManager.get_active_intervals(
            tstart - 3 * day, tstart - 2 * day, pbehavior))
        expected_intervals = []
        self.assertEqual(intervals, expected_intervals)

    def test_get_intervals_with_pbehaviors_by_eid(self):
        day = 24 * 3600

        tstart1 = 1530288000  # 2018/06/29 18:00:00
        tstop1 = tstart1 + 3600

        tstart2 = tstart1 + 1800
        tstop2 = tstop1 + 1800

        pbehavior1 = deepcopy(self.pbehavior)
        pbehavior2 = deepcopy(self.pbehavior)
        pbehavior1.update({
            'eids': [1],
            'rrule': 'FREQ=DAILY',
            'tstart': tstart1,
            'tstop': tstop1
        })
        pbehavior2.update({
            'eids': [1],
            'rrule': 'FREQ=DAILY',
            'tstart': tstart2,
            'tstop': tstop2
        })

        self.pbm.collection.insert([pbehavior1, pbehavior2])

        expected_intervals = [
            (tstart1, tstart1, False),
            (tstart1, tstop2, True),
            (tstop2, tstart1 + day, False),
            (tstart1 + day, tstop2 + day, True),
            (tstop2 + day, tstart1 + 2 * day, False),
            (tstart1 + 2 * day, tstop2 + 2 * day, True),
            (tstop2 + 2 * day, tstart1 + 3 * day, False),
            (tstart1 + 3 * day, tstop2 + 3 * day, True),
            (tstop2 + 3 * day, tstart1 + 4 * day, False),
            (tstart1 + 4 * day, tstop2 + 4 * day, True),
            (tstop2 + 4 * day, tstart1 + 5 * day, False),
        ]
        intervals = list(self.pbm.get_intervals_with_pbehaviors_by_eid(
            tstart1, tstart1 + 5 * day, 1))
        self.assertEqual(intervals, expected_intervals)

        # Entity without pbehaviors
        expected_intervals = [
            (tstart1, tstart1 + 5 * day, False),
        ]
        intervals = list(self.pbm.get_intervals_with_pbehaviors_by_eid(
            tstart1, tstart1 + 5 * day, 2))
        self.assertEqual(intervals, expected_intervals)

    def test_is_pbh_expired(self):
        pbehavior1 = deepcopy(self.pbehavior)
        now = datetime.utcnow()
        pbehavior1[PBehavior.TSTART] = timegm((now - timedelta(minutes=1)).timetuple())
        pbehavior1[PBehavior.TSTOP] = timegm((now + timedelta(days=1)).timetuple())
        pbehavior1[PBehavior.RRULE] = ''
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm(now.timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=2)).timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=23)).timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now - timedelta(hours=23)).timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now - timedelta(hours=2)).timetuple())))
        self.assertTrue(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=25)).timetuple())))
        self.assertTrue(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=24)).timetuple())))
        self.assertTrue(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=48)).timetuple())))

        # rrule with count
        pbehavior1[PBehavior.TIMEZONE] = "UTC"
        pbehavior1[PBehavior.RRULE] = 'FREQ=HOURLY;INTERVAL=2;COUNT=3'
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm(now.timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=3)).timetuple())))
        self.assertTrue(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=5)).timetuple())))
        self.assertTrue(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=6)).timetuple())))
        self.assertTrue(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=8)).timetuple())))

        # rrule without count
        pbehavior1[PBehavior.RRULE] = 'FREQ=HOURLY;INTERVAL=2'
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm(now.timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(hours=5)).timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(days=5)).timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(minutes=5)).timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(days=365)).timetuple())))
        self.assertFalse(self.pbm.is_pbh_expired(pbehavior1, timegm((now + timedelta(days=31)).timetuple())))

    def test_generate_event(self):
        day = 24 * 3600
        now = int(time.time())  # start-time set to nearly day for fast recurrent rule calculation
        tstart1 = now - now % day - 2 * day  #
        tstop1 = tstart1 + 3600

        tstart2 = tstart1 + 1800
        tstop2 = tstop1 + 1800
        alarm_coll = self.pbm.alarmAdapter.collection
        alarm_coll.remove()

        entity_id = 'é'

        self.pbm.context._put_entities([{
            '_id': entity_id,
            'name': 'pbehavior-engine-test1',
            'depends': ["connector/connector_name"],
            'type': 'pbehavior-metric-test',
            PBehavior.FILTER: {},
            'infos': {}
        }])

        pbehavior1 = deepcopy(self.pbehavior)
        pbehavior1.update({
            "_id": "259f5636-132e-11e9-a604-0242ac10a037",
            "filter": "{\"_id\": \"xxxxxx/scenario\"}",
            "name": "downtime",
            "author": "xxx",
            "enabled": True,
            "type_": "pause",
            "comments": [
                {
                    "message": "Scénario  non fonctionnel",
                    "_id": "43043416-db8c-421c-947f-7a21100bd6f7",
                    "author": "xxx"
                }
            ],
            "connector": "canopsis",
            "reason": "Problème Scénario",
            "connector_name": "canopsis",
            "rrule": 'invalid rrule ff',
            "tstart": 1546942445,
            "tstop": 2147483647,
            "eids": [entity_id]
        })
        self.pbm.collection.remove()
        self.pbm.collection.insert([pbehavior1])
        now = int(time.time())
        events = list(self.pbm.generate_pbh_event(now))
        self.assertEqual(len(events), 0)

        pbehavior1.update({
            "_id": "259f5636-132e-11e9-a604-0242ac10a037",
            "filter": "{\"_id\": 1}",
            "name": "downtime",
            "author": "xxx",
            "enabled": True,
            "type_": "pause",
            "comments": [
                {
                    "message": "Scénario  non fonctionnel",
                    "_id": "43043416-db8c-421c-947f-7a21100bd6f7",
                    "author": "xxx"
                }
            ],
            "connector": "canopsis",
            "reason": "Problème Scénario",
            "connector_name": "canopsis",
            "rrule": '',
            "tstart": 1546942445,
            "tstop": 2147483647,
            "eids": [entity_id]
            }
        )

        alarm_coll.insert({
            "_id": "alarm_2",
            "d": entity_id,
            "v": {
                "state": {
                    "_t": 'stateinc',
                    "t": 1587429072,
                    "a": 'webinar.webinar',
                    "m": 'noveo alarm',
                    "val": 3
                },
                "status": {
                    "_t": 'statusinc',
                    "t": 1587429072,
                    "a": 'webinar.webinar',
                    "m": 'noveo alarm',
                    "val": 1
                },
                "resolved": None
            },
            "component": "com1",
            "connector": "conn1",
            "connector_name": "conn_name1"
        })

        self.pbm.collection.update({'_id': pbehavior1['_id']}, {"$set": pbehavior1})
        now = int(time.time())
        events = list(self.pbm.generate_pbh_event(now))
        self.assertEqual(len(events), 1)

        pbehavior1.update({
            'name': 'hourly test',
            'eids': [entity_id],
            'rrule': 'FREQ=HOURLY',
            'tstart': tstart1,
            'tstop': tstop1
        })
        self.pbm.collection.update({'_id': pbehavior1['_id']}, {"$set": pbehavior1})
        now = int(time.time())
        events = list(self.pbm.generate_pbh_event(now))
        # pbhenter
        self.assertEqual(len(events), 2)
        self.assertEqual(events[0]['event_type'], 'pbhleave')
        self.assertEqual(events[0]['display_name'], 'downtime')
        self.assertEqual(events[0]['timestamp'], 2147483647)
        old_now = now

        next_hour = old_now + 3599
        events = list(self.pbm.generate_pbh_event(next_hour))
        # pbhleave for above pbhenter because pbehavior reach its due date
        self.assertEqual(len(events), 2)
        self.assertEqual(events[0]['event_type'], 'pbhleave')
        self.assertEqual(events[0]['display_name'], 'hourly test')
        self.assertEqual(events[0]['timestamp'], old_now + 3600 - old_now % 3600)
        # new pbhenter
        self.assertEqual(events[1]['event_type'], 'pbhenter')
        self.assertEqual(events[1]['display_name'], 'hourly test')
        self.assertEqual(events[1]['timestamp'], next_hour - next_hour % 3600)

        # modify rrule
        pbehavior1.update({
            'name': 'minutely test',
            'eids': [entity_id],
            'rrule': 'FREQ=MINUTELY;INTERVAL=15',
        })
        self.pbm.collection.update({'_id': pbehavior1['_id']}, {"$set": pbehavior1})
        pivot_time = next_hour - next_hour % 3600 + 13 * 60 * 3
        events = list(self.pbm.generate_pbh_event(pivot_time))
        self.assertEqual(len(events), 2)
        # send pbhleave for last pbhenter
        self.assertEqual(events[0]['event_type'], 'pbhleave')
        self.assertEqual(events[0]['display_name'], 'hourly test')
        self.assertEqual(events[0]['timestamp'], next_hour - next_hour % 3600 + 3600)
        # send new pbhenter for new rrule
        self.assertEqual(events[1]['event_type'], 'pbhenter')
        self.assertEqual(events[1]['display_name'], 'minutely test')
        # because pivot time is 39th minute of hour
        # so with FREQ=MINUTELY;INTERVAL=15 --> last start time is: 30th minute of hour
        self.assertEqual(events[1]['timestamp'], next_hour - next_hour % 3600 + 30 * 60)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
