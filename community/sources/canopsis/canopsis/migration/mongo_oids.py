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

from canopsis.logger import Logger
from canopsis.common.middleware import Middleware
from canopsis.migration.manager import MigrationModule

DEFAULT_STORAGE = 'storage://'


class MongoOIDsModule(MigrationModule):

    def __init__(self, storage=DEFAULT_STORAGE, *args, **kwargs):
        super(MongoOIDsModule, self).__init__(*args, **kwargs)

        self.logger = Logger.get('migrationmodule', MigrationModule.LOG_PATH)

        self.storage = Middleware.get_middleware_by_uri(storage)
        self.storage.connect()

    def init(self):
        pass

    def update(self):
        if self.get_version('mongo_oids') < 1:
            self.logger.info(u'Migrating to version 1')

            self.update_to_version_1()
            self.set_version('mongo_oids', 1)

        if self.get_version('mongo_oids') < 2:
            self.logger.info(u'Migrating to version 2')
            self.update_periodic2periodical()
            self.set_version('mongo_oids', 2)

    def update_periodic2periodical(self):
        """Change collection names of periodic to timed and timed to periodical."""

        db = self.storage._database

        # dict of collections by name
        collections2rename = {}
        hasperiodic = False

        for collectionname in db.collection_names():

            if collectionname.startswith('periodic_'):
                hasperiodic = True

                collections2rename[collectionname] = db.get_collection(collectionname)

            if collectionname.startswith('timed_'):

                collections2rename[collectionname] = db.get_collection(collectionname)

        if hasperiodic:

            for name in collections2rename:
                collection = collections2rename[name]

                newname = name.replace('timed_', 'periodical_')
                newname = newname.replace('periodic_', 'timed_')

                collection.rename(newname)

    def update_to_version_1(self):
        collections = self.storage._database.collection_names(
            include_system_collections=False
        )

        # exclude GridFS collections
        collections = filter(
            lambda item: any([
                not item.endswith('.chunks'),
                not item.endswith('.files'),
                '{0}.chunks'.format(item) not in collections,
                '{0}.files'.format(item) not in collections
            ]),
            collections
        )

        for collection in collections:
            self.storage.table = collection

            cursor = self.storage.find_elements()

            docs = [
                doc
                for doc in cursor
                if not isinstance(doc[self.storage.ID], basestring)
            ]

            self.logger.info(
                '-- collection: {0} (documents to migrate: {1})'.format(
                    collection,
                    len(docs)
                )
            )

            oids = []
            for doc in docs:
                oids.append(doc[self.storage.ID])
                doc[self.storage.ID] = str(doc[self.storage.ID])

                self.storage.put_element(element=doc)

            self.storage.remove_elements(ids=oids)
