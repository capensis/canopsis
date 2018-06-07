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

from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR
from canopsis.statsng.enums import StatRequestFields


def compute_stats(logger, request_body):
    """
    Compute the response to a request to the stats API.
    """
    tstart = request_body.pop(StatRequestFields.tstart, None)
    tstop = request_body.pop(StatRequestFields.tstop, None)
    group_by = request_body.pop(StatRequestFields.group_by, [])
    limit = request_body.pop(StatRequestFields.limit, None)
    offset = request_body.pop(StatRequestFields.offset, None)
    entity_filter = request_body.pop(StatRequestFields.filter, None)

    try:
        stats = request_body.pop(StatRequestFields.stats)
    except KeyError:
        message = 'missing field: {0}'.format(StatRequestFields.stats)
        logger.exception(message)
        return gen_json_error({'description': message}, HTTP_ERROR)

    if request_body:
        message = 'unexpected fields: {0}'.format(
            ', '.join(request_body.keys()))
        logger.error(message)
        return gen_json_error({'description': message}, HTTP_ERROR)

    return gen_json({})


def exports(ws):

    @ws.application.post(
        '/api/v2/stats'
    )
    def stats():
        try:
            request_body = request.json
        except ValueError as verror:
            ws.logger.exception('malformed JSON: {0}'.format(verror))
            return gen_json_error(
                {'description': 'malformed JSON: {0}'.format(verror)},
                HTTP_ERROR)

        if request_body is None:
            ws.logger.exception('empty request')
            return gen_json_error(
                {'description': 'empty request'},
                HTTP_ERROR)

        return compute_stats(ws.logger, request_body)

    @ws.application.post(
        '/api/v2/stats/<stat_name>'
    )
    def stat(stat_name):
        try:
            request_body = request.json
        except ValueError as verror:
            ws.logger.exception('malformed JSON: {0}'.format(verror))
            return gen_json_error(
                {'description': 'malformed JSON: {0}'.format(verror)},
                HTTP_ERROR)

        if request_body is None:
            ws.logger.exception('empty request')
            return gen_json_error(
                {'description': 'empty request'},
                HTTP_ERROR)

        # TODO: send 404 for unknown stat_name

        if StatRequestFields.stats in request_body:
            message = 'unexpected fields: {0}'.format(StatRequestFields.stats)
            ws.logger.exception(message)
            return gen_json_error({'description': message}, HTTP_ERROR)

        request_body[StatRequestFields.stats] = [stat_name]

        return compute_stats(ws.logger, request_body)
