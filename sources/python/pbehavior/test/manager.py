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

from copy import deepcopy
from calendar import timegm
from datetime import datetime, timedelta
from json import dumps
from uuid import uuid4

from mock import patch, PropertyMock
from unittest import main, TestCase


from canopsis.context.manager import Context
from canopsis.middleware.core import Middleware
from canopsis.pbehavior.manager import PBehaviorManager


class TestManager(TestCase):
    def setUp(self):
        pbehavior_storage = Middleware.get_middleware_by_uri(
            'storage-default-testpbehavior://'
        )

        self.pbm = PBehaviorManager()
        self.context = Context(data_scope='test_context')

        self.pbm[PBehaviorManager.PBEHAVIOR_STORAGE] = pbehavior_storage

        self.pbehavior_id = str(uuid4())
        self.comment_id = str(uuid4())
        self.pbehavior = {
            'name': 'test',
            'filter': dumps({'connector': 'nagios-test-connector'}),
            'comments': [{
                '_id': self.comment_id,
                'author': 'test_author_comment',
                'ts': timegm(datetime.utcnow().timetuple()),
                'message': 'test_message'
            }],
            'tstart': timegm(datetime.utcnow().timetuple()),
            'tstop': timegm((datetime.utcnow() + timedelta(days=1)).timetuple()),
            'rrule': 'FREQ=DAILY',
            'enabled': True,
            'connector': 'test_connector',
            'connector_name': 'test_connector_name',
            'author': 'test_author'
        }

        data = deepcopy(self.pbehavior)
        data.update({'_id': self.pbehavior_id})
        self.pbm.pbehavior_storage.put_element(element=data)

    def tearDown(self):
        self.pbm.pbehavior_storage.remove_elements()
        self.context[Context.CTX_STORAGE].remove_elements()

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
        self.pbm.update(self.pbehavior_id, name='test_name2')
        pb = self.pbm.get(self.pbehavior_id)
        self.assertTrue(pb is not None)
        self.assertEqual(pb['name'], 'test_name2')

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

        self.pbm.pbehavior_storage.put_element(element=self.pbehavior)
        self.pbm.pbehavior_storage.put_element(element=pbehavior_1)
        self.pbm.pbehavior_storage.put_element(element=pbehavior_2)
        pbs = self.pbm.get_pbehaviors(2)

        self.assertTrue(isinstance(pbs, list))
        self.assertEqual(len(pbs), 2)

        pbs_2 = sorted(pbs, key=lambda el: el['tstart'], reverse=True)
        self.assertEqual(pbs, pbs_2)

    @patch('canopsis.pbehavior.manager.PBehaviorManager.context', new_callable=PropertyMock)
    def test_compute_pbehaviors_filters(self, mock_context):
        mock_context.return_value = Context(data_scope='test_context')
        entities = [{
            'entity_id': 1,
            'name': 'engine-test1',
            'type': 'metric-test',
            'connector': 'canopsis-test-connector',
            'connector_name': 'canopsis-test',
        }, {
            'entity_id': 2,
            'name': 'big-engine-test2',
            'type': 'metric-test',
            'connector': 'canopsis-test-connector',
            'connector_name': 'canopsis-test',
        }, {
            'entity_id': 3,
            'name': 'test_context3',
            'type': 'resource-test',
            'connector': 'nagios-test-connector',
            'connector_name': 'nagios-test',
        }]
        self.context[Context.CTX_STORAGE].put_elements(entities)
        self.pbm.compute_pbehaviors_filters()
        pb = self.pbm.get(self.pbehavior_id)

        self.assertTrue(pb is not None)
        self.assertTrue('eids' in pb)
        self.assertTrue(isinstance(pb['eids'], list))
        self.assertEqual(len(pb['eids']), 1)
        self.assertEqual(pb['eids'][0], 3)

        pb = deepcopy(self.pbehavior)
        pb.update({
                    'filter': {
                        '$or': [
                            {'type': 'resource-test'},
                            {'name': {
                                '$in':
                                    ['engine-test1', 'big-engine-test2']
                                }
                            }]
                        }
                   })
        pb_id = self.pbm.create(**pb)
        self.pbm.compute_pbehaviors_filters()
        pb = self.pbm.get(pb_id)

        self.assertTrue(pb is not None)
        self.assertTrue('eids' in pb)
        self.assertTrue(isinstance(pb['eids'], list))
        self.assertEqual(len(pb['eids']), 3)


if __name__ == '__main__':
    main()
