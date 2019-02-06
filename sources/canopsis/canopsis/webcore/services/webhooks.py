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
from canopsis.webhooks import WebhookManager
from canopsis.webcore.utils import (gen_json, gen_json_error,
                                    HTTP_NOT_FOUND, HTTP_ERROR)


def exports(ws):

    webhook_manager = WebhookManager(WebhookManager.default_collection())

    @ws.application.get(
        '/api/v2/webhook/<webhook_id>'
    )
    def get_webhook_by_id(webhook_id):
        try:
            document = webhook_manager.get_webhook_by_id(webhook_id)
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the webhook data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json(document)

    @ws.application.post(
        '/api/v2/webhook'
    )
    @ws.application.put(
        '/api/v2/webhook/<webhook_id>'
    )
    def upsert(webhook_id=None):
        try:
            webhook = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'Invalid JSON'},
                HTTP_ERROR
            )

        if webhook is None or not isinstance(webhook, dict):
            return gen_json_error(
                {'description': 'nothing to upsert'}, HTTP_ERROR)

        try:
            ok = webhook_manager.upsert_webhook(webhook, webhook_id)
        except CollectionError as ce:
            ws.logger.error('Webhook error : {}'.format(ce))
            return gen_json_error(
                {'description': 'error while upserting an webhook'},
                HTTP_ERROR
            )

        if not ok:
            return gen_json_error({'description': 'failed to upsert webhook'},
                                  HTTP_ERROR)

        return gen_json({})

    @ws.application.delete(
        '/api/v2/webhook/<webhook_id>'
    )
    def delete_webhook_by_id(webhook_id):
        try:
            ok = webhook_manager.delete_webhook_by_id(webhook_id)
        except PyMongoError:
            return gen_json_error(
                {"description": "Can not retrieve the webhook data from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        return gen_json({"status": ok})
