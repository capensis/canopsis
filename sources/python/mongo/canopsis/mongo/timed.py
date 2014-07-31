#!/usr/bin/env python
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

from canopsis.mongo import Storage
from canopsis.storage.timed import TimedStorage
from canopsis.timeserie.timewindow import get_offset_timewindow


class TimedStorage(Storage, TimedStorage):

    class Key:

        DATA_ID = 'd'
        VALUE = 'v'
        TIMESTAMP = 't'

    TIMESTAMP_BY_ID = \
        [(Key.DATA_ID, Storage.ASC), (Key.TIMESTAMP, Storage.DESC)]

    def get(
        self, data_ids, timewindow=None, limit=0, skip=0, sort=None,
        *args, **kwargs
    ):

        result = dict()

        # set a where clause for the search
        where = {
                    TimedStorage.Key.DATA_ID: {'$in': data_ids}
                }

        # if timewindow is not None, get latest timestamp before
        # timewindow.stop()
        if timewindow is not None:
            timestamp = timewindow.stop()
            where[TimedStorage.Key.TIMESTAMP] = {'$lte': timewindow.stop()}

        # do the query
        cursor = self._find(document=where)

        # if timewindow is None or contains only one point, get only last
        # document respectively before now or before the one point
        if limit:
            cursor.limit(limit)
        if skip:
            cursor.skip(skip)
        if sort is not None:
            Storage._update_sort(sort)
            cursor.sort(sort)

        # apply a specific index
        cursor.hint(TimedStorage.TIMESTAMP_BY_ID)

        # iterate on all documents
        for document in cursor:
            timestamp = document[TimedStorage.Key.TIMESTAMP]
            value = document[TimedStorage.Key.VALUE]
            data_id = document[TimedStorage.Key.DATA_ID]

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

    def count(self, data_id, *args, **kwargs):

        query = {
            TimedStorage.Key.DATA_ID: data_id
        }

        cursor = self._find(document=query)
        cursor.hint(TimedStorage.TIMESTAMP_BY_ID)
        result = cursor.count()

        return result

    def put(self, data_id, value, timestamp, *args, **kwargs):

        timewindow = get_offset_timewindow(offset=timestamp)

        data = self.get(
            data_ids=[data_id], timewindow=timewindow, limit=1)

        data_value = None

        if data:
            data_value = data[data_id][0][TimedStorage.Index.VALUE]

        if value != data_value:  # new entry to insert

            values_to_insert = {
                    TimedStorage.Key.DATA_ID: data_id,
                    TimedStorage.Key.TIMESTAMP: timestamp,
                    TimedStorage.Key.VALUE: value
            }
            self._insert(document=values_to_insert)

    def remove(self, data_ids, timewindow=None, *args, **kwargs):

        where = {
            TimedStorage.Key.DATA_ID: {'$in': data_ids}
        }

        if timewindow is not None:
            where[TimedStorage.Key.TIMESTAMP] = \
                {'$gte': timewindow.start(), '$lte': timewindow.stop()}

        self._remove(document=where)

    def _get_indexes(self):

        result = super(TimedStorage, self)._get_indexes()

        result.append(TimedStorage.TIMESTAMP_BY_ID)

        return result
