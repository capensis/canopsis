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

import unittest
from uuid import uuid4
from canopsis.common import root_path

from canopsis.linklist.manager import Linklist
from canopsis.logger import Logger
from canopsis.common.middleware import Middleware
import xmlrunner


class CheckManagerTest(unittest.TestCase):
    """
    Base class for all check manager tests.
    """

    def setUp(self):
        logger = Logger.get('linklist', Linklist.LOG_PATH)
        self.linklist_storage = Middleware.get_middleware_by_uri(
            'storage-default-testlinklist://'
        )
        self.manager = Linklist(logger=logger,
                                storage=self.linklist_storage)

        self.name = 'testlinklist'
        self.id_ = str(uuid4())
        self.ids = [self.id_]
        self.document_content = {
            'id': self.id_,
            'name': self.name,
            'linklist': ['http://canopsis.org'],
            'mfilter': '{"$and": [{"connector": "collectd"}]}'
        }

    def tearDown(self):
        self.linklist_storage.remove_elements()

    def linklist_count_equals(self, count):
        result = list(self.manager.find(ids=self.ids))
        self.assertEqual(len(result), count)


class LinkListTest(CheckManagerTest):

    def test_put(self):
        self.manager.put(
            self.document_content
        )

        self.linklist_count_equals(1)

    def test_get(self):
        self.manager.put(
            self.document_content
        )

        self.manager.put({
            'name': self.name + '1',
            'linklist': ['http://canopsis.org'],
            'mfilter': '{"$and": [{"connector": "collectd"}]}'
        })

        self.linklist_count_equals(1)

        result = self.manager.find()
        self.assertGreaterEqual(len(list(result)), 2)

        result = self.manager.find(limit=1)
        self.assertEqual(len(list(result)), 1)

    def test_remove(self):
        self.linklist_count_equals(0)

        self.manager.put(
            self.document_content
        )

        self.linklist_count_equals(1)

        self.manager.remove(self.ids)

        self.linklist_count_equals(0)


if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
