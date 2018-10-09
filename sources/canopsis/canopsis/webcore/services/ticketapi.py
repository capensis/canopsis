# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from canopsis.common.collection import CollectionError
from canopsis.common.converters import id_filter
from canopsis.models.ticketapi import TicketApi
from canopsis.ticketapi.manager import TicketApiManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR


def exports(ws):

    ws.application.router.add_filter('id_filter', id_filter)

    ticketapi_manager = TicketApiManager(*TicketApiManager.provide_default_basics())

    @ws.application.get(
        '/api/v2/ticketapi/<ticketapi_id:id_filter>'
    )
    def get_ticketapi(ticketapi_id):
        """
        Get an existing ticketapi.

        :param ticketapi_id: ID of the ticketApiConfig
        :type ticketapi_id: str
        :returns: <TicketApi>
        :rtype: dict
        """
        ticketapi = ticketapi_manager.get_id(id_=ticketapi_id)
        if not isinstance(ticketapi, TicketApi):
            return gen_json_error({'description': 'failed to get ticketapi'},
                                  HTTP_ERROR)

        return gen_json(ticketapi.to_dict())

    @ws.application.post(
        '/api/v2/ticketapi'
    )
    def create_ticketapi():
        """
        Create a new ticketapi.

        :returns: nothing
        """
        try:
            element = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        if element is None or not isinstance(element, dict):
            return gen_json_error(
                {'description': 'nothing to insert'}, HTTP_ERROR)
        try:
            TicketApi(**TicketApi.convert_keys(element))
        except TypeError:
            return gen_json_error(
                {'description': 'invalid ticketapi format'}, HTTP_ERROR)

        try:
            ok = ticketapi_manager.create(ticketapi=element)
        except CollectionError as ce:
            ws.logger.error('TicketApi creation error : {}'.format(ce))
            return gen_json_error(
                {'description': 'error while creating an ticketapi'},
                HTTP_ERROR
            )

        if not ok:
            return gen_json_error({'description': 'failed to create ticketapi'},
                                  HTTP_ERROR)

        return gen_json({})

    @ws.application.put(
        '/api/v2/ticketapi/<ticketapi_id:id_filter>'
    )
    def update_ticketapi(ticketapi_id):
        """
        Update an existing ticketapi.

        :param ticketapi_id: ID of the ticketapi
        :type ticketapi_id: str
        :returns: nothing
        """
        try:
            element = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        if element is None or not isinstance(element, dict) or len(element) <= 0:
            return gen_json_error(
                {'description': 'wrong update dict'}, HTTP_ERROR)
        try:
            TicketApi(**TicketApi.convert_keys(element))
        except TypeError:
            return gen_json_error(
                {'description': 'invalid ticketapi format'}, HTTP_ERROR)

        try:
            ok = ticketapi_manager.update_id(id_=ticketapi_id, ticketapi=element)
        except CollectionError as ce:
            ws.logger.error('TicketApi update error : {}'.format(ce))
            return gen_json_error(
                {'description': 'error while updating an ticketapi'},
                HTTP_ERROR
            )
        if not ok:
            return gen_json_error({'description': 'failed to update ticketapi'},
                                  HTTP_ERROR)

        return gen_json({})

    @ws.application.delete(
        '/api/v2/ticketapi/<ticketapi_id:id_filter>'
    )
    def delete_id(ticketapi_id):
        """
        Delete a ticketApiConfig, based on his id.

        :param ticketapi_id: ID of the ticketApiConfig
        :type ticketapi_id: str

        :rtype: bool
        """
        ws.logger.info('Delete ticketapi : {}'.format(ticketapi_id))

        return gen_json(ticketapi_manager.delete_id(id_=ticketapi_id))
