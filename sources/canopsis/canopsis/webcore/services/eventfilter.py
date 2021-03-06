#!/usr/bin/env python
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
from pymongo.errors import AutoReconnect

from canopsis.common.collection import CollectionError
from canopsis.eventfilter.enums import RuleField
from canopsis.eventfilter.manager import RuleManager, InvalidRuleError
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR, \
    HTTP_NOT_FOUND


def exports(ws):

    rule_manager = RuleManager(ws.logger)

    @ws.application.get(
        '/api/v2/eventfilter/rules'
    )
    def list_rules():
        return gen_json(rule_manager.list())

    @ws.application.get(
        '/api/v2/eventfilter/rules/<rule_id>'
    )
    def get_rule(rule_id):
        try:
            rule = rule_manager.get_by_id(rule_id)
        except AutoReconnect as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        if rule:
            return gen_json(rule)
        else:
            return gen_json_error({
                'description': 'No rule with id: {0}'.format(rule_id)
            }, HTTP_NOT_FOUND)

    @ws.application.post(
        '/api/v2/eventfilter/rules'
    )
    def create_rule():
        try:
            request_body = request.json
        except ValueError as verror:
            return gen_json_error({
                'description': 'Malformed JSON: {0}'.format(verror)
            }, HTTP_ERROR)

        if not request_body:
            return gen_json_error({
                'description': 'Empty request'
            }, HTTP_ERROR)

        try:
            rule_id = rule_manager.create(request_body)
        except (InvalidRuleError, CollectionError) as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        return gen_json({
            RuleField.id: rule_id
        })

    @ws.application.delete(
        '/api/v2/eventfilter/rules/<rule_id>'
    )
    def remove_rule(rule_id):
        try:
            rule = rule_manager.get_by_id(rule_id)
        except AutoReconnect as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        if not rule:
            return gen_json_error({
                'description': 'No rule with id: {0}'.format(rule_id)
            }, HTTP_NOT_FOUND)

        try:
            rule_manager.remove_with_id(rule_id)
        except CollectionError as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        return gen_json({})

    @ws.application.put(
        '/api/v2/eventfilter/rules/<rule_id>'
    )
    def update_rule(rule_id):
        try:
            rule = rule_manager.get_by_id(rule_id)
        except AutoReconnect as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        if not rule:
            return gen_json_error({
                'description': 'No rule with id: {0}'.format(rule_id)
            }, HTTP_NOT_FOUND)

        try:
            request_body = request.json
        except ValueError as verror:
            return gen_json_error({
                'description': 'Malformed JSON: {0}'.format(verror)
            }, HTTP_ERROR)

        if not request_body:
            return gen_json_error({
                'description': 'Empty request'
            }, HTTP_ERROR)

        try:
            rule_manager.update(rule_id, request_body)
        except (InvalidRuleError, CollectionError) as e:
            return gen_json_error({
                'description': e.message
            }, HTTP_ERROR)

        return gen_json({})
