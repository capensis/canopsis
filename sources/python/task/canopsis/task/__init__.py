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

from canopsis.common.utils import lookup

from inspect import isroutine

"""
Event processing rule module.

Provides tools to process event rules.

A rule is a couple of (condition, action) tasks where condition can be None.

A rule respects those types::

   - task: task to execute.
   - dict:
      + condition (optional): condition task to check.
      + action: action task to run if condition does not exist or is True.

A task uses a python function. Therefore it is possible to use an absolute
path or to register a function in rule tasks with the function/decorator
``register_task``. The related function must takes in parameter the ``event``
to process, a dict ``ctx`` which exists on a rule life and a ``**kwargs`` which
contains parameters filled related to task parameters and the rule api.

A task respects those types::
   - str: task path to execute.
   - dict:
      + task_path: task path to execute.
      + params: dict of task parameters.

For example, let ``my.my_condition`` and ``my.my_action`` respectively
custom condition and action.

A typical parameterized rule configuration is as follow:
{
    "condition": "my.my_condition",
    "action": "my.my_action"
}

Or this one if you have the parameter ``foo`` equals to 2 to use in the action:
{
    "condition": "my.my_condition",
    "action": {
        "task_path": "my.my_action",
        "params": {
            "foo": 2
        }
    }
}

Or even this one without condition and where the action is registered with the
name ``my_action``.
{
    "action": "my_action"
}
"""

RULE = 'rule'

CONDITION_FIELD = 'condition'  #: condition field name in rule conf
ACTION_FIELD = 'action'  #: actions field name in rule conf

TASK_PARAMS = 'params'  #: task params field name in task conf

TASK_PATH = 'task_path'  #: task path field name in task conf

RULES = 'rules'  #: rules rule name
SWITCH = 'switch'  #: switch rule name


def task(event, ctx, **kwargs):
    """
    Default task signature.
    """
    pass


class RuleError(Exception):
    """
    Handle rule error.
    """

    pass


class ConditionError(RuleError):
    """
    Handle condition error.
    """

    pass


class ActionError(RuleError):
    """
    Handle action error
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

    if path not in __TASK_PATHS:
        result = lookup(path=path, cached=cached)
    else:
        result = __TASK_PATHS[path]

    return result


def register_tasks(force=False, **tasks):
    """
    Register a set of input task by name.

    :param bool force: force registration (default False).
    :param dict tasks: set of couple(name, function)

    :raises RuleError: if not force and task already exist
    """

    for path in tasks:
        task = tasks[path]
        if not force and path in __TASK_PATHS:
            raise RuleError('Rule %s already registered' % path)
        else:
            __TASK_PATHS[path] = task


def register_task(name=None, force=False):
    """
    Decorator which registers function in registered tasks with function name
    """

    if isroutine(name):
        # if no parameter has been given
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
            __TASK_PATHS.pop(path, None)
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
            # Embed importerror in RuleError
            raise RuleError(ie)

    # get dedicated callable task with params
    elif TASK_PATH in task_conf:
        task_path = task_conf[TASK_PATH]
        try:
            task = get_task(path=task_path, cached=cached)

        except ImportError as ie:
            # embed import error in Rule Error
            raise RuleError(ie)

        else:
            # if task has been founded
            if task is not None:
                # try to get params
                if TASK_PARAMS in task_conf:
                    params = task_conf[TASK_PARAMS]

    # result is the couple (task, params)
    result = task, params

    return result


def process_rule(event, rule, ctx=None, cached=True, raiseError=False):
    """
    Apply input rule on input event in checking if the rule condition matches
    with the event and if True, execute rule actions.

    :param rule: rule to apply on input event. contains both condition and
        actions.
    :type rule: dict

    :param event: event to check and to process.
    :type event: dict

    :param cached: indicates to actions to use cache instead of
        importing them dynamically.
    :type cached: bool

    :param bool raiseError: If True (False by default), raise the first error
        encountered while executing conditions (ConditionError) or actions
        (ActionError).

    :return: a tuple of (condition, action) results.
    :rtype: tuple
    """

    # create a context which will be shared among condition and actions
    if ctx is None:
        ctx = {}

    # do actions if a condition may exist (not if rule is an iterable)
    do_actions = not isinstance(rule, dict) or CONDITION_FIELD not in rule

    # if a condition may be founded
    if not do_actions:

        do_actions = run_task(
            event=event,
            ctx=ctx,
            task_conf=rule,
            task_name=CONDITION_FIELD,
            cached=cached,
            exception_type=ConditionError,
            raiseError=raiseError)

    action_result = None

    # if actions have to be performed
    if do_actions is True:

        action_result = run_task(
            event=event,
            ctx=ctx,
            task_conf=rule,
            task_name=ACTION_FIELD,
            cached=cached,
            exception_type=ActionError,
            raiseError=raiseError)

    result = do_actions, action_result

    return result


def run_task(
    event, ctx, task_conf, task_name, raiseError, exception_type, cached
):
    """
    Run a single task related to an event, a ctx, a task_conf, and a task_name.

    If an error occures, input exception_type is raised.

    :param dict event: event to process.
    :param dict ctx: task execution context.
    :param task_conf: task configuration.
    :type task_conf: str or dict.
    :param str task_name task name to execute.
    :param bool raiseError: if True, raise any task error, else result if the
        raised error.
    :param type exception_type: exception type to raise if an error occured
    :param bool cached: use cache memory to save task references from input
        task name.
    """

    result = None

    try:
        task, params = get_task_with_params(
            task_conf=task_conf, task_name=task_name, cached=cached)
    except RuleError as e:
        # if action does not exist, do nothing
        pass
    else:
        # initialize params if None
        try:
            result = task(
                event=event, ctx=ctx, raiseError=raiseError, **params)
        except Exception as e:
            error = exception_type(e)
            if raiseError:
                raise error
            result = error

    return result


@register_task
def rules(event, ctx, rules, cached=True, raiseError=False, **kwargs):
    """
    Rule which run all input rules.

    :param rules: rules to process
    :return: condition, action where condition is a logical and on all
        conditions, and action is a list of all action results.
    """

    condition = False
    action = []

    for rule in rules:

        #  execute all rules while condition is False
        result_condition, result_action = process_rule(
            event=event,
            ctx=ctx,
            rule=rule,
            cached=cached,
            raiseError=raiseError,
            **kwargs)

        condition |= result_condition
        action.append(result_action)

    result = condition, action

    return result


@register_task
def switch(event, ctx, rules, cached=True, raiseError=False, **kwargs):
    """
    Execute first rule among input
    """

    result = False, None

    for rule in rules:

        # try to execute all rules while condition is False
        condition, action = process_rule(
            event=event,
            ctx=ctx,
            rule=rule,
            cached=cached,
            raiseError=raiseError,
            **kwargs)

        # stop when condition is True
        if condition is True:
            result = condition, action
            break

    return result
