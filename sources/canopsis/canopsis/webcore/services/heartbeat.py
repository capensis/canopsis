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

from pymongo.errors import PyMongoError
from bottle import request

from canopsis.webcore.utils import (gen_json, gen_json_error,
                                    HTTP_ERROR, HTTP_NOT_FOUND)
from canopsis.models.heartbeat import HeartBeat
from canopsis.heartbeat.heartbeat import (HeartbeatManager,
                                          HeartbeatPatternExistsError)
from canopsis.common.collection import CollectionError


def gen_database_error():
    return gen_json_error(
                {"description": "can not retrieve the canopsis version from "
                                "database, contact your administrator."},
                HTTP_ERROR)


def exports(ws):

    @ws.application.post(
        "/api/v2/heartbeat/"
    )
    def create_heartbeat():
        """Create a new heartbeat. Read the body of the request to extract
        the heartbeat as a json.

        :rtype: a dict with the status (name) of the request
        and ID of created Heartbeat as description.
        """
        try:
            json = request.json
        except ValueError:
            return gen_json_error({'description': "invalid json."},
                                  HTTP_ERROR)
        try:
            model = HeartBeat(json)
        except ValueError:
            return gen_json_error(
                {"description": "invalid heartbeat payload."}, HTTP_ERROR)

        try:
            manager = HeartbeatManager(
                *HeartbeatManager.provide_default_basics())

            heartbeat_id = manager.create_heartbeat(model)

        except HeartbeatPatternExistsError:
            return gen_json_error(
                {"description": "heartbeat pattern already exists"},
                HTTP_ERROR)

        except (PyMongoError, CollectionError):
            return gen_database_error()

        return gen_json({
            "name": "heartbeat created",
            "description": heartbeat_id
        })

    @ws.application.get(
        "/api/v2/heartbeat/"
    )
    def list_heartbeats():
        """ Return every heartbeats stored in database.

        :rtype: a json representation as a list of every heartbeats stored in
        or a dict with the status (name) and the description of the issue
        encountered.
        """
        try:
            manager = HeartbeatManager(
                *HeartbeatManager.provide_default_basics())
            return gen_json([x for x in manager.list_heartbeat_collection()])
        except PyMongoError:
            return gen_database_error()

    @ws.application.delete(
        '/api/v2/heartbeat/<heartbeat_id:id_filter>'
    )
    def remove_heartbeat(heartbeat_id):
        """
        Remove a Heartbeat by ID.

        :param `str` heartbeat_id: Heartbeat ID.
        :returns: ``200 OK`` if success or ``404 Not Found`` if a Heartbeat
                  with a given ID is not found or ``400 Bad Request``
                  if database error.
        """
        try:
            manager = HeartbeatManager(
                *HeartbeatManager.provide_default_basics())

            if not manager.find_heartbeat_document(heartbeat_id):
                return gen_json_error({"name": "heartbeat not found",
                                       "description": heartbeat_id},
                                      HTTP_NOT_FOUND)

            manager.remove_heartbeat_document(heartbeat_id)

        except (PyMongoError, CollectionError):
            return gen_database_error()

        return gen_json({
            "name": "heartbeat removed",
            "description": heartbeat_id
        })

    @ws.application.get(
        '/api/v2/heartbeat/<heartbeat_id:id_filter>'
    )
    def get_heartbeat(heartbeat_id):
        """
        Get a Heartbeat by ID.

        :param `str` heartbeat_id: Heartbeat ID.
        :returns: ``200 OK`` and a Heartbeat document as response body
                  if success or ``404 Not Found`` if a Heartbeat
                  with a given ID is not found or ``400 Bad Request``
                  if database error.
        """
        try:
            manager = HeartbeatManager(
                *HeartbeatManager.provide_default_basics())

            heartbeat = manager.find_heartbeat_document(heartbeat_id)
            if not heartbeat:
                return gen_json_error({"name": "heartbeat not found",
                                       "description": heartbeat_id},
                                      HTTP_NOT_FOUND)

        except (PyMongoError, CollectionError):
            return gen_database_error()

        return gen_json(heartbeat)
