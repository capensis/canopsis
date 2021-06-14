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
#from mock import MagicMock
#from time import time
#from logging import ERROR
#from canopsis.engines.cleaner import engine


#class DataCleanerTest(TestCase):

#    def setUp(self):
#        self.engine = engine(logging_level=ERROR)
#
#    def test_collection_used(self):
#
#        self.assertIn('events', self.engine.clean_collection)
#        self.assertIn('events_log', self.engine.clean_collection)
#
#    def test_db_connection(self):
#
#        # event collection
#        collection_events = self.engine.clean_collection['events']
#        self.assertIsNotNone(collection_events)
#
#        # events log collection
#        collection_events_log = self.engine.clean_collection['events_log']
#        self.assertIsNotNone(collection_events_log)
#
#        # object collection
#        self.assertIsNotNone(self.engine.object)
#
#    def test_get_configuration(self):
#        configuration = self.engine.get_configuration()
#        self.assertIsNotNone(configuration)
#        # Is engine configuration loaded properly
#        self.assertEqual(
#            configuration['crecord_type'],
#            'datacleaner'
#        )
#
#    def test_retention_date(self):
#        # Canopsis install default configuration
#        # This is a minimal configuration
#
#        one_year = 3600 * 24 * 365
#
#        self.engine.get_configuration = MagicMock(return_value={
#            'retention_duration': 0,
#            'use_secure_delay': True
#        })
#
#        # Expects your processor is fast enough
#        security_duration_limit = int(time() - one_year)
#        retention_date = self.engine.get_retention_date()
#        self.assertEqual(retention_date, security_duration_limit)
#
#        # Tests with no secure delay
#        more_than_one_year = one_year + 10
#        self.engine.get_configuration = MagicMock(return_value={
#            'retention_duration': more_than_one_year,  # greater than one year
#            'use_secure_delay': True
#        })
#        compare_duration = int(time() - more_than_one_year)
#        # security_duration_limit = int(time() - 3600 * 24 * 365)
#        retention_date = self.engine.get_retention_date()
#        self.assertEqual(retention_date, compare_duration)
#
#        # Tests with no secure delay /!\
#        less_than_one_year = one_year - 10
#        self.engine.get_configuration = MagicMock(return_value={
#            'retention_duration': less_than_one_year,  # greater than one year
#            'use_secure_delay': True
#        })
#        security_duration_limit = int(time() - 3600 * 24 * 365)
#        retention_date = self.engine.get_retention_date()
#        self.assertEqual(retention_date, security_duration_limit)
#
#        # Tests with no secure delay off
#        more_than_one_year = one_year + 10
#        self.engine.get_configuration = MagicMock(return_value={
#            'retention_duration': more_than_one_year,  # greater than one year
#            'use_secure_delay': False
#        })
#        compare_duration = int(time() - more_than_one_year)
#        retention_date = self.engine.get_retention_date()
#        self.assertEqual(retention_date, compare_duration)
#
#        # Tests with no secure delay off /!\
#        less_than_one_year = one_year - 10
#        self.engine.get_configuration = MagicMock(return_value={
#            'retention_duration': less_than_one_year,  # greater than one year
#            'use_secure_delay': False
#        })
#        compare_duration = int(time() - less_than_one_year)
#        retention_date = self.engine.get_retention_date()
#        self.assertEqual(retention_date, compare_duration)
#
#    def test_beat_effective_clean(self):
#
#        one_year = 3600 * 24 * 365
#
#        self.engine.get_configuration = MagicMock(return_value={
#            'retention_duration': one_year,
#            'use_secure_delay': True
#        })
#
#        def mock_remove_method(query):
#
#            # Is test done on timstamp key
#            self.assertIn('timestamp', query)
#            # There should not be any other field in the query
#            self.assertEqual(len(query), 1)
#
#            # Query have to be a compared timestamp lower or equal to value
#            self.assertIn('$lte', query['timestamp'])
#
#            # There should no have any other parameter for timestamp
#            self.assertEqual(len(query['timestamp']), 1)
#
#            # delay should be a number of seconds
#            self.assertTrue(
#                query['timestamp']['$lte'] == int(time() - one_year)
#            )
#
#        for collection in ['events', 'events_log']:
#            collection_backend = self.engine.clean_collection[collection]
#            collection_backend.remove = mock_remove_method
#
#            self.engine.beat()
#
#
#if __name__ == "__main__":
#    main()
