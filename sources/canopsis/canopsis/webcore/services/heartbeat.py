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
from bottle import request
from canopsis.heartbeat.manager import HeartBeatService
from canopsis.heartbeat.manager import HeartBeatServiceException
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR
from canopsis.models.heartbeat import HeartBeat

hb_service = HeartBeatService(*HeartBeatService.provide_default_basics())


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

        try:
            hb_service.create(HeartBeat(**json))
            return gen_json({"name": "heartbeat created"})

        except HeartBeatServiceException as exc:
            ws.logger.exception("Can not create heartbeat.")
            return gen_json_error({"name": "Can not create heartbeat",
                                   "description": exc.message},
                                  HTTP_ERROR)
        except Exception:
            ws.logger.exception("Can not create heartbeat.")
            return gen_json_error(
                {"name": "Can not create heartbeat",
                 "description": "Contact the administrator."},
                HTTP_ERROR)


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
            return hb_service.get_heartbeats()
        except Exception:
            ws.logger.exception("Can not retreive hearbeats from database.")
            return gen_json_error({'description': 'something went wrong.'},
                                  HTTP_ERROR)
