#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

import unittest
import sys, os
import logging

sys.path.append(os.path.expanduser('~/opt/amqp2engines/engines/'))

import filters
from cengine import DROP


class KnownValues(unittest.TestCase):
	def setUp(self):
		self.engine = filters.engine(logging_level=logging.DEBUG)
		self.engine.beat()

	def test_01_Init(self):
		self.engine.drop_event_count = 0
		self.engine.pass_event_count = 0
		self.engine.configuration = {
			'rules': [
				{
				 'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'connector': 'changeme'},
				 'actions': [{'type':'override',
					      'field':'connector',
					      'value':'it_works'
					      },
					     {'type':'pass'}],
				 'name': 'change-connector-name'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'connector': 'nagios'},
				 'actions': [{'type':'pass'}],
				 'name': 'check-connector-pass'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'connector': 'collectd'},
				 'actions': [{'type':'drop'}],
				 'name': 'check-connector-drop'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'connector': 'priority'},
				 'actions': [{'type':'pass'}],
				 'name': 'check-connector-pass2'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'test_field': { '$gt': 1378713357 } },
				 'actions': [{'type':'drop'}],
				 'name': 'check-gt-drop'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {"tags": {"$in": ["collectd2event"]} },
				 'name': 'check-in-default'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'connector': 'nagios'},
				 'actions': [{'type':'pass'}],
				 'name': 'chec-connector-pass3'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'connector': 'second_rule'},
				 'actions': [{'type':'pass'}],
				 'name': 'chec-connector-pass4'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'connector': 'priority'},
				 'actions': [{'type':'drop'}],
				 'name': 'check-connector-drop2'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'test_field': { '$eq': 'cengine' } },
				 'actions': [{'type':'pass'}],
				 'name': 'check-eq-pass2'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'tags': { '$in': ["tag2"]} },
				 'actions': [{'type':'remove',
					      'key': 'tags',
					      'element': 'tag2'},
					     {'type':'pass'}],
				 'name': 'change-tag-pass'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'perfdatas': { '$in': ["perf1"]} },
				 'actions': [{'type':'remove',
					      'key': 'perfdatas',
					      'element': 'perf1'},
					     {'type':'pass'}],
				 'name': 'remove-perdata-pass'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'remove_me': "True"},
				 'actions': [{'type':'remove',
					      'key': 'remove_me'},
					     {'type':'pass'}],
				 'name': 'remove-eventfield-pass'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'connector': "add_here"},
				 'actions': [{'type':'override',
					      'field': 'added_field',
					      'value': "this_was_added_at_runtime"},
					     {'type':'pass'}],
				 'name': 'add-eventfield-pass'},

				{'description': 'unit test rule',
				 '_id': 'unittestidrule',
				 'mfilter': {'test_field': { '$gt': 1378713357 } },
				 'actions': [{'type':'drop'}],
				 'name': 'check-gt-drop2'},
			],
			'priority' : 2,
			'default_action': 'drop',
			'configuration': 'white_black_lists',
		}

		event = {'connector': '',
			 'connector_name': '',
			 'event_type': '',
			 'source_type': '',
			 'component': '',
			 'tags': [],
			 'rk': ''
			}


		event['connector'] = 'changeme'
		event = self.engine.work(event)

		self.assertEqual("it_works", event['connector'])

		# Test normal behaviors
		event['connector'] = 'nagios'

		self.assertEqual(self.engine.work(event), event)

		event['connector'] = 'collectd'
		self.assertEqual(self.engine.work(event), DROP)


		# second rule matched
		event['connector'] = 'second_rule'
		self.assertEqual(self.engine.work(event), event)

		# Test default actions
		event['connector'] = 'default_drop'
		self.assertEqual(self.engine.work(event), DROP)

		# Change default action
		self.engine.configuration['default_action'] = 'pass'
		event['connector'] = 'default_pass'
		self.assertEqual(self.engine.work(event), event)

		# rule priority validation sorted is the same used in beat method in the engine
		event['connector'] = 'priority'
		self.assertEqual(self.engine.work(event), event)

		# remove 'tag2'
		event['connector'] = ""
		event['tags'] = ['tag1', 'tag2', 'tag3']
		event = self.engine.work(event)
		self.assertEqual(event['tags'][1], "tag3")

		# remove 'perf1'
		event['perfdatas'] = {'perf1': 13374242,
				      'perf2': 42421337,
				      'perf3': 42}
		event = self.engine.work(event)
		self.assertEqual(('perf1' in event['perfdatas']), False)

		# add field
		event['connector'] = "add_here"
		event = self.engine.work(event)
		self.assertEqual(event['added_field'], "this_was_added_at_runtime")

		# remove field
		event['connector'] = ""
		event['remove_me'] = 'True'
		event['tags'] = []
		del event['perfdatas']
		event = self.engine.work(event)
		self.assertEqual(('remove_me' in event), False)

		# No configuration, default configuration is loaded
		self.engine.configuration = {}
		self.assertEqual(self.engine.work(event), event)


if __name__ == "__main__":
	unittest.main()
