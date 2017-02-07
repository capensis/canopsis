# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category


CONF_PATH = 'context_graph/manager.conf'
CATEGORY = 'CONTEXTGRAPH'

@conf_paths(CONF_PATH)
@add_category(CATEGORY)
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

    def check_re(self, re_id):
        """_check_re

        :param re_id:
        """

    def check_conn(self, conn_id):
        """_check_conn

        :param conn_id:
        """

    def check_links(self, conn_id, comp_id, re_id):
        """_check_links

        :param conn_id:
        :param comp_id:
        :param re_id:
        """

    def _checks_conn_re_link(self, conn_id, re_id):
        """_checks_conn_re_link

        :param conn_id:
        :param re_id:
        """
    def _check_conn_comp_link(self, conn_id, comp_id):
        """_check_conn_comp_link

        :param conn_id:
        :param comp_id:
        """

    def add_comp(self, comp):
        """add_comp

        :param comp:
        """

    def add_re(self, re):
        """add_re

        :param re:
        """

    def add_conn(self, conn):
        """add_conn

        :param conn:
        """
