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

from canopsis.middleware.core import Middleware
from canopsis.common.ws import route
from bottle import HTTPError


def exports(ws):
    @route(ws.application.get, name='storage')
    def get_elements(protocol, data_type, data_scope, _id=None, **kwargs):
        storage = Middleware.get_middleware(
            protocol, data_type, data_scope,
            **kwargs
        )
        storage.connect()

        return storage.get_elements(ids=_id) or []

    @route(
        ws.application.post,
        name='storage',
        payload=[
            'query', 'projection',
            'limit', 'skip', 'sort',
            'with_count'
        ]
    )
    def find_elements(
        protocol, data_type, data_scope,
        query=None, projection=None,
        limit=0, skip=0, sort=None,
        with_count=False,
        **kwargs
    ):
        storage = Middleware.get_middleware(
            protocol, data_type, data_scope,
            **kwargs
        )
        storage.connect()

        total = storage.count_elements(query=query)
        result = storage.find_elements(
            query=query, projection=projection,
            limit=limit, skip=skip, sort=sort,
            with_count=with_count
        )

        return result, total

    @route(ws.application.put, name='storage', payload=['element'])
    def put_element(
        protocol, data_type, data_scope,
        _id=None, element=None,
        **kwargs
    ):
        storage = Middleware.get_middleware(
            protocol, data_type, data_scope,
            **kwargs
        )
        storage.connect()

        if not storage.put_element(element, _id=_id):
            return HTTPError(500, 'Impossible to put element in storage')

    @route(ws.application.delete, name='storage', payload=['_filter'])
    def remove_elements(
        protocol, data_type, data_scope,
        _id=None, _filter=None,
        **kwargs
    ):
        storage = Middleware.get_middleware(
            protocol, data_type, data_scope,
            **kwargs
        )
        storage.connect()

        storage.remove_elements(ids=_id, _filter=_filter)
