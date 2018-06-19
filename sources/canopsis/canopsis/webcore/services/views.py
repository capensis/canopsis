#!/usr/bin/env python
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

from canopsis.views.enums import ViewField
from canopsis.views.adapter import ViewAdapter
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR, \
    HTTP_NOT_FOUND


def exports(ws):

    adapter = ViewAdapter()

    @ws.application.get(
        '/api/v2/views/<view_id>'
    )
    def get_view(view_id):
        view = adapter.get_by_id(view_id)

        if view:
            return gen_json(view)
        else:
            return gen_json_error({
                'description': 'No view with id: {0}'.format(view_id)
            }, HTTP_NOT_FOUND)


    @ws.application.post(
        '/api/v2/views'
    )
    def create_view():
        try:
            request_body = request.json
        except ValueError as verror:
            return gen_json_error({
                'description': 'Malformed JSON: {0}'.format(verror)
            }, HTTP_ERROR)

        view_id = adapter.create(request_body)
        return gen_json({'id': view_id})


    @ws.application.delete(
        '/api/v2/views/<view_id>'
    )
    def remove_view(view_id):
        # TODO: should there be an error if the view does not exist?
        adapter.remove_with_id(view_id)
        return gen_json({})


    @ws.application.put(
        '/api/v2/views/<view_id>'
    )
    def update_view(view_id):
        try:
            request_body = request.json
        except ValueError as verror:
            return gen_json_error({
                'description': 'Malformed JSON: {0}'.format(verror)
            }, HTTP_ERROR)

        # TODO: should there be an error if the view does not exist?
        adapter.update(view_id, request_body)
        return gen_json({})
