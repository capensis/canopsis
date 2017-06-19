# -*- coding:utf-8 -*- 
from canopsis.task.core import register_task
from canopsis.selector.manager import Selector

selector_manager = Selector()

@register_task
def beat_processing(engine, logger=None, **kwargs):
    """
        selector's beat processing to launch sla calculs
    """
    logger.critical('beat_selector')
    selector_manager.calcul_slas()

