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

from pymongo.errors import PyMongoError

from canopsis.event import forger

from canopsis.common.collection import CollectionError
from canopsis.watcherng import WatcherManager
from canopsis.webcore.utils import (gen_json, gen_json_error,
                                    HTTP_NOT_FOUND, HTTP_ERROR)
from canopsis.common.amqp import AmqpPublisher
from canopsis.common.amqp import get_default_connection as \
    get_default_amqp_conn


def exports(ws):

    watcher_manager = WatcherManager(WatcherManager.default_collection())
    amqp_pub = AmqpPublisher(get_default_amqp_conn(), ws.logger)

    @ws.application.get(
        '/api/v2/watcherng'
    )
    def get_watcher_list():
        """
        Return the list of all watchers.

        :returns: <Watcherng>
        :rtype: list
        """
        try:
            document = watcher_manager.get_watcher_list()
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the watchers list from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json(document)

    @ws.application.get(
        '/api/v2/watcherng/<watcher_id>'
    )
    def get_watcher_by_id(watcher_id):
        """
        Return a watcher given the id.

        :param watcher_id: ID of the watcher
        :type watcher_id: str
        :returns: <Watcherng>
        :rtype: dict
        """
        try:
            document = watcher_manager.get_watcher_by_id(watcher_id)
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the watcher data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        if document is None:
            return gen_json_error(
                {"description": "No watcher found with ID " + watcher_id},
                HTTP_ERROR)

        return gen_json(document)

    @ws.application.post(
        '/api/v2/watcherng'
    )
    def create_watcher():
        """
        Create a new watcher.

        :returns: ID of the watcher
        :rtype: string
        """
        try:
            watcher = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'Invalid JSON'},
                HTTP_ERROR
            )

        if watcher is None or not isinstance(watcher, dict):
            return gen_json_error(
                {'description': 'Nothing to create'}, HTTP_ERROR)

        if watcher['type'] != 'watcher':
            return gen_json_error(
                {'description': 'Entity is not a watcher'}, HTTP_ERROR)

        if 'entities' not in watcher or 'state' not in watcher or 'output_template' not in watcher:
            return gen_json_error(
                {'description': 'Watcher is missing important specific fields'}, HTTP_ERROR)

        if '_id' not in watcher:
            watcher['_id'] = str(uuid.uuid4())

        try:
            wid = watcher_manager.create_watcher(watcher)

            event = forger(
                connector="watcher",
                connector_name="watcher",
                event_type="updatewatcher",
                source_type="component",
                component=wid)
            amqp_pub.canopsis_event(event)

            return wid
        except CollectionError as ce:
            ws.logger.error('Watcherng creation error : {}'.format(ce))
            return gen_json_error(
                {'description': 'Error while creating a watcher'},
                HTTP_ERROR
            )

    @ws.application.put(
        '/api/v2/watcherng/<watcher_id>'
    )
    def update_watcher_by_id(watcher_id):
        """
        Update an existing watcher.

        :param watcher_id: ID of the watcher
        :type watcher_id: str
        :rtype: dict
        """
        try:
            watcher = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'Invalid JSON'},
                HTTP_ERROR
            )

        if watcher is None or not isinstance(watcher, dict):
            return gen_json_error(
                {'description': 'Nothing to update'}, HTTP_ERROR)

        try:
            ok = watcher_manager.update_watcher_by_id(watcher, watcher_id)
        except CollectionError as ce:
            ws.logger.error('Watcherng update error : {}'.format(ce))
            return gen_json_error(
                {'description': 'Error while updating a watcher'},
                HTTP_ERROR
            )

        if not ok:
            return gen_json_error(
                {'description': 'Failed to update watcher'},
                HTTP_ERROR
            )

        event = forger(
            connector="watcher",
            connector_name="watcher",
            event_type="updatewatcher",
            source_type="component",
            component=watcher_id)
        amqp_pub.canopsis_event(event)

        return gen_json({})

    @ws.application.delete(
        '/api/v2/watcherng/<watcher_id>'
    )
    def delete_watcher_by_id(watcher_id):
        """
        Delete an existing watcher, given its id.

        :param watcher_id: ID of the watcher
        :type watcher_id: str
        :rtype: dict
        """
        try:
            ok = watcher_manager.delete_watcher_by_id(watcher_id)
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the watcher data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json({"status": ok})
