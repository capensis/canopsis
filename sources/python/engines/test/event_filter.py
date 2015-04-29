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

from canopsis.common.init import Init

from unittest import TestCase, main

from logging import DEBUG, INFO

from canopsis.engines import DROP
from canopsis.engines.event_filter import engine

# TODO, reset theses tests because they are not accurate and not clean

conf = {'rules': [
        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'connector': 'changeme'},
         'actions': [{'type': 'override',
                  'field': 'connector',
                  'value': 'it_works'},
                 {'type': 'pass'}],
         'name': 'change-connector-name'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'connector': 'nagios'},
         'actions': [{'type': 'pass'}],
         'name': 'check-connector-pass'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'connector': 'collectd'},
         'actions': [{'type': 'drop'}],
         'name': 'check-connector-drop'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'connector': 'priority'},
         'actions': [{'type': 'pass'}],
         'name': 'check-connector-pass2'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'test_field': {'$gt': 1378713357}},
         'actions': [{'type': 'drop'}],
         'name': 'check-gt-drop'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {"tags": {"$in": ["collectd2event"]}},
         'name': 'check-in-default'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'connector': 'nagios'},
         'actions': [{'type': 'pass'}],
         'name': 'chec-connector-pass3'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'connector': 'second_rule'},
         'actions': [{'type': 'pass'}],
         'name': 'chec-connector-pass4'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'connector': 'priority'},
         'actions': [{'type': 'drop'}],
         'name': 'check-connector-drop2'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'test_field': {'$eq': 'Engine'}},
         'actions': [{'type': 'pass'}],
         'name': 'check-eq-pass2'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'tags': {'$in': ["tag2"]}},
         'actions': [{'type': 'remove',
                  'key': 'tags',
                  'element': 'tag2'},
                 {'type': 'pass'}],
         'name': 'change-tag-pass'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'perfdatas': {'$in': ["perf1"]}},
         'actions': [{'type': 'remove',
                  'key': 'perfdatas',
                  'element': 'perf1'},
                 {'type': 'pass'}],
         'name': 'remove-perdata-pass'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'remove_me': "True"},
         'actions': [{'type': 'remove',
                  'key': 'remove_me'},
                 {'type': 'pass'}],
         'name': 'remove-eventfield-pass'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'connector': "add_here"},
         'actions': [{'type': 'override',
                  'field': 'added_field',
                  'value': "this_was_added_at_runtime"},
                 {'type': 'pass'}],
         'name': 'add-eventfield-pass'},

        {'description': 'unit test rule',
         '_id': 'unittestidrule',
         'mfilter': {'test_field': {'$gt': 1378713357}},
         'actions': [{'type': 'drop'}],
         'name': 'check-gt-drop2'},

        {'description': 're route',
         '_id': 'unittestidrule',
         'mfilter': {'connector': 'route_me'},
         'actions': [{'type': 'route',
                  'route': 'new_route_defined'}],
         'name': 're-route'},

        {'description': 'rm & add hostgroup',
         '_id': 'unittestidrule',
         'mfilter': {"hostgroups": {"$in": ["linux mint"]}},
         'actions': [{'type': 'remove',
                      'key': 'hostgroups',
                      'element': 'linux mint'},
                     {'type': 'override',
                      'field': 'hostgroups',
                      'value': 'debian jessie'}],
         'name': 'rm-hostgroup-add-hostgroup-pass'}
        ],
    'priority': 2,
    'default_action': 'drop',
    'configuration': 'white_black_lists'}


class KnownValues(TestCase):

    def setUp(self):
        self.engine = engine(
            **{'name': 'passdropity', 'logging_level': INFO})

        self.engine.next_amqp_queues = ["consolidation"]
        self.engine.drop_event_count = 0
        self.engine.pass_event_count = 0
        self.engine.configuration = conf

    def test_change_field(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        event['connector'] = 'changeme'
        event = self.engine.work(event)
        if 'hostgroups' not in event:
            self.assertEqual("it_works", event['connector'])

    def test_normal_behavior(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        # Test normal behaviors
        event['connector'] = 'nagios'
        self.assertEqual(self.engine.work(event), event)

        event['connector'] = 'collectd'
        self.assertEqual(self.engine.work(event), DROP)

    def test_match_field(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        # second rule matched
        event['connector'] = 'second_rule'
        self.assertEqual(self.engine.work(event), event)

    def test_default_action(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'rk': ''
        }

        self.engine.configuration['default_action'] = 'drop'
        # Test default actions
        self.assertEqual(self.engine.work(event), DROP)

    def test_changed_daction(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        # Change default action
        self.engine.configuration['default_action'] = 'pass'
        event['connector'] = 'default_pass'
        self.assertEqual(self.engine.work(event), event)

    def test_remove_field_list(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        # remove 'tag2'
        event['tags'] = ['tag1', 'tag2', 'tag3']
        event = self.engine.work(event)
        self.assertEqual(event['tags'][1], "tag3")

    def test_remove_field_dict(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        # remove 'perf1'
        event['perfdatas'] = {
            'perf1': 13374242,
            'perf2': 42421337,
            'perf3': 42
        }
        event = self.engine.work(event)
        self.assertEqual(('perf1' in event['perfdatas']), False)

    def test_add_field(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        # add field
        event['connector'] = "add_here"
        event = self.engine.work(event)
        self.assertEqual(event['added_field'], "this_was_added_at_runtime")

    def test_remove_field(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        # remove field
        event['remove_me'] = 'True'
        event = self.engine.work(event)
        self.assertEqual(('remove_me' in event), False)

    def test_change_route(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        event['connector'] = 'route_me'
        event['remove_me'] = 'False'
        event = self.engine.work(event)
        self.assertEqual(self.engine.next_amqp_queues[0], "new_route_defined")

    def test_no_conf(self):

        event = {
            'connector': '',
            'connector_name': '',
            'event_type': '',
            'source_type': '',
            'component': '',
            'tags': [],
            'rk': ''
        }

        # No configuration, default configuration is loaded
        self.engine.configuration = {}
        self.assertEqual(self.engine.work(event), event)


if __name__ == "__main__":
    main()
