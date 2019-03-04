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

from bottle import request

from pymongo.errors import PyMongoError

from canopsis.common.collection import CollectionError
from canopsis.watcherng import WatcherManager
from canopsis.webcore.utils import (gen_json, gen_json_error,
                                    HTTP_NOT_FOUND, HTTP_ERROR)


def exports(ws):

    watcher_manager = WatcherManager(WatcherManager.default_collection(), ws.logger)

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

        try:
            return watcher_manager.create_watcher(watcher)
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
