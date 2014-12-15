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

__version__ = "0.1"

from canopsis.common.init import basestring
from canopsis.common.utils import lookup, path

from inspect import isroutine

"""
Task module.

Provides tools to process tasks.

A task uses a python function. Therefore it is possible to use an absolute
path, an id or to register a function in tasks with the function/decorator
``register_task``. The related function may take in parameter a dict ``ctx``
which exists on a task life and a ``**kwargs`` which
contains parameters filled related to task parameters.

A task respects those types::
   - str: task name to execute.
   - dict:
      + id: task name to execute.
      + params: dict of task parameters.
"""

TASK = 'task'  #: task parameter name

TASK_PARAMS = 'params'  #: task params field name in task conf

TASK_ID = 'id'  #: task id field name in task conf


def task(**kwargs):
    """
    Default task signature.
    """

    pass


class TaskError(Exception):
    """
    Default task error.
    """

    pass


# global set of tasks by ids
__TASKS_BY_ID = {}


def get_task(_id, cache=True):
    """
    Get task related to an id which could be:

    - a registered task id.
    - a python path to a function.

    :param str id: task id to get.
    :param bool cache: use cache system to quick access to task
        (True by default).

    :raises ImportError: if task is not found in runtime.
    """

    result = None

    if _id in __TASKS_BY_ID:
        result = __TASKS_BY_ID[_id]
    else:
        result = lookup(path=_id, cached=cache)

    return result


def register_tasks(force=False, **tasks):
    """
    Register a set of input task by name.

    :param bool force: force registration (default False).
    :param dict tasks: set of couple(name, function)

    :return: old tasks by id.
    :rtype: dict
    :raises: TaskError if not force and task already exists.
    """

    result = {}

    for _id in tasks:
        task = tasks[_id]
        old_task = None
        try:  # is task existing ?
            old_task = get_task(_id)
        except ImportError:
            pass
        # if old task does not exist, save new task in global cache
        if old_task is None:
            __TASKS_BY_ID[_id] = task
        else:
            if force:  # if force, overwrite old task in cache
                __TASKS_BY_ID[_id] = task
            # save old task in the result
            result[_id] = old_task

    # raise old tasks if not force
    if result and not force:
        raise TaskError(
            'Rule(s) {} are already registered'.format(result)
        )

    return result


def register_task(_id=None, force=False):
    """
    Decorator which registers function in registered tasks with function _id
    """

    if isroutine(_id):  # if _id is a routine <=> no parameter is given
        # _id is routine _id
        result = _id  # task is the routine
        _id = path(_id)  # task id is its python path
        register_tasks(force=force, **{_id: result})

    else:  # if _id is a str or None
        def register_task(function, _id=_id):
            """
            Register input function as a task
            """

            if _id is None:
                _id = path(function)

            register_tasks(force=force, **{_id: function})

            return function

        result = register_task

    return result


def unregister_tasks(*ids):
    """
    Unregister input ids. If ids is empty, clear all registered tasks.

    :param tuple ids: tuple of task ids
    """

    if ids:
        for id in ids:
            if id in __TASKS_BY_ID:
                del __TASKS_BY_ID[id]
    else:
        __TASKS_BY_ID.clear()


def get_task_with_params(conf, cache=True):
    """
    Get callable task processing with params.

    :param conf: task conf from where getting task.
    :type conf: str or dict
    :param str task_name: task name to find from input conf if not None.
    :param bool cache: try to get a cache task or not.

    :return: tuple of (callable task, task parameters)
    :raises: ImportError if task is not registered.
    """

    task, params = None, {}

    # get dedicated callable task without params
    if isinstance(conf, basestring):
        task = get_task(_id=conf, cache=cache)

    # get dedicated callable task with params
    elif TASK_ID in conf:
        task_path = conf[TASK_ID]
        try:
            task = get_task(_id=task_path, cache=cache)

        except ImportError:
            raise

        else:
            # if task has been founded
            if task is not None:
                # try to get params
                if TASK_PARAMS in conf:
                    params = conf[TASK_PARAMS]

    # result is the couple (task, params)
    result = task, params

    return result


def run_task(
    conf=None,
    ctx=None,
    raiseError=True, exception_type=TaskError,
    cache=True
):
    """
    Run a single task related to a conf, and a task_name, a task ctx and
        not functional parameters.

    If an error occures, input exception_type is raised with raised error
        inside.

    :param conf: task configuration. None by default.
    :type conf: str or dict.
    :param dict ctx: task execution context. Empty dict by default.
    :param bool raiseError: if True (default), raise any task error, else
        result if the raised error.
    :param type exception_type: (default TaskError) exception type to raise if
        an error occured.
    :param bool cache: (True by default) use cache memory to save task
        references from input task name.
    """

    result = None

    # initialize the ctx
    if ctx is None:
        ctx = {}

    try:
        task, params = get_task_with_params(conf=conf, cache=cache)
    except TaskError as e:
        # if action does not exist and raiseError is False, do nothing
        if raiseError:
            raise
    else:
        try:  # process task
            result = task(ctx=ctx, **params)
        except Exception as e:  # if an error occured
            error = exception_type(e)  # embed it
            if raiseError:  # if raiseError, raise it
                raise error
            result = error  # or result is the error

    return result


def new_conf(_id, **params):
    """
    Generate a new task conf related to input _id and params.

    :param str _id: task id
    :param dict params: task parameters.

    :return: task conf depending on params:
        - empty: _id
        - not empty: {TASK_ID: _id, TASK_PARAMS: params}
    :rtype: str or dict
    """

    result = None

    if not params:
        result = _id

    else:
        result = {
            TASK_ID: _id,
            TASK_PARAMS: params
        }

    return result
