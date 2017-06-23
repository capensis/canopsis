# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

"""
Python webcore utility library.
"""

from __future__ import unicode_literals

from bottle import response
import json

HTTP_OK = 200
HTTP_ERROR = 400
HTTP_NOT_FOUND = 404
HTTP_NOT_ALLOWED = 405

CONTENT_TYPE_JSON = "application/json"

ERR_ALLOWED_KEYS = ["name", "description"]

MSG_ELT_NOT_DICT = "The element inside a JSON should be a dict"


def __gen_json(result, status, allowed_keys=None):
    """Generate a standard JSON, set the response to the given status and
    add a Content-Type header with "application/json".

    If allowed_keys is set, only the keys contained in allowed_keys will
    be added in the JSON or in the element inside the JSON

    :param dict,list result: the data the API send back
    :param int status: the status to return
    :param list allowed_keys: the list of the key the json should contains
    :return str: a string representation of the JSON.
    """
    json_body = {}

    response.content_type = CONTENT_TYPE_JSON
    if isinstance(status, int):
        response.status = status

    if isinstance(result, list):
        json_body = []

        for element in result:
            json_element = {}
            if not isinstance(element, dict):
                raise ValueError(MSG_ELT_NOT_DICT)

            if allowed_keys is not None:
                for key in allowed_keys:
                    json_element[key] = element.get(key, "")
            else:
                json_element = element.copy()

            json_body.append(json_element)

    elif isinstance(result, dict):
        if allowed_keys is None:
            json_body = result
        else:
            json_body = {}
            for key in allowed_keys:
                json_body[key] = result.get(key, "")

    return json.dumps(json_body)


def gen_json_error(result, status):
    """Generate a standard error JSON, set the response to the given
    status and header "Content-type" to "application/json".

    :param dict,list result: the data the API send back
    :param int status: the status to return
    :return str: a string representation of the JSON.
    """
    return __gen_json(result, status, ERR_ALLOWED_KEYS)


def gen_json(result, status=200, allowed_keys=None):
    """Generate a standard JSON, set the response to the given
    status and header "Content-type" to "application/json".

    :param dict,list result: the data the API send back
    :param int status: the status to return, default 200
    :param list allowed_keys: the list of the keys allowed in the JSON
    :return str: a string representation of the JSON.
    """
    return __gen_json(result, status, allowed_keys=allowed_keys)
