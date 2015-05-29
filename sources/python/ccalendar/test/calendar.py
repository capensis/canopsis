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
from canopsis.ccalendar.manager import Calendar
from uuid import uuid4

DEBUG = False


class CalendarManagerTest(TestCase):
    """
    Base class for all calendar manager tests.
    """

    def setUp(self):
        """
        initialize a manager.
        """

        self.manager = Calendar()
        self.name = 'testcalendar'
        self.id = str(uuid4())
        self.ids = [self.id]
        self.document_content = {
            'calendar_property': [{
                'label': 'label1',
                'value': 'value1'
            }]
        }

    def clean(self):
        self.manager.remove(self.ids)

    def get_calendar(self):
        return self.manager.find(ids=self.ids)

    def calendar_count_equals(self, count):
        result = list(self.get_calendar())
        if DEBUG:
            print result
        result = len(result)
        self.assertEqual(result, count)


class CalendarTest(CalendarManagerTest):

    def test_put(self):
        self.clean()

        self.manager.put(
            self.id,
            self.document_content
        )

        self.calendar_count_equals(1)

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

        self.calendar_count_equals(1)

        result = self.manager.find()
        self.assertGreaterEqual(len(list(result)), 2)

        result = self.manager.find(limit=1)
        self.assertEqual(len(list(result)), 1)

    def test_remove(self):
        self.clean()

        self.calendar_count_equals(0)

        self.manager.put(
            self.id,
            self.document_content
        )

        self.calendar_count_equals(1)

        self.manager.remove(self.ids)

        self.calendar_count_equals(0)


if __name__ == '__main__':
    main()
