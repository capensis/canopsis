# -*- coding: utf-8 -*-

from __future__ import unicode_literals
from bottle import request, response, install
from json import loads
from six import string_types
import os.path

from canopsis.common.converters import id_filter
from canopsis.common.ws import route
from canopsis.metaalarmrule.manager import MetaAlarmRuleManager
from canopsis.webcore.utils import HTTP_ERROR, HTTP_NOT_FOUND

VALID_PARAMS = [
    '_id', 'name', 'type', 'patterns', 'config', 'auto_resolve'
]

VALID_CONFIG_PARAMS = {
    'time_interval': ['timebased', 'complex', 'valuegroup'],
    'threshold_count': ['complex', 'valuegroup'],
    'threshold_rate': ['complex'],
    'alarm_patterns': ['attribute', 'complex', 'valuegroup'],
    'entity_patterns': ['attribute', 'complex', 'valuegroup'],
    'event_patterns': ['attribute', 'complex', 'valuegroup'],
    'attribute_patterns': ['attribute', 'complex'],
    'value_path': ['valuegroup']
}

VALID_RULE_TYPES = [
    'relation', 'timebased', 'attribute', 'complex', 'valuegroup'
]

import yaml
from bottle_swagger import SwaggerPlugin

def init_swagger():
    this_dir = os.path.dirname(os.path.abspath(__file__))
    with open("{}/swagger/swagger.yml".format(this_dir)) as f:
        swagger_def = yaml.load(f)

    swagger_plugin = SwaggerPlugin(swagger_def, ignore_undefined_api_routes=True, serve_swagger_ui=True)
    install(swagger_plugin)

def _set_status(status):
    if isinstance(status, int):
        response.status = status

class RouteHandlerMetaAlarmRule(object):
    def __init__(self, ma_rule_manager):
        self.ma_rule_manager = ma_rule_manager

    def _sanitize(self, name, rule_type, patterns, config, ma_rule_id, auto_resolve=False):
        if rule_type not in VALID_RULE_TYPES:
            raise ValueError("rule type invalid value {}".format(rule_type))
        if name is not None and not isinstance(name, string_types):
            raise ValueError("name has invalid value: {}".format(name))
        if ma_rule_id is not None and not isinstance(ma_rule_id, string_types):
            raise ValueError("_id has invalid value: {}".format(ma_rule_id))

        if isinstance(auto_resolve, bool):
            raise ValueError("invalid auto_resolve value type {}".format(auto_resolve))

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
        return name, rule_type, patterns, config, ma_rule_id

    def create(self, name, rule_type, patterns, config, ma_rule_id=None, auto_resolve=False):
        name, rule_type, patterns, config, ma_rule_id = self._sanitize(name, rule_type, patterns, config, ma_rule_id, auto_resolve)
        result = self.ma_rule_manager.create(name, rule_type, patterns, config, ma_rule_id=ma_rule_id)
        return result

    def update(self, _id, name, rule_type, patterns, config, auto_resolve=False):
        name, rule_type, patterns, config, _id = self._sanitize(name, rule_type, patterns, config, _id, False)
        result = self.ma_rule_manager.update(_id, name, rule_type, patterns, config)
        return result

    def read(self, rule_id):
        return self.ma_rule_manager.read(rule_id)

    def read_all(self):
        return self.ma_rule_manager.read_all()

    def delete(self, rule_id):
        return self.ma_rule_manager.delete(rule_id)


def exports(ws):
    try:
        init_swagger()
    except Exception as exc:
        ws.logger.exception("init_swagger exception {}".format(exc))
    else:
        ws.logger.info("init_swagger done")

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
            _set_status(HTTP_ERROR)
            return {'description': 'invalid JSON'}

        invalid_keys = []

        for key in elements.keys():
            if key not in VALID_PARAMS:
                invalid_keys.append(key)
                elements.pop(key)

        if len(invalid_keys) != 0:
            ws.logger.error('Invalid keys {} in payload'.format(invalid_keys))

        ma_rule_id = elements.get("_id")
        if ma_rule_id == "":
            ma_rule_id = None

        try:
            return rh_ma_rule.create(
                elements["name"], elements["type"], elements.get("patterns"), elements.get("config"), 
                ma_rule_id=ma_rule_id, auto_resolve=elements.get("auto_resolve", False))
        except (TypeError, KeyError):
            _set_status(HTTP_ERROR)
            return {'description': 'The fields \'name\' and \'type\' are required.'}
        except ValueError as exc:
            _set_status(HTTP_ERROR)
            return {'description': 'The fields \'name\' and \'type\' are required.'}

    @ws.application.put('/api/v2/metaalarmrule/<rule_id>')
    def update(rule_id):
        """
        Create a metaalarmrule.
        """
        try:
            elements = request.json
        except ValueError:
            _set_status(HTTP_ERROR)
            return {'description': 'invalid JSON'}
        invalid_keys = []

        for key in elements.keys():
            if key not in VALID_PARAMS:
                invalid_keys.append(key)
                elements.pop(key)

        if len(invalid_keys) != 0:
            ws.logger.error('Invalid keys {} in payload'.format(invalid_keys))

        try:
            success = rh_ma_rule.update(rule_id, elements["name"], elements["type"], elements.get("patterns"),
                                     elements.get("config"), auto_resolve=elements.get("auto_resolve", False))
        except Exception as exc:
            _set_status(HTTP_ERROR)
            return {'description': '{}'.format(exc)}
        return {"is_success": success}

    @ws.application.get('/api/v2/metaalarmrule/<rule_id>')
    def read(rule_id=None):
        r = rh_ma_rule.read(rule_id)
        if r is None:
            _set_status(HTTP_NOT_FOUND)
            return {'description': 'Rule ID={} not found'.format(rule_id)}
        if r.get("patterns") is None:
            r["patterns"] = {}
        if r.get("config") is None:
            r["config"] = {}
        return r

    @ws.application.get('/api/v2/metaalarmrule')
    def read_all():
        """
        :return:
        """
        return rh_ma_rule.read_all()

    @ws.application.delete('/api/v2/metaalarmrule/<rule_id>')
    def delete_rule(rule_id):
        """Delete the meta-alarm rule that match the rule_id

        :param rule_id: the meta-alarm rule id
        :return: a dict with two field. "acknowledged" that True when
        delete is a sucess. False, otherwise.
        :rtype: dict
        """
        ws.logger.info('Delete meta-alarm rule: {}'.format(rule_id))

        return rh_ma_rule.delete(rule_id)
