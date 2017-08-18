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
from canopsis.event.eventslogmanager import EventsLog
from canopsis.common.utils import singleton_per_scope


class EventsLogManagerTest(TestCase):
    """Base class for eventslogmanager tests
    """

    def setUp(self):
        """initialize a manager
        """
        self.manager = singleton_per_scope(EventsLog)

    def tearDown(self):
        pass


class EventsLogTest(EventsLogManagerTest):

    def test_get_eventlog_count_by_period(self):

        def mock_find(query={}, limit=100, with_count=True):
            return None, 5

        self.manager[EventsLog.EVENTSLOG_STORAGE].find_elements = mock_find

        result = self.manager.get_eventlog_count_by_period(
            1433113200,
            1435705200
        )

        self.assertEquals(result[0].get('count'), 5)

if __name__ == '__main__':
    main()
