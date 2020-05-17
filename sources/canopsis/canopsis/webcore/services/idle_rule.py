# -*- coding: utf-8 -*-

from __future__ import unicode_literals
from bottle import request
from json import loads
from six import string_types
import time
import os.path

from canopsis.common.errors import NotFoundError
from canopsis.idle_rule.manager import IdleRuleManager
from canopsis.models.idle_rule import IdleRule
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR, HTTP_NOT_FOUND


class RouteHandlerIdleRule(object):
    def __init__(self, idle_rule_manager):
        self.idle_rule_manager = idle_rule_manager

    @staticmethod
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

    def read(self, rule_id):
        return self.idle_rule_manager.get_by_id(rule_id)

    def read_all(self):
        return self.idle_rule_manager.read_all()

    def delete(self, rule_id):
        return self.idle_rule_manager.delete(rule_id)


def exports(ws):

    idle_rule_manager = IdleRuleManager(
        *IdleRuleManager.provide_default_basics())
    rh_idle_rule = RouteHandlerIdleRule(idle_rule_manager=idle_rule_manager)

    @ws.application.post('/api/v2/idle-rule')
    def create():
        """
        Create a idle_rule.
        """
        try:
            elements = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        try:
            rule = IdleRule.new_from_dict(
                elements, rh_idle_rule.get_username(), int(time.time()))
        except (TypeError, ValueError, KeyError) as exception:
            return gen_json_error(
                {'description': 'invalid idle rule: {}'.format(
                    exception.message)},
                HTTP_ERROR)

        try:
            idle_rule_manager.create(rule)
        except ValueError as exception:
            return gen_json_error(
                {'description': 'failed to create idle rule: {}'.format(
                    exception.message)},
                HTTP_ERROR)

        return gen_json(rule.as_dict())

    @ws.application.put('/api/v2/idle-rule/<rule_id>')
    def update(rule_id):
        """
        Create a idle_rule.
        """
        try:
            elements = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        try:
            rule = IdleRule.new_from_dict(
                elements, rh_idle_rule.get_username(), int(time.time()))
        except (TypeError, ValueError, KeyError) as exception:
            return gen_json_error(
                {'description': 'invalid idle rule: {}'.format(
                    exception.message)},
                HTTP_ERROR)

        try:
            success = idle_rule_manager.update(rule_id, rule)
        except ValueError as exception:
            return gen_json_error(
                {'description': 'failed to update idle rule: {}'.format(
                    exception.message)},
                HTTP_ERROR)
        except NotFoundError as exception:
            return gen_json_error(
                {"description": exception.message},
                HTTP_NOT_FOUND)

        if not success:
            return gen_json_error(
                {"description": "failed to update idle rule"},
                HTTP_ERROR)

        return gen_json(rule.as_dict())

    @ws.application.get('/api/v2/idle-rule/<rule_id>')
    def read(rule_id=None):
        return gen_json(rh_idle_rule.read(rule_id))

    @ws.application.get('/api/v2/idle-rule')
    def read_all():
        return gen_json(rh_idle_rule.read_all())

    @ws.application.delete('/api/v2/idle-rule/<rule_id>')
    def delete_rule(rule_id):
        try:
            result = rh_idle_rule.delete(rule_id)
        except NotFoundError as exception:
            return gen_json_error(
                {"description": exception.message},
                HTTP_NOT_FOUND)

        return gen_json(result)
