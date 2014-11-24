# -*- coding: utf-8 -*-
# --------------------------------
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

from canopsis.storage import Storage


class TimedStorage(Storage):
    """
    Store dedicated to manage timed data.
    It saves one value at one timestamp.
    Two consecutives timestamp values can not be same values.
    """

    __datatype__ = 'timed'

    class Index:

        TIMESTAMP = 0
        VALUE = 1
        DATA_ID = 2

    DATA_ID = 'data_id'
    VALUE = 'value'
    TIMESTAMP = 'timestamp'

    def get(
        self, data_ids, timewindow=None, limit=0, skip=0, sort=None
    ):
        """
        Get a dictionary of sorted list of triplet of dictionaries such as :

        dict(
            tuple(
                timestamp,
                dict(data_type, data_value), data id))

        If timewindow is None, result is all timed document.

        :return:
        :rtype: dict of tuple(float, dict, str)
        """

        raise NotImplementedError()

    def count(self, data_id):
        """
        Get number of timed documents for input data_id.
        """

        raise NotImplementedError()

    def put(self, data_id, value, timestamp):
        """
        Put a dictionary of value by name in collection.
        """

        raise NotImplementedError()

    def remove(self, data_ids, timewindow=None):
        """
        Remove timed_data existing on input timewindow.
        """

        raise NotImplementedError()
