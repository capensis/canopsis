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

from .condition import check, CONDITION_FIELD
from .action import do_action, ACTIONS_FIELD


class RuleError(Exception):
    """
    Handle rule errors.
    """

    pass


def apply_rule(rule, event, cached_action=True):
    """
    Apply input rule on input event in checking if the rule condition matches
    with the event and if True, execute rule actions.

    :param rule: rule to apply on input event. contains both condition and \
        actions.
    :type rule: dict
    :param event: event to check and to process.
    :type event: dict
    :param cached_action: indicates to actions to use cache instead of \
        importing them dynamically.
    :type cached_action: bool

    :return: ordered list of rule action results if rule condition is checked.
    :rtype: list
    """

    result = []

    condition = rule[CONDITION_FIELD] if CONDITION_FIELD in rule else None

    if condition is None or check(condition, event):

        actions = rule[ACTIONS_FIELD] if ACTIONS_FIELD in rule else []

        for action in actions:
            action_result = do_action(action, event, cached_action)

            result.append(action_result)

    return result
