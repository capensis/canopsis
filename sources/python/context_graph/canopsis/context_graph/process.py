# -*- coding: utf-8 -*-

"""Module in chage of defining the graph context and updating
 it when it's needed."""

from __future__ import unicode_literals

from canopsis.task.core import register_task
from canopsis.context_graph.manager import ContextGraph
from canopsis.common.utils import singleton_per_scope

from canopsis.context_graph.manager import ContextGraph


context_graph_manager = ContextGraph()

@register_task
def event_processing(
        engine, event, manager=None, logger=None, ctx=None, tm=None, cm=None,
        **kwargs
):
    """event_processing

    :param engine:
    :param event:
    :param manager:
    :param logger:
    :param ctx:
    :param tm:
    :param cm:
    :param **kwargs:
    """
    comp_id = event['component']
    re_id = '{0}/{1}'.format(event['resource'], comp_id)
    conn_id = '{0}/{1}'.format(event['connector'], event['connector_name'])


    if not context_graph_manager.check_comp(comp_id):
        depends = [conn_id]
        if 'resource' in event.keys():
            depends.append(re_id)
            context_graph_manager.add_re({
                '_id': re_id,
                'name': event['resource'],
                'type': 'resource',
                'depends': [conn_id],
                'impact': [comp_id],
                'measurements': [],
                'infos': {}
            })
        context_graph_manager.add_comp({
            '_id': comp_id,
            'name': comp_id,
            'type': 'component',
            'depends': depends,
            'impact': [],
            'measurements': [],
            'infos': {}
        })
    if not context_graph_manager.check_re(re_id):
        context_graph_manager.add_re({
                '_id': re_id,
                'name': event['resource'],
                'type': 'resource',
                'depends': [conn_id],
                'impact': [comp_id],
                'measurements': [],
                'infos': {}
        })
        context_graph_manager.manage_comp_to_re_link(re_id, comp_id)
    if not context_graph_manager.check_conn(conn_id):
        logger.error('add connector')
        context_graph_manager.add_conn({
            '_id': conn_id,
            'name': event['connector_name'],
            'type': 'connector',
            'depends': [],
            'impact': [comp_id, re_id],
            'measurements': [],
            'infos': {}
        })
        context_graph_manager.manage_re_to_conn_link(conn_id, re_id)
        context_graph_manager.manage_comp_to_conn_link(conn_id, comp_id)

