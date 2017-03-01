# -*- coding: utf-8 -*-

"""Module in chage of defining the graph context and updating
 it when it's needed."""

from __future__ import unicode_literals

from canopsis.task.core import register_task
from canopsis.context_graph.manager import ContextGraph


context_graph_manager = ContextGraph()


cache = set()

LOGGER = None


def update_cache():
    global cache
    cache = context_graph_manager.get_all_entities()


def create_entity(
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
    LOGGER.error('6')
    """Raise TypeError if the type of the entities entities does not match
    the expected type.
    :param entities: the entities to check.
    :param expected: the expected type.
    :raises TypeError: if the entity does not match the expected one."""
    if entities["type"] != expected:
        raise TypeError("Entities {0} does not match {1}".format(
            entities["_id"], expected))
    return True


def update_depends_links(ent_from, ent_to):
    """Update the links depends from ent_from to ent_to. Basicaly, append
    the id of ent_to on the the "depends" of ent_from.
    :param ent_from: the entities that will be updated
    :param ent_to: the entities which id will be used to update ent_from"""
    if ent_to["_id"] not in ent_from["depends"]:
        ent_from["depends"].append(ent_to["_id"])


def update_impact_links(ent_from, ent_to):
    """Update the links impact from ent_from to ent_to. Basicaly, append
    the id of ent_to on the the "impact" of ent_from.
    :ent_from: the entities that will be updated
    :ent_to: the entities which id will be used to update :ent_from:"""
    if ent_to["_id"] not in ent_from["impact"]:
        ent_from["impact"].append(ent_to["_id"])


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


def determine_presence(ids, data):
    """Determine if the ids are in data. If the ids is present in data, the
    element of the tuple that show the presence of the ids will be set to True.
    False otherwise. If the ids given is set to None, the matching element will
    be set to None.
    :param ids: a set of dict of ids
    :param data: a set of ids
    :return: a tuple of boolean of None following pattern
    (connector presence, component presence, resource presence).

    .. note :: currently an only a resource can be None in an canopsis event.
    """

    conn_here = ids['conn_id'] in data
    comp_here = ids['comp_id'] in data

    if ids['re_id'] is not None:
        res_here = ids['re_id'] in data
    else:
        res_here = None

    LOGGER.error('determine presence ids: {0}, in {1} donne {2}'.format(ids, data, (conn_here, comp_here, res_here)))

    return (conn_here, comp_here, res_here)


def add_missing_ids(presence, ids):
    """Update the cache"""
    if not presence[0]:  # Update connector
        cache.add(ids["conn_id"])

    if not presence[1]:  # Update component
        cache.add(ids["comp_id"])

    if not presence[2]:  # Update resource
        cache.add(ids['re_id'])


def update_context_case1(ids):
    """Case 1 update entities"""

    LOGGER.debug("Case 1.")
    comp = create_entity(
        ids['comp_id'],
        ids['comp_id'],
        'component',
        depends=[ids['conn_id'], ids['re_id']],
        impact=[]
    )
    re = create_entity(
        ids['re_id'],
        ids['re_id'],
        'resource',
        depends=[ids['conn_id']],
        impact=[ids['comp_id']]
    )
    conn = create_entity(
        ids['conn_id'],
        ids['conn_id'],
        'connector',
        depends=[],
        impact=[ids['re_id'], ids['comp_id']]
    )
    context_graph_manager.put_entities([comp, re, conn])


def update_context_case1_re_none(ids):
    LOGGER.debug("Case 1 re none.")
    comp = create_entity(
        ids['comp_id'],
        ids['comp_id'],
        'component',
        depends=[ids['conn_id']],
        impact=[]
    )
    conn = create_entity(
        ids['conn_id'],
        ids['conn_id'],
        'connector',
        depends=[],
        impact=[ids['comp_id']]
    )
    context_graph_manager.put_entities([comp, conn])


def update_context_case2(ids, in_db):
    """Case 2 update entities"""

    LOGGER.debug("Case 2")

    comp = create_entity(
        ids['comp_id'],
        ids['comp_id'],
        'component',
        depends=[ids['re_id']],
        impact=[]
    )
    re = create_entity(
        ids['re_id'],
        ids['re_id'],
        'resource',
        depends=[],
        impact=[ids['comp_id']]
    )
    update_links_conn_res(in_db[0], re)
    update_links_conn_comp(in_db[0], comp)
    context_graph_manager.put_entities([comp, re, in_db[0]])


def update_context_case2_re_none(ids, in_db):
    """Case 2 update entities"""

    LOGGER.debug("Case 2 re none ")

    comp = create_entity(
        ids['comp_id'],
        ids['comp_id'],
        'component',
        depends=[],
        impact=[]
    )
    update_links_conn_comp(in_db[0], comp)
    context_graph_manager.put_entities([comp, in_db[0]])


def update_context_case3(ids, in_db):
    """Case 3 update entities"""
    LOGGER.debug("Case 3")
    comp = {}
    conn = {}
    for i in in_db:
        if i['type'] == 'connector':
            conn = i
        elif i['type'] == 'component':
            comp = i
    re = create_entity(
        ids['re_id'],
        ids['re_id'],
        'resource',
        depends=[],
        impact=[]
    )
    update_links_res_comp(re, comp)
    update_links_conn_res(conn, re)
    LOGGER.error(' GREP cas 3 put {1}'.format([comp, re, conn]))
    context_graph_manager.put_entities([comp, re, conn])


def update_case4(entities, ids):
    """Case 4 update entities"""
    LOGGER.debug("Case 4 : nothing to do here.")


def update_context_case5(ids, in_db):
    """Case 5 update entities"""
    # (False, True, False) or (False, True, None)

    LOGGER.debug("Update context case 6.")

    resource = None

    for entity in in_db:
        if entity["type"] == "resource":
            resource = entity
        elif entity["type"] == "component":
            component = entity

    connector = create_entity(ids["conn_id"],
                              ids["conn_id"],
                              "connector")

    if resource is not None:
        resource = create_entity(ids["re_id"],
                                 ids["re_id"],
                                 "return esource")
        update_links_res_comp(resource, component)
        update_links_conn_res(connector, resource)

    else:
        update_links_conn_res(connector, resource)

    update_links_conn_comp(connector, component)
    context_graph_manager.put_entities([connector, component, resource])


def update_context_case6(ids, in_db):
    LOGGER.debug("Update context case 6.")

    resource = None

    for entity in in_db:
        if entity["type"] == "resource":
            resource = entity
        elif entity["type"] == "component":
            component = entity

    connector = create_entity(ids["conn_id"], ids["conn_id"], "connector")
    update_links_conn_comp(connector, component)

    if ids["re_id"] is not None:
        update_links_conn_res(connector, resource)

    context_graph_manager.put_entities([connector, component, resource])


def update_context(presence, ids, in_db):
    if presence == (False, False, False):
        update_context_case1(ids)
    elif presence == (False, False, None):
        # Case 1
        update_context_case1(ids)
    elif presence == (True, False, False) or presence == (True, False, None):
        # Case 2
        update_context_case2(ids, in_db)

    elif presence == (True, True, False):
        # Case 3
        update_context_case3(ids, in_db)

    elif presence == (True, True, True):
        # Case 4
        pass

    elif presence == (False, True, False) or presence == (False, True, None):
        # Case 5
        update_context_case5(ids, in_db)

    elif presence == (False, True, True) or presence == (False, True, None):
        # Case 6
        update_context_case6(ids, in_db)

    else:
        LOGGER.warning(
            "No case for the given presence : {0} and ids {1}".format(
                presence, ids))
        raise ValueError("No case for the given ids and data.")


def gen_ids(event):
    ret_val = {
        'comp_id': '{0}'.format(event['component']),
        're_id': None,
        'conn_id': '{0}/{1}'.format(
            event['connector'],
            event['connector_name'])
    }
    if 'resource' in event.keys():
        ret_val['re_id'] = '{0}/{1}'.format(event['resource'],
                                            event['component'])
    return ret_val


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
    LOGGER.error('1')
    ids = gen_ids(event)

    LOGGER.error('2')
    presence = determine_presence(ids, cache)
    LOGGER.error('3')
    LOGGER.error(presence)

    if presence == (True, True, True) or presence == (True, True, None):
        # Everything is in cache, so we skip
        logger.error('CAS 4 DANS LE CACHE!!!!')
        return None
    LOGGER.error('4')
    add_missing_ids(presence, ids)

    logger.error('5')

    entities_in_db = context_graph_manager.get_entity(ids.values())
    data = set()
    for i in entities_in_db:
        data.add(i['_id'])

    logger.error('6')
    presence = determine_presence(ids, data)

    logger.error('7')
    if presence == (True, True, True) or presence == (True, True, None):
        # Everything is in cache, so we skip
        return None
    logger.error('8')
    update_context(presence, ids, entities_in_db)
    LOGGER.debug("*** The end. ***")


@register_task
def beat(engine, logger=None, **kwargs):
    update_cache()
