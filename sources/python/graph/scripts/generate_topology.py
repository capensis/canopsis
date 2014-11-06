#!/usr/bin/env python
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

from canopsis.context.manager import Context
from canopsis.topology.manager import Topology

from uuid import uuid4 as uuid
from argparse import ArgumentParser


def generate_context_topology(name='test'):
    """
    Generate a context topology related to some random parameters

    :param str name: topology name
    """

    # initialize context and topology
    context = Context()
    topology = Topology()

    # clean old topology
    topology.delete(ids=name)
    topology[Topology.STORAGE].remove_elements()

    # create a root node
    root_node = {
        Topology.ID: str(uuid()),
        Topology.ENTITY_ID: topology_name,
        Topology.TOPOLOGY_ID: topology_name,
        Topology.NEXT: []
    }

    topology_graph = {
        Topology.ID: topology_name,  # topology name
        Topology.ROOT: root_node[Topology.ID]  # root id
    }

    # get all components and resources
    components = context.find('component')
    resources = context.find('resource')

    nodes_by_entity_id = {}

    def generate_node(entity):
        entity_id = entity['_id']
        node_id = str(uuid())
        result = {
            Topology.ID: node_id,
            Topology.ENTITY_ID: entity_id,
            Topology.NEXT: [],
            Topology.TOPOLOGY_ID: topology_name
        }
        nodes_by_entity_id[entity_id] = result
        return result

    # generate nodes from components
    for component in components:
        generate_node(component)

    # generate nodes from resources
    for resource in resources:
        node = generate_node(resource)
        component = resource.copy()
        component[Context.TYPE] = 'component'
        component[Context.NAME] = component['component']
        del component['component']
        entity_id = context.get_entity_id(component)
        # add resource in parent component node next nodes
        try:
            next = nodes_by_entity_id[entity_id][Topology.NEXT]
        except Exception as e:
            # in some awkward cases, entity_id does not exist in nodes_by_entity_id
            # TODO: fix this case x)
            print(e)
            continue
        next.append(node[Topology.ID])
        # add root in resource.next
        node[Topology.NEXT].append(root_node[Topology.ID])

    # save nodes in storage
    for entity_id in nodes_by_entity_id:
        node = nodes_by_entity_id[entity_id]
        topology.push_node(node)
    # add root node
    topology.push_node(root_node)
    # save topology graph
    topology.push(topology_graph)


if __name__ == '__main__':

    parser = ArgumentParser(description='Generate a topology')
    parser.add_argument(dest='name', help='topology name to generate')
    parser.add_argument(
        '-o', dest='operations', help='number of operations to add', default=0)
    parser.add_argument(
        '-n', dest='nodes', default=1,
        help='maximal number of nodes to bind to operations')
    args = parser.parse_args()

    topology_name = args.name
    generate_context_topology(topology_name)
