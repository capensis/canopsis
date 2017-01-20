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

from canopsis.common.utils import path
from canopsis.task.core import (
    __TASKS_BY_ID as TASKS_BY_ID, TASK_PARAMS, TASK_ID, TaskError,
    get_task, new_conf,  # get_task_with_params,
    run_task, register_tasks, register_task, unregister_tasks,
    tasks, RESULT, ERROR
)


def test_exception(**kwargs):

    raise Exception()


COUNT = 'count'


class TaskUnregisterTest(TestCase):
    """
    Test unregister_tasks function
    """

    def setUp(self):
        """
        Create two tasks for future unregistration.
        """

        # clear dico and add only self names
        self.names = 'a', 'b'
        for name in self.names:
            register_tasks(force=True, **{name: None})

    def test_unregister(self):
        """
        Unregister one by one
        """

        for name in self.names:
            unregister_tasks(name)
            self.assertNotIn(name, TASKS_BY_ID)

    def test_unregister_all(self):
        """
        Unregister all tasks at a time.
        """

        unregister_tasks(*self.names)
        for name in self.names:
            self.assertNotIn(name, TASKS_BY_ID)

    def test_unregister_clear(self):
        """
        Unregister all tasks with an empty parameter.
        """

        _TASKS_BY_ID = TASKS_BY_ID.copy()
        unregister_tasks()
        self.assertFalse(TASKS_BY_ID)
        TASKS_BY_ID.update(_TASKS_BY_ID)


class TaskRegistrationTest(TestCase):
    """
    Test to register tasks.
    """

    def setUp(self):
        """
        """
        # clean task paths
        self.tasks = {'a': 1, 'b': 2, 'c': 3}
        register_tasks(force=True, **self.tasks)

    def test_register(self):
        """
        Check for registered task in registered tasks
        """
        for task in self.tasks:
            self.assertIn(task, TASKS_BY_ID)

    def test_register_raise(self):
        """
        Test to catch TaskError while registring already present tasks.
        """

        with self.assertRaises(TaskError):
            register_tasks(**self.tasks)

    def test_register_force(self):
        """
        Test to register existing tasks with force.
        """

        self.assertNotIn('d', TASKS_BY_ID)
        new_tasks = {'a': 'a', 'c': 'c', 'd': 'd'}
        old_tasks = register_tasks(
            force=True, **new_tasks
        )
        for new_task in new_tasks:
            self.assertEqual(get_task(new_task), new_tasks[new_task])
        self.assertNotIn('b', old_tasks)
        self_tasks_wo_b = self.tasks
        del self_tasks_wo_b['b']
        self.assertEqual(old_tasks, self_tasks_wo_b)


class GetTaskTest(TestCase):
    """
    Test get task function.
    """

    def test_get_unregisteredtask(self):
        """
        Test to get unregistered task.
        """

        getTaskTest = path(GetTaskTest)
        task = get_task(getTaskTest)
        self.assertEqual(task, GetTaskTest)

    def test_get_registeredtask(self):
        """
        Test to get registered task.
        """

        _id = 'a'
        register_tasks(force=True, **{_id: GetTaskTest})
        task = get_task(_id=_id)
        self.assertEqual(task, GetTaskTest)


class TaskRegistrationDecoratorTest(TestCase):
    """
    Test registration decorator
    """

    def test_register_without_parameters(self):

        def register():
            pass
        register_task(force=True)(register)
        self.assertIn(path(register), TASKS_BY_ID)

    def test_register(self):

        @register_task(force=True)
        def register():
            pass
        self.assertIn(path(register), TASKS_BY_ID)

    def test_registername(self):

        _id = 'toto'

        @register_task(_id, force=True)
        def register():
            pass
        self.assertIn(_id, TASKS_BY_ID)

    def test_raise(self):

        def toto():
            pass
        _id = path(toto)

        register_tasks(force=True, **{_id: 6})

        self.assertRaises(TaskError, register_task, toto)


