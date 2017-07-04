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

from calendar import timegm
from copy import deepcopy
from datetime import datetime, timedelta
from json import dumps
from uuid import uuid4

from unittest import main

from canopsis.pbehavior.manager import PBehavior

from base import BaseTest


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
            'author': 'test_author'
        }

        data = deepcopy(self.pbehavior)
        data.update({'_id': self.pbehavior_id})
        self.pbm.pbehavior_storage.put_element(element=data)

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
        pb_new.update({'filter': {'author': {'$in': ['author1, author2']}}})
        pb = self.pbm.create(**pb_new)
        self.assertTrue(pb is not None)

    def test_read(self):
        pb = self.pbm.read(_id=self.pbehavior_id)
        pbs = self.pbm.read()
        self.assertTrue(pb is not None)
        self.assertTrue(isinstance(pbs, list))
        self.assertEqual(len(pbs), 1)

    def test_update(self):
        self.pbm.update(self.pbehavior_id, name='test_name2',
                        connector=None, connector_name=None)
        pb = self.pbm.get(self.pbehavior_id)
        self.assertTrue(pb is not None)
        self.assertEqual(pb['name'], 'test_name2')
        self.assertEqual(pb['connector'], 'test_connector')
        self.assertEqual(pb['connector_name'], 'test_connector_name')

    def test_delete(self):
        self.pbm.delete(_id=self.pbehavior_id)
        pb = self.pbm.get(self.pbehavior_id)
        self.assertTrue(pb is None)

    def test_create_pbehavior_comment(self):
        self.pbm.create_pbehavior_comment(self.pbehavior_id, 'author', 'msg')
        pb = self.pbm.get(self.pbehavior_id)
        self.assertTrue('comments' in pb)
        self.assertTrue(isinstance(pb['comments'], list))
        self.assertEqual(len(pb['comments']), 2)

        self.pbm._update_pbehavior(self.pbehavior_id, {'$set': {'comments': None}})
        self.pbm.create_pbehavior_comment(self.pbehavior_id, 'author', 'msg')
        pb = self.pbm.get(self.pbehavior_id)
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

        pb = self.pbm.get(
            self.pbehavior_id,
            query={'comments': {'$elemMatch': {'_id': self.comment_id}}}
        )
        self.assertTrue(isinstance(pb['comments'], list))
        self.assertEqual(pb['comments'][0]['author'], new_author)
        self.assertEqual(pb['comments'][0]['message'], new_message)

        pb2 = self.pbm.get(
            'id_does_not_exist',
            query={'comments': {'$elemMatch': {'_id': self.comment_id}}}
        )
        self.assertIsNone(pb2)

    def test_delete_pbehavior_comment(self):
        self.pbm.create_pbehavior_comment(self.pbehavior_id, 'author', 'msg')
        pb = self.pbm.get(self.pbehavior_id)
        self.assertEqual(len(pb['comments']), 2)

        self.pbm.delete_pbehavior_comment(self.pbehavior_id, self.comment_id)
        pb = self.pbm.get(self.pbehavior_id)
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

        self.pbm.pbehavior_storage.put_elements(
            elements=(self.pbehavior, pbehavior_1, pbehavior_2)
        )
        pbs = self.pbm.get_pbehaviors(2)

        self.assertTrue(isinstance(pbs, list))
        self.assertEqual(len(pbs), 2)

        pbs_2 = sorted(pbs, key=lambda el: el['tstart'], reverse=True)
        self.assertEqual(pbs, pbs_2)

    def test_compute_pbehaviors_filters(self):
        self.pbm.context._put_entities(self.entities)
        self.pbm.compute_pbehaviors_filters()
        pb = self.pbm.get(self.pbehavior_id)

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
        pb = self.pbm.get(pb_id)

        self.assertTrue(pb is not None)
        self.assertTrue('eids' in pb)
        self.assertTrue(isinstance(pb['eids'], list))

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

        self.pbm.pbehavior_storage.put_elements(
            elements=(pbehavior_1, pbehavior_2, pbehavior_3, pbehavior_4)
        )

        self.entities[0]['timestamp'] = timegm((now + timedelta(days=2)).timetuple())
        self.entities[1]['timestamp'] = timegm(now.timetuple())
        self.entities[2]['timestamp'] = timegm((now - timedelta(days=2)).timetuple())

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
        self.pbm.pbehavior_storage.put_elements(
            elements=(pbehavior_1,)
        )

        self.entities[0]['timestamp'] = timegm((now - timedelta(days=2)).timetuple())
        self.entities[1]['timestamp'] = timegm((now - timedelta(seconds=2)).timetuple())
        self.pbm.context._put_entities(self.entities)

        # Check is a passed pbehavior is detected as not triggered
        result = self.pbm._check_pbehavior(
            self.entity_id_1, [pb_name1]
        )
        self.assertFalse(result)

    def test_get_active_pbheviors(self):
        now = datetime.utcnow()
        pbehavior_1 = deepcopy(self.pbehavior)
        pbehavior_2 = deepcopy(self.pbehavior)
        pbehavior_1.update({'eids': [self.entity_id_1],
                            'name': 'pb1',
                            'tstart': timegm(now.timetuple()),
                            'tstop': timegm((now + timedelta(days=2)).timetuple()),
                            'rrule': None
                            })
        pbehavior_2.update({'eids': [self.entity_id_3],
                            'name': 'pb2',
                            'tstart': timegm(now.timetuple())})

        self.pbm.pbehavior_storage.put_elements(
            elements=(pbehavior_1, pbehavior_2)
        )

        self.pbm.context._put_entities(self.entities)

        tab = self.pbm.get_active_pbheviors([self.entity_id_1, self.entity_id_2])
        names = [x['name'] for x in tab]
        self.assertEqual(names, ['pb1'])

if __name__ == '__main__':
    main()
