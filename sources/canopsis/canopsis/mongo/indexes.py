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

"""This file is copied to canopsis libs folder, so there should no have direct canopsis import here"""

collections_indexes = {
    'perfdata2': [
        [('co', 1), ('re', 1), ('me', 1)],
        [('re', 1), ('me', 1)],
        [('me', 1)],
        [('tg', 1)]
    ],
    'events': [
        [('connector_name', 1), ('resource', 1), ('component', 1), ('state', 1), ('state_type', 1), ('event_type', 1)],
        [('tags', 1), ('source_type', 1)],
        [('component', 1), ('resource', 1), ('event_type', 1)],
        [('resource', 1), ('event_type', 1)],
        [('event_type', 1)],
        [('crecord_type', 1)],
    ],
    'events_log': [
        [('connector_name', 1),('resource', 1),('component', 1),('state', 1),('event_type', 1)],
        [('component', 1),('resource', 1),('event_type', 1)],
        [('resource', 1),('event_type', 1)],
        [('event_type', 1)],
        [('state_type', 1)],
        [('tags', 1)],
        [('referer', 1)],
        [('timestamp', 1)]
    ],
    'entities': [
        [('type', 1),('name', 1)],
        [('type', 1),('timestamp', 1)],
        [('type', 1),('component', 1),('resource', 1),('id', 1)],
        [('type', 1),('timestamp', 1),('component', 1),('resource', 1)],
        [('type', 1),('nodeid', 1)],
        [('crecord_type', 1),('objclass', 1)]
    ]

}


def init():

    from canopsis.old.account import Account
    from canopsis.old.storage import get_storage

    storage = get_storage(account=Account(user="root", group="root"), namespace='object')

    for collection in collections_indexes:
        storage.get_backend(collection).drop_indexes()
        for index in collections_indexes[collection]:
            storage.get_backend(collection).ensure_index(index)
        logger.info(" + {} Indexe(s) recreated for collection {}".format(len(collections_indexes[collection]), collection))


def update():
    init()
