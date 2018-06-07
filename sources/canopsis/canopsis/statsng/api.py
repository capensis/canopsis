# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

from canopsis.statsng.enums import StatRequestFields


class StatsAPIError(Exception):
    """
    A StatsAPIError is an Exception that can be raised by a StatsAPI object.

    It should be handled in `webcore/services/statsng.py`, and returned as a
    JSON object in the response.
    """
    def __init__(self, message):
        super(StatsAPIError, self).__init__(message)
        self.message = message


class StatRequest(object):
    def __init__(self):
        self.stats = None
        self.tstart = None
        self.tstop = None
        self.group_by = []
        self.entity_filter = None

    @staticmethod
    def from_request_body(body):
        """
        Create a StatRequest from the body of a request.

        :param dict body: The parsed body of a request.
        :rtype: StatRequest
        :raises: StatsAPIError
        """
        request = StatRequest()

        request.stats = body.pop(StatRequestFields.stats, None)
        request.tstart = body.pop(StatRequestFields.tstart, None)
        request.tstop = body.pop(StatRequestFields.tstop, None)
        request.group_by = body.pop(StatRequestFields.group_by, [])
        request.entity_filter = body.pop(StatRequestFields.filter, None)

        if body:
            raise StatsAPIError('Unexpected fields: {0}'.format(
                ', '.join(body.keys())))

        return request

class StatsAPI(object):
    """
    A StatsAPI object handles the computation of statistics.
    """
    def __init__(self, logger):
        self.logger = logger

    def handle_request(self, request):
        """
        Handle a request to the statistics API.

        :param StatRequest request:
        :rtype dict:
        :raises: StatsAPIError
        """
        return {}
