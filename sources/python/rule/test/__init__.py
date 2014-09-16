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

from canopsis.rule import (
    CONDITION_FIELD, ACTIONS_FIELD, TASK_PATH, TASK_PARAMS, RuleError,
    ConditionError, ActionError, get_task_with_params, process_rule)
from canopsis.common.utils import path


def test_condition_true(event, ctx):
    return True


def test_condition_false(event, ctx):
    return False


MESSAGE = 'message'


def test_exception(event, ctx):
    raise Exception()

COUNT = 'count'


def test_condition_count(event, ctx):
    ctx[COUNT] = 0
    return True


def test_action(event, ctx):
    ctx[COUNT] += 1
    return ctx[COUNT]


def test_wrong_params():
    pass


class GetTaskWithParamsTest(TestCase):

    def setUp(self):

        self.wrong_function = 'test.test'

        self.existing_function = 'canopsis.rule.get_task_with_params'

    def test_none_task_from_str(self):

        task_conf = self.wrong_function

        self.assertRaises(RuleError, get_task_with_params, task_conf=task_conf)

    def test_none_task_from_dict(self):

        task_conf = {TASK_PATH: self.wrong_function}

        self.assertRaises(RuleError, get_task_with_params, task_conf=task_conf)

    def test_none_task_from_dict_with_task_name(self):

        task_name = 'test'

        task_conf = {task_name: self.wrong_function}

        self.assertRaises(
            RuleError,
            get_task_with_params,
            task_conf=task_conf, task_name=task_name)

    def test_none_task_from_dict_with_task_name_and_dict(self):

        task_name = 'test'

        task_conf = {task_name: {TASK_PATH: self.wrong_function}}

        self.assertRaises(
            RuleError,
            get_task_with_params,
            task_conf=task_conf, task_name=task_name)

    def test_task_from_str(self):

        task_conf = self.existing_function

        task, params = get_task_with_params(task_conf=task_conf)

        self.assertEqual((task, params), (get_task_with_params, None))

    def test_task_from_dict(self):

        task_conf = {TASK_PATH: self.existing_function}

        task, params = get_task_with_params(task_conf=task_conf)

        self.assertEqual((task, params), (get_task_with_params, None))

    def test_task_from_dict_with_task_name(self):

        task_name = 'test'

        task_conf = {task_name: self.existing_function}

        task, params = get_task_with_params(
            task_conf=task_conf, task_name=task_name)

        self.assertEqual((task, params), (get_task_with_params, None))

    def test_task_from_dict_with_task_name_and_dict(self):

        task_name = 'test'

        task_conf = {task_name: {TASK_PATH: self.existing_function}}

        task, params = get_task_with_params(
            task_conf=task_conf, task_name=task_name)

        self.assertEqual((task, params), (get_task_with_params, None))

    def test_task_from_dict_with_params(self):

        param = 1

        task_conf = {
            TASK_PATH: self.existing_function,
            TASK_PARAMS: param}

        task, params = get_task_with_params(task_conf=task_conf)

        self.assertEqual((task, params), (get_task_with_params, param))

    def test_task_from_dict_with_name_and_params(self):

        param = 1

        task_name = 'test'

        task_conf = {
            task_name:
            {
                TASK_PATH: self.existing_function,
                TASK_PARAMS: param
            }
        }

        task, params = get_task_with_params(
            task_conf=task_conf, task_name=task_name)

        self.assertEqual((task, params), (get_task_with_params, param))

    def test_cache(self):

        task_conf = self.existing_function

        task_not_cached_0, _ = get_task_with_params(
            task_conf=task_conf, cached=False)

        task_not_cached_1, _ = get_task_with_params(
            task_conf=task_conf, cached=False)

        self.assertTrue(task_not_cached_0 is task_not_cached_1)

        task_cached_0, _ = get_task_with_params(task_conf=task_conf)

        task_cached_1, _ = get_task_with_params(task_conf=task_conf)

        self.assertTrue(task_cached_0 is task_cached_1)


class TestProcessRule(TestCase):

    def setUp(self):
        self.event = {}
        self.test_action = path(test_action)
        self.test_wrong_params = path(test_wrong_params)
        self.test_exception = path(test_exception)
        self.test_condition_count = path(test_condition_count)
        self.test_condition_false = path(test_condition_false)
        self.test_condition_true = path(test_condition_true)
        self.ctx = {COUNT: 0}

    def test_action(self):

        action = self.test_action

        condition, result = process_rule(event=self.event, rule=action, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0], 1)

    def test_actions(self):

        action = [self.test_action]

        condition, result = process_rule(event=self.event, rule=action, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0], 1)

    def test_no_condition_action(self):

        rule = {ACTIONS_FIELD: self.test_action}

        condition, result = process_rule(event=self.event, rule=rule, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0], 1)

    def test_no_condition_actions(self):

        actions = (self.test_action,)
        rule = {ACTIONS_FIELD: actions}

        condition, result = process_rule(event=self.event, rule=rule, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0], 1)

    def test_condition_error(self):

        rule = {CONDITION_FIELD: self.test_exception}

        self.assertRaises(
            ConditionError,
            process_rule,
            event=self.event, rule=rule, raiseError=True)

    def test_condition_false(self):

        rule = {CONDITION_FIELD: self.test_condition_false}

        condition, result = process_rule(event=self.event, rule=rule)

        self.assertFalse(condition)

    def test_condition_true(self):

        rule = {
            CONDITION_FIELD: self.test_condition_true,
            ACTIONS_FIELD: self.test_action
        }

        condition, result = process_rule(
            event=self.event, rule=rule, raiseError=True, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertTrue(result)
        self.assertEqual(result[0], 1)

    def test_action_error(self):

        rule = {
            ACTIONS_FIELD: self.test_exception
        }

        self.assertRaises(
            ActionError,
            process_rule,
            event=self.event, rule=rule, ctx=self.ctx, raiseError=True
        )

    def test_action_error_noraiseError(self):

        rule = {
            ACTIONS_FIELD: self.test_exception
        }

        condition, result = process_rule(
            event=self.event, rule=rule, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertTrue(result)
        self.assertTrue(type(result[0]) is ActionError)

    def test_wrong_parameters(self):

        rule = {ACTIONS_FIELD: self.test_wrong_params}

        condition, result = process_rule(
            event=self.event, rule=rule, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertTrue(result)
        self.assertTrue(type(result[0]) is ActionError)

if __name__ == '__main__':
    main()
