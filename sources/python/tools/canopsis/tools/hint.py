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

from canopsis.mongo.indexes import collections_indexes
import json


def get_keys(query_chunk):
    #recursive search in dictionary for keys that are not mongo keywords
    keys_found = []
    for key in query_chunk:
        if not key.startswith('$'):
            keys_found.append(key)
        if type(query_chunk[key]) == dict:
            keys_found += get_keys(query_chunk[key])
        if type(query_chunk[key]) == list:
            for sub_query_item in query_chunk[key]:
                keys_found += get_keys(sub_query_item)
    return keys_found


def get_hint(collection, query_keys):
    #test if collection gets indexes
    if collection in collections_indexes:
        #seek if all query field exists in any index
        for index in collections_indexes[collection]:
            match = True
            index_keys = [key[0] for key in index]
            if len(query_keys) == len(index_keys):
                for query_key in query_keys:
                    if query_key not in index_keys:
                        match = False
                if match:
                    return index


def set_hint(mfilter, collection, cursor, storage, logger):

    query_keys = list(set(get_keys(mfilter)))

    indexes = get_hint(collection, query_keys)

    if indexes:
        try:
            cursor.hint(indexes)
            logger.debug('Indexes found will use hint on cursor')
        except Exception as e:
            logger.warning('Unable to set hint for cursor {}'.format(e))
    else:
        logger.info('No index found for this query, will establish states')

        statistics = storage.get_backend('indexes_statistics')
        statistics.update({'_id': '.'.join(query_keys)},
            {'$set': {
                'fields': query_keys,
                'last_filter': json.dumps(mfilter)
            }, '$inc': {'count': 1}
            }, upsert=True)
