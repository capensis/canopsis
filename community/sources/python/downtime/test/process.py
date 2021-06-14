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

# TODO 4-01-2017

#from unittest import TestCase, main
#
#from canopsis.pbehavior.manager import PBehaviorManager
#from canopsis.event.manager import Event
#from canopsis.context.manager import Context
#
#from canopsis.downtime.process import event_processing, beat_processing
#from canopsis.downtime.process import DOWNTIME
#
#
#class DowntimeProcessingTest(TestCase):
#
#    def setUp(self):
#
#        self.downtimes = PBehaviorManager(data_scope='test_pbehavior')
#        self.events = Event(data_scope='test_events')
#        self.context = Context(data_scope='test_context')
#
#    def tearDown(self):
#
#        self.downtimes.remove()
#        self.events.remove()
#        self.context.remove()
#
#
#class EventProcessingTest(DowntimeProcessingTest):
#
#    def setUp(self):
#
#        super(EventProcessingTest, self).setUp()
#
#        self.test_event = {
#            'connector': 'unittest',
#            'connector_name': self.__class__.__name__,
#            'event_type': 'check',
#            'source_type': 'resource',
#            'component': 'component0',
#            'resource': 'resource0'
#        }
#
#        self.test_rk = self.events.get_rk(self.test_event)
#
#    def _process(self):
#
#        event = event_processing(
#            self, self.test_event,
#            downtimes=self.downtimes,
#            events=self.events
#        )
#        dbevent = self.events.get(self.test_rk)
#
#        return event, dbevent
#
#    def test_without_downtime(self):
#
#        event, dbevent = self._process()
#
#        self.assertIsNotNone(dbevent)
#        self.assertFalse(event[DOWNTIME])
#        self.assertFalse(dbevent[DOWNTIME])
#
#    def test_with_downtime(self):
#        # TODO: add downtime
#
#        event, dbevent = self._process()
#
#        self.assertIsNotNone(dbevent)
#        self.assertTrue(event[DOWNTIME])
#        self.assertTrue(dbevent[DOWNTIME])
#
#
#class BeatProcessingTest(DowntimeProcessingTest):
#
#    def _process(self):
#
#        beat_processing(
#            self,
#            downtimes=self.downtimes,
#            events=self.events,
#            context=self.context
#        )
#
#        return self.events.find(query={DOWNTIME: True})
#
#    def test_no_downtime(self):
#
#        result = self._process()
#
#        self.assertFalse(result)
#
#    def test_with_downtimes(self):
#        # TODO: add downtime
#
#        result = self._process()
#
#        self.assertTrue(result)
#
#
#if __name__ == '__main__':
#    main()
