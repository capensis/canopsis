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


class TaskError(Exception):
    """
    Default task error.
    """


# global set of tasks by ids
__TASKS_BY_ID = {}


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
        # if _id has been already registered
        if _id in __TASKS_BY_ID:
            if force:
                # save old task in result
                result[_id] = __TASKS_BY_ID[_id]
            else:
                # raise old tasks if not force
                raise TaskError(
                    'Rule {} is already registered'.format(_id)
                )
        # save new task
        __TASKS_BY_ID[_id] = task

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

# register the function register_tasks
register_task()(register_tasks)


@register_task
def get_task(_id, cache=True, cacheonly=False):
    """
    Get task related to an id which could be:

    - a registered task id.
    - a python path to a function.

    :param str id: task id to get.
    :param bool cache: use cache system to quick access to task
        (True by default).
    :param bool cacheonly: if True, do not try to lookup for resolving the task
        searching.

    :raises ImportError: if task is not found in runtime.
    """

    result = None

    if _id in __TASKS_BY_ID:
        result = __TASKS_BY_ID[_id]
    elif not cacheonly:
        result = lookup(path=_id, cached=cache)

    return result


@register_task
def task(**kwargs):
    """
    Default task signature.
    """


@register_task
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


@register_task
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
                    params = conf[TASK_PARAMS].copy()

    # result is the couple (task, params)
    result = task, params

    return result


@register_task
def run_task(conf=None, ctx=None, raiseerror=True, cache=True, **kwargs):
    """
    Run a single task related to a conf, and a task_name, a task ctx and
        not functional parameters.

    If an error occures, input exception_type is raised with raised error
        inside.

    :param conf: task configuration. None by default.
    :type conf: str or dict.
    :param dict ctx: task execution context. Empty dict by default.
    :param bool raiseerror: if True (default), raise any task error, else
        result if the raised error.
    :param bool cache: (True by default) use cache memory to save task
        references from input task name.
    :param kwargs: additional task parameters.
    """

    result = None

    # initialize the ctx
    if ctx is None:
        ctx = {}

    task, params = get_task_with_params(conf=conf, cache=cache)
    # add kwargs in params
    params.update(kwargs)

    try:  # process task
        result = task(ctx=ctx, **params)
    except Exception as error:  # if an error occured
        if raiseerror:  # if raiseerror, raise it
            raise
        result = error  # or result is the error

    return result


@register_task
def new_conf(task, **params):
    """
    Generate a new task conf related to input task id and params.

    :param task: task identifier.
    :type task: str or routine
    :param dict params: task parameters.

    :return: {TASK_ID: _id, TASK_PARAMS: params}
    :rtype: dict
    """

    result = None

    # if task is a task routine, find the corresponding task id
    if isroutine(task):
        for task_id in __TASKS_BY_ID:
            _task = __TASKS_BY_ID[task_id]
            if _task == task:
                task = task_id

    result = {
        TASK_ID: task,
        TASK_PARAMS: params
    }

    return result


#: action result
RESULT = 'result'
#: action error
ERROR = 'error'


@register_task
def tasks(confs=None, raiseerror=False, **kwargs):
    """
    run a list of tasks in processing several input confs and returns a list
        of dict {'result': task result, 'error': task error}.

    :param confs: task confs to process.
    :type confs: str or list or dict
    :param bool raiseerror: if True (default False) raise the first encountered
        error during processing of tasks.
    :param kwargs: additional parameters to process in all tasks.

    :return: a list containing dict of {RESULT: result, ERROR: error}.
    :rtype: list
    """

    result = []

    if confs is not None:
        # ensure confs is a list
        if isinstance(confs, (basestring, dict)):
            confs = [confs]
        for conf in confs:
            task, params = get_task_with_params(conf)
            # add kwargs in params
            params.update(kwargs)
            # initialize both result and error
            task_result, task_error = None, None
            try:
                task_result = task(**params)
            except Exception as task_error:
                if raiseerror:
                    raise

            result_to_append = {
                RESULT: task_result,
                ERROR: task_error
            }
            result.append(result_to_append)

    return result
