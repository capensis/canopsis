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
from canopsis.common.ini_parser import IniParser
from canopsis.common.utils import is_mongo_successfull
from canopsis.middleware.core import Middleware

CONF_PATH = 'common/associative_table.conf'
AT_CAT = 'ASSOCIATIVE_TABLE'
AT_K_STORAGE = 'associative_table_storage_uri'

NAME = 'name'
CONTENT = 'content'


class AssociativeTableManager():
    """
    AssociativeTable, grouped in a single collection and indexed with a table
    name.
    """

    def __init__(self, logger, storage=None, path=CONF_PATH, *args, **kwargs):
        self.logger = logger
        self.config = IniParser(path=path, logger=self.logger)

        if storage is None:
            section = self.config.get(AT_CAT)
            if AT_K_STORAGE in section:
                self.storage = Middleware.get_middleware_by_uri(
                    section[AT_K_STORAGE]
                )
            else:
                self.logger.error('Cannot read {} parameter in configuration'
                                  .format(AT_K_STORAGE))
        else:
            self.storage = storage

    def create(self, table_name):
        """
        Create a new AssociativeTable object

        :param str table_name: the table name
        :rtype: <AssociativeTable>
        """
        base = {
            NAME: table_name,
            CONTENT: {}
        }
        self.logger.info('Creating associative table "{}".'.format(table_name))
        self.storage._backend.insert(base)

        return AssociativeTable(table_name=table_name, content={})

    def get(self, table_name):
        """
        Search for this table name in the collection.

        :param str table_name: the table name
        :rtype: <AssociativeTable> or None
        """
        query = {
            NAME: {"$eq": table_name}
        }
        table = self.storage._backend.find(query)

        if table.count() > 0:
            content = list(table.limit(1))[0].get(CONTENT, {})
            return AssociativeTable(table_name=table_name,
                                    content=content)

        self.logger.info('Impossible to find associative table "{}".'
                         .format(table_name))
        return None

    def save(self, atable):
        """
        Update an AssociativeTable in db.

        :param object atable: the table to update
        :rtype: bool
        """
        find = {NAME: {"$eq": atable.table_name}}
        update = {
            NAME: atable.table_name,
            CONTENT: atable.get_all()
        }
        mongo_dict = self.storage._backend.update(find, update)

        return is_mongo_successfull(mongo_dict)

    def delete(self, table_name):
        """
        Delete an associative table object.

        :param str table_name: the name of the table.
        :rtype: bool
        """
        query = {
            NAME: {"$eq": table_name}
        }
        self.logger.info('Deleting associative table: {}'.format(table_name))
        result = self.storage._backend.remove(query)

        return is_mongo_successfull(result)
