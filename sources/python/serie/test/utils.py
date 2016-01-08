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

from unittest import TestCase, main

from canopsis.serie.utils import build_filter_from_regex, NAME_BY_KEYS


class TestSerieUtils(TestCase):
    def test_build_filter_from_regex_co(self):
        mfilter = build_filter_from_regex('co:test')
        expected = {NAME_BY_KEYS['co:']: {'$regex': 'test'}}

        self.assertEqual(mfilter, expected)

    def test_build_filter_from_regex_re(self):
        mfilter = build_filter_from_regex('re:test')
        expected = {NAME_BY_KEYS['re:']: {'$regex': 'test'}}

        self.assertEqual(mfilter, expected)

    def test_build_filter_from_regex_me(self):
        mfilter = build_filter_from_regex('me:test')
        expected = {NAME_BY_KEYS['me:']: {'$regex': 'test'}}

        self.assertEqual(mfilter, expected)

    def test_build_filter_from_regex_co_re_me(self):
        mfilter = build_filter_from_regex('co:test re:test me:test')
        expected = {'$and': [
            {NAME_BY_KEYS['co:']: {'$regex': 'test'}},
            {NAME_BY_KEYS['re:']: {'$regex': 'test'}},
            {NAME_BY_KEYS['me:']: {'$regex': 'test'}}
        ]}

        self.assertEqual(mfilter, expected)

    def test_build_filter_from_regex_co_co(self):
        mfilter = build_filter_from_regex('co:test1 co:test2')
        expected = {'$or': [
            {NAME_BY_KEYS['co:']: {'$regex': 'test1'}},
            {NAME_BY_KEYS['co:']: {'$regex': 'test2'}}
        ]}

        self.assertEqual(mfilter, expected)

    def test_build_filter_from_regex_co_co_re(self):
        mfilter = build_filter_from_regex('co:test1 co:test2 re:test')
        expected = {'$and': [
            {'$or': [
                {NAME_BY_KEYS['co:']: {'$regex': 'test1'}},
                {NAME_BY_KEYS['co:']: {'$regex': 'test2'}}
            ]},
            {NAME_BY_KEYS['re:']: {'$regex': 'test'}}
        ]}

        self.assertEqual(mfilter, expected)

    def test_build_filter_from_regex_co_co_re_re(self):
        mfilter = build_filter_from_regex(
            'co:test1 co:test2 re:test1 re:test2'
        )
        expected = {'$and': [
            {'$or': [
                {NAME_BY_KEYS['co:']: {'$regex': 'test1'}},
                {NAME_BY_KEYS['co:']: {'$regex': 'test2'}}
            ]},
            {'$or': [
                {NAME_BY_KEYS['re:']: {'$regex': 'test1'}},
                {NAME_BY_KEYS['re:']: {'$regex': 'test2'}}
            ]}
        ]}

        self.assertEqual(mfilter, expected)


if __name__ == '__main__':
    main()
