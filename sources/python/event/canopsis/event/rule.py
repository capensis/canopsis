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

"""
Event processing library.

Provides tools to process an event related to a rule.

A rule is a list of actions or a couple of (condition, actions).
"""

from canopsis.common.path import resolve_element

CONDITION_FIELD = 'condition'
ACTIONS_FIELD = 'actions'

TASK_PATH = 'path'


class EventProcessingError(Exception):
    """
    Handle event processing errors.
    """

    pass

__GLOBAL_EVENT_PROCESSING = {}


def event_processing(event, ctx, **params):
    """
    Event processing signature to respect in order to process an event

    A condition may returns a boolean value.
    """

    pass


def get_event_task(path, cached=True):
    """
    Get an event processing related to an absolute input action path.

    :param str path: absolute task path to get

    :param bool cached: use runtime cache to get the action_path if previously
        loaded.

    :return: callable action which takes in parameter a context, an event, or
        None if action does not exist in runtime
    :rtype: callable which respects the signature of function event_procesing

    :raises EventProcessingError: if path is unknown from runtime.
    """

    result = None

    if cached and path in __GLOBAL_EVENT_PROCESSING:
        result = __GLOBAL_EVENT_PROCESSING[path]

    else:
        try:
            result = resolve_element(path)
        except ImportError:
            raise EventProcessingError(
                'path %s is unknown in runtime' % path)
        if result is not None and cached:
            __GLOBAL_EVENT_PROCESSING[path] = result

    return result


def get_processing_task(rule, task_name=None, cached=True):
    """
    Get callable task processing with params.

    :param dict rule: rule from where get task.

    :param str task_name:

    :param bool cached: try to get a cached task or not.

    :return: tuple of (callable task, task parameters)
    """

    result = None, None

    task = None

    if task_name is None:
        task = rule

    elif task_name in rule:
        task = rule[task_name]

    else:
        raise EventProcessingError(
            "No task name %s found in rule %s" % (task_name, rule))

    # get dedicated callable task with params
    if TASK_PATH in result:
        path = task[TASK_PATH]
        result = get_event_task(path=path, cached=cached), task

    return result


def process_event(event, rule, cached=True):
    """
    Apply input rule on input event in checking if the rule condition matches
    with the event and if True, execute rule actions.

    :param rule: rule to apply on input event. contains both condition and \
        actions.
    :type rule: dict

    :param event: event to check and to process.
    :type event: dict

    :param cached: indicates to actions to use cache instead of \
        importing them dynamically.
    :type cached: bool

    :return: ordered list of rule action results if rule condition is checked.
    :rtype: list
    """

    result = []

    # create a context which will be shared among condition and actions
    ctx = {}

    # get condition
    condition_task, parameters = get_processing_task(
        rule=rule, task_name=CONDITION_FIELD, cached=cached)

    # do actions if True
    do_actions = False

    if condition_task is None:
        do_actions = True

    else:
        do_actions = condition_task(event=event, ctx=ctx, **parameters)

    # if actions have to be performed
    if do_actions:
        # get actions
        actions = rule[ACTIONS_FIELD] if ACTIONS_FIELD in rule else ()

        # for all action
        for action in actions:
            # get action task with parameters
            action_task, parameters = get_processing_task(
                rule=action, cached=cached)
            # and run action
            try:
                action_result = action_task(event=event, ctx=ctx, **parameters)
            except Exception as e:
                result.append(EventProcessingError(e))
            else:
                # in adding the result in function result
                result.append(action_result)

    return result
