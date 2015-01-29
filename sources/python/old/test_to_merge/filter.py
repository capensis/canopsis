# !/usr/bin/env python
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

import unittest
from canopsis.old.cfilter import Filter


class KnownValues(unittest.TestCase):

    def setUp(self):

        self.filter = Filter()

    def test_add_filter_by_id(self):

        print('It should generate an empty filter')
        cfilter = self.filter.add_filter_by_id([], True)
        self.assertEqual(cfilter, {})

        print('It should generate a filter with a single key equal test')
        cfilter = self.filter.add_filter_by_id(['keyone'], True)
        self.assertEqual(cfilter, {'_id': {'$eq': 'keyone'}})

        print('It should generate a filter with a single key not equal test')
        cfilter = self.filter.add_filter_by_id(['keyone'], False)
        self.assertEqual(cfilter, {'_id': {'$ne': 'keyone'}})

        print(
            'It should generate a filter with a _id key with nested' +
            ' dict with noop key and given dict as value')
        value = ['keyone', 'keytwo']
        cfilter = self.filter.add_filter_by_id(value, True)
        test = {'_id': {'$in': ['keyone', 'keytwo']}}
        self.assertEqual(cfilter, test)

    def test_make_mfilter(self):

        mfilter = {'key': 'value'}
        value = ['key1', 'key2']

        print('It should return an empty filter')
        cfilter = self.filter.make_filter()
        self.assertEqual(cfilter, {})

        print('It should return given filter')
        cfilter = self.filter.make_filter(mfilter=mfilter)
        self.assertIn('key', cfilter)
        self.assertEqual(cfilter['key'], 'value')
        self.assertEqual(len(cfilter.keys()), 1)

        print('It should return an include filter')
        cfilter = self.filter.make_filter(includes=value)
        self.assertIn('_id', cfilter)
        self.assertEqual(len(cfilter.keys()), 1)
        self.assertIn('$in', cfilter['_id'])
        self.assertEqual(len(cfilter['_id'].keys()), 1)
        self.assertEqual(cfilter['_id']['$in'], value)

        print('It should return an exclude filter')
        cfilter = self.filter.make_filter(excludes=value)
        self.assertIn('_id', cfilter)
        self.assertEqual(len(cfilter.keys()), 1)
        self.assertIn('$nin', cfilter['_id'])
        self.assertEqual(len(cfilter['_id'].keys()), 1)
        self.assertEqual(cfilter['_id']['$nin'], value)

        print('It should return a compund filter between mfilter and excludes')
        cfilter = self.filter.make_filter(mfilter=mfilter, excludes=value)
        self.assertIn('$and', cfilter)
        self.assertEqual(len(cfilter['$and']), 2)
        self.assertEqual(cfilter['$and'][0], mfilter)
        self.assertEqual(cfilter['$and'][1], {'_id': {'$nin': value}})

        print('It should return a compund filter between mfilter and includes')
        cfilter = self.filter.make_filter(mfilter=mfilter, includes=value)
        self.assertIn('$or', cfilter)
        self.assertEqual(len(cfilter['$or']), 2)
        self.assertEqual(cfilter, {'$or': [mfilter, {'_id': {'$in': value}}]})

        print('It should return a compund filter between excludes includes')
        cfilter = self.filter.make_filter(
            excludes=['singlekey'],
            includes=value)
        test = {'$and': [
            {'_id': {'$in': value}},
            {'_id': {'$ne': 'singlekey'}}]}
        self.assertEqual(cfilter, test)


if __name__ == "__main__":
    unittest.main(verbosity=2)
