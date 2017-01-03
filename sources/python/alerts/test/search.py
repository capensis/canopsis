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

from unittest import TestCase, main
from re import compile as compile_

from canopsis.alerts.search.interpreter import interpret


class TestSearch(TestCase):
    grammar_file = '/opt/canopsis/etc/alerts/search/grammar.bnf'

    def test_interpret_base(self):
        cases = [
            {
                'search': 'a = "a"',
                'expected_filter': {'a': {'$eq': 'a'}}
            },
            {
                'search': 'a = "a" AND NOT b = "b"',
                'expected_filter': {
                    '$and': [
                        {'a': {'$eq': 'a'}}, {'b': {'$not': {'$eq': 'b'}}}
                    ]
                }
            },
            {
                'search': 'ALL a = "a"',
                'expected_scope': 'all',
                'expected_filter': {'a': {'$eq': 'a'}}
            }
        ]

        for case in cases:
            scope, filter_ = interpret(
                case['search'],
                grammar_file=self.grammar_file
            )

            self.assertEqual(scope, case.get('expected_scope', 'this'))
            self.assertEqual(filter_, case['expected_filter'])

    def test_interpret_types(self):
        cases = [
            {
                'search': 'a = "a"',
                'expected_filter': {'a': {'$eq': 'a'}}
            },
            {
                'search': "a = 'a'",
                'expected_filter': {'a': {'$eq': 'a'}}
            },
            {
                'search': 'a = 1',
                'expected_filter': {'a': {'$eq': 1}}
            },
            {
                'search': 'a = +1',
                'expected_filter': {'a': {'$eq': 1}}
            },
            {
                'search': 'a = -1',
                'expected_filter': {'a': {'$eq': -1}}
            },
            {
                'search': 'a = 1.2345',
                'expected_filter': {'a': {'$eq': 1.2345}}
            },
            {
                'search': 'a = +1.2345',
                'expected_filter': {'a': {'$eq': 1.2345}}
            },
            {
                'search': 'a = -1.2345',
                'expected_filter': {'a': {'$eq': -1.2345}}
            },
            {
                'search': 'a = TRUE',
                'expected_filter': {'a': {'$eq': True}}
            },
            {
                'search': 'a = FALSE',
                'expected_filter': {'a': {'$eq': False}}
            },
            {
                'search': 'a = NULL',
                'expected_filter': {'a': {'$eq': None}}
            }
        ]

        for case in cases:
            scope, filter_ = interpret(
                case['search'],
                grammar_file=self.grammar_file
            )

            self.assertEqual(filter_, case['expected_filter'])

    def test_interpret_compop(self):
        cases = [
            {
                'search': 'a <= 1',
                'expected_filter': {'a': {'$lte': 1}}
            },
            {
                'search': 'a < 1',
                'expected_filter': {'a': {'$lt': 1}}
            },
            {
                'search': 'a = 1',
                'expected_filter': {'a': {'$eq': 1}}
            },
            {
                'search': 'a != 1',
                'expected_filter': {'a': {'$ne': 1}}
            },
            {
                'search': 'a >= 1',
                'expected_filter': {'a': {'$gte': 1}}
            },
            {
                'search': 'a > 1',
                'expected_filter': {'a': {'$gt': 1}}
            },
            {
                'search': 'a IN 1',
                'expected_filter': {'a': {'$in': 1}}
            },
            {
                'search': 'a NIN 1',
                'expected_filter': {'a': {'$nin': 1}}
            },
            {
                'search': 'a LIKE 1',
                'expected_filter': {'a': {'$regex': 1}}
            },
            {
                'search': 'NOT a = 1',
                'expected_filter': {'a': {'$not': {'$eq': 1}}}
            },
            {
                'search': 'NOT a LIKE "1"',
                'expected_filter': {'a': {'$not': compile_(u'1')}}
            }
        ]

        for case in cases:
            scope, filter_ = interpret(
                case['search'],
                grammar_file=self.grammar_file
            )

            self.assertEqual(filter_, case['expected_filter'])

    def test_interpret_condop(self):
        cases = [
            {
                'search': 'a = 1 AND b = 2',
                'expected_filter': {
                    '$and': [{'a': {'$eq': 1}}, {'b': {'$eq': 2}}]
                }
            },
            {
                'search': 'a = 1 OR b = 2',
                'expected_filter': {
                    '$or': [{'a': {'$eq': 1}}, {'b': {'$eq': 2}}]
                }
            },
            {
                'search': 'a = 1 AND b = 2 OR c = 3',
                'expected_filter': {
                    '$or': [
                        {
                            '$and': [
                                {'a': {'$eq': 1}},
                                {'b': {'$eq': 2}}
                            ]
                        },
                        {'c': {'$eq': 3}}
                    ]
                }
            },
            {
                'search': 'a = 1 OR b = 2 AND c = 3',
                'expected_filter': {
                    '$and': [
                        {
                            '$or': [
                                {'a': {'$eq': 1}}, {'b': {'$eq': 2}}
                            ]
                        },
                        {'c': {'$eq': 3}}
                    ]
                }
            },
            {
                'search': (
                    'ALL a LIKE 1 OR b <= 2 AND c = 3 AND d != 4 OR '
                    'f > "five"'
                ),
                'expected_scope': 'all',

                'expected_filter': {
                    '$or': [
                        {
                            '$and': [
                                {
                                    '$and': [
                                        {
                                            '$or': [
                                                {'a': {'$regex': 1}},
                                                {'b': {'$lte': 2}}
                                            ]
                                        },
                                        {'c': {'$eq': 3}}
                                    ]
                                },
                                {'d': {'$ne': 4}}
                            ]
                        },
                        {'f': {'$gt': 'five'}}
                    ]
                }

            }
        ]

        for case in cases:
            scope, filter_ = interpret(
                case['search'],
                grammar_file=self.grammar_file
            )

            self.assertEqual(scope, case.get('expected_scope', 'this'))
            self.assertEqual(filter_, case['expected_filter'])


if __name__ == '__main__':
    main()
