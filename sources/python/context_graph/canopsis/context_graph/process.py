# -*- coding: utf-8 -*-

"""Module in chage of defining the graph context and updating
 it when it's needed."""

from __future__ import unicode_literals

from canopsis.task.core import register_task
from canopsis.context_graph.manager import ContextGraph
import time


context_graph_manager = ContextGraph()

cache_comp = set(['comp_ids'])
cache_conn = set(['conn_ids'])
cache_re = set(['re_ids'])


def update_cache():
    global cache_re
    global cache_comp
    global cache_conn
    cache = context_graph_manager.get_all_entities()
    cache_re = cache['re_ids']
    cache_comp = cache['comp_ids']
    cache_conn = cache['conn_ids']


def create_entity(logger, id, name, etype, depends=[], impact=[], measurements=[], infos={}):
    return {'_id': id,
            'type': etype,
            'name': name,
            'depends': depends,
            'impact': impact,
            'measurements': measurements,
            'infos': infos
            }


def update_depends_link(logger, ent1_id, ent2_id, context):
    """
    Update depends link between from entity ent1 to ent2 in the context.
    """

    try:
        ent_from = context[ent1_id]
    except KeyError:
        pass

    if ent2_id not in ent_from["depends"]:
        ent_from["depends"].append(ent2_id)


def update_link(logger, ent1_id, ent2_id, context):
    """
    Update depends link from entity ent1 to ent2 in the context.
    Basicaly, add ent2_id in the field "impact" of the entity
    identified by ent1_id the then add ent1_id in the field "depends"
    of ent2_id.
    """
    try:
        ent1 = context[ent1_id]
    except KeyError:
        pass
    try:
        ent2 = context[ent2_id]
    except KeyError:
        pass

    if ent1_id not in ent2["impact"]:
        ent2["impact"].append(ent1_id)

    if ent2_id not in ent1["depends"]:
        ent1["depends"].append(ent2_id)


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

    time_get_ctx = None
    time_process_ctx = None
    time_update_ctx = None

    comp_id = event['component']

    re_id = None
    if 'resource' in event.keys():
        re_id = '{0}/{1}'.format(event['resource'], comp_id)

    conn_id = '{0}/{1}'.format(event['connector'], event['connector_name'])

    start = time.time()
    related_ctx = context_graph_manager.get_entity([comp_id, re_id, conn_id])
    end = time.time()
    time_get_ctx = end - start

    start = time.time()

    comp_there = False
    re_there = False
    conn_there = False

    for entity in related_ctx:
        if entity["_id"] == comp_id:
            comp_there = True

        if entity["_id"] == re_id:
            re_there = True

        if entity["_id"] == conn_id:
            conn_there = True

    context = {}

    for doc in related_ctx:
        context[doc["_id"]] = doc

    if not comp_there:
        depends = [conn_id]
        if 'resource' in event.keys():
            depends.append(re_id)
            re = create_entity(logger,
                               re_id,
                               event['resource'],
                               'resource',
                               [conn_id],
                               [comp_id],
                               [],
                               {})  # FIXME handles info if needed
            context[re_id] = re
        # FIXME handles info if needed
        comp = create_entity(logger,
                             comp_id, comp_id, 'component', depends, [], [], {})
        context[comp_id] = comp

    # If comp did not exist, a resource is created above
    if re_id is not None and comp_there is True:
        if not re_there:
            re = create_entity(logger,
                               re_id,
                               event['resource'],
                               'resource',
                               [conn_id],
                               [comp_id],
                               [],
                               {})
            context[re_id] = re

        # update link between re_id and comp_id if needed
        update_depends_link(logger, re_id, comp_id, context)

    if not conn_there:
        impact = [comp_id]
        if re_id is not None:
            impact.append(re_id)
            conn = create_entity(logger,
                                 conn_id,
                                 event['connector_name'],
                                 'connector',
                                 [],
                                 impact,
                                 [],
                                 {})

            context[conn_id] = conn

        # update link from re to conn
        update_depends_link(logger, re_id, conn_id, context)

        # update link from conn to comp
        update_depends_link(logger, comp_id, conn_id, context)

    update_link(logger, comp_id, conn_id, context)

    if re_id is not None:
        update_link(logger, re_id, conn_id, context)

    end = time.time()
    time_process_ctx = end - start

    start = time.time()
    if not conn_there or not comp_there or not re_there:
        context_graph_manager.put_entities(context.values())
    end = time.time()
    time_update_ctx = end - start

    #logger.critical("get context : {0}\tprocess context :{1}\tupdate context :{2}".format(time_get_ctx, time_process_ctx, time_update_ctx))
    update_cache()


@register_task
def beat(engine, logger=None, **kwargs):
    update_cache()
    logger.error(cache_comp)
    logger.error(cache_re)
    logger.error(cache_conn)
