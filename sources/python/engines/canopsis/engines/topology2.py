# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

from canopsis.engines import Engine
from canopsis.topology.manager import Topology
from canopsis.context.manager import Context

from time import time


class engine(Engine):
    etype = 'topology2'

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        # load topology manager
        self.topology = Topology()
        self.context = Context()

    def work(self, event, *args, **kwargs):

        return event

    def operation(
        self, node, at_least, at_most, check_state, in_state, default_state=0
    ):
        """
        Do a generic state operation
        """
        result = node

        state = node['state']

        # do something only if check_state is different than state
        if state != check_state:
            entity_id = self.context.get_entity_id(node)
            source_nodes = self.topology.find_nodes_by_next(next=entity_id)
            len_source_nodes = len(source_nodes)
            for source_node in source_nodes:
                at_most = len_source_nodes - at_most
                if check_state == source_node['state']:
                    at_least -= 1
                    at_most -= 1
                if at_least <= 0 or at_most <= 0:
                    result['state'] = in_state
                    break

            if result['state'] != in_state:
                result['state'] = default_state

        return result

    def all(self, node, check_state, in_state, default_state):

        node_id = self.context.get_entity_id(node)
        source_nodes = self.topology.find_nodes_by_next(next=node_id)
        len_source_nodes = len(source_nodes)

        return self.operation(
            node, at_least=len_source_nodes, at_most=len_source_nodes,
            check_state=check_state, default_state=default_state,
            in_state=in_state)

    def any(self, node, check_state, in_state, default_state):

        node_id = self.context.get_entity_id(node)
        source_nodes = self.topology.find_nodes_by_next(next=node_id)
        len_source_nodes = len(source_nodes)

        return self.operation(
            node, at_least=1, at_most=len_source_nodes,
            check_state=check_state, default_state=default_state,
            in_state=in_state)

    def worst_state(self, node):

        result = node

        entity = node['entity']
        entity_id = self.context.get_entity_id(node)

        if entity['state'] >= node['state']:
            result['state'] = entity['state']

        else:
            source_nodes = self.topology.find_nodes_by_next(next=entity_id)
            result['state'] = max(
                source_node['state'] for source_node in source_nodes)

        return result

    def best_state(self, node):

        result = node

        entity = node['entity']
        entity_id = self.context.get_entity_id(node)

        if entity['state'] <= node['state']:
            result['state'] = entity['state']

        else:
            source_nodes = self.topology.find_nodes_by_next(next=entity_id)
            result['state'] = min(
                source_node['state'] for source_node in source_nodes)

        return result
