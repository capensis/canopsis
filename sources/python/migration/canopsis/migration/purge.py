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

from canopsis.migration.manager import MigrationModule
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.model import Parameter

from canopsis.old.account import Account
from canopsis.old.storage import Storage


CONF_PATH = 'migration/purge.conf'
CATEGORY = 'PURGE'
CONTENT = [
    Parameter('collections', parser=Parameter.array())
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class PurgeModule(MigrationModule):

    @property
    def collections(self):
        if not hasattr(self, '_collections'):
            self.collections = None

        return self._collections

    @collections.setter
    def collections(self, value):
        if value is None:
            value = []

        self._collections = value

    def __init__(self, collections=None, *args, **kwargs):
        super(PurgeModule, self).__init__(*args, **kwargs)

        self.storage = Storage(account=Account(user='root', group='root'))

        if collections is not None:
            self.collections = collections

    def init(self):
        for collection in self.collections:
            self.logger.info('Drop collection: {0}'.format(collection))
            self.storage.drop_namespace(collection)

    def update(self):
        pass
