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
from canopsis.common.utils import lookup

from inspect import isroutine

"""
Task module.

Provides tools to process tasks rules.

A task uses a python function. Therefore it is possible to use an absolute
path or to register a function in rule tasks with the function/decorator
``register_task``. The related function must takes in parameter the ``event``
to process, a dict ``ctx`` which exists on a rule life and a ``**kwargs`` which
contains parameters filled related to task parameters and the rule api.

A task respects those types::
   - str: task path to execute.
   - dict: optional parameters
      + task_path: task path to execute.
      + params: dict of task parameters.

Therefore, a task could be as a str as a dict containing above structure.
"""

TASK = 'task'  #: task parameter name

TASK_PARAMS = 'params'  #: task params field name in task conf

TASK_PATH = 'task_path'  #: task path field name in task conf


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


__TASK_PATHS = {}


def get_task(path, cached=True):
    """
    Get task related to a path which could be:

    - a registered task.
    - a python function.

    :param str path: task path to get.
    :param bool cached: use cache (True by default).

    :raises ImportError: if task is not found in runtime.
    """

    result = None

    if path in __TASK_PATHS:
        result = __TASK_PATHS[path]
    else:
        result = lookup(path=path, cached=cached)
        __TASK_PATHS[path] = result

    return result


def register_tasks(force=False, **tasks):
    """
    Register a set of input task by name.

    :param bool force: force registration (default False).
    :param dict tasks: set of couple(name, function)

    :raises: TaskError if not force and task already exists.
    """

    already_registered = []

    for path in tasks:
        try:
            task = get_task(path)
        except ImportError:
            task = None
        if task is not None and not force:
            already_registered.append(path)
        else:
            __TASK_PATHS[path] = task

    if already_registered:
        raise TaskError(
            'Rule(s) {} are already registered'.format(already_registered)
        )


def register_task(name=None, force=False):
    """
    Decorator which registers function in registered tasks with function name
    """

    if isroutine(name):  # if name is a routine <=> no parameter is given
        # name is routine name
        result = name
        name = name.__name__
        register_tasks(force=force, **{name: result})

    else:  # if name is a str or None
        def register_task(function, name=name):
            """
            Register input function as a task
            """

            if name is None:
                name = function.__name__

            register_tasks(force=force, **{name: function})

            return function

        result = register_task

    return result


def unregister_tasks(*paths):
    """
    Unregister input paths. If paths is empty, clear all registered tasks.

    :param tuple paths: tuple of task paths
    """

    if paths:
        for path in paths:
            if path in __TASK_PATHS:
                del __TASK_PATHS[path]
    else:
        __TASK_PATHS.clear()


def get_task_with_params(task_conf, task_name=None, cached=True):
    """
    Get callable task processing with params.

    :param task_conf: task conf from where getting task.
    :type task_conf: str or dict

    :param str task_name: task name to find from input task_conf if not None

    :param bool cached: try to get a cached task or not.

    :return: tuple of (callable task, task parameters)
    """

    task, params = None, {}

    # if task_name is not None, try to find it in task_conf
    if task_name is not None:

        # in ensuring than task_conf is a dict
        if isinstance(task_conf, dict):
            # if task_name exists in task_conf, get it
            if task_name in task_conf:
                task_conf = task_conf[task_name]

    # get dedicated callable task without params
    if isinstance(task_conf, basestring):
        try:
            task = get_task(path=task_conf, cached=cached)
        except ImportError as ie:
            # Embed importerror in TaskError
            raise TaskError(ie)

    # get dedicated callable task with params
    elif TASK_PATH in task_conf:
        task_path = task_conf[TASK_PATH]
        try:
            task = get_task(path=task_path, cached=cached)

        except ImportError as ie:
            # embed import error in Rule Error
            raise TaskError(ie)

        else:
            # if task has been founded
            if task is not None:
                # try to get params
                if TASK_PARAMS in task_conf:
                    params = task_conf[TASK_PARAMS]

    # result is the couple (task, params)
    result = task, params

    return result


def run_task(
    task_conf=None, task_name=None,
    ctx=None,
    raiseError=True, exception_type=TaskError,
    cached=True
):
    """
    Run a single task related to a task_conf, and a task_name, a task ctx and
        not functional parameters.

    If an error occures, input exception_type is raised with raised error
        inside.

    :param task_conf: task configuration. None by default.
    :type task_conf: str or dict.
    :param str task_name task name to execute. None by default.
    :param dict ctx: task execution context. Empty dict by default.
    :param bool raiseError: if True (default), raise any task error, else
        result if the raised error.
    :param type exception_type: (default TaskError) exception type to raise if
        an error occured.
    :param bool cached: (True by default) use cache memory to save task
        references from input task name.
    """

    result = None

    # initialize the ctx
    if ctx is None:
        ctx = {}

    try:
        task, params = get_task_with_params(
            task_conf=task_conf, task_name=task_name, cached=cached)
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
