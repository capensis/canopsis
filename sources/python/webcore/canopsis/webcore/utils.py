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

import json

APP_TYPE = "application/json"

ERR_ALLOWED_KEYS = ["name", "description"]


def __gen_json(response, result, status, allowed_keys=None):
    """Generate a standard JSON, set the response to the given status and
    add a Content-Type header with "application/json".

    If allowed_keys is set, only the keys contained in allowed_keys will
    be added in the json.

    :param response: the bottle response
    :param dict result: the data the API send back
    :param int status: the status to return
    :param list allowed_keys: the list of the key the json should contains
    :return str: a string representation of the JSON.
    """

    response.status = status
    response.content_type = APP_TYPE

    if allowed_keys is None:
        json_ = result
    else:
        json_ = {}
        for key in allowed_keys:
            json_[key] = result.get(key, "")

    return json.dumps(json_)


def gen_json_error(response, result, status):
    """Generate a standard error JSON, set the response to the given
    status and header "Content-type" to "application/json".

    :param response: the bottle response
    :param int status: the status to return
    :return str: a string representation of the JSON.
    """
    return __gen_json(response, result, status, ERR_ALLOWED_KEYS)


def gen_json(response, result, status=200):
    """Generate a standard JSON, set the response to the given
    status and header "Content-type" to "application/json".

    :param response: the bottle response
    :param int status: the status to return
    :return str: a string representation of the JSON.
    """
    return __gen_json(response, result, status)
