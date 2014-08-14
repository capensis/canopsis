# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.mongo import MongoStorage
from canopsis.storage.timedtyped import TimedTypedStorage
from canopsis.timeserie.timewindow import get_offset_timewindow


class MongoTimedTypedStorage(MongoStorage, TimedTypedStorage):
    """
    MongoStorage dedicated to manage timed typed data.
    """

    __datatype__ = 'timedtyped'  #: register this class to timed data types

    class Key:

        DATA_ID = 'd'
        VALUE = 'v'
        TIMESTAMP = 't'
        TYPE = 'k'

    TYPE_BY_TIMESTAMP_BY_ID = \
        [
            (Key.DATA_ID, MongoStorage.ASC),
            (Key.TIMESTAMP, MongoStorage.DESC),
            (Key.TYPE, MongoStorage.ASC)]
    TIMESTAMP_BY_TYPE_BY_ID = \
        [
            (Key.DATA_ID, MongoStorage.ASC),
            (Key.TYPE, MongoStorage.ASC),
            (Key.TIMESTAMP, MongoStorage.DESC)
        ]
    TIMESTAMP_BY_TYPE = \
        [
            (Key.TYPE, MongoStorage.ASC),
            (Key.TIMESTAMP, MongoStorage.DESC)
        ]

    def get(
        self, data_ids=None, data_type=None, timewindow=None,
        limit=0, skip=0, sort=None, *args, **kwargs
    ):

        result = dict()

        # set a where clause for the search
        where = dict()

        if data_ids is not None:
            where[TimedTypedStorage.Key.DATA_ID] = {'$in': data_ids}

        if data_type is not None:
            where[TimedTypedStorage.Key.TYPE] = data_type

        # if timewindow is not None, get latest timestamp before
        # timewindow.stop()
        if timewindow is not None:
            timestamp = timewindow.stop()
            where[TimedTypedStorage.Key.TIMESTAMP] = \
                {'$lte': timewindow.stop()}

        # do the query
        cursor = self._find(document=where)

        # if timewindow is None or contains only one point, get only last
        # document respectively before now or before the one point
        if limit:
            cursor.limit(limit)
        if skip:
            cursor.skip(skip)
        if sort is not None:
            MongoStorage._update_sort(sort)
            cursor.sort(sort)

        # apply a specific index
        if data_type is None:
            cursor.hint(TimedTypedStorage.TYPE_BY_TIMESTAMP_BY_ID)

        elif data_ids is None:
            cursor.hint(TimedTypedStorage.TIMESTAMP_BY_TYPE)

        else:
            cursor.hint(TimedTypedStorage.TIMESTAMP_BY_TYPE_BY_ID)

        # iterate on all documents
        for document in cursor:
            timestamp = document[TimedTypedStorage.Key.TIMESTAMP]
            value = document[TimedTypedStorage.Key.VALUE]
            data_id = document[TimedTypedStorage.Key.DATA_ID]

            # a value to get is composed of a timestamp, values and document id
            value_to_append = (timestamp, value, document['_id'])

            if data_id not in result:
                result[data_id] = [value_to_append]

            else:
                result[data_id].append(value_to_append)

            if timewindow is not None and timestamp not in timewindow:
                # stop when a document is just before the start timewindow
                break

        return result

    def count(self, data_id=None, data_type=None, *args, **kwargs):

        query = dict()

        if data_id is not None:
            query[TimedTypedStorage.Key.DATA_ID] = data_id

        if data_type is not None:
            query[TimedTypedStorage.Key.TYPE] = data_type

        cursor = self._find(document=query)

        # apply a specific index
        if data_type is None:
            cursor.hint(TimedTypedStorage.TYPE_BY_TIMESTAMP_BY_ID)

        elif data_id is None:
            cursor.hint(TimedTypedStorage.TIMESTAMP_BY_TYPE)

        else:
            cursor.hint(TimedTypedStorage.TIMESTAMP_BY_TYPE_BY_ID)

        result = cursor.count()

        return result

    def put(self, data_id, data_type, value, timestamp, *args, **kwargs):

        timewindow = get_offset_timewindow(offset=timestamp)

        data = self.get(
            data_ids=[data_id], data_type=data_type, timewindow=timewindow,
            limit=1)

        data_value = None

        if data:
            data_value = data[data_id][0][TimedTypedStorage.Index.VALUE]

        if value != data_value:  # new entry to insert

            values_to_insert = {
                    TimedTypedStorage.Key.DATA_ID: data_id,
                    TimedTypedStorage.Key.TIMESTAMP: timestamp,
                    TimedTypedStorage.Key.VALUE: value,
                    TimedTypedStorage.Key.TYPE: data_type
            }
            self._insert(document=values_to_insert)

    def remove(
        self, data_ids=None, data_type=None, timewindow=None, *args, **kwargs
    ):

        query = dict()

        if data_ids is not None:
            query[TimedTypedStorage.Key.DATA_ID] = {'$in': data_ids}

        if data_type is not None:
            query[TimedTypedStorage.Key.TYPE] = data_type

        if timewindow is not None:
            query[TimedTypedStorage.Key.TIMESTAMP] = \
                {'$gte': timewindow.start(), '$lte': timewindow.stop()}

        self._remove(document=query)

    def _get_indexes(self):

        result = super(TimedTypedStorage, self)._get_indexes()

        result.append(TimedTypedStorage.TYPE_BY_TIMESTAMP_BY_ID)
        result.append(TimedTypedStorage.TIMESTAMP_BY_TYPE)
        result.append(TimedTypedStorage.TIMESTAMP_BY_TYPE_BY_ID)

        return result
