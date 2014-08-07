============================================================
Rule: package managing condition/actions on event processing
============================================================

.. contents:
    maxdepth: 2

.. module:: canopsis.rule
    :synopsis: rule library for processing events

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Objective
=========

This package manage rules such as couples of (condition, actions).

A condition is a filter such as mongo filter (http://docs.mongodb.org/manual/reference/operator/) which is applied on an input event. If condition is checking, then related actions are performed. An action is a dictionary which contains at least a name (corresponding to a python function) and additional action parameters (event fields to updatein the case of an event updating actionfor example).

Package contents
================

.. data:: __version__ = "0.1"

    current module version

.. class:: RuleError

    Handle rule errors.

.. function:: apply_rule(rule, event, cached_action=True)

    Apply input rule on input event in checking if the rule condition matches
    with the event and if True, execute rule actions.

    :param rule: rule to apply on input event. contains both condition and \
        actions
    :type rule: dict
    :param event: event to check and to process.
    :type event: dict
    :param cached_action: indicates to actions to use cache instead of importing them dynamically
    :type cached_action: bool

    :return: ordered list of rule action results if rule condition is checked.
    :rtype: list

.. module:: canopsis.rule.condition

Manage rule conditions in handling a mongo filter type (http://docs.mongodb.org/manual/reference/operator/).

.. data:: CONDITION_FIELDS = 'condition'

    rule condition field name

.. function:: check(condition, event)

    Check if input condition matches input event.

.. module:: canopsis.rule.action

Manage rule actions in resolving action name with absolute python function path

.. data:: ACTIONS_FIELDS = 'actions'

    rule actions field name

.. data:: ACTION_NAME_FIELD = 'name'

    action field name

.. class:: ActionError

    Handle action execution errors

.. function:: do_action(action, event, cached_action=True)

    Do an action function related to input name.

    An action should take in parameters:
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
