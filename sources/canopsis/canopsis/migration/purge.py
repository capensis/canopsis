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

from canopsis.confng import Configuration, Json
from canopsis.logger import Logger
from canopsis.migration.manager import MigrationModule
from canopsis.old.account import Account
from canopsis.old.storage import Storage

DEFAULT_COLLECTIONS = ["cache", "events", "events_log", "object", "userpreferences"]


class PurgeModule(MigrationModule):

    CONF_PATH = 'etc/migration/purge.conf'
    CATEGORY = 'PURGE'

    def __init__(self, collections=None, *args, **kwargs):
        super(PurgeModule, self).__init__(*args, **kwargs)

        self.logger = Logger.get('migrationmodule', MigrationModule.LOG_PATH)
        self.config = Configuration.load(PurgeModule.CONF_PATH, Json)
        conf = self.config.get(self.CATEGORY, {})

        self.storage = Storage(account=Account(user='root', group='root'))

        if collections is not None:
            self.collections = collections
        else:
            self.collections = conf.get('collections', DEFAULT_COLLECTIONS)

    def init(self, yes=False):
        for collection in self.collections:
            self.logger.info(u'Drop collection: {0}'.format(collection))
            self.storage.drop_namespace(collection)

    def update(self, yes=False):
        pass
