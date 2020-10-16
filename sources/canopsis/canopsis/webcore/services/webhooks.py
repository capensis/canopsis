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
from canopsis.common.collection import CollectionError
from canopsis.webhooks import WebhookManager
from canopsis.webcore.utils import (gen_json, gen_json_error, HTTP_ERROR)
import time


def exports(ws):

    webhook_manager = WebhookManager(WebhookManager.default_collection())

    @ws.application.post(
        '/api/v2/webhook'
    )
    def create_webhook():
        """
        Create a new webhook.

        :returns: ID of the webhook
        :rtype: string
        """
        try:
            webhook = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'Invalid JSON'},
                HTTP_ERROR
            )

        if webhook is None or not isinstance(webhook, dict):
            return gen_json_error(
                {'description': 'Nothing to create'}, HTTP_ERROR)

        if '_id' not in webhook:
            webhook['_id'] = str(uuid.uuid4())
        else:
            webhook['_id'] = webhook['_id'].strip()

        now = int(time.time())
        webhook['creation_date'] = now
        webhook['last_update_date'] = now

        try:
            _id = webhook_manager.create_webhook(webhook)
            return {'_id': _id}
        except CollectionError as ce:
            ws.logger.error('Webhook creation error : {}'.format(ce))
            return gen_json_error(
                {'description': 'Error while creating an webhook'},
                HTTP_ERROR
            )

    @ws.application.put(
        '/api/v2/webhook/<webhook_id:re:.+>'
    )
    def update_webhook_by_id(webhook_id):
        """
        Update an existing webhook.

        :param webhook_id: ID of the webhook
        :type webhook_id: str
        :rtype: dict
        """
        try:
            webhook = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'Invalid JSON'},
                HTTP_ERROR
            )

        if webhook is None or not isinstance(webhook, dict):
            return gen_json_error(
                {'description': 'Nothing to update'}, HTTP_ERROR)

        try:
            ok = webhook_manager.update_webhook_by_id(webhook, webhook_id)
        except CollectionError as ce:
            ws.logger.error('Webhook update error : {}'.format(ce))
            return gen_json_error(
                {'description': 'Error while updating an webhook'},
                HTTP_ERROR
            )

        if not ok:
            return gen_json_error(
                {'description': 'Failed to update webhook'},
                HTTP_ERROR
            )

        return gen_json({})

    @ws.application.delete(
        '/api/v2/webhook/<webhook_id:re:.+>'
    )
    def delete_webhook_by_id(webhook_id):
        """
        Delete an existing webhook, given its id.

        :param webhook_id: ID of the webhook
        :type webhook_id: str
        :rtype: dict
        """
        try:
            ok = webhook_manager.delete_webhook_by_id(webhook_id)
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the webhook data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json({"status": ok})

    @ws.application.get('/api/v2/webhook')
    def read():
        """
        Get a pbehavior.
        """
        _id = request.params.get('_id', None)
        search = request.params.get('search', None)
        try:
            limit = int(request.params.get('limit', None))
        except:
            limit = None
        try:
            skip = int(request.params.get('skip', None))
        except:
            skip = None
        return webhook_manager.read(_id, search, limit, skip)
