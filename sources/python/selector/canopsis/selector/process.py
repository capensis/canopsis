# -*- coding:utf-8 -*- 
from canopsis.task.core import register_task
from canopsis.selector.manager import Selector

selector_manager = Selector()

@register_task
def beat_processing(engine, logger=None, **kwargs):
    """beat_processing
    selector's beat processing

    :param engine: selector engine
    :param logger: engin logger
    :param **kwargs:
    """
    selector_manager.compute_slas()

