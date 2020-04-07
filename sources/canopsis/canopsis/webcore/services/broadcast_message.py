# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2019 "Capensis" [http://www.capensis.fr]
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

import uuid
from bottle import request
from six import string_types
from pymongo.errors import PyMongoError

from canopsis.common.collection import CollectionError
from canopsis.broadcast_message import BroadcastMessageManager
from canopsis.webcore.utils import (gen_json, gen_json_error,
                                    HTTP_NOT_FOUND, HTTP_ERROR)
from dateutil.parser import parse

FIELDS = {"message", "color", "start", "end"}


def sanitize_payload(payload):
    for k in payload:
        if k not in FIELDS:
            payload.pop(k)
    if not set(payload.keys()) == FIELDS:
        raise Exception("payload contains not enough fields")
    if not isinstance(payload["message"], string_types) or \
        not isinstance(payload["color"], string_types) or \
        not isinstance(payload["start"], (int, float)) or \
        not isinstance(payload["end"], (int, float)):

        raise Exception("invalid field type")
    return payload


def exports(ws):
    message_manager = BroadcastMessageManager(BroadcastMessageManager.default_collection())

    @ws.application.get(
        '/api/v2/broadcast-message'
    )
    def get_message_list():
        """
        Return the list of all messages.

        :returns: <message>
        :rtype: list
        """
        try:
            document = message_manager.get_message_list()
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the messages list from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json(document)

    @ws.application.get(
        '/api/v2/broadcast-message/<message_id>'
    )
    def get_message_by_id(message_id):
        """
        Return a message given the id.

        :param message_id: ID of the message
        :type message_id: str
        :returns: <message>
        :rtype: dict
        """
        try:
            document = message_manager.get_message_by_id(message_id)
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the message data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        if document is None:
            return gen_json_error(
                {"description": "No message found with ID " + message_id},
                HTTP_ERROR)

        return gen_json(document)

    @ws.application.post(
        '/api/v2/broadcast-message'
    )
    def create_message():
        """
        Create a new message.

        :returns: ID of the message
        :rtype: string
        """
        try:
            message = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'Invalid JSON'},
                HTTP_ERROR
            )

        if message is None or not isinstance(message, dict):
            return gen_json_error(
                {'description': 'Nothing to create'}, HTTP_ERROR)

        try:
            message = sanitize_payload(message)
        except Exception as e:
            ws.logger.error('message creation error : {}'.format(e))
            return gen_json_error(
                {'description': 'Invalid payload'},
                HTTP_ERROR
            )

        if '_id' not in message:
            message['_id'] = str(uuid.uuid4())

        try:
            _id = message_manager.create_message(message)
            return {'_id': _id}
        except CollectionError as ce:
            ws.logger.error('message creation error : {}'.format(ce))
            return gen_json_error(
                {'description': 'Error while creating an message'},
                HTTP_ERROR
            )

    @ws.application.put(
        '/api/v2/broadcast-message/<message_id>'
    )
    def update_message_by_id(message_id):
        """
        Update an existing message.

        :param message_id: ID of the message
        :type message_id: str
        :rtype: dict
        """
        try:
            message = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'Invalid JSON'},
                HTTP_ERROR
            )

        if message is None or not isinstance(message, dict):
            return gen_json_error(
                {'description': 'Nothing to update'}, HTTP_ERROR)

        try:
            message = sanitize_payload(message)
        except Exception as e:
            ws.logger.error('message creation error : {}'.format(e))
            return gen_json_error(
                {'description': 'Invalid payload'},
                HTTP_ERROR
            )

        try:
            ok = message_manager.update_message_by_id(message, message_id)
        except CollectionError as ce:
            ws.logger.error('message update error : {}'.format(ce))
            return gen_json_error(
                {'description': 'Error while updating an message'},
                HTTP_ERROR
            )

        if not ok:
            return gen_json_error(
                {'description': 'Failed to update message'},
                HTTP_ERROR
            )

        return gen_json({})

    @ws.application.delete(
        '/api/v2/broadcast-message/<message_id>'
    )
    def delete_message_by_id(message_id):
        """
        Delete an existing message, given its id.

        :param message_id: ID of the message
        :type message_id: str
        :rtype: dict
        """
        try:
            ok = message_manager.delete_message_by_id(message_id)
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the message data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json({"status": ok})

    @ws.application.get(
        '/api/v2/broadcast-message/active'
    )
    def get_active():
        """
        Delete an existing message, given its id.

        :param message_id: ID of the message
        :type message_id: str
        :rtype: dict
        """
        try:
            active_msg = message_manager.get_active()
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the message data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json(active_msg)
