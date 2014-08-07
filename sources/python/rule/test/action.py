#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.rule.action \
    import do_action, _GLOBAL_ACTIONS, ActionError


def test_action(event, **kwargs):
    result = kwargs.copy()
    if 'count' not in result:
        result['count'] = 0
    result['count'] += 1
    return result


class TestException(Exception):
    pass


def test_action_error(event, **kwargs):
    raise TestException()


class ActionTest(TestCase):

    def setUp(self):
        self.action = {'name': 'test.action.test_action'}

    def test_do_uncached_action(self):

        self.assertEqual(len(_GLOBAL_ACTIONS), 0)

        do_action(action=self.action, event=None, cached_action=False)

        self.assertEqual(len(_GLOBAL_ACTIONS), 0)

    def test_cached_action(self):

        self.assertEqual(len(_GLOBAL_ACTIONS), 0)

        do_action(action=self.action, event=None, cached_action=False)

        self.assertEqual(len(_GLOBAL_ACTIONS), 1)

    def test_unamed_action(self):

        action = {}

        error = False

        try:
            do_action(action=action, event=None)
        except ActionError:
            error = True

        self.assertTrue(error)

    def test_unknown_action(self):

        action = {'name': 'plop'}

        error = False

        try:
            do_action(action=action, event=None)
        except ActionError:
            error = True

        self.assertTrue(error)

    def test_action_error(self):

        action = {'name': 'test_action_error'}

        error = False

        try:
            do_action(action=action, event=None)
        except TestException:
            error = True

        self.assertTrue(error)

if __name__ == '__main__':
    main()
