# -*- coding: utf-8 -*-
"""Module in chage of defining the graph context and updating
 it when it's needed."""

from __future__ import unicode_literals

import time

from canopsis.task.core import register_task
from canopsis.context_graph.manager import ContextGraph
from canopsis.perfdata.manager import PerfData

context_graph_manager = None
pertfdata_manager = PerfData(*PerfData.provide_default_basics())

cache = set()

LOGGER = None

def update_cache():
    """Update the entity cache "cache"
    """
    global cache
    cache = context_graph_manager.get_all_entities_id()


def create_entity(id,
                  name,
                  etype,
                  depends=[],
                  impact=[],
                  measurements={},
                  infos={},
                  **kwargs):
    """Create an entity with following information and put it state at enable.
    :param id: the entity id
    :type id: a string
    :param name: the entity name
    :type name: a string
    :param etype: the entity type
    :type etype: a string
    :param depends: every entities that depends of the current entity
    :type depends: a list
    :param impact: every entities that depends of the current entity
    :type impact: a list
    :param measurements: measurements link to the current entity
    :type measurements: a dict
    :param infos: information related to the entity
    :type infos: a dict

    :return: a dict
    """

    ent = {
            '_id': id,
            'type': etype,
            'name': name,
            'depends': depends,
            'impact': impact,
            'measurements': measurements,
            'infos': infos
        }

    for key in kwargs:
        ent[key] = kwargs[key]

    if etype == 'component':
        ent.pop("measurements")

    context_graph_manager._enable_entity(ent)

    return ent


def check_type(entities, expected):
    """Raise TypeError if the type of the entities entities does not match
    the expected type.
    :param entities: the entities to check.
    :type entities: a dict
    :param expected: the expected type.
    :type expected: a string
    :raises TypeError: if the entity does not match the expected one.
    :return: true if the entity match the given type
    """
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

    update_impact_links(conn, res)
    update_depends_links(res, conn)


def update_links_conn_comp(conn, comp):
    """Update depends and impact links between the conn connector and the
    comp component. Raise a TypeError if conn is not a connector and/or comp
    is not a component.
    :param conn: the connector to update
    :param comp: the component to update
    :raises TypeError: if the entities does not match the expected one."""
    check_type(conn, "connector")
    check_type(comp, "component")

    update_impact_links(conn, comp)
    update_depends_links(comp, conn)


def update_links_res_comp(res, comp):
    """Update depends and impact links between the res resource and the
    comp component. Raise a TypeError if res is not a resource and/or comp
    is not a component.
    :param res: the resource to update
    :param comp: the component to update
    :raises TypeError: if the entities does not match the expected one."""
    check_type(res, "resource")
    check_type(comp, "component")

    update_impact_links(res, comp)
    update_depends_links(comp, res)


def determine_presence(ids, data):
    """Determine if the ids are in data. If the ids is present in data, the
    element of the tuple that show the presence of the ids will be set to True.
    False otherwise. If the ids given is set to None, the matching element will
    be set to None.
    :param ids: a set of dict of ids
    :param data: a set of ids
    :return: a tuple of boolean of None following pattern
    (connector presence, component presence, resource presence).

    .. note :: currently only a resource can be None in an canopsis event.
    """

    conn_here = ids['conn_id'] in data
    comp_here = ids['comp_id'] in data

    if ids['re_id'] is not None:
        res_here = ids['re_id'] in data
    else:
        res_here = None

    return (conn_here, comp_here, res_here)


def add_missing_ids(presence, ids):
    """Add missing id in the cache. See determine_presence function."""
    if not presence[0]:  # Update connector
        cache.add(ids["conn_id"])

    if not presence[1]:  # Update component
        cache.add(ids["comp_id"])

    if presence[2] is not None and not presence[2]:  # Update resource
        cache.add(ids['re_id'])


def create_ent_metric(event):
    """Create a metric entity from an event.
    :param event: the event to use to create the metric entity
    :type event: a dict.
    :return: an event
    :return type: a dict"""
    result = []

    # Parse the perf_data field.
    if "perf_data" in event and event["perf_data"] is not None:
        parser = PerfData(event["perf_data"])
        result += parser.perf_data_array()

    for perf in event["perf_data_array"]:
        id_ = "/metric/{0}/{1}/{2}/{3}".format(
            event["connector"], event["connector_name"], event["component"],
            perf["metric"])
        ent_metric = create_entity(
            id=id_,
            name=perf["metric"],
            etype="metric",
            depends=[],
            impact=[event["resource"]],
            measurements={},
            infos={},
            resource=event["resource"],
            component=event["component"],
            connector=event["connector"],
            connector_name=event["connector_name"])
        result.append(ent_metric)

    return result


