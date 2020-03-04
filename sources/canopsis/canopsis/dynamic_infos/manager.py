# -*- coding: utf-8 -*-
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
The dynamic_infos.manager module defines the DynamicInfosManager class, which
can be used to manage DynamicInfosRule.

These rules are used by the dynamic-infos go engines, to add informations to
the alarms.
"""

from __future__ import unicode_literals
import re

from canopsis.common.collection import MongoCollection
from canopsis.common.errors import NotFoundError
from canopsis.common.mongo_store import MongoStore
from canopsis.logger import Logger
from canopsis.models.dynamic_infos import DynamicInfo, DynamicInfosRule


class DynamicInfosManager(object):
    """The DynamicInfosManager allows to create, update or remove
    DynamicInfosRules from the database."""
    LOG_PATH = 'var/log/dynamic_infos.log'
    COLLECTION = 'dynamic_infos'

    SEARCHABLE_FIELDS = [
        DynamicInfosRule.ID, DynamicInfosRule.NAME,
        DynamicInfosRule.DESCRIPTION,
        DynamicInfosRule.AUTHOR,
        "{}.{}".format(DynamicInfosRule.INFOS, DynamicInfo.NAME),
        "{}.{}".format(DynamicInfosRule.INFOS, DynamicInfo.VALUE),
    ]

    def __init__(self, logger, mongo_collection):
        self.logger = logger
        self.collection = mongo_collection

    @classmethod
    def provide_default_basics(cls):
        """Provide logger and collection.

        ! Do not use in tests !

        :rtype: Tuple[logging.Logger,
                      canopsis.common.collection.MongoCollection]
        """
        logger = Logger.get('dynamic_infos', cls.LOG_PATH)
        store = MongoStore.get_default()
        collection = store.get_collection(name=cls.COLLECTION)
        mongo_collection = MongoCollection(collection)

        return (logger, mongo_collection)

    def count(self, search="", search_fields=None):
        """Return the number of DynamicInfosRules.

        If search is defined, only the rules where one of the fields listed in
        search_fields contains the value of search will be counted.

        :param Optional[str] search: The value to search for
        :param Optional[List[str]] search_fields: The names of the fields to
            search on. The searchable fields are id, name, description,
            infos.name and infos.value. The search uses all of those fields by
            default.
        :rtype: int
        """
        query = self._get_search_query(search, search_fields)
        return self.collection.find(query).count()

    def list(self, search="", search_fields=None, limit=0, offset=0):
        """Return a list of the DynamicInfosRules as dictionaries.

        If search is defined, only the rules where one of the fields listed in
        search_fields contains the value of search will be returned.

        :param Optional[str] search: The value to search for
        :param Optional[List[str]] search_fields: The names of the fields to
            search on. The searchable fields are id, name, description,
            infos.name and infos.value. The search uses all of those fields by
            default.
        :rtype: List[Dict[str, Any]]
        """
        query = self._get_search_query(search, search_fields)
        return list(self.collection.find(query).limit(limit).skip(offset))

    def _get_search_query(self, search="", search_fields=None):
        """Return a MongoDB query used to search DynamicInfosRules.

        :param Optional[str] search: The value to search for
        :param Optional[List[str]] search_fields: The names of the fields to
            search on. The searchable fields are id, name, description,
            infos.name and infos.value. The search uses all of those fields by
            default.
        :rtype: Dict[str, Any]
        :raises: ValueError if search_fields contains an invalid field
        """
        if not search:
            return {}

        if not search_fields:
            search_fields = DynamicInfosManager.SEARCHABLE_FIELDS

        invalid_fields = (
            set(search_fields) - set(DynamicInfosManager.SEARCHABLE_FIELDS))
        if invalid_fields:
            raise ValueError(
                "the following search fields are invalid: {}".format(
                    ", ".join(invalid_fields)))

        escaped_search = re.escape(search)
        return {
            "$or": [
                {field: {'$regex': escaped_search, '$options': 'i'}}
                for field in search_fields
            ]
        }

    def get_by_id(self, rule_id):
        """Return a rule as a dictionary given its id.

        If the rule does not exist, None is returned instead.

        :param str rule_id:
        :rtype: Dict[str, Any]
        """
        return self.collection.find_one({
            '_id': rule_id
        })

    def create(self, rule):
        """Create a new rule.

        :param DynamicInfosRule rule:
        :raises: ValueError if there is already a rule with the same if as the
            one provided in the parameters.
        """
        if self.get_by_id(rule.id) is not None:
            raise ValueError("duplicate id {}".format(rule.id))

        self.collection.insert(rule.as_dict())

    def update(self, rule_id, rule):
        """Update an existing rule.

        :param str rule_id:
        :param DynamicInfosRule rule:
        :raises: NotFoundError if there is no rule with the id provided in the
            parameters.
        :raises: ValueError if the rule_id parameter does not match the id of
            the rule.
        :return bool: A boolean indicating whether the update succeeded.
        """
        previous_value = self.get_by_id(rule_id)
        if previous_value is None:
            raise NotFoundError("no dynamic infos rule with id {}".format(
                rule_id))

        if rule_id != rule.id:
            raise ValueError("the _id field should not be modified")

        if DynamicInfosRule.AUTHOR in previous_value:
            rule.author = previous_value[DynamicInfosRule.AUTHOR]
        if DynamicInfosRule.CREATION_DATE in previous_value:
            rule.creation_date = previous_value[DynamicInfosRule.CREATION_DATE]

        resp = self.collection.update(
            {DynamicInfosRule.ID: rule_id},
            rule.as_dict())
        return self.collection.is_successfull(resp)

    def delete(self, rule_id):
        """Delete an existing rule.

        :param str rule_id:
        :raises: NotFoundError if there is no rule with the id provided in the
            parameters.
        :return bool: A boolean indicating whether the update succeeded.
        """
        if self.get_by_id(rule_id) is None:
            raise NotFoundError("no dynamic infos rule with id {}".format(
                rule_id))

        resp = self.collection.remove({DynamicInfosRule.ID: rule_id})
        return self.collection.is_successfull(resp)
