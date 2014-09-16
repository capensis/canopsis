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

from collections import Iterable

from canopsis.common.utils import resolve_element
from canopsis.common.utils import force_iterable

"""
Event processing rule module.

Provides tools to process event rules.

A rule is a couple of (condition, action+) where condition can be None.
"""

CONDITION_FIELD = 'condition'  #: condition field name in rule conf
ACTIONS_FIELD = 'actions'  #: actions field name in rule conf

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

    task, params = None, None

    # if task_name is not None, find it in task_conf
    if task_name is not None:

        # in ensuring than task_conf is a dict
        if isinstance(task_conf, dict) and task_name in task_conf:
            task_conf = task_conf[task_name]

        else:  # else raise a Rule error because task_name does not exist
            raise RuleError(
                'Task %s is not in task_conf %s' % (task_name, task_conf))

    # get dedicated callable task without params
    if isinstance(task_conf, str):
        try:
            task = resolve_element(path=task_conf, cached=cached)

        except ImportError as ie:
            # Embed importerror in RuleError
            raise RuleError(ie)

    # get dedicated callable task with params
    elif TASK_PATH in task_conf:
        task_path = task_conf[TASK_PATH]
        try:
            task = resolve_element(path=task_path, cached=cached)

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
        encountered while executing actions.

    :return: a tuple of (condition result, action ordered result)
    :rtype: list
    """

    result = False, []

    # create a context which will be shared among condition and actions
    if ctx is None:
        ctx = {}

    # do actions if a condition may exist (not if rule is an iterable)
    do_actions = not isinstance(rule, dict)

    # if a condition may be founded
    if isinstance(rule, dict):

        # get condition
        try:
            condition_task, params = get_task_with_params(
                task_conf=rule, task_name=CONDITION_FIELD, cached=cached)
        except RuleError:
            # if condition does not exist, do_actions is True
            do_actions = True
        else:
            if params is None:
                params = {}
            try:
                do_actions = condition_task(event=event, ctx=ctx, **params)
            except Exception as e:

                if raiseError:
                    raise ConditionError(e)

    action_results = []

    # if actions have to be performed
    if do_actions:

        action_confs = rule

        # if rule is a dictionary
        if isinstance(action_confs, dict):
            # get action_confs
            action_confs = rule[ACTIONS_FIELD] if ACTIONS_FIELD in rule else ()

        if isinstance(action_confs, str):
            action_confs = [action_confs]

        # for all action
        for action_conf in action_confs:
            # get action_conf task with params

            action_task, params = get_task_with_params(
                task_conf=action_conf, cached=cached)

            # if params is None, params is an empty dict
            if params is None:
                params = {}

            # and run action
            try:
                action_result = action_task(event=event, ctx=ctx, **params)
            except Exception as e:
                error = ActionError(e)

                if raiseError:
                    raise error
                else:
                    action_results.append(error)
            else:
                # in adding the action result in action_results
                action_results.append(action_result)

    result = do_actions, action_results

    return result
