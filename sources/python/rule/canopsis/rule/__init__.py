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

__version__ = "0.1"

from canopsis.common.utils import lookup

"""
Event processing rule module.

Provides tools to process event rules.

A rule is a couple of (condition, action) where condition can be None.
"""

RULE = 'rule'

CONDITION_FIELD = 'condition'  #: condition field name in rule conf
ACTION_FIELD = 'action'  #: actions field name in rule conf
ELSE_FIELD = 'else'  #: else actions field name in rule conf

TASK_PARAMS = 'params'  #: task params field name in task conf

TASK_PATH = 'task_path'  #: task path field name in task conf


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
    if isinstance(task_conf, str):
        try:
            task = lookup(path=task_conf, cached=cached)
        except ImportError as ie:
            # Embed importerror in RuleError
            raise RuleError(ie)

    # get dedicated callable task with params
    elif TASK_PATH in task_conf:
        task_path = task_conf[TASK_PATH]
        try:
            task = lookup(path=task_path, cached=cached)

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
    :rtype: list
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

    else:  # go to the else field

        action_result = run_task(
            event=event,
            ctx=ctx,
            task_conf=rule,
            task_name=ELSE_FIELD,
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
