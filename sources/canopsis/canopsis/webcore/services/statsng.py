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

from bottle import request

from canopsis.statsng.api import StatRequest, StatsAPI
from canopsis.statsng.enums import StatRequestFields
from canopsis.statsng.errors import StatsAPIError, UnknownStatNameError
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR, \
    HTTP_NOT_FOUND


def parse_request():
    """
    Parse the body of a request.

    :rtype: StatRequest
    :raises: StatsAPIError
    """
    try:
        request_body = request.json
    except ValueError as verror:
        raise StatsAPIError('Malformed JSON: {0}'.format(verror))

    if request_body is None:
        raise StatsAPIError('Empty request')

    return StatRequest.from_request_body(request_body)


def exports(ws):

    api = StatsAPI(ws.logger)

    @ws.application.post(
        '/api/v2/stats'
    )
    def stats():
        try:
            stat_request = parse_request()

            if not stat_request.stats:
                raise StatsAPIError(
                    "The stats field is required and should not be empty.")

            return gen_json(api.handle_request(stat_request))
        except StatsAPIError as error:
            ws.logger.exception(error.message)
            return gen_json_error({'description': error.message}, HTTP_ERROR)

    @ws.application.post(
        '/api/v2/stats/<stat_name>'
    )
    def stat(stat_name):
        try:
            stat_request = parse_request()

            if stat_request.stats is not None:
                raise StatsAPIError(
                    'Unexpected fields: {0}'.format(StatRequestFields.stats))

            stat_request.stats = [stat_name]

            return gen_json(api.handle_request(stat_request))
        except UnknownStatNameError as error:
            ws.logger.exception(error.message)
            return gen_json_error({'description': error.message},
                                  HTTP_NOT_FOUND)
        except StatsAPIError as error:
            ws.logger.exception(error.message)
            return gen_json_error({'description': error.message}, HTTP_ERROR)
