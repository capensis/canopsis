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

from pymongo.errors import PyMongoError

from canopsis.action.manager import ActionManager
from canopsis.common.collection import CollectionError
from canopsis.common.converters import id_filter
from canopsis.models.action import Action
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR, HTTP_UNAUTHORIZED
import time


def exports(ws):
    ws.application.router.add_filter('id_filter', id_filter)

    action_manager = ActionManager(*ActionManager.provide_default_basics())

    def _sanitize_pbehavior_parameter(action):
        parameters = action.get('parameters', {})
        if 'tstart' in parameters and 'tstop' in parameters and \
            isinstance(parameters['tstart'], (int, float)) and \
            isinstance(parameters['tstop'], (int, float)):
            return

        if 'start_on_trigger' in parameters and 'duration' in parameters and \
            isinstance(parameters['duration'], dict) and \
            'seconds' in parameters['duration'] and 'unit' in parameters['duration'] and \
            isinstance(parameters['duration']['seconds'], (int, float)) and \
            isinstance(parameters['duration']['unit'], basestring) and \
            isinstance(parameters['start_on_trigger'], bool):
            return

        raise ValueError("invalid pbehavior parameters")

    @ws.application.get(
        '/api/v2/actions'
    )
    def get_action_list():
        """
        Return the list of all actions.

        :returns: <Action>
        :rtype: list
        """
        try:
            actions = action_manager.get_action_list()
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the actions list from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json(actions)

    @ws.application.get(
        '/api/v2/actions/<action_id:id_filter>'
    )
    def get_action(action_id):
        """
        Get an existing action.

        :param action_id: ID of the alarm-action
        :type action_id: str
        :returns: <Action>
        :rtype: dict
        """
        action = action_manager.get_id(id_=action_id)
        if not isinstance(action, Action):
            return gen_json_error({'description': 'failed to get action'},
                                  HTTP_ERROR)

        return gen_json(action.to_dict())

    @ws.application.post(
        '/api/v2/actions'
    )
    def create_action():
        """
        Create a new action.

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

        # If element.get(...) fails, it returns None, which is evaluated as False
        # Empty elements such as {}, [] and "" are also evaluated to False
        # So element.get() will check if each value exists AND if the value isn't empty
        if element.get("hook") and (element.get("fields") or element.get("regex")):
            return gen_json_error(
                {'description': 'Sent data has a hook and fields/regex defined at the same time'}, HTTP_ERROR)

        delay = element.get("delay")
        if delay and not action_manager.is_delay_valid(delay):
            return gen_json_error({"description": "delay value is invalid"}, HTTP_ERROR)

        now = int(time.time())
        element['creation_date'] = now
        element['last_update_date'] = now

        try:
            if element.get('type', '') == 'pbehavior':
                _sanitize_pbehavior_parameter(element)
            Action(**Action.convert_keys(element))
        except TypeError:
            return gen_json_error(
                {'description': 'invalid action format'}, HTTP_ERROR)
        except ValueError:
            return gen_json_error(
                {'description': 'invalid pbehavior parameters'}, HTTP_ERROR)

        try:
            ok = action_manager.create(action=element)
        except CollectionError as ce:
            ws.logger.error('Action creation error : {}'.format(ce))
            return gen_json_error(
                {'description': 'error while creating an action'},
                HTTP_ERROR
            )

        if not ok:
            return gen_json_error({'description': 'failed to create action'},
                                  HTTP_ERROR)

        return gen_json({'_id': element['_id']})

    @ws.application.put(
        '/api/v2/actions/<action_id:id_filter>'
    )
    def update_action(action_id):
        """
        Update an existing alarm action.

        :param action_id: ID of the action
        :type action_id: str
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
            Action(**Action.convert_keys(element))
        except TypeError:
            return gen_json_error(
                {'description': 'invalid action format'}, HTTP_ERROR)

        try:
            ok = action_manager.update_id(id_=action_id, action=element)
        except CollectionError as ce:
            ws.logger.error('Action update error : {}'.format(ce))
            return gen_json_error(
                {'description': 'error while updating an action'},
                HTTP_ERROR
            )
        if not ok:
            return gen_json_error({'description': 'failed to update action'},
                                  HTTP_ERROR)

        return gen_json({})

    @ws.application.delete(
        '/api/v2/actions/<action_id:id_filter>'
    )
    def delete_id(action_id):
        """
        Delete a action, based on his id.

        :param action_id: ID of the action
        :type action_id: str

        :rtype: bool
        """
        ws.logger.info('Delete action : {}'.format(action_id))

        return gen_json(action_manager.delete_id(id_=action_id))
