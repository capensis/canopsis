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

from unittest import main, TestCase

from canopsis.pbehavior.manager import PBehaviorManager

from icalendar import Event

from time import time


class PBehaviorManagerTest(TestCase):
    """Test PBehaviorManager.
    """

    def setUp(self):

        # create a new PBehaviorManager
        self.manager = PBehaviorManager(data_scope='test_pbehavior')

    def tearDown(self):
        # drop behaviors
        self.manager.remove()


class GetDocumentProperties(PBehaviorManagerTest):
    """Test method _get_document_properties.
    """

    def test_empty(self):
        """Test with an empty document.
        """

        document = {}
        result = self.manager._get_document_properties(document=document)
        self.assertEqual(result, {PBehaviorManager.BEHAVIORS: []})

    def test(self):
        """Test with pbehaviors.
        """

        document = {PBehaviorManager.BEHAVIORS: ['test']}
        result = self.manager._get_document_properties(document=document)
        self.assertEqual(result, document)


class GetVeventProperties(PBehaviorManagerTest):
    """Test method _get_vevent_properties.
    """

    def test_empty(self):
        """Test with an empty vevent.
        """

        vevent = Event()
        result = self.manager._get_vevent_properties(vevent=vevent)
        self.assertEqual(result, {PBehaviorManager.BEHAVIORS: []})

    def test(self):
        """Test with a vevent containing pbehaviors.
        """

        iCalvevent = 'BEGIN:VEVENT\n{0}:["test"]\nEND:VEVENT'.format(
            PBehaviorManager.BEHAVIOR_TYPE
        )
        vevent = Event.from_ical(iCalvevent)
        result = self.manager._get_vevent_properties(vevent=vevent)
        self.assertEqual(result, {PBehaviorManager.BEHAVIORS: ['test']})


class GetQuery(PBehaviorManagerTest):
    """Test method get_query.
    """

    def test_empty(self):

        result = PBehaviorManager.get_query(behaviors=[])
        self.assertEqual(result, {PBehaviorManager.BEHAVIORS: []})

    def test(self):

        behaviors = ['test']
        result = PBehaviorManager.get_query(behaviors=behaviors)
        self.assertEqual(result, {PBehaviorManager.BEHAVIORS: behaviors})


class GetEnding(PBehaviorManagerTest):
    """Test method getending.
    """

    def setUp(self):

        super(GetEnding, self).setUp()

        self.source = 'test'
        self.behaviors = ['behavior']
        self.document = PBehaviorManager.get_document(
            source=self.source,
            behaviors=self.behaviors
        )

    def test_alltime(self):
        """Test getending method where time is everytime.
        """

        # check all time
        self.manager.put(vevents=[self.document])
        ending = self.manager.getending(source=self.source)
        self.assertEqual(
            ending,
            {self.behaviors[0]: self.document[PBehaviorManager.DTEND]}
        )

    def test_alltimes(self):
        """Test getending method where several times exist in everytime.
        """

        future = time() + 10000
        count = 5

        documents = [
            self.document.copy()
            for i in range(count)
        ]
        for i in range(len(documents)):
            documents[i][PBehaviorManager.UID] = str(i)
            documents[i][PBehaviorManager.DTEND] = future - i
        # check all time
        self.manager.put(vevents=documents)
        ending = self.manager.getending(source=self.source)
        self.assertEqual(
            ending,
            {self.behaviors[0]: future}
        )

    def test_wrong_source(self):
        """Test when source does not exist.
        """

        ending = self.manager.getending(source='wrong source')
        self.assertFalse(ending)

    def test_excluded_time(self):
        """Test in a time which does not exist at ts.
        """

        self.document[PBehaviorManager.DTEND] = time()
        self.manager.put(vevents=[self.document])
        ending = self.manager.getending(source=self.source)
        self.assertFalse(ending)

    def test_included_time(self):
        """Test to get ending date when dtstart and dtend are given.
        """

        self.document[PBehaviorManager.DTSTART] = time()
        self.document[PBehaviorManager.DTEND] = time() + 10000

        self.manager.put(vevents=[self.document])
        ending = self.manager.getending(source=self.source)
        self.assertEqual(
            ending, {self.behaviors[0]: self.document[PBehaviorManager.DTEND]}
        )

    def test_rrule(self):
        """Test to get ending date when an rrule.
        """

        self.document[PBehaviorManager.DTSTART] = time()
        self.document[PBehaviorManager.DTEND] = time() + 10000
        self.document[PBehaviorManager.RRULE] = "FREQ=DAILY"
        self.manager.put(vevents=[self.document])
        ending = self.manager.getending(source=self.source)
        self.assertEqual(
            ending, {self.behaviors[0]: self.document[PBehaviorManager.DTEND]}
        )

if __name__ == '__main__':
    main()
