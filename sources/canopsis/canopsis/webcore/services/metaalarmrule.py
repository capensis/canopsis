# -*- coding: utf-8 -*-

from __future__ import unicode_literals
from bottle import request
from json import loads
from six import string_types
import os.path

from canopsis.common.converters import id_filter
from canopsis.common.ws import route
from canopsis.metaalarmrule.manager import MetaAlarmRuleManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR

VALID_PARAMS = [
    '_id', 'name', 'type', 'filter', 'enabled',
]

VALID_RULE_TYPES = [
    'relation', 'time', 'attribute', 'complex',
]


class RouteHandlerMetaAlarmRule(object):
    def __init__(self, ma_rule_manager):
        self.ma_rule_manager = ma_rule_manager

    def create(self, name, rule_type, filter, enabled):
        if rule_type not in VALID_RULE_TYPES:
            raise ValueError("rule type invalid value {}".format(rule_type))
        if isinstance(filter, string_types):
            try:
                filter = loads(filter)
            except ValueError:
                raise ValueError("Can't decode rule filter parameter: {}"
                                 .format(filter))
        result = self.ma_rule_manager.create(name, rule_type, filter, enabled)
        return result

    def read(self, rule_id):
        return self.ma_rule_manager.read(rule_id)

    def delete(self, rule_id):
        return self.ma_rule_manager.delete(rule_id)


def exports(ws):

    ws.application.router.add_filter('id_filter', id_filter)

    ma_rule_manager = MetaAlarmRuleManager(
        *MetaAlarmRuleManager.provide_default_basics())
    rh_ma_rule = RouteHandlerMetaAlarmRule(ma_rule_manager=ma_rule_manager)

    @ws.application.post('/api/v2/metaalarmrule')
    def create():
        """
        Create a metaalarmrule.
        """
        try:
            elements = request.json
        except ValueError:
            return gen_json_error(
                {'description': 'invalid JSON'},
                HTTP_ERROR
            )

        invalid_keys = []

        for key in elements.keys():
            if key not in VALID_PARAMS:
                invalid_keys.append(key)
                elements.pop(key)

        enabled = elements.get("enabled", False)
        if not isinstance(enabled, bool):
            if isinstance(enabled, string_types):
                enabled = enabled.lower()
            enabled = enabled in (1, "1", "true")

        if len(invalid_keys) != 0:
            ws.logger.error('Invalid keys {} in payload'.format(invalid_keys))

        try:
            return rh_ma_rule.create(elements["name"], elements["type"], elements["filter"], enabled)
        except TypeError:
            return gen_json_error(
                {'description': 'The fields name, type, filter are required.'},
                HTTP_ERROR
            )
        except ValueError as exc:
            return gen_json_error(
                {'description': '{}'.format(exc)},
                HTTP_ERROR
            )

    @ws.application.get('/api/v2/metaalarmrule/<entity_id:id_filter>')
    def read(entity_id=None):
        return gen_json(rh_ma_rule.read(entity_id))

    @ws.application.delete('/api/v2/metaalarmrule/<entity_id:id_filter>')
    def delete(entity_id):
        """Delete the meta-alarm rule that match the entity_id

        :param entity_id: the meta-alarm rule id
        :return: a dict with two field. "acknowledged" that True when
        delete is a sucess. False, otherwise.
        :rtype: dict
        """
        ws.logger.info('Delete meta-alarm rule: {}'.format(entity_id))

        return gen_json(rh_ma_rule.delete(entity_id))

