#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
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
    CONDITION_FIELD, ACTION_FIELD, TASK_PATH, TASK_PARAMS, RuleError,
    ConditionError, ActionError, get_task_with_params, process_rule,
    get_task, register_tasks, unregister_tasks, __TASK_PATHS as TASK_PATHS,
    register_task, RULES, rules, switch)

from canopsis.common.utils import path


def test_condition_true(event, ctx, **kwargs):
    return True


def test_condition_false(event, ctx, **kwargs):
    return False


MESSAGE = 'message'


def test_exception(event, ctx, **kwargs):
    raise Exception()

COUNT = 'count'


def test_condition_count(event, ctx, **kwargs):
    ctx[COUNT] = 0
    return True


def test_action(event, ctx, **kwargs):
    ctx[COUNT] += 1
    return ctx[COUNT]


def test_wrong_params():
    pass


class TaskUnregisterTest(TestCase):

    def setUp(self):
        # clear dico and add only self names
        TASK_PATHS.clear()
        self.names = 'a', 'b'
        for name in self.names:
            register_tasks(**{name: None})

    def test_unregister(self):
        for name in self.names:
            unregister_tasks(name)
            self.assertNotIn(name, TASK_PATHS)

    def test_unregister_all(self):
        unregister_tasks(*self.names)
        self.assertFalse(TASK_PATHS)

    def test_unregister_clear(self):
        unregister_tasks()
        self.assertFalse(TASK_PATHS)


class TaskRegistrationTest(TestCase):

    def setUp(self):
        # clean task paths
        TASK_PATHS.clear()
        self.tasks = {'a': None, 'b': None}
        register_tasks(**self.tasks)

    def test_register(self):
        """
        Check for registered task in registered tasks
        """
        for task in self.tasks:
            self.assertIn(task, TASK_PATHS)

    def test_register_raise(self):

        self.assertRaises(RuleError, register_tasks, **self.tasks)

    def test_register_force(self):

        register_tasks(force=True, **self.tasks)


class GetTaskTest(TestCase):

    def setUp(self):
        # clean all task paths
        TASK_PATHS.clear()

    def test_get_unregisteredtask(self):

        getTaskTest = path(GetTaskTest)
        task = get_task(getTaskTest)
        self.assertIs(task, GetTaskTest)

    def test_get_registeredtask(self):
        task_path = 'a'
        register_tasks(**{task_path: GetTaskTest})
        task = get_task(path=task_path)
        self.assertIs(task, GetTaskTest)


class TaskRegistrationDecoratorTest(TestCase):

    def setUp(self):
        TASK_PATHS.clear()

    def test_register_without_parameters(self):

        @register_task
        def register():
            pass
        self.assertIn('register', TASK_PATHS)

    def test_register(self):

        @register_task()
        def register():
            pass
        self.assertIn('register', TASK_PATHS)

    def test_registername(self):

        name = 'toto'

        @register_task(name)
        def register():
            pass
        self.assertIn(name, TASK_PATHS)

    def test_raise(self):
        name = 'toto'
        register_tasks(**{name: None})

        error = False
        try:
            @register_task
            def toto():
                pass
        except Exception:
            error = True
        self.assertTrue(error)


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

        self.assertEqual((task, params), (get_task_with_params, {}))

    def test_task_from_dict(self):

        task_conf = {TASK_PATH: self.existing_function}

        task, params = get_task_with_params(task_conf=task_conf)

        self.assertEqual((task, params), (get_task_with_params, {}))

    def test_task_from_dict_with_task_name(self):

        task_name = 'test'

        task_conf = {task_name: self.existing_function}

        task, params = get_task_with_params(
            task_conf=task_conf, task_name=task_name)

        self.assertEqual((task, params), (get_task_with_params, {}))

    def test_task_from_dict_with_task_name_and_dict(self):

        task_name = 'test'

        task_conf = {task_name: {TASK_PATH: self.existing_function}}

        task, params = get_task_with_params(
            task_conf=task_conf, task_name=task_name)

        self.assertEqual((task, params), (get_task_with_params, {}))

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

        condition, result = process_rule(
            event=self.event, rule=action, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertEqual(result, 1)

    def test_actions(self):

        action = {
            TASK_PATH: 'canopsis.rule.action.actions',
            TASK_PARAMS: {
                'actions': self.test_action
            }
        }

        condition, result = process_rule(
            event=self.event, rule=action, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertEqual(result, [1])

    def test_no_condition_action(self):

        rule = {ACTION_FIELD: self.test_action}

        condition, result = process_rule(
            event=self.event, rule=rule, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertEqual(result, 1)

    def test_no_condition_actions(self):

        action = {
            TASK_PATH: 'canopsis.rule.action.actions',
            TASK_PARAMS: {
                'actions': self.test_action
            }
        }
        rule = {ACTION_FIELD: action}

        condition, result = process_rule(
            event=self.event, rule=rule, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertEqual(result, [1])

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
            ACTION_FIELD: self.test_action
        }

        condition, result = process_rule(
            event=self.event, rule=rule, raiseError=True, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertEqual(result, 1)

    def test_action_error(self):

        rule = {
            ACTION_FIELD: self.test_exception
        }

        self.assertRaises(
            ActionError,
            process_rule,
            event=self.event, rule=rule, ctx=self.ctx, raiseError=True
        )

    def test_action_error_noraiseError(self):

        rule = {
            ACTION_FIELD: self.test_exception
        }

        condition, result = process_rule(
            event=self.event, rule=rule, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertTrue(type(result) is ActionError)

    def test_wrong_parameters(self):

        rule = {ACTION_FIELD: self.test_wrong_params}

        condition, result = process_rule(
            event=self.event, rule=rule, ctx=self.ctx)

        self.assertTrue(condition)
        self.assertTrue(type(result) is ActionError)

    def test_rules(self):

        self.count = 10

        rules_task = path(rules)

        # construct rules
        rule = {
            TASK_PATH: rules_task,
            TASK_PARAMS: {
                RULES: [self.test_action for i in range(self.count)]
            }
        }

        process_rule(event=self.event, rule=rule, ctx=self.ctx)

        self.assertEqual(self.ctx[COUNT], self.count)

    def test_switch(self):

        self.count = 10

        switch_task = path(switch)

        # construct rules
        rule = {
            TASK_PATH: switch_task,
            TASK_PARAMS: {
                RULES: [self.test_action for i in range(self.count)]
            }
        }

        process_rule(event=self.event, rule=rule, ctx=self.ctx)

        self.assertEqual(self.ctx[COUNT], 1)


if __name__ == '__main__':
    main()
