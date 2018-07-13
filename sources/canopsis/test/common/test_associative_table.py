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

import logging
import unittest
from canopsis.common import root_path
from canopsis.common.associative_table.associative_table import AssociativeTable
from canopsis.common.associative_table.manager import AssociativeTableManager
from canopsis.common.middleware import Middleware
import xmlrunner

class AssociativeTableTest(unittest.TestCase):
    """Test the associative table module.
    """
    def setUp(self):
        self.logger = logging.getLogger()
        self.logger.setLevel(logging.DEBUG)
        self.storage = Middleware.get_middleware_by_uri(
            'storage-default-testassociativetable://'
        )

        self.assoc_manager = AssociativeTableManager(logger=self.logger,
                                                     collection=self.storage._backend)

    def tearDown(self):
        """Teardown"""
        self.storage.remove_elements()

    def test_associative_table(self):
        """
        Test the associative table object.
        """
        key1 = 'do'
        key2 = 'ré'

        assoc = self.assoc_manager.create('test')
        self.assertTrue(isinstance(assoc, AssociativeTable))
        self.assertEqual(assoc.get(key1), None)

        assoc.set(key1, 'mineur')
        get = assoc.get(key1)
        self.assertEqual(get, 'mineur')

        assoc.set(key2, ['majeur'])
        get = assoc.get(key2)
        self.assertListEqual(assoc.get(key2), ['majeur'])

    def test_associative_table_manager(self):
        masterkey = 'masterkey'

        assoc_table = self.assoc_manager.get(masterkey)
        self.assertEqual(assoc_table, None)

        assoc_table = self.assoc_manager.create(masterkey)
        self.assertDictEqual(assoc_table.get_all(), {})

        assoc_table.set('titi', 'grosminet')
        res = self.assoc_manager.save(assoc_table)
        self.assertTrue(res)

        assoc_table2 = self.assoc_manager.get(masterkey)
        self.assertDictEqual(assoc_table2.get_all(), {'titi': 'grosminet'})

        res = self.assoc_manager.delete(masterkey)
        self.assertTrue(res)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
