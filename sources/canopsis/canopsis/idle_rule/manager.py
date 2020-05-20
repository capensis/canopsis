# -*- coding: utf-8 -*-

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
The idle_rule.manager module defines the IdleRuleManager class, which
can be used to manage IdleRule.

These rules are used by the axe go engine to change state or cancel idle alarms.
"""

from __future__ import unicode_literals

from canopsis.common.collection import MongoCollection
from canopsis.common.errors import NotFoundError
from canopsis.common.mongo_store import MongoStore
from canopsis.logger import Logger
from canopsis.models.idle_rule import IdleRule


class IdleRuleManager(object):
    """The IdleRuleManager allows to create, update or remove
    IdleRule records from the database."""
    LOG_PATH = 'var/log/idle_rule.log'
    COLLECTION = 'idle_rule'

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
        logger = Logger.get('idle_rule', cls.LOG_PATH)
        store = MongoStore.get_default()
        collection = store.get_collection(name=cls.COLLECTION)
        mongo_collection = MongoCollection(collection)

        return (logger, mongo_collection)

    def read_all(self):
        """Return list of all rules.
        :rtype: list of dictionaries
        """
        return list(self.collection.find({}))

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

        :param IdleRule rule:
        :raises: ValueError if there is already a rule with the same id as the
            one provided in the parameters.
        """
        if rule.id and self.get_by_id(rule.id) is not None:
            raise ValueError("duplicate id {}".format(rule.id))

        self.collection.insert(rule.as_dict())

    def update(self, rule_id, rule):
        """Update an existing rule.

        :param str rule_id:
        :param IdleRule rule:
        :raises: NotFoundError if there is no rule with the id provided in the
            parameters.
        :raises: ValueError if the rule_id parameter does not match the id of
            the rule.
        :return bool: A boolean indicating whether the update succeeded.
        """
        previous_value = self.get_by_id(rule_id)
        if previous_value is None:
            raise NotFoundError("no idle rule with id {}".format(
                rule_id))

        if rule_id != rule.id:
            raise ValueError("the _id field should not be modified")
        if IdleRule.CREATION_DATE in previous_value:
            rule.creation_date = previous_value[IdleRule.CREATION_DATE]

        resp = self.collection.update(
            {IdleRule.ID: rule_id},
            rule.as_dict())
        return self.collection.is_successfull(resp)

    def delete(self, rule_id):
        """Delete an existing rule.

        :param str rule_id:
        :raises: NotFoundError if there is no rule with the id provided in the
            parameters.
        :return bool: A boolean indicating whether the update succeeded.
        """

        def _check_response(response):
            ack = self.collection.is_successfull(response)
            if 'n' in response and response['n'] == 0:
                raise NotFoundError("no idle rule with id {}".format(
                    rule_id))

            return {
                'acknowledged': ack,
                'deletedCount': response['n']
            }

        return _check_response(self.collection.remove({"_id": rule_id}))
