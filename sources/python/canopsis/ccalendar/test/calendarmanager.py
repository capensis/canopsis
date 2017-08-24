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
from canopsis.ccalendar.manager import CalendarManager


class CalendarManagerTest(TestCase):
    """
    Base class for all calendar manager tests.

    A event calendar is herited from a vevent,
    that is why this class is herited from
    VeventManagerTest
    """

    def setUp(self):
        """
        initialize a manager.
        """

        self.manager = CalendarManager()
        self.document_content = self.manager.get_document(
            category="2", output="true date or not",
            dtstart=1434359519, dtend=1434705119
        )
        self.vevent_content = self.manager.get_vevent(
            self.document_content
        )
        print("type and value of vevent_content", self.vevent_content)

    def tearDown(self):
        pass

    def clean(self):
        # call of remove_by_source without param to clean DB
        self.manager.remove_by_source()

    def get_calendar(self):
        return self.manager.find(ids=self.ids)


class CalendarTest(CalendarManagerTest):

    def test_get_document_properties(self):
        self.manager._get_document_properties(
            self.document_content
        )

    def test_get_vevent_properties(self):
        self.manager._get_vevent_properties(
            self.vevent_content
        )

    # TODO 4-01-2017
    #def test_put(self):
        #self.clean()

    #    vevents = []
    #    vevents.append(self.vevent_content)

    #    self.manager.put(
    #        vevents
    #    )

    # TODO 4-01-2017
    #def test_get_by_uids(self):

    #    result = self.manager.get_by_uids(
    #        uids=[
    #            "a41ce502-91f6-44d6-bc5d-027582af8e58",
    #            "4303402a-d395-4fe3-abe9-0c8b156030dd"
    #        ], with_count=True
    #    )

    def test_values(self):
        query = {'category': '1', 'output': 'test'}
        result = self.manager.values(
            query=query,
            limit=5,
            with_count=False
        )
        print("values result", result)

    def test_remove(self):

        # remove method ok, bad uids are not considered
        self.manager.remove(uids=[
            "da2a1474-9dbe-424e-a05f-dbd771d5dafb",
            "2455d885-d828-4934-8865-597512c3f09b",
            "2455d885-d828-4934-8865-597512c3f09c"
        ])

if __name__ == '__main__':
    main()
