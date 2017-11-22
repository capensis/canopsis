#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

from json import loads
import logging
from unittest import TestCase, main

from canopsis.backbone.event import Event


class BackboneEventTest(TestCase):

    def setUp(self):
        self.logger = logging.getLogger('alarms')

        self.event = Event(
            connector='cleese',
            connector_name='chapman',
            component='idle',
            resource='palin',
            source_type='resource',
            output='How Not to be seen',
        )

    def test_is_valid(self):
        self.assertTrue(self.event.is_valid())

        bad_event = Event(
            connector='cleese',
            connector_name='chapman',
            component='idle',
            source_type='Dinsdale'
        )
        self.assertFalse(bad_event.is_valid())

        bad_event = Event(
            connector='cleese',
            connector_name='chapman',
            component='idle',
            source_type='resource'
        )
        self.assertFalse(bad_event.is_valid())

    def test_ensure_utf8_format(self):
        try:
            self.event.ensure_utf8_format()
        except Exception:
            self.fail('ensure_utf8_format has failed but it shouldn\'t')

    def test_to_json(self):
        res = self.event.to_json()
        self.assertIsInstance(res, str)

        json = loads(res)
        del json['timestamp']
        self.assertDictEqual(json, {
            'connector': 'cleese',
            'resource': 'palin',
            'event_type': 'check',
            'component': 'idle',
            'alarm': None,
            'context': None,
            'display_name': '',
            'source_type': 'resource',
            'state': 0,
            'connector_name': 'chapman',
            'output': 'How Not to be seen',
            'long_output': '',
            'perf_data': '',
            'perf_data_array': [],
            'more': {},
        })

if __name__ == '__main__':
    main()
