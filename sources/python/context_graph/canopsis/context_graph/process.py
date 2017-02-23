# -*- coding: utf-8 -*-

"""Module in chage of defining the graph context and updating
 it when it's needed."""

from __future__ import unicode_literals

from canopsis.task.core import register_task
from canopsis.context_graph.manager import ContextGraph


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


def create_entity(
        logger,
        id,
        name,
        etype,
        depends=[],
        impact=[],
        measurements=[],
        infos={}):
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


def update_case1(entities, ids):
    """Case 1 update entities"""
    comp_there = False
    re_there = False
    conn_there = False
    for i in entities:
        if i['type'] == 'component':
            comp_there = True
        elif i['type'] == 'resource':
            re_there = True
        elif i['type'] == 'connector':
            conn_there = True



def update_case2(entities, ids):
    """Case 2 update entities"""
    comp_there = False
    re_there = False
    conn_there = False
    for i in entities:
        if i['type'] == 'component':
            comp_there = True
        elif i['type'] == 'resource':
            re_there = True
        elif i['type'] == 'connector':
            conn_there = True


def update_case3(entities, ids):
    """Case 3 update entities"""


def update_case4(entities, ids):
    """Case 4 update entities"""
    pass


def update_case5(entities, ids):
    """Case 5 update entities"""
    

def update_case6(entities, ids):
    """Case 6 update entities"""
    conn_there = False:
    for i in entities:
        if i['_id'] in ids:
            conn_there = True
    if not conn_there:
        #insert conn + maj impact depends in comp and re


def update_entities(case, ids):
    tab_entities = context_graph_manager.get_entity(ids)
    if case == 1:
        update_case1(tab_entities, ids)
    elif case == 2:
        update_case2(tab_entities, ids)
    elif case == 3:
        update_case3(tab_entities, ids)
    elif case == 4:
        update_case4(tab_entities, ids)
    elif case == 5:
        update_case5(tab_entities, ids)
    elif case == 6:
        update_case6(tab_entities, ids)
    else:
        # FIXME : did a logger here is a good idea ?
        raise ValueError("Unknown case : {0}.".format("case"))


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

    # Possible cases :
    # 0 -> Not in cache
    # 1 −> In cache
    #
    #
    #  Connector    Ressource    Component
    #     0             0            0     -> case 1
    #     0             0            1     -> case 2
    #     1             0            1     -> case 3
    #     1             1            1     -> case 4
    #     1             0            0     -> case 5
    #     1             1            0     -> case 6
    #
    #  Case 1 :
    #    Nothing exist in the cache, create every entities in database.
    #
    #  Case 2 :
    #    The component and the resource are not in the cache, so we need
    #    to create a component and if the event had a resource
    #    (resource != None) create a resource too, then update links between
    #    this component and resource.
    #
    #  Case 3 :
    #    Create a component and if the event had a resource create a resource
    #    too, then update links between this component and resource.
    #
    #  Case 4 :
    #    Every entity are in the cache so they are into the database, nothing
    #    to do here.
    #
    #  Case 5 :
    #    Create a component and the resource then update the links between
    #    the component and the resource.
    #
    #  Case 6 :
    #    Create a component then update the links between the component and
    #    the resource.

    case = 0  # 0 -> exception raised
    comp_id = event['component']

    re_id = None
    if 'resource' in event.keys():
        re_id = '{0}/{1}'.format(event['resource'], comp_id)

    conn_id = '{0}/{1}'.format(event['connector'], event['connector_name'])

    # add comment with the 6 possibles cases and an explaination

    # cache and case determination
    ids = set([])

    if conn_id in cache_conn:
        if comp_id in cache_comp:
            if re_id is not None:
                if re_id not in cache_re:
                    # 3
                    ids.add(re_id)
                    ids.add(comp_id)
                    ids.add(conn_id)
                    cache_re.add(re_id)
                    case = 3
                # else:
                    # 4 => pass
        else:
            # 2
            case = 2
            ids.add(comp_id)
            ids.add(conn_id)
            cache_comp.add(comp_id)
            if re_id is not None:
                if re_id not in cache_re:
                    ids.add(re_id)
                    cache_re.add(re_id)
    else:
        # 6
        case = 6
        cache_conn.add(conn_id)
        ids.add(conn_id)
        ids.add(comp_id)
        if comp_id in cache_comp:
            if re_id is not None:
                if re_id not in cache_re:
                    # 5
                    case = 5
                    ids.add(re_id)
                    cache_re.add(re_id)
        else:
            # 1
            case = 1
            cache_comp.add(comp_id)
            if re_id is not None:
                cache_re.add(re_id)
                ids.add(re_id)

    # retrieves required entities from database
    update_entities(case, ids)