class GetTaskWithParamsTest(TestCase):
    """
    Test get task with params function.
    """

    def setUp(self):

        self.wrong_function = 'test.test'

        self.existing_function = 'canopsis.task.get_task_with_params'

    # TODO 4-01-2017
    # def test_none_task_from_str(self):

    #     conf = self.wrong_function

    #     self.assertRaises(ImportError, get_task_with_params, conf=conf)

    # TODO 4-01-2017
    # def test_none_task_from_dict(self):

    #     conf = {TASK_ID: self.wrong_function}

    #     self.assertRaises(ImportError, get_task_with_params, conf=conf)

    # TODO 4-01-2017
    # def test_task_from_str(self):

    #     conf = self.existing_function

    #     task, params = get_task_with_params(conf=conf)

    #     self.assertEqual((task, params), (get_task_with_params, {}))

    # TODO 4-01-2017
    # def test_task_from_dict(self):

    #     conf = {TASK_ID: self.existing_function}

    #     task, params = get_task_with_params(conf=conf)

    #     self.assertEqual((task, params), (get_task_with_params, {}))

    # TODO 4-01-2017
    # def test_task_from_dict_with_params(self):

    #     param = {'a': 1}

    #     conf = {
    #         TASK_ID: self.existing_function,
    #         TASK_PARAMS: param}

    #     task, params = get_task_with_params(conf=conf)

    #     self.assertEqual((task, params), (get_task_with_params, param))

    # TODO 4-01-2017
    # def test_cache(self):

    #     conf = self.existing_function

    #     task_not_cached_0, _ = get_task_with_params(
    #         conf=conf, cache=False)

    #     task_not_cached_1, _ = get_task_with_params(
    #         conf=conf, cache=False)

    #     self.assertTrue(task_not_cached_0 is task_not_cached_1)

    #     task_cached_0, _ = get_task_with_params(conf=conf)

    #     task_cached_1, _ = get_task_with_params(conf=conf)

    #     self.assertTrue(task_cached_0 is task_cached_1)


class RunTaskTest(TestCase):
    """
    Test run task.
    """

    def setUp(self):

        @register_task('test', force=True)
        def test(**kwargs):
            return self

        @register_task('test_exception', force=True)
        def test_exception(**kwargs):
            raise Exception()

        @register_task('test_params', force=True)
        def test_params(ctx, **kwargs):
            return kwargs['a'] + kwargs['b'] + ctx['a'] + 1

    def test_simple(self):
        """
        Test simple task.
        """

        result = run_task('test')
        self.assertIs(result, self)

    def test_exception(self):
        """
        Test task which raises an exception.
        """

        self.assertRaises(Exception, run_task, 'test_exception')

    def test_exception_without_raise(self):
        """
        Test task which raises an exception.
        """

        result = run_task('test_exception', raiseerror=False)
        self.assertTrue(isinstance(result, Exception))

    def test_simple_params(self):
        """
        Test task with params
        """

        conf = new_conf('test_params', **{'a': 1, 'b': 2})
        result = run_task(conf, ctx={'a': 1})
        self.assertEqual(result, 5)


class NewConfTest(TestCase):
    """
    Test new conf.
    """

    def test_id(self):
        """
        Test to generate a new conf with only an id.
        """

        conf = new_conf('a')
        self.assertEqual(conf, {'id': 'a', 'params': {}})

    def test_with_empty_params(self):
        """
        Test to generate a new conf with empty params.
        """

        conf = new_conf('a', **{})

        self.assertEqual(conf, {'id': 'a', 'params': {}})

    def test_with_params(self):
        """
        Test to generate a new conf with params.
        """

        params = {'a': 1}
        conf = new_conf('a', **params)

        self.assertEqual(conf[TASK_ID], 'a')
        self.assertEqual(conf[TASK_PARAMS], params)

    def test_with_routine(self):
        """
        Test to generate a new conf related to a task routine.
        """

        conf = new_conf(run_task)

        self.assertEqual(conf['id'], path(run_task))

    def test_with_routine_and_params(self):
        """
        Test to generate a new conf related to a task routine and params.
        """

        params = {'a': 1}
        conf = new_conf(run_task, **params)

        self.assertEqual(conf[TASK_ID], path(run_task))
        self.assertEqual(conf[TASK_PARAMS], params)


class TasksTest(TestCase):
    """
    Test tasks function.
    """

    def setUp(self):

        @register_task('action', force=True)
        def action(**kwargs):
            pass

        self.error = NotImplementedError()

        @register_task('error', force=True)
        def action_raise(**kwargs):
            raise self.error

        @register_task('one', force=True)
        def action_one(**kwargs):
            return 1

    def test_empty(self):
        """
        Test to execute empty tasks.
        """

        results = tasks()
        self.assertEqual(len(results), 0)

    def test_unary(self):
        """
        Test to execute one void action.
        """

        results = tasks(confs='one')
        self.assertEqual(results, [{RESULT: 1, ERROR: None}])

    def test_exception(self):
        """
        Test to run an action which raises an exception.
        """

        results = tasks(confs='error')
        self.assertEqual(len(results), 1)
        self.assertIsNone(results[0][RESULT])
        self.assertTrue(isinstance(results[0][ERROR], NotImplementedError))

    def test_raiseerror(self):
        """
        Test to raise an error from a task execution.
        """

        self.assertRaises(
            NotImplementedError,
            tasks,
            confs='error',
            raiseerror=True
        )

    def test_many(self):
        """
        Test to run tasks without errors.
        """

        confs = ["action", "error", "one"]
        _results = [
            {RESULT: None, ERROR: None},
            {RESULT: None, ERROR: self.error},
            {RESULT: 1, ERROR: None}
        ]
        results = tasks(confs=confs)
        self.assertEqual(len(results), len(confs))
        self.assertEqual(_results, results)


if __name__ == '__main__':
    main()
