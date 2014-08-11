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


from canopsis.rule import apply_rule
from canopsis.rule.condition import CONDITION_FIELD
from canopsis.rule.action import ACTIONS_FIELD


def test_action(event, **kwargs):
    result = kwargs.copy()
    if 'count' not in result:
        result['count'] = 0
    result['count'] += 1
    return result


class RuleTest(TestCase):
    def setUp(self):
        self.action = {'name': 'test.test_action'}
        self.true_condition = {'name': 'event'}
        self.false_condition = {'name': 'not_event'}

    def test_empty(self):
        rule = {}
        result = apply_rule(rule=rule, ctx=None, event={})

        self.assertEqual(len(result), 0)

    def test_false_rule(self):

        rule = {
            CONDITION_FIELD: {'name': 'event'},
            ACTIONS_FIELD: [{'name': 'test.test_action'}]
        }

        result = apply_rule(rule=rule, ctx=None, event={'name': 'event'})
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0]['count'], 1)
        self.assertNotEqual(result, self.rule[ACTIONS_FIELD])

    def test_true_rule(self):

        rule = {
            CONDITION_FIELD: {'name': 'event'},
            ACTIONS_FIELD: [{'name': 'test.test_action'}]
        }

        event = {'name': 'not_event'}

        result = apply_rule(rule=rule, ctx=None, event=event)
        self.assertEqual(len(result), 0)

        rule[CONDITION_FIELD] = {'name': 'event'}
        result = apply_rule(rule=rule, ctx=None, event=event)

        self.assertEqual(len(result), 0)

if __name__ == '__main__':
    main()
