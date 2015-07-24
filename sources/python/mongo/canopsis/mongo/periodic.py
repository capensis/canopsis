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

from canopsis.mongo.core import MongoStorage
from canopsis.storage.periodic import PeriodicStorage
from canopsis.timeserie.timewindow import Period

from md5 import new as md5

from operator import itemgetter

from datetime import datetime

from time import mktime


class MongoPeriodicStorage(MongoStorage, PeriodicStorage):
    """
    MongoStorage dedicated to manage periodic data.
    """

    class Index:

        DATA_ID = 'i'
        TIMESTAMP = 't'
        VALUES = 'v'
        PERIOD = 'p'
        LAST_UPDATE = 'l'

        QUERY = [(DATA_ID, 1), (PERIOD, 1), (TIMESTAMP, 1)]

    def count(self, data_id, period, timewindow=None, *args, **kwargs):

        data = self.get(
            data_id=data_id,
            timewindow=timewindow,
            period=period
        )

        result = len(data)

        return result

    def size(
        self, data_id=None, period=None, timewindow=None,
        *args, **kwargs
    ):

        where = {
            MongoPeriodicStorage.Index.DATA_ID: data_id
        }

        if timewindow is not None:
            where[MongoPeriodicStorage.Index.TIMESTAMP] = {
                '$gte': timewindow.start(),
                '$lte': timewindow.stop()
            }

        if period is not None:
            where[MongoPeriodicStorage.Index.PERIOD] = period

        cursor = self._find(document=where)
        cursor.hint(MongoPeriodicStorage.Index.QUERY)

        result = cursor.count()

        return result

    def get(
        self, data_id, period, timewindow=None, limit=0,
        *args, **kwargs
    ):

        query = self._get_documents_query(
            data_id=data_id,
            timewindow=timewindow,
            period=period
        )

        projection = {
            MongoPeriodicStorage.Index.TIMESTAMP: 1,
            MongoPeriodicStorage.Index.VALUES: 1
        }

        if period is None:
            projection[MongoPeriodicStorage.Index.PERIOD] = 1

        cursor = self._find(document=query, projection=projection)

        cursor.hint(MongoPeriodicStorage.Index.QUERY)

        result = []

        if limit != 0:
            cursor = cursor[:limit]

        for document in cursor:

            timestamp = int(document[MongoPeriodicStorage.Index.TIMESTAMP])

            values = document[MongoPeriodicStorage.Index.VALUES]

            for t in values:
                value = values[t]
                value_timestamp = timestamp + int(t)

                if timewindow is None or value_timestamp in timewindow:
                    result.append((value_timestamp, value))

        result.sort(key=itemgetter(0))

        return result

    def put(self, data_id, period, points, cache=False, *args, **kwargs):

        # initialize a dictionary of perfdata value by value field
        # and id_timestamp
        doc_props_by_id_ts = {}
        # previous variable contains a dict of entries to put in
        # the related document

        # prepare data to insert/update in db
        for ts, value in points:

            ts = int(ts)
            id_timestamp = int(period.round_timestamp(ts, normalize=True))

            document_properties = doc_props_by_id_ts.setdefault(
                id_timestamp, {}
            )

            if '_id' not in document_properties:
                document_properties['_id'] = \
                    MongoPeriodicStorage._get_document_id(
                        data_id=data_id,
                        timestamp=id_timestamp, period=period)
                document_properties[MongoPeriodicStorage.Index.LAST_UPDATE] = \
                    ts

            else:
                last_update = MongoPeriodicStorage.Index.LAST_UPDATE
                if document_properties[last_update] < ts:
                    document_properties[last_update] = ts

            field_name = "{0}.{1}".format(
                MongoPeriodicStorage.Index.VALUES, ts - id_timestamp)

            document_properties[field_name] = value

        for id_timestamp in doc_props_by_id_ts:
            document_properties = doc_props_by_id_ts[id_timestamp]

            # remove _id and last_update
            _id = document_properties.pop('_id')

            _set = {
                MongoPeriodicStorage.Index.DATA_ID: data_id,
                MongoPeriodicStorage.Index.PERIOD: period.unit_values,
                MongoPeriodicStorage.Index.TIMESTAMP: id_timestamp
            }
            _set.update(document_properties)

            document_properties['_id'] = _id

            result = self._update(
                spec={'_id': _id}, document={'$set': _set}, cache=cache
            )

        return result

    def remove(
        self, data_id, period=None, timewindow=None, cache=False,
        *args, **kwargs
    ):

        query = self._get_documents_query(
            data_id=data_id, timewindow=timewindow, period=period
        )

        if timewindow is not None:

            projection = {
                MongoPeriodicStorage.Index.TIMESTAMP: 1,
                MongoPeriodicStorage.Index.VALUES: 1
            }

            documents = self._find(document=query, projection=projection)

            for document in documents:
                timestamp = document.get(MongoPeriodicStorage.Index.TIMESTAMP)
                values = document.get(MongoPeriodicStorage.Index.VALUES)
                values_to_save = {
                    t: values[t] for t in values
                    if (timestamp + int(t)) not in timewindow
                }
                _id = document.get('_id')

                if len(values_to_save) > 0:
                    self._update(
                        spec={'_id': _id},
                        document={
                            '$set': {
                                MongoPeriodicStorage.Index.VALUES:
                                values_to_save}
                        },
                        cache=cache)
                else:
                    self._remove(document=_id, cache=cache)

        else:
            self._remove(document=query, cache=cache)

    def all_indexes(self, *args, **kwargs):

        result = super(MongoPeriodicStorage, self).all_indexes(*args, **kwargs)

        result.append(MongoPeriodicStorage.Index.QUERY)

        return result

    def _get_documents_query(self, data_id, timewindow, period):
        """
        Get mongo documents query about data_id, timewindow and period.

        If period is None and timewindow is not None, period takes default
        period value for data_id.
        """

        result = {
            MongoPeriodicStorage.Index.DATA_ID: data_id
        }

        if period is not None:  # manage specific period
            result[MongoPeriodicStorage.Index.PERIOD] = period.unit_values

        if timewindow is not None:  # manage specific timewindow
            start_timestamp, stop_timestamp = \
                MongoPeriodicStorage._get_id_timestamps(
                    timewindow=timewindow, period=period)
            result[MongoPeriodicStorage.Index.TIMESTAMP] = {
                '$gte': start_timestamp,
                '$lte': stop_timestamp}

        return result

    @staticmethod
    def _get_id_timestamps(timewindow, period):
        """
        Get id timestamps related to input timewindow and period.
        """

        # get minimal timestamp
        start_timestamp = int(
            period.round_timestamp(timewindow.start(), normalize=True))
        # and maximal timestamp
        stop_timestamp = int(
            period.round_timestamp(timewindow.stop(), normalize=True))
        stop_datetime = datetime.fromtimestamp(stop_timestamp)
        delta = period.get_delta()
        stop_datetime += delta
        stop_timestamp = mktime(stop_datetime.timetuple())

        result = start_timestamp, stop_timestamp

        return result

    @staticmethod
    def _get_document_id(data_id, timestamp, period):
        """
        Get periodic document id related to input data_id, timestamp and period
        """

        md5_result = md5()

        # add data_id in id
        md5_result.update(data_id.encode('ascii', 'ignore'))

        # add id_timestamp in id
        md5_result.update(str(timestamp).encode('ascii', 'ignore'))

        # add period in id
        unit_with_value = period.get_max_unit()
        if unit_with_value is None:
            raise MongoPeriodicStorage.Error(
                "period {0} must contain at least one valid unit among {1}".
                format(period, Period.UNITS))

        md5_result.update(
            unit_with_value[Period.UNIT].encode('ascii', 'ignore'))

        # resolve md5
        result = md5_result.hexdigest()

        return result