def update_context_case1(ids, event):
    """Case 1 update entities. No resource, component or connector exist in
    the context.
    :param ids: the tuple of ids.
    :type ids: a tuple
    :param event: the event to process
    :type event: a dict.
    :return: a list of entities (as a dict)
    """

    result = []

    LOGGER.debug("Case 1.")
    comp = create_entity(
        ids['comp_id'],
        event['component'],
        'component',
        depends=[ids['conn_id'], ids['re_id']],
        impact=[])

    conn = create_entity(
        ids['conn_id'],
        event['connector_name'],
        'connector',
        depends=[],
        impact=[ids['re_id'], ids['comp_id']])

    result.append(comp)
    result.append(conn)

    if event["event_type"] == "perf":
        result += create_ent_metric(event)

    re = create_entity(
        ids['re_id'],
        event['resource'],
        'resource',
        depends=[ids['conn_id']],
        impact=[ids['comp_id']])
    result.append(re)


    return result


def update_context_case1_re_none(ids, event):
    """Case 1 update entities. No component or connector exist in
    the context and no resource are define in the event.
    :param ids: the tuple of ids.
    :type ids: a tuple
    :return: a list of entities (as a dict)
    """
    LOGGER.debug("Case 1 re none.")
    comp = create_entity(
        ids['comp_id'],
        event['component'],
        'component',
        depends=[ids['conn_id']],
        impact=[])
    conn = create_entity(
        ids['conn_id'],
        event['connector_name'],
        'connector',
        depends=[],
        impact=[ids['comp_id']])
    return [comp, conn]


def update_context_case2(ids, in_db, event):
    """Case 2 update entities. A component exist in the context
    but no connector nor resource.
    :param ids: the tuple of ids.
    :param type: a tuple
    :param in_db: a list of the rest of entity in the event.
    :type in_db: a tuple
    :return: a list of entities (as a dict)"""

    LOGGER.debug("Case 2")

    comp = create_entity(
        ids['comp_id'],
        event['component'],
        'component',
        depends=[ids['re_id']],
        impact=[])
    re = create_entity(
        ids['re_id'],
        event['resource'],
        'resource',
        depends=[],
        impact=[ids['comp_id']])
    update_links_conn_res(in_db[0], re)
    update_links_conn_comp(in_db[0], comp)
    return [comp, re, in_db[0]]


def update_context_case2_re_none(ids, in_db, event):
    """Case 2 update entities. A component exist in the context
    but no connector and no resource are in the event.
    :param ids: the tuple of ids.
    :param type: a tuple
    :param in_db: a list of the rest of entity in the event.
    :type in_db: a tuple
    :return: a list of entities (as a dict)"""

    LOGGER.debug("Case 2 re none ")

    comp = create_entity(
        ids['comp_id'], event['component'], 'component', depends=[], impact=[])
    update_links_conn_comp(in_db[0], comp)
    return [comp, in_db[0]]


def update_context_case3(ids, in_db, event):
    """Case 3 update entities. A component and connector exist in the context
    but no resource.
    :param ids: the tuple of ids.
    :param type: a tuple
    :param in_db: a list of the rest of entity in the event.
    :type in_db: a tuple
    :return: a list of entities (as a dict)"""

    LOGGER.debug("Case 3")
    comp = {}
    conn = {}
    for i in in_db:
        if i['type'] == 'connector':
            conn = i
        elif i['type'] == 'component':
            comp = i
    re = create_entity(
        ids['re_id'], event['resource'], 'resource', depends=[], impact=[])
    update_links_res_comp(re, comp)
    update_links_conn_res(conn, re)
    return [comp, re, conn]


def update_context_case5(ids, in_db, event):
    """Case 5 update entities. A connector exist in the context
    but no component.
    :param ids: the tuple of ids.
    :param type: a tuple
    :param in_db: a list of the rest of entity in the event.
    :type in_db: a tuple
    :return: a list of entities (as a dict)"""

    LOGGER.debug("Update context case 6.")

    resource = None
    component = None

    for entity in in_db:
        if entity["type"] == "resource":
            resource = entity
        elif entity["type"] == "component":
            component = entity

    connector = create_entity(
        ids["conn_id"], event["connector_name"], "connector", impact=[], depends=[])

    if ids["re_id"] is not None:
        #res_impact = [component["_id"]]
        #res_depends = [connector["_id"]]
        resource = create_entity(
            ids["re_id"], event["resource"], "resource", impact=[], depends=[])
        update_links_res_comp(resource, component)
        update_links_conn_res(connector, resource)
        update_links_conn_comp(connector, component)
        return [connector, component, resource]

    LOGGER.debug("Ressource  None")
    update_links_conn_comp(connector, component)
    return [connector, component]


