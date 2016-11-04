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
from canopsis.middleware.core import Middleware


CONF_PATH = 'migration/mongo_oids.conf'
CATEGORY = 'MONGO_OIDS'
CONTENT = [
    Parameter('storage')
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class MongoOIDsModule(MigrationModule):

    @property
    def storage(self):
        if not hasattr(self, '_storage'):
            self.storage = None

        return self._storage

    @storage.setter
    def storage(self, value):
        if value is None:
            value = 'storage://'

        if isinstance(value, basestring):
            value = Middleware.get_middleware_by_uri(value)
            value.connect()

        self._storage = value

    def __init__(self, storage=None, *args, **kwargs):
        super(MongoOIDsModule, self).__init__(*args, **kwargs)

        if storage is not None:
            self.storage = storage

    def init(self):
        pass

    def update(self):
        if self.get_version('mongo_oids') < 1:
            self.logger.info('Migrating to version 1')

            self.update_to_version_1()
            self.set_version('mongo_oids', 1)

        if self.get_version('mongo_oids') < 2:
            self.logger.info('Migrating to version 2')
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
