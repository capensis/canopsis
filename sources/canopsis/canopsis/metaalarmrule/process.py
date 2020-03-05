# -*- coding: utf-8 -*-

from canopsis.metaalarmrule.manager import MetaAlarmRuleManager
from canopsis.common.utils import singleton_per_scope

def init_managers():
    """
    Init managers [sic].
    """
    config, mar_logger, mar_collection = MetaAlarmRuleManager.provide_default_basics()
    ma_rule_kwargs = {'config': config,
                 'logger': mar_logger,
                 'ma_rule_collection': mar_collection}
    ma_rule_manager = singleton_per_scope(MetaAlarmRuleManager, kwargs=ma_rule_kwargs)

    return ma_rule_manager

_ma_rule_manager = init_managers()
