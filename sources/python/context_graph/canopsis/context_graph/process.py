# -*- coding: utf-8 -*-

"""Module in chage of defining the graph context and updating
 it when it's needed."""

from __future__ import unicode_literals

from canopsis.task.core import register_task
from canopsis.context_graph.manager import ContextGraph
from canopsis.common.utils import singleton_per_scope

from canopsis.context_graph.manager import ContextGraph


context_graph_manager = ContextGraph()


def create_entity(logger, id, name, etype, depends=[], impact=[], measurements=[], infos={}):
    logger.critical("Create entity of _id {0}".format(id))
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

    ent_from = context[ent1_id]

    # for ent in ent_from["depends"]:
    #     if not ent in ent_to["depends"]:
    #         ent_to["depends"].append(ent)

    if ent2_id not in ent_from["depends"]:
        logger.critical(
            "Add {0} on the field depends of {1}.".format(ent2_id, ent1_id))
        ent_from["depends"].append(ent2_id)


def update_link(logger, ent1_id, ent2_id, context):
    """
    Update depends link from entity ent1 to ent2 in the context.
    Basicaly, add ent2_id in the field "impact" of the entity
    identified by ent1_id the then add ent1_id in the field "depends"
    of ent2_id.
    """
    ent1 = context[ent1_id]
    ent2 = context[ent2_id]

    if ent1_id not in ent2["impact"]:
        logger.critical(
            "Add {0} on the field impacts of {1}.".format(ent1_id, ent2_id))
        ent2["impact"].append(ent1_id)

    if ent2_id not in ent1["depends"]:
        logger.critical(
            "Add {0} on the field impacts of {1}.".format(ent2_id, ent3_id))
        ent1["depends"].append(ent2_id)


@register_task
def event_processing(
        engine, event, manager=None, logger=None, ctx=None, tm=None, cm=None,
        **kwargs
):

    import pprint
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

    logger.critical("Context graph : process")

    logger.critical("Event : {0}".format(pprint.saferepr(event)))

    # Retreive id from id
    comp_id = event['component']

    re_id = None
    if 'resource' in event.keys():
        re_id = '{0}/{1}'.format(event['resource'], comp_id)

    conn_id = '{0}/{1}'.format(event['connector'], event['connector_name'])

    logger.critical("Resource id : {0}.".format(re_id))
    logger.critical("Component id : {0}.".format(comp_id))
    logger.critical("Connector id : {0}.".format(conn_id))

    related_ctx = context_graph_manager.get_entity([comp_id, re_id, conn_id])

    comp_there = False
    re_there = False
    conn_there = False

    for entity in related_ctx:
        if entity["_id"] == comp_id:
            logger.critical(
                "Component {0} already exist in database".format(comp_id))
            comp_there = True

        if entity["_id"] == re_id:
            logger.critical(
                "Resource {0} already exist in database".format(re_id))
            re_there = True

        if entity["_id"] == conn_id:
            logger.critical(
                "Connector {0} already exist in database".format(conn_id))
            conn_there = True

    context = {}

    for doc in related_ctx:
        context[doc["_id"]] = doc

    logger.critical("Local context : {0}".format(context))

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
        logger.critical("Re_id is None and comp_there is True")
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
        logger.critical("Re_id is None")
        update_link(logger, re_id, conn_id, context)

    context_graph_manager.put_entities(context.values())
    logger.critical("The end.")
