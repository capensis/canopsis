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

    def manage_comp_to_re_link(self, re_id, comp_id):
        comp = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id': comp_id}))
        for i in comp:
            if not re_id in i['depends']:
                tmp = i
                tmp['depends'].append(re_id)
                self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)

    def manage_re_to_conn_link(self, conn_id, re_id):
        re = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id': re_id}))
        for i in re:
            if not conn_id in i['depends']:
                tmp = i
                tmp['depends'].append(conn_id)
                self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)


    def manage_comp_to_conn_link(self, conn_id, comp_id):
        comp = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id': comp_id}))
        for i in comp:
            if not conn_id in i['depends']:
                tmp = i
                tmp['depends'].append(conn_id)
                self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)

    def _check_conn_comp_link(self, conn_id, comp_id):
        """_check_conn_comp_link

        :param conn_id:
        :param comp_id:
        """
        conn = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id':conn_id}))
        comp = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id':comp_id}))
        for i in conn:
            for j in comp:
                if not comp_id in i['impact']:
                    tmp = i
                    tmp['impact'].append(comp_id)
                    self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)
                if not conn_id in j['depends']:
                    tmp = j
                    tmp['depends'].append(conn_id)
                    self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)

    def _check_conn_re_link(self, conn_id, re_id):
        """_checks_conn_re_link

        :param conn_id:
        :param re_id:
        """
        conn = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id':conn_id}))
        re = list(self[ContextGraph.ENTITIES_STORAGE].get_elements(query={'_id':re_id}))
        for i in conn:
            for j in re:
                if not re_id in i['impact']:
                    tmp = i
                    tmp['impact'].append(re_id)
                    self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)
                if not conn_id in j['depends']:
                    tmp = j
                    tmp['depends'].append(conn_id)
                    self[ContextGraph.ENTITIES_STORAGE].put_element(element=tmp)

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
