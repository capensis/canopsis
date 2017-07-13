# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

from canopsis.common.associative_table.associative_table import AssociativeTable
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.model import Parameter

CONF_PATH = 'common/associative_table.conf'
AT_CAT = 'ASSOCIATIVE_TABLE'

NAME = 'name'
KEY = 'key'
CONTENT = 'content'
BASE = {
    KEY: '',
    CONTENT: {}
}

AT_CONTENT = [
    Parameter('associative_table_storage_uri')
]


@conf_paths(CONF_PATH)
@add_category(AT_CAT, content=AT_CONTENT)
class AssociativeTableManager(MiddlewareRegistry):
    """
    AssociativeTable, grouped in a single collection and indexed with a table
    name.
    """

    ASSOC_STORAGE = 'associative_table_storage_uri'

    def __init__(self, storage=None, *args, **kwargs):
        super(AssociativeTableManager, self).__init__(*args, **kwargs)

        if storage is not None:
            self[AssociativeTableManager.ASSOC_STORAGE] = storage

    @property
    def storage(self):
        """
        Simple access to the storage.
        """
        return self[AssociativeTableManager.ASSOC_STORAGE]

    def get(self, table_name):
        """
        Search for this table name in the collection.

        :param str table_name: the table name
        :rtype: <AssociativeTable>
        """
        query = {
            NAME: {"$eq": table_name}
        }
        table = self.storage._backend.find(query)
        print(table)

        if table.count() == 0:
            print('Impossible to find associative table {}. Creating one...'
                  .format(table_name))
            base = {NAME: BASE}
            base[NAME][KEY] = table_name
            table = self.storage._backend.insert_one(base)
            print(table)

        return AssociativeTable(table_name=table_name, content=table)

    def save(self, atable):
        """
        Write stored list.

        :param object atable: the table to update
        :returns: mongo response
        """
        _id = {NAME: atable.table_name}

        return self.storage._backend.update(_id, atable.content)
