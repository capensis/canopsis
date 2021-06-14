# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2019 "Capensis" [http://www.capensis.com]
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
The webcore.services.dynamic_infos module defines the /api/v2/dynamic-infos
API.
"""

from __future__ import unicode_literals

import time
from bottle import request

from canopsis.common.collection import CollectionError
from canopsis.common.converters import id_filter
from canopsis.common.errors import NotFoundError
from canopsis.dynamic_infos.manager import DynamicInfosManager
from canopsis.models.dynamic_infos import DynamicInfosRule
from canopsis.webcore.utils import (
    gen_json, gen_json_error, HTTP_ERROR, HTTP_NOT_FOUND
)
import copy
import canopsis.alerts.translate_search as ts

def get_username():
    """Returns the username of the logged-in user, or ''."""
    try:
        session = request.environ.get('beaker.session', {})
        user = session.get('user', '')

        # The content of user depends on the authentication method. If the user
        # logged in with HTTP authentication, it contains the username. If they
        # logged in with the loggin form, it contains a dictionnary.
        if isinstance(user, basestring):
            return user
        return user.get('_id', '')
    except AttributeError:
        return ''


def replace_pattern_key(filter_, replacing_key):
    if isinstance(filter_, list):
        for i, fil in enumerate(filter_):
            filter_[i] = replace_pattern_key(fil, replacing_key)

    elif isinstance(filter_, dict):
        for key, value in filter_.items():
            new_value = replace_pattern_key(value, replacing_key)
            filter_[key] = new_value

            new_key = key
            if new_key == "pattern":
                new_key= replacing_key
            filter_[new_key] = filter_.pop(key)

    return filter_


def translate_search(list_path, search):
    def f():
        filters = []
        for path in list_path:
            _, bnf_filter = ts.parse_search(search)
            filters.append(replace_pattern_key(bnf_filter, path))
        return {"$or": filters}

    return f()


def exports(ws):

    ws.application.router.add_filter('id_filter', id_filter)

    manager = DynamicInfosManager(*DynamicInfosManager.provide_default_basics())

    @ws.application.get(
        '/api/v2/dynamic-infos'
    )
    def list_rules():
        search = request.query.search or ""
        search_fields = [
            field.strip()
            for field in request.query.search_fields.split(',')
            if field.strip()
        ]

        bnf_search = None
        try:
            _, bnf_search = ts.parse_search(search)
        except:
            pass

        if bnf_search:
            list_pattern = manager.list_all_patterns()
            bnf_search = translate_search(list_pattern, search)

        try:
            limit = int(request.query.limit or '0')
        except ValueError:
            return gen_json_error(
                {"description": "limit should be an integer"},
                HTTP_ERROR)
        try:
            skip = int(request.query.skip or '0')
        except ValueError:
            return gen_json_error(
                {"description": "skip should be an integer"},
                HTTP_ERROR)

        sort_key = request.query.sort_key
        sort_dir = request.query.sort_dir
        if not sort_dir:
            sort_dir = manager.SORT_DIRECTION_ASC
        if sort_dir not in [None, manager.SORT_DIRECTION_ASC, manager.SORT_DIRECTION_DESC]:
            return gen_json_error(
                {"description": "sort_dir should be either ASC or DESC"},
                HTTP_ERROR)

        try:
            count = manager.count(search, search_fields, bnf_search)
            rules = manager.list(search, search_fields, limit, skip, sort_key, sort_dir, bnf_search)
        except CollectionError:
            return gen_json_error(
                {"description": "Cannot retrieve the dynamic infos list from "
                                "the database, please contact your "
                                "administrator."},
                HTTP_ERROR)
        except ValueError as exception:
            return gen_json_error(
                {"description": exception.message},
                HTTP_ERROR)

        return gen_json({
            'count': count,
            'rules': rules,
        })

    @ws.application.get(
        '/api/v2/dynamic-infos/<rule_id:id_filter>'
    )
    def get_by_id(rule_id):
        rule = manager.get_by_id(rule_id)
        if rule is None:
            return gen_json_error(
                {"description": "no dynamic infos rule with id {}".format(
                    rule_id)},
                HTTP_NOT_FOUND)

        return gen_json(rule)

    @ws.application.post(
        '/api/v2/dynamic-infos'
    )
    def create():
        try:
            body = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        try:
            rule = DynamicInfosRule.new_from_dict(
                body, get_username(), int(time.time()))
        except (TypeError, ValueError, KeyError) as exception:
            return gen_json_error(
                {'description': 'invalid dynamic infos: {}'.format(
                    exception.message)},
                HTTP_ERROR)

        try:
            manager.create(rule)
        except (ValueError, Exception) as exception:
            return gen_json_error(
                {'description': 'failed to create dynamic infos: {}'.format(
                    exception.message)},
                HTTP_ERROR)

        return gen_json(rule.as_dict())

    @ws.application.put(
        '/api/v2/dynamic-infos/<rule_id:id_filter>'
    )
    def update(rule_id):
        try:
            body = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        try:
            rule = DynamicInfosRule.new_from_dict(
                body, get_username(), int(time.time()))
        except (TypeError, ValueError, KeyError) as exception:
            return gen_json_error(
                {'description': 'invalid dynamic infos: {}'.format(
                    exception.message)},
                HTTP_ERROR)

        try:
            success = manager.update(rule_id, rule)
        except ValueError as exception:
            return gen_json_error(
                {'description': 'failed to update dynamic infos: {}'.format(
                    exception.message)},
                HTTP_ERROR)
        except NotFoundError as exception:
            return gen_json_error(
                {"description": exception.message},
                HTTP_NOT_FOUND)

        if not success:
            return gen_json_error(
                {"description": "failed to update dynamic infos"},
                HTTP_ERROR)

        return gen_json(rule.as_dict())


    @ws.application.delete(
        '/api/v2/dynamic-infos/<rule_id:id_filter>'
    )
    def delete(rule_id):
        try:
            success = manager.delete(rule_id)
        except NotFoundError as exception:
            return gen_json_error(
                {"description": exception.message},
                HTTP_NOT_FOUND)
        if not success:
            return gen_json_error(
                {"description": "failed to delete dynamic infos"},
                HTTP_ERROR)

        return gen_json({})
