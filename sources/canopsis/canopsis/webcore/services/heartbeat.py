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
from canopsis.heartbeat.manager import HeartBeatService
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR
from canopsis.models.heartbeat import HeartBeat


from canopsis.heartbeat.heartbeat import (HeartbeatManager,
                                          HeartbeatPatternExistsError)
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import CollectionError

hb_service = HeartBeatService(*HeartBeatService.provide_default_basics())


def get_heartbeat_collection():
    store = MongoStore.get_default()
    return store.get_collection(name=HeartbeatManager.COLLECTION)


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
        :rtype: a dict with the status (name) of the request and if needed a
        description.
        """

        try:
            json = request.json
        except ValueError:
            return gen_json_error({'description': "invalid json."},
                                  HTTP_ERROR)

        if not HeartBeat.is_valid_heartbeat(json):
            return gen_json_error(
                {"description": "invalid heartbeat payload."}, HTTP_ERROR)

        try:
            collection = get_heartbeat_collection()
        except PyMongoError:
            return gen_database_error()

        manager = HeartbeatManager(collection)
        try:
            heartbeat_id = manager.insert_heartbeat_document(
                json[HeartBeat.PATTERN_KEY],
                json[HeartBeat.EXPECTED_INTERVAL_KEY])

        except HeartbeatPatternExistsError:
            return gen_json_error(
                {"description": "heartbeat pattern already exists"},
                HTTP_ERROR)

        except (PyMongoError, CollectionError):
            return gen_database_error()

        heartbeat = manager.find_heartbeat_document(heartbeat_id)
        return gen_json(heartbeat)

    @ws.application.get(
        "/api/v2/heartbeat/"
    )
    def get_heartbeats():
        """ Return every heartbeats stored in database.
        :rtype: a json representation as a dict of every heartbeats stored in
        or a dict with the status (name) and the description of the issue
        encountered.
        """
        try:
            collection = get_heartbeat_collection()
            manager = HeartbeatManager(collection)
            return gen_json([x for x in manager.list_heartbeat_collection()])
        except PyMongoError:
            return gen_database_error()
