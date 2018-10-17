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

from canopsis.views.enums import ViewField, GroupField
from canopsis.views.adapters import ViewAdapter, GroupAdapter, \
    InvalidViewError, InvalidGroupError, NonEmptyGroupError, \
    InvalidFilterError
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR, \
    HTTP_NOT_FOUND


def exports(ws):

    view_adapter = ViewAdapter(ws.logger)
    group_adapter = GroupAdapter(ws.logger)

    @ws.application.get(
        '/api/v2/views'
    )
    def list_views():
        name = request.query.getunicode(ViewField.name)
        title = request.query.getunicode(ViewField.title)

        try:
            return gen_json(view_adapter.list(name, title))
        except InvalidFilterError as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

    @ws.application.get(
        '/api/v2/views/<view_id>'
    )
    def get_view(view_id):
        view = view_adapter.get_by_id(view_id)

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

        if not request_body:
            return gen_json_error({
                'description': 'Empty request'
            }, HTTP_ERROR)

        try:
            view_id = view_adapter.create(request_body)
        except InvalidViewError as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        return gen_json({
            ViewField.id: view_id
        })

    @ws.application.delete(
        '/api/v2/views/<view_id>'
    )
    def remove_view(view_id):
        if not view_adapter.get_by_id(view_id):
            return gen_json_error({
                'description': 'No view with id: {0}'.format(view_id)
            }, HTTP_NOT_FOUND)

        view_adapter.remove_with_id(view_id)

        return gen_json({})

    @ws.application.put(
        '/api/v2/views/<view_id>'
    )
    def update_view(view_id):
        if not view_adapter.get_by_id(view_id):
            return gen_json_error({
                'description': 'No view with id: {0}'.format(view_id)
            }, HTTP_NOT_FOUND)

        try:
            request_body = request.json
        except ValueError as verror:
            return gen_json_error({
                'description': 'Malformed JSON: {0}'.format(verror)
            }, HTTP_ERROR)

        if not request_body:
            return gen_json_error({
                'description': 'Empty request'
            }, HTTP_ERROR)

        try:
            view_adapter.update(view_id, request_body)
        except InvalidViewError as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        return gen_json({})

    @ws.application.get(
        '/api/v2/views/groups'
    )
    def list_groups():
        name = request.query.getunicode(GroupField.name)
        return gen_json(group_adapter.list(name))

    @ws.application.get(
        '/api/v2/views/groups/<group_id>'
    )
    def get_group(group_id):
        name = request.query.getunicode(ViewField.name)
        title = request.query.getunicode(ViewField.title)

        try:
            group = group_adapter.get_by_id(group_id, name, title)
        except InvalidFilterError as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        if not group:
            return gen_json_error({
                'description': 'No group with id: {0}'.format(group_id)
            }, HTTP_NOT_FOUND)

        return gen_json(group)

    @ws.application.post(
        '/api/v2/views/groups'
    )
    def create_group():
        try:
            request_body = request.json
        except ValueError as verror:
            return gen_json_error({
                'description': 'Malformed JSON: {0}'.format(verror)
            }, HTTP_ERROR)

        if not request_body:
            return gen_json_error({
                'description': 'Empty request'
            }, HTTP_ERROR)

        try:
            group_id = group_adapter.create(request_body)
        except InvalidGroupError as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        return gen_json({
            GroupField.id: group_id
        })

    @ws.application.delete(
        '/api/v2/views/groups/<group_id>'
    )
    def remove_group(group_id):
        if not group_adapter.exists(group_id):
            return gen_json_error({
                'description': 'No group with id: {0}'.format(group_id)
            }, HTTP_NOT_FOUND)

        try:
            group_adapter.remove_with_id(group_id)
        except NonEmptyGroupError:
            return gen_json_error({
                'description': 'The group is not empty'
            }, HTTP_ERROR)

        return gen_json({})

    @ws.application.put(
        '/api/v2/views/groups/<group_id>'
    )
    def update_group(group_id):
        if not group_adapter.exists(group_id):
            return gen_json_error({
                'description': 'No group with id: {0}'.format(group_id)
            }, HTTP_NOT_FOUND)

        try:
            request_body = request.json
        except ValueError as verror:
            return gen_json_error({
                'description': 'Malformed JSON: {0}'.format(verror)
            }, HTTP_ERROR)

        if not request_body:
            return gen_json_error({
                'description': 'Empty request'
            }, HTTP_ERROR)

        try:
            group_adapter.update(group_id, request_body)
        except InvalidGroupError as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        return gen_json({})