def update_context_case6(ids, in_db, event):
    """Case 6 update entities. A connector and a resource exist in the context
    but no component.
    :param ids: the tuple of ids.
    :param type: a tuple
    :param in_db: a list of the rest of entity in the event.
    :type in_db: a tuple
    :return: a list of entities (as a dict)"""

    LOGGER.debug("Update context case 6.")

    resource = None

    for entity in in_db:
        if entity["type"] == "resource":
            resource = entity
        elif entity["type"] == "component":
            component = entity

    connector = create_entity(
        ids["conn_id"], event["connector_name"], "connector", impact=[], depends=[])
    update_links_conn_comp(connector, component)

    if resource is not None:
        update_links_conn_res(connector, resource)
        return [connector, component, resource]

    return [connector, component]


def update_context(presence, ids, in_db, event):
    """Update the context.
    :param presence: information about a entity exist in the database
    :param ids: a list of ids.
    :param in_db: the related entities from the context
    :param event: the current event
    """
    extra_infos = {}
    for field in context_graph_manager.extra_fields:
        if field in event.keys():
            extra_infos[field] = event[field]

    event_id = get_event_id(ids, event)

    if "infos" in event:
        event_info = event["infos"]
    else:
        event_info = {}

    to_update = None
    if presence == (False, False, False):
        # Case 1
        to_update = update_context_case1(ids, event)

    elif presence == (False, False, None):
        # Case 1
        to_update = update_context_case1_re_none(ids, event)

    elif presence == (True, False, False):
        # Case 2
        to_update = update_context_case2(ids, in_db, event)

    elif presence == (True, False, None):
        to_update = update_context_case2_re_none(ids, in_db, event)

    elif presence == (True, True, False):
        # Case 3
        to_update = update_context_case3(ids, in_db, event)

    elif presence == (True, True, True):
        # Case 4
        pass

    elif presence == (False, True, False) or presence == (False, True, None):
        # Case 5
        to_update = update_context_case5(ids, in_db, event)

    elif presence == (False, True, True) or presence == (False, True, None):
        # Case 6
        to_update = update_context_case6(ids, in_db, event)

    else:
        LOGGER.warning("No case for the given presence : {0} and ids {1}".
                       format(presence, ids))
        raise ValueError("No case for the given ids and data.")

    evt_entity = None
    for entity in to_update:

        # If there is no "enabled_history" field in "infos", we assume
        # the entity was just create
        if "enable_history" not in entity:
            ContextGraph._enable_entity(entity, event["timestamp"])

        if entity["_id"] == event_id:
            evt_entity = entity

    for key in extra_infos:
        evt_entity['infos'][key] = extra_infos[key]

    for key in event_info:
        evt_entity["infos"][key] = event_info[key]

    context_graph_manager._put_entities(to_update)


def gen_ids(event):
    """Generate the id of entity present in an event
    :param event: the current event
    """
    ret_val = {
        'comp_id': '{0}'.format(event['component']),
        're_id': None,
        'conn_id': '{0}/{1}'.format(event['connector'],
                                    event['connector_name'])
    }
    if 'resource' in event.keys():
        ret_val['re_id'] = '{0}/{1}'.format(event['resource'],
                                            event['component'])
    return ret_val


def get_event_id(ids, event):
    """Retreive the event id from the dict generated by gen_ids function
    with the given event.
    """
    # ascertain the type of the entity related to the event
    if "resource" in event:
        return ids["re_id"]

    if "component" in event:
        return ids["comp_id"]

    return ids["conn_id"]


@register_task
def event_processing(engine,
                     event,
                     manager=None,
                     logger=None,
                     ctx=None,
                     tm=None,
                     cm=None,
                     **kwargs):
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

    global context_graph_manager
    if context_graph_manager is None:
        context_graph_manager = ContextGraph(logger)

    if event['event_type'] not in context_graph_manager.event_types:
        return None

    global LOGGER
    LOGGER = logger
    ids = gen_ids(event)

    presence = determine_presence(ids, cache)

    if presence == (True, True, True) or presence == (True, True, None):
        # Everything is in cache, so we skip
        return None
    add_missing_ids(presence, ids)

    entities_in_db = context_graph_manager.get_entities_by_id(ids.values())
    data = set()
    for i in entities_in_db:
        data.add(i['_id'])

    presence = determine_presence(ids, data)

    if presence == (True, True, True) or presence == (True, True, None):
        # Everything is in cache, so we skip
        return None

    update_context(presence, ids, entities_in_db, event)
    LOGGER.debug("*** The end. ***")


@register_task
def beat(engine, logger=None, **kwargs):
    update_cache()
