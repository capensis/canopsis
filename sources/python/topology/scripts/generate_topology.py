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
from canopsis.topology.manager import TopologyManager
from canopsis.topology.elements import Edge, Node, Topology

from argparse import ArgumentParser


def generate_context_topology(name='context'):
    """
    Generate a context topology where nodes are components and resources,
    and edges are dependencies from components to resources, or from resources
    to the topology.

    :param str name: topology name.
    """

    # initialize context and topology
    context = Context()
    manager = TopologyManager()

    # clean old topology
    manager.del_elts(ids=name)

    topology = manager.get_graphs(ids=name, add_elts=True)
    if topology is not None:  # if topology already exists, delete content
        for elt_id in topology._gelts:
            elt = topology._gelts[elt_id]
            elt.delete(manager=manager)
        topology.delete(manager=manager)
    # init the topology
    topology = Topology(_id=name)

    def addElt(elt):
        """
        Add input elt in topology.

        :param GraphElement elt: elt to add to topology.
        """

        topology.add_elts(elt.id)
        elt.save(manager)

    components = context.find('component')
    for component in components:
        component_id = context.get_entity_id(component)
        component_node = Node(entity=component_id)
        addElt(component_node)

        ctx, _ = context.get_entity_context_and_name(component)

        resources = context.find('resource', context=ctx)
        if resources:  # link component to all its resources with the same edge
            edge = Edge(sources=component_node.id, targets=[])
            addElt(edge)  # add edge in topology
            for resource in resources:
                resource_id = context.get_entity_id(resource)
                resource_node = Node(entity=resource_id)
                addElt(resource_node)  # add edge in topology
                edge.targets.append(resource_id)
                root_edge = Edge(sources=resource_node.id, targets=topology.id)
                addElt(root_edge)  # add edge in topology
        else:  # if no resources, link the component to the topology
            edge = Edge(sources=component_node.id, targets=topology.id)
            addElt(edge)  # add edge in topology

    # save topology
    topology.save(manager=manager)


if __name__ == '__main__':

    parser = ArgumentParser(description='Generate a topology')
    parser.add_argument(
        dest='name',
        help='topology name to generate (default: context)',
        default='context'
    )
    args = parser.parse_args()

    topology_name = args.name
    generate_context_topology(topology_name)
