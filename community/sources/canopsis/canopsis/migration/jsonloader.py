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

import json
import os
from socket import getfqdn

from canopsis.common.utils import ensure_iterable
from canopsis.logger import Logger
from canopsis.migration.manager import MigrationModule
from canopsis.old.account import Account
from canopsis.old.storage import get_storage

DEFAULT_JSON_PATH = '~/opt/mongodb/load.d'


class JSONLoaderModule(MigrationModule):

    def __init__(self, json_path=None, *args, **kwargs):
        super(JSONLoaderModule, self).__init__(*args, **kwargs)

        self.logger = Logger.get('migrationmodule', MigrationModule.LOG_PATH)

        if json_path is not None:
            self.json_path = json_path
        else:
            self.json_path = os.path.expanduser(DEFAULT_JSON_PATH)

        self.storage = get_storage(
            account=Account(user='root', group='root'),
            namespace='object'
        )

    def init(self, yes=False):
        substitutes = [
            ('[[HOSTNAME]]', getfqdn())
        ]

        for name in os.listdir(self.json_path):
            path = os.path.join(self.json_path, name)

            if os.path.isdir(path) and name.startswith('json_'):
                col = name[len('json_'):]

                for docname in os.listdir(path):
                    docpath = os.path.join(path, docname)

                    try:
                        with open(docpath) as f:
                            json_data = f.read()

                            for pattern, value in substitutes:
                                json_data = json_data.replace(pattern, value)

                            data = json.loads(json_data)

                    except Exception as err:
                        self.logger.error(
                            'Unable to load JSON file "{0}": {1}'.format(
                                docname,
                                err
                            )
                        )

                        data = []

                    self.load_documents(data, col, docname)

    def update(self, yes=False):
        self.init()

    def load_documents(self, data, collection, filename):
        storage = self.storage.get_backend(collection)
        data = ensure_iterable(data)

        for doc in data:
            if 'loader_id' not in doc:
                self.logger.error(
                    'Missing "loader_id" key in document, skipping'
                )
                self.logger.debug(str(doc))

                continue

            mfilter = {'loader_id': doc['loader_id']}
            doc_exists = storage.find(mfilter).count()

            if doc_exists:
                if not doc.get('loader_no_update', True):
                    storage.update(mfilter, doc, upsert=True)

                else:
                    self.logger.info(u'Document "{0}" not updatable'.format(
                        doc['loader_id']
                    ))

            else:
                storage.update(mfilter, doc, upsert=True)
