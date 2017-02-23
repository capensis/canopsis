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

LOGGER = None


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


def check_type(entities, expected):
    """Raise TypeError if the type of the entities entities does not match
    the expected type.
    :param entities: the entities to check.
    :param expected: the expected type.
    :raises TypeError: if the entity does not match the expected one."""
    if entities["type"] != expected:
        raise TypeError("Entities {0} does not match {1}".format(
            entities["id"], expected))


def update_depends_links(ent_from, ent_to):
    """Update the links depends from ent_from to ent_to. Basicaly, append
    the id of ent_to on the the "depends" of ent_from.
    :param ent_from: the entities that will be updated
    :param ent_to: the entities which id will be used to update ent_from"""
    if ent_to["_id"] not in ent_from["depends"]:
        ent_from["depends"] = ent_to["_id"]


def update_impact_links(ent_from, ent_to):
    """Update the links impact from ent_from to ent_to. Basicaly, append
    the id of ent_to on the the "impact" of ent_from.
    :ent_from: the entities that will be updated
    :ent_to: the entities which id will be used to update :ent_from:"""
    if ent_to["_id"] not in ent_from["impact"]:
        ent_from["impact"] = ent_to["_id"]


def update_links_conn_res(conn, res):
    """Update depends and impact links between the conn connector and the
    res resource. Raise a TypeError if conn is not a connector and/or res
    is not a resource.
    :param conn: the connector to update
    :param res: the resource to update
    :raises TypeError: if the entities does not match the expected one."""
    check_type(conn, "connector")
    check_type(res, "resource")

    update_depends_links(conn, res)
    update_impact_links(res, conn)


def update_links_conn_comp(conn, comp):
    """Update depends and impact links between the conn connector and the
    comp component. Raise a TypeError if conn is not a connector and/or comp
    is not a component.
    :param conn: the connector to update
    :param comp: the component to update
    :raises TypeError: if the entities does not match the expected one."""
    check_type(conn, "connector")
    check_type(comp, "component")

    update_depends_links(conn, comp)
    update_impact_links(comp, conn)


def update_links_res_comp(res, comp):
    """Update depends and impact links between the res resource and the
    comp component. Raise a TypeError if res is not a resource and/or comp
    is not a component.
    :param res: the resource to update
    :param comp: the component to update
    :raises TypeError: if the entities does not match the expected one."""
    check_type(res, "resource")
    check_type(comp, "component")

    update_depends_links(res, comp)
    update_impact_links(comp, res)


def update_case1(entities, ids):
    """Case 1 update entities"""
    LOGGER.debug("Case 1.")
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
    LOGGER.debug("Case 4 : nothing to do here.")


def update_case5(entities, ids):
    """Case 5 update entities"""
    conn_there = False
    re_there = False
    for i in entities:
        if i['type'] == 'connector':
            conn_there = True
        if i['type'] == 'resource':
            re_there = True
    if conn_there:
        if re_there:
            pass
        else:
            # put re + update comp depends + update conn impact
    else:
        if re_there:
            # put comp + maj impac in conn + update re depends
        else:
            # put comp + put re + update conn impact for comp and re

        


def update_case6(entities, ids):
    """Case 6 update entities"""
    conn_there = False
    for i in entities:
        if i['type'] == 'connector':
            conn_there = True
    if not conn_there:
        pass
        # insert conn + maj impact depends in comp and re


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

    global LOGGER
    LOGGER = logger

    # Possible cases :
    # 0 -> Not in cache
    # 1 −> In cache
    #
    #
    #  Connector    Ressource    Component
    #     0             0            0     -> case 1
    #     1             0            0     -> case 2
    #     1             0            1     -> case 3
    #     1             1            1     -> case 4
    #     0             0            1     -> case 5
    #     0             1            1     -> case 6
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
    #    If the event have a resource create a resource
    #    too, then update links between this component and resource and between
    #    connector and resource.
    #
    #  Case 4 :
    #    Every entity are in the cache so they are into the database, nothing
    #    to do here.
    #
    #  Case 5 :
    #    Create a connector and if the event have a resource create a resource
    #    and if needed update the links between the component and the resource.
    #
    #  Case 6 :
    #    Create a connector then update the links between the connector and
    #    the resource and between the connector and component.

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
