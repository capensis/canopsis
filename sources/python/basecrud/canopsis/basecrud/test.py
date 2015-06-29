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

from unittest import TestCase, main
from canopsis.basecrud.manager import BaseCrud
from uuid import uuid4

DEBUG = False


class BaseCrudManagerTest(TestCase):
    """
    Base class for all check manager tests.
    """

    def setUp(self):
        """
        initialize a manager.
        """

        self.manager = BaseCrud()
        self.name = 'testsession'
        self.id = str(uuid4())
        self.ids = [self.id]
        self.document_content = {
            'content': [{
                'label': 'test',
                'value': 'http://www.canopsis.org'
            }]
        }

    def clean(self):
        self.manager.remove(self.ids)

    def count_equals(self, count):
        result = list(self.manager.find(ids=self.ids))
        if DEBUG:
            print result
        result = len(result)
        self.assertEqual(result, count)


class BaseCrudTest(BaseCrudManagerTest):

    def test_put(self):
        self.clean()

        self.manager.put(
            self.id,
            self.document_content
        )

        self.count_equals(1)

    def test_get(self):
        self.clean()

        self.manager.put(
            self.id,
            self.document_content
        )

        self.manager.put(
            self.id + '1',
            self.document_content
        )

        self.count_equals(1)

        result = self.manager.find()
        self.assertGreaterEqual(len(list(result)), 2)

        result = self.manager.find(limit=1)
        self.assertEqual(len(list(result)), 1)

    def test_remove(self):
        self.clean()

        self.count_equals(0)

        self.manager.put(
            self.id,
            self.document_content
        )

        self.count_equals(1)

        self.manager.remove(self.ids)

        self.count_equals(0)


if __name__ == '__main__':
    # Cannot test base class as no storage is associated
    pass
