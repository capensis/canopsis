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

from calendar import timegm
from datetime import datetime, timedelta
from json import dumps

from mock import Mock, patch
from unittest import main

from canopsis.pbehavior.process import event_processing, beat_processing,\
    PBEHAVIOR_CREATE, PBEHAVIOR_DELETE
from canopsis.context_graph.manager import ContextGraph
from test_base import BaseTest, MockEngine
import unittest
from canopsis.common import root_path
import xmlrunner

class TestProcess(BaseTest):

    def setUp(self):
        super(TestProcess, self).setUp()

    def test_event_processing(self):
        event = {
            "event_type": "pbehavior",
            "pbehavior_name": "downtime",
            "start": timegm(datetime.utcnow().timetuple()),
            "end": timegm((datetime.utcnow() + timedelta(days=1)).timetuple()),
            "action": PBEHAVIOR_CREATE,
            "connector": "test_connector",
            "connector_name": "test_connector_name",
            "author": "test_author",
            "component": 'test_component',
            "source_type": "resource",
            "resource": "a_resource",
            "action": PBEHAVIOR_CREATE
        }

        query = {
            'name': event['pbehavior_name'],
            'filter': dumps({'_id': ContextGraph.get_id(event)}),
            'tstart': event['start'], 'tstop': event['end'],
            'connector': event['connector'],
            'connector_name': event['connector_name'],
            'author': event['author']
        }

        event_processing(MockEngine(), event, pbm=self.pbm, logger=Mock())
        pbehavior = list(self.pbm.pb_storage.get_elements(query=query))
        self.assertEqual(len(pbehavior), 1)
        self.assertDictContainsSubset(query, pbehavior[0])

        event.update({'action': PBEHAVIOR_DELETE})
        event_processing(MockEngine(), event, pbm=self.pbm, logger=Mock())
        pbehavior = list(self.pbm.pb_storage.get_elements(query=query))
        self.assertEqual(len(pbehavior), 0)

    @patch('canopsis.pbehavior.manager.PBehaviorManager.compute_pbehaviors_filters')
    def test_beat_processing(self, mock_compute):
        beat_processing(MockEngine(), logger=Mock())
        self.assertEqual(mock_compute.call_count, 1)
        # method compute_pbehaviors_filters is tested in method TestManager.test_compute_pbehaviors_filters


if __name__ == '__main__':
    output = root_path + "tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
