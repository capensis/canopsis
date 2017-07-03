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

from canopsis.pbehavior.process import (event_processing, beat_processing,
                                        get_entity_id,
                                        PBEHAVIOR_CREATE, PBEHAVIOR_DELETE)

from base import BaseTest


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
            "resource": "a_resource"
        }

        query = {
            'name': event['pbehavior_name'],
            'filter': dumps({'entity_id': get_entity_id(event)}),
            'tstart': event['start'], 'tstop': event['end'],
            'connector': event['connector'],
            'connector_name': event['connector_name'],
            'author': event['author']
        }

        event_processing(None, event, pbm=self.pbm, logger=Mock())
        pbehavior = list(self.pbm.pbehavior_storage.get_elements(query=query))
        self.assertEqual(len(pbehavior), 1)
        self.assertDictContainsSubset(query, pbehavior[0])

        event.update({'action': PBEHAVIOR_DELETE})
        event_processing(None, event, pbm=self.pbm, logger=Mock())
        pbehavior = list(self.pbm.pbehavior_storage.get_elements(query=query))
        self.assertEqual(len(pbehavior), 0)

    @patch('canopsis.pbehavior.manager.PBehaviorManager.compute_pbehaviors_filters')
    def test_beat_processing(self, mock_compute):
        beat_processing(None, logger=Mock())
        self.assertEqual(mock_compute.call_count, 1)
        # method compute_pbehaviors_filters is tested in method TestManager.test_compute_pbehaviors_filters

    def test_get_entity_id(self):
        event = {
            'source_type': 'component',
            'component': 't_component',
            'connector': 't_connector',
            'connector_name': 't_connector_name',
            'resource': 't_resource',
            'selector': 't_selector'
        }
        entity_id = get_entity_id(event)
        self.assertEqual(
            '/component/t_connector/t_connector_name/t_component',
            entity_id
        )

        event.update({'source_type': 'resource'})
        entity_id = get_entity_id(event)
        self.assertEqual(
            '/resource/t_connector/t_connector_name/t_component/t_resource',
            entity_id
        )

        event.update({'source_type': 'selector'})
        entity_id = get_entity_id(event)
        self.assertEqual(
            '/selector/t_connector/t_connector_name/t_selector',
            entity_id
        )


if __name__ == '__main__':
    main()
