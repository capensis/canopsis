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


class PeriodicStorage(Storage):
    """
    Storage dedicated to manage periodic data.
    """

    __datatype__ = 'periodic'

    TIMESTAMP = 'timestamp'
    VALUES = 'values'
    PERIOD = 'period'
    LAST_UPDATE = 'last_update'

    def count(self, data_id, period, timewindow=None):
        """
        Get number of periodic documents for input data_id.
        """

        raise NotImplementedError()

    def size(
        self, data_id=None, period=None, timewindow=None,
        *args, **kwargs
    ):
        """
        Get size occupied by research filter data_id
        """

        raise NotImplementedError()

    def get(
        self, data_id, period, timewindow=None, limit=0, skip=0, sort=None
    ):
        """
        Get a list of points.
        """

        raise NotImplementedError()

    def put(self, data_id, period, points, cache=False):
        """
        Put periodic points in periodic collection with specific period values.

        points is an iterable of (timestamp, value)

        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()

    def remove(self, data_id, period=None, timewindow=None, cache=False):
        """
        Remove periodic data related to data_id, timewindow and period.
        If timewindow is None, remove all periodic_data with input period.
        If period is None, remove all data whatever related period.

        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()
