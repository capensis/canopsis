# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category


@conf_paths('context_graph/manager.conf')
@add_category('CONTEXTGRAPH')
class ContextGraph(MiddlewareRegistry):
    """ContextGraph"""

    ENTITIES_STORAGE = 'entities_storage'
    ORGANISATIONS_STORAGE = 'organisations_storage'
    USERS_STORAGE = 'measurements_storage'

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

    def get_entity(self, id):
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

        return list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query=query))

    def put_entities(self, entities):
        """
        Store entities into database.
        """
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

    def add_comp(self, comp):
        """add_comp

        :param comp:
        """
        self[ContextGraph.ENTITIES_STORAGE].put_element(element=comp)

    def add_re(self, re):
        """add_re

        :param re:
        """
        self[ContextGraph.ENTITIES_STORAGE].put_element(element=re)

    def add_conn(self, conn):
        """add_conn

        :param conn:
        """
        self[ContextGraph.ENTITIES_STORAGE].put_element(element=conn)

    def get_all_entities(self):
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

    def update_entity(self, id_, entity):
        """Update an entity identified by id_ with the given entity."""
        # TODO add traitement to check if every required field are present

    def delete_enity(self, id_):
        """Delete an entity identified by id_ from the context."""

    def get_entities(self,
                     query={},
                     projection={},
                     limit=0,
                     sort=False,
                     with_count=False):
        """Retreives entities matching the query and the projection.
        """
        #TODO handle projection, limit, sort, with_count

        if isinstance(query, dict):
            raise TypeError("Query must be a dict")

        if isinstance(projection, dict):
            raise TypeError("Projection must be a dict")

        result = self[ContextGraph.ENTITIES_STORAGE].get_elements(query=query)

        return list(result)
