# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from canopsis.common.utils import force_iterable
from canopsis.configuration import conf_paths, add_category
from canopsis.middleware.manager import Manager
from canopsis.storage.scoped import ScopedStorage
from canopsis.storage.filter import Filter

CONF_PATH = 'topology/topology.conf'
CATEGORY = 'TOPOLOGY'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class TopologyManager(Manager):
    """
    Manage access to topologies
    """

    STORAGE = 'topology'
    DATA_SCOPE = STORAGE

    ENTITY_ID = 'entity'
    NEXT = 'next'
    NODES = 'nodes'
    TYPE = 'type'

    def __init__(self, data_scope=DATA_SCOPE, *args, **kwargs):

        super(TopologyManager, self).__init__(
            data_scope=data_scope, *args, **kwargs)

    @staticmethod
    def _get_topology_scope():
        """
        Get topology scope
        """
        return {TopologyManager.TYPE: TopologyManager.DATA_SCOPE}

    @staticmethod
    def _get_topology_node_scope(topology_id=None):
        """
        Get topology node scope related to a topology_id
        """

        result = {'type': 'topology_node'}
        if topology_id is not None:
            result[TopologyManager.STORAGE] = topology_id
        return result

    def _add_nodes(self, topologies):
        """
        Add nodes into input topologies.
        """

        if topologies:
            topologies = force_iterable(topologies)
            for topology in topologies:
                nodes = self.get_nodes(ids=topology[ScopedStorage.ID])
                topology['nodes'] = nodes

    def get_topologies(self, ids=None, add_nodes=False):
        """
        Get one or more topologies.

        :return: depending on ids:

            - an id: one topology or None if corresponding topology does not exist.
            - list of id: list of topologies or empty list if ids do not correspond to existing topologies.
            - None: all existing topologies.
        """
        scope = self._get_topology_scope()
        result = self[TopologyManager.STORAGE].get(scope=scope, ids=ids)
        if add_nodes:
            self._add_nodes(result)
        return result

    def find(self, regex=None, add_nodes=False):
        """
        Find topologies where id match with inpur regex_id
        """

        scope = self._get_topology_scope()
        _filter = Filter()
        _filter.add_regex(
            name=ScopedStorage.ID, value=regex, case_sensitive=True)
        result = self[TopologyManager.STORAGE].find(
            scope=scope, _filter=_filter)
        if add_nodes:
            self._add_nodes(result)
        return result

    def delete(self, ids=None):
        """
        Delete one or more topologies depending on input ids:

            - an id: try to delete topology where id correspond to id
            - list of ids: try to delete topologies where id are in input ids
            - None: delete all topologies
        """

        scope = self._get_topology_scope()
        self[TopologyManager.STORAGE].remove(scope=scope, ids=ids)

    def update(self, _id, topology):
        """
        Update one topology.
        """

        scope = self._get_topology_scope()
        self[TopologyManager.STORAGE].put(scope=scope, _id=_id, data=topology)

    def get_nodes(self, ids=None):
        """
        Get nodes
        """

        scope = self._get_topology_node_scope()
        result = self[TopologyManager.STORAGE].get(scope=scope, ids=ids)

        return result

    def find_nodes_by_entity_id(self, entity_id):

        scope = self._get_topology_node_scope()
        _filter = Filter()
        _filter[TopologyManager.ENTITY_ID] = entity_id
        result = self[TopologyManager.STORAGE].find(
            scope=scope, _filter=_filter)

        return result

    def find_nodes_by_next(self, next=None):

        scope = self._get_topology_node_scope()
        _filter = Filter()
        _filter[TopologyManager.NEXT] = next
        result = self[TopologyManager.STORAGE].find(
            scope=scope, _filter=_filter)

        return result
