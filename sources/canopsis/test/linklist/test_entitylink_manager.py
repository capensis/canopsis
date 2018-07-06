#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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
from uuid import uuid4
from canopsis.common import root_path
from canopsis.context_graph.manager import ContextGraph
from canopsis.entitylink.manager import Entitylink
from canopsis.event import forger
from canopsis.logger import Logger
from canopsis.common.middleware import Middleware
import xmlrunner


class CheckManagerTest(unittest.TestCase):
    """
    Base class for all check manager tests.
    """

    def setUp(self):
        logger = Logger.get('linklist', Entitylink.LOG_PATH)
        self.entity_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentitylink://'
        )

        # init a context
        self.context_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )
        self.context_graph = ContextGraph(logger)
        self.context_graph.ent_storage = self.context_storage

        self.manager = Entitylink(logger=logger,
                                  storage=self.entity_storage,
                                  context_graph=self.context_graph)

        self.name = 'testentitylink'
        self.id_ = str(uuid4())
        self.ids = [self.id_]
        self.document_content = {
            'computed_links': [{
                'label': 'link1',
                'url': 'http://www.canopsis.org'
            }]
        }

        self.manager.put(
            self.id_,
            self.document_content
        )

        entity = {
            '_id': "a_component",
            'type': 'component',
            'name': 'conn-name1',
            'depends': [],
            'impact': [],
            'measurements': [],
            'infos': {}
        }
        self.context_graph.create_entity(entity)

    def tearDown(self):
        self.entity_storage.remove_elements()
        self.context_storage.remove_elements()

    def linklist_count_equals(self, count):
        result = list(self.manager.find(ids=self.ids))
        self.assertEqual(len(result), count)


class LinkListTest(CheckManagerTest):

    def test_put(self):
        self.linklist_count_equals(1)

    def test_get(self):
        self.manager.put(
            self.id_ + '1',
            self.document_content
        )

        self.linklist_count_equals(1)

        result = self.manager.find()
        self.assertGreaterEqual(len(list(result)), 2)

        result = self.manager.find(limit=1)
        self.assertEqual(len(list(result)), 1)

    def test_remove(self):
        self.linklist_count_equals(1)

        self.manager.remove(self.ids)

        self.linklist_count_equals(0)

    def test_get_id_from_event(self):
        event = forger(
            event_type="check",
            component="a_component",
            #state=1,
            #output="output",
        )
        res = self.manager.get_id_from_event(event=event)
        self.assertEqual(res, 'a_component')

    def test_get_or_create_from_event(self):
        event = forger(
            event_type="check",
            component="a_component",
        )

        res = self.manager.get_or_create_from_event(event=event)
        self.assertTrue(isinstance(res, dict))
        self.assertEqual(res['_id'], 'a_component')

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
