#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

import validictory

from canopsis.old.storage import get_storage
from canopsis.old.account import Account


class NoSchemaError(Exception):
    def __init__(self, schema_id, *args, **kwargs):
        super(NoSchemaError, self).__init__(*args, **kwargs)

        self.schema_id = schema_id

    def __str__(self):
        return 'Schema {0} not found in database'.format(self.schema_id)

    def __unicode__(self):
        return u'Schema {0} not found in database'.format(self.schema_id)


cache = {}


def get(schema_id):
    """
        Get schema from its ID.
        Will look in database if the schema isn't loaded in cache.

        :param schema_id: Schema identifier (value of _id field in Mongo document).
        :type schema_id: str

        :returns: schema field of Mongo document.
    """

    if schema_id not in cache:
        db = get_storage('schemas', account=Account(user='root', group='root')).get_backend()
        doc = db.find_one(schema_id)
        del db

        if not doc:
            raise NoSchemaError(schema_id)

        cache[schema_id] = doc['schema']

    return cache[schema_id]


def validate(dictionary, schema_id):
    """
        Validate a dictionary using a schema.

        :param dictionary: Dictionary to validate.
        :type dictionary: dict

        :param schema_id: Schema identifier (value of _id field in Mongo document).
        :type schema_id: str

        :returns: True if the validation succeed, False otherwise.
        WARNING: disabled, always returns True.
    """

    schema = get(schema_id)

    try:
        validictory.validate(dictionary, schema, required_by_default=False)
        return True

    except validictory.ValidationError:
        return True
