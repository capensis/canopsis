#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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
Manager for event filter rules.
"""

from __future__ import unicode_literals

from uuid import uuid4

from canopsis.common.collection import MongoCollection
from canopsis.common.mongo_store import MongoStore
from canopsis.eventfilter.enums import RuleField, RuleType, RuleOutcome, \
    ENRICHMENT_FIELDS

RULE_COLLECTION = 'eventfilter'


class InvalidRuleError(Exception):
    """
    An InvalidRuleError is an exception that is raised when a rule is invalid.
    """
    def __init__(self, message):
        super(InvalidRuleError, self).__init__(message)
        self.message = message


class RuleManager(object):
    """
    Manager for event filter rules.
    """
    def __init__(self, logger):
        self.logger = logger
        self.rule_collection = MongoCollection(
            MongoStore.get_default().get_collection(RULE_COLLECTION))

    def get_by_id(self, rule_id):
        """
        Get an event filter rule given its id.

        :param str rule_id: the id of the rule.
        :rtype: Dict[str, Any]
        """
        return self.rule_collection.find_one({
            RuleField.id: rule_id
        })

    def create(self, rule):
        """
        Create a new rule and return its id.

        :param Dict[str, Any] rule:
        :rtype: str
        :raises: InvalidRuleError if the rule is invalid. CollectionError if
        the creation fails.
        """
        rule_id = str(uuid4())

        rule[RuleField.id] = rule_id
        self.validate(rule_id, rule)

        self.rule_collection.insert(rule)
        return rule_id

    def remove_with_id(self, rule_id):
        """
        Remove a rule given its id.

        :param str rule_id: the id of the rule. CollectionError if the
        creation fails.
        """
        self.rule_collection.remove({
            RuleField.id: rule_id
        })

    def list(self):
        """
        Return a list of all the rules.

        :rtype: List[Dict[str, Any]]
        """
        return list(self.rule_collection.find({}))

    def update(self, rule_id, rule):
        """
        Update a rule given its id.

        :param str rule_id: the id of the rule.
        :param Dict[str, Any] rule:
        :raises: InvalidRuleError if the rule is invalid. CollectionError if
        the creation fails.
        """
        self.validate(rule_id, rule)

        self.rule_collection.update({
            RuleField.id: rule_id
        }, rule, upsert=False)

    def validate(self, rule_id, rule):
        """
        Check that the rule is valid.

        The pattern and external_data fields are not validated by this method.

        :param Dict[str, Any] view:
        :raises: InvalidRuleError if it is invalid.
        """
        # Validate id field
        if rule.get(RuleField.id, rule_id) != rule_id:
            raise InvalidRuleError(
                'The {0} field should not be modified.'.format(RuleField.id))

        # Check that there are no unexpected fields in the rule
        unexpected_fields = set(rule.keys()).difference(RuleField.values)
        if unexpected_fields:
            raise InvalidRuleError(
                'Unexpected fields: {0}.'.format(', '.join(unexpected_fields)))

        # Validate the description field
        if not isinstance(rule.get(RuleField.description, ""), basestring):
            raise InvalidRuleError(
                'The {0} field should be a string.'.format(
                    RuleField.description))

        # Validate the type field
        if RuleField.type not in rule:
            raise InvalidRuleError(
                'The {0} field is required.'.format(RuleField.type))

        if rule.get(RuleField.type) not in RuleType.values:
            raise InvalidRuleError(
                'The {0} field should be one of: {1}.'.format(
                    RuleField.type,
                    ', '.join(RuleType.values)))

        # Validate the priority field
        if not isinstance(rule.get(RuleField.priority, 0), int):
            raise InvalidRuleError(
                'The {0} field should be an integer.'.format(
                    RuleField.priority))

        # Validate the enabled field
        if not isinstance(rule.get(RuleField.enabled, True), bool):
            raise InvalidRuleError(
                'The {0} field should be a boolean.'.format(
                    RuleField.enabled))

        if rule.get(RuleField.type) != RuleType.enrichment:
            # Check that the enrichment fields are not defined for
            # non-enrichment rules.
            unexpected_fields = set(rule.keys()).intersection(
                ENRICHMENT_FIELDS)
            if unexpected_fields:
                raise InvalidRuleError(
                    'The following fields should only be defined for '
                    'enrichment rules: {0}.'.format(
                        ', '.join(unexpected_fields)))

        else:
            # Validate the actions field of the enrichment rules.
            if RuleField.actions not in rule:
                raise InvalidRuleError(
                    'The {0} field is required for enrichment rules.'.format(
                        RuleField.actions))

            if not isinstance(rule.get(RuleField.actions), list):
                raise InvalidRuleError(
                    'The {0} field should be a list.'.format(
                        RuleField.actions))

            if not rule.get(RuleField.actions):
                raise InvalidRuleError(
                    'The {0} field should contain at least one action.'.format(
                        RuleField.actions))

            # Validate the on_success field of the enrichment rules.
            outcome = rule.get(RuleField.on_success)
            if outcome and outcome not in RuleOutcome.values:
                raise InvalidRuleError(
                    'The {0} field should be one of: {1}.'.format(
                        RuleField.on_success,
                        ', '.join(RuleOutcome.values)))

            # Validate the on_failure field of the enrichment rules.
            outcome = rule.get(RuleField.on_failure)
            if outcome and outcome not in RuleOutcome.values:
                raise InvalidRuleError(
                    'The {0} field should be one of: {1}.'.format(
                        RuleField.on_failure,
                        ', '.join(RuleOutcome.values)))
