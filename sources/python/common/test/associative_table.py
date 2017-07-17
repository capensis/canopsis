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

from unittest import TestCase, main

from canopsis.common.associative_table.associative_table import AssociativeTable
from canopsis.common.associative_table.manager import AssociativeTableManager
from canopsis.middleware.core import Middleware


class AssociativeTableTest(TestCase):
    """Test the associative table module.
    """
    def setUp(self):
        self.storage = Middleware.get_middleware_by_uri(
            'storage-default-testassociativetable://'
        )

        self.assoc_manager = AssociativeTableManager(storage=self.storage)

    def tearDown(self):
        """Teardown"""
        self.storage.remove_elements()

    def test_associative_table(self):
        """
        Test the associative table object.
        """
        key1 = 'do'
        key2 = 'ré'

        assoc = self.assoc_manager.get('test')
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
        self.assertDictEqual(assoc_table.get_all(), {})

        assoc_table.set('titi', 'grosminet')
        res = self.assoc_manager.save(assoc_table)
        self.assertTrue(res['updatedExisting'])
        self.assertEqual(res['n'], 1)

        assoc_table2 = self.assoc_manager.get(masterkey)
        self.assertDictEqual(assoc_table2.get_all(), {'titi': 'grosminet'})

if __name__ == '__main__':
    main()
