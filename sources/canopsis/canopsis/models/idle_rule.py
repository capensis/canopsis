# --------------------------------
# Copyright (c) 2020 "Capensis" [http://www.capensis.com]
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
The models.idle_rule module defines the IdleRule model, that
represents the rule used by the axe go engine to change state or cancel idle alarms.
"""

from __future__ import unicode_literals
from uuid import uuid4


class IdleRule(object):
    """A IdleRule is an object representing a rule that changes idle alarms.

    These rules can be managed with the /api/v2/idle_rule API, and are used
    by the axe go engine.

    :param str id_: A unique id
    :param str name: The name of the rule
    :param str type: The name of the rule
    :param str duration: The idle duration
    :param dict operation: The operation of the rule
    :param str author: The author of the rule
    :param int creation_date: The date the rule was created as a timestamp.
    :param int last_modified_date: The date the rule was last modified as a
        timestamp.
    :param Optional[str] description: An optional description
    :param Optional[List] alarm_patterns: A list of alarm patterns
    :param Optional[List] entity_patterns: A list of entity patterns
    """

    ID = '_id'
    NAME = 'name'
    AUTHOR = 'author'
    CREATION_DATE = 'creation_date'
    LAST_MODIFIED_DATE = 'last_modified_date'
    DESCRIPTION = 'description'
    TYPE = 'type'
    DURATION = 'duration'
    OPERATION = 'operation'
    ALARM_PATTERNS = 'alarm_patterns'
    ENTITY_PATTERNS = 'entity_patterns'

    REQUIRED_FIELDS = [
        TYPE, DURATION, OPERATION
    ]
    VALID_FIELDS = frozenset([
        ID, NAME, AUTHOR, CREATION_DATE, LAST_MODIFIED_DATE, DESCRIPTION,
        TYPE, DURATION, OPERATION, ALARM_PATTERNS, ENTITY_PATTERNS
    ])

    def __init__(self, id_, name, rule_type, duration, operation, author, creation_date, last_modified_date,
                 description=None, alarm_patterns=None, entity_patterns=None):
        self.id = id_
        self.name = name
        self.author = author
        self.rule_type = rule_type
        self.duration = duration
        self.operation = operation
        self.creation_date = creation_date
        self.last_modified_date = last_modified_date
        self.description = ""
        self.alarm_patterns = alarm_patterns
        self.entity_patterns = entity_patterns

        if description is not None:
            self.description = description

        self._check_valid()

    def __str__(self):
        return '{}'.format(self.id)

    def __repr__(self):
        return '<IdleRule {}>'.format(self.__str__())

    @classmethod
    def new_from_dict(cls, idle_rule, author, date):
        """Create a new IdleRule from a dictionary.


        :param Dict[str, Any] idle_rule: The idle rule as a
            dictionary
        :param str author: The author of the rule
        :param int date: The date the rule was last modified as a timestamp
        :rtype: IdleRule
        :raises: TypeError, ValueError or KeyError if the IdleRules is
        not valid.
        """
        if not isinstance(idle_rule, dict):
            raise ValueError("expected a dictionary, not {}".format(
                idle_rule))

        try:
            id_ = idle_rule[cls.ID]
        except KeyError:
            id_ = str(uuid4())

        for field in cls.REQUIRED_FIELDS:
            if field not in idle_rule:
                raise KeyError("the {} field is required".format(field))

        invalid_fields = set(idle_rule.keys()) - cls.VALID_FIELDS
        if invalid_fields:
            raise ValueError("the following keys are invalid: {}".format(
                ", ".join(invalid_fields)))

        return cls(
            id_=id_,
            name=idle_rule[cls.NAME],
            rule_type=idle_rule[cls.TYPE],
            duration=idle_rule[cls.DURATION],
            operation=idle_rule[cls.OPERATION],
            author=author,
            creation_date=date,
            last_modified_date=date,
            description=idle_rule.get(cls.DESCRIPTION),
            alarm_patterns=idle_rule.get(cls.ALARM_PATTERNS),
            entity_patterns=idle_rule.get(cls.ENTITY_PATTERNS))

    def _check_valid(self):
        """Check that the IdleRule is valid.

        :raises: TypeError or ValueError if the IdleRules is not
        valid.
        """
        if not isinstance(self.id, basestring):
            raise TypeError(
                "{} should be a string".format(IdleRule.ID))
        if not self.id:
            raise ValueError(
                "{} should not be empty".format(IdleRule.ID))

        for val, field in (
                (self.rule_type, IdleRule.TYPE), (self.duration, IdleRule.DURATION)):
            if not isinstance(val, basestring):
                raise ValueError(
                    "{} should be a string".format(field))
            if not val:
                raise ValueError(
                    "{} should not be empty".format(field))

        if not self.operation:
            raise ValueError(
                "{} should not be empty".format(IdleRule.OPERATION))

        for val, field in (
                (self.name, IdleRule.NAME), (self.author, IdleRule.AUTHOR), (self.description, IdleRule.DESCRIPTION)):
            if not isinstance(val, basestring):
                raise ValueError(
                    "{} should be a string".format(field))

        if not isinstance(self.creation_date, int):
            raise ValueError(
                "{} should be a timestamp".format(IdleRule.CREATION_DATE))

        if not isinstance(self.last_modified_date, int):
            raise ValueError(
                "{} should be a timestamp".format(IdleRule.LAST_MODIFIED_DATE))

        if self.alarm_patterns is not None:
            if not isinstance(self.alarm_patterns, list):
                raise ValueError(
                    "{} should be a list or null".format(
                        IdleRule.ALARM_PATTERNS))
            for pattern in self.alarm_patterns:
                if not isinstance(pattern, dict):
                    raise ValueError(
                        "{} should only contain dictionaries, not {}".format(
                            IdleRule.ALARM_PATTERNS, pattern))

        if self.entity_patterns is not None:
            if not isinstance(self.entity_patterns, list):
                raise ValueError(
                    "{} should be a list or null".format(
                        IdleRule.ENTITY_PATTERNS))
            for pattern in self.entity_patterns:
                if not isinstance(pattern, dict):
                    raise ValueError(
                        "{} should only contain dictionaries, not {}".format(
                            IdleRule.ENTITY_PATTERNS, pattern))

    def as_dict(self):
        """Return the IdleRule as a dictionary that can be stored in
        MongoDB.

        :rtype: Dict[str, Any]
        """
        return {
            IdleRule.ID: self.id,
            IdleRule.NAME: self.name,
            IdleRule.TYPE: self.rule_type,
            IdleRule.DURATION: self.duration,
            IdleRule.OPERATION: self.operation,
            IdleRule.AUTHOR: self.author,
            IdleRule.CREATION_DATE: self.creation_date,
            IdleRule.LAST_MODIFIED_DATE: self.last_modified_date,
            IdleRule.DESCRIPTION: self.description,
            IdleRule.ALARM_PATTERNS: self.alarm_patterns,
            IdleRule.ENTITY_PATTERNS: self.entity_patterns,
        }

