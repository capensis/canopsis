# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.event import forger


@conf_paths('context_graph/manager.conf')
@add_category('CONTEXTGRAPH')
class ContextGraph(MiddlewareRegistry):
    """ContextGraph"""

    ENTITIES_STORAGE = 'entities_storage'
    ORGANISATIONS_STORAGE = 'organisations_storage'
    USERS_STORAGE = 'measurements_storage'

    RESOURCE = "resource"
    COMPONENT = "component"
    CONNECTOR = "connector"

    @classmethod
    def get_id(cls, event):
        """Return the id extracted from the event as a string
        :param event: the event from which we extract the id
        :return type: boolean a string
        """
        if '_id' in event:
            return event['_id']

        source_type = ''
        if 'source_type' in event:
            source_type = event['source_type']
        elif 'type' in event:
            source_type = event['type']

        id_ = ""

        if source_type == cls.COMPONENT:
            id_ = event["component"]
        elif source_type == cls.RESOURCE:
            id_ = "{0}/{1}".format(event["resource"], event["component"])
        elif source_type == cls.CONNECTOR:
            id_ = "{0}/{1}".format(event["connector"], event["connector_name"])
        else:
            error_desc = "Event type should be 'connector', 'resource' or\
            'component' not {0}.".format(source_type)
            raise ValueError(error_desc)

        return id_


    def __init__(self, *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """
        super(ContextGraph, self).__init__(*args, **kwargs)

    def check_comp(self, comp_id):
        """_check_comp

        check if the component exists in database

        :param comp_id: id of component
        :return type: boolean
        """
        return len(list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id': comp_id}))) > 0

    def get_entities_by_id(self, id):
        """
        Retreive the entity identified by his id. If id is a list of id,
        get_entity return every entities who match the ids present in the list

        :param id the id of the entity. id can be a list
        """
        query = {"_id": None}
        if isinstance(id, type([])):
            ids = []
            for i in id:
                ids.append(i)
            query["_id"] = {"$in": ids}
        else:
            query["_id"] = id

        result = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query=query))
        if len(result) == 1:
            return result[0]
        elif len(result) == 0:
            # Return {} due to the plentiful of iteration on the keys of entity
            return {}
        else:
            return result

    def put_entities(self, entities):
        """
        Store entities into database.
        """
        if not isinstance(entities, list):
            entities = [entities]
        self[ContextGraph.ENTITIES_STORAGE].put_elements(entities)

    def check_re(self, re_id):
        """_check_re

        :param re_id:
        """
        return len(list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id': re_id}))) > 0

    def check_conn(self, conn_id):
        """_check_conn

        :param conn_id:
        """
        return len(list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id': conn_id}))) > 0

    def check_links(self, conn_id, comp_id, re_id):
        """_check_links

        :param conn_id:
        :param comp_id:
        :param re_id:
        """
        raise NotImplementedError

    def manage_comp_to_re_link(self, re_id, comp_id):
        """Update component-resource link"""
        comp = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': comp_id}))
        for i in comp:
            if re_id not in i['depends']:
                tmp = i
                tmp['depends'].append(re_id)
                self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)

    def manage_re_to_conn_link(self, conn_id, re_id):
        """Update resource-connector link"""
        re = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': re_id}))
        for i in re:
            if conn_id not in i['depends']:
                tmp = i
                tmp['depends'].append(conn_id)
                self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)

    def manage_comp_to_conn_link(self, conn_id, comp_id):
        """Update component-connector link"""
        comp = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': comp_id}))
        for i in comp:
            if conn_id not in i['depends']:
                tmp = i
                tmp['depends'].append(conn_id)
                self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)

    def _check_conn_comp_link(self, conn_id, comp_id):
        """_check_conn_comp_link

        :param conn_id:
        :param comp_id:
        """
        conn = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': conn_id}))
        comp = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': comp_id}))
        for i in conn:
            for j in comp:
                if comp_id not in i['impact']:
                    tmp = i
                    tmp['impact'].append(comp_id)
                    self[ContextGraph.ENTITIES_STORAGE].put_element(
                        element=tmp)
                if conn_id not in j['depends']:
                    tmp = j
                    tmp['depends'].append(conn_id)
                    self[ContextGraph.ENTITIES_STORAGE].put_element(
                        element=tmp)

    def _check_conn_re_link(self, conn_id, re_id):
        """_checks_conn_re_link

        :param conn_id:
        :param re_id:
        """
        conn = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': conn_id}))
        re = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': re_id}))
        for i in conn:
            for j in re:
                if re_id not in i['impact']:
                    tmp = i
                    tmp['impact'].append(re_id)
                    self[ContextGraph.ENTITIES_STORAGE].put_element(
                        element=tmp)
                if conn_id not in j['depends']:
                    tmp = j
                    tmp['depends'].append(conn_id)
                    self[ContextGraph.ENTITIES_STORAGE].put_element(
                        element=tmp)

    def _check_comp_re_link(self, comp_id, re_id):
        """_check_com_re_link

        :param comp_id:
        :param re_id:
        """


    def get_all_entities_id(self):
        """
            get all entities ids by types
        """
        entities = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={}))
        ret_val = set([])
        for i in entities:
            ret_val.add(i['_id'])
        # pritn(ret_val)
        return ret_val

    def create_entity(self, entity):
        """Create an entity in the contexte with the given entity."""
        # TODO add traitement to check if every required field are present
        if entity['depends'] != []:
            # update
            query = {'$or': []}
            for i in entity['depends']:
                query['$or'].append({'_id': i})
            tmp = self[ContextGraph.ENTITIES_STORAGE].get_elements(
                query=query
            )
            for j in tmp:
                try:
                    j['impact'].append(entity['_id'])
                    self[ContextGraph.ENTITIES_STORAGE].put_element(j)
                except ValueError:
                    pass

        if entity['impact'] != []:
            # update
            query = {'$or': []}
            for i in entity['impact']:
                query['$or'].append({'_id': i})
            tmp = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
                query=query
            ))
            for j in tmp:
                try:
                    j['depends'].append(entity['_id'])
                    self[ContextGraph.ENTITIES_STORAGE].put_element(j)
                except ValueError:
                    pass

        self[ContextGraph.ENTITIES_STORAGE].put_element(entity)

    def update_entity(self, entity):
        """Update an entity identified by id_ with the given entity."""
        # TODO add traitement to check if every required field are present
        self[ContextGraph.ENTITIES_STORAGE].put_element(entity)

    def delete_entity(self, id_):
        """Delete an entity identified by id_ from the context."""
        entity = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query={'_id': id_}
        ))[0]
        if entity['depends'] != []:
            # update entity in depends list
            query = {'$or': []}
            for i in entity['depends']:
                query['$or'].append({'_id': i})

            tmp = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
                query=query
            ))
            for j in tmp:
                try:
                    j['impact'].remove(id_)
                    self[ContextGraph.ENTITIES_STORAGE].put_element(j)
                except ValueError:
                    pass
        if entity['impact'] != []:
            # update entity in impact list
            query={'$or': []}
            for i in entity['impact']:
                query['$or'].append({'_id': i})
            tmp = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
                query=query
            ))
            for j in tmp:
                try:
                    j['depends'].remove(id_)
                except ValueError:
                    # a corriger mais je tente comme ça pour les tests
                    pass
                self[ContextGraph.ENTITIES_STORAGE].put_element(j)

        self[ContextGraph.ENTITIES_STORAGE].remove_elements(ids=[id_])

    def get_entities(self,
                     query={},
                     projection={},
                     limit=0,
                     start=0,
                     sort=False,
                     with_count=False):
        """Retreives entities matching the query and the projection.
        """
        # TODO handle projection, limit, sort, with_count

        if not isinstance(query, dict):
            raise TypeError("Query must be a dict")

        if not isinstance(projection, dict):
            raise TypeError("Projection must be a dict")

        result = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(
            query=query,
            limit=limit,
            skip=start,
            sort=sort,
            projection=projection,
            with_count=with_count
        ))

        return result


    def get_event(self, entity, event_type='check', **kwargs):
        """Get an event from an entity.

        :param dict entity: entity to convert to an event.
        :param str event_type: specific event_type. Default check.
        :param dict kwargs: additional fields to put in the event.
        :rtype: dict
        """

        # keys from entity that should not be in event
        delete_keys = ["_id", "depends", "impact", "type"]

        kwargs['event_type'] = event_type

        # In some cases, name is present but is component in fact
        if 'name' in entity:
            if 'component' not in entity:
                entity['component'] = entity['name']
            entity.pop('name')

        # remove key that should not be in the event
        tmp = entity.copy()
        for key in delete_keys:
            tmp.pop(key)

        # fill kwargs with entity values
        for field in tmp:
            kwargs[field] = tmp[field]

        # forge the event
        result = forger(**kwargs)

        return result
