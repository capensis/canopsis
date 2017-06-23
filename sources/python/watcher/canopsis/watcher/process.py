# -*- coding:utf-8 -*-
from canopsis.task.core import register_task
from canopsis.watcher.manager import Watcher

watcher_manager = Watcher()


@register_task
def beat_processing(engine, logger=None, **kwargs):
    """beat_processing
    watcher's beat processing

    :param engine: watcher engine
    :param logger: engin logger
    :param **kwargs:
    """
    #watcher_manager.compute_slas()
