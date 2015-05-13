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
from datetime import datetime
from dateutil.relativedelta import relativedelta


class PBehaviorManagerTest(TestCase):
    """Test PBehaviorManager.
    """

    def setUp(self):
        # create a new PBehaviorManager
        self.manager = PBehaviorManager(data_scope='test_pbehavior')
        # UTs check period with(out) behaviors on periods in one hour.
        self.now = datetime.now()
        # use same rrule which is checked every day
        rrule = "FREQ=DAILY"
        # use same duration which lasts 1 jour
        duration = "PT1H"  # 1 hour duration
        # construct events
        self.events = []
        # with count events
        self.count = 5
        # and (60 mn / (self.count - 1)) different minutes between date
        minutes = int(60 / self.count)
        # get an ical template with dtstart, rrule andd duration
        ical = "BEGIN:VEVENT\nDTSTART:{0}\nRRULE:{1}\nDURATION:{2}\nEND:VEVENT"
        # construct (self.count) events in order to leave event period once
        for i in range(self.count):
            rd = relativedelta(minutes=minutes * (i + 1))
            dtstart = self.now + rd
            dtstart = "{0}{1}{2}".format(
                dtstart.year, dtstart.month, dtstart.day
            )
            event_ical = ical.format(dtstart, rrule, duration)
            event = Event.from_ical(event_ical)
            self.events.append(event)

        # construct behaviors such as ["1", "2", ..., "self.count - 1"]
        self.behaviors = [str(i) for i in range(self.count)]
        # construct documents such as:
        # {id: i, values: [{period: events[j], behaviors: behaviors[j]}]}
        self.values = [
            {
                PBehaviorManager.ID: self.behaviors[i],
                PBehaviorManager.VALUES: [
                    {
                        PBehaviorManager.PERIOD: self.events[j].to_ical(),
                        PBehaviorManager.BEHAVIORS: self.behaviors[:j]
                    } for j in range(self.count)
                ]
            } for i in range(self.count)
        ]

        for document in self.documents:
            self.manager.put(
                entity_id=document[PBehaviorManager.ID], document=document
            )

    def tearDown(self):
        # drop behaviors
        self.manager.remove(entity_ids=self.behaviors)

    def test_CRUD(self):
        """Test CRUD methods.
        """

        # check to retrieve one document
        documents = self.manager.get(self.behaviors[0])
        self.assertIsInstance(documents, dict)

        # check to retrieve several documents
        documents = self.manager.get(self.behaviors)
        self.assertEqual(len(documents), self.count)

        # remove one document
        self.manager.remove(entity_ids=self.behaviors[0])

        # check to retrieve one document which does not exist
        documents = self.manager.get(self.behaviors[0])
        self.assertIsNone(documents)

        # check to retrieve several documents
        documents = self.manager.get(self.behaviors)
        self.assertEqual(len(documents), self.count - 1)

        # remove all documents except one
        self.manager.remove(entity_ids=self.behaviors[:-1])

        # check to retrieve several documents
        documents = self.manager.get(self.behaviors)
        self.assertEqual(len(documents), 1)

    def test__get_ending(self):
        """Test _get_ending method.
        """

        # remove one minute from self.now in order to not match with any period
        rd = relativedelta(minutes=1)
        dtts = self.now - rd

        # check each behavior
        for behavior in self.behaviors:
            # check to get endings outside periods
            endings = self.manager._get_ending(
                behaviors={behavior}, documents=self.documents, dtts=dtts
            )
            self.assertFalse(endings)
            # check to get endings inside periods
            endings = self.manager._get_ending(
                behaviors={behavior}, documents=self.documents,
                dtts=self.now + rd
            )
            self.assertTrue(endings)
        # check set of behaviors
        for i in range(self.count):
            # get first i behaviors
            sub_behaviors = set(self.behaviors[:i])
            # check to get endings outside periods
            endings = self.manager._get_ending(
                behaviors=sub_behaviors, documents=self.documents, dtts=dtts
            )
            self.assertFalse(endings)
            # check to get endings inside periods
            endings = self.manager._get_ending(
                behaviors=sub_behaviors, documents=self.documents,
                dtts=self.now
            )
            self.assertEqual(len(endings), len(sub_behaviors))

    def test_get_ending(self):
        """Test get_ending method.
        """

    def test_whois(self):
        """Test whois method.
        """

if __name__ == '__main__':
    main()
