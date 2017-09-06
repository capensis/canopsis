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
from canopsis.migration.manager import MigrationModule
from canopsis.old.account import Account
from canopsis.old.storage import get_storage


class IndexesModule(object):

    INDEXES = {
        'object': [
            [('crecord_type', 1)],
            [('crecord_type', 1), ('crecord_name', 1)],
        ],
        'events': [
            [
                ('connector_name', 1),
                ('event_type', 1),
                ('component', 1),
                ('resource', 1),
                ('state_type', 1),
                ('state', 1)
            ], [
                ('source_type', 1),
                ('state', 1)
            ], [
                ('event_type', 1),
                ('component', 1),
                ('resource', 1)
            ], [
                ('event_type', 1),
                ('resource', 1)
            ],
            [('event_type', 1)],
            [('source_type', 1)],
            [('domain', 1)],
            [('perimeter', 1)],
            [('connector', 1)],
            [('component', 1)],
            [('resource', 1)],
            [('status', 1)],
            [('state', 1)],
            [('ack', 1)],

        ],
        'events_log': [
            [
                ('connector_name', 1),
                ('event_type', 1),
                ('component', 1),
                ('resource', 1),
                ('state_type', 1),
                ('state', 1)
            ], [('source_type', 1), ('tags', 1)],
            [('event_type', 1), ('component', 1), ('resource', 1)],
            [('rk', 1), ('timestamp', 1)],
            [('event_type', 1), ('resource', 1)],
            [('event_type', 1)],
            [('tags', 1)],
            [('referer', 1)],
            [('event_type', 1)],
            [('domain', 1)],
            [('perimeter', 1)],
            [('connector', 1)],
            [('component', 1)],
            [('resource', 1)],
            [('status', 1)],
            [('state', 1)],
            [('ack', 1)],
        ],
        'ack': [
            [('rk', 1), ('solved', 1)]
        ]
    }

    def __init__(self):
        self.logger = Logger.get('migrationmodule', MigrationModule.LOG_PATH)
        self.storage = get_storage(
            account=Account(user='root', group='root'),
            namespace='object'
        )

    def init(self):
        for collection in IndexesModule.INDEXES:
            self.logger.info(u'Indexing collection: {0}'.format(collection))
            col = self.storage.get_backend(collection)
            col.drop_indexes()

            for index in IndexesModule.INDEXES[collection]:
                col.ensure_index(index)

    def update(self):
        answer = self.ask(
            'Add/Update indexes (update may take time)?',
            default=False
        )

        if answer:
            self.init()
