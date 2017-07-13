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
import logging

from canopsis.common.ini_parser import IniParser
from canopsis.common.associative_table.associative_table import AssociativeTable
from canopsis.middleware.core import Middleware
#from canopsis.configuration.configurable.decorator import (
#    conf_paths, add_category
#)
#from canopsis.middleware.registry import MiddlewareRegistry
#from canopsis.configuration.model import Parameter

CONF_PATH = 'common/associative_table.conf'
AT_CAT = 'ASSOCIATIVE_TABLE'
AT_K_STORAGE = 'associative_table_storage_uri'


class AssociativeTableManager():
    """
    AssociativeTable, grouped in a single collection and indexed with a table
    name.
    """

    def __init__(self, storage=None, *args, **kwargs):
        self.logger = logging.getLogger('common')

        self.config = IniParser(path=CONF_PATH, logger=self.logger)

        self.storage = storage
        if storage is None:
            section = self.config.get(AT_CAT)
            if AT_K_STORAGE in section:
                self.storage = Middleware.get_middleware_by_uri(
                    section[AT_K_STORAGE]
                )
            else:
                self.logger.error('Cannot read {} parameter in configuration'
                                  .format(AT_K_STORAGE))

    def get(self, table_name):
        """
        Search for this table name in the collection.

        :param str table_name: the table name
        :rtype: <AssociativeTable>
        """
        query = {
            table_name: {"$exists": True}
        }
        table = self.storage._backend.find(query)

        if table.count() > 0:
            content = list(table.limit(1))[0]
            content.pop('_id')
            return AssociativeTable(table_name=table_name, content=content)

        self.logger.info('Impossible to find associative table "{}". '
                         'Creating new one...'.format(table_name))
        base = {
            table_name: {}
        }
        self.storage._backend.insert(base)

        return AssociativeTable(table_name=table_name, content={})

    def save(self, atable):
        """
        Write stored list.

        :param object atable: the table to update
        :returns: mongo response
        """
        return self.storage._backend.update(atable.table_name, atable.content)
