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


class TimedStorage(Storage):
    """Storage dedicated to manage timed data."""

    __datatype__ = 'timed'

    TIMESTAMP = 'timestamp'
    VALUES = 'values'
    PERIOD = 'period'
    LAST_UPDATE = 'last_update'

    def count(self, data_id, timewindow=None, tags=None):
        """Get number of timed documents for input data_id."""

        raise NotImplementedError()

    def size(self, data_id=None, timewindow=None, tags=None, *args, **kwargs):
        """Get size occupied by research filter data_id."""

        raise NotImplementedError()

    def get(
            self, data_id, timewindow=None, limit=0, skip=0, sort=None,
            tags=None
    ):
        """Get a list of points."""

        raise NotImplementedError()

    def put(self, data_id, points, tags=None, cache=False):
        """Put timed points in a timed collection with specific timed values.

        :param points: iterable of (timestamp, value).
        :param list tags: indexed tags of points.
        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()

    def remove(self, data_id, timewindow=None, cache=False, tags=None):
        """Remove timed data related to data_id and timewindow.

        :param canopsis.timeserie.timewindow.TimeWindow timewindow: Default
            remove all timed data with input period.
        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()
