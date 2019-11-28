# --------------------------------
# Copyright (c) 2019 "Capensis" [http://www.capensis.com]
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
The models.dynamic_infos module defines the DynamicInfosRule model, that
represents the rule used by the dynamic-infos go engine to add informations to
the alarms.
"""

from __future__ import unicode_literals
from uuid import uuid4


class DynamicInfosRule(object):
    """A DynamicInfosRule is an object representing a rule that adds dynamic
    informations to the alarms.

    These rules can be managed with the /api/v2/dynamic_infos API, and are used
    by the dynamic-infos go engine.

    :param str id_: A unique id
    :param str name: The name of the rule
    :param str author: The author of the rule
    :param int creation_date: The date the rule was created as a timestamp.
    :param int last_modified_date: The date the rule was last modified as a
        timestamp.
    :param Optional[str] description: An optional description
    :param List[DynamicInfo] infos: The informations that will be added by the
        rule
    :param Optional[List] alarm_patterns: A list of alarm patterns
    :param Optional[List] entity_patterns: A list of entity patterns
    """

    ID = '_id'
    NAME = 'name'
    AUTHOR = 'author'
    CREATION_DATE = 'creation_date'
    LAST_MODIFIED_DATE = 'last_modified_date'
    DESCRIPTION = 'description'
    INFOS = 'infos'
    ALARM_PATTERNS = 'alarm_patterns'
    ENTITY_PATTERNS = 'entity_patterns'

    REQUIRED_FIELDS = [
        NAME, INFOS
    ]
    VALID_FIELDS = frozenset([
        ID, NAME, AUTHOR, CREATION_DATE, LAST_MODIFIED_DATE, DESCRIPTION,
        INFOS, ALARM_PATTERNS, ENTITY_PATTERNS
    ])

    def __init__(self, id_, name, author, creation_date, last_modified_date,
                 description=None, infos=None, alarm_patterns=None,
                 entity_patterns=None):
        self.id = id_
        self.name = name
        self.author = author
        self.creation_date = creation_date
        self.last_modified_date = last_modified_date
        self.description = ""
        self.infos = []
        self.alarm_patterns = alarm_patterns
        self.entity_patterns = entity_patterns

        if description is not None:
            self.description = description
        if infos is not None:
            self.infos = infos

        self._check_valid()

    def __str__(self):
        return '{}'.format(self.id)

    def __repr__(self):
        return '<DynamicInfosRule {}>'.format(self.__str__())

    @classmethod
    def new_from_dict(cls, dynamic_infos_rule, author, date):
        """Create a new DynamicInfosRule from a dictionary.


        :param Dict[str, Any] dynamic_infos_rule: The dynamic info rule as a
            dictionary
        :param str author: The author of the rule
        :param int date: The date the rule was last modified as a timestamp
        :rtype: DynamicInfosRule
        :raises: TypeError, ValueError or KeyError if the DynamicInfosRules is
        not valid.
        """
        if not isinstance(dynamic_infos_rule, dict):
            raise ValueError("expected a dictionnary, not {}".format(
                dynamic_infos_rule))

        try:
            id_ = dynamic_infos_rule[cls.ID]
        except KeyError:
            id_ = str(uuid4())

        for field in cls.REQUIRED_FIELDS:
            if field not in dynamic_infos_rule:
                raise KeyError("the {} field is required".format(field))

        invalid_fields = set(dynamic_infos_rule.keys()) - cls.VALID_FIELDS
        if invalid_fields:
            raise ValueError("the following keys are invalid: {}".format(
                ", ".join(invalid_fields)))

        if not isinstance(dynamic_infos_rule[DynamicInfosRule.INFOS], list):
            raise ValueError(
                "{} should be a list".format(DynamicInfosRule.INFOS))
        infos = [
            DynamicInfo.new_from_dict(info)
            for info in dynamic_infos_rule[DynamicInfosRule.INFOS]
        ]

        return cls(
            id_=id_,
            name=dynamic_infos_rule[cls.NAME],
            author=author,
            creation_date=date,
            last_modified_date=date,
            description=dynamic_infos_rule.get(cls.DESCRIPTION),
            infos=infos,
            alarm_patterns=dynamic_infos_rule.get(cls.ALARM_PATTERNS),
            entity_patterns=dynamic_infos_rule.get(cls.ENTITY_PATTERNS))

    def _check_valid(self):
        """Check that the DynamicInfosRule is valid.

        :raises: TypeError or ValueError if the DynamicInfosRules is not
        valid.
        """
        if not isinstance(self.id, basestring):
            raise TypeError(
                "{} should be a string".format(DynamicInfosRule.ID))
        if not self.id:
            raise ValueError(
                "{} should not be empty".format(DynamicInfosRule.ID))

        if not isinstance(self.name, basestring):
            raise ValueError(
                "{} should be a string".format(DynamicInfosRule.NAME))

        if not isinstance(self.author, basestring):
            raise ValueError(
                "{} should be a string".format(DynamicInfosRule.AUTHOR))

        if not isinstance(self.creation_date, int):
            raise ValueError(
                "{} should be a timestamp".format(DynamicInfosRule.CREATION_DATE))

        if not isinstance(self.last_modified_date, int):
            raise ValueError(
                "{} should be a timestamp".format(DynamicInfosRule.LAST_MODIFIED_DATE))

        if not isinstance(self.description, basestring):
            raise ValueError(
                "{} should be a string".format(DynamicInfosRule.DESCRIPTION))

        if not isinstance(self.infos, list):
            raise ValueError(
                "{} should be a list".format(DynamicInfosRule.INFOS))
        for info in self.infos:
            if not isinstance(info, DynamicInfo):
                raise ValueError(
                    "{} should only contain DynamicInfos, not {}".format(
                        DynamicInfosRule.INFOS, info))

        if self.alarm_patterns is not None:
            if not isinstance(self.alarm_patterns, list):
                raise ValueError(
                    "{} should be a list or null".format(
                        DynamicInfosRule.ALARM_PATTERNS))
            for pattern in self.alarm_patterns:
                if not isinstance(pattern, dict):
                    raise ValueError(
                        "{} should only contain dictionaries, not {}".format(
                            DynamicInfosRule.ALARM_PATTERNS, pattern))

        if self.entity_patterns is not None:
            if not isinstance(self.entity_patterns, list):
                raise ValueError(
                    "{} should be a list or null".format(
                        DynamicInfosRule.ENTITY_PATTERNS))
            for pattern in self.entity_patterns:
                if not isinstance(pattern, dict):
                    raise ValueError(
                        "{} should only contain dictionaries, not {}".format(
                            DynamicInfosRule.ENTITY_PATTERNS, pattern))

    def as_dict(self):
        """Return the DynamicInfosRule as a dictionnary that can be stored in
        MongoDB.

        :rtype: Dict[str, Any]
        """
        return {
            DynamicInfosRule.ID: self.id,
            DynamicInfosRule.NAME: self.name,
            DynamicInfosRule.AUTHOR: self.author,
            DynamicInfosRule.CREATION_DATE: self.creation_date,
            DynamicInfosRule.LAST_MODIFIED_DATE: self.last_modified_date,
            DynamicInfosRule.DESCRIPTION: self.description,
            DynamicInfosRule.INFOS: [info.as_dict() for info in self.infos],
            DynamicInfosRule.ALARM_PATTERNS: self.alarm_patterns,
            DynamicInfosRule.ENTITY_PATTERNS: self.entity_patterns,
        }


class DynamicInfo(object):
    """A DynamicInfo is an object representing an information that can be added
    to an alarm by a DynamicInfoRule.

    :param str name: The name of the information
    :param str value: The value or the information
    """
    NAME = 'name'
    VALUE = 'value'
    FIELDS = frozenset([
        NAME, VALUE
    ])

    def __init__(self, name, value):
        self.name = name
        self.value = value

        self._check_valid()

    def __str__(self):
        return '{}'.format(self.name)

    def __repr__(self):
        return '<DynamicInfo {}>'.format(self.__str__())

    @classmethod
    def new_from_dict(cls, info):
        """Create a new DynamicInfo from a dictionary.

        :rtype: DynamicInfo
        :raises: TypeError, ValueError or KeyError if the DynamicInfo is not
        valid.
        """
        for field in cls.FIELDS:
            if field not in info:
                raise KeyError(
                    'each element of the infos array should have a field '
                    '"{}"'.format(field))

        invalid_fields = set(info.keys()) - cls.FIELDS
        if invalid_fields:
            raise ValueError(
                "the following fields are invalid in the elements of the "
                "infos array: {}".format(", ".join(invalid_fields)))

        return cls(
            name=info[cls.NAME],
            value=info[cls.VALUE])

    def _check_valid(self):
        """Check that the DynamicInfo is valid.

        :raises: TypeError or ValueError if the DynamicInfosRules is not
        valid.
        """
        if not isinstance(self.name, basestring):
            raise TypeError(
                "{} should be a string".format(DynamicInfo.NAME))
        if not self.name:
            raise ValueError(
                "{} should not be empty".format(DynamicInfo.NAME))

        if not isinstance(self.value, basestring):
            raise TypeError(
                "{} should be a string".format(DynamicInfo.VALUE))

    def as_dict(self):
        """Return the DynamicInfo as a dictionnary that can be stored in
        MongoDB.

        :rtype: Dict[str, str]
        """
        return {
            DynamicInfo.NAME: self.name,
            DynamicInfo.VALUE: self.value,
        }
