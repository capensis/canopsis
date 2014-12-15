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


"""
Event rule module.

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

__version__ = "0.1"

from canopsis.common.init import basestring
from canopsis.task import register_task, run_task

RULE = 'rule'

CONDITION_FIELD = 'condition'  #: condition field name in rule conf
ACTION_FIELD = 'action'  #: actions field name in rule conf

RULES = 'rules'  #: rules rule name


class RuleError(Exception):
    """
    Handle rule error.
    """

    pass


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
            raiseError=raiseError
        )

    action_result = None

    # if actions have to be performed
    if do_actions is True:

        action_result = run_task(
            event=event,
            ctx=ctx,
            task_conf=rule,
            task_name=ACTION_FIELD,
            cached=cached,
            raiseError=raiseError
        )

    result = do_actions, action_result

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
