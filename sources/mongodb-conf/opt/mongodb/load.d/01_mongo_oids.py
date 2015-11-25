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

from canopsis.middleware.core import Middleware


def init():
    pass


def update():
    print('Migration from ObjectId to str')
    storage = Middleware.get_middleware_by_uri('mongodb://')
    storage.connect()

    collections = storage._database.collection_names(
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
        storage.table = collection

        cursor = storage.find_elements()

        docs = [
            doc
            for doc in cursor
            if not isinstance(doc[storage.ID], basestring)
        ]

        print('-- collection: {0} (documents to migrate: {1})'.format(
            collection, len(docs)
        ))

        oids = []
        for doc in docs:
            oids.append(doc[storage.ID])
            doc[storage.ID] = str(doc[storage.ID])

            storage.put_element(element=doc)

        storage.remove_elements(ids=oids)
