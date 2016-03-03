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

from canopsis.storage.core import Storage


class PeriodicalStorage(Storage):
    """Store dedicated to manage periodical data.

    A periodical data exist over a period.

    It saves one value at one timestamp.
    Two consecutives timestamp values can not be same values.
    """

    __datatype__ = 'periodical'

    class Index:

        TIMESTAMP = 0
        VALUE = 1
        DATA_ID = 2

    DATA_ID = 'data_id'
    VALUE = 'value'
    TIMESTAMP = 'timestamp'

    def get(
            self, data_ids, timewindow=None, _filter=None,
            limit=0, skip=0, sort=None
    ):
        """Get data values ordered by timestamp in descresent order and data ids
        if data_ids is a list of data ids.

        If timewindow is None, result is all timed document.

        :param str(s) data_ids: data id(s).
        :param TimeWindow timedwindow: data value time window. Default is
            current timestamp offsset.
        :param _filter: value filter.
        :param int limit: maximal number of data to retrieve.
        :param int skip: number of documents to skip.
        :param dict sort: sort documents by something else than timestamp.

        :return: depends on type of data_ids

            - string: list of {
                    TimedStorage.TIMESTAMP: timestamp,
                    TimedStorage.VALUE: value
                }
            - list: dictionary of {
                    data_id: self.get(data_id, ...)
                }

        :rtype: dict or list
        """

        raise NotImplementedError()

    def find(self, timewindow=None, _filter=None):
        """Find data values ordered by timestamp in descresent order and data
        ids.

        If timewindow is None, result is all timed document.

        :param TimeWindow timedwindow: data value time window. Default is
            current timestamp offsset.
        :param _filter: value filter.

        :return: dictionary of {
                data_id: list of {
                    TimedStorage.TIMESTAMP: timestamp,
                    TimedStorage.VALUE: value
                }
            }
        :rtype: dict
        """

        raise NotImplementedError()

    def count(self, data_ids, timewindow=None, _filter=None):
        """Get number of period of timed documents for input data_id.

        :param str(s) data_ids: data_id(s) to count.
        :param TimeWindow timewindow: count timewindow.
        :param dict _filter: filter applied on values.
        """

        raise NotImplementedError()

    def put(self, data_id, value, timestamp, cache=False):
        """Put a dictionary of value by name in collection.

        :param str data_id: related data_id.
        :param value: value to associate to data_id at timestamp time.
        :param float timestamp: timestamp to use for associating value to
            data_id. If timestamp already exists, update existing value with
            the same timestamp. Otherwise, check if value is different than the
            previous value and add a document if False.
        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()

    def remove(self, data_ids, timewindow=None, cache=False):
        """Remove periodical data existing on input timewindow.

        :param str(s) data_ids: related data_id(s) to remove.
        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()
