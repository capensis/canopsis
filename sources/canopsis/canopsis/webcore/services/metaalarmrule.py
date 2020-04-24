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
    '_id', 'name', 'type', 'patterns', 'config',
]

VALID_CONFIG_PARAMS = {
    'time_interval': ['timebased', 'complex'],
    'threshold_count': ['complex'],
    'threshold_rate': ['complex'],
    'alarm_patterns': ['attribute', 'complex'],
    'entity_patterns': ['attribute', 'complex'],
    'event_patterns': ['attribute', 'complex'],
    'attribute_patterns': ['attribute', 'complex'],
}

VALID_RULE_TYPES = [
    'relation', 'timebased', 'attribute', 'complex',
]


class RouteHandlerMetaAlarmRule(object):
    def __init__(self, ma_rule_manager):
        self.ma_rule_manager = ma_rule_manager

    def _sanitize(self, rule_type, patterns, config):
        if rule_type not in VALID_RULE_TYPES:
            raise ValueError("rule type invalid value {}".format(rule_type))
        if isinstance(patterns, string_types):
            try:
                patterns = loads(patterns)
            except ValueError:
                raise ValueError("Can't decode rule patterns parameter: {}"
                                 .format(patterns))

        if isinstance(config, string_types):
            try:
                config = loads(config)
            except ValueError:
                raise ValueError("Can't decode rule config parameter: {}"
                                 .format(config))

        if isinstance(config, dict):
            for config_type, valid_rule_types in VALID_CONFIG_PARAMS.items():
                if config_type in config and rule_type not in valid_rule_types:
                    raise ValueError("invalid rule_type {} with config {}"
                                     .format(rule_type, config_type))
        elif config is not None:
            raise ValueError("invalid config value type {}".format(config))
        return rule_type, patterns, config

    def create(self, name, rule_type, patterns, config):
        rule_type, patterns, config = self._sanitize(rule_type, patterns, config)
        result = self.ma_rule_manager.create(name, rule_type, patterns, config)
        return result

    def update(self, _id, name, rule_type, patterns, config):
        rule_type, patterns, config = self._sanitize(rule_type, patterns, config)
        result = self.ma_rule_manager.update(_id, name, rule_type, patterns, config)
        return result

    def read(self, rule_id):
        return self.ma_rule_manager.read(rule_id)

    def read_all(self):
        return self.ma_rule_manager.read_all()

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

        if len(invalid_keys) != 0:
            ws.logger.error('Invalid keys {} in payload'.format(invalid_keys))

        try:
            return rh_ma_rule.create(elements["name"], elements["type"], elements.get("patterns"), elements.get("config"))
        except (TypeError, KeyError):
            return gen_json_error(
                {'description': 'The fields \'name\' and \'type\' are required.'},
                HTTP_ERROR
            )
        except ValueError as exc:
            return gen_json_error(
                {'description': '{}'.format(exc)},
                HTTP_ERROR
            )

    @ws.application.put('/api/v2/metaalarmrule/<rule_id>')
    def update(rule_id):
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

        if len(invalid_keys) != 0:
            ws.logger.error('Invalid keys {} in payload'.format(invalid_keys))

        try:
            success = rh_ma_rule.update(rule_id, elements["name"], elements["type"], elements.get("patterns"),
                                     elements.get("config"))
        except Exception as exc:
            return gen_json_error(
                {'description': '{}'.format(exc)},
                HTTP_ERROR
            )
        return gen_json({"is_success": success})

    @ws.application.get('/api/v2/metaalarmrule/<rule_id:id_filter>')
    def read(rule_id=None):
        return gen_json(rh_ma_rule.read(rule_id))

    @ws.application.get('/api/v2/metaalarmrule')
    def read_all():
        """
        :return:
        """
        return gen_json(rh_ma_rule.read_all())

    @ws.application.delete('/api/v2/metaalarmrule/<rule_id:id_filter>')
    def delete(rule_id):
        """Delete the meta-alarm rule that match the rule_id

        :param rule_id: the meta-alarm rule id
        :return: a dict with two field. "acknowledged" that True when
        delete is a sucess. False, otherwise.
        :rtype: dict
        """
        ws.logger.info('Delete meta-alarm rule: {}'.format(rule_id))

        return gen_json(rh_ma_rule.delete(rule_id))
