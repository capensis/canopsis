# -*- coding: utf-8 -*-

from __future__ import unicode_literals
from bottle import request, response
from json import loads
from six import string_types
import time

from canopsis.common.errors import NotFoundError
from canopsis.idle_rule.manager import IdleRuleManager
from canopsis.models.idle_rule import IdleRule
from canopsis.webcore.utils import HTTP_ERROR, HTTP_NOT_FOUND


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
            response.status = HTTP_ERROR
            return {'description': 'invalid JSON'}

        try:
            rule = IdleRule.new_from_dict(
                elements, rh_idle_rule.get_username(), int(time.time()))
        except (TypeError, ValueError, KeyError) as exception:
            response.status = HTTP_ERROR
            return {'description': 'invalid idle rule: {}'.format(
                    exception.message)}

        try:
            idle_rule_manager.create(rule)
        except ValueError as exception:
            response.status = HTTP_ERROR
            return {'description': 'failed to create idle rule: {}'.format(
                    exception.message)}

        return rule.as_dict()

    @ws.application.put('/api/v2/idle-rule/<rule_id>')
    def update(rule_id):
        """
        Create a idle_rule.
        """
        try:
            elements = request.json
        except ValueError:
            response.status = HTTP_ERROR
            return {'description': 'invalid JSON'},

        try:
            rule = IdleRule.new_from_dict(
                elements, rh_idle_rule.get_username(), int(time.time()))
        except (TypeError, ValueError, KeyError) as exception:
            response.status = HTTP_ERROR
            return {'description': 'invalid idle rule: {}'.format(
                    exception.message)}

        try:
            success = idle_rule_manager.update(rule_id, rule)
        except ValueError as exception:
            response.status = HTTP_ERROR
            return {'description': 'failed to update idle rule: {}'.format(
                    exception.message)}
        except NotFoundError as exception:
            response.status = HTTP_NOT_FOUND
            return {"description": exception.message}

        if not success:
            response.status = HTTP_ERROR
            return {"description": "failed to update idle rule"}

        return rule.as_dict()

    @ws.application.get('/api/v2/idle-rule/<rule_id>')
    def read(rule_id=None):
        result = rh_idle_rule.read(rule_id)
        if result is None:
            response.status = HTTP_NOT_FOUND
            return {
                'name': rule_id,
                'description': 'Rule not found',
            }
        return result

    @ws.application.get('/api/v2/idle-rule')
    def read_all():
        return rh_idle_rule.read_all()

    @ws.application.delete('/api/v2/idle-rule/<rule_id>')
    def delete_rule(rule_id):
        try:
            result = rh_idle_rule.delete(rule_id)
        except NotFoundError as exception:
            response.status = HTTP_NOT_FOUND
            return {"description": exception.message}

        return result
