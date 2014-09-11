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

from canopsis.common.utils import resolve_element

# rule actions field
ACTIONS_FIELD = 'actions'

# cache dictionary which stores action by their names
_GLOBAL_ACTIONS = {}

ACTION_NAME_FIELD = 'name'


class ActionError(Exception):
    """
    Manage action errors.
    """

    pass


def get_action(action_path, cached=True):
    """
    Get a callable related to absolute input action path.

    :param action_path: absolute action path to get
    :type action_path: str

    :param cached: use runtime cache to get the action_path if previously
        loaded.
    :type cached: bool

    :return: callable action which takes in parameter a context, an event, or
        None if action does not exist in runtime
    :rtype: callable

    :raise: ActionError if:
        - action is unknown from runtime.
    """

    result = None

    if cached and action_path in _GLOBAL_ACTIONS:
        result = _GLOBAL_ACTIONS[action_path]

    else:
        try:
            result = resolve_element(action_path)
        except ImportError:
            raise ActionError('action %s is unknown in runtime' % action_path)
        if result is not None and cached:
            _GLOBAL_ACTIONS[action_path] = result

    return result


def do_action(action, ctx, event, cached_action=True):
    """
    Do an action function related to input name.

    An action should take in parameters:
    - a ctx.
    - an event.
    - a kwargs such as action parameters.

    :param action: action configuration to run.
    :type action: dict

    :param event: event to process with input action.
    :type event: dict

    :param cached_action: use cache in order to resolve an action.
    :type cached_action: bool

    :return: action processing result.

    :raise: ActionError if:
        - action is unknown from runtime.
        - action does not have a name.
        - action execution raises an Exception.
    """

    result = None

    # start to find related action
    if ACTION_NAME_FIELD not in action:
        # raise an ActionError if action does not exist
        raise ActionError('action %s must have a name' % action)

    # get action name
    name = action[ACTION_NAME_FIELD]

    action_fn = get_action(name, cached_action)

    try:
        # call related action with input event and action such as a kwargs
        result = action_fn(ctx, event, **action)
    except Exception as e:
        raise ActionError(e)

    return result
