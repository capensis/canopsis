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
from bottle import request, install
from six import string_types
from pymongo.errors import PyMongoError

from canopsis.common.collection import CollectionError
from canopsis.playlist import ViewPlaylistManager
from canopsis.webcore.utils import (gen_json, gen_json_error,
                                    HTTP_NOT_FOUND, HTTP_ERROR)
from canopsis.webcore.services.internal import sanitize_popup_timeout


FIELDS = {"name", "interval", "fullscreen", "enabled", "tabs_list"}


def sanitize_payload(payload):
    for k in payload:
        if k not in FIELDS:
            payload.pop(k)
    if not set(payload.keys()) == FIELDS:
        raise Exception("payload contains not enough fields")
    if not isinstance(payload["name"], string_types) or \
            not isinstance(payload["fullscreen"], bool) or \
            not isinstance(payload['enabled'], bool) or \
            not isinstance(payload["tabs_list"], list):

        raise Exception("invalid field type")

    sanitize_popup_timeout(payload['interval'])
    return payload


def exports(ws):
    playlist_manager = ViewPlaylistManager(
        ViewPlaylistManager.default_collection())

    @ws.application.get(
        '/api/v2/playlist'
    )
    def get_playlist_list():
        """
        Return the list of all messages.

        :returns: <message>
        :rtype: list
        """
        try:
            document = playlist_manager.get_playlist_list()
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the playlists list from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json(document)

    @ws.application.get(
        '/api/v2/playlist/<playlist_id>'
    )
    def get_message_by_id(playlist_id):
        """
        Return a message given the id.

        :param playlist_id: ID of the message
        :type playlist_id: str
        :returns: <message>
        :rtype: dict
        """
        try:
            document = playlist_manager.get_playlist_by_id(playlist_id)
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the playlist data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        if document is None:
            return gen_json_error(
                {"description": "No message found with ID " + playlist_id},
                HTTP_ERROR)

        return gen_json(document)

    @ws.application.post(
        '/api/v2/playlist'
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
            _id = playlist_manager.create_playlist(message)
            return {'_id': _id}
        except CollectionError as ce:
            ws.logger.error('message creation error : {}'.format(ce))
            return gen_json_error(
                {'description': 'Error while creating an message'},
                HTTP_ERROR
            )

    @ws.application.put(
        '/api/v2/playlist/<playlist_id>'
    )
    def update_message_by_id(playlist_id):
        """
        Update an existing message.

        :param playlist_id: ID of the message
        :type playlist_id: str
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
            ok = playlist_manager.update_playlist_by_id(message, playlist_id)
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
        '/api/v2/playlist/<playlist_id>'
    )
    def delete_message_by_id(playlist_id):
        """
        Delete an existing message, given its id.

        :param playlist_id: ID of the message
        :type playlist_id: str
        :rtype: dict
        """
        try:
            ok = playlist_manager.delete_playlist_by_id(playlist_id)
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the playlist data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json({"status": ok})
